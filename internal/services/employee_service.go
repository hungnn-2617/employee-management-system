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

func (s *EmployeeService) GetByID(ctx context.Context, id string) (*models.Employee, error) {
	ctx, cancel := context.WithTimeout(ctx, config.ContextTimeout)
	defer cancel()

	if id == "" {
		return nil, utils.ErrInvalidID
	}

	employee, err := s.employeeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if employee == nil {
		return nil, ErrEmployeeNotFound
	}

	return employee, nil
}

func (s *EmployeeService) Update(ctx context.Context, id string, req *models.UpdateEmployeeRequest) (*models.Employee, error) {
	ctx, cancel := context.WithTimeout(ctx, config.ContextTimeout)
	defer cancel()

	if id == "" {
		return nil, utils.ErrInvalidID
	}

	// Validate request
	if err := utils.ValidateUpdateEmployeeRequest(req); err != nil {
		return nil, err
	}

	// Get existing employee
	employee, err := s.employeeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if employee == nil {
		return nil, ErrEmployeeNotFound
	}

	// Update fields if provided
	if req.Name != nil {
		employee.Name = *req.Name
	}
	if req.Age != nil {
		employee.Age = *req.Age
	}
	if req.Position != nil {
		employee.Position = *req.Position
	}
	if req.DepartmentID != nil {
		// Check if new department exists
		exists, err := s.departmentRepo.Exists(ctx, *req.DepartmentID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, ErrDepartmentNotFound
		}
		employee.DepartmentID = *req.DepartmentID
	}
	if req.Salary != nil {
		employee.Salary = *req.Salary
	}

	employee.UpdatedAt = time.Now()

	if err := s.employeeRepo.Update(ctx, employee); err != nil {
		return nil, err
	}

	return employee, nil
}

func (s *EmployeeService) GetAll(ctx context.Context, filter models.EmployeeFilter) (*models.EmployeeListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, config.ContextTimeout)
	defer cancel()

	if filter.Limit <= 0 {
		filter.Limit = 10
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	employees, err := s.employeeRepo.GetAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	totalCount, err := s.employeeRepo.Count(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &models.EmployeeListResponse{
		TotalCount: totalCount,
		Employees:  employees,
	}, nil
}
