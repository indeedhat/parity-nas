package config

import (
	"os"
)

const TtyKey = "tty"

type TtyCfg struct {
	Version uint `icl:"version"`

	Shell    string `icl:"shell"`
	StartDir string `icl:"start_dir"`
}

// Tty initializes a TtyCfg struct
//
// If a config file exists then it will attempt to load it otherwise a new instance will be created
func Tty() (*TtyCfg, error) {
	var c TtyCfg

	if err := loadConfig(ServerKey, &c); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}

		c = TtyCfg{
			Version:  1,
			Shell:    "bash",
			StartDir: "/root",
		}
	}

	return &c, nil
}
