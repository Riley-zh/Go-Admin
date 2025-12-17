package handler

import (
	"net/http"
	"strconv"

	"go-admin/internal/service"
	"go-admin/pkg/response"

	"github.com/gin-gonic/gin"
)

// FileHandler handles file-related HTTP requests
type FileHandler struct {
	fileService *service.FileService
}

// NewFileHandler creates a new file handler
func NewFileHandler() *FileHandler {
	return &FileHandler{
		fileService: service.NewFileService(),
	}
}

// UploadFile handles file upload requests
func (h *FileHandler) UploadFile(c *gin.Context) {
	// Get uploaded file from form
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to get uploaded file")
		return
	}

	// Get user ID from context (assuming it's set by auth middleware)
	userID := c.GetUint("user_id")
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Upload file
	uploadedFile, err := h.fileService.UploadFile(file, userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to upload file: "+err.Error())
		return
	}

	response.Success(c, "File uploaded successfully", uploadedFile)
}

// GetFileByID handles requests to get a file by ID
func (h *FileHandler) GetFileByID(c *gin.Context) {
	// Parse file ID from URL parameter
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid file ID")
		return
	}

	// Get file
	file, err := h.fileService.GetFileByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "File not found")
		return
	}

	response.Success(c, "File retrieved successfully", file)
}

// ListFiles handles requests to list files with pagination
func (h *FileHandler) ListFiles(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// Ensure page and pageSize are positive
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	// List files
	files, total, err := h.fileService.ListFiles(page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to list files")
		return
	}

	response.Success(c, "Files retrieved successfully", gin.H{
		"data":      files,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// DeleteFile handles requests to delete a file
func (h *FileHandler) DeleteFile(c *gin.Context) {
	// Parse file ID from URL parameter
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid file ID")
		return
	}

	// Delete file
	if err := h.fileService.DeleteFile(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete file: "+err.Error())
		return
	}

	response.Success(c, "File deleted successfully", nil)
}

// DownloadFile handles requests to download a file
func (h *FileHandler) DownloadFile(c *gin.Context) {
	// Parse file ID from URL parameter
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid file ID")
		return
	}

	// Get file metadata
	file, err := h.fileService.GetFileByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "File not found")
		return
	}

	// Serve file
	c.FileAttachment(file.Path, file.Name)
}
