package bootstrap

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/GantangSatria/MyBank-BE/config"
)

func NewDB(cfg *config.Config) *sql.DB {
	db, err := sql.Open("mysql", cfg.DB.DSN())
	if err != nil {
		log.Fatalf("[db] failed to open database: %v", err)
	}

	// Verifikasi koneksi aktif
	if err := db.Ping(); err != nil {
		log.Fatalf("[db] failed to ping database: %v", err)
	}

	// Connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)

	log.Println("[db] database connected successfully")
	return db
}