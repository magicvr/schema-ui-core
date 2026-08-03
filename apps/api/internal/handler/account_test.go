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
	// GOAL-006 S5 · the seeded admin holds the list-edit-lifecycle menu grant.
	if got := session.Features["menu_list_edit_lifecycle"]; !got {
		t.Fatalf("admin menu_list_edit_lifecycle = %v, want true", got)
	}
}

// TestAccountsMeFeaturesViewer asserts a viewer resolves the menu feature to
// false (read-only role, no menu grant) through /me (GOAL-006 S5 / V-MENU-02).
func TestAccountsMeFeaturesViewer(t *testing.T) {
	env := newAuthTestEnv(t)
	env.addUser(t, "viewer", "pw", []string{"viewer"})
	token := env.login(t, "viewer", "pw")

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
	if got := session.Features["menu_list_edit_lifecycle"]; got {
		t.Fatalf("viewer menu_list_edit_lifecycle = %v, want false", got)
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
	// S5 parity: the dev session carries the admin menu grant feature.
	features, _ := body["features"].(map[string]any)
	if got := features["menu_list_edit_lifecycle"]; got != true {
		t.Fatalf("dev features.menu_list_edit_lifecycle = %v, want true", got)
	}
	// GOAL-011 (A-004 F-003): the dev session mirrors the admin users/roles
	// permissions and menus so the fallback never silently drops the gates.
	perms, _ := user["permissions"].([]any)
	permSet := map[string]bool{}
	for _, p := range perms {
		if s, ok := p.(string); ok {
			permSet[s] = true
		}
	}
	for _, want := range []string{"users.read", "users.write", "roles.read", "roles.write"} {
		if !permSet[want] {
			t.Fatalf("dev permissions missing %s in %v", want, perms)
		}
	}
	for _, want := range []string{"menu_users", "menu_roles"} {
		if got := features[want]; got != true {
			t.Fatalf("dev features.%s = %v, want true", want, got)
		}
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
