package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/sachin-duhan/postal-go/common/types"
)

// BenchmarkClientCreation benchmarks client creation
func BenchmarkClientCreation(b *testing.B) {
	baseURL := "https://postal.example.com"
	apiKey := "test-api-key"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client, err := NewClient(baseURL, apiKey)
		if err != nil {
			b.Fatalf("failed to create client: %v", err)
		}
		_ = client
	}
}

// BenchmarkSendMessage benchmarks sending a single message
func BenchmarkSendMessage(b *testing.B) {
	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(types.Result{
			MessageID: "benchmark-msg",
			Status:    "success",
		})
	}))
	defer ts.Close()

	// Create client
	client, err := NewClient(ts.URL, "test-key")
	if err != nil {
		b.Fatalf("failed to create client: %v", err)
	}

	message := &types.Message{
		To:       []string{"recipient@example.com"},
		From:     "sender@example.com",
		Subject:  "Benchmark Test",
		HTMLBody: "<h1>Benchmark</h1><p>This is a benchmark test message.</p>",
	}

	ctx := context.Background()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := client.SendMessage(ctx, message)
		if err != nil {
			b.Fatalf("SendMessage() error = %v", err)
		}
	}
}

// BenchmarkSendMessageConcurrent benchmarks sending messages concurrently
func BenchmarkSendMessageConcurrent(b *testing.B) {
	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(types.Result{
			MessageID: "concurrent-msg",
			Status:    "success",
		})
	}))
	defer ts.Close()

	// Create client
	client, err := NewClient(ts.URL, "test-key")
	if err != nil {
		b.Fatalf("failed to create client: %v", err)
	}

	message := &types.Message{
		To:       []string{"recipient@example.com"},
		From:     "sender@example.com",
		Subject:  "Concurrent Benchmark",
		HTMLBody: "<h1>Concurrent Test</h1>",
	}

	ctx := context.Background()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := client.SendMessage(ctx, message)
			if err != nil {
				b.Fatalf("SendMessage() error = %v", err)
			}
		}
	})
}

// BenchmarkSendRawMessage benchmarks sending raw messages
func BenchmarkSendRawMessage(b *testing.B) {
	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(types.Result{
			MessageID: "raw-benchmark-msg",
			Status:    "success",
		})
	}))
	defer ts.Close()

	// Create client
	client, err := NewClient(ts.URL, "test-key")
	if err != nil {
		b.Fatalf("failed to create client: %v", err)
	}

	rawMessage := &types.RawMessage{
		Mail: "From: sender@example.com\r\nTo: recipient@example.com\r\nSubject: Raw Benchmark\r\n\r\nRaw message body",
		To:   []string{"recipient@example.com"},
		From: "sender@example.com",
	}

	ctx := context.Background()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := client.SendRawMessage(ctx, rawMessage)
		if err != nil {
			b.Fatalf("SendRawMessage() error = %v", err)
		}
	}
}

// BenchmarkMessageWithAttachments benchmarks messages with various attachment sizes
func BenchmarkMessageWithAttachments(b *testing.B) {
	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(types.Result{
			MessageID: "attachment-msg",
			Status:    "success",
		})
	}))
	defer ts.Close()

	// Create client
	client, err := NewClient(ts.URL, "test-key")
	if err != nil {
		b.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()

	// Test different attachment sizes
	attachmentSizes := []int{1024, 10240, 102400} // 1KB, 10KB, 100KB

	for _, size := range attachmentSizes {
		b.Run(byteCountBinary(int64(size)), func(b *testing.B) {
			// Create attachment data
			data := make([]byte, size)
			for i := range data {
				data[i] = byte('A')
			}

			message := &types.Message{
				To:      []string{"recipient@example.com"},
				From:    "sender@example.com",
				Subject: "Attachment Benchmark",
				Body:    "Message with attachment",
				Attachments: []types.Attachment{
					{
						Name:        "test-file.txt",
						ContentType: "text/plain",
						Data:        string(data), // Simplified for benchmark
					},
				},
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := client.SendMessage(ctx, message)
				if err != nil {
					b.Fatalf("SendMessage() error = %v", err)
				}
			}
		})
	}
}

