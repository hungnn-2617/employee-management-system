package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"employee-management-system/internal/models"
)

type MySQLEmployeeRepository struct {
	db *sql.DB
}

func NewMySQLEmployeeRepository(db *sql.DB) *MySQLEmployeeRepository {
	return &MySQLEmployeeRepository{db: db}
}

func (r *MySQLEmployeeRepository) Create(ctx context.Context, employee *models.Employee) error {
	query := `
		INSERT INTO employees (id, name, age, position, department_id, salary, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		employee.ID,
		employee.Name,
		employee.Age,
		employee.Position,
		employee.DepartmentID,
		employee.Salary,
		employee.CreatedAt,
		employee.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert employee: %w", err)
	}

	return nil
}
