package database

import (
	"database/sql"
	"fmt"
	"time"

	"go-admin/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

// Stats returns database connection pool statistics
func Stats() sql.DBStats {
	if db != nil {
		sqlDB, err := db.DB()
		if err == nil {
			return sqlDB.Stats()
		}
	}
	return sql.DBStats{}
}

// Init initializes the database connection
func Init(cfg config.DBConfig) (*gorm.DB, error) {
	var err error

	// Connect to database
	db, err = gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get generic database object
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database object: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)           // Maximum number of connections in the idle connection pool
	sqlDB.SetMaxOpenConns(100)          // Maximum number of open connections to the database
	sqlDB.SetConnMaxLifetime(time.Hour) // Maximum amount of time a connection may be reused

	return db, nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return db
}

// Close closes the database connection
func Close() error {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			return fmt.Errorf("failed to get database object: %w", err)
		}
		return sqlDB.Close()
	}
	return nil
}
