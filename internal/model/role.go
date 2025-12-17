package model

import (
	"time"

	"gorm.io/gorm"
)

// Role represents a role in the system
type Role struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Name        string `gorm:"size:50;uniqueIndex;not null" json:"name"`
	Description string `gorm:"size:255" json:"description"`
	Status      int    `gorm:"default:1" json:"status"` // 1: active, 0: inactive
}

// GetID returns the ID of the role
func (r *Role) GetID() uint {
	return r.ID
}

// GetStatus returns the status of the role
func (r *Role) GetStatus() int {
	return r.Status
}

// TableName specifies the table name
func (Role) TableName() string {
	return "roles"
}

// Permission represents a permission in the system
type Permission struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Name        string `gorm:"size:50;uniqueIndex;not null" json:"name"`
	Description string `gorm:"size:255" json:"description"`
	Resource    string `gorm:"size:100;not null" json:"resource"` // e.g., "user", "role", "permission"
	Action      string `gorm:"size:50;not null" json:"action"`    // e.g., "create", "read", "update", "delete"
}

// TableName specifies the table name
func (Permission) TableName() string {
	return "permissions"
}

// RolePermission represents the relationship between roles and permissions
type RolePermission struct {
	ID           uint `gorm:"primarykey" json:"id"`
	RoleID       uint `gorm:"not null;index" json:"role_id"`
	PermissionID uint `gorm:"not null;index" json:"permission_id"`
}

// TableName specifies the table name
func (RolePermission) TableName() string {
	return "role_permissions"
}
