package models

import "time"

type RevokedToken struct {
	ID        int       `json:"id"`
	TokenJTI  string    `json:"token_jti"`
	ExpiresAt time.Time `json:"expires_at"`
	RevokedAt time.Time `json:"revoked_at"`
}
