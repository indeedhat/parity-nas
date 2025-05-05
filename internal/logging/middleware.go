package logging

import (
	"net/http"
	"time"

	"github.com/indeedhat/parity-nas/internal/servermux"
)

func LoggingMiddleware(logger *Logger) func(servermux.RequestHandler) servermux.RequestHandler {
	return func(next servermux.RequestHandler) servermux.RequestHandler {
		return func(ctx *servermux.Context) error {
			start := time.Now()

			rw := &responseWrapper{ResponseWriter: ctx.Writer()}
			ctx.ReplaceWriter(rw)

			defer func() {
				logger.Zerolog().Info().
					Str("category", logger.Category()).
					Str("method", ctx.Request().Method).
					Stringer("url", ctx.Request().URL).
					Int("status", rw.status).
					Int("size", rw.size).
					Dur("duration", time.Since(start)).
					Msg("")
			}()

			return next(ctx)
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
