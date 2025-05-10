package config

import (
	"os"

	"github.com/indeedhat/parity-nas/pkg/config"
)

const SystemStatusKey = "system_status"

type SystemStatusCfg struct {
	Version uint `icl:"version"`

	PollRate uint8 `icl:"poll_rate"`
}

// SystemStatus initializes a SystemStatusCfg struct
//
// If a config file exists then it will attempt to load it otherwise a new instance will be created
func SystemStatus() (*SystemStatusCfg, error) {
	var c SystemStatusCfg

	if err := config.Load(SystemStatusKey, &c); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}

		c = SystemStatusCfg{
			Version:  1,
			PollRate: 2,
		}
	}

	return &c, nil
}
