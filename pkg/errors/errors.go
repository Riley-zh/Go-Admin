package errors

import (
	"net/http"
)

// Error represents an application error
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Error implements the error interface
func (e *Error) Error() string {
	return e.Message
}

// New creates a new error
func New(code int, message, details string) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// BadRequest creates a 400 Bad Request error
func BadRequest(message, details string) *Error {
	return &Error{
		Code:    http.StatusBadRequest,
		Message: message,
		Details: details,
	}
}

// Unauthorized creates a 401 Unauthorized error
func Unauthorized(message, details string) *Error {
	return &Error{
		Code:    http.StatusUnauthorized,
		Message: message,
		Details: details,
	}
}

// Forbidden creates a 403 Forbidden error
func Forbidden(message, details string) *Error {
	return &Error{
		Code:    http.StatusForbidden,
		Message: message,
		Details: details,
	}
}

// NotFound creates a 404 Not Found error
func NotFound(message, details string) *Error {
	return &Error{
		Code:    http.StatusNotFound,
		Message: message,
		Details: details,
	}
}

// InternalServerError creates a 500 Internal Server Error
func InternalServerError(message, details string) *Error {
	return &Error{
		Code:    http.StatusInternalServerError,
		Message: message,
		Details: details,
	}
}

// Conflict creates a 409 Conflict error
func Conflict(message, details string) *Error {
	return &Error{
		Code:    http.StatusConflict,
		Message: message,
		Details: details,
	}
}
