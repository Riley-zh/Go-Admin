package service

import (
	"time"

	"go-admin/internal/database"
	"go-admin/internal/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// AuditService 审计服务
type AuditService struct {
	db *gorm.DB
}

// AuditLog 审计日志结构
type AuditLog struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id"`
	ActionType  string    `json:"action_type"`
	Resource    string    `json:"resource"`
	IP          string    `json:"ip"`
	UserAgent   string    `json:"user_agent"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// TableName 指定表名
func (AuditLog) TableName() string {
	return "audit_logs"
}

// NewAuditService 创建新的审计服务
func NewAuditService() *AuditService {
	return &AuditService{
		db: database.GetDB(),
	}
}

// Log 记录审计日志
func (s *AuditService) Log(userID uint, actionType, resource, description string, c *gin.Context) {
	auditLog := &AuditLog{
		UserID:      userID,
		ActionType:  actionType,
		Resource:    resource,
		IP:          c.ClientIP(),
		UserAgent:   c.GetHeader("User-Agent"),
		Description: description,
		CreatedAt:   time.Now(),
	}

	// 异步记录审计日志
	go func() {
		if err := s.db.Create(auditLog).Error; err != nil {
			logger.Error("Failed to record audit log", zap.Error(err))
		}
	}()
}

// GetAuditLogs 获取审计日志列表
func (s *AuditService) GetAuditLogs(page, pageSize int) ([]AuditLog, int64, error) {
	var logs []AuditLog
	var total int64

	offset := (page - 1) * pageSize

	// 获取总数
	if err := s.db.Model(&AuditLog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	if err := s.db.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetAuditLogByID 根据ID获取审计日志
func (s *AuditService) GetAuditLogByID(id uint) (*AuditLog, error) {
	var log AuditLog
	if err := s.db.First(&log, id).Error; err != nil {
		return nil, err
	}
	return &log, nil
}

// DeleteAuditLog 删除审计日志
func (s *AuditService) DeleteAuditLog(id uint) error {
	return s.db.Delete(&AuditLog{}, id).Error
}

// ClearAuditLogs 清空审计日志
func (s *AuditService) ClearAuditLogs() error {
	return s.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&AuditLog{}).Error
}

// GetUserAuditLogs 获取特定用户的审计日志
func (s *AuditService) GetUserAuditLogs(userID, page, pageSize int) ([]AuditLog, int64, error) {
	var logs []AuditLog
	var total int64

	offset := (page - 1) * pageSize

	// 获取总数
	if err := s.db.Model(&AuditLog{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	if err := s.db.Where("user_id = ?", userID).Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
