package config

import (
	"errors"
	"io/fs"
)

const SystemStatusKey = "system_status"

type SystemStatusCfg struct {
	Version uint `hcl:"version"`

	PollRate uint8 `hcl:"poll_rate"`
}

// SystemStatus initializes a SystemStatusCfg struct
//
// If a config file exists then it will attempt to load it otherwise a new instance will be created
func SystemStatus() (*SystemStatusCfg, error) {
	var c SystemStatusCfg

	if err := loadConfig(SystemStatusKey, &c); err != nil {
		if !errors.Is(fs.ErrNotExist, err) {
			return nil, err
		}

		c = SystemStatusCfg{
			Version:  1,
			PollRate: 2,
		}
	}

	return &c, nil
}
