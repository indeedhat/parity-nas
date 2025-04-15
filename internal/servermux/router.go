package servermux

import (
	"net/http"
)

type RequestHandler func(Context) error
type Middleware func(RequestHandler) RequestHandler

type Router struct {
	mux        *http.ServeMux
	middleware []Middleware
	basePath   string
}

func NewRouter(middleware ...Middleware) Router {
	return Router{
		mux:        http.DefaultServeMux,
		middleware: middleware,
	}
}

func (r Router) ServerMux() *http.ServeMux {
	return r.mux
}

func (r Router) Get(path string, handler RequestHandler, middleware ...Middleware) {
	r.mux.HandleFunc("GET "+r.basePath+path, r.wrap(handler, middleware...))
}

func (r Router) Post(path string, handler RequestHandler, middleware ...Middleware) {
	r.mux.HandleFunc("POST "+r.basePath+path, r.wrap(handler, middleware...))
}

func (r Router) Put(path string, handler RequestHandler, middleware ...Middleware) {
	r.mux.HandleFunc("PUT "+r.basePath+path, r.wrap(handler, middleware...))
}

func (r Router) Patch(path string, handler RequestHandler, middleware ...Middleware) {
	r.mux.HandleFunc("PATCH "+r.basePath+path, r.wrap(handler, middleware...))
}

func (r Router) Delete(path string, handler RequestHandler, middleware ...Middleware) {
	r.mux.HandleFunc("DELETE "+r.basePath+path, r.wrap(handler, middleware...))
}

func (r Router) Group(path string, middleware ...Middleware) Router {
	return Router{
		mux:        r.mux,
		basePath:   r.basePath + path,
		middleware: append(r.middleware, middleware...),
	}
}

func (r Router) apply(handler RequestHandler, middleware ...Middleware) RequestHandler {
	stack := append(r.middleware, middleware...)

	for i := len(stack) - 1; i >= 0; i-- {
		handler = stack[i](handler)
	}

	return handler
}

func (r Router) wrap(handler RequestHandler, middleware ...Middleware) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := NewContext(w, req)
		err := r.apply(handler, middleware...)(ctx)

		switch resp := err.(type) {
		case nil:
			return
		case Response:
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(resp.Code())
			w.Write(resp.Data())
		default:
			errResp := ctx.InternalError(err.Error()).(Response)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(errResp.Code())
			w.Write(errResp.Data())
		}
	}
}
