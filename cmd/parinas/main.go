package main

import (
	"net/http"

	"github.com/indeedhat/parity-nas/internal/env"
	"github.com/indeedhat/parity-nas/internal/routes"
	"github.com/rs/cors"
)

func main() {
	mux := routes.Build()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{env.Get(env.CorsAllowHost)},
		AllowCredentials: true,
		ExposedHeaders:   []string{"Auth_token"},
	})

	http.ListenAndServe(":8080", c.Handler(mux))
}
