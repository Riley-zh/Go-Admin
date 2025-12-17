package repository

import (
	"errors"

	"go-admin/internal/database"
	"go-admin/internal/model"

	"gorm.io/gorm"
)

// PermissionRepository defines the permission repository interface
type PermissionRepository interface {
	Create(permission *model.Permission) error
	GetByID(id uint) (*model.Permission, error)
	GetByName(name string) (*model.Permission, error)
	Update(permission *model.Permission) error
	Delete(id uint) error
	List(page, pageSize int) ([]*model.Permission, int64, error)
	GetPermissionsByRoleID(roleID uint) ([]*model.Permission, error)
	GetPermissionsByUserID(userID uint) ([]*model.Permission, error)
	CheckUserPermission(userID uint, resource, action string) (bool, error)
}

// permissionRepository implements PermissionRepository interface
type permissionRepository struct {
	db *gorm.DB
}

// NewPermissionRepository creates a new permission repository
func NewPermissionRepository() PermissionRepository {
	return &permissionRepository{
		db: database.GetDB(),
	}
}

// Create creates a new permission
func (r *permissionRepository) Create(permission *model.Permission) error {
	return r.db.Create(permission).Error
}

// GetByID gets a permission by ID
func (r *permissionRepository) GetByID(id uint) (*model.Permission, error) {
	var permission model.Permission
	err := r.db.Where("id = ?", id).First(&permission).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &permission, nil
}

// GetByName gets a permission by name
func (r *permissionRepository) GetByName(name string) (*model.Permission, error) {
	var permission model.Permission
	err := r.db.Where("name = ?", name).First(&permission).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &permission, nil
}

// Update updates a permission
func (r *permissionRepository) Update(permission *model.Permission) error {
	return r.db.Save(permission).Error
}

// Delete deletes a permission (soft delete)
func (r *permissionRepository) Delete(id uint) error {
	return r.db.Where("id = ?", id).Delete(&model.Permission{}).Error
}

// List lists permissions with pagination
func (r *permissionRepository) List(page, pageSize int) ([]*model.Permission, int64, error) {
	var permissions []*model.Permission
	var total int64

	offset := (page - 1) * pageSize
	err := r.db.Model(&model.Permission{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Offset(offset).Limit(pageSize).Find(&permissions).Error
	if err != nil {
		return nil, 0, err
	}

	return permissions, total, nil
}

// GetPermissionsByRoleID gets permissions by role ID
func (r *permissionRepository) GetPermissionsByRoleID(roleID uint) ([]*model.Permission, error) {
	var permissions []*model.Permission
	err := r.db.
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

// GetPermissionsByUserID gets permissions by user ID
func (r *permissionRepository) GetPermissionsByUserID(userID uint) ([]*model.Permission, error) {
	var permissions []*model.Permission
	err := r.db.
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("JOIN user_roles ON role_permissions.role_id = user_roles.role_id").
		Where("user_roles.user_id = ?", userID).
		Group("permissions.id").
		Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

// CheckUserPermission checks if a user has a specific permission
func (r *permissionRepository) CheckUserPermission(userID uint, resource, action string) (bool, error) {
	var count int64
	err := r.db.
		Model(&model.Permission{}).
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("JOIN user_roles ON role_permissions.role_id = user_roles.role_id").
		Where("user_roles.user_id = ? AND permissions.resource = ? AND permissions.action = ?", userID, resource, action).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
