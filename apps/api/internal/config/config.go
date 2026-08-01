package config

import (
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config is the R2 runtime configuration: HTTP + logging, plus the auth
// (JWT / refresh / SQLite) and dev-session surface defined by GOAL-005 D-004.
type Config struct {
	AppName      string
	AppEnv       string
	HTTPAddr     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	LogLevelName string

	AuthJWTSecret         string
	AuthAccessTTL         time.Duration
	AuthRefreshTTL        time.Duration
	DBPath                string
	AdminInitialPassword  string
	AuthDevSessionEnabled bool
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

		AuthJWTSecret:         envOr("AUTH_JWT_SECRET", ""),
		AuthAccessTTL:         durationEnv("AUTH_ACCESS_TTL", 15*time.Minute),
		AuthRefreshTTL:        durationEnv("AUTH_REFRESH_TTL", 30*24*time.Hour),
		DBPath:                envOr("DB_PATH", "./data/schema-ui.db"),
		AdminInitialPassword:  envOr("ADMIN_INITIAL_PASSWORD", ""),
		AuthDevSessionEnabled: boolEnv("AUTH_DEV_SESSION_ENABLED", false),
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

func boolEnv(key string, fallback bool) bool {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}
	return b
}
