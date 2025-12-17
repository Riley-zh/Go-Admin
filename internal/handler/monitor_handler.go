package handler

import (
	"strconv"

	"go-admin/internal/service"
	"go-admin/pkg/response"

	"github.com/gin-gonic/gin"
)

// MonitorHandler handles monitoring-related HTTP requests
type MonitorHandler struct {
	monitorService *service.MonitorService
}

// NewMonitorHandler creates a new monitor handler
func NewMonitorHandler() *MonitorHandler {
	return &MonitorHandler{
		monitorService: service.NewMonitorService(),
	}
}

// GetSystemInfo handles requests to get system information
func (h *MonitorHandler) GetSystemInfo(c *gin.Context) {
	systemInfo := h.monitorService.GetSystemInfo()
	response.Success(c, "System information retrieved successfully", systemInfo)
}

// GetSystemMetrics handles requests to get current system metrics
func (h *MonitorHandler) GetSystemMetrics(c *gin.Context) {
	metrics := h.monitorService.GetSystemMetrics()
	response.Success(c, "System metrics retrieved successfully", metrics)
}

// GetRecentMetrics handles requests to get recent system metrics for charting
func (h *MonitorHandler) GetRecentMetrics(c *gin.Context) {
	// Parse count parameter (default to 10)
	count, err := strconv.Atoi(c.DefaultQuery("count", "10"))
	if err != nil || count <= 0 {
		count = 10
	}

	// Limit count to prevent excessive data
	if count > 100 {
		count = 100
	}

	metrics := h.monitorService.GetRecentMetrics(count)
	response.Success(c, "Recent metrics retrieved successfully", metrics)
}
