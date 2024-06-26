# xpress

<p align="center">
  xpress is a package built on top of Go standard HTTP package. Designed for easy use.
</p>

## **Attention**

xpress requires **Go version `1.22` or higher** to run. You can check your current Go version by running `go version` in your terminal

## ⚙️ Installation

```bash
go get -u github.com/gopkgsquad/xpress
```

## Quickstart
```go
package main

import "github.com/gopkgsquad/xpress"

func main() {
    // Initialize a new xpress router
    router := xpress.NewRouter()

    // Define a route for the GET method on the root path '/'
    router.HFunc("GET /", func(w http.ResponseWriter, r *http.Request) {})

    srv := &http.Server{
	Addr:    ":3000",
	Handler: router,
    }

    // Start the server on port 3000
    if err := srv.ListenAndServe(); err != nil {
	log.Fatalf("error:", err.Error())
    }
}
```

## Examples

#### [**Basic Routing**]
```go
func main() {
    // Initialize a new xpress router
    router := xpress.NewRouter()

    // Define a route for the POST method on the root path '/'
    router.HFunc("POST /", func(w http.ResponseWriter, r *http.Request) {})
    // Define a route for the GET method on the root path '/'
    router.HFunc("GET /", func(w http.ResponseWriter, r *http.Request) {})
    // Define a route for the GET method on the root path '/' with {id} as a params
    router.HFunc("GET /{id}", func(w http.ResponseWriter, r *http.Request) {})
    // Define a route for the PUT method on the root path '/' with {id} as a params
    router.HFunc("PUT /{id}", func(w http.ResponseWriter, r *http.Request) {})
    // Define a route for the DELETE method on the root path '/' with {id} as a params
    router.HFunc("DELETE /{id}", func(w http.ResponseWriter, r *http.Request) {})

    srv := &http.Server{
	Addr:    ":3000",
	Handler: router,
    }

    // Start the server on port 3000
    if err := srv.ListenAndServe(); err != nil {
	log.Fatalf("error:", err.Error())
    }
}

```

#### [**Live Reloading**]
```go
package main

import "github.com/gopkgsquad/xpress"

func main() {
    // Initialize a new xpress router
    router := xpress.NewRouter()

    // Define a route for the GET method on the root path '/'
    router.HFunc("GET /", func(w http.ResponseWriter, r *http.Request) {})

    srv := &http.Server{
	Addr:    ":3000",
	Handler: router,
    }

    // Start the watcher for live reloading, first parameter is http.Server
    // and second parameter is frequency to check file changes
    // here in NewWatcher you also pass 3rd param true/false to show log is
    // coming from which file and line eg xpress.NewWatcher(srv, time.Second*2, false).Start()
    // default value is true

    // you can also call router.StartServer(srv) to start the application with live reaload - available in v1.0.8 or higher
    xpress.NewWatcher(srv, time.Second*2).Start()
}
```

#### [**Middleware Setup**]
```go
func LoggerMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("FROM Global Middleware")
	h.Serve(w, r)
	})
}

func Authenticate(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // your authentication logic goes here
		log.Println("FROM Authenticate Middleware")
        h.Serve(w, r)
	})
}

func main() {
    // Initialize a new xpress router
    router := xpress.NewRouter()

    // you can pass n number of middleware eg router.U(LoggerMiddleware1, LoggerMiddleware2, ...)
    router.U(LoggerMiddleware)

    // here in router.M as well you can pass n number of middleware eg router.M(Authenticate, ValidateRequest, ...)
    router.M(Authenticate).HFunc("POST /", func(w http.ResponseWriter, r *http.Request) {})


   srv := &http.Server{
	Addr:    ":3000",
	Handler: router,
    }

    // Start the server on port 3000
    if err := srv.ListenAndServe(); err != nil {
	log.Fatalf("error:", err.Error())
    }
}

```

### [**Grouping Routes**]
```go
func LoggerMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("FROM Global Middleware")
        h.Serve(w, r)
	})
}

func Authenticate(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // your authentication logic goes here
		log.Println("FROM Authenticate Middleware")
        h.Serve(w, r)
	})
}

func middleware1(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // your authentication logic goes here
		log.Println("FROM Authenticate Middleware")
        h.Serve(w, r)
	})
}

func middleware2(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // your authentication logic goes here
		log.Println("FROM Authenticate Middleware")
        h.Serve(w, r)
	})
}

func main() {
    // Initialize a new xpress router
    router := xpress.NewRouter()

    privateRoute := router.Group("/admin").U(LoggerMiddleware, Authenticate)
    publicRoute := router.Group("/users").U(LoggerMiddleware)

    privateRoute.M(m1, m2).HFunc("POST /", func(w http.ResponseWriter, r *http.Request) {})

    // if M will not have any middleware it won't throw any error but for better practice
    // you can remove .M if you're not passing any middleware
    publicRoute.M().HFunc("POST /", func(w http.ResponseWriter, r *http.Request) {})

    srv := &http.Server{
	Addr:    ":3000",
	Handler: router,
    }

    // Start the server on port 3000
    if err := srv.ListenAndServe(); err != nil {
	log.Fatalf("error:", err.Error())
    }
}

```

### [**Route Registration**]
```go
func LoggerMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("FROM Global Middleware")
        h.Serve(w, r)
	})
}

func Authenticate(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // your authentication logic goes here
		log.Println("FROM Authenticate Middleware")
        h.Serve(w, r)
	})
}

func middleware1(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // your authentication logic goes here
		log.Println("FROM Authenticate Middleware")
        h.Serve(w, r)
	})
}

func middleware2(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // your authentication logic goes here
		log.Println("FROM Authenticate Middleware")
        h.Serve(w, r)
	})
}

func main() {
    // Initialize a new xpress router
    router := xpress.NewRouter()

    modules := make([]xpress.Module, 0)
    modules = append(modules, xpress.Module{
		Name: "/admin",
		Middleware: []func(http.Handler) http.Handler{Authenticate},
		Routes: []xpress.Route{
            {
				Method: "POST", Path: "/roles",
				Handler: func(w http.ResponseWriter, r *http.Request) {},
				Middleware: []func(http.Handler) http.Handler{middleware1, middleware2},
			},
            {
				Method: "GET", Path: "/roles",
				Handler: func(w http.ResponseWriter, r *http.Request) {},
				Middleware: []func(http.Handler) http.Handler{middleware2, middleware2},
			},
            {
				Method: "GET", Path: "/roles/{id}",
				Handler: func(w http.ResponseWriter, r *http.Request) {},
				Middleware: []func(http.Handler) http.Handler{middleware1},
			},
            {
				Method: "PUT", Path: "/roles/{id}",
				Handler: func(w http.ResponseWriter, r *http.Request) {},
				Middleware: []func(http.Handler) http.Handler{middleware2},
			},
            {
				Method: "DELETE", Path: "/roles/{id}",
				Handler: func(w http.ResponseWriter, r *http.Request) {},
				Middleware: []func(http.Handler) http.Handler{middleware2, middleware2},
			},
        },
	})

    xpress.RegisterRoutes(modules, router)

    srv := &http.Server{
	Addr:    ":3000",
	Handler: router,
    }

    // Start the server on port 3000
    if err := srv.ListenAndServe(); err != nil {
	log.Fatalf("error:", err.Error())
    }
}

```
