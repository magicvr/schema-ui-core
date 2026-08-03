package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
)

// R5 S6 (I-008-003 §6) · records write endpoints append operation-log rows with
// the correct event, actor and record id.
func TestOperationLogRecordsWrites(t *testing.T) {
	env := newAuthTestEnv(t)

	// create
	body := `{"name":"OpLog Co","status":"active","owner":"alice"}`
	req := bearer(t, adminToken(t, env), http.MethodPost, "/api/records", body)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want 201: %v", rr.Code, rr.Body.String())
	}
	var created map[string]any
	if err := json.NewDecoder(rr.Body).Decode(&created); err != nil {
		t.Fatalf("decode create: %v", err)
	}
	recID, _ := created["id"].(string)
	if recID == "" {
		t.Fatalf("create id missing in %v", created)
	}

	// update
	patch := `{"status":"archived"}`
	req = bearer(t, adminToken(t, env), http.MethodPatch, "/api/records/"+recID, patch)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("update status = %d, want 200", rr.Code)
	}

	// delete
	req = bearer(t, adminToken(t, env), http.MethodDelete, "/api/records/"+recID, "")
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d, want 204", rr.Code)
	}

	ops, err := env.st.ListOperations(10)
	if err != nil {
		t.Fatalf("ListOperations: %v", err)
	}
	var recordOps []store.Operation
	for _, op := range ops {
		if strings.HasPrefix(op.Event, "records.") {
			recordOps = append(recordOps, op)
		}
	}
	want := []string{store.EventRecordDelete, store.EventRecordUpdate, store.EventRecordCreate}
	if len(recordOps) != len(want) {
		t.Fatalf("record ops = %d, want %d (create/update/delete)", len(recordOps), len(want))
	}
	for i, ev := range want {
		op := recordOps[i] // newest first: delete, update, create
		if op.Event != ev {
			t.Fatalf("recordOps[%d].event = %q, want %q", i, op.Event, ev)
		}
		if op.ActorID != "user-admin" {
			t.Fatalf("recordOps[%d].actor_id = %q, want user-admin", i, op.ActorID)
		}
		if op.ActorName != "Admin" {
			t.Fatalf("recordOps[%d].actor_name = %q, want Admin", i, op.ActorName)
		}
		if op.RecordID == nil || *op.RecordID != recID {
			t.Fatalf("recordOps[%d].record_id = %v, want %s", i, op.RecordID, recID)
		}
	}
	// create detail carries the record name summary (no secrets).
	createOp := recordOps[2]
	if createOp.Detail == nil || !strings.Contains(*createOp.Detail, "OpLog Co") {
		t.Fatalf("create detail = %v, want name summary", createOp.Detail)
	}
}

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

	// refresh (public route)
	code, _ = sendJSON(t, env.mux, http.MethodPost, "/api/auth/refresh",
		`{"refreshToken":`+quote(refresh)+`}`)
	if code != http.StatusOK {
		t.Fatalf("refresh status = %d, want 200", code)
	}

	// logout (public route)
	code, _ = sendJSON(t, env.mux, http.MethodPost, "/api/auth/logout",
		`{"refreshToken":`+quote(refresh)+`}`)
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
		// I-008-003 §3: every auth event carries the frozen username detail.
		if op.Detail == nil || !strings.Contains(*op.Detail, `"username":"admin"`) {
			t.Fatalf("authOps[%d].detail = %v, want username summary", i, op.Detail)
		}
	}
}

// R5 S6 (I-008-003 §5) · failed writes do not append operation-log rows.
func TestOperationLogNoRowsOnFailedWrite(t *testing.T) {
	env := newAuthTestEnv(t)

	// Invalid create body → 400, no log row.
	req := bearer(t, adminToken(t, env), http.MethodPost, "/api/records", `{"name":""}`)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("invalid create status = %d, want 400", rr.Code)
	}

	// Anonymous write → 401, no log row.
	req = httptest.NewRequest(http.MethodPost, "/api/records",
		strings.NewReader(`{"name":"X","status":"a","owner":"b"}`))
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
	var recordOps []store.Operation
	for _, op := range ops {
		if strings.HasPrefix(op.Event, "records.") {
			recordOps = append(recordOps, op)
		}
	}
	if len(recordOps) != 0 {
		t.Fatalf("record ops after failed writes = %d, want 0", len(recordOps))
	}
}
