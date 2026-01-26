package services

import (
	"context"
	"testing"
	"time"

	"employee-management-system/internal/models"
)

type MockEmployeeRepository struct {
	employees map[string]*models.Employee
}

func NewMockEmployeeRepository() *MockEmployeeRepository {
	return &MockEmployeeRepository{
		employees: make(map[string]*models.Employee),
	}
}

func (m *MockEmployeeRepository) Create(ctx context.Context, employee *models.Employee) error {
	m.employees[employee.ID] = employee
	return nil
}

func (m *MockEmployeeRepository) GetByID(ctx context.Context, id string) (*models.Employee, error) {
	if emp, ok := m.employees[id]; ok {
		return emp, nil
	}
	return nil, nil
}

func (m *MockEmployeeRepository) GetAll(ctx context.Context, filter models.EmployeeFilter) ([]*models.Employee, error) {
	var result []*models.Employee
	for _, emp := range m.employees {
		if emp.DeletedAt == nil {
			if filter.DepartmentID == "" || emp.DepartmentID == filter.DepartmentID {
				result = append(result, emp)
			}
		}
	}
	return result, nil
}

func (m *MockEmployeeRepository) Count(ctx context.Context, filter models.EmployeeFilter) (int64, error) {
	count := int64(0)
	for _, emp := range m.employees {
		if emp.DeletedAt == nil {
			if filter.DepartmentID == "" || emp.DepartmentID == filter.DepartmentID {
				count++
			}
		}
	}
	return count, nil
}

func (m *MockEmployeeRepository) Update(ctx context.Context, employee *models.Employee) error {
	m.employees[employee.ID] = employee
	return nil
}

func (m *MockEmployeeRepository) Delete(ctx context.Context, id string) error {
	if emp, ok := m.employees[id]; ok {
		now := time.Now()
		emp.DeletedAt = &now
	}
	return nil
}

func (m *MockEmployeeRepository) HardDelete(ctx context.Context, id string) error {
	delete(m.employees, id)
	return nil
}

func (m *MockEmployeeRepository) Search(ctx context.Context, keyword string, limit, offset int) ([]*models.Employee, int64, error) {
	return nil, 0, nil
}

func (m *MockEmployeeRepository) GetByDepartmentID(ctx context.Context, departmentID string, limit, offset int) ([]*models.Employee, int64, error) {
	return nil, 0, nil
}

// MockDepartmentRepository is a mock implementation of DepartmentRepository
type MockDepartmentRepository struct {
	departments map[string]*models.Department
}

func NewMockDepartmentRepository() *MockDepartmentRepository {
	return &MockDepartmentRepository{
		departments: make(map[string]*models.Department),
	}
}

func (m *MockDepartmentRepository) Create(ctx context.Context, department *models.Department) error {
	m.departments[department.ID] = department
	return nil
}

func (m *MockDepartmentRepository) GetByID(ctx context.Context, id string) (*models.Department, error) {
	if dept, ok := m.departments[id]; ok {
		return dept, nil
	}
	return nil, nil
}

func (m *MockDepartmentRepository) GetAll(ctx context.Context) ([]*models.Department, error) {
	var result []*models.Department
	for _, dept := range m.departments {
		if dept.DeletedAt == nil {
			result = append(result, dept)
		}
	}
	return result, nil
}

func (m *MockDepartmentRepository) Update(ctx context.Context, department *models.Department) error {
	m.departments[department.ID] = department
	return nil
}

func (m *MockDepartmentRepository) Delete(ctx context.Context, id string) error {
	if dept, ok := m.departments[id]; ok {
		now := time.Now()
		dept.DeletedAt = &now
	}
	return nil
}

func (m *MockDepartmentRepository) HardDelete(ctx context.Context, id string) error {
	delete(m.departments, id)
	return nil
}

func (m *MockDepartmentRepository) Exists(ctx context.Context, id string) (bool, error) {
	if dept, ok := m.departments[id]; ok && dept.DeletedAt == nil {
		return true, nil
	}
	return false, nil
}

// Test cases for EmployeeService
func TestEmployeeService_Create(t *testing.T) {
	empRepo := NewMockEmployeeRepository()
	deptRepo := NewMockDepartmentRepository()
	service := NewEmployeeService(empRepo, deptRepo)

	// Setup: Create a department first
	ctx := context.Background()
	deptRepo.departments["dept-1"] = &models.Department{
		ID:        "dept-1",
		Name:      "Engineering",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name    string
		req     *models.CreateEmployeeRequest
		wantErr bool
	}{
		{
			name: "valid employee",
			req: &models.CreateEmployeeRequest{
				Name:         "John Doe",
				Age:          30,
				Position:     "Developer",
				DepartmentID: "dept-1",
				Salary:       5000,
			},
			wantErr: false,
		},
		{
			name: "invalid age",
			req: &models.CreateEmployeeRequest{
				Name:         "John Doe",
				Age:          0,
				Position:     "Developer",
				DepartmentID: "dept-1",
				Salary:       5000,
			},
			wantErr: true,
		},
		{
			name: "invalid salary",
			req: &models.CreateEmployeeRequest{
				Name:         "John Doe",
				Age:          30,
				Position:     "Developer",
				DepartmentID: "dept-1",
				Salary:       -100,
			},
			wantErr: true,
		},
		{
			name: "empty name",
			req: &models.CreateEmployeeRequest{
				Name:         "",
				Age:          30,
				Position:     "Developer",
				DepartmentID: "dept-1",
				Salary:       5000,
			},
			wantErr: true,
		},
		{
			name: "invalid department",
			req: &models.CreateEmployeeRequest{
				Name:         "John Doe",
				Age:          30,
				Position:     "Developer",
				DepartmentID: "non-existent",
				Salary:       5000,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			emp, err := service.Create(ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && emp == nil {
				t.Error("Create() returned nil employee when no error expected")
			}
		})
	}
}

func TestEmployeeService_GetByID(t *testing.T) {
	empRepo := NewMockEmployeeRepository()
	deptRepo := NewMockDepartmentRepository()
	service := NewEmployeeService(empRepo, deptRepo)

	ctx := context.Background()

	// Setup: Add an employee
	empRepo.employees["emp-1"] = &models.Employee{
		ID:           "emp-1",
		Name:         "John Doe",
		Age:          30,
		Position:     "Developer",
		DepartmentID: "dept-1",
		Salary:       5000,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "existing employee",
			id:      "emp-1",
			wantErr: false,
		},
		{
			name:    "non-existing employee",
			id:      "non-existent",
			wantErr: true,
		},
		{
			name:    "empty id",
			id:      "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			emp, err := service.GetByID(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && emp == nil {
				t.Error("GetByID() returned nil employee when no error expected")
			}
		})
	}
}

func TestDepartmentService_Create(t *testing.T) {
	deptRepo := NewMockDepartmentRepository()
	service := NewDepartmentService(deptRepo)

	ctx := context.Background()

	tests := []struct {
		name    string
		req     *models.CreateDepartmentRequest
		wantErr bool
	}{
		{
			name:    "valid department",
			req:     &models.CreateDepartmentRequest{Name: "Engineering"},
			wantErr: false,
		},
		{
			name:    "empty name",
			req:     &models.CreateDepartmentRequest{Name: ""},
			wantErr: true,
		},
		{
			name:    "whitespace name",
			req:     &models.CreateDepartmentRequest{Name: "   "},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dept, err := service.Create(ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && dept == nil {
				t.Error("Create() returned nil department when no error expected")
			}
		})
	}
}
