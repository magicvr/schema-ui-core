package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
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

	ProfileName       string
	ModulesEnabled    []string
	ProfileSource     string
	ProfilePrecedence []string
	ProfileError      error
}

// Load reads configuration from the environment with safe local defaults.
func Load() *Config {
	cfg := &Config{
		AppName:      envOr("APP_NAME", "schema-ui-core-api"),
		AppEnv:       envOr("APP_ENV", ""),
		HTTPAddr:     envOr("HTTP_ADDR", ":25080"),
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
		ProfileName:           envOr("APP_PROFILE", string(kernel.ProfileMVP)),
	}

	explicitModules, err := kernel.ParseModuleList(os.Getenv("APP_MODULES_ENABLED"))
	if err != nil {
		cfg.ProfileError = err
		return cfg
	}
	resolved, err := kernel.ResolveProfile(cfg.ProfileName, explicitModules)
	if err != nil {
		cfg.ProfileError = err
		return cfg
	}
	cfg.ProfileName = string(resolved.Name)
	cfg.ModulesEnabled = append([]string(nil), resolved.Modules...)
	cfg.ProfileSource = resolved.Source
	cfg.ProfilePrecedence = append([]string(nil), resolved.Precedence...)
	return cfg
}

// ValidateProd fails startup for non-development environments when the static
// development session fallback is enabled (GOAL-008 A-005 F-002). The dev
// session substitutes a high-privilege static identity for unauthenticated
// requests, so it may only ever run in local development; any other APP_ENV
// with AUTH_DEV_SESSION_ENABLED=true is a hard startup error, not a warning.
//
// A-002 F-002-005 (GOAL-009 S5): non-development environments additionally
// require a JWT signing secret with a minimum length and both letters and
// digits, so a short or guessable HS256 key cannot silently start production.
// Development keeps the explicit low bar (documented insecure dev key).
func (c *Config) ValidateProd() error {
	if c.ProfileError != nil {
		return fmt.Errorf("invalid module profile: %w", c.ProfileError)
	}
	if c.AppEnv == "" {
		return fmt.Errorf("APP_ENV must be set explicitly (development for local runs, production for deployments); refusing to guess")
	}
	if c.AppEnv == "development" {
		return nil
	}
	if c.AuthDevSessionEnabled {
		return fmt.Errorf("AUTH_DEV_SESSION_ENABLED must be false when APP_ENV=%q", c.AppEnv)
	}
	if len(c.AuthJWTSecret) < minJWTSecretLen {
		return fmt.Errorf(
			"AUTH_JWT_SECRET must be at least %d characters when APP_ENV=%q",
			minJWTSecretLen, c.AppEnv,
		)
	}
	if !containsLettersAndDigits(c.AuthJWTSecret) {
		return fmt.Errorf(
			"AUTH_JWT_SECRET must contain both letters and digits when APP_ENV=%q",
			c.AppEnv,
		)
	}
	return nil
}

// minJWTSecretLen is the minimum HS256 signing-key length enforced outside
// development (A-002 F-002-005).
const minJWTSecretLen = 32

// containsLettersAndDigits rejects all-digit / all-letter / single-class keys
// (weak entropy) while keeping the rule simple and checkable.
func containsLettersAndDigits(s string) bool {
	hasLetter := false
	hasDigit := false
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			hasLetter = true
		}
		if r >= '0' && r <= '9' {
			hasDigit = true
		}
		if hasLetter && hasDigit {
			return true
		}
	}
	return false
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
