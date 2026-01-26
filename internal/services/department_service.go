package services

import (
	"context"
	"time"

	"github.com/google/uuid"

	"employee-management-system/internal/config"
	"employee-management-system/internal/models"
	"employee-management-system/internal/repositories"
	"employee-management-system/internal/utils"
)

// DepartmentService handles business logic for departments
type DepartmentService struct {
	departmentRepo repositories.DepartmentRepository
}

// NewDepartmentService creates a new department service
func NewDepartmentService(departmentRepo repositories.DepartmentRepository) *DepartmentService {
	return &DepartmentService{
		departmentRepo: departmentRepo,
	}
}

// GetByID retrieves a department by ID
func (s *DepartmentService) GetByID(ctx context.Context, id string) (*models.Department, error) {
	ctx, cancel := context.WithTimeout(ctx, config.ContextTimeout)
	defer cancel()

	if id == "" {
		return nil, utils.ErrInvalidID
	}

	department, err := s.departmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if department == nil {
		return nil, ErrDepartmentNotFound
	}

	return department, nil
}

// Create creates a new department
func (s *DepartmentService) Create(ctx context.Context, req *models.CreateDepartmentRequest) (*models.Department, error) {
	ctx, cancel := context.WithTimeout(ctx, config.ContextTimeout)
	defer cancel()

	// Validate request
	if err := utils.ValidateCreateDepartmentRequest(req); err != nil {
		return nil, err
	}

	// Create department
	now := time.Now()
	department := &models.Department{
		ID:        uuid.New().String(),
		Name:      req.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.departmentRepo.Create(ctx, department); err != nil {
		return nil, err
	}

	return department, nil
}

func (s *DepartmentService) Update(ctx context.Context, id string, req *models.UpdateDepartmentRequest) (*models.Department, error) {
	ctx, cancel := context.WithTimeout(ctx, config.ContextTimeout)
	defer cancel()

	if id == "" {
		return nil, utils.ErrInvalidID
	}

	// Get existing department
	department, err := s.departmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if department == nil {
		return nil, ErrDepartmentNotFound
	}

	// Update fields if provided
	if req.Name != nil {
		if *req.Name == "" {
			return nil, utils.ErrInvalidName
		}
		department.Name = *req.Name
	}

	department.UpdatedAt = time.Now()

	if err := s.departmentRepo.Update(ctx, department); err != nil {
		return nil, err
	}

	return department, nil
}

func (s *DepartmentService) GetAll(ctx context.Context) (*models.DepartmentListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, config.ContextTimeout)
	defer cancel()

	departments, err := s.departmentRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return &models.DepartmentListResponse{
		TotalCount:  int64(len(departments)),
		Departments: departments,
	}, nil
}

func (s *DepartmentService) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, config.ContextTimeout)
	defer cancel()

	if id == "" {
		return utils.ErrInvalidID
	}

	// Check if department exists
	department, err := s.departmentRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if department == nil {
		return ErrDepartmentNotFound
	}

	return s.departmentRepo.Delete(ctx, id)
}

// GetByDepartmentID retrieves employees by department ID
func (s *EmployeeService) GetByDepartmentID(ctx context.Context, departmentID string, limit, offset int) (*models.EmployeeListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, config.ContextTimeout)
	defer cancel()

	if departmentID == "" {
		return nil, utils.ErrInvalidID
	}

	exists, err := s.departmentRepo.Exists(ctx, departmentID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrDepartmentNotFound
	}

	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	employees, totalCount, err := s.employeeRepo.GetByDepartmentID(ctx, departmentID, limit, offset)
	if err != nil {
		return nil, err
	}

	return &models.EmployeeListResponse{
		TotalCount: totalCount,
		Employees:  employees,
	}, nil
}
