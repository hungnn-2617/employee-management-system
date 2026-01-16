package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"employee-management-system/internal/config"
)

func NewMySQLConnection(cfg config.DatabaseConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// NewMySQLConnectionWithRetry creates a MySQL connection with retry logic
func NewMySQLConnectionWithRetry(cfg config.DatabaseConfig, maxRetries int, retryInterval time.Duration) (*sql.DB, error) {
	var db *sql.DB
	var err error

	for i := 0; i < maxRetries; i++ {
		db, err = NewMySQLConnection(cfg)
		if err == nil {
			return db, nil
		}

		fmt.Printf("Failed to connect to database (attempt %d/%d): %v\n", i+1, maxRetries, err)
		time.Sleep(retryInterval)
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
}
