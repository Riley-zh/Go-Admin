package logger

import (
	"testing"

	"go-admin/config"

	"github.com/stretchr/testify/assert"
)

func TestLogger_Init(t *testing.T) {
	// Create a test configuration
	cfg := config.LogConfig{
		Level:  "info",
		Output: "console",
	}

	// Initialize logger
	err := Init(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, logger)
	assert.NotNil(t, sugar)
}

func TestLogger_Get(t *testing.T) {
	// Create a test configuration
	cfg := config.LogConfig{
		Level:  "info",
		Output: "console",
	}

	// Initialize logger
	err := Init(cfg)
	assert.NoError(t, err)

	// Get logger
	l := Get()
	assert.NotNil(t, l)
	assert.Equal(t, logger, l)
}

func TestLogger_Sugar(t *testing.T) {
	// Create a test configuration
	cfg := config.LogConfig{
		Level:  "info",
		Output: "console",
	}

	// Initialize logger
	err := Init(cfg)
	assert.NoError(t, err)

	// Get sugared logger
	s := Sugar()
	assert.NotNil(t, s)
	assert.Equal(t, sugar, s)
}

func TestLogger_LogFunctions(t *testing.T) {
	// Create a test configuration
	cfg := config.LogConfig{
		Level:  "debug",
		Output: "console",
	}

	// Initialize logger
	err := Init(cfg)
	assert.NoError(t, err)

	// Test log functions
	Debug("Test debug message")
	Info("Test info message")
	Warn("Test warn message")
	Error("Test error message")

	// These should not panic
	assert.NotPanics(t, func() {
		Debug("Test debug message")
		Info("Test info message")
		Warn("Test warn message")
		Error("Test error message")
	})
}

func TestLogger_GetLevel(t *testing.T) {
	// Create a test configuration
	cfg := config.LogConfig{
		Level:  "info",
		Output: "console",
	}

	// Initialize logger
	err := Init(cfg)
	assert.NoError(t, err)

	// Get level
	level := GetLevel()
	assert.Equal(t, "info", level)
}

func TestLogger_SetLevel(t *testing.T) {
	// Create a test configuration
	cfg := config.LogConfig{
		Level:  "info",
		Output: "console",
	}

	// Initialize logger
	err := Init(cfg)
	assert.NoError(t, err)

	// Set level
	err = SetLevel("debug")
	assert.NoError(t, err)

	// Get level
	level := GetLevel()
	assert.Equal(t, "debug", level)

	// Test invalid level
	err = SetLevel("invalid")
	assert.Error(t, err)
}
