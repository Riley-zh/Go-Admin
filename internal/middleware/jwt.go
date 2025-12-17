package middleware

import (
	"net/http"
	"strings"

	"go-admin/internal/cache"
	"go-admin/internal/service"

	"github.com/gin-gonic/gin"
)

// JWTMiddleware represents the JWT middleware
type JWTMiddleware struct {
	authService service.AuthService
}

// NewJWTMiddleware creates a new JWT middleware
func NewJWTMiddleware() *JWTMiddleware {
	return &JWTMiddleware{
		authService: service.NewAuthService(),
	}
}

// Handle is the middleware function for JWT authentication
func (m *JWTMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Check if header has Bearer prefix
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		// Extract token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Check if token is in blacklist
		cacheInstance := cache.GetInstance()
		if _, exists := cacheInstance.Get("blacklist:" + tokenString); exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is invalid"})
			c.Abort()
			return
		}

		// Validate token
		user, err := m.authService.GetUserByToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Set user in context
		c.Set("user", user)
		c.Set("userID", user.ID)
		c.Set("username", user.Username)

		// Continue to next handler
		c.Next()
	}
}
