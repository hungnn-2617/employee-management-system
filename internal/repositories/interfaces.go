package repositories

import (
	"context"

	"employee-management-system/internal/models"
)

type EmployeeRepository interface {
	Create(ctx context.Context, employee *models.Employee) error
	GetByID(ctx context.Context, id string) (*models.Employee, error)
	Update(ctx context.Context, employee *models.Employee) error
	GetAll(ctx context.Context, filter models.EmployeeFilter) ([]*models.Employee, error)
	Count(ctx context.Context, filter models.EmployeeFilter) (int64, error)
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context, keyword string, limit, offset int) ([]*models.Employee, int64, error)
	GetByDepartmentID(ctx context.Context, departmentID string, limit, offset int) ([]*models.Employee, int64, error)
}

type DepartmentRepository interface {
	Create(ctx context.Context, department *models.Department) error
	GetByID(ctx context.Context, id string) (*models.Department, error)
	Exists(ctx context.Context, id string) (bool, error)
	Update(ctx context.Context, department *models.Department) error
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context) ([]*models.Department, error)
}
