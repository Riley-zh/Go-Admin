package handler

import (
	"net/http"
	"strconv"

	"go-admin/internal/model"
	"go-admin/internal/service"
	"go-admin/pkg/errors"

	"github.com/gin-gonic/gin"
)

// UserHandler represents the user handler
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler() *UserHandler {
	return &UserHandler{
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
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Create user
	user, err := h.userService.CreateUser(req.Username, req.Password, req.Email, req.Nickname)
	if err != nil {
		if appErr, ok := err.(*errors.Error); ok {
			c.JSON(appErr.Code, gin.H{
				"error":   appErr.Message,
				"details": appErr.Details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create user",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    user,
	})
}

// GetUserByID handles getting a user by ID
func (h *UserHandler) GetUserByID(c *gin.Context) {
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

	// Get user
	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		if appErr, ok := err.(*errors.Error); ok {
			c.JSON(appErr.Code, gin.H{
				"error":   appErr.Message,
				"details": appErr.Details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get user",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// UpdateUser handles updating a user
func (h *UserHandler) UpdateUser(c *gin.Context) {
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

	// Validate request
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Create user model
	user := &model.User{
		ID:       uint(id),
		Email:    req.Email,
		Nickname: req.Nickname,
	}

	// Update user
	err = h.userService.UpdateUser(user)
	if err != nil {
		if appErr, ok := err.(*errors.Error); ok {
			c.JSON(appErr.Code, gin.H{
				"error":   appErr.Message,
				"details": appErr.Details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update user",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
	})
}

// DeleteUser handles deleting a user
func (h *UserHandler) DeleteUser(c *gin.Context) {
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

	// Delete user
	err = h.userService.DeleteUser(uint(id))
	if err != nil {
		if appErr, ok := err.(*errors.Error); ok {
			c.JSON(appErr.Code, gin.H{
				"error":   appErr.Message,
				"details": appErr.Details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete user",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}

// ListUsers handles listing users with pagination
func (h *UserHandler) ListUsers(c *gin.Context) {
	// Get pagination parameters
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// Check if roles should be included
	includeRoles := c.DefaultQuery("include_roles", "false") == "true"

	var users interface{}
	var total int64

	if includeRoles {
		// List users with roles to prevent N+1 query problem
		users, total, err = h.userService.ListUsersWithRoles(page, pageSize)
	} else {
		// List users without roles (original behavior)
		users, total, err = h.userService.ListUsers(page, pageSize)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to list users",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users":     users,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// ChangePassword handles changing user password
func (h *UserHandler) ChangePassword(c *gin.Context) {
	// Get user ID from context (set by JWT middleware)
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	// Validate request
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Change password
	err := h.userService.ChangePassword(userID, req.OldPassword, req.NewPassword)
	if err != nil {
		if appErr, ok := err.(*errors.Error); ok {
			c.JSON(appErr.Code, gin.H{
				"error":   appErr.Message,
				"details": appErr.Details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to change password",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password changed successfully",
	})
}
