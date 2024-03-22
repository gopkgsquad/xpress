package xpress

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gopkgsquad/gloader"
	"github.com/gopkgsquad/glogger"
)

// ANSI color escape codes
const (
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorReset   = "\033[0m"
	ColorMagenta = "\033[35m"
	ColorCyan    = "\033[36m"
)

var _ Router = &MuxRouter{}

// MuxRouter implements the Router interface
type MuxRouter struct {
	mux         *http.ServeMux
	middlewares []func(http.Handler) http.Handler
	prefix      string
	logger      *glogger.Logger
}

// newMuxRouter creates a new instance of MyRouter
func newMuxRouter() *MuxRouter {
	logger := glogger.NewLogger(os.Stdout, glogger.LogLevelInfo)
	return &MuxRouter{
		mux:    http.NewServeMux(),
		prefix: "",
		logger: logger,
	}
}

// HFunc registers an HTTP handler for a specific method and path
func (r *MuxRouter) HFunc(pattern string, handler http.HandlerFunc) {
	if pattern == "" || len(strings.SplitN(pattern, " ", 2)) < 2 || handler == nil {
		r.logger.Fatal("pattern or handler is nil or invalid")
	}

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
		logger:      r.logger,
	}
}

// Group creates a new router with a prefix
func (r *MuxRouter) Group(prefix string) *MuxRouter {
	return &MuxRouter{
		mux:         r.mux,
		middlewares: r.middlewares,
		prefix:      r.prefix + prefix,
		logger:      r.logger,
	}
}

func (r *MuxRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	capturer := &statusCapturer{ResponseWriter: w}

	r.mux.ServeHTTP(capturer, req)

	duration := time.Since(start)

	statusColor := getColorForStatus(capturer.status)
	methodColor := ColorCyan
	pathColor := ColorGreen
	timeColor := ColorMagenta
	remoteAddrColor := ColorYellow

	// Log request details with colors
	r.logger.Infof("%s | %s | %s | %s | %s",
		colorize(statusColor, strconv.Itoa(capturer.status)),
		colorize(methodColor, req.Method),
		colorize(pathColor, req.URL.Path),
		colorize(timeColor, formatResponseTime(duration)),
		colorize(remoteAddrColor, strings.SplitN(req.RemoteAddr, ":", 2)[0]),
	)
}

func (r *MuxRouter) StartServer(srv *http.Server) {
	gloader.NewWatcher(srv, time.Second, r.logger).Start()
}

// getColorForStatus returns the color for the given HTTP status code
func getColorForStatus(status int) string {
	switch {
	case status >= 200 && status < 300:
		return ColorGreen // Success status code
	case status >= 400 && status < 500:
		return ColorYellow // Client error status code
	case status >= 500:
		return ColorRed // Server error status code
	default:
		return ColorReset // Other status codes
	}
}

// colorize adds color to a string
func colorize(color, text string) string {
	return fmt.Sprintf("%s%s%s", color, text, ColorReset)
}
