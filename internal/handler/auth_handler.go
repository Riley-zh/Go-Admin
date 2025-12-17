package handler

import (
	"go-admin/internal/service"
	"go-admin/pkg/errors"

	"github.com/gin-gonic/gin"
)

// AuthHandler represents the auth handler
type AuthHandler struct {
	*BaseHandler
	authService service.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		BaseHandler: NewBaseHandler(),
		authService: service.NewAuthService(),
	}
}

// RegisterRequest represents the register request body
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Nickname string `json:"nickname" binding:"max=100"`
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	// Validate request
	var req RegisterRequest
	if !h.BindAndValidate(c, &req) {
		return
	}

	// Register user
	user, err := h.authService.Register(req.Username, req.Password, req.Email, req.Nickname)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleCreated(c, "User registered successfully", gin.H{"user": user})
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	// Validate request
	var req LoginRequest
	if !h.BindAndValidate(c, &req) {
		return
	}

	// Authenticate user
	token, user, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, gin.H{
		"message": "Login successful",
		"token":   token,
		"user":    user,
	})
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// Get token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		h.HandleError(c, errors.Unauthorized("Authorization header is required", ""))
		return
	}

	// Extract token
	tokenString := authHeader[len("Bearer "):]

	// Logout user
	err := h.authService.Logout(tokenString)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, gin.H{"message": "Logout successful"})
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// Get token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		h.HandleError(c, errors.Unauthorized("Authorization header is required", ""))
		return
	}

	// Extract token
	tokenString := authHeader[len("Bearer "):]

	// Refresh token
	newToken, err := h.authService.RefreshToken(tokenString)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, gin.H{
		"message": "Token refreshed successfully",
		"token":   newToken,
	})
}
