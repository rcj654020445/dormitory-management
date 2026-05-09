// Package types defines error codes and error types.
// Layer 0: No internal dependencies allowed.
package types

import (
	"fmt"
)

// Error codes for the application.
const (
	ErrCodeOK           = 0
	ErrCodeBadRequest   = 400
	ErrCodeUnauthorized = 401
	ErrCodeForbidden    = 403
	ErrCodeNotFound     = 404
	ErrCodeConflict     = 409
	ErrCodeInternal     = 500
)

// Error codes as strings for API responses.
const (
	ErrCodeStrOK           = "OK"
	ErrCodeStrBadRequest   = "BAD_REQUEST"
	ErrCodeStrUnauthorized = "UNAUTHORIZED"
	ErrCodeStrForbidden    = "FORBIDDEN"
	ErrCodeStrNotFound     = "NOT_FOUND"
	ErrCodeStrConflict     = "CONFLICT"
	ErrCodeStrInternal     = "INTERNAL_ERROR"
)

// AppError represents an application error with code and message.
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s", e.Message, e.Details)
	}
	return e.Message
}

// NewBadRequestError creates a bad request error.
func NewBadRequestError(details string) *AppError {
	return &AppError{Code: ErrCodeBadRequest, Message: ErrCodeStrBadRequest, Details: details}
}

// NewUnauthorizedError creates an unauthorized error.
func NewUnauthorizedError(details string) *AppError {
	return &AppError{Code: ErrCodeUnauthorized, Message: ErrCodeStrUnauthorized, Details: details}
}

// NewForbiddenError creates a forbidden error.
func NewForbiddenError(details string) *AppError {
	return &AppError{Code: ErrCodeForbidden, Message: ErrCodeStrForbidden, Details: details}
}

// NewNotFoundError creates a not found error.
func NewNotFoundError(resource string) *AppError {
	return &AppError{Code: ErrCodeNotFound, Message: ErrCodeStrNotFound, Details: resource + " not found"}
}

// NewConflictError creates a conflict error.
func NewConflictError(resource string) *AppError {
	return &AppError{Code: ErrCodeConflict, Message: ErrCodeStrConflict, Details: resource + " already exists"}
}

// NewInternalError creates an internal server error.
func NewInternalError(details string) *AppError {
	return &AppError{Code: ErrCodeInternal, Message: ErrCodeStrInternal, Details: details}
}
