package config

import (
	"os"

	"github.com/indeedhat/parity-nas/pkg/config"
)

const ServerKey = "server"

type ServerCfg struct {
	Version uint `icl:"version"`

	MaxBodySize int64 `icl:"max_body_size"`
}

// Server initializes a ServerConfig struct
//
// If a config file exists then it will attempt to load it otherwise a new instance will be created
func Server() (*ServerCfg, error) {
	var c ServerCfg

	if err := config.Load(ServerKey, &c); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}

		c = ServerCfg{
			Version:     1,
			MaxBodySize: 1 << 20, // 1MB
		}
	}

	return &c, nil
}
