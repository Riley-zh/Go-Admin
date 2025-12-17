package repository

import (
	"go-admin/internal/database"

	"gorm.io/gorm"
)

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return database.GetDB()
}
