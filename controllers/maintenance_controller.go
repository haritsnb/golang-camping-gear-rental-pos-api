package controllers

import (
	"app/helpers"
	"app/models"
	"app/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MaintenanceController struct {
	service services.MaintenanceService
}

func NewMaintenanceController(s services.MaintenanceService) *MaintenanceController {
	return &MaintenanceController{service: s}
}

// GetAll godoc
// @Summary      Daftar Semua Servis/Maintenance
// @Description  Memuat seluruh tiket servis dan perbaikan alat camping
// @Tags         Maintenance
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} helpers.APIResponse{data=[]models.Maintenance} "Sukses"
// @Failure      401,500 {object} helpers.APIResponse "Error"
// @Router       /maintenance [get]
func (ctl *MaintenanceController) GetAll(c *gin.Context) {
	list, err := ctl.service.GetAll(c.Request.Context())
	if err != nil {
		helpers.InternalServerError(c, err.Error())
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Daftar pemeliharaan/servis alat", list)
}

// GetByID godoc
// @Summary      Detail Servis Berdasarkan ID
// @Description  Memuat rincian tiket perbaikan alat dan status pengerjaannya
// @Tags         Maintenance
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Maintenance ID"
// @Success      200 {object} helpers.APIResponse{data=models.Maintenance} "Sukses"
// @Failure      401,404 {object} helpers.APIResponse "Not Found"
// @Router       /maintenance/{id} [get]
func (ctl *MaintenanceController) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	m, err := ctl.service.GetByID(c.Request.Context(), id)
	if err != nil {
		helpers.ResponseError(c, http.StatusNotFound, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Detail pemeliharaan alat", m)
}

// GetDeleted godoc
// @Summary      Daftar Servis Terhapus (Trash)
// @Tags         Maintenance
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} helpers.APIResponse{data=[]models.Maintenance} "Sukses"
// @Router       /maintenance/deleted [get]
func (ctl *MaintenanceController) GetDeleted(c *gin.Context) {
	list, err := ctl.service.GetDeleted(c.Request.Context())
	if err != nil {
		helpers.InternalServerError(c, err.Error())
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Daftar tiket servis terhapus", list)
}

// Create godoc
// @Summary      Input Tiket Servis Manual
// @Description  Mendaftarkan unit rusak ke modul servis dan mengubah status unit fisik menjadi 'maintenance' (Admin/Staff Only)
// @Tags         Maintenance
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body models.MaintenanceCreateDTO true "Payload Maintenance"
// @Success      201 {object} helpers.APIResponse{data=models.Maintenance} "Created"
// @Failure      400,401,403,500 {object} helpers.APIResponse "Bad Request"
// @Router       /maintenance [post]
func (ctl *MaintenanceController) Create(c *gin.Context) {
	var req models.MaintenanceCreateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, "Payload maintenance tidak valid", err.Error())
		return
	}

	userID := c.GetInt("user_id")
	m, err := ctl.service.Create(c.Request.Context(), req, userID)
	if err != nil {
		helpers.InternalServerError(c, err.Error())
		return
	}
	helpers.ResponseSuccess(c, http.StatusCreated, "Data servis/pemeliharaan alat berhasil ditambahkan", m)
}

// Complete godoc
// @Summary      Selesaikan Perbaikan Unit
// @Description  Menyelesaikan tiket servis dan secara otomatis mengembalikan status unit fisik menjadi 'available' (Admin/Staff Only)
// @Tags         Maintenance
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Maintenance ID"
// @Success      200 {object} helpers.APIResponse "Sukses"
// @Failure      400,401,403 {object} helpers.APIResponse "Bad Request"
// @Router       /maintenance/{id}/complete [patch]
func (ctl *MaintenanceController) Complete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetInt("user_id")

	if err := ctl.service.Complete(c.Request.Context(), id, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Servis selesai. Status unit fisik kembali TERSEDIA (available).", nil)
}

// Update godoc
// @Summary      Update Tiket Servis (Flexible Partial)
// @Description  Memperbarui deskripsi kerusakan, estimasi biaya, atau tanggal mulai servis (Admin/Staff Only)
// @Tags         Maintenance
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path int true "Maintenance ID"
// @Param        request body models.MaintenanceUpdateDTO true "Payload Update Maintenance"
// @Success      200 {object} helpers.APIResponse "Sukses"
// @Failure      400,401,403 {object} helpers.APIResponse "Bad Request"
// @Router       /maintenance/{id} [put]
func (ctl *MaintenanceController) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req models.MaintenanceUpdateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, "Data input tidak valid", err.Error())
		return
	}

	userID := c.GetInt("user_id")
	if err := ctl.service.Update(c.Request.Context(), id, req, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Data maintenance berhasil diperbarui", nil)
}

// Delete godoc
// @Summary      Soft Delete Tiket Servis
// @Description  Menonaktifkan catatan tiket pemeliharaan (Admin Only)
// @Tags         Maintenance
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Maintenance ID"
// @Success      200 {object} helpers.APIResponse "Sukses"
// @Failure      400,401,403 {object} helpers.APIResponse "Bad Request"
// @Router       /maintenance/{id} [delete]
func (ctl *MaintenanceController) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetInt("user_id")
	if err := ctl.service.Delete(c.Request.Context(), id, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Data maintenance berhasil dihapus (soft delete)", nil)
}

// Restore godoc
// @Summary      Restore Tiket Servis
// @Description  Memulihkan kembali tiket servis yang berstatus soft-deleted (Admin Only)
// @Tags         Maintenance
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Maintenance ID"
// @Success      200 {object} helpers.APIResponse "Sukses"
// @Failure      400,401,403 {object} helpers.APIResponse "Bad Request"
// @Router       /maintenance/{id} [post]
func (ctl *MaintenanceController) Restore(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetInt("user_id")
	if err := ctl.service.Restore(c.Request.Context(), id, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Data maintenance berhasil dipulihkan (restore)", nil)
}

// ForceDelete godoc
// @Summary      Delete Permanent Tiket Servis (Force Delete)
// @Description  Menghapus fisik data tiket servis secara permanen dari database (Admin Only)
// @Tags         Maintenance
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Maintenance ID"
// @Success      200 {object} helpers.APIResponse "Sukses"
// @Failure      400,401,403 {object} helpers.APIResponse "Bad Request"
// @Router       /maintenance/{id}/force [delete]
func (ctl *MaintenanceController) ForceDelete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := ctl.service.ForceDelete(c.Request.Context(), id); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Data maintenance berhasil dihapus secara permanen dari database", nil)
}
