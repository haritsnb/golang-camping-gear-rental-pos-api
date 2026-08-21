package middlewares

import (
	"app/config"
	"app/helpers"
	"app/repositories"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	RoleName string `json:"role_name"`
	jwt.RegisteredClaims
}

func AuthMiddleware(revokedRepo repositories.RevokedTokenRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			helpers.ResponseError(c, http.StatusUnauthorized, "Header otorisasi diperlukan", nil)
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			helpers.ResponseError(c, http.StatusUnauthorized, "Format token invalid (gunakan Bearer <token>)", nil)
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims := &JWTClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.RequireEnv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			helpers.ResponseError(c, http.StatusUnauthorized, "Token tidak valid atau kadaluarsa", nil)
			c.Abort()
			return
		}

		// Periksa apakah token ada di daftar revoked/blacklist
		isRevoked, err := revokedRepo.IsRevoked(c.Request.Context(), claims.ID)
		if err != nil || isRevoked {
			helpers.ResponseError(c, http.StatusUnauthorized, "Sesi login telah berakhir / token telah di-logout", nil)
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role_name", claims.RoleName)
		c.Set("jti", claims.ID)
		c.Set("exp", claims.ExpiresAt.Time)

		c.Next()
	}
}
