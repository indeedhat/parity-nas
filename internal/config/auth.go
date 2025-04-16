package config

import (
	"os"
)

const AuthKey = "auth"

type AuthCfgUser struct {
	Username   string `icl:".param"`
	Permission uint8  `icl:"permission"`
}

// Read returns if the user has read permissions for the cluster
func (u AuthCfgUser) Read() bool {
	return u.Permission&4 == 4
}

// Write returns if the user has write permissions for the cluster
func (u AuthCfgUser) Write() bool {
	return u.Permission&2 == 2
}

// Admin returns if the user has admin permissions for the cluster
func (u AuthCfgUser) Admin() bool {
	return u.Permission&1 == 1
}

type AuthCfgUsers []AuthCfgUser

// Find attempts to find a user by their name
func (u AuthCfgUsers) Find(name string) *AuthCfgUser {
	for _, user := range u {
		if name == user.Username {
			return &user
		}
	}

	return nil
}

type AuthCfg struct {
	Version uint `icl:"version"`

	Users AuthCfgUsers `icl:"user"`
}

// Auth initializes a AuthCfg struct
//
// If a config file exists then it will attempt to load it otherwise a new instance will be created
func Auth() (*AuthCfg, error) {
	var c AuthCfg

	if err := loadConfig(AuthKey, &c); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}

		c = AuthCfg{
			Version: 1,
			Users: []AuthCfgUser{
				{
					Username:   "root",
					Permission: 7,
				},
			},
		}
	}

	return &c, nil
}
