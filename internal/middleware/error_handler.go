package middleware

import (
	"net/http"

	"go-admin/internal/logger"
	"go-admin/pkg/errors"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ErrorHandlerMiddleware represents the error handler middleware
type ErrorHandlerMiddleware struct{}

// NewErrorHandlerMiddleware creates a new error handler middleware
func NewErrorHandlerMiddleware() *ErrorHandlerMiddleware {
	return &ErrorHandlerMiddleware{}
}

// Handle is the middleware function for unified error handling
func (m *ErrorHandlerMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request
		c.Next()

		// Check for errors
		if len(c.Errors) > 0 {
			// Get the last error
			lastErr := c.Errors.Last()

			// Log the error
			logger.Error("Request error",
				zap.String("error", lastErr.Error()),
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
			)

			// Check if it's our custom error type
			if appErr, ok := lastErr.Err.(*errors.Error); ok {
				c.JSON(appErr.Code, gin.H{
					"error":   appErr.Message,
					"details": appErr.Details,
				})
				return
			}

			// Handle other types of errors
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Internal server error",
				"details": lastErr.Error(),
			})
			return
		}
	}
}
