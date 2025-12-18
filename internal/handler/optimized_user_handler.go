package handler

import (
	"context"
	"fmt"
	"time"

	"go-admin/internal/service"
	"go-admin/pkg/api"
	"go-admin/pkg/middleware"
	"go-admin/pkg/validation"

	"github.com/gin-gonic/gin"
)

// OptimizedUserHandler represents an optimized user handler using the new API client
type OptimizedUserHandler struct {
	*BaseHandler
	userService service.UserService
	apiClient   *api.APIClient
	validator   *validation.Validator
}

// NewOptimizedUserHandler creates a new optimized user handler
func NewOptimizedUserHandler(apiClient *api.APIClient) *OptimizedUserHandler {
	return &OptimizedUserHandler{
		BaseHandler: NewBaseHandler(),
		userService: service.NewUserService(),
		apiClient:   apiClient,
		validator:   validation.NewValidator(apiClient.GetHTTPClient()),
	}
}

// OptimizedCreateUserRequest represents the create user request with enhanced validation
type OptimizedCreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50" validate:"username" example:"johndoe"`
	Password string `json:"password" binding:"required,min=8,max=50" validate:"password" example:"Password123!"`
	Email    string `json:"email" binding:"required,email" example:"johndoe@example.com"`
	Nickname string `json:"nickname" binding:"max=100" example:"John Doe"`
	Phone    string `json:"phone" validate:"phone" example:"+1234567890"`
	Avatar   string `json:"avatar" validate:"url" example:"https://example.com/avatar.jpg"`
}

// OptimizedUpdateUserRequest represents the update user request with enhanced validation
type OptimizedUpdateUserRequest struct {
	Nickname string `json:"nickname" binding:"max=100" example:"John Doe"`
	Phone    string `json:"phone" validate:"phone" example:"+1234567890"`
	Avatar   string `json:"avatar" validate:"url" example:"https://example.com/avatar.jpg"`
}

// CreateUser creates a new user with optimized API handling
// @Summary Create a new user
// @Description Create a new user with enhanced validation and optimized JSON processing
// @Tags users
// @Accept json
// @Produce json
// @Param request body OptimizedCreateUserRequest true "User details"
// @Success 201 {object} map[string]interface{} "User created successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 409 {object} map[string]interface{} "Conflict - User already exists"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /users [post]
func (h *OptimizedUserHandler) CreateUser(c *gin.Context) {
	var req OptimizedCreateUserRequest

	// Use the middleware validator for enhanced validation
	validator := middleware.GetValidator(c)
	if err := validator.Validate(&req); err != nil {
		c.JSON(400, gin.H{
			"error":   "Validation failed",
			"details": err.Error(),
		})
		return
	}

	// Additional validation with custom rules
	customRules := map[string]string{
		"username": "unique",
		"email":    "unique",
	}

	if err := validator.ValidateWithRules(&req, customRules); err != nil {
		c.JSON(400, gin.H{
			"error":   "Validation failed",
			"details": err.Error(),
		})
		return
	}

	user, err := h.userService.CreateUser(req.Username, req.Password, req.Email, req.Nickname)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleCreated(c, "User created successfully", user)
}

// UpdateUser updates a user with optimized API handling
// @Summary Update a user
// @Description Update a user with enhanced validation and optimized JSON processing
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body OptimizedUpdateUserRequest true "User details"
// @Success 200 {object} map[string]interface{} "User updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 404 {object} map[string]interface{} "Not Found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /users/{id} [put]
func (h *OptimizedUserHandler) UpdateUser(c *gin.Context) {
	id, err := h.ParseIDParam(c, "id")
	if err != nil {
		h.HandleError(c, err)
		return
	}

	var req OptimizedUpdateUserRequest

	// Use the middleware validator for enhanced validation
	validator := middleware.GetValidator(c)
	if err := validator.Validate(&req); err != nil {
		c.JSON(400, gin.H{
			"error":   "Validation failed",
			"details": err.Error(),
		})
		return
	}

	// Get the existing user
	user, err := h.userService.GetUserByID(id)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	// Update user fields
	user.Nickname = req.Nickname
	user.Avatar = req.Avatar

	err = h.userService.UpdateUser(user)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccessWithMessage(c, "User updated successfully", user)
}

// ExternalAPIExample demonstrates using the optimized API client for external API calls
// @Summary Example of external API call
// @Description Demonstrates using the optimized API client for external API calls
// @Tags external
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "External API response"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /external/api-example [get]
func (h *OptimizedUserHandler) ExternalAPIExample(c *gin.Context) {
	// Example of using the optimized API client for external API calls
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	// Define the response schema
	var response struct {
		Data struct {
			ID       int    `json:"id"`
			Name     string `json:"name"`
			Username string `json:"username"`
			Email    string `json:"email"`
		} `json:"data"`
		Success bool   `json:"success"`
		Message string `json:"message"`
	}

	// Make the API call with validation
	err := h.apiClient.GetWithValidation(ctx, "/external/users/1", &response, &response)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, response)
}

// BatchAPIExample demonstrates using the batch processor for multiple API calls
// @Summary Example of batch API calls
// @Description Demonstrates using the batch processor for multiple API calls
// @Tags external
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Batch API responses"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /external/batch-example [post]
func (h *OptimizedUserHandler) BatchAPIExample(c *gin.Context) {
	var request struct {
		UserIDs []int `json:"user_ids" binding:"required"`
	}

	if !h.BindAndValidate(c, &request) {
		return
	}

	// Create batch processor
	batchProcessor := api.NewBatchProcessor(h.apiClient, 10, 5)

	// Prepare batch requests
	requests := make([]api.Request, len(request.UserIDs))
	for i, userID := range request.UserIDs {
		requests[i] = api.Request{
			Method: "GET",
			Path:   fmt.Sprintf("/external/users/%d", userID),
		}
	}

	// Process batch
	ctx, cancel := context.WithTimeout(c.Request.Context(), 60*time.Second)
	defer cancel()

	responses, err := batchProcessor.ProcessBatch(ctx, requests)
	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, responses)
}

// StreamingAPIExample demonstrates using the streaming processor for large API responses
// @Summary Example of streaming API response
// @Description Demonstrates using the streaming processor for large API responses
// @Tags external
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Streaming API response"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /external/streaming-example [get]
func (h *OptimizedUserHandler) StreamingAPIExample(c *gin.Context) {
	// Create streaming processor
	streamingProcessor := api.NewStreamingProcessor(h.apiClient)

	// Process streaming response
	var results []interface{}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 60*time.Second)
	defer cancel()

	err := streamingProcessor.ProcessStream(ctx, "/external/users/stream", func(item interface{}) error {
		results = append(results, item)
		return nil
	})

	if err != nil {
		h.HandleError(c, err)
		return
	}

	h.HandleSuccess(c, results)
}
