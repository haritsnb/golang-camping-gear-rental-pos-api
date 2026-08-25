package services

import (
	"app/models"
	"app/repositories"
	"context"
	"database/sql"
	"errors"
)

type CustomerService interface {
	GetAll(ctx context.Context) ([]models.Customer, error)
	GetByID(ctx context.Context, id int) (*models.Customer, error)
	GetDeleted(ctx context.Context) ([]models.Customer, error)
	Create(ctx context.Context, dto models.CustomerCreateDTO, photoURL *string, userID int) (*models.Customer, error)
	Update(ctx context.Context, id int, dto models.CustomerUpdateDTO, photoURL *string, userID int) error
	SetBlacklist(ctx context.Context, id int, isBlacklist bool, userID int) error
	Delete(ctx context.Context, id, userID int) error
	Restore(ctx context.Context, id, userID int) error
	ForceDelete(ctx context.Context, id int) error
}

type customerService struct {
	repo repositories.CustomerRepository
}

func NewCustomerService(repo repositories.CustomerRepository) CustomerService {
	return &customerService{repo: repo}
}

func (s *customerService) GetAll(ctx context.Context) ([]models.Customer, error) {
	return s.repo.GetAll(ctx)
}

func (s *customerService) GetByID(ctx context.Context, id int) (*models.Customer, error) {
	customer, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("pelanggan tidak ditemukan")
		}
		return nil, err
	}
	return customer, nil
}

func (s *customerService) GetDeleted(ctx context.Context) ([]models.Customer, error) {
	return s.repo.GetDeleted(ctx)
}

func (s *customerService) Create(ctx context.Context, dto models.CustomerCreateDTO, photoURL *string, userID int) (*models.Customer, error) {
	existing, err := s.repo.GetByIdentityNumber(ctx, dto.IdentityNumber)
	if err == nil && existing != nil {
		return nil, errors.New("nomor identitas (KTP/SIM/Paspor) sudah terdaftar")
	}

	customer := &models.Customer{
		Name:             dto.Name,
		IdentityType:     dto.IdentityType,
		IdentityNumber:   dto.IdentityNumber,
		IdentityPhotoURL: photoURL,
		Phone:            dto.Phone,
		EmergencyContact: dto.EmergencyContact,
		Email:            dto.Email,
		Address:          dto.Address,
		Notes:            dto.Notes,
		AuditTrail: models.AuditTrail{
			CreatedBy:  &userID,
			ModifiedBy: &userID,
		},
	}

	if err := s.repo.Create(ctx, customer); err != nil {
		return nil, err
	}

	return customer, nil
}

func (s *customerService) Update(ctx context.Context, id int, dto models.CustomerUpdateDTO, photoURL *string, userID int) error {
	// Validasi nomor identitas jika diubah
	if dto.IdentityNumber != nil && *dto.IdentityNumber != "" {
		existing, _ := s.repo.GetByIdentityNumber(ctx, *dto.IdentityNumber)
		if existing != nil && existing.ID != id {
			return errors.New("nomor identitas baru sudah digunakan oleh pelanggan lain")
		}
	}

	err := s.repo.Update(ctx, id, dto, photoURL, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("pelanggan tidak ditemukan")
		}
		return err
	}

	return nil
}

func (s *customerService) SetBlacklist(ctx context.Context, id int, isBlacklist bool, userID int) error {
	err := s.repo.SetBlacklist(ctx, id, isBlacklist, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("pelanggan tidak ditemukan")
		}
		return err
	}
	return nil
}

func (s *customerService) Delete(ctx context.Context, id, userID int) error {
	err := s.repo.Delete(ctx, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("pelanggan tidak ditemukan atau sudah dihapus")
		}
		return err
	}
	return nil
}

func (s *customerService) Restore(ctx context.Context, id, userID int) error {
	err := s.repo.Restore(ctx, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("pelanggan tidak ditemukan atau tidak dalam status terhapus")
		}
		return err
	}
	return nil
}

func (s *customerService) ForceDelete(ctx context.Context, id int) error {
	err := s.repo.ForceDelete(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("pelanggan tidak ditemukan")
		}
		return errors.New("gagal menghapus permanen: pelanggan masih memiliki riwayat transaksi rental")
	}
	return nil
}
