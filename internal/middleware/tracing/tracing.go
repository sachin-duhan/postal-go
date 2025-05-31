package tracing

import (
	"log"
	"net/http"
	"time"

	"github.com/sachin-duhan/postal-go/internal/middleware"
)

// TracingHandler handles trace events
type TracingHandler interface {
	OnRequestStart(req *http.Request)
	OnRequestEnd(req *http.Request, resp *http.Response, duration time.Duration, err error)
}

// Config configures the tracing middleware
type Config struct {
	ServiceName string
	Handler     TracingHandler
}

// DefaultHandler is a basic implementation of TracingHandler
type DefaultHandler struct {
	logger *log.Logger
}

func NewDefaultHandler(logger *log.Logger) *DefaultHandler {
	if logger == nil {
		logger = log.Default()
	}
	return &DefaultHandler{logger: logger}
}

func (h *DefaultHandler) OnRequestStart(req *http.Request) {
	h.logger.Printf("[TRACE] %s request started: %s %s",
		req.Method, req.URL.String(), req.Header.Get("X-Request-ID"))
}

func (h *DefaultHandler) OnRequestEnd(req *http.Request, resp *http.Response, duration time.Duration, err error) {
	status := 0
	if resp != nil {
		status = resp.StatusCode
	}

	if err != nil {
		h.logger.Printf("[TRACE] %s request failed after %v: %s %s [%d] - %v",
			req.Method, duration, req.URL.String(), req.Header.Get("X-Request-ID"), status, err)
		return
	}

	h.logger.Printf("[TRACE] %s request completed in %v: %s %s [%d]",
		req.Method, duration, req.URL.String(), req.Header.Get("X-Request-ID"), status)
}

// New returns a middleware that adds tracing
func New(cfg Config) middleware.Middleware {
	if cfg.Handler == nil {
		cfg.Handler = NewDefaultHandler(nil)
	}

	return func(next http.RoundTripper) http.RoundTripper {
		return &transport{
			next:    next,
			handler: cfg.Handler,
		}
	}
}

type transport struct {
	next    http.RoundTripper
	handler TracingHandler
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()
	t.handler.OnRequestStart(req)

	resp, err := t.next.RoundTrip(req)
	duration := time.Since(start)

	t.handler.OnRequestEnd(req, resp, duration, err)
	return resp, err
}
