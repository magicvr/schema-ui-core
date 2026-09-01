package cache

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

var (
	_ kernel.Cache        = (*Memory)(nil)
	_ kernel.CacheView    = (*memoryView)(nil)
	_ kernel.ExpiryPolicy = AbsoluteExpiry{}
	_ kernel.ExpiryPolicy = SlidingExpiry{}
	_ kernel.ExpiryPolicy = nextMidnightPolicy{}
)

var testCtx = context.Background()

// stepClock is a manual test clock; advance mutates its instant.
type stepClock struct{ t time.Time }

func (c *stepClock) now() time.Time          { return c.t }
func (c *stepClock) advance(d time.Duration) { c.t = c.t.Add(d) }

func newTestMemory(t *testing.T, maxEntries int) (*Memory, *stepClock) {
	t.Helper()
	clock := &stepClock{t: time.Date(2026, 9, 1, 9, 0, 0, 0, time.UTC)}
	m, err := newMemoryWithClock(maxEntries, clock.now)
	if err != nil {
		t.Fatalf("newMemoryWithClock: %v", err)
	}
	return m, clock
}

func mustView(t *testing.T, m *Memory, ns string) kernel.CacheView {
	t.Helper()
	v, err := m.Namespace(kernel.CacheNamespace(ns))
	if err != nil {
		t.Fatalf("Namespace(%q): %v", ns, err)
	}
	return v
}

func TestNewMemoryValidation(t *testing.T) {
	if _, err := NewMemory(0); err == nil {
		t.Fatal("NewMemory(0) must fail closed")
	}
	if _, err := NewMemory(-1); err == nil {
		t.Fatal("NewMemory(-1) must fail closed")
	}
	if _, err := NewMemory(10); err != nil {
		t.Fatalf("NewMemory(10): %v", err)
	}
}

func TestMemoryNamespaceValidation(t *testing.T) {
	m, _ := newTestMemory(t, 10)
	for _, ns := range []kernel.CacheNamespace{"", "Wallet", "-wallet", "wallet-", "wallet--x", "钱包"} {
		if v, err := m.Namespace(ns); err == nil || v != nil {
			t.Errorf("Namespace(%q) must fail closed, got view=%v err=%v", ns, v, err)
		} else if !errors.Is(err, kernel.ErrInvalidCacheNamespace) {
			t.Errorf("Namespace(%q) error = %v, want ErrInvalidCacheNamespace", ns, err)
		}
	}
}

