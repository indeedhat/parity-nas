package logging

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/indeedhat/parity-nas/pkg/logging"
	"github.com/indeedhat/parity-nas/pkg/server_mux"
)

func LoggingMiddleware(logger *logging.Logger) servermux.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrw := &responseWrapper{ResponseWriter: rw}

			defer func() {
				logger.WithData(logging.Data{
					"method":   r.Method,
					"url":      r.URL.String(),
					"status":   wrw.status,
					"size":     wrw.size,
					"duration": time.Since(start),
				}).Info("")
			}()

			next(wrw, r)
		}
	}
}

type responseWrapper struct {
	http.ResponseWriter
	headerWritten bool
	status        int
	size          int
}

func (w *responseWrapper) WriteHeader(code int) {
	if w.headerWritten {
		return
	}

	w.status = code
	w.headerWritten = true
	w.ResponseWriter.WriteHeader(code)
}

func (w *responseWrapper) Write(buf []byte) (int, error) {
	w.WriteHeader(http.StatusOK)
	n, err := w.ResponseWriter.Write(buf)
	w.size += n
	return n, err
}

func (w *responseWrapper) Flush() {
	w.ResponseWriter.(http.Flusher).Flush()
}

func (w *responseWrapper) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

func (w *responseWrapper) ReadFrom(r io.Reader) (int64, error) {
	w.WriteHeader(http.StatusOK)

	n, err := w.ResponseWriter.(io.ReaderFrom).ReadFrom(r)
	w.size += int(n)
	return n, err
}

var (
	_ http.Flusher  = (*responseWrapper)(nil)
	_ http.Hijacker = (*responseWrapper)(nil)
	_ io.ReaderFrom = (*responseWrapper)(nil)
	_ http.Flusher  = (*responseWrapper)(nil)
)
