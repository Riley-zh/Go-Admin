package handler

import (
	"go-admin/internal/database"
	"go-admin/pkg/response"

	"github.com/gin-gonic/gin"
)

// DBHandler handles database-related HTTP requests
type DBHandler struct{}

// NewDBHandler creates a new database handler
func NewDBHandler() *DBHandler {
	return &DBHandler{}
}

// GetDBStats handles requests to get database connection pool statistics
func (h *DBHandler) GetDBStats(c *gin.Context) {
	stats := database.Stats()

	response.Success(c, "Database statistics retrieved successfully", gin.H{
		"max_open_connections": stats.MaxOpenConnections,
		"open_connections":     stats.OpenConnections,
		"in_use":               stats.InUse,
		"idle":                 stats.Idle,
		"wait_count":           stats.WaitCount,
		"wait_duration":        stats.WaitDuration,
		"max_idle_closed":      stats.MaxIdleClosed,
		"max_lifetime_closed":  stats.MaxLifetimeClosed,
		"max_idle_time_closed": stats.MaxIdleTimeClosed,
	})
}
