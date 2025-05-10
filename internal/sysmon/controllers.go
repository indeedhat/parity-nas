package sysmon

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/indeedhat/parity-nas/internal/config"
	"github.com/indeedhat/parity-nas/pkg/server_mux"
)

// LiveMonitorController creates an event stream connection to pass back system stats over
func LiveMonitorController(rw http.ResponseWriter, r *http.Request) {
	statusCfg, err := config.SystemStatus()
	if err != nil {
		servermux.InternalErrorf(rw, "failed to load status config %s", err)
		return
	}

	mountCfg, err := config.Mount()
	if err != nil {
		servermux.InternalErrorf(rw, "failed to load mount config %s", err)
		return
	}

	netIfaceCfg, err := config.NetInterface()
	if err != nil {
		servermux.InternalErrorf(rw, "failed to load netif config %s", err)
		return
	}

	monitor := NewMonitor(Config{
		PollRate:      statusCfg.PollRate,
		Mounts:        mountCfg.Tracked,
		NetInterfaces: netIfaceCfg.Tracked,
	})

	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")

	ticker := time.NewTicker(time.Second * time.Duration(statusCfg.PollRate))

loop:
	for {
		select {
		case <-r.Context().Done():
			break loop
		case <-ticker.C:
			data, _ := json.Marshal(monitor.Read())

			rw.Write([]byte("data: "))
			rw.Write(data)
			rw.Write([]byte("\n\n"))

			if f, ok := rw.(http.Flusher); ok {
				f.Flush()
			}
		}
	}
}
