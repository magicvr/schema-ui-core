package config

import (
	"strings"
	"testing"
)

// R2 surface (workspace-026 GOAL-003 D-001): cache.max_entries — the
// process-wide bounded-entry budget of the in-memory cache provider. Default
// 10000; explicit invalid values fail closed (LoadError), never degrade to
// the default.
func TestCacheMaxEntriesConfig(t *testing.T) {
	t.Run("untouched config carries the default", func(t *testing.T) {
		writeConfig(t, "app:\n  env: development\n")
		cfg := Load()
		if cfg.LoadError != nil {
			t.Fatalf("LoadError: %v", cfg.LoadError)
		}
		if cfg.CacheMaxEntries != DefaultCacheMaxEntries {
			t.Fatalf("CacheMaxEntries = %d, want default %d", cfg.CacheMaxEntries, DefaultCacheMaxEntries)
		}
	})

	t.Run("yaml cache.max_entries parses", func(t *testing.T) {
		writeConfig(t, "app:\n  env: development\ncache:\n  max_entries: 5000\n")
		cfg := Load()
		if cfg.LoadError != nil {
			t.Fatalf("LoadError: %v", cfg.LoadError)
		}
		if cfg.CacheMaxEntries != 5000 {
			t.Fatalf("CacheMaxEntries = %d, want 5000", cfg.CacheMaxEntries)
		}
	})

	t.Run("env overrides yaml", func(t *testing.T) {
		t.Setenv("CACHE_MAX_ENTRIES", "777")
		writeConfig(t, "app:\n  env: development\ncache:\n  max_entries: 5000\n")
		cfg := Load()
		if cfg.LoadError != nil {
			t.Fatalf("LoadError: %v", cfg.LoadError)
		}
		if cfg.CacheMaxEntries != 777 {
			t.Fatalf("CacheMaxEntries = %d, want env 777", cfg.CacheMaxEntries)
		}
	})

	t.Run("invalid env fails closed naming the key", func(t *testing.T) {
		for _, raw := range []string{"not-a-number", "0", "-5"} {
			t.Setenv("CACHE_MAX_ENTRIES", raw)
			writeConfig(t, "app:\n  env: development\n")
			cfg := Load()
			if cfg.LoadError == nil || !strings.Contains(cfg.LoadError.Error(), "CACHE_MAX_ENTRIES") {
				t.Fatalf("CACHE_MAX_ENTRIES=%q must fail closed naming the key, got %v", raw, cfg.LoadError)
			}
		}
	})

	t.Run("invalid yaml value fails closed naming the key", func(t *testing.T) {
		writeConfig(t, "app:\n  env: development\ncache:\n  max_entries: 0\n")
		cfg := Load()
		if cfg.LoadError == nil || !strings.Contains(cfg.LoadError.Error(), "cache.max_entries") {
			t.Fatalf("cache.max_entries=0 must fail closed naming the key, got %v", cfg.LoadError)
		}
	})

	t.Run("negative budget is a startup gate in ValidateProd", func(t *testing.T) {
		c := &Config{AppEnv: "development", CacheMaxEntries: -1}
		if err := c.ValidateProd(); err == nil || !strings.Contains(err.Error(), "cache.max_entries") {
			t.Fatalf("negative CacheMaxEntries must fail ValidateProd, got %v", err)
		}
	})
}
