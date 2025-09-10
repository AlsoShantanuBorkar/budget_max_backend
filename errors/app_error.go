package errors

import (
	"fmt"
	"net/http"
)

// AppError defines a standardized error structure for the application
// with a code, message, and (optionally) an underlying error.
type AppError struct {
    Code    int // e.g. "USER_NOT_FOUND", "DB_ERROR", "VALIDATION_ERROR"
    Message string
    Err     error
}

func (e *AppError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
    }
    return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
    return e.Err
}

// Factory functions for common error codes
func NewNotFoundError(entity string, err error) *AppError {
    return &AppError{
        Code:    http.StatusNotFound,
        Message: fmt.Sprintf("%s not found", entity),
        Err:     err,
    }
}

func NewDBError(err error) *AppError {
    return &AppError{
        Code:    http.StatusInternalServerError,
        Message: "Database error",
        Err:     err,
    }
}

func NewValidationError(msg string, err error) *AppError {
    return &AppError{
        Code:    http.StatusBadRequest,
        Message: msg,
        Err:     err,
    }
}

func NewUnauthorizedError(msg string, err error) *AppError {
    return &AppError{
        Code:    http.StatusUnauthorized,
        Message: msg,
        Err:     err,
    }
}

func NewConflictError(msg string, err error) *AppError {
    return &AppError{
        Code:    http.StatusConflict,
        Message: msg,
        Err:     err,
    }
}

func NewInternalError(err error) *AppError {
    return &AppError{
        Code:    http.StatusInternalServerError,
        Message: "Internal Server Error",
        Err:     err,
    }
}

func NewBadRequestError(msg string, err error) *AppError {
    return &AppError{
        Code:    http.StatusBadRequest,
        Message: msg,
        Err:     err,
    }
}

func NewForbiddenError(msg string, err error) *AppError {
    return &AppError{
        Code:    http.StatusForbidden,
        Message: msg,
        Err:     err,
    }
}

func NewTooManyRequestsError(msg string, err error) *AppError {
    return &AppError{
        Code:    http.StatusTooManyRequests,
        Message: msg,
        Err:     err,
    }
}