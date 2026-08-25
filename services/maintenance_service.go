package services

import (
	"app/models"
	"app/repositories"
	"context"
	"database/sql"
	"errors"
	"time"
)

type MaintenanceService interface {
	GetAll(ctx context.Context) ([]models.Maintenance, error)
	GetByID(ctx context.Context, id int) (*models.Maintenance, error)
	GetDeleted(ctx context.Context) ([]models.Maintenance, error)
	Create(ctx context.Context, req models.MaintenanceCreateDTO, userID int) (*models.Maintenance, error)
	Update(ctx context.Context, id int, req models.MaintenanceUpdateDTO, userID int) error
	Complete(ctx context.Context, id int, userID int) error
	Delete(ctx context.Context, id, userID int) error
	Restore(ctx context.Context, id, userID int) error
	ForceDelete(ctx context.Context, id int) error
}

type maintenanceService struct {
	db              *sql.DB
	repo            repositories.MaintenanceRepository
	productUnitRepo repositories.ProductUnitRepository
}

func NewMaintenanceService(db *sql.DB, repo repositories.MaintenanceRepository, puRepo repositories.ProductUnitRepository) MaintenanceService {
	return &maintenanceService{db: db, repo: repo, productUnitRepo: puRepo}
}

func (s *maintenanceService) GetAll(ctx context.Context) ([]models.Maintenance, error) {
	return s.repo.GetAll(ctx)
}

func (s *maintenanceService) GetByID(ctx context.Context, id int) (*models.Maintenance, error) {
	m, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("data maintenance tidak ditemukan")
		}
		return nil, err
	}
	return m, nil
}

func (s *maintenanceService) GetDeleted(ctx context.Context) ([]models.Maintenance, error) {
	return s.repo.GetDeleted(ctx)
}

func (s *maintenanceService) Create(ctx context.Context, req models.MaintenanceCreateDTO, userID int) (*models.Maintenance, error) {
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		startDate = time.Now()
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	m := &models.Maintenance{
		ProductUnitID:    req.ProductUnitID,
		IssueDescription: req.IssueDescription,
		Cost:             req.Cost,
		StartDate:        startDate,
		Status:           "in_progress",
		AuditTrail: models.AuditTrail{
			CreatedBy:  &userID,
			ModifiedBy: &userID,
		},
	}

	if err := s.repo.CreateWithTx(ctx, tx, m); err != nil {
		return nil, err
	}

	if err := s.productUnitRepo.UpdateStatusWithTx(ctx, tx, req.ProductUnitID, "maintenance", userID); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return m, nil
}

func (s *maintenanceService) Update(ctx context.Context, id int, req models.MaintenanceUpdateDTO, userID int) error {
	err := s.repo.Update(ctx, id, req, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("data maintenance tidak ditemukan")
		}
		return err
	}
	return nil
}

func (s *maintenanceService) Complete(ctx context.Context, id int, userID int) error {
	m, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errors.New("data maintenance tidak ditemukan")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := s.repo.Complete(ctx, id, userID); err != nil {
		return err
	}

	if err := s.productUnitRepo.UpdateStatusWithTx(ctx, tx, m.ProductUnitID, "available", userID); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *maintenanceService) Delete(ctx context.Context, id, userID int) error {
	err := s.repo.Delete(ctx, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("data maintenance tidak ditemukan atau sudah dihapus")
		}
		return err
	}
	return nil
}

func (s *maintenanceService) Restore(ctx context.Context, id, userID int) error {
	err := s.repo.Restore(ctx, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("data maintenance tidak ditemukan atau tidak dalam status terhapus")
		}
		return err
	}
	return nil
}

func (s *maintenanceService) ForceDelete(ctx context.Context, id int) error {
	err := s.repo.ForceDelete(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("data maintenance tidak ditemukan")
		}
		return err
	}
	return nil
}
