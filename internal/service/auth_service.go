package service

import (
	"errors"
	"time"

	"go-admin/internal/cache"
	"go-admin/internal/logger"
	"go-admin/internal/model"
	"go-admin/internal/repository"
	"go-admin/pkg/utils"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
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
	ID       string `json:"jti,omitempty"` // JWT ID for token identification and blacklisting
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
	// Validate token first to get expiration time
	token, err := s.ValidateToken(tokenString)
	if err != nil {
		// Even if token is invalid, add it to blacklist as a safety measure
		cacheInstance := cache.GetInstance()
		if err := cacheInstance.Set("blacklist:"+tokenString, true, 24*time.Hour); err != nil {
			logger.Error("Failed to add invalid token to blacklist", zap.Error(err), zap.String("token", tokenString[:min(len(tokenString), 10)]+"..."))
		}
		return nil
	}

	// Get token expiration time to set appropriate blacklist duration
	claims, ok := token.Claims.(*AuthClaims)
	if !ok || !token.Valid {
		// Add to blacklist with default duration
		cacheInstance := cache.GetInstance()
		if err := cacheInstance.Set("blacklist:"+tokenString, true, 24*time.Hour); err != nil {
			logger.Error("Failed to add invalid token to blacklist", zap.Error(err), zap.String("token", tokenString[:min(len(tokenString), 10)]+"..."))
		}
		return nil
	}

	// Calculate remaining time until token expires
	var blacklistDuration time.Duration
	if claims.ExpiresAt != nil {
		blacklistDuration = claims.ExpiresAt.Time.Sub(time.Now())
		// Add a buffer time to ensure token is definitely invalid
		blacklistDuration += 5 * time.Minute
	} else {
		// Default to 24 hours if no expiration is set
		blacklistDuration = 24 * time.Hour
	}

	// Add token to blacklist cache with calculated duration
	cacheInstance := cache.GetInstance()
	if err := cacheInstance.Set("blacklist:"+tokenString, true, blacklistDuration); err != nil {
		logger.Error("Failed to add token to blacklist", zap.Error(err), zap.String("token", tokenString[:min(len(tokenString), 10)]+"..."))
	}
	
	// Also add the JTI (JWT ID) if available to prevent token reuse
	if claims.ID != "" {
		if err := cacheInstance.Set("blacklist:jti:"+claims.ID, true, blacklistDuration); err != nil {
			logger.Error("Failed to add JTI to blacklist", zap.Error(err), zap.String("jti", claims.ID))
		}
	}

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

	// Check if JTI is in blacklist (additional security check)
	if claims.ID != "" {
		if _, exists := cacheInstance.Get("blacklist:jti:" + claims.ID); exists {
			return "", errors.New("token is invalid")
		}
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

	// Add old token to blacklist to prevent reuse
	// Calculate remaining time until old token expires
	var blacklistDuration time.Duration
	if claims.ExpiresAt != nil {
		blacklistDuration = claims.ExpiresAt.Time.Sub(time.Now())
		// Add a buffer time to ensure token is definitely invalid
		blacklistDuration += 5 * time.Minute
	} else {
		// Default to 24 hours if no expiration is set
		blacklistDuration = 24 * time.Hour
	}

	// Add old token to blacklist cache with calculated duration
	if err := cacheInstance.Set("blacklist:"+tokenString, true, blacklistDuration); err != nil {
		logger.Error("Failed to add old token to blacklist", zap.Error(err), zap.String("token", tokenString[:min(len(tokenString), 10)]+"..."))
	}
	
	// Also add the JTI (JWT ID) if available to prevent token reuse
	if claims.ID != "" {
		if err := cacheInstance.Set("blacklist:jti:"+claims.ID, true, blacklistDuration); err != nil {
			logger.Error("Failed to add JTI to blacklist", zap.Error(err), zap.String("jti", claims.ID))
		}
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
	secret, err := utils.GetJWTSecret()
	if err != nil {
		return nil, err
	}
	
	token, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

// generateToken generates JWT token for a user
func (s *authService) generateToken(user *model.User) (string, error) {
	// Generate a unique JWT ID for token identification and blacklisting
	jti := utils.GenerateUUID()
	
	claims := AuthClaims{
		UserID:   user.ID,
		Username: user.Username,
		ID:       jti, // Add JWT ID for token tracking and blacklisting
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "go-admin",
			ID:        jti, // Also set in RegisteredClaims
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret, err := utils.GetJWTSecret()
	if err != nil {
		return "", err
	}
	return token.SignedString(secret)
}
