package handler

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// loginRateLimiter is a small in-memory sliding-window limiter for failed
// login attempts (D2, hardened D-001 P1): per client identity (IP + username)
// at most max failures per window. It is process-local and best-effort — it
// does not protect against distributed attacks, but it stops trivial online
// brute force on a single instance.
type loginRateLimiter struct {
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

func newLoginRateLimiter(window time.Duration, max, capacity int) *loginRateLimiter {
	if capacity <= 0 {
		capacity = 1 << 16
	}
	return &loginRateLimiter{
		window:   window,
		max:      max,
		attempts: make(map[string][]time.Time),
		capacity: capacity,
	}
}

// allow reports whether the key may attempt a login now. It never creates a
// new map entry (W4 P0-1): only record() registers a key, so the capacity
// eviction in record() is actually reachable and the map stays bounded. An
// absent key (no failures yet) is always allowed. For an existing key, stale
// window entries are pruned in place before the limit check.
func (l *loginRateLimiter) allow(key string, now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	cutoff := now.Add(-l.window)
	list, exists := l.attempts[key]
	if !exists {
		return true
	}
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

// retryAfterSeconds is the remaining window after the oldest in-window
// failure (W15-F10). Zero when the key is allowed.
func (l *loginRateLimiter) retryAfterSeconds(key string, now time.Time) int {
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
	remain := oldest.Add(l.window).Sub(now)
	if remain <= 0 {
		return 1
	}
	return int(remain.Round(time.Second) / time.Second)
}

// record registers one failed attempt for the key, creating the map entry if
// needed. Bounded: when the map exceeds capacity the oldest key is evicted, so
// an attacker cannot exhaust memory by spraying distinct client identities.
func (l *loginRateLimiter) record(key string, now time.Time) {
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

// clear drops every failure for the key. Called after a successful login so a
// legitimate client never accumulates a poisoned bucket (D-001 P1).
func (l *loginRateLimiter) clear(key string) {
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

// trustedProxyCIDRs is the explicit reverse-proxy allow-list (W7 F-008):
// X-Real-IP is trusted ONLY from a direct peer inside one of these CIDRs.
// Defaults to loopback alone; the composition root installs the configured
// list (Config.HTTPTrustedProxies) at startup. Empty/loopback-only is the
// fail-safe: an operator who wants per-client rate limiting behind a proxy
// must explicitly list that proxy's network.
var trustedProxyCIDRs = []*net.IPNet{
	mustCIDR("127.0.0.1/8"),
}

// SetTrustedProxyCIDRs replaces the trusted-proxy allow-list from config.
// Invalid CIDR strings are a startup error (fail-closed) rather than being
// silently ignored.
func SetTrustedProxyCIDRs(cidrs []string) error {
	nets := make([]*net.IPNet, 0, len(cidrs)+1)
	addLoopback := true
	for _, raw := range cidrs {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}
		_, ipnet, err := net.ParseCIDR(raw)
		if err != nil {
			return fmt.Errorf("invalid trusted proxy CIDR %q: %w", raw, err)
		}
		nets = append(nets, ipnet)
		if ipnet.Contains(net.ParseIP("127.0.0.1")) {
			addLoopback = false
		}
	}
	if addLoopback {
		nets = append(nets, mustCIDR("127.0.0.1/8"))
	}
	trustedProxyCIDRs = nets
	return nil
}

func mustCIDR(cidr string) *net.IPNet {
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		panic("handler: invalid hardcoded CIDR " + cidr + ": " + err.Error())
	}
	return ipnet
}

// loginClientIP returns the client identity used for login rate limiting: the
// X-Real-IP header when the direct peer is a configured trusted reverse proxy,
// otherwise the direct peer address. The header is never trusted from an
// untrusted peer: a client-facing server must not be spoofable by setting
// X-Real-IP itself (D-001 P1; W7 F-008 — trust is explicit CIDRs, not all
// private addresses).
func loginClientIP(r *http.Request) string {
	peer := clientIP(r)
	if trustedReverseProxy(peer) {
		if real := r.Header.Get("X-Real-IP"); strings.TrimSpace(real) != "" {
			return strings.TrimSpace(real)
		}
	}
	return peer
}

// trustedReverseProxy reports whether the direct peer is a configured trusted
// reverse proxy (an explicit CIDR from SetTrustedProxyCIDRs, defaulting to
// loopback only).
func trustedReverseProxy(peer string) bool {
	ip := net.ParseIP(peer)
	if ip == nil {
		return false
	}
	for _, ipnet := range trustedProxyCIDRs {
		if ipnet.Contains(ip) {
			return true
		}
	}
	return false
}

// clientIP returns the direct peer IP; proxy-forwarded headers are not trusted
// from an unknown peer (they are client-controlled), so behind a reverse proxy
// all clients share the proxy IP and the limiter degrades to a global cap —
// still a useful brute-force brake, and safe against spoofed X-Forwarded-For.
func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
