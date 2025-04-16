package config

import (
	"io/fs"

	"github.com/indeedhat/icl"
	"github.com/indeedhat/parity-nas/internal/env"
)

const fileExt = ".icl"

type baseConfig struct {
	Version uint `icl:"version"`
}

// Save attempts to save a config struct to file by its type key
func Save(key string, v any) error {
	path := configPath(key)

	return icl.MarshalFile(v, path)
}

// configPath produces a save path for a config by type key
func configPath(key string) string {
	if path := env.ConfigPath.Get(); path != "" {
		return path + key + fileExt
	}

	return "/etc/parinas/" + key + fileExt
}

// loadConfig attempts to load a config struct from file by its type key
func loadConfig(key string, v any) error {
	path := configPath(key)

	err := icl.UnMarshalFile(path, v)
	if e, ok := err.(*fs.PathError); ok {
		err = e.Err
	}

	return err
}
