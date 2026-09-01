package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
)

func loginBody(t *testing.T, env *authTestEnv, username, password string) (int, map[string]any) {
	t.Helper()
	return sendJSON(t, env.mux, http.MethodPost, "/api/auth/login",
		`{"username":`+quote(username)+`,"password":`+quote(password)+`}`)
}

func TestAuthLoginSuccess(t *testing.T) {
	env := newAuthTestEnv(t)
	code, body := loginBody(t, env, testSeedUsername, testSeedPassword)
	if code != http.StatusOK {
		t.Fatalf("status = %d, want %d: %v", code, http.StatusOK, body)
	}
	if body["accessToken"] == "" || body["refreshToken"] == "" {
		t.Fatalf("tokens missing in %v", body)
	}
	user, _ := body["user"].(map[string]any)
	if user["id"] != "user-admin" {
		t.Fatalf("user.id = %v, want user-admin", user["id"])
	}
}

func TestAuthLoginBadPassword(t *testing.T) {
	env := newAuthTestEnv(t)
	code, body := loginBody(t, env, testSeedUsername, "wrong")
	if code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", code, http.StatusUnauthorized)
	}
	if body["error"] != "UNAUTHORIZED" {
		t.Fatalf("error = %v, want UNAUTHORIZED", body["error"])
	}
}

func TestAuthLoginUnknownUser(t *testing.T) {
	env := newAuthTestEnv(t)
	code, body := loginBody(t, env, "nobody", "pw")
	if code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", code, http.StatusUnauthorized)
	}
	if body["error"] != "UNAUTHORIZED" {
		t.Fatalf("error = %v, want UNAUTHORIZED (no enumeration)", body["error"])
	}
}

// TestAccountLockLifecycle covers the production login-lock source (GOAL-004
// S4-6 → GOAL-014 D-002 layered model): the login surface never distinguishes
// a locked account from an unknown user — both are a generic 401
// UNAUTHORIZED so an attacker cannot enumerate account existence via
// lock/disable codes. Repeated failures from ONE source open the
// per-(account|source) lock; correct-password logins from that source answer
// the same generic envelope, and the lock expires automatically.
func TestAccountLockLifecycle(t *testing.T) {
	env := newAuthTestEnv(t)

	// 5 consecutive failures: 401 each, no lock yet.
	for range 5 {
		code, body := loginBody(t, env, testSeedUsername, "wrong")
		if code != http.StatusUnauthorized {
			t.Fatalf("failure status = %d, want 401: %v", code, body)
		}
	}
	// The 5th failure opened the SOURCE lock: even the correct password is
	// now a generic 401 UNAUTHORIZED (W7 F-009 — no distinct lock/disable
	// oracle).
	code, body := loginBody(t, env, testSeedUsername, testSeedPassword)
	if code != http.StatusUnauthorized {
		t.Fatalf("locked login status = %d, want 401: %v", code, body)
	}
	if body["error"] != "UNAUTHORIZED" {
		t.Fatalf("error = %v, want UNAUTHORIZED", body["error"])
	}

	// Expiry: move past the lock window and the SAME source recovers.
	u, err := env.authRepository.UserByUsername(testSeedUsername)
	if err != nil {
		t.Fatal(err)
	}
	// httptest.NewRequest defaults RemoteAddr to 192.0.2.1 — the source the
	// failures above were recorded under.
	const httpSourceIP = "192.0.2.1"
	if locked, err := env.authRepository.LoginLockedFor(u.ID, httpSourceIP, time.Now()); err != nil || !locked {
		t.Fatalf("source lock expected open: locked=%v err=%v", locked, err)
	}
	afterWindow := time.Now().Add(auth.LockWindow + time.Minute)
	access, refresh, _, err := env.a.Login(testSeedUsername, testSeedPassword, afterWindow, httpSourceIP)
	if err != nil {
		t.Fatalf("login after expiry: %v", err)
	}
	if access == "" || refresh == "" {
		t.Fatal("tokens missing after lock expiry")
	}
	// A successful login resets every failure counter.
	u2, err := env.authRepository.UserByUsername(testSeedUsername)
	if err != nil {
		t.Fatal(err)
	}
	if u2.FailedLoginCount != 0 || u2.LockedUntil != 0 {
		t.Fatalf("counter not reset: failed=%d lockedUntil=%d", u2.FailedLoginCount, u2.LockedUntil)
	}
}

