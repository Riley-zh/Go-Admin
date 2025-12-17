package common

import (
	"net/http"
	"strconv"

	"go-admin/pkg/errors"

	"github.com/gin-gonic/gin"
)

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

// HandleError handles error responses in a consistent way
func HandleError(c *gin.Context, err error) {
	if appErr, ok := err.(*errors.Error); ok {
		c.JSON(appErr.Code, gin.H{
			"error":   appErr.Message,
			"details": appErr.Details,
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"error":   "Internal server error",
		"details": err.Error(),
	})
}

// HandleValidationError handles validation errors
func HandleValidationError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error":   "Invalid request format",
		"details": err.Error(),
	})
}

// HandleSuccess handles success responses
func HandleSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

// HandleSuccessWithMessage handles success responses with a message
func HandleSuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"message": message,
		"data":    data,
	})
}

// HandleCreated handles created responses
func HandleCreated(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, gin.H{
		"message": message,
		"data":    data,
	})
}

// HandleDeleted handles deleted responses
func HandleDeleted(c *gin.Context, message string) {
	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}

// HandlePaginationResponse handles paginated responses
func HandlePaginationResponse(c *gin.Context, items interface{}, total int64, params PaginationParams) {
	c.JSON(http.StatusOK, gin.H{
		"items":     items,
		"total":     total,
		"page":      params.Page,
		"page_size": params.PageSize,
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