func TestMemorySetGetDeleteBasics(t *testing.T) {
	m, _ := newTestMemory(t, 10)
	v := mustView(t, m, "wallet")

	if got, ok := v.Get(testCtx, "k"); ok || got != nil {
		t.Fatalf("miss before Set: got=%v ok=%v", got, ok)
	}
	if err := v.Set(testCtx, "k", []byte("v1"), AbsoluteExpiry{TTL: 0}); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if got, ok := v.Get(testCtx, "k"); !ok || string(got) != "v1" {
		t.Fatalf("hit: got=%q ok=%v", got, ok)
	}
	// Empty value hits with (empty, true) — distinct from miss.
	if err := v.Set(testCtx, "empty", []byte{}, AbsoluteExpiry{TTL: 0}); err != nil {
		t.Fatalf("Set empty: %v", err)
	}
	if got, ok := v.Get(testCtx, "empty"); !ok || got == nil || len(got) != 0 {
		t.Fatalf("empty-value hit: got=%v ok=%v", got, ok)
	}
	// Delete then miss; delete of absent key is idempotent.
	if err := v.Delete(testCtx, "k"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if got, ok := v.Get(testCtx, "k"); ok || got != nil {
		t.Fatalf("hit after Delete: got=%v ok=%v", got, ok)
	}
	if err := v.Delete(testCtx, "k"); err != nil {
		t.Fatalf("idempotent Delete: %v", err)
	}
}

func TestMemoryFailClosedValidation(t *testing.T) {
	m, _ := newTestMemory(t, 10)
	v := mustView(t, m, "wallet")

	if err := v.Set(testCtx, "", []byte("v"), AbsoluteExpiry{}); !errors.Is(err, kernel.ErrInvalidCacheKey) {
		t.Errorf("Set invalid key: %v", err)
	}
	if err := v.Set(testCtx, "k", nil, AbsoluteExpiry{}); !errors.Is(err, kernel.ErrInvalidCacheValue) {
		t.Errorf("Set nil value: %v", err)
	}
	if err := v.Set(testCtx, "k", []byte("v"), nil); !errors.Is(err, kernel.ErrInvalidCachePolicy) {
		t.Errorf("Set nil policy: %v", err)
	}
	// ValidateCacheSet order: key beats value beats policy.
	if err := v.Set(testCtx, "", nil, nil); !errors.Is(err, kernel.ErrInvalidCacheKey) {
		t.Errorf("order key-first: %v", err)
	}
	if err := v.Delete(testCtx, ""); !errors.Is(err, kernel.ErrInvalidCacheKey) {
		t.Errorf("Delete invalid key: %v", err)
	}
	// Get has no error channel: invalid key is a miss.
	if got, ok := v.Get(testCtx, ""); got != nil || ok {
		t.Errorf("Get invalid key must be a miss, got %v %v", got, ok)
	}
}

func TestMemoryCopySemantics(t *testing.T) {
	m, _ := newTestMemory(t, 10)
	v := mustView(t, m, "wallet")

	buf := []byte("original")
	if err := v.Set(testCtx, "k", buf, AbsoluteExpiry{TTL: 0}); err != nil {
		t.Fatal(err)
	}
	buf[0] = 'X' // caller mutates its buffer after Set
	if got, _ := v.Get(testCtx, "k"); string(got) != "original" {
		t.Fatalf("Set must copy input; got %q", got)
	}
	got, _ := v.Get(testCtx, "k")
	got[0] = 'Y' // caller mutates the returned slice
	if again, _ := v.Get(testCtx, "k"); string(again) != "original" {
		t.Fatalf("Get must return a fresh copy; got %q", again)
	}
}

func TestMemoryNamespaceIsolation(t *testing.T) {
	m, _ := newTestMemory(t, 10)
	a := mustView(t, m, "wallet")
	b := mustView(t, m, "session2fa")
	if err := a.Set(testCtx, "k", []byte("a"), AbsoluteExpiry{TTL: 0}); err != nil {
		t.Fatal(err)
	}
	if got, ok := b.Get(testCtx, "k"); ok || got != nil {
		t.Fatalf("namespace isolation broken: b sees %q ok=%v", got, ok)
	}
	if got, ok := a.Get(testCtx, "k"); !ok || string(got) != "a" {
		t.Fatalf("a lost its value: %q ok=%v", got, ok)
	}
}

func TestMemoryAbsoluteExpiry(t *testing.T) {
	m, clock := newTestMemory(t, 10)
	v := mustView(t, m, "wallet")
	if err := v.Set(testCtx, "k", []byte("v"), AbsoluteExpiry{TTL: 5 * time.Second}); err != nil {
		t.Fatal(err)
	}
	clock.advance(4 * time.Second)
	if _, ok := v.Get(testCtx, "k"); !ok {
		t.Fatal("hit within TTL expected")
	}
	clock.advance(2 * time.Second) // now t+6: expired, and the hit did NOT refresh
	if got, ok := v.Get(testCtx, "k"); ok || got != nil {
		t.Fatal("absolute expiry must not refresh on hits")
	}
}

func TestMemorySlidingExpiry(t *testing.T) {
	m, clock := newTestMemory(t, 10)
	v := mustView(t, m, "wallet")
	if err := v.Set(testCtx, "k", []byte("v"), SlidingExpiry{Window: 5 * time.Second}); err != nil {
		t.Fatal(err)
	}
	clock.advance(4 * time.Second)
	if _, ok := v.Get(testCtx, "k"); !ok {
		t.Fatal("hit within window expected")
	}
	clock.advance(4 * time.Second) // now t+8: the hit at t+4 refreshed to t+9
	if _, ok := v.Get(testCtx, "k"); !ok {
		t.Fatal("sliding expiry must refresh on hits")
	}
	// The hit at t+8 refreshed the expiry to t+13; t+14 must be a miss.
	clock.advance(6 * time.Second) // now t+14
	if got, ok := v.Get(testCtx, "k"); ok || got != nil {
		t.Fatal("entry must expire after refreshed window passes")
	}
}

func TestMemoryZeroTTLNeverExpires(t *testing.T) {
	m, clock := newTestMemory(t, 10)
	v := mustView(t, m, "wallet")
	if err := v.Set(testCtx, "k", []byte("v"), AbsoluteExpiry{TTL: 0}); err != nil {
		t.Fatal(err)
	}
	clock.advance(24 * time.Hour)
	if _, ok := v.Get(testCtx, "k"); !ok {
		t.Fatal("zero TTL must never expire")
	}
	// Sliding with a non-positive Window must also never expire (and never
	// refresh into an expiry): F-003 coverage.
	if err := v.Set(testCtx, "s", []byte("s"), SlidingExpiry{Window: 0}); err != nil {
		t.Fatal(err)
	}
	clock.advance(48 * time.Hour)
	if _, ok := v.Get(testCtx, "s"); !ok {
		t.Fatal("zero sliding window must never expire")
	}
}

// nextMidnightPolicy is the CUSTOM pluggable-strategy sample (VP-026 判据 #2:
// strategies are injected via the interface, not hardcoded). It is stateless:
// entries expire at the next local midnight, absolute-style.
type nextMidnightPolicy struct{}

func (nextMidnightPolicy) ExpireAt(now time.Time) time.Time {
	return time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
}

func (nextMidnightPolicy) Refresh(_ time.Time, previous time.Time) (time.Time, bool) {
	return previous, false
}

func TestMemoryCustomPolicyPluggable(t *testing.T) {
	m, clock := newTestMemory(t, 10)
	v := mustView(t, m, "wallet")
	if err := v.Set(testCtx, "k", []byte("v"), nextMidnightPolicy{}); err != nil {
		t.Fatal(err)
	}
	clock.advance(14 * time.Hour) // 23:00 same day: still before midnight
	if _, ok := v.Get(testCtx, "k"); !ok {
		t.Fatal("hit before midnight expected")
	}
	clock.advance(2 * time.Hour) // 01:00 next day
	if got, ok := v.Get(testCtx, "k"); ok || got != nil {
		t.Fatal("entry must expire after midnight")
	}
}

func TestMemoryFIFOEvictionBound(t *testing.T) {
	m, _ := newTestMemory(t, 3)
	v := mustView(t, m, "wallet")
	// Fill to capacity and overflow: oldest insertion is evicted each time.
	for _, k := range []string{"k1", "k2", "k3", "k4"} {
		if err := v.Set(testCtx, k, []byte(k), AbsoluteExpiry{TTL: 0}); err != nil {
			t.Fatalf("Set %s: %v", k, err)
		}
	}
	if got, ok := v.Get(testCtx, "k1"); ok || got != nil {
		t.Fatal("k1 must be FIFO-evicted")
	}
	for _, k := range []string{"k2", "k3", "k4"} {
		if _, ok := v.Get(testCtx, k); !ok {
			t.Fatalf("%s must survive", k)
		}
	}
}

func TestMemoryEvictionOverwriteKeepsPosition(t *testing.T) {
	m, _ := newTestMemory(t, 2)
	v := mustView(t, m, "wallet")
	if err := v.Set(testCtx, "a", []byte("a1"), AbsoluteExpiry{TTL: 0}); err != nil {
		t.Fatal(err)
	}
	if err := v.Set(testCtx, "b", []byte("b1"), AbsoluteExpiry{TTL: 0}); err != nil {
		t.Fatal(err)
	}
	// Overwrite a: FIFO position is preserved (no move-to-back).
	if err := v.Set(testCtx, "a", []byte("a2"), AbsoluteExpiry{TTL: 0}); err != nil {
		t.Fatal(err)
	}
	if err := v.Set(testCtx, "c", []byte("c1"), AbsoluteExpiry{TTL: 0}); err != nil {
		t.Fatal(err)
	}
	if got, ok := v.Get(testCtx, "a"); ok || got != nil {
		t.Fatal("FIFO must evict the OLDEST INSERTED key even after overwrite")
	}
	if _, ok := v.Get(testCtx, "b"); !ok {
		t.Fatal("b must survive")
	}
	if _, ok := v.Get(testCtx, "c"); !ok {
		t.Fatal("c must survive")
	}
}

func TestMemoryLiveOverwriteKeepsPositionAcrossPolicyChange(t *testing.T) {
	// F-004 disposition (D-001 勘误): a LIVE overwrite keeps its FIFO
	// position even when the policy instance changes — only a dead (expired)
	// entry re-inserts. No interface comparison is performed (uncomparable
	// policy types can never panic the provider).
	m, _ := newTestMemory(t, 2)
	v := mustView(t, m, "wallet")
	if err := v.Set(testCtx, "a", []byte("a"), AbsoluteExpiry{TTL: 0}); err != nil {
		t.Fatal(err)
	}
	if err := v.Set(testCtx, "b", []byte("b"), AbsoluteExpiry{TTL: 0}); err != nil {
		t.Fatal(err)
	}
	// Live overwrite of a with a DIFFERENT policy (hour-long absolute):
	// position must be preserved.
	if err := v.Set(testCtx, "a", []byte("a2"), AbsoluteExpiry{TTL: time.Hour}); err != nil {
		t.Fatal(err)
	}
	if err := v.Set(testCtx, "c", []byte("c"), AbsoluteExpiry{TTL: 0}); err != nil {
		t.Fatal(err)
	}
	if got, ok := v.Get(testCtx, "a"); ok || got != nil {
		t.Fatal("a must still be the oldest (policy change must not move it to the back)")
	}
	if _, ok := v.Get(testCtx, "b"); !ok {
		t.Fatal("b must survive")
	}
	if _, ok := v.Get(testCtx, "c"); !ok {
		t.Fatal("c must survive")
	}
}

func TestMemoryGlobalBudgetAcrossNamespaces(t *testing.T) {
	// F-001 adjudication: the maxEntries budget is PROCESS-WIDE. Two
	// namespaces share one budget; the globally-oldest entry is evicted.
	m, _ := newTestMemory(t, 3)
	a := mustView(t, m, "wallet")
	b := mustView(t, m, "session2fa")
	for _, k := range []string{"a1", "a2", "a3"} {
		if err := a.Set(testCtx, k, []byte(k), AbsoluteExpiry{TTL: 0}); err != nil {
			t.Fatalf("a.Set %s: %v", k, err)
		}
	}
	if err := b.Set(testCtx, "b1", []byte("b1"), AbsoluteExpiry{TTL: 0}); err != nil {
		t.Fatal(err)
	}
	// a1 was the globally-oldest: evicted even though namespace a is at 2.
	if got, ok := a.Get(testCtx, "a1"); ok || got != nil {
		t.Fatal("a1 must be evicted by the process-wide budget")
	}
	if _, ok := a.Get(testCtx, "a2"); !ok {
		t.Fatal("a2 must survive")
	}
	if _, ok := a.Get(testCtx, "a3"); !ok {
		t.Fatal("a3 must survive")
	}
	if _, ok := b.Get(testCtx, "b1"); !ok {
		t.Fatal("b1 must survive")
	}
	if m.Len() != 3 {
		t.Fatalf("process-wide total = %d, want 3", m.Len())
	}
}

func TestMemoryGlobalFIFOInterleave(t *testing.T) {
	// Global insertion order drives eviction across namespaces: a1, a2, b1
	// then a3 evicts a1 (globally oldest), not b1.
	m, _ := newTestMemory(t, 3)
	a := mustView(t, m, "wallet")
	b := mustView(t, m, "session2fa")
	for _, k := range []string{"a1", "a2"} {
		if err := a.Set(testCtx, k, []byte(k), AbsoluteExpiry{TTL: 0}); err != nil {
			t.Fatal(err)
		}
	}
	if err := b.Set(testCtx, "b1", []byte("b1"), AbsoluteExpiry{TTL: 0}); err != nil {
		t.Fatal(err)
	}
	if err := a.Set(testCtx, "a3", []byte("a3"), AbsoluteExpiry{TTL: 0}); err != nil {
		t.Fatal(err)
	}
	if got, ok := a.Get(testCtx, "a1"); ok || got != nil {
		t.Fatal("a1 must be evicted (global FIFO)")
	}
	if _, ok := b.Get(testCtx, "b1"); !ok {
		t.Fatal("b1 must survive (inserted before a3)")
	}
}

func TestMemoryLazyCleanupFreesCapacity(t *testing.T) {
	m, clock := newTestMemory(t, 3)
	v := mustView(t, m, "wallet")
	if err := v.Set(testCtx, "a", []byte("a"), AbsoluteExpiry{TTL: time.Second}); err != nil {
		t.Fatal(err)
	}
	if err := v.Set(testCtx, "b", []byte("b"), AbsoluteExpiry{TTL: 0}); err != nil {
		t.Fatal(err)
	}
	if err := v.Set(testCtx, "c", []byte("c"), AbsoluteExpiry{TTL: time.Second}); err != nil {
		t.Fatal(err)
	}
	clock.advance(2 * time.Second)
	// Lazy removal happens on the read path: a and c are swept now.
	if got, ok := v.Get(testCtx, "a"); ok || got != nil {
		t.Fatal("expired a must miss and be removed")
	}
	if got, ok := v.Get(testCtx, "c"); ok || got != nil {
		t.Fatal("expired c must miss and be removed")
	}
	// Only b (never-expiring) remains: d must be insertable WITHOUT eviction
	// of b — proving the lazy sweeps freed capacity.
	if err := v.Set(testCtx, "d", []byte("d"), AbsoluteExpiry{TTL: 0}); err != nil {
		t.Fatal(err)
	}
	if _, ok := v.Get(testCtx, "b"); !ok {
		t.Fatal("b must survive (lazy cleanup freed the capacity, not eviction)")
	}
	if _, ok := v.Get(testCtx, "d"); !ok {
		t.Fatal("d must be present")
	}
}

func TestMemoryExpiredEntriesStillBoundTheTotal(t *testing.T) {
	m, clock := newTestMemory(t, 3)
	v := mustView(t, m, "wallet")
	for _, k := range []string{"a", "b", "c"} {
		if err := v.Set(testCtx, k, []byte(k), AbsoluteExpiry{TTL: time.Second}); err != nil {
			t.Fatal(err)
		}
	}
	clock.advance(2 * time.Second) // all expired but NOT swept (no touch)
	for _, k := range []string{"d", "e", "f", "g"} {
		if err := v.Set(testCtx, k, []byte(k), AbsoluteExpiry{TTL: 0}); err != nil {
			t.Fatalf("Set %s: %v", k, err)
		}
	}
	// FIFO evictions must keep total <= 3 even with expired entries present.
	// After 4 sets (d..g) from a full 3-entry space, a..c and d are gone.
	for _, k := range []string{"a", "b", "c", "d"} {
		if got, ok := v.Get(testCtx, k); ok || got != nil {
			t.Fatalf("%s must have been evicted (bound holds over expired entries)", k)
		}
	}
	for _, k := range []string{"e", "f", "g"} {
		if _, ok := v.Get(testCtx, k); !ok {
			t.Fatalf("%s must survive", k)
		}
	}
}

func TestMemoryConcurrentAccess(t *testing.T) {
	m, _ := newTestMemory(t, 64)
	v := mustView(t, m, "wallet")
	other := mustView(t, m, "session2fa")
	var wg sync.WaitGroup
	for g := 0; g < 8; g++ {
		wg.Add(1)
		go func(seed int) {
			defer wg.Done()
			for i := 0; i < 200; i++ {
				key := string(rune('a' + (seed+i)%4))
				switch (seed + i) % 4 {
				case 0:
					_ = v.Set(testCtx, key, []byte("v"), AbsoluteExpiry{TTL: time.Minute})
				case 1:
					_, _ = v.Get(testCtx, key)
				case 2:
					_ = v.Delete(testCtx, key)
				default:
					_, _ = other.Get(testCtx, key)
				}
			}
		}(g)
	}
	wg.Wait()
	// Post-race sanity: the store must still answer coherently.
	if _, err := m.Namespace("wallet"); err != nil {
		t.Fatalf("Namespace after race: %v", err)
	}
}

func TestMemoryConcurrentBudgetBound(t *testing.T) {
	// F-003: after concurrent unique-key Sets, the process-wide total must
	// still respect maxEntries (the bound survives the race).
	m, _ := newTestMemory(t, 25)
	v := mustView(t, m, "wallet")
	var wg sync.WaitGroup
	for g := 0; g < 8; g++ {
		wg.Add(1)
		go func(seed int) {
			defer wg.Done()
			for i := 0; i < 50; i++ {
				key := fmt.Sprintf("g%d-%d", seed, i)
				_ = v.Set(testCtx, key, []byte(key), AbsoluteExpiry{TTL: time.Minute})
			}
		}(g)
	}
	wg.Wait()
	if total := m.Len(); total > 25 {
		t.Fatalf("process-wide total = %d after 400 unique Sets, want <= 25", total)
	}
}
