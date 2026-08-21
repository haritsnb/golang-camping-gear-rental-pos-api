package config

import (
	"database/sql"
	"fmt"
	"log"

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
		log.Fatal("Gagal membuat instance database: ", err)
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		log.Fatal("Gagal terhubung ke database (Ping gagal): ", err)
	}

	fmt.Println("Berhasil terhubung ke database PostgreSQL!")

	Migrate(db)

	return db
}

func Migrate(db *sql.DB) {
	migrations := &migrate.FileMigrationSource{
		Dir: "database/migrations",
	}
	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		log.Fatalf("Gagal menjalankan migrasi sql-migrate: %v", err)
	}
	log.Printf("Sukses menerapkan %d file migrasi database!\n", n)
}
