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

func (ctl *RentalController) GetAll(c *gin.Context) {
	status := c.Query("status")
	rentals, err := ctl.service.GetAll(c.Request.Context(), status)
	if err != nil {
		helpers.InternalServerError(c, err.Error())
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Daftar transaksi rental", rentals)
}

func (ctl *RentalController) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	rental, err := ctl.service.GetByID(c.Request.Context(), id)
	if err != nil {
		helpers.ResponseError(c, http.StatusNotFound, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Detail invoice rental", rental)
}

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

func (ctl *RentalController) Handover(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetInt("user_id")

	if err := ctl.service.Handover(c.Request.Context(), id, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Handover selesai. Status sewa kini AKTIF dan unit RENTED", nil)
}

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

func (ctl *RentalController) Cancel(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetInt("user_id")

	if err := ctl.service.Cancel(c.Request.Context(), id, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Booking rental berhasil dibatalkan", nil)
}

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

func (ctl *RentalController) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetInt("user_id")
	if err := ctl.service.Delete(c.Request.Context(), id, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Rental berhasil dihapus (soft delete)", nil)
}

func (ctl *RentalController) Restore(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetInt("user_id")
	if err := ctl.service.Restore(c.Request.Context(), id, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Rental berhasil dipulihkan (restore)", nil)
}

func (ctl *RentalController) ForceDelete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := ctl.service.ForceDelete(c.Request.Context(), id); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Rental berhasil dihapus secara permanen dari database", nil)
}
