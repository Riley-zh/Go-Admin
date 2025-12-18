package response

import (
	"net/http"
	"strconv"

	"go-admin/pkg/errors"

	"github.com/gin-gonic/gin"
)

// Response represents the standard API response structure
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success returns a successful response
func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: message,
		Data:    data,
	})
}

// Error returns an error response
func Error(c *gin.Context, code int, message string) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
	})
}

// HandleError handles error responses in a consistent way
func HandleError(c *gin.Context, err error) {
	if appErr, ok := err.(*errors.Error); ok {
		c.JSON(appErr.Code, Response{
			Code:    appErr.Code,
			Message: appErr.Message,
			Data:    appErr.Details,
		})
		return
	}
	c.JSON(http.StatusInternalServerError, Response{
		Code:    http.StatusInternalServerError,
		Message: "Internal server error",
		Data:    err.Error(),
	})
}

// HandleValidationError handles validation errors
func HandleValidationError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    http.StatusBadRequest,
		Message: "Invalid request format",
		Data:    err.Error(),
	})
}

// HandleSuccess handles success responses
func HandleSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

// HandleSuccessWithMessage handles success responses with a message
func HandleSuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: message,
		Data:    data,
	})
}

// HandleCreated handles created responses
func HandleCreated(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Code:    http.StatusCreated,
		Message: message,
		Data:    data,
	})
}

// HandleDeleted handles deleted responses
func HandleDeleted(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: message,
	})
}

// PaginationParams represents pagination parameters
type PaginationParams struct {
	Page     int
	PageSize int
}

// GetPaginationParams extracts pagination parameters from the request
func GetPaginationParams(c *gin.Context) PaginationParams {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
	}
}

// ParseIDParam parses ID parameter from the request
func ParseIDParam(c *gin.Context, paramName string) (uint, error) {
	idStr := c.Param(paramName)
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

// HandlePaginationResponse handles paginated responses
func HandlePaginationResponse(c *gin.Context, items interface{}, total int64, params PaginationParams) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "success",
		Data: map[string]interface{}{
			"items":     items,
			"total":     total,
			"page":      params.Page,
			"page_size": params.PageSize,
		},
	})
}

// BindAndValidate binds and validates request
func BindAndValidate(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		HandleValidationError(c, err)
		return false
	}
	return true
}

// PageResponse represents paginated response structure
type PageResponse struct {
	Data     interface{} `json:"data"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}
