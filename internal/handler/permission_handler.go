package handler

import (
	"go-admin/internal/model"
	"go-admin/internal/service"
	"go-admin/pkg/response"

	"github.com/gin-gonic/gin"
)

// PermissionHandler represents the permission handler
type PermissionHandler struct {
	*BaseHandler
	permissionService service.PermissionService
}

// NewPermissionHandler creates a new permission handler
func NewPermissionHandler() *PermissionHandler {
	return &PermissionHandler{
		BaseHandler:       NewBaseHandler(),
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
	if !h.BindAndValidate(c, &req) {
		return
	}

	// Create permission
	permission, err := h.permissionService.CreatePermission(req.Name, req.Description, req.Resource, req.Action)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleCreated(c, "Permission created successfully", gin.H{"permission": permission})
}

// GetPermissionByID handles getting a permission by ID
func (h *PermissionHandler) GetPermissionByID(c *gin.Context) {
	// Get permission ID from path parameter
	id, err := h.ParseIDParam(c, "id")
	if err != nil {
		h.HandleValidationError(c, err)
		return
	}

	// Get permission
	permission, err := h.permissionService.GetPermissionByID(id)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, gin.H{"permission": permission})
}

// UpdatePermission handles updating a permission
func (h *PermissionHandler) UpdatePermission(c *gin.Context) {
	// Get permission ID from path parameter
	id, err := h.ParseIDParam(c, "id")
	if err != nil {
		h.HandleValidationError(c, err)
		return
	}

	// Validate request
	var req UpdatePermissionRequest
	if !h.BindAndValidate(c, &req) {
		return
	}

	// Create permission model
	permission := &model.Permission{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Resource:    req.Resource,
		Action:      req.Action,
	}

	// Update permission
	err = h.permissionService.UpdatePermission(permission)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, gin.H{"message": "Permission updated successfully"})
}

// DeletePermission handles deleting a permission
func (h *PermissionHandler) DeletePermission(c *gin.Context) {
	// Get permission ID from path parameter
	id, err := h.ParseIDParam(c, "id")
	if err != nil {
		h.HandleValidationError(c, err)
		return
	}

	// Delete permission
	err = h.permissionService.DeletePermission(id)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleDeleted(c, "Permission deleted successfully")
}

// ListPermissions handles listing permissions with pagination
func (h *PermissionHandler) ListPermissions(c *gin.Context) {
	// Get pagination parameters
	pagination := response.GetPaginationParams(c)

	// List permissions
	permissions, total, err := h.permissionService.ListPermissions(pagination.Page, pagination.PageSize)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, gin.H{
		"permissions": permissions,
		"pagination": gin.H{
			"page":        pagination.Page,
			"page_size":   pagination.PageSize,
			"total":       total,
			"total_pages": (total + int64(pagination.PageSize) - 1) / int64(pagination.PageSize),
		},
	})
}

// AssignPermission handles assigning a permission to a role
func (h *PermissionHandler) AssignPermission(c *gin.Context) {
	// Validate request
	var req AssignPermissionRequest
	if !h.BindAndValidate(c, &req) {
		return
	}

	// Assign permission to role
	err := h.permissionService.AssignPermissionToRole(req.RoleID, req.PermissionID)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, gin.H{"message": "Permission assigned successfully"})
}

// RemovePermission handles removing a permission from a role
func (h *PermissionHandler) RemovePermission(c *gin.Context) {
	// Validate request
	var req AssignPermissionRequest
	if !h.BindAndValidate(c, &req) {
		return
	}

	// Remove permission from role
	err := h.permissionService.RemovePermissionFromRole(req.RoleID, req.PermissionID)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, gin.H{"message": "Permission removed successfully"})
}

// GetPermissionsByRoleID handles getting permissions by role ID
func (h *PermissionHandler) GetPermissionsByRoleID(c *gin.Context) {
	// Get role ID from path parameter
	roleID, err := h.ParseIDParam(c, "id")
	if err != nil {
		h.HandleValidationError(c, err)
		return
	}

	// Get permissions by role ID
	permissions, err := h.permissionService.GetPermissionsByRoleID(roleID)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, gin.H{"permissions": permissions})
}

// GetPermissionsByUserID handles getting permissions by user ID
func (h *PermissionHandler) GetPermissionsByUserID(c *gin.Context) {
	// Get user ID from path parameter
	userID, err := h.ParseIDParam(c, "id")
	if err != nil {
		h.HandleValidationError(c, err)
		return
	}

	// Get permissions by user ID
	permissions, err := h.permissionService.GetPermissionsByUserID(userID)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, gin.H{"permissions": permissions})
}
