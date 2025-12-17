package service

import (
	"fmt"
	"time"

	"go-admin/internal/model"
	"go-admin/internal/repository"
)

// NotificationService handles notification business logic
type NotificationService struct {
	notificationRepo *repository.NotificationRepository
}

// NewNotificationService creates a new notification service
func NewNotificationService() *NotificationService {
	return &NotificationService{
		notificationRepo: repository.NewNotificationRepository(),
	}
}

// CreateNotification creates a new notification
func (s *NotificationService) CreateNotification(title, content, notificationType, status string, userID uint, startDate, endDate *time.Time) (*model.Notification, error) {
	notification := &model.Notification{
		Title:     title,
		Content:   content,
		Type:      notificationType,
		Status:    status,
		CreatedBy: userID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	if err := s.notificationRepo.Create(notification); err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	return notification, nil
}

// GetNotificationByID retrieves a notification by its ID
func (s *NotificationService) GetNotificationByID(id uint) (*model.Notification, error) {
	return s.notificationRepo.GetByID(id)
}

// ListNotifications retrieves notifications with pagination and filters
func (s *NotificationService) ListNotifications(page, pageSize int, status, notificationType string) ([]model.Notification, int64, error) {
	return s.notificationRepo.List(page, pageSize, status, notificationType)
}

// UpdateNotification updates a notification
func (s *NotificationService) UpdateNotification(id uint, title, content, notificationType, status string, startDate, endDate *time.Time) (*model.Notification, error) {
	// First get the existing notification
	notification, err := s.notificationRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get notification: %w", err)
	}

	// Update fields
	notification.Title = title
	notification.Content = content
	notification.Type = notificationType
	notification.Status = status
	notification.StartDate = startDate
	notification.EndDate = endDate

	if err := s.notificationRepo.Update(notification); err != nil {
		return nil, fmt.Errorf("failed to update notification: %w", err)
	}

	return notification, nil
}

// DeleteNotification removes a notification by its ID
func (s *NotificationService) DeleteNotification(id uint) error {
	if err := s.notificationRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}
	return nil
}

// GetActiveNotifications retrieves active notifications
func (s *NotificationService) GetActiveNotifications() ([]model.Notification, error) {
	return s.notificationRepo.GetActiveNotifications()
}
