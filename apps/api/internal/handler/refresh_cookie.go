package handler

import (
	"net/http"
	"strings"
)

// refreshCookieName is the httpOnly cookie key for the refresh token (W17
// GOAL-018 D-001: mitigates XSS exfiltration risk vs. localStorage).
const refreshCookieName = "refresh_token"

// setRefreshCookie writes the refresh token as an httpOnly cookie (D-001 I-001).
// Attributes: HttpOnly=true, Secure (prod/dev adaptive), SameSite=Lax,
// Path=/api/auth, Max-Age=30 days. The response JSON body STILL contains the
// refreshToken field for non-browser client compatibility (D-001 I-002).
func setRefreshCookie(w http.ResponseWriter, token string, secure bool) {
	cookie := &http.Cookie{
		Name:     refreshCookieName,
		Value:    token,
		Path:     "/api/auth",
		MaxAge:   2592000, // 30 days (matches existing refresh token TTL)
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode, // D-001: Lax balances security + top-level navigation (email links)
	}
	http.SetCookie(w, cookie)
}

// clearRefreshCookie removes the refresh token cookie (logout).
func clearRefreshCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     refreshCookieName,
		Value:    "",
		Path:     "/api/auth",
		MaxAge:   -1, // immediate expiry
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)
}

// extractRefreshToken reads the refresh token from three sources (D-001 I-002):
// 1. Cookie (priority 1: browser SPA with httpOnly cookie)
// 2. X-Refresh-Token header (priority 2: non-browser clients)
// 3. JSON body refreshToken field (priority 3: legacy/test compatibility)
// Returns empty string if none found. The caller decodes the JSON body if needed.
func extractRefreshToken(r *http.Request) string {
	// Priority 1: Cookie (browser with httpOnly)
	if cookie, err := r.Cookie(refreshCookieName); err == nil && cookie.Value != "" {
		return cookie.Value
	}
	// Priority 2: Header (non-browser clients, e.g. mobile SDKs)
	if header := r.Header.Get("X-Refresh-Token"); header != "" {
		return header
	}
	// Priority 3: Body (legacy/test — caller must decode JSON if this path is taken)
	return ""
}

// isDevMode detects development environment (HTTP localhost) to disable Secure
// cookie attribute (D-001 I-004: browsers reject Secure cookies over HTTP).
func isDevMode(r *http.Request) bool {
	// Check if the request came over HTTP (not HTTPS)
	if r.TLS == nil && r.Header.Get("X-Forwarded-Proto") != "https" {
		host := r.Host
		if host == "" {
			host = r.Header.Get("Host")
		}
		// localhost or 127.0.0.1 over HTTP = dev mode
		if strings.HasPrefix(host, "localhost:") || strings.HasPrefix(host, "127.0.0.1:") {
			return true
		}
	}
	return false
}
