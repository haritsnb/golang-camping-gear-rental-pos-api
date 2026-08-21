package services

import (
	"app/config"
	"app/models"
	"app/repositories"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx context.Context, username, password string) (string, *models.User, error)
	Logout(ctx context.Context, jti string, exp time.Time) error
	GetProfile(ctx context.Context, userID int) (*models.User, error)
}

type authService struct {
	userRepo    repositories.UserRepository
	revokedRepo repositories.RevokedTokenRepository
}

func NewAuthService(u repositories.UserRepository, r repositories.RevokedTokenRepository) AuthService {
	return &authService{userRepo: u, revokedRepo: r}
}

func (s *authService) Login(ctx context.Context, username, password string) (string, *models.User, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", nil, errors.New("username atau password salah")
	}

	if !user.IsActive {
		return "", nil, errors.New("akun Anda dinonaktifkan")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("username atau password salah")
	}

	jti := fmt.Sprintf("%d-%d", user.ID, time.Now().UnixNano())
	expTime := time.Now().Add(24 * time.Hour)

	claims := jwt.MapClaims{
		"jti":       jti,
		"user_id":   user.ID,
		"username":  user.Username,
		"role_name": user.RoleName,
		"exp":       expTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.RequireEnv("JWT_SECRET")))
	if err != nil {
		return "", nil, err
	}

	return tokenString, user, nil
}

func (s *authService) Logout(ctx context.Context, jti string, exp time.Time) error {
	return s.revokedRepo.Revoke(ctx, &models.RevokedToken{
		TokenJTI:  jti,
		ExpiresAt: exp,
		RevokedAt: time.Now(),
	})
}

func (s *authService) GetProfile(ctx context.Context, userID int) (*models.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}