// BenchmarkMultipleRecipients benchmarks messages with varying recipient counts
func BenchmarkMultipleRecipients(b *testing.B) {
	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(types.Result{
			MessageID: "multi-recipient-msg",
			Status:    "success",
		})
	}))
	defer ts.Close()

	// Create client
	client, err := NewClient(ts.URL, "test-key")
	if err != nil {
		b.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()

	// Test different recipient counts
	recipientCounts := []int{1, 10, 50, 100}

	for _, count := range recipientCounts {
		b.Run(getRecipientCountName(count), func(b *testing.B) {
			// Generate recipients
			recipients := make([]string, count)
			for i := 0; i < count; i++ {
				recipients[i] = generateRecipient(i)
			}

			message := &types.Message{
				To:       recipients,
				From:     "sender@example.com",
				Subject:  "Multiple Recipients Benchmark",
				HTMLBody: "<h1>Broadcast Message</h1>",
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := client.SendMessage(ctx, message)
				if err != nil {
					b.Fatalf("SendMessage() error = %v", err)
				}
			}
		})
	}
}

// BenchmarkConcurrentClients benchmarks multiple clients sending messages concurrently
func BenchmarkConcurrentClients(b *testing.B) {
	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(types.Result{
			MessageID: "concurrent-client-msg",
			Status:    "success",
		})
	}))
	defer ts.Close()

	numClients := 10
	clients := make([]Client, numClients)

	// Create multiple clients
	for i := 0; i < numClients; i++ {
		client, err := NewClient(ts.URL, "test-key")
		if err != nil {
			b.Fatalf("failed to create client %d: %v", i, err)
		}
		clients[i] = client
	}

	message := &types.Message{
		To:       []string{"recipient@example.com"},
		From:     "sender@example.com",
		Subject:  "Concurrent Clients Test",
		HTMLBody: "<h1>Concurrent Test</h1>",
	}

	ctx := context.Background()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		clientIndex := 0
		for pb.Next() {
			client := clients[clientIndex%numClients]
			clientIndex++
			
			_, err := client.SendMessage(ctx, message)
			if err != nil {
				b.Fatalf("SendMessage() error = %v", err)
			}
		}
	})
}

// BenchmarkWithMiddleware benchmarks client with middleware overhead
func BenchmarkWithMiddleware(b *testing.B) {
	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(types.Result{
			MessageID: "middleware-msg",
			Status:    "success",
		})
	}))
	defer ts.Close()

	// Create client
	client, err := NewClient(ts.URL, "test-key")
	if err != nil {
		b.Fatalf("failed to create client: %v", err)
	}

	// Add multiple middleware layers
	for i := 0; i < 5; i++ {
		client = client.WithMiddleware(func(next http.RoundTripper) http.RoundTripper {
			return roundTripperFunc(func(r *http.Request) (*http.Response, error) {
				// Simulate middleware overhead
				time.Sleep(time.Microsecond)
				return next.RoundTrip(r)
			})
		})
	}

	message := &types.Message{
		To:       []string{"recipient@example.com"},
		From:     "sender@example.com",
		Subject:  "Middleware Benchmark",
		HTMLBody: "<h1>Middleware Test</h1>",
	}

	ctx := context.Background()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := client.SendMessage(ctx, message)
		if err != nil {
			b.Fatalf("SendMessage() error = %v", err)
		}
	}
}

// BenchmarkClientConfigUpdate benchmarks config updates
func BenchmarkClientConfigUpdate(b *testing.B) {
	client, err := NewClient("https://postal.example.com", "test-key")
	if err != nil {
		b.Fatalf("failed to create client: %v", err)
	}

	config := &Config{
		Timeout:        60 * time.Second,
		MaxRetries:     5,
		RetryInterval:  2 * time.Second,
		MaxConcurrency: 20,
		Debug:          true,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = client.WithConfig(config)
	}
}

