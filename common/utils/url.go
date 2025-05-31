package utils

import (
	"fmt"
	"net/url"
	"strings"
)

// URLBuilder helps construct valid Postal API URLs
type URLBuilder struct {
	baseURL string
}

// NewURLBuilder creates a new URLBuilder
func NewURLBuilder(baseURL string) (*URLBuilder, error) {
	if _, err := ValidateURL(baseURL); err != nil {
		return nil, err
	}
	return &URLBuilder{baseURL: strings.TrimSuffix(baseURL, "/")}, nil
}

// BuildPath joins the base URL with the given path
func (b *URLBuilder) BuildPath(path string) string {
	return fmt.Sprintf("%s/%s", b.baseURL, strings.TrimPrefix(path, "/"))
}

// ValidateURL checks if the URL is valid and returns parsed URL
func ValidateURL(rawURL string) (*url.URL, error) {
	if rawURL == "" {
		return nil, fmt.Errorf("URL cannot be empty")
	}

	// Add scheme if missing
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + rawURL
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL format: %w", err)
	}

	if parsedURL.Host == "" {
		return nil, fmt.Errorf("URL must have a host")
	}

	return parsedURL, nil
}

// StandardizeURL ensures the URL has a scheme and is properly formatted
func StandardizeURL(rawURL string) (string, error) {
	parsedURL, err := ValidateURL(rawURL)
	if err != nil {
		return "", err
	}

	// Use HTTP for localhost/127.0.0.1
	if strings.Contains(parsedURL.Host, "localhost") || strings.Contains(parsedURL.Host, "127.0.0.1") {
		parsedURL.Scheme = "http"
	}

	return parsedURL.String(), nil
}
