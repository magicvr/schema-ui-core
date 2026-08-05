package config

import (
	"strings"
	"testing"
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
