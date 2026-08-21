package models

import "time"

type Rental struct {
	ID                 int          `json:"id"`
	InvoiceNumber      string       `json:"invoice_number"`
	CustomerID         int          `json:"customer_id"`
	CustomerName       string       `json:"customer_name,omitempty"`
	UserID             int          `json:"user_id"`
	UserName           string       `json:"user_name,omitempty"`
	BookingDate        time.Time    `json:"booking_date"`
	StartDate          time.Time    `json:"start_date"`
	ExpectedReturnDate time.Time    `json:"expected_return_date"`
	ActualReturnDate   *time.Time   `json:"actual_return_date,omitempty"`
	TotalRentalFee     float64      `json:"total_rental_fee"`
	TotalDeposit       float64      `json:"total_deposit"`
	TotalPenaltyFee    float64      `json:"total_penalty_fee"`
	GrandTotal         float64      `json:"grand_total"`
	Status             string       `json:"status"` // booked, active, returned, completed, cancelled
	Items              []RentalItem `json:"items,omitempty"`
	AuditTrail
}

type BookingItemDTO struct {
	ProductUnitID int `json:"product_unit_id" binding:"required"`
}

type RentalBookingDTO struct {
	CustomerID         int              `json:"customer_id" binding:"required"`
	StartDate          string           `json:"start_date" binding:"required"` // Format: 2006-01-02 15:04:05
	ExpectedReturnDate string           `json:"expected_return_date" binding:"required"`
	Items              []BookingItemDTO `json:"items" binding:"required,min=1"`
	DownPayment        float64          `json:"down_payment"`
	PaymentMethod      string           `json:"payment_method"`
}

type ReturnItemConditionDTO struct {
	RentalItemID   int     `json:"rental_item_id" binding:"required"`
	ConditionIn    string  `json:"condition_in" binding:"required"` // good, damaged, lost
	ItemPenaltyFee float64 `json:"item_penalty_fee"`
	ReturnNotes    *string `json:"return_notes"`
}

type RentalReturnDTO struct {
	ActualReturnDate string                   `json:"actual_return_date" binding:"required"`
	Items            []ReturnItemConditionDTO `json:"items" binding:"required,min=1"`
}

type RentalSettlementDTO struct {
	PaymentMethod string  `json:"payment_method" binding:"required"`
	Notes         *string `json:"notes"`
}

type RentalUpdateDTO struct {
	StartDate          *string  `json:"start_date"`
	ExpectedReturnDate *string  `json:"expected_return_date"`
	TotalRentalFee     *float64 `json:"total_rental_fee"`
	TotalDeposit       *float64 `json:"total_deposit"`
	TotalPenaltyFee    *float64 `json:"total_penalty_fee"`
	GrandTotal         *float64 `json:"grand_total"`
	Status             *string  `json:"status"`
}
