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

// GetAll godoc
// @Summary      Daftar Semua Pelanggan
// @Description  Memuat seluruh daftar pelanggan terdaftar beserta status KYC dan blacklist
// @Tags         Customers
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} helpers.APIResponse{data=[]models.Customer} "Sukses"
// @Failure      401,500 {object} helpers.APIResponse "Error"
// @Router       /customers [get]
func (ctl *CustomerController) GetAll(c *gin.Context) {
	customers, err := ctl.service.GetAll(c.Request.Context())
	if err != nil {
		helpers.InternalServerError(c, err.Error())
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Daftar pelanggan berhasil dimuat", customers)
}

// GetByID godoc
// @Summary      Detail Pelanggan Berdasarkan ID
// @Description  Memuat data lengkap pelanggan dan URL foto identitasnya
// @Tags         Customers
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Customer ID"
// @Success      200 {object} helpers.APIResponse{data=models.Customer} "Sukses"
// @Failure      401,404 {object} helpers.APIResponse "Not Found"
// @Router       /customers/{id} [get]
func (ctl *CustomerController) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	cust, err := ctl.service.GetByID(c.Request.Context(), id)
	if err != nil {
		helpers.ResponseError(c, http.StatusNotFound, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Detail pelanggan ditemukan", cust)
}

// GetDeleted godoc
// @Summary      Daftar Pelanggan Terhapus (Trash)
// @Tags         Customers
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} helpers.APIResponse{data=[]models.Customer} "Sukses"
// @Router       /customers/deleted [get]
func (ctl *CustomerController) GetDeleted(c *gin.Context) {
	customers, err := ctl.service.GetDeleted(c.Request.Context())
	if err != nil {
		helpers.InternalServerError(c, err.Error())
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Daftar pelanggan terhapus", customers)
}

// Create godoc
// @Summary      Registrasi Pelanggan & Upload KYC
// @Description  Mendaftarkan data pelanggan baru dan mengunggah foto kartu identitas (KTP/SIM/Paspor)
// @Tags         Customers
// @Security     BearerAuth
// @Accept       multipart/form-data
// @Produce      json
// @Param        name formData string true "Nama Lengkap"
// @Param        identity_type formData string true "Tipe Identitas (KTP/SIM/Paspor)"
// @Param        identity_number formData string true "Nomor Unik Identitas"
// @Param        phone formData string true "Nomor HP / WA Aktif"
// @Param        address formData string true "Alamat Domisili"
// @Param        emergency_contact formData string false "Kontak Darurat"
// @Param        email formData string false "Email Pelanggan"
// @Param        notes formData string false "Catatan Tambahan"
// @Param        identity_photo formData file false "Foto Fisik Kartu Identitas"
// @Success      201 {object} helpers.APIResponse{data=models.Customer} "Created"
// @Failure      400,401,500 {object} helpers.APIResponse "Bad Request / Upload Error"
// @Router       /customers [post]
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

// Update godoc
// @Summary      Update Pelanggan (Flexible Partial + Opsional Foto)
// @Description  Memperbarui biodata pelanggan dan opsional mengganti file foto kartu identitas
// @Tags         Customers
// @Security     BearerAuth
// @Accept       multipart/form-data
// @Accept       json
// @Produce      json
// @Param        id path int true "Customer ID"
// @Param        name formData string false "Nama Lengkap"
// @Param        identity_type formData string false "Tipe Identitas"
// @Param        identity_number formData string false "Nomor Identitas"
// @Param        phone formData string false "Nomor HP"
// @Param        address formData string false "Alamat"
// @Param        emergency_contact formData string false "Kontak Darurat"
// @Param        email formData string false "Email"
// @Param        notes formData string false "Catatan"
// @Param        identity_photo formData file false "Foto Identitas Baru"
// @Success      200 {object} helpers.APIResponse "Sukses"
// @Failure      400,401,500 {object} helpers.APIResponse "Bad Request"
// @Router       /customers/{id} [put]
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

// Blacklist godoc
// @Summary      Ubah Status Blacklist Pelanggan
// @Description  Memasukkan atau menghapus pelanggan dari daftar blacklist (Admin Only)
// @Tags         Customers
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path int true "Customer ID"
// @Param        request body object{is_blacklisted=bool} true "Payload Blacklist"
// @Success      200 {object} helpers.APIResponse "Sukses"
// @Failure      400,401,403 {object} helpers.APIResponse "Bad Request"
// @Router       /customers/{id}/blacklist [patch]
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

// Delete godoc
// @Summary      Soft Delete Pelanggan
// @Description  Menonaktifkan data pelanggan (Admin Only)
// @Tags         Customers
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Customer ID"
// @Success      200 {object} helpers.APIResponse "Sukses"
// @Failure      400,401,403 {object} helpers.APIResponse "Bad Request"
// @Router       /customers/{id} [delete]
func (ctl *CustomerController) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetInt("user_id")

	if err := ctl.service.Delete(c.Request.Context(), id, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Customer berhasil dihapus (soft delete)", nil)
}

// Restore godoc
// @Summary      Restore Pelanggan
// @Description  Memulihkan kembali data pelanggan yang berstatus soft-deleted (Admin Only)
// @Tags         Customers
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Customer ID"
// @Success      200 {object} helpers.APIResponse "Sukses"
// @Failure      400,401,403 {object} helpers.APIResponse "Bad Request"
// @Router       /customers/{id} [post]
func (ctl *CustomerController) Restore(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetInt("user_id")

	if err := ctl.service.Restore(c.Request.Context(), id, userID); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Customer berhasil dipulihkan (restore)", nil)
}

// ForceDelete godoc
// @Summary      Delete Permanent Pelanggan (Force Delete)
// @Description  Menghapus fisik data pelanggan dari database (Admin Only)
// @Tags         Customers
// @Security     BearerAuth
// @Produce      json
// @Param        id path int true "Customer ID"
// @Success      200 {object} helpers.APIResponse "Sukses"
// @Failure      400,401,403 {object} helpers.APIResponse "Bad Request / Foreign Key Error"
// @Router       /customers/{id}/force [delete]
func (ctl *CustomerController) ForceDelete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := ctl.service.ForceDelete(c.Request.Context(), id); err != nil {
		helpers.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Customer berhasil dihapus secara permanen dari database", nil)
}
