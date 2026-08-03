package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
)

// GOAL-011 S2 · roles list/detail expose the seed system roles with system=true
// and millisecond timestamps.
func TestRolesListAndDetail(t *testing.T) {
	env := newAuthTestEnv(t)
	code, list := getResource(t, env, "/api/roles?pageSize=100")
	if code != http.StatusOK {
		t.Fatalf("list status = %d, want 200: %v", code, list)
	}
	items, ok := list["items"].([]any)
	if !ok || len(items) < 3 {
		t.Fatalf("list items = %v, want >= 3 seed roles", list["items"])
	}
	byID := map[string]map[string]any{}
	for _, item := range items {
		row, ok := item.(map[string]any)
		if !ok {
			continue
		}
		if id, _ := row["id"].(string); id != "" {
			byID[id] = row
		}
	}
	admin, ok := byID["role-admin"]
	if !ok {
		t.Fatalf("role-admin missing in list: %v", byID)
	}
	if admin["system"] != true {
		t.Fatalf("role-admin system = %v, want true", admin["system"])
	}
	if admin["key"] != "admin" || admin["updatedAt"] == nil {
		t.Fatalf("role-admin row = %v", admin)
	}

	code, _ = getResource(t, env, "/api/roles/role-admin")
	if code != http.StatusOK {
		t.Fatalf("detail status = %d, want 200", code)
	}
	code, body := getResource(t, env, "/api/roles/role-nope")
	if code != http.StatusNotFound || body["error"] != "ROLE_NOT_FOUND" {
		t.Fatalf("not found = %d %v, want 404 ROLE_NOT_FOUND", code, body)
	}
}

// GOAL-011 S2 · roles write lifecycle with system/in-use/invalid-key protection.
func TestRolesWriteLifecycleAndProtection(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)
	post := func(body string) (int, map[string]any) {
		req := bearer(t, token, http.MethodPost, "/api/roles", body)
		rr := httptest.NewRecorder()
		env.mux.ServeHTTP(rr, req)
		var out map[string]any
		_ = json.NewDecoder(rr.Body).Decode(&out)
		return rr.Code, out
	}

	// create a user role
	code, created := post(`{"key":"ops","name":"Operator"}`)
	if code != http.StatusCreated {
		t.Fatalf("create status = %d, want 201: %v", code, created)
	}
	if created["id"] != "role-ops" || created["system"] != false {
		t.Fatalf("created = %v, want id role-ops system false", created)
	}

	// duplicate key → 409 ROLE_KEY_TAKEN
	code, out := post(`{"key":"ops","name":"Dup"}`)
	if code != http.StatusConflict || out["error"] != "ROLE_KEY_TAKEN" {
		t.Fatalf("duplicate = %d %v, want 409 ROLE_KEY_TAKEN", code, out)
	}
	// invalid key → 400 INVALID_ROLE_KEY
	code, out = post(`{"key":"OPS!","name":"Bad"}`)
	if code != http.StatusBadRequest || out["error"] != "INVALID_ROLE_KEY" {
		t.Fatalf("bad key = %d %v, want 400 INVALID_ROLE_KEY", code, out)
	}

	// update name
	req := bearer(t, token, http.MethodPatch, "/api/roles/role-ops", `{"name":"Ops"}`)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("update status = %d, want 200", rr.Code)
	}
	// system role update → 409 ROLE_SYSTEM
	req = bearer(t, token, http.MethodPatch, "/api/roles/role-admin", `{"name":"Root"}`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	_ = json.NewDecoder(rr.Body).Decode(&out)
	if rr.Code != http.StatusConflict || out["error"] != "ROLE_SYSTEM" {
		t.Fatalf("system update = %d %v, want 409 ROLE_SYSTEM", rr.Code, out)
	}

	// assign a user to ops then delete → 409 ROLE_IN_USE
	req = bearer(t, token, http.MethodPost, "/api/users",
		`{"username":"opsuser","name":"Ops User","password":"pw12345","roles":["ops"]}`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("create user: %d %s", rr.Code, rr.Body.String())
	}
	req = bearer(t, token, http.MethodDelete, "/api/roles/role-ops", "")
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	_ = json.NewDecoder(rr.Body).Decode(&out)
	if rr.Code != http.StatusConflict || out["error"] != "ROLE_IN_USE" {
		t.Fatalf("in-use delete = %d %v, want 409 ROLE_IN_USE", rr.Code, out)
	}

	// system role delete → 409 ROLE_SYSTEM
	req = bearer(t, token, http.MethodDelete, "/api/roles/role-admin", "")
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	_ = json.NewDecoder(rr.Body).Decode(&out)
	if rr.Code != http.StatusConflict || out["error"] != "ROLE_SYSTEM" {
		t.Fatalf("system delete = %d %v, want 409 ROLE_SYSTEM", rr.Code, out)
	}
}

// GOAL-011 S2 · a free role can be deleted (204) then 404.
func TestRolesDeleteFreeRole(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)
	req := bearer(t, token, http.MethodPost, "/api/roles", `{"key":"temporary","name":"Temp"}`)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("create: %d %s", rr.Code, rr.Body.String())
	}
	req = bearer(t, token, http.MethodDelete, "/api/roles/role-temporary", "")
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d, want 204", rr.Code)
	}
	code, _ := getResource(t, env, "/api/roles/role-temporary")
	if code != http.StatusNotFound {
		t.Fatalf("detail after delete = %d, want 404", code)
	}
}

// GOAL-011 S2 · anonymous 401; viewer reads roles but cannot write (403).
func TestRolesAuthGates(t *testing.T) {
	env := newAuthTestEnv(t)

	req := httptest.NewRequest(http.MethodGet, "/api/roles", nil)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("anonymous status = %d, want 401", rr.Code)
	}

	env.addUser(t, "viewer3", "pw", []string{"viewer"})
	vToken := env.login(t, "viewer3", "pw")
	req = bearer(t, vToken, http.MethodGet, "/api/roles", "")
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("viewer list status = %d, want 200", rr.Code)
	}
	req = bearer(t, vToken, http.MethodPost, "/api/roles", `{"key":"x","name":"X"}`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("viewer write status = %d, want 403", rr.Code)
	}
}

// GOAL-011 S2 · roles write endpoints append roles.* operation-log events.
func TestRolesOperationLogEvents(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)

	req := bearer(t, token, http.MethodPost, "/api/roles", `{"key":"oprole","name":"Op Role"}`)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("create: %d %s", rr.Code, rr.Body.String())
	}
	req = bearer(t, token, http.MethodPatch, "/api/roles/role-oprole", `{"name":"Op Role 2"}`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("patch: %d %s", rr.Code, rr.Body.String())
	}
	req = bearer(t, token, http.MethodDelete, "/api/roles/role-oprole", "")
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("delete: %d %s", rr.Code, rr.Body.String())
	}

	ops, err := env.st.ListOperations(10)
	if err != nil {
		t.Fatalf("ListOperations: %v", err)
	}
	var roleOps []store.Operation
	for _, op := range ops {
		if strings.HasPrefix(op.Event, "roles.") {
			roleOps = append(roleOps, op)
		}
	}
	want := []string{store.EventRoleDelete, store.EventRoleUpdate, store.EventRoleCreate}
	if len(roleOps) != len(want) {
		t.Fatalf("role ops = %d, want %d", len(roleOps), len(want))
	}
	for i, ev := range want {
		if roleOps[i].Event != ev {
			t.Fatalf("roleOps[%d].event = %q, want %q", i, roleOps[i].Event, ev)
		}
	}
}
