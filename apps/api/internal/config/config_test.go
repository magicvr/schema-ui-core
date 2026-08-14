package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLoadResolvesProfileAndModuleOverrides(t *testing.T) {
	t.Run("default mvp profile", func(t *testing.T) {
		t.Setenv("APP_PROFILE", "")
		t.Setenv("APP_MODULES_ENABLED", "")
		cfg := Load()
		if cfg.ProfileError != nil {
			t.Fatal(cfg.ProfileError)
		}
		if cfg.ProfileName != "mvp" || cfg.ProfileSource != "profile.default" || len(cfg.ModulesEnabled) == 0 {
			t.Fatalf("unexpected profile config: %+v", cfg)
		}
	})

	t.Run("explicit modules override profile defaults", func(t *testing.T) {
		t.Setenv("APP_PROFILE", "admin")
		t.Setenv("APP_MODULES_ENABLED", "core.server-registration")
		cfg := Load()
		if cfg.ProfileError != nil {
			t.Fatal(cfg.ProfileError)
		}
		if cfg.ProfileSource != "modules.enabled" || len(cfg.ModulesEnabled) != 1 || cfg.ModulesEnabled[0] != "core.server-registration" {
			t.Fatalf("unexpected explicit module config: %+v", cfg)
		}
	})

	t.Run("custom profile requires explicit modules", func(t *testing.T) {
		t.Setenv("APP_PROFILE", "custom")
		t.Setenv("APP_MODULES_ENABLED", "")
		cfg := Load()
		if cfg.ProfileError == nil {
			t.Fatal("custom profile without modules must fail closed")
		}
		if !strings.Contains(cfg.ProfileError.Error(), "APP_MODULES_ENABLED") {
			t.Fatalf("custom profile error = %q, want actual environment key", cfg.ProfileError)
		}
	})
}

// TestValidateProd covers the production guard added in response to GOAL-008
// A-005 F-002 (dev-session fallback) and A-002 F-002-005 (JWT secret minimum
// length/entropy): both are local-development-only or non-negotiable settings
// that must fail startup in any non-development environment.
func TestValidateProd(t *testing.T) {
	t.Run("unset APP_ENV fails closed (C3)", func(t *testing.T) {
		c := &Config{AppEnv: ""}
		if err := c.ValidateProd(); err == nil {
			t.Fatal("unset APP_ENV must be a startup error, not a silent development fallback")
		}
	})

	t.Run("development may enable dev session", func(t *testing.T) {
		c := &Config{AppEnv: "development", AuthDevSessionEnabled: true}
		if err := c.ValidateProd(); err != nil {
			t.Fatalf("development + dev session should pass, got: %v", err)
		}
	})

	t.Run("production with dev session fails closed", func(t *testing.T) {
		c := &Config{AppEnv: "production", AuthDevSessionEnabled: true}
		if err := c.ValidateProd(); err == nil {
			t.Fatal("production + dev session must be a startup error")
		}
	})

	t.Run("production without dev session passes", func(t *testing.T) {
		c := &Config{AppEnv: "production", AuthDevSessionEnabled: false, AuthJWTSecret: strongSecret}
		if err := c.ValidateProd(); err != nil {
			t.Fatalf("production without dev session should pass, got: %v", err)
		}
	})

	t.Run("non-development non-production env also fails closed", func(t *testing.T) {
		c := &Config{AppEnv: "staging", AuthDevSessionEnabled: true}
		if err := c.ValidateProd(); err == nil {
			t.Fatal("staging + dev session must be a startup error")
		}
	})

	t.Run("production with a short JWT secret fails closed", func(t *testing.T) {
		c := &Config{AppEnv: "production", AuthDevSessionEnabled: false, AuthJWTSecret: "short-secret"}
		if err := c.ValidateProd(); err == nil {
			t.Fatal("production with a short AUTH_JWT_SECRET must be a startup error")
		}
	})

	t.Run("production with an all-letter JWT secret fails closed", func(t *testing.T) {
		c := &Config{
			AppEnv:                "production",
			AuthDevSessionEnabled: false,
			AuthJWTSecret:         "abcdefghijklmnopqrstuvwxyzabcdefghij",
		}
		if err := c.ValidateProd(); err == nil {
			t.Fatal("production with an all-letter AUTH_JWT_SECRET must be a startup error")
		}
	})

	t.Run("production with an all-digit JWT secret fails closed", func(t *testing.T) {
		c := &Config{
			AppEnv:                "production",
			AuthDevSessionEnabled: false,
			AuthJWTSecret:         "12345678901234567890123456789012",
		}
		if err := c.ValidateProd(); err == nil {
			t.Fatal("production with an all-digit AUTH_JWT_SECRET must be a startup error")
		}
	})

	t.Run("development keeps the low JWT secret bar", func(t *testing.T) {
		c := &Config{AppEnv: "development", AuthJWTSecret: "dev"}
		if err := c.ValidateProd(); err != nil {
			t.Fatalf("development should keep the low bar, got: %v", err)
		}
	})
}

