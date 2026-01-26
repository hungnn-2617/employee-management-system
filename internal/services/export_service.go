package services

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"employee-management-system/internal/config"
	"employee-management-system/internal/models"
	"employee-management-system/internal/repositories"
)

type ExportService struct {
	employeeRepo repositories.EmployeeRepository
	mu           sync.Mutex
	exportDir    string
}

func NewExportService(employeeRepo repositories.EmployeeRepository, exportDir string) *ExportService {
	return &ExportService{
		employeeRepo: employeeRepo,
		exportDir:    exportDir,
	}
}

type ExportResult struct {
	Format   string `json:"format"`
	FilePath string `json:"filePath"`
	Success  bool   `json:"success"`
	Error    string `json:"error,omitempty"`
	Duration string `json:"duration"`
}

func (s *ExportService) exportToCSV(employees []*models.Employee) ExportResult {
	start := time.Now()
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("%s/employees_%s.csv", s.exportDir, timestamp)

	result := ExportResult{
		Format:   "CSV",
		FilePath: filename,
	}

	s.mu.Lock()
	file, err := os.Create(filename)
	s.mu.Unlock()

	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("failed to create file: %v", err)
		result.Duration = time.Since(start).String()
		return result
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"ID", "Name", "Age", "Position", "DepartmentID", "Salary", "CreatedAt", "UpdatedAt"}
	if err := writer.Write(header); err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("failed to write CSV header: %v", err)
		result.Duration = time.Since(start).String()
		return result
	}

	for _, emp := range employees {
		row := []string{
			emp.ID,
			emp.Name,
			strconv.Itoa(emp.Age),
			emp.Position,
			emp.DepartmentID,
			strconv.FormatFloat(emp.Salary, 'f', 2, 64),
			emp.CreatedAt.Format(time.RFC3339),
			emp.UpdatedAt.Format(time.RFC3339),
		}
		if err := writer.Write(row); err != nil {
			result.Success = false
			result.Error = fmt.Sprintf("failed to write CSV row: %v", err)
			result.Duration = time.Since(start).String()
			return result
		}
	}

	result.Success = true
	result.Duration = time.Since(start).String()
	return result
}

func (s *ExportService) ExportToCSV(ctx context.Context) (*ExportResult, error) {
	ctx, cancel := context.WithTimeout(ctx, config.ContextTimeout)
	defer cancel()

	filter := models.EmployeeFilter{
		Limit:  10000,
		Offset: 0,
	}

	employees, err := s.employeeRepo.GetAll(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get employees: %w", err)
	}

	s.mu.Lock()
	if err := os.MkdirAll(s.exportDir, 0755); err != nil {
		s.mu.Unlock()
		return nil, fmt.Errorf("failed to create export directory: %w", err)
	}
	s.mu.Unlock()

	result := s.exportToCSV(employees)
	return &result, nil
}
