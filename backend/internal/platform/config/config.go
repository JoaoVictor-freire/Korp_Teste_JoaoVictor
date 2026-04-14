package config

import (
	"os"
	"strconv"
	"strings"
)

type HTTPConfig struct {
	Host string
	Port string
}

func NewHTTPConfig(prefix string, defaultPort string) HTTPConfig {
	host := getEnv(prefix+"_HOST", "localhost")
	port := getEnv(prefix+"_PORT", defaultPort)

	return HTTPConfig{
		Host: host,
		Port: port,
	}
}

func (c HTTPConfig) Address() string {
	return c.Host + ":" + c.Port
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func GetEnvAsInt(key string, fallback int) int {
	rawValue := os.Getenv(key)
	if rawValue == "" {
		return fallback
	}

	value, err := strconv.Atoi(rawValue)
	if err != nil {
		return fallback
	}

	return value
}

func GetEnv(key string, fallback string) string {
	return getEnv(key, fallback)
}

func GetEnvTrimmed(key string, fallback string) string {
	return strings.TrimSpace(getEnv(key, fallback))
}
