package routes

import (
	"log"
	"net/http"

	"github.com/indeedhat/parity-nas/internal/routes/context"
)

type RequestHandler func(context.Context) error
type Middleware func(RequestHandler) RequestHandler

type Router struct {
	mux        *http.ServeMux
	middleware []Middleware
	basePath   string
}

func newRouter(middleware ...Middleware) Router {
	return Router{
		mux:        http.DefaultServeMux,
		middleware: middleware,
	}
}

func (r Router) Get(path string, handler RequestHandler, middleware ...Middleware) {
	r.mux.HandleFunc("GET "+r.basePath+path, r.wrap(handler, middleware...))
}

func (r Router) Post(path string, handler RequestHandler, middleware ...Middleware) {
	log.Println("POST " + r.basePath + path)
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
		ctx := context.New(w, req)
		err := r.apply(handler, middleware...)(ctx)

		switch resp := err.(type) {
		case nil:
			return
		case context.Response:
			w.WriteHeader(resp.Code())
			w.Header().Set("Content-Type", "application/json")
			w.Write(resp.Data())
		default:
			errResp := ctx.InternalError(err.Error()).(context.Response)
			w.WriteHeader(errResp.Code())
			w.Header().Set("Content-Type", "application/json")
			w.Write(errResp.Data())
		}
	}
}
