package controllers

import (
	"app/helpers"
	"app/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReportController struct {
	service services.ReportService
}

func NewReportController(s services.ReportService) *ReportController {
	return &ReportController{service: s}
}

// GetRevenueReport godoc
// @Summary      Laporan Keuangan & Laba Bersih
// @Description  Memuat rekap omset sewa, denda, biaya servis/maintenance alat, dan laba bersih
// @Tags         Reports
// @Security     BearerAuth
// @Produce      json
// @Param        start_date query string false "Tanggal Awal (YYYY-MM-DD)" example("2026-08-01")
// @Param        end_date query string false "Tanggal Akhir (YYYY-MM-DD)" example("2026-08-31")
// @Success      200 {object} helpers.APIResponse{data=models.RevenueReport} "Sukses"
// @Failure      401,500 {object} helpers.APIResponse "Error"
// @Router       /reports/revenue [get]
func (ctl *ReportController) GetRevenueReport(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	data, err := ctl.service.GetRevenueReport(c.Request.Context(), startDate, endDate)
	if err != nil {
		helpers.InternalServerError(c, "Gagal memuat laporan keuangan: "+err.Error())
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Laporan keuangan berhasil dimuat", data)
}

// GetRentalSummaryReport godoc
// @Summary      Laporan Ringkasan Transaksi Rental
// @Description  Memuat statistik jumlah status transaksi sewa (booked, active, completed, cancelled)
// @Tags         Reports
// @Security     BearerAuth
// @Produce      json
// @Param        start_date query string false "Tanggal Awal (YYYY-MM-DD)" example("2026-08-01")
// @Param        end_date query string false "Tanggal Akhir (YYYY-MM-DD)" example("2026-08-31")
// @Success      200 {object} helpers.APIResponse{data=models.RentalSummaryReport} "Sukses"
// @Failure      401,500 {object} helpers.APIResponse "Error"
// @Router       /reports/rentals [get]
func (ctl *ReportController) GetRentalSummaryReport(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	data, err := ctl.service.GetRentalSummaryReport(c.Request.Context(), startDate, endDate)
	if err != nil {
		helpers.InternalServerError(c, "Gagal memuat ringkasan rental: "+err.Error())
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Laporan ringkasan rental berhasil dimuat", data)
}

// GetTopProductsReport godoc
// @Summary      Laporan Produk Terlaris
// @Description  Memuat daftar peringkat alat camping yang paling sering disewa dan total omsetnya
// @Tags         Reports
// @Security     BearerAuth
// @Produce      json
// @Param        start_date query string false "Tanggal Awal (YYYY-MM-DD)" example("2026-08-01")
// @Param        end_date query string false "Tanggal Akhir (YYYY-MM-DD)" example("2026-08-31")
// @Param        limit query int false "Jumlah Data (Default: 10)" example(10)
// @Success      200 {object} helpers.APIResponse{data=[]models.TopProductReport} "Sukses"
// @Failure      401,500 {object} helpers.APIResponse "Error"
// @Router       /reports/top-products [get]
func (ctl *ReportController) GetTopProductsReport(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	data, err := ctl.service.GetTopProductsReport(c.Request.Context(), startDate, endDate, limit)
	if err != nil {
		helpers.InternalServerError(c, "Gagal memuat produk terlaris: "+err.Error())
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Laporan produk terlaris berhasil dimuat", data)
}

// GetInventoryReport godoc
// @Summary      Laporan Aset & Inventaris
// @Description  Memuat ringkasan seluruh unit fisik inventaris (tersedia, disewa, servis, hilang) beserta total valuasi aset
// @Tags         Reports
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} helpers.APIResponse{data=models.InventoryReport} "Sukses"
// @Failure      401,500 {object} helpers.APIResponse "Error"
// @Router       /reports/inventory [get]
func (ctl *ReportController) GetInventoryReport(c *gin.Context) {
	data, err := ctl.service.GetInventoryReport(c.Request.Context())
	if err != nil {
		helpers.InternalServerError(c, "Gagal memuat laporan inventaris: "+err.Error())
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Laporan inventaris berhasil dimuat", data)
}
