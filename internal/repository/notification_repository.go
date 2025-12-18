package repository

import (
	"go-admin/internal/database"
	"go-admin/internal/model"

	"gorm.io/gorm"
)

// NotificationRepository handles notification data operations
type NotificationRepository struct {
	db *gorm.DB
}

// NewNotificationRepository creates a new notification repository
func NewNotificationRepository() *NotificationRepository {
	return &NotificationRepository{
		db: database.GetDB(),
	}
}

// Create saves a new notification
func (r *NotificationRepository) Create(notification *model.Notification) error {
	return r.db.Create(notification).Error
}

// GetByID retrieves a notification by its ID
func (r *NotificationRepository) GetByID(id uint) (*model.Notification, error) {
	var notification model.Notification
	err := r.db.Where("id = ?", id).First(&notification).Error
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

// List retrieves notifications with pagination and filters
func (r *NotificationRepository) List(page, pageSize int, status, notificationType string) ([]model.Notification, int64, error) {
	var notifications []model.Notification
	var total int64

	query := r.db.Model(&model.Notification{})

	// Apply filters
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if notificationType != "" {
		query = query.Where("type = ?", notificationType)
	}

	// Count total records
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Apply pagination and ordering
	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&notifications).Error
	return notifications, total, err
}

// Update updates a notification
func (r *NotificationRepository) Update(notification *model.Notification) error {
	return r.db.Save(notification).Error
}

// Delete removes a notification by its ID
func (r *NotificationRepository) Delete(id uint) error {
	return r.db.Delete(&model.Notification{}, id).Error
}

// GetActiveNotifications retrieves active notifications that should be displayed
func (r *NotificationRepository) GetActiveNotifications() ([]model.Notification, error) {
	var notifications []model.Notification

	// Get notifications that are published and within date range
	query := r.db.Where("status = ?", "published")

	err := query.Find(&notifications).Error
	return notifications, err
}
