package config

import (
	"app/database/seeders"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
)

func Database() *sql.DB {
	host := RequireEnv("DB_HOST")
	port := RequireEnv("DB_PORT")
	user := RequireEnv("DB_USER")
	password := RequireEnv("DB_PASSWORD")
	dbname := RequireEnv("DB_NAME")
	sslmode := RequireEnv("DB_SSL")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Gagal membuat instance koneksi database: ", err)
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		log.Fatal("Gagal terhubung ke database PostgreSQL (Ping gagal): ", err)
	}

	fmt.Println("==================================================")
	fmt.Println("Berhasil terhubung ke database PostgreSQL!")

	// Migrasi DDL (Tabel, Relasi, Index)
	Migrate(db)

	// Seeder SQL untuk Master Data Kategori, Brand, Produk, dsb.
	Seed(db)

	// Seeder Go untuk Super Admin (Enkripsi Bcrypt Realtime)
	if err := seeders.SeedSuperAdmin(db); err != nil {
		log.Fatalf("Gagal menjalankan File Seeder Users: %v", err)
	}

	fmt.Println("==================================================")

	return db
}

func Migrate(db *sql.DB) {
	migrate.SetTable("gorp_migrations")
	migrations := &migrate.FileMigrationSource{
		Dir: "database/migrations",
	}

	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		log.Fatalf("Gagal menjalankan migrasi skema: %v", err)
	}
	if n > 0 {
		log.Printf("MIGRASI SUKSES: Menerapkan %d file migrasi skema database.\n", n)
	}
}

func Seed(db *sql.DB) {
	if _, err := os.Stat("database/seeders"); os.IsNotExist(err) {
		return
	}

	migrate.SetTable("gorp_seeders")
	seedersSource := &migrate.FileMigrationSource{
		Dir: "database/seeders",
	}

	n, err := migrate.Exec(db, "postgres", seedersSource, migrate.Up)
	if err != nil {
		log.Fatalf("Gagal menjalankan seeder SQL: %v", err)
	}
	if n > 0 {
		log.Printf("SEEDER SQL SUKSES: Menerapkan %d file seeder master data.\n", n)
	}
}
