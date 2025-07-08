package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sachin-duhan/postal-go/common/types"
)

func TestNewTransport(t *testing.T) {
	tests := []struct {
		name    string
		baseURL string
		apiKey  string
		wantErr bool
	}{
		{
			name:    "valid HTTPS URL",
			baseURL: "https://postal.example.com",
			apiKey:  "test-key",
			wantErr: false,
		},
		{
			name:    "valid HTTP URL",
			baseURL: "http://localhost:5000",
			apiKey:  "test-key",
			wantErr: false,
		},
		{
			name:    "URL with path",
			baseURL: "https://postal.example.com/api",
			apiKey:  "test-key",
			wantErr: false,
		},
		{
			name:    "invalid URL",
			baseURL: "://invalid-url",
			apiKey:  "test-key",
			wantErr: true,
		},
		{
			name:    "empty URL",
			baseURL: "",
			apiKey:  "test-key",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &http.Client{}
			transport, err := NewTransport(tt.baseURL, tt.apiKey, client)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTransport() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && transport == nil {
				t.Error("NewTransport() returned nil transport")
			}
			if !tt.wantErr && transport.apiKey != tt.apiKey {
				t.Errorf("NewTransport() apiKey = %v, want %v", transport.apiKey, tt.apiKey)
			}
		})
	}
}

func TestTransportDo(t *testing.T) {
	tests := []struct {
		name           string
		request        *Request
		serverResponse func(w http.ResponseWriter, r *http.Request)
		wantErr        bool
		wantResult     *types.Result
		errContains    string
	}{
		{
			name: "successful request",
			request: &Request{
				Method: http.MethodPost,
				Path:   "send/message",
				Body: map[string]interface{}{
					"to":      []string{"test@example.com"},
					"from":    "sender@example.com",
					"subject": "Test",
					"body":    "Test body",
				},
			},
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				// Verify request method
				if r.Method != http.MethodPost {
					t.Errorf("expected POST, got %s", r.Method)
				}
				// Verify path
				if !strings.Contains(r.URL.Path, "send/message") {
					t.Errorf("expected path to contain 'send/message', got %s", r.URL.Path)
				}
				// Verify headers
				if r.Header.Get("Content-Type") != "application/json" {
					t.Error("expected Content-Type: application/json")
				}
				if r.Header.Get("X-Server-API-Key") != "test-key" {
					t.Error("expected X-Server-API-Key header")
				}

				w.WriteHeader(200)
				json.NewEncoder(w).Encode(types.Result{
					MessageID: "12345",
					Status:    "success",
				})
			},
			wantErr: false,
			wantResult: &types.Result{
				MessageID: "12345",
				Status:    "success",
			},
		},
		{
			name: "request with custom headers",
			request: &Request{
				Method: http.MethodPost,
				Path:   "send/message",
				Body:   map[string]string{"test": "data"},
				Headers: map[string]string{
					"X-Custom-Header": "custom-value",
				},
			},
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("X-Custom-Header") != "custom-value" {
					t.Error("expected custom header")
				}
				w.WriteHeader(200)
				json.NewEncoder(w).Encode(types.Result{Status: "success"})
			},
			wantErr: false,
		},
		{
			name: "400 error response",
			request: &Request{
				Method: http.MethodPost,
				Path:   "send/message",
				Body:   map[string]string{"invalid": "data"},
			},
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(400)
				json.NewEncoder(w).Encode(types.PostalError{
					Code:    "validation_error",
					Message: "Invalid request data",
				})
			},
			wantErr:     true,
			errContains: "validation_error",
		},
		{
			name: "401 unauthorized",
			request: &Request{
				Method: http.MethodPost,
				Path:   "send/message",
				Body:   map[string]string{"test": "data"},
			},
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(401)
				json.NewEncoder(w).Encode(types.PostalError{
					Code:    "unauthorized",
					Message: "Invalid API key",
				})
			},
			wantErr:     true,
			errContains: "unauthorized",
		},
		{
			name: "429 rate limit",
			request: &Request{
				Method: http.MethodPost,
				Path:   "send/message",
				Body:   map[string]string{"test": "data"},
			},
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(429)
				json.NewEncoder(w).Encode(types.PostalError{
					Code:    "rate_limit",
					Message: "Rate limit exceeded",
				})
			},
			wantErr:     true,
			errContains: "rate_limit",
		},
		{
			name: "500 server error",
			request: &Request{
				Method: http.MethodPost,
				Path:   "send/message",
				Body:   map[string]string{"test": "data"},
			},
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(500)
				json.NewEncoder(w).Encode(types.PostalError{
					Code:    "server_error",
					Message: "Internal server error",
				})
			},
			wantErr:     true,
			errContains: "server_error",
		},
		{
			name: "malformed error response",
			request: &Request{
				Method: http.MethodPost,
				Path:   "send/message",
				Body:   map[string]string{"test": "data"},
			},
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(400)
				w.Write([]byte("invalid json"))
			},
			wantErr:     true,
			errContains: "failed to parse error response",
		},
		{
			name: "malformed success response",
			request: &Request{
				Method: http.MethodPost,
				Path:   "send/message",
				Body:   map[string]string{"test": "data"},
			},
			serverResponse: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				w.Write([]byte("invalid json"))
			},
			wantErr:     true,
			errContains: "failed to parse response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			ts := httptest.NewServer(http.HandlerFunc(tt.serverResponse))
			defer ts.Close()

			// Create transport
			client := &http.Client{}
			transport, err := NewTransport(ts.URL, "test-key", client)
			if err != nil {
				t.Fatalf("failed to create transport: %v", err)
			}

			// Execute request
			ctx := context.Background()
			result, err := transport.Do(ctx, tt.request)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("Transport.Do() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && tt.errContains != "" && err != nil {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Transport.Do() error = %v, want error containing %v", err, tt.errContains)
				}
			}

			// Check result
			if !tt.wantErr && result == nil {
				t.Error("Transport.Do() returned nil result")
			}
			if tt.wantResult != nil && result != nil {
				if result.MessageID != tt.wantResult.MessageID {
					t.Errorf("Transport.Do() MessageID = %v, want %v", result.MessageID, tt.wantResult.MessageID)
				}
				if result.Status != tt.wantResult.Status {
					t.Errorf("Transport.Do() Status = %v, want %v", result.Status, tt.wantResult.Status)
				}
			}
		})
	}
}

