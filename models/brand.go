package models

type Brand struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	IsActive    bool    `json:"is_active"`
	AuditTrail
}

type BrandCreateDTO struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
	IsActive    *bool   `json:"is_active"`
}

type BrandUpdateDTO struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	IsActive    *bool   `json:"is_active"`
}
