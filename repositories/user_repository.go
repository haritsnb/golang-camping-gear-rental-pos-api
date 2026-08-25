package repositories

import (
	"app/models"
	"context"
	"database/sql"
)

type UserRepository interface {
	GetAll(ctx context.Context) ([]models.User, error)
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetDeleted(ctx context.Context) ([]models.User, error)
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, id int, req models.UserUpdateDTO, passwordHash *string, modifiedBy int) error
	Delete(ctx context.Context, id, deletedBy int) error
	Restore(ctx context.Context, id, restoredBy int) error // <-- Baru
	ForceDelete(ctx context.Context, id int) error         // <-- Baru
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) GetAll(ctx context.Context) ([]models.User, error) {
	query := `
		SELECT 
			u.id, 
			u.role_id, 
			r.name, 
			u.name, 
			u.username, 
			COALESCE(u.phone, ''), 
			u.is_active, 
			u.created_at, 
			u.modified_at
		FROM users u
		JOIN roles r ON r.id = u.role_id
		WHERE u.deleted_at IS NULL 
		ORDER BY u.id DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		var phone string
		if err := rows.Scan(&u.ID, &u.RoleID, &u.RoleName, &u.Name, &u.Username, &phone, &u.IsActive, &u.CreatedAt, &u.ModifiedAt); err != nil {
			return nil, err
		}
		if phone != "" {
			u.Phone = &phone
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *userRepo) GetByID(ctx context.Context, id int) (*models.User, error) {
	query := `
		SELECT 
			u.id, 
			u.role_id, 
			r.name, 
			u.name, 
			u.username, 
			u.password_hash, 
			COALESCE(u.phone, ''), 
			u.is_active, 
			u.created_at, 
			u.modified_at
		FROM users u
		JOIN roles r ON r.id = u.role_id
		WHERE u.id = $1 AND u.deleted_at IS NULL`

	u := &models.User{}
	var phone string
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.RoleID, &u.RoleName, &u.Name, &u.Username, &u.PasswordHash, &phone, &u.IsActive, &u.CreatedAt, &u.ModifiedAt,
	)
	if err != nil {
		return nil, err
	}
	if phone != "" {
		u.Phone = &phone
	}
	return u, nil
}

func (r *userRepo) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
		SELECT 
			u.id, 
			u.role_id, 
			r.name, 
			u.name, 
			u.username, 
			u.password_hash, 
			COALESCE(u.phone, ''), 
			u.is_active, 
			u.created_at, 
			u.modified_at
		FROM users u
		JOIN roles r ON r.id = u.role_id
		WHERE u.username = $1 AND u.deleted_at IS NULL`

	u := &models.User{}
	var phone string
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&u.ID, &u.RoleID, &u.RoleName, &u.Name, &u.Username, &u.PasswordHash, &phone, &u.IsActive, &u.CreatedAt, &u.ModifiedAt,
	)
	if err != nil {
		return nil, err
	}
	if phone != "" {
		u.Phone = &phone
	}
	return u, nil
}

func (r *userRepo) GetDeleted(ctx context.Context) ([]models.User, error) {
	query := `
		SELECT u.id, u.role_id, r.name, u.name, u.username, COALESCE(u.phone, ''), u.is_active, u.created_at, u.modified_at, u.deleted_at
		FROM users u
		JOIN roles r ON r.id = u.role_id
		WHERE u.deleted_at IS NOT NULL ORDER BY u.deleted_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.User
	for rows.Next() {
		var u models.User
		var phone string
		if err := rows.Scan(&u.ID, &u.RoleID, &u.RoleName, &u.Name, &u.Username, &phone, &u.IsActive, &u.CreatedAt, &u.ModifiedAt, &u.DeletedAt); err != nil {
			return nil, err
		}
		if phone != "" {
			u.Phone = &phone
		}
		list = append(list, u)
	}
	return list, nil
}

func (r *userRepo) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (role_id, name, username, password_hash, phone, is_active, created_by, modified_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
		RETURNING id, created_at, modified_at`

	return r.db.QueryRowContext(ctx, query,
		user.RoleID, user.Name, user.Username, user.PasswordHash, user.Phone, user.IsActive, user.CreatedBy, user.ModifiedBy,
	).Scan(&user.ID, &user.CreatedAt, &user.ModifiedAt)
}

func (r *userRepo) Update(ctx context.Context, id int, req models.UserUpdateDTO, passwordHash *string, modifiedBy int) error {
	query := `
		UPDATE users 
		SET 
			role_id       = COALESCE($1, role_id),
			name          = COALESCE($2, name),
			username      = COALESCE($3, username),
			password_hash = COALESCE($4, password_hash),
			phone         = COALESCE($5, phone),
			is_active     = COALESCE($6, is_active),
			modified_at   = NOW(),
			modified_by   = $7
		WHERE id = $8 AND deleted_at IS NULL`

	res, err := r.db.ExecContext(ctx, query,
		req.RoleID,
		req.Name,
		req.Username,
		passwordHash,
		req.Phone,
		req.IsActive,
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

func (r *userRepo) Delete(ctx context.Context, id, deletedBy int) error {
	query := `UPDATE users SET deleted_at = NOW(), deleted_by = $1 WHERE id = $2 AND deleted_at IS NULL`
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

func (r *userRepo) Restore(ctx context.Context, id, restoredBy int) error {
	query := `
		UPDATE users 
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

func (r *userRepo) ForceDelete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`
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
