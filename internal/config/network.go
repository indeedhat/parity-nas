package config

import (
	"errors"
	"io/fs"
)

const NetInterfaceKey = "disk"

type NetInterfaceCfg struct {
	baseConfig

	Tracked []string `hcl:"tracked_ifaces"`
}

// NetInterface initializes a NetInterfaceCfg struct
//
// If a config file exists then it will attempt to load it otherwise a new instance will be created
func NetInterface() (*NetInterfaceCfg, error) {
	var c *NetInterfaceCfg

	if err := loadConfig(SystemStatusKey, &c); err != nil {
		if !errors.Is(fs.ErrNotExist, err) {
			return nil, err
		}

		c = &NetInterfaceCfg{
			baseConfig: baseConfig{Version: 1},
		}
	}

	return c, nil
}
