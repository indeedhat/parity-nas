package routes

import (
	"log"
	"net/http"

	"github.com/indeedhat/parity-nas/internal/auth"
	"github.com/indeedhat/parity-nas/internal/routes/context"
)

func isLoggedIn(next RequestHandler) RequestHandler {
	return func(ctx context.Context) error {
		jwt := auth.ExtractJwtFromAuthHeader(ctx)
		if jwt == "" {
			log.Println("empty jwt")
			return ctx.Error(http.StatusUnauthorized, "Not authorized")
		}

		claims, err := auth.VerifyJwt(jwt)
		if err != nil {
			log.Println("bad jwt")
			return ctx.Error(http.StatusUnauthorized, "Not authorized")
		}

		ctx.Set("user-claims", claims)

		log.Println("jwt ok")
		// TODO: look up and verify the user
		return next(ctx)
	}
}

func isGuest(next RequestHandler) RequestHandler {
	return func(ctx context.Context) error {
		jwt := auth.ExtractJwtFromAuthHeader(ctx)
		if jwt == "" {
			return next(ctx)
		}

		if _, err := auth.VerifyJwt(jwt); err == nil {
			return next(ctx)
		}

		return ctx.Error(http.StatusForbidden, "Already logged in")
	}
}
