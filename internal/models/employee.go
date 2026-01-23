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

type UpdateEmployeeRequest struct {
	Name         *string  `json:"name,omitempty"`
	Age          *int     `json:"age,omitempty"`
	Position     *string  `json:"position,omitempty"`
	DepartmentID *string  `json:"departmentId,omitempty"`
	Salary       *float64 `json:"salary,omitempty"`
}

type EmployeeListResponse struct {
	TotalCount int64       `json:"totalCount"`
	Employees  []*Employee `json:"employees"`
}

type EmployeeFilter struct {
	DepartmentID string
	Keyword      string
	Limit        int
	Offset       int
}
