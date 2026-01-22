package handlers

import (
	"employee-management-system/internal/models"
	"employee-management-system/internal/services"
	"employee-management-system/internal/utils"
	"encoding/json"
	"net/http"
	"strconv"
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
	case path == "" && r.Method == http.MethodGet:
		h.GetAll(w, r)
	case path == "" && r.Method == http.MethodPost:
		h.Create(w, r)
	case strings.HasPrefix(path, "/") && r.Method == http.MethodGet:
		h.GetByID(w, r, strings.TrimPrefix(path, "/"))
	case strings.HasPrefix(path, "/") && r.Method == http.MethodPut:
		h.Update(w, r, strings.TrimPrefix(path, "/"))
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

// GetByID handles GET /employees/:id
func (h *EmployeeHandler) GetByID(w http.ResponseWriter, r *http.Request, id string) {
	employee, err := h.employeeService.GetByID(r.Context(), id)
	if err != nil {
		switch err {
		case services.ErrEmployeeNotFound:
			utils.WriteError(w, http.StatusNotFound, "Employee not found")
		case utils.ErrInvalidID:
			utils.WriteError(w, http.StatusBadRequest, "Invalid employee ID")
		default:
			utils.WriteError(w, http.StatusInternalServerError, "Failed to get employee")
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, employee)
}

// Update handles PUT /employees/{id}
func (h *EmployeeHandler) Update(w http.ResponseWriter, r *http.Request, id string) {
	var req models.UpdateEmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	employee, err := h.employeeService.Update(r.Context(), id, &req)
	if err != nil {
		switch err {
		case services.ErrEmployeeNotFound:
			utils.WriteError(w, http.StatusNotFound, "Employee not found")
		case services.ErrDepartmentNotFound:
			utils.WriteError(w, http.StatusNotFound, "Department not found")
		case utils.ErrInvalidID:
			utils.WriteError(w, http.StatusBadRequest, "Invalid employee ID")
		case utils.ErrInvalidName, utils.ErrInvalidAge, utils.ErrInvalidSalary,
			utils.ErrInvalidPosition, utils.ErrInvalidDepartmentID:
			utils.WriteError(w, http.StatusBadRequest, err.Error())
		default:
			utils.WriteError(w, http.StatusInternalServerError, "Failed to update employee")
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, employee)
}

// GetAll handles GET /employees
func (h *EmployeeHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	limit, _ := strconv.Atoi(query.Get("limit"))
	offset, _ := strconv.Atoi(query.Get("offset"))
	departmentID := query.Get("departmentId")

	filter := models.EmployeeFilter{
		DepartmentID: departmentID,
		Limit:        limit,
		Offset:       offset,
	}

	response, err := h.employeeService.GetAll(r.Context(), filter)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to get employees")
		return
	}

	utils.WriteJSON(w, http.StatusOK, response)
}
