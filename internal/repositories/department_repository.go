package repositories

import (
	"context"
	"database/sql"
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
