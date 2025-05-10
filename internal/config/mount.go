package config

import (
	"os"

	"github.com/indeedhat/parity-nas/pkg/config"
)

const MountKey = "disk"

type MountCfg struct {
	Version uint `icl:"version"`

	Tracked []string `icl:"tracked_disks"`
}

// Mount initializes a MountCfg struct
//
// If a config file exists then it will attempt to load it otherwise a new instance will be created
func Mount() (*MountCfg, error) {
	var c MountCfg

	if err := config.Load(MountKey, &c); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}

		c = MountCfg{
			Version: 1,
		}
	}

	return &c, nil
}
