package controllers

import (
	"app/helpers"
	"app/models"
	"app/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaymentController struct {
	service services.PaymentService
}

func NewPaymentController(s services.PaymentService) *PaymentController {
	return &PaymentController{service: s}
}

// GetByRentalID godoc
// @Summary      Histori Pembayaran Berdasarkan Rental ID
// @Description  Memuat seluruh log pembayaran transaksi sewa (DP, pelunasan, deposit, refund, denda)
// @Tags         Payments
// @Security     BearerAuth
// @Produce      json
// @Param        rental_id path int true "Rental ID"
// @Success      200 {object} helpers.APIResponse{data=[]models.Payment} "Sukses"
// @Failure      401,500 {object} helpers.APIResponse "Error"
// @Router       /payments/rental/{rental_id} [get]
func (ctl *PaymentController) GetByRentalID(c *gin.Context) {
	rentalID, _ := strconv.Atoi(c.Param("rental_id"))
	payments, err := ctl.service.GetByRentalID(c.Request.Context(), rentalID)
	if err != nil {
		helpers.InternalServerError(c, err.Error())
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Histori pembayaran rental berhasil dimuat", payments)
}

// GetByID godoc
// @Summary      Detail Pembayaran Berdasarkan ID
// @Description  Memuat detail satu bukti transaksi pembayaran
// @Tags         Payments
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Payment ID"
// @Success      200 {object} helpers.APIResponse{data=models.Payment} "Sukses"
// @Failure      401,404 {object} helpers.APIResponse "Not Found"
// @Router       /payments/{id} [get]
func (ctl *PaymentController) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	payment, err := ctl.service.GetByID(c.Request.Context(), id)
	if err != nil {
		helpers.ResponseError(c, http.StatusNotFound, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Detail pembayaran ditemukan", payment)
}

// Create godoc
// @Summary      Catat Pembayaran Manual
// @Description  Mencatat pembayaran tambahan/pelunasan manual untuk rental tertentu
// @Tags         Payments
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body models.PaymentCreateDTO true "Payload Pembayaran"
// @Success      201 {object} helpers.APIResponse{data=models.Payment} "Created"
// @Failure      400,401,500 {object} helpers.APIResponse "Bad Request"
// @Router       /payments [post]
func (ctl *PaymentController) Create(c *gin.Context) {
	var req models.PaymentCreateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, "Payload pembayaran tidak valid", err.Error())
		return
	}

	userID := c.GetInt("user_id")
	payment, err := ctl.service.Create(c.Request.Context(), req, userID)
	if err != nil {
		helpers.InternalServerError(c, err.Error())
		return
	}
	helpers.ResponseSuccess(c, http.StatusCreated, "Pembayaran berhasil dicatat", payment)
}

// Update godoc
// @Summary      Update Pembayaran (Flexible Partial)
// @Description  Memperbarui nominal, metode, atau nomor referensi pembayaran (Admin Only)
// @Tags         Payments
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path int true "Payment ID"
// @Param        request body models.PaymentUpdateDTO true "Payload Update Pembayaran"
// @Success      200 {object} helpers.APIResponse "Sukses"
// @Failure      400,401,403 {object} helpers.APIResponse "Bad Request"
// @Router       /payments/{id} [put]
func (ctl *PaymentController) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req models.PaymentUpdateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, "Data input tidak valid", err.Error())
		return
	}

	userID := c.GetInt("user_id")
	if err := ctl.service.Update(c.Request.Context(), id, req, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Pembayaran berhasil diperbarui", nil)
}

// Delete godoc
// @Summary      Soft Delete Pembayaran
// @Description  Menonaktifkan data pembayaran (Admin Only)
// @Tags         Payments
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Payment ID"
// @Success      200 {object} helpers.APIResponse "Sukses"
// @Failure      400,401,403 {object} helpers.APIResponse "Bad Request"
// @Router       /payments/{id} [delete]
func (ctl *PaymentController) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetInt("user_id")
	if err := ctl.service.Delete(c.Request.Context(), id, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Pembayaran berhasil dihapus (soft delete)", nil)
}

// Restore godoc
// @Summary      Restore Pembayaran
// @Description  Memulihkan kembali data pembayaran yang berstatus soft-deleted (Admin Only)
// @Tags         Payments
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Payment ID"
// @Success      200 {object} helpers.APIResponse "Sukses"
// @Failure      400,401,403 {object} helpers.APIResponse "Bad Request"
// @Router       /payments/{id} [post]
func (ctl *PaymentController) Restore(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetInt("user_id")
	if err := ctl.service.Restore(c.Request.Context(), id, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Pembayaran berhasil dipulihkan (restore)", nil)
}

// ForceDelete godoc
// @Summary      Delete Permanent Pembayaran (Force Delete)
// @Description  Menghapus fisik data pembayaran secara permanen dari database (Admin Only)
// @Tags         Payments
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Payment ID"
// @Success      200 {object} helpers.APIResponse "Sukses"
// @Failure      400,401,403 {object} helpers.APIResponse "Bad Request"
// @Router       /payments/{id}/force [delete]
func (ctl *PaymentController) ForceDelete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := ctl.service.ForceDelete(c.Request.Context(), id); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Pembayaran berhasil dihapus secara permanen dari database", nil)
}
