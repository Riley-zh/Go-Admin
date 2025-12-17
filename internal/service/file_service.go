package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"go-admin/internal/model"
	"go-admin/internal/repository"
)

// FileService handles file business logic
type FileService struct {
	fileRepo *repository.FileRepository
}

// NewFileService creates a new file service
func NewFileService() *FileService {
	return &FileService{
		fileRepo: repository.NewFileRepository(),
	}
}

// UploadFile handles file upload logic
func (s *FileService) UploadFile(fileHeader *multipart.FileHeader, userID uint) (*model.File, error) {
	// Open uploaded file
	src, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Create upload directory if not exists
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Generate file path
	filePath := filepath.Join(uploadDir, fileHeader.Filename)

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// Copy file
	if _, err := io.Copy(dst, src); err != nil {
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}

	// Save file metadata to database
	file := &model.File{
		Name:      fileHeader.Filename,
		Path:      filePath,
		Size:      fileHeader.Size,
		MimeType:  fileHeader.Header.Get("Content-Type"),
		CreatedBy: userID,
	}

	if err := s.fileRepo.Create(file); err != nil {
		// If database save fails, remove the uploaded file
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to save file metadata: %w", err)
	}

	return file, nil
}

// GetFileByID retrieves a file by its ID
func (s *FileService) GetFileByID(id uint) (*model.File, error) {
	return s.fileRepo.GetByID(id)
}

// ListFiles retrieves files with pagination
func (s *FileService) ListFiles(page, pageSize int) ([]model.File, int64, error) {
	return s.fileRepo.List(page, pageSize)
}

// DeleteFile removes a file by its ID
func (s *FileService) DeleteFile(id uint) error {
	// First get the file to get its path
	file, err := s.fileRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get file: %w", err)
	}

	// Delete file from filesystem
	if err := os.Remove(file.Path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file from filesystem: %w", err)
	}

	// Delete file record from database
	if err := s.fileRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete file from database: %w", err)
	}

	return nil
}
