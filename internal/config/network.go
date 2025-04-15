package config

import (
	"errors"
	"io/fs"
)

const NetInterfaceKey = "net"

type NetInterfaceCfg struct {
	Version uint `icl:"version"`

	Tracked []string `icl:"tracked_ifaces"`
}

// NetInterface initializes a NetInterfaceCfg struct
//
// If a config file exists then it will attempt to load it otherwise a new instance will be created
func NetInterface() (*NetInterfaceCfg, error) {
	var c NetInterfaceCfg

	if err := loadConfig(NetInterfaceKey, &c); err != nil {
		if !errors.Is(fs.ErrNotExist, err) {
			return nil, err
		}

		c = NetInterfaceCfg{
			Version: 1,
		}
	}

	return &c, nil
}
