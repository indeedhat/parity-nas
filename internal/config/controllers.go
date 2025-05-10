package config

import (
	"net/http"

	servermux "github.com/indeedhat/parity-nas/pkg/server_mux"
)

type configEntry struct {
	Config any
	Error  error
}

// ViewConfigController is a debug only controller that will return the current state of the full
// icl config tree
func ViewConfigController(rw http.ResponseWriter, r *http.Request) {
	config := make(map[string]any)

	net, err := NetInterface()
	config["network"] = configEntry{net, err}

	mount, err := Mount()
	config["mount"] = configEntry{mount, err}

	sysmon, err := SystemStatus()
	config["sysmon"] = configEntry{sysmon, err}

	server, err := Server()
	config["server"] = configEntry{server, err}

	servermux.Ok(rw, config)
}
