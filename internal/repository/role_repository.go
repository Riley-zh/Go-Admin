package repository

import (
	"go-admin/internal/database"
	"go-admin/internal/model"

	"gorm.io/gorm"
)

// RoleRepository defines the role repository interface
type RoleRepository interface {
	BaseRepository[*model.Role]
	GetRolesByUserID(userID uint) ([]*model.Role, error)
	GetUserRoles(userID uint) ([]*model.Role, error)
	GetRoleHierarchy(roleID uint) ([]*model.Role, error)
	GetRoleChildren(roleID uint) ([]*model.Role, error)
}

// roleRepository implements RoleRepository interface
type roleRepository struct {
	BaseRepository[*model.Role]
	db *gorm.DB
}

// NewRoleRepository creates a new role repository
func NewRoleRepository() RoleRepository {
	return &roleRepository{
		BaseRepository: NewBaseRepository(&model.Role{}),
		db:             database.GetDB(),
	}
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
