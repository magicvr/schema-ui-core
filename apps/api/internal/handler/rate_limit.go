package handler

import (
	"net"
	"net/http"
	"sync"
	"time"
)

// loginRateLimiter is a small in-memory sliding-window limiter for failed
// login attempts (D2): per client IP, at most max failures per window. It is
// process-local and best-effort — it does not protect against distributed
// attacks, but it stops trivial online brute force on a single instance.
type loginRateLimiter struct {
	mu       sync.Mutex
	window   time.Duration
	max      int
	attempts map[string][]time.Time
}

func newLoginRateLimiter(window time.Duration, max int) *loginRateLimiter {
	return &loginRateLimiter{window: window, max: max, attempts: make(map[string][]time.Time)}
}

// allow reports whether the key may attempt a login now.
func (l *loginRateLimiter) allow(key string, now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	cutoff := now.Add(-l.window)
	list := l.attempts[key]
	kept := list[:0]
	for _, t := range list {
		if t.After(cutoff) {
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

// record registers one failed attempt for the key.
func (l *loginRateLimiter) record(key string, now time.Time) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.attempts[key] = append(l.attempts[key], now)
}

// clientIP returns the direct peer IP; proxy-forwarded headers are not trusted
// (they are client-controlled), so behind a reverse proxy all clients share the
// proxy IP and the limiter degrades to a global cap — still a useful brute-force
// brake, and safe against spoofed X-Forwarded-For.
func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
