package repositories

import (
	"app/models"
	"context"
	"database/sql"
)

type BrandRepository interface {
	GetAll(ctx context.Context) ([]models.Brand, error)
	GetByID(ctx context.Context, id int) (*models.Brand, error)
	GetByName(ctx context.Context, name string) (*models.Brand, error)
	GetDeleted(ctx context.Context) ([]models.Brand, error)
	Create(ctx context.Context, b *models.Brand) error
	Update(ctx context.Context, id int, req models.BrandUpdateDTO, modifiedBy int) error
	Delete(ctx context.Context, id, deletedBy int) error
	Restore(ctx context.Context, id, restoredBy int) error
	ForceDelete(ctx context.Context, id int) error
}

type brandRepo struct{ db *sql.DB }

func NewBrandRepository(db *sql.DB) BrandRepository { return &brandRepo{db: db} }

func (r *brandRepo) GetAll(ctx context.Context) ([]models.Brand, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, description, is_active, created_at, modified_at FROM brands WHERE deleted_at IS NULL ORDER BY name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.Brand
	for rows.Next() {
		var b models.Brand
		if err := rows.Scan(&b.ID, &b.Name, &b.Description, &b.IsActive, &b.CreatedAt, &b.ModifiedAt); err != nil {
			return nil, err
		}
		list = append(list, b)
	}
	return list, nil
}

func (r *brandRepo) GetByID(ctx context.Context, id int) (*models.Brand, error) {
	b := &models.Brand{}
	err := r.db.QueryRowContext(ctx, `SELECT id, name, description, is_active, created_at, modified_at FROM brands WHERE id = $1 AND deleted_at IS NULL`, id).
		Scan(&b.ID, &b.Name, &b.Description, &b.IsActive, &b.CreatedAt, &b.ModifiedAt)
	return b, err
}

func (r *brandRepo) GetByName(ctx context.Context, name string) (*models.Brand, error) {
	b := &models.Brand{}
	err := r.db.QueryRowContext(ctx, `SELECT id, name, description, is_active, created_at, modified_at FROM brands WHERE name = $1 AND deleted_at IS NULL`, name).
		Scan(&b.ID, &b.Name, &b.Description, &b.IsActive, &b.CreatedAt, &b.ModifiedAt)
	return b, err
}

func (r *brandRepo) GetDeleted(ctx context.Context) ([]models.Brand, error) {
	query := `SELECT id, name, description, is_active, created_at, modified_at, deleted_at FROM brands WHERE deleted_at IS NOT NULL ORDER BY deleted_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.Brand
	for rows.Next() {
		var b models.Brand
		if err := rows.Scan(&b.ID, &b.Name, &b.Description, &b.IsActive, &b.CreatedAt, &b.ModifiedAt, &b.DeletedAt); err != nil {
			return nil, err
		}
		list = append(list, b)
	}
	return list, nil
}

func (r *brandRepo) Create(ctx context.Context, b *models.Brand) error {
	return r.db.QueryRowContext(ctx, `INSERT INTO brands (name, description, is_active, created_by, modified_by) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, modified_at`,
		b.Name, b.Description, b.IsActive, b.CreatedBy, b.ModifiedBy).Scan(&b.ID, &b.CreatedAt, &b.ModifiedAt)
}

func (r *brandRepo) Update(ctx context.Context, id int, req models.BrandUpdateDTO, modifiedBy int) error {
	query := `
		UPDATE brands 
		SET 
			name        = COALESCE($1, name),
			description = COALESCE($2, description),
			is_active   = COALESCE($3, is_active),
			modified_at = NOW(),
			modified_by = $4
		WHERE id = $5 AND deleted_at IS NULL`

	res, err := r.db.ExecContext(ctx, query, req.Name, req.Description, req.IsActive, modifiedBy, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *brandRepo) Delete(ctx context.Context, id, deletedBy int) error {
	res, err := r.db.ExecContext(ctx, `UPDATE brands SET deleted_at = NOW(), deleted_by = $1 WHERE id = $2 AND deleted_at IS NULL`, deletedBy, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *brandRepo) Restore(ctx context.Context, id, restoredBy int) error {
	query := `UPDATE brands SET deleted_at = NULL, deleted_by = NULL, modified_at = NOW(), modified_by = $1 WHERE id = $2 AND deleted_at IS NOT NULL`
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

func (r *brandRepo) ForceDelete(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM brands WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
