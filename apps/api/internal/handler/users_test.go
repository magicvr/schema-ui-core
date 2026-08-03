package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
)

// GOAL-011 S2 · users list/detail expose the seeded admin without password_hash.
func TestUsersListAndDetail(t *testing.T) {
	env := newAuthTestEnv(t)
	code, list := getResource(t, env, "/api/users?pageSize=100")
	if code != http.StatusOK {
		t.Fatalf("list status = %d, want 200: %v", code, list)
	}
	items, ok := list["items"].([]any)
	if !ok || len(items) == 0 {
		t.Fatalf("list items = %v, want at least the seed admin", list["items"])
	}
	for _, item := range items {
		row, ok := item.(map[string]any)
		if !ok {
			continue
		}
		if _, has := row["password_hash"]; has {
			t.Fatalf("password_hash leaked in list row: %v", row)
		}
		if _, has := row["password"]; has {
			t.Fatalf("password leaked in list row: %v", row)
		}
	}
	// detail of the seed admin
	code, detail := getResource(t, env, "/api/users/user-admin")
	if code != http.StatusOK {
		t.Fatalf("detail status = %d, want 200: %v", code, detail)
	}
	if detail["username"] != "admin" {
		t.Fatalf("username = %v, want admin", detail["username"])
	}
	if _, has := detail["password_hash"]; has {
		t.Fatalf("password_hash leaked in detail: %v", detail)
	}
	if detail["updatedAt"] == nil || detail["createdAt"] == nil {
		t.Fatalf("detail missing timestamps: %v", detail)
	}
}

// GOAL-011 S2 · users write lifecycle: create → detail → patch → delete, with
// roles assignment and operation-log events.
func TestUsersCreateUpdateDeleteLifecycle(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)

	// create (editor role exists in seed; viewer too)
	req := bearer(t, token, http.MethodPost, "/api/users",
		`{"username":"alice","name":"Alice","password":"secret123","roles":["editor","viewer"]}`)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want 201: %s", rr.Code, rr.Body.String())
	}
	var created map[string]any
	_ = json.NewDecoder(rr.Body).Decode(&created)
	id, _ := created["id"].(string)
	if !strings.HasPrefix(id, "usr-") {
		t.Fatalf("created id = %q, want usr- prefix", id)
	}
	if created["username"] != "alice" || created["name"] != "Alice" {
		t.Fatalf("created = %v, want alice/Alice", created)
	}
	if _, has := created["password_hash"]; has {
		t.Fatalf("password_hash leaked on create: %v", created)
	}

	// detail
	code, detail := getResource(t, env, "/api/users/"+id)
	if code != http.StatusOK || detail["roles"] == nil {
		t.Fatalf("detail status = %d roles=%v, want 200 with roles", code, detail["roles"])
	}

	// patch name + roles (promote to admin — there are now two admins)
	req = bearer(t, token, http.MethodPatch, "/api/users/"+id,
		`{"name":"Alice A.","roles":["admin","editor"]}`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("patch status = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	// delete → 204
	req = bearer(t, token, http.MethodDelete, "/api/users/"+id, "")
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d, want 204: %s", rr.Code, rr.Body.String())
	}
	code, _ = getResource(t, env, "/api/users/"+id)
	if code != http.StatusNotFound {
		t.Fatalf("detail after delete status = %d, want 404", code)
	}
	// not-found code is USER_NOT_FOUND
	code, body := getResource(t, env, "/api/users/usr-nonexistent")
	if code != http.StatusNotFound || body["error"] != "USER_NOT_FOUND" {
		t.Fatalf("not found = %d %v, want 404 USER_NOT_FOUND", code, body)
	}
}

