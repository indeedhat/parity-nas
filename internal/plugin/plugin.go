package plugin

import (
	"errors"
	"os"
	"os/exec"

	"github.com/indeedhat/parity-nas/internal/config"
	"github.com/indeedhat/parity-nas/internal/logging"
)

type PluginManager struct {
	log *logging.Logger
	cfg *config.PluginCfg
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

	if err := m.downloadArchive(entry); err != nil {
		return err
	}
	// NB: this will also clean up the extracted files from the next stage if it gets that far
	defer m.cleanupArchive(entry)

	if err := m.extractArchive(entry); err != nil {
		return err
	}

	if err := m.buildGoBinary(entry); err != nil {
		return err
	}

	return nil
}

func (m PluginManager) buildGoBinary(entry config.PluginEntry) error {
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", entry.SharedObjectPath(m.cfg), ".")
	cmd.Dir = entry.ArchiveExtractPath(m.cfg)

	return cmd.Run()
}

func (m PluginManager) initializePlugin(entry config.PluginEntry) error {
	return nil
}

func (m PluginManager) checkForExisting(entry config.PluginEntry) bool {
	_, err := os.Stat(entry.SharedObjectPath(m.cfg))

	return errors.Is(err, os.ErrNotExist)
}
