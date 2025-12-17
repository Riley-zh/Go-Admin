package handler

import (
	"net/http"
	"strconv"

	"go-admin/internal/service"
	"go-admin/pkg/response"

	"github.com/gin-gonic/gin"
)

// TaskHandler handles task-related HTTP requests
type TaskHandler struct {
	taskService *service.TaskService
}

// NewTaskHandler creates a new task handler
func NewTaskHandler() *TaskHandler {
	return &TaskHandler{
		taskService: service.NewTaskService(),
	}
}

// CreateTask handles requests to create a new task
func (h *TaskHandler) CreateTask(c *gin.Context) {
	// Define request structure
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		CronExpr    string `json:"cron_expr" binding:"required"`
		Handler     string `json:"handler" binding:"required"`
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

	// Create task
	task, err := h.taskService.CreateTask(
		req.Name,
		req.Description,
		req.CronExpr,
		req.Handler,
		userID,
	)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create task: "+err.Error())
		return
	}

	response.Success(c, "Task created successfully", task)
}

// GetTaskByID handles requests to get a task by ID
func (h *TaskHandler) GetTaskByID(c *gin.Context) {
	// Parse task ID from URL parameter
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid task ID")
		return
	}

	// Get task
	task, err := h.taskService.GetTaskByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "Task not found")
		return
	}

	response.Success(c, "Task retrieved successfully", task)
}

// ListTasks handles requests to list tasks with pagination
func (h *TaskHandler) ListTasks(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// Parse filter parameters
	status := c.Query("status")

	// Ensure page and pageSize are positive
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	// List tasks
	tasks, total, err := h.taskService.ListTasks(page, pageSize, status)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to list tasks")
		return
	}

	response.Success(c, "Tasks retrieved successfully", gin.H{
		"data":      tasks,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// UpdateTask handles requests to update a task
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	// Parse task ID from URL parameter
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid task ID")
		return
	}

	// Define request structure
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		CronExpr    string `json:"cron_expr"`
		Handler     string `json:"handler"`
		Status      string `json:"status"`
	}

	// Bind JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	// Update task
	task, err := h.taskService.UpdateTask(
		uint(id),
		req.Name,
		req.Description,
		req.CronExpr,
		req.Handler,
		req.Status,
	)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update task: "+err.Error())
		return
	}

	response.Success(c, "Task updated successfully", task)
}

// DeleteTask handles requests to delete a task
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	// Parse task ID from URL parameter
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid task ID")
		return
	}

	// Delete task
	if err := h.taskService.DeleteTask(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete task: "+err.Error())
		return
	}

	response.Success(c, "Task deleted successfully", nil)
}

// RunTaskImmediately handles requests to run a task immediately
func (h *TaskHandler) RunTaskImmediately(c *gin.Context) {
	// Parse task ID from URL parameter
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid task ID")
		return
	}

	// Run task immediately
	if err := h.taskService.RunTaskImmediately(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to run task: "+err.Error())
		return
	}

	response.Success(c, "Task started successfully", nil)
}
