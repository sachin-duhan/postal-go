package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	client "github.com/sachin-duhan/postal-go"
	"github.com/sachin-duhan/postal-go/common/types"
)

// TestMain sets up and tears down the test environment
func TestMain(m *testing.M) {
	// Run tests
	code := m.Run()
	os.Exit(code)
}

func TestClientIntegration_SendMessage(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Create mock postal server
	server := NewMockPostalServer()
	defer server.Close()

	// Create client
	postalClient, err := client.NewClient(server.URL, "test-api-key")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()

	tests := []struct {
		name        string
		message     *types.Message
		wantErr     bool
		wantStatus  string
	}{
		{
			name: "successful message send",
			message: &types.Message{
				To:       []string{"recipient@example.com"},
				From:     "sender@example.com",
				Subject:  "Integration Test",
				HTMLBody: "<h1>Test Message</h1>",
			},
			wantErr:    false,
			wantStatus: "success",
		},
		{
			name: "message with attachments",
			message: &types.Message{
				To:      []string{"recipient@example.com"},
				From:    "sender@example.com",
				Subject: "Test with Attachment",
				Body:    "Message with attachment",
				Attachments: []types.Attachment{
					{
						Name:        "test.txt",
						ContentType: "text/plain",
						Data:        "VGVzdCBjb250ZW50", // Base64 for "Test content"
					},
				},
			},
			wantErr:    false,
			wantStatus: "success",
		},
		{
			name: "message with multiple recipients",
			message: &types.Message{
				To:       []string{"recipient1@example.com", "recipient2@example.com"},
				CC:       []string{"cc@example.com"},
				BCC:      []string{"bcc@example.com"},
				From:     "sender@example.com",
				Subject:  "Multiple Recipients",
				HTMLBody: "Test message for multiple recipients",
			},
			wantErr:    false,
			wantStatus: "success",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := postalClient.SendMessage(ctx, tt.message)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("SendMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !tt.wantErr {
				if result == nil {
					t.Error("SendMessage() returned nil result")
					return
				}
				if result.Status != tt.wantStatus {
					t.Errorf("SendMessage() status = %v, want %v", result.Status, tt.wantStatus)
				}
				if result.MessageID == "" {
					t.Error("SendMessage() returned empty message ID")
				}
			}
		})
	}
}

func TestClientIntegration_SendRawMessage(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Create mock postal server
	server := NewMockPostalServer()
	defer server.Close()

	// Create client
	postalClient, err := client.NewClient(server.URL, "test-api-key")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()

	rawMsg := &types.RawMessage{
		Mail: "From: sender@example.com\r\nTo: recipient@example.com\r\nSubject: Raw Message Test\r\n\r\nThis is a raw message body.",
		To:   []string{"recipient@example.com"},
		From: "sender@example.com",
	}

	result, err := postalClient.SendRawMessage(ctx, rawMsg)
	if err != nil {
		t.Fatalf("SendRawMessage() error = %v", err)
	}

	if result.Status != "success" {
		t.Errorf("SendRawMessage() status = %v, want success", result.Status)
	}
	if result.MessageID == "" {
		t.Error("SendRawMessage() returned empty message ID")
	}
}

func TestClientIntegration_ConcurrentSending(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Create mock postal server
	server := NewMockPostalServer()
	defer server.Close()

	// Create client
	postalClient, err := client.NewClient(server.URL, "test-api-key")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()
	numMessages := 50
	var wg sync.WaitGroup
	errors := make(chan error, numMessages)
	results := make(chan *types.Result, numMessages)

	// Send messages concurrently
	for i := 0; i < numMessages; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			
			msg := &types.Message{
				To:       []string{fmt.Sprintf("recipient%d@example.com", i)},
				From:     "sender@example.com",
				Subject:  fmt.Sprintf("Concurrent Test %d", i),
				HTMLBody: fmt.Sprintf("<h1>Message %d</h1>", i),
			}
			
			result, err := postalClient.SendMessage(ctx, msg)
			if err != nil {
				errors <- err
			} else {
				results <- result
			}
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(errors)
	close(results)

	// Check for errors
	var errorCount int
	for err := range errors {
		t.Errorf("concurrent send error: %v", err)
		errorCount++
	}

	// Check results
	var successCount int
	for result := range results {
		if result.Success() {
			successCount++
		}
	}

	if errorCount > 0 {
		t.Errorf("had %d errors out of %d messages", errorCount, numMessages)
	}
	if successCount != numMessages-errorCount {
		t.Errorf("expected %d successful results, got %d", numMessages-errorCount, successCount)
	}
}

