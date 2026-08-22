package controllers

import (
	"app/helpers"
	"app/models"
	"app/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BrandController struct {
	service services.BrandService
}

func NewBrandController(s services.BrandService) *BrandController {
	return &BrandController{service: s}
}

// GetAll godoc
// @Summary      Daftar Semua Brands
// @Description  Memuat seluruh brand peralatan outdoor yang terdaftar
// @Tags         Brands
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} helpers.APIResponse{data=[]models.Brand} "Sukses"
// @Failure      401,500 {object} helpers.APIResponse "Error"
// @Router       /brands [get]
func (ctl *BrandController) GetAll(c *gin.Context) {
	brands, err := ctl.service.GetAll(c.Request.Context())
	if err != nil {
		helpers.InternalServerError(c, err.Error())
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Daftar brand berhasil dimuat", brands)
}

// GetByID godoc
// @Summary      Detail Brand Berdasarkan ID
// @Description  Memuat rincian informasi brand berdasarkan ID
// @Tags         Brands
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Brand ID"
// @Success      200 {object} helpers.APIResponse{data=models.Brand} "Sukses"
// @Failure      401,404 {object} helpers.APIResponse "Not Found"
// @Router       /brands/{id} [get]
func (ctl *BrandController) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	brand, err := ctl.service.GetByID(c.Request.Context(), id)
	if err != nil {
		helpers.ResponseError(c, http.StatusNotFound, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Detail brand ditemukan", brand)
}

// Create godoc
// @Summary      Tambah Brand Baru
// @Description  Menambahkan brand outdoor baru ke katalog (Admin/Staff Only)
// @Tags         Brands
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body models.BrandCreateDTO true "Payload Brand"
// @Success      201 {object} helpers.APIResponse{data=models.Brand} "Created"
// @Failure      400,401,403 {object} helpers.APIResponse "Bad Request"
// @Router       /brands [post]
func (ctl *BrandController) Create(c *gin.Context) {
	var req models.BrandCreateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, "Data input tidak valid", err.Error())
		return
	}

	userID := c.GetInt("user_id")
	brand, err := ctl.service.Create(c.Request.Context(), req, userID)
	if err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusCreated, "Brand berhasil dibuat", brand)
}

// Update godoc
// @Summary      Update Brand (Flexible Partial)
// @Description  Memperbarui nama, deskripsi, atau status aktif brand (Admin/Staff Only)
// @Tags         Brands
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path int true "Brand ID"
// @Param        request body models.BrandUpdateDTO true "Payload Update Brand"
// @Success      200 {object} helpers.APIResponse "Sukses"
// @Failure      400,401,403 {object} helpers.APIResponse "Bad Request"
// @Router       /brands/{id} [put]
func (ctl *BrandController) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req models.BrandUpdateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, "Data input tidak valid", err.Error())
		return
	}

	userID := c.GetInt("user_id")
	if err := ctl.service.Update(c.Request.Context(), id, req, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Brand berhasil diperbarui", nil)
}

// Delete godoc
// @Summary      Soft Delete Brand
// @Description  Menonaktifkan data brand (Admin Only)
// @Tags         Brands
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Brand ID"
// @Success      200 {object} helpers.APIResponse "Sukses"
// @Failure      400,401,403 {object} helpers.APIResponse "Bad Request"
// @Router       /brands/{id} [delete]
func (ctl *BrandController) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetInt("user_id")
	if err := ctl.service.Delete(c.Request.Context(), id, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Brand berhasil dihapus (soft delete)", nil)
}

// Restore godoc
// @Summary      Restore Brand
// @Description  Memulihkan kembali brand yang berstatus soft-deleted (Admin Only)
// @Tags         Brands
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Brand ID"
// @Success      200 {object} helpers.APIResponse "Sukses"
// @Failure      400,401,403 {object} helpers.APIResponse "Bad Request"
// @Router       /brands/{id} [post]
func (ctl *BrandController) Restore(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetInt("user_id")
	if err := ctl.service.Restore(c.Request.Context(), id, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Brand berhasil dipulihkan (restore)", nil)
}

// ForceDelete godoc
// @Summary      Delete Permanent Brand (Force Delete)
// @Description  Menghapus fisik data brand secara permanen dari database (Admin Only)
// @Tags         Brands
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Brand ID"
// @Success      200 {object} helpers.APIResponse "Sukses"
// @Failure      400,401,403 {object} helpers.APIResponse "Bad Request / Foreign Key Error"
// @Router       /brands/{id}/force [delete]
func (ctl *BrandController) ForceDelete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := ctl.service.ForceDelete(c.Request.Context(), id); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Brand berhasil dihapus secara permanen dari database", nil)
}