func TestTransportRequestBody(t *testing.T) {
	// Test that request body is properly marshaled
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}

		var requestData map[string]interface{}
		if err := json.Unmarshal(body, &requestData); err != nil {
			t.Fatalf("failed to unmarshal request body: %v", err)
		}

		// Verify the request data
		if requestData["test"] != "data" {
			t.Errorf("expected test=data, got %v", requestData["test"])
		}

		w.WriteHeader(200)
		json.NewEncoder(w).Encode(types.Result{Status: "success"})
	}))
	defer ts.Close()

	client := &http.Client{}
	transport, err := NewTransport(ts.URL, "test-key", client)
	if err != nil {
		t.Fatalf("failed to create transport: %v", err)
	}

	req := &Request{
		Method: http.MethodPost,
		Path:   "test",
		Body: map[string]interface{}{
			"test": "data",
		},
	}

	ctx := context.Background()
	_, err = transport.Do(ctx, req)
	if err != nil {
		t.Errorf("Transport.Do() error = %v", err)
	}
}

func TestTransportInvalidRequestBody(t *testing.T) {
	client := &http.Client{}
	transport, err := NewTransport("https://example.com", "test-key", client)
	if err != nil {
		t.Fatalf("failed to create transport: %v", err)
	}

	// Create request with body that cannot be marshaled
	req := &Request{
		Method: http.MethodPost,
		Path:   "test",
		Body:   make(chan int), // channels cannot be marshaled to JSON
	}

	ctx := context.Background()
	_, err = transport.Do(ctx, req)
	if err == nil {
		t.Error("expected error for invalid request body")
	}
	if !strings.Contains(err.Error(), "failed to marshal request body") {
		t.Errorf("expected marshal error, got %v", err)
	}
}

func TestTransportContextCancellation(t *testing.T) {
	// Create server that delays response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This should not complete because context will be cancelled
		select {
		case <-r.Context().Done():
			return
		}
	}))
	defer ts.Close()

	client := &http.Client{}
	transport, err := NewTransport(ts.URL, "test-key", client)
	if err != nil {
		t.Fatalf("failed to create transport: %v", err)
	}

	// Create cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	req := &Request{
		Method: http.MethodPost,
		Path:   "test",
		Body:   map[string]string{"test": "data"},
	}

	_, err = transport.Do(ctx, req)
	if err == nil {
		t.Error("expected error due to context cancellation")
	}
}

type mockRoundTripper struct {
	called bool
	rt     http.RoundTripper
}

func (m *mockRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	m.called = true
	if m.rt != nil {
		return m.rt.RoundTrip(r)
	}
	return &http.Response{
		StatusCode: 200,
		Header:     make(http.Header),
		Body:       io.NopCloser(bytes.NewBufferString(`{"status": "success"}`)),
	}, nil
}

func TestTransportAddMiddleware(t *testing.T) {
	client := &http.Client{}
	transport, err := NewTransport("https://example.com", "test-key", client)
	if err != nil {
		t.Fatalf("failed to create transport: %v", err)
	}

	// Create mock middleware
	mock := &mockRoundTripper{}
	middleware := func(next http.RoundTripper) http.RoundTripper {
		mock.rt = next
		return mock
	}

	// Add middleware
	transport.AddMiddleware(middleware)

	// Verify middleware was added
	if len(transport.middleware) != 1 {
		t.Errorf("expected 1 middleware, got %d", len(transport.middleware))
	}
}

func BenchmarkTransportDo(b *testing.B) {
	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(types.Result{
			MessageID: "benchmark-msg",
			Status:    "success",
		})
	}))
	defer ts.Close()

	// Create transport
	client := &http.Client{}
	transport, err := NewTransport(ts.URL, "test-key", client)
	if err != nil {
		b.Fatalf("failed to create transport: %v", err)
	}

	req := &Request{
		Method: http.MethodPost,
		Path:   "send/message",
		Body: map[string]interface{}{
			"to":      []string{"test@example.com"},
			"from":    "sender@example.com",
			"subject": "Benchmark",
			"body":    "Benchmark body",
		},
	}

	ctx := context.Background()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := transport.Do(ctx, req)
		if err != nil {
			b.Fatalf("Transport.Do() error = %v", err)
		}
	}
}

func BenchmarkTransportRequestMarshaling(b *testing.B) {
	req := &Request{
		Method: http.MethodPost,
		Path:   "send/message",
		Body: map[string]interface{}{
			"to":      []string{"test@example.com"},
			"from":    "sender@example.com",
			"subject": "Benchmark",
			"body":    "Benchmark body",
			"attachments": []map[string]string{
				{
					"name":         "test.txt",
					"content_type": "text/plain",
					"data":         "base64data",
				},
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(req.Body)
		if err != nil {
			b.Fatalf("json.Marshal() error = %v", err)
		}
	}
}