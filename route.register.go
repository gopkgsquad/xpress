package xpress

import "net/http"

type Module struct {
	Name       string
	Routes     []Route
	Middleware []func(http.Handler) http.Handler
}

type Route struct {
	Method     string
	Path       string
	Handler    http.HandlerFunc
	Middleware []func(http.Handler) http.Handler
}

func RegisterRoutes(modules []Module, router *MuxRouter) {
	for _, module := range modules {
		for _, route := range module.Routes {
			nrouter := router.Group(module.Name).M(module.Middleware...)
			registerRoute(route.Method+" "+route.Path, route, nrouter)
		}
	}
}

func registerRoute(pattern string, route Route, router IRouter) {
	router.M(route.Middleware...).HFunc(pattern, route.Handler)
}
