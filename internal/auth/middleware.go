package auth

import (
	"net/http"

	"github.com/indeedhat/parity-nas/pkg/server_mux"
)

var (
	PermissionAny   uint8 = 0
	PermissionAdmin uint8 = 1
	PermissionWrite uint8 = 2
	PermissionRead  uint8 = 4
)

// IsLoggedInMiddleware will only accept requests from users with a valid login JWT
func IsLoggedInMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return UserHasPermissionMiddleware(PermissionAny)(next)
}

// UserHasPermissionMiddleware checks if the logged in user has a specific permission level
func UserHasPermissionMiddleware(level uint8) servermux.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, r *http.Request) {
			jwt := extractJwtFromAuthHeader(r)
			if jwt == "" {
				if jwt = r.URL.Query().Get("bearer"); jwt == "" {
					servermux.WriteError(rw, http.StatusUnauthorized, "Not authorized")
					return
				}
			}

			claims, err := verifyJwt(jwt)
			if err != nil {
				servermux.WriteError(rw, http.StatusUnauthorized, "Not authorized")
				return
			}

			r = r.WithContext(r.Context().(servermux.Context).WithData("user-claims", claims))

			if claims.Permission&level != level {
				servermux.WriteError(rw, http.StatusForbidden, "Forbidden")
				return
			}

			next(rw, r)
		}
	}
}

// IsGuestMiddleware will only accept requests from users withot a valid login JWT
func IsGuestMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		jwt := extractJwtFromAuthHeader(r)
		if jwt == "" {
			next(rw, r)
			return
		}

		if _, err := verifyJwt(jwt); err == nil {
			next(rw, r)
			return
		}

		servermux.WriteError(rw, http.StatusForbidden, "Already logged in")
	}
}
