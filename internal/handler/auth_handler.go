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
	Username string `json:"username" binding:"required,min=3,max=50" example:"johndoe"`
	Password string `json:"password" binding:"required,min=6,max=50" example:"password123"`
	Email    string `json:"email" binding:"required,email" example:"johndoe@example.com"`
	Nickname string `json:"nickname" binding:"max=100" example:"John Doe"`
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"johndoe"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account with username, password and email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration details"
// @Success 201 {object} map[string]interface{} "User registered successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 409 {object} map[string]interface{} "Conflict - User already exists"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /auth/register [post]
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

// Login godoc
// @Summary User login
// @Description Authenticate a user with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} map[string]interface{} "Login successful"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	// Validate request
	var req LoginRequest
	if !h.BindAndValidate(c, &req) {
		return
	}

	// Authenticate user
	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	token, user, err := h.authService.Login(req.Username, req.Password, clientIP, userAgent)
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

// Logout godoc
// @Summary User logout
// @Description Logout a user and invalidate the JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Logout successful"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /auth/logout [post]
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

// RefreshToken godoc
// @Summary Refresh JWT token
// @Description Refresh an existing JWT token to extend its validity
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Token refreshed successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /auth/refresh [post]
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
	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	newToken, err := h.authService.RefreshToken(tokenString, clientIP, userAgent)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, gin.H{
		"message": "Token refreshed successfully",
		"token":   newToken,
	})
}
