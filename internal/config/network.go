package config

import (
	"os"

	"github.com/indeedhat/parity-nas/pkg/config"
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

	if err := config.Load(NetInterfaceKey, &c); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}

		c = NetInterfaceCfg{
			Version: 1,
		}
	}

	return &c, nil
}
