package repositories

import (
	"app/models"
	"context"
	"database/sql"
)

type ReportRepository interface {
	GetRevenueReport(ctx context.Context, startDate, endDate string) (*models.RevenueReport, error)
	GetRentalSummaryReport(ctx context.Context, startDate, endDate string) (*models.RentalSummaryReport, error)
	GetTopProductsReport(ctx context.Context, startDate, endDate string, limit int) ([]models.TopProductReport, error)
	GetInventoryReport(ctx context.Context) (*models.InventoryReport, error)
}

type reportRepo struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) ReportRepository {
	return &reportRepo{db: db}
}

func (r *reportRepo) GetRevenueReport(ctx context.Context, startDate, endDate string) (*models.RevenueReport, error) {
	// Hitung Penerimaan Sewa & Denda dari Transaksi Rental Selesai (completed)
	queryRental := `
		SELECT 
			COALESCE(SUM(total_rental_fee), 0),
			COALESCE(SUM(total_penalty_fee), 0),
			COUNT(id)
		FROM rentals
		WHERE status = 'completed'
		  AND deleted_at IS NULL
		  AND DATE(booking_date) BETWEEN $1 AND $2`

	var rentalRevenue, penaltyRevenue float64
	var totalTrx int
	err := r.db.QueryRowContext(ctx, queryRental, startDate, endDate).Scan(&rentalRevenue, &penaltyRevenue, &totalTrx)
	if err != nil {
		return nil, err
	}

	// Hitung Pengeluaran Servis Alat (Maintenance Cost)
	queryMaint := `
		SELECT COALESCE(SUM(cost), 0)
		FROM maintenance
		WHERE deleted_at IS NULL
		  AND DATE(start_date) BETWEEN $1 AND $2`

	var maintCost float64
	err = r.db.QueryRowContext(ctx, queryMaint, startDate, endDate).Scan(&maintCost)
	if err != nil {
		return nil, err
	}

	gross := rentalRevenue + penaltyRevenue
	net := gross - maintCost

	return &models.RevenueReport{
		StartDate:            startDate,
		EndDate:              endDate,
		TotalRentalRevenue:   rentalRevenue,
		TotalPenaltyRevenue:  penaltyRevenue,
		TotalGrossRevenue:    gross,
		TotalMaintenanceCost: maintCost,
		TotalNetProfit:       net,
		TotalTransactions:    totalTrx,
	}, nil
}

func (r *reportRepo) GetRentalSummaryReport(ctx context.Context, startDate, endDate string) (*models.RentalSummaryReport, error) {
	query := `
		SELECT 
			COUNT(id) FILTER (WHERE status = 'booked') AS total_booked,
			COUNT(id) FILTER (WHERE status = 'active') AS total_active,
			COUNT(id) FILTER (WHERE status = 'returned') AS total_returned,
			COUNT(id) FILTER (WHERE status = 'completed') AS total_completed,
			COUNT(id) FILTER (WHERE status = 'cancelled') AS total_cancelled,
			COUNT(id) AS total_transaction
		FROM rentals
		WHERE deleted_at IS NULL
		  AND DATE(booking_date) BETWEEN $1 AND $2`

	rep := &models.RentalSummaryReport{
		StartDate: startDate,
		EndDate:   endDate,
	}

	err := r.db.QueryRowContext(ctx, query, startDate, endDate).Scan(
		&rep.TotalBooked,
		&rep.TotalActive,
		&rep.TotalReturned,
		&rep.TotalCompleted,
		&rep.TotalCancelled,
		&rep.TotalTransaction,
	)
	if err != nil {
		return nil, err
	}

	return rep, nil
}

func (r *reportRepo) GetTopProductsReport(ctx context.Context, startDate, endDate string, limit int) ([]models.TopProductReport, error) {
	query := `
		SELECT 
			p.id,
			p.name,
			c.name AS category_name,
			b.name AS brand_name,
			COUNT(ri.id) AS total_rented,
			COALESCE(SUM(ri.subtotal), 0) AS total_revenue
		FROM rental_items ri
		JOIN rentals ren ON ren.id = ri.rental_id
		JOIN product_units pu ON pu.id = ri.product_unit_id
		JOIN products p ON p.id = pu.product_id
		JOIN categories c ON c.id = p.category_id
		JOIN brands b ON b.id = p.brand_id
		WHERE ren.status IN ('active', 'returned', 'completed')
		  AND ren.deleted_at IS NULL
		  AND ri.deleted_at IS NULL
		  AND DATE(ren.booking_date) BETWEEN $1 AND $2
		GROUP BY p.id, p.name, c.name, b.name
		ORDER BY total_rented DESC, total_revenue DESC
		LIMIT $3`

	rows, err := r.db.QueryContext(ctx, query, startDate, endDate, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.TopProductReport
	for rows.Next() {
		var item models.TopProductReport
		if err := rows.Scan(&item.ProductID, &item.ProductName, &item.CategoryName, &item.BrandName, &item.TotalRented, &item.TotalRevenue); err != nil {
			return nil, err
		}
		list = append(list, item)
	}

	return list, nil
}

func (r *reportRepo) GetInventoryReport(ctx context.Context) (*models.InventoryReport, error) {
	query := `
		SELECT 
			COALESCE(SUM(p.lost_compensation_fee), 0) AS total_assets_value,
			COUNT(pu.id) AS total_units,
			COUNT(pu.id) FILTER (WHERE pu.status = 'available') AS available_units,
			COUNT(pu.id) FILTER (WHERE pu.status = 'booked') AS booked_units,
			COUNT(pu.id) FILTER (WHERE pu.status = 'rented') AS rented_units,
			COUNT(pu.id) FILTER (WHERE pu.status = 'maintenance') AS maintenance_units,
			COUNT(pu.id) FILTER (WHERE pu.status = 'lost') AS lost_units
		FROM product_units pu
		JOIN products p ON p.id = pu.product_id
		WHERE pu.deleted_at IS NULL AND p.deleted_at IS NULL`

	rep := &models.InventoryReport{}
	err := r.db.QueryRowContext(ctx, query).Scan(
		&rep.TotalAssetsValue,
		&rep.TotalUnits,
		&rep.AvailableUnits,
		&rep.BookedUnits,
		&rep.RentedUnits,
		&rep.MaintenanceUnits,
		&rep.LostUnits,
	)
	if err != nil {
		return nil, err
	}

	return rep, nil
}
