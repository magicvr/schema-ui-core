package ratelimit

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// limiter returns the provider-built *Memory for direct introspection
// (legacy in-package semantics assertions moved from handler auth_test.go).
func limiter(t *testing.T, window time.Duration, max, capacity int) *Memory {
	t.Helper()
	l, ok := NewProvider().NewRateLimiter(window, max, capacity).(*Memory)
	if !ok {
		t.Fatalf("provider returned %T, want *Memory", l)
	}
	return l
}

// D-002 §1: Allow never registers — this is the legacy W4 P0-1 regression
// (allow() before record(); if allow registered the key, record()'s capacity
// eviction would be dead code and a username spray would grow the map).
func TestMemoryAllowDoesNotRegisterKey(t *testing.T) {
	l := limiter(t, 15*time.Minute, 1, 2)
	now := time.Now().UTC()

	if !l.Allow("10.0.0.1|spray", now) {
		t.Fatal("fresh key must be allowed")
	}
	if len(l.attempts) != 0 {
		t.Fatalf("Allow must not register a key, got %d entries", len(l.attempts))
	}
	if len(l.order) != 0 {
		t.Fatalf("Allow must not touch the eviction order, got %d", len(l.order))
	}

	// Simulate the real login path: Allow then Record for many distinct
	// usernames. Capacity 2 means only the two newest keys survive.
	for _, user := range []string{"a", "b", "c", "d"} {
		key := "10.0.0.1|" + user
		if !l.Allow(key, now) {
			t.Fatalf("fresh key %s must be allowed", user)
		}
		l.Record(key, now)
	}
	if len(l.attempts) != 2 {
		t.Fatalf("sprayed map must stay at capacity 2, got %d entries", len(l.attempts))
	}
	if !l.Allow("10.0.0.1|a", now) {
		t.Fatal("oldest evicted key must be allowed again")
	}
	if !l.Allow("10.0.0.1|b", now) {
		t.Fatal("second-oldest evicted key must be allowed again")
	}
	if l.Allow("10.0.0.1|d", now) {
		t.Fatal("newest key must still hold its failure")
	}
}

// Sliding window semantics: per-key budget, window expiry, clear on success,
// bounded capacity eviction (legacy TestLoginRateLimiterUnit).
func TestMemorySlidingWindowSemantics(t *testing.T) {
	l := limiter(t, 15*time.Minute, 2, 1<<16)
	now := time.Now().UTC()
	l.Record("10.0.0.1|admin", now)
	if !l.Allow("10.0.0.1|admin", now) {
		t.Fatal("first failure under the limit must still allow")
	}
	l.Record("10.0.0.1|admin", now)
	if l.Allow("10.0.0.1|admin", now) {
		t.Fatal("attempt with the window full must be blocked")
	}
	if !l.Allow("10.0.0.2|admin", now) {
		t.Fatal("a different IP must not inherit another IP's failures")
	}
	if !l.Allow("10.0.0.1|other", now) {
		t.Fatal("a different username on the same IP must not inherit the failures")
	}
	if !l.Allow("10.0.0.1|admin", now.Add(16*time.Minute)) {
		t.Fatal("attempt after the window must be allowed again")
	}

	// A successful login clears the failure bucket (D-001 P1).
	l.Record("10.0.0.1|admin", now)
	l.Record("10.0.0.1|admin", now)
	if l.Allow("10.0.0.1|admin", now) {
		t.Fatal("bucket full before clear must block")
	}
	l.Clear("10.0.0.1|admin")
	if !l.Allow("10.0.0.1|admin", now) {
		t.Fatal("after clear the key must be allowed")
	}

	// Bounded map: spraying distinct identities evicts the oldest key (D-001 P1).
	small := limiter(t, 15*time.Minute, 1, 3)
	small.Record("k1", now)
	small.Record("k2", now)
	small.Record("k3", now)
	small.Record("k4", now) // evicts k1 (oldest)
	if !small.Allow("k1", now) {
		t.Fatal("evicted oldest key must be allowed again")
	}
	if small.Allow("k4", now) {
		t.Fatal("newest key must still hold its failure")
	}
}

