package router

import (
	"net/http"

	"employee-management-system/internal/handlers"
	"employee-management-system/internal/middleware"
)

type Router struct {
	employeeHandler   *handlers.EmployeeHandler
	departmentHandler *handlers.DepartmentHandler
	exportHandler     *handlers.ExportHandler
}

func NewRouter(
	employeeHandler *handlers.EmployeeHandler,
	departmentHandler *handlers.DepartmentHandler,
	exportHandler *handlers.ExportHandler,
) *Router {
	return &Router{
		employeeHandler:   employeeHandler,
		departmentHandler: departmentHandler,
		exportHandler:     exportHandler,
	}
}

func (router *Router) SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	// Employee routes
	mux.Handle("/employees", router.employeeHandler)
	mux.Handle("/employees/", router.employeeHandler)

	// Department routes
	mux.Handle("/departments", router.departmentHandler)
	mux.Handle("/departments/", router.departmentHandler)

	// Export routes
	mux.Handle("/export/", router.exportHandler)

	// Health check endpoint
	mux.HandleFunc("/health", healthCheckHandler)

	// Apply middleware chain
	handler := middleware.Chain(
		mux,
		middleware.RecoverMiddleware,
		middleware.LoggingMiddleware,
		middleware.CORSMiddleware,
		middleware.BasicAuthMiddleware,
	)

	return handler
}

// healthCheckHandler handles the health check endpoint
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
