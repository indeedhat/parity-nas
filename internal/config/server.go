package config

import (
	"io/fs"
	"os"
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

	if err := loadConfig(ServerKey, &c); err != nil {
		switch e := err.(type) {
		case *fs.PathError:
			err = e.Err
		}

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
