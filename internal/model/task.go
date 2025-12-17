package model

import (
	"time"

	"gorm.io/gorm"
)

// Task represents a scheduled task
type Task struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null;size:255" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	CronExpr    string         `gorm:"not null;size:100" json:"cron_expr"`     // Cron expression
	Handler     string         `gorm:"not null;size:255" json:"handler"`       // Handler function name
	Status      string         `gorm:"size:20;default:'active'" json:"status"` // active, inactive, error
	LastRun     *time.Time     `json:"last_run,omitempty"`
	NextRun     *time.Time     `json:"next_run,omitempty"`
	CreatedBy   uint           `gorm:"not null" json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
