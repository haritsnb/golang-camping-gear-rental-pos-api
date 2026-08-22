package controllers

import (
	"app/helpers"
	"app/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService services.AuthService
}

func NewAuthController(a services.AuthService) *AuthController {
	return &AuthController{authService: a}
}

type LoginDTO struct {
	Username string `json:"username" binding:"required" example:"admin"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// Login godoc
// @Summary      Login User
// @Description  Autentikasi kredensial pengguna dan mengembalikan Bearer JWT Token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body LoginDTO true "Payload Kredensial Login"
// @Success      200 {object} helpers.APIResponse{data=object{token=string,user=models.User}} "Login Berhasil"
// @Failure      400,401 {object} helpers.APIResponse "Bad Request / Unauthorized"
// @Router       /auth/login [post]
func (ctl *AuthController) Login(c *gin.Context) {
	var req LoginDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, "Payload request tidak valid", err.Error())
		return
	}

	token, user, err := ctl.authService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		helpers.ResponseError(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Login berhasil", gin.H{
		"token": token,
		"user":  user,
	})
}

// Logout godoc
// @Summary      Logout User
// @Description  Revoke token JWT aktif pengguna saat ini dan memasukkannya ke tabel blacklist
// @Tags         Auth
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} helpers.APIResponse "Logout Berhasil"
// @Failure      401,500 {object} helpers.APIResponse "Unauthorized / Internal Server Error"
// @Router       /auth/logout [post]
func (ctl *AuthController) Logout(c *gin.Context) {
	jti := c.GetString("jti")
	expVal, _ := c.Get("exp")
	exp := expVal.(time.Time)

	if err := ctl.authService.Logout(c.Request.Context(), jti, exp); err != nil {
		helpers.InternalServerError(c, "Gagal melakukan logout")
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Logout berhasil, token di-revoke", nil)
}

// Me godoc
// @Summary      Get Profile Pengguna
// @Description  Memuat data profil pengguna yang sedang login berdasarkan klaim token JWT
// @Tags         Auth
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} helpers.APIResponse{data=models.User} "Profil Ditemukan"
// @Failure      401,404 {object} helpers.APIResponse "Unauthorized / Not Found"
// @Router       /auth/me [get]
func (ctl *AuthController) Me(c *gin.Context) {
	userID := c.GetInt("user_id")
	user, err := ctl.authService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		helpers.ResponseError(c, http.StatusNotFound, "Profil user tidak ditemukan", nil)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Berhasil memuat profil", user)
}
