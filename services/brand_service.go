package services

import (
	"app/models"
	"app/repositories"
	"context"
	"database/sql"
	"errors"
)

type BrandService interface {
	GetAll(ctx context.Context) ([]models.Brand, error)
	GetByID(ctx context.Context, id int) (*models.Brand, error)
	GetDeleted(ctx context.Context) ([]models.Brand, error)
	Create(ctx context.Context, req models.BrandCreateDTO, userID int) (*models.Brand, error)
	Update(ctx context.Context, id int, req models.BrandUpdateDTO, userID int) error
	Delete(ctx context.Context, id, userID int) error
	Restore(ctx context.Context, id, userID int) error
	ForceDelete(ctx context.Context, id int) error
}

type brandService struct {
	repo repositories.BrandRepository
}

func NewBrandService(repo repositories.BrandRepository) BrandService {
	return &brandService{repo: repo}
}

func (s *brandService) GetAll(ctx context.Context) ([]models.Brand, error) {
	return s.repo.GetAll(ctx)
}

func (s *brandService) GetByID(ctx context.Context, id int) (*models.Brand, error) {
	b, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("brand tidak ditemukan")
		}
		return nil, err
	}
	return b, nil
}

func (s *brandService) GetDeleted(ctx context.Context) ([]models.Brand, error) {
	return s.repo.GetDeleted(ctx)
}

func (s *brandService) Create(ctx context.Context, req models.BrandCreateDTO, userID int) (*models.Brand, error) {
	existing, err := s.repo.GetByName(ctx, req.Name)
	if err == nil && existing != nil {
		return nil, errors.New("nama brand sudah digunakan")
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	b := &models.Brand{
		Name:        req.Name,
		Description: req.Description,
		IsActive:    isActive,
		AuditTrail: models.AuditTrail{
			CreatedBy:  &userID,
			ModifiedBy: &userID,
		},
	}
	err = s.repo.Create(ctx, b)
	return b, err
}

func (s *brandService) Update(ctx context.Context, id int, req models.BrandUpdateDTO, userID int) error {
	if req.Name != nil && *req.Name != "" {
		existing, _ := s.repo.GetByName(ctx, *req.Name)
		if existing != nil && existing.ID != id {
			return errors.New("nama brand sudah digunakan")
		}
	}

	err := s.repo.Update(ctx, id, req, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("brand tidak ditemukan")
		}
		return err
	}
	return nil
}

func (s *brandService) Delete(ctx context.Context, id, userID int) error {
	err := s.repo.Delete(ctx, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("brand tidak ditemukan atau sudah dihapus")
		}
		return err
	}
	return nil
}

func (s *brandService) Restore(ctx context.Context, id, userID int) error {
	err := s.repo.Restore(ctx, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("brand tidak ditemukan atau tidak dalam status terhapus")
		}
		return err
	}
	return nil
}

func (s *brandService) ForceDelete(ctx context.Context, id int) error {
	err := s.repo.ForceDelete(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("brand tidak ditemukan")
		}
		return errors.New("gagal menghapus permanen: brand masih digunakan oleh data produk yang ada")
	}
	return nil
}
