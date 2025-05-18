package config

import (
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/indeedhat/parity-nas/pkg/config"
)

const PluginKey = "plugins"

type PluginCfg struct {
	Version uint `icl:"version"`

	SavePath string        `icl:"save_path"`
	TempPath string        `icl:"temp_path"`
	Plugins  []PluginEntry `icl:"plugin"`
}

type PluginEntry struct {
	GithubLink string `icl:".param"`
	Version    string `icl:"version"`
}

func (e PluginEntry) Name() string {
	parts := strings.Split(e.GithubLink, "/")
	return parts[len(parts)-1]
}

func (e PluginEntry) ArchiveUrl() string {
	u, err := url.Parse(e.GithubLink)
	if err != nil {
		return ""
	}

	u.Scheme = "https"
	final, _ := url.JoinPath(u.String(), "/archive/refs/tags/"+e.Version+".zip")

	return final
}

func (e PluginEntry) ArchiveSavePath(cfg *PluginCfg) string {
	return e.ArchiveExtractPath(cfg) + ".zip"
}

func (e PluginEntry) ArchiveExtractPath(cfg *PluginCfg) string {
	return path.Join(cfg.TempPath, e.Name()+"-"+e.Version)
}

func (e PluginEntry) SharedObjectPath(cfg *PluginCfg) string {
	return path.Join(cfg.SavePath, e.Name()+"-"+e.Version) + ".so"
}

// Server initializes a ServerConfig struct
//
// If a config file exists then it will attempt to load it otherwise a new instance will be created
func Plugins() (*PluginCfg, error) {
	var c PluginCfg

	if err := config.Load(PluginKey, &c); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}

		c = PluginCfg{
			Version: 1,
		}
	}

	return &c, nil
}
