package config

import (
	"log/slog"
	"os"
	"strings"
	"time"
)

// Config is the minimal R1 runtime configuration (no auth / DB).
type Config struct {
	AppName      string
	AppEnv       string
	HTTPAddr     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	LogLevelName string
}

// Load reads configuration from the environment with safe local defaults.
func Load() *Config {
	return &Config{
		AppName:      envOr("APP_NAME", "schema-ui-core-api"),
		AppEnv:       envOr("APP_ENV", "development"),
		HTTPAddr:     envOr("HTTP_ADDR", ":8080"),
		ReadTimeout:  durationEnv("HTTP_READ_TIMEOUT", 5*time.Second),
		WriteTimeout: durationEnv("HTTP_WRITE_TIMEOUT", 10*time.Second),
		IdleTimeout:  durationEnv("HTTP_IDLE_TIMEOUT", 60*time.Second),
		LogLevelName: envOr("LOG_LEVEL", "info"),
	}
}

// LogLevel maps LOG_LEVEL to slog.
func (c *Config) LogLevel() slog.Level {
	switch strings.ToLower(c.LogLevelName) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func envOr(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

func durationEnv(key string, fallback time.Duration) time.Duration {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return fallback
	}
	return d
}
