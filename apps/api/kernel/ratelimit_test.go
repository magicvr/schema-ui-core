package kernel

import (
	"testing"
	"time"
)

// Compile-time port-surface assertions (D-002 §10): the interface shapes are
// implementable by stubs, i.e. the contract is not an unimplementable fiction.
type stubRateLimiter struct{}

func (stubRateLimiter) Allow(string, time.Time) bool { return true }
func (stubRateLimiter) Record(string, time.Time)     {}
func (stubRateLimiter) RetryAfterSeconds(string, time.Time) int {
	return 0
}
func (stubRateLimiter) Clear(string) {}

type stubRateLimiterProvider struct{}

func (stubRateLimiterProvider) NewRateLimiter(time.Duration, int, int) RateLimiter {
	return stubRateLimiter{}
}

var (
	_ RateLimiter         = stubRateLimiter{}
	_ RateLimiterProvider = stubRateLimiterProvider{}
)

func TestDefaultRateLimiterCapacity(t *testing.T) {
	if got, want := DefaultRateLimiterCapacity, 1<<16; got != want {
		t.Fatalf("DefaultRateLimiterCapacity = %d, want %d", got, want)
	}
}

func TestRateLimiterInWindow(t *testing.T) {
	now := time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)
	window := 15 * time.Minute

	cases := []struct {
		name   string
		t      time.Time
		window time.Duration
		want   bool
	}{
		{"failure just now", now, window, true},
		{"one nanosecond inside", now.Add(-window + time.Nanosecond), window, true},
		{"exactly on cutoff", now.Add(-window), window, false},
		{"one nanosecond beyond cutoff", now.Add(-window - time.Nanosecond), window, false},
		{"far past", now.Add(-2 * window), window, false},
		{"future failure", now.Add(time.Hour), window, true},
		{"zero window same instant", now, 0, false},
		{"zero window a nanosecond earlier", now.Add(-time.Nanosecond), 0, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := RateLimiterInWindow(tc.t, tc.window, now); got != tc.want {
				t.Fatalf("RateLimiterInWindow(%v, %v, %v) = %v, want %v", tc.t, tc.window, now, got, tc.want)
			}
		})
	}
}

func TestRateLimiterRetryAfterSeconds(t *testing.T) {
	now := time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)
	window := 15 * time.Minute

	cases := []struct {
		name   string
		oldest time.Time
		want   int
	}{
		{"full window remaining", now, 900},
		{"five minutes remaining", now.Add(-10 * time.Minute), 300},
		{"thirty seconds remaining", now.Add(-15*time.Minute + 30*time.Second), 30},
		{"window just elapsed", now.Add(-15 * time.Minute), 1},
		{"beyond window", now.Add(-16 * time.Minute), 1},
		{"sub-second rounds down", now.Add(-15*time.Minute + 400*time.Millisecond), 0},
		{"sub-second rounds up", now.Add(-15*time.Minute + 600*time.Millisecond), 1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := RateLimiterRetryAfterSeconds(tc.oldest, window, now); got != tc.want {
				t.Fatalf("RateLimiterRetryAfterSeconds(%v, %v, %v) = %d, want %d", tc.oldest, window, now, got, tc.want)
			}
		})
	}
}
