package middlewares

import (
	"app/helpers"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func RoleGuard(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("role_name")
		if !exists {
			helpers.ResponseError(c, http.StatusForbidden, "Akses ditolak: role tidak ditemukan", nil)
			c.Abort()
			return
		}

		userRole := roleVal.(string)
		for _, r := range allowedRoles {
			if strings.EqualFold(r, userRole) {
				c.Next()
				return
			}
		}

		helpers.ResponseError(c, http.StatusForbidden, "Akses ditolak: hak akses tidak mencukupi", nil)
		c.Abort()
	}
}
