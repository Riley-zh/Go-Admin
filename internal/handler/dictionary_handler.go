package handler

import (
	"net/http"
	"strconv"

	"go-admin/internal/model"
	"go-admin/internal/service"
	"go-admin/pkg/errors"

	"github.com/gin-gonic/gin"
)

// DictionaryHandler represents the dictionary handler
type DictionaryHandler struct {
	dictService service.DictionaryService
}

// NewDictionaryHandler creates a new dictionary handler
func NewDictionaryHandler() *DictionaryHandler {
	return &DictionaryHandler{
		dictService: service.NewDictionaryService(),
	}
}

// CreateDictionaryRequest represents the create dictionary request body
type CreateDictionaryRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Title       string `json:"title" binding:"required,min=1,max=200"`
	Description string `json:"description" binding:"max=500"`
}

// UpdateDictionaryRequest represents the update dictionary request body
type UpdateDictionaryRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Title       string `json:"title" binding:"required,min=1,max=200"`
	Description string `json:"description" binding:"max=500"`
}

// CreateDictionaryItemRequest represents the create dictionary item request body
type CreateDictionaryItemRequest struct {
	Label  string `json:"label" binding:"required,min=1,max=200"`
	Value  string `json:"value" binding:"required,min=1,max=200"`
	Sort   int    `json:"sort" binding:"gte=0"`
	Status int    `json:"status" binding:"oneof=0 1"`
}

// UpdateDictionaryItemRequest represents the update dictionary item request body
type UpdateDictionaryItemRequest struct {
	Label  string `json:"label" binding:"required,min=1,max=200"`
	Value  string `json:"value" binding:"required,min=1,max=200"`
	Sort   int    `json:"sort" binding:"gte=0"`
	Status int    `json:"status" binding:"oneof=0 1"`
}

// CreateDictionary handles creating a new dictionary
func (h *DictionaryHandler) CreateDictionary(c *gin.Context) {
	// Validate request
	var req CreateDictionaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Create dictionary
	dictionary, err := h.dictService.CreateDictionary(req.Name, req.Title, req.Description)
	if err != nil {
		if appErr, ok := err.(*errors.Error); ok {
			c.JSON(appErr.Code, gin.H{
				"error":   appErr.Message,
				"details": appErr.Details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create dictionary",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Dictionary created successfully",
		"dictionary": dictionary,
	})
}

// GetDictionaryByID handles getting a dictionary by ID
func (h *DictionaryHandler) GetDictionaryByID(c *gin.Context) {
	// Get dictionary ID from path parameter
	idStr := c.Param("dictId")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid dictionary ID",
			"details": "字典ID格式不正确",
		})
		return
	}

	// Get dictionary
	dictionary, err := h.dictService.GetDictionaryByID(uint(id))
	if err != nil {
		if appErr, ok := err.(*errors.Error); ok {
			c.JSON(appErr.Code, gin.H{
				"error":   appErr.Message,
				"details": appErr.Details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get dictionary",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"dictionary": dictionary,
	})
}

// UpdateDictionary handles updating a dictionary
func (h *DictionaryHandler) UpdateDictionary(c *gin.Context) {
	// Get dictionary ID from path parameter
	idStr := c.Param("dictId")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid dictionary ID",
			"details": "字典ID格式不正确",
		})
		return
	}

	// Validate request
	var req UpdateDictionaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Create dictionary model
	dictionary := &model.Dictionary{
		ID:          uint(id),
		Name:        req.Name,
		Title:       req.Title,
		Description: req.Description,
	}

	// Update dictionary
	err = h.dictService.UpdateDictionary(dictionary)
	if err != nil {
		if appErr, ok := err.(*errors.Error); ok {
			c.JSON(appErr.Code, gin.H{
				"error":   appErr.Message,
				"details": appErr.Details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update dictionary",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Dictionary updated successfully",
	})
}

// DeleteDictionary handles deleting a dictionary
func (h *DictionaryHandler) DeleteDictionary(c *gin.Context) {
	// Get dictionary ID from path parameter
	idStr := c.Param("dictId")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid dictionary ID",
			"details": "字典ID格式不正确",
		})
		return
	}

	// Delete dictionary
	err = h.dictService.DeleteDictionary(uint(id))
	if err != nil {
		if appErr, ok := err.(*errors.Error); ok {
			c.JSON(appErr.Code, gin.H{
				"error":   appErr.Message,
				"details": appErr.Details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete dictionary",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Dictionary deleted successfully",
	})
}

// ListDictionaries handles listing dictionaries with pagination
func (h *DictionaryHandler) ListDictionaries(c *gin.Context) {
	// Get pagination parameters
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// List dictionaries
	dictionaries, total, err := h.dictService.ListDictionaries(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to list dictionaries",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"dictionaries": dictionaries,
		"total":        total,
		"page":         page,
		"page_size":    pageSize,
	})
}

