package ratelimit

// Memory provider for the kernel rate-limiter port (VP-027 / workspace-027
// GOAL-003 D-001, R2): the legacy loginRateLimiter semantics (handler
// package) evolved into the port implementation — sliding-window failure
// budget, Allow that never registers a key, Record-only map growth with FIFO
// capacity eviction, Retry-After with the frozen minimum-second behavior,
// Clear on success. Window pruning and Retry-After computation delegate to
// the kernel executable predicates (D-002 §3/§5) so every provider stays
// bit-identical to the frozen contract. Safe for concurrent use; no
// background goroutine (D-002 §7). D-001 P1: spraying distinct keys cannot
// grow the map — only Record reaches the capacity eviction.

import (
	"sync"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

var (
	_ kernel.RateLimiterProvider = (*Provider)(nil)
	_ kernel.RateLimiter         = (*Memory)(nil)
)

// Provider is the in-memory factory (D-002 §1). The composition root owns a
// single instance (fx.Provide); a future Redis-tier provider implements the
// same kernel interfaces when RT-Q05 triggers.
type Provider struct{}

// NewProvider returns the memory provider factory.
func NewProvider() *Provider { return &Provider{} }

// NewRateLimiter builds a process-local sliding-window limiter with the given
// window, failure budget max and distinct-key capacity. A capacity of zero or
// less falls back to kernel.DefaultRateLimiterCapacity (1 << 16).
func (p *Provider) NewRateLimiter(window time.Duration, max, capacity int) kernel.RateLimiter {
	if capacity <= 0 {
		capacity = kernel.DefaultRateLimiterCapacity
	}
	return &Memory{
		window:   window,
		max:      max,
		attempts: make(map[string][]time.Time),
		capacity: capacity,
	}
}

// Memory implements kernel.RateLimiter over an in-process sliding window
// (legacy loginRateLimiter semantics, D-002 §3/§4).
type Memory struct {
	mu     sync.Mutex
	window time.Duration
	max    int
	// attempts holds the recent failure timestamps per client identity.
	attempts map[string][]time.Time
	// order tracks key insertion so the map stays bounded: when capacity is
	// reached the oldest key is evicted (best-effort memory guard, D-001 P1).
	order    []string
	capacity int
}

// Allow reports whether key may attempt now. It never creates a new map entry
// (D-002 §1): an absent key (no failures yet) is always allowed; for an
// existing key, stale window entries are pruned in place before the limit
// check.
func (l *Memory) Allow(key string, now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	list, exists := l.attempts[key]
	if !exists {
		return true
	}
	kept := list[:0]
	for _, t := range list {
		if kernel.RateLimiterInWindow(t, l.window, now) {
			kept = append(kept, t)
		}
	}
	if len(kept) >= l.max {
		l.attempts[key] = kept
		return false
	}
	l.attempts[key] = kept
	return true
}

// RetryAfterSeconds is the remaining window after the oldest in-window
// failure (D-002 §5; the kernel predicate is the single semantic authority).
// Zero when the key is allowed; callers invoke it only after Allow returned
// false.
func (l *Memory) RetryAfterSeconds(key string, now time.Time) int {
	l.mu.Lock()
	defer l.mu.Unlock()
	list := l.attempts[key]
	if len(list) == 0 {
		return 0
	}
	oldest := list[0]
	for _, t := range list[1:] {
		if t.Before(oldest) {
			oldest = t
		}
	}
	return kernel.RateLimiterRetryAfterSeconds(oldest, l.window, now)
}

// Record registers one failed attempt for the key, creating the map entry if
// needed. Bounded: when the map exceeds capacity the oldest key is evicted,
// so an attacker cannot exhaust memory by spraying distinct client
// identities (D-001 P1).
func (l *Memory) Record(key string, now time.Time) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if _, exists := l.attempts[key]; !exists {
		if len(l.attempts) >= l.capacity {
			if len(l.order) > 0 {
				oldest := l.order[0]
				l.order = l.order[1:]
				delete(l.attempts, oldest)
			}
		}
		l.order = append(l.order, key)
	}
	l.attempts[key] = append(l.attempts[key], now)
}

// Clear drops every failure for the key. Called after a successful attempt so
// a legitimate client never accumulates a poisoned bucket (D-001 P1).
func (l *Memory) Clear(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if _, exists := l.attempts[key]; !exists {
		return
	}
	delete(l.attempts, key)
	for i, k := range l.order {
		if k == key {
			l.order = append(l.order[:i], l.order[i+1:]...)
			break
		}
	}
}
