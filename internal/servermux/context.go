package servermux

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Context struct {
	rw   http.ResponseWriter
	req  *http.Request
	data map[string]any
	cfg  ServerConfig
}

// NewContext creates a new instance of the servermux context from the provided ResponseWriter and Request params
func NewContext(cfg ServerConfig, rw http.ResponseWriter, r *http.Request) Context {
	return Context{
		rw:   rw,
		req:  r,
		data: make(map[string]any),
		cfg:  cfg,
	}
}

// WithWriter allows you to create a new Context from the current one with a different ResponseWriter
// This allows you to wrap the ResponseWriter from within middleware
func (c Context) WithWriter(rw http.ResponseWriter) Context {
	return Context{
		rw:   rw,
		req:  c.req,
		data: c.data,
		cfg:  c.cfg,
	}
}

// Request returns the underlying Request instance
func (c Context) Request() *http.Request {
	return c.req
}

// Writer returns the underlying ResponseWriter instance
func (c Context) Writer() http.ResponseWriter {
	return c.rw
}

// Body reads in the Request body as a byte array
func (c Context) Body() []byte {
	data, err := io.ReadAll(io.LimitReader(c.req.Body, c.cfg.MaxBodySize))
	if err != nil {
		return nil
	}

	c.req.Body = io.NopCloser(bytes.NewBuffer(data))
	return data
}

// UnmarshalBody unmarshales the request body into the provided data structure
//
// NB: this is JSON only
func (c Context) UnmarshalBody(v any) error {
	// 1 MB limit
	data, err := io.ReadAll(io.LimitReader(c.req.Body, 2<<20))
	if err != nil {
		return nil
	}

	c.req.Body = io.NopCloser(bytes.NewBuffer(data))

	return json.Unmarshal(data, v)
}

// Validate runs the provided struct against its validation tags
func (c Context) Validate(v any) error {
	checker := validator.New()
	return checker.Struct(v)
}

// Error is a convenience methdod for returning an error from a controller
func (c Context) Error(code int, v any) error {
	if msg, ok := v.(string); ok {
		return c.Response(code, errorResponse{msg})
	}
	return c.Response(code, v)
}

// InternalError is a convenience method for returning a 500 error from a controller
func (c Context) InternalError(msg string) error {
	return c.Response(http.StatusInternalServerError, errorResponse{msg})
}

// InternalErrorf is a convenience method for returning a 500 error from a controller that inclueds
// string formatting
func (c Context) InternalErrorf(msg string, a ...any) error {
	return c.Response(
		http.StatusInternalServerError,
		errorResponse{fmt.Sprintf(msg, a...)},
	)
}

// Response is a convenience method for constructing a response to return from a controller
func (c Context) Response(code int, v any) error {
	if v == nil {
		return Response{code: code}
	}

	data, err := json.Marshal(v)
	if err != nil {
		return Response{http.StatusInternalServerError, []byte(`{"error":"failed to generate response json"}`)}
	}

	return Response{code, data}
}

// Ok is a convenience method for returning a 200 response from a controller
func (c Context) Ok(v any) error {
	return c.Response(http.StatusOK, v)
}

// NoContent is a convenience method for returning a 201 response from a controller
func (c Context) NoContent() error {
	return c.Response(http.StatusOK, nil)
}

// Get a value from the context's key value store
func (c Context) Get(key string) (any, bool) {
	val, ok := c.data[key]
	return val, ok
}

// Set a value in the context's key value store
func (c Context) Set(key string, value any) {
	c.data[key] = value
}

type Response struct {
	code int
	data []byte
}

// Error implements the error interface
func (r Response) Error() string {
	return string(r.data)
}

// Data returns the underlying responses byte array
func (r Response) Data() []byte {
	return r.data
}

// Code returns the internal http status code for the response
func (r Response) Code() int {
	return r.code
}

type errorResponse struct {
	Error string `json:"error"`
}