// GOAL-011 S2 · create field validation and role-ref validation.
func TestUsersCreateValidation(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)
	send := func(body string) (int, map[string]any) {
		req := bearer(t, token, http.MethodPost, "/api/users", body)
		rr := httptest.NewRecorder()
		env.mux.ServeHTTP(rr, req)
		var out map[string]any
		_ = json.NewDecoder(rr.Body).Decode(&out)
		return rr.Code, out
	}

	if code, out := send(`{"name":"NoPass","roles":[]}`); code != http.StatusBadRequest || out["error"] != "INVALID_CREATE_FIELD" {
		t.Fatalf("missing password = %d %v, want 400 INVALID_CREATE_FIELD", code, out)
	}
	if code, out := send(`{"username":"u","name":"U","password":"x","roles":"admin"}`); code != http.StatusBadRequest || out["error"] != "INVALID_ROLE_REF" {
		t.Fatalf("roles not array = %d %v, want 400 INVALID_ROLE_REF", code, out)
	}
	if code, out := send(`{"username":"u","name":"U","password":"x","roles":["ghost"]}`); code != http.StatusBadRequest || out["error"] != "INVALID_ROLE_REF" {
		t.Fatalf("unknown role = %d %v, want 400 INVALID_ROLE_REF", code, out)
	}
}

// GOAL-011 S2 · duplicate username → 409 USERNAME_TAKEN.
func TestUsersCreateDuplicateUsername(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)
	req := bearer(t, token, http.MethodPost, "/api/users",
		`{"username":"admin","name":"Dup","password":"x12345"}`)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusConflict {
		t.Fatalf("duplicate status = %d, want 409: %s", rr.Code, rr.Body.String())
	}
	var out map[string]any
	_ = json.NewDecoder(rr.Body).Decode(&out)
	if out["error"] != "USERNAME_TAKEN" {
		t.Fatalf("duplicate error = %v, want USERNAME_TAKEN", out["error"])
	}
}

// GOAL-011 S2 · self-delete and self-demote are rejected (SELF_OPERATION).
func TestUsersSelfProtection(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)

	// cannot delete yourself
	req := bearer(t, token, http.MethodDelete, "/api/users/user-admin", "")
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusConflict {
		t.Fatalf("self-delete status = %d, want 409: %s", rr.Code, rr.Body.String())
	}
	var out map[string]any
	_ = json.NewDecoder(rr.Body).Decode(&out)
	if out["error"] != "SELF_OPERATION" {
		t.Fatalf("self-delete error = %v, want SELF_OPERATION", out["error"])
	}

	// cannot remove your own admin role
	req = bearer(t, token, http.MethodPatch, "/api/users/user-admin", `{"roles":["editor"]}`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	_ = json.NewDecoder(rr.Body).Decode(&out)
	if rr.Code != http.StatusConflict || out["error"] != "SELF_OPERATION" {
		t.Fatalf("self-demote = %d %v, want 409 SELF_OPERATION", rr.Code, out)
	}
}

// GOAL-011 S2 · a non-admin actor (dev session, users.write but not an admin
// user) deleting the only admin → 409 LAST_ADMIN at the HTTP layer
// (A-004 F-002).
func TestUsersLastAdminHTTP(t *testing.T) {
	env := newDevSessionTestEnv(t)
	// The dev-session identity (dev-001) holds users.write but is NOT a real
	// admin user; deleting user-admin (the only admin) leaves zero admins.
	req := httptest.NewRequest(http.MethodDelete, "/api/users/user-admin", nil)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusConflict {
		t.Fatalf("delete only admin status = %d, want 409: %s", rr.Code, rr.Body.String())
	}
	var out map[string]any
	_ = json.NewDecoder(rr.Body).Decode(&out)
	if out["error"] != "LAST_ADMIN" {
		t.Fatalf("error = %v, want LAST_ADMIN", out["error"])
	}
}

