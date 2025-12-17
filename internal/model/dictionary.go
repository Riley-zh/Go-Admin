package model

import (
	"time"

	"gorm.io/gorm"
)

// Dictionary represents a dictionary entry
type Dictionary struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Name        string `gorm:"size:100;not null;uniqueIndex" json:"name"` // Dictionary name (unique)
	Title       string `gorm:"size:200;not null" json:"title"`            // Dictionary title
	Description string `gorm:"size:500" json:"description"`               // Dictionary description
	Status      int    `gorm:"default:1" json:"status"`                   // 1: active, 0: inactive
}

// DictionaryItem represents a dictionary item
type DictionaryItem struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	DictionaryID uint   `gorm:"not null;index" json:"dictionary_id"` // Foreign key to Dictionary
	Label        string `gorm:"size:200;not null" json:"label"`      // Item label
	Value        string `gorm:"size:200;not null" json:"value"`      // Item value
	Sort         int    `gorm:"default:0" json:"sort"`               // Sort order
	Status       int    `gorm:"default:1" json:"status"`             // 1: active, 0: inactive
}

// TableName specifies the table name for Dictionary
func (Dictionary) TableName() string {
	return "dictionaries"
}

// TableName specifies the table name for DictionaryItem
func (DictionaryItem) TableName() string {
	return "dictionary_items"
}
