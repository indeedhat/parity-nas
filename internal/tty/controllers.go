package tty

import (
	"context"
	"net/http"
	"os/exec"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
	"github.com/indeedhat/parity-nas/internal/config"
	"github.com/indeedhat/parity-nas/internal/servermux"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: verify host
		return true
	},
}

// TtyController creates a websocket connection for an interactive shell session
func TtyController(ctx *servermux.Context) error {
	cfg, err := config.Tty()
	if err != nil {
		return ctx.InternalError("Failed to load config")
	}

	conn, err := upgrader.Upgrade(ctx.Writer(), ctx.Request(), nil)
	if err != nil {
		return ctx.InternalError(err.Error())
	}
	defer conn.Close()

	cmd := exec.Command(cfg.Shell)
	cmd.Dir = cfg.StartDir

	ptmx, err := pty.Start(cmd)
	defer ptmx.Close()
	defer cmd.Process.Kill()

	ptyCtx, cancel := context.WithCancel(context.Background())

	client := newClient(ptyCtx, cancel, conn, ptmx)
	client.Run()

	select {
	case <-ctx.Request().Context().Done():
	case <-ptyCtx.Done():
	}

	return nil
}
