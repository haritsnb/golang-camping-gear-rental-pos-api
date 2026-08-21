package models

type Category struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	AuditTrail
}

type CategoryCreateDTO struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
}

type CategoryUpdateDTO struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}
