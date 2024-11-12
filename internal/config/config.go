package config

import (
	"os"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/indeedhat/parity-nas/internal/env"
)

type baseConfig struct {
	Version uint `hcl:"version"`
}

// Save attempts to save a config struct to file by its type key
func Save(key string, v any) error {
	path := configPath(key)

	f := hclwrite.NewFile()
	gohcl.EncodeIntoBody(v, f.Body())

	return os.WriteFile(path, f.Bytes(), 0644)
}

// configPath produces a save path for a config by type key
func configPath(key string) string {
	if path := env.ConfigPath.Get(); path != "" {
		return path + key + ".hcl"
	}

	return "/etc/parinas/" + key + ".hcl"
}

// loadConfig attempts to load a config struct from file by its type key
func loadConfig(key string, v any) error {
	path := configPath(key)

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return hclsimple.Decode(path, data, nil, v)
}
