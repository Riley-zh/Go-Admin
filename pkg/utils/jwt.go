package utils

import (
	"os"
)

// GetJWTSecret returns the JWT secret from environment variables or default value
func GetJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "go-admin-secret-default"
	}
	return []byte(secret)
}
