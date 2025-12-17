package middleware

import (
	"errors"
	"net/http"

	"go-admin/internal/logger"
	apperrors "go-admin/pkg/errors"

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

			// Log the error with improved error handling
			logger.Error("Request error",
				zap.Error(lastErr), // 使用zap.Error()而不是zap.String()
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("client_ip", c.ClientIP()),
			)

			// Check if it's our custom error type - 使用errors.As()进行类型断言
			var appErr *apperrors.Error
			if errors.As(lastErr, &appErr) {
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
