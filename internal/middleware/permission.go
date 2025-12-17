package middleware

import (
	"net/http"
	"strings"

	"go-admin/internal/repository"
	"go-admin/internal/service"

	"github.com/gin-gonic/gin"
)

// PermissionMiddleware represents the permission middleware
type PermissionMiddleware struct {
	authService    service.AuthService
	permissionRepo repository.PermissionRepository
}

// NewPermissionMiddleware creates a new permission middleware
func NewPermissionMiddleware() *PermissionMiddleware {
	return &PermissionMiddleware{
		authService:    service.NewAuthService(),
		permissionRepo: repository.NewPermissionRepository(),
	}
}

// RequirePermission is the middleware function for permission checking
func (m *PermissionMiddleware) RequirePermission(resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context (set by JWT middleware)
		userIDValue, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		userID, ok := userIDValue.(uint)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
			c.Abort()
			return
		}

		// Check permission
		hasPermission, err := m.permissionRepo.CheckUserPermission(userID, resource, action)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check permission"})
			c.Abort()
			return
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		// Continue to next handler
		c.Next()
	}
}

// RequireRole is the middleware function for role checking
func (m *PermissionMiddleware) RequireRole(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context (set by JWT middleware)
		userIDValue, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		userID, ok := userIDValue.(uint)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
			c.Abort()
			return
		}

		// Get user roles
		roleRepo := repository.NewRoleRepository()
		userRoles, err := roleRepo.GetRolesByUserID(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user roles"})
			c.Abort()
			return
		}

		// Check if user has any of the required roles
		hasRequiredRole := false
		for _, userRole := range userRoles {
			for _, requiredRole := range requiredRoles {
				if strings.EqualFold(userRole.Name, requiredRole) {
					hasRequiredRole = true
					break
				}
			}
			if hasRequiredRole {
				break
			}
		}

		if !hasRequiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient role permissions"})
			c.Abort()
			return
		}

		// Continue to next handler
		c.Next()
	}
}
