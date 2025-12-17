package model

import (
	"time"

	"gorm.io/gorm"
)

// Log represents a system log entry
type Log struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Level       string `gorm:"size:20;not null" json:"level"`   // Log level (INFO, WARN, ERROR, etc.)
	Method      string `gorm:"size:10" json:"method"`           // HTTP method (GET, POST, etc.)
	Path        string `gorm:"size:255" json:"path"`            // Request path
	StatusCode  int    `json:"status_code"`                     // HTTP status code
	ClientIP    string `gorm:"size:50" json:"client_ip"`        // Client IP address
	UserAgent   string `gorm:"size:500" json:"user_agent"`      // User agent
	RequestID   string `gorm:"size:50;index" json:"request_id"` // Request ID for tracing
	UserID      uint   `gorm:"index" json:"user_id"`            // User ID (if authenticated)
	Username    string `gorm:"size:50" json:"username"`         // Username (if authenticated)
	Message     string `gorm:"type:text" json:"message"`        // Log message
	RequestBody string `gorm:"type:text" json:"request_body"`   // Request body (for debugging)
	ErrorDetail string `gorm:"type:text" json:"error_detail"`   // Error details (for ERROR level)
	Response    string `gorm:"type:text" json:"response"`       // Response data (for debugging)
	Latency     int64  `json:"latency"`                         // Request latency in milliseconds
}

// TableName specifies the table name
func (Log) TableName() string {
	return "logs"
}
