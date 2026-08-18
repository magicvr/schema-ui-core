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
		cfg := Load()
		if cfg.ProfileError != nil {
			t.Fatal(cfg.ProfileError)
		}
		if cfg.ProfileName != "mvp" || cfg.ProfileSource != "profile.default" || len(cfg.ModulesEnabled) == 0 {
			t.Fatalf("unexpected profile config: %+v", cfg)
		}
	})

	t.Run("app.modules.list overrides profile defaults (T-06)", func(t *testing.T) {
		writeConfig(t, `app:
  profile: mvp
  modules:
    list: [core.server-registration, admin.users]
`)
		cfg := Load()
		if cfg.ProfileError != nil {
			t.Fatal(cfg.ProfileError)
		}
		if cfg.ProfileSource != "modules.list" {
			t.Fatalf("ProfileSource = %q, want modules.list", cfg.ProfileSource)
		}
		if len(cfg.ModulesEnabled) != 2 || cfg.ModulesEnabled[0] != "core.server-registration" || cfg.ModulesEnabled[1] != "admin.users" {
			t.Fatalf("ModulesEnabled = %v", cfg.ModulesEnabled)
		}
		if cfg.ProfileName != "custom" {
			t.Fatalf("ProfileName = %q, want custom", cfg.ProfileName)
		}
	})

	t.Run("app.modules.preset builtin name (T-06)", func(t *testing.T) {
		writeConfig(t, `app:
  profile: mvp
  modules:
    preset: admin
`)
		cfg := Load()
		if cfg.ProfileError != nil {
			t.Fatal(cfg.ProfileError)
		}
		if cfg.ProfileSource != "modules.preset" || cfg.ProfileName != "admin" {
			t.Fatalf("preset resolution = %s/%s, want modules.preset/admin", cfg.ProfileSource, cfg.ProfileName)
		}
		if len(cfg.ModulesEnabled) == 0 {
			t.Fatal("preset admin must enable modules")
		}
	})

	t.Run("app.modules.preset custom file (T-06)", func(t *testing.T) {
		dir := t.TempDir()
		preset := filepath.Join(dir, "custom.yaml")
		if err := os.WriteFile(preset, []byte(`modules:
  - core.server-registration
  - admin.roles
`), 0o644); err != nil {
			t.Fatal(err)
		}
		writeConfig(t, `app:
  profile: mvp
  modules:
    preset: `+preset+`
`)
		cfg := Load()
		if cfg.ProfileError != nil {
			t.Fatal(cfg.ProfileError)
		}
		if cfg.ProfileSource != "modules.preset" || len(cfg.ModulesEnabled) != 2 {
			t.Fatalf("preset file resolution: source=%s modules=%v", cfg.ProfileSource, cfg.ModulesEnabled)
		}
	})

	t.Run("app.modules preset and list are mutually exclusive (T-06)", func(t *testing.T) {
		writeConfig(t, `app:
  profile: mvp
  modules:
    preset: admin
    list: [core.server-registration]
`)
		cfg := Load()
		if cfg.ProfileError == nil || !strings.Contains(cfg.ProfileError.Error(), "mutually exclusive") {
			t.Fatalf("both preset+list must fail closed, got %v", cfg.ProfileError)
		}
	})

	t.Run("custom profile without app.modules fails closed (T-06)", func(t *testing.T) {
		writeConfig(t, `app:
  profile: custom
`)
		cfg := Load()
		if cfg.ProfileError == nil {
			t.Fatal("custom profile without app.modules must fail closed")
		}
		if !strings.Contains(cfg.ProfileError.Error(), "app.modules") {
			t.Fatalf("custom profile error = %q, want app.modules mention", cfg.ProfileError)
		}
	})

	t.Run("legacy env module selectors are ignored (T-06)", func(t *testing.T) {
		t.Setenv("APP_PROFILE", "demo")
		t.Setenv("APP_MODULES_ENABLED", "core.server-registration")
		cfg := Load()
		if cfg.ProfileError != nil {
			t.Fatal(cfg.ProfileError)
		}
		if cfg.ProfileName != "mvp" {
			t.Fatalf("APP_PROFILE must be ignored (YAML-only), got %q", cfg.ProfileName)
		}
	})
}

