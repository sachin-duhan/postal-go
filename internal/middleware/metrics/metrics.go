package metrics

import (
	"net/http"
	"time"

	"github.com/sachin-duhan/postal-go/internal/middleware"
)

// Collector interface for collecting metrics
type Collector interface {
	ObserveRequestDuration(method, path string, duration time.Duration)
	IncRequestCount(method, path string, statusCode int)
	ObserveResponseSize(method, path string, bytes int64)
}

// New returns a middleware that collects metrics
func New(collector Collector) middleware.Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return &transport{
			next:      next,
			collector: collector,
		}
	}
}

type transport struct {
	next      http.RoundTripper
	collector Collector
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()

	resp, err := t.next.RoundTrip(req)

	duration := time.Since(start)
	method := req.Method
	path := req.URL.Path

	if err != nil {
		t.collector.IncRequestCount(method, path, 0) // 0 indicates error
		return resp, err
	}

	t.collector.ObserveRequestDuration(method, path, duration)
	t.collector.IncRequestCount(method, path, resp.StatusCode)

	if resp.ContentLength > 0 {
		t.collector.ObserveResponseSize(method, path, resp.ContentLength)
	}

	return resp, nil
}
