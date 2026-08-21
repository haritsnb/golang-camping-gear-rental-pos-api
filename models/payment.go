package models

import "time"

type Payment struct {
	ID              int       `json:"id"`
	RentalID        int       `json:"rental_id"`
	UserID          int       `json:"user_id"`
	PaymentType     string    `json:"payment_type"` // down_payment, full_payment, deposit_in, deposit_refund, penalty_charge
	PaymentMethod   string    `json:"payment_method"`
	Amount          float64   `json:"amount"`
	ReferenceNumber *string   `json:"reference_number,omitempty"`
	Notes           *string   `json:"notes,omitempty"`
	PaidAt          time.Time `json:"paid_at"`
	AuditTrail
}

type PaymentCreateDTO struct {
	RentalID        int     `json:"rental_id" binding:"required"`
	PaymentType     string  `json:"payment_type" binding:"required"`
	PaymentMethod   string  `json:"payment_method" binding:"required"`
	Amount          float64 `json:"amount" binding:"required"`
	ReferenceNumber *string `json:"reference_number"`
	Notes           *string `json:"notes"`
}

type PaymentUpdateDTO struct {
	PaymentType     *string  `json:"payment_type"`
	PaymentMethod   *string  `json:"payment_method"`
	Amount          *float64 `json:"amount"`
	ReferenceNumber *string  `json:"reference_number"`
	Notes           *string  `json:"notes"`
}
