package repository

import (
	"go-admin/internal/model"
)

// ResourceRepository defines the resource repository interface
type ResourceRepository interface {
	Create(resource *model.Resource) error
	Update(resource *model.Resource) error
	Delete(id uint) error
	GetByID(id uint) (*model.Resource, error)
	GetByName(name string) (*model.Resource, error)
	List(query *ResourceQuery) ([]*model.Resource, int64, error)
	GetChildren(parentID uint) ([]*model.Resource, error)
}

// ActionRepository defines the action repository interface
type ActionRepository interface {
	Create(action *model.Action) error
	Update(action *model.Action) error
	Delete(id uint) error
	GetByID(id uint) (*model.Action, error)
	GetByName(name string) (*model.Action, error)
	List(query *ActionQuery) ([]*model.Action, int64, error)
}

// PermissionRepository defines the permission repository interface
type PermissionRepository interface {
	Create(permission *model.PermissionExtended) error
	Update(permission *model.PermissionExtended) error
	Delete(id uint) error
	GetByID(id uint) (*model.PermissionExtended, error)
	GetByRoleResourceAction(roleID, resourceID, actionID uint) (*model.PermissionExtended, error)
	GetByRoleID(roleID uint) ([]*model.PermissionExtended, error)
	GetByResourceID(resourceID uint) ([]*model.PermissionExtended, error)
	GetByActionID(actionID uint) ([]*model.PermissionExtended, error)
	CheckUserPermission(userID uint, resource, action string) (bool, error)
	
	// Permission model operations (for Permission model - used by handler)
	CreatePermission(permission *model.Permission) error
	UpdatePermission(permission *model.Permission) error
	DeletePermission(permission *model.Permission) error
	GetPermissionByID(id uint) (*model.Permission, error)
	GetPermissionByName(name string) (*model.Permission, error)
	ListPermissions(page, pageSize int) ([]*model.Permission, int64, error)
	
	// Role-Permission operations
	AssignPermissionToRole(roleID, permissionID uint) error
	RemovePermissionFromRole(roleID, permissionID uint) error
	GetPermissionsByRoleID(roleID uint) ([]*model.Permission, error)
	GetPermissionsByUserID(userID uint) ([]*model.Permission, error)
	
	// User attributes
	CreateUserAttribute(attribute *model.UserAttribute) error
	UpdateUserAttribute(attribute *model.UserAttribute) error
	GetUserAttribute(userID uint, key string) (*model.UserAttribute, error)
	GetUserAttributes(userID uint) (map[string]interface{}, error)
	
	// Resource attributes
	CreateResourceAttribute(attribute *model.ResourceAttribute) error
	UpdateResourceAttribute(attribute *model.ResourceAttribute) error
	GetResourceAttribute(resourceID uint, key string) (*model.ResourceAttribute, error)
	GetResourceAttributes(resourceID uint) (map[string]interface{}, error)
	
	// Role hierarchy
	CreateRoleHierarchy(hierarchy *model.RoleHierarchy) error
	DeleteRoleHierarchy(id uint) error
	GetRoleHierarchy(parentID, childID uint) (*model.RoleHierarchy, error)
	GetRoleHierarchyByParent(parentID uint) ([]*model.RoleHierarchy, error)
	GetRoleHierarchyByChild(childID uint) ([]*model.RoleHierarchy, error)
	
	// Audit logging
	CreateAuditLog(log *model.PermissionAuditLog) error
	GetAuditLogs(userID uint, limit int) ([]*model.PermissionAuditLog, error)
}