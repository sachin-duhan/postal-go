package client

import (
	"net/http"
	"time"
)

// Config holds the client configuration
type Config struct {
	Timeout        time.Duration
	MaxRetries     int
	RetryInterval  time.Duration
	MaxConcurrency int
	Debug          bool
	Transport      *http.Transport
}

// Option is a function that configures the client
type Option func(*clientImpl)

// Middleware represents a function that wraps the client's transport layer
type Middleware func(http.RoundTripper) http.RoundTripper

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		Timeout:        30 * time.Second,
		MaxRetries:     3,
		RetryInterval:  time.Second,
		MaxConcurrency: 10,
		Debug:          false,
		Transport:      http.DefaultTransport.(*http.Transport).Clone(),
	}
}
