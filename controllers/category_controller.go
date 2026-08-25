package controllers

import (
	"app/helpers"
	"app/models"
	"app/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoryController struct {
	service services.CategoryService
}

func NewCategoryController(s services.CategoryService) *CategoryController {
	return &CategoryController{service: s}
}

// GetAll godoc
// @Summary      Daftar Semua Kategori
// @Description  Mengambil seluruh data kategori produk alat camping
// @Tags         Categories
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} helpers.APIResponse{data=[]models.Category} "Sukses"
// @Failure      401,500 {object} helpers.APIResponse "Error"
// @Router       /categories [get]
func (ctl *CategoryController) GetAll(c *gin.Context) {
	categories, err := ctl.service.GetAll(c.Request.Context())
	if err != nil {
		helpers.InternalServerError(c, err.Error())
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Daftar kategori berhasil dimuat", categories)
}

// GetByID godoc
// @Summary      Detail Kategori Berdasarkan ID
// @Description  Mengambil rincian informasi kategori berdasarkan ID
// @Tags         Categories
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Category ID"
// @Success      200 {object} helpers.APIResponse{data=models.Category} "Sukses"
// @Failure      401,404 {object} helpers.APIResponse "Not Found"
// @Router       /categories/{id} [get]
func (ctl *CategoryController) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	cat, err := ctl.service.GetByID(c.Request.Context(), id)
	if err != nil {
		helpers.ResponseError(c, http.StatusNotFound, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Detail kategori ditemukan", cat)
}

// GetDeleted godoc
// @Summary      Daftar Kategori Terhapus (Trash)
// @Tags         Categories
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} helpers.APIResponse{data=[]models.Category} "Sukses"
// @Router       /categories/deleted [get]
func (ctl *CategoryController) GetDeleted(c *gin.Context) {
	cats, err := ctl.service.GetDeleted(c.Request.Context())
	if err != nil {
		helpers.InternalServerError(c, err.Error())
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Daftar kategori terhapus", cats)
}

// Create godoc
// @Summary      Tambah Kategori Baru
// @Description  Menambahkan kategori perlengkapan camping baru (Admin/Staff Only)
// @Tags         Categories
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body models.CategoryCreateDTO true "Payload Kategori"
// @Success      201 {object} helpers.APIResponse{data=models.Category} "Created"
// @Failure      400,401,403 {object} helpers.APIResponse "Bad Request"
// @Router       /categories [post]
func (ctl *CategoryController) Create(c *gin.Context) {
	var req models.CategoryCreateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, "Data input tidak valid", err.Error())
		return
	}

	userID := c.GetInt("user_id")
	cat, err := ctl.service.Create(c.Request.Context(), req, userID)
	if err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusCreated, "Kategori berhasil dibuat", cat)
}

// Update godoc
// @Summary      Update Kategori (Flexible Partial)
// @Description  Memperbarui nama atau deskripsi kategori (Admin/Staff Only)
// @Tags         Categories
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path int true "Category ID"
// @Param        request body models.CategoryUpdateDTO true "Payload Update Kategori"
// @Success      200 {object} helpers.APIResponse "Sukses"
// @Failure      400,401,403 {object} helpers.APIResponse "Bad Request"
// @Router       /categories/{id} [put]
func (ctl *CategoryController) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req models.CategoryUpdateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, "Data input tidak valid", err.Error())
		return
	}

	userID := c.GetInt("user_id")
	if err := ctl.service.Update(c.Request.Context(), id, req, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Kategori berhasil diperbarui", nil)
}

// Delete godoc
// @Summary      Soft Delete Kategori
// @Description  Menonaktifkan kategori produk (Admin Only)
// @Tags         Categories
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Category ID"
// @Success      200 {object} helpers.APIResponse "Sukses"
// @Failure      400,401,403 {object} helpers.APIResponse "Bad Request"
// @Router       /categories/{id} [delete]
func (ctl *CategoryController) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetInt("user_id")
	if err := ctl.service.Delete(c.Request.Context(), id, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Kategori berhasil dihapus (soft delete)", nil)
}

// Restore godoc
// @Summary      Restore Kategori
// @Description  Memulihkan kembali kategori yang berstatus soft-deleted (Admin Only)
// @Tags         Categories
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Category ID"
// @Success      200 {object} helpers.APIResponse "Sukses"
// @Failure      400,401,403 {object} helpers.APIResponse "Bad Request"
// @Router       /categories/{id} [post]
func (ctl *CategoryController) Restore(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetInt("user_id")
	if err := ctl.service.Restore(c.Request.Context(), id, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Kategori berhasil dipulihkan (restore)", nil)
}

// ForceDelete godoc
// @Summary      Delete Permanent Kategori (Force Delete)
// @Description  Menghapus fisik data kategori secara permanen dari database (Admin Only)
// @Tags         Categories
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Category ID"
// @Success      200 {object} helpers.APIResponse "Sukses"
// @Failure      400,401,403 {object} helpers.APIResponse "Bad Request / Foreign Key Error"
// @Router       /categories/{id}/force [delete]
func (ctl *CategoryController) ForceDelete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := ctl.service.ForceDelete(c.Request.Context(), id); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Kategori berhasil dihapus secara permanen dari database", nil)
}
