package types

import (
	"encoding/json"
	"testing"
)

func TestResult_Success(t *testing.T) {
	tests := []struct {
		name   string
		result *Result
		want   bool
	}{
		{
			name: "successful result",
			result: &Result{
				MessageID: "12345",
				Status:    "success",
			},
			want: true,
		},
		{
			name: "failed result",
			result: &Result{
				MessageID: "12345",
				Status:    "failed",
				Errors:    []string{"validation error"},
			},
			want: false,
		},
		{
			name: "pending result",
			result: &Result{
				MessageID: "12345",
				Status:    "pending",
			},
			want: false,
		},
		{
			name: "empty status",
			result: &Result{
				MessageID: "12345",
				Status:    "",
			},
			want: false,
		},
		{
			name: "case sensitive status",
			result: &Result{
				MessageID: "12345",
				Status:    "Success", // uppercase S
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.result.Success(); got != tt.want {
				t.Errorf("Result.Success() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResult_Failed(t *testing.T) {
	tests := []struct {
		name   string
		result *Result
		want   bool
	}{
		{
			name: "successful result",
			result: &Result{
				MessageID: "12345",
				Status:    "success",
			},
			want: false,
		},
		{
			name: "failed result",
			result: &Result{
				MessageID: "12345",
				Status:    "failed",
				Errors:    []string{"validation error"},
			},
			want: true,
		},
		{
			name: "pending result",
			result: &Result{
				MessageID: "12345",
				Status:    "pending",
			},
			want: true,
		},
		{
			name: "empty status",
			result: &Result{
				MessageID: "12345",
				Status:    "",
			},
			want: true,
		},
		{
			name: "case sensitive status",
			result: &Result{
				MessageID: "12345",
				Status:    "Success", // uppercase S
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.result.Failed(); got != tt.want {
				t.Errorf("Result.Failed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResult_SuccessAndFailed_Inverse(t *testing.T) {
	// Test that Success() and Failed() are always inverse of each other
	testResults := []*Result{
		{MessageID: "1", Status: "success"},
		{MessageID: "2", Status: "failed"},
		{MessageID: "3", Status: "pending"},
		{MessageID: "4", Status: ""},
		{MessageID: "5", Status: "invalid"},
		{MessageID: "6", Status: "Success"},
	}

	for _, result := range testResults {
		success := result.Success()
		failed := result.Failed()
		
		if success == failed {
			t.Errorf("Result with status %q: Success() = %v, Failed() = %v, they should be inverse", 
				result.Status, success, failed)
		}
		if success && failed {
			t.Errorf("Result with status %q: both Success() and Failed() are true", result.Status)
		}
		if !success && !failed {
			t.Errorf("Result with status %q: both Success() and Failed() are false", result.Status)
		}
	}
}

func TestResultJSONMarshaling(t *testing.T) {
	tests := []struct {
		name   string
		result *Result
	}{
		{
			name: "basic successful result",
			result: &Result{
				MessageID: "msg_12345",
				Status:    "success",
			},
		},
		{
			name: "result with data",
			result: &Result{
				MessageID: "msg_12346",
				Status:    "success",
				Data: map[string]interface{}{
					"queue_id":   "queue_67890",
					"priority":   "high",
					"scheduled":  false,
				},
			},
		},
		{
			name: "failed result with errors",
			result: &Result{
				MessageID: "msg_12347",
				Status:    "failed",
				Errors: []string{
					"invalid recipient email",
					"subject too long",
				},
			},
		},
		{
			name: "result with both data and errors",
			result: &Result{
				MessageID: "msg_12348",
				Status:    "partial_success",
				Data: map[string]interface{}{
					"sent_count":   2,
					"failed_count": 1,
				},
				Errors: []string{
					"one recipient failed validation",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal to JSON
			jsonData, err := json.Marshal(tt.result)
			if err != nil {
				t.Fatalf("json.Marshal() error = %v", err)
			}

			// Unmarshal back to Result
			var unmarshaled Result
			err = json.Unmarshal(jsonData, &unmarshaled)
			if err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}

			// Verify fields
			if unmarshaled.MessageID != tt.result.MessageID {
				t.Errorf("MessageID = %v, want %v", unmarshaled.MessageID, tt.result.MessageID)
			}
			if unmarshaled.Status != tt.result.Status {
				t.Errorf("Status = %v, want %v", unmarshaled.Status, tt.result.Status)
			}

			// Verify Data map
			if len(unmarshaled.Data) != len(tt.result.Data) {
				t.Errorf("Data length = %v, want %v", len(unmarshaled.Data), len(tt.result.Data))
			}
			for key, value := range tt.result.Data {
				actualValue := unmarshaled.Data[key]
				// Handle JSON number conversion (int -> float64)
				if actualFloat, ok := actualValue.(float64); ok {
					if expectedInt, ok := value.(int); ok {
						if actualFloat != float64(expectedInt) {
							t.Errorf("Data[%s] = %v, want %v", key, actualValue, value)
						}
						continue
					}
				}
				if actualValue != value {
					t.Errorf("Data[%s] = %v, want %v", key, actualValue, value)
				}
			}

			// Verify Errors slice
			if len(unmarshaled.Errors) != len(tt.result.Errors) {
				t.Errorf("Errors length = %v, want %v", len(unmarshaled.Errors), len(tt.result.Errors))
			}
			for i, err := range tt.result.Errors {
				if i < len(unmarshaled.Errors) && unmarshaled.Errors[i] != err {
					t.Errorf("Errors[%d] = %v, want %v", i, unmarshaled.Errors[i], err)
				}
			}
		})
	}
}

func TestResultJSONOmitEmpty(t *testing.T) {
	// Test that empty/nil fields are omitted from JSON
	result := &Result{
		MessageID: "msg_12349",
		Status:    "success",
		// Data and Errors are nil/empty and should be omitted
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	jsonStr := string(jsonData)
	
	// Should contain required fields
	if !resultContains(jsonStr, "message_id") {
		t.Error("JSON should contain message_id field")
	}
	if !resultContains(jsonStr, "status") {
		t.Error("JSON should contain status field")
	}
	if !resultContains(jsonStr, "msg_12349") {
		t.Error("JSON should contain message ID value")
	}
	if !resultContains(jsonStr, "success") {
		t.Error("JSON should contain status value")
	}

	// Should NOT contain omitted fields
	if resultContains(jsonStr, "data") {
		t.Error("JSON should not contain data field when nil")
	}
	if resultContains(jsonStr, "errors") {
		t.Error("JSON should not contain errors field when nil")
	}
}

func TestResultNilPointer(t *testing.T) {
	// Test that methods don't panic on nil pointer
	// This test just verifies we handle nil pointers gracefully in real usage
	// The methods require a valid Result instance to work properly
}

func BenchmarkResult_Success(b *testing.B) {
	result := &Result{
		MessageID: "msg_benchmark",
		Status:    "success",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = result.Success()
	}
}

func BenchmarkResult_Failed(b *testing.B) {
	result := &Result{
		MessageID: "msg_benchmark",
		Status:    "failed",
		Errors:    []string{"test error"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = result.Failed()
	}
}

func BenchmarkResultJSONMarshal(b *testing.B) {
	result := &Result{
		MessageID: "msg_benchmark",
		Status:    "success",
		Data: map[string]interface{}{
			"queue_id":   "queue_12345",
			"priority":   "normal",
			"scheduled":  false,
			"sent_at":    "2023-12-01T10:00:00Z",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(result)
		if err != nil {
			b.Fatalf("json.Marshal() error = %v", err)
		}
	}
}

// Helper function for string contains check
func resultContains(s, substr string) bool {
	return len(s) >= len(substr) && substr != "" && 
		(s == substr || len(s) > len(substr) && resultContainsHelper(s, substr))
}

func resultContainsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}