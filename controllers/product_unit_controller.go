package controllers

import (
	"app/helpers"
	"app/models"
	"app/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductUnitController struct {
	service services.ProductUnitService
}

func NewProductUnitController(s services.ProductUnitService) *ProductUnitController {
	return &ProductUnitController{service: s}
}

func (ctl *ProductUnitController) GetAll(c *gin.Context) {
	units, err := ctl.service.GetAll(c.Request.Context())
	if err != nil {
		helpers.InternalServerError(c, err.Error())
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Daftar unit fisik produk berhasil dimuat", units)
}

func (ctl *ProductUnitController) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	unit, err := ctl.service.GetByID(c.Request.Context(), id)
	if err != nil {
		helpers.ResponseError(c, http.StatusNotFound, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Detail unit fisik ditemukan", unit)
}

func (ctl *ProductUnitController) ScanUnitCode(c *gin.Context) {
	unitCode := c.Param("unit_code")
	unit, err := ctl.service.GetByUnitCode(c.Request.Context(), unitCode)
	if err != nil {
		helpers.ResponseError(c, http.StatusNotFound, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Scan unit fisik berhasil", unit)
}

func (ctl *ProductUnitController) Create(c *gin.Context) {
	var req models.ProductUnitCreateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, "Data input tidak valid", err.Error())
		return
	}

	userID := c.GetInt("user_id")
	unit, err := ctl.service.Create(c.Request.Context(), req, userID)
	if err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusCreated, "Unit fisik produk berhasil didaftarkan", unit)
}

func (ctl *ProductUnitController) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req models.ProductUnitUpdateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, "Data input tidak valid", err.Error())
		return
	}

	userID := c.GetInt("user_id")
	if err := ctl.service.Update(c.Request.Context(), id, req, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Unit fisik berhasil diperbarui", nil)
}

func (ctl *ProductUnitController) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetInt("user_id")
	if err := ctl.service.Delete(c.Request.Context(), id, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Unit fisik berhasil dihapus (soft delete)", nil)
}

func (ctl *ProductUnitController) Restore(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetInt("user_id")
	if err := ctl.service.Restore(c.Request.Context(), id, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Unit fisik berhasil dipulihkan (restore)", nil)
}

func (ctl *ProductUnitController) ForceDelete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := ctl.service.ForceDelete(c.Request.Context(), id); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Unit fisik berhasil dihapus secara permanen dari database", nil)
}
