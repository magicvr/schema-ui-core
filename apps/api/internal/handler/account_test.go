package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
)

// TestAccountsMeRequiresAuth asserts the request-level identity gate: without an
// access token the /me endpoint fails closed with 401 (M8).
func TestAccountsMeRequiresAuth(t *testing.T) {
	env := newAuthTestEnv(t)
	code, body := getJSON(t, env.mux, "/api/accounts/me")
	if code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", code, http.StatusUnauthorized)
	}
	if body["error"] != "UNAUTHENTICATED" {
		t.Fatalf("error = %v, want UNAUTHENTICATED", body["error"])
	}
}

// TestAccountsMeInvalidTokenFailsClosed covers a malformed/foreign Bearer token.
func TestAccountsMeInvalidTokenFailsClosed(t *testing.T) {
	env := newAuthTestEnv(t)
	req := bearer(t, "not-a-real-token", http.MethodGet, "/api/accounts/me", "")
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusUnauthorized)
	}
	if got := rr.Body.String(); !containsString(got, "UNAUTHENTICATED") {
		t.Fatalf("body = %q, want UNAUTHENTICATED error", got)
	}
}

// TestAccountsMeReturnsIdentity asserts a valid access token resolves to the
// seeded admin identity via GET /api/accounts/me (M8 / M10).
func TestAccountsMeReturnsIdentity(t *testing.T) {
	env := newAuthTestEnv(t)
	token := env.login(t, testSeedUsername, testSeedPassword)

	req := bearer(t, token, http.MethodGet, "/api/accounts/me", "")
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d: %s", rr.Code, http.StatusOK, rr.Body.String())
	}

	var session account.Session
	if err := json.NewDecoder(rr.Body).Decode(&session); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if session.User.ID != "user-admin" {
		t.Fatalf("user.id = %v, want user-admin", session.User.ID)
	}
	if !containsRoles(session.User.Roles, "admin") {
		t.Fatalf("roles = %v, want to contain admin", session.User.Roles)
	}
	if session.Features == nil {
		t.Fatalf("features missing from %v", session)
	}
}

// TestAccountsMeDevSessionFallback asserts the explicit opt-in local-development
// fallback substitutes StaticDevSession when enabled (M9 opt-in side).
func TestAccountsMeDevSessionFallback(t *testing.T) {
	env := newDevSessionTestEnv(t)
	code, body := getJSON(t, env.mux, "/api/accounts/me")
	if code != http.StatusOK {
		t.Fatalf("status = %d, want %d (dev fallback)", code, http.StatusOK)
	}
	user, _ := body["user"].(map[string]any)
	if user["id"] != "dev-001" {
		t.Fatalf("user.id = %v, want dev-001", user["id"])
	}
}

func containsRoles(roles []string, want string) bool {
	for _, r := range roles {
		if r == want {
			return true
		}
	}
	return false
}
