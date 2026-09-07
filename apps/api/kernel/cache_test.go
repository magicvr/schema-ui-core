package kernel

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

// Compile-time port-surface guard (F-005b): the frozen signatures must stay
// implementable by any provider — a stub suffices to lock the method sets.
type stubCacheView struct{}

func (stubCacheView) Get(context.Context, string) ([]byte, bool)              { return nil, false }
func (stubCacheView) Set(context.Context, string, []byte, ExpiryPolicy) error { return nil }
func (stubCacheView) Delete(context.Context, string) error                    { return nil }

var (
	_ CacheView    = stubCacheView{}
	_ ExpiryPolicy = absPolicyStub{}
)

type absPolicyStub struct{}

func (absPolicyStub) ExpireAt(now time.Time) time.Time { return now.Add(time.Hour) }
func (absPolicyStub) Refresh(_ time.Time, previous time.Time) (time.Time, bool) {
	return previous, false
}

// Contract-level fast tests for the kernel cache port (VP-026 GOAL-002 D-002
// §2/§3/§8): validation helpers and sentinel errors. Provider behavior tests
// land with the R2 in-memory provider.

func TestValidCacheNamespace(t *testing.T) {
	tests := []struct {
		name string
		ns   CacheNamespace
		want bool
	}{
		{"single letter", CacheNamespace("a"), true},
		{"single digit", CacheNamespace("0"), true},
		{"letter then hyphen", CacheNamespace("ab-cd"), true},
		{"max length", CacheNamespace(strings.Repeat("a", 64)), true},
		{"hyphen inside", CacheNamespace("wallet-balance"), true},
		{"digits inside", CacheNamespace("session2fa"), true},

		{"empty", CacheNamespace(""), false},
		{"uppercase", CacheNamespace("Wallet"), false},
		{"leading hyphen", CacheNamespace("-wallet"), false},
		{"trailing hyphen", CacheNamespace("wallet-"), false},
		{"double hyphen", CacheNamespace("wallet--balance"), false},
		{"underscore", CacheNamespace("wallet_balance"), false},
		{"dot", CacheNamespace("wallet.balance"), false},
		{"space", CacheNamespace("wallet balance"), false},
		{"unicode", CacheNamespace("钱包"), false},
		{"over max length", CacheNamespace(strings.Repeat("a", 65)), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidCacheNamespace(tt.ns); got != tt.want {
				t.Errorf("ValidCacheNamespace(%q) = %v, want %v", tt.ns, got, tt.want)
			}
		})
	}
}

func TestValidCacheKey(t *testing.T) {
	long256 := strings.Repeat("k", 256)
	long257 := strings.Repeat("k", 257)
	tests := []struct {
		name string
		key  string
		want bool
	}{
		{"single char", "a", true},
		{"hyphen and colon", "user:42-profile", true},
		{"utf8", "用户:42", true},
		{"max length", long256, true},
		{"spaces", "key with spaces", true},

		{"empty", "", false},
		{"over max length", long257, false},
		{"tab", "a\tb", false},
		{"newline", "a\nb", false},
		{"del byte", "a\x7fb", false},
		{"nul byte", "a\x00b", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidCacheKey(tt.key); got != tt.want {
				t.Errorf("ValidCacheKey(%q) = %v, want %v", tt.key, got, tt.want)
			}
		})
	}
}

func TestCacheSentinelErrors(t *testing.T) {
	sentinels := []error{
		ErrInvalidCacheNamespace,
		ErrInvalidCacheKey,
		ErrInvalidCacheValue,
		ErrInvalidCachePolicy,
	}
	for _, sentinel := range sentinels {
		if sentinel == nil {
			t.Fatal("cache sentinel error must not be nil")
		}
		// errors.Is must resolve through a %w wrap.
		if !errors.Is(fmt.Errorf("wrap: %w", sentinel), sentinel) {
			t.Errorf("errors.Is through wrap failed for %v", sentinel)
		}
	}
	// Sentinels must be pairwise distinct.
	seen := map[string]error{}
	for _, sentinel := range sentinels {
		if previous, exists := seen[sentinel.Error()]; exists {
			t.Errorf("sentinel collision: %v and %v share message", previous, sentinel)
		}
		seen[sentinel.Error()] = sentinel
	}
}

func TestCacheEntryExpired(t *testing.T) {
	now := time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name      string
		expiresAt time.Time
		now       time.Time
		want      bool
	}{
		{"zero never expires", time.Time{}, now, false},
		{"far future not expired", now.Add(time.Hour), now, false},
		{"exactly at expiry is expired", now, now, true},
		{"past is expired", now.Add(-time.Second), now, true},
		{"before expiry not expired", now, now.Add(-time.Second), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CacheEntryExpired(tt.expiresAt, tt.now); got != tt.want {
				t.Errorf("CacheEntryExpired(%v, %v) = %v, want %v", tt.expiresAt, tt.now, got, tt.want)
			}
		})
	}
}

func TestValidateCacheSet(t *testing.T) {
	policy := absPolicyStub{}
	tests := []struct {
		name   string
		key    string
		value  []byte
		policy ExpiryPolicy
		want   error
	}{
		{"valid", "k", []byte("v"), policy, nil},
		{"empty value allowed", "k", []byte{}, policy, nil},
		{"invalid key", "", []byte("v"), policy, ErrInvalidCacheKey},
		{"nil value", "k", nil, policy, ErrInvalidCacheValue},
		{"nil policy", "k", []byte("v"), nil, ErrInvalidCachePolicy},
		// Order: key is checked before value/policy.
		{"invalid key and nil value reports key", "", nil, policy, ErrInvalidCacheKey},
		{"invalid key and nil policy reports key", "", []byte("v"), nil, ErrInvalidCacheKey},
		{"nil value and nil policy reports value", "k", nil, nil, ErrInvalidCacheValue},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateCacheSet(tt.key, tt.value, tt.policy)
			if !errors.Is(got, tt.want) {
				t.Errorf("ValidateCacheSet(%q, %v, %v) = %v, want %v", tt.key, tt.value, tt.policy, got, tt.want)
			}
		})
	}
}
