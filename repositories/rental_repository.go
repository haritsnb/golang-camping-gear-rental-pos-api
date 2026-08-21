package repositories

import (
	"app/models"
	"context"
	"database/sql"
)

type RentalRepository interface {
	GetAll(ctx context.Context, status string) ([]models.Rental, error)
	GetByID(ctx context.Context, id int) (*models.Rental, error)
	GetByIDWithTx(ctx context.Context, tx *sql.Tx, id int) (*models.Rental, error)
	CreateWithTx(ctx context.Context, tx *sql.Tx, r *models.Rental) error
	UpdateWithTx(ctx context.Context, tx *sql.Tx, r *models.Rental) error
	UpdateFlexible(ctx context.Context, id int, req models.RentalUpdateDTO, modifiedBy int) error
	DeleteWithTx(ctx context.Context, tx *sql.Tx, id, deletedBy int) error
	RestoreWithTx(ctx context.Context, tx *sql.Tx, id, restoredBy int) error
	ForceDeleteWithTx(ctx context.Context, tx *sql.Tx, id int) error
}

type rentalRepo struct{ db *sql.DB }

func NewRentalRepository(db *sql.DB) RentalRepository { return &rentalRepo{db: db} }

func (r *rentalRepo) GetAll(ctx context.Context, status string) ([]models.Rental, error) {
	query := `
		SELECT r.id, r.invoice_number, r.customer_id, c.name, r.user_id, u.name, r.booking_date, r.start_date, r.expected_return_date, r.actual_return_date, r.total_rental_fee, r.total_deposit, r.total_penalty_fee, r.grand_total, r.status, r.created_at, r.modified_at
		FROM rentals r
		JOIN customers c ON c.id = r.customer_id
		JOIN users u ON u.id = r.user_id
		WHERE r.deleted_at IS NULL`

	var args []interface{}
	if status != "" {
		query += " AND r.status = $1"
		args = append(args, status)
	}
	query += " ORDER BY r.id DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rentals []models.Rental
	for rows.Next() {
		var ren models.Rental
		if err := rows.Scan(&ren.ID, &ren.InvoiceNumber, &ren.CustomerID, &ren.CustomerName, &ren.UserID, &ren.UserName, &ren.BookingDate, &ren.StartDate, &ren.ExpectedReturnDate, &ren.ActualReturnDate, &ren.TotalRentalFee, &ren.TotalDeposit, &ren.TotalPenaltyFee, &ren.GrandTotal, &ren.Status, &ren.CreatedAt, &ren.ModifiedAt); err != nil {
			return nil, err
		}
		rentals = append(rentals, ren)
	}
	return rentals, nil
}

func (r *rentalRepo) GetByID(ctx context.Context, id int) (*models.Rental, error) {
	query := `
		SELECT r.id, r.invoice_number, r.customer_id, c.name, r.user_id, u.name, r.booking_date, r.start_date, r.expected_return_date, r.actual_return_date, r.total_rental_fee, r.total_deposit, r.total_penalty_fee, r.grand_total, r.status, r.created_at, r.modified_at
		FROM rentals r
		JOIN customers c ON c.id = r.customer_id
		JOIN users u ON u.id = r.user_id
		WHERE r.id = $1 AND r.deleted_at IS NULL`
	ren := &models.Rental{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&ren.ID, &ren.InvoiceNumber, &ren.CustomerID, &ren.CustomerName, &ren.UserID, &ren.UserName, &ren.BookingDate, &ren.StartDate, &ren.ExpectedReturnDate, &ren.ActualReturnDate, &ren.TotalRentalFee, &ren.TotalDeposit, &ren.TotalPenaltyFee, &ren.GrandTotal, &ren.Status, &ren.CreatedAt, &ren.ModifiedAt,
	)
	return ren, err
}

