package controllers

import (
	"app/helpers"
	"app/models"
	"app/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	service services.ProductService
}

func NewProductController(s services.ProductService) *ProductController {
	return &ProductController{service: s}
}

// GetAll godoc
// @Summary      Katalog Produk
// @Description  Memuat seluruh katalog produk perlengkapan camping beserta tarif harian, deposit, dan denda
// @Tags         Products
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} helpers.APIResponse{data=[]models.Product} "Sukses"
// @Failure      401,500 {object} helpers.APIResponse "Error"
// @Router       /products [get]
func (ctl *ProductController) GetAll(c *gin.Context) {
	products, err := ctl.service.GetAll(c.Request.Context())
	if err != nil {
		helpers.InternalServerError(c, err.Error())
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Daftar produk berhasil dimuat", products)
}

// GetByID godoc
// @Summary      Detail Produk Berdasarkan ID
// @Description  Memuat rincian informasi dan tarif produk sewa
// @Tags         Products
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Product ID"
// @Success      200 {object} helpers.APIResponse{data=models.Product} "Sukses"
// @Failure      401,404 {object} helpers.APIResponse "Not Found"
// @Router       /products/{id} [get]
func (ctl *ProductController) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	prod, err := ctl.service.GetByID(c.Request.Context(), id)
	if err != nil {
		helpers.ResponseError(c, http.StatusNotFound, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Detail produk ditemukan", prod)
}

// Create godoc
// @Summary      Tambah Produk Baru
// @Description  Mendaftarkan produk sewa baru ke katalog sistem (Admin/Staff Only)
// @Tags         Products
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body models.ProductCreateDTO true "Payload Produk"
// @Success      201 {object} helpers.APIResponse{data=models.Product} "Created"
// @Failure      400,401,403 {object} helpers.APIResponse "Bad Request"
// @Router       /products [post]
func (ctl *ProductController) Create(c *gin.Context) {
	var req models.ProductCreateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, "Data input tidak valid", err.Error())
		return
	}

	userID := c.GetInt("user_id")
	prod, err := ctl.service.Create(c.Request.Context(), req, userID)
	if err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusCreated, "Produk berhasil dibuat", prod)
}

// Update godoc
// @Summary      Update Produk (Flexible Partial)
// @Description  Memperbarui harga sewa, denda, deposit, atau kategori produk secara fleksibel (Admin/Staff Only)
// @Tags         Products
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path int true "Product ID"
// @Param        request body models.ProductUpdateDTO true "Payload Update Produk"
// @Success      200 {object} helpers.APIResponse "Sukses"
// @Failure      400,401,403 {object} helpers.APIResponse "Bad Request"
// @Router       /products/{id} [put]
func (ctl *ProductController) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req models.ProductUpdateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, "Data input tidak valid", err.Error())
		return
	}

	userID := c.GetInt("user_id")
	if err := ctl.service.Update(c.Request.Context(), id, req, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Produk berhasil diperbarui", nil)
}

// Delete godoc
// @Summary      Soft Delete Produk
// @Description  Menonaktifkan produk dari katalog sewa (Admin Only)
// @Tags         Products
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Product ID"
// @Success      200 {object} helpers.APIResponse "Sukses"
// @Failure      400,401,403 {object} helpers.APIResponse "Bad Request"
// @Router       /products/{id} [delete]
func (ctl *ProductController) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetInt("user_id")
	if err := ctl.service.Delete(c.Request.Context(), id, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Produk berhasil dihapus (soft delete)", nil)
}

// Restore godoc
// @Summary      Restore Produk
// @Description  Memulihkan kembali produk yang berstatus soft-deleted (Admin Only)
// @Tags         Products
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Product ID"
// @Success      200 {object} helpers.APIResponse "Sukses"
// @Failure      400,401,403 {object} helpers.APIResponse "Bad Request"
// @Router       /products/{id} [post]
func (ctl *ProductController) Restore(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetInt("user_id")
	if err := ctl.service.Restore(c.Request.Context(), id, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Produk berhasil dipulihkan (restore)", nil)
}

// ForceDelete godoc
// @Summary      Delete Permanent Produk (Force Delete)
// @Description  Menghapus fisik data produk secara permanen dari database (Admin Only)
// @Tags         Products
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Product ID"
// @Success      200 {object} helpers.APIResponse "Sukses"
// @Failure      400,401,403 {object} helpers.APIResponse "Bad Request / Foreign Key Error"
// @Router       /products/{id}/force [delete]
func (ctl *ProductController) ForceDelete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := ctl.service.ForceDelete(c.Request.Context(), id); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Produk berhasil dihapus secara permanen dari database", nil)
}
