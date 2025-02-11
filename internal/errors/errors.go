package errors

import (
	"fmt"
	"net/http"
)

// AppError represents an application-specific error
type AppError struct {
	Message string
	Code    int
	Err     error
}

// Error returns the error message
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// ValidationError returns a new AppError for validation failures
func ValidationError(message string) *AppError {
	return &AppError{
		Message: message,
		Code:    http.StatusBadRequest,
	}
}

// NotFoundError returns a new AppError for resource not found
func NotFoundError(message string) *AppError {
	return &AppError{
		Message: message,
		Code:    http.StatusNotFound,
	}
}

// ConflictError returns a new AppError for conflicting resources
func ConflictError(message string) *AppError {
	return &AppError{
		Message: message,
		Code:    http.StatusConflict,
	}
}

// InternalError wraps internal server errors
func InternalError(message string, err error) *AppError {
	return &AppError{
		Message: message,
		Code:    http.StatusInternalServerError,
		Err:     err,
	}
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) (*AppError, bool) {
	appErr, ok := err.(*AppError)
	return appErr, ok
}
