package services

import (
	"app/models"
	"app/repositories"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"time"
)

type RentalService interface {
	GetAll(ctx context.Context, status string) ([]models.Rental, error)
	GetByID(ctx context.Context, id int) (*models.Rental, error)
	GetDeleted(ctx context.Context) ([]models.Rental, error)
	Booking(ctx context.Context, req models.RentalBookingDTO, userID int) (*models.Rental, error)
	Handover(ctx context.Context, rentalID int, userID int) error
	Return(ctx context.Context, rentalID int, req models.RentalReturnDTO, userID int) (*models.Rental, error)
	Settlement(ctx context.Context, rentalID int, req models.RentalSettlementDTO, userID int) error
	Cancel(ctx context.Context, rentalID int, userID int) error
	Update(ctx context.Context, id int, req models.RentalUpdateDTO, userID int) error
	Delete(ctx context.Context, id, userID int) error
	Restore(ctx context.Context, id, userID int) error
	ForceDelete(ctx context.Context, id int) error
}

type rentalService struct {
	db              *sql.DB
	rentalRepo      repositories.RentalRepository
	rentalItemRepo  repositories.RentalItemRepository
	productUnitRepo repositories.ProductUnitRepository
	productRepo     repositories.ProductRepository
	customerRepo    repositories.CustomerRepository
	paymentRepo     repositories.PaymentRepository
	maintenanceRepo repositories.MaintenanceRepository
}

func NewRentalService(
	db *sql.DB,
	rRepo repositories.RentalRepository,
	riRepo repositories.RentalItemRepository,
	puRepo repositories.ProductUnitRepository,
	pRepo repositories.ProductRepository,
	cRepo repositories.CustomerRepository,
	payRepo repositories.PaymentRepository,
	mRepo repositories.MaintenanceRepository,
) RentalService {
	return &rentalService{
		db:              db,
		rentalRepo:      rRepo,
		rentalItemRepo:  riRepo,
		productUnitRepo: puRepo,
		productRepo:     pRepo,
		customerRepo:    cRepo,
		paymentRepo:     payRepo,
		maintenanceRepo: mRepo,
	}
}

func (s *rentalService) GetAll(ctx context.Context, status string) ([]models.Rental, error) {
	return s.rentalRepo.GetAll(ctx, status)
}

func (s *rentalService) GetByID(ctx context.Context, id int) (*models.Rental, error) {
	rental, err := s.rentalRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("transaksi rental tidak ditemukan")
		}
		return nil, err
	}
	items, err := s.rentalItemRepo.GetByRentalID(ctx, id)
	if err != nil {
		return nil, err
	}
	rental.Items = items
	return rental, nil
}

func (s *rentalService) GetDeleted(ctx context.Context) ([]models.Rental, error) {
	return s.rentalRepo.GetDeleted(ctx)
}

