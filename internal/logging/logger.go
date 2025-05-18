package logging

import (
	"log"
	"os"
	"time"

	"github.com/indeedhat/parity-nas/internal/config"
	"github.com/indeedhat/parity-nas/pkg/logging"
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

	buffer = newBuffer(int(cfg.MemoryBufferLen))
	console := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	multiWriter := zerolog.MultiLevelWriter(buffer, fh, console)
	zl = zerolog.New(multiWriter).With().Timestamp().Logger()
}

var (
	fh     *os.File
	zl     zerolog.Logger
	buffer *LogBuffer
)

// Close the file handler used by the logger
func Close() {
	if fh != nil {
		fh.Close()
	}
}

func New(category string) *logging.Logger {
	return logging.New(zl, category)
}
