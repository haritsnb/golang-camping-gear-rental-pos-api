package services

import (
	"app/models"
	"app/repositories"
	"context"
	"database/sql"
	"errors"
)

type ProductService interface {
	GetAll(ctx context.Context) ([]models.Product, error)
	GetByID(ctx context.Context, id int) (*models.Product, error)
	GetDeleted(ctx context.Context) ([]models.Product, error)
	Create(ctx context.Context, req models.ProductCreateDTO, userID int) (*models.Product, error)
	Update(ctx context.Context, id int, req models.ProductUpdateDTO, userID int) error
	Delete(ctx context.Context, id, userID int) error
	Restore(ctx context.Context, id, userID int) error
	ForceDelete(ctx context.Context, id int) error
}

type productService struct {
	repo repositories.ProductRepository
}

func NewProductService(repo repositories.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) GetAll(ctx context.Context) ([]models.Product, error) {
	return s.repo.GetAll(ctx)
}

func (s *productService) GetByID(ctx context.Context, id int) (*models.Product, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("produk tidak ditemukan")
		}
		return nil, err
	}
	return p, nil
}

func (s *productService) GetDeleted(ctx context.Context) ([]models.Product, error) {
	return s.repo.GetDeleted(ctx)
}

func (s *productService) Create(ctx context.Context, req models.ProductCreateDTO, userID int) (*models.Product, error) {
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	p := &models.Product{
		CategoryID:          req.CategoryID,
		BrandID:             req.BrandID,
		Name:                req.Name,
		RentalPricePerDay:   req.RentalPricePerDay,
		DefaultDeposit:      req.DefaultDeposit,
		LateFeePerDay:       req.LateFeePerDay,
		LostCompensationFee: req.LostCompensationFee,
		IsActive:            isActive,
		AuditTrail: models.AuditTrail{
			CreatedBy:  &userID,
			ModifiedBy: &userID,
		},
	}
	err := s.repo.Create(ctx, p)
	return p, err
}

func (s *productService) Update(ctx context.Context, id int, req models.ProductUpdateDTO, userID int) error {
	err := s.repo.Update(ctx, id, req, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("produk tidak ditemukan")
		}
		return err
	}
	return nil
}

func (s *productService) Delete(ctx context.Context, id, userID int) error {
	err := s.repo.Delete(ctx, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("produk tidak ditemukan atau sudah dihapus")
		}
		return err
	}
	return nil
}

func (s *productService) Restore(ctx context.Context, id, userID int) error {
	err := s.repo.Restore(ctx, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("produk tidak ditemukan atau tidak dalam status terhapus")
		}
		return err
	}
	return nil
}

func (s *productService) ForceDelete(ctx context.Context, id int) error {
	err := s.repo.ForceDelete(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("produk tidak ditemukan")
		}
		return errors.New("gagal menghapus permanen: produk masih memiliki unit fisik terkait")
	}
	return nil
}
