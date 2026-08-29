package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/modules/operationlog"
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
		`{"username":"opsuser","name":"Ops User","password":"pw123456","roles":["ops"]}`)
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

func TestRolesGrantLifecycleAndEffectiveProjection(t *testing.T) {
	env := newAuthTestEnv(t)
	token := adminToken(t, env)
	req := bearer(t, token, http.MethodPost, "/api/roles",
		`{"key":"support","name":"Support","permissions":["users.read"],"menuItems":["menu-users"]}`)
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("create role with grants = %d, want 201: %s", rr.Code, rr.Body.String())
	}
	var created map[string]any
	_ = json.NewDecoder(rr.Body).Decode(&created)
	if created["editable"] != true || created["deletable"] != true || created["assignedUsers"] != float64(0) {
		t.Fatalf("created management flags = %v", created)
	}
	if _, err := env.authRepository.CreateRoleWithGrants(
		"role-manager", "Role Manager", []string{"roles.read", "roles.write"}, nil, time.Now().UTC(),
	); err != nil {
		t.Fatalf("create non-admin role manager: %v", err)
	}
	env.addUser(t, "role-manager", "role-manager-password", []string{"role-manager"})
	managerToken := env.login(t, "role-manager", "role-manager-password")
	req = bearer(t, managerToken, http.MethodPatch, "/api/roles/role-support", `{"permissions":["roles.read"]}`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	var forbidden map[string]any
	_ = json.NewDecoder(rr.Body).Decode(&forbidden)
	if rr.Code != http.StatusForbidden || forbidden["error"] != "ROLE_GRANT_FORBIDDEN" {
		t.Fatalf("non-admin grant update = %d %v, want 403 ROLE_GRANT_FORBIDDEN", rr.Code, forbidden)
	}

	req = bearer(t, token, http.MethodPost, "/api/users",
		`{"username":"support-user","name":"Support User","password":"support-password","roles":["support"]}`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("create support user = %d: %s", rr.Code, rr.Body.String())
	}
	_, _, accountUser, err := env.a.Login("support-user", "support-password", time.Now().UTC())
	if err != nil {
		t.Fatalf("support login: %v", err)
	}
	if !containsRoles(accountUser.Permissions, "users.read") {
		t.Fatalf("support permissions = %v, want users.read", accountUser.Permissions)
	}
	features, err := env.a.Features(accountUser)
	if err != nil || !features["menu_users"] {
		t.Fatalf("support features = %v err=%v, want menu_users", features, err)
	}

	code, detail := getResource(t, env, "/api/roles/role-support")
	if code != http.StatusOK || detail["assignedUsers"] != float64(1) || detail["deletable"] != false {
		t.Fatalf("assigned role detail = %d %v", code, detail)
	}

	req = bearer(t, token, http.MethodPatch, "/api/roles/role-support",
		`{"permissions":["roles.read"],"menuItems":[]}`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("replace grants = %d, want 200: %s", rr.Code, rr.Body.String())
	}
	_, _, accountUser, err = env.a.Login("support-user", "support-password", time.Now().UTC())
	if err != nil {
		t.Fatalf("support relogin: %v", err)
	}
	if !containsRoles(accountUser.Permissions, "roles.read") || containsRoles(accountUser.Permissions, "users.read") {
		t.Fatalf("replaced permissions = %v, want roles.read only", accountUser.Permissions)
	}
	features, err = env.a.Features(accountUser)
	if err != nil || features["menu_users"] {
		t.Fatalf("features after menu grant removal = %v err=%v", features, err)
	}

	req = bearer(t, token, http.MethodPatch, "/api/roles/role-support", `{"permissions":["ghost.permission"]}`)
	rr = httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	var out map[string]any
	_ = json.NewDecoder(rr.Body).Decode(&out)
	if rr.Code != http.StatusBadRequest || out["error"] != "INVALID_PERMISSION_REF" {
		t.Fatalf("invalid permission = %d %v, want 400 INVALID_PERMISSION_REF", rr.Code, out)
	}
}

// GOAL-011 S4 · every roles route is authenticated; viewer can read list/detail
// but every write route is forbidden.
func TestRolesAuthGates(t *testing.T) {
	env := newAuthTestEnv(t)
	routes := []struct {
		method string
		path   string
		body   string
		read   bool
	}{
		{http.MethodGet, "/api/roles", "", true},
		{http.MethodGet, "/api/roles/role-admin", "", true},
		{http.MethodPost, "/api/roles", `{"key":"x","name":"X"}`, false},
		{http.MethodPatch, "/api/roles/role-admin", `{"name":"Root"}`, false},
		{http.MethodDelete, "/api/roles/role-admin", "", false},
	}
	for _, tc := range routes {
		code, out := sendJSON(t, env.mux, tc.method, tc.path, tc.body)
		if code != http.StatusUnauthorized || out["error"] != "UNAUTHENTICATED" {
			t.Fatalf("anonymous %s %s = %d %v, want 401 UNAUTHENTICATED", tc.method, tc.path, code, out)
		}
	}

	env.addUser(t, "viewer3", "pw", []string{"viewer"})
	vToken := env.login(t, "viewer3", "pw")
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

	ops, err := env.operations.ListOperations(10)
	if err != nil {
		t.Fatalf("ListOperations: %v", err)
	}
	var roleOps []operationlog.Operation
	for _, op := range ops {
		if strings.HasPrefix(op.Event, "roles.") {
			roleOps = append(roleOps, op)
		}
	}
	want := []string{operationlog.EventRoleDelete, operationlog.EventRoleUpdate, operationlog.EventRoleCreate}
	if len(roleOps) != len(want) {
		t.Fatalf("role ops = %d, want %d", len(roleOps), len(want))
	}
	for i, ev := range want {
		if roleOps[i].Event != ev {
			t.Fatalf("roleOps[%d].event = %q, want %q", i, roleOps[i].Event, ev)
		}
	}
	for i, op := range roleOps {
		if op.ActorID != "user-admin" || op.ActorName != "Admin" {
			t.Fatalf("roleOps[%d].actor = %q/%q, want user-admin/Admin", i, op.ActorID, op.ActorName)
		}
		if op.RecordID == nil || *op.RecordID != "role-oprole" {
			t.Fatalf("roleOps[%d].record_id = %v, want role-oprole", i, op.RecordID)
		}
		if op.Detail != nil && strings.Contains(*op.Detail, "password") {
			t.Fatalf("detail leaked a secret: %v", *op.Detail)
		}
		if op.Event == operationlog.EventRoleDelete {
			if op.Detail != nil {
				t.Fatalf("delete detail = %q, want nil", *op.Detail)
			}
		} else {
			envelope, err := operationlog.ParseDetail(*op.Detail)
			if err != nil || envelope.After["key"] != "oprole" {
				t.Fatalf("%s detail = %v, want R2 envelope key=oprole: %v", op.Event, op.Detail, err)
			}
		}
	}
}
