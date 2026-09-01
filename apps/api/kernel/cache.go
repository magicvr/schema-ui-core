package kernel

// Kernel cache port (VP-026 / workspace-026 GOAL-002 D-002, R1).
//
// The port is the only cache contract for the kernel and every module: a
// namespace-scoped Get/Set/Delete of []byte values with pluggable expiry
// policies. Public types carry neither provider handles nor serialization
// policy: callers address values by (namespace, key) pairs inside a scoped
// view, and the active provider (in-memory default R2, Redis seam R3) resolves
// storage. This mirrors the kernel.ObjectStore storage-port precedent: domain
// code consumes the port, never a backend handle.
//
// Contract frozen by workspace-026 GOAL-002 D-002:
//
//   - Non-generic []byte payload: miss is distinguished from a stored empty
//     value by the ok bool; nil values cannot be stored (fail closed).
//   - Namespace is an explicit scoped view (Cache.Namespace) with open-set
//     shape validation (fail closed); the Redis key-prefix mapping
//     (<ns>:<key>) is reserved for the R3 seam document.
//   - Expiry is lazy: expired entries are removed on the read/write paths
//     only; no background goroutine, no new lifecycle, no shutdown-drain
//     obligation (VP-021).
//   - All CacheView methods are safe for concurrent use.

import (
	"context"
	"errors"
	"regexp"
	"time"
)

// CacheNamespace isolates cache key spaces. The set is OPEN (future business
// modules create their own namespaces) but every value must pass
// ValidCacheNamespace or adapters fail closed (GOAL-002 D-002 §2).
type CacheNamespace string

const (
	// CacheNamespaceMaxLen is the upper bound on namespace length (in bytes).
	CacheNamespaceMaxLen = 64
)

// cacheNamespacePattern is the only namespace shape the port accepts: one or
// more lower-case alphanumeric segments joined by single hyphens — must not
// start/end with a hyphen and must not contain consecutive hyphens. Enforcing
// it at the port keeps crafted namespaces from turning into invalid Redis key
// prefixes or cross-module collisions, fail-closed.
var cacheNamespacePattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// ValidCacheNamespace reports whether ns is a well-formed cache namespace.
// Unknown namespaces are NOT a programming error here (open set) — any shape-
// valid value is accepted; invalid shapes are rejected before touching the
// provider.
func ValidCacheNamespace(ns CacheNamespace) bool {
	return len(ns) > 0 && len(ns) <= CacheNamespaceMaxLen && cacheNamespacePattern.MatchString(string(ns))
}

const (
	// CacheKeyMaxLen is the upper bound on key length (in bytes).
	CacheKeyMaxLen = 256
)

// ValidCacheKey reports whether key is a well-formed cache key: non-empty,
// at most CacheKeyMaxLen bytes, and free of control characters. UTF-8
// multibyte characters are allowed (the byte scan only rejects bytes below
// 0x20 and the DEL byte 0x7f).
func ValidCacheKey(key string) bool {
	if len(key) == 0 || len(key) > CacheKeyMaxLen {
		return false
	}
	for i := 0; i < len(key); i++ {
		if key[i] < 0x20 || key[i] == 0x7f {
			return false
		}
	}
	return true
}

// Sentinel errors (GOAL-002 D-002 §8). Callers use errors.Is; miss is NOT an
// error (Get carries the ok bool).
var (
	// ErrInvalidCacheNamespace guards Cache.Namespace against malformed
	// namespace values.
	ErrInvalidCacheNamespace = errors.New("kernel: invalid cache namespace")
	// ErrInvalidCacheKey guards Set/Delete against malformed keys.
	ErrInvalidCacheKey = errors.New("kernel: invalid cache key")
	// ErrInvalidCacheValue guards Set against nil values.
	ErrInvalidCacheValue = errors.New("kernel: invalid cache value")
	// ErrInvalidCachePolicy guards Set against a nil expiry policy.
	ErrInvalidCachePolicy = errors.New("kernel: invalid cache expiry policy")
)

// CacheEntryExpired reports whether an entry expiring at expiresAt is expired
// at now (GOAL-002 D-002 §5): a zero expiry means "never expires" and is never
// expired; an entry is expired the moment now reaches its expiry instant
// (now.Before(expiresAt) is false). Adapters MUST use this predicate for
// expiry checks so lazy cleanup matches the frozen semantics on every
// provider.
func CacheEntryExpired(expiresAt, now time.Time) bool {
	return !expiresAt.IsZero() && !now.Before(expiresAt)
}

// ValidateCacheSet is the executable fail-closed entry for Set arguments
// (GOAL-002 D-002 §4/§8): key shape, then value, then policy — in that order.
// Providers MUST call it before touching storage, and MUST NOT store a value
// that fails it. Get has no error channel; callers surface invalid keys
// through this check on the Set/Delete side.
func ValidateCacheSet(key string, value []byte, policy ExpiryPolicy) error {
	if !ValidCacheKey(key) {
		return ErrInvalidCacheKey
	}
	if value == nil {
		return ErrInvalidCacheValue
	}
	if policy == nil {
		return ErrInvalidCachePolicy
	}
	return nil
}

// ExpiryPolicy is the pluggable entry-expiry contract (GOAL-002 D-002 §5).
// Implementations must be stateless and safe for concurrent use; the two
// frozen base policies (absolute, sliding) ship with the R2 in-memory
// provider. A zero time.Time expiry means "never expires".
type ExpiryPolicy interface {
	// ExpireAt returns the expiry instant for a Set performed at now.
	ExpireAt(now time.Time) time.Time
	// Refresh is called after a Get hit at now on an entry expiring at
	// previous. It returns the entry's new expiry instant and whether the
	// expiry was extended. Absolute policies return (previous, false);
	// sliding policies return (now+window, true).
	Refresh(now time.Time, previous time.Time) (time.Time, bool)
}

// Cache is the kernel cache port (R1). Callers obtain a namespace-scoped
// view via Namespace — the only way to touch values — so namespace isolation
// is structural, not conventional. Implementations must validate the
// namespace (ValidCacheNamespace) and fail closed on violations.
type Cache interface {
	// Namespace returns the scoped view for ns. Invalid namespaces yield
	// ErrInvalidCacheNamespace; valid namespaces always succeed and the
	// returned view is safe for concurrent use.
	Namespace(ns CacheNamespace) (CacheView, error)
}

// CacheView is a single-namespace scope of the cache port. All methods are
// safe for concurrent use (GOAL-002 D-002 §7). Keys are validated before
// reaching the provider: Set/Delete reject invalid keys with
// ErrInvalidCacheKey; Get has no error channel and treats an invalid key as a
// miss (nil, false) — callers surface such programming errors through the
// Set/Delete side.
type CacheView interface {
	// Get returns the value stored under key and whether it was a hit.
	// A miss (absent or expired) returns (nil, false); a stored empty value
	// (non-nil zero-length slice) hits with (empty slice, true) — ok
	// distinguishes miss from a stored empty value. The returned slice is a
	// fresh copy: adapters must never share their internal buffer with the
	// caller.
	Get(ctx context.Context, key string) ([]byte, bool)
	// Set stores value under key with the given expiry policy. value must be
	// non-nil (ErrInvalidCacheValue) and policy must be non-nil
	// (ErrInvalidCachePolicy). The input slice is copied at the boundary.
	Set(ctx context.Context, key string, value []byte, policy ExpiryPolicy) error
	// Delete removes key. Deleting an absent key succeeds (idempotent).
	Delete(ctx context.Context, key string) error
}
