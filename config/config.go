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
	REFRESH_TOKEN_SEC  int
}

var authConfigInstance *AuthConfig

func GetAuthConfigInstance() *AuthConfig {
	if authConfigInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if authConfigInstance == nil {
			tokenDuration, err := strconv.Atoi(Env("TOKEN_DURATION_SEC"))
			if err != nil {
				fmt.Println("Error parsing TOKEN_DURATION")
			}

			refreshDuration, err := strconv.Atoi(Env("REFRESH_TOKEN_SEC"))
			if err != nil {
				fmt.Println("Error parsing REFRESH_TOKEN_SEC")
			}

			authConfigInstance = &AuthConfig{
				TOKEN_KEY:          Env("TOKEN_KEY"),
				REFRESH_TOKEN_KEY:  Env("REFRESH_TOKEN_KEY"),
				TOKEN_DURATION_SEC: tokenDuration,
				REFRESH_TOKEN_SEC:  refreshDuration, // 1 week
			}
		}
	}

	return authConfigInstance
}
