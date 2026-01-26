package repositories

import (
	"context"
	"database/sql"
	"employee-management-system/internal/models"
	"fmt"
	"time"
)

type MySQLDepartmentRepository struct {
	db *sql.DB
}

func NewMySQLDepartmentRepository(db *sql.DB) *MySQLDepartmentRepository {
	return &MySQLDepartmentRepository{db: db}
}

func (r *MySQLDepartmentRepository) Exists(ctx context.Context, id string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM departments WHERE id = ?)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *MySQLDepartmentRepository) Create(ctx context.Context, department *models.Department) error {
	query := `
		INSERT INTO departments (id, name, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		department.ID,
		department.Name,
		department.CreatedAt,
		department.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert department: %w", err)
	}

	return nil
}

func (r *MySQLDepartmentRepository) GetByID(ctx context.Context, id string) (*models.Department, error) {
	query := `
		SELECT id, name, created_at, updated_at, deleted_at
		FROM departments
		WHERE id = ? AND deleted_at IS NULL
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	department := &models.Department{}
	err = stmt.QueryRowContext(ctx, id).Scan(
		&department.ID,
		&department.Name,
		&department.CreatedAt,
		&department.UpdatedAt,
		&department.DeletedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get department: %w", err)
	}

	return department, nil
}

func (r *MySQLDepartmentRepository) Update(ctx context.Context, department *models.Department) error {
	query := `
		UPDATE departments
		SET name = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx,
		department.Name,
		time.Now(),
		department.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update department: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("department not found or already deleted")
	}

	return nil
}

func (r *MySQLDepartmentRepository) GetAll(ctx context.Context) ([]*models.Department, error) {
	query := `
		SELECT id, name, created_at, updated_at, deleted_at
		FROM departments
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query departments: %w", err)
	}
	defer rows.Close()

	var departments []*models.Department
	for rows.Next() {
		department := &models.Department{}
		err := rows.Scan(
			&department.ID,
			&department.Name,
			&department.CreatedAt,
			&department.UpdatedAt,
			&department.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan department: %w", err)
		}
		departments = append(departments, department)
	}

	return departments, nil
}

func (r *MySQLDepartmentRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE departments SET deleted_at = ? WHERE id = ? AND deleted_at IS NULL`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to delete department: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("department not found or already deleted")
	}

	return nil
}
