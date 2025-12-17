package middleware

import (
	"net/http"
	"strings"

	"go-admin/internal/logger"
	"go-admin/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// EnhancedPermissionMiddleware represents the enhanced permission middleware
type EnhancedPermissionMiddleware struct {
	permissionService service.PermissionService
}

// NewEnhancedPermissionMiddleware creates a new enhanced permission middleware
func NewEnhancedPermissionMiddleware() *EnhancedPermissionMiddleware {
	return &EnhancedPermissionMiddleware{
		permissionService: service.NewPermissionService(),
	}
}

// RequirePermission is the middleware function for permission checking
func (m *EnhancedPermissionMiddleware) RequirePermission(resource, action string) gin.HandlerFunc {
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

		// Build context for permission checking
		context := m.buildPermissionContext(c)

		// Check permission
		hasPermission, err := m.permissionService.CheckPermission(c.Request.Context(), userID, resource, action, context)
		if err != nil {
			logger.Error("Failed to check permission", zap.Error(err), zap.Uint("userID", userID), zap.String("resource", resource), zap.String("action", action))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check permission"})
			c.Abort()
			return
		}

		if !hasPermission {
			logger.Warn("Permission denied", zap.Uint("userID", userID), zap.String("resource", resource), zap.String("action", action))
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			c.Abort()
			return
		}

		// Permission granted, continue
		c.Next()
	}
}

// RequirePermissionWithContext is the middleware function for permission checking with custom context
func (m *EnhancedPermissionMiddleware) RequirePermissionWithContext(resource, action string, contextBuilder func(*gin.Context) map[string]interface{}) gin.HandlerFunc {
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

		// Build context for permission checking
		context := m.buildPermissionContext(c)
		
		// Add custom context if provided
		if contextBuilder != nil {
			customContext := contextBuilder(c)
			for k, v := range customContext {
				context[k] = v
			}
		}

		// Check permission
		hasPermission, err := m.permissionService.CheckPermission(c.Request.Context(), userID, resource, action, context)
		if err != nil {
			logger.Error("Failed to check permission", zap.Error(err), zap.Uint("userID", userID), zap.String("resource", resource), zap.String("action", action))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check permission"})
			c.Abort()
			return
		}

		if !hasPermission {
			logger.Warn("Permission denied", zap.Uint("userID", userID), zap.String("resource", resource), zap.String("action", action))
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			c.Abort()
			return
		}

		// Permission granted, continue
		c.Next()
	}
}

// RequireAnyPermission checks if user has any of the specified permissions
func (m *EnhancedPermissionMiddleware) RequireAnyPermission(permissions []PermissionRequirement) gin.HandlerFunc {
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

		// Build context for permission checking
		context := m.buildPermissionContext(c)

		// Check each permission
		for _, perm := range permissions {
			hasPermission, err := m.permissionService.CheckPermission(c.Request.Context(), userID, perm.Resource, perm.Action, context)
			if err != nil {
				logger.Error("Failed to check permission", zap.Error(err), zap.Uint("userID", userID), zap.String("resource", perm.Resource), zap.String("action", perm.Action))
				continue
			}

			if hasPermission {
				// Permission granted, continue
				c.Next()
				return
			}
		}

		// No permission granted
		logger.Warn("No matching permission found", zap.Uint("userID", userID))
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		c.Abort()
	}
}

// RequireAllPermissions checks if user has all of the specified permissions
func (m *EnhancedPermissionMiddleware) RequireAllPermissions(permissions []PermissionRequirement) gin.HandlerFunc {
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

		// Build context for permission checking
		context := m.buildPermissionContext(c)

		// Check each permission
		for _, perm := range permissions {
			hasPermission, err := m.permissionService.CheckPermission(c.Request.Context(), userID, perm.Resource, perm.Action, context)
			if err != nil {
				logger.Error("Failed to check permission", zap.Error(err), zap.Uint("userID", userID), zap.String("resource", perm.Resource), zap.String("action", perm.Action))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check permission"})
				c.Abort()
				return
			}

			if !hasPermission {
				logger.Warn("Permission denied", zap.Uint("userID", userID), zap.String("resource", perm.Resource), zap.String("action", perm.Action))
				c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
				c.Abort()
				return
			}
		}

		// All permissions granted, continue
		c.Next()
	}
}

// PermissionRequirement represents a permission requirement
type PermissionRequirement struct {
	Resource string
	Action   string
}

// buildPermissionContext builds the permission context from the request
func (m *EnhancedPermissionMiddleware) buildPermissionContext(c *gin.Context) map[string]interface{} {
	context := make(map[string]interface{})

	// Add request information
	context["method"] = c.Request.Method
	context["path"] = c.Request.URL.Path
	context["client_ip"] = c.ClientIP()
	context["user_agent"] = c.Request.UserAgent()

	// Add time information (can be used for time-based permissions)
	// context["time"] = time.Now().Format("15:04") // HH:MM format

	// Add query parameters
	if len(c.Request.URL.Query()) > 0 {
		queryParams := make(map[string]interface{})
		for key, values := range c.Request.URL.Query() {
			if len(values) == 1 {
				queryParams[key] = values[0]
			} else {
				queryParams[key] = values
			}
		}
		context["query_params"] = queryParams
	}

	// Add path parameters
	pathParams := make(map[string]interface{})
	for _, param := range c.Params {
		pathParams[param.Key] = param.Value
	}
	if len(pathParams) > 0 {
		context["path_params"] = pathParams
	}

	// Add headers (excluding sensitive ones)
	headers := make(map[string]interface{})
	for key, values := range c.Request.Header {
		keyLower := strings.ToLower(key)
		// Skip sensitive headers
		if keyLower == "authorization" || keyLower == "cookie" || keyLower == "x-api-key" {
			continue
		}
		if len(values) == 1 {
			headers[key] = values[0]
		} else {
			headers[key] = values
		}
	}
	if len(headers) > 0 {
		context["headers"] = headers
	}

	return context
}

// DataPermissionMiddleware provides data-level permission checking
type DataPermissionMiddleware struct {
	permissionService service.PermissionService
}

// NewDataPermissionMiddleware creates a new data permission middleware
func NewDataPermissionMiddleware() *DataPermissionMiddleware {
	return &DataPermissionMiddleware{
		permissionService: service.NewPermissionService(),
	}
}

// FilterDataByPermission filters data based on user permissions
func (m *DataPermissionMiddleware) FilterDataByPermission(resource, action string, dataFilter func(*gin.Context, uint) (interface{}, error)) gin.HandlerFunc {
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

		// Check permission first
		context := make(map[string]interface{})
		hasPermission, err := m.permissionService.CheckPermission(c.Request.Context(), userID, resource, action, context)
		if err != nil {
			logger.Error("Failed to check permission", zap.Error(err), zap.Uint("userID", userID), zap.String("resource", resource), zap.String("action", action))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check permission"})
			c.Abort()
			return
		}

		if !hasPermission {
			logger.Warn("Permission denied", zap.Uint("userID", userID), zap.String("resource", resource), zap.String("action", action))
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			c.Abort()
			return
		}

		// Apply data filtering
		filteredData, err := dataFilter(c, userID)
		if err != nil {
			logger.Error("Failed to filter data", zap.Error(err), zap.Uint("userID", userID))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to filter data"})
			c.Abort()
			return
		}

		// Store filtered data in context for later use
		c.Set("filteredData", filteredData)
		c.Next()
	}
}