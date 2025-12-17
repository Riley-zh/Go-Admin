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
	Name        string `json:"name" binding:"required,min=1,max=50" example:"admin"`
	Description string `json:"description" binding:"max=255" example:"Administrator role with full access"`
}

// UpdateRoleRequest represents the update role request body
type UpdateRoleRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=50" example:"admin"`
	Description string `json:"description" binding:"max=255" example:"Administrator role with full access"`
}

// AssignRoleRequest represents the assign role request body
type AssignRoleRequest struct {
	UserID uint `json:"user_id" binding:"required" example:"1"`
	RoleID uint `json:"role_id" binding:"required" example:"1"`
}

// RemoveRoleRequest represents the remove role request body
type RemoveRoleRequest struct {
	UserID uint `json:"user_id" binding:"required" example:"1"`
	RoleID uint `json:"role_id" binding:"required" example:"1"`
}

// CreateRole godoc
// @Summary Create a new role
// @Description Create a new role with name and description
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateRoleRequest true "Role details"
// @Success 201 {object} map[string]interface{} "Role created successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 409 {object} map[string]interface{} "Conflict - Role already exists"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /roles [post]
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
// GetRoleByID godoc
// @Summary Get role by ID
// @Description Get a role by its ID
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Role ID"
// @Success 200 {object} map[string]interface{} "Role retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 404 {object} map[string]interface{} "Role not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /roles/{id} [get]
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
// UpdateRole godoc
// @Summary Update a role
// @Description Update an existing role's information
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Role ID"
// @Param request body UpdateRoleRequest true "Updated role details"
// @Success 200 {object} map[string]interface{} "Role updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 404 {object} map[string]interface{} "Role not found"
// @Failure 409 {object} map[string]interface{} "Conflict - Role name already exists"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /roles/{id} [put]
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
// DeleteRole godoc
// @Summary Delete a role
// @Description Delete a role by its ID
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Role ID"
// @Success 200 {object} map[string]interface{} "Role deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 404 {object} map[string]interface{} "Role not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /roles/{id} [delete]
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
// ListRoles godoc
// @Summary List roles
// @Description Get a list of roles with pagination
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} map[string]interface{} "Roles retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /roles [get]
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
// AssignRole godoc
// @Summary Assign role to user
// @Description Assign a role to a user
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body AssignRoleRequest true "Assign role request"
// @Success 200 {object} map[string]interface{} "Role assigned successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 404 {object} map[string]interface{} "User or role not found"
// @Failure 409 {object} map[string]interface{} "User already has this role"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /roles/assign [post]
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
// RemoveRole godoc
// @Summary Remove role from user
// @Description Remove a role from a user
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body RemoveRoleRequest true "Remove role request"
// @Success 200 {object} map[string]interface{} "Role removed successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 404 {object} map[string]interface{} "User or role not found"
// @Failure 409 {object} map[string]interface{} "User doesn't have this role"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /roles/remove [post]
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
// GetRolesByUserID godoc
// @Summary Get roles by user ID
// @Description Get all roles assigned to a user by user ID
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_id path string true "User ID"
// @Success 200 {object} map[string]interface{} "Roles retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /roles/user/{user_id} [get]
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
