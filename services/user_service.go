package services

import (
	"app/models"
	"app/repositories"
	"context"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetAll(ctx context.Context) ([]models.User, error)
	GetByID(ctx context.Context, id int) (*models.User, error)
	Create(ctx context.Context, req models.UserCreateDTO, currentUserID int) (*models.User, error)
	Update(ctx context.Context, id int, req models.UserUpdateDTO, currentUserID int) error
	Delete(ctx context.Context, id, currentUserID int) error
	Restore(ctx context.Context, id, currentUserID int) error     // <-- Baru
	ForceDelete(ctx context.Context, id, currentUserID int) error // <-- Baru
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetAll(ctx context.Context) ([]models.User, error) {
	return s.repo.GetAll(ctx)
}

func (s *userService) GetByID(ctx context.Context, id int) (*models.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user tidak ditemukan")
		}
		return nil, err
	}
	return user, nil
}

func (s *userService) Create(ctx context.Context, req models.UserCreateDTO, currentUserID int) (*models.User, error) {
	existingUser, _ := s.repo.GetByUsername(ctx, req.Username)
	if existingUser != nil {
		return nil, errors.New("username sudah digunakan oleh akun lain")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("gagal mengenkripsi password")
	}

	user := &models.User{
		RoleID:       req.RoleID,
		Name:         req.Name,
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		Phone:        req.Phone,
		IsActive:     true,
		AuditTrail: models.AuditTrail{
			CreatedBy:  &currentUserID,
			ModifiedBy: &currentUserID,
		},
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Update(ctx context.Context, id int, req models.UserUpdateDTO, currentUserID int) error {
	if req.Username != nil && *req.Username != "" {
		existingUser, _ := s.repo.GetByUsername(ctx, *req.Username)
		if existingUser != nil && existingUser.ID != id {
			return errors.New("username baru sudah digunakan oleh akun lain")
		}
	}

	var passwordHash *string
	if req.Password != nil && *req.Password != "" {
		if len(*req.Password) < 6 {
			return errors.New("password minimal harus 6 karakter")
		}
		hashed, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return errors.New("gagal mengenkripsi password baru")
		}
		strHash := string(hashed)
		passwordHash = &strHash
	}

	err := s.repo.Update(ctx, id, req, passwordHash, currentUserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("user tidak ditemukan")
		}
		return err
	}

	return nil
}

func (s *userService) Delete(ctx context.Context, id, currentUserID int) error {
	if id == currentUserID {
		return errors.New("tidak dapat menghapus akun yang sedang digunakan saat ini")
	}

	err := s.repo.Delete(ctx, id, currentUserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("user tidak ditemukan atau sudah dihapus")
		}
		return err
	}

	return nil
}

func (s *userService) Restore(ctx context.Context, id, currentUserID int) error {
	err := s.repo.Restore(ctx, id, currentUserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("user tidak ditemukan atau user tidak dalam status terhapus")
		}
		return err
	}

	return nil
}

func (s *userService) ForceDelete(ctx context.Context, id, currentUserID int) error {
	if id == currentUserID {
		return errors.New("tidak dapat menghapus permanen akun yang sedang Anda gunakan")
	}

	err := s.repo.ForceDelete(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("user tidak ditemukan")
		}

		return errors.New("gagal menghapus permanen: user masih memiliki riwayat transaksi/data audit terkait")
	}

	return nil
}
