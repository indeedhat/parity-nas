package servermux

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// Body reads in the Request body as a byte array
func Body(r *http.Request) []byte {
	cfg := r.Context().(Context).Config()

	data, err := io.ReadAll(io.LimitReader(r.Body, cfg.MaxBodySize))
	if err != nil {
		return nil
	}

	r.Body = io.NopCloser(bytes.NewBuffer(data))

	return data
}

// UnmarshalBody unmarshales the request body into the provided data structure
//
// NB: this is JSON only
func UnmarshalBody(r *http.Request, v any) error {
	data := Body(r)
	if data == nil {
		return errors.New("could not read request body")
	}

	return json.Unmarshal(data, v)
}

// Validate runs the provided struct against its validation tags
func Validate(v any) error {
	checker := validator.New()
	return checker.Struct(v)
}

// InternalError is a convenience method for returning a 500 error from a controller
func InternalError(rw http.ResponseWriter, msg string) {
	WriteResponse(rw, http.StatusInternalServerError, errorResponse{msg})
}

// InternalErrorf is a convenience method for returning a 500 error from a controller that inclueds
// string formatting
func InternalErrorf(rw http.ResponseWriter, msg string, a ...any) {
	WriteResponse(
		rw,
		http.StatusInternalServerError,
		errorResponse{fmt.Sprintf(msg, a...)},
	)
}

func WriteError(rw http.ResponseWriter, code int, v any) {
	switch val := v.(type) {
	case string:
		WriteResponse(rw, code, errors.New(val))
	default:
		WriteResponse(rw, code, v)
	}
}

// Response is a convenience method for constructing a response to return from a controller
func WriteResponse(rw http.ResponseWriter, code int, v any) {
	var resp *response

	switch val := v.(type) {
	case response:
		resp = &val
	case nil:
		resp = &response{code: code}
	case error:
		v = errorResponse{val.Error()}
	}

	if resp == nil {
		data, err := json.Marshal(v)
		if err != nil {
			resp = &response{http.StatusInternalServerError, []byte(`{"error":"failed to generate response json"}`)}
		}

		resp = &response{code, data}
	}

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(resp.code)
	rw.Write(resp.data)
}

// Ok is a convenience method for returning a 200 response from a controller
func Ok(rw http.ResponseWriter, v any) {
	WriteResponse(rw, http.StatusOK, v)
}

// NoContent is a convenience method for returning a 201 response from a controller
func NoContent(rw http.ResponseWriter) {
	WriteResponse(rw, http.StatusOK, nil)
}

type response struct {
	code int
	data []byte
}

type errorResponse struct {
	Error string `json:"error"`
}
