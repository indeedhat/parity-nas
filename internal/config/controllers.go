package config

import (
	"github.com/indeedhat/parity-nas/internal/servermux"
)

type configEntry struct {
	Config any
	Error  error
}

func ViewConfigController(ctx servermux.Context) error {
	config := make(map[string]any)

	net, err := NetInterface()
	config["network"] = configEntry{net, err}

	mount, err := Mount()
	config["mount"] = configEntry{mount, err}

	sysmon, err := SystemStatus()
	config["sysmon"] = configEntry{sysmon, err}

	return ctx.Ok(config)
}