// CreateDictionaryItem handles creating a new dictionary item
func (h *DictionaryHandler) CreateDictionaryItem(c *gin.Context) {
	// Get dictionary ID from path parameter
	dictIDStr := c.Param("dictId")
	dictID, err := strconv.ParseUint(dictIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid dictionary ID",
			"details": "字典ID格式不正确",
		})
		return
	}

	// Validate request
	var req CreateDictionaryItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Create dictionary item
	item, err := h.dictService.CreateDictionaryItem(uint(dictID), req.Label, req.Value, req.Sort, req.Status)
	if err != nil {
		if appErr, ok := err.(*errors.Error); ok {
			c.JSON(appErr.Code, gin.H{
				"error":   appErr.Message,
				"details": appErr.Details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create dictionary item",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Dictionary item created successfully",
		"item":    item,
	})
}

// GetDictionaryItemByID handles getting a dictionary item by ID
func (h *DictionaryHandler) GetDictionaryItemByID(c *gin.Context) {
	// Get item ID from path parameter
	idStr := c.Param("itemId")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid dictionary item ID",
			"details": "字典项ID格式不正确",
		})
		return
	}

	// Get dictionary item
	item, err := h.dictService.GetDictionaryItemByID(uint(id))
	if err != nil {
		if appErr, ok := err.(*errors.Error); ok {
			c.JSON(appErr.Code, gin.H{
				"error":   appErr.Message,
				"details": appErr.Details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get dictionary item",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"item": item,
	})
}

// UpdateDictionaryItem handles updating a dictionary item
func (h *DictionaryHandler) UpdateDictionaryItem(c *gin.Context) {
	// Get item ID from path parameter
	idStr := c.Param("itemId")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid dictionary item ID",
			"details": "字典项ID格式不正确",
		})
		return
	}

	// Validate request
	var req UpdateDictionaryItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Get dictionary ID from path parameter
	dictIDStr := c.Param("dictId")
	dictID, err := strconv.ParseInt(dictIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid dictionary ID",
			"details": "字典ID格式不正确",
		})
		return
	}

	// Create dictionary item object
	item := &model.DictionaryItem{
		ID:           uint(id),
		DictionaryID: uint(dictID),
		Label:        req.Label,
		Value:        req.Value,
		Sort:         req.Sort,
		Status:       req.Status,
	}

	// Update dictionary item
	err = h.dictService.UpdateDictionaryItem(item)
	if err != nil {
		if appErr, ok := err.(*errors.Error); ok {
			c.JSON(appErr.Code, gin.H{
				"error":   appErr.Message,
				"details": appErr.Details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update dictionary item",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Dictionary item updated successfully",
	})
}

// DeleteDictionaryItem handles deleting a dictionary item
func (h *DictionaryHandler) DeleteDictionaryItem(c *gin.Context) {
	// Get item ID from path parameter
	idStr := c.Param("itemId")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid dictionary item ID",
			"details": "字典项ID格式不正确",
		})
		return
	}

	// Delete dictionary item
	err = h.dictService.DeleteDictionaryItem(uint(id))
	if err != nil {
		if appErr, ok := err.(*errors.Error); ok {
			c.JSON(appErr.Code, gin.H{
				"error":   appErr.Message,
				"details": appErr.Details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete dictionary item",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Dictionary item deleted successfully",
	})
}

// ListDictionaryItems handles listing dictionary items with pagination
func (h *DictionaryHandler) ListDictionaryItems(c *gin.Context) {
	// Get dictionary ID from path parameter
	dictIDStr := c.Param("dictId")
	dictID, err := strconv.ParseInt(dictIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid dictionary ID",
			"details": "字典ID格式不正确",
		})
		return
	}

	// Get pagination parameters
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// List dictionary items
	items, total, err := h.dictService.ListDictionaryItems(int(dictID), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to list dictionary items",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":     items,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// ListAllDictionaryItems handles listing all dictionary items
func (h *DictionaryHandler) ListAllDictionaryItems(c *gin.Context) {
	// Get dictionary ID from path parameter
	dictIDStr := c.Param("dictId")
	dictID, err := strconv.ParseInt(dictIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid dictionary ID",
			"details": "字典ID格式不正确",
		})
		return
	}

	// List all dictionary items
	items, err := h.dictService.ListAllDictionaryItems(int(dictID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to list dictionary items",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": items,
	})
}
