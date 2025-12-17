package model

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// Resource represents a resource in the system
type Resource struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Name        string `gorm:"size:100;uniqueIndex;not null" json:"name"`
	Description string `gorm:"size:255" json:"description"`
	Type        string `gorm:"size:50;not null" json:"type"` // system, module, menu, api, data
	ParentID    *uint  `gorm:"index" json:"parent_id"`
	Path        string `gorm:"size:255" json:"path"` // resource path or identifier
	Status      int    `gorm:"default:1" json:"status"` // 1: active, 0: inactive
}

// TableName specifies the table name
func (Resource) TableName() string {
	return "resources"
}

// Action represents an action that can be performed on resources
type Action struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Name        string `gorm:"size:50;uniqueIndex;not null" json:"name"`
	Description string `gorm:"size:255" json:"description"`
	Category    string `gorm:"size:50" json:"category"` // crud, system, business
}

// TableName specifies the table name
func (Action) TableName() string {
	return "actions"
}

// RoleHierarchy represents role inheritance relationships
type RoleHierarchy struct {
	ID         uint `gorm:"primarykey" json:"id"`
	ParentID   uint `gorm:"not null;index" json:"parent_id"`
	ChildID    uint `gorm:"not null;index" json:"child_id"`
	Permission bool `gorm:"default:true" json:"permission"` // whether to inherit permissions
}

// TableName specifies the table name
func (RoleHierarchy) TableName() string {
	return "role_hierarchies"
}

// PermissionExtended represents extended permissions with conditions
type PermissionExtended struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	RoleID     uint   `gorm:"not null;index" json:"role_id"`
	ResourceID uint   `gorm:"not null;index" json:"resource_id"`
	ActionID   uint   `gorm:"not null;index" json:"action_id"`
	Conditions string `gorm:"type:text" json:"conditions,omitempty"` // JSON string for ABAC conditions
	Priority   int    `gorm:"default:0" json:"priority"` // higher priority takes precedence
	Status     int    `gorm:"default:1" json:"status"`   // 1: active, 0: inactive
}

// TableName specifies the table name
func (PermissionExtended) TableName() string {
	return "permissions_extended"
}

// PermissionCondition represents conditions for ABAC
type PermissionCondition struct {
	ResourceAttributes map[string]interface{} `json:"resource_attributes,omitempty"` // e.g., {"department": "IT", "level": "confidential"}
	UserAttributes     map[string]interface{} `json:"user_attributes,omitempty"`     // e.g., {"department": "IT", "role_level": "manager"}
	Environment        map[string]interface{} `json:"environment,omitempty"`         // e.g., {"time": "09:00-18:00", "ip_range": "192.168.1.0/24"}
	Expression         string                 `json:"expression,omitempty"`        // SpEL or custom expression
}

// MarshalConditions marshals conditions to JSON string
func (pc *PermissionCondition) MarshalConditions() (string, error) {
	if pc == nil {
		return "", nil
	}
	data, err := json.Marshal(pc)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// UnmarshalConditions unmarshals conditions from JSON string
func (pc *PermissionCondition) UnmarshalConditions(data string) error {
	if data == "" {
		return nil
	}
	return json.Unmarshal([]byte(data), pc)
}

// UserAttribute represents user attributes for ABAC
type UserAttribute struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	UserID uint   `gorm:"not null;index" json:"user_id"`
	Key    string `gorm:"size:100;not null" json:"key"`
	Value  string `gorm:"size:255;not null" json:"value"`
	Type   string `gorm:"size:50;not null" json:"type"` // string, number, boolean, date
}

// TableName specifies the table name
func (UserAttribute) TableName() string {
	return "user_attributes"
}

// ResourceAttribute represents resource attributes for ABAC
type ResourceAttribute struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	ResourceID uint   `gorm:"not null;index" json:"resource_id"`
	Key        string `gorm:"size:100;not null" json:"key"`
	Value      string `gorm:"size:255;not null" json:"value"`
	Type       string `gorm:"size:50;not null" json:"type"` // string, number, boolean, date
}

// TableName specifies the table name
func (ResourceAttribute) TableName() string {
	return "resource_attributes"
}

// PermissionAuditLog represents permission audit logs
type PermissionAuditLog struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`

	UserID       uint   `gorm:"not null;index" json:"user_id"`
	ResourceID   uint   `gorm:"index" json:"resource_id"`
	ActionID     uint   `gorm:"index" json:"action_id"`
	PermissionID uint   `gorm:"index" json:"permission_id"`
	Operation    string `gorm:"size:50;not null" json:"operation"` // grant, revoke, check, deny
	Result       bool   `gorm:"not null" json:"result"`
	Reason       string `gorm:"type:text" json:"reason,omitempty"`
	IPAddress    string `gorm:"size:45" json:"ip_address,omitempty"`
	UserAgent    string `gorm:"size:255" json:"user_agent,omitempty"`
	Context      string `gorm:"type:text" json:"context,omitempty"` // JSON string with additional context
}

// TableName specifies the table name
func (PermissionAuditLog) TableName() string {
	return "permission_audit_logs"
}