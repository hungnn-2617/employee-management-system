package utils

import (
	"errors"
	"strings"

	"employee-management-system/internal/models"
)

var (
	ErrInvalidName         = errors.New("name is required and cannot be empty")
	ErrInvalidAge          = errors.New("age must be greater than 0")
	ErrInvalidSalary       = errors.New("salary must be greater than 0")
	ErrInvalidPosition     = errors.New("position is required and cannot be empty")
	ErrInvalidDepartmentID = errors.New("departmentId is required and cannot be empty")
	ErrInvalidID           = errors.New("invalid ID format")
)

func ValidateCreateEmployeeRequest(req *models.CreateEmployeeRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return ErrInvalidName
	}
	if req.Age <= 0 {
		return ErrInvalidAge
	}
	if req.Salary <= 0 {
		return ErrInvalidSalary
	}
	if strings.TrimSpace(req.Position) == "" {
		return ErrInvalidPosition
	}
	if strings.TrimSpace(req.DepartmentID) == "" {
		return ErrInvalidDepartmentID
	}
	return nil
}

func ValidateUpdateEmployeeRequest(req *models.UpdateEmployeeRequest) error {
	if req.Name != nil && strings.TrimSpace(*req.Name) == "" {
		return ErrInvalidName
	}
	if req.Age != nil && *req.Age <= 0 {
		return ErrInvalidAge
	}
	if req.Salary != nil && *req.Salary <= 0 {
		return ErrInvalidSalary
	}
	if req.Position != nil && strings.TrimSpace(*req.Position) == "" {
		return ErrInvalidPosition
	}
	if req.DepartmentID != nil && strings.TrimSpace(*req.DepartmentID) == "" {
		return ErrInvalidDepartmentID
	}
	return nil
}
