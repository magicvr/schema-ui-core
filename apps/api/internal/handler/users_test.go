package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
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

// D1 hardening: PATCH with "roles": null means "no role change" — it must not
// clear the user's roles (null ≠ explicit empty array).
func TestUsersPatchRolesNullKeepsRoles(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)
	req := bearer(t, token, http.MethodPost, "/api/users",
		`{"username":"alice","name":"Alice","password":"secret123","roles":["editor"]}`)
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

	// Explicit empty array clears roles (existing contract).
	req = bearer(t, token, http.MethodPatch, "/api/users/"+id, `{"roles":[]}`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("patch empty roles status = %d", rr.Code)
	}
	_, detail := getResource(t, env, "/api/users/"+id)
	if got := detail["roles"]; got != nil {
		if arr, ok := got.([]any); !ok || len(arr) != 0 {
			t.Fatalf("roles after explicit [] = %v, want empty", got)
		}
	}

	// Re-assign a role, then patch with roles: null → roles must be unchanged.
	req = bearer(t, token, http.MethodPatch, "/api/users/"+id, `{"roles":["viewer"]}`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("reassign status = %d", rr.Code)
	}
	req = bearer(t, token, http.MethodPatch, "/api/users/"+id, `{"roles":null,"name":"Alice N."}`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("patch null roles status = %d: %s", rr.Code, rr.Body.String())
	}
	_, detail = getResource(t, env, "/api/users/"+id)
	if detail["name"] != "Alice N." {
		t.Fatalf("name = %v, want updated Alice N.", detail["name"])
	}
	roles, _ := detail["roles"].([]any)
	if len(roles) != 1 || roles[0] != "viewer" {
		t.Fatalf("roles after null patch = %v, want [viewer] (unchanged)", detail["roles"])
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

	if code, out := send(`{"username":"nopass","name":"NoPass","roles":[]}`); code != http.StatusBadRequest || out["error"] != "INVALID_PASSWORD" {
		t.Fatalf("missing password = %d %v, want 400 INVALID_PASSWORD", code, out)
	}
	if code, out := send(`{"username":"u","name":"U","password":"password","roles":42}`); code != http.StatusBadRequest || out["error"] != "INVALID_ROLE_REF" {
		t.Fatalf("roles not array = %d %v, want 400 INVALID_ROLE_REF", code, out)
	}
	if code, out := send(`{"username":"u","name":"U","password":"password","roles":["ghost"]}`); code != http.StatusBadRequest || out["error"] != "INVALID_ROLE_REF" {
		t.Fatalf("unknown role = %d %v, want 400 INVALID_ROLE_REF", code, out)
	}
}

// GOAL-011 S2 · duplicate username → 409 USERNAME_TAKEN.
func TestUsersCreateDuplicateUsername(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)
	req := bearer(t, token, http.MethodPost, "/api/users",
		`{"username":"admin","name":"Dup","password":"x1234567"}`)
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
		{http.MethodPost, "/api/users", `{"username":"x","name":"X","password":"y123456"}`, false},
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

