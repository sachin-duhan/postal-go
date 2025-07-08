package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"time"

	"github.com/sachin-duhan/postal-go/common/types"
)

// MockPostalServer provides a mock implementation of the Postal API server
type MockPostalServer struct {
	*httptest.Server
	messageCounter int
	mu            sync.Mutex
	config        MockServerConfig
}

// MockServerConfig configures the mock server behavior
type MockServerConfig struct {
	// Delay adds artificial latency to responses
	Delay time.Duration
	
	// FailureRate sets the percentage of requests that should fail (0.0 to 1.0)
	FailureRate float64
	
	// ValidAPIKeys defines which API keys are considered valid
	ValidAPIKeys []string
	
	// CustomResponses allows overriding responses for specific patterns
	CustomResponses map[string]MockResponse
}

// MockResponse defines a custom response for the mock server
type MockResponse struct {
	StatusCode int
	Body       interface{}
	Headers    map[string]string
}

// NewMockPostalServer creates a new mock Postal server with default configuration
func NewMockPostalServer() *MockPostalServer {
	return NewMockPostalServerWithConfig(MockServerConfig{
		ValidAPIKeys: []string{"test-api-key", "valid-key"},
	})
}

// NewMockPostalServerWithConfig creates a new mock server with custom configuration
func NewMockPostalServerWithConfig(config MockServerConfig) *MockPostalServer {
	mps := &MockPostalServer{
		config: config,
	}
	
	if len(mps.config.ValidAPIKeys) == 0 {
		mps.config.ValidAPIKeys = []string{"test-api-key"}
	}
	
	if mps.config.CustomResponses == nil {
		mps.config.CustomResponses = make(map[string]MockResponse)
	}
	
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/send/message", mps.handleSendMessage)
	mux.HandleFunc("/api/v1/send/raw", mps.handleSendRaw)
	mux.HandleFunc("/health", mps.handleHealth)
	
	mps.Server = httptest.NewServer(mux)
	return mps
}

// NewMockPostalServerWithDelay creates a mock server that adds delay to responses
func NewMockPostalServerWithDelay(delay time.Duration) *MockPostalServer {
	return NewMockPostalServerWithConfig(MockServerConfig{
		Delay:        delay,
		ValidAPIKeys: []string{"test-api-key"},
	})
}

// NewMockPostalServerWithErrors creates a mock server that returns various error responses
func NewMockPostalServerWithErrors() *MockPostalServer {
	customResponses := map[string]MockResponse{
		"ratelimit": {
			StatusCode: 429,
			Body: types.PostalError{
				Code:    "rate_limit",
				Message: "Rate limit exceeded",
			},
		},
		"unauthorized": {
			StatusCode: 401,
			Body: types.PostalError{
				Code:    "unauthorized",
				Message: "Invalid API key",
			},
		},
		"servererror": {
			StatusCode: 500,
			Body: types.PostalError{
				Code:    "server_error",
				Message: "Internal server error",
			},
		},
		"validation": {
			StatusCode: 400,
			Body: types.PostalError{
				Code:    "validation_error",
				Message: "Invalid request data",
			},
		},
	}
	
	return NewMockPostalServerWithConfig(MockServerConfig{
		ValidAPIKeys:    []string{"test-api-key"},
		CustomResponses: customResponses,
	})
}

func (mps *MockPostalServer) handleSendMessage(w http.ResponseWriter, r *http.Request) {
	mps.addDelay()
	
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Check API key
	if !mps.isValidAPIKey(r.Header.Get("X-Server-API-Key")) {
		mps.writeErrorResponse(w, 401, types.PostalError{
			Code:    "unauthorized",
			Message: "Invalid API key",
		})
		return
	}
	
	// Parse request body
	var msg types.Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		mps.writeErrorResponse(w, 400, types.PostalError{
			Code:    "invalid_json",
			Message: "Invalid JSON in request body",
		})
		return
	}
	
	// Check for custom responses based on recipient
	if len(msg.To) > 0 {
		for pattern, response := range mps.config.CustomResponses {
			if mps.containsPattern(msg.To[0], pattern) {
				mps.writeCustomResponse(w, response)
				return
			}
		}
	}
	
	// Simulate random failures if configured
	if mps.shouldFail() {
		mps.writeErrorResponse(w, 500, types.PostalError{
			Code:    "random_failure",
			Message: "Simulated random failure",
		})
		return
	}
	
	// Generate successful response
	mps.mu.Lock()
	mps.messageCounter++
	msgID := fmt.Sprintf("msg_%d_%d", mps.messageCounter, time.Now().Unix())
	mps.mu.Unlock()
	
	result := types.Result{
		MessageID: msgID,
		Status:    "success",
		Data: map[string]interface{}{
			"queue_id":    fmt.Sprintf("queue_%d", mps.messageCounter),
			"priority":    "normal",
			"scheduled":   false,
			"recipients":  len(msg.To),
			"created_at":  time.Now().Format(time.RFC3339),
		},
	}
	
	mps.writeJSONResponse(w, 200, result)
}

