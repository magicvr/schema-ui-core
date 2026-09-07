package cache

// Base expiry policies for the kernel cache port (VP-026 / workspace-026
// GOAL-002 D-002 §5, GOAL-003 D-001). Both policies are stateless values:
// safe for concurrent use and freely shareable across entries.

import (
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

var (
	_ kernel.ExpiryPolicy = AbsoluteExpiry{}
	_ kernel.ExpiryPolicy = SlidingExpiry{}
)

// AbsoluteExpiry expires an entry a fixed TTL after it was Set; Get hits do
// NOT refresh the expiry. A non-positive TTL means "never expires" (zero
// time.Time), matching the port contract's zero-expiry rule (D-002 §5).
type AbsoluteExpiry struct {
	TTL time.Duration
}

// ExpireAt implements kernel.ExpiryPolicy.
func (a AbsoluteExpiry) ExpireAt(now time.Time) time.Time {
	if a.TTL <= 0 {
		return time.Time{}
	}
	return now.Add(a.TTL)
}

// Refresh implements kernel.ExpiryPolicy: absolute policies never extend.
func (a AbsoluteExpiry) Refresh(_ time.Time, previous time.Time) (time.Time, bool) {
	return previous, false
}

// SlidingExpiry expires an entry Window after its last access: both Set and
// every Get hit refresh the expiry (D-002 §5). A non-positive Window means
// "never expires and never refreshes".
type SlidingExpiry struct {
	Window time.Duration
}

// ExpireAt implements kernel.ExpiryPolicy.
func (s SlidingExpiry) ExpireAt(now time.Time) time.Time {
	if s.Window <= 0 {
		return time.Time{}
	}
	return now.Add(s.Window)
}

// Refresh implements kernel.ExpiryPolicy: sliding policies extend from now.
func (s SlidingExpiry) Refresh(now time.Time, _ time.Time) (time.Time, bool) {
	if s.Window <= 0 {
		return time.Time{}, false
	}
	return now.Add(s.Window), true
}
