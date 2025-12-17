package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go-admin/internal/service"
	"go-admin/pkg/response"

	"github.com/gin-gonic/gin"
)

// ImportExportHandler handles import/export-related HTTP requests
type ImportExportHandler struct {
	importExportService *service.ImportExportService
}

// NewImportExportHandler creates a new import/export handler
func NewImportExportHandler() *ImportExportHandler {
	return &ImportExportHandler{
		importExportService: service.NewImportExportService(),
	}
}

// ExportUsers handles requests to export users data to Excel
func (h *ImportExportHandler) ExportUsers(c *gin.Context) {
	// For demonstration, we'll create sample user data
	// In a real application, you would fetch this from the database
	headers := []string{"ID", "Username", "Email", "Created At"}
	data := [][]interface{}{
		{1, "john_doe", "john@example.com", "2023-01-15"},
		{2, "jane_smith", "jane@example.com", "2023-02-20"},
		{3, "bob_johnson", "bob@example.com", "2023-03-10"},
	}

	// Export to Excel
	buffer, err := h.importExportService.ExportToExcel(headers, data, "Users")
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to export users: "+err.Error())
		return
	}

	// Set headers for file download
	filename := fmt.Sprintf("users_export_%s.xlsx", time.Now().Format("20060102_150405"))
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	// Send the Excel file
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buffer.Bytes())
}

// ImportUsers handles requests to import users data from Excel
func (h *ImportExportHandler) ImportUsers(c *gin.Context) {
	// Get uploaded file from form
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to get uploaded file")
		return
	}

	// Check file extension
	if file.Filename[len(file.Filename)-5:] != ".xlsx" && file.Filename[len(file.Filename)-4:] != ".xls" {
		response.Error(c, http.StatusBadRequest, "Only Excel files (.xlsx, .xls) are allowed")
		return
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to open uploaded file")
		return
	}
	defer src.Close()

	// Parse has_headers parameter
	hasHeadersStr := c.DefaultQuery("has_headers", "true")
	hasHeaders, err := strconv.ParseBool(hasHeadersStr)
	if err != nil {
		hasHeaders = true
	}

	// Import from Excel
	rows, err := h.importExportService.ImportFromExcel(src, hasHeaders)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to import users: "+err.Error())
		return
	}

	// For demonstration, we'll just return the parsed data
	// In a real application, you would process and save this data to the database
	response.Success(c, "Users imported successfully", gin.H{
		"row_count": len(rows),
		"data":      rows,
	})
}

// ExportData handles requests to export generic data to Excel
func (h *ImportExportHandler) ExportData(c *gin.Context) {
	// Parse parameters from request
	sheetName := c.DefaultQuery("sheet_name", "Data")

	// For demonstration, we'll create sample data
	// In a real application, you would get this data from request body or database
	headers := []string{"Column 1", "Column 2", "Column 3"}
	data := [][]interface{}{
		{"Value 1", "Value 2", "Value 3"},
		{"Value 4", "Value 5", "Value 6"},
		{"Value 7", "Value 8", "Value 9"},
	}

	// Export to Excel
	buffer, err := h.importExportService.ExportToExcel(headers, data, sheetName)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to export data: "+err.Error())
		return
	}

	// Set headers for file download
	filename := fmt.Sprintf("data_export_%s.xlsx", time.Now().Format("20060102_150405"))
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	// Send the Excel file
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buffer.Bytes())
}