// Retry-After follows the frozen kernel predicate: seconds remaining after
// the oldest in-window failure; remaining <= 0 maps to 1. RetryAfterSeconds
// does NOT prune (D-002 v0.1.1 §3, A-002 F-006): a key whose failures are all
// beyond the window reports the minimum 1 until Allow prunes it.
func TestMemoryRetryAfter(t *testing.T) {
	l := limiter(t, 15*time.Minute, 2, 1<<16)
	now := time.Now().UTC()
	if got := l.RetryAfterSeconds("k", now); got != 0 {
		t.Fatalf("no failures RetryAfterSeconds = %d, want 0", got)
	}
	l.Record("k", now.Add(-10*time.Minute))
	l.Record("k", now)
	if got, want := l.RetryAfterSeconds("k", now), 5*60; got != want {
		t.Fatalf("RetryAfterSeconds = %d, want %d", got, want)
	}
	l.Record("k2", now.Add(-20*time.Minute)) // beyond window: not pruned by RetryAfter
	if got := l.RetryAfterSeconds("k2", now); got != 1 {
		t.Fatalf("window-expired key RetryAfterSeconds = %d, want 1 (no pruning; remain<=0 → 1)", got)
	}
	// Allow prunes: after an Allow pass the stale entry is gone and the key
	// reports 0 again (D-002 v0.1.1 §3: pruning happens on Allow only).
	_ = l.Allow("k2", now)
	if got := l.RetryAfterSeconds("k2", now); got != 0 {
		t.Fatalf("after Allow pruning RetryAfterSeconds = %d, want 0", got)
	}
}

// D-002 §4: capacity <= 0 falls back to kernel.DefaultRateLimiterCapacity;
// the resulting map is bounded and evicts the oldest key.
func TestProviderDefaultCapacityFallback(t *testing.T) {
	l, ok := NewProvider().NewRateLimiter(15*time.Minute, 1, 0).(*Memory)
	if !ok {
		t.Fatalf("provider returned %T, want *Memory", l)
	}
	if got, want := l.capacity, kernel.DefaultRateLimiterCapacity; got != want {
		t.Fatalf("capacity = %d, want %d", got, want)
	}
	now := time.Now().UTC()
	// Exceed the map bound by one: the very first key must be evicted.
	for i := 0; i <= kernel.DefaultRateLimiterCapacity; i++ {
		l.Record("k"+itoa(i), now)
	}
	if len(l.attempts) != kernel.DefaultRateLimiterCapacity {
		t.Fatalf("map size = %d, want %d", len(l.attempts), kernel.DefaultRateLimiterCapacity)
	}
	if got := l.RetryAfterSeconds("k0", now); got != 0 {
		t.Fatalf("evicted oldest key RetryAfterSeconds = %d, want 0", got)
	}
}

// D-002 §6: all methods are safe for concurrent use (-race).
func TestMemoryConcurrent(t *testing.T) {
	l := NewProvider().NewRateLimiter(time.Minute, 3, 1<<16)
	now := time.Now().UTC()
	var wg sync.WaitGroup
	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range 500 {
				key := "ip|user"
				_ = l.Allow(key, now)
				l.Record(key, now)
				_ = l.AllowRecord(key, now)
				_ = l.RetryAfterSeconds(key, now)
				if j%7 == 0 {
					l.Clear(key)
				}
			}
		}()
	}
	wg.Wait()
}

