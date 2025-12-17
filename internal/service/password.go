package service

import (
	"golang.org/x/crypto/bcrypt"
)

// hashPassword hashes a password using bcrypt
func hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// checkPassword compares a hashed password with a plaintext password
func checkPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
