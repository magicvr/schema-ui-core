package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
)

// users/roles write operation-log events are covered by TestUsersOperationLogEvents
// and TestRolesOperationLogEvents (GOAL-011 S2). records.* events are historical
// only after 0006 retirement (GOAL-011 S3) — no records API remains to write them.

// R5 S6 (I-008-003 §6) · auth login/logout/refresh append operation-log rows.
func TestOperationLogAuthEvents(t *testing.T) {
	env := newAuthTestEnv(t)

	// login (public route, via env.mux)
	code, out := sendJSON(t, env.mux, http.MethodPost, "/api/auth/login",
		`{"username":"admin","password":"test-password"}`)
	if code != http.StatusOK {
		t.Fatalf("login status = %d, want 200: %v", code, out)
	}
	refresh, _ := out["refreshToken"].(string)
	if refresh == "" {
		t.Fatalf("refreshToken missing")
	}

	// refresh (public route): rotation issues a new refresh token, the old one
	// is revoked. Logout below uses the ROTATED token so the auth.logout row
	// records a first (non-idempotent) logout path (R-015).
	code, out = sendJSON(t, env.mux, http.MethodPost, "/api/auth/refresh",
		`{"refreshToken":`+quote(refresh)+`}`)
	if code != http.StatusOK {
		t.Fatalf("refresh status = %d, want 200", code)
	}
	rotated, _ := out["refreshToken"].(string)
	if rotated == "" || rotated == refresh {
		t.Fatalf("rotated refreshToken missing or unchanged")
	}

	// logout (public route) with the rotated, still-valid token.
	code, _ = sendJSON(t, env.mux, http.MethodPost, "/api/auth/logout",
		`{"refreshToken":`+quote(rotated)+`}`)
	if code != http.StatusNoContent {
		t.Fatalf("logout status = %d, want 204", code)
	}

	ops, err := env.st.ListOperations(10)
	if err != nil {
		t.Fatalf("ListOperations: %v", err)
	}
	var authOps []store.Operation
	for _, op := range ops {
		if strings.HasPrefix(op.Event, "auth.") {
			authOps = append(authOps, op)
		}
	}
	want := []string{store.EventAuthLogout, store.EventAuthRefresh, store.EventAuthLogin}
	if len(authOps) != len(want) {
		t.Fatalf("auth ops = %d, want %d (login/refresh/logout)", len(authOps), len(want))
	}
	for i, ev := range want {
		op := authOps[i]
		if op.Event != ev {
			t.Fatalf("authOps[%d].event = %q, want %q", i, op.Event, ev)
		}
		if op.ActorID != "user-admin" {
			t.Fatalf("authOps[%d].actor_id = %q, want user-admin", i, op.ActorID)
		}
		// I-008-003 §3: every auth event carries the frozen username detail —
		// exact JSON shape, username only, no sensitive fields (R-015).
		if op.Detail == nil {
			t.Fatalf("authOps[%d].detail = nil, want username summary", i)
		}
		var d struct {
			Username string `json:"username"`
		}
		if err := json.Unmarshal([]byte(*op.Detail), &d); err != nil {
			t.Fatalf("authOps[%d].detail %q not JSON: %v", i, *op.Detail, err)
		}
		if d.Username != "admin" {
			t.Fatalf("authOps[%d].detail.username = %q, want admin", i, d.Username)
		}
		var extra map[string]any
		if err := json.Unmarshal([]byte(*op.Detail), &extra); err == nil && len(extra) != 1 {
			t.Fatalf("authOps[%d].detail = %v, want exactly {username} (no token/password/secret)", i, extra)
		}
	}
}

// R5 S6 (I-008-003 §5) · failed writes do not append operation-log rows.
func TestOperationLogNoRowsOnFailedWrite(t *testing.T) {
	env := newAuthTestEnv(t)

	// Invalid create body → 400, no log row (GOAL-011: users resource).
	req := bearer(t, adminToken(t, env), http.MethodPost, "/api/users", `{"name":"NoPass"}`)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("invalid create status = %d, want 400", rr.Code)
	}

	// Anonymous write → 401, no log row.
	req = httptest.NewRequest(http.MethodPost, "/api/users",
		strings.NewReader(`{"username":"x","name":"X","password":"y12345"}`))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("anon create status = %d, want 401", rr.Code)
	}

	ops, err := env.st.ListOperations(10)
	if err != nil {
		t.Fatalf("ListOperations: %v", err)
	}
	var userOps []store.Operation
	for _, op := range ops {
		if strings.HasPrefix(op.Event, "users.") {
			userOps = append(userOps, op)
		}
	}
	if len(userOps) != 0 {
		t.Fatalf("user ops after failed writes = %d, want 0", len(userOps))
	}
}
