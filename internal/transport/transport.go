package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sachin-duhan/postal-go/common/types"
	"github.com/sachin-duhan/postal-go/common/utils"
	"github.com/sachin-duhan/postal-go/internal/middleware"
)

// Transport handles HTTP communication with the Postal API
type Transport struct {
	urlBuilder *utils.URLBuilder
	apiKey     string
	httpClient *http.Client
	middleware []middleware.Middleware
}

// Request represents an API request
type Request struct {
	Method  string
	Path    string
	Body    interface{}
	Headers map[string]string
}

// NewTransport creates a new Transport instance
func NewTransport(baseURL, apiKey string, client *http.Client) (*Transport, error) {
	// Validate and standardize the URL
	standardURL, err := utils.StandardizeURL(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	urlBuilder, err := utils.NewURLBuilder(standardURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create URL builder: %w", err)
	}

	return &Transport{
		urlBuilder: urlBuilder,
		apiKey:     apiKey,
		httpClient: client,
	}, nil
}

// Do executes an API request
func (t *Transport) Do(ctx context.Context, req *Request) (*types.Result, error) {
	url := t.urlBuilder.BuildPath(req.Path)

	body, err := json.Marshal(req.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, req.Method, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set default headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Server-API-Key", t.apiKey)

	// Set custom headers
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	// Apply middleware chain without modifying the client
	client := t.httpClient
	if len(t.middleware) > 0 {
		// Create a copy of the client to avoid race conditions
		clientCopy := *t.httpClient
		rt := t.httpClient.Transport
		if rt == nil {
			rt = http.DefaultTransport
		}
		clientCopy.Transport = middleware.Chain(t.middleware...)(rt)
		client = &clientCopy
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Handle error responses
	if resp.StatusCode >= 400 {
		var postalErr types.PostalError
		if err := json.Unmarshal(respBody, &postalErr); err != nil {
			return nil, fmt.Errorf("failed to parse error response: %w", err)
		}
		postalErr.StatusCode = resp.StatusCode
		return nil, &postalErr
	}

	// Parse success response
	var result types.Result
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

// AddMiddleware adds middleware to the transport
func (t *Transport) AddMiddleware(m middleware.Middleware) {
	t.middleware = append(t.middleware, m)
}
