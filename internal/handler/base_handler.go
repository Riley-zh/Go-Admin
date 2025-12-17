package handler

import (
	"go-admin/pkg/common"

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
	common.HandleError(c, err)
}

// HandleValidationError handles validation errors
func (h *BaseHandler) HandleValidationError(c *gin.Context, err error) {
	common.HandleValidationError(c, err)
}

// HandleSuccess handles success responses
func (h *BaseHandler) HandleSuccess(c *gin.Context, data interface{}) {
	common.HandleSuccess(c, data)
}

// HandleSuccessWithMessage handles success responses with a message
func (h *BaseHandler) HandleSuccessWithMessage(c *gin.Context, message string, data interface{}) {
	common.HandleSuccessWithMessage(c, message, data)
}

// HandleCreated handles created responses
func (h *BaseHandler) HandleCreated(c *gin.Context, message string, data interface{}) {
	common.HandleCreated(c, message, data)
}

// HandleDeleted handles deleted responses
func (h *BaseHandler) HandleDeleted(c *gin.Context, message string) {
	common.HandleDeleted(c, message)
}

// HandlePaginationResponse handles paginated responses
func (h *BaseHandler) HandlePaginationResponse(c *gin.Context, items interface{}, total int64, params common.PaginationParams) {
	common.HandlePaginationResponse(c, items, total, params)
}

// BindAndValidate binds and validates request
func (h *BaseHandler) BindAndValidate(c *gin.Context, req interface{}) bool {
	return common.BindAndValidate(c, req)
}

// ParseIDParam parses ID parameter from the request
func (h *BaseHandler) ParseIDParam(c *gin.Context, paramName string) (uint, error) {
	return common.ParseIDParam(c, paramName)
}

// GetPaginationParams extracts pagination parameters from the request
func (h *BaseHandler) GetPaginationParams(c *gin.Context) common.PaginationParams {
	return common.GetPaginationParams(c)
}
