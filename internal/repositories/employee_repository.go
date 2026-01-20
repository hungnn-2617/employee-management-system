package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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

func (r *MySQLEmployeeRepository) Update(ctx context.Context, employee *models.Employee) error {
	query := `
	UPDATE employees
	SET name = ?, age = ?, position = ?, department_id = ?, salary = ?, updated_at = ?
	WHERE id = ? AND deleted_at IS NULL
	`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()
	result, err := stmt.ExecContext(ctx,
		employee.Name,
		employee.Age,
		employee.Position,
		employee.DepartmentID,
		employee.Salary,
		time.Now(),
		employee.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update employee: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("employee not found or already deleted")
	}

	return nil
}

func (r *MySQLEmployeeRepository) GetByID(ctx context.Context, id string) (*models.Employee, error) {
	query := `
		SELECT id, name, age, position, department_id, salary, created_at, updated_at, deleted_at
		FROM employees
		WHERE id = ? AND deleted_at IS NULL
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	employee := &models.Employee{}
	err = stmt.QueryRowContext(ctx, id).Scan(
		&employee.ID,
		&employee.Name,
		&employee.Age,
		&employee.Position,
		&employee.DepartmentID,
		&employee.Salary,
		&employee.CreatedAt,
		&employee.UpdatedAt,
		&employee.DeletedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get employee: %w", err)
	}

	return employee, nil
}
