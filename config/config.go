package config

import (
	"log"
	"os"
	"strings"
)

func RequireEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("CONFIG ERROR: Variabel '%s' TIDAK DITEMUKAN di file .env!", key)
	}

	value = strings.TrimSpace(value)
	if value == "" {
		log.Fatalf("CONFIG ERROR: Variabel '%s' ADA di file .env, tetapi nilainya KOSONG!", key)
	}

	return value
}
