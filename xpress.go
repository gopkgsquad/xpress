package xpress

import (
	"net/http"
)

type Router interface {
	HFunc(string, http.HandlerFunc)
	U(...func(http.Handler) http.Handler)
	M(...func(http.Handler) http.Handler) Router

	Group(string) *MuxRouter
	StartServer(*http.Server)
}

func NewRouter() *MuxRouter {
	return newMuxRouter()
}
