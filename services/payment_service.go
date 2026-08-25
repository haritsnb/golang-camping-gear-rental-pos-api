package services

import (
	"app/models"
	"app/repositories"
	"context"
	"database/sql"
	"errors"
)

type PaymentService interface {
	GetByRentalID(ctx context.Context, rentalID int) ([]models.Payment, error)
	GetByID(ctx context.Context, id int) (*models.Payment, error)
	GetDeleted(ctx context.Context) ([]models.Payment, error)
	Create(ctx context.Context, req models.PaymentCreateDTO, userID int) (*models.Payment, error)
	Update(ctx context.Context, id int, req models.PaymentUpdateDTO, userID int) error
	Delete(ctx context.Context, id, userID int) error
	Restore(ctx context.Context, id, userID int) error
	ForceDelete(ctx context.Context, id int) error
}

type paymentService struct {
	repo repositories.PaymentRepository
}

func NewPaymentService(repo repositories.PaymentRepository) PaymentService {
	return &paymentService{repo: repo}
}

func (s *paymentService) GetByRentalID(ctx context.Context, rentalID int) ([]models.Payment, error) {
	return s.repo.GetByRentalID(ctx, rentalID)
}

func (s *paymentService) GetByID(ctx context.Context, id int) (*models.Payment, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("data pembayaran tidak ditemukan")
		}
		return nil, err
	}
	return p, nil
}

func (s *paymentService) GetDeleted(ctx context.Context) ([]models.Payment, error) {
	return s.repo.GetDeleted(ctx)
}

func (s *paymentService) Create(ctx context.Context, req models.PaymentCreateDTO, userID int) (*models.Payment, error) {
	p := &models.Payment{
		RentalID:        req.RentalID,
		UserID:          userID,
		PaymentType:     req.PaymentType,
		PaymentMethod:   req.PaymentMethod,
		Amount:          req.Amount,
		ReferenceNumber: req.ReferenceNumber,
		Notes:           req.Notes,
		AuditTrail: models.AuditTrail{
			CreatedBy:  &userID,
			ModifiedBy: &userID,
		},
	}
	err := s.repo.Create(ctx, p)
	return p, err
}

func (s *paymentService) Update(ctx context.Context, id int, req models.PaymentUpdateDTO, userID int) error {
	err := s.repo.Update(ctx, id, req, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("data pembayaran tidak ditemukan")
		}
		return err
	}
	return nil
}

func (s *paymentService) Delete(ctx context.Context, id, userID int) error {
	err := s.repo.Delete(ctx, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("data pembayaran tidak ditemukan atau sudah dihapus")
		}
		return err
	}
	return nil
}

func (s *paymentService) Restore(ctx context.Context, id, userID int) error {
	err := s.repo.Restore(ctx, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("data pembayaran tidak ditemukan atau tidak dalam status terhapus")
		}
		return err
	}
	return nil
}

func (s *paymentService) ForceDelete(ctx context.Context, id int) error {
	err := s.repo.ForceDelete(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("data pembayaran tidak ditemukan")
		}
		return err
	}
	return nil
}
