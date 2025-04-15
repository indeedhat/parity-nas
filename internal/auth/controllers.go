package auth

import (
	"net/http"

	"github.com/indeedhat/parity-nas/internal/servermux"
)

type loginRequest struct {
	User   string `json:"user" validate:"required"`
	Passwd string `json:"passwd" validate:"required"`
}

func LoginController(ctx servermux.Context) error {
	var req loginRequest
	if err := ctx.UnmarshalBody(&req); err != nil {
		return ctx.Error(http.StatusUnprocessableEntity, "Unprocessale Content")
	}

	if err := ctx.Validate(req); err != nil {
		return ctx.Error(http.StatusUnprocessableEntity, err)
	}

	// TODO: need to actually handle user auth

	jwt, err := GenerateUserJwt("1", req.User)
	if err != nil {
		return ctx.InternalError("Failed to process login")
	}

	ctx.Writer().Header().Set("auth_token", "jwt."+jwt)
	return ctx.NoContent()
}

func VerifyLoginController(ctx servermux.Context) error {
	// NB: This controller only really exists to allow the auth middleware to run, it doesn't
	//     actually need to do anything itself
	return ctx.NoContent()
}
