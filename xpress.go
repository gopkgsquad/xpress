package xpress

import (
	"io"
	"log"
	"net/http"
)

// ANSI color escape codes
const (
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorReset   = "\033[0m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
)

type ILogger interface {
	Debug(msg interface{})
	Info(msg interface{})
	Warning(msg interface{})
	Error(msg interface{})
	Fatal(msg interface{})
	Infof(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	LogHTTPRequest(level LogLevel, message interface{})
}

type IRouter interface {
	HFunc(string, http.HandlerFunc)
	U(...func(http.Handler) http.Handler)
	M(...func(http.Handler) http.Handler) IRouter

	Group(string) *MuxRouter
}

var logFormat LogFormat

func NewRouter(logger ILogger) *MuxRouter {
	if logFormat == FormatJSON {
		return newMuxRouter(logger, FormatJSON)
	}
	return newMuxRouter(logger, FormatColorized)
}

// NewLogger creates a new instance of the custom logger.
func NewLogger(out io.Writer, level LogLevel, format LogFormat, renderCallerInfo ...bool) ILogger {
	logFormat = format
	var render bool
	if len(renderCallerInfo) > 0 {
		render = renderCallerInfo[0]
	} else {
		render = true
	}

	return &Logger{
		level:            level,
		logger:           log.New(out, "", 0),
		timeColor:        colorCyan,
		fileColor:        colorMagenta,
		renderCallerInfo: &render,
		format:           format,
	}
}
