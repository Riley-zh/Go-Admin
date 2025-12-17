package handler

import (
	"time"

	"go-admin/internal/service"
	"go-admin/pkg/response"

	"github.com/gin-gonic/gin"
)

// MetricsHandler handles metrics-related HTTP requests
type MetricsHandler struct {
	metricsService *service.MetricsService
}

// NewMetricsHandler creates a new metrics handler
func NewMetricsHandler() *MetricsHandler {
	return &MetricsHandler{
		metricsService: service.NewMetricsService(),
	}
}

// GetMetrics handles requests to get current system metrics
func (h *MetricsHandler) GetMetrics(c *gin.Context) {
	metrics := h.metricsService.CollectMetrics()
	response.Success(c, "Metrics collected successfully", metrics)
}

// GetHealth handles requests to get health status
func (h *MetricsHandler) GetHealth(c *gin.Context) {
	health := h.metricsService.GetHealthStatus()
	response.Success(c, "Health status retrieved successfully", health)
}

// GetHealthDetailed handles requests to get detailed health status
func (h *MetricsHandler) GetHealthDetailed(c *gin.Context) {
	health := h.metricsService.GetHealthStatus()

	// Add more detailed information
	health["uptime"] = time.Since(time.Now()) // In a real implementation, you would track actual uptime

	response.Success(c, "Detailed health status retrieved successfully", health)
}
