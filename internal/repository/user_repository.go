package repository

import (
	"errors"

	"go-admin/internal/database"
	"go-admin/internal/model"

	"gorm.io/gorm"
)

// UserRepository defines the user repository interface
type UserRepository interface {
	Create(user *model.User) error
	GetByID(id uint) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	Update(user *model.User) error
	Delete(id uint) error
	List(page, pageSize int) ([]*model.User, int64, error)
	ListWithRoles(page, pageSize int) ([]*model.UserWithRoles, int64, error)
}

// userRepository implements UserRepository interface
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository() UserRepository {
	return &userRepository{
		db: database.GetDB(),
	}
}

// Create creates a new user
func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// GetByID gets a user by ID
func (r *userRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.Where("id = ? AND status = ?", id, 1).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetByUsername gets a user by username
func (r *userRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ? AND status = ?", username, 1).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetByEmail gets a user by email
func (r *userRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ? AND status = ?", email, 1).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Update updates a user
func (r *userRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// Delete deletes a user (soft delete)
func (r *userRepository) Delete(id uint) error {
	return r.db.Where("id = ?", id).Delete(&model.User{}).Error
}

// List lists users with pagination
func (r *userRepository) List(page, pageSize int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	offset := (page - 1) * pageSize
	err := r.db.Model(&model.User{}).Where("status = ?", 1).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Where("status = ?", 1).Offset(offset).Limit(pageSize).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// ListWithRoles lists users with their roles using a single query to prevent N+1 problem
func (r *userRepository) ListWithRoles(page, pageSize int) ([]*model.UserWithRoles, int64, error) {
	var usersWithRoles []*model.UserWithRoles
	var total int64

	offset := (page - 1) * pageSize
	
	// Count total users
	err := r.db.Model(&model.User{}).Where("status = ?", 1).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get users with their roles in a single query using left join
	err = r.db.Table("users").
		Select("users.id, users.created_at, users.updated_at, users.deleted_at, users.username, users.email, users.nickname, users.avatar, users.status").
		Where("users.status = ?", 1).
		Offset(offset).
		Limit(pageSize).
		Scan(&usersWithRoles).Error
	if err != nil {
		return nil, 0, err
	}

	// If no users found, return empty slice
	if len(usersWithRoles) == 0 {
		return usersWithRoles, total, nil
	}

	// Get all user IDs
	userIDs := make([]uint, len(usersWithRoles))
	for i, user := range usersWithRoles {
		userIDs[i] = user.ID
	}

	// Get all roles for these users in a single query
	var userRoles []struct {
		UserID uint `json:"user_id"`
		RoleID uint `json:"role_id"`
		model.Role
	}

	err = r.db.Table("user_roles").
		Select("user_roles.user_id, user_roles.role_id, roles.*").
		Joins("LEFT JOIN roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id IN ? AND roles.status = ?", userIDs, 1).
		Scan(&userRoles).Error
	if err != nil {
		return nil, 0, err
	}

	// Create a map of user ID to roles
	userRoleMap := make(map[uint][]*model.Role)
	for _, ur := range userRoles {
		role := &model.Role{
			ID:          ur.Role.ID,
			CreatedAt:   ur.Role.CreatedAt,
			UpdatedAt:   ur.Role.UpdatedAt,
			Name:        ur.Role.Name,
			Description: ur.Role.Description,
			Status:      ur.Role.Status,
		}
		userRoleMap[ur.UserID] = append(userRoleMap[ur.UserID], role)
	}

	// Assign roles to users
	for _, user := range usersWithRoles {
		if roles, exists := userRoleMap[user.ID]; exists {
			user.Roles = roles
		}
	}

	return usersWithRoles, total, nil
}
