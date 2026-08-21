package main

import (
	"app/config"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Println("Peringatan: File .env tidak ditemukan, membaca dari environment sistem.")
	}
}

func main() {
	db := config.Database()
	defer db.Close()

	route := gin.Default()

	// Static route untuk file upload foto identitas
	route.Static("/storages", "./storages")

	// Panggil registrasi route dari controller/router di sini...

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8081"
	}
	route.Run(":" + port)
}