func TestClientIntegration_ErrorHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Create mock postal server with error responses
	server := NewMockPostalServerWithErrors()
	defer server.Close()

	// Create client
	postalClient, err := client.NewClient(server.URL, "test-api-key")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()

	tests := []struct {
		name        string
		message     *types.Message
		wantErr     bool
		errContains string
	}{
		{
			name: "rate limit error",
			message: &types.Message{
				To:      []string{"ratelimit@example.com"},
				From:    "sender@example.com",
				Subject: "Rate Limit Test",
				Body:    "This should trigger rate limit",
			},
			wantErr:     true,
			errContains: "rate_limit",
		},
		{
			name: "unauthorized error",
			message: &types.Message{
				To:      []string{"unauthorized@example.com"},
				From:    "sender@example.com",
				Subject: "Unauthorized Test",
				Body:    "This should trigger unauthorized",
			},
			wantErr:     true,
			errContains: "unauthorized",
		},
		{
			name: "server error",
			message: &types.Message{
				To:      []string{"servererror@example.com"},
				From:    "sender@example.com",
				Subject: "Server Error Test",
				Body:    "This should trigger server error",
			},
			wantErr:     true,
			errContains: "server_error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := postalClient.SendMessage(ctx, tt.message)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("SendMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if tt.wantErr && tt.errContains != "" {
				if !contains(err.Error(), tt.errContains) {
					t.Errorf("SendMessage() error = %v, want error containing %v", err, tt.errContains)
				}
			}
		})
	}
}

func TestClientIntegration_ContextCancellation(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Create mock postal server with delay
	server := NewMockPostalServerWithDelay(200 * time.Millisecond)
	defer server.Close()

	// Create client
	postalClient, err := client.NewClient(server.URL, "test-api-key")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Create context with short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	msg := &types.Message{
		To:      []string{"recipient@example.com"},
		From:    "sender@example.com",
		Subject: "Timeout Test",
		Body:    "This should timeout",
	}

	_, err = postalClient.SendMessage(ctx, msg)
	if err == nil {
		t.Error("expected timeout error")
	}
}

// MockPostalServer creates a test server that mimics Postal API responses
type MockPostalServer struct {
	*httptest.Server
	messageCounter int
	mu            sync.Mutex
}

func NewMockPostalServer() *MockPostalServer {
	mps := &MockPostalServer{}
	mux := http.NewServeMux()
	
	mux.HandleFunc("/api/v1/send/message", mps.handleSendMessage)
	mux.HandleFunc("/api/v1/send/raw", mps.handleSendRaw)
	
	mps.Server = httptest.NewServer(mux)
	return mps
}

func (mps *MockPostalServer) handleSendMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Verify API key
	if r.Header.Get("X-Server-API-Key") != "test-api-key" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	mps.mu.Lock()
	mps.messageCounter++
	msgID := fmt.Sprintf("msg_%d", mps.messageCounter)
	counter := mps.messageCounter
	mps.mu.Unlock()

	result := types.Result{
		MessageID: msgID,
		Status:    "success",
		Data: map[string]interface{}{
			"queue_id": fmt.Sprintf("queue_%d", counter),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (mps *MockPostalServer) handleSendRaw(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Verify API key
	if r.Header.Get("X-Server-API-Key") != "test-api-key" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	mps.mu.Lock()
	mps.messageCounter++
	msgID := fmt.Sprintf("raw_msg_%d", mps.messageCounter)
	mps.mu.Unlock()

	result := types.Result{
		MessageID: msgID,
		Status:    "success",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func NewMockPostalServerWithErrors() *httptest.Server {
	mux := http.NewServeMux()
	
	mux.HandleFunc("/api/v1/send/message", func(w http.ResponseWriter, r *http.Request) {
		var msg types.Message
		json.NewDecoder(r.Body).Decode(&msg)
		
		// Trigger different errors based on recipient
		if len(msg.To) > 0 {
			recipient := msg.To[0]
			switch {
			case contains(recipient, "ratelimit"):
				w.WriteHeader(429)
				json.NewEncoder(w).Encode(types.PostalError{
					Code:    "rate_limit",
					Message: "Rate limit exceeded",
				})
			case contains(recipient, "unauthorized"):
				w.WriteHeader(401)
				json.NewEncoder(w).Encode(types.PostalError{
					Code:    "unauthorized",
					Message: "Invalid API key",
				})
			case contains(recipient, "servererror"):
				w.WriteHeader(500)
				json.NewEncoder(w).Encode(types.PostalError{
					Code:    "server_error",
					Message: "Internal server error",
				})
			default:
				w.WriteHeader(200)
				json.NewEncoder(w).Encode(types.Result{
					MessageID: "test_msg",
					Status:    "success",
				})
			}
		}
	})
	
	return httptest.NewServer(mux)
}

func NewMockPostalServerWithDelay(delay time.Duration) *httptest.Server {
	mux := http.NewServeMux()
	
	mux.HandleFunc("/api/v1/send/message", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(types.Result{
			MessageID: "delayed_msg",
			Status:    "success",
		})
	})
	
	return httptest.NewServer(mux)
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && substr != "" && 
		(s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}