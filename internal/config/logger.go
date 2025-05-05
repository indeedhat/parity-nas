package config

import (
	"os"
)

const LoggerKey = "logger"

type LoggerCfg struct {
	Version uint `icl:"version"`

	SavePath        string `icl:"save_path"`
	MemoryBufferLen uint   `icl:"memory_buffer_len"`
}

// Logger initializes a LoggerCfg struct
//
// If a config file exists then it will attempt to load it otherwise a new instance will be created
func Logger() (*LoggerCfg, error) {
	var c LoggerCfg

	if err := loadConfig(LoggerKey, &c); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}

		c = LoggerCfg{
			Version:  1,
			SavePath: "parinas.log",
		}
	}

	return &c, nil
}
