package cache

// In-memory provider for the kernel cache port (VP-026 / workspace-026
// GOAL-003 D-001, F-001 adjudication): the maxEntries budget is PROCESS-WIDE —
// after every Set the TOTAL entry count across ALL namespaces (including
// not-yet-swept expired entries) is <= maxEntries, with FIFO eviction of the
// globally-oldest entry. Lazy TTL cleanup runs on the read/write paths only;
// no background goroutine, no lifecycle (D-002 §5). Safe for concurrent use.

import (
	"container/list"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

var _ kernel.Cache = (*Memory)(nil)

// Memory implements kernel.Cache over an in-process store.
//
// Bound guarantee (D-001 + F-001 adjudication): after every Set the TOTAL
// number of stored entries across every namespace — including expired ones
// that have not been swept yet — is <= maxEntries. When the budget is
// reached, the globally OLDEST inserted entry is evicted (FIFO), regardless
// of which namespace inserted it; a busy namespace can therefore displace
// idle namespaces' old entries (that is the cost of a shared process budget).
// Expired entries are removed lazily: a Get that finds an expired entry, a
// Set that overwrites one, and Delete all drop it from both the map and the
// global FIFO list.
//
// Concurrency: a single mutex serializes every operation. Values are copied
// at the boundary (Set takes a copy; Get returns a copy).
type Memory struct {
	mu         sync.Mutex
	maxEntries int
	count      int // TOTAL entries across all namespaces
	now        func() time.Time
	spaces     map[kernel.CacheNamespace]*cacheSpace
	order      *list.List // global FIFO across namespaces; elements hold *cacheEntry
}

// cacheSpace is one namespace's key map. Insertion order lives on the
// Memory-level list (process-wide budget), not here.
type cacheSpace struct {
	entries map[string]*cacheEntry
}

// cacheEntry is one stored value with its expiry bookkeeping.
type cacheEntry struct {
	key       string
	ns        kernel.CacheNamespace
	value     []byte
	expiresAt time.Time
	policy    kernel.ExpiryPolicy
	element   *list.Element
}

// NewMemory builds a bounded in-memory cache provider. maxEntries must be
// positive (fail closed; the config layer enforces the same rule earlier —
// the cache.max_entries key rejects non-positive values at load).
func NewMemory(maxEntries int) (*Memory, error) {
	if maxEntries <= 0 {
		return nil, fmt.Errorf("cache: maxEntries must be positive (got %d)", maxEntries)
	}
	return &Memory{
		maxEntries: maxEntries,
		now:        time.Now,
		spaces:     make(map[kernel.CacheNamespace]*cacheSpace),
		order:      list.New(),
	}, nil
}

// newMemoryWithClock is the test seam for clock-injected expiry tests
// (D-002 §5: providers use time.Now in production).
func newMemoryWithClock(maxEntries int, now func() time.Time) (*Memory, error) {
	m, err := NewMemory(maxEntries)
	if err != nil {
		return nil, err
	}
	if now != nil {
		m.now = now
	}
	return m, nil
}

// Namespace implements kernel.Cache: validates the namespace fail-closed and
// returns a scoped view. A missing namespace is created on first access.
func (m *Memory) Namespace(ns kernel.CacheNamespace) (kernel.CacheView, error) {
	if !kernel.ValidCacheNamespace(ns) {
		return nil, kernel.ErrInvalidCacheNamespace
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	space, ok := m.spaces[ns]
	if !ok {
		space = &cacheSpace{entries: make(map[string]*cacheEntry)}
		m.spaces[ns] = space
	}
	return &memoryView{m: m, ns: ns, space: space}, nil
}

// MaxEntries reports the configured process-wide bounded-entry budget. The
// port contract itself does not expose capacity (D-002 §6); this accessor
// serves wiring tests and operator diagnostics only.
func (m *Memory) MaxEntries() int { return m.maxEntries }

// Len reports the total number of stored entries across all namespaces
// (including not-yet-swept expired ones) — the live count against the
// process-wide budget. Test/diagnostic accessor; not part of the port.
func (m *Memory) Len() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.count
}

// memoryView is a single-namespace scope. It is only valid while its Memory
// stays alive; all operations take the parent mutex, so views are safe for
// concurrent use (D-002 §7).
type memoryView struct {
	m     *Memory
	ns    kernel.CacheNamespace
	space *cacheSpace
}

// Get implements kernel.CacheView. A stored empty value (non-nil zero-length
// slice) hits with (empty, true); an expired entry is removed lazily and
// reported as a miss; an invalid key is a miss (no error channel, D-002 §3).
func (v *memoryView) Get(_ context.Context, key string) ([]byte, bool) {
	if !kernel.ValidCacheKey(key) {
		return nil, false
	}
	v.m.mu.Lock()
	defer v.m.mu.Unlock()
	e, ok := v.space.entries[key]
	if !ok {
		return nil, false
	}
	now := v.m.now()
	if kernel.CacheEntryExpired(e.expiresAt, now) {
		v.removeLocked(e)
		return nil, false
	}
	if newExpiry, refresh := e.policy.Refresh(now, e.expiresAt); refresh {
		e.expiresAt = newExpiry
	}
	return copyBytes(e.value), true
}

// Set implements kernel.CacheView. Validation (kernel.ValidateCacheSet:
// key -> value -> policy) runs before any storage mutation. A LIVE overwrite
// (same key, not expired) keeps its global FIFO position even when the policy
// changes; overwriting a DEAD (expired) entry drops it and re-inserts at the
// back. Capacity is enforced BEFORE insertion: while the process-wide total
// is at maxEntries, the globally-oldest entry is evicted (F-001
// adjudication: process-global budget).
func (v *memoryView) Set(_ context.Context, key string, value []byte, policy kernel.ExpiryPolicy) error {
	if err := kernel.ValidateCacheSet(key, value, policy); err != nil {
		return err
	}
	v.m.mu.Lock()
	defer v.m.mu.Unlock()
	now := v.m.now()
	if e, ok := v.space.entries[key]; ok {
		if kernel.CacheEntryExpired(e.expiresAt, now) {
			// Dead entry: drop and re-insert below (position refresh).
			v.removeLocked(e)
		} else {
			// Live overwrite keeps insertion order (FIFO semantics), including
			// across policy changes — the entry is the same key, refreshed.
			e.value = copyBytes(value)
			e.expiresAt = policy.ExpireAt(now)
			e.policy = policy
			return nil
		}
	}
	for v.m.count >= v.m.maxEntries {
		v.evictLocked()
	}
	entry := &cacheEntry{
		key:       key,
		ns:        v.ns,
		value:     copyBytes(value),
		expiresAt: policy.ExpireAt(now),
		policy:    policy,
	}
	entry.element = v.m.order.PushBack(entry)
	v.space.entries[key] = entry
	v.m.count++
	return nil
}

// Delete implements kernel.CacheView: idempotent (deleting an absent key
// succeeds). An expired entry is removed from both structures.
func (v *memoryView) Delete(_ context.Context, key string) error {
	if !kernel.ValidCacheKey(key) {
		return kernel.ErrInvalidCacheKey
	}
	v.m.mu.Lock()
	defer v.m.mu.Unlock()
	if e, ok := v.space.entries[key]; ok {
		v.removeLocked(e)
	}
	return nil
}

// evictLocked removes the globally OLDEST inserted entry (FIFO, F-001:
// process-wide budget). Caller holds m.mu and the store must be non-empty.
func (v *memoryView) evictLocked() {
	front := v.m.order.Front()
	if front == nil {
		return
	}
	v.removeLocked(front.Value.(*cacheEntry))
}

// removeLocked drops an entry from the global FIFO list and its namespace map,
// decrementing the process-wide count.
func (v *memoryView) removeLocked(e *cacheEntry) {
	v.m.order.Remove(e.element)
	delete(v.m.spaces[e.ns].entries, e.key)
	v.m.count--
}

// copyBytes copies src preserving empty-value semantics: an empty-but-non-nil
// input stays non-nil empty — the port contract distinguishes a stored empty
// value from a miss (D-002 §1/§4). Plain append([]byte(nil), src...) would
// collapse an empty src into nil. A nil input (rejected by ValidateCacheSet
// on the Set side; defensive here) becomes a non-nil empty slice.
func copyBytes(src []byte) []byte {
	out := make([]byte, len(src))
	copy(out, src)
	return out
}
