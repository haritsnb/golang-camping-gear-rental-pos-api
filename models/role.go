package models

type Role struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	AuditTrail
}

type RoleCreateDTO struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
}

type RoleUpdateDTO struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}
