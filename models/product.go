package models

type Product struct {
	ID                  int     `json:"id"`
	CategoryID          int     `json:"category_id"`
	CategoryName        string  `json:"category_name,omitempty"`
	BrandID             int     `json:"brand_id"`
	BrandName           string  `json:"brand_name,omitempty"`
	Name                string  `json:"name"`
	RentalPricePerDay   float64 `json:"rental_price_per_day"`
	DefaultDeposit      float64 `json:"default_deposit"`
	LateFeePerDay       float64 `json:"late_fee_per_day"`
	LostCompensationFee float64 `json:"lost_compensation_fee"`
	IsActive            bool    `json:"is_active"`
	AuditTrail
}

type ProductCreateDTO struct {
	CategoryID          int     `json:"category_id" binding:"required"`
	BrandID             int     `json:"brand_id" binding:"required"`
	Name                string  `json:"name" binding:"required"`
	RentalPricePerDay   float64 `json:"rental_price_per_day" binding:"required"`
	DefaultDeposit      float64 `json:"default_deposit"`
	LateFeePerDay       float64 `json:"late_fee_per_day" binding:"required"`
	LostCompensationFee float64 `json:"lost_compensation_fee" binding:"required"`
	IsActive            *bool   `json:"is_active"`
}

type ProductUpdateDTO struct {
	CategoryID          *int     `json:"category_id"`
	BrandID             *int     `json:"brand_id"`
	Name                *string  `json:"name"`
	RentalPricePerDay   *float64 `json:"rental_price_per_day"`
	DefaultDeposit      *float64 `json:"default_deposit"`
	LateFeePerDay       *float64 `json:"late_fee_per_day"`
	LostCompensationFee *float64 `json:"lost_compensation_fee"`
	IsActive            *bool    `json:"is_active"`
}
