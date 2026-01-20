package models

import "time"

type Employee struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Age          int        `json:"age"`
	Position     string     `json:"position"`
	DepartmentID string     `json:"departmentId"`
	Salary       float64    `json:"salary"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	DeletedAt    *time.Time `json:"deletedAt,omitempty"`
}

type CreateEmployeeRequest struct {
	Name         string  `json:"name"`
	Age          int     `json:"age"`
	Position     string  `json:"position"`
	DepartmentID string  `json:"departmentId"`
	Salary       float64 `json:"salary"`
}
