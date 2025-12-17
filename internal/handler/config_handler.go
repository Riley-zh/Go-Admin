package handler

import (
	"go-admin/config"

	"github.com/gin-gonic/gin"
)

// ConfigHandler handles config-related HTTP requests
type ConfigHandler struct {
	*BaseHandler
}

// NewConfigHandler creates a new config handler
func NewConfigHandler() *ConfigHandler {
	return &ConfigHandler{
		BaseHandler: NewBaseHandler(),
	}
}

// GetConfig handles requests to get current configuration
func (h *ConfigHandler) GetConfig(c *gin.Context) {
	cfg := config.Get()

	// Return configuration (excluding sensitive information)
	h.HandleSuccess(c, gin.H{
		"app": gin.H{
			"name": cfg.App.Name,
			"env":  cfg.App.Env,
			"port": cfg.App.Port,
		},
		"db": gin.H{
			"host": cfg.DB.Host,
			"port": cfg.DB.Port,
			"user": cfg.DB.User,
			"name": cfg.DB.Name,
			// Note: Password is intentionally omitted for security
		},
		"log": gin.H{
			"level":  cfg.Log.Level,
			"output": cfg.Log.Output,
		},
		"cache": gin.H{
			"maxsize":    cfg.Cache.MaxSize,
			"gcinterval": cfg.Cache.GCInterval,
		},
	})
}

// ReloadConfig handles requests to reload configuration
func (h *ConfigHandler) ReloadConfig(c *gin.Context) {
	// In this implementation, configuration is automatically reloaded when the file changes
	// This endpoint is provided for manual triggering if needed
	cfg, err := config.Load()
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, gin.H{
		"message": "Configuration reloaded successfully",
		"app": gin.H{
			"name": cfg.App.Name,
			"env":  cfg.App.Env,
			"port": cfg.App.Port,
		},
	})
}
