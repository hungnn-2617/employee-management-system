package models

import "time"

type Department struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}

type CreateDepartmentRequest struct {
	Name string `json:"name"`
}

type UpdateDepartmentRequest struct {
	Name *string `json:"name,omitempty"`
}
