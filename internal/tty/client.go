package tty

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
)

const (
	msgTypeIo     = "io"
	msgTypeNotice = "notice"
	msgTypeResize = "resize"
)

type resizeMsg struct {
	Cols uint16 `json:"cols"`
	Rows uint16 `json:"rows"`
}

type client struct {
	ws        *websocket.Conn
	ctx       context.Context
	ctxCancel context.CancelFunc
	ptmx      *os.File
}

func newClient(
	ctx context.Context,
	cancel context.CancelFunc,
	conn *websocket.Conn,
	ptmx *os.File,
) *client {
	return &client{
		ws:        conn,
		ctx:       ctx,
		ctxCancel: cancel,
		ptmx:      ptmx,
	}
}

func (c *client) Run() {
	go c.ptyToSock()
	go c.sockToPty()
}

func (c *client) ptyToSock() {
	buf := make([]byte, 1024)

	for {
		n, err := c.ptmx.Read(buf)
		if err != nil {
			c.closeWithNotice("ws read failure")
			return
		}

		if err := c.ws.WriteMessage(websocket.TextMessage, buf[:n]); err != nil {
			c.closeWithNotice("tty write failure")
			return
		}
	}
}

func (c *client) sockToPty() {
	for {
		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			log.Print(err)
			c.closeWithNotice("tty read failure")
			return
		}

		parts := strings.SplitN(string(msg), ":", 2)
		if len(parts) != 2 {
			c.closeWithNotice("msg parse failure")
			return
		}

		switch parts[0] {
		case msgTypeIo:
			if _, err := c.ptmx.Write(msg); err != nil {
				c.closeWithNotice("ws write failure")
				return
			}
		case msgTypeNotice:
			// pass
		case msgTypeResize:
			var r resizeMsg
			if err := json.Unmarshal(msg, &r); err != nil {
				c.closeWithNotice("resize msg parse failure")
			}

			err := pty.Setsize(c.ptmx, &pty.Winsize{
				Cols: r.Cols,
				Rows: r.Rows,
			})
			if err != nil {
				c.closeWithNotice("tty resize failure")
				return
			}
		default:
			c.closeWithNotice("msg parse failure")
			return
		}
	}
}

func (c *client) closeWithNotice(msg string) {
	log.Printf("close: %s", msg)
	c.ws.WriteMessage(websocket.TextMessage, []byte(msgTypeNotice+":"+msg))

	c.ctxCancel()
}
