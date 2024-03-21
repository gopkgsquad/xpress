package xpress

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

var _ Router = &MuxRouter{}

// MuxRouter implements the Router interface
type MuxRouter struct {
	mux         *http.ServeMux
	middlewares []func(http.Handler) http.Handler
	prefix      string
}

// newMuxRouter creates a new instance of MyRouter
func newMuxRouter() *MuxRouter {
	return &MuxRouter{
		mux:    http.NewServeMux(),
		prefix: "",
	}
}

// HFunc registers an HTTP handler for a specific method and path
func (r *MuxRouter) HFunc(pattern string, handler http.HandlerFunc) {
	npattern := strings.Split(pattern, " ")[0] + " " + r.prefix + strings.Split(pattern, " ")[1]
	r.mux.HandleFunc(npattern, func(w http.ResponseWriter, req *http.Request) {
		r.chain(handler).ServeHTTP(w, req)
	})
}

// U adds middleware to the router
func (r *MuxRouter) U(middlewares ...func(http.Handler) http.Handler) {
	r.middlewares = append(r.middlewares, middlewares...)
}

// M creates a new router with middleware
func (r *MuxRouter) M(middlewares ...func(http.Handler) http.Handler) Router {
	return &MuxRouter{
		mux:         r.mux,
		middlewares: append(r.middlewares, middlewares...),
		prefix:      r.prefix,
	}
}

// Group creates a new router with a prefix
func (r *MuxRouter) Group(prefix string) *MuxRouter {
	return &MuxRouter{
		mux:         r.mux,
		middlewares: r.middlewares,
		prefix:      r.prefix + prefix,
	}
}

func (r *MuxRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	capturer := &statusCapturer{ResponseWriter: w}

	r.mux.ServeHTTP(capturer, req)

	duration := time.Since(start)

	msg := fmt.Sprintf("%v | %v | %s | %s | %s", capturer.status, formatResponseTime(duration), strings.SplitN(req.RemoteAddr, ":", 2)[0], req.Method, req.URL.Path)

	NewLogger().Info(msg)
}
