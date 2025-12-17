package handler

import (
	"net/http"
	"strconv"
	"time"

	"go-admin/internal/service"
	"go-admin/pkg/response"

	"github.com/gin-gonic/gin"
)

// NotificationHandler handles notification-related HTTP requests
type NotificationHandler struct {
	notificationService *service.NotificationService
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{
		notificationService: service.NewNotificationService(),
	}
}

// CreateNotification handles requests to create a new notification
func (h *NotificationHandler) CreateNotification(c *gin.Context) {
	// Define request structure
	var req struct {
		Title     string `json:"title" binding:"required"`
		Content   string `json:"content" binding:"required"`
		Type      string `json:"type" binding:"required"`
		Status    string `json:"status"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	// Bind JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	// Get user ID from context (assuming it's set by auth middleware)
	userID := c.GetUint("user_id")
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Parse dates if provided
	var startDate, endDate *time.Time
	if req.StartDate != "" {
		parsedStartDate, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "Invalid start date format")
			return
		}
		startDate = &parsedStartDate
	}
	if req.EndDate != "" {
		parsedEndDate, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "Invalid end date format")
			return
		}
		endDate = &parsedEndDate
	}

	// Set default status if not provided
	status := req.Status
	if status == "" {
		status = "draft"
	}

	// Create notification
	notification, err := h.notificationService.CreateNotification(
		req.Title,
		req.Content,
		req.Type,
		status,
		userID,
		startDate,
		endDate,
	)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create notification: "+err.Error())
		return
	}

	response.Success(c, "Notification created successfully", notification)
}

// GetNotificationByID handles requests to get a notification by ID
func (h *NotificationHandler) GetNotificationByID(c *gin.Context) {
	// Parse notification ID from URL parameter
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid notification ID")
		return
	}

	// Get notification
	notification, err := h.notificationService.GetNotificationByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "Notification not found")
		return
	}

	response.Success(c, "Notification retrieved successfully", notification)
}

// ListNotifications handles requests to list notifications with pagination and filters
func (h *NotificationHandler) ListNotifications(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// Parse filter parameters
	status := c.Query("status")
	notificationType := c.Query("type")

	// Ensure page and pageSize are positive
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	// List notifications
	notifications, total, err := h.notificationService.ListNotifications(page, pageSize, status, notificationType)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to list notifications")
		return
	}

	response.Success(c, "Notifications retrieved successfully", gin.H{
		"data":      notifications,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// UpdateNotification handles requests to update a notification
func (h *NotificationHandler) UpdateNotification(c *gin.Context) {
	// Parse notification ID from URL parameter
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid notification ID")
		return
	}

	// Define request structure
	var req struct {
		Title     string `json:"title"`
		Content   string `json:"content"`
		Type      string `json:"type"`
		Status    string `json:"status"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	// Bind JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	// Parse dates if provided
	var startDate, endDate *time.Time
	if req.StartDate != "" {
		parsedStartDate, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "Invalid start date format")
			return
		}
		startDate = &parsedStartDate
	}
	if req.EndDate != "" {
		parsedEndDate, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "Invalid end date format")
			return
		}
		endDate = &parsedEndDate
	}

	// Update notification
	notification, err := h.notificationService.UpdateNotification(
		uint(id),
		req.Title,
		req.Content,
		req.Type,
		req.Status,
		startDate,
		endDate,
	)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update notification: "+err.Error())
		return
	}

	response.Success(c, "Notification updated successfully", notification)
}

// DeleteNotification handles requests to delete a notification
func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	// Parse notification ID from URL parameter
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid notification ID")
		return
	}

	// Delete notification
	if err := h.notificationService.DeleteNotification(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete notification: "+err.Error())
		return
	}

	response.Success(c, "Notification deleted successfully", nil)
}

// GetActiveNotifications handles requests to get active notifications
func (h *NotificationHandler) GetActiveNotifications(c *gin.Context) {
	// Get active notifications
	notifications, err := h.notificationService.GetActiveNotifications()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get active notifications")
		return
	}

	response.Success(c, "Active notifications retrieved successfully", notifications)
}
