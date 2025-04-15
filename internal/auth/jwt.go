package auth

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/indeedhat/parity-nas/internal/env"
	"github.com/indeedhat/parity-nas/internal/servermux"
)

var ErrInvalidJWT = errors.New("Invalid jwt")

type UserClaims struct {
	jwt.RegisteredClaims

	UserName string `json:"nme"`
	UserId   string `json:"uid"`
}

// GenerateJWT will generate a new JWT for the given account model
func GenerateJWT(claims jwt.Claims) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte(env.JwtSecret.Get()))
}

// GenerateUserJwt genertes a new JWT specifically for a user login session
func GenerateUserJwt(id, name string) (string, error) {
	return GenerateJWT(UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID: strconv.Itoa(int(time.Now().Unix())),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(
				time.Duration(env.JwtTTl.Get()) * time.Second,
			)),
		},
		UserName: name,
		UserId:   id,
	})
}

// extractJwtFromAuthHeader will verify that the Authorization header both exists and is in the
// Bearer format, if so it will extract the token (hopefully this should be a valid JWT)
func extractJwtFromAuthHeader(ctx servermux.Context) string {
	authHeader := ctx.Request().Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		return ""
	}

	return parts[1]
}

// VerifyJwt will check that the JWT is both a jwt and valid
func verifyJwt(jwtString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != "HS256" {
			return nil, ErrInvalidJWT
		}

		return []byte(env.JwtSecret.Get()), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		return nil, ErrInvalidJWT
	} else {
		return claims, nil
	}
}
