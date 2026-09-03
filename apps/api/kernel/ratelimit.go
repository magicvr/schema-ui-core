package kernel

// Kernel rate-limiter port (VP-027 / workspace-027 GOAL-002 D-002, R1;
// VP-032 / workspace-032 GOAL-002 D-002 additive AllowRecord).
//
// The port is the only rate-limiter contract for the kernel and every module:
// a process-local sliding-window failure budget keyed by opaque strings, with
// the Allow/Record/Clear split that keeps Allow side-effect free (D-001 P1 —
// spraying distinct keys cannot grow the map). AllowRecord is the TOCTOU-free
// equivalent of Allow-then-Record under one lock (VP-032). Handlers obtain
// limiters from an injected kernel.RateLimiterProvider factory; public types
// carry neither provider handles nor key structure (the existing
// `IP|identifier`, `op|IP|user` and bare-IP key shapes are caller-side
// conventions, VP-027 D-002 §2).
//
// Contract frozen by workspace-027 GOAL-002 D-002 + workspace-032 GOAL-002 D-002:
//
//   - Allow never registers a key; Record is the only path that creates map
//     entries, so the capacity eviction in Record stays reachable.
//   - AllowRecord is atomic Allow-then-Record: false does not register; true
//     uses Record's map-growth / eviction path.
//   - Window semantics are sliding with in-path pruning; the executable
//     predicates RateLimiterInWindow and RateLimiterRetryAfterSeconds below
//     are the single semantic authority every provider MUST use (W12 D-002:
//     Retry-After behavior preserved).
//   - capacity <= 0 in the factory falls back to DefaultRateLimiterCapacity.
//   - All methods are safe for concurrent use; no background goroutine, so no
//     new lifecycle (VP-021 shutdown obligations do not trigger).

import "time"

// RateLimiter is the kernel rate-limiter port (R1). Implementations MUST be
// safe for concurrent use and MUST keep Allow side-effect free (no key
// registration, no record). The now parameter is an explicit clock injection
// for deterministic tests; production callers pass time.Now().UTC().
type RateLimiter interface {
	// Allow reports whether key may attempt now. It never creates a new map
	// entry: an absent key (no failures yet) is always allowed. Stale
	// in-window entries are pruned in place before the limit check.
	Allow(key string, now time.Time) bool
	// Record registers one failed attempt for key, creating the entry if
	// needed. Bounded: implementations evict the oldest key when capacity is
	// reached (D-001 P1 memory guard).
	Record(key string, now time.Time)
	// AllowRecord atomically reports whether key may attempt now and, if so,
	// records that attempt. Sequential semantics match Allow followed by
	// Record (VP-032 / workspace-032 GOAL-002 D-002): an absent key is always
	// allowed then registered; a key whose in-window count is already >= max
	// is denied and not recorded. Callers MUST invoke RetryAfterSeconds only
	// after AllowRecord (or Allow) returned false. New call sites SHOULD use
	// AllowRecord; Allow and Record remain for compatibility.
	AllowRecord(key string, now time.Time) bool
	// RetryAfterSeconds is the remaining window after the oldest in-window
	// failure. Callers MUST invoke it only after Allow or AllowRecord
	// returned false (W12 D-002: Retry-After header semantic); the value
	// follows RateLimiterRetryAfterSeconds.
	RetryAfterSeconds(key string, now time.Time) int
	// Clear drops every failure for key. Called after a successful attempt so
	// a legitimate client never accumulates a poisoned bucket.
	Clear(key string)
}

// RateLimiterProvider is the kernel-level factory for rate limiters (R1).
// Handlers inject this interface instead of constructing provider types, so
// the in-memory provider (R2) and any future Redis-tier provider implement
// the same contract with zero consumer changes.
type RateLimiterProvider interface {
	// NewRateLimiter returns a rate limiter with the given window, failure
	// budget max and distinct-key capacity. A capacity <= 0 falls back to
	// DefaultRateLimiterCapacity.
	NewRateLimiter(window time.Duration, max, capacity int) RateLimiter
}

// DefaultRateLimiterCapacity is the fallback distinct-key bound when a
// capacity of zero or less is passed to RateLimiterProvider.NewRateLimiter
// (D-001 P1 memory guard; matches the legacy 1<<16 default).
const DefaultRateLimiterCapacity = 1 << 16

// RateLimiterInWindow is the frozen in-window predicate (D-002 §3): a failure
// timestamp t stays inside the sliding window at now iff t.After(now-window).
// Providers MUST use this predicate for pruning so window semantics match the
// legacy loginRateLimiter on every provider (a timestamp exactly on the
// cutoff is not kept).
func RateLimiterInWindow(t time.Time, window time.Duration, now time.Time) bool {
	return t.After(now.Add(-window))
}

// RateLimiterRetryAfterSeconds is the frozen Retry-After computation (D-002
// §5, W12 D-002): seconds remaining after the oldest in-window failure, with
// a minimum of 1 when the window has elapsed (remain <= 0), then
// remain.Round(time.Second)/time.Second. Providers MUST use this function so
// every Retry-After value matches the legacy behavior exactly.
func RateLimiterRetryAfterSeconds(oldest time.Time, window time.Duration, now time.Time) int {
	remain := oldest.Add(window).Sub(now)
	if remain <= 0 {
		return 1
	}
	return int(remain.Round(time.Second) / time.Second)
}
