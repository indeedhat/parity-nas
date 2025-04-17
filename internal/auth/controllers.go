package auth

import (
	"net/http"

	"github.com/indeedhat/parity-nas/internal/servermux"
)

type loginRequest struct {
	User   string `json:"user" validate:"required"`
	Passwd string `json:"passwd" validate:"required"`
}

// LoginController handles user login attempts
func LoginController(ctx servermux.Context) error {
	var req loginRequest
	if err := ctx.UnmarshalBody(&req); err != nil {
		return ctx.Error(http.StatusUnprocessableEntity, "Unprocessale Content")
	}

	if err := ctx.Validate(req); err != nil {
		return ctx.Error(http.StatusUnprocessableEntity, err)
	}

	user := attemptSystemLogin(req.User, req.Passwd)
	if user == nil {
		return ctx.Error(http.StatusUnauthorized, "login failed")
	}

	jwt, err := GenerateUserJwt("1", req.User, user.Permission)
	if err != nil {
		return ctx.InternalError("Failed to process login")
	}

	ctx.Writer().Header().Set("auth_token", "jwt."+jwt)
	return ctx.NoContent()
}

// VerifyLoginController returns the current status of the login
//
// actually it does nothing, it just allows the middleware to run
func VerifyLoginController(ctx servermux.Context) error {
	return ctx.NoContent()
}
