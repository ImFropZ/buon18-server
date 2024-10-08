package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

var lock = &sync.Mutex{}

type Config struct {
	PORT int

	DB_CONNECTION_STRING string
	TOKEN_KEY            string
	REFRESH_TOKEN_KEY    string
	TOKEN_DURATION_SEC   int
	REFRESH_TOKEN_SEC    int

	LOGGING_DIR string

	// -- Valkey
	VALKEY_ADDRESSES   []string
	VALKEY_PWD         string
	CACHE_DURATION_SEC int

	// -- Trusted Proxies
	TRUSTED_PROXIES []string

	// -- CORS
	ACCESS_CONTROL_ALLOW_ORIGIN      string
	ACCESS_CONTROL_ALLOW_CREDENTIALS string
	ACCESS_CONTROL_ALLOW_HEADERS     string
	ACCESS_CONTROL_ALLOW_METHODS     string
	ACCESS_CONTROL_EXPOSE_HEADERS    string
	ACCESS_CONTROL_MAX_AGE           string

	// -- TLS
	CERT_FILE string
	KEY_FILE  string
}

var configInstance *Config

func GetConfigInstance() *Config {
	if configInstance == nil {
		lock.Lock()
		defer lock.Unlock()

		// -- Port
		port, err := strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			fmt.Println("Error parsing PORT")
			port = 80
		}

		// -- Cache Duration
		cacheDuration, err := strconv.Atoi(Env("CACHE_DURATION_SEC"))
		if err != nil {
			fmt.Println("Error parsing CACHE_DURATION_SEC")
		}

		// -- Token Duration
		tokenDuration, err := strconv.Atoi(Env("TOKEN_DURATION_SEC"))
		if err != nil {
			fmt.Println("Error parsing TOKEN_DURATION")
		}

		refreshDuration, err := strconv.Atoi(Env("REFRESH_TOKEN_SEC"))
		if err != nil {
			fmt.Println("Error parsing REFRESH_TOKEN_SEC")
		}

		// -- Trusted Proxies
		proxies := Env("TRUSTED_PROXIES")
		trustedProxies := []string{}
		if proxies != "" {
			trustedProxies = append(trustedProxies, proxies)
		}

		configInstance = &Config{
			PORT: port,

			// -- Logging
			LOGGING_DIR: Env("LOGGING_DIR"),

			// -- Database
			DB_CONNECTION_STRING: validateEnvString("DB_CONNECTION_STRING"),

			// -- Valkey
			VALKEY_ADDRESSES:   strings.Split(Env("VALKEY_ADDRESSES"), ","),
			VALKEY_PWD:         Env("VALKEY_PWD"),
			CACHE_DURATION_SEC: cacheDuration,

			// -- Auth
			TOKEN_KEY:          validateEnvString("TOKEN_KEY"),
			REFRESH_TOKEN_KEY:  validateEnvString("REFRESH_TOKEN_KEY"),
			TOKEN_DURATION_SEC: tokenDuration,
			REFRESH_TOKEN_SEC:  refreshDuration,

			// -- Trusted Proxies
			TRUSTED_PROXIES: trustedProxies,

			// -- CORS
			ACCESS_CONTROL_ALLOW_ORIGIN:      Env("ACCESS_CONTROL_ALLOW_ORIGIN"),
			ACCESS_CONTROL_ALLOW_CREDENTIALS: Env("ACCESS_CONTROL_ALLOW_CREDENTIALS"),
			ACCESS_CONTROL_ALLOW_HEADERS:     Env("ACCESS_CONTROL_ALLOW_HEADERS"),
			ACCESS_CONTROL_ALLOW_METHODS:     Env("ACCESS_CONTROL_ALLOW_METHODS"),
			ACCESS_CONTROL_EXPOSE_HEADERS:    Env("ACCESS_CONTROL_EXPOSE_HEADERS"),
			ACCESS_CONTROL_MAX_AGE:           Env("ACCESS_CONTROL_MAX_AGE"),

			// -- TLS
			CERT_FILE: Env("CERT_FILE"),
			KEY_FILE:  Env("KEY_FILE"),
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
