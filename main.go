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

	// LAYER SERVICES
	authService := services.NewAuthService(userRepo, revokedRepo)
	roleService := services.NewRoleService(roleRepo)
	userService := services.NewUserService(userRepo)

	// LAYER CONTROLLERS
	authController := controllers.NewAuthController(authService)
	roleController := controllers.NewRoleController(roleService)
	userController := controllers.NewUserController(userService)

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
		}
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8081"
	}
	log.Printf("Aplikasi POS Rental Camping siap diakses di http://localhost:%s", port)
	router.Run(":" + port)
}
