package repository

import (
	"errors"

	"go-admin/internal/database"
	"go-admin/internal/model"

	"gorm.io/gorm"
)

// LogRepository defines the log repository interface
type LogRepository interface {
	Create(log *model.Log) error
	GetByID(id uint) (*model.Log, error)
	List(page, pageSize int, level, method, path, username string) ([]*model.Log, int64, error)
	Delete(id uint) error
	DeleteByCondition(level, method, path, username string, days int) (int64, error)
}

// logRepository implements LogRepository interface
type logRepository struct {
	db *gorm.DB
}

// NewLogRepository creates a new log repository
func NewLogRepository() LogRepository {
	return &logRepository{
		db: database.GetDB(),
	}
}

// Create creates a new log entry
func (r *logRepository) Create(log *model.Log) error {
	return r.db.Create(log).Error
}

// GetByID gets a log entry by ID
func (r *logRepository) GetByID(id uint) (*model.Log, error) {
	var log model.Log
	err := r.db.First(&log, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &log, nil
}

// List lists log entries with pagination and filtering
func (r *logRepository) List(page, pageSize int, level, method, path, username string) ([]*model.Log, int64, error) {
	var logs []*model.Log
	var total int64

	offset := (page - 1) * pageSize
	query := r.db.Model(&model.Log{})

	// Apply filters
	if level != "" {
		query = query.Where("level = ?", level)
	}
	if method != "" {
		query = query.Where("method = ?", method)
	}
	if path != "" {
		query = query.Where("path LIKE ?", "%"+path+"%")
	}
	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}

	// Get total count
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get logs with pagination
	err = query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// Delete deletes a log entry by ID
func (r *logRepository) Delete(id uint) error {
	return r.db.Delete(&model.Log{}, id).Error
}

// DeleteByCondition deletes log entries by condition
func (r *logRepository) DeleteByCondition(level, method, path, username string, days int) (int64, error) {
	query := r.db.Model(&model.Log{})

	// Apply filters
	if level != "" {
		query = query.Where("level = ?", level)
	}
	if method != "" {
		query = query.Where("method = ?", method)
	}
	if path != "" {
		query = query.Where("path LIKE ?", "%"+path+"%")
	}
	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if days > 0 {
		query = query.Where("created_at < DATE_SUB(NOW(), INTERVAL ? DAY)", days)
	}

	result := query.Delete(&model.Log{})
	return result.RowsAffected, result.Error
}
