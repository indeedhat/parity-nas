package config

import (
	"errors"
	"io/fs"
)

const MountKey = "disk"

type MountCfg struct {
	baseConfig

	Tracked []string `hcl:"tracked_disks"`
}

// Mount initializes a MountCfg struct
//
// If a config file exists then it will attempt to load it otherwise a new instance will be created
func Mount() (*MountCfg, error) {
	var c *MountCfg

	if err := loadConfig(SystemStatusKey, &c); err != nil {
		if !errors.Is(fs.ErrNotExist, err) {
			return nil, err
		}

		c = &MountCfg{
			baseConfig: baseConfig{Version: 1},
		}
	}

	return c, nil
}
