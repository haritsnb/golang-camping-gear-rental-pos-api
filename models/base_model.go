package models

import "time"

type AuditTrail struct {
	CreatedAt  time.Time  `json:"created_at"`
	CreatedBy  *int       `json:"created_by,omitempty"`
	ModifiedAt time.Time  `json:"modified_at"`
	ModifiedBy *int       `json:"modified_by,omitempty"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
	DeletedBy  *int       `json:"deleted_by,omitempty"`
}
