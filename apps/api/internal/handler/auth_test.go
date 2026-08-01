package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
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
	if body["error"] != "UNAUTHORIZED" {
		t.Fatalf("error = %v, want UNAUTHORIZED", body["error"])
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
