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

// GetAll godoc
// @Summary      Daftar Semua Unit Fisik
// @Description  Memuat seluruh unit fisik produk yang memiliki kode barcode/QR dan serial number
// @Tags         Product Units
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} helpers.APIResponse{data=[]models.ProductUnit} "Sukses"
// @Failure      401,500 {object} helpers.APIResponse "Error"
// @Router       /product-units [get]
func (ctl *ProductUnitController) GetAll(c *gin.Context) {
	units, err := ctl.service.GetAll(c.Request.Context())
	if err != nil {
		helpers.InternalServerError(c, err.Error())
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Daftar unit fisik produk berhasil dimuat", units)
}

// GetByID godoc
// @Summary      Detail Unit Fisik Berdasarkan ID
// @Description  Memuat detail unit fisik, riwayat kondisi, dan status ketersediaannya
// @Tags         Product Units
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Product Unit ID"
// @Success      200 {object} helpers.APIResponse{data=models.ProductUnit} "Sukses"
// @Failure      401,404 {object} helpers.APIResponse "Not Found"
// @Router       /product-units/{id} [get]
func (ctl *ProductUnitController) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	unit, err := ctl.service.GetByID(c.Request.Context(), id)
	if err != nil {
		helpers.ResponseError(c, http.StatusNotFound, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Detail unit fisik ditemukan", unit)
}

// GetDeleted godoc
// @Summary      Daftar Unit Fisik Terhapus (Trash)
// @Tags         Product Units
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} helpers.APIResponse{data=[]models.ProductUnit} "Sukses"
// @Router       /product-units/deleted [get]
func (ctl *ProductUnitController) GetDeleted(c *gin.Context) {
	units, err := ctl.service.GetDeleted(c.Request.Context())
	if err != nil {
		helpers.InternalServerError(c, err.Error())
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Daftar unit fisik terhapus", units)
}

// ScanUnitCode godoc
// @Summary      Scan Barcode / QR Code Unit Fisik
// @Description  Memindai kode barcode unik unit fisik saat transaksi checkout atau serah terima alat
// @Tags         Product Units
// @Security     BearerAuth
// @Produce      json
// @Param        unit_code path string true "Unit Code Barcode (Contoh: TND-EIG-001)"
// @Success      200 {object} helpers.APIResponse{data=models.ProductUnit} "Sukses"
// @Failure      401,404 {object} helpers.APIResponse "Not Found"
// @Router       /product-units/scan/{unit_code} [get]
func (ctl *ProductUnitController) ScanUnitCode(c *gin.Context) {
	unitCode := c.Param("unit_code")
	unit, err := ctl.service.GetByUnitCode(c.Request.Context(), unitCode)
	if err != nil {
		helpers.ResponseError(c, http.StatusNotFound, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Scan unit fisik berhasil", unit)
}

// Create godoc
// @Summary      Tambah Unit Fisik Baru
// @Description  Mendaftarkan unit fisik serialized dengan kode unik dan nomor seri (Admin/Staff Only)
// @Tags         Product Units
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body models.ProductUnitCreateDTO true "Payload Unit Fisik"
// @Success      201 {object} helpers.APIResponse{data=models.ProductUnit} "Created"
// @Failure      400,401,403 {object} helpers.APIResponse "Bad Request"
// @Router       /product-units [post]
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

// Update godoc
// @Summary      Update Unit Fisik (Flexible Partial)
// @Description  Memperbarui nomor seri, kondisi, status, atau catatan unit fisik (Admin/Staff Only)
// @Tags         Product Units
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path int true "Product Unit ID"
// @Param        request body models.ProductUnitUpdateDTO true "Payload Update Unit Fisik"
// @Success      200 {object} helpers.APIResponse "Sukses"
// @Failure      400,401,403 {object} helpers.APIResponse "Bad Request"
// @Router       /product-units/{id} [put]
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

// Delete godoc
// @Summary      Soft Delete Unit Fisik
// @Description  Menonaktifkan unit fisik dari daftar inventaris (Admin Only)
// @Tags         Product Units
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Product Unit ID"
// @Success      200 {object} helpers.APIResponse "Sukses"
// @Failure      400,401,403 {object} helpers.APIResponse "Bad Request"
// @Router       /product-units/{id} [delete]
func (ctl *ProductUnitController) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetInt("user_id")
	if err := ctl.service.Delete(c.Request.Context(), id, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Unit fisik berhasil dihapus (soft delete)", nil)
}

// Restore godoc
// @Summary      Restore Unit Fisik
// @Description  Memulihkan kembali unit fisik yang berstatus soft-deleted (Admin Only)
// @Tags         Product Units
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Product Unit ID"
// @Success      200 {object} helpers.APIResponse "Sukses"
// @Failure      400,401,403 {object} helpers.APIResponse "Bad Request"
// @Router       /product-units/{id} [post]
func (ctl *ProductUnitController) Restore(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetInt("user_id")
	if err := ctl.service.Restore(c.Request.Context(), id, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Unit fisik berhasil dipulihkan (restore)", nil)
}

// ForceDelete godoc
// @Summary      Delete Permanent Unit Fisik (Force Delete)
// @Description  Menghapus fisik data unit dari database (Admin Only)
// @Tags         Product Units
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Product Unit ID"
// @Success      200 {object} helpers.APIResponse "Sukses"
// @Failure      400,401,403 {object} helpers.APIResponse "Bad Request / Foreign Key Error"
// @Router       /product-units/{id}/force [delete]
func (ctl *ProductUnitController) ForceDelete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := ctl.service.ForceDelete(c.Request.Context(), id); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Unit fisik berhasil dihapus secara permanen dari database", nil)
}
