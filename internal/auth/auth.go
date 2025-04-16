package auth

import (
	"errors"
	"os/exec"
	"strings"

	"github.com/indeedhat/parity-nas/internal/config"
	"github.com/openwall/yescrypt-go"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUnknownAlgo     = errors.New("Unknown hashing algorithm")
	ErrInvalidPassword = errors.New("Invalid Password")
)

// attemptSystemLogin attempts to login using system credentials
func attemptSystemLogin(user, pass string) *config.AuthCfgUser {
	cfg, err := config.Auth()
	if err != nil {
		return nil
	}

	userEntry := cfg.Users.Find(user)
	if userEntry == nil {
		return nil
	}

	cmd := exec.Command("getent", "shadow", user)
	out, err := cmd.Output()

	if err != nil {
		return nil
	}

	parts := strings.Split(string(out), ":")
	if len(parts) < 2 {
		return nil
	}

	if err := verifyPassword(parts[1], pass); err != nil {
		return nil
	}

	return userEntry
}

// verifyPassword extracts the hashing algorithm and runs the appropriate verification check
func verifyPassword(hash, pass string) error {
	parts := strings.Split(hash, "$")
	if len(parts) < 2 {
		return nil
	}

	switch parts[1] {
	// yescrypt
	case "y":
		hashed, err := yescrypt.Hash([]byte(pass), []byte(hash))
		if err != nil {
			return err
		}
		if hash != string(hashed) {
			return ErrInvalidPassword
		}

	// bcrypt
	case "1", "2", "sha1", "5", "6", "2a", "2x", "2y", "2b":
		return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))

	default:
		return ErrUnknownAlgo
	}

	return nil
}