// VP-032 D-002 §1: sequential AllowRecord ≡ Allow-then-Record, including
// the deny path (no extra timestamp) and the absent-key register path.
func TestMemoryAllowRecordSequentialEquivalence(t *testing.T) {
	now := time.Date(2026, 9, 3, 12, 0, 0, 0, time.UTC)
	split := limiter(t, 15*time.Minute, 2, 1<<16)
	atomicL := limiter(t, 15*time.Minute, 2, 1<<16)

	step := func(key string, n time.Time) {
		t.Helper()
		want := split.Allow(key, n)
		if want {
			split.Record(key, n)
		}
		got := atomicL.AllowRecord(key, n)
		if got != want {
			t.Fatalf("AllowRecord(%s) = %v, want %v (Allow-then-Record)", key, got, want)
		}
	}

	step("k", now)     // 1/2
	step("k", now)     // 2/2
	step("k", now)     // deny
	step("other", now) // independent key
	if !atomicL.AllowRecord("k", now.Add(16*time.Minute)) {
		t.Fatal("AllowRecord after window must allow (and record)")
	}
	if got := len(atomicL.attempts["k"]); got != 1 {
		t.Fatalf("after window AllowRecord must replace pruned list, got %d", got)
	}
}

// VP-032 D-002 §2: AllowRecord false must not register a previously absent
// key. (Unreachable for max > 0 on a fresh key — absent is always allowed —
// so seed to the budget first, then deny.)
func TestMemoryAllowRecordDenyDoesNotGrow(t *testing.T) {
	l := limiter(t, 15*time.Minute, 1, 8)
	now := time.Now().UTC()
	if !l.AllowRecord("k", now) {
		t.Fatal("first AllowRecord must allow")
	}
	if l.AllowRecord("k", now) {
		t.Fatal("second AllowRecord at max=1 must deny")
	}
	if got := len(l.attempts["k"]); got != 1 {
		t.Fatalf("deny must not append, got %d timestamps", got)
	}
	before := len(l.attempts)
	if l.AllowRecord("k", now) {
		t.Fatal("repeated deny must stay denied")
	}
	if len(l.attempts) != before {
		t.Fatal("deny must not grow the map")
	}
}

// VP-032 D-002 §3: N concurrent AllowRecord on one key cannot penetrate max.
func TestMemoryAllowRecordConcurrentBudget(t *testing.T) {
	const max = 8
	const n = 64
	l := NewProvider().NewRateLimiter(time.Minute, max, 1<<16)
	now := time.Now().UTC()
	var allowed atomic.Int32
	var wg sync.WaitGroup
	wg.Add(n)
	for range n {
		go func() {
			defer wg.Done()
			if l.AllowRecord("k", now) {
				allowed.Add(1)
			}
		}()
	}
	wg.Wait()
	if got := allowed.Load(); got != max {
		t.Fatalf("concurrent AllowRecord allowed = %d, want %d (no TOCTOU penetration)", got, max)
	}
}

// GOAL-003 D-002 §2: Reserve occupies a slot that counts immediately, and
// Cancel removes exactly that slot while preserving all other history —
// the property key-wide Clear cannot provide.
func TestMemoryReserveCancelsOnlyItsOwnSlot(t *testing.T) {
	l := limiter(t, 15*time.Minute, 3, 1<<16)
	now := time.Now().UTC()

	// Seed two counted failures.
	l.Record("k", now)
	l.Record("k", now)

	token, ok := l.Reserve("k", now)
	if !ok {
		t.Fatal("Reserve under budget must allow")
	}
	if l.AllowRecord("k", now) {
		t.Fatal("after Reserve the bucket is at max and must deny")
	}

	// Cancel only the reserved slot: the two seeded failures survive.
	l.Cancel("k", token)
	if got := len(l.attempts["k"]); got != 2 {
		t.Fatalf("after Cancel history length = %d, want 2 (only the reserved slot rolled back)", got)
	}
	if !l.Allow("k", now) {
		t.Fatal("after Cancel the bucket must be below max again")
	}
	// A second Reserve on a fresh slot, cancelled, must also leave history.
	tok2, ok := l.Reserve("k", now)
	if !ok {
		t.Fatal("second Reserve must allow")
	}
	l.Cancel("k", tok2)
	if got := len(l.attempts["k"]); got != 2 {
		t.Fatalf("after second Cancel history length = %d, want 2", got)
	}
}

