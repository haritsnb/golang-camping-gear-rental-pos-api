package seeders

import (
	"database/sql"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func SeedSuperAdmin(db *sql.DB) error {
	const (
		adminUsername = "admin"
		adminPassword = "password123"
		adminRoleName = "admin"
	)

	// Cek role di database
	var roleID int
	err := db.QueryRow("SELECT id FROM roles WHERE name = $1 AND deleted_at IS NULL", adminRoleName).Scan(&roleID)
	if err != nil {
		// Jika role belum ada, buat role admin
		err = db.QueryRow(`
			INSERT INTO roles (name, description, created_at, modified_at)
			VALUES ($1, 'Super Administrator dengan hak akses penuh', NOW(), NOW())
			RETURNING id`, adminRoleName).Scan(&roleID)
		if err != nil {
			return fmt.Errorf("gagal membuat role default admin: %v", err)
		}
	}

	// Generate Bcrypt Hash langsung via Go
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("gagal mengenkripsi password: %v", err)
	}

	// Cek user di database
	var existingUserID int
	err = db.QueryRow("SELECT id FROM users WHERE username = $1", adminUsername).Scan(&existingUserID)

	if err == sql.ErrNoRows {
		// User belum ada -> Lakukan INSERT
		var newUserID int
		queryInsert := `
			INSERT INTO users (role_id, name, username, password_hash, phone, is_active, created_at, modified_at)
			VALUES ($1, 'Super Administrator', $2, $3, '081234567890', TRUE, NOW(), NOW())
			RETURNING id`
		err = db.QueryRow(queryInsert, roleID, adminUsername, string(hashedPassword)).Scan(&newUserID)
		if err != nil {
			return fmt.Errorf("gagal insert super admin: %v", err)
		}

		// Update self-audit trail (created_by & modified_by)
		_, _ = db.Exec("UPDATE users SET created_by = $1, modified_by = $1 WHERE id = $1", newUserID)
		_, _ = db.Exec("UPDATE roles SET created_by = $1, modified_by = $1 WHERE id = $2", newUserID, roleID)

		log.Printf("GO SEEDER: Berhasil membuat akun Super Admin (Username: '%s' | Password: '%s')\n", adminUsername, adminPassword)
	} else if err == nil {
		// User sudah ada -> Update password_hash, role_id, dan status aktif
		queryUpdate := `
			UPDATE users 
			SET role_id = $1, password_hash = $2, is_active = TRUE, deleted_at = NULL, modified_at = NOW()
			WHERE id = $3`
		_, err = db.Exec(queryUpdate, roleID, string(hashedPassword), existingUserID)
		if err != nil {
			return fmt.Errorf("gagal update password super admin: %v", err)
		}

		log.Printf("GO SEEDER: Password Super Admin '%s' berhasil diperbarui & disinkronisasi.\n", adminUsername)
	} else {
		return fmt.Errorf("gagal memeriksa user admin: %v", err)
	}

	return nil
}
