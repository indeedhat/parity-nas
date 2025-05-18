package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/indeedhat/parity-nas/internal"
	"github.com/indeedhat/parity-nas/internal/config"
	"github.com/indeedhat/parity-nas/internal/env"
	"github.com/indeedhat/parity-nas/internal/logging"
	"github.com/indeedhat/parity-nas/internal/plugin"
	"github.com/indeedhat/parity-nas/pkg/server_mux"
	"github.com/rs/cors"
)

var flagUsePlugins bool

func main() {
	logger := logging.New("parinas")

	flag.BoolVar(&flagUsePlugins, "with-plugin-support", false, "enable plugin support")
	flag.Parse()

	log.Print(flagUsePlugins)

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
		if err := pm.Init(); err != nil {
			logger.Fatalf("Failed to initialize plugins: %s", err)
		}
	}

	logger.Infof("ListenAndServer: %v", http.ListenAndServe(":8080", c.Handler(mux)))
}
