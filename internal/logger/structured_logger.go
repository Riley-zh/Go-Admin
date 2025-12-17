package logger

import (
	"context"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// StructuredLogger provides enhanced structured logging capabilities
type StructuredLogger struct {
	logger *zap.Logger
	fields map[string]interface{}
}

// NewStructuredLogger creates a new structured logger with optional initial fields
func NewStructuredLogger(fields ...zap.Field) *StructuredLogger {
	return &StructuredLogger{
		logger: Get(),
		fields: make(map[string]interface{}),
	}
}

// WithField adds a field to the logger
func (l *StructuredLogger) WithField(key string, value interface{}) *StructuredLogger {
	l.fields[key] = value
	return l
}

// WithFields adds multiple fields to the logger
func (l *StructuredLogger) WithFields(fields map[string]interface{}) *StructuredLogger {
	for k, v := range fields {
		l.fields[k] = v
	}
	return l
}

// WithContext adds context information to the logger
func (l *StructuredLogger) WithContext(ctx context.Context) *StructuredLogger {
	// Extract common context values
	if requestID := ctx.Value("requestID"); requestID != nil {
		l.fields["requestID"] = requestID
	}
	if userID := ctx.Value("userID"); userID != nil {
		l.fields["userID"] = userID
	}
	if traceID := ctx.Value("traceID"); traceID != nil {
		l.fields["traceID"] = traceID
	}
	return l
}

// WithError adds an error field to the logger
func (l *StructuredLogger) WithError(err error) *StructuredLogger {
	if err != nil {
		l.fields["error"] = err.Error()
	}
	return l
}

// WithDuration adds a duration field to the logger
func (l *StructuredLogger) WithDuration(d time.Duration) *StructuredLogger {
	l.fields["duration"] = d.String()
	return l
}

// Debug logs a debug message with the accumulated fields
func (l *StructuredLogger) Debug(msg string) {
	l.log(zap.DebugLevel, msg)
}

// Info logs an info message with the accumulated fields
func (l *StructuredLogger) Info(msg string) {
	l.log(zap.InfoLevel, msg)
}

// Warn logs a warning message with the accumulated fields
func (l *StructuredLogger) Warn(msg string) {
	l.log(zap.WarnLevel, msg)
}

// Error logs an error message with the accumulated fields
func (l *StructuredLogger) Error(msg string) {
	l.log(zap.ErrorLevel, msg)
}

// Fatal logs a fatal message with the accumulated fields and exits
func (l *StructuredLogger) Fatal(msg string) {
	l.log(zap.FatalLevel, msg)
}

// log performs the actual logging with the accumulated fields
func (l *StructuredLogger) log(level zapcore.Level, msg string) {
	if l.logger == nil {
		return
	}

	// Convert map to zap fields
	fields := make([]zap.Field, 0, len(l.fields))
	for k, v := range l.fields {
		switch val := v.(type) {
		case string:
			fields = append(fields, zap.String(k, val))
		case int:
			fields = append(fields, zap.Int(k, val))
		case int64:
			fields = append(fields, zap.Int64(k, val))
		case float64:
			fields = append(fields, zap.Float64(k, val))
		case bool:
			fields = append(fields, zap.Bool(k, val))
		case time.Duration:
			fields = append(fields, zap.Duration(k, val))
		case time.Time:
			fields = append(fields, zap.Time(k, val))
		case error:
			fields = append(fields, zap.Error(val))
		default:
			fields = append(fields, zap.Any(k, val))
		}
	}

	// Log with the appropriate level
	switch level {
	case zapcore.DebugLevel:
		l.logger.Debug(msg, fields...)
	case zapcore.InfoLevel:
		l.logger.Info(msg, fields...)
	case zapcore.WarnLevel:
		l.logger.Warn(msg, fields...)
	case zapcore.ErrorLevel:
		l.logger.Error(msg, fields...)
	case zapcore.FatalLevel:
		l.logger.Fatal(msg, fields...)
	}
}

// Reset clears all accumulated fields
func (l *StructuredLogger) Reset() *StructuredLogger {
	l.fields = make(map[string]interface{})
	return l
}

// Clone creates a new StructuredLogger with the same fields
func (l *StructuredLogger) Clone() *StructuredLogger {
	newLogger := &StructuredLogger{
		logger: l.logger,
		fields: make(map[string]interface{}),
	}
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}
	return newLogger
}

// GetFields returns a copy of the current fields
func (l *StructuredLogger) GetFields() map[string]interface{} {
	fields := make(map[string]interface{}, len(l.fields))
	for k, v := range l.fields {
		fields[k] = v
	}
	return fields
}