package models

type User struct {
	ID           int     `json:"id"`
	RoleID       int     `json:"role_id"`
	RoleName     string  `json:"role_name,omitempty"`
	Name         string  `json:"name"`
	Username     string  `json:"username"`
	PasswordHash string  `json:"-"`
	Phone        *string `json:"phone,omitempty"`
	IsActive     bool    `json:"is_active"`
	AuditTrail
}

type UserCreateDTO struct {
	RoleID   int     `json:"role_id" binding:"required"`
	Name     string  `json:"name" binding:"required"`
	Username string  `json:"username" binding:"required"`
	Password string  `json:"password" binding:"required,min=6"`
	Phone    *string `json:"phone"`
}

type UserUpdateDTO struct {
	RoleID   *int    `json:"role_id"`
	Name     *string `json:"name"`
	Username *string `json:"username"`
	Phone    *string `json:"phone"`
	IsActive *bool   `json:"is_active"`
	Password *string `json:"password"`
}
