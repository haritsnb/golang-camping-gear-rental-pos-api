package repositories

import (
	"app/models"
	"context"
	"database/sql"
)

type ProductUnitRepository interface {
	GetAll(ctx context.Context) ([]models.ProductUnit, error)
	GetByID(ctx context.Context, id int) (*models.ProductUnit, error)
	GetByUnitCode(ctx context.Context, code string) (*models.ProductUnit, error)
	GetByIDWithTx(ctx context.Context, tx *sql.Tx, id int) (*models.ProductUnit, error)
	GetDeleted(ctx context.Context) ([]models.ProductUnit, error)
	UpdateStatusWithTx(ctx context.Context, tx *sql.Tx, id int, status string, modifiedBy int) error
	Create(ctx context.Context, u *models.ProductUnit) error
	Update(ctx context.Context, id int, req models.ProductUnitUpdateDTO, modifiedBy int) error
	Delete(ctx context.Context, id, deletedBy int) error
	Restore(ctx context.Context, id, restoredBy int) error
	ForceDelete(ctx context.Context, id int) error
}

type productUnitRepo struct{ db *sql.DB }

func NewProductUnitRepository(db *sql.DB) ProductUnitRepository { return &productUnitRepo{db: db} }

func (r *productUnitRepo) GetAll(ctx context.Context) ([]models.ProductUnit, error) {
	query := `
		SELECT pu.id, pu.product_id, p.name, pu.unit_code, pu.serial_number, pu.condition, pu.status, pu.notes, pu.created_at, pu.modified_at
		FROM product_units pu
		JOIN products p ON p.id = pu.product_id
		WHERE pu.deleted_at IS NULL ORDER BY pu.id DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var units []models.ProductUnit
	for rows.Next() {
		var u models.ProductUnit
		if err := rows.Scan(&u.ID, &u.ProductID, &u.ProductName, &u.UnitCode, &u.SerialNumber, &u.Condition, &u.Status, &u.Notes, &u.CreatedAt, &u.ModifiedAt); err != nil {
			return nil, err
		}
		units = append(units, u)
	}
	return units, nil
}

func (r *productUnitRepo) GetByID(ctx context.Context, id int) (*models.ProductUnit, error) {
	query := `
		SELECT pu.id, pu.product_id, p.name, pu.unit_code, pu.serial_number, pu.condition, pu.status, pu.notes, pu.created_at, pu.modified_at
		FROM product_units pu
		JOIN products p ON p.id = pu.product_id
		WHERE pu.id = $1 AND pu.deleted_at IS NULL`
	u := &models.ProductUnit{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&u.ID, &u.ProductID, &u.ProductName, &u.UnitCode, &u.SerialNumber, &u.Condition, &u.Status, &u.Notes, &u.CreatedAt, &u.ModifiedAt)
	return u, err
}

func (r *productUnitRepo) GetByUnitCode(ctx context.Context, code string) (*models.ProductUnit, error) {
	u := &models.ProductUnit{}
	query := `
		SELECT pu.id, pu.product_id, p.name, pu.unit_code, pu.serial_number, pu.condition, pu.status, pu.notes, pu.created_at, pu.modified_at
		FROM product_units pu
		JOIN products p ON p.id = pu.product_id
		WHERE pu.unit_code = $1 AND pu.deleted_at IS NULL`
	err := r.db.QueryRowContext(ctx, query, code).
		Scan(&u.ID, &u.ProductID, &u.ProductName, &u.UnitCode, &u.SerialNumber, &u.Condition, &u.Status, &u.Notes, &u.CreatedAt, &u.ModifiedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *productUnitRepo) GetByIDWithTx(ctx context.Context, tx *sql.Tx, id int) (*models.ProductUnit, error) {
	query := `
		SELECT pu.id, pu.product_id, p.name, pu.unit_code, pu.serial_number, pu.condition, pu.status, pu.notes, pu.created_at, pu.modified_at
		FROM product_units pu
		JOIN products p ON p.id = pu.product_id
		WHERE pu.id = $1 AND pu.deleted_at IS NULL FOR UPDATE`
	u := &models.ProductUnit{}
	err := tx.QueryRowContext(ctx, query, id).Scan(&u.ID, &u.ProductID, &u.ProductName, &u.UnitCode, &u.SerialNumber, &u.Condition, &u.Status, &u.Notes, &u.CreatedAt, &u.ModifiedAt)
	return u, err
}

func (r *productUnitRepo) GetDeleted(ctx context.Context) ([]models.ProductUnit, error) {
	query := `
		SELECT pu.id, pu.product_id, p.name, pu.unit_code, pu.serial_number, pu.condition, pu.status, pu.notes, pu.created_at, pu.modified_at, pu.deleted_at
		FROM product_units pu
		JOIN products p ON p.id = pu.product_id
		WHERE pu.deleted_at IS NOT NULL ORDER BY pu.deleted_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.ProductUnit
	for rows.Next() {
		var u models.ProductUnit
		if err := rows.Scan(&u.ID, &u.ProductID, &u.ProductName, &u.UnitCode, &u.SerialNumber, &u.Condition, &u.Status, &u.Notes, &u.CreatedAt, &u.ModifiedAt, &u.DeletedAt); err != nil {
			return nil, err
		}
		list = append(list, u)
	}
	return list, nil
}

