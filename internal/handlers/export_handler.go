package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"employee-management-system/internal/services"
	"employee-management-system/internal/utils"
)

type ExportHandler struct {
	exportService *services.ExportService
}

func NewExportHandler(exportService *services.ExportService) *ExportHandler {
	return &ExportHandler{
		exportService: exportService,
	}
}

func (h *ExportHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	path := r.URL.Path

	switch path {
	case "/export/csv":
		h.ExportCSV(w, r)
	default:
		utils.WriteError(w, http.StatusNotFound, "Not found")
	}
}

func (h *ExportHandler) ExportCSV(w http.ResponseWriter, r *http.Request) {
	result, err := h.exportService.ExportToCSV(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to export to CSV: "+err.Error())
		return
	}

	filePath := result.FilePath
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to read exported file: "+err.Error())
		return
	}

	fileName := filepath.Base(filePath)
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	w.Header().Set("Content-Length", string(rune(len(fileData))))

	w.WriteHeader(http.StatusOK)
	w.Write(fileData)
}
