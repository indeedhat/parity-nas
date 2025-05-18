package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	parinas "github.com/indeedhat/parity-nas/internal"
	"github.com/indeedhat/parity-nas/internal/config"
	"github.com/indeedhat/parity-nas/internal/env"
	"github.com/indeedhat/parity-nas/internal/logging"
	"github.com/indeedhat/parity-nas/internal/plugin"
	servermux "github.com/indeedhat/parity-nas/pkg/server_mux"
	"github.com/rs/cors"
)

var flagUsePlugins bool

func main() {
	logger := logging.New("parinas")

	flag.BoolVar(&flagUsePlugins, "with-plugin-support", false, "enable plugin support")
	flag.Parse()

	serverCfg, err := config.Server()
	if err != nil {
		logger.Fatalf("Failed to load server config: %s", err)
	}

	proxyCfg, err := config.WebProxy()
	if err != nil {
		logger.Fatalf("Failed to load web proxy config: %s", err)
	}

	router := servermux.NewRouter(servermux.ServerConfig{
		MaxBodySize: serverCfg.MaxBodySize,
	})

	mux := parinas.BuildRoutes(
		router,
		proxyCfg,
	)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{env.CorsAllowHost.Get()},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Authorization"},
		ExposedHeaders:   []string{"Auth_token"},
	})

	if flagUsePlugins {
		pluginCfg, err := config.Plugins()
		if err != nil {
			logger.Fatalf("Failed to load plugin config: %s", err)
		}

		pm := plugin.NewManager(pluginCfg, logger, router)
		defer pm.Close()

		if err := pm.Init(); err != nil {
			logger.Fatalf("Failed to initialize plugins: %s", err)
		}
	}

	svr := &http.Server{
		Addr:    ":8080",
		Handler: c.Handler(mux),
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Info("Starting server on :8080")
		logger.Infof("ListenAndServer: %v", svr.ListenAndServe())
	}()

	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := svr.Shutdown(ctx); err != nil {
		logger.Info("Server forced to shutdown after timeout")
	}
}
