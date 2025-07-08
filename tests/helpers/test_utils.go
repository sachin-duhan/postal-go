package helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/sachin-duhan/postal-go/common/types"
)

// TestUtils provides utility functions for testing
type TestUtils struct{}

// NewTestUtils creates a new instance of test utilities
func NewTestUtils() *TestUtils {
	return &TestUtils{}
}

// AssertNoError fails the test if err is not nil
func (u *TestUtils) AssertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// AssertError fails the test if err is nil
func (u *TestUtils) AssertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected error but got nil")
	}
}

// AssertErrorContains fails the test if err is nil or doesn't contain the expected string
func (u *TestUtils) AssertErrorContains(t *testing.T, err error, expected string) {
	t.Helper()
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if !strings.Contains(err.Error(), expected) {
		t.Fatalf("error %q does not contain expected string %q", err.Error(), expected)
	}
}

// AssertEqual fails the test if actual != expected
func (u *TestUtils) AssertEqual(t *testing.T, actual, expected interface{}) {
	t.Helper()
	if actual != expected {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}

// AssertNotEqual fails the test if actual == expected
func (u *TestUtils) AssertNotEqual(t *testing.T, actual, expected interface{}) {
	t.Helper()
	if actual == expected {
		t.Fatalf("expected values to be different, but both are %v", actual)
	}
}

// AssertTrue fails the test if condition is false
func (u *TestUtils) AssertTrue(t *testing.T, condition bool) {
	t.Helper()
	if !condition {
		t.Fatal("expected condition to be true")
	}
}

// AssertFalse fails the test if condition is true
func (u *TestUtils) AssertFalse(t *testing.T, condition bool) {
	t.Helper()
	if condition {
		t.Fatal("expected condition to be false")
	}
}

// AssertNotNil fails the test if value is nil
func (u *TestUtils) AssertNotNil(t *testing.T, value interface{}) {
	t.Helper()
	if value == nil {
		t.Fatal("expected non-nil value")
	}
}

// AssertNil fails the test if value is not nil
func (u *TestUtils) AssertNil(t *testing.T, value interface{}) {
	t.Helper()
	if value != nil {
		t.Fatalf("expected nil value, got %v", value)
	}
}

// AssertSliceEqual fails the test if slices are not equal
func (u *TestUtils) AssertSliceEqual(t *testing.T, actual, expected []string) {
	t.Helper()
	if len(actual) != len(expected) {
		t.Fatalf("slice lengths differ: expected %d, got %d", len(expected), len(actual))
	}
	for i, v := range actual {
		if v != expected[i] {
			t.Fatalf("slice element %d differs: expected %v, got %v", i, expected[i], v)
		}
	}
}

// AssertMapEqual fails the test if maps are not equal
func (u *TestUtils) AssertMapEqual(t *testing.T, actual, expected map[string]string) {
	t.Helper()
	if len(actual) != len(expected) {
		t.Fatalf("map lengths differ: expected %d, got %d", len(expected), len(actual))
	}
	for k, v := range expected {
		if actual[k] != v {
			t.Fatalf("map key %q differs: expected %v, got %v", k, v, actual[k])
		}
	}
}

// AssertJSONEqual compares two objects by marshaling them to JSON
func (u *TestUtils) AssertJSONEqual(t *testing.T, actual, expected interface{}) {
	t.Helper()
	
	actualJSON, err := json.Marshal(actual)
	if err != nil {
		t.Fatalf("failed to marshal actual value: %v", err)
	}
	
	expectedJSON, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("failed to marshal expected value: %v", err)
	}
	
	if !bytes.Equal(actualJSON, expectedJSON) {
		t.Fatalf("JSON values differ:\nActual:   %s\nExpected: %s", actualJSON, expectedJSON)
	}
}

// CreateTestContext creates a context with a reasonable timeout for tests
func (u *TestUtils) CreateTestContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 30*time.Second)
}

// CreateTestContextWithTimeout creates a context with a custom timeout
func (u *TestUtils) CreateTestContextWithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

// HTTPResponseBuilder helps build HTTP responses for testing
type HTTPResponseBuilder struct {
	statusCode int
	headers    map[string]string
	body       interface{}
}

// NewHTTPResponseBuilder creates a new HTTP response builder
func NewHTTPResponseBuilder() *HTTPResponseBuilder {
	return &HTTPResponseBuilder{
		statusCode: 200,
		headers:    make(map[string]string),
	}
}

// WithStatusCode sets the HTTP status code
func (b *HTTPResponseBuilder) WithStatusCode(code int) *HTTPResponseBuilder {
	b.statusCode = code
	return b
}

// WithHeader adds a header to the response
func (b *HTTPResponseBuilder) WithHeader(key, value string) *HTTPResponseBuilder {
	b.headers[key] = value
	return b
}

// WithJSONBody sets the response body as JSON
func (b *HTTPResponseBuilder) WithJSONBody(body interface{}) *HTTPResponseBuilder {
	b.body = body
	return b
}

// Build creates the HTTP response
func (b *HTTPResponseBuilder) Build() *http.Response {
	var bodyReader io.ReadCloser = http.NoBody
	
	if b.body != nil {
		jsonData, err := json.Marshal(b.body)
		if err == nil {
			bodyReader = io.NopCloser(bytes.NewReader(jsonData))
		}
	}
	
	response := &http.Response{
		StatusCode: b.statusCode,
		Status:     http.StatusText(b.statusCode),
		Header:     make(http.Header),
		Body:       bodyReader,
	}
	
	// Set default content type for JSON
	if b.body != nil {
		response.Header.Set("Content-Type", "application/json")
	}
	
	// Add custom headers
	for key, value := range b.headers {
		response.Header.Set(key, value)
	}
	
	return response
}

