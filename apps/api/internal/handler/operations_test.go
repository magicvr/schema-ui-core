package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
	"github.com/magicvr/schema-ui-core/apps/api/internal/requestid"
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

	ops, err := env.operations.ListOperations(10)
	if err != nil {
		t.Fatalf("ListOperations: %v", err)
	}
	var authOps []operationlog.Operation
	for _, op := range ops {
		if strings.HasPrefix(op.Event, "auth.") {
			authOps = append(authOps, op)
		}
	}
	want := []string{operationlog.EventAuthLogout, operationlog.EventAuthRefresh, operationlog.EventAuthLogin}
	if len(authOps) != len(want) {
		t.Fatalf("auth ops = %d, want %d (login/refresh/logout)", len(authOps), len(want))
	}
	// W25 (A-001 响应 F-006, self): the historical ORDER assertion relied on
	// single-connection serialization of same-millisecond writes. Since the
	// sqlite store pools connections (2026-08-23), same-tick commits may land
	// in any order and the DESC read tie-break falls to the random id suffix —
	// ordering is no longer a contract. Assert the event SET plus every
	// per-operation property instead of positional order.
	seen := map[string]bool{}
	for _, ev := range want {
		seen[ev] = false
	}
	for i, op := range authOps {
		if _, known := seen[op.Event]; !known {
			t.Fatalf("authOps[%d].event = %q, want one of login/refresh/logout", i, op.Event)
		}
		seen[op.Event] = true
		if op.ActorID != "user-admin" {
			t.Fatalf("authOps[%d].actor_id = %q, want user-admin", i, op.ActorID)
		}
		// R2: every auth event carries the versioned audit detail with username
		// only; credentials remain excluded/redacted.
		if op.Detail == nil {
			t.Fatalf("authOps[%d].detail = nil, want username summary", i)
		}
		detail, err := operationlog.ParseDetail(*op.Detail)
		if err != nil {
			t.Fatalf("authOps[%d].detail %q not R2 envelope: %v", i, *op.Detail, err)
		}
		if detail.After["username"] != "admin" {
			t.Fatalf("authOps[%d].detail.after.username = %v, want admin", i, detail.After["username"])
		}
		if op.SessionID == "" {
			t.Fatalf("authOps[%d].session_id empty", i)
		}
		for _, forbidden := range []string{"password", "accessToken", "refreshToken", "secret"} {
			if strings.Contains(*op.Detail, forbidden) {
				t.Fatalf("authOps[%d].detail contains sensitive key %q: %s", i, forbidden, *op.Detail)
			}
		}
	}
}

func TestR1CorrelationIDPersistsOnAuthOperation(t *testing.T) {
	env := newAuthTestEnv(t)
	h := requestid.Middleware(env.mux)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{"username":"admin","password":"test-password"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(requestid.HeaderName, "r1-auth-001")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("login status = %d: %s", rr.Code, rr.Body.String())
	}
	ops, _, err := env.operations.ListOperationsFiltered(operationlog.OperationFilter{
		Event: operationlog.EventAuthLogin, Sort: "createdAt", Order: "desc", Page: 1, PageSize: 10,
	})
	if err != nil {
		t.Fatalf("list auth operations: %v", err)
	}
	if len(ops) != 1 || ops[0].CorrelationID != "r1-auth-001" {
		t.Fatalf("auth operation correlation = %+v, want r1-auth-001", ops)
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
		strings.NewReader(`{"username":"x","name":"X","password":"y123456"}`))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("anon create status = %d, want 401", rr.Code)
	}

	ops, err := env.operations.ListOperations(10)
	if err != nil {
		t.Fatalf("ListOperations: %v", err)
	}
	var userOps []operationlog.Operation
	for _, op := range ops {
		if strings.HasPrefix(op.Event, "users.") {
			userOps = append(userOps, op)
		}
	}
	if len(userOps) != 0 {
		t.Fatalf("user ops after failed writes = %d, want 0", len(userOps))
	}
}

// R4 C3.4 / FR-005: operationlog append failure must not flip a successful
// business write (best-effort contract). Injecting a forced store error, a
// users.create / roles.create write still succeeds.
func TestOperationLogFailurePreservesBusinessSuccess(t *testing.T) {
	env := newAuthTestEnv(t)
	env.operations.SetOperationLogError(errors.New("forced operation log failure"))
	token := adminToken(t, env)

	// users.create → 201 despite log failure.
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, token, http.MethodPost, "/api/users",
		`{"username":"logfail","name":"Log Fail","password":"passw0rd-ok"}`))
	if rr.Code != http.StatusCreated {
		t.Fatalf("create user with log failure = %d, want 201", rr.Code)
	}

	// roles.create → 201 despite log failure.
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, bearer(t, token, http.MethodPost, "/api/roles",
		`{"key":"auditor","name":"Auditor"}`))
	if rr.Code != http.StatusCreated {
		t.Fatalf("create role with log failure = %d, want 201", rr.Code)
	}
}

