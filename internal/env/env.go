// Define constants for all the environment variables that are passed to the webapp
// via the .env config file
package env

import (
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
)

const (
	// Auth
	JwtSecret = "JWT_SECRET"
	// Time since jwt generation that will cause the jwt to be refreshed
	JwtRefreshAge = "JWT_REFRESH_AGE"
	JwtTTl        = "JWT_TTL"
)

const (
	// Web server stuffs
	WebDomain     = "WEB_ROOT"
	CorsAllowHost = "CORS_ALLOW_HOST"
)

func Get(key string, fallback ...string) string {
	val := os.Getenv(key)

	if val == "" && len(fallback) > 0 {
		return fallback[0]
	}

	return val
}

func GetEnvInt(key string, fallback ...int) int {
	val := os.Getenv(key)

	if val == "" && len(fallback) > 0 {
		return fallback[0]
	}

	parsed, err := strconv.Atoi(val)
	if err != nil && len(fallback) > 0 {
		return fallback[0]
	}

	return parsed
}
