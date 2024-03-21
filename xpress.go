package xpress

import "net/http"

type Router interface {
	HFunc(pattern string, handler http.HandlerFunc)
	U(middlewares ...func(http.Handler) http.Handler)
	M(middlewares ...func(http.Handler) http.Handler) Router

	Group(prefix string) *MuxRouter
}

func NewRouter() *MuxRouter {
	return newMuxRouter()
}
