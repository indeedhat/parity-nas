package config

import (
	"os"
)

const WebProxyKey = "web_proxy"

type WebProxyCfg struct {
	Version uint `icl:"version"`

	Prefix string `icl:"prefix"`

	Handlers []WebProxyHandler `icl:"handler"`
}

type WebProxyHandler struct {
	Prefix   string           `icl:".param"`
	Scheme   string           `icl:"scheme"`
	Host     string           `icl:"host"`
	Port     uint16           `icl:"port"`
	Mutators WebProxyMutators `icl:"mutators"`
}

type WebProxyMutators struct {
	BaseTag   bool `icl:"base_tag"`
	HtmlLinks bool `icl:"html_links"`
	JsImports bool `icl:"js_imports"`
}

// WebProxy initializes a WebProxyCfg struct
//
// If a config file exists then it will attempt to load it otherwise a new instance will be created
func WebProxy() (*WebProxyCfg, error) {
	var c WebProxyCfg

	if err := loadConfig(WebProxyKey, &c); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}

		c = WebProxyCfg{
			Version: 1,
			Prefix:  "_",
		}
	}

	return &c, nil
}

// FIndHandler searches the proxy handler slice for a handler with the matching prefix
//
// If there are multiple handlers with the same prefix then the last in the list will be returned
func (c *WebProxyCfg) FIndHandler(prefix string) *WebProxyHandler {
	for i := len(c.Handlers) - 1; i >= 0; i-- {
		if prefix == c.Handlers[i].Prefix {
			return &c.Handlers[i]
		}
	}

	return nil
}
