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

func (r *MySQLEmployeeRepository) Search(ctx context.Context, keyword string, limit, offset int) ([]*models.Employee, int64, error) {
	// Count query
	countQuery := `
		SELECT COUNT(*)
		FROM employees
		WHERE deleted_at IS NULL AND (name LIKE ? OR position LIKE ?)
	`
	searchPattern := "%" + keyword + "%"

	var totalCount int64
	err := r.db.QueryRowContext(ctx, countQuery, searchPattern, searchPattern).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count search results: %w", err)
	}

	// Search query
	query := `
		SELECT id, name, age, position, department_id, salary, created_at, updated_at, deleted_at
		FROM employees
		WHERE deleted_at IS NULL AND (name LIKE ? OR position LIKE ?)
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, searchPattern, searchPattern, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search employees: %w", err)
	}
	defer rows.Close()

	var employees []*models.Employee
	for rows.Next() {
		employee := &models.Employee{}
		err := rows.Scan(
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
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan employee: %w", err)
		}
		employees = append(employees, employee)
	}

	return employees, totalCount, nil
}

func (r *MySQLEmployeeRepository) GetAll(ctx context.Context, filter models.EmployeeFilter) ([]*models.Employee, error) {
	query := `
		SELECT id, name, age, position, department_id, salary, created_at, updated_at, deleted_at
		FROM employees
		WHERE deleted_at IS NULL
	`
	args := []interface{}{}

	if filter.DepartmentID != "" {
		query += " AND department_id = ?"
		args = append(args, filter.DepartmentID)
	}

	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, filter.Limit, filter.Offset)

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query employees: %w", err)
	}
	defer rows.Close()

	var employees []*models.Employee
	for rows.Next() {
		employee := &models.Employee{}
		err := rows.Scan(
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
		if err != nil {
			return nil, fmt.Errorf("failed to scan employee: %w", err)
		}
		employees = append(employees, employee)
	}

	return employees, nil
}

func (r *MySQLEmployeeRepository) Count(ctx context.Context, filter models.EmployeeFilter) (int64, error) {
	query := `SELECT COUNT(*) FROM employees WHERE deleted_at IS NULL`
	args := []interface{}{}

	if filter.DepartmentID != "" {
		query += " AND department_id = ?"
		args = append(args, filter.DepartmentID)
	}

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	var count int64
	err = stmt.QueryRowContext(ctx, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count employees: %w", err)
	}

	return count, nil
}
