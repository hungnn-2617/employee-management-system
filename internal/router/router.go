package router

import (
	"employee-management-system/internal/handlers"
	"net/http"
)

type Router struct {
	employeeHandler *handlers.EmployeeHandler
}

func NewRouter(
	employeeHandler *handlers.EmployeeHandler,
) *Router {
	return &Router{
		employeeHandler: employeeHandler,
	}
}

func (router *Router) SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	// Employee routes
	mux.Handle("/employees", router.employeeHandler)
	mux.Handle("/employees/", router.employeeHandler)

	mux.HandleFunc("/health", healthCheckHandler)

	return corsMiddleware(mux)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
