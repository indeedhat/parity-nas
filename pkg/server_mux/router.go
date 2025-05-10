package servermux

import (
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

type ServerConfig struct {
	MaxBodySize int64
}

type Router struct {
	mux        *http.ServeMux
	middleware []Middleware
	basePath   string
	cfg        ServerConfig
}

// NewRouter creates a new router instance with the provided middleware stack assigned
func NewRouter(serverCfg ServerConfig, middleware ...Middleware) Router {
	return Router{
		mux:        http.DefaultServeMux,
		middleware: middleware,
		cfg:        serverCfg,
	}
}

// ServerMux returns the underlying http.ServeMux instance
func (r Router) ServerMux() *http.ServeMux {
	return r.mux
}

// Get registers a handler on the GET method on the provided uri
func (r Router) Get(path string, handler http.HandlerFunc, middleware ...Middleware) {
	r.mux.HandleFunc("GET "+r.basePath+path, r.wrap(handler, middleware...))
}

// Post registers a handler on the POST method on the provided uri
func (r Router) Post(path string, handler http.HandlerFunc, middleware ...Middleware) {
	r.mux.HandleFunc("POST "+r.basePath+path, r.wrap(handler, middleware...))
}

// Put registers a handler on the PUT method on the provided uri
func (r Router) Put(path string, handler http.HandlerFunc, middleware ...Middleware) {
	r.mux.HandleFunc("PUT "+r.basePath+path, r.wrap(handler, middleware...))
}

// Patch registers a handler on the PATCH method on the provided uri
func (r Router) Patch(path string, handler http.HandlerFunc, middleware ...Middleware) {
	r.mux.HandleFunc("PATCH "+r.basePath+path, r.wrap(handler, middleware...))
}

// Delete registers a handler on the DELETE method on the provided uri
func (r Router) Delete(path string, handler http.HandlerFunc, middleware ...Middleware) {
	r.mux.HandleFunc("DELETE "+r.basePath+path, r.wrap(handler, middleware...))
}

// All registers a handler on all request methods on the provided uri
func (r Router) All(path string, handler http.HandlerFunc, middleware ...Middleware) {
	r.mux.HandleFunc(path, r.wrap(handler, middleware...))
}

// Group creates a sub router and assigns a base path and middleware to all routes assigned within it
func (r Router) Group(path string, middleware ...Middleware) Router {
	return Router{
		mux:        r.mux,
		basePath:   r.basePath + path,
		middleware: append(r.middleware, middleware...),
		cfg:        r.cfg,
	}
}

// apply middleware to the handler
func (r Router) apply(handler http.HandlerFunc, middleware ...Middleware) http.HandlerFunc {
	stack := append(r.middleware, middleware...)

	for i := range stack {
		handler = stack[len(stack)-1-i](handler)
	}

	return handler
}

// wrap handler in a response parsing closure
func (r Router) wrap(handler http.HandlerFunc, middleware ...Middleware) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		ctx := NewContext(req.Context(), r.cfg)
		req = req.WithContext(ctx)

		r.apply(handler, middleware...)(rw, req)
	}
}
