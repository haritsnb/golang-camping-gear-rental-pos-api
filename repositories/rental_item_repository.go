package repositories

import (
	"app/models"
	"context"
	"database/sql"
)

type RentalItemRepository interface {
	GetByRentalID(ctx context.Context, rentalID int) ([]models.RentalItem, error)
	GetByRentalIDWithTx(ctx context.Context, tx *sql.Tx, rentalID int) ([]models.RentalItem, error)
	CreateWithTx(ctx context.Context, tx *sql.Tx, item *models.RentalItem) error
	UpdateReturnWithTx(ctx context.Context, tx *sql.Tx, item *models.RentalItem) error
	DeleteByRentalIDWithTx(ctx context.Context, tx *sql.Tx, rentalID, deletedBy int) error
	RestoreByRentalIDWithTx(ctx context.Context, tx *sql.Tx, rentalID, restoredBy int) error
	ForceDeleteByRentalIDWithTx(ctx context.Context, tx *sql.Tx, rentalID int) error
}

type rentalItemRepo struct{ db *sql.DB }

func NewRentalItemRepository(db *sql.DB) RentalItemRepository { return &rentalItemRepo{db: db} }

func (r *rentalItemRepo) GetByRentalID(ctx context.Context, rentalID int) ([]models.RentalItem, error) {
	query := `
		SELECT ri.id, ri.rental_id, ri.product_unit_id, pu.unit_code, p.name, ri.daily_rate, ri.duration_days, ri.subtotal, ri.condition_out, ri.condition_in, ri.item_penalty_fee, ri.return_notes
		FROM rental_items ri
		JOIN product_units pu ON pu.id = ri.product_unit_id
		JOIN products p ON p.id = pu.product_id
		WHERE ri.rental_id = $1 AND ri.deleted_at IS NULL`
	rows, err := r.db.QueryContext(ctx, query, rentalID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.RentalItem
	for rows.Next() {
		var it models.RentalItem
		if err := rows.Scan(&it.ID, &it.RentalID, &it.ProductUnitID, &it.UnitCode, &it.ProductName, &it.DailyRate, &it.DurationDays, &it.Subtotal, &it.ConditionOut, &it.ConditionIn, &it.ItemPenaltyFee, &it.ReturnNotes); err != nil {
			return nil, err
		}
		items = append(items, it)
	}
	return items, nil
}

func (r *rentalItemRepo) GetByRentalIDWithTx(ctx context.Context, tx *sql.Tx, rentalID int) ([]models.RentalItem, error) {
	query := `
		SELECT ri.id, ri.rental_id, ri.product_unit_id, pu.unit_code, p.name, ri.daily_rate, ri.duration_days, ri.subtotal, ri.condition_out, ri.condition_in, ri.item_penalty_fee, ri.return_notes
		FROM rental_items ri
		JOIN product_units pu ON pu.id = ri.product_unit_id
		JOIN products p ON p.id = pu.product_id
		WHERE ri.rental_id = $1 AND ri.deleted_at IS NULL FOR UPDATE`
	rows, err := tx.QueryContext(ctx, query, rentalID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.RentalItem
	for rows.Next() {
		var it models.RentalItem
		if err := rows.Scan(&it.ID, &it.RentalID, &it.ProductUnitID, &it.UnitCode, &it.ProductName, &it.DailyRate, &it.DurationDays, &it.Subtotal, &it.ConditionOut, &it.ConditionIn, &it.ItemPenaltyFee, &it.ReturnNotes); err != nil {
			return nil, err
		}
		items = append(items, it)
	}
	return items, nil
}

func (r *rentalItemRepo) CreateWithTx(ctx context.Context, tx *sql.Tx, item *models.RentalItem) error {
	query := `
		INSERT INTO rental_items (rental_id, product_unit_id, daily_rate, duration_days, subtotal, condition_out, created_by, modified_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	return tx.QueryRowContext(ctx, query,
		item.RentalID, item.ProductUnitID, item.DailyRate, item.DurationDays, item.Subtotal, item.ConditionOut, item.CreatedBy, item.ModifiedBy,
	).Scan(&item.ID)
}

func (r *rentalItemRepo) UpdateReturnWithTx(ctx context.Context, tx *sql.Tx, item *models.RentalItem) error {
	query := `
		UPDATE rental_items 
		SET condition_in = $1, item_penalty_fee = $2, return_notes = $3, modified_at = NOW(), modified_by = $4
		WHERE id = $5`
	_, err := tx.ExecContext(ctx, query, item.ConditionIn, item.ItemPenaltyFee, item.ReturnNotes, item.ModifiedBy, item.ID)
	return err
}

func (r *rentalItemRepo) DeleteByRentalIDWithTx(ctx context.Context, tx *sql.Tx, rentalID, deletedBy int) error {
	_, err := tx.ExecContext(ctx, `UPDATE rental_items SET deleted_at = NOW(), deleted_by = $1 WHERE rental_id = $2 AND deleted_at IS NULL`, deletedBy, rentalID)
	return err
}

func (r *rentalItemRepo) RestoreByRentalIDWithTx(ctx context.Context, tx *sql.Tx, rentalID, restoredBy int) error {
	_, err := tx.ExecContext(ctx, `UPDATE rental_items SET deleted_at = NULL, deleted_by = NULL, modified_at = NOW(), modified_by = $1 WHERE rental_id = $2 AND deleted_at IS NOT NULL`, restoredBy, rentalID)
	return err
}

func (r *rentalItemRepo) ForceDeleteByRentalIDWithTx(ctx context.Context, tx *sql.Tx, rentalID int) error {
	_, err := tx.ExecContext(ctx, `DELETE FROM rental_items WHERE rental_id = $1`, rentalID)
	return err
}
