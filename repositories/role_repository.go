package repositories

import (
	"app/models"
	"context"
	"database/sql"
)

type RoleRepository interface {
	GetAll(ctx context.Context) ([]models.Role, error)
	GetByID(ctx context.Context, id int) (*models.Role, error)
	GetByName(ctx context.Context, name string) (*models.Role, error)
	GetDeleted(ctx context.Context) ([]models.Role, error)
	Create(ctx context.Context, role *models.Role) error
	Update(ctx context.Context, id int, req models.RoleUpdateDTO, modifiedBy int) error
	Delete(ctx context.Context, id, deletedBy int) error
	Restore(ctx context.Context, id, restoredBy int) error
	ForceDelete(ctx context.Context, id int) error
}

type roleRepo struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) RoleRepository {
	return &roleRepo{db: db}
}

func (r *roleRepo) GetAll(ctx context.Context) ([]models.Role, error) {
	query := `SELECT id, name, description, created_at, modified_at FROM roles WHERE deleted_at IS NULL ORDER BY id ASC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var role models.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.ModifiedAt); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (r *roleRepo) GetByID(ctx context.Context, id int) (*models.Role, error) {
	role := &models.Role{}
	err := r.db.QueryRowContext(ctx, `SELECT id, name, description, created_at, modified_at FROM roles WHERE id = $1 AND deleted_at IS NULL`, id).
		Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.ModifiedAt)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (r *roleRepo) GetByName(ctx context.Context, name string) (*models.Role, error) {
	role := &models.Role{}
	err := r.db.QueryRowContext(ctx, `SELECT id, name, description, created_at, modified_at FROM roles WHERE name = $1 AND deleted_at IS NULL`, name).
		Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.ModifiedAt)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (r *roleRepo) GetDeleted(ctx context.Context) ([]models.Role, error) {
	query := `SELECT id, name, description, created_at, modified_at, deleted_at FROM roles WHERE deleted_at IS NOT NULL ORDER BY deleted_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.Role
	for rows.Next() {
		var role models.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.ModifiedAt, &role.DeletedAt); err != nil {
			return nil, err
		}
		list = append(list, role)
	}
	return list, nil
}

func (r *roleRepo) Create(ctx context.Context, role *models.Role) error {
	query := `INSERT INTO roles (name, description, created_by, modified_by) VALUES ($1, $2, $3, $4) RETURNING id, created_at, modified_at`
	return r.db.QueryRowContext(ctx, query, role.Name, role.Description, role.CreatedBy, role.ModifiedBy).Scan(&role.ID, &role.CreatedAt, &role.ModifiedAt)
}

func (r *roleRepo) Update(ctx context.Context, id int, req models.RoleUpdateDTO, modifiedBy int) error {
	query := `
		UPDATE roles 
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

func (r *roleRepo) Delete(ctx context.Context, id, deletedBy int) error {
	query := `UPDATE roles SET deleted_at = NOW(), deleted_by = $1 WHERE id = $2 AND deleted_at IS NULL`
	res, err := r.db.ExecContext(ctx, query, deletedBy, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *roleRepo) Restore(ctx context.Context, id, restoredBy int) error {
	query := `
		UPDATE roles 
		SET deleted_at = NULL, deleted_by = NULL, modified_at = NOW(), modified_by = $1 
		WHERE id = $2 AND deleted_at IS NOT NULL`

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

func (r *roleRepo) ForceDelete(ctx context.Context, id int) error {
	query := `DELETE FROM roles WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