// BenchmarkMemoryAllocation benchmarks memory allocation patterns
func BenchmarkMemoryAllocation(b *testing.B) {
	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(types.Result{
			MessageID: "memory-test-msg",
			Status:    "success",
		})
	}))
	defer ts.Close()

	// Create client
	client, err := NewClient(ts.URL, "test-key")
	if err != nil {
		b.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Create new message each iteration to test allocation
		message := &types.Message{
			To:       []string{"recipient@example.com"},
			From:     "sender@example.com",
			Subject:  "Memory Test",
			HTMLBody: "<h1>Memory allocation test</h1>",
			Headers: map[string]string{
				"X-Test": "value",
			},
		}

		_, err := client.SendMessage(ctx, message)
		if err != nil {
			b.Fatalf("SendMessage() error = %v", err)
		}
	}
}

// BenchmarkHighThroughput benchmarks high throughput scenarios
func BenchmarkHighThroughput(b *testing.B) {
	// Create test server with minimal processing
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"message_id":"throughput-msg","status":"success"}`))
	}))
	defer ts.Close()

	// Create client
	client, err := NewClient(ts.URL, "test-key")
	if err != nil {
		b.Fatalf("failed to create client: %v", err)
	}

	message := &types.Message{
		To:      []string{"recipient@example.com"},
		From:    "sender@example.com",
		Subject: "Throughput Test",
		Body:    "High throughput test message",
	}

	ctx := context.Background()
	
	// Test sequential throughput
	b.Run("Sequential", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := client.SendMessage(ctx, message)
			if err != nil {
				b.Fatalf("SendMessage() error = %v", err)
			}
		}
	})

	// Test parallel throughput
	b.Run("Parallel", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, err := client.SendMessage(ctx, message)
				if err != nil {
					b.Fatalf("SendMessage() error = %v", err)
				}
			}
		})
	})
}

// Helper functions for benchmarks

func byteCountBinary(b int64) string {
	const unit = 1024
	if b < unit {
		return formatInt(b) + "B"
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return formatFloat(float64(b)/float64(div)) + "KMGTPE"[exp:exp+1] + "B"
}

func formatInt(i int64) string {
	if i < 10 {
		return string(rune('0' + i))
	}
	return "10+"
}

func formatFloat(f float64) string {
	if f < 10.0 {
		return "1-9"
	}
	return "10+"
}

func getRecipientCountName(count int) string {
	return string(rune('0'+count/10)) + "recipients"
}

func generateRecipient(index int) string {
	return "recipient" + string(rune('0'+(index%10))) + "@example.com"
}

// Stress test with resource monitoring
func BenchmarkStressTest(b *testing.B) {
	// Create test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate some processing time
		time.Sleep(time.Millisecond)
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(types.Result{
			MessageID: "stress-test-msg",
			Status:    "success",
		})
	}))
	defer ts.Close()

	// Create multiple clients
	numClients := 100
	clients := make([]Client, numClients)
	for i := 0; i < numClients; i++ {
		client, err := NewClient(ts.URL, "test-key")
		if err != nil {
			b.Fatalf("failed to create client %d: %v", i, err)
		}
		clients[i] = client
	}

	message := &types.Message{
		To:       []string{"recipient@example.com"},
		From:     "sender@example.com",
		Subject:  "Stress Test",
		HTMLBody: "<h1>Stress Test</h1>",
	}

	ctx := context.Background()
	var wg sync.WaitGroup
	var mu sync.Mutex
	var successCount, errorCount int

	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(clientIndex int) {
			defer wg.Done()
			
			client := clients[clientIndex%numClients]
			_, err := client.SendMessage(ctx, message)
			
			mu.Lock()
			if err != nil {
				errorCount++
			} else {
				successCount++
			}
			mu.Unlock()
		}(i)
	}
	
	wg.Wait()
	
	b.Logf("Success: %d, Errors: %d", successCount, errorCount)
}