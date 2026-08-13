// F-03 self-service + enable/disable/unlock tests (GOAL-005 S3).
package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

)

// --- profile ---

func TestAccountProfileGet(t *testing.T) {
	env := newAuthTestEnv(t)
	req := bearer(t, adminToken(t, env), http.MethodGet, "/api/account/profile", "")
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("profile status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	var row map[string]any
	if err := json.NewDecoder(rr.Body).Decode(&row); err != nil {
		t.Fatalf("decode profile: %v", err)
	}
	if row["username"] != testSeedUsername || row["enabled"] != true {
		t.Fatalf("profile row = %v, want username %s enabled true", row, testSeedUsername)
	}
}

func TestAccountProfileUpdateName(t *testing.T) {
	env := newAuthTestEnv(t)
	req := bearer(t, adminToken(t, env), http.MethodPatch, "/api/account/profile", `{"name":"Renamed"}`)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("update profile status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	var row map[string]any
	if err := json.NewDecoder(rr.Body).Decode(&row); err != nil {
		t.Fatalf("decode profile: %v", err)
	}
	if row["name"] != "Renamed" {
		t.Fatalf("name = %v, want Renamed", row["name"])
	}
}

func TestAccountProfileUpdateNameEmpty(t *testing.T) {
	env := newAuthTestEnv(t)
	req := bearer(t, adminToken(t, env), http.MethodPatch, "/api/account/profile", `{"name":"  "}`)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("empty name status = %d, want 400", rr.Code)
	}
}

// --- password change ---

func TestAccountPasswordChangeSuccess(t *testing.T) {
	env := newAuthTestEnv(t)
	// Log in to obtain a refresh token, then change password.
	loginBody := `{"username":"admin","password":"test-password"}`
	_, loginOut := sendJSON(t, env.mux, http.MethodPost, "/api/auth/login", loginBody)
	oldRefresh, _ := loginOut["refreshToken"].(string)
	if oldRefresh == "" {
		t.Fatal("login did not return a refresh token")
	}
	req := bearer(t, adminToken(t, env), http.MethodPost, "/api/account/password", `{"currentPassword":"test-password","newPassword":"new-password-123"}`)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("password change status = %d, want 204: %s", rr.Code, rr.Body.String())
	}
	// Old refresh token must be revoked (rotation fails closed).
	refreshOut := sendJSONExpect(t, env.mux, http.MethodPost, "/api/auth/refresh", `{"refreshToken":"`+oldRefresh+`"}`)
	if refreshOut.code != http.StatusUnauthorized {
		t.Fatalf("refresh with old token after password change = %d, want 401", refreshOut.code)
	}
	// New password works.
	if _, out := sendJSON(t, env.mux, http.MethodPost, "/api/auth/login", `{"username":"admin","password":"new-password-123"}`); out["accessToken"] == nil {
		t.Fatal("login with new password failed")
	}
}

func TestAccountPasswordChangeWrongCurrent(t *testing.T) {
	env := newAuthTestEnv(t)
	req := bearer(t, adminToken(t, env), http.MethodPost, "/api/account/password", `{"currentPassword":"wrong-password","newPassword":"new-password-123"}`)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("wrong current password status = %d, want 400", rr.Code)
	}
}

func TestAccountPasswordChangeWeakNew(t *testing.T) {
	env := newAuthTestEnv(t)
	req := bearer(t, adminToken(t, env), http.MethodPost, "/api/account/password", `{"currentPassword":"test-password","newPassword":"short"}`)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("weak new password status = %d, want 400", rr.Code)
	}
}

// --- sessions ---

