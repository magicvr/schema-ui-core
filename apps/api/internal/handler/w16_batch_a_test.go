package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// W16-F01: an API-created user starts with must_change_password=1; business
// APIs are gated until the user replaces the initial password, and the forced
// change returns a fresh token pair.
func TestForcedPasswordChangeGateAndReissue(t *testing.T) {
	env := newAuthTestEnv(t)
	admin := adminToken(t, env)

	// Create a user through the API (must_change_password defaults to true).
	req := bearer(t, admin, http.MethodPost, "/api/users",
		`{"username":"forced","name":"Forced","password":"initial-pass","roles":["editor"]}`)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("create status = %d: %s", rr.Code, rr.Body.String())
	}
	var created map[string]any
	_ = json.NewDecoder(rr.Body).Decode(&created)
	id, _ := created["id"].(string)
	if id == "" {
		t.Fatalf("created user missing id: %v", created)
	}
	if created["mustChangePassword"] != true {
		t.Fatalf("created user mustChangePassword = %v, want true", created["mustChangePassword"])
	}

	// Login succeeds and the user snapshot carries the forced flag.
	code, loginOut := sendJSON(t, env.mux, http.MethodPost, "/api/auth/login",
		`{"username":"forced","password":"initial-pass"}`)
	if code != http.StatusOK {
		t.Fatalf("login = %d: %v", code, loginOut)
	}
	token, _ := loginOut["accessToken"].(string)
	if token == "" {
		t.Fatal("login accessToken missing")
	}
	userMap, _ := loginOut["user"].(map[string]any)
	if userMap["mustChangePassword"] != true {
		t.Fatalf("login user mustChangePassword = %v, want true", userMap["mustChangePassword"])
	}

	// The /me surface is allowed so the frontend can resolve the session.
	meProbe := bearer(t, token, http.MethodGet, "/api/accounts/me", "")
	meRR := httptest.NewRecorder()
	env.mux.ServeHTTP(meRR, meProbe)
	if meRR.Code != http.StatusOK {
		t.Fatalf("forced /me status = %d, want 200", meRR.Code)
	}

	// Business API is blocked.
	probe := bearer(t, token, http.MethodGet, "/api/users/"+id, "")
	probeRR := httptest.NewRecorder()
	env.mux.ServeHTTP(probeRR, probe)
	if probeRR.Code != http.StatusForbidden {
		t.Fatalf("gated business status = %d, want 403", probeRR.Code)
	}
	if !strings.Contains(probeRR.Body.String(), "MUST_CHANGE_PASSWORD") {
		t.Fatalf("gated body = %s, want MUST_CHANGE_PASSWORD", probeRR.Body.String())
	}

	// Same-password replacement is rejected (independent audit F-002).
	same := bearer(t, token, http.MethodPost, "/api/account/password",
		`{"currentPassword":"initial-pass","newPassword":"initial-pass"}`)
	sameRR := httptest.NewRecorder()
	env.mux.ServeHTTP(sameRR, same)
	if sameRR.Code != http.StatusBadRequest {
		t.Fatalf("same-password change status = %d, want 400", sameRR.Code)
	}

	// Forced password change returns 200 with a fresh token pair.
	change := bearer(t, token, http.MethodPost, "/api/account/password",
		`{"currentPassword":"initial-pass","newPassword":"new-secret-123"}`)
	changeRR := httptest.NewRecorder()
	env.mux.ServeHTTP(changeRR, change)
	if changeRR.Code != http.StatusOK {
		t.Fatalf("forced change status = %d, want 200: %s", changeRR.Code, changeRR.Body.String())
	}
	var changed tokenResponse
	if err := json.NewDecoder(changeRR.Body).Decode(&changed); err != nil {
		t.Fatalf("decode forced change response: %v", err)
	}
	if changed.AccessToken == "" || changed.RefreshToken == "" {
		t.Fatal("forced change did not return new tokens")
	}
	if changed.User.MustChangePassword {
		t.Fatal("forced change user still has mustChangePassword=true")
	}

	// The fresh token can access business APIs.
	probe2 := bearer(t, changed.AccessToken, http.MethodGet, "/api/users/"+id, "")
	probe2RR := httptest.NewRecorder()
	env.mux.ServeHTTP(probe2RR, probe2)
	if probe2RR.Code != http.StatusOK {
		t.Fatalf("post-change business status = %d, want 200", probe2RR.Code)
	}
}

// W16-F07: revoke-others bumps token_version, revokes every old refresh token,
// and returns a fresh pair for the current caller.
func TestRevokeOthersReissuesTokensAndRevokesOtherSessions(t *testing.T) {
	env := newAuthTestEnv(t)
	// Two independent sign-ins for the same admin account.
	code1, out1 := sendJSON(t, env.mux, http.MethodPost, "/api/auth/login",
		`{"username":"admin","password":"test-password"}`)
	if code1 != http.StatusOK {
		t.Fatalf("first login = %d", code1)
	}
	access1, _ := out1["accessToken"].(string)
	code2, out2 := sendJSON(t, env.mux, http.MethodPost, "/api/auth/login",
		`{"username":"admin","password":"test-password"}`)
	if code2 != http.StatusOK {
		t.Fatalf("second login = %d", code2)
	}
	access2, _ := out2["accessToken"].(string)
	refresh2, _ := out2["refreshToken"].(string)

	// Call revoke-others with the first session.
	req := bearer(t, access1, http.MethodPost, "/api/account/sessions/revoke-others", "")
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("revoke-others status = %d: %s", rr.Code, rr.Body.String())
	}
	var body tokenResponse
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("decode revoke-others response: %v", err)
	}
	if body.AccessToken == "" || body.RefreshToken == "" {
		t.Fatal("revoke-others did not return fresh tokens")
	}

	// The other session's access token is immediately invalid (token_version bump).
	probe := bearer(t, access2, http.MethodGet, "/api/account/profile", "")
	probeRR := httptest.NewRecorder()
	env.mux.ServeHTTP(probeRR, probe)
	if probeRR.Code != http.StatusUnauthorized {
		t.Fatalf("other access token status = %d, want 401", probeRR.Code)
	}

	// The other session's refresh token is revoked.
	refreshCode, _ := sendJSON(t, env.mux, http.MethodPost, "/api/auth/refresh",
		`{"refreshToken":"`+refresh2+`"}`)
	if refreshCode != http.StatusUnauthorized {
		t.Fatalf("other refresh token status = %d, want 401", refreshCode)
	}

	// The current caller's new pair works.
	probe2 := bearer(t, body.AccessToken, http.MethodGet, "/api/account/profile", "")
	probe2RR := httptest.NewRecorder()
	env.mux.ServeHTTP(probe2RR, probe2)
	if probe2RR.Code != http.StatusOK {
		t.Fatalf("new access token status = %d, want 200", probe2RR.Code)
	}
	// The old current refresh token was also revoked; the new one rotates.
}
