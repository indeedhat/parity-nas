package main

import (
	"log"
	"net/http"

	"github.com/indeedhat/parity-nas/internal"
	"github.com/indeedhat/parity-nas/internal/config"
	"github.com/indeedhat/parity-nas/internal/env"
	"github.com/indeedhat/parity-nas/internal/servermux"
	"github.com/rs/cors"
)

func main() {
	serverCfg, err := config.Server()
	if err != nil {
		log.Fatalf("Failed to load server config: %s", err)
	}

	proxyCfg, err := config.WebProxy()
	if err != nil {
		log.Fatalf("Failed to load web proxy config: %s", err)
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

	http.ListenAndServe(":8080", c.Handler(mux))
}
