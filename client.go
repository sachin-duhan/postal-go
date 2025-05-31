package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sachin-duhan/postal-go/common/types"
	"github.com/sachin-duhan/postal-go/common/validation"
	"github.com/sachin-duhan/postal-go/internal/transport"
)

// Client represents the interface for interacting with the Postal API
type Client interface {
	// SendMessage sends an email using the message builder pattern
	SendMessage(ctx context.Context, msg *types.Message) (*types.Result, error)

	// SendRawMessage sends a pre-formatted email message
	SendRawMessage(ctx context.Context, raw *types.RawMessage) (*types.Result, error)

	// WithMiddleware adds middleware to the client
	WithMiddleware(middleware ...Middleware) Client

	// WithConfig updates the client configuration
	WithConfig(cfg *Config) Client
}

// clientImpl is the concrete implementation of the Client interface
type clientImpl struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	config     *Config
	middleware []Middleware
	transport  *transport.Transport
}

// NewClient creates a new Postal API client
func NewClient(baseURL, apiKey string, opts ...Option) (Client, error) {
	client := &clientImpl{
		baseURL:    baseURL,
		apiKey:     apiKey,
		httpClient: &http.Client{},
		config:     DefaultConfig(),
	}

	// Initialize transport
	transport, err := transport.NewTransport(baseURL, apiKey, client.httpClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create transport: %w", err)
	}
	client.transport = transport

	// Apply options
	for _, opt := range opts {
		opt(client)
	}

	return client, nil
}

// SendMessage implements Client
func (c *clientImpl) SendMessage(ctx context.Context, msg *types.Message) (*types.Result, error) {
	if err := validation.ValidateMessage(msg); err != nil {
		return nil, err
	}

	req := &transport.Request{
		Method: http.MethodPost,
		Path:   "send/message",
		Body:   msg,
	}

	return c.transport.Do(ctx, req)
}

// SendRawMessage implements Client
func (c *clientImpl) SendRawMessage(ctx context.Context, raw *types.RawMessage) (*types.Result, error) {
	if err := validation.ValidateRawMessage(raw); err != nil {
		return nil, err
	}

	req := &transport.Request{
		Method: http.MethodPost,
		Path:   "send/raw",
		Body:   raw,
	}

	return c.transport.Do(ctx, req)
}

// WithMiddleware implements Client
func (c *clientImpl) WithMiddleware(middleware ...Middleware) Client {
	c.middleware = append(c.middleware, middleware...)
	return c
}

// WithConfig implements Client
func (c *clientImpl) WithConfig(cfg *Config) Client {
	c.config = cfg
	c.httpClient.Timeout = cfg.Timeout
	return c
}

// Ensure clientImpl implements Client interface
var _ Client = (*clientImpl)(nil)
