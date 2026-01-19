package handlers

import (
	"employee-management-system/internal/models"
	"employee-management-system/internal/services"
	"employee-management-system/internal/utils"
	"encoding/json"
	"net/http"
	"strings"
)

type EmployeeHandler struct {
	employeeService *services.EmployeeService
}

func NewEmployeeHandler(employeeService *services.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{
		employeeService: employeeService,
	}
}

func (h *EmployeeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/employees")
	path = strings.TrimSuffix(path, "/")

	switch {
	case path == "" && r.Method == http.MethodPost:
		h.Create(w, r)
	default:
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// Create handles POST /employees
func (h *EmployeeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateEmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	employee, err := h.employeeService.Create(r.Context(), &req)
	if err != nil {
		switch err {
		case services.ErrDepartmentNotFound:
			utils.WriteError(w, http.StatusNotFound, "Department not found")
		case utils.ErrInvalidName, utils.ErrInvalidAge, utils.ErrInvalidSalary,
			utils.ErrInvalidPosition, utils.ErrInvalidDepartmentID:
			utils.WriteError(w, http.StatusBadRequest, err.Error())
		default:
			utils.WriteError(w, http.StatusInternalServerError, "Failed to create employee")
		}
		return
	}

	utils.WriteJSON(w, http.StatusCreated, employee)
}
