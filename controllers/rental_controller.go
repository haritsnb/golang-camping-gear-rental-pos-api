package controllers

import (
	"app/helpers"
	"app/models"
	"app/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RentalController struct {
	service services.RentalService
}

func NewRentalController(s services.RentalService) *RentalController {
	return &RentalController{service: s}
}

// GetAll godoc
// @Summary      Daftar Semua Transaksi Rental
// @Description  Memuat daftar transaksi sewa alat camping dengan filter status opsional (?status=booked/active/returned/completed/cancelled)
// @Tags         Rentals
// @Security     BearerAuth
// @Produce      json
// @Param        status query string false "Filter Status Transaksi" Enums(booked, active, returned, completed, cancelled)
// @Success      200 {object} helpers.APIResponse{data=[]models.Rental} "Sukses"
// @Failure      401,500 {object} helpers.APIResponse "Error"
// @Router       /rentals [get]
func (ctl *RentalController) GetAll(c *gin.Context) {
	status := c.Query("status")
	rentals, err := ctl.service.GetAll(c.Request.Context(), status)
	if err != nil {
		helpers.InternalServerError(c, err.Error())
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Daftar transaksi rental", rentals)
}

// GetByID godoc
// @Summary      Detail Invoice Rental
// @Description  Memuat rincian transaksi invoice sewa beserta rincian item unit fisik yang disewa
// @Tags         Rentals
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Rental ID"
// @Success      200 {object} helpers.APIResponse{data=models.Rental} "Sukses"
// @Failure      401,404 {object} helpers.APIResponse "Not Found"
// @Router       /rentals/{id} [get]
func (ctl *RentalController) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	rental, err := ctl.service.GetByID(c.Request.Context(), id)
	if err != nil {
		helpers.ResponseError(c, http.StatusNotFound, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Detail invoice rental", rental)
}

// GetDeleted godoc
// @Summary      Daftar Rental Terhapus (Trash)
// @Tags         Rentals
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} helpers.APIResponse{data=[]models.Rental} "Sukses"
// @Router       /rentals/deleted [get]
func (ctl *RentalController) GetDeleted(c *gin.Context) {
	rentals, err := ctl.service.GetDeleted(c.Request.Context())
	if err != nil {
		helpers.InternalServerError(c, err.Error())
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Daftar transaksi rental terhapus", rentals)
}

// Booking godoc
// @Summary      Booking Sewa Alat (Fase 1)
// @Description  Membuat pesanan sewa baru, mengunci unit fisik, menghitung durasi hari, deposit, dan mencatat DP
// @Tags         Rentals
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body models.RentalBookingDTO true "Payload Booking Sewa"
// @Success      201 {object} helpers.APIResponse{data=models.Rental} "Created"
// @Failure      400,401,500 {object} helpers.APIResponse "Bad Request / Blacklisted Customer / Unit Unavailable"
// @Router       /rentals/booking [post]
func (ctl *RentalController) Booking(c *gin.Context) {
	var req models.RentalBookingDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, "Payload booking tidak valid", err.Error())
		return
	}

	userID := c.GetInt("user_id")
	rental, err := ctl.service.Booking(c.Request.Context(), req, userID)
	if err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	helpers.ResponseSuccess(c, http.StatusCreated, "Booking berhasil dibuat", rental)
}

// Handover godoc
// @Summary      Handover Alat Fisik (Fase 2)
// @Description  Menyerahkan unit alat ke pelanggan saat pengambilan. Status sewa menjadi 'active' dan unit fisik menjadi 'rented'
// @Tags         Rentals
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Rental ID"
// @Success      200 {object} helpers.APIResponse "Sukses"
// @Failure      400,401,500 {object} helpers.APIResponse "Bad Request"
// @Router       /rentals/{id}/handover [post]
func (ctl *RentalController) Handover(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetInt("user_id")

	if err := ctl.service.Handover(c.Request.Context(), id, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Handover selesai. Status sewa kini AKTIF dan unit RENTED", nil)
}

// Return godoc
// @Summary      Pengembalian Alat & Cek Denda (Fase 3)
// @Description  Mencatat pengembalian alat aktual, menghitung denda telat otomatis, dan mengevaluasi denda kerusakan fisik tiap unit
// @Tags         Rentals
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path int true "Rental ID"
// @Param        request body models.RentalReturnDTO true "Kondisi Fisik Unit Kembali"
// @Success      200 {object} helpers.APIResponse{data=models.Rental} "Sukses"
// @Failure      400,401,500 {object} helpers.APIResponse "Bad Request"
// @Router       /rentals/{id}/return [post]
func (ctl *RentalController) Return(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req models.RentalReturnDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, "Payload pengembalian tidak valid", err.Error())
		return
	}

	userID := c.GetInt("user_id")
	rental, err := ctl.service.Return(c.Request.Context(), id, req, userID)
	if err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Pengembalian berhasil dicatat. Lanjutkan ke tahap settlement.", rental)
}

