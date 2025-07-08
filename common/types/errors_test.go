package types

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestPostalError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *PostalError
		expected string
	}{
		{
			name: "error without details",
			err: &PostalError{
				Code:    "validation_error",
				Message: "Invalid email format",
			},
			expected: "validation_error: Invalid email format",
		},
		{
			name: "error with details",
			err: &PostalError{
				Code:    "validation_error",
				Message: "Invalid email format",
				Details: map[string]interface{}{
					"field": "email",
					"value": "invalid-email",
				},
			},
			expected: "validation_error: Invalid email format (details: map[field:email value:invalid-email])",
		},
		{
			name: "error with empty details",
			err: &PostalError{
				Code:    "server_error",
				Message: "Internal server error",
				Details: map[string]interface{}{},
			},
			expected: "server_error: Internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("PostalError.Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNewPostalError(t *testing.T) {
	tests := []struct {
		name       string
		code       string
		message    string
		statusCode int
	}{
		{
			name:       "validation error",
			code:       "validation_error",
			message:    "Invalid request",
			statusCode: 400,
		},
		{
			name:       "unauthorized error",
			code:       "unauthorized",
			message:    "Invalid API key",
			statusCode: 401,
		},
		{
			name:       "rate limit error",
			code:       "rate_limit",
			message:    "Too many requests",
			statusCode: 429,
		},
		{
			name:       "server error",
			code:       "server_error",
			message:    "Internal error",
			statusCode: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewPostalError(tt.code, tt.message, tt.statusCode)
			
			if err.Code != tt.code {
				t.Errorf("NewPostalError() Code = %v, want %v", err.Code, tt.code)
			}
			if err.Message != tt.message {
				t.Errorf("NewPostalError() Message = %v, want %v", err.Message, tt.message)
			}
			if err.StatusCode != tt.statusCode {
				t.Errorf("NewPostalError() StatusCode = %v, want %v", err.StatusCode, tt.statusCode)
			}
			if err.Details == nil {
				t.Error("NewPostalError() Details should be initialized")
			}
			if len(err.Details) != 0 {
				t.Error("NewPostalError() Details should be empty")
			}
		})
	}
}

func TestPostalError_WithDetails(t *testing.T) {
	err := NewPostalError("validation_error", "Invalid request", 400)
	
	details := map[string]interface{}{
		"field": "email",
		"value": "invalid@",
		"reason": "missing domain",
	}
	
	updatedErr := err.WithDetails(details)
	
	// Should return the same instance
	if updatedErr != err {
		t.Error("WithDetails() should return the same instance")
	}
	
	// Should have the details
	if len(err.Details) != 3 {
		t.Errorf("WithDetails() Details length = %v, want 3", len(err.Details))
	}
	
	for key, value := range details {
		if err.Details[key] != value {
			t.Errorf("WithDetails() Details[%s] = %v, want %v", key, err.Details[key], value)
		}
	}
}

func TestIsRateLimit(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "rate limit error",
			err:  ErrRateLimit,
			want: true,
		},
		{
			name: "wrapped rate limit error",
			err:  errors.Join(ErrRateLimit, errors.New("additional context")),
			want: true,
		},
		{
			name: "different error",
			err:  ErrUnauthorized,
			want: false,
		},
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
		{
			name: "custom error",
			err:  errors.New("custom error"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsRateLimit(tt.err); got != tt.want {
				t.Errorf("IsRateLimit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsUnauthorized(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "unauthorized error",
			err:  ErrUnauthorized,
			want: true,
		},
		{
			name: "wrapped unauthorized error",
			err:  errors.Join(ErrUnauthorized, errors.New("invalid key")),
			want: true,
		},
		{
			name: "different error",
			err:  ErrRateLimit,
			want: false,
		},
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
		{
			name: "custom error",
			err:  errors.New("custom error"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsUnauthorized(tt.err); got != tt.want {
				t.Errorf("IsUnauthorized() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsServerError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "server error",
			err:  ErrServerError,
			want: true,
		},
		{
			name: "wrapped server error",
			err:  errors.Join(ErrServerError, errors.New("database down")),
			want: true,
		},
		{
			name: "different error",
			err:  ErrUnauthorized,
			want: false,
		},
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
		{
			name: "custom error",
			err:  errors.New("custom error"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsServerError(tt.err); got != tt.want {
				t.Errorf("IsServerError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorConstants(t *testing.T) {
	// Test that all error constants are properly defined
	if ErrInvalidConfig == nil {
		t.Error("ErrInvalidConfig should not be nil")
	}
	if ErrUnauthorized == nil {
		t.Error("ErrUnauthorized should not be nil")
	}
	if ErrRateLimit == nil {
		t.Error("ErrRateLimit should not be nil")
	}
	if ErrServerError == nil {
		t.Error("ErrServerError should not be nil")
	}
	if ErrInvalidMessage == nil {
		t.Error("ErrInvalidMessage should not be nil")
	}

	// Test error messages
	expectedMessages := map[error]string{
		ErrInvalidConfig:  "invalid configuration",
		ErrUnauthorized:   "unauthorized: invalid API key",
		ErrRateLimit:      "rate limit exceeded",
		ErrServerError:    "postal server error",
		ErrInvalidMessage: "invalid message",
	}

	for err, expectedMsg := range expectedMessages {
		if err.Error() != expectedMsg {
			t.Errorf("Error message for %T = %v, want %v", err, err.Error(), expectedMsg)
		}
	}
}

func TestPostalErrorJSONMarshaling(t *testing.T) {
	err := &PostalError{
		Code:       "validation_error",
		Message:    "Invalid email format",
		StatusCode: 400,
		Details: map[string]interface{}{
			"field": "email",
			"value": "invalid@",
		},
	}

	// Test that PostalError can be properly marshaled/unmarshaled using standard json package
	// This is important for API error responses
	data, marshalErr := json.Marshal(err)
	if marshalErr != nil {
		t.Fatalf("Failed to marshal PostalError: %v", marshalErr)
	}

	var unmarshaled PostalError
	unmarshalErr := json.Unmarshal(data, &unmarshaled)
	if unmarshalErr != nil {
		t.Fatalf("Failed to unmarshal PostalError: %v", unmarshalErr)
	}

	if unmarshaled.Code != err.Code {
		t.Errorf("Unmarshaled Code = %v, want %v", unmarshaled.Code, err.Code)
	}
	if unmarshaled.Message != err.Message {
		t.Errorf("Unmarshaled Message = %v, want %v", unmarshaled.Message, err.Message)
	}
	// Note: StatusCode won't be in JSON due to json:"-" tag
}

func BenchmarkPostalError_Error(b *testing.B) {
	err := &PostalError{
		Code:    "validation_error",
		Message: "Invalid email format",
		Details: map[string]interface{}{
			"field":  "email",
			"value":  "invalid@email",
			"reason": "missing domain part",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

func BenchmarkIsRateLimit(b *testing.B) {
	err := ErrRateLimit
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = IsRateLimit(err)
	}
}

func BenchmarkNewPostalError(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewPostalError("validation_error", "Invalid request", 400)
	}
}