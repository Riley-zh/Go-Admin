package service

import (
	"go-admin/internal/model"
	"go-admin/internal/repository"
	"go-admin/pkg/errors"
)

// LogService defines the log service interface
type LogService interface {
	CreateLog(log *model.Log) error
	GetLogByID(id uint) (*model.Log, error)
	ListLogs(page, pageSize int, level, method, path, username string) ([]*model.Log, int64, error)
	DeleteLog(id uint) error
	ClearLogs(level, method, path, username string, days int) (int64, error)
}

// logService implements LogService interface
type logService struct {
	logRepo repository.LogRepository
}

// NewLogService creates a new log service
func NewLogService() LogService {
	return &logService{
		logRepo: repository.NewLogRepository(),
	}
}

// CreateLog creates a new log entry
func (s *logService) CreateLog(log *model.Log) error {
	return s.logRepo.Create(log)
}

// GetLogByID gets a log entry by ID
func (s *logService) GetLogByID(id uint) (*model.Log, error) {
	log, err := s.logRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if log == nil {
		return nil, errors.NotFound("Log not found", "日志不存在")
	}

	return log, nil
}

// ListLogs lists log entries with pagination and filtering
func (s *logService) ListLogs(page, pageSize int, level, method, path, username string) ([]*model.Log, int64, error) {
	return s.logRepo.List(page, pageSize, level, method, path, username)
}

// DeleteLog deletes a log entry by ID
func (s *logService) DeleteLog(id uint) error {
	// Check if log exists
	existingLog, err := s.logRepo.GetByID(id)
	if err != nil {
		return err
	}
	if existingLog == nil {
		return errors.NotFound("Log not found", "日志不存在")
	}

	// Delete log
	return s.logRepo.Delete(id)
}

// ClearLogs clears log entries by condition
func (s *logService) ClearLogs(level, method, path, username string, days int) (int64, error) {
	return s.logRepo.DeleteByCondition(level, method, path, username, days)
}
