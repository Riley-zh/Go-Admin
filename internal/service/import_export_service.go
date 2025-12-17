package service

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/xuri/excelize/v2"
)

// ImportExportService handles data import and export operations
type ImportExportService struct{}

// NewImportExportService creates a new import/export service
func NewImportExportService() *ImportExportService {
	return &ImportExportService{}
}

// ExportToExcel exports data to Excel format
func (s *ImportExportService) ExportToExcel(headers []string, data [][]interface{}, sheetName string) (*bytes.Buffer, error) {
	// Create a new Excel file
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Create a new sheet or use the default one
	if sheetName != "" {
		f.SetSheetName("Sheet1", sheetName)
	} else {
		sheetName = "Sheet1"
	}

	// Write headers
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	// Write data
	for rowIdx, rowData := range data {
		for colIdx, cellData := range rowData {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
			f.SetCellValue(sheetName, cell, cellData)
		}
	}

	// Save to buffer
	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("failed to write Excel to buffer: %w", err)
	}

	return buffer, nil
}

// ImportFromExcel imports data from Excel file
func (s *ImportExportService) ImportFromExcel(file multipart.File, hasHeaders bool) ([][]string, error) {
	// Read the file content
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Open Excel file from bytes
	f, err := excelize.OpenReader(bytes.NewReader(content))
	if err != nil {
		return nil, fmt.Errorf("failed to open Excel file: %w", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Get all sheets
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("no sheets found in Excel file")
	}

	// Use the first sheet
	sheetName := sheets[0]

	// Read all rows
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to read rows from sheet: %w", err)
	}

	// Process rows based on whether headers are included
	var result [][]string
	if hasHeaders && len(rows) > 0 {
		// Skip the first row (headers)
		result = rows[1:]
	} else {
		result = rows
	}

	return result, nil
}
