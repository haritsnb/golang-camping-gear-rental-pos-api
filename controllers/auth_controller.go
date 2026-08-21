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
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

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

func (ctl *AuthController) Me(c *gin.Context) {
	userID := c.GetInt("user_id")
	user, err := ctl.authService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		helpers.ResponseError(c, http.StatusNotFound, "Profil user tidak ditemukan", nil)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Berhasil memuat profil", user)
}
