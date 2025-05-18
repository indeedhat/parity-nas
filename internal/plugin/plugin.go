package plugin

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"plugin"

	parinas "github.com/indeedhat/parity-nas/internal"
	"github.com/indeedhat/parity-nas/internal/config"
	"github.com/indeedhat/parity-nas/pkg/logging"
	servermux "github.com/indeedhat/parity-nas/pkg/server_mux"
)

// these types for the plugin functions are kind of useless, i cannot use them in type checking
// but they are here to provide a reference for the possible function signatures
type InitFunc func(logger *logging.Logger) error
type RouterFunc func(router servermux.Router, logger *logging.Logger) error
type CloserFunc func() error

type PluginManager struct {
	log     *logging.Logger
	cfg     *config.PluginCfg
	router  servermux.Router
	closers map[string]func() error
}

func NewManager(cfg *config.PluginCfg, log *logging.Logger, router servermux.Router) *PluginManager {
	return &PluginManager{
		log:     log,
		cfg:     cfg,
		router:  router,
		closers: make(map[string]func() error),
	}
}

func (m *PluginManager) Init() error {
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

func (m *PluginManager) Close() map[string]error {
	errors := make(map[string]error, len(m.closers))

	for name, closer := range m.closers {
		errors[name] = closer()
	}

	return errors
}

func (m *PluginManager) installPlugin(entry config.PluginEntry) error {
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

func (m *PluginManager) buildGoBinary(entry config.PluginEntry) error {
	cwd, _ := os.Getwd()
	soPath := path.Join(cwd, entry.SharedObjectPath(m.cfg))

	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", soPath, "main.go")
	cmd.Dir = entry.ArchiveExtractPath(m.cfg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (m *PluginManager) initializePlugin(entry config.PluginEntry) error {
	m.log.Infof("initializing plugin %s", entry.Name())

	p, err := plugin.Open(entry.SharedObjectPath(m.cfg))
	if err != nil {
		return err
	}

	f, err := p.Lookup("Init")
	if err != nil {
		return err
	}

	logger := m.log.WithCategory("plugin").WithAttr("plugin", entry.Name())

	initPlugin, ok := f.(func(logger *logging.Logger) error)
	if !ok {
		return errors.New("invalid signature for Init function")
	}
	if err := initPlugin(logger); err != nil {
		return err
	}

	routers := map[string]uint8{
		"PublicRoutes": parinas.PermissionPublic,
		"GuestRoutes":  parinas.PermissionGuest,
		"UserRoutes":   parinas.PermissionUser,
		"AdminRoutes":  parinas.PermissionAdmin,
	}
	for name, permission := range routers {
		f, err := p.Lookup(name)
		if err != nil {
			continue
		}

		initFunc, ok := f.(func(router servermux.Router, logger *logging.Logger) error)
		if !ok {
			return fmt.Errorf("invalid signature for %s function", name)
		}

		if err := initFunc(parinas.PluginRouter(m.router, permission, entry.Name()), logger); err != nil {
			return err
		}
	}

	if f, err = p.Lookup("Close"); err != nil {
		return nil
	}

	closerFunc, ok := f.(func() error)
	if !ok {
		return fmt.Errorf("invalid signature for closer function")
	}

	m.closers[entry.Name()] = closerFunc

	return nil
}

func (m *PluginManager) checkForExisting(entry config.PluginEntry) bool {
	_, err := plugin.Open(entry.SharedObjectPath(m.cfg))

	return err == nil
}
