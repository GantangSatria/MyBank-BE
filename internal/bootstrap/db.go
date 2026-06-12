package bootstrap

import (
	"database/sql"
	"log"
	"time"

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

	// Retry ping database
	maxRetries := 20

	for i := 1; i <= maxRetries; i++ {
		err := db.Ping()

		if err == nil {
			log.Println("[db] database connected successfully")

			startKeepAlive(db)

			return db
		}

		log.Printf(
			"[db] waiting database connection (%d/%d): %v",
			i,
			maxRetries,
			err,
		)

		time.Sleep(5 * time.Second)
	}

	log.Fatal("[db] failed to connect database after retries")

	return nil
}

// helper free plan wkwk
func startKeepAlive(db *sql.DB) {
	ticker := time.NewTicker(30 * time.Minute)

	go func() {
		for range ticker.C {
			var result int

			err := db.QueryRow("SELECT 1").Scan(&result)

			if err != nil {
				log.Printf("[db] keepalive failed: %v", err)
				continue
			}

			log.Println("[db] keepalive success")
		}
	}()
}