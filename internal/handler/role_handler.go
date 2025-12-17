package handler

import (
	"go-admin/internal/model"
	"go-admin/internal/service"
	"go-admin/pkg/common"

	"github.com/gin-gonic/gin"
)

// RoleHandler represents the role handler
type RoleHandler struct {
	*BaseHandler
	roleService service.RoleService
}

// NewRoleHandler creates a new role handler
func NewRoleHandler() *RoleHandler {
	return &RoleHandler{
		BaseHandler: NewBaseHandler(),
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
	if !h.BindAndValidate(c, &req) {
		return
	}

	// Create role
	role, err := h.roleService.CreateRole(req.Name, req.Description)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleCreated(c, "Role created successfully", gin.H{"role": role})
}

// GetRoleByID handles getting a role by ID
func (h *RoleHandler) GetRoleByID(c *gin.Context) {
	// Get role ID from path parameter
	id, err := h.ParseIDParam(c, "id")
	if err != nil {
		h.HandleValidationError(c, err)
		return
	}

	// Get role
	role, err := h.roleService.GetRoleByID(id)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, gin.H{"role": role})
}

// UpdateRole handles updating a role
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	// Get role ID from path parameter
	id, err := h.ParseIDParam(c, "id")
	if err != nil {
		h.HandleValidationError(c, err)
		return
	}

	// Validate request
	var req UpdateRoleRequest
	if !h.BindAndValidate(c, &req) {
		return
	}

	// Create role model
	role := &model.Role{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
	}

	// Update role
	err = h.roleService.UpdateRole(role)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccessWithMessage(c, "Role updated successfully", nil)
}

// DeleteRole handles deleting a role
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	// Get role ID from path parameter
	id, err := h.ParseIDParam(c, "id")
	if err != nil {
		h.HandleValidationError(c, err)
		return
	}

	// Delete role
	err = h.roleService.DeleteRole(id)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleDeleted(c, "Role deleted successfully")
}

// ListRoles handles listing roles with pagination
func (h *RoleHandler) ListRoles(c *gin.Context) {
	// Get pagination parameters
	pagination := common.GetPaginationParams(c)

	// List roles
	roles, total, err := h.roleService.ListRoles(pagination.Page, pagination.PageSize)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, gin.H{
		"roles": roles,
		"pagination": gin.H{
			"page":        pagination.Page,
			"page_size":   pagination.PageSize,
			"total":       total,
			"total_pages": (total + int64(pagination.PageSize) - 1) / int64(pagination.PageSize),
		},
	})
}

// AssignRole handles assigning a role to a user
func (h *RoleHandler) AssignRole(c *gin.Context) {
	// Validate request
	var req AssignRoleRequest
	if !h.BindAndValidate(c, &req) {
		return
	}

	// Assign role to user
	err := h.roleService.AssignRoleToUser(req.UserID, req.RoleID)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, gin.H{"message": "Role assigned successfully"})
}

// RemoveRole handles removing a role from a user
func (h *RoleHandler) RemoveRole(c *gin.Context) {
	// Validate request
	var req AssignRoleRequest
	if !h.BindAndValidate(c, &req) {
		return
	}

	// Remove role from user
	err := h.roleService.RemoveRoleFromUser(req.UserID, req.RoleID)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, gin.H{"message": "Role removed successfully"})
}

// GetRolesByUserID handles getting roles by user ID
func (h *RoleHandler) GetRolesByUserID(c *gin.Context) {
	// Get user ID from path parameter
	userID, err := h.ParseIDParam(c, "id")
	if err != nil {
		h.HandleValidationError(c, err)
		return
	}

	// Get roles
	roles, err := h.roleService.GetRolesByUserID(userID)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, gin.H{"roles": roles})
}
