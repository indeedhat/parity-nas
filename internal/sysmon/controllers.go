package sysmon

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/indeedhat/parity-nas/internal/config"
	"github.com/indeedhat/parity-nas/internal/servermux"
)

// LiveMonitorController creates an event stream connection to pass back system stats over
func LiveMonitorController(ctx servermux.Context) error {
	statusCfg, e1 := config.SystemStatus()
	mountCfg, e2 := config.Mount()
	netIfaceCfg, e3 := config.NetInterface()

	if statusCfg == nil || mountCfg == nil || netIfaceCfg == nil {
		return ctx.InternalError(fmt.Sprintf("failed to load config %s %s %s", e1, e2, e3))
	}

	monitor := NewMonitor(Config{
		PollRate:      statusCfg.PollRate,
		Mounts:        mountCfg.Tracked,
		NetInterfaces: netIfaceCfg.Tracked,
	})

	w := ctx.Writer()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ticker := time.NewTicker(time.Second * time.Duration(statusCfg.PollRate))

loop:
	for {
		select {
		case <-ctx.Request().Context().Done():
			break loop
		case <-ticker.C:
			data, _ := json.Marshal(monitor.Read())

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