// strongSecret satisfies the production AUTH_JWT_SECRET rule (≥32 chars, mixed).
const strongSecret = "a9k2m4n6p8q0r2s4t6u8v0w2x4y6z8a9b1c3d5"

// W7 (GOAL-008): YAML authority with env override, interpolation, and
// fail-closed behavior. These tests write a temp YAML and point CONFIG_FILE
// at it; process env is cleared/restored per test via t.Setenv.

func writeConfig(t *testing.T, yamlText string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(p, []byte(yamlText), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("CONFIG_FILE", p)
	return p
}

func TestLoadYAMLLayer(t *testing.T) {
	t.Run("CONFIG_FILE with plain values and env override", func(t *testing.T) {
		y := `app:
  name: yaml-name
  env: development
  profile: mvp
http:
  addr: ":9999"
  read_timeout: 3s
  write_timeout: 4s
  idle_timeout: 9s
log:
  level: debug
auth:
  jwt_secret: yaml-secret
  access_ttl: 10m
  refresh_ttl: 48h
  dev_session_enabled: true
db:
  path: /tmp/yaml.db
admin:
  initial_password: yaml-admin
upload:
  allowed_types: "image/png,text/plain"
  max_files_per_user: 42
  max_bytes_per_user: 1048576
`
		writeConfig(t, y)
		// env overrides YAML when set
		t.Setenv("HTTP_ADDR", ":4321")
		t.Setenv("UPLOAD_MAX_FILES_PER_USER", "7")
		t.Setenv("AUTH_DEV_SESSION_ENABLED", "false")
		cfg := Load()
		if cfg.LoadError != nil {
			t.Fatalf("LoadError: %v", cfg.LoadError)
		}
		if cfg.AppName != "yaml-name" {
			t.Errorf("AppName = %q, want yaml-name", cfg.AppName)
		}
		if cfg.AppEnv != "development" {
			t.Errorf("AppEnv = %q, want development", cfg.AppEnv)
		}
		if cfg.HTTPAddr != ":4321" {
			t.Errorf("HTTPAddr = %q, want env override :4321", cfg.HTTPAddr)
		}
		if cfg.ReadTimeout != 3*time.Second {
			t.Errorf("ReadTimeout = %v, want 3s", cfg.ReadTimeout)
		}
		if cfg.LogLevelName != "debug" {
			t.Errorf("LogLevelName = %q, want debug", cfg.LogLevelName)
		}
		if cfg.AuthJWTSecret != "yaml-secret" {
			t.Errorf("AuthJWTSecret = %q, want yaml-secret", cfg.AuthJWTSecret)
		}
		if cfg.AuthAccessTTL != 10*time.Minute {
			t.Errorf("AuthAccessTTL = %v, want 10m", cfg.AuthAccessTTL)
		}
		if cfg.AuthDevSessionEnabled {
			t.Error("AuthDevSessionEnabled = true, want env override false")
		}
		if cfg.DBPath != "/tmp/yaml.db" {
			t.Errorf("DBPath = %q, want /tmp/yaml.db", cfg.DBPath)
		}
		if cfg.UploadAllowedTypes != "image/png,text/plain" {
			t.Errorf("UploadAllowedTypes = %q", cfg.UploadAllowedTypes)
		}
		if cfg.UploadMaxFilesPerUser != 7 {
			t.Errorf("UploadMaxFilesPerUser = %d, want env override 7", cfg.UploadMaxFilesPerUser)
		}
		if cfg.UploadMaxBytesPerUser != 1048576 {
			t.Errorf("UploadMaxBytesPerUser = %d, want 1048576", cfg.UploadMaxBytesPerUser)
		}
		if cfg.ProfileError != nil {
			t.Fatalf("ProfileError: %v", cfg.ProfileError)
		}
		if cfg.ProfileName != "mvp" {
			t.Errorf("ProfileName = %q, want mvp", cfg.ProfileName)
		}
	})

	t.Run("explicit CONFIG_FILE missing fails closed", func(t *testing.T) {
		t.Setenv("CONFIG_FILE", filepath.Join(t.TempDir(), "nope.yaml"))
		cfg := Load()
		if cfg.LoadError == nil {
			t.Fatal("missing explicit CONFIG_FILE must be a LoadError (fail-closed)")
		}
	})

	t.Run("bare env reference without value fails closed", func(t *testing.T) {
		y := `app:
  env: development
auth:
  jwt_secret: ${MISSING_JWT_W7}
`
		writeConfig(t, y)
		os.Unsetenv("MISSING_JWT_W7")
		cfg := Load()
		if cfg.LoadError == nil {
			t.Fatal("bare ${MISSING_JWT_W7} must fail closed")
		}
		if !strings.Contains(cfg.LoadError.Error(), "MISSING_JWT_W7") {
			t.Fatalf("LoadError = %v, want mention of MISSING_JWT_W7", cfg.LoadError)
		}
	})

	t.Run("default value applies when env unset", func(t *testing.T) {
		y := `app:
  env: development
auth:
  jwt_secret: ${W7_DEFAULTED:-fallback-secret}
`
		writeConfig(t, y)
		os.Unsetenv("W7_DEFAULTED")
		cfg := Load()
		if cfg.LoadError != nil {
			t.Fatalf("LoadError: %v", cfg.LoadError)
		}
		if cfg.AuthJWTSecret != "fallback-secret" {
			t.Errorf("AuthJWTSecret = %q, want fallback-secret", cfg.AuthJWTSecret)
		}
	})

	t.Run("interpolation uses env when set", func(t *testing.T) {
		y := `app:
  env: development
auth:
  jwt_secret: ${W7_INTERP:-fallback}
`
		writeConfig(t, y)
		t.Setenv("W7_INTERP", "env-wins")
		cfg := Load()
		if cfg.LoadError != nil {
			t.Fatalf("LoadError: %v", cfg.LoadError)
		}
		if cfg.AuthJWTSecret != "env-wins" {
			t.Errorf("AuthJWTSecret = %q, want env-wins", cfg.AuthJWTSecret)
		}
	})

	t.Run("unknown YAML keys fail closed", func(t *testing.T) {
		y := `app:
  env: development
  bogus_key: true
`
		writeConfig(t, y)
		cfg := Load()
		if cfg.LoadError == nil {
			t.Fatal("unknown YAML key must be a LoadError (typos fail loudly)")
		}
	})

	t.Run("CONFIG_ENV_FILE supplies secrets without overriding process env", func(t *testing.T) {
		y := `app:
  env: development
auth:
  jwt_secret: ${W7_ENVFILE_SECRET}
`
		writeConfig(t, y)
		envFile := filepath.Join(t.TempDir(), "secrets.env")
		if err := os.WriteFile(envFile, []byte("W7_ENVFILE_SECRET=from-file\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		t.Setenv("CONFIG_ENV_FILE", envFile)
		os.Unsetenv("W7_ENVFILE_SECRET")
		cfg := Load()
		if cfg.LoadError != nil {
			t.Fatalf("LoadError: %v", cfg.LoadError)
		}
		if cfg.AuthJWTSecret != "from-file" {
			t.Errorf("AuthJWTSecret = %q, want from-file", cfg.AuthJWTSecret)
		}
	})

	t.Run("CONFIG_ENV_FILE never overrides an existing process env", func(t *testing.T) {
		y := `app:
  env: development
auth:
  jwt_secret: ${W7_ENVFILE_SECRET2}
`
		writeConfig(t, y)
		envFile := filepath.Join(t.TempDir(), "secrets.env")
		if err := os.WriteFile(envFile, []byte("W7_ENVFILE_SECRET2=from-file\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		t.Setenv("CONFIG_ENV_FILE", envFile)
		t.Setenv("W7_ENVFILE_SECRET2", "process-wins")
		cfg := Load()
		if cfg.LoadError != nil {
			t.Fatalf("LoadError: %v", cfg.LoadError)
		}
		if cfg.AuthJWTSecret != "process-wins" {
			t.Errorf("AuthJWTSecret = %q, want process-wins (process env beats env file)", cfg.AuthJWTSecret)
		}
	})

	t.Run("explicit CONFIG_ENV_FILE missing fails closed", func(t *testing.T) {
		t.Setenv("CONFIG_ENV_FILE", filepath.Join(t.TempDir(), "nope.env"))
		cfg := Load()
		if cfg.LoadError == nil {
			t.Fatal("missing explicit CONFIG_ENV_FILE must be a LoadError (fail-closed)")
		}
	})

	t.Run("ValidateProd surfaces LoadError", func(t *testing.T) {
		t.Setenv("CONFIG_FILE", filepath.Join(t.TempDir(), "nope.yaml"))
		cfg := Load()
		if cfg.LoadError == nil {
			t.Fatal("missing explicit CONFIG_FILE must be a LoadError")
		}
		if err := cfg.ValidateProd(); err == nil {
			t.Fatal("ValidateProd must fail when LoadError is set")
		}
	})
}
