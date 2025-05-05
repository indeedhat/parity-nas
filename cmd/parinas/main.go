package main

import (
	"net/http"

	"github.com/indeedhat/parity-nas/internal"
	"github.com/indeedhat/parity-nas/internal/config"
	"github.com/indeedhat/parity-nas/internal/env"
	"github.com/indeedhat/parity-nas/internal/logging"
	"github.com/indeedhat/parity-nas/internal/servermux"
	"github.com/rs/cors"
)

func main() {
	logger := logging.New("parinas")

	serverCfg, err := config.Server()
	if err != nil {
		logger.Fatalf("Failed to load server config: %s", err)
	}

	proxyCfg, err := config.WebProxy()
	if err != nil {
		logger.Fatalf("Failed to load web proxy config: %s", err)
	}

	mux := parinas.BuildRoutes(
		servermux.ServerConfig{
			MaxBodySize: serverCfg.MaxBodySize,
		},
		proxyCfg,
	)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{env.CorsAllowHost.Get()},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Authorization"},
		ExposedHeaders:   []string{"Auth_token"},
	})

	logger.Infof("ListenAndServer: %v", http.ListenAndServe(":8080", c.Handler(mux)))
}