// GOAL-003 D-002 §1: Reserve deny registers nothing and returns a zero token.
func TestMemoryReserveDenyDoesNotGrow(t *testing.T) {
	l := limiter(t, 15*time.Minute, 1, 8)
	now := time.Now().UTC()
	if _, ok := l.Reserve("k", now); !ok {
		t.Fatal("first Reserve must allow")
	}
	token, ok := l.Reserve("k", now)
	if ok {
		t.Fatal("Reserve at max=1 must deny")
	}
	if token != 0 {
		t.Fatalf("denied Reserve token = %d, want 0", token)
	}
	if got := len(l.attempts["k"]); got != 1 {
		t.Fatalf("denied Reserve must not append, got %d slots", got)
	}
}

// GOAL-003 D-002 §2: Cancel is a no-op for an unknown token, after Clear, or
// when the key never existed; cancelling the last slot removes the key from
// both the map and the eviction order.
func TestMemoryReserveCancelNoOpAndCleanup(t *testing.T) {
	l := limiter(t, 15*time.Minute, 3, 8)
	now := time.Now().UTC()

	// Unknown key / unknown token: no-op, no panic, no entry created.
	l.Cancel("ghost", 123)
	if len(l.attempts) != 0 {
		t.Fatal("Cancel on an absent key must not create entries")
	}

	tok, ok := l.Reserve("k", now)
	if !ok {
		t.Fatal("Reserve must allow")
	}
	l.Cancel("k", 999) // unknown token
	if got := len(l.attempts["k"]); got != 1 {
		t.Fatalf("Cancel with unknown token must be a no-op, got %d slots", got)
	}
	l.Clear("k")
	l.Cancel("k", tok) // after Clear the slot is gone — no-op
	if _, exists := l.attempts["k"]; exists {
		t.Fatal("Cancel after Clear must not resurrect the key")
	}

	// Cancelling the last slot removes the key and its order entry.
	tok2, ok := l.Reserve("k2", now)
	if !ok {
		t.Fatal("Reserve on k2 must allow")
	}
	l.Cancel("k2", tok2)
	if _, exists := l.attempts["k2"]; exists {
		t.Fatal("last-slot Cancel must delete the key")
	}
	for _, k := range l.order {
		if k == "k2" {
			t.Fatal("cancelled key must be removed from the eviction order")
		}
	}
}

// GOAL-003 D-002 §1: N concurrent Reserve on one key cannot penetrate max
// (the same atomicity criterion as AllowRecord).
func TestMemoryReserveConcurrentBudget(t *testing.T) {
	const max = 8
	const n = 64
	l := NewProvider().NewRateLimiter(time.Minute, max, 1<<16)
	now := time.Now().UTC()
	var allowed atomic.Int32
	var wg sync.WaitGroup
	wg.Add(n)
	for range n {
		go func() {
			defer wg.Done()
			if _, ok := l.Reserve("k", now); ok {
				allowed.Add(1)
			}
		}()
	}
	wg.Wait()
	if got := allowed.Load(); got != max {
		t.Fatalf("concurrent Reserve allowed = %d, want %d (no TOCTOU penetration)", got, max)
	}
}

// GOAL-003 D-002 §1: Reserve/Cancel interleave with AllowRecord/Record/Allow
// under the same budget and stay race-clean (-race).
func TestMemoryReserveCancelConcurrent(t *testing.T) {
	l := NewProvider().NewRateLimiter(time.Minute, 4, 1<<16)
	now := time.Now().UTC()
	var wg sync.WaitGroup
	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range 300 {
				key := "ip|user"
				if tok, ok := l.Reserve(key, now); ok {
					if j%2 == 0 {
						l.Cancel(key, tok)
					}
				}
				_ = l.AllowRecord(key, now)
				_ = l.Allow(key, now)
				if j%11 == 0 {
					l.Clear(key)
				}
			}
		}()
	}
	wg.Wait()
}

func itoa(i int) string {
	const digits = "0123456789"
	if i == 0 {
		return "0"
	}
	var buf [16]byte
	pos := len(buf)
	for i > 0 {
		pos--
		buf[pos] = digits[i%10]
		i /= 10
	}
	return string(buf[pos:])
}
