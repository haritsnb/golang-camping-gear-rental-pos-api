package models

import "time"

type Maintenance struct {
	ID               int        `json:"id"`
	ProductUnitID    int        `json:"product_unit_id"`
	UnitCode         string     `json:"unit_code,omitempty"`
	IssueDescription string     `json:"issue_description"`
	Cost             float64    `json:"cost"`
	StartDate        time.Time  `json:"start_date"`
	CompletedDate    *time.Time `json:"completed_date,omitempty"`
	Status           string     `json:"status"` // in_progress, completed
	AuditTrail
}

type MaintenanceCreateDTO struct {
	ProductUnitID    int     `json:"product_unit_id" binding:"required"`
	IssueDescription string  `json:"issue_description" binding:"required"`
	Cost             float64 `json:"cost"`
	StartDate        string  `json:"start_date" binding:"required"` // 2006-01-02
}

type MaintenanceUpdateDTO struct {
	IssueDescription *string  `json:"issue_description"`
	Cost             *float64 `json:"cost"`
	StartDate        *string  `json:"start_date"`
	Status           *string  `json:"status"`
}
