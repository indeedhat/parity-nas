package logging

import (
	"context"
	"log"
	"log/slog"
)

type splitHandler []slog.Handler

func (s splitHandler) Enabled(ctx context.Context, l slog.Level) bool {
	for _, h := range s {
		if h.Enabled(ctx, l) {
			return true
		}
	}
	return false
}

func (s splitHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, h := range s {
		if err := h.Handle(ctx, r); err != nil {
			log.Printf("failed to write slog: %s", err)
		}
	}

	return nil
}

func (s splitHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	attrHandler := make(splitHandler, 0, len(s))

	for _, h := range s {
		attrHandler = append(attrHandler, h.WithAttrs(attrs))
	}

	return attrHandler
}

func (s splitHandler) WithGroup(category string) slog.Handler {
	groupHandler := make(splitHandler, 0, len(s))

	for _, h := range s {
		groupHandler = append(groupHandler, h.WithGroup(category))
	}

	return groupHandler
}

var _ slog.Handler = (*splitHandler)(nil)

type leveler struct {
}
