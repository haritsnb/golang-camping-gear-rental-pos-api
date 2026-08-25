package repositories

import (
	"app/models"
	"context"
	"database/sql"
)

type ProductRepository interface {
	GetAll(ctx context.Context) ([]models.Product, error)
	GetByID(ctx context.Context, id int) (*models.Product, error)
	GetDeleted(ctx context.Context) ([]models.Product, error)
	Create(ctx context.Context, p *models.Product) error
	Update(ctx context.Context, id int, req models.ProductUpdateDTO, modifiedBy int) error
	Delete(ctx context.Context, id, deletedBy int) error
	Restore(ctx context.Context, id, restoredBy int) error
	ForceDelete(ctx context.Context, id int) error
}

type productRepo struct{ db *sql.DB }

func NewProductRepository(db *sql.DB) ProductRepository { return &productRepo{db: db} }

func (r *productRepo) GetAll(ctx context.Context) ([]models.Product, error) {
	query := `
		SELECT 
			p.id, p.category_id, c.name, p.brand_id, b.name, p.name, 
			p.rental_price_per_day, p.default_deposit, p.late_fee_per_day, 
			p.lost_compensation_fee, p.is_active,
			COUNT(pu.id) FILTER (WHERE pu.deleted_at IS NULL) AS total_units,
			COUNT(pu.id) FILTER (WHERE pu.deleted_at IS NULL AND pu.status = 'available') AS available_units,
			COUNT(pu.id) FILTER (WHERE pu.deleted_at IS NULL AND pu.status = 'rented') AS rented_units,
			COUNT(pu.id) FILTER (WHERE pu.deleted_at IS NULL AND pu.status = 'maintenance') AS maintenance_units,
			p.created_at, p.modified_at
		FROM products p
		JOIN categories c ON c.id = p.category_id
		JOIN brands b ON b.id = p.brand_id
		LEFT JOIN product_units pu ON pu.product_id = p.id AND pu.deleted_at IS NULL
		WHERE p.deleted_at IS NULL 
		GROUP BY p.id, c.name, b.name
		ORDER BY p.id DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(
			&p.ID, &p.CategoryID, &p.CategoryName, &p.BrandID, &p.BrandName, &p.Name,
			&p.RentalPricePerDay, &p.DefaultDeposit, &p.LateFeePerDay, &p.LostCompensationFee, &p.IsActive,
			&p.TotalUnits, &p.AvailableUnits, &p.RentedUnits, &p.MaintenanceUnits,
			&p.CreatedAt, &p.ModifiedAt,
		); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *productRepo) GetByID(ctx context.Context, id int) (*models.Product, error) {
	query := `
		SELECT 
			p.id, p.category_id, c.name, p.brand_id, b.name, p.name, 
			p.rental_price_per_day, p.default_deposit, p.late_fee_per_day, 
			p.lost_compensation_fee, p.is_active,
			COUNT(pu.id) FILTER (WHERE pu.deleted_at IS NULL) AS total_units,
			COUNT(pu.id) FILTER (WHERE pu.deleted_at IS NULL AND pu.status = 'available') AS available_units,
			COUNT(pu.id) FILTER (WHERE pu.deleted_at IS NULL AND pu.status = 'rented') AS rented_units,
			COUNT(pu.id) FILTER (WHERE pu.deleted_at IS NULL AND pu.status = 'maintenance') AS maintenance_units,
			p.created_at, p.modified_at
		FROM products p
		JOIN categories c ON c.id = p.category_id
		JOIN brands b ON b.id = p.brand_id
		LEFT JOIN product_units pu ON pu.product_id = p.id AND pu.deleted_at IS NULL
		WHERE p.id = $1 AND p.deleted_at IS NULL
		GROUP BY p.id, c.name, b.name`

	p := &models.Product{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.CategoryID, &p.CategoryName, &p.BrandID, &p.BrandName, &p.Name,
		&p.RentalPricePerDay, &p.DefaultDeposit, &p.LateFeePerDay, &p.LostCompensationFee, &p.IsActive,
		&p.TotalUnits, &p.AvailableUnits, &p.RentedUnits, &p.MaintenanceUnits,
		&p.CreatedAt, &p.ModifiedAt,
	)
	return p, err
}

func (r *productRepo) GetDeleted(ctx context.Context) ([]models.Product, error) {
	query := `
		SELECT p.id, p.category_id, c.name, p.brand_id, b.name, p.name, p.rental_price_per_day, p.default_deposit, p.late_fee_per_day, p.lost_compensation_fee, p.is_active, 0, 0, 0, 0, p.created_at, p.modified_at, p.deleted_at
		FROM products p
		JOIN categories c ON c.id = p.category_id
		JOIN brands b ON b.id = p.brand_id
		WHERE p.deleted_at IS NOT NULL ORDER BY p.deleted_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.CategoryID, &p.CategoryName, &p.BrandID, &p.BrandName, &p.Name, &p.RentalPricePerDay, &p.DefaultDeposit, &p.LateFeePerDay, &p.LostCompensationFee, &p.IsActive, &p.TotalUnits, &p.AvailableUnits, &p.RentedUnits, &p.MaintenanceUnits, &p.CreatedAt, &p.ModifiedAt, &p.DeletedAt); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}

func (r *productRepo) Create(ctx context.Context, p *models.Product) error {
	query := `
		INSERT INTO products (category_id, brand_id, name, rental_price_per_day, default_deposit, late_fee_per_day, lost_compensation_fee, is_active, created_by, modified_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) 
		RETURNING id, created_at, modified_at`

	return r.db.QueryRowContext(ctx, query, p.CategoryID, p.BrandID, p.Name, p.RentalPricePerDay, p.DefaultDeposit, p.LateFeePerDay, p.LostCompensationFee, p.IsActive, p.CreatedBy, p.ModifiedBy).Scan(&p.ID, &p.CreatedAt, &p.ModifiedAt)
}

func (r *productRepo) Update(ctx context.Context, id int, req models.ProductUpdateDTO, modifiedBy int) error {
	query := `
		UPDATE products 
		SET 
			category_id           = COALESCE($1, category_id),
			brand_id              = COALESCE($2, brand_id),
			name                  = COALESCE($3, name),
			rental_price_per_day   = COALESCE($4, rental_price_per_day),
			default_deposit       = COALESCE($5, default_deposit),
			late_fee_per_day      = COALESCE($6, late_fee_per_day),
			lost_compensation_fee = COALESCE($7, lost_compensation_fee),
			is_active             = COALESCE($8, is_active),
			modified_at           = NOW(),
			modified_by           = $9
		WHERE id = $10 AND deleted_at IS NULL`

	res, err := r.db.ExecContext(ctx, query,
		req.CategoryID,
		req.BrandID,
		req.Name,
		req.RentalPricePerDay,
		req.DefaultDeposit,
		req.LateFeePerDay,
		req.LostCompensationFee,
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

func (r *productRepo) Delete(ctx context.Context, id, deletedBy int) error {
	res, err := r.db.ExecContext(ctx, `UPDATE products SET deleted_at = NOW(), deleted_by = $1 WHERE id = $2 AND deleted_at IS NULL`, deletedBy, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *productRepo) Restore(ctx context.Context, id, restoredBy int) error {
	query := `UPDATE products SET deleted_at = NULL, deleted_by = NULL, modified_at = NOW(), modified_by = $1 WHERE id = $2 AND deleted_at IS NOT NULL`
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

func (r *productRepo) ForceDelete(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM products WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
