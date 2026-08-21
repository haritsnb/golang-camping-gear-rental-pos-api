package services

import (
	"app/models"
	"app/repositories"
	"context"
	"database/sql"
	"errors"
)

type ProductUnitService interface {
	GetAll(ctx context.Context) ([]models.ProductUnit, error)
	GetByID(ctx context.Context, id int) (*models.ProductUnit, error)
	GetByUnitCode(ctx context.Context, code string) (*models.ProductUnit, error)
	Create(ctx context.Context, req models.ProductUnitCreateDTO, userID int) (*models.ProductUnit, error)
	Update(ctx context.Context, id int, req models.ProductUnitUpdateDTO, userID int) error
	Delete(ctx context.Context, id, userID int) error
	Restore(ctx context.Context, id, userID int) error
	ForceDelete(ctx context.Context, id int) error
}

type productUnitService struct {
	repo repositories.ProductUnitRepository
}

func NewProductUnitService(repo repositories.ProductUnitRepository) ProductUnitService {
	return &productUnitService{repo: repo}
}

func (s *productUnitService) GetAll(ctx context.Context) ([]models.ProductUnit, error) {
	return s.repo.GetAll(ctx)
}

func (s *productUnitService) GetByID(ctx context.Context, id int) (*models.ProductUnit, error) {
	unit, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("unit fisik tidak ditemukan")
		}
		return nil, err
	}
	return unit, nil
}

func (s *productUnitService) GetByUnitCode(ctx context.Context, code string) (*models.ProductUnit, error) {
	unit, err := s.repo.GetByUnitCode(ctx, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("unit fisik dengan barcode/QR tersebut tidak ditemukan")
		}
		return nil, err
	}
	return unit, nil
}

func (s *productUnitService) Create(ctx context.Context, req models.ProductUnitCreateDTO, userID int) (*models.ProductUnit, error) {
	existing, _ := s.repo.GetByUnitCode(ctx, req.UnitCode)
	if existing != nil {
		return nil, errors.New("kode unit fisik (unit_code) sudah terdaftar")
	}

	cond := "good"
	if req.Condition != "" {
		cond = req.Condition
	}
	u := &models.ProductUnit{
		ProductID:    req.ProductID,
		UnitCode:     req.UnitCode,
		SerialNumber: req.SerialNumber,
		Condition:    cond,
		Status:       "available",
		Notes:        req.Notes,
		AuditTrail: models.AuditTrail{
			CreatedBy:  &userID,
			ModifiedBy: &userID,
		},
	}
	err := s.repo.Create(ctx, u)
	return u, err
}

func (s *productUnitService) Update(ctx context.Context, id int, req models.ProductUnitUpdateDTO, userID int) error {
	if req.UnitCode != nil && *req.UnitCode != "" {
		existing, _ := s.repo.GetByUnitCode(ctx, *req.UnitCode)
		if existing != nil && existing.ID != id {
			return errors.New("kode unit fisik baru sudah digunakan oleh unit lain")
		}
	}

	err := s.repo.Update(ctx, id, req, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("unit fisik tidak ditemukan")
		}
		return err
	}
	return nil
}

func (s *productUnitService) Delete(ctx context.Context, id, userID int) error {
	err := s.repo.Delete(ctx, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("unit fisik tidak ditemukan atau sudah dihapus")
		}
		return err
	}
	return nil
}

func (s *productUnitService) Restore(ctx context.Context, id, userID int) error {
	err := s.repo.Restore(ctx, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("unit fisik tidak ditemukan atau tidak dalam status terhapus")
		}
		return err
	}
	return nil
}

func (s *productUnitService) ForceDelete(ctx context.Context, id int) error {
	err := s.repo.ForceDelete(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("unit fisik tidak ditemukan")
		}
		return errors.New("gagal menghapus permanen: unit fisik masih terikat dengan riwayat transaksi rental atau servis")
	}
	return nil
}
