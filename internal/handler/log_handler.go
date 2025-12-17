package handler

import (
	"net/http"
	"strconv"

	"go-admin/internal/logger"
	"go-admin/internal/service"
	"go-admin/pkg/response"

	"github.com/gin-gonic/gin"
)

// LogHandler handles log-related HTTP requests
type LogHandler struct {
	logService service.LogService
}

// NewLogHandler creates a new log handler
func NewLogHandler() *LogHandler {
	return &LogHandler{
		logService: service.NewLogService(),
	}
}

// GetLogByID handles requests to get a log by ID
func (h *LogHandler) GetLogByID(c *gin.Context) {
	// Parse log ID from URL parameter
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid log ID")
		return
	}

	// Get log
	log, err := h.logService.GetLogByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "Log not found")
		return
	}

	response.Success(c, "Log retrieved successfully", log)
}

// ListLogs handles requests to list logs with pagination
func (h *LogHandler) ListLogs(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// Ensure page and pageSize are positive
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	// List logs
	logs, total, err := h.logService.ListLogs(page, pageSize, "", "", "", "")
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to list logs")
		return
	}

	response.Success(c, "Logs retrieved successfully", gin.H{
		"data":      logs,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// DeleteLog handles requests to delete a log
func (h *LogHandler) DeleteLog(c *gin.Context) {
	// Parse log ID from URL parameter
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid log ID")
		return
	}

	// Delete log
	if err := h.logService.DeleteLog(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete log: "+err.Error())
		return
	}

	response.Success(c, "Log deleted successfully", nil)
}

// ClearLogs handles requests to clear all logs
func (h *LogHandler) ClearLogs(c *gin.Context) {
	// Clear logs
	_, err := h.logService.ClearLogs("", "", "", "", 0)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to clear logs: "+err.Error())
		return
	}

	response.Success(c, "Logs cleared successfully", nil)
}

// LogLevelHandler handles log level-related HTTP requests
type LogLevelHandler struct{}

// NewLogLevelHandler creates a new log level handler
func NewLogLevelHandler() *LogLevelHandler {
	return &LogLevelHandler{}
}

// GetLogLevel handles requests to get current log level
func (h *LogLevelHandler) GetLogLevel(c *gin.Context) {
	level := logger.GetLevel()

	response.Success(c, "Current log level retrieved successfully", gin.H{
		"level": level,
	})
}

// SetLogLevel handles requests to set log level
func (h *LogLevelHandler) SetLogLevel(c *gin.Context) {
	// Define request structure
	var req struct {
		Level string `json:"level" binding:"required"`
	}

	// Bind JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	// Set log level
	if err := logger.SetLevel(req.Level); err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to set log level: "+err.Error())
		return
	}

	response.Success(c, "Log level updated successfully", gin.H{
		"level": req.Level,
	})
}
