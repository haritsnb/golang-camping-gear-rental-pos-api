package services

import (
	"app/models"
	"app/repositories"
	"context"
	"time"
)

type ReportService interface {
	GetRevenueReport(ctx context.Context, startDate, endDate string) (*models.RevenueReport, error)
	GetRentalSummaryReport(ctx context.Context, startDate, endDate string) (*models.RentalSummaryReport, error)
	GetTopProductsReport(ctx context.Context, startDate, endDate string, limit int) ([]models.TopProductReport, error)
	GetInventoryReport(ctx context.Context) (*models.InventoryReport, error)
}

type reportService struct {
	repo repositories.ReportRepository
}

func NewReportService(repo repositories.ReportRepository) ReportService {
	return &reportService{repo: repo}
}

// Memastikan format tanggal default (Awal bulan s/d Hari ini jika kosong)
func parseDateRange(start, end string) (string, string) {
	now := time.Now()
	if start == "" {
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")
	}
	if end == "" {
		end = now.Format("2006-01-02")
	}
	return start, end
}

func (s *reportService) GetRevenueReport(ctx context.Context, startDate, endDate string) (*models.RevenueReport, error) {
	start, end := parseDateRange(startDate, endDate)
	return s.repo.GetRevenueReport(ctx, start, end)
}

func (s *reportService) GetRentalSummaryReport(ctx context.Context, startDate, endDate string) (*models.RentalSummaryReport, error) {
	start, end := parseDateRange(startDate, endDate)
	return s.repo.GetRentalSummaryReport(ctx, start, end)
}

func (s *reportService) GetTopProductsReport(ctx context.Context, startDate, endDate string, limit int) ([]models.TopProductReport, error) {
	start, end := parseDateRange(startDate, endDate)
	if limit <= 0 {
		limit = 10
	}
	return s.repo.GetTopProductsReport(ctx, start, end, limit)
}

func (s *reportService) GetInventoryReport(ctx context.Context) (*models.InventoryReport, error) {
	return s.repo.GetInventoryReport(ctx)
}
