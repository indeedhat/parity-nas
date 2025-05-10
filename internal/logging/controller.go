package logging

import (
	"net/http"
)

// LiveMonitorLogs creates an event stream connection to pass back system logs over
func LiveMonitorLogsController(rw http.ResponseWriter, r *http.Request) {
	l := New("logs")

	readCh := make(chan []byte)

	n := buffer.Connect(readCh, -1)
	defer buffer.Disconnect(readCh)

	go l.WithData(map[string]any{
		"log_count": n,
	}).Info("Live log view opened")

	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")

loop:
	for {
		select {
		case <-r.Context().Done():
			break loop
		case data := <-readCh:
			rw.Write([]byte("data: "))
			rw.Write(data)
			rw.Write([]byte("\n\n"))

			if f, ok := rw.(http.Flusher); ok {
				f.Flush()
			}
		}
	}
}