// GOAL-011 S4 · every users route is authenticated; viewer can read list/detail
// but every write route is forbidden.
func TestUsersAuthGates(t *testing.T) {
	env := newAuthTestEnv(t)
	routes := []struct {
		method string
		path   string
		body   string
		read   bool
	}{
		{http.MethodGet, "/api/users", "", true},
		{http.MethodGet, "/api/users/user-admin", "", true},
		{http.MethodPost, "/api/users", `{"username":"x","name":"X","password":"y12345"}`, false},
		{http.MethodPatch, "/api/users/user-admin", `{"name":"Root"}`, false},
		{http.MethodDelete, "/api/users/user-admin", "", false},
	}
	for _, tc := range routes {
		code, out := sendJSON(t, env.mux, tc.method, tc.path, tc.body)
		if code != http.StatusUnauthorized || out["error"] != "UNAUTHENTICATED" {
			t.Fatalf("anonymous %s %s = %d %v, want 401 UNAUTHENTICATED", tc.method, tc.path, code, out)
		}
	}

	env.addUser(t, "viewer2", "pw", []string{"viewer"})
	vToken := env.login(t, "viewer2", "pw")
	for _, tc := range routes {
		req := bearer(t, vToken, tc.method, tc.path, tc.body)
		rr := httptest.NewRecorder()
		env.mux.ServeHTTP(rr, req)
		want := http.StatusForbidden
		if tc.read {
			want = http.StatusOK
		}
		if rr.Code != want {
			t.Fatalf("viewer %s %s = %d, want %d: %s", tc.method, tc.path, rr.Code, want, rr.Body.String())
		}
		if !tc.read {
			var out map[string]any
			_ = json.NewDecoder(rr.Body).Decode(&out)
			if out["error"] != "FORBIDDEN" {
				t.Fatalf("viewer %s %s error = %v, want FORBIDDEN", tc.method, tc.path, out["error"])
			}
		}
	}
}

// GOAL-011 S2 · a created user can log in with the password provided at create
// (password is hashed, never stored/echoed).
func TestUsersPasswordWriteOnly(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)
	req := bearer(t, token, http.MethodPost, "/api/users",
		`{"username":"bob","name":"Bob","password":"hunter2x","roles":[]}`)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want 201: %s", rr.Code, rr.Body.String())
	}
	// login as bob works → password was bcrypt-hashed
	if _, _, _, err := env.a.Login("bob", "hunter2x", time.Now().UTC()); err != nil {
		t.Fatalf("bob login: %v", err)
	}
}

// GOAL-011 S2 · users write endpoints append users.* operation-log events with
// the actor and a username detail (never a secret).
func TestUsersOperationLogEvents(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)

	req := bearer(t, token, http.MethodPost, "/api/users",
		`{"username":"opuser","name":"Op User","password":"pw12345","roles":[]}`)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("create: %d %s", rr.Code, rr.Body.String())
	}
	var created map[string]any
	_ = json.NewDecoder(rr.Body).Decode(&created)
	id, _ := created["id"].(string)

	req = bearer(t, token, http.MethodPatch, "/api/users/"+id, `{"name":"Op User 2"}`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("patch: %d %s", rr.Code, rr.Body.String())
	}

	req = bearer(t, token, http.MethodDelete, "/api/users/"+id, "")
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("delete: %d %s", rr.Code, rr.Body.String())
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
	want := []string{store.EventUserDelete, store.EventUserUpdate, store.EventUserCreate}
	if len(userOps) != len(want) {
		t.Fatalf("user ops = %d, want %d", len(userOps), len(want))
	}
	for i, ev := range want {
		if userOps[i].Event != ev {
			t.Fatalf("userOps[%d].event = %q, want %q", i, userOps[i].Event, ev)
		}
	}
	for i, op := range userOps {
		if op.ActorID != "user-admin" {
			t.Fatalf("actor = %q, want user-admin", op.ActorID)
		}
		if op.ActorName != "Admin" {
			t.Fatalf("actor_name = %q, want Admin", op.ActorName)
		}
		if op.RecordID == nil || *op.RecordID != id {
			t.Fatalf("userOps[%d].record_id = %v, want %s", i, op.RecordID, id)
		}
		if op.Detail != nil && strings.Contains(*op.Detail, "password") {
			t.Fatalf("detail leaked a secret: %v", *op.Detail)
		}
		if op.Event == store.EventUserDelete {
			if op.Detail != nil {
				t.Fatalf("delete detail = %q, want nil", *op.Detail)
			}
		} else if op.Detail == nil || *op.Detail != `{"username":"opuser"}` {
			t.Fatalf("%s detail = %v, want username-only JSON", op.Event, op.Detail)
		}
	}
}
