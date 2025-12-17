package middleware

import (
	"bytes"
	"io"
	"time"

	"go-admin/internal/logger"

	"github.com/gin-gonic/gin"
)

// ResponseWriter is a wrapper around gin.ResponseWriter that captures response body
type ResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write captures the response body
func (w ResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Body returns the captured response body
func (w ResponseWriter) Body() []byte {
	return w.body.Bytes()
}

// RequestLoggerMiddleware provides detailed request logging
func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip logging for health check endpoints
		if c.Request.URL.Path == "/health" || c.Request.URL.Path == "/ping" {
			c.Next()
			return
		}

		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Read request body for logging (limit size to avoid memory issues)
		var requestBody []byte
		if c.Request.Body != nil && c.Request.Method != "GET" && c.Request.Method != "HEAD" {
			requestBody, _ = io.ReadAll(io.LimitReader(c.Request.Body, 1024*1024)) // Limit to 1MB
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Wrap response writer to capture response
		w := &ResponseWriter{
			ResponseWriter: c.Writer,
			body:          bytes.NewBufferString(""),
		}
		c.Writer = w

		// Process request
		c.Next()

		// Calculate metrics
		end := time.Now()
		latency := end.Sub(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := w.Status()
		responseSize := w.Size()
		userAgent := c.Request.UserAgent()

		// Combine path and query
		if raw != "" {
			path = path + "?" + raw
		}

		// Create structured logger with common fields
		log := logger.NewStructuredLogger().
			WithField("status", statusCode).
			WithField("method", method).
			WithField("path", path).
			WithField("ip", clientIP).
			WithField("latency", latency).
			WithField("userAgent", userAgent).
			WithField("responseSize", responseSize)

		// Add user ID if available
		if userID, exists := c.Get("userID"); exists {
			log = log.WithField("userID", userID)
		}

		// Add request ID if available
		if requestID, exists := c.Get("requestID"); exists {
			log = log.WithField("requestID", requestID)
		}

		// Add request body for non-GET requests (be careful with sensitive data)
		if len(requestBody) > 0 && method != "GET" && method != "HEAD" {
			// Mask sensitive data
			bodyStr := string(requestBody)
			if method == "POST" && (path == "/api/v1/login" || path == "/api/v1/register") {
				// Mask password fields
				bodyStr = maskSensitiveData(bodyStr)
			}
			log = log.WithField("requestBody", bodyStr)
		}

		// Add response body for debugging (be careful with sensitive data)
		if statusCode >= 400 && w.Body() != nil {
			responseBody := w.Body()
			if len(responseBody) > 0 && len(responseBody) < 1024 { // Only log small responses
				log = log.WithField("responseBody", string(responseBody))
			}
		}

		// Add error if any
		if len(c.Errors) > 0 {
			log = log.WithField("error", c.Errors.String())
		}

		// Log based on status code
		switch {
		case statusCode >= 500:
			log.Error("Server error")
		case statusCode >= 400:
			log.Warn("Client error")
		default:
			log.Info("Request")
		}
	}
}

// maskSensitiveData masks sensitive fields in JSON for logging
func maskSensitiveData(data string) string {
	// Simple implementation - in production, you might want a more sophisticated approach
	// This masks common password field names
	masked := data
	passwordFields := []string{"\"password\":\"", "\"pwd\":\"", "\"pass\":\""}
	
	for _, field := range passwordFields {
		start := 0
		for {
			idx := findSubstring(masked, field, start)
			if idx == -1 {
				break
			}
			
			// Find the end of the password value
			valueStart := idx + len(field)
			endQuote := findSubstring(masked, "\"", valueStart)
			if endQuote == -1 {
				break
			}
			
			// Replace the password value with asterisks
			masked = masked[:valueStart] + "*****" + masked[endQuote:]
			start = valueStart + 6 // Move past the masked password
		}
	}
	
	return masked
}

// findSubstring finds the first occurrence of substr in s starting from index start
func findSubstring(s, substr string, start int) int {
	if start >= len(s) {
		return -1
	}
	
	for i := start; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	
	return -1
}