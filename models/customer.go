package models

type Customer struct {
	ID               int     `json:"id"`
	Name             string  `json:"name"`
	IdentityType     string  `json:"identity_type"`
	IdentityNumber   string  `json:"identity_number"`
	IdentityPhotoURL *string `json:"identity_photo_url,omitempty"`
	Phone            string  `json:"phone"`
	EmergencyContact *string `json:"emergency_contact,omitempty"`
	Email            *string `json:"email,omitempty"`
	Address          string  `json:"address"`
	IsBlacklisted    bool    `json:"is_blacklisted"`
	Notes            *string `json:"notes,omitempty"`
	AuditTrail
}

type CustomerCreateDTO struct {
	Name             string  `form:"name" json:"name" binding:"required"`
	IdentityType     string  `form:"identity_type" json:"identity_type" binding:"required"`
	IdentityNumber   string  `form:"identity_number" json:"identity_number" binding:"required"`
	Phone            string  `form:"phone" json:"phone" binding:"required"`
	EmergencyContact *string `form:"emergency_contact" json:"emergency_contact"`
	Email            *string `form:"email" json:"email"`
	Address          string  `form:"address" json:"address" binding:"required"`
	Notes            *string `form:"notes" json:"notes"`
}

type CustomerUpdateDTO struct {
	Name             *string `form:"name" json:"name"`
	IdentityType     *string `form:"identity_type" json:"identity_type"`
	IdentityNumber   *string `form:"identity_number" json:"identity_number"`
	Phone            *string `form:"phone" json:"phone"`
	EmergencyContact *string `form:"emergency_contact" json:"emergency_contact"`
	Email            *string `form:"email" json:"email"`
	Address          *string `form:"address" json:"address"`
	IsBlacklisted    *bool   `form:"is_blacklisted" json:"is_blacklisted"`
	Notes            *string `form:"notes" json:"notes"`
}
