package composition

import (
	"strings"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/internal/cache"
	"github.com/magicvr/schema-ui-core/apps/api/internal/config"
)

// VP-026 / workspace-026 GOAL-003 D-001: lock the cache wiring — ONE
// kernel.Cache (in-memory provider) built from cache.max_entries, with
// fail-closed rejection of non-positive budgets on hand-built Configs.
func TestNewCacheWiring(t *testing.T) {
	t.Run("max_entries propagates to the memory provider", func(t *testing.T) {
		cachePort, err := newCache(&config.Config{CacheMaxEntries: 5000})
		if err != nil || cachePort == nil {
			t.Fatalf("newCache: cache-nil=%t err=%v", cachePort == nil, err)
		}
		mem, ok := cachePort.(*cache.Memory)
		if !ok {
			t.Fatalf("newCache must produce *cache.Memory, got %T", cachePort)
		}
		if mem.MaxEntries() != 5000 {
			t.Fatalf("maxEntries = %d, want 5000", mem.MaxEntries())
		}
	})

	t.Run("zero-value config falls back to the load default", func(t *testing.T) {
		// Mirrors the db/objects convention: loader-bypassed zero-value
		// Configs (test harnesses, fx) mean "use defaults", not an error.
		cachePort, err := newCache(&config.Config{})
		if err != nil || cachePort == nil {
			t.Fatalf("newCache: cache-nil=%t err=%v", cachePort == nil, err)
		}
		mem, ok := cachePort.(*cache.Memory)
		if !ok {
			t.Fatalf("newCache must produce *cache.Memory, got %T", cachePort)
		}
		if mem.MaxEntries() != config.DefaultCacheMaxEntries {
			t.Fatalf("maxEntries = %d, want default %d", mem.MaxEntries(), config.DefaultCacheMaxEntries)
		}
	})

	t.Run("negative budget fails closed (defense-in-depth)", func(t *testing.T) {
		cachePort, err := newCache(&config.Config{CacheMaxEntries: -1})
		if err == nil || cachePort != nil {
			t.Fatalf("cache-nil=%t err=%v; want nil/nil/error", cachePort != nil, err)
		}
		if !strings.Contains(err.Error(), "cache.max_entries") {
			t.Fatalf("error must name the config key: %v", err)
		}
	})
}