// EmailValidator provides email validation utilities for testing
type EmailValidator struct{}

// NewEmailValidator creates a new email validator
func NewEmailValidator() *EmailValidator {
	return &EmailValidator{}
}

// IsValid checks if an email address is valid (simple validation)
func (v *EmailValidator) IsValid(email string) bool {
	if email == "" {
		return false
	}
	
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	
	if len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}
	
	return strings.Contains(parts[1], ".")
}

// GenerateTestEmails generates a list of test email addresses
func (v *EmailValidator) GenerateTestEmails(count int) []string {
	emails := make([]string, count)
	for i := 0; i < count; i++ {
		emails[i] = generateTestEmail(i)
	}
	return emails
}

func generateTestEmail(index int) string {
	domains := []string{"example.com", "test.org", "sample.net"}
	domain := domains[index%len(domains)]
	return strings.ToLower(strings.ReplaceAll(strings.ReplaceAll("user INDEX@DOMAIN", "INDEX", string(rune('0'+index%10))), "DOMAIN", domain))
}

// MessageValidator provides message validation utilities for testing
type MessageValidator struct{}

// NewMessageValidator creates a new message validator
func NewMessageValidator() *MessageValidator {
	return &MessageValidator{}
}

// ValidateBasicFields checks if a message has all required basic fields
func (v *MessageValidator) ValidateBasicFields(msg *types.Message) []string {
	var errors []string
	
	if len(msg.To) == 0 {
		errors = append(errors, "missing To field")
	}
	
	if msg.From == "" {
		errors = append(errors, "missing From field")
	}
	
	if msg.Subject == "" {
		errors = append(errors, "missing Subject field")
	}
	
	if msg.Body == "" && msg.HTMLBody == "" {
		errors = append(errors, "missing body content")
	}
	
	return errors
}

// ValidateEmailAddresses checks if all email addresses in a message are valid
func (v *MessageValidator) ValidateEmailAddresses(msg *types.Message) []string {
	var errors []string
	validator := NewEmailValidator()
	
	// Validate To addresses
	for _, email := range msg.To {
		if !validator.IsValid(email) {
			errors = append(errors, "invalid To email: "+email)
		}
	}
	
	// Validate CC addresses
	for _, email := range msg.CC {
		if !validator.IsValid(email) {
			errors = append(errors, "invalid CC email: "+email)
		}
	}
	
	// Validate BCC addresses
	for _, email := range msg.BCC {
		if !validator.IsValid(email) {
			errors = append(errors, "invalid BCC email: "+email)
		}
	}
	
	// Validate From address
	if msg.From != "" && !validator.IsValid(msg.From) {
		errors = append(errors, "invalid From email: "+msg.From)
	}
	
	return errors
}

// PerformanceTimer helps measure execution times in tests
type PerformanceTimer struct {
	startTime time.Time
}

// NewPerformanceTimer creates a new performance timer
func NewPerformanceTimer() *PerformanceTimer {
	return &PerformanceTimer{
		startTime: time.Now(),
	}
}

// Elapsed returns the time elapsed since the timer was created
func (pt *PerformanceTimer) Elapsed() time.Duration {
	return time.Since(pt.startTime)
}

// Reset resets the timer
func (pt *PerformanceTimer) Reset() {
	pt.startTime = time.Now()
}

// AssertMaxDuration fails the test if elapsed time exceeds maxDuration
func (pt *PerformanceTimer) AssertMaxDuration(t *testing.T, maxDuration time.Duration) {
	t.Helper()
	elapsed := pt.Elapsed()
	if elapsed > maxDuration {
		t.Fatalf("operation took too long: %v > %v", elapsed, maxDuration)
	}
}

// ConcurrencyTester helps test concurrent operations
type ConcurrencyTester struct {
	numWorkers int
	operations []func() error
}

// NewConcurrencyTester creates a new concurrency tester
func NewConcurrencyTester(numWorkers int) *ConcurrencyTester {
	return &ConcurrencyTester{
		numWorkers: numWorkers,
	}
}

// AddOperation adds an operation to be executed concurrently
func (ct *ConcurrencyTester) AddOperation(op func() error) {
	ct.operations = append(ct.operations, op)
}

// Execute runs all operations concurrently and returns any errors
func (ct *ConcurrencyTester) Execute() []error {
	if len(ct.operations) == 0 {
		return nil
	}
	
	errChan := make(chan error, len(ct.operations))
	
	// Run operations concurrently
	for _, op := range ct.operations {
		go func(operation func() error) {
			errChan <- operation()
		}(op)
	}
	
	// Collect results
	var errors []error
	for i := 0; i < len(ct.operations); i++ {
		if err := <-errChan; err != nil {
			errors = append(errors, err)
		}
	}
	
	return errors
}

// StringUtils provides string manipulation utilities for testing
type StringUtils struct{}

// NewStringUtils creates a new string utilities instance
func NewStringUtils() *StringUtils {
	return &StringUtils{}
}

// Contains checks if a string contains a substring
func (su *StringUtils) Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// ContainsAny checks if a string contains any of the given substrings
func (su *StringUtils) ContainsAny(s string, substrs ...string) bool {
	for _, substr := range substrs {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}

// ContainsAll checks if a string contains all of the given substrings
func (su *StringUtils) ContainsAll(s string, substrs ...string) bool {
	for _, substr := range substrs {
		if !strings.Contains(s, substr) {
			return false
		}
	}
	return true
}

// RandomString generates a random string of specified length
func (su *StringUtils) RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[time.Now().Nanosecond()%len(charset)]
	}
	return string(result)
}