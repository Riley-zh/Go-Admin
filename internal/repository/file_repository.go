package repository

import (
	"go-admin/internal/model"

	"gorm.io/gorm"
)

// FileRepository handles file data operations
type FileRepository struct {
	db *gorm.DB
}

// NewFileRepository creates a new file repository
func NewFileRepository() *FileRepository {
	return &FileRepository{
		db: GetDB(),
	}
}

// Create saves a new file record
func (r *FileRepository) Create(file *model.File) error {
	return r.db.Create(file).Error
}

// GetByID retrieves a file by its ID
func (r *FileRepository) GetByID(id uint) (*model.File, error) {
	var file model.File
	err := r.db.Where("id = ?", id).First(&file).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

// List retrieves files with pagination
func (r *FileRepository) List(page, pageSize int) ([]model.File, int64, error) {
	var files []model.File
	var total int64

	offset := (page - 1) * pageSize
	err := r.db.Model(&model.File{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&files).Error
	return files, total, err
}

// Delete removes a file by its ID
func (r *FileRepository) Delete(id uint) error {
	return r.db.Delete(&model.File{}, id).Error
}

// Update updates a file record
func (r *FileRepository) Update(file *model.File) error {
	return r.db.Save(file).Error
}
