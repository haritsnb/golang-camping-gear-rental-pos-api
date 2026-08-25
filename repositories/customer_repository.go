package repositories

import (
	"app/models"
	"context"
	"database/sql"
)

type CustomerRepository interface {
	GetAll(ctx context.Context) ([]models.Customer, error)
	GetByID(ctx context.Context, id int) (*models.Customer, error)
	GetByIdentityNumber(ctx context.Context, identityNumber string) (*models.Customer, error)
	GetDeleted(ctx context.Context) ([]models.Customer, error)
	Create(ctx context.Context, c *models.Customer) error
	Update(ctx context.Context, id int, req models.CustomerUpdateDTO, photoURL *string, modifiedBy int) error
	SetBlacklist(ctx context.Context, id int, isBlacklist bool, modifiedBy int) error
	Delete(ctx context.Context, id, deletedBy int) error
	Restore(ctx context.Context, id, restoredBy int) error
	ForceDelete(ctx context.Context, id int) error
}

type customerRepo struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) CustomerRepository {
	return &customerRepo{db: db}
}

func (r *customerRepo) GetAll(ctx context.Context) ([]models.Customer, error) {
	query := `
		SELECT 
			id, name, identity_type, identity_number, identity_photo_url, phone, 
			emergency_contact, email, address, is_blacklisted, notes, created_at, modified_at
		FROM customers 
		WHERE deleted_at IS NULL 
		ORDER BY id DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []models.Customer
	for rows.Next() {
		var c models.Customer
		if err := rows.Scan(
			&c.ID, &c.Name, &c.IdentityType, &c.IdentityNumber, &c.IdentityPhotoURL,
			&c.Phone, &c.EmergencyContact, &c.Email, &c.Address, &c.IsBlacklisted,
			&c.Notes, &c.CreatedAt, &c.ModifiedAt,
		); err != nil {
			return nil, err
		}
		customers = append(customers, c)
	}
	return customers, nil
}

func (r *customerRepo) GetByID(ctx context.Context, id int) (*models.Customer, error) {
	query := `
		SELECT 
			id, name, identity_type, identity_number, identity_photo_url, phone, 
			emergency_contact, email, address, is_blacklisted, notes, created_at, modified_at
		FROM customers 
		WHERE id = $1 AND deleted_at IS NULL`

	c := &models.Customer{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&c.ID, &c.Name, &c.IdentityType, &c.IdentityNumber, &c.IdentityPhotoURL,
		&c.Phone, &c.EmergencyContact, &c.Email, &c.Address, &c.IsBlacklisted,
		&c.Notes, &c.CreatedAt, &c.ModifiedAt,
	)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *customerRepo) GetByIdentityNumber(ctx context.Context, identityNumber string) (*models.Customer, error) {
	query := `
		SELECT 
			id, name, identity_type, identity_number, identity_photo_url, phone, 
			emergency_contact, email, address, is_blacklisted, notes, created_at, modified_at
		FROM customers 
		WHERE identity_number = $1 AND deleted_at IS NULL`

	c := &models.Customer{}
	err := r.db.QueryRowContext(ctx, query, identityNumber).Scan(
		&c.ID, &c.Name, &c.IdentityType, &c.IdentityNumber, &c.IdentityPhotoURL,
		&c.Phone, &c.EmergencyContact, &c.Email, &c.Address, &c.IsBlacklisted,
		&c.Notes, &c.CreatedAt, &c.ModifiedAt,
	)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *customerRepo) GetDeleted(ctx context.Context) ([]models.Customer, error) {
	query := `
		SELECT id, name, identity_type, identity_number, identity_photo_url, phone, emergency_contact, email, address, is_blacklisted, notes, created_at, modified_at, deleted_at
		FROM customers WHERE deleted_at IS NOT NULL ORDER BY deleted_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.Customer
	for rows.Next() {
		var c models.Customer
		if err := rows.Scan(&c.ID, &c.Name, &c.IdentityType, &c.IdentityNumber, &c.IdentityPhotoURL, &c.Phone, &c.EmergencyContact, &c.Email, &c.Address, &c.IsBlacklisted, &c.Notes, &c.CreatedAt, &c.ModifiedAt, &c.DeletedAt); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, nil
}

func (r *customerRepo) Create(ctx context.Context, c *models.Customer) error {
	query := `
		INSERT INTO customers (
			name, identity_type, identity_number, identity_photo_url, phone, 
			emergency_contact, email, address, notes, created_by, modified_by
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) 
		RETURNING id, created_at, modified_at`

	return r.db.QueryRowContext(ctx, query,
		c.Name, c.IdentityType, c.IdentityNumber, c.IdentityPhotoURL, c.Phone,
		c.EmergencyContact, c.Email, c.Address, c.Notes, c.CreatedBy, c.ModifiedBy,
	).Scan(&c.ID, &c.CreatedAt, &c.ModifiedAt)
}

// Update menggunakan COALESCE: jika parameter nil, data lama tetap dipertahankan
func (r *customerRepo) Update(ctx context.Context, id int, req models.CustomerUpdateDTO, photoURL *string, modifiedBy int) error {
	query := `
		UPDATE customers 
		SET 
			name               = COALESCE($1, name),
			identity_type      = COALESCE($2, identity_type),
			identity_number    = COALESCE($3, identity_number),
			identity_photo_url = COALESCE($4, identity_photo_url),
			phone              = COALESCE($5, phone),
			emergency_contact  = COALESCE($6, emergency_contact),
			email              = COALESCE($7, email),
			address            = COALESCE($8, address),
			is_blacklisted     = COALESCE($9, is_blacklisted),
			notes              = COALESCE($10, notes),
			modified_at        = NOW(),
			modified_by        = $11
		WHERE id = $12 AND deleted_at IS NULL`

	res, err := r.db.ExecContext(ctx, query,
		req.Name,
		req.IdentityType,
		req.IdentityNumber,
		photoURL,
		req.Phone,
		req.EmergencyContact,
		req.Email,
		req.Address,
		req.IsBlacklisted,
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

func (r *customerRepo) SetBlacklist(ctx context.Context, id int, isBlacklist bool, modifiedBy int) error {
	query := `UPDATE customers SET is_blacklisted = $1, modified_at = NOW(), modified_by = $2 WHERE id = $3 AND deleted_at IS NULL`
	res, err := r.db.ExecContext(ctx, query, isBlacklist, modifiedBy, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *customerRepo) Delete(ctx context.Context, id, deletedBy int) error {
	query := `UPDATE customers SET deleted_at = NOW(), deleted_by = $1 WHERE id = $2 AND deleted_at IS NULL`
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

func (r *customerRepo) Restore(ctx context.Context, id, restoredBy int) error {
	query := `
		UPDATE customers 
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

func (r *customerRepo) ForceDelete(ctx context.Context, id int) error {
	query := `DELETE FROM customers WHERE id = $1`
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
