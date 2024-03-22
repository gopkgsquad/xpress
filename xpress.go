package xpress

import (
	"net/http"
	"os"
	"time"

	"github.com/gopkgsquad/gloader"
	"github.com/gopkgsquad/glogger"
)

type Router interface {
	HFunc(string, http.HandlerFunc)
	U(...func(http.Handler) http.Handler)
	M(...func(http.Handler) http.Handler) Router

	Group(string) *MuxRouter
	StartServer(*http.Server)
}

func NewWatcher(srv *http.Server, interval time.Duration) *gloader.Watcher {
	logger := glogger.NewLogger(os.Stdout, glogger.LogLevelInfo, false)
	return gloader.NewWatcher(srv, interval, logger)
}

func NewRouter() *MuxRouter {
	return newMuxRouter()
}
