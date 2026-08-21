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

func (ctl *MaintenanceController) GetAll(c *gin.Context) {
	list, err := ctl.service.GetAll(c.Request.Context())
	if err != nil {
		helpers.InternalServerError(c, err.Error())
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Daftar pemeliharaan/servis alat", list)
}

func (ctl *MaintenanceController) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	m, err := ctl.service.GetByID(c.Request.Context(), id)
	if err != nil {
		helpers.ResponseError(c, http.StatusNotFound, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Detail pemeliharaan alat", m)
}

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

func (ctl *MaintenanceController) Complete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetInt("user_id")

	if err := ctl.service.Complete(c.Request.Context(), id, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Servis selesai. Status unit fisik kembali TERSEDIA (available).", nil)
}

func (ctl *MaintenanceController) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetInt("user_id")
	if err := ctl.service.Delete(c.Request.Context(), id, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Data maintenance berhasil dihapus (soft delete)", nil)
}

func (ctl *MaintenanceController) Restore(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetInt("user_id")
	if err := ctl.service.Restore(c.Request.Context(), id, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Data maintenance berhasil dipulihkan (restore)", nil)
}

func (ctl *MaintenanceController) ForceDelete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := ctl.service.ForceDelete(c.Request.Context(), id); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Data maintenance berhasil dihapus secara permanen dari database", nil)
}
