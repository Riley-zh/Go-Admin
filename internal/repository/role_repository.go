package repository

import (
	"errors"

	"go-admin/internal/database"
	"go-admin/internal/model"

	"gorm.io/gorm"
)

// RoleRepository defines the role repository interface
type RoleRepository interface {
	Create(role *model.Role) error
	GetByID(id uint) (*model.Role, error)
	GetByName(name string) (*model.Role, error)
	Update(role *model.Role) error
	Delete(id uint) error
	List(page, pageSize int) ([]*model.Role, int64, error)
	GetRolesByUserID(userID uint) ([]*model.Role, error)
	GetUserRoles(userID uint) ([]*model.Role, error)
	GetRoleHierarchy(roleID uint) ([]*model.Role, error)
	GetRoleChildren(roleID uint) ([]*model.Role, error)
}

// roleRepository implements RoleRepository interface
type roleRepository struct {
	db *gorm.DB
}

// NewRoleRepository creates a new role repository
func NewRoleRepository() RoleRepository {
	return &roleRepository{
		db: database.GetDB(),
	}
}

// Create creates a new role
func (r *roleRepository) Create(role *model.Role) error {
	return r.db.Create(role).Error
}

// GetByID gets a role by ID
func (r *roleRepository) GetByID(id uint) (*model.Role, error) {
	var role model.Role
	err := r.db.Where("id = ? AND status = ?", id, 1).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

// GetByName gets a role by name
func (r *roleRepository) GetByName(name string) (*model.Role, error) {
	var role model.Role
	err := r.db.Where("name = ? AND status = ?", name, 1).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

// Update updates a role
func (r *roleRepository) Update(role *model.Role) error {
	return r.db.Save(role).Error
}

// Delete deletes a role (soft delete)
func (r *roleRepository) Delete(id uint) error {
	return r.db.Where("id = ?", id).Delete(&model.Role{}).Error
}

// List lists roles with pagination
func (r *roleRepository) List(page, pageSize int) ([]*model.Role, int64, error) {
	var roles []*model.Role
	var total int64

	offset := (page - 1) * pageSize
	err := r.db.Model(&model.Role{}).Where("status = ?", 1).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Where("status = ?", 1).Offset(offset).Limit(pageSize).Find(&roles).Error
	if err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

// GetRolesByUserID gets roles by user ID
func (r *roleRepository) GetRolesByUserID(userID uint) ([]*model.Role, error) {
	var roles []*model.Role
	err := r.db.
		Joins("JOIN user_roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ? AND roles.status = ?", userID, 1).
		Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// GetUserRoles gets roles by user ID (alias for GetRolesByUserID)
func (r *roleRepository) GetUserRoles(userID uint) ([]*model.Role, error) {
	return r.GetRolesByUserID(userID)
}

// GetRoleHierarchy gets all ancestor roles for a given role ID
func (r *roleRepository) GetRoleHierarchy(roleID uint) ([]*model.Role, error) {
	var roles []*model.Role
	// This is a simplified implementation
	// In a real application, you would traverse the role hierarchy tree
	query := `
		SELECT r.* FROM roles r
		JOIN role_hierarchy rh ON r.id = rh.parent_id
		WHERE rh.child_id = ? AND r.status = ?
	`
	err := r.db.Raw(query, roleID, 1).Scan(&roles).Error
	return roles, err
}

// GetRoleChildren gets all direct child roles for a given role ID
func (r *roleRepository) GetRoleChildren(roleID uint) ([]*model.Role, error) {
	var roles []*model.Role
	// This is a simplified implementation
	// In a real application, you would get direct children from the role hierarchy
	query := `
		SELECT r.* FROM roles r
		JOIN role_hierarchy rh ON r.id = rh.child_id
		WHERE rh.parent_id = ? AND r.status = ?
	`
	err := r.db.Raw(query, roleID, 1).Scan(&roles).Error
	return roles, err
}
