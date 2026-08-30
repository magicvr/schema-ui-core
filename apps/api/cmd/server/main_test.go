package main

import (
	"io"
	"log/slog"
	"strings"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/internal/config"
)

// W15 F-003 (GOAL-016 A-001): the production bootstrap seed must satisfy the
// frozen 8–72 byte password policy before it is hashed; development keeps the
// documented "admin" fallback.
func TestResolveSeedHashPolicy(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	t.Run("development keeps the admin fallback", func(t *testing.T) {
		cfg := &config.Config{AppEnv: "development"}
		hash, err := resolveSeedHash(cfg, logger)
		if err != nil {
			t.Fatalf("dev fallback should hash: %v", err)
		}
		if hash == "" {
			t.Fatal("dev fallback must produce a hash")
		}
	})

	t.Run("production missing seed fails closed", func(t *testing.T) {
		cfg := &config.Config{AppEnv: "production"}
		if _, err := resolveSeedHash(cfg, logger); err == nil {
			t.Fatal("production without ADMIN_INITIAL_PASSWORD must fail")
		}
	})

	t.Run("production weak seed fails closed with policy error", func(t *testing.T) {
		for _, weak := range []string{"weak", "1234567", strings.Repeat("x", 73)} {
			cfg := &config.Config{AppEnv: "production", AdminInitialPassword: weak}
			_, err := resolveSeedHash(cfg, logger)
			if err == nil {
				t.Errorf("production seed %q must fail the 8–72 policy", weak)
			} else if !strings.Contains(err.Error(), "password policy") {
				t.Errorf("production seed %q error = %v, want policy wording", weak, err)
			}
		}
	})

	t.Run("production compliant seed passes", func(t *testing.T) {
		cfg := &config.Config{AppEnv: "production", AdminInitialPassword: "seed-long-enough-1"}
		if _, err := resolveSeedHash(cfg, logger); err != nil {
			t.Fatalf("compliant production seed should hash: %v", err)
		}
	})
}