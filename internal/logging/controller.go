package logging

import (
	"net/http"

	"github.com/indeedhat/parity-nas/internal/servermux"
)

// LiveMonitorLogs creates an event stream connection to pass back system logs over
func LiveMonitorLogsController(ctx *servermux.Context) error {
	readCh := make(chan []byte)

	buffer.Connect(readCh, -1)
	defer buffer.Disconnect(readCh)

	w := ctx.Writer()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

loop:
	for {
		select {
		case <-ctx.Request().Context().Done():
			break loop
		case data := <-readCh:
			w.Write([]byte("data: "))
			w.Write(data)
			w.Write([]byte("\n\n"))

			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		}
	}

	return nil
}
