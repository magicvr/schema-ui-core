package handler

// Client source-IP helpers for per-client rate-limit key construction
// (VP-027 / workspace-027 GOAL-003 D-001, R2): the trusted-reverse-proxy
// allow-list and the loginClientIP/clientIP resolution moved here from the
// legacy rate_limit.go when the limiter itself was ported to
// internal/ratelimit. These helpers are handler-layer key conventions (D-002
// §2) and deliberately live outside the provider package: they parse HTTP
// request metadata, which is not a provider concern.

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

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

// LoginClientIP returns the client identity used for per-client rate
// limiting: the X-Real-IP header when the direct peer is a configured trusted
// reverse proxy, otherwise the direct peer address. The header is never
// trusted from an untrusted peer: a client-facing server must not be
// spoofable by setting X-Real-IP itself (D-001 P1; W7 F-008 — trust is
// explicit CIDRs, not all private addresses).
func LoginClientIP(r *http.Request) string {
	peer := clientIP(r)
	if trustedReverseProxy(peer) {
		if real := r.Header.Get("X-Real-IP"); strings.TrimSpace(real) != "" {
			return strings.TrimSpace(real)
		}
	}
	return peer
}

func loginClientIP(r *http.Request) string {
	return LoginClientIP(r)
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
