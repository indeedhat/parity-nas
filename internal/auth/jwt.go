package auth

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/indeedhat/parity-nas/internal/env"
)

var ErrInvalidJWT = errors.New("Invalid jwt")

type UserClaims struct {
	jwt.RegisteredClaims

	UserName   string `json:"nme"`
	UserId     string `json:"uid"`
	Permission uint8  `json:"lvl"`
}

// GenerateJWT will generate a new JWT for the given account model
func GenerateJWT(claims jwt.Claims) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte(env.JwtSecret.Get()))
}

// GenerateUserJwt genertes a new JWT specifically for a user login session
func GenerateUserJwt(id, name string, permission uint8) (string, error) {
	return GenerateJWT(UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID: strconv.Itoa(int(time.Now().Unix())),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(
				time.Duration(env.JwtTTl.Get()) * time.Second,
			)),
		},
		UserName:   name,
		UserId:     id,
		Permission: permission,
	})
}

// extractJwtFromAuthHeader will verify that the Authorization header both exists and is in the
// Bearer format, if so it will extract the token (hopefully this should be a valid JWT)
func extractJwtFromAuthHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
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
func verifyJwt(jwtString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(jwtString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != "HS256" {
			return nil, ErrInvalidJWT
		}

		return []byte(env.JwtSecret.Get()), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidJWT
	}

	return claims, nil
}
