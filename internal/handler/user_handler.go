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
	Username string `json:"username" binding:"required,min=3,max=50" example:"johndoe"`
	Password string `json:"password" binding:"required,min=6,max=50" example:"password123"`
	Email    string `json:"email" binding:"required,email" example:"johndoe@example.com"`
	Nickname string `json:"nickname" binding:"max=100" example:"John Doe"`
}

// UpdateUserRequest represents the update user request body
type UpdateUserRequest struct {
	Email    string `json:"email" binding:"omitempty,email" example:"johndoe@example.com"`
	Nickname string `json:"nickname" binding:"max=100" example:"John Doe"`
}

// ChangePasswordRequest represents the change password request body
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=6,max=50" example:"oldpassword123"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=50" example:"newpassword123"`
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user account (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateUserRequest true "User details"
// @Success 201 {object} map[string]interface{} "User created successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 409 {object} map[string]interface{} "Conflict - User already exists"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /users [post]
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
// GetUserByID godoc
// @Summary Get user by ID
// @Description Retrieve a user's information by their ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{} "User retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /users/{id} [get]
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

// UpdateUser godoc
// @Summary Update a user
// @Description Update a user's information
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param request body UpdateUserRequest true "Updated user information"
// @Success 200 {object} map[string]interface{} "User updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /users/{id} [put]
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
// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{} "User deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /users/{id} [delete]
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
// ListUsers godoc
// @Summary List users
// @Description Get a list of users with pagination
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param include_roles query bool false "Include user roles" default(false)
// @Success 200 {object} map[string]interface{} "Users retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /users [get]
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
// ChangePassword godoc
// @Summary Change user password
// @Description Change the password of the authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ChangePasswordRequest true "Password change details"
// @Success 200 {object} map[string]interface{} "Password changed successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /users/change-password [put]
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
