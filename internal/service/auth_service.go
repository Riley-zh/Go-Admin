package service

import (
	"errors"
	"time"

	"go-admin/internal/cache"
	"go-admin/internal/model"
	"go-admin/internal/repository"
	"go-admin/pkg/utils"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthService defines the auth service interface
type AuthService interface {
	Register(username, password, email, nickname string) (*model.User, error)
	Login(username, password string) (string, *model.User, error)
	Logout(tokenString string) error
	RefreshToken(tokenString string) (string, error)
	GetUserByToken(tokenString string) (*model.User, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
}

// authService implements AuthService interface
type authService struct {
	userRepo repository.UserRepository
}

// AuthClaims represents the claims in JWT token
type AuthClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// NewAuthService creates a new auth service
func NewAuthService() AuthService {
	return &authService{
		userRepo: repository.NewUserRepository(),
	}
}

// Register registers a new user
func (s *authService) Register(username, password, email, nickname string) (*model.User, error) {
	// Check if username already exists
	existingUser, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists
	existingUser, err = s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
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

// Login authenticates a user and generates JWT token
func (s *authService) Login(username, password string) (string, *model.User, error) {
	// Get user by username
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return "", nil, err
	}
	if user == nil {
		return "", nil, errors.New("invalid username or password")
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", nil, errors.New("invalid username or password")
	}

	// Generate JWT token
	token, err := s.generateToken(user)
	if err != nil {
		return "", nil, err
	}

	// Hide password in response
	user.Password = ""

	return token, user, nil
}

// Logout invalidates the JWT token
func (s *authService) Logout(tokenString string) error {
	// Add token to blacklist cache
	cacheInstance := cache.GetInstance()
	cacheInstance.Set("blacklist:"+tokenString, true, 24*time.Hour)
	return nil
}

// RefreshToken generates a new token based on the old token
func (s *authService) RefreshToken(tokenString string) (string, error) {
	// Validate token
	token, err := s.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	// Check if token is in blacklist
	cacheInstance := cache.GetInstance()
	if _, exists := cacheInstance.Get("blacklist:" + tokenString); exists {
		return "", errors.New("token is invalid")
	}

	// Extract claims
	claims, ok := token.Claims.(*AuthClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token")
	}

	// Get user
	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("user not found")
	}

	// Generate new token
	newToken, err := s.generateToken(user)
	if err != nil {
		return "", err
	}

	return newToken, nil
}

// GetUserByToken gets user by token
func (s *authService) GetUserByToken(tokenString string) (*model.User, error) {
	// Validate token
	token, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Check if token is in blacklist
	cacheInstance := cache.GetInstance()
	if _, exists := cacheInstance.Get("blacklist:" + tokenString); exists {
		return nil, errors.New("token is invalid")
	}

	// Extract claims
	claims, ok := token.Claims.(*AuthClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Get user
	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Hide password
	user.Password = ""

	return user, nil
}

// ValidateToken validates JWT token
func (s *authService) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		return utils.GetJWTSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

// generateToken generates JWT token for a user
func (s *authService) generateToken(user *model.User) (string, error) {
	claims := AuthClaims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "go-admin",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(utils.GetJWTSecret())
}
