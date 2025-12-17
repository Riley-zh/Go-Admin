package model

import (
	"time"

	"gorm.io/gorm"
)

// Notification represents a notification or announcement entity
type Notification struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `gorm:"not null;size:255" json:"title"`
	Content   string         `gorm:"type:text" json:"content"`
	Type      string         `gorm:"size:50;default:'announcement'" json:"type"` // announcement, notification
	Status    string         `gorm:"size:20;default:'draft'" json:"status"`      // draft, published, archived
	CreatedBy uint           `gorm:"not null" json:"created_by"`
	StartDate *time.Time     `json:"start_date,omitempty"`
	EndDate   *time.Time     `json:"end_date,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