func TestUsersRoleAssignmentAuthorization(t *testing.T) {
	env := newAuthTestEnv(t)
	now := time.Now().UTC()
	if _, err := env.authRepository.CreateRoleWithGrants(
		"manager", "User manager",
		[]string{"roles.assign", "roles.read", "users.read", "users.write"}, nil, now,
	); err != nil {
		t.Fatalf("create manager role: %v", err)
	}
	if _, err := env.authRepository.CreateRoleWithGrants(
		"writer", "Role writer", []string{"roles.write", "users.read"}, nil, now,
	); err != nil {
		t.Fatalf("create writer role: %v", err)
	}
	if _, err := env.authRepository.CreateRoleWithGrants(
		"users-writer", "Users writer", []string{"users.read", "users.write"}, nil, now,
	); err != nil {
		t.Fatalf("create users-writer role: %v", err)
	}
	env.addUser(t, "manager", "manager-password", []string{"manager"})
	env.addUser(t, "writer", "writer-password", []string{"users-writer"})

	managerToken := env.login(t, "manager", "manager-password")
	req := bearer(t, managerToken, http.MethodPost, "/api/users",
		`{"username":"delegate","name":"Delegate","password":"delegate-password","roles":["viewer"]}`)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("delegated subset create = %d, want 201: %s", rr.Code, rr.Body.String())
	}
	var created map[string]any
	_ = json.NewDecoder(rr.Body).Decode(&created)
	id, _ := created["id"].(string)

	req = bearer(t, managerToken, http.MethodPatch, "/api/users/"+id, `{"roles":["writer"]}`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	var out map[string]any
	_ = json.NewDecoder(rr.Body).Decode(&out)
	if rr.Code != http.StatusForbidden || out["error"] != "ROLE_ASSIGNMENT_FORBIDDEN" {
		t.Fatalf("permission superset assignment = %d %v, want 403 ROLE_ASSIGNMENT_FORBIDDEN", rr.Code, out)
	}

	req = bearer(t, managerToken, http.MethodPatch, "/api/users/"+id, `{"roles":["admin"]}`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	_ = json.NewDecoder(rr.Body).Decode(&out)
	if rr.Code != http.StatusForbidden || out["error"] != "ROLE_ASSIGNMENT_FORBIDDEN" {
		t.Fatalf("non-admin admin assignment = %d %v, want 403 ROLE_ASSIGNMENT_FORBIDDEN", rr.Code, out)
	}

	writerToken := env.login(t, "writer", "writer-password")
	req = bearer(t, writerToken, http.MethodPatch, "/api/users/"+id, `{"roles":["viewer"]}`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	_ = json.NewDecoder(rr.Body).Decode(&out)
	if rr.Code != http.StatusForbidden || out["error"] != "ROLE_ASSIGNMENT_FORBIDDEN" {
		t.Fatalf("users.write-only assignment = %d %v, want 403 ROLE_ASSIGNMENT_FORBIDDEN", rr.Code, out)
	}
}

func TestUsersPasswordPolicyPreservesBytesAndRevokesRefresh(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)
	postBody := func(body string) (int, map[string]any) {
		req := bearer(t, token, http.MethodPost, "/api/users", body)
		rr := httptest.NewRecorder()
		env.mux.ServeHTTP(rr, req)
		var out map[string]any
		_ = json.NewDecoder(rr.Body).Decode(&out)
		return rr.Code, out
	}
	post := func(password string) (int, map[string]any) {
		return postBody(`{"username":"password-user","name":"Password User","password":` + quote(password) + `}`)
	}
	for _, invalid := range []string{"short7", "        ", strings.Repeat("x", 73)} {
		if code, out := post(invalid); code != http.StatusBadRequest || out["error"] != "INVALID_PASSWORD" {
			t.Fatalf("invalid password len=%d = %d %v, want 400 INVALID_PASSWORD", len([]byte(invalid)), code, out)
		}
	}
	for _, body := range []string{
		`{"username":"password-user","name":"Password User"}`,
		`{"username":"password-user","name":"Password User","password":null}`,
		`{"username":"password-user","name":"Password User","password":12345678}`,
	} {
		if code, out := postBody(body); code != http.StatusBadRequest || out["error"] != "INVALID_PASSWORD" {
			t.Fatalf("non-string or missing password = %d %v, want 400 INVALID_PASSWORD", code, out)
		}
	}

	password := "  exact-password  "
	code, created := post(password)
	if code != http.StatusCreated {
		t.Fatalf("create exact password = %d %v, want 201", code, created)
	}
	if _, _, _, err := env.a.Login("password-user", password, time.Now().UTC()); err != nil {
		t.Fatalf("login with exact password: %v", err)
	}
	if _, _, _, err := env.a.Login("password-user", strings.TrimSpace(password), time.Now().UTC()); !errors.Is(err, auth.ErrInvalidCredentials) {
		t.Fatalf("trimmed password login err = %v, want ErrInvalidCredentials", err)
	}

	now := time.Now().UTC()
	access, refresh, _, err := env.a.Login("password-user", password, now)
	if err != nil {
		t.Fatalf("login before rotation: %v", err)
	}
	id, _ := created["id"].(string)
	req := bearer(t, token, http.MethodPatch, "/api/users/"+id, `{"password":12345678}`)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	var invalidPatch map[string]any
	_ = json.NewDecoder(rr.Body).Decode(&invalidPatch)
	if rr.Code != http.StatusBadRequest || invalidPatch["error"] != "INVALID_PASSWORD" {
		t.Fatalf("non-string password patch = %d %v, want 400 INVALID_PASSWORD", rr.Code, invalidPatch)
	}
	newPassword := "  replacement-password  "
	req = bearer(t, token, http.MethodPatch, "/api/users/"+id, `{"password":`+quote(newPassword)+`}`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("password patch = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	if _, _, _, err := env.a.Refresh(refresh, now.Add(time.Minute)); !errors.Is(err, auth.ErrTokenRevoked) {
		t.Fatalf("old refresh after password change = %v, want ErrTokenRevoked", err)
	}
	req = bearer(t, access, http.MethodGet, "/api/accounts/me", "")
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("existing access token before TTL expiry = %d, want 200", rr.Code)
	}
	if _, _, _, err := env.a.Login("password-user", password, now.Add(time.Minute)); !errors.Is(err, auth.ErrInvalidCredentials) {
		t.Fatalf("old password login err = %v, want ErrInvalidCredentials", err)
	}
	if _, _, _, err := env.a.Login("password-user", newPassword, now.Add(time.Minute)); err != nil {
		t.Fatalf("new password login: %v", err)
	}
}

// GOAL-011 S2 · users write endpoints append users.* operation-log events with
// the actor and a username detail (never a secret).
func TestUsersOperationLogEvents(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)

	req := bearer(t, token, http.MethodPost, "/api/users",
		`{"username":"opuser","name":"Op User","password":"pw123456","roles":[]}`)
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
	want := []string{operationlog.EventUserDelete, operationlog.EventUserUpdate, operationlog.EventUserCreate}
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
		if op.Event == operationlog.EventUserDelete {
			if op.Detail != nil {
				t.Fatalf("delete detail = %q, want nil", *op.Detail)
			}
		} else if op.Detail == nil || *op.Detail != `{"username":"opuser"}` {
			t.Fatalf("%s detail = %v, want username-only JSON", op.Event, op.Detail)
		}
	}
}
