package middleware

import (
	"net/http"
)

// Middleware represents a function that wraps an http.RoundTripper
type Middleware func(http.RoundTripper) http.RoundTripper

// Chain combines multiple middleware into a single middleware
func Chain(middleware ...Middleware) Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		for i := len(middleware) - 1; i >= 0; i-- {
			next = middleware[i](next)
		}
		return next
	}
}

// RoundTripperFunc converts a function to http.RoundTripper
type RoundTripperFunc func(*http.Request) (*http.Response, error)

func (f RoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