func TestLoadRuntimeModePrecedenceAndFailClosed(t *testing.T) {
	t.Run("default is normal", func(t *testing.T) {
		cfg := Load()
		if cfg.LoadError != nil || cfg.RuntimeMode != RuntimeModeNormal {
			t.Fatalf("runtime mode = %q, load error = %v", cfg.RuntimeMode, cfg.LoadError)
		}
	})
	t.Run("yaml value", func(t *testing.T) {
		writeConfig(t, "runtime:\n  mode: read-only\n")
		cfg := Load()
		if cfg.LoadError != nil || cfg.RuntimeMode != RuntimeModeReadOnly {
			t.Fatalf("runtime mode = %q, load error = %v", cfg.RuntimeMode, cfg.LoadError)
		}
	})
	t.Run("environment overrides yaml", func(t *testing.T) {
		writeConfig(t, "runtime:\n  mode: maintenance\n")
		t.Setenv("RUNTIME_MODE", "degraded")
		cfg := Load()
		if cfg.LoadError != nil || cfg.RuntimeMode != RuntimeModeDegraded {
			t.Fatalf("runtime mode = %q, load error = %v", cfg.RuntimeMode, cfg.LoadError)
		}
	})
	for _, tc := range []struct {
		name  string
		yaml  string
		env   string
	}{
		{name: "empty yaml", yaml: "runtime:\n  mode: \"\"\n"},
		{name: "unknown yaml", yaml: "runtime:\n  mode: paused\n"},
		{name: "empty env", yaml: "", env: ""},
		{name: "unknown env", yaml: "", env: "paused"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.yaml != "" {
				writeConfig(t, tc.yaml)
			}
			t.Setenv("RUNTIME_MODE", tc.env)
			cfg := Load()
			if cfg.LoadError == nil {
				t.Fatalf("invalid runtime mode %q must set LoadError", tc.env)
			}
		})
	}
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

