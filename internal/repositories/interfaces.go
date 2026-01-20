package repositories

import (
	"context"

	"employee-management-system/internal/models"
)

type EmployeeRepository interface {
	Create(ctx context.Context, employee *models.Employee) error
}

type DepartmentRepository interface {
	Exists(ctx context.Context, id string) (bool, error)
}