func (r *productUnitRepo) UpdateStatusWithTx(ctx context.Context, tx *sql.Tx, id int, status string, modifiedBy int) error {
	query := `UPDATE product_units SET status = $1, modified_at = NOW(), modified_by = $2 WHERE id = $3 AND deleted_at IS NULL`
	_, err := tx.ExecContext(ctx, query, status, modifiedBy, id)
	return err
}

func (r *productUnitRepo) Create(ctx context.Context, u *models.ProductUnit) error {
	query := `
		INSERT INTO product_units (product_id, unit_code, serial_number, condition, status, notes, created_by, modified_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, created_at, modified_at`
	return r.db.QueryRowContext(ctx, query, u.ProductID, u.UnitCode, u.SerialNumber, u.Condition, u.Status, u.Notes, u.CreatedBy, u.ModifiedBy).Scan(&u.ID, &u.CreatedAt, &u.ModifiedAt)
}

func (r *productUnitRepo) Update(ctx context.Context, id int, req models.ProductUnitUpdateDTO, modifiedBy int) error {
	query := `
		UPDATE product_units 
		SET 
			product_id    = COALESCE($1, product_id),
			unit_code     = COALESCE($2, unit_code),
			serial_number = COALESCE($3, serial_number),
			condition     = COALESCE($4, condition),
			status        = COALESCE($5, status),
			notes         = COALESCE($6, notes),
			modified_at   = NOW(),
			modified_by   = $7
		WHERE id = $8 AND deleted_at IS NULL`

	res, err := r.db.ExecContext(ctx, query,
		req.ProductID,
		req.UnitCode,
		req.SerialNumber,
		req.Condition,
		req.Status,
		req.Notes,
		modifiedBy,
		id,
	)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *productUnitRepo) Delete(ctx context.Context, id, deletedBy int) error {
	res, err := r.db.ExecContext(ctx, `UPDATE product_units SET deleted_at = NOW(), deleted_by = $1 WHERE id = $2 AND deleted_at IS NULL`, deletedBy, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *productUnitRepo) Restore(ctx context.Context, id, restoredBy int) error {
	query := `UPDATE product_units SET deleted_at = NULL, deleted_by = NULL, modified_at = NOW(), modified_by = $1 WHERE id = $2 AND deleted_at IS NOT NULL`
	res, err := r.db.ExecContext(ctx, query, restoredBy, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *productUnitRepo) ForceDelete(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM product_units WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
