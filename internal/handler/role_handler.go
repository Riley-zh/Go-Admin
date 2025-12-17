package handler

import (
	"net/http"
	"strconv"

	"go-admin/internal/model"
	"go-admin/internal/service"
	"go-admin/pkg/errors"

	"github.com/gin-gonic/gin"
)

// RoleHandler represents the role handler
type RoleHandler struct {
	roleService service.RoleService
}

// NewRoleHandler creates a new role handler
func NewRoleHandler() *RoleHandler {
	return &RoleHandler{
		roleService: service.NewRoleService(),
	}
}

// CreateRoleRequest represents the create role request body
type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=50"`
	Description string `json:"description" binding:"max=255"`
}

// UpdateRoleRequest represents the update role request body
type UpdateRoleRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=50"`
	Description string `json:"description" binding:"max=255"`
}

// AssignRoleRequest represents the assign role request body
type AssignRoleRequest struct {
	UserID uint `json:"user_id" binding:"required"`
	RoleID uint `json:"role_id" binding:"required"`
}

// CreateRole handles creating a new role
func (h *RoleHandler) CreateRole(c *gin.Context) {
	// Validate request
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Create role
	role, err := h.roleService.CreateRole(req.Name, req.Description)
	if err != nil {
		if appErr, ok := err.(*errors.Error); ok {
			c.JSON(appErr.Code, gin.H{
				"error":   appErr.Message,
				"details": appErr.Details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create role",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Role created successfully",
		"role":    role,
	})
}

// GetRoleByID handles getting a role by ID
func (h *RoleHandler) GetRoleByID(c *gin.Context) {
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

	// Get role
	role, err := h.roleService.GetRoleByID(uint(id))
	if err != nil {
		if appErr, ok := err.(*errors.Error); ok {
			c.JSON(appErr.Code, gin.H{
				"error":   appErr.Message,
				"details": appErr.Details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get role",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"role": role,
	})
}

// UpdateRole handles updating a role
func (h *RoleHandler) UpdateRole(c *gin.Context) {
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

	// Validate request
	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Create role model
	role := &model.Role{
		ID:          uint(id),
		Name:        req.Name,
		Description: req.Description,
	}

	// Update role
	err = h.roleService.UpdateRole(role)
	if err != nil {
		if appErr, ok := err.(*errors.Error); ok {
			c.JSON(appErr.Code, gin.H{
				"error":   appErr.Message,
				"details": appErr.Details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update role",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Role updated successfully",
	})
}

// DeleteRole handles deleting a role
func (h *RoleHandler) DeleteRole(c *gin.Context) {
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

	// Delete role
	err = h.roleService.DeleteRole(uint(id))
	if err != nil {
		if appErr, ok := err.(*errors.Error); ok {
			c.JSON(appErr.Code, gin.H{
				"error":   appErr.Message,
				"details": appErr.Details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete role",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Role deleted successfully",
	})
}

// ListRoles handles listing roles with pagination
func (h *RoleHandler) ListRoles(c *gin.Context) {
	// Get pagination parameters
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// List roles
	roles, total, err := h.roleService.ListRoles(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to list roles",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"roles":     roles,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// AssignRole handles assigning a role to a user
func (h *RoleHandler) AssignRole(c *gin.Context) {
	// Validate request
	var req AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Assign role to user
	err := h.roleService.AssignRoleToUser(req.UserID, req.RoleID)
	if err != nil {
		if appErr, ok := err.(*errors.Error); ok {
			c.JSON(appErr.Code, gin.H{
				"error":   appErr.Message,
				"details": appErr.Details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to assign role",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Role assigned successfully",
	})
}

// RemoveRole handles removing a role from a user
func (h *RoleHandler) RemoveRole(c *gin.Context) {
	// Validate request
	var req AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Remove role from user
	err := h.roleService.RemoveRoleFromUser(req.UserID, req.RoleID)
	if err != nil {
		if appErr, ok := err.(*errors.Error); ok {
			c.JSON(appErr.Code, gin.H{
				"error":   appErr.Message,
				"details": appErr.Details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to remove role",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Role removed successfully",
	})
}

// GetRolesByUserID handles getting roles by user ID
func (h *RoleHandler) GetRolesByUserID(c *gin.Context) {
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

	// Get roles by user ID
	roles, err := h.roleService.GetRolesByUserID(uint(id))
	if err != nil {
		if appErr, ok := err.(*errors.Error); ok {
			c.JSON(appErr.Code, gin.H{
				"error":   appErr.Message,
				"details": appErr.Details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get roles",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"roles": roles,
	})
}
