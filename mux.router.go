package xpress

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var _ IRouter = &MuxRouter{}

// MuxRouter implements the Router interface
type MuxRouter struct {
	mux         *http.ServeMux
	middlewares []func(http.Handler) http.Handler
	prefix      string
	logger      ILogger
	LogFormat   LogFormat
}

// newMuxRouter creates a new instance of MyRouter
func newMuxRouter(logger ILogger, format LogFormat) *MuxRouter {
	return &MuxRouter{
		mux:       http.NewServeMux(),
		prefix:    "",
		logger:    logger,
		LogFormat: format,
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
func (r *MuxRouter) M(middlewares ...func(http.Handler) http.Handler) IRouter {
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
	methodColor := colorCyan
	pathColor := colorGreen
	timeColor := colorMagenta
	remoteAddrColor := colorYellow

	host, _, _ := net.SplitHostPort(req.RemoteAddr)
	reqPath := req.URL.Path
	if req.URL.RawQuery != "" {
		reqPath += "?" + req.URL.RawQuery
	}

	if r.LogFormat == FormatColorized {
		r.logger.LogHTTPRequest(LogLevelInfo, fmt.Sprintf("%s | %s | %s | %s | %s",
			colorize(statusColor, strconv.Itoa(capturer.status)),
			colorize(methodColor, req.Method),
			colorize(pathColor, reqPath),
			colorize(timeColor, formatResponseTime(duration)),
			colorize(remoteAddrColor, host),
		))
	} else if r.LogFormat == FormatJSON {
		logStruct := map[string]interface{}{
			"status":        strconv.Itoa(capturer.status),
			"method":        req.Method,
			"path":          reqPath,
			"response_time": formatResponseTime(duration),
			"remote_addr":   host,
		}
		r.logger.LogHTTPRequest(LogLevelInfo, logStruct)
	}
}

// getColorForStatus returns the color for the given HTTP status code
func getColorForStatus(status int) string {
	switch {
	case status >= 200 && status < 300:
		return colorGreen // Success status code
	case status >= 400 && status < 500:
		return colorYellow // Client error status code
	case status >= 500:
		return colorRed // Server error status code
	default:
		return colorReset // Other status codes
	}
}

// colorize adds color to a string
func colorize(color, text string) string {
	return fmt.Sprintf("%s%s%s", color, text, colorReset)
}
