package controllers

import (
	"app/helpers"
	"app/models"
	"app/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoleController struct {
	service services.RoleService
}

func NewRoleController(s services.RoleService) *RoleController {
	return &RoleController{service: s}
}

// GetAll godoc
// @Summary      Daftar Semua Roles
// @Description  Mengambil seluruh data roles yang aktif (Admin Only)
// @Tags         Roles
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} helpers.APIResponse{data=[]models.Role} "Sukses"
// @Failure      401,403,500 {object} helpers.APIResponse "Error Otorisasi / Server"
// @Router       /roles [get]
func (ctl *RoleController) GetAll(c *gin.Context) {
	roles, err := ctl.service.GetAll(c.Request.Context())
	if err != nil {
		helpers.InternalServerError(c, err.Error())
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Daftar role berhasil dimuat", roles)
}

// GetByID godoc
// @Summary      Detail Role Berdasarkan ID
// @Description  Mengambil data detail role berdasarkan parameter ID (Admin Only)
// @Tags         Roles
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Role ID"
// @Success      200 {object} helpers.APIResponse{data=models.Role} "Sukses"
// @Failure      401,403,404 {object} helpers.APIResponse "Not Found / Error"
// @Router       /roles/{id} [get]
func (ctl *RoleController) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	role, err := ctl.service.GetByID(c.Request.Context(), id)
	if err != nil {
		helpers.ResponseError(c, http.StatusNotFound, "Role tidak ditemukan", nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Detail role ditemukan", role)
}

// GetDeleted godoc
// @Summary      Daftar Role Terhapus (Trash)
// @Description  Memuat seluruh data role yang berstatus soft-deleted (Admin Only)
// @Tags         Roles
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} helpers.APIResponse{data=[]models.Role} "Sukses"
// @Router       /roles/deleted [get]
func (ctl *RoleController) GetDeleted(c *gin.Context) {
	roles, err := ctl.service.GetDeleted(c.Request.Context())
	if err != nil {
		helpers.InternalServerError(c, err.Error())
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Daftar role terhapus", roles)
}

// Create godoc
// @Summary      Tambah Role Baru
// @Description  Membuat role baru dalam sistem (Admin Only)
// @Tags         Roles
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body models.RoleCreateDTO true "Payload Role"
// @Success      201 {object} helpers.APIResponse{data=models.Role} "Created"
// @Failure      400,401,403 {object} helpers.APIResponse "Bad Request / Forbidden"
// @Router       /roles [post]
func (ctl *RoleController) Create(c *gin.Context) {
	var req models.RoleCreateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, "Data input tidak valid", err.Error())
		return
	}

	userID := c.GetInt("user_id")
	createdRole, err := ctl.service.Create(c.Request.Context(), req, userID)
	if err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusCreated, "Role berhasil dibuat", createdRole)
}

// Update godoc
// @Summary      Update Role (Flexible Partial)
// @Description  Memperbarui data nama/deskripsi role secara parsial (Admin Only)
// @Tags         Roles
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path int true "Role ID"
// @Param        request body models.RoleUpdateDTO true "Payload Update Role"
// @Success      200 {object} helpers.APIResponse "Sukses"
// @Failure      400,401,403 {object} helpers.APIResponse "Bad Request"
// @Router       /roles/{id} [put]
func (ctl *RoleController) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req models.RoleUpdateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, "Data input tidak valid", err.Error())
		return
	}

	userID := c.GetInt("user_id")
	if err := ctl.service.Update(c.Request.Context(), id, req, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Role berhasil diperbarui", nil)
}

// Delete godoc
// @Summary      Soft Delete Role
// @Description  Menonaktifkan role ke status terhapus sementara (Admin Only)
// @Tags         Roles
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Role ID"
// @Success      200 {object} helpers.APIResponse "Sukses"
// @Failure      400,401,403 {object} helpers.APIResponse "Bad Request"
// @Router       /roles/{id} [delete]
func (ctl *RoleController) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetInt("user_id")

	if err := ctl.service.Delete(c.Request.Context(), id, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Role berhasil dihapus (soft delete)", nil)
}

// Restore godoc
// @Summary      Restore Role
// @Description  Memulihkan kembali role yang berstatus soft-deleted (Admin Only)
// @Tags         Roles
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Role ID"
// @Success      200 {object} helpers.APIResponse "Sukses"
// @Failure      400,401,403 {object} helpers.APIResponse "Bad Request"
// @Router       /roles/{id} [post]
func (ctl *RoleController) Restore(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetInt("user_id")

	if err := ctl.service.Restore(c.Request.Context(), id, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Role berhasil dipulihkan (restore)", nil)
}

// ForceDelete godoc
// @Summary      Delete Permanent Role (Force Delete)
// @Description  Menghapus fisik data role secara permanen dari database (Admin Only)
// @Tags         Roles
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Role ID"
// @Success      200 {object} helpers.APIResponse "Sukses"
// @Failure      400,401,403 {object} helpers.APIResponse "Bad Request / Foreign Key Constraint"
// @Router       /roles/{id}/force [delete]
func (ctl *RoleController) ForceDelete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := ctl.service.ForceDelete(c.Request.Context(), id); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Role berhasil dihapus secara permanen dari database", nil)
}
