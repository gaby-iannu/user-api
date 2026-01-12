package config

import (
	"os"
	"time"
)

type Config struct {
	Port         string
	DatabaseURL  string
	LogLevel     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func Load() *Config {
	return &Config{
		Port:         getEnv("PORT", "8080"),
		DatabaseURL:  getEnv("DATABASE_URL", ""),
		LogLevel:     getEnv("LOG_LEVEL", "info"),
		ReadTimeout:  getDuration("READ_TIMEOUT", 5*time.Second),
		WriteTimeout: getDuration("WRITE_TIMEOUT", 10*time.Second),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
	}
	return defaultValue
}
