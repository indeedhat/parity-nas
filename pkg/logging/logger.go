package logging

import (
	"fmt"

	"github.com/rs/zerolog"
)

type Logger struct {
	zerolog  *zerolog.Logger
	category string
}

func New(zl zerolog.Logger, category string) *Logger {
	return &Logger{
		zerolog:  &zl,
		category: category,
	}
}

func (l *Logger) WithCategory(category string) *Logger {
	return &Logger{
		zerolog:  l.zerolog,
		category: category,
	}
}

func (l *Logger) Zerolog() *zerolog.Logger {
	return l.zerolog
}

func (l *Logger) Category() string {
	return l.category
}

func (l *Logger) Trace(msg string) {
	l.zerolog.Trace().Str("category", l.category).Msg(msg)
}

func (l *Logger) Tracef(msg string, a ...any) {
	l.zerolog.Trace().Str("category", l.category).Msg(fmt.Sprintf(msg, a...))
}

func (l *Logger) Debug(msg string) {
	l.zerolog.Debug().Str("category", l.category).Msg(msg)
}

func (l *Logger) Debugf(msg string, a ...any) {
	l.zerolog.Debug().Str("category", l.category).Msg(fmt.Sprintf(msg, a...))
}

func (l *Logger) Info(msg string) {
	l.zerolog.Info().Str("category", l.category).Msg(msg)
}

func (l *Logger) Infof(msg string, a ...any) {
	l.zerolog.Info().Str("category", l.category).Msg(fmt.Sprintf(msg, a...))
}

func (l *Logger) Warning(msg string) {
	l.zerolog.Warn().Str("category", l.category).Msg(msg)
}

func (l *Logger) Warningf(msg string, a ...any) {
	l.zerolog.Warn().Str("category", l.category).Msg(fmt.Sprintf(msg, a...))
}

func (l *Logger) Error(msg string) {
	l.zerolog.Error().Str("category", l.category).Msg(msg)
}

func (l *Logger) Errorf(msg string, a ...any) {
	l.zerolog.Error().Str("category", l.category).Msg(fmt.Sprintf(msg, a...))
}

func (l *Logger) Fatal(msg string) {
	l.zerolog.Fatal().Str("category", l.category).Msg(msg)
}

func (l *Logger) Fatalf(msg string, a ...any) {
	l.zerolog.Fatal().Str("category", l.category).Msg(fmt.Sprintf(msg, a...))
}

func (l *Logger) WithData(data map[string]any) *DataLogEntry {
	return &DataLogEntry{
		Logger: l,
		data:   data,
	}
}

type DataLogEntry struct {
	*Logger
	data map[string]any
}

func (l *DataLogEntry) Trace(msg string) {
	l.zerolog.Trace().Str("category", l.category).Fields(l.data).Msg(msg)
}

func (l *DataLogEntry) Tracef(msg string, a ...any) {
	l.zerolog.Trace().Str("category", l.category).Fields(l.data).Msg(fmt.Sprintf(msg, a...))
}

func (l *DataLogEntry) Debug(msg string) {
	l.zerolog.Debug().Str("category", l.category).Fields(l.data).Msg(msg)
}

func (l *DataLogEntry) Debugf(msg string, a ...any) {
	l.zerolog.Debug().Str("category", l.category).Fields(l.data).Msg(fmt.Sprintf(msg, a...))
}

func (l *DataLogEntry) Info(msg string) {
	l.zerolog.Info().Str("category", l.category).Fields(l.data).Msg(msg)
}

func (l *DataLogEntry) Infof(msg string, a ...any) {
	l.zerolog.Info().Str("category", l.category).Fields(l.data).Msg(fmt.Sprintf(msg, a...))
}

func (l *DataLogEntry) Warning(msg string) {
	l.zerolog.Warn().Str("category", l.category).Fields(l.data).Msg(msg)
}

func (l *DataLogEntry) Warningf(msg string, a ...any) {
	l.zerolog.Warn().Str("category", l.category).Fields(l.data).Msg(fmt.Sprintf(msg, a...))
}

func (l *DataLogEntry) Error(msg string) {
	l.zerolog.Error().Str("category", l.category).Fields(l.data).Msg(msg)
}

func (l *DataLogEntry) Errorf(msg string, a ...any) {
	l.zerolog.Error().Str("category", l.category).Fields(l.data).Msg(fmt.Sprintf(msg, a...))
}

func (l *DataLogEntry) Fatal(msg string) {
	l.zerolog.Fatal().Str("category", l.category).Fields(l.data).Msg(msg)
}

func (l *DataLogEntry) Fatalf(msg string, a ...any) {
	l.zerolog.Fatal().Str("category", l.category).Fields(l.data).Msg(fmt.Sprintf(msg, a...))
}
