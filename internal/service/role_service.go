package service

import (
	"go-admin/internal/model"
	"go-admin/internal/repository"
	"go-admin/pkg/errors"
)

// RoleService defines the role service interface
type RoleService interface {
	CreateRole(name, description string) (*model.Role, error)
	GetRoleByID(id uint) (*model.Role, error)
	GetRoleByName(name string) (*model.Role, error)
	UpdateRole(role *model.Role) error
	DeleteRole(id uint) error
	ListRoles(page, pageSize int) ([]*model.Role, int64, error)
	AssignRoleToUser(userID, roleID uint) error
	RemoveRoleFromUser(userID, roleID uint) error
	GetRolesByUserID(userID uint) ([]*model.Role, error)
}

// roleService implements RoleService interface
type roleService struct {
	roleRepo repository.RoleRepository
	userRepo repository.UserRepository
}

// NewRoleService creates a new role service
func NewRoleService() RoleService {
	return &roleService{
		roleRepo: repository.NewRoleRepository(),
		userRepo: repository.NewUserRepository(),
	}
}

// CreateRole creates a new role
func (s *roleService) CreateRole(name, description string) (*model.Role, error) {
	// Check if role name already exists
	existingRole, err := s.roleRepo.GetByName(name)
	if err != nil {
		return nil, err
	}
	if existingRole != nil {
		return nil, errors.Conflict("Role name already exists", "角色名称已存在")
	}

	// Create role
	role := &model.Role{
		Name:        name,
		Description: description,
		Status:      1,
	}

	err = s.roleRepo.Create(role)
	if err != nil {
		return nil, err
	}

	return role, nil
}

// GetRoleByID gets a role by ID
func (s *roleService) GetRoleByID(id uint) (*model.Role, error) {
	role, err := s.roleRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.NotFound("Role not found", "角色不存在")
	}

	return role, nil
}

// GetRoleByName gets a role by name
func (s *roleService) GetRoleByName(name string) (*model.Role, error) {
	role, err := s.roleRepo.GetByName(name)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.NotFound("Role not found", "角色不存在")
	}

	return role, nil
}

// UpdateRole updates a role
func (s *roleService) UpdateRole(role *model.Role) error {
	// Check if role exists
	existingRole, err := s.roleRepo.GetByID(role.ID)
	if err != nil {
		return err
	}
	if existingRole == nil {
		return errors.NotFound("Role not found", "角色不存在")
	}

	// Check if role name already exists (excluding current role)
	if role.Name != existingRole.Name {
		otherRole, err := s.roleRepo.GetByName(role.Name)
		if err != nil {
			return err
		}
		if otherRole != nil {
			return errors.Conflict("Role name already exists", "角色名称已存在")
		}
	}

	// Update role
	return s.roleRepo.Update(role)
}

// DeleteRole deletes a role
func (s *roleService) DeleteRole(id uint) error {
	// Check if role exists
	existingRole, err := s.roleRepo.GetByID(id)
	if err != nil {
		return err
	}
	if existingRole == nil {
		return errors.NotFound("Role not found", "角色不存在")
	}

	// Delete role
	return s.roleRepo.Delete(id)
}

// ListRoles lists roles with pagination
func (s *roleService) ListRoles(page, pageSize int) ([]*model.Role, int64, error) {
	return s.roleRepo.List(page, pageSize)
}

// AssignRoleToUser assigns a role to a user
func (s *roleService) AssignRoleToUser(userID, roleID uint) error {
	// Check if user exists
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.NotFound("User not found", "用户不存在")
	}

	// Check if role exists
	role, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.NotFound("Role not found", "角色不存在")
	}

	// Check if the user-role relationship already exists
	// This would typically be done with a database query, but for simplicity,
	// we'll let the database handle duplicate constraints

	// Create user-role relationship
	userRole := &model.UserRole{
		UserID: userID,
		RoleID: roleID,
	}

	db := repository.GetDB()
	err = db.Create(userRole).Error
	if err != nil {
		// Check if it's a duplicate entry error
		// In a real implementation, you'd check the specific database error
		return errors.Conflict("Role already assigned to user", "角色已分配给该用户")
	}

	return nil
}

// RemoveRoleFromUser removes a role from a user
func (s *roleService) RemoveRoleFromUser(userID, roleID uint) error {
	// Check if user exists
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.NotFound("User not found", "用户不存在")
	}

	// Check if role exists
	role, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.NotFound("Role not found", "角色不存在")
	}

	// Delete user-role relationship
	db := repository.GetDB()
	result := db.Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&model.UserRole{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.NotFound("Role not assigned to user", "角色未分配给该用户")
	}

	return nil
}

// GetRolesByUserID gets roles by user ID
func (s *roleService) GetRolesByUserID(userID uint) ([]*model.Role, error) {
	// Check if user exists
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.NotFound("User not found", "用户不存在")
	}

	return s.roleRepo.GetRolesByUserID(userID)
}
