package main

import (
	"app/config"
	"app/controllers"
	"app/docs"
	_ "app/docs"
	"app/middlewares"
	"app/repositories"
	"app/services"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Outdoor Gear Rental POS API
// @version         1.0
// @description     Backend RESTful API Sistem Kasir & Manajemen Rental Alat Camping menggunakan Golang + Gin + PostgreSQL + JWT Blacklist + KYC Upload.
// @termsOfService  http://swagger.io/terms/

// @contact.name    Harits Nala Barrun
// @contact.email   developer.haritsnb@gmail.com

// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath        /api

// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization
// @description     Masukkan token dengan format: Bearer <token_jwt>

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("Info: File .env tidak ditemukan, membaca konfigurasi environment sistem.")
	}
}

func main() {
	// Set host dinamis agar bisa berjalan di localhost maupun server produksi
	docs.SwaggerInfo.Host = ""
	docs.SwaggerInfo.BasePath = "/api"

	// DATABASE
	db := config.Database()
	defer db.Close()

	// LAYER REPOSITORIES
	revokedRepo := repositories.NewRevokedTokenRepository(db)
	roleRepo := repositories.NewRoleRepository(db)
	userRepo := repositories.NewUserRepository(db)
	custRepo := repositories.NewCustomerRepository(db)
	catRepo := repositories.NewCategoryRepository(db)
	brandRepo := repositories.NewBrandRepository(db)
	prodRepo := repositories.NewProductRepository(db)
	unitRepo := repositories.NewProductUnitRepository(db)
	rentalRepo := repositories.NewRentalRepository(db)
	rentalItemRepo := repositories.NewRentalItemRepository(db)
	payRepo := repositories.NewPaymentRepository(db)
	maintRepo := repositories.NewMaintenanceRepository(db)
	reportRepo := repositories.NewReportRepository(db)

	// LAYER SERVICES
	authService := services.NewAuthService(userRepo, revokedRepo)
	roleService := services.NewRoleService(roleRepo)
	userService := services.NewUserService(userRepo)
	custService := services.NewCustomerService(custRepo)
	catService := services.NewCategoryService(catRepo)
	brandService := services.NewBrandService(brandRepo)
	prodService := services.NewProductService(prodRepo)
	unitService := services.NewProductUnitService(unitRepo)
	payService := services.NewPaymentService(payRepo)
	maintService := services.NewMaintenanceService(db, maintRepo, unitRepo)
	rentalService := services.NewRentalService(db, rentalRepo, rentalItemRepo, unitRepo, prodRepo, custRepo, payRepo, maintRepo)
	reportService := services.NewReportService(reportRepo)

	// LAYER CONTROLLERS
	authController := controllers.NewAuthController(authService)
	roleController := controllers.NewRoleController(roleService)
	userController := controllers.NewUserController(userService)
	custController := controllers.NewCustomerController(custService)
	catController := controllers.NewCategoryController(catService)
	brandController := controllers.NewBrandController(brandService)
	prodController := controllers.NewProductController(prodService)
	unitController := controllers.NewProductUnitController(unitService)
	rentalController := controllers.NewRentalController(rentalService)
	payController := controllers.NewPaymentController(payService)
	maintController := controllers.NewMaintenanceController(maintService)
	reportController := controllers.NewReportController(reportService)

	router := gin.Default()

	// Static route untuk melayani file di folder test:
	router.Static("/test", "./test")

	// Static route untuk melayani file upload foto KYC
	router.Static("/storages", "./storages")

	api := router.Group("/api")
	{
		api.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		authGroup := api.Group("/auth")
		{
			authGroup.POST("/login", authController.Login)
			authGroup.POST("/logout", middlewares.AuthMiddleware(revokedRepo), authController.Logout)
			authGroup.GET("/me", middlewares.AuthMiddleware(revokedRepo), authController.Me)
		}

		protected := api.Group("", middlewares.AuthMiddleware(revokedRepo))
		{
			rolesGroup := protected.Group("/roles", middlewares.RoleGuard("admin"))
			{
				rolesGroup.GET("", roleController.GetAll)
				rolesGroup.GET("/:id", roleController.GetByID)
				rolesGroup.GET("/deleted", roleController.GetDeleted)
				rolesGroup.POST("", roleController.Create)
				rolesGroup.POST("/:id", roleController.Restore)
				rolesGroup.PUT("/:id", roleController.Update)
				rolesGroup.DELETE("/:id", roleController.Delete)
				rolesGroup.DELETE("/:id/force", roleController.ForceDelete)
			}

			usersGroup := protected.Group("/users", middlewares.RoleGuard("admin"))
			{
				usersGroup.GET("", userController.GetAll)
				usersGroup.GET("/:id", userController.GetByID)
				usersGroup.GET("/deleted", userController.GetDeleted)
				usersGroup.POST("", userController.Create)
				usersGroup.POST("/:id", userController.Restore)
				usersGroup.PUT("/:id", userController.Update)
				usersGroup.DELETE("/:id", userController.Delete)
				usersGroup.DELETE("/:id/force", userController.ForceDelete)
			}

			custGroup := protected.Group("/customers")
			{
				custGroup.GET("", custController.GetAll)
				custGroup.GET("/:id", custController.GetByID)
				custGroup.GET("/deleted", custController.GetDeleted)
				custGroup.POST("", custController.Create)
				custGroup.POST("/:id", custController.Restore)
				custGroup.PUT("/:id", custController.Update)
				custGroup.PATCH("/:id/blacklist", middlewares.RoleGuard("admin"), custController.Blacklist)
				custGroup.DELETE("/:id", middlewares.RoleGuard("admin"), custController.Delete)
				custGroup.DELETE("/:id/force", middlewares.RoleGuard("admin"), custController.ForceDelete)
			}

			catGroup := protected.Group("/categories")
			{
				catGroup.GET("", catController.GetAll)
				catGroup.GET("/:id", catController.GetByID)
				catGroup.GET("/deleted", middlewares.RoleGuard("admin", "staff"), catController.GetDeleted)
				catGroup.POST("", middlewares.RoleGuard("admin", "staff"), catController.Create)
				catGroup.POST("/:id", middlewares.RoleGuard("admin"), catController.Restore)
				catGroup.PUT("/:id", middlewares.RoleGuard("admin", "staff"), catController.Update)
				catGroup.DELETE("/:id", middlewares.RoleGuard("admin"), catController.Delete)
				catGroup.DELETE("/:id/force", middlewares.RoleGuard("admin"), catController.ForceDelete)
			}

			brandGroup := protected.Group("/brands")
			{
				brandGroup.GET("", brandController.GetAll)
				brandGroup.GET("/:id", brandController.GetByID)
				brandGroup.GET("/deleted", middlewares.RoleGuard("admin", "staff"), brandController.GetDeleted)
				brandGroup.POST("", middlewares.RoleGuard("admin", "staff"), brandController.Create)
				brandGroup.POST("/:id", middlewares.RoleGuard("admin"), brandController.Restore)
				brandGroup.PUT("/:id", middlewares.RoleGuard("admin", "staff"), brandController.Update)
				brandGroup.DELETE("/:id", middlewares.RoleGuard("admin"), brandController.Delete)
				brandGroup.DELETE("/:id/force", middlewares.RoleGuard("admin"), brandController.ForceDelete)
			}

			prodGroup := protected.Group("/products")
			{
				prodGroup.GET("", prodController.GetAll)
				prodGroup.GET("/:id", prodController.GetByID)
				prodGroup.GET("/deleted", middlewares.RoleGuard("admin", "staff"), prodController.GetDeleted)
				prodGroup.POST("", middlewares.RoleGuard("admin", "staff"), prodController.Create)
				prodGroup.POST("/:id", middlewares.RoleGuard("admin"), prodController.Restore)
				prodGroup.PUT("/:id", middlewares.RoleGuard("admin", "staff"), prodController.Update)
				prodGroup.DELETE("/:id", middlewares.RoleGuard("admin"), prodController.Delete)
				prodGroup.DELETE("/:id/force", middlewares.RoleGuard("admin"), prodController.ForceDelete)
			}

			unitGroup := protected.Group("/product-units")
			{
				unitGroup.GET("", unitController.GetAll)
				unitGroup.GET("/:id", unitController.GetByID)
				unitGroup.GET("/deleted", middlewares.RoleGuard("admin", "staff"), unitController.GetDeleted)
				unitGroup.GET("/scan/:unit_code", unitController.ScanUnitCode)
				unitGroup.POST("", middlewares.RoleGuard("admin", "staff"), unitController.Create)
				unitGroup.POST("/:id", middlewares.RoleGuard("admin"), unitController.Restore)
				unitGroup.PUT("/:id", middlewares.RoleGuard("admin", "staff"), unitController.Update)
				unitGroup.DELETE("/:id", middlewares.RoleGuard("admin"), unitController.Delete)
				unitGroup.DELETE("/:id/force", middlewares.RoleGuard("admin"), unitController.ForceDelete)
			}

			rentalGroup := protected.Group("/rentals")
			{
				rentalGroup.GET("", rentalController.GetAll)
				rentalGroup.GET("/:id", rentalController.GetByID)
				rentalGroup.GET("/deleted", middlewares.RoleGuard("admin", "cashier"), rentalController.GetDeleted)
				rentalGroup.POST("/booking", rentalController.Booking)
				rentalGroup.POST("/:id/handover", rentalController.Handover)
				rentalGroup.POST("/:id/return", rentalController.Return)
				rentalGroup.POST("/:id/settlement", rentalController.Settlement)
				rentalGroup.POST("/:id/cancel", rentalController.Cancel)
				rentalGroup.POST("/:id", middlewares.RoleGuard("admin"), rentalController.Restore)
				rentalGroup.PUT("/:id", middlewares.RoleGuard("admin"), rentalController.Update)
				rentalGroup.DELETE("/:id", middlewares.RoleGuard("admin"), rentalController.Delete)
				rentalGroup.DELETE("/:id/force", middlewares.RoleGuard("admin"), rentalController.ForceDelete)
			}

			payGroup := protected.Group("/payments")
			{
				payGroup.GET("/rental/:rental_id", payController.GetByRentalID)
				payGroup.GET("/:id", payController.GetByID)
				payGroup.GET("/deleted", middlewares.RoleGuard("admin"), payController.GetDeleted)
				payGroup.POST("", payController.Create)
				payGroup.POST("/:id", middlewares.RoleGuard("admin"), payController.Restore)
				payGroup.PUT("/:id", middlewares.RoleGuard("admin"), payController.Update)
				payGroup.DELETE("/:id", middlewares.RoleGuard("admin"), payController.Delete)
				payGroup.DELETE("/:id/force", middlewares.RoleGuard("admin"), payController.ForceDelete)
			}

			maintGroup := protected.Group("/maintenance")
			{
				maintGroup.GET("", maintController.GetAll)
				maintGroup.GET("/:id", maintController.GetByID)
				maintGroup.GET("/deleted", middlewares.RoleGuard("admin", "staff"), maintController.GetDeleted)
				maintGroup.POST("", middlewares.RoleGuard("admin", "staff"), maintController.Create)
				maintGroup.POST("/:id", middlewares.RoleGuard("admin"), maintController.Restore)
				maintGroup.PUT("/:id", middlewares.RoleGuard("admin", "staff"), maintController.Update)
				maintGroup.PATCH("/:id/complete", middlewares.RoleGuard("admin", "staff"), maintController.Complete)
				maintGroup.DELETE("/:id", middlewares.RoleGuard("admin"), maintController.Delete)
				maintGroup.DELETE("/:id/force", middlewares.RoleGuard("admin"), maintController.ForceDelete)
			}

			reportGroup := protected.Group("/reports", middlewares.RoleGuard("admin", "staff"))
			{
				reportGroup.GET("/revenue", reportController.GetRevenueReport)
				reportGroup.GET("/rentals", reportController.GetRentalSummaryReport)
				reportGroup.GET("/top-products", reportController.GetTopProductsReport)
				reportGroup.GET("/inventory", reportController.GetInventoryReport)
			}
		}
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8081"
	}
	log.Printf("Aplikasi POS Rental Camping siap diakses di http://localhost:%s", port)
	router.Run(":" + port)
}
