package handler

import (
	"net/http"

	"go-admin/internal/middleware"
	"go-admin/pkg/response"

	"github.com/gin-gonic/gin"
)

// SecurityHandler handles security-related HTTP requests
type SecurityHandler struct {
	csrfMiddleware *middleware.CSRFMiddleware
}

// NewSecurityHandler creates a new security handler
func NewSecurityHandler() *SecurityHandler {
	return &SecurityHandler{
		csrfMiddleware: middleware.NewCSRFMiddleware(),
	}
}

// GetCSRFToken handles requests to get a CSRF token
func (h *SecurityHandler) GetCSRFToken(c *gin.Context) {
	token, err := h.csrfMiddleware.GenerateToken()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to generate CSRF token")
		return
	}

	// Set token in cookie
	c.SetCookie("csrf_token", token, 3600, "/", "", false, true)

	response.Success(c, "CSRF token generated successfully", gin.H{
		"token": token,
	})
}

// GetRateLimitConfig handles requests to get rate limit configuration
func (h *SecurityHandler) GetRateLimitConfig(c *gin.Context) {
	// For now, we'll return a static configuration
	// In a real application, this would come from configuration
	response.Success(c, "Rate limit configuration retrieved successfully", gin.H{
		"requests_per_minute": 60,
		"burst":               10,
	})
}
