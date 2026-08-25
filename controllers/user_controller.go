package controllers

import (
	"app/helpers"
	"app/models"
	"app/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service services.UserService
}

func NewUserController(s services.UserService) *UserController {
	return &UserController{service: s}
}

// GetAll godoc
// @Summary      Daftar Semua Pengguna
// @Description  Mengambil seluruh data user akun operasional toko (Admin Only)
// @Tags         Users
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} helpers.APIResponse{data=[]models.User} "Sukses"
// @Failure      401,403,500 {object} helpers.APIResponse "Error"
// @Router       /users [get]
func (ctl *UserController) GetAll(c *gin.Context) {
	users, err := ctl.service.GetAll(c.Request.Context())
	if err != nil {
		helpers.InternalServerError(c, err.Error())
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Daftar user berhasil dimuat", users)
}

// GetByID godoc
// @Summary      Detail User Berdasarkan ID
// @Description  Mengambil data detail akun user berdasarkan ID (Admin Only)
// @Tags         Users
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "User ID"
// @Success      200 {object} helpers.APIResponse{data=models.User} "Sukses"
// @Failure      401,403,404 {object} helpers.APIResponse "Not Found"
// @Router       /users/{id} [get]
func (ctl *UserController) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := ctl.service.GetByID(c.Request.Context(), id)
	if err != nil {
		helpers.ResponseError(c, http.StatusNotFound, "User tidak ditemukan", nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Detail user ditemukan", user)
}

// GetDeleted godoc
// @Summary      Daftar User Terhapus (Trash)
// @Tags         Users
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} helpers.APIResponse{data=[]models.User} "Sukses"
// @Router       /users/deleted [get]
func (ctl *UserController) GetDeleted(c *gin.Context) {
	users, err := ctl.service.GetDeleted(c.Request.Context())
	if err != nil {
		helpers.InternalServerError(c, err.Error())
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Daftar user terhapus", users)
}

// Create godoc
// @Summary      Tambah Akun User Baru
// @Description  Mendaftarkan akun user baru (Kasir/Staff/Admin) beserta enkripsi bcrypt password (Admin Only)
// @Tags         Users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body models.UserCreateDTO true "Payload Registrasi User"
// @Success      201 {object} helpers.APIResponse{data=models.User} "Created"
// @Failure      400,401,403 {object} helpers.APIResponse "Bad Request"
// @Router       /users [post]
func (ctl *UserController) Create(c *gin.Context) {
	var req models.UserCreateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, "Data input tidak valid", err.Error())
		return
	}

	userID := c.GetInt("user_id")
	createdUser, err := ctl.service.Create(c.Request.Context(), req, userID)
	if err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusCreated, "User berhasil didaftarkan", createdUser)
}

// Update godoc
// @Summary      Update User (Flexible Partial)
// @Description  Memperbarui data profil, role, username, password, atau status aktif user secara fleksibel (Admin Only)
// @Tags         Users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path int true "User ID"
// @Param        request body models.UserUpdateDTO true "Payload Update User"
// @Success      200 {object} helpers.APIResponse "Sukses"
// @Failure      400,401,403 {object} helpers.APIResponse "Bad Request"
// @Router       /users/{id} [put]
func (ctl *UserController) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req models.UserUpdateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, "Data input tidak valid", err.Error())
		return
	}

	userID := c.GetInt("user_id")
	if err := ctl.service.Update(c.Request.Context(), id, req, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "User berhasil diperbarui", nil)
}

// Delete godoc
// @Summary      Soft Delete User
// @Description  Menonaktifkan akun user ke status terhapus sementara (Admin Only)
// @Tags         Users
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "User ID"
// @Success      200 {object} helpers.APIResponse "Sukses"
// @Failure      400,401,403 {object} helpers.APIResponse "Bad Request"
// @Router       /users/{id} [delete]
func (ctl *UserController) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetInt("user_id")

	if err := ctl.service.Delete(c.Request.Context(), id, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "User berhasil dihapus (soft delete)", nil)
}

// Restore godoc
// @Summary      Restore User
// @Description  Memulihkan kembali akun user yang berstatus soft-deleted (Admin Only)
// @Tags         Users
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "User ID"
// @Success      200 {object} helpers.APIResponse "Sukses"
// @Failure      400,401,403 {object} helpers.APIResponse "Bad Request"
// @Router       /users/{id} [post]
func (ctl *UserController) Restore(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetInt("user_id")

	if err := ctl.service.Restore(c.Request.Context(), id, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "User berhasil dipulihkan (restore)", nil)
}

// ForceDelete godoc
// @Summary      Delete Permanent User (Force Delete)
// @Description  Menghapus fisik data akun user secara permanen dari database (Admin Only)
// @Tags         Users
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "User ID"
// @Success      200 {object} helpers.APIResponse "Sukses"
// @Failure      400,401,403 {object} helpers.APIResponse "Bad Request / Foreign Key Error"
// @Router       /users/{id}/force [delete]
func (ctl *UserController) ForceDelete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetInt("user_id")

	if err := ctl.service.ForceDelete(c.Request.Context(), id, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "User berhasil dihapus secara permanen dari database", nil)
}
