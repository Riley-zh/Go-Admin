package handler

import (
	"net/http"
	"strconv"

	"go-admin/internal/model"
	"go-admin/internal/service"
	"go-admin/pkg/errors"

	"github.com/gin-gonic/gin"
)

// MenuHandler represents the menu handler
type MenuHandler struct {
	menuService service.MenuService
}

// NewMenuHandler creates a new menu handler
func NewMenuHandler() *MenuHandler {
	return &MenuHandler{
		menuService: service.NewMenuService(),
	}
}

// CreateMenuRequest represents the create menu request body
type CreateMenuRequest struct {
	Name       string `json:"name" binding:"required,min=1,max=50"`
	Title      string `json:"title" binding:"max=100"`
	Icon       string `json:"icon" binding:"max=50"`
	Path       string `json:"path" binding:"max=255"`
	Component  string `json:"component" binding:"max=255"`
	Redirect   string `json:"redirect" binding:"max=255"`
	Permission string `json:"permission" binding:"max=100"`
	ParentID   int    `json:"parent_id" binding:"gte=0"`
	Sort       int    `json:"sort" binding:"gte=0"`
	Status     int    `json:"status" binding:"oneof=0 1"`
	Hidden     int    `json:"hidden" binding:"oneof=0 1"`
}

// UpdateMenuRequest represents the update menu request body
type UpdateMenuRequest struct {
	Name       string `json:"name" binding:"required,min=1,max=50"`
	Title      string `json:"title" binding:"max=100"`
	Icon       string `json:"icon" binding:"max=50"`
	Path       string `json:"path" binding:"max=255"`
	Component  string `json:"component" binding:"max=255"`
	Redirect   string `json:"redirect" binding:"max=255"`
	Permission string `json:"permission" binding:"max=100"`
	ParentID   int    `json:"parent_id" binding:"gte=0"`
	Sort       int    `json:"sort" binding:"gte=0"`
	Status     int    `json:"status" binding:"oneof=0 1"`
	Hidden     int    `json:"hidden" binding:"oneof=0 1"`
}

// CreateMenu handles creating a new menu
func (h *MenuHandler) CreateMenu(c *gin.Context) {
	// Validate request
	var req CreateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Create menu
	menu, err := h.menuService.CreateMenu(
		req.Name, req.Title, req.Icon, req.Path, req.Component, req.Redirect, req.Permission,
		req.ParentID, req.Sort, req.Status, req.Hidden,
	)
	if err != nil {
		if appErr, ok := err.(*errors.Error); ok {
			c.JSON(appErr.Code, gin.H{
				"error":   appErr.Message,
				"details": appErr.Details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create menu",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Menu created successfully",
		"menu":    menu,
	})
}

// GetMenuByID handles getting a menu by ID
func (h *MenuHandler) GetMenuByID(c *gin.Context) {
	// Get menu ID from path parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid menu ID",
			"details": "菜单ID格式不正确",
		})
		return
	}

	// Get menu
	menu, err := h.menuService.GetMenuByID(uint(id))
	if err != nil {
		if appErr, ok := err.(*errors.Error); ok {
			c.JSON(appErr.Code, gin.H{
				"error":   appErr.Message,
				"details": appErr.Details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get menu",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"menu": menu,
	})
}

// UpdateMenu handles updating a menu
func (h *MenuHandler) UpdateMenu(c *gin.Context) {
	// Get menu ID from path parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid menu ID",
			"details": "菜单ID格式不正确",
		})
		return
	}

	// Validate request
	var req UpdateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Create menu model
	menu := &model.Menu{
		ID:         uint(id),
		Name:       req.Name,
		Title:      req.Title,
		Icon:       req.Icon,
		Path:       req.Path,
		Component:  req.Component,
		Redirect:   req.Redirect,
		Permission: req.Permission,
		ParentID:   uint(req.ParentID),
		Sort:       req.Sort,
		Status:     req.Status,
		Hidden:     req.Hidden,
	}

	// Update menu
	err = h.menuService.UpdateMenu(menu)
	if err != nil {
		if appErr, ok := err.(*errors.Error); ok {
			c.JSON(appErr.Code, gin.H{
				"error":   appErr.Message,
				"details": appErr.Details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update menu",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Menu updated successfully",
	})
}

// DeleteMenu handles deleting a menu
func (h *MenuHandler) DeleteMenu(c *gin.Context) {
	// Get menu ID from path parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid menu ID",
			"details": "菜单ID格式不正确",
		})
		return
	}

	// Delete menu
	err = h.menuService.DeleteMenu(uint(id))
	if err != nil {
		if appErr, ok := err.(*errors.Error); ok {
			c.JSON(appErr.Code, gin.H{
				"error":   appErr.Message,
				"details": appErr.Details,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete menu",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Menu deleted successfully",
	})
}

// ListMenus handles listing menus with pagination
func (h *MenuHandler) ListMenus(c *gin.Context) {
	// Get pagination parameters
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// List menus
	menus, total, err := h.menuService.ListMenus(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to list menus",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"menus":     menus,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetMenuTree handles getting menu tree
func (h *MenuHandler) GetMenuTree(c *gin.Context) {
	// Get menu tree
	menuTree, err := h.menuService.GetMenuTree()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get menu tree",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"menu_tree": menuTree,
	})
}
