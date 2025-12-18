package app

import (
	"time"

	"go-admin/internal/handler"
	"go-admin/pkg/api"
	"go-admin/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// initializeOptimizedAPI initializes the optimized API components
func initializeOptimizedAPI(router *gin.Engine) {
	// Create API client configuration
	apiConfig := &api.Config{
		BaseURL: "",
		DefaultHeaders: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		Timeout:       30 * time.Second,
		MaxRetries:    3,
		EnableLogging: true,
	}

	// Initialize the default API client
	api.InitDefaultAPIClient(apiConfig)

	// Add optimized API middleware
	router.Use(middleware.NewOptimizedAPIMiddleware(apiConfig).Handle())
	router.Use(middleware.APITimeoutMiddleware(30 * time.Second))
	router.Use(middleware.APIMetricsMiddleware())
	router.Use(middleware.APISecurityMiddleware())
	router.Use(middleware.APICorrelationIDMiddleware())
	router.Use(middleware.ResponseOptimizationMiddleware())
}

// registerOptimizedRoutes registers routes with optimized handlers
func registerOptimizedRoutes(router *gin.Engine) {
	// API version 1 group with optimized middleware
	v1 := router.Group("/api/v1")
	{
		// Optimized user handlers
		apiClient := api.DefaultAPIClient()
		optimizedUserHandler := handler.NewOptimizedUserHandler(apiClient)

		v1.POST("/users-optimized", optimizedUserHandler.CreateUser)
		v1.PUT("/users-optimized/:id", optimizedUserHandler.UpdateUser)

		// External API examples
		v1.GET("/external/api-example", optimizedUserHandler.ExternalAPIExample)
		v1.POST("/external/batch-example", optimizedUserHandler.BatchAPIExample)
		v1.GET("/external/streaming-example", optimizedUserHandler.StreamingAPIExample)
	}
}
