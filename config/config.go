package config

import (
	"fmt"
	"strconv"
	"sync"
)

var lock = &sync.Mutex{}

type AuthConfig struct {
	TOKEN_KEY          string
	REFRESH_TOKEN_KEY  string
	TOKEN_DURATION_SEC int
}

var authConfigInstance *AuthConfig

func GetAuthConfigInstance() *AuthConfig {
	if authConfigInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if authConfigInstance == nil {
			fmt.Println("Creating AuthConfig instance now.")
			tokenDuration, err := strconv.Atoi(Env("TOKEN_DURATION_SEC"))
			if err != nil {
				fmt.Println("Error parsing TOKEN_DURATION")
			}

			authConfigInstance = &AuthConfig{
				TOKEN_KEY:          Env("TOKEN_KEY"),
				REFRESH_TOKEN_KEY:  Env("REFRESH_TOKEN_KEY"),
				TOKEN_DURATION_SEC: tokenDuration,
			}
		} else {
			fmt.Println("AuthConfig instance already created.")
		}
	} else {
		fmt.Println("AuthConfig instance already created.")
	}

	return authConfigInstance
}
