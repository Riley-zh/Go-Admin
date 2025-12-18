package middleware

import (
	"bytes"
	"context"
	"io"
	"time"

	"go-admin/pkg/api"
	"go-admin/pkg/jsonutils"
	"go-admin/pkg/validation"

	"github.com/gin-gonic/gin"
)

// OptimizedAPIMiddleware provides optimized API handling with efficient JSON processing
type OptimizedAPIMiddleware struct {
	apiClient *api.APIClient
	validator *validation.ValidationMiddleware
}

// NewOptimizedAPIMiddleware creates a new optimized API middleware
func NewOptimizedAPIMiddleware(config *api.Config) *OptimizedAPIMiddleware {
	apiClient := api.NewAPIClient(config)
	validator := validation.NewValidationMiddleware()

	return &OptimizedAPIMiddleware{
		apiClient: apiClient,
		validator: validator,
	}
}

// Handle creates the middleware handler function
func (m *OptimizedAPIMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Store the API client in the context for use in handlers
		c.Set("api_client", m.apiClient)
		c.Set("validator", m.validator)

		// Optimize JSON processing
		if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			if err := m.optimizeJSONProcessing(c); err != nil {
				c.JSON(400, gin.H{
					"error":   "Invalid JSON format",
					"details": err.Error(),
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// optimizeJSONProcessing optimizes JSON processing for requests
func (m *OptimizedAPIMiddleware) optimizeJSONProcessing(c *gin.Context) error {
	// Read the request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}

	// Validate JSON format
	if len(body) > 0 {
		if err := jsonutils.ValidateJSON(body); err != nil {
			return err
		}

		// Restore the request body for further processing
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	}

	return nil
}

// GetAPIClient retrieves the API client from the context
func GetAPIClient(c *gin.Context) *api.APIClient {
	if client, exists := c.Get("api_client"); exists {
		return client.(*api.APIClient)
	}
	return api.DefaultAPIClient()
}

// GetValidator retrieves the validator from the context
func GetValidator(c *gin.Context) *validation.ValidationMiddleware {
	if validator, exists := c.Get("validator"); exists {
		return validator.(*validation.ValidationMiddleware)
	}
	return validation.NewValidationMiddleware()
}

// RequestValidationMiddleware provides request validation
func RequestValidationMiddleware(validationTarget interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		validator := GetValidator(c)

		// Bind and validate the request
		if err := c.ShouldBindJSON(validationTarget); err != nil {
			c.JSON(400, gin.H{
				"error":   "Invalid request format",
				"details": err.Error(),
			})
			c.Abort()
			return
		}

		if err := validator.ValidateStruct(validationTarget); err != nil {
			c.JSON(400, gin.H{
				"error":   "Validation failed",
				"details": err.Error(),
			})
			c.Abort()
			return
		}

		// Store the validated object in context
		c.Set("validated", validationTarget)
		c.Next()
	}
}

// ResponseOptimizationMiddleware optimizes JSON responses
func ResponseOptimizationMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Create a response writer to intercept the response
		writer := &responseWriter{
			ResponseWriter: c.Writer,
			buffer:         bytes.NewBuffer(nil),
		}

		c.Writer = writer
		c.Next()

		// Optimize the response if it's JSON
		if writer.isJSON() {
			optimizeJSONResponse(writer)
		}
	})
}

// responseWriter is a custom response writer to intercept responses
type responseWriter struct {
	gin.ResponseWriter
	buffer *bytes.Buffer
}

// Write intercepts the write operation to buffer the response
func (w *responseWriter) Write(data []byte) (int, error) {
	w.buffer.Write(data)
	return w.ResponseWriter.Write(data)
}

// isJSON checks if the response is JSON
func (w *responseWriter) isJSON() bool {
	contentType := w.Header().Get("Content-Type")
	return contentType == "application/json" || contentType == "application/json; charset=utf-8"
}

// optimizeJSONResponse optimizes the JSON response
func optimizeJSONResponse(writer *responseWriter) {
	// This is a placeholder for JSON response optimization
	// In a real implementation, you might:
	// 1. Minify the JSON
	// 2. Apply compression
	// 3. Validate the JSON format
	// 4. Add caching headers
}

// APITimeoutMiddleware adds timeout handling for API requests
func APITimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// APIRateLimitMiddleware provides rate limiting for API requests
func APIRateLimitMiddleware(requests int, window time.Duration) gin.HandlerFunc {
	// This would integrate with the existing rate limiter
	return func(c *gin.Context) {
		// Placeholder for API-specific rate limiting
		c.Next()
	}
}

// APIMetricsMiddleware provides metrics collection for API requests
func APIMetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		// Record metrics
		duration := time.Since(start)
		status := c.Writer.Status()

		// This would integrate with the existing metrics collector
		_ = path
		_ = method
		_ = duration
		_ = status
	}
}

// APIResponseCacheMiddleware provides response caching for API requests
func APIResponseCacheMiddleware(duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only cache GET requests
		if c.Request.Method != "GET" {
			c.Next()
			return
		}

		// Check if response is already cached
		// This would integrate with the existing cache system

		c.Next()

		// Cache the response if successful
		if c.Writer.Status() == 200 {
			// Cache logic here
		}
	}
}

// APICompressionMiddleware provides response compression
func APICompressionMiddleware() gin.HandlerFunc {
	// This would integrate with the existing compression middleware
	return func(c *gin.Context) {
		c.Next()
	}
}

// APISecurityMiddleware provides security headers for API responses
func APISecurityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Add security headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		c.Next()
	}
}

// APICorrelationIDMiddleware adds correlation IDs to API requests
func APICorrelationIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		correlationID := c.GetHeader("X-Correlation-ID")
		if correlationID == "" {
			correlationID = generateCorrelationID()
		}

		c.Header("X-Correlation-ID", correlationID)
		c.Set("correlation_id", correlationID)

		c.Next()
	}
}

// generateCorrelationID generates a correlation ID for tracking requests
func generateCorrelationID() string {
	// This would generate a unique correlation ID
	// For now, just use a timestamp-based ID
	return time.Now().Format("20060102150405.000000000")
}
