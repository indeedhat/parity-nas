package config

import (
	"github.com/indeedhat/parity-nas/internal/servermux"
)

type configEntry struct {
	Config any
	Error  error
}

// ViewConfigController is a debug only controller that will return the current state of the full
// icl config tree
func ViewConfigController(ctx servermux.Context) error {
	config := make(map[string]any)

	net, err := NetInterface()
	config["network"] = configEntry{net, err}

	mount, err := Mount()
	config["mount"] = configEntry{mount, err}

	sysmon, err := SystemStatus()
	config["sysmon"] = configEntry{sysmon, err}

	server, err := Server()
	config["server"] = configEntry{server, err}

	return ctx.Ok(config)
}
