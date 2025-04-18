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
	JwtSecret stringEnv = "JWT_SECRET"
	// Time since jwt generation that will cause the jwt to be refreshed
	JwtRefreshAge intEnv = "JWT_REFRESH_AGE"
	JwtTTl        intEnv = "JWT_TTL"
)

const (
	// Web server stuffs
	WebDomain     stringEnv = "WEB_ROOT"
	CorsAllowHost stringEnv = "CORS_ALLOW_HOST"
)

const (
	// Config related stuffs
	ConfigPath stringEnv = "CONFIG_PATH"
)

const (
	// Debug related stuffs
	DebugMode boolEnv = "DEBUG_MODE"
)

type EnVar[T any] interface {
	Get(...T) T
}

type stringEnv string

// Get the value of the stringEnv
func (k stringEnv) Get(fallback ...string) string {
	val := os.Getenv(string(k))

	if val == "" && len(fallback) > 0 {
		return fallback[0]
	}

	return val
}

var _ EnVar[string] = (*stringEnv)(nil)

type intEnv string

// Get the value of the intEnv
func (k intEnv) Get(fallback ...int) int {
	val := os.Getenv(string(k))

	if val == "" && len(fallback) > 0 {
		return fallback[0]
	}

	parsed, err := strconv.Atoi(val)
	if err != nil && len(fallback) > 0 {
		return fallback[0]
	}

	return parsed
}

var _ EnVar[int] = (*intEnv)(nil)

type boolEnv string

// Get the value of the boolEnv
func (k boolEnv) Get(fallback ...bool) bool {
	val := os.Getenv(string(k))

	if val == "" && len(fallback) > 0 {
		return fallback[0]
	}

	parsed, err := strconv.ParseBool(val)
	if err != nil && len(fallback) > 0 {
		return fallback[0]
	}

	return parsed
}

var _ EnVar[bool] = (*boolEnv)(nil)
