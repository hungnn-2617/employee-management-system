package main

import (
	"log"
	"time"

	"employee-management-system/internal/config"
	"employee-management-system/internal/database"
)

func main() {
	cfg := config.Load()

	log.Println("Connecting to database...")
	db, err := database.NewMySQLConnectionWithRetry(cfg.Database, 5, 3*time.Second)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Database connected successfully!")
}
