package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

// GetJWTSecret returns the JWT secret from environment variables
// In production, JWT_SECRET must be set and meet security requirements
func GetJWTSecret() ([]byte, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}
	
	// Validate secret strength
	if err := validateJWTSecret(secret); err != nil {
		return nil, fmt.Errorf("JWT_SECRET does not meet security requirements: %v", err)
	}
	
	return []byte(secret), nil
}

// validateJWTSecret checks if the JWT secret meets security requirements
func validateJWTSecret(secret string) error {
	// Check minimum length (at least 32 characters)
	if len(secret) < 32 {
		return fmt.Errorf("JWT secret must be at least 32 characters long")
	}
	
	// Check for complexity - must contain uppercase, lowercase, numbers, and special characters
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(secret)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(secret)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(secret)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(secret)
	
	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return fmt.Errorf("JWT secret must contain uppercase, lowercase, numbers, and special characters")
	}
	
	// Check for common weak secrets
	weakSecrets := []string{
		"password", "secret", "admin", "123456", "qwerty", 
		"letmein", "welcome", "monkey", "dragon", "football",
	}
	secretLower := regexp.MustCompile(`[a-z]`).ReplaceAllString(strings.ToLower(secret), "")
	for _, weak := range weakSecrets {
		if strings.Contains(secretLower, weak) {
			return fmt.Errorf("JWT secret contains common weak patterns")
		}
	}
	
	return nil
}

// GenerateSecureJWTSecret generates a cryptographically secure JWT secret
func GenerateSecureJWTSecret() (string, error) {
	// Generate 64 bytes of random data
	bytes := make([]byte, 64)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	
	// Encode to base64 URL-safe string
	secret := base64.URLEncoding.EncodeToString(bytes)
	
	return secret, nil
}

// GenerateUUID generates a UUID v4
func GenerateUUID() string {
	return uuid.New().String()
}
