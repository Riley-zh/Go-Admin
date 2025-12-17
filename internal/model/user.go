package model

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Username string `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Password string `gorm:"size:255;not null" json:"password"`
	Email    string `gorm:"size:100;uniqueIndex" json:"email"`
	Nickname string `gorm:"size:100" json:"nickname"`
	Avatar   string `gorm:"size:255" json:"avatar"`
	Status   int    `gorm:"default:1" json:"status"` // 1: active, 0: inactive
}

// GetID returns the ID of the user
func (u *User) GetID() uint {
	return u.ID
}

// GetStatus returns the status of the user
func (u *User) GetStatus() int {
	return u.Status
}

// UserWithRoles represents a user with their roles
type UserWithRoles struct {
	User
	Roles []*Role `json:"roles,omitempty"`
}

// TableName specifies the table name
func (User) TableName() string {
	return "users"
}

// UserRole represents the relationship between users and roles
type UserRole struct {
	ID     uint `gorm:"primarykey" json:"id"`
	UserID uint `gorm:"not null;index" json:"user_id"`
	RoleID uint `gorm:"not null;index" json:"role_id"`
}

// TableName specifies the table name
func (UserRole) TableName() string {
	return "user_roles"
}
