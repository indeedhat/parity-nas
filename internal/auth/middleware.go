package auth

import (
	"net/http"

	"github.com/indeedhat/parity-nas/internal/servermux"
)

var (
	PermissionAny   uint8 = 0
	PermissionAdmin uint8 = 1
	PermissionWrite uint8 = 2
	PermissionRead  uint8 = 4
)

// IsLoggedInMiddleware will only accept requests from users with a valid login JWT
func IsLoggedInMiddleware(next servermux.RequestHandler) servermux.RequestHandler {
	return UserHasPermissionMiddleware(PermissionAny)(next)
}

// UserHasPermissionMiddleware checks if the logged in user has a specific permission level
func UserHasPermissionMiddleware(level uint8) func(servermux.RequestHandler) servermux.RequestHandler {
	return func(next servermux.RequestHandler) servermux.RequestHandler {
		return func(ctx *servermux.Context) error {
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

			if claims.Permission&level != level {
				return ctx.Error(http.StatusForbidden, "Forbidden")
			}

			return next(ctx)
		}
	}
}

// IsGuestMiddleware will only accept requests from users withot a valid login JWT
func IsGuestMiddleware(next servermux.RequestHandler) servermux.RequestHandler {
	return func(ctx *servermux.Context) error {
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