// W14 F-03 (GOAL-016): structured filters (event/actor/date range) and CSV
// export on the activity log.
func TestOperationLogStructuredFiltersAndExport(t *testing.T) {
	env := newAuthTestEnv(t)
	now := time.Now().UTC()
	for _, op := range []operationlog.Operation{
		{ID: "op-filter-1", Event: operationlog.EventAuthLogin, ActorID: "user-admin", ActorName: "Admin", CorrelationID: "r2-read-001", CreatedAt: now.Add(-2 * time.Hour)},
		{ID: "op-filter-2", Event: operationlog.EventUserCreate, ActorID: "user-admin", ActorName: "Admin", CorrelationID: "r2-read-002", SessionID: "sess-filter-2", CreatedAt: now.Add(-time.Hour)},
		{ID: "op-filter-3", Event: operationlog.EventAuthLogout, ActorID: "user-editor", ActorName: "Editor", CreatedAt: now},
	} {
		if err := env.operations.RecordOperation(op); err != nil {
			t.Fatalf("record operation %s: %v", op.ID, err)
		}
	}

	token := adminToken(t, env)
	code, body := getResourceAs(t, env, token, "/api/operations?event=users.create")
	if code != http.StatusOK || body["total"] != float64(1) {
		t.Fatalf("structured list = %d %v, want total 1", code, body)
	}
	items, _ := body["items"].([]any)
	if len(items) != 1 || items[0].(map[string]any)["correlationId"] != "r2-read-002" || items[0].(map[string]any)["sessionId"] != "sess-filter-2" {
		t.Fatalf("operation list correlation/session = %v, want r2-read-002 / sess-filter-2", items)
	}

	code, body = getResourceAs(t, env, token, "/api/operations/op-filter-1")
	if code != http.StatusOK || body["correlationId"] != "r2-read-001" {
		t.Fatalf("operation detail correlation = %d %v, want r2-read-001", code, body)
	}

	code, body = getResourceAs(t, env, token, "/api/operations?actorName=Editor")
	if code != http.StatusOK || body["total"] != float64(1) {
		t.Fatalf("actor filter list = %d %v, want total 1", code, body)
	}

	// Invalid date filter must surface as 400 (DomainError path in list).
	code, body = getResourceAs(t, env, token, "/api/operations?from=not-a-date")
	if code != http.StatusBadRequest || body["error"] != "INVALID_DATE_FILTER" {
		t.Fatalf("invalid date filter = %d %v, want 400 INVALID_DATE_FILTER", code, body)
	}

	// CSV export applies the same filters and is attachment-disposed.
	req := bearer(t, token, http.MethodGet, "/api/operations/export?event=users.create", "")
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("export = %d: %s", rr.Code, rr.Body.String())
	}
	if ct := rr.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/csv") {
		t.Fatalf("export content-type = %q, want text/csv", ct)
	}
	if rr.Header().Get("Content-Disposition") == "" {
		t.Fatal("export missing Content-Disposition")
	}
	if !strings.Contains(rr.Body.String(), "users.create") {
		t.Fatalf("export body missing users.create row: %q", rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "correlationId") || !strings.Contains(rr.Body.String(), "r2-read-002") {
		t.Fatalf("export missing correlation column/value: %q", rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "sessionId") || !strings.Contains(rr.Body.String(), "sess-filter-2") {
		t.Fatalf("export missing session column/value: %q", rr.Body.String())
	}
}

func TestR2CorrelationIDPersistsOnUsersOperation(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)
	h := requestid.Middleware(env.mux)
	req := bearer(t, token, http.MethodPost, "/api/users", `{"username":"r2-user","name":"R2 User","password":"passw0rd-ok"}`)
	req.Header.Set(requestid.HeaderName, "r2-user-001")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("create status = %d: %s", rr.Code, rr.Body.String())
	}
	ops, _, err := env.operations.ListOperationsFiltered(operationlog.OperationFilter{
		Event: operationlog.EventUserCreate, Sort: "createdAt", Order: "desc", Page: 1, PageSize: 10,
	})
	if err != nil {
		t.Fatalf("list user operations: %v", err)
	}
	if len(ops) != 1 || ops[0].CorrelationID != "r2-user-001" {
		t.Fatalf("user operation correlation = %+v, want r2-user-001", ops)
	}
	if ops[0].SessionID == "" {
		t.Fatal("user operation missing session_id")
	}
	if ops[0].Detail == nil {
		t.Fatal("user operation missing structured detail")
	}
	detail, err := operationlog.ParseDetail(*ops[0].Detail)
	if err != nil || detail.Action != "create" || detail.After["username"] != "r2-user" {
		t.Fatalf("user operation detail = %+v, err=%v", detail, err)
	}
}

// Anti-resurrection guard (REC-004): after 0006 records_retire the retired
// /api/records product routes are NOT registered on the mux. Any HTTP method
// against /api/records must fail closed with 404 rather than reaching a
// handler or appending operation-log rows.
func TestRetiredRecordsRoutesUnregistered(t *testing.T) {
	env := newAuthTestEnv(t)

	for _, method := range []string{
		http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete,
	} {
		// Unauthenticated attempt.
		req := httptest.NewRequest(method, "/api/records", strings.NewReader(`{}`))
		rr := httptest.NewRecorder()
		env.mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusNotFound {
			t.Fatalf("%s /api/records status = %d, want 404 (route unregistered)", method, rr.Code)
		}

		// Authenticated attempt must still fail closed (no route, not 403).
		req = bearer(t, adminToken(t, env), method, "/api/records", `{}`)
		rr = httptest.NewRecorder()
		env.mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusNotFound {
			t.Fatalf("authed %s /api/records status = %d, want 404 (route unregistered)", method, rr.Code)
		}

		// Detail path too.
		req = httptest.NewRequest(method, "/api/records/rec-1", strings.NewReader(`{}`))
		rr = httptest.NewRecorder()
		env.mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusNotFound {
			t.Fatalf("%s /api/records/{id} status = %d, want 404 (route unregistered)", method, rr.Code)
		}
	}

	// No operation-log rows may be appended by any of the above (no handlers exist).
	ops, err := env.operations.ListOperations(10)
	if err != nil {
		t.Fatalf("ListOperations: %v", err)
	}
	for _, op := range ops {
		if strings.HasPrefix(op.Event, "records.") {
			t.Fatalf("unexpected records.* operation-log row %q after retired-route attempts", op.Event)
		}
	}
}
