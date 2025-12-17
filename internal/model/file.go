package model

import (
	"time"

	"gorm.io/gorm"
)

// File represents a file entity
type File struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"not null;size:255" json:"name"`
	Path      string         `gorm:"not null;size:500" json:"path"`
	Size      int64          `gorm:"not null" json:"size"`
	MimeType  string         `gorm:"size:100" json:"mime_type"`
	CreatedBy uint           `gorm:"not null" json:"created_by"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