func (r *rentalRepo) GetByIDWithTx(ctx context.Context, tx *sql.Tx, id int) (*models.Rental, error) {
	query := `
		SELECT r.id, r.invoice_number, r.customer_id, c.name, r.user_id, u.name, r.booking_date, r.start_date, r.expected_return_date, r.actual_return_date, r.total_rental_fee, r.total_deposit, r.total_penalty_fee, r.grand_total, r.status, r.created_at, r.modified_at
		FROM rentals r
		JOIN customers c ON c.id = r.customer_id
		JOIN users u ON u.id = r.user_id
		WHERE r.id = $1 AND r.deleted_at IS NULL FOR UPDATE`
	ren := &models.Rental{}
	err := tx.QueryRowContext(ctx, query, id).Scan(
		&ren.ID, &ren.InvoiceNumber, &ren.CustomerID, &ren.CustomerName, &ren.UserID, &ren.UserName, &ren.BookingDate, &ren.StartDate, &ren.ExpectedReturnDate, &ren.ActualReturnDate, &ren.TotalRentalFee, &ren.TotalDeposit, &ren.TotalPenaltyFee, &ren.GrandTotal, &ren.Status, &ren.CreatedAt, &ren.ModifiedAt,
	)
	return ren, err
}

func (r *rentalRepo) CreateWithTx(ctx context.Context, tx *sql.Tx, ren *models.Rental) error {
	query := `
		INSERT INTO rentals (invoice_number, customer_id, user_id, booking_date, start_date, expected_return_date, total_rental_fee, total_deposit, total_penalty_fee, grand_total, status, created_by, modified_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id`
	return tx.QueryRowContext(ctx, query,
		ren.InvoiceNumber, ren.CustomerID, ren.UserID, ren.BookingDate, ren.StartDate, ren.ExpectedReturnDate, ren.TotalRentalFee, ren.TotalDeposit, ren.TotalPenaltyFee, ren.GrandTotal, ren.Status, ren.CreatedBy, ren.ModifiedBy,
	).Scan(&ren.ID)
}

func (r *rentalRepo) UpdateWithTx(ctx context.Context, tx *sql.Tx, ren *models.Rental) error {
	query := `
		UPDATE rentals
		SET actual_return_date = $1, total_rental_fee = $2, total_deposit = $3, total_penalty_fee = $4, grand_total = $5, status = $6, modified_at = NOW(), modified_by = $7
		WHERE id = $8`
	_, err := tx.ExecContext(ctx, query,
		ren.ActualReturnDate, ren.TotalRentalFee, ren.TotalDeposit, ren.TotalPenaltyFee, ren.GrandTotal, ren.Status, ren.ModifiedBy, ren.ID,
	)
	return err
}

func (r *rentalRepo) UpdateFlexible(ctx context.Context, id int, req models.RentalUpdateDTO, modifiedBy int) error {
	query := `
		UPDATE rentals 
		SET 
			start_date           = COALESCE($1, start_date),
			expected_return_date = COALESCE($2, expected_return_date),
			total_rental_fee     = COALESCE($3, total_rental_fee),
			total_deposit        = COALESCE($4, total_deposit),
			total_penalty_fee    = COALESCE($5, total_penalty_fee),
			grand_total          = COALESCE($6, grand_total),
			status               = COALESCE($7, status),
			modified_at          = NOW(),
			modified_by          = $8
		WHERE id = $9 AND deleted_at IS NULL`

	res, err := r.db.ExecContext(ctx, query,
		req.StartDate, req.ExpectedReturnDate, req.TotalRentalFee, req.TotalDeposit,
		req.TotalPenaltyFee, req.GrandTotal, req.Status, modifiedBy, id,
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

func (r *rentalRepo) DeleteWithTx(ctx context.Context, tx *sql.Tx, id, deletedBy int) error {
	res, err := tx.ExecContext(ctx, `UPDATE rentals SET deleted_at = NOW(), deleted_by = $1 WHERE id = $2 AND deleted_at IS NULL`, deletedBy, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *rentalRepo) RestoreWithTx(ctx context.Context, tx *sql.Tx, id, restoredBy int) error {
	query := `UPDATE rentals SET deleted_at = NULL, deleted_by = NULL, modified_at = NOW(), modified_by = $1 WHERE id = $2 AND deleted_at IS NOT NULL`
	res, err := tx.ExecContext(ctx, query, restoredBy, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *rentalRepo) ForceDeleteWithTx(ctx context.Context, tx *sql.Tx, id int) error {
	res, err := tx.ExecContext(ctx, `DELETE FROM rentals WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
