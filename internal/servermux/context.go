package servermux

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Context struct {
	rw   http.ResponseWriter
	req  *http.Request
	data map[string]any
}

func NewContext(rw http.ResponseWriter, r *http.Request) Context {
	return Context{
		rw:   rw,
		req:  r,
		data: make(map[string]any)}
}

func (c Context) Request() *http.Request {
	return c.req
}

func (c Context) Writer() http.ResponseWriter {
	return c.rw
}

func (c Context) Body() []byte {
	// 1 MB limit
	data, err := io.ReadAll(io.LimitReader(c.req.Body, 2<<20))
	if err != nil {
		return nil
	}

	c.req.Body = io.NopCloser(bytes.NewBuffer(data))
	return data
}

func (c Context) UnmarshalBody(v any) error {
	// 1 MB limit
	data, err := io.ReadAll(io.LimitReader(c.req.Body, 2<<20))
	if err != nil {
		return nil
	}

	c.req.Body = io.NopCloser(bytes.NewBuffer(data))

	return json.Unmarshal(data, v)
}

func (c Context) Validate(v any) error {
	checker := validator.New()
	return checker.Struct(v)
}

func (c Context) Error(code int, v any) error {
	if msg, ok := v.(string); ok {
		return c.Response(code, errorResponse{msg})
	}
	return c.Response(code, v)
}

func (c Context) InternalError(msg string) error {
	return c.Response(http.StatusInternalServerError, errorResponse{msg})
}

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

func (c Context) Ok(v any) error {
	return c.Response(http.StatusOK, v)
}

func (c Context) NoContent() error {
	return c.Response(http.StatusOK, nil)
}

func (c Context) Get(key string) (any, bool) {
	val, ok := c.data[key]
	return val, ok
}

func (c Context) Set(key string, value any) {
	c.data[key] = value
}

type Response struct {
	code int
	data []byte
}

func (r Response) Error() string {
	return string(r.data)
}

func (r Response) Data() []byte {
	return r.data
}

func (r Response) Code() int {
	return r.code
}

type errorResponse struct {
	Error string `json:"error"`
}
