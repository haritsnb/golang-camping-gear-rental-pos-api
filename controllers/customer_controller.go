package controllers

import (
	"app/helpers"
	"app/models"
	"app/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CustomerController struct {
	service services.CustomerService
}

func NewCustomerController(s services.CustomerService) *CustomerController {
	return &CustomerController{service: s}
}

func (ctl *CustomerController) GetAll(c *gin.Context) {
	customers, err := ctl.service.GetAll(c.Request.Context())
	if err != nil {
		helpers.InternalServerError(c, err.Error())
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Daftar pelanggan berhasil dimuat", customers)
}

func (ctl *CustomerController) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	cust, err := ctl.service.GetByID(c.Request.Context(), id)
	if err != nil {
		helpers.ResponseError(c, http.StatusNotFound, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Detail pelanggan ditemukan", cust)
}

func (ctl *CustomerController) Create(c *gin.Context) {
	var dto models.CustomerCreateDTO
	if err := c.ShouldBind(&dto); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, "Data input tidak valid", err.Error())
		return
	}

	var photoURL *string
	file, err := c.FormFile("identity_photo")
	if err == nil {
		path, errSave := helpers.SaveUploadedFile(c, file, "storages/identities")
		if errSave != nil {
			helpers.ResponseError(c, http.StatusBadRequest, "Gagal mengunggah foto identitas", errSave.Error())
			return
		}
		photoURL = &path
	}

	userID := c.GetInt("user_id")
	createdCustomer, err := ctl.service.Create(c.Request.Context(), dto, photoURL, userID)
	if err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	helpers.ResponseSuccess(c, http.StatusCreated, "Customer KYC berhasil didaftarkan", createdCustomer)
}

func (ctl *CustomerController) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var dto models.CustomerUpdateDTO
	if err := c.ShouldBind(&dto); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, "Data input tidak valid", err.Error())
		return
	}

	var photoURL *string
	file, err := c.FormFile("identity_photo")
	if err == nil {
		path, errSave := helpers.SaveUploadedFile(c, file, "storages/identities")
		if errSave != nil {
			helpers.ResponseError(c, http.StatusBadRequest, "Gagal mengunggah foto identitas baru", errSave.Error())
			return
		}
		photoURL = &path
	}

	userID := c.GetInt("user_id")
	if err := ctl.service.Update(c.Request.Context(), id, dto, photoURL, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Customer berhasil diperbarui", nil)
}

func (ctl *CustomerController) Blacklist(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		IsBlacklisted bool `json:"is_blacklisted"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, "Format JSON salah", err.Error())
		return
	}

	userID := c.GetInt("user_id")
	if err := ctl.service.SetBlacklist(c.Request.Context(), id, req.IsBlacklisted, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Status blacklist customer berhasil diperbarui", nil)
}

func (ctl *CustomerController) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetInt("user_id")

	if err := ctl.service.Delete(c.Request.Context(), id, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Customer berhasil dihapus (soft delete)", nil)
}

func (ctl *CustomerController) Restore(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetInt("user_id")

	if err := ctl.service.Restore(c.Request.Context(), id, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Customer berhasil dipulihkan (restore)", nil)
}

func (ctl *CustomerController) ForceDelete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := ctl.service.ForceDelete(c.Request.Context(), id); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Customer berhasil dihapus secara permanen dari database", nil)
}
