package client

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sachin-duhan/postal-go/common/types"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		baseURL string
		apiKey  string
		opts    []Option
		wantErr bool
	}{
		{
			name:    "valid client creation",
			baseURL: "https://postal.example.com",
			apiKey:  "test-api-key",
			wantErr: false,
		},
		{
			name:    "valid client with http URL",
			baseURL: "http://localhost:5000",
			apiKey:  "test-api-key",
			wantErr: false,
		},
		{
			name:    "invalid URL",
			baseURL: "://invalid-url",
			apiKey:  "test-api-key",
			wantErr: true,
		},
		{
			name:    "empty URL",
			baseURL: "",
			apiKey:  "test-api-key",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.baseURL, tt.apiKey, tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && client == nil {
				t.Error("NewClient() returned nil client")
			}
		})
	}
}

func TestSendMessage(t *testing.T) {
	tests := []struct {
		name           string
		message        *types.Message
		mockResponse   string
		mockStatusCode int
		wantErr        bool
		errContains    string
	}{
		{
			name: "valid message",
			message: &types.Message{
				To:       []string{"recipient@example.com"},
				From:     "sender@example.com",
				Subject:  "Test Subject",
				HTMLBody: "<p>Test Body</p>",
			},
			mockResponse:   `{"message_id": "12345", "status": "success"}`,
			mockStatusCode: 200,
			wantErr:        false,
		},
		{
			name: "message with attachments",
			message: &types.Message{
				To:      []string{"recipient@example.com"},
				From:    "sender@example.com",
				Subject: "Test with Attachment",
				Body:    "Test Body",
				Attachments: []types.Attachment{
					{
						Name:        "test.txt",
						ContentType: "text/plain",
						Data:        "VGVzdCBjb250ZW50", // Base64 for "Test content"
					},
				},
			},
			mockResponse:   `{"message_id": "12346", "status": "success"}`,
			mockStatusCode: 200,
			wantErr:        false,
		},
		{
			name: "validation error - missing to",
			message: &types.Message{
				From:    "sender@example.com",
				Subject: "Test Subject",
				Body:    "Test Body",
			},
			wantErr:     true,
			errContains: "recipient (To) is required",
		},
		{
			name: "validation error - missing from",
			message: &types.Message{
				To:      []string{"recipient@example.com"},
				Subject: "Test Subject",
				Body:    "Test Body",
			},
			wantErr:     true,
			errContains: "sender (From) is required",
		},
		{
			name: "validation error - missing subject",
			message: &types.Message{
				To:   []string{"recipient@example.com"},
				From: "sender@example.com",
				Body: "Test Body",
			},
			wantErr:     true,
			errContains: "subject is required",
		},
		{
			name: "validation error - missing body",
			message: &types.Message{
				To:      []string{"recipient@example.com"},
				From:    "sender@example.com",
				Subject: "Test Subject",
			},
			wantErr:     true,
			errContains: "either plain body or HTML body is required",
		},
		{
			name: "validation error - invalid email",
			message: &types.Message{
				To:      []string{"invalid-email"},
				From:    "sender@example.com",
				Subject: "Test Subject",
				Body:    "Test Body",
			},
			wantErr:     true,
			errContains: "invalid recipient email",
		},
		{
			name: "server error response",
			message: &types.Message{
				To:       []string{"recipient@example.com"},
				From:     "sender@example.com",
				Subject:  "Test Subject",
				HTMLBody: "Test Body",
			},
			mockResponse:   `{"code": "rate_limit", "message": "Rate limit exceeded"}`,
			mockStatusCode: 429,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request
				if r.Method != http.MethodPost {
					t.Errorf("expected POST method, got %s", r.Method)
				}
				if r.URL.Path != "/api/v1/send/message" {
					t.Errorf("expected path /api/v1/send/message, got %s", r.URL.Path)
				}
				if r.Header.Get("X-Server-API-Key") != "test-key" {
					t.Errorf("expected API key header")
				}

				// Send response
				w.WriteHeader(tt.mockStatusCode)
				if tt.mockResponse != "" {
					w.Write([]byte(tt.mockResponse))
				}
			}))
			defer ts.Close()

			// Create client
			client, err := NewClient(ts.URL, "test-key")
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}

			// Send message
			ctx := context.Background()
			result, err := client.SendMessage(ctx, tt.message)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("SendMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && tt.errContains != "" && err != nil {
				if !contains(err.Error(), tt.errContains) {
					t.Errorf("SendMessage() error = %v, want error containing %v", err, tt.errContains)
				}
			}
			if !tt.wantErr && result == nil {
				t.Error("SendMessage() returned nil result")
			}
		})
	}
}

