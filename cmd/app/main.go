package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"employee-management-system/internal/config"
	"employee-management-system/internal/database"
	"employee-management-system/internal/handlers"
	"employee-management-system/internal/repositories"
	"employee-management-system/internal/router"
	"employee-management-system/internal/services"
)

func main() {
	cfg := config.Load()

	log.Println("Connecting to database...")
	db, err := database.NewMySQLConnectionWithRetry(cfg.Database, 5, 3*time.Second)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Database connected successfully")

	employeeRepo := repositories.NewMySQLEmployeeRepository(db)
	departmentRepo := repositories.NewMySQLDepartmentRepository(db)

	employeeService := services.NewEmployeeService(employeeRepo, departmentRepo)
	departmentService := services.NewDepartmentService(departmentRepo)
	exportService := services.NewExportService(employeeRepo, "./exports")

	employeeHandler := handlers.NewEmployeeHandler(employeeService)
	departmentHandler := handlers.NewDepartmentHandler(departmentService, employeeService)
	exportHandler := handlers.NewExportHandler(exportService)

	appRouter := router.NewRouter(employeeHandler, departmentHandler, exportHandler)
	handler := appRouter.SetupRoutes()

	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      handler,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		log.Printf("Server starting on port %s", cfg.Server.Port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	fmt.Println("Server stopped")
}
