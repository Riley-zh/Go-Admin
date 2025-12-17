package middleware

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CSRFMiddleware represents the CSRF protection middleware
type CSRFMiddleware struct {
	tokenLength int
	expiration  time.Duration
}

// NewCSRFMiddleware creates a new CSRF middleware
func NewCSRFMiddleware() *CSRFMiddleware {
	return &CSRFMiddleware{
		tokenLength: 32,
		expiration:  1 * time.Hour,
	}
}

// GenerateToken generates a new CSRF token
func (m *CSRFMiddleware) GenerateToken() (string, error) {
	bytes := make([]byte, m.tokenLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// ValidateToken validates a CSRF token
func (m *CSRFMiddleware) ValidateToken(token, expectedToken string) bool {
	// Constant time comparison to prevent timing attacks
	return subtle.ConstantTimeCompare([]byte(token), []byte(expectedToken)) == 1
}

// Protect is the middleware function for CSRF protection
func (m *CSRFMiddleware) Protect() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip CSRF protection for GET, HEAD, OPTIONS requests
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		// Get token from header
		token := c.GetHeader("X-CSRF-Token")
		if token == "" {
			// If not in header, try to get from form
			token = c.PostForm("csrf_token")
		}

		// Get expected token from cookie
		expectedToken, err := c.Cookie("csrf_token")
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "CSRF token missing"})
			c.Abort()
			return
		}

		// Validate token
		if !m.ValidateToken(token, expectedToken) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid CSRF token"})
			c.Abort()
			return
		}

		// Continue to next handler
		c.Next()
	}
}

// GetCSRFToken returns a new CSRF token
func (m *CSRFMiddleware) GetCSRFToken(c *gin.Context) {
	token, err := m.GenerateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate CSRF token"})
		return
	}

	// Set token in cookie
	c.SetCookie("csrf_token", token, int(m.expiration.Seconds()), "/", "", false, true)

	// Return token in response
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
