package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Load(t *testing.T) {
	// Test loading configuration
	cfg, err := Load()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	// Check that required fields are set
	assert.NotEmpty(t, cfg.App.Name)
	assert.NotEmpty(t, cfg.DB.Host)
	assert.NotEmpty(t, cfg.JWT.Secret)
}

func TestConfig_Get(t *testing.T) {
	// Load configuration first
	_, err := Load()
	assert.NoError(t, err)

	// Get configuration
	cfg := Get()
	assert.NotNil(t, cfg)

	// Check that required fields are set
	assert.NotEmpty(t, cfg.App.Name)
	assert.NotEmpty(t, cfg.DB.Host)
	assert.NotEmpty(t, cfg.JWT.Secret)
}

func TestConfig_IsDevelopment(t *testing.T) {
	// Load configuration first
	_, err := Load()
	assert.NoError(t, err)

	// Get configuration
	cfg := Get()
	assert.NotNil(t, cfg)

	// Test IsDevelopment function
	// By default, it should be development
	assert.True(t, cfg.IsDevelopment())
}

func TestConfig_DSN(t *testing.T) {
	// Load configuration first
	_, err := Load()
	assert.NoError(t, err)

	// Get configuration
	cfg := Get()
	assert.NotNil(t, cfg)

	// Test DSN function
	dsn := cfg.DB.DSN()
	assert.NotEmpty(t, dsn)
	assert.Contains(t, dsn, cfg.DB.User)
	assert.Contains(t, dsn, cfg.DB.Password)
	assert.Contains(t, dsn, cfg.DB.Host)
	assert.Contains(t, dsn, cfg.DB.Port)
	assert.Contains(t, dsn, cfg.DB.Name)
}

func TestConfig_Validate(t *testing.T) {
	// Create a configuration with missing required fields
	cfg := &Configuration{
		App: AppConfig{
			Name: "",
		},
		DB: DBConfig{
			Host: "",
		},
		JWT: JWTConfig{
			Secret: "",
		},
	}

	// Validation should fail
	err := cfg.validate()
	assert.Error(t, err)
	// Note: Since validation stops at the first error, we can only check the first one
	assert.Contains(t, err.Error(), "app.name is required")

	// Create a valid configuration
	cfg = &Configuration{
		App: AppConfig{
			Name: "test-app",
		},
		DB: DBConfig{
			Host: "localhost",
		},
		JWT: JWTConfig{
			Secret: "test-secret",
		},
	}

	// Validation should pass
	err = cfg.validate()
	assert.NoError(t, err)
}