func (mps *MockPostalServer) handleSendRaw(w http.ResponseWriter, r *http.Request) {
	mps.addDelay()
	
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Check API key
	if !mps.isValidAPIKey(r.Header.Get("X-Server-API-Key")) {
		mps.writeErrorResponse(w, 401, types.PostalError{
			Code:    "unauthorized",
			Message: "Invalid API key",
		})
		return
	}
	
	// Parse request body
	var rawMsg types.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&rawMsg); err != nil {
		mps.writeErrorResponse(w, 400, types.PostalError{
			Code:    "invalid_json",
			Message: "Invalid JSON in request body",
		})
		return
	}
	
	// Check for custom responses
	if len(rawMsg.To) > 0 {
		for pattern, response := range mps.config.CustomResponses {
			if mps.containsPattern(rawMsg.To[0], pattern) {
				mps.writeCustomResponse(w, response)
				return
			}
		}
	}
	
	// Simulate random failures
	if mps.shouldFail() {
		mps.writeErrorResponse(w, 500, types.PostalError{
			Code:    "random_failure",
			Message: "Simulated random failure",
		})
		return
	}
	
	// Generate successful response
	mps.mu.Lock()
	mps.messageCounter++
	msgID := fmt.Sprintf("raw_msg_%d_%d", mps.messageCounter, time.Now().Unix())
	mps.mu.Unlock()
	
	result := types.Result{
		MessageID: msgID,
		Status:    "success",
		Data: map[string]interface{}{
			"queue_id":   fmt.Sprintf("raw_queue_%d", mps.messageCounter),
			"type":       "raw",
			"created_at": time.Now().Format(time.RFC3339),
		},
	}
	
	mps.writeJSONResponse(w, 200, result)
}

func (mps *MockPostalServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   "mock-1.0.0",
	}
	mps.writeJSONResponse(w, 200, response)
}

func (mps *MockPostalServer) addDelay() {
	if mps.config.Delay > 0 {
		time.Sleep(mps.config.Delay)
	}
}

func (mps *MockPostalServer) isValidAPIKey(apiKey string) bool {
	for _, validKey := range mps.config.ValidAPIKeys {
		if apiKey == validKey {
			return true
		}
	}
	return false
}

func (mps *MockPostalServer) shouldFail() bool {
	if mps.config.FailureRate <= 0 {
		return false
	}
	// Simple pseudo-random failure based on time
	return float64(time.Now().Nanosecond()%1000)/1000.0 < mps.config.FailureRate
}

func (mps *MockPostalServer) containsPattern(text, pattern string) bool {
	return strings.Contains(strings.ToLower(text), strings.ToLower(pattern))
}

func (mps *MockPostalServer) writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (mps *MockPostalServer) writeErrorResponse(w http.ResponseWriter, statusCode int, err types.PostalError) {
	err.StatusCode = statusCode
	mps.writeJSONResponse(w, statusCode, err)
}

func (mps *MockPostalServer) writeCustomResponse(w http.ResponseWriter, response MockResponse) {
	// Set custom headers
	for key, value := range response.Headers {
		w.Header().Set(key, value)
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response.Body)
}

// GetMessageCount returns the number of messages processed by the mock server
func (mps *MockPostalServer) GetMessageCount() int {
	mps.mu.Lock()
	defer mps.mu.Unlock()
	return mps.messageCounter
}

// ResetMessageCount resets the message counter
func (mps *MockPostalServer) ResetMessageCount() {
	mps.mu.Lock()
	defer mps.mu.Unlock()
	mps.messageCounter = 0
}

// AddCustomResponse adds a custom response for a specific pattern
func (mps *MockPostalServer) AddCustomResponse(pattern string, response MockResponse) {
	mps.config.CustomResponses[pattern] = response
}

// RemoveCustomResponse removes a custom response pattern
func (mps *MockPostalServer) RemoveCustomResponse(pattern string) {
	delete(mps.config.CustomResponses, pattern)
}

// SetFailureRate sets the failure rate for the mock server
func (mps *MockPostalServer) SetFailureRate(rate float64) {
	if rate < 0 {
		rate = 0
	}
	if rate > 1 {
		rate = 1
	}
	mps.config.FailureRate = rate
}

// MockTransport provides a mock HTTP transport for testing without a server
type MockTransport struct {
	responses map[string]*http.Response
	mu       sync.RWMutex
}

// NewMockTransport creates a new mock transport
func NewMockTransport() *MockTransport {
	return &MockTransport{
		responses: make(map[string]*http.Response),
	}
}

// RoundTrip implements the http.RoundTripper interface
func (mt *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	mt.mu.RLock()
	defer mt.mu.RUnlock()
	
	key := fmt.Sprintf("%s %s", req.Method, req.URL.Path)
	if response, exists := mt.responses[key]; exists {
		return response, nil
	}
	
	// Default response if no specific response is configured
	return &http.Response{
		StatusCode: 404,
		Status:     "404 Not Found",
		Header:     make(http.Header),
		Body:       http.NoBody,
	}, nil
}

// SetResponse sets a mock response for a specific method and path
func (mt *MockTransport) SetResponse(method, path string, response *http.Response) {
	mt.mu.Lock()
	defer mt.mu.Unlock()
	key := fmt.Sprintf("%s %s", method, path)
	mt.responses[key] = response
}

// ClearResponses clears all configured responses
func (mt *MockTransport) ClearResponses() {
	mt.mu.Lock()
	defer mt.mu.Unlock()
	mt.responses = make(map[string]*http.Response)
}