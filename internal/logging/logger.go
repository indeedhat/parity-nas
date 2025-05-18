package logging

import (
	"log"
	"log/slog"
	"os"

	"github.com/indeedhat/parity-nas/internal/config"
	"github.com/indeedhat/parity-nas/pkg/logging"
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

	opts := &slog.HandlerOptions{Level: slog.LevelDebug}
	sl = slog.New(splitHandler{
		slog.NewTextHandler(os.Stdout, opts),
		slog.NewJSONHandler(fh, opts),
		slog.NewJSONHandler(buffer, opts),
	})
}

var (
	sl     *slog.Logger
	fh     *os.File
	buffer *LogBuffer
)

// Close the file handler used by the logger
func Close() {
	if fh != nil {
		fh.Close()
	}
}

func New(category string) *logging.Logger {
	return logging.New(sl, category)
}
