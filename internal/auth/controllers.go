package auth

import (
	"net/http"

	"github.com/indeedhat/parity-nas/pkg/server_mux"
)

type loginRequest struct {
	User   string `json:"user" validate:"required"`
	Passwd string `json:"passwd" validate:"required"`
}

// LoginController handles user login attempts
func LoginController(rw http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := servermux.UnmarshalBody(r, &req); err != nil {
		servermux.WriteError(rw, http.StatusUnprocessableEntity, "Unprocessale Content")
		return
	}

	if err := servermux.Validate(req); err != nil {
		servermux.WriteError(rw, http.StatusUnprocessableEntity, err)
		return
	}

	user := attemptSystemLogin(req.User, req.Passwd)
	if user == nil {
		servermux.WriteError(rw, http.StatusUnauthorized, "login failed")
	}

	jwt, err := GenerateUserJwt("1", req.User, user.Permission)
	if err != nil {
		servermux.InternalError(rw, "Failed to process login")
	}

	rw.Header().Set("auth_token", "jwt."+jwt)
	servermux.NoContent(rw)
}

// VerifyLoginController returns the current status of the login
//
// actually it does nothing, it just allows the middleware to run
func VerifyLoginController(rw http.ResponseWriter, _ *http.Request) {
	servermux.NoContent(rw)
}