func TestSendRawMessage(t *testing.T) {
	tests := []struct {
		name           string
		message        *types.RawMessage
		mockResponse   string
		mockStatusCode int
		wantErr        bool
		errContains    string
	}{
		{
			name: "valid raw message",
			message: &types.RawMessage{
				Mail: "From: sender@example.com\r\nTo: recipient@example.com\r\nSubject: Test\r\n\r\nBody",
				To:   []string{"recipient@example.com"},
				From: "sender@example.com",
			},
			mockResponse:   `{"message_id": "12347", "status": "success"}`,
			mockStatusCode: 200,
			wantErr:        false,
		},
		{
			name: "validation error - missing mail content",
			message: &types.RawMessage{
				To:   []string{"recipient@example.com"},
				From: "sender@example.com",
			},
			wantErr:     true,
			errContains: "raw mail content is required",
		},
		{
			name: "validation error - missing to",
			message: &types.RawMessage{
				Mail: "From: sender@example.com\r\nSubject: Test\r\n\r\nBody",
				From: "sender@example.com",
			},
			wantErr:     true,
			errContains: "recipient (To) is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request
				if r.URL.Path != "/api/v1/send/raw" {
					t.Errorf("expected path /api/v1/send/raw, got %s", r.URL.Path)
				}

				// Send response
				w.WriteHeader(tt.mockStatusCode)
				if tt.mockResponse != "" {
					w.Write([]byte(tt.mockResponse))
				}
			}))
			defer ts.Close()

			// Create client
			client, err := NewClient(ts.URL, "test-key")
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}

			// Send raw message
			ctx := context.Background()
			result, err := client.SendRawMessage(ctx, tt.message)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("SendRawMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && tt.errContains != "" && err != nil {
				if !contains(err.Error(), tt.errContains) {
					t.Errorf("SendRawMessage() error = %v, want error containing %v", err, tt.errContains)
				}
			}
			if !tt.wantErr && result == nil {
				t.Error("SendRawMessage() returned nil result")
			}
		})
	}
}

func TestClientWithConfig(t *testing.T) {
	client, err := NewClient("https://postal.example.com", "test-key")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test WithConfig
	newConfig := &Config{
		Timeout:        60 * time.Second,
		MaxRetries:     5,
		RetryInterval:  2 * time.Second,
		MaxConcurrency: 20,
		Debug:          true,
	}

	updatedClient := client.WithConfig(newConfig)
	if updatedClient == nil {
		t.Error("WithConfig() returned nil")
	}

	// Verify it returns the same client (method chaining)
	if updatedClient != client {
		t.Error("WithConfig() should return the same client instance")
	}
}

func TestClientWithMiddleware(t *testing.T) {
	client, err := NewClient("https://postal.example.com", "test-key")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	testMiddleware := func(next http.RoundTripper) http.RoundTripper {
		return roundTripperFunc(func(r *http.Request) (*http.Response, error) {
			return next.RoundTrip(r)
		})
	}

	updatedClient := client.WithMiddleware(testMiddleware)
	if updatedClient == nil {
		t.Error("WithMiddleware() returned nil")
	}

	// Verify it returns the same client (method chaining)
	if updatedClient != client {
		t.Error("WithMiddleware() should return the same client instance")
	}
}

func TestConcurrentSending(t *testing.T) {
	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate processing time
		time.Sleep(10 * time.Millisecond)
		w.WriteHeader(200)
		w.Write([]byte(`{"message_id": "12348", "status": "success"}`))
	}))
	defer ts.Close()

	// Create client
	client, err := NewClient(ts.URL, "test-key")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Send messages concurrently
	numMessages := 10
	errors := make(chan error, numMessages)
	ctx := context.Background()

	for i := 0; i < numMessages; i++ {
		go func(i int) {
			msg := &types.Message{
				To:       []string{"recipient@example.com"},
				From:     "sender@example.com",
				Subject:  "Test Subject",
				HTMLBody: "Test Body",
			}
			_, err := client.SendMessage(ctx, msg)
			errors <- err
		}(i)
	}

	// Collect results
	for i := 0; i < numMessages; i++ {
		if err := <-errors; err != nil {
			t.Errorf("concurrent send failed: %v", err)
		}
	}
}

func TestContextCancellation(t *testing.T) {
	// Create test server with delay
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate long processing
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(200)
		w.Write([]byte(`{"message_id": "12349", "status": "success"}`))
	}))
	defer ts.Close()

	// Create client
	client, err := NewClient(ts.URL, "test-key")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Create context with immediate cancellation
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Try to send message
	msg := &types.Message{
		To:       []string{"recipient@example.com"},
		From:     "sender@example.com",
		Subject:  "Test Subject",
		HTMLBody: "Test Body",
	}

	_, err = client.SendMessage(ctx, msg)
	if err == nil {
		t.Error("expected error due to context cancellation")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled error, got %v", err)
	}
}

// Helper functions
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(s)] != "" && substr != "" &&
		(s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || contains(s[1:], substr)))
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}
