package middleware

import (
	"employee-management-system/internal/utils"
	"encoding/base64"
	"net/http"
	"strings"
)

var validUsers = map[string]string{
	"admin": "admin123",
	"user":  "user123",
}

func BasicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.Header().Set("WWW-Authenticate", `Basic realm="Employee Management System"`)
			utils.WriteError(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		if !strings.HasPrefix(authHeader, "Basic ") {
			utils.WriteError(w, http.StatusUnauthorized, "Invalid authorization format")
			return
		}

		encoded := strings.TrimPrefix(authHeader, "Basic ")
		decoded, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			utils.WriteError(w, http.StatusUnauthorized, "Invalid authorization encoding")
			return
		}

		credentials := strings.SplitN(string(decoded), ":", 2)
		if len(credentials) != 2 {
			utils.WriteError(w, http.StatusUnauthorized, "Invalid credentials format")
			return
		}

		username := credentials[0]
		password := credentials[1]

		if storedPassword, exists := validUsers[username]; !exists || storedPassword != password {
			utils.WriteError(w, http.StatusUnauthorized, "Invalid username or password")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func OptionalBasicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			next.ServeHTTP(w, r)
			return
		}

		if !strings.HasPrefix(authHeader, "Basic ") {
			utils.WriteError(w, http.StatusUnauthorized, "Invalid authorization format")
			return
		}

		encoded := strings.TrimPrefix(authHeader, "Basic ")
		decoded, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			utils.WriteError(w, http.StatusUnauthorized, "Invalid authorization encoding")
			return
		}

		credentials := strings.SplitN(string(decoded), ":", 2)
		if len(credentials) != 2 {
			utils.WriteError(w, http.StatusUnauthorized, "Invalid credentials format")
			return
		}

		username := credentials[0]
		password := credentials[1]

		if storedPassword, exists := validUsers[username]; !exists || storedPassword != password {
			utils.WriteError(w, http.StatusUnauthorized, "Invalid username or password")
			return
		}

		next.ServeHTTP(w, r)
	})
}
