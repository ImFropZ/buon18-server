package config

import (
	"fmt"
	"strconv"
	"sync"
)

var lock = &sync.Mutex{}

type Config struct {
	DB_CONNECTION_STRING string
	TOKEN_KEY            string
	REFRESH_TOKEN_KEY    string
	TOKEN_DURATION_SEC   int
	REFRESH_TOKEN_SEC    int
}

var configInstance *Config

func GetConfigInstance() *Config {
	if configInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if configInstance == nil {
			tokenDuration, err := strconv.Atoi(Env("TOKEN_DURATION_SEC"))
			if err != nil {
				fmt.Println("Error parsing TOKEN_DURATION")
			}

			refreshDuration, err := strconv.Atoi(Env("REFRESH_TOKEN_SEC"))
			if err != nil {
				fmt.Println("Error parsing REFRESH_TOKEN_SEC")
			}

			configInstance = &Config{
				// -- Database
				DB_CONNECTION_STRING: validateEnvString("DB_CONNECTION_STRING"),

				// -- Auth
				TOKEN_KEY:          validateEnvString("TOKEN_KEY"),
				REFRESH_TOKEN_KEY:  validateEnvString("REFRESH_TOKEN_KEY"),
				TOKEN_DURATION_SEC: tokenDuration,
				REFRESH_TOKEN_SEC:  refreshDuration, // 1 week
			}
		}
	}

	return configInstance
}

func validateEnvString(key string) (value string) {
	value = Env(key)
	if value == "" {
		panic(fmt.Sprintf("%s enviroment variable is required", key))
	}
	return
}