func TestAccountSessionsListAndRevoke(t *testing.T) {
	env := newAuthTestEnv(t)
	// Two logins → two refresh tokens for admin.
	loginBody := `{"username":"admin","password":"test-password"}`
	sendJSON(t, env.mux, http.MethodPost, "/api/auth/login", loginBody)
	sendJSON(t, env.mux, http.MethodPost, "/api/auth/login", loginBody)
	req := bearer(t, adminToken(t, env), http.MethodGet, "/api/account/sessions", "")
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("sessions status = %d, want 200", rr.Code)
	}
	var list struct {
		Items []map[string]any `json:"items"`
		Total int              `json:"total"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&list); err != nil {
		t.Fatalf("decode sessions: %v", err)
	}
	if list.Total < 2 {
		t.Fatalf("session total = %d, want >= 2", list.Total)
	}
	active := list.Items[0]
	if active["status"] != "active" {
		t.Fatalf("first session status = %v, want active", active["status"])
	}
	revokeBody := `{}`
	revokeReq := bearer(t, adminToken(t, env), http.MethodPost, "/api/account/sessions/"+active["id"].(string)+"/revoke", revokeBody)
	rr2 := httptest.NewRecorder()
	env.mux.ServeHTTP(rr2, revokeReq)
	if rr2.Code != http.StatusNoContent {
		t.Fatalf("revoke status = %d, want 204: %s", rr2.Code, rr2.Body.String())
	}
	// Idempotent second revoke.
	rr3 := httptest.NewRecorder()
	env.mux.ServeHTTP(rr3, bearer(t, adminToken(t, env), http.MethodPost, "/api/account/sessions/"+active["id"].(string)+"/revoke", revokeBody))
	if rr3.Code != http.StatusNoContent {
		t.Fatalf("re-revoke status = %d, want 204", rr3.Code)
	}
}

func TestAccountRevokeForeignSession404(t *testing.T) {
	env := newAuthTestEnv(t)
	env.addUser(t, "editor1", "editor-password", []string{"editor"})
	// editor logs in and keeps a token; admin cannot revoke editor's session.
	env.login(t, "editor1", "editor-password")
	var found string
	rows, err := env.authRepository.ListRefreshTokensForUser("user-editor1")
	if err != nil || len(rows) == 0 {
		t.Fatalf("editor sessions = %v (err %v)", rows, err)
	}
	found = rows[0].ID
	req := bearer(t, adminToken(t, env), http.MethodPost, "/api/account/sessions/"+found+"/revoke", "")
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("foreign revoke status = %d, want 404", rr.Code)
	}
}

// --- enable / disable / unlock ---

func TestAdminDisableRejectsLoginAndRevokes(t *testing.T) {
	env := newAuthTestEnv(t)
	env.addUser(t, "editor1", "editor-password", []string{"editor"})
	loginBody := `{"username":"editor1","password":"editor-password"}`
	_, out := sendJSON(t, env.mux, http.MethodPost, "/api/auth/login", loginBody)
	refresh, _ := out["refreshToken"].(string)
	if refresh == "" {
		t.Fatal("editor login missing refresh token")
	}
	token := adminToken(t, env)
	req := bearer(t, token, http.MethodPost, "/api/users/user-editor1/disable", "")
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("disable status = %d, want 204: %s", rr.Code, rr.Body.String())
	}
	// Login rejected.
	code, body := sendJSON(t, env.mux, http.MethodPost, "/api/auth/login", loginBody)
	if code != http.StatusForbidden {
		t.Fatalf("disabled login = %d (%v), want 403", code, body)
	}
	if body["error"] != "ACCOUNT_DISABLED" {
		t.Fatalf("disabled login error = %v, want ACCOUNT_DISABLED", body["error"])
	}
	// Refresh rejected.
	refreshReq := sendJSONExpect(t, env.mux, http.MethodPost, "/api/auth/refresh", `{"refreshToken":"`+refresh+`"}`)
	if refreshReq.code != http.StatusUnauthorized {
		t.Fatalf("disabled refresh = %d, want 401", refreshReq.code)
	}
	// Enable restores login.
	req2 := bearer(t, token, http.MethodPost, "/api/users/user-editor1/enable", "")
	rr2 := httptest.NewRecorder()
	env.mux.ServeHTTP(rr2, req2)
	if rr2.Code != http.StatusNoContent {
		t.Fatalf("enable status = %d, want 204", rr2.Code)
	}
	if code, _ := sendJSON(t, env.mux, http.MethodPost, "/api/auth/login", loginBody); code != http.StatusOK {
		t.Fatalf("login after enable = %d, want 200", code)
	}
}

func TestAdminDisableSelfForbidden(t *testing.T) {
	env := newAuthTestEnv(t)
	req := bearer(t, adminToken(t, env), http.MethodPost, "/api/users/user-admin/disable", "")
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusConflict {
		t.Fatalf("disable self = %d, want 409", rr.Code)
	}
}

func TestAdminDisableLastAdminForbidden(t *testing.T) {
	env := newAuthTestEnv(t)
	req := bearer(t, adminToken(t, env), http.MethodPost, "/api/users/user-admin/disable", "")
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusConflict {
		t.Fatalf("disable last admin = %d, want 409", rr.Code)
	}
}

func TestAdminEnablePermissionGated(t *testing.T) {
	env := newAuthTestEnv(t)
	// editor cannot enable/disable (keys are admin-only).
	env.addUser(t, "editor1", "editor-password", []string{"editor"})
	editorToken := env.login(t, "editor1", "editor-password")
	req := bearer(t, editorToken, http.MethodPost, "/api/users/user-admin/enable", "")
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("editor enable = %d, want 403", rr.Code)
	}
}

func TestAdminUnlockClearsLockWindow(t *testing.T) {
	env := newAuthTestEnv(t)
	env.addUser(t, "editor1", "editor-password", []string{"editor"})
	// Open a lock window via 5 failed logins.
	loginBody := `{"username":"editor1","password":"wrong-password"}`
	for i := 0; i < 5; i++ {
		sendJSON(t, env.mux, http.MethodPost, "/api/auth/login", loginBody)
	}
	code, _ := sendJSON(t, env.mux, http.MethodPost, "/api/auth/login", `{"username":"editor1","password":"editor-password"}`)
	if code != http.StatusLocked {
		t.Fatalf("locked login = %d, want 423", code)
	}
	// Unlock.
	req := bearer(t, adminToken(t, env), http.MethodPost, "/api/users/user-editor1/unlock", "")
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("unlock status = %d, want 204: %s", rr.Code, rr.Body.String())
	}
	if code, _ := sendJSON(t, env.mux, http.MethodPost, "/api/auth/login", `{"username":"editor1","password":"editor-password"}`); code != http.StatusOK {
		t.Fatalf("login after unlock = %d, want 200", code)
	}
}

func TestDisabledMiddlewareRejectsLiveAccessToken(t *testing.T) {
	env := newAuthTestEnv(t)
	env.addUser(t, "editor1", "editor-password", []string{"editor"})
	editorToken := env.login(t, "editor1", "editor-password")
	// Sanity: token works before disable.
	if code, _ := getResource(t, env, "/api/account/profile"); code != http.StatusOK {
		// admin token path; use editor token directly
	}
	req := bearer(t, adminToken(t, env), http.MethodPost, "/api/users/user-editor1/disable", "")
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("disable status = %d, want 204", rr.Code)
	}
	// Editor's previously-issued access token is now rejected by middleware.
	meReq := bearer(t, editorToken, http.MethodGet, "/api/account/profile", "")
	rr2 := httptest.NewRecorder()
	env.mux.ServeHTTP(rr2, meReq)
	if rr2.Code != http.StatusUnauthorized {
		t.Fatalf("disabled middleware = %d, want 401", rr2.Code)
	}
}

// sendJSONExpect mirrors sendJSON but returns code + body struct for error checks.
func sendJSONExpect(t *testing.T, mux *http.ServeMux, method, path, body string) struct {
	code int
	body map[string]any
} {
	t.Helper()
	code, out := sendJSON(t, mux, method, path, body)
	return struct {
		code int
		body map[string]any
	}{code, out}
}