func (s *rentalService) Booking(ctx context.Context, req models.RentalBookingDTO, userID int) (*models.Rental, error) {
	cust, err := s.customerRepo.GetByID(ctx, req.CustomerID)
	if err != nil {
		return nil, errors.New("data customer tidak ditemukan")
	}
	if cust.IsBlacklisted {
		return nil, errors.New("customer masuk dalam daftar BLACKLIST, transaksi ditolak")
	}

	startDate, err := time.Parse("2006-01-02 15:04:05", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("format start_date salah (YYYY-MM-DD HH:MM:SS): %v", err)
	}
	expDate, err := time.Parse("2006-01-02 15:04:05", req.ExpectedReturnDate)
	if err != nil {
		return nil, fmt.Errorf("format expected_return_date salah (YYYY-MM-DD HH:MM:SS): %v", err)
	}

	diffHours := expDate.Sub(startDate).Hours()
	if diffHours <= 0 {
		return nil, errors.New("tanggal kembali harus lebih besar dari tanggal mulai")
	}
	durationDays := int(math.Ceil(diffHours / 24.0))
	if durationDays <= 0 {
		durationDays = 1
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var totalRentalFee, totalDeposit float64
	type itemData struct {
		unitID    int
		dailyRate float64
		subtotal  float64
	}
	var preparedItems []itemData

	for _, item := range req.Items {
		unit, err := s.productUnitRepo.GetByIDWithTx(ctx, tx, item.ProductUnitID)
		if err != nil {
			return nil, fmt.Errorf("unit ID %d tidak ditemukan", item.ProductUnitID)
		}
		if unit.Status != "available" {
			return nil, fmt.Errorf("unit %s saat ini sedang tidak tersedia (status: %s)", unit.UnitCode, unit.Status)
		}

		prod, err := s.productRepo.GetByID(ctx, unit.ProductID)
		if err != nil {
			return nil, err
		}

		subtotal := prod.RentalPricePerDay * float64(durationDays)
		totalRentalFee += subtotal
		totalDeposit += prod.DefaultDeposit

		preparedItems = append(preparedItems, itemData{
			unitID:    unit.ID,
			dailyRate: prod.RentalPricePerDay,
			subtotal:  subtotal,
		})

		if err := s.productUnitRepo.UpdateStatusWithTx(ctx, tx, unit.ID, "booked", userID); err != nil {
			return nil, err
		}
	}

	invoiceNo := fmt.Sprintf("INV-%d-%d", time.Now().Unix(), req.CustomerID)
	rental := &models.Rental{
		InvoiceNumber:      invoiceNo,
		CustomerID:         req.CustomerID,
		UserID:             userID,
		BookingDate:        time.Now(),
		StartDate:          startDate,
		ExpectedReturnDate: expDate,
		TotalRentalFee:     totalRentalFee,
		TotalDeposit:       totalDeposit,
		TotalPenaltyFee:    0,
		GrandTotal:         totalRentalFee + totalDeposit,
		Status:             "booked",
		AuditTrail: models.AuditTrail{
			CreatedBy:  &userID,
			ModifiedBy: &userID,
		},
	}

	if err := s.rentalRepo.CreateWithTx(ctx, tx, rental); err != nil {
		return nil, err
	}

	for _, pi := range preparedItems {
		rItem := &models.RentalItem{
			RentalID:      rental.ID,
			ProductUnitID: pi.unitID,
			DailyRate:     pi.dailyRate,
			DurationDays:  durationDays,
			Subtotal:      pi.subtotal,
			ConditionOut:  "good",
			AuditTrail: models.AuditTrail{
				CreatedBy:  &userID,
				ModifiedBy: &userID,
			},
		}
		if err := s.rentalItemRepo.CreateWithTx(ctx, tx, rItem); err != nil {
			return nil, err
		}
	}

	if req.DownPayment > 0 {
		payment := &models.Payment{
			RentalID:      rental.ID,
			UserID:        userID,
			PaymentType:   "down_payment",
			PaymentMethod: req.PaymentMethod,
			Amount:        req.DownPayment,
			AuditTrail: models.AuditTrail{
				CreatedBy:  &userID,
				ModifiedBy: &userID,
			},
		}
		if err := s.paymentRepo.CreateWithTx(ctx, tx, payment); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return rental, nil
}

func (s *rentalService) Handover(ctx context.Context, rentalID int, userID int) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	rental, err := s.rentalRepo.GetByIDWithTx(ctx, tx, rentalID)
	if err != nil {
		return errors.New("data rental tidak ditemukan")
	}
	if rental.Status != "booked" {
		return fmt.Errorf("handover gagal: status rental saat ini '%s' (harus 'booked')", rental.Status)
	}

	items, err := s.rentalItemRepo.GetByRentalIDWithTx(ctx, tx, rentalID)
	if err != nil {
		return err
	}

	for _, item := range items {
		if err := s.productUnitRepo.UpdateStatusWithTx(ctx, tx, item.ProductUnitID, "rented", userID); err != nil {
			return err
		}
	}

	rental.Status = "active"
	rental.ModifiedBy = &userID
	if err := s.rentalRepo.UpdateWithTx(ctx, tx, rental); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *rentalService) Return(ctx context.Context, rentalID int, req models.RentalReturnDTO, userID int) (*models.Rental, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	rental, err := s.rentalRepo.GetByIDWithTx(ctx, tx, rentalID)
	if err != nil {
		return nil, errors.New("data rental tidak ditemukan")
	}
	if rental.Status != "active" {
		return nil, fmt.Errorf("return gagal: status rental '%s' (harus 'active')", rental.Status)
	}

	actualReturn, err := time.Parse("2006-01-02 15:04:05", req.ActualReturnDate)
	if err != nil {
		return nil, errors.New("format actual_return_date salah")
	}

	var totalPenaltyFee float64

	if actualReturn.After(rental.ExpectedReturnDate) {
		lateHours := actualReturn.Sub(rental.ExpectedReturnDate).Hours()
		lateDays := int(math.Ceil(lateHours / 24.0))

		items, _ := s.rentalItemRepo.GetByRentalIDWithTx(ctx, tx, rentalID)
		for _, it := range items {
			unit, _ := s.productUnitRepo.GetByIDWithTx(ctx, tx, it.ProductUnitID)
			prod, _ := s.productRepo.GetByID(ctx, unit.ProductID)
			totalPenaltyFee += prod.LateFeePerDay * float64(lateDays)
		}
	}

	for _, itReq := range req.Items {
		itemPenalty := itReq.ItemPenaltyFee
		totalPenaltyFee += itemPenalty

		rItem := &models.RentalItem{
			ID:             itReq.RentalItemID,
			ConditionIn:    &itReq.ConditionIn,
			ItemPenaltyFee: itemPenalty,
			ReturnNotes:    itReq.ReturnNotes,
			AuditTrail: models.AuditTrail{
				ModifiedBy: &userID,
			},
		}
		if err := s.rentalItemRepo.UpdateReturnWithTx(ctx, tx, rItem); err != nil {
			return nil, err
		}
	}

	rental.ActualReturnDate = &actualReturn
	rental.TotalPenaltyFee = totalPenaltyFee
	rental.GrandTotal = rental.TotalRentalFee + rental.TotalPenaltyFee
	rental.Status = "returned"
	rental.ModifiedBy = &userID

	if err := s.rentalRepo.UpdateWithTx(ctx, tx, rental); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return rental, nil
}

func (s *rentalService) Settlement(ctx context.Context, rentalID int, req models.RentalSettlementDTO, userID int) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	rental, err := s.rentalRepo.GetByIDWithTx(ctx, tx, rentalID)
	if err != nil {
		return errors.New("data rental tidak ditemukan")
	}
	if rental.Status != "returned" {
		return fmt.Errorf("settlement gagal: rental belum dalam status 'returned'")
	}

	netDepositBalance := rental.TotalDeposit - rental.TotalPenaltyFee

	if netDepositBalance > 0 {
		pay := &models.Payment{
			RentalID:      rental.ID,
			UserID:        userID,
			PaymentType:   "deposit_refund",
			PaymentMethod: req.PaymentMethod,
			Amount:        netDepositBalance,
			Notes:         req.Notes,
			AuditTrail: models.AuditTrail{
				CreatedBy:  &userID,
				ModifiedBy: &userID,
			},
		}
		if err := s.paymentRepo.CreateWithTx(ctx, tx, pay); err != nil {
			return err
		}
	} else if netDepositBalance < 0 {
		pay := &models.Payment{
			RentalID:      rental.ID,
			UserID:        userID,
			PaymentType:   "penalty_charge",
			PaymentMethod: req.PaymentMethod,
			Amount:        math.Abs(netDepositBalance),
			Notes:         req.Notes,
			AuditTrail: models.AuditTrail{
				CreatedBy:  &userID,
				ModifiedBy: &userID,
			},
		}
		if err := s.paymentRepo.CreateWithTx(ctx, tx, pay); err != nil {
			return err
		}
	}

	items, err := s.rentalItemRepo.GetByRentalIDWithTx(ctx, tx, rentalID)
	if err != nil {
		return err
	}

	for _, item := range items {
		newStatus := "available"
		if item.ConditionIn != nil {
			switch *item.ConditionIn {
			case "damaged":
				newStatus = "maintenance"
				maint := &models.Maintenance{
					ProductUnitID:    item.ProductUnitID,
					IssueDescription: fmt.Sprintf("Kerusakan pasca rental %s. Catatan: %v", rental.InvoiceNumber, item.ReturnNotes),
					StartDate:        time.Now(),
					Status:           "in_progress",
					AuditTrail: models.AuditTrail{
						CreatedBy:  &userID,
						ModifiedBy: &userID,
					},
				}
				_ = s.maintenanceRepo.CreateWithTx(ctx, tx, maint)
			case "lost":
				newStatus = "lost"
			}
		}
		if err := s.productUnitRepo.UpdateStatusWithTx(ctx, tx, item.ProductUnitID, newStatus, userID); err != nil {
			return err
		}
	}

	rental.Status = "completed"
	rental.ModifiedBy = &userID
	if err := s.rentalRepo.UpdateWithTx(ctx, tx, rental); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *rentalService) Cancel(ctx context.Context, rentalID int, userID int) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	rental, err := s.rentalRepo.GetByIDWithTx(ctx, tx, rentalID)
	if err != nil {
		return errors.New("rental tidak ditemukan")
	}
	if rental.Status != "booked" {
		return errors.New("hanya booking yang belum diambil yang dapat dibatalkan")
	}

	items, _ := s.rentalItemRepo.GetByRentalIDWithTx(ctx, tx, rentalID)
	for _, it := range items {
		_ = s.productUnitRepo.UpdateStatusWithTx(ctx, tx, it.ProductUnitID, "available", userID)
	}

	rental.Status = "cancelled"
	rental.ModifiedBy = &userID
	if err := s.rentalRepo.UpdateWithTx(ctx, tx, rental); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *rentalService) Update(ctx context.Context, id int, req models.RentalUpdateDTO, userID int) error {
	err := s.rentalRepo.UpdateFlexible(ctx, id, req, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("transaksi rental tidak ditemukan")
		}
		return err
	}
	return nil
}

func (s *rentalService) Delete(ctx context.Context, id, userID int) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := s.rentalRepo.DeleteWithTx(ctx, tx, id, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("transaksi rental tidak ditemukan atau sudah dihapus")
		}
		return err
	}
	if err := s.rentalItemRepo.DeleteByRentalIDWithTx(ctx, tx, id, userID); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *rentalService) Restore(ctx context.Context, id, userID int) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := s.rentalRepo.RestoreWithTx(ctx, tx, id, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("transaksi rental tidak ditemukan atau tidak dalam status terhapus")
		}
		return err
	}
	if err := s.rentalItemRepo.RestoreByRentalIDWithTx(ctx, tx, id, userID); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *rentalService) ForceDelete(ctx context.Context, id int) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Hapus riwayat items
	_ = s.rentalItemRepo.ForceDeleteByRentalIDWithTx(ctx, tx, id)

	// 2. Hapus rental permanen
	if err := s.rentalRepo.ForceDeleteWithTx(ctx, tx, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("transaksi rental tidak ditemukan")
		}
		return errors.New("gagal menghapus permanen: rental masih terikat dengan data pembayaran")
	}

	return tx.Commit()
}
