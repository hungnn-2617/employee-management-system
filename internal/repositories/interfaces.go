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
}

type DepartmentRepository interface {
	Exists(ctx context.Context, id string) (bool, error)
}
