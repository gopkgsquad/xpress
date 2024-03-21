package xpress

import (
	"fmt"
	"net/http"
	"time"
)

type statusCapturer struct {
	http.ResponseWriter
	status int
}

func (s *statusCapturer) WriteHeader(statusCode int) {
	s.status = statusCode
	s.ResponseWriter.WriteHeader(statusCode)
}

func (r *MuxRouter) chain(endpoint http.HandlerFunc) http.Handler {
	// Return ahead of time if there aren't any middlewares for the chain
	if len(r.middlewares) == 0 {
		return endpoint
	}

	// Wrap the end handler with the middleware chain
	h := r.middlewares[len(r.middlewares)-1](endpoint)
	for i := len(r.middlewares) - 2; i >= 0; i-- {
		h = r.middlewares[i](h)
	}

	return h
}

func formatResponseTime(duration time.Duration) string {
	if duration < time.Millisecond {
		return fmt.Sprintf("%.2fµs", float64(duration.Microseconds()))
	}
	return fmt.Sprintf("%.2fms", float64(duration.Milliseconds()))
}