// GOAL-014 D-002 (W13 F-007): a per-(account|source) lock from repeated
// failures denies THAT client's logins only and — deliberately — does NOT
// revoke any session: forced logout was the weaponizable edge of the old
// global-lock model. The pre-existing refresh token stays live and unrevoked.
// (A source lock never writes users.locked_until, so Refresh is unaffected;
// the global ceiling keeps its separate locked-account rotation denial.)
func TestAccountSourceLockKeepsSessions(t *testing.T) {
	env := newAuthTestEnv(t)
	_, refresh, _, err := env.a.Login(testSeedUsername, testSeedPassword, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	for range 5 {
		loginBody(t, env, testSeedUsername, "wrong")
	}
	// Same-source correct-password login is denied (the pair lock).
	code, _ := loginBody(t, env, testSeedUsername, testSeedPassword)
	if code != http.StatusUnauthorized {
		t.Fatalf("locked-source login = %d, want 401 (W7 F-009 envelope)", code)
	}
	rt, err := env.authRepository.RefreshTokenByHash(auth.HashToken(refresh))
	if err != nil {
		t.Fatal(err)
	}
	if rt.RevokedAt != nil {
		t.Fatal("source lock must NOT revoke the live refresh token (GOAL-014 D-002 §3)")
	}
}

func TestAuthLoginBadBody(t *testing.T) {
	env := newAuthTestEnv(t)
	for _, body := range []string{"", "not json", `{"username":"admin"}`} {
		code, _ := sendJSON(t, env.mux, http.MethodPost, "/api/auth/login", body)
		if code != http.StatusBadRequest {
			t.Fatalf("%q: status = %d, want %d", body, code, http.StatusBadRequest)
		}
	}
}

func refreshBody(t *testing.T, env *authTestEnv, refresh string) (int, map[string]any) {
	t.Helper()
	return sendJSON(t, env.mux, http.MethodPost, "/api/auth/refresh",
		`{"refreshToken":`+quote(refresh)+`}`)
}

func TestAuthRefreshRotates(t *testing.T) {
	env := newAuthTestEnv(t)
	_, login := loginBody(t, env, testSeedUsername, testSeedPassword)
	oldRefresh, _ := login["refreshToken"].(string)

	code, body := refreshBody(t, env, oldRefresh)
	if code != http.StatusOK {
		t.Fatalf("status = %d, want %d: %v", code, http.StatusOK, body)
	}
	newRefresh, _ := body["refreshToken"].(string)
	if newRefresh == "" || newRefresh == oldRefresh {
		t.Fatalf("new refresh = %q, want non-empty and different", newRefresh)
	}
	// Old refresh must now be rejected (rotation revokes it).
	code2, _ := refreshBody(t, env, oldRefresh)
	if code2 != http.StatusUnauthorized {
		t.Fatalf("old refresh status = %d, want %d (revoked)", code2, http.StatusUnauthorized)
	}
}

func TestAuthRefreshInvalid(t *testing.T) {
	env := newAuthTestEnv(t)
	code, body := refreshBody(t, env, "bogus-token")
	if code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", code, http.StatusUnauthorized)
	}
	if body["error"] != "REFRESH_TOKEN_EXPIRED" {
		t.Fatalf("error = %v, want REFRESH_TOKEN_EXPIRED", body["error"])
	}
}

func TestAuthLogoutRevokesRefresh(t *testing.T) {
	env := newAuthTestEnv(t)
	_, login := loginBody(t, env, testSeedUsername, testSeedPassword)
	refresh, _ := login["refreshToken"].(string)

	code, _ := sendJSON(t, env.mux, http.MethodPost, "/api/auth/logout",
		`{"refreshToken":`+quote(refresh)+`}`)
	if code != http.StatusNoContent {
		t.Fatalf("logout status = %d, want %d", code, http.StatusNoContent)
	}
	// Refresh after logout is rejected; logout is idempotent.
	code2, _ := refreshBody(t, env, refresh)
	if code2 != http.StatusUnauthorized {
		t.Fatalf("refresh after logout = %d, want %d", code2, http.StatusUnauthorized)
	}
	code3, _ := sendJSON(t, env.mux, http.MethodPost, "/api/auth/logout",
		`{"refreshToken":`+quote(refresh)+`}`)
	if code3 != http.StatusNoContent {
		t.Fatalf("second logout = %d, want %d (idempotent)", code3, http.StatusNoContent)
	}
}

func TestAuthLogoutUnknownIdempotent(t *testing.T) {
	env := newAuthTestEnv(t)
	code, _ := sendJSON(t, env.mux, http.MethodPost, "/api/auth/logout", `{"refreshToken":"bogus"}`)
	if code != http.StatusNoContent {
		t.Fatalf("logout status = %d, want %d", code, http.StatusNoContent)
	}
}

