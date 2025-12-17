package service

import (
	"errors"
	"go-admin/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(id uint) (*model.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(username string) (*model.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (*model.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) List(page, pageSize int) ([]*model.User, int64, error) {
	args := m.Called(page, pageSize)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*model.User), int64(args.Int(1)), args.Error(2)
}

func (m *MockUserRepository) ListWithRoles(page, pageSize int) ([]*model.UserWithRoles, int64, error) {
	args := m.Called(page, pageSize)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*model.UserWithRoles), int64(args.Int(1)), args.Error(2)
}

func TestUserService_CreateUser(t *testing.T) {
	// Create a mock user repository
	mockRepo := new(MockUserRepository)

	// Create a user service with the mock repository
	userService := &userService{
		userRepo: mockRepo,
	}

	// Test creating a user with valid data
	mockRepo.On("GetByUsername", "testuser").Return(nil, nil).Once()
	mockRepo.On("GetByEmail", "test@example.com").Return(nil, nil).Once()
	mockRepo.On("Create", mock.AnythingOfType("*model.User")).Return(nil).Once()

	user, err := userService.CreateUser("testuser", "testpassword", "test@example.com", "Test User")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "Test User", user.Nickname)
	// Password should be hashed, so it shouldn't be the plain text password
	assert.NotEqual(t, "testpassword", user.Password)

	// Test creating a user with duplicate username
	mockRepo.On("GetByUsername", "testuser_dup").Return(&model.User{ID: 1}, nil).Once()

	_, err = userService.CreateUser("testuser_dup", "testpassword2", "test2@example.com", "Test User 2")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Username already exists")

	// Test creating a user with duplicate email
	mockRepo.On("GetByUsername", "testuser_email").Return(nil, nil).Once()
	mockRepo.On("GetByEmail", "test@example.com").Return(&model.User{ID: 1}, nil).Once()

	_, err = userService.CreateUser("testuser_email", "testpassword", "test@example.com", "Test User 2")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Email already exists")

	// Ensure all expectations were met
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUserByID(t *testing.T) {
	// Create a mock user repository
	mockRepo := new(MockUserRepository)

	// Create a user service with the mock repository
	userService := &userService{
		userRepo: mockRepo,
	}

	// Test getting an existing user
	expectedUser := &model.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Nickname: "Test User",
		Password: "hashed_password",
	}
	mockRepo.On("GetByID", uint(1)).Return(expectedUser, nil).Once()

	user, err := userService.GetUserByID(1)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, uint(1), user.ID)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "Test User", user.Nickname)
	assert.Equal(t, "", user.Password) // Password should be hidden

	// Test getting a non-existent user
	mockRepo.On("GetByID", uint(999)).Return((*model.User)(nil), nil).Once()

	nonExistentUser, err := userService.GetUserByID(999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "User not found")
	assert.Nil(t, nonExistentUser)

	// Ensure all expectations were met
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUserByUsername(t *testing.T) {
	// Create a mock user repository
	mockRepo := new(MockUserRepository)

	// Create a user service with the mock repository
	userService := &userService{
		userRepo: mockRepo,
	}

	// Test getting an existing user by username
	expectedUser := &model.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Nickname: "Test User",
		Password: "hashed_password",
	}
	mockRepo.On("GetByUsername", "testuser").Return(expectedUser, nil).Once()

	user, err := userService.GetUserByUsername("testuser")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, uint(1), user.ID)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "Test User", user.Nickname)
	assert.Equal(t, "", user.Password) // Password should be hidden

	// Test getting a non-existent user
	mockRepo.On("GetByUsername", "nonexistent").Return((*model.User)(nil), nil).Once()

	nonExistentUser, err := userService.GetUserByUsername("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "User not found")
	assert.Nil(t, nonExistentUser)

	// Ensure all expectations were met
	mockRepo.AssertExpectations(t)
}

func TestUserService_UpdateUser(t *testing.T) {
	// Create a mock user repository
	mockRepo := new(MockUserRepository)

	// Create a user service with the mock repository
	userService := &userService{
		userRepo: mockRepo,
	}

	// Test updating an existing user
	existingUser := &model.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Nickname: "Test User",
	}
	mockRepo.On("GetByID", uint(1)).Return(existingUser, nil).Once()
	mockRepo.On("Update", mock.AnythingOfType("*model.User")).Return(nil).Once()

	userToUpdate := &model.User{
		ID:       1,
		Nickname: "Updated User",
		Avatar:   "updated_avatar_url",
	}
	err := userService.UpdateUser(userToUpdate)
	assert.NoError(t, err)

	// Test updating a non-existent user
	mockRepo.On("GetByID", uint(999)).Return((*model.User)(nil), nil).Once()

	fakeUser := &model.User{ID: 999}
	err = userService.UpdateUser(fakeUser)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "User not found")

	// Ensure all expectations were met
	mockRepo.AssertExpectations(t)
}

func TestUserService_ChangePassword(t *testing.T) {
	// Create a mock user repository
	mockRepo := new(MockUserRepository)

	// Create a user service with the mock repository
	userService := &userService{
		userRepo: mockRepo,
	}

	// Test changing password with correct old password
	// Using a correct bcrypt hash for "testpassword"
	existingUser := &model.User{
		ID:       1,
		Username: "testuser",
		Password: "$2a$10$GbkecLrL1LXY8YkiEL4jT.cviQcnthMcy8iOGhccEjSk2j2XRoAAG", // hashed "testpassword"
	}
	mockRepo.On("GetByID", uint(1)).Return(existingUser, nil).Once()
	mockRepo.On("Update", mock.AnythingOfType("*model.User")).Return(nil).Once()

	err := userService.ChangePassword(1, "testpassword", "newpassword")
	assert.NoError(t, err)

	// Test changing password with incorrect old password
	mockRepo.On("GetByID", uint(2)).Return(existingUser, nil).Once()

	err = userService.ChangePassword(2, "wrongpassword", "newpassword")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid old password")

	// Test changing password for non-existent user
	mockRepo.On("GetByID", uint(999)).Return((*model.User)(nil), nil).Once()

	err = userService.ChangePassword(999, "oldpassword", "newpassword")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "User not found")

	// Test database error
	mockRepo.On("GetByID", uint(3)).Return((*model.User)(nil), errors.New("database error")).Once()

	err = userService.ChangePassword(3, "oldpassword", "newpassword")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")

	// Ensure all expectations were met
	mockRepo.AssertExpectations(t)
}
