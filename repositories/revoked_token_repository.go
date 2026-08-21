package repositories

import (
	"app/models"
	"context"
	"database/sql"
)

type RevokedTokenRepository interface {
	Revoke(ctx context.Context, token *models.RevokedToken) error
	IsRevoked(ctx context.Context, jti string) (bool, error)
}

type revokedTokenRepo struct {
	db *sql.DB
}

func NewRevokedTokenRepository(db *sql.DB) RevokedTokenRepository {
	return &revokedTokenRepo{db: db}
}

func (r *revokedTokenRepo) Revoke(ctx context.Context, token *models.RevokedToken) error {
	query := `INSERT INTO revoked_tokens (token_jti, expires_at, revoked_at) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, token.TokenJTI, token.ExpiresAt, token.RevokedAt)
	return err
}

func (r *revokedTokenRepo) IsRevoked(ctx context.Context, jti string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM revoked_tokens WHERE token_jti = $1)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, jti).Scan(&exists)
	return exists, err
}
