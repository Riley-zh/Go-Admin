package handler

import (
	"go-admin/internal/model"
	"go-admin/internal/service"
	"go-admin/pkg/errors"

	"github.com/gin-gonic/gin"
)

// UserHandler represents the user handler
type UserHandler struct {
	*BaseHandler
	userService service.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler() *UserHandler {
	return &UserHandler{
		BaseHandler: NewBaseHandler(),
		userService: service.NewUserService(),
	}
}

// CreateUserRequest represents the create user request body
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Nickname string `json:"nickname" binding:"max=100"`
}

// UpdateUserRequest represents the update user request body
type UpdateUserRequest struct {
	Email    string `json:"email" binding:"omitempty,email"`
	Nickname string `json:"nickname" binding:"max=100"`
}

// ChangePasswordRequest represents the change password request body
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=6,max=50"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=50"`
}

// CreateUser handles creating a new user
func (h *UserHandler) CreateUser(c *gin.Context) {
	// Validate request
	var req CreateUserRequest
	if !h.BindAndValidate(c, &req) {
		return
	}

	// Create user
	user, err := h.userService.CreateUser(req.Username, req.Password, req.Email, req.Nickname)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleCreated(c, "User created successfully", gin.H{"user": user})
}

// GetUserByID handles getting a user by ID
func (h *UserHandler) GetUserByID(c *gin.Context) {
	// Get user ID from path parameter
	id, err := h.ParseIDParam(c, "id")
	if err != nil {
		h.HandleValidationError(c, err)
		return
	}

	// Get user
	user, err := h.userService.GetUserByID(id)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, gin.H{"user": user})
}

// UpdateUser handles updating a user
func (h *UserHandler) UpdateUser(c *gin.Context) {
	// Get user ID from path parameter
	id, err := h.ParseIDParam(c, "id")
	if err != nil {
		h.HandleValidationError(c, err)
		return
	}

	// Validate request
	var req UpdateUserRequest
	if !h.BindAndValidate(c, &req) {
		return
	}

	// Create user model
	user := &model.User{
		ID:       id,
		Email:    req.Email,
		Nickname: req.Nickname,
	}

	// Update user
	err = h.userService.UpdateUser(user)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccessWithMessage(c, "User updated successfully", nil)
}

// DeleteUser handles deleting a user
func (h *UserHandler) DeleteUser(c *gin.Context) {
	// Get user ID from path parameter
	id, err := h.ParseIDParam(c, "id")
	if err != nil {
		h.HandleValidationError(c, err)
		return
	}

	// Delete user
	err = h.userService.DeleteUser(id)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccessWithMessage(c, "User deleted successfully", nil)
}

// ListUsers handles listing users with pagination
func (h *UserHandler) ListUsers(c *gin.Context) {
	// Get pagination parameters
	params := h.GetPaginationParams(c)

	// Check if roles should be included
	includeRoles := c.DefaultQuery("include_roles", "false") == "true"

	var users interface{}
	var total int64
	var err error

	if includeRoles {
		// List users with roles to prevent N+1 query problem
		users, total, err = h.userService.ListUsersWithRoles(params.Page, params.PageSize)
	} else {
		// List users without roles (original behavior)
		users, total, err = h.userService.ListUsers(params.Page, params.PageSize)
	}

	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandlePaginationResponse(c, gin.H{"users": users}, total, params)
}

// ChangePassword handles changing user password
func (h *UserHandler) ChangePassword(c *gin.Context) {
	// Get user ID from context (set by JWT middleware)
	userIDValue, exists := c.Get("userID")
	if !exists {
		h.HandleError(c, errors.Unauthorized("User not authenticated", "用户未认证"))
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		h.HandleError(c, errors.InternalServerError("Invalid user ID", "无效的用户ID"))
		return
	}

	// Validate request
	var req ChangePasswordRequest
	if !h.BindAndValidate(c, &req) {
		return
	}

	// Change password
	err := h.userService.ChangePassword(userID, req.OldPassword, req.NewPassword)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccessWithMessage(c, "Password changed successfully", nil)
}
