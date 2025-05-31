package types

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidConfig represents configuration validation errors
	ErrInvalidConfig = errors.New("invalid configuration")

	// ErrUnauthorized represents authentication errors
	ErrUnauthorized = errors.New("unauthorized: invalid API key")

	// ErrRateLimit represents rate limiting errors
	ErrRateLimit = errors.New("rate limit exceeded")

	// ErrServerError represents internal server errors
	ErrServerError = errors.New("postal server error")

	// ErrInvalidMessage represents message validation errors
	ErrInvalidMessage = errors.New("invalid message")
)

// PostalError represents a detailed API error
type PostalError struct {
	Code       string                 `json:"code"`
	Message    string                 `json:"message"`
	Details    map[string]interface{} `json:"details,omitempty"`
	StatusCode int                    `json:"-"`
}

// Error implements the error interface
func (e *PostalError) Error() string {
	if len(e.Details) > 0 {
		return fmt.Sprintf("%s: %s (details: %v)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// IsRateLimit checks if the error is a rate limit error
func IsRateLimit(err error) bool {
	return errors.Is(err, ErrRateLimit)
}

// IsUnauthorized checks if the error is an authentication error
func IsUnauthorized(err error) bool {
	return errors.Is(err, ErrUnauthorized)
}

// IsServerError checks if the error is a server error
func IsServerError(err error) bool {
	return errors.Is(err, ErrServerError)
}

// NewPostalError creates a new PostalError with the given details
func NewPostalError(code string, message string, statusCode int) *PostalError {
	return &PostalError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Details:    make(map[string]interface{}),
	}
}

// WithDetails adds details to the error
func (e *PostalError) WithDetails(details map[string]interface{}) *PostalError {
	e.Details = details
	return e
}
