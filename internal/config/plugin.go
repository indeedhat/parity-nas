package config

import (
	"os"
)

const PluginKey = "plugins"

type PluginCfg struct {
	Version uint `icl:"version"`

	SavePath string        `icl:"save_path"`
	Plugins  []PluginEntry `icl:"plugin"`
}

type PluginEntry struct {
	GithubLink string         `icl:".param"`
	Version    string         `icl:"version"`
	Data       map[string]any `icl:"data"`
}

// Server initializes a ServerConfig struct
//
// If a config file exists then it will attempt to load it otherwise a new instance will be created
func Plugins() (*PluginCfg, error) {
	var c PluginCfg

	if err := loadConfig(PluginKey, &c); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}

		c = PluginCfg{
			Version: 1,
		}
	}

	return &c, nil
}
