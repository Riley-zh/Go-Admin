package service

import (
	"go-admin/internal/model"
	"go-admin/internal/repository"
	"go-admin/pkg/errors"
)

// UserService defines the user service interface
type UserService interface {
	CreateUser(username, password, email, nickname string) (*model.User, error)
	GetUserByID(id uint) (*model.User, error)
	GetUserByUsername(username string) (*model.User, error)
	UpdateUser(user *model.User) error
	DeleteUser(id uint) error
	ListUsers(page, pageSize int) ([]*model.User, int64, error)
	ListUsersWithRoles(page, pageSize int) ([]*model.UserWithRoles, int64, error)
	ChangePassword(userID uint, oldPassword, newPassword string) error
}

// userService implements UserService interface
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService() UserService {
	return &userService{
		userRepo: repository.NewUserRepository(),
	}
}

// CreateUser creates a new user
func (s *userService) CreateUser(username, password, email, nickname string) (*model.User, error) {
	// Check if username already exists
	existingUser, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.Conflict("Username already exists", "用户名已存在")
	}

	// Check if email already exists
	existingUser, err = s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.Conflict("Email already exists", "邮箱已存在")
	}

	// Hash password
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &model.User{
		Username: username,
		Password: string(hashedPassword),
		Email:    email,
		Nickname: nickname,
		Status:   1,
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	// Hide password in response
	user.Password = ""

	return user, nil
}

// GetUserByID gets a user by ID
func (s *userService) GetUserByID(id uint) (*model.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.NotFound("User not found", "用户不存在")
	}

	// Hide password
	user.Password = ""

	return user, nil
}

// GetUserByUsername gets a user by username
func (s *userService) GetUserByUsername(username string) (*model.User, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.NotFound("User not found", "用户不存在")
	}

	// Hide password
	user.Password = ""

	return user, nil
}

// UpdateUser updates a user
func (s *userService) UpdateUser(user *model.User) error {
	// Check if user exists
	existingUser, err := s.userRepo.GetByID(user.ID)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return errors.NotFound("User not found", "用户不存在")
	}

	// Update user
	return s.userRepo.Update(user)
}

// DeleteUser deletes a user
func (s *userService) DeleteUser(id uint) error {
	// Check if user exists
	existingUser, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return errors.NotFound("User not found", "用户不存在")
	}

	// Delete user
	return s.userRepo.Delete(id)
}

// ListUsers lists users with pagination
func (s *userService) ListUsers(page, pageSize int) ([]*model.User, int64, error) {
	users, total, err := s.userRepo.List(page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// Hide passwords
	for _, user := range users {
		user.Password = ""
	}

	return users, total, nil
}

// ListUsersWithRoles lists users with their roles using optimized queries to prevent N+1 problem
func (s *userService) ListUsersWithRoles(page, pageSize int) ([]*model.UserWithRoles, int64, error) {
	usersWithRoles, total, err := s.userRepo.ListWithRoles(page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// Hide passwords
	for _, user := range usersWithRoles {
		user.Password = ""
	}

	return usersWithRoles, total, nil
}

// ChangePassword changes user password
func (s *userService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	// Get user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.NotFound("User not found", "用户不存在")
	}

	// Check old password
	if !checkPassword(user.Password, oldPassword) {
		return errors.BadRequest("Invalid old password", "原密码错误")
	}

	// Hash new password
	hashedPassword, err := hashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update password
	user.Password = string(hashedPassword)
	return s.userRepo.Update(user)
}
