package main

import (
	"net/http"

	"github.com/indeedhat/parity-nas/internal"
	"github.com/indeedhat/parity-nas/internal/env"
	"github.com/rs/cors"
)

func main() {
	mux := parinas.BuildRoutes()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{env.CorsAllowHost.Get()},
		AllowCredentials: true,
		ExposedHeaders:   []string{"Auth_token"},
	})

	http.ListenAndServe(":8080", c.Handler(mux))
}
