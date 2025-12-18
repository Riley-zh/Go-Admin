package handler

import (
	"go-admin/pkg/response"

	"github.com/gin-gonic/gin"
)

// BaseHandler provides common handler functionality
type BaseHandler struct{}

// NewBaseHandler creates a new base handler
func NewBaseHandler() *BaseHandler {
	return &BaseHandler{}
}

// HandleError handles error responses
func (h *BaseHandler) HandleError(c *gin.Context, err error) {
	response.HandleError(c, err)
}

// HandleValidationError handles validation errors
func (h *BaseHandler) HandleValidationError(c *gin.Context, err error) {
	response.HandleValidationError(c, err)
}

// HandleSuccess handles success responses
func (h *BaseHandler) HandleSuccess(c *gin.Context, data interface{}) {
	response.HandleSuccess(c, data)
}

// HandleSuccessWithMessage handles success responses with a message
func (h *BaseHandler) HandleSuccessWithMessage(c *gin.Context, message string, data interface{}) {
	response.HandleSuccessWithMessage(c, message, data)
}

// HandleCreated handles created responses
func (h *BaseHandler) HandleCreated(c *gin.Context, message string, data interface{}) {
	response.HandleCreated(c, message, data)
}

// HandleDeleted handles deleted responses
func (h *BaseHandler) HandleDeleted(c *gin.Context, message string) {
	response.HandleDeleted(c, message)
}

// HandlePaginationResponse handles paginated responses
func (h *BaseHandler) HandlePaginationResponse(c *gin.Context, items interface{}, total int64, params response.PaginationParams) {
	response.HandlePaginationResponse(c, items, total, params)
}

// BindAndValidate binds and validates request
func (h *BaseHandler) BindAndValidate(c *gin.Context, req interface{}) bool {
	return response.BindAndValidate(c, req)
}

// ParseIDParam parses ID parameter from the request
func (h *BaseHandler) ParseIDParam(c *gin.Context, paramName string) (uint, error) {
	return response.ParseIDParam(c, paramName)
}

// GetPaginationParams extracts pagination parameters from the request
func (h *BaseHandler) GetPaginationParams(c *gin.Context) response.PaginationParams {
	return response.GetPaginationParams(c)
}