// TestBrandingConfig covers the W9 brand-asset processing policy: defaults,
// YAML layer and env overrides.
func TestBrandingConfig(t *testing.T) {
	t.Run("defaults", func(t *testing.T) {
		cfg := Load()
		if cfg.LoadError != nil {
			t.Fatalf("LoadError: %v", cfg.LoadError)
		}
		if cfg.BrandingMaxBytes != 4<<20 {
			t.Errorf("BrandingMaxBytes = %d, want %d", cfg.BrandingMaxBytes, 4<<20)
		}
		if cfg.BrandingLogoMaxDimension != 512 || cfg.BrandingFaviconDimension != 64 {
			t.Errorf("dimensions = %d/%d, want 512/64", cfg.BrandingLogoMaxDimension, cfg.BrandingFaviconDimension)
		}
		if cfg.BrandingJPEGQuality != 82 {
			t.Errorf("BrandingJPEGQuality = %d, want 82", cfg.BrandingJPEGQuality)
		}
	})

	t.Run("yaml layer", func(t *testing.T) {
		y := "app:\n  env: development\nbranding:\n  max_bytes: 2097152\n  logo_max_dimension: 256\n  favicon_dimension: 32\n  jpeg_quality: 75\n"
		writeConfig(t, y)
		cfg := Load()
		if cfg.LoadError != nil {
			t.Fatalf("LoadError: %v", cfg.LoadError)
		}
		if cfg.BrandingMaxBytes != 2097152 || cfg.BrandingLogoMaxDimension != 256 ||
			cfg.BrandingFaviconDimension != 32 || cfg.BrandingJPEGQuality != 75 {
			t.Errorf("yaml branding = %+v", cfg)
		}
	})

	t.Run("out-of-range jpeg quality falls back to default", func(t *testing.T) {
		y := "app:\n  env: development\nbranding:\n  jpeg_quality: 120\n"
		writeConfig(t, y)
		cfg := Load()
		if cfg.LoadError != nil {
			t.Fatalf("LoadError: %v", cfg.LoadError)
		}
		if cfg.BrandingJPEGQuality != 82 {
			t.Errorf("BrandingJPEGQuality = %d, want default 82", cfg.BrandingJPEGQuality)
		}
	})

	t.Run("env overrides", func(t *testing.T) {
		t.Setenv("BRANDING_MAX_BYTES", "1048576")
		t.Setenv("BRANDING_LOGO_MAX_DIMENSION", "128")
		t.Setenv("BRANDING_FAVICON_DIMENSION", "16")
		t.Setenv("BRANDING_JPEG_QUALITY", "60")
		cfg := Load()
		if cfg.LoadError != nil {
			t.Fatalf("LoadError: %v", cfg.LoadError)
		}
		if cfg.BrandingMaxBytes != 1048576 || cfg.BrandingLogoMaxDimension != 128 ||
			cfg.BrandingFaviconDimension != 16 || cfg.BrandingJPEGQuality != 60 {
			t.Errorf("env branding = %+v", cfg)
		}
	})
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

// GOAL-013 D-002 §4: navigation.order parsing (YAML list + NAVIGATION_ORDER
// env override + malformed fallback).
func TestLoadNavigationOrder(t *testing.T) {
	t.Run("yaml order list is parsed", func(t *testing.T) {
		y := `app:
  env: development
navigation:
  order:
    - menu_roles
    - menu_dashboard
`
		writeConfig(t, y)
		cfg := Load()
		if cfg.LoadError != nil {
			t.Fatalf("LoadError: %v", cfg.LoadError)
		}
		want := []string{"menu_roles", "menu_dashboard"}
		if !strings.EqualFold(strings.Join(cfg.NavigationOrder, ","), strings.Join(want, ",")) {
			t.Fatalf("NavigationOrder = %v, want %v", cfg.NavigationOrder, want)
		}
	})

	t.Run("empty order yields nil default", func(t *testing.T) {
		y := `app:
  env: development
navigation:
  order: []
`
		writeConfig(t, y)
		cfg := Load()
		if cfg.LoadError != nil {
			t.Fatalf("LoadError: %v", cfg.LoadError)
		}
		if len(cfg.NavigationOrder) != 0 {
			t.Fatalf("NavigationOrder = %v, want empty (default)", cfg.NavigationOrder)
		}
	})

	t.Run("env NAVIGATION_ORDER overrides yaml", func(t *testing.T) {
		y := `app:
  env: development
navigation:
  order:
    - menu_roles
`
		writeConfig(t, y)
		t.Setenv("NAVIGATION_ORDER", "menu_account, menu_users")
		cfg := Load()
		if cfg.LoadError != nil {
			t.Fatalf("LoadError: %v", cfg.LoadError)
		}
		want := []string{"menu_account", "menu_users"}
		if !strings.EqualFold(strings.Join(cfg.NavigationOrder, ","), strings.Join(want, ",")) {
			t.Fatalf("NavigationOrder = %v, want %v", cfg.NavigationOrder, want)
		}
	})

	t.Run("malformed navigation.order is not fatal (fallback to default)", func(t *testing.T) {
		y := `app:
  env: development
navigation:
  order: not-a-list
`
		writeConfig(t, y)
		cfg := Load()
		if cfg.LoadError != nil {
			t.Fatalf("LoadError: %v, want graceful fallback", cfg.LoadError)
		}
		if len(cfg.NavigationOrder) != 0 {
			t.Fatalf("NavigationOrder = %v, want empty after fallback", cfg.NavigationOrder)
		}
	})

	t.Run("non-string entry falls back to default", func(t *testing.T) {
		y := `app:
  env: development
navigation:
  order:
    - menu_roles
    - 42
`
		writeConfig(t, y)
		cfg := Load()
		if cfg.LoadError != nil {
			t.Fatalf("LoadError: %v, want graceful fallback", cfg.LoadError)
		}
		if len(cfg.NavigationOrder) != 0 {
			t.Fatalf("NavigationOrder = %v, want empty after fallback", cfg.NavigationOrder)
		}
	})
}

// A-003 findings regression tests:
//  F-002: omitted YAML keys keep code defaults (no zeroing).
//  F-003: inline " #" inside a quoted value is not a comment.
//  F-005: empty/comment-only files mean all defaults; multi-document YAML is rejected.
func TestLoadA003Findings(t *testing.T) {
	t.Run("F-002 omitted keys keep defaults", func(t *testing.T) {
		y := `app:
  env: development
`
		writeConfig(t, y)
		cfg := Load()
		if cfg.LoadError != nil {
			t.Fatalf("LoadError: %v", cfg.LoadError)
		}
		if cfg.HTTPAddr != ":25080" {
			t.Errorf("HTTPAddr = %q, want default :25080 (omitted key)", cfg.HTTPAddr)
		}
		if cfg.DBPath != "./data/schema-ui.db" {
			t.Errorf("DBPath = %q, want default ./data/schema-ui.db (omitted key)", cfg.DBPath)
		}
		if cfg.ReadTimeout != 5*time.Second {
			t.Errorf("ReadTimeout = %v, want default 5s (omitted key)", cfg.ReadTimeout)
		}
		if cfg.UploadMaxFilesPerUser != 1000 {
			t.Errorf("UploadMaxFilesPerUser = %d, want default 1000", cfg.UploadMaxFilesPerUser)
		}
		if cfg.LogLevelName != "info" {
			t.Errorf("LogLevelName = %q, want default info", cfg.LogLevelName)
		}
	})

	t.Run("F-003 quoted hash is not a comment", func(t *testing.T) {
		y := `app:
  env: development
  name: "My App #1"
`
		writeConfig(t, y)
		cfg := Load()
		if cfg.LoadError != nil {
			t.Fatalf("LoadError: %v", cfg.LoadError)
		}
		if cfg.AppName != "My App #1" {
			t.Errorf("AppName = %q, want My App #1 (quoted hash must survive)", cfg.AppName)
		}
	})

	t.Run("F-005 empty file means defaults", func(t *testing.T) {
		p := filepath.Join(t.TempDir(), "empty.yaml")
		if err := os.WriteFile(p, []byte("# only a comment\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		t.Setenv("CONFIG_FILE", p)
		cfg := Load()
		if cfg.LoadError != nil {
			t.Fatalf("LoadError: %v, want empty file to fall back to defaults", cfg.LoadError)
		}
		if cfg.HTTPAddr != ":25080" {
			t.Errorf("HTTPAddr = %q, want default :25080", cfg.HTTPAddr)
		}
		if cfg.AppEnv != "" {
			t.Errorf("AppEnv = %q, want empty default", cfg.AppEnv)
		}
	})

	t.Run("F-005 multi-document YAML rejected", func(t *testing.T) {
		y := `app:
  env: development
---
app:
  bogus_key: true
`
		writeConfig(t, y)
		cfg := Load()
		if cfg.LoadError == nil {
			t.Fatal("multi-document YAML must be a LoadError (second doc escapes KnownFields)")
		}
		if !strings.Contains(cfg.LoadError.Error(), "multiple YAML documents") {
			t.Fatalf("LoadError = %v, want mention of multiple documents", cfg.LoadError)
		}
	})
}

