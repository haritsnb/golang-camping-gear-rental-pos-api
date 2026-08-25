package services

import (
	"app/models"
	"app/repositories"
	"context"
	"database/sql"
	"errors"
)

type RoleService interface {
	GetAll(ctx context.Context) ([]models.Role, error)
	GetByID(ctx context.Context, id int) (*models.Role, error)
	GetDeleted(ctx context.Context) ([]models.Role, error)
	Create(ctx context.Context, req models.RoleCreateDTO, currentUserID int) (*models.Role, error)
	Update(ctx context.Context, id int, req models.RoleUpdateDTO, currentUserID int) error
	Delete(ctx context.Context, id, currentUserID int) error
	Restore(ctx context.Context, id, currentUserID int) error
	ForceDelete(ctx context.Context, id int) error
}

type roleService struct {
	repo repositories.RoleRepository
}

func NewRoleService(repo repositories.RoleRepository) RoleService {
	return &roleService{repo: repo}
}

func (s *roleService) GetAll(ctx context.Context) ([]models.Role, error) {
	return s.repo.GetAll(ctx)
}

func (s *roleService) GetByID(ctx context.Context, id int) (*models.Role, error) {
	role, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("role tidak ditemukan")
		}
		return nil, err
	}
	return role, nil
}

func (s *roleService) GetDeleted(ctx context.Context) ([]models.Role, error) {
	return s.repo.GetDeleted(ctx)
}

func (s *roleService) Create(ctx context.Context, req models.RoleCreateDTO, currentUserID int) (*models.Role, error) {
	existing, _ := s.repo.GetByName(ctx, req.Name)
	if existing != nil {
		return nil, errors.New("nama role sudah digunakan")
	}

	role := &models.Role{
		Name:        req.Name,
		Description: req.Description,
		AuditTrail: models.AuditTrail{
			CreatedBy:  &currentUserID,
			ModifiedBy: &currentUserID,
		},
	}

	if err := s.repo.Create(ctx, role); err != nil {
		return nil, err
	}

	return role, nil
}

func (s *roleService) Update(ctx context.Context, id int, req models.RoleUpdateDTO, currentUserID int) error {
	// Validasi jika name dikirimkan, tidak boleh duplikat dengan role lain
	if req.Name != nil && *req.Name != "" {
		existing, _ := s.repo.GetByName(ctx, *req.Name)
		if existing != nil && existing.ID != id {
			return errors.New("nama role sudah digunakan")
		}
	}

	err := s.repo.Update(ctx, id, req, currentUserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("role tidak ditemukan")
		}
		return err
	}

	return nil
}

func (s *roleService) Delete(ctx context.Context, id, currentUserID int) error {
	if id == 1 {
		return errors.New("role default 'admin' sistem tidak dapat dihapus")
	}

	err := s.repo.Delete(ctx, id, currentUserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("role tidak ditemukan atau sudah dihapus")
		}
		return err
	}

	return nil
}

func (s *roleService) Restore(ctx context.Context, id, currentUserID int) error {
	err := s.repo.Restore(ctx, id, currentUserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("role tidak ditemukan atau tidak dalam status terhapus")
		}
		return err
	}

	return nil
}

func (s *roleService) ForceDelete(ctx context.Context, id int) error {
	if id == 1 {
		return errors.New("role default 'admin' sistem tidak dapat dihapus permanen")
	}

	err := s.repo.ForceDelete(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("role tidak ditemukan")
		}

		return errors.New("gagal menghapus permanen: role masih digunakan oleh akun user yang terdaftar")
	}

	return nil
}
