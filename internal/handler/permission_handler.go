package handler

import (
	"net/http"
	"strconv"

	"go-admin/internal/model"
	"go-admin/internal/service"
	"go-admin/pkg/errors"

	"github.com/gin-gonic/gin"
)

// PermissionHandler represents the permission handler
type PermissionHandler struct {
	permissionService service.PermissionService
}

// NewPermissionHandler creates a new permission handler
func NewPermissionHandler() *PermissionHandler {
	return &PermissionHandler{
		permissionService: service.NewPermissionService(),
	}
}

// CreatePermissionRequest represents the create permission request body
type CreatePermissionRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=50"`
	Description string `json:"description" binding:"max=255"`
	Resource    string `json:"resource" binding:"required,min=1,max=100"`
	Action      string `json:"action" binding:"required,min=1,max=50"`
}

// UpdatePermissionRequest represents the update permission request body
type UpdatePermissionRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=50"`
	Description string `json:"description" binding:"max=255"`
	Resource    string `json:"resource" binding:"required,min=1,max=100"`
	Action      string `json:"action" binding:"required,min=1,max=50"`
}

// AssignPermissionRequest represents the assign permission request body
type AssignPermissionRequest struct {
	RoleID       uint `json:"role_id" binding:"required"`
	PermissionID uint `json:"permission_id" binding:"required"`
}

// CreatePermission handles creating a new permission
func (h *PermissionHandler) CreatePermission(c *gin.Context) {
	// Validate request
	var req CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Create permission
	permission, err := h.permissionService.CreatePermission(req.Name, req.Description, req.Resource, req.Action)
	if err != nil {
		if appErr, ok := err.(*errors.Error); ok {
			c.JSON(appErr.Code, gin.H{
				"error":   appErr.Message,
				"details": appErr.Details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create permission",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Permission created successfully",
		"permission": permission,
	})
}

// GetPermissionByID handles getting a permission by ID
func (h *PermissionHandler) GetPermissionByID(c *gin.Context) {
	// Get permission ID from path parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid permission ID",
			"details": "权限ID格式不正确",
		})
		return
	}

	// Get permission
	permission, err := h.permissionService.GetPermissionByID(uint(id))
	if err != nil {
		if appErr, ok := err.(*errors.Error); ok {
			c.JSON(appErr.Code, gin.H{
				"error":   appErr.Message,
				"details": appErr.Details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get permission",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"permission": permission,
	})
}

// UpdatePermission handles updating a permission
func (h *PermissionHandler) UpdatePermission(c *gin.Context) {
	// Get permission ID from path parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid permission ID",
			"details": "权限ID格式不正确",
		})
		return
	}

	// Validate request
	var req UpdatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Create permission model
	permission := &model.Permission{
		ID:          uint(id),
		Name:        req.Name,
		Description: req.Description,
		Resource:    req.Resource,
		Action:      req.Action,
	}

	// Update permission
	err = h.permissionService.UpdatePermission(permission)
	if err != nil {
		if appErr, ok := err.(*errors.Error); ok {
			c.JSON(appErr.Code, gin.H{
				"error":   appErr.Message,
				"details": appErr.Details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update permission",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Permission updated successfully",
	})
}

// DeletePermission handles deleting a permission
func (h *PermissionHandler) DeletePermission(c *gin.Context) {
	// Get permission ID from path parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid permission ID",
			"details": "权限ID格式不正确",
		})
		return
	}

	// Delete permission
	err = h.permissionService.DeletePermission(uint(id))
	if err != nil {
		if appErr, ok := err.(*errors.Error); ok {
			c.JSON(appErr.Code, gin.H{
				"error":   appErr.Message,
				"details": appErr.Details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete permission",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Permission deleted successfully",
	})
}

// ListPermissions handles listing permissions with pagination
func (h *PermissionHandler) ListPermissions(c *gin.Context) {
	// Get pagination parameters
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// List permissions
	permissions, total, err := h.permissionService.ListPermissions(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to list permissions",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"permissions": permissions,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
	})
}

// AssignPermission handles assigning a permission to a role
func (h *PermissionHandler) AssignPermission(c *gin.Context) {
	// Validate request
	var req AssignPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Assign permission to role
	err := h.permissionService.AssignPermissionToRole(req.RoleID, req.PermissionID)
	if err != nil {
		if appErr, ok := err.(*errors.Error); ok {
			c.JSON(appErr.Code, gin.H{
				"error":   appErr.Message,
				"details": appErr.Details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to assign permission",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Permission assigned successfully",
	})
}

// RemovePermission handles removing a permission from a role
func (h *PermissionHandler) RemovePermission(c *gin.Context) {
	// Validate request
	var req AssignPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Remove permission from role
	err := h.permissionService.RemovePermissionFromRole(req.RoleID, req.PermissionID)
	if err != nil {
		if appErr, ok := err.(*errors.Error); ok {
			c.JSON(appErr.Code, gin.H{
				"error":   appErr.Message,
				"details": appErr.Details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to remove permission",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Permission removed successfully",
	})
}

// GetPermissionsByRoleID handles getting permissions by role ID
func (h *PermissionHandler) GetPermissionsByRoleID(c *gin.Context) {
	// Get role ID from path parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid role ID",
			"details": "角色ID格式不正确",
		})
		return
	}

	// Get permissions by role ID
	permissions, err := h.permissionService.GetPermissionsByRoleID(uint(id))
	if err != nil {
		if appErr, ok := err.(*errors.Error); ok {
			c.JSON(appErr.Code, gin.H{
				"error":   appErr.Message,
				"details": appErr.Details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get permissions",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"permissions": permissions,
	})
}

// GetPermissionsByUserID handles getting permissions by user ID
func (h *PermissionHandler) GetPermissionsByUserID(c *gin.Context) {
	// Get user ID from path parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid user ID",
			"details": "用户ID格式不正确",
		})
		return
	}

	// Get permissions by user ID
	permissions, err := h.permissionService.GetPermissionsByUserID(uint(id))
	if err != nil {
		if appErr, ok := err.(*errors.Error); ok {
			c.JSON(appErr.Code, gin.H{
				"error":   appErr.Message,
				"details": appErr.Details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get permissions",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"permissions": permissions,
	})
}
