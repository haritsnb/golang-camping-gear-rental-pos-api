package repositories

import (
	"app/models"
	"context"
	"database/sql"
)

type CategoryRepository interface {
	GetAll(ctx context.Context) ([]models.Category, error)
	GetByID(ctx context.Context, id int) (*models.Category, error)
	GetByName(ctx context.Context, name string) (*models.Category, error)
	GetDeleted(ctx context.Context) ([]models.Category, error)
	Create(ctx context.Context, c *models.Category) error
	Update(ctx context.Context, id int, req models.CategoryUpdateDTO, modifiedBy int) error
	Delete(ctx context.Context, id, deletedBy int) error
	Restore(ctx context.Context, id, restoredBy int) error
	ForceDelete(ctx context.Context, id int) error
}

type categoryRepo struct{ db *sql.DB }

func NewCategoryRepository(db *sql.DB) CategoryRepository { return &categoryRepo{db: db} }

func (r *categoryRepo) GetAll(ctx context.Context) ([]models.Category, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, description, created_at, modified_at FROM categories WHERE deleted_at IS NULL ORDER BY name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.ModifiedAt); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, nil
}

func (r *categoryRepo) GetByID(ctx context.Context, id int) (*models.Category, error) {
	c := &models.Category{}
	err := r.db.QueryRowContext(ctx, `SELECT id, name, description, created_at, modified_at FROM categories WHERE id = $1 AND deleted_at IS NULL`, id).
		Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.ModifiedAt)
	if err != nil {
		return nil, err // <-- WAJIB return nil
	}
	return c, nil
}

func (r *categoryRepo) GetByName(ctx context.Context, name string) (*models.Category, error) {
	c := &models.Category{}
	err := r.db.QueryRowContext(ctx, `SELECT id, name, description, created_at, modified_at FROM categories WHERE name = $1 AND deleted_at IS NULL`, name).
		Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.ModifiedAt)
	if err != nil {
		return nil, err // <-- WAJIB return nil
	}
	return c, nil
}

func (r *categoryRepo) GetDeleted(ctx context.Context) ([]models.Category, error) {
	query := `SELECT id, name, description, created_at, modified_at, deleted_at FROM categories WHERE deleted_at IS NOT NULL ORDER BY deleted_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt, &c.ModifiedAt, &c.DeletedAt); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, nil
}

func (r *categoryRepo) Create(ctx context.Context, c *models.Category) error {
	return r.db.QueryRowContext(ctx, `INSERT INTO categories (name, description, created_by, modified_by) VALUES ($1, $2, $3, $4) RETURNING id, created_at, modified_at`,
		c.Name, c.Description, c.CreatedBy, c.ModifiedBy).Scan(&c.ID, &c.CreatedAt, &c.ModifiedAt)
}

func (r *categoryRepo) Update(ctx context.Context, id int, req models.CategoryUpdateDTO, modifiedBy int) error {
	query := `
		UPDATE categories 
		SET 
			name        = COALESCE($1, name),
			description = COALESCE($2, description),
			modified_at = NOW(),
			modified_by = $3
		WHERE id = $4 AND deleted_at IS NULL`

	res, err := r.db.ExecContext(ctx, query, req.Name, req.Description, modifiedBy, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *categoryRepo) Delete(ctx context.Context, id, deletedBy int) error {
	res, err := r.db.ExecContext(ctx, `UPDATE categories SET deleted_at = NOW(), deleted_by = $1 WHERE id = $2 AND deleted_at IS NULL`, deletedBy, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *categoryRepo) Restore(ctx context.Context, id, restoredBy int) error {
	query := `UPDATE categories SET deleted_at = NULL, deleted_by = NULL, modified_at = NOW(), modified_by = $1 WHERE id = $2 AND deleted_at IS NOT NULL`
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

func (r *categoryRepo) ForceDelete(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM categories WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
