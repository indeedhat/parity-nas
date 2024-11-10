package auth

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/indeedhat/parity-nas/internal/env"
	"github.com/indeedhat/parity-nas/internal/routes/context"
)

type UserClaims struct {
	jwt.RegisteredClaims

	UserName string `json:"nme"`
	UserId   string `json:"uid"`
}

var ErrInvalidJWT = errors.New("Invalid jwt")

var jwtSecret = os.Getenv(env.JwtSecret)

// ExtractJwtFromAuthHeader will verify that the Authorization header both exists and is in the
// Bearer format, if so it will extract the token (hopefully this should be a valid JWT)
func ExtractJwtFromAuthHeader(ctx context.Context) string {
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
func VerifyJwt(jwtString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != "HS256" {
			return nil, ErrInvalidJWT
		}

		return []byte(jwtSecret), nil
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

// GenerateJWT will generate a new JWT for the given account model
func GenerateJWT(claims jwt.Claims) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte(jwtSecret))
}

func GenerateUserJwt(id, name string) (string, error) {
	return GenerateJWT(UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID: strconv.Itoa(int(time.Now().Unix())),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(
				time.Duration(env.GetEnvInt(env.JwtTTl)),
			)),
		},
		UserName: name,
		UserId:   id,
	})
}
