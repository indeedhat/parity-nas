package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/indeedhat/parity-nas/internal/config"
	"github.com/indeedhat/parity-nas/internal/routes/context"
	"github.com/indeedhat/parity-nas/internal/status"
)

func LiveMonitor(ctx context.Context) error {
	statusCfg, _ := config.SystemStatus()
	mountCfg, _ := config.Mount()
	netIfaceCfg, _ := config.NetInterface()

	if statusCfg == nil || mountCfg == nil || netIfaceCfg == nil {
		return ctx.InternalError("failed to load config")
	}

	monitor := status.NewMonitor(status.Config{
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

	defer func() {
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}()

	ticker := time.NewTicker(time.Second * time.Duration(statusCfg.PollRate))
	for range ticker.C {
		data, _ := json.Marshal(monitor.Read())

		w.Write(data)
		w.Write([]byte("\n\n"))
	}

	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	return ctx.NoContent()
}
