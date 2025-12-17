package model

import (
	"time"

	"gorm.io/gorm"
)

// Menu represents a menu item in the system
type Menu struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Name       string `gorm:"size:50;not null" json:"name"`
	Title      string `gorm:"size:100" json:"title"`
	Icon       string `gorm:"size:50" json:"icon"`
	Path       string `gorm:"size:255" json:"path"`
	Component  string `gorm:"size:255" json:"component"`
	Redirect   string `gorm:"size:255" json:"redirect"`
	Permission string `gorm:"size:100" json:"permission"` // Associated permission identifier
	ParentID   uint   `gorm:"default:0" json:"parent_id"` // 0 means root level
	Sort       int    `gorm:"default:0" json:"sort"`      // Sort order
	Status     int    `gorm:"default:1" json:"status"`    // 1: active, 0: inactive
	Hidden     int    `gorm:"default:0" json:"hidden"`    // 1: hidden, 0: visible
}

// TableName specifies the table name
func (Menu) TableName() string {
	return "menus"
}
