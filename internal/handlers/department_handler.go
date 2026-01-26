package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"employee-management-system/internal/models"
	"employee-management-system/internal/services"
	"employee-management-system/internal/utils"
)

type DepartmentHandler struct {
	departmentService *services.DepartmentService
	employeeService   *services.EmployeeService
}

func NewDepartmentHandler(departmentService *services.DepartmentService, employeeService *services.EmployeeService) *DepartmentHandler {
	return &DepartmentHandler{
		departmentService: departmentService,
		employeeService:   employeeService,
	}
}

func (h *DepartmentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/departments")
	path = strings.TrimSuffix(path, "/")

	switch {
	case path == "" && r.Method == http.MethodGet:
		h.GetAll(w, r)
	case path == "" && r.Method == http.MethodPost:
		h.Create(w, r)
	case strings.HasSuffix(path, "/employees") && r.Method == http.MethodGet:
		// Handle /departments/:id/employees
		id := strings.TrimSuffix(strings.TrimPrefix(path, "/"), "/employees")
		h.GetEmployees(w, r, id)
	case strings.HasPrefix(path, "/") && r.Method == http.MethodGet:
		h.GetByID(w, r, strings.TrimPrefix(path, "/"))
	case strings.HasPrefix(path, "/") && r.Method == http.MethodPut:
		h.Update(w, r, strings.TrimPrefix(path, "/"))
	case strings.HasPrefix(path, "/") && r.Method == http.MethodDelete:
		h.Delete(w, r, strings.TrimPrefix(path, "/"))
	default:
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// Create handles POST /departments
func (h *DepartmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateDepartmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	department, err := h.departmentService.Create(r.Context(), &req)
	if err != nil {
		switch err {
		case utils.ErrInvalidName:
			utils.WriteError(w, http.StatusBadRequest, err.Error())
		default:
			utils.WriteError(w, http.StatusInternalServerError, "Failed to create department")
		}
		return
	}

	utils.WriteJSON(w, http.StatusCreated, department)
}

// Update handles PUT /departments/:id
func (h *DepartmentHandler) Update(w http.ResponseWriter, r *http.Request, id string) {
	var req models.UpdateDepartmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	department, err := h.departmentService.Update(r.Context(), id, &req)
	if err != nil {
		switch err {
		case services.ErrDepartmentNotFound:
			utils.WriteError(w, http.StatusNotFound, "Department not found")
		case utils.ErrInvalidID:
			utils.WriteError(w, http.StatusBadRequest, "Invalid department ID")
		case utils.ErrInvalidName:
			utils.WriteError(w, http.StatusBadRequest, err.Error())
		default:
			utils.WriteError(w, http.StatusInternalServerError, "Failed to update department")
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, department)
}

// GetAll handles GET /departments
func (h *DepartmentHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	response, err := h.departmentService.GetAll(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to get departments")
		return
	}

	utils.WriteJSON(w, http.StatusOK, response)
}

// GetByID handles GET /departments/:id
func (h *DepartmentHandler) GetByID(w http.ResponseWriter, r *http.Request, id string) {
	department, err := h.departmentService.GetByID(r.Context(), id)
	if err != nil {
		switch err {
		case services.ErrDepartmentNotFound:
			utils.WriteError(w, http.StatusNotFound, "Department not found")
		case utils.ErrInvalidID:
			utils.WriteError(w, http.StatusBadRequest, "Invalid department ID")
		default:
			utils.WriteError(w, http.StatusInternalServerError, "Failed to get department")
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, department)
}

// Delete handles DELETE /departments/:id
func (h *DepartmentHandler) Delete(w http.ResponseWriter, r *http.Request, id string) {
	err := h.departmentService.Delete(r.Context(), id)
	if err != nil {
		switch err {
		case services.ErrDepartmentNotFound:
			utils.WriteError(w, http.StatusNotFound, "Department not found")
		case utils.ErrInvalidID:
			utils.WriteError(w, http.StatusBadRequest, "Invalid department ID")
		default:
			utils.WriteError(w, http.StatusInternalServerError, "Failed to delete department")
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Department deleted successfully"})
}

func (h *DepartmentHandler) GetEmployees(w http.ResponseWriter, r *http.Request, departmentID string) {
	query := r.URL.Query()
	limit, _ := strconv.Atoi(query.Get("limit"))
	offset, _ := strconv.Atoi(query.Get("offset"))

	response, err := h.employeeService.GetByDepartmentID(r.Context(), departmentID, limit, offset)
	if err != nil {
		switch err {
		case services.ErrDepartmentNotFound:
			utils.WriteError(w, http.StatusNotFound, "Department not found")
		case utils.ErrInvalidID:
			utils.WriteError(w, http.StatusBadRequest, "Invalid department ID")
		default:
			utils.WriteError(w, http.StatusInternalServerError, "Failed to get employees")
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, response)
}
