package logging

import (
	"log"
	"os"
	"time"

	"github.com/indeedhat/parity-nas/internal/config"
	"github.com/rs/zerolog"
)

func init() {
	cfg, err := config.Logger()
	if err != nil {
		log.Fatal("failed to initialize default logger")
	}

	fh, err := os.OpenFile(cfg.SavePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal("failed to open log file")
	}

	buffer := newBuffer(int(cfg.MemoryBufferLen))
	console := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	multiWriter := zerolog.MultiLevelWriter(buffer, fh, console)
	zl = zerolog.New(multiWriter).With().Timestamp().Logger()
}

var (
	fh *os.File
	zl zerolog.Logger
)

// Close the file handler used by the logger
func Close() {
	if fh != nil {
		fh.Close()
	}
}

type Logger struct {
	zerolog  *zerolog.Logger
	category string
}

func New(category string) *Logger {
	return &Logger{
		zerolog:  &zl,
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

func (l *Logger) Debug(msg string) {
	l.zerolog.Debug().Str("category", l.category).Msg(msg)
}

func (l *Logger) Info(msg string) {
	l.zerolog.Info().Str("category", l.category).Msg(msg)
}

func (l *Logger) Warning(msg string) {
	l.zerolog.Warn().Str("category", l.category).Msg(msg)
}

func (l *Logger) Error(msg string) {
	l.zerolog.Error().Str("category", l.category).Msg(msg)
}

func (l *Logger) Fatal(msg string) {
	l.zerolog.Fatal().Str("category", l.category).Msg(msg)
}
