package repositories

import (
	"app/models"
	"context"
	"database/sql"
)

type PaymentRepository interface {
	GetByRentalID(ctx context.Context, rentalID int) ([]models.Payment, error)
	GetByID(ctx context.Context, id int) (*models.Payment, error)
	GetDeleted(ctx context.Context) ([]models.Payment, error)
	CreateWithTx(ctx context.Context, tx *sql.Tx, p *models.Payment) error
	Create(ctx context.Context, p *models.Payment) error
	Update(ctx context.Context, id int, req models.PaymentUpdateDTO, modifiedBy int) error
	Delete(ctx context.Context, id, deletedBy int) error
	Restore(ctx context.Context, id, restoredBy int) error
	ForceDelete(ctx context.Context, id int) error
}

type paymentRepo struct{ db *sql.DB }

func NewPaymentRepository(db *sql.DB) PaymentRepository { return &paymentRepo{db: db} }

func (r *paymentRepo) GetByRentalID(ctx context.Context, rentalID int) ([]models.Payment, error) {
	query := `SELECT id, rental_id, user_id, payment_type, payment_method, amount, reference_number, notes, paid_at FROM payments WHERE rental_id = $1 AND deleted_at IS NULL ORDER BY id ASC`
	rows, err := r.db.QueryContext(ctx, query, rentalID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var payments []models.Payment
	for rows.Next() {
		var p models.Payment
		if err := rows.Scan(&p.ID, &p.RentalID, &p.UserID, &p.PaymentType, &p.PaymentMethod, &p.Amount, &p.ReferenceNumber, &p.Notes, &p.PaidAt); err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}
	return payments, nil
}

func (r *paymentRepo) GetByID(ctx context.Context, id int) (*models.Payment, error) {
	query := `SELECT id, rental_id, user_id, payment_type, payment_method, amount, reference_number, notes, paid_at FROM payments WHERE id = $1 AND deleted_at IS NULL`
	p := &models.Payment{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&p.ID, &p.RentalID, &p.UserID, &p.PaymentType, &p.PaymentMethod, &p.Amount, &p.ReferenceNumber, &p.Notes, &p.PaidAt)
	return p, err
}

func (r *paymentRepo) GetDeleted(ctx context.Context) ([]models.Payment, error) {
	query := `SELECT id, rental_id, user_id, payment_type, payment_method, amount, reference_number, notes, paid_at, deleted_at FROM payments WHERE deleted_at IS NOT NULL ORDER BY deleted_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.Payment
	for rows.Next() {
		var p models.Payment
		if err := rows.Scan(&p.ID, &p.RentalID, &p.UserID, &p.PaymentType, &p.PaymentMethod, &p.Amount, &p.ReferenceNumber, &p.Notes, &p.PaidAt, &p.DeletedAt); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}

func (r *paymentRepo) CreateWithTx(ctx context.Context, tx *sql.Tx, p *models.Payment) error {
	query := `INSERT INTO payments (rental_id, user_id, payment_type, payment_method, amount, reference_number, notes, paid_at, created_by, modified_by) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), $8, $9) RETURNING id`
	return tx.QueryRowContext(ctx, query, p.RentalID, p.UserID, p.PaymentType, p.PaymentMethod, p.Amount, p.ReferenceNumber, p.Notes, p.CreatedBy, p.ModifiedBy).Scan(&p.ID)
}

func (r *paymentRepo) Create(ctx context.Context, p *models.Payment) error {
	query := `INSERT INTO payments (rental_id, user_id, payment_type, payment_method, amount, reference_number, notes, paid_at, created_by, modified_by) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), $8, $9) RETURNING id`
	return r.db.QueryRowContext(ctx, query, p.RentalID, p.UserID, p.PaymentType, p.PaymentMethod, p.Amount, p.ReferenceNumber, p.Notes, p.CreatedBy, p.ModifiedBy).Scan(&p.ID)
}

func (r *paymentRepo) Update(ctx context.Context, id int, req models.PaymentUpdateDTO, modifiedBy int) error {
	query := `
		UPDATE payments 
		SET 
			payment_type     = COALESCE($1, payment_type),
			payment_method   = COALESCE($2, payment_method),
			amount           = COALESCE($3, amount),
			reference_number = COALESCE($4, reference_number),
			notes            = COALESCE($5, notes),
			modified_at      = NOW(),
			modified_by      = $6
		WHERE id = $7 AND deleted_at IS NULL`

	res, err := r.db.ExecContext(ctx, query, req.PaymentType, req.PaymentMethod, req.Amount, req.ReferenceNumber, req.Notes, modifiedBy, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *paymentRepo) Delete(ctx context.Context, id, deletedBy int) error {
	res, err := r.db.ExecContext(ctx, `UPDATE payments SET deleted_at = NOW(), deleted_by = $1 WHERE id = $2 AND deleted_at IS NULL`, deletedBy, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *paymentRepo) Restore(ctx context.Context, id, restoredBy int) error {
	query := `UPDATE payments SET deleted_at = NULL, deleted_by = NULL, modified_at = NOW(), modified_by = $1 WHERE id = $2 AND deleted_at IS NOT NULL`
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

func (r *paymentRepo) ForceDelete(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM payments WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
