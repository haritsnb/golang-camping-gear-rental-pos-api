package repositories

import (
	"app/models"
	"context"
	"database/sql"
)

type MaintenanceRepository interface {
	GetAll(ctx context.Context) ([]models.Maintenance, error)
	GetByID(ctx context.Context, id int) (*models.Maintenance, error)
	Create(ctx context.Context, m *models.Maintenance) error
	CreateWithTx(ctx context.Context, tx *sql.Tx, m *models.Maintenance) error
	Update(ctx context.Context, id int, req models.MaintenanceUpdateDTO, modifiedBy int) error
	Complete(ctx context.Context, id int, modifiedBy int) error
	Delete(ctx context.Context, id, deletedBy int) error
	Restore(ctx context.Context, id, restoredBy int) error
	ForceDelete(ctx context.Context, id int) error
}

type maintenanceRepo struct{ db *sql.DB }

func NewMaintenanceRepository(db *sql.DB) MaintenanceRepository { return &maintenanceRepo{db: db} }

func (r *maintenanceRepo) GetAll(ctx context.Context) ([]models.Maintenance, error) {
	query := `
		SELECT m.id, m.product_unit_id, pu.unit_code, m.issue_description, m.cost, m.start_date, m.completed_date, m.status, m.created_at, m.modified_at
		FROM maintenance m
		JOIN product_units pu ON pu.id = m.product_unit_id
		WHERE m.deleted_at IS NULL ORDER BY m.id DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.Maintenance
	for rows.Next() {
		var m models.Maintenance
		if err := rows.Scan(&m.ID, &m.ProductUnitID, &m.UnitCode, &m.IssueDescription, &m.Cost, &m.StartDate, &m.CompletedDate, &m.Status, &m.CreatedAt, &m.ModifiedAt); err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, nil
}

func (r *maintenanceRepo) GetByID(ctx context.Context, id int) (*models.Maintenance, error) {
	query := `
		SELECT m.id, m.product_unit_id, pu.unit_code, m.issue_description, m.cost, m.start_date, m.completed_date, m.status, m.created_at, m.modified_at
		FROM maintenance m
		JOIN product_units pu ON pu.id = m.product_unit_id
		WHERE m.id = $1 AND m.deleted_at IS NULL`
	m := &models.Maintenance{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&m.ID, &m.ProductUnitID, &m.UnitCode, &m.IssueDescription, &m.Cost, &m.StartDate, &m.CompletedDate, &m.Status, &m.CreatedAt, &m.ModifiedAt)
	return m, err
}

func (r *maintenanceRepo) Create(ctx context.Context, m *models.Maintenance) error {
	query := `INSERT INTO maintenance (product_unit_id, issue_description, cost, start_date, status, created_by, modified_by) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	return r.db.QueryRowContext(ctx, query, m.ProductUnitID, m.IssueDescription, m.Cost, m.StartDate, m.Status, m.CreatedBy, m.ModifiedBy).Scan(&m.ID)
}

func (r *maintenanceRepo) CreateWithTx(ctx context.Context, tx *sql.Tx, m *models.Maintenance) error {
	query := `INSERT INTO maintenance (product_unit_id, issue_description, cost, start_date, status, created_by, modified_by) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	return tx.QueryRowContext(ctx, query, m.ProductUnitID, m.IssueDescription, m.Cost, m.StartDate, m.Status, m.CreatedBy, m.ModifiedBy).Scan(&m.ID)
}

func (r *maintenanceRepo) Update(ctx context.Context, id int, req models.MaintenanceUpdateDTO, modifiedBy int) error {
	query := `
		UPDATE maintenance 
		SET 
			issue_description = COALESCE($1, issue_description),
			cost              = COALESCE($2, cost),
			start_date        = COALESCE($3, start_date),
			status            = COALESCE($4, status),
			modified_at       = NOW(),
			modified_by       = $5
		WHERE id = $6 AND deleted_at IS NULL`

	res, err := r.db.ExecContext(ctx, query, req.IssueDescription, req.Cost, req.StartDate, req.Status, modifiedBy, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *maintenanceRepo) Complete(ctx context.Context, id int, modifiedBy int) error {
	query := `UPDATE maintenance SET status = 'completed', completed_date = NOW(), modified_at = NOW(), modified_by = $1 WHERE id = $2 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, modifiedBy, id)
	return err
}

func (r *maintenanceRepo) Delete(ctx context.Context, id, deletedBy int) error {
	res, err := r.db.ExecContext(ctx, `UPDATE maintenance SET deleted_at = NOW(), deleted_by = $1 WHERE id = $2 AND deleted_at IS NULL`, deletedBy, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *maintenanceRepo) Restore(ctx context.Context, id, restoredBy int) error {
	query := `UPDATE maintenance SET deleted_at = NULL, deleted_by = NULL, modified_at = NOW(), modified_by = $1 WHERE id = $2 AND deleted_at IS NOT NULL`
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

func (r *maintenanceRepo) ForceDelete(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM maintenance WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
