package main

import (
	"app/config"
	"app/controllers"
	"app/middlewares"
	"app/repositories"
	"app/services"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("Info: File .env tidak ditemukan, membaca konfigurasi environment sistem.")
	}
}

func main() {
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

	// LAYER SERVICES
	authService := services.NewAuthService(userRepo, revokedRepo)
	roleService := services.NewRoleService(roleRepo)
	userService := services.NewUserService(userRepo)
	custService := services.NewCustomerService(custRepo)
	catService := services.NewCategoryService(catRepo)
	brandService := services.NewBrandService(brandRepo)
	prodService := services.NewProductService(prodRepo)
	unitService := services.NewProductUnitService(unitRepo)

	// LAYER CONTROLLERS
	authController := controllers.NewAuthController(authService)
	roleController := controllers.NewRoleController(roleService)
	userController := controllers.NewUserController(userService)
	custController := controllers.NewCustomerController(custService)
	catController := controllers.NewCategoryController(catService)
	brandController := controllers.NewBrandController(brandService)
	prodController := controllers.NewProductController(prodService)
	unitController := controllers.NewProductUnitController(unitService)

	router := gin.Default()

	// Static route untuk melayani file upload foto KYC
	router.Static("/storages", "./storages")

	api := router.Group("/api/v1")
	{
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
				unitGroup.GET("/scan/:unit_code", unitController.ScanUnitCode)
				unitGroup.POST("", middlewares.RoleGuard("admin", "staff"), unitController.Create)
				unitGroup.POST("/:id", middlewares.RoleGuard("admin"), unitController.Restore)
				unitGroup.PUT("/:id", middlewares.RoleGuard("admin", "staff"), unitController.Update)
				unitGroup.DELETE("/:id", middlewares.RoleGuard("admin"), unitController.Delete)
				unitGroup.DELETE("/:id/force", middlewares.RoleGuard("admin"), unitController.ForceDelete)
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