// Settlement godoc
// @Summary      Settlement & Rekonsiliasi Deposit (Fase 4)
// @Description  Menyelesaikan transaksi sewa, refund sisa deposit atau tagih denda, dan mengupdate status unit fisik (available / maintenance)
// @Tags         Rentals
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path int true "Rental ID"
// @Param        request body models.RentalSettlementDTO true "Metode Penyelesaian"
// @Success      200 {object} helpers.APIResponse "Sukses"
// @Failure      400,401,500 {object} helpers.APIResponse "Bad Request"
// @Router       /rentals/{id}/settlement [post]
func (ctl *RentalController) Settlement(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req models.RentalSettlementDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, "Payload settlement tidak valid", err.Error())
		return
	}

	userID := c.GetInt("user_id")
	if err := ctl.service.Settlement(c.Request.Context(), id, req, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Settlement berhasil. Transaksi sewa SELESAI.", nil)
}

// Cancel godoc
// @Summary      Batalkan Booking Sewa
// @Description  Membatalkan pesanan booking yang belum diserahterimakan dan mengembalikan status unit fisik menjadi available
// @Tags         Rentals
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Rental ID"
// @Success      200 {object} helpers.APIResponse "Sukses"
// @Failure      400,401,500 {object} helpers.APIResponse "Bad Request"
// @Router       /rentals/{id}/cancel [post]
func (ctl *RentalController) Cancel(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetInt("user_id")

	if err := ctl.service.Cancel(c.Request.Context(), id, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Booking rental berhasil dibatalkan", nil)
}

// Update godoc
// @Summary      Update Transaksi Rental (Flexible Partial)
// @Description  Memperbarui tanggal, total biaya, atau status transaksi rental secara fleksibel (Admin Only)
// @Tags         Rentals
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path int true "Rental ID"
// @Param        request body models.RentalUpdateDTO true "Payload Update Rental"
// @Success      200 {object} helpers.APIResponse "Sukses"
// @Failure      400,401,403 {object} helpers.APIResponse "Bad Request"
// @Router       /rentals/{id} [put]
func (ctl *RentalController) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req models.RentalUpdateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, "Data input tidak valid", err.Error())
		return
	}

	userID := c.GetInt("user_id")
	if err := ctl.service.Update(c.Request.Context(), id, req, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Data rental berhasil diperbarui", nil)
}

// Delete godoc
// @Summary      Soft Delete Transaksi Rental
// @Description  Menonaktifkan transaksi rental dan riwayat rental itemnya (Admin Only)
// @Tags         Rentals
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Rental ID"
// @Success      200 {object} helpers.APIResponse "Sukses"
// @Failure      400,401,403 {object} helpers.APIResponse "Bad Request"
// @Router       /rentals/{id} [delete]
func (ctl *RentalController) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetInt("user_id")
	if err := ctl.service.Delete(c.Request.Context(), id, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Rental berhasil dihapus (soft delete)", nil)
}

// Restore godoc
// @Summary      Restore Transaksi Rental
// @Description  Memulihkan kembali transaksi rental dan rental item yang berstatus soft-deleted (Admin Only)
// @Tags         Rentals
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Rental ID"
// @Success      200 {object} helpers.APIResponse "Sukses"
// @Failure      400,401,403 {object} helpers.APIResponse "Bad Request"
// @Router       /rentals/{id} [post]
func (ctl *RentalController) Restore(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetInt("user_id")
	if err := ctl.service.Restore(c.Request.Context(), id, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Rental berhasil dipulihkan (restore)", nil)
}

// ForceDelete godoc
// @Summary      Delete Permanent Rental (Force Delete)
// @Description  Menghapus fisik data transaksi rental beserta itemnya secara permanen dari database (Admin Only)
// @Tags         Rentals
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Rental ID"
// @Success      200 {object} helpers.APIResponse "Sukses"
// @Failure      400,401,403 {object} helpers.APIResponse "Bad Request / Foreign Key Error"
// @Router       /rentals/{id}/force [delete]
func (ctl *RentalController) ForceDelete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := ctl.service.ForceDelete(c.Request.Context(), id); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Rental berhasil dihapus secara permanen dari database", nil)
}
