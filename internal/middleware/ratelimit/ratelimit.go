package ratelimit

import (
	"net/http"

	"golang.org/x/time/rate"

	"github.com/sachin-duhan/postal-go/internal/middleware"
)

// Config configures the rate limit middleware
type Config struct {
	RequestsPerSecond float64
	Burst             int
	Enabled           bool
}

// New returns a middleware that limits request rate
func New(cfg Config) middleware.Middleware {
	if !cfg.Enabled {
		return func(next http.RoundTripper) http.RoundTripper {
			return next
		}
	}

	return func(next http.RoundTripper) http.RoundTripper {
		return &transport{
			next:    next,
			limiter: rate.NewLimiter(rate.Limit(cfg.RequestsPerSecond), cfg.Burst),
		}
	}
}

type transport struct {
	next    http.RoundTripper
	limiter *rate.Limiter
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	err := t.limiter.Wait(ctx)
	if err != nil {
		return nil, err
	}
	return t.next.RoundTrip(req)
}
