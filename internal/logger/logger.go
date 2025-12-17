package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go-admin/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger       *zap.Logger
	sugar        *zap.SugaredLogger
	currentLevel zapcore.Level
)

// Init initializes the logger with the given configuration
func Init(cfg config.LogConfig) error {
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
		return fmt.Errorf("failed to parse log level: %w", err)
	}

	// Configure encoder
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var encoder zapcore.Encoder
	if cfg.Output == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// Configure writer
	var writers []zapcore.WriteSyncer

	// Console output
	if cfg.Output == "console" || cfg.Output == "both" {
		writers = append(writers, zapcore.AddSync(os.Stdout))
	}

	// File output
	if cfg.Output == "file" || cfg.Output == "both" {
		// Create logs directory if it doesn't exist
		logDir := "logs"
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return fmt.Errorf("failed to create logs directory: %w", err)
		}

		logFile := filepath.Join(logDir, fmt.Sprintf("%s.log", time.Now().Format("2006-01-02")))
		lumberjackLogger := &lumberjack.Logger{
			Filename:   logFile,
			MaxSize:    100, // megabytes
			MaxAge:     7,   // days
			MaxBackups: 3,
			LocalTime:  true,
			Compress:   true,
		}
		writers = append(writers, zapcore.AddSync(lumberjackLogger))
	}

	// If no output is specified, default to console
	if len(writers) == 0 {
		writers = append(writers, zapcore.AddSync(os.Stdout))
	}

	writeSyncer := zapcore.NewMultiWriteSyncer(writers...)

	// Save current level
	currentLevel = level

	// Create core
	core := zapcore.NewCore(encoder, writeSyncer, level)

	// Create logger with options
	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	sugar = logger.Sugar()

	return nil
}

// Sync flushes any buffered log entries
func Sync() {
	if logger != nil {
		logger.Sync()
	}
	if sugar != nil {
		sugar.Sync()
	}
}

// Get returns the underlying zap logger
func Get() *zap.Logger {
	return logger
}

// Sugar returns the sugared logger
func Sugar() *zap.SugaredLogger {
	return sugar
}

// DefaultStructuredLogger creates a structured logger with default fields
func DefaultStructuredLogger() *StructuredLogger {
	return &StructuredLogger{
		logger: logger,
		fields: make(map[string]interface{}),
	}
}

// GetLevel returns the current log level
func GetLevel() string {
	return currentLevel.String()
}

// SetLevel sets the log level dynamically
func SetLevel(levelStr string) error {
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(levelStr)); err != nil {
		return fmt.Errorf("failed to parse log level: %w", err)
	}

	currentLevel = level
	return nil
}

// Debug logs a debug message
func Debug(msg string, fields ...zap.Field) {
	if logger != nil {
		logger.Debug(msg, fields...)
	}
}

// Info logs an info message
func Info(msg string, fields ...zap.Field) {
	if logger != nil {
		logger.Info(msg, fields...)
	}
}

// Warn logs a warning message
func Warn(msg string, fields ...zap.Field) {
	if logger != nil {
		logger.Warn(msg, fields...)
	}
}

// Error logs an error message
func Error(msg string, fields ...zap.Field) {
	if logger != nil {
		logger.Error(msg, fields...)
	}
}

// Fatal logs a fatal message and exits
func Fatal(msg string, fields ...zap.Field) {
	if logger != nil {
		logger.Fatal(msg, fields...)
	}
}

// GinLogger returns a gin.HandlerFunc for logging requests
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		comment := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if raw != "" {
			path = path + "?" + raw
		}

		fields := []zap.Field{
			zap.Int("status", statusCode),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("ip", clientIP),
			zap.Duration("latency", latency),
		}

		if len(comment) > 0 {
			fields = append(fields, zap.String("error", comment))
		}

		switch {
		case statusCode >= 500:
			Error("Server error", fields...)
		case statusCode >= 400:
			Warn("Client error", fields...)
		default:
			Info("Request", fields...)
		}
	}
}
