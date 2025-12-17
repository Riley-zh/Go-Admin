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
	rateLimitConfig middleware.RateLimitConfig
}

// NewSecurityHandler creates a new security handler
func NewSecurityHandler() *SecurityHandler {
	return &SecurityHandler{
		csrfMiddleware: middleware.NewCSRFMiddleware(),
		rateLimitConfig: middleware.DefaultRateLimitConfig(),
	}
}

// GetCSRFToken handles requests to get a CSRF token
func (h *SecurityHandler) GetCSRFToken(c *gin.Context) {
	token, err := h.csrfMiddleware.GenerateToken()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to generate CSRF token")
		return
	}

	// Set token in cookie with enhanced security attributes
	// Secure=true to only send over HTTPS (in production)
	// HttpOnly=true to prevent XSS attacks from accessing the cookie
	isSecure := c.Request.TLS != nil // Use Secure flag if HTTPS
	c.SetCookie("csrf_token", token, 3600, "/", "", isSecure, true)

	// Note: SameSite attribute is not directly supported in Gin's SetCookie method
	// In a production environment, you might want to use a custom cookie setting method
	// that supports SameSite=Strict for enhanced CSRF protection

	response.Success(c, "CSRF token generated successfully", gin.H{
		"token": token,
	})
}

// GetRateLimitConfig handles requests to get rate limit configuration
func (h *SecurityHandler) GetRateLimitConfig(c *gin.Context) {
	response.Success(c, "Rate limit configuration retrieved successfully", gin.H{
		"requests_per_window": h.rateLimitConfig.Requests,
		"window_seconds":      h.rateLimitConfig.Window.Seconds(),
		"use_distributed_cache": h.rateLimitConfig.UseDistributedCache,
	})
}
