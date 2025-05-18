package logging

import (
	"fmt"
	"log/slog"

	"golang.org/x/net/context"
)

type Logger struct {
	slog     *slog.Logger
	category string
	ctx      context.Context
}

func New(sl *slog.Logger, category string) *Logger {
	return &Logger{
		slog:     sl.With(slog.String("category", category)),
		category: category,
		ctx:      context.Background(),
	}
}

func NewWithContext(ctx context.Context, sl *slog.Logger, category string) *Logger {
	return &Logger{
		slog:     sl.With(slog.String("category", category)),
		category: category,
		ctx:      ctx,
	}
}

func (l *Logger) WithCategory(category string) *Logger {
	return &Logger{
		slog:     l.slog.With(slog.String("category", category)),
		category: category,
		ctx:      l.ctx,
	}
}

func (l *Logger) WithCategoryContext(ctx context.Context, category string) *Logger {
	return &Logger{
		slog:     l.slog.With(slog.String("category", category)),
		category: category,
		ctx:      ctx,
	}
}

func (l *Logger) Slog() *slog.Logger {
	return l.slog
}

func (l *Logger) Category() string {
	return l.category
}

func (l *Logger) Debug(msg string) {
	l.slog.DebugContext(l.ctx, msg)
}

func (l *Logger) Debugf(msg string, a ...any) {
	l.slog.DebugContext(l.ctx, fmt.Sprintf(msg, a...))
}

func (l *Logger) Info(msg string) {
	l.slog.InfoContext(l.ctx, msg)
}

func (l *Logger) Infof(msg string, a ...any) {
	l.slog.InfoContext(l.ctx, fmt.Sprintf(msg, a...))
}

func (l *Logger) Warn(msg string) {
	l.slog.WarnContext(l.ctx, msg)
}

func (l *Logger) Warningf(msg string, a ...any) {
	l.slog.WarnContext(l.ctx, fmt.Sprintf(msg, a...))
}

func (l *Logger) Error(msg string) {
	l.slog.ErrorContext(l.ctx, msg)
}

func (l *Logger) Errorf(msg string, a ...any) {
	l.slog.ErrorContext(l.ctx, fmt.Sprintf(msg, a...))
}

func (l *Logger) Fatal(msg string) {
	l.slog.ErrorContext(l.ctx, msg, "fatal_error", true)
	panic(msg)
}

func (l *Logger) Fatalf(msg string, a ...any) {
	l.slog.With("fatal_error", true).ErrorContext(l.ctx, fmt.Sprintf(msg, a...))
	panic(msg)
}

func (l *Logger) WithData(data Data) *Logger {
	return &Logger{
		slog:     l.slog.With(data.asList()...),
		category: l.category,
		ctx:      l.ctx,
	}
}

type Data map[string]any

func (d Data) asList() []any {
	list := make([]any, 0, len(d)*2)

	for k, v := range d {
		list = append(list, k, v)
	}

	return list
}
