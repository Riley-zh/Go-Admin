package middleware

import (
	"net/http"
	"runtime/debug"

	"go-admin/internal/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RecoveryMiddleware represents the recovery middleware
type RecoveryMiddleware struct{}

// NewRecoveryMiddleware creates a new recovery middleware
func NewRecoveryMiddleware() *RecoveryMiddleware {
	return &RecoveryMiddleware{}
}

// Handle is the middleware function for global exception handling
func (m *RecoveryMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic with standardized format
				logger.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("stack", string(debug.Stack())),
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
					zap.String("client_ip", c.ClientIP()),
					zap.String("user_agent", c.Request.UserAgent()),
				)

				// Return error response
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Internal server error",
					"message": "An unexpected error occurred. Please try again later.",
				})

				// Abort the request
				c.Abort()
			}
		}()

		// Continue to next handler
		c.Next()
	}
}
