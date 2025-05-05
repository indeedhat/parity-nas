package servermux

import (
	"net/http"
)

type RequestHandler func(*Context) error
type Middleware func(RequestHandler) RequestHandler

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
func NewRouter(serverCfg ServerConfig, logger Middleware, middleware ...Middleware) Router {
	return Router{
		mux: http.DefaultServeMux,
		middleware: append(
			[]Middleware{logger, responseWriterMiddleware},
			middleware...,
		),
		cfg: serverCfg,
	}
}

// ServerMux returns the underlying http.ServeMux instance
func (r Router) ServerMux() *http.ServeMux {
	return r.mux
}

// Get registers a handler on the GET method on the provided uri
func (r Router) Get(path string, handler RequestHandler, middleware ...Middleware) {
	r.mux.HandleFunc("GET "+r.basePath+path, r.wrap(handler, middleware...))
}

// Post registers a handler on the POST method on the provided uri
func (r Router) Post(path string, handler RequestHandler, middleware ...Middleware) {
	r.mux.HandleFunc("POST "+r.basePath+path, r.wrap(handler, middleware...))
}

// Put registers a handler on the PUT method on the provided uri
func (r Router) Put(path string, handler RequestHandler, middleware ...Middleware) {
	r.mux.HandleFunc("PUT "+r.basePath+path, r.wrap(handler, middleware...))
}

// Patch registers a handler on the PATCH method on the provided uri
func (r Router) Patch(path string, handler RequestHandler, middleware ...Middleware) {
	r.mux.HandleFunc("PATCH "+r.basePath+path, r.wrap(handler, middleware...))
}

// Delete registers a handler on the DELETE method on the provided uri
func (r Router) Delete(path string, handler RequestHandler, middleware ...Middleware) {
	r.mux.HandleFunc("DELETE "+r.basePath+path, r.wrap(handler, middleware...))
}

// All registers a handler on all request methods on the provided uri
func (r Router) All(path string, handler RequestHandler, middleware ...Middleware) {
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
func (r Router) apply(handler RequestHandler, middleware ...Middleware) RequestHandler {
	stack := append(r.middleware, middleware...)

	for i := len(stack) - 1; i >= 0; i-- {
		handler = stack[i](handler)
	}

	return handler
}

// wrap handler in a response parsing closure
func (r Router) wrap(handler RequestHandler, middleware ...Middleware) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := NewContext(r.cfg, w, req)
		r.apply(handler, middleware...)(ctx)
	}
}

func responseWriterMiddleware(next RequestHandler) RequestHandler {
	return func(ctx *Context) error {
		err := next(ctx)

		switch resp := err.(type) {
		case nil:
			// pass
		case Response:
			ctx.Writer().Header().Set("Content-Type", "application/json; charset=utf-8")
			ctx.Writer().WriteHeader(resp.Code())
			ctx.Writer().Write(resp.Data())
		default:
			errResp := ctx.InternalError(err.Error()).(Response)
			ctx.Writer().Header().Set("Content-Type", "application/json; charset=utf-8")
			ctx.Writer().WriteHeader(errResp.Code())
			ctx.Writer().Write(errResp.Data())
		}

		return nil
	}
}
