package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	client "github.com/sachin-duhan/postal-go"
	"github.com/sachin-duhan/postal-go/common/types"
)

// customRoundTripper implements http.RoundTripper
type customRoundTripper struct {
	next     http.RoundTripper
	callback func(*http.Request) (*http.Response, error)
}

func (rt *customRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return rt.callback(req)
}

func newRoundTripper(next http.RoundTripper, callback func(*http.Request) (*http.Response, error)) http.RoundTripper {
	return &customRoundTripper{next: next, callback: callback}
}

// loggingMiddleware creates a middleware that logs request/response details
func loggingMiddleware() client.Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return newRoundTripper(next, func(req *http.Request) (*http.Response, error) {
			start := time.Now()
			log.Printf("[REQUEST] %s %s", req.Method, req.URL)

			resp, err := next.RoundTrip(req)
			if err != nil {
				log.Printf("[ERROR] Request failed: %v", err)
				return resp, err
			}

			duration := time.Since(start)
			log.Printf("[RESPONSE] Status: %d, Duration: %v", resp.StatusCode, duration)
			return resp, nil
		})
	}
}

// retryMiddleware creates a middleware that implements retry logic
func retryMiddleware(maxRetries int, retryInterval time.Duration) client.Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return newRoundTripper(next, func(req *http.Request) (*http.Response, error) {
			var lastErr error
			for attempt := 0; attempt <= maxRetries; attempt++ {
				if attempt > 0 {
					log.Printf("[RETRY] Attempt %d of %d", attempt, maxRetries)
					time.Sleep(retryInterval)
				}

				// Clone the request body for retries if needed
				if req.Body != nil {
					req.Body = http.NoBody
				}

				resp, err := next.RoundTrip(req)
				if err == nil && resp.StatusCode < 500 {
					return resp, nil
				}

				lastErr = fmt.Errorf("attempt %d failed: %v", attempt, err)
				if resp != nil {
					resp.Body.Close()
				}
			}
			return nil, fmt.Errorf("all retry attempts failed: %v", lastErr)
		})
	}
}

// headerMiddleware adds custom headers to all requests
func headerMiddleware(headers map[string]string) client.Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return newRoundTripper(next, func(req *http.Request) (*http.Response, error) {
			for key, value := range headers {
				req.Header.Set(key, value)
			}
			return next.RoundTrip(req)
		})
	}
}

func main() {
	// Initialize the client with custom configuration
	postalClient, err := client.NewClient(
		"https://postal.example.com", // Replace with your Postal server URL
		"your-api-key",               // Replace with your API key
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Configure the client
	config := &client.Config{
		Timeout:        30 * time.Second,
		MaxRetries:     5,
		RetryInterval:  2 * time.Second,
		MaxConcurrency: 10,
		Debug:          true,
	}

	// Create custom headers
	customHeaders := map[string]string{
		"X-Application-Name": "postal-go-example",
		"X-Client-Version":   "1.0.0",
	}

	// Configure client with middleware and config
	postalClient = postalClient.
		WithConfig(config).
		WithMiddleware(
			headerMiddleware(customHeaders),
			loggingMiddleware(),
			retryMiddleware(3, time.Second),
		)

	// Example 1: Send a complex message with attachments and custom headers
	message := &types.Message{
		To:      []string{"recipient1@example.com", "recipient2@example.com"},
		From:    "sender@yourdomain.com",
		Subject: "Advanced Postal-Go Example",
		Body:    "This is a test email with attachments and custom headers.",
		HTMLBody: `
			<html>
				<body>
					<h1>Advanced Example</h1>
					<p>This email demonstrates:</p>
					<ul>
						<li>Multiple recipients</li>
						<li>HTML content</li>
						<li>Custom headers</li>
						<li>Middleware chain (logging, retries, headers)</li>
					</ul>
				</body>
			</html>
		`,
		Headers: map[string]string{
			"X-Custom-Header":  "custom-value",
			"X-Priority":       "1",
			"List-Unsubscribe": "<mailto:unsubscribe@yourdomain.com>",
		},
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Send the message
	result1, err := postalClient.SendMessage(ctx, message)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
	} else {
		log.Printf("Message sent successfully! Result: %+v", result1)
	}

	// Example 2: Send a raw message
	rawMessage := &types.RawMessage{
		To:   []string{"recipient@example.com"},
		From: "sender@yourdomain.com",
		Mail: `From: sender@yourdomain.com
To: recipient@example.com
Subject: Raw Message Example
Content-Type: multipart/alternative; boundary="boundary123"

--boundary123
Content-Type: text/plain; charset="UTF-8"

This is a raw email message sent using Postal-Go.

--boundary123
Content-Type: text/html; charset="UTF-8"

<html>
<body>
	<h1>Raw Message</h1>
	<p>This is a raw email message sent using <strong>Postal-Go</strong>.</p>
</body>
</html>

--boundary123--
`,
	}

	result2, err := postalClient.SendRawMessage(ctx, rawMessage)
	if err != nil {
		log.Printf("Failed to send raw message: %v", err)
	} else {
		log.Printf("Raw message sent successfully! Result: %+v", result2)
	}
}
