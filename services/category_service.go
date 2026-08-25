package services

import (
	"app/models"
	"app/repositories"
	"context"
	"database/sql"
	"errors"
)

type CategoryService interface {
	GetAll(ctx context.Context) ([]models.Category, error)
	GetByID(ctx context.Context, id int) (*models.Category, error)
	GetDeleted(ctx context.Context) ([]models.Category, error)
	Create(ctx context.Context, req models.CategoryCreateDTO, userID int) (*models.Category, error)
	Update(ctx context.Context, id int, req models.CategoryUpdateDTO, userID int) error
	Delete(ctx context.Context, id, userID int) error
	Restore(ctx context.Context, id, userID int) error
	ForceDelete(ctx context.Context, id int) error
}

type categoryService struct {
	repo repositories.CategoryRepository
}

func NewCategoryService(repo repositories.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) GetAll(ctx context.Context) ([]models.Category, error) {
	return s.repo.GetAll(ctx)
}

func (s *categoryService) GetByID(ctx context.Context, id int) (*models.Category, error) {
	cat, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("kategori tidak ditemukan")
		}
		return nil, err
	}
	return cat, nil
}

func (s *categoryService) GetDeleted(ctx context.Context) ([]models.Category, error) {
	return s.repo.GetDeleted(ctx)
}

func (s *categoryService) Create(ctx context.Context, req models.CategoryCreateDTO, userID int) (*models.Category, error) {
	existing, err := s.repo.GetByName(ctx, req.Name)
	if err == nil && existing != nil {
		return nil, errors.New("nama kategori sudah digunakan")
	}

	cat := &models.Category{
		Name:        req.Name,
		Description: req.Description,
		AuditTrail: models.AuditTrail{
			CreatedBy:  &userID,
			ModifiedBy: &userID,
		},
	}
	err = s.repo.Create(ctx, cat)
	return cat, err
}

func (s *categoryService) Update(ctx context.Context, id int, req models.CategoryUpdateDTO, userID int) error {
	if req.Name != nil && *req.Name != "" {
		existing, _ := s.repo.GetByName(ctx, *req.Name)
		if existing != nil && existing.ID != id {
			return errors.New("nama kategori sudah digunakan")
		}
	}

	err := s.repo.Update(ctx, id, req, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("kategori tidak ditemukan")
		}
		return err
	}
	return nil
}

func (s *categoryService) Delete(ctx context.Context, id, userID int) error {
	err := s.repo.Delete(ctx, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("kategori tidak ditemukan atau sudah dihapus")
		}
		return err
	}
	return nil
}

func (s *categoryService) Restore(ctx context.Context, id, userID int) error {
	err := s.repo.Restore(ctx, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("kategori tidak ditemukan atau tidak dalam status terhapus")
		}
		return err
	}
	return nil
}

func (s *categoryService) ForceDelete(ctx context.Context, id int) error {
	err := s.repo.ForceDelete(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("kategori tidak ditemukan")
		}
		return errors.New("gagal menghapus permanen: kategori masih digunakan oleh data produk yang ada")
	}
	return nil
}
