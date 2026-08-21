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

func (ctl *PaymentController) GetByRentalID(c *gin.Context) {
	rentalID, _ := strconv.Atoi(c.Param("rental_id"))
	payments, err := ctl.service.GetByRentalID(c.Request.Context(), rentalID)
	if err != nil {
		helpers.InternalServerError(c, err.Error())
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Histori pembayaran rental berhasil dimuat", payments)
}

func (ctl *PaymentController) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	payment, err := ctl.service.GetByID(c.Request.Context(), id)
	if err != nil {
		helpers.ResponseError(c, http.StatusNotFound, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Detail pembayaran ditemukan", payment)
}

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

func (ctl *PaymentController) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetInt("user_id")
	if err := ctl.service.Delete(c.Request.Context(), id, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Pembayaran berhasil dihapus (soft delete)", nil)
}

func (ctl *PaymentController) Restore(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetInt("user_id")
	if err := ctl.service.Restore(c.Request.Context(), id, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Pembayaran berhasil dipulihkan (restore)", nil)
}

func (ctl *PaymentController) ForceDelete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := ctl.service.ForceDelete(c.Request.Context(), id); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Pembayaran berhasil dihapus secara permanen dari database", nil)
}
