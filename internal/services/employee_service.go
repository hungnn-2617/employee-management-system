package services

import (
	"context"
	"errors"
	"time"

	"employee-management-system/internal/config"
	"employee-management-system/internal/models"
	"employee-management-system/internal/repositories"
	"employee-management-system/internal/utils"

	"github.com/google/uuid"
)

var (
	ErrEmployeeNotFound   = errors.New("employee not found")
	ErrDepartmentNotFound = errors.New("department not found")
	ErrInvalidInput       = errors.New("invalid input data")
)

type EmployeeService struct {
	employeeRepo   repositories.EmployeeRepository
	departmentRepo repositories.DepartmentRepository
}

func NewEmployeeService(employeeRepo repositories.EmployeeRepository, departmentRepo repositories.DepartmentRepository) *EmployeeService {
	return &EmployeeService{
		employeeRepo:   employeeRepo,
		departmentRepo: departmentRepo,
	}
}

func (s *EmployeeService) Create(ctx context.Context, req *models.CreateEmployeeRequest) (*models.Employee, error) {
	ctx, cancel := context.WithTimeout(ctx, config.ContextTimeout)
	defer cancel()

	if err := utils.ValidateCreateEmployeeRequest(req); err != nil {
		return nil, err
	}

	exists, err := s.departmentRepo.Exists(ctx, req.DepartmentID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrDepartmentNotFound
	}

	now := time.Now()
	employee := &models.Employee{
		ID:           uuid.New().String(),
		Name:         req.Name,
		Age:          req.Age,
		Position:     req.Position,
		DepartmentID: req.DepartmentID,
		Salary:       req.Salary,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.employeeRepo.Create(ctx, employee); err != nil {
		return nil, err
	}

	return employee, nil
}
