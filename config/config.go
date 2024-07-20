package config

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

var lock = &sync.Mutex{}

type Config struct {
	DB_CONNECTION_STRING string
	TOKEN_KEY            string
	REFRESH_TOKEN_KEY    string
	TOKEN_DURATION_SEC   int
	REFRESH_TOKEN_SEC    int

	// -- Trusted Proxies
	TRUSTED_PROXIES []string

	// -- CORS
	ALLOW_ORIGINS  []string
	ALLOW_METHODS  []string
	ALLOW_HEADERS  []string
	EXPOSE_HEADERS []string
	MAX_AGE        time.Duration

	// -- TLS
	CERT_FILE string
	KEY_FILE  string
}

var configInstance *Config

func GetConfigInstance() *Config {
	if configInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if configInstance == nil {
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

			// -- CORS
			aOrigins := Env("ALLOW_ORIGINS")
			allowOrigins := []string{}
			if aOrigins == "" {
				allowOrigins = append(allowOrigins, "*")
			} else {
				allowOrigins = append(allowOrigins, strings.Split(aOrigins, ",")...)
			}

			aMethods := Env("ALLOW_METHODS")
			allowMethods := []string{}
			if aMethods == "" {
				allowMethods = append(allowMethods, "GET", "POST", "PATCH", "DELETE", "OPTIONS")
			} else {
				allowMethods = append(allowMethods, strings.Split(aMethods, ",")...)
			}

			aHeaders := Env("ALLOW_HEADERS")
			allowHeaders := []string{}
			if aHeaders == "" {
				allowHeaders = append(allowHeaders, "*")
			} else {
				allowHeaders = append(allowHeaders, strings.Split(aHeaders, ",")...)
			}

			eHeaders := Env("EXPOSE_HEADERS")
			exposeHeaders := []string{}
			if eHeaders == "" {
				exposeHeaders = append(exposeHeaders, "Content-Length")
			} else {
				exposeHeaders = append(exposeHeaders, strings.Split(eHeaders, ",")...)
			}

			mAge := Env("MAX_AGE")
			maxAge := 120
			if mAge != "" {
				maxAge, err = strconv.Atoi(mAge)
				if err != nil {
					fmt.Println("Error parsing MAX_AGE")
				}
			}

			configInstance = &Config{
				// -- Database
				DB_CONNECTION_STRING: validateEnvString("DB_CONNECTION_STRING"),

				// -- Auth
				TOKEN_KEY:          validateEnvString("TOKEN_KEY"),
				REFRESH_TOKEN_KEY:  validateEnvString("REFRESH_TOKEN_KEY"),
				TOKEN_DURATION_SEC: tokenDuration,
				REFRESH_TOKEN_SEC:  refreshDuration,

				// -- Trusted Proxies
				TRUSTED_PROXIES: trustedProxies,

				// -- CORS
				ALLOW_ORIGINS:  allowOrigins,
				ALLOW_METHODS:  allowMethods,
				ALLOW_HEADERS:  allowHeaders,
				EXPOSE_HEADERS: exposeHeaders,
				MAX_AGE:        time.Duration(maxAge) * time.Second,

				// -- TLS
				CERT_FILE: Env("CERT_FILE"),
				KEY_FILE:  Env("KEY_FILE"),
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
