package auth

import (
	"net/http"

	"github.com/indeedhat/parity-nas/internal/servermux"
)

// IsLoggedInMiddleware will only accept requests from users with a valid login JWT
func IsLoggedInMiddleware(next servermux.RequestHandler) servermux.RequestHandler {
	return func(ctx servermux.Context) error {
		jwt := extractJwtFromAuthHeader(ctx)
		if jwt == "" {
			if jwt = ctx.Request().URL.Query().Get("bearer"); jwt == "" {
				return ctx.Error(http.StatusUnauthorized, "Not authorized")
			}
		}

		claims, err := verifyJwt(jwt)
		if err != nil {
			return ctx.Error(http.StatusUnauthorized, "Not authorized")
		}

		ctx.Set("user-claims", claims)

		// TODO: look up and verify the user
		return next(ctx)
	}
}

// IsGuestMiddleware will only accept requests from users withot a valid login JWT
func IsGuestMiddleware(next servermux.RequestHandler) servermux.RequestHandler {
	return func(ctx servermux.Context) error {
		jwt := extractJwtFromAuthHeader(ctx)
		if jwt == "" {
			return next(ctx)
		}

		if _, err := verifyJwt(jwt); err == nil {
			return next(ctx)
		}

		return ctx.Error(http.StatusForbidden, "Already logged in")
	}
}
