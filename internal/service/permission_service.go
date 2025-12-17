package service

import (
	"go-admin/internal/model"
	"go-admin/internal/repository"
	"go-admin/pkg/errors"
)

// PermissionService defines the permission service interface
type PermissionService interface {
	CreatePermission(name, description, resource, action string) (*model.Permission, error)
	GetPermissionByID(id uint) (*model.Permission, error)
	GetPermissionByName(name string) (*model.Permission, error)
	UpdatePermission(permission *model.Permission) error
	DeletePermission(id uint) error
	ListPermissions(page, pageSize int) ([]*model.Permission, int64, error)
	AssignPermissionToRole(roleID, permissionID uint) error
	RemovePermissionFromRole(roleID, permissionID uint) error
	GetPermissionsByRoleID(roleID uint) ([]*model.Permission, error)
	GetPermissionsByUserID(userID uint) ([]*model.Permission, error)
	CheckUserPermission(userID uint, resource, action string) (bool, error)
}

// permissionService implements PermissionService interface
type permissionService struct {
	permissionRepo repository.PermissionRepository
	roleRepo       repository.RoleRepository
}

// NewPermissionService creates a new permission service
func NewPermissionService() PermissionService {
	return &permissionService{
		permissionRepo: repository.NewPermissionRepository(),
		roleRepo:       repository.NewRoleRepository(),
	}
}

// CreatePermission creates a new permission
func (s *permissionService) CreatePermission(name, description, resource, action string) (*model.Permission, error) {
	// Check if permission name already exists
	existingPermission, err := s.permissionRepo.GetByName(name)
	if err != nil {
		return nil, err
	}
	if existingPermission != nil {
		return nil, errors.Conflict("Permission name already exists", "权限名称已存在")
	}

	// Create permission
	permission := &model.Permission{
		Name:        name,
		Description: description,
		Resource:    resource,
		Action:      action,
	}

	err = s.permissionRepo.Create(permission)
	if err != nil {
		return nil, err
	}

	return permission, nil
}

// GetPermissionByID gets a permission by ID
func (s *permissionService) GetPermissionByID(id uint) (*model.Permission, error) {
	permission, err := s.permissionRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if permission == nil {
		return nil, errors.NotFound("Permission not found", "权限不存在")
	}

	return permission, nil
}

// GetPermissionByName gets a permission by name
func (s *permissionService) GetPermissionByName(name string) (*model.Permission, error) {
	permission, err := s.permissionRepo.GetByName(name)
	if err != nil {
		return nil, err
	}
	if permission == nil {
		return nil, errors.NotFound("Permission not found", "权限不存在")
	}

	return permission, nil
}

// UpdatePermission updates a permission
func (s *permissionService) UpdatePermission(permission *model.Permission) error {
	// Check if permission exists
	existingPermission, err := s.permissionRepo.GetByID(permission.ID)
	if err != nil {
		return err
	}
	if existingPermission == nil {
		return errors.NotFound("Permission not found", "权限不存在")
	}

	// Check if permission name already exists (excluding current permission)
	if permission.Name != existingPermission.Name {
		otherPermission, err := s.permissionRepo.GetByName(permission.Name)
		if err != nil {
			return err
		}
		if otherPermission != nil {
			return errors.Conflict("Permission name already exists", "权限名称已存在")
		}
	}

	// Update permission
	return s.permissionRepo.Update(permission)
}

// DeletePermission deletes a permission
func (s *permissionService) DeletePermission(id uint) error {
	// Check if permission exists
	existingPermission, err := s.permissionRepo.GetByID(id)
	if err != nil {
		return err
	}
	if existingPermission == nil {
		return errors.NotFound("Permission not found", "权限不存在")
	}

	// Delete permission
	return s.permissionRepo.Delete(id)
}

// ListPermissions lists permissions with pagination
func (s *permissionService) ListPermissions(page, pageSize int) ([]*model.Permission, int64, error) {
	return s.permissionRepo.List(page, pageSize)
}

// AssignPermissionToRole assigns a permission to a role
func (s *permissionService) AssignPermissionToRole(roleID, permissionID uint) error {
	// Check if role exists
	role, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.NotFound("Role not found", "角色不存在")
	}

	// Check if permission exists
	permission, err := s.permissionRepo.GetByID(permissionID)
	if err != nil {
		return err
	}
	if permission == nil {
		return errors.NotFound("Permission not found", "权限不存在")
	}

	// Check if the role-permission relationship already exists
	// This would typically be done with a database query, but for simplicity,
	// we'll let the database handle duplicate constraints

	// Create role-permission relationship
	rolePermission := &model.RolePermission{
		RoleID:       roleID,
		PermissionID: permissionID,
	}

	db := repository.GetDB()
	err = db.Create(rolePermission).Error
	if err != nil {
		// Check if it's a duplicate entry error
		// In a real implementation, you'd check the specific database error
		return errors.Conflict("Permission already assigned to role", "权限已分配给该角色")
	}

	return nil
}

// RemovePermissionFromRole removes a permission from a role
func (s *permissionService) RemovePermissionFromRole(roleID, permissionID uint) error {
	// Check if role exists
	role, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.NotFound("Role not found", "角色不存在")
	}

	// Check if permission exists
	permission, err := s.permissionRepo.GetByID(permissionID)
	if err != nil {
		return err
	}
	if permission == nil {
		return errors.NotFound("Permission not found", "权限不存在")
	}

	// Delete role-permission relationship
	db := repository.GetDB()
	result := db.Where("role_id = ? AND permission_id = ?", roleID, permissionID).Delete(&model.RolePermission{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.NotFound("Permission not assigned to role", "权限未分配给该角色")
	}

	return nil
}

// GetPermissionsByRoleID gets permissions by role ID
func (s *permissionService) GetPermissionsByRoleID(roleID uint) ([]*model.Permission, error) {
	// Check if role exists
	role, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.NotFound("Role not found", "角色不存在")
	}

	return s.permissionRepo.GetPermissionsByRoleID(roleID)
}

// GetPermissionsByUserID gets permissions by user ID
func (s *permissionService) GetPermissionsByUserID(userID uint) ([]*model.Permission, error) {
	return s.permissionRepo.GetPermissionsByUserID(userID)
}

// CheckUserPermission checks if a user has a specific permission
func (s *permissionService) CheckUserPermission(userID uint, resource, action string) (bool, error) {
	return s.permissionRepo.CheckUserPermission(userID, resource, action)
}
