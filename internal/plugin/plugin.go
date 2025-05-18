package plugin

import (
	"errors"
	"os"
	"os/exec"
	"path"
	"plugin"

	"github.com/indeedhat/parity-nas/internal/config"
	"github.com/indeedhat/parity-nas/pkg/logging"
	servermux "github.com/indeedhat/parity-nas/pkg/server_mux"
)

type PluginManager struct {
	log    *logging.Logger
	cfg    *config.PluginCfg
	router servermux.Router
}

func NewManager(cfg *config.PluginCfg, log *logging.Logger, router servermux.Router) PluginManager {
	return PluginManager{log, cfg, router}
}

func (m PluginManager) Init() error {
	for _, entry := range m.cfg.Plugins {
		if err := m.installPlugin(entry); err != nil {
			return err
		}

		if err := m.initializePlugin(entry); err != nil {
			return err
		}
	}

	return nil
}

func (m PluginManager) installPlugin(entry config.PluginEntry) error {
	if m.checkForExisting(entry) {
		return nil
	}
	m.log.Infof("installing plugin %s", entry.Name())

	if err := m.downloadArchive(entry); err != nil {
		m.log.Error("download fail")
		return err
	}
	// NB: this will also clean up the extracted files from the next stage if it gets that far
	defer m.cleanupArchive(entry)

	if err := m.extractArchive(entry); err != nil {
		m.log.Error("extract fail")
		return err
	}

	if err := m.buildGoBinary(entry); err != nil {
		m.log.Error("build fail")
		return err
	}

	return nil
}

func (m PluginManager) buildGoBinary(entry config.PluginEntry) error {
	cwd, _ := os.Getwd()
	soPath := path.Join(cwd, entry.SharedObjectPath(m.cfg))

	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", soPath, "main.go")
	cmd.Dir = entry.ArchiveExtractPath(m.cfg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (m PluginManager) initializePlugin(entry config.PluginEntry) error {
	m.log.Infof("initializing plugin %s", entry.Name())

	p, err := plugin.Open(entry.SharedObjectPath(m.cfg))
	if err != nil {
		return err
	}

	f, err := p.Lookup("Init")
	if err != nil {
		return err
	}

	initPlugin, ok := f.(func(router servermux.Router, logger *logging.Logger) error)
	if !ok {
		return errors.New("invalid signature for Init function")
	}

	return initPlugin(m.router, m.log.WithCategory(entry.Name()))
}

func (m PluginManager) checkForExisting(entry config.PluginEntry) bool {
	_, err := os.Stat(entry.SharedObjectPath(m.cfg))

	return !errors.Is(err, os.ErrNotExist)
}
