// Package customerrors defines the custom error types used across the
// repository, service and handler layers so HTTP status codes can be
// derived in one central place.

package customerrors

import (
	"errors"
	"fmt"
)

// ErrNotFound is the sentinel error for missing resources. Concrete
// NotFoundError values match it via errors.Is.
var ErrNotFound = errors.New("resource not found")

// NotFoundError reports that a requested resource does not exist (HTTP 404).
type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string { return e.Message }

// Is lets errors.Is(err, ErrNotFound) succeed for any NotFoundError.
func (e *NotFoundError) Is(target error) bool { return target == ErrNotFound }

// NotFoundf builds a NotFoundError with a formatted message.
func NotFoundf(format string, args ...any) *NotFoundError {
	return &NotFoundError{Message: fmt.Sprintf(format, args...)}
}

// ValidationError reports invalid client input (HTTP 400 bad request).
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string { return e.Message }

// Validationf builds a ValidationError with a formatted message.
func Validationf(format string, args ...any) *ValidationError {
	return &ValidationError{Message: fmt.Sprintf(format, args...)}
}

// ConflictError reports an operation blocked by existing relationships,
// e.g. deleting a genre that still has movies (HTTP 400 per spec).
type ConflictError struct {
	Message string
}

func (e *ConflictError) Error() string { return e.Message }

// Conflictf builds a ConflictError with a formatted message.
func Conflictf(format string, args ...any) *ConflictError {
	return &ConflictError{Message: fmt.Sprintf(format, args...)}
}
