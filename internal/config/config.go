package config

import (
	"os"
	"time"
)

type Config struct {
	AppName         string
	AppEnv          string
	HTTPPort        string
	ShutdownTimeout time.Duration
}

func Load() Config {
	return Config{
		AppName:         getEnv("APP_NAME", "hcm-go"),
		AppEnv:          getEnv("APP_ENV", "development"),
		HTTPPort:        getEnv("HTTP_PORT", "8080"),
		ShutdownTimeout: getDurationEnv("SHUTDOWN_TIMEOUT", 10*time.Second),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

func getDurationEnv(key string, fallback time.Duration) time.Duration {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}

	duration, err := time.ParseDuration(raw)
	if err != nil {
		return fallback
	}

	return duration
}