// TestAuthFlowEndToEnd walks login → /me → refresh → logout as the acceptance
// matrix M1/M3/M4/M6 does.
func TestAuthFlowEndToEnd(t *testing.T) {
	env := newAuthTestEnv(t)

	// Login.
	code, login := loginBody(t, env, testSeedUsername, testSeedPassword)
	if code != http.StatusOK {
		t.Fatalf("login = %d, want 200", code)
	}
	access, _ := login["accessToken"].(string)
	refresh, _ := login["refreshToken"].(string)

	// /me with access token returns identity.
	req := bearer(t, access, http.MethodGet, "/api/accounts/me", "")
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("/me = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	var me map[string]any
	if err := json.NewDecoder(rr.Body).Decode(&me); err != nil {
		t.Fatalf("decode /me: %v", err)
	}
	user, _ := me["user"].(map[string]any)
	if user["id"] != "user-admin" {
		t.Fatalf("/me user.id = %v, want user-admin", user["id"])
	}

	// Refresh.
	code, refreshed := refreshBody(t, env, refresh)
	if code != http.StatusOK {
		t.Fatalf("refresh = %d, want 200", code)
	}
	newAccess, _ := refreshed["accessToken"].(string)
	if newAccess == "" {
		t.Fatalf("refreshed accessToken missing")
	}

	// Old access token is still valid until its TTL, new access works too.
	req = bearer(t, newAccess, http.MethodGet, "/api/accounts/me", "")
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("/me with new access = %d, want 200", rr.Code)
	}

	// Logout revokes the rotated refresh.
	code, _ = sendJSON(t, env.mux, http.MethodPost, "/api/auth/logout",
		`{"refreshToken":`+quote(refreshed["refreshToken"].(string))+`}`)
	if code != http.StatusNoContent {
		t.Fatalf("logout = %d, want 204", code)
	}
}

func TestAuthDevSessionDoesNotBypassLoginEndpoint(t *testing.T) {
	// The dev-session fallback must not make /api/auth/login succeed without a
	// password: login always validates against the store.
	env := newDevSessionTestEnv(t)
	code, body := loginBody(t, env, testSeedUsername, "wrong")
	if code != http.StatusUnauthorized {
		t.Fatalf("login status = %d, want %d", code, http.StatusUnauthorized)
	}
	if body["error"] != "UNAUTHORIZED" {
		t.Fatalf("error = %v, want UNAUTHORIZED", body["error"])
	}
}

// D2 hardening: repeated failed logins from one client IP are rate-limited
// (429 after the configured threshold), while the limiter is per-IP.
// GOAL-004 S4-6: the first 5 consecutive failures trip the account lock
// (401 UNAUTHORIZED after W7 F-009, no lock enumeration) before the 21-attempt per-IP limiter; the rate-limit
// path is exercised with a nonexistent user, which never locks.
func TestLoginRateLimit(t *testing.T) {
	env := newAuthTestEnv(t)

	// A nonexistent user never locks an account: 20 failed attempts are
	// allowed, the 21st is rejected with 429.
	for range 20 {
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login",
			strings.NewReader(`{"username":"no-such-user","password":"wrong-password"}`))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		env.mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusUnauthorized {
			t.Fatalf("attempt status = %d, want 401", rr.Code)
		}
	}
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login",
		strings.NewReader(`{"username":"no-such-user","password":"wrong-password"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusTooManyRequests {
		t.Fatalf("over-limit login status = %d, want 429", rr.Code)
	}
	var out map[string]string
	_ = json.NewDecoder(rr.Body).Decode(&out)
	if out["error"] != "RATE_LIMITED" {
		t.Fatalf("over-limit code = %q, want RATE_LIMITED", out["error"])
	}
	if rr.Header().Get("Retry-After") == "" {
		t.Fatal("RATE_LIMITED missing Retry-After")
	}

	// A correct password is still rejected under lockout (fail-closed).
	// The nonexistent user is rate-limited per-IP but never account-locked.
	req = httptest.NewRequest(http.MethodPost, "/api/auth/login",
		strings.NewReader(`{"username":"no-such-user","password":"test-password"}`))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusTooManyRequests {
		t.Fatalf("correct password under lockout status = %d, want 429", rr.Code)
	}
}

// D-001 P1: behind a trusted reverse proxy (loopback/private peer) the
// X-Real-IP header identifies the real client; it is never trusted from an
// untrusted peer.
func TestLoginClientIPTrustsXRealIPOnlyFromTrustedPeer(t *testing.T) {
	trusted := httptest.NewRequest(http.MethodPost, "/api/auth/login",
		strings.NewReader(`{"username":"admin","password":"x"}`))
	trusted.RemoteAddr = "127.0.0.1:5555"
	trusted.Header.Set("X-Real-IP", "203.0.113.7")
	if got := loginClientIP(trusted); got != "203.0.113.7" {
		t.Fatalf("loopback peer X-Real-IP = %q, want 203.0.113.7", got)
	}

	spoofed := httptest.NewRequest(http.MethodPost, "/api/auth/login",
		strings.NewReader(`{"username":"admin","password":"x"}`))
	spoofed.RemoteAddr = "203.0.113.99:40000"
	spoofed.Header.Set("X-Real-IP", "198.51.100.1")
	if got := loginClientIP(spoofed); got != "203.0.113.99" {
		t.Fatalf("untrusted peer X-Real-IP = %q, want direct peer 203.0.113.99", got)
	}
}
