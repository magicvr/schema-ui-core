package users

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/handler"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	authsession "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession"
	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
)

func newTestEnv(t *testing.T) (*auth.Authenticator, *store.Store) {
	t.Helper()
	hash, err := auth.HashPassword("test-password", 4)
	if err != nil {
		t.Fatalf("hash seed password: %v", err)
	}
	st, err := testsupport.OpenStore(filepath.Join(t.TempDir(), "test.db"), "admin", hash, true)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	a := auth.New([]byte("test-secret"), 15*time.Minute, 30*24*time.Hour, st, false)
	return a, st
}

func planWithUsers(t *testing.T) kernel.Plan {
	t.Helper()
	registry, err := kernel.NewRegistry(kernel.BuiltinModules())
	if err != nil {
		t.Fatalf("registry: %v", err)
	}
	plan, err := registry.Resolve([]string{
		"core.server-registration", "core.auth-session", "core.manifest-route",
		"core.navigation-capability", "core.schema-render", "core.operationlog",
		"admin.users",
	})
	if err != nil {
		t.Fatalf("resolve plan: %v", err)
	}
	return plan
}

func TestUsersProviderRegistersSurfaces(t *testing.T) {
	a, st := newTestEnv(t)
	provider := New(a, authsession.NewRepository(st), st)
	set, err := kernel.RegisterContributions(context.Background(), planWithUsers(t), []kernel.Provider{provider})
	if err != nil {
		t.Fatalf("RegisterContributions: %v", err)
	}

	wantRoutes := []string{"GET /api/users", "GET /api/users/{id}", "POST /api/users", "PATCH /api/users/{id}", "DELETE /api/users/{id}"}
	if len(set.Routes) != len(wantRoutes) {
		t.Fatalf("routes = %d, want %d", len(set.Routes), len(wantRoutes))
	}
	for _, key := range wantRoutes {
		if !slices.ContainsFunc(set.Routes, func(r kernel.RouteContribution) bool { return r.Key == key }) {
			t.Fatalf("route %q missing from provider set", key)
		}
	}
	if len(set.Pages) != 1 || set.Pages[0].PageID != "users" || set.Pages[0].Owner != ModuleID {
		t.Fatalf("pages = %+v, want single users page owned by %s", set.Pages, ModuleID)
	}
	wantPerms := []string{"users.read", "users.write"}
	if len(set.Permissions) != len(wantPerms) {
		t.Fatalf("permissions = %d, want %d", len(set.Permissions), len(wantPerms))
	}
	for _, perm := range wantPerms {
		if !slices.ContainsFunc(set.Permissions, func(p kernel.PermissionContribution) bool { return p.Permission == perm }) {
			t.Fatalf("permission %q missing", perm)
		}
	}
	if len(set.Navigation) != 1 || set.Navigation[0].NodeID != "menu_users" {
		t.Fatalf("navigation = %+v, want menu_users", set.Navigation)
	}
	if len(set.Fragments) != 1 || set.Fragments[0].FragmentID != "users" {
		t.Fatalf("fragments = %+v, want users fragment", set.Fragments)
	}
}

// TestUsersProviderServesAuthenticatedCRUD validates the C3.3 production path:
// the provider surface (core + provider routes, mirroring composition) serves
// the users resource with the frozen auth/permission behavior (anonymous 401,
// authenticated list 200, unknown detail 404). This is the in-test cutover
// evidence replacing the removed central Register comparison.
func TestUsersProviderServesAuthenticatedCRUD(t *testing.T) {
	a, st := newTestEnv(t)
	plan := planWithUsers(t)
	provider := New(a, authsession.NewRepository(st), st)
	set, err := kernel.RegisterContributions(context.Background(), plan, []kernel.Provider{provider})
	if err != nil {
		t.Fatalf("RegisterContributions: %v", err)
	}
	mux := http.NewServeMux()
	handler.Register(mux, a, st, plan) // core auth/health/schema
	for _, route := range set.Routes {
		mux.Handle(route.Method+" "+route.Pattern, route.Handler)
	}

	// Anonymous access fails closed.
	anon := httptest.NewRecorder()
	mux.ServeHTTP(anon, httptest.NewRequest(http.MethodGet, "/api/users", nil))
	if anon.Code != http.StatusUnauthorized {
		t.Fatalf("anonymous GET /api/users = %d, want 401", anon.Code)
	}

	// Admin login, then authenticated list.
	login := httptest.NewRecorder()
	mux.ServeHTTP(login, httptest.NewRequest(http.MethodPost, "/api/auth/login",
		strings.NewReader(`{"username":"admin","password":"test-password"}`)))
	if login.Code != http.StatusOK {
		t.Fatalf("login = %d, want 200", login.Code)
	}
	var body struct {
		AccessToken string `json:"accessToken"`
	}
	if err := json.NewDecoder(login.Body).Decode(&body); err != nil || body.AccessToken == "" {
		t.Fatalf("login body missing accessToken: %v", err)
	}

	list := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	req.Header.Set("Authorization", "Bearer "+body.AccessToken)
	mux.ServeHTTP(list, req)
	if list.Code != http.StatusOK {
		t.Fatalf("authenticated GET /api/users = %d, want 200", list.Code)
	}

	// Unknown detail → 404.
	detail := httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/users/unknown-id", nil)
	req.Header.Set("Authorization", "Bearer "+body.AccessToken)
	mux.ServeHTTP(detail, req)
	if detail.Code != http.StatusNotFound {
		t.Fatalf("authenticated GET /api/users/unknown-id = %d, want 404", detail.Code)
	}
}

// TestUsersProviderFullCRUD covers the C3.4 behavior matrix on the provider
// finalize path (create → list → detail → patch → delete) with operationlog rows
// appended, proving the production surface preserves the frozen CRUD semantics.
func TestUsersProviderFullCRUD(t *testing.T) {
	a, st := newTestEnv(t)
	plan := planWithUsers(t)
	provider := New(a, authsession.NewRepository(st), st)
	set, err := kernel.RegisterContributions(context.Background(), plan, []kernel.Provider{provider})
	if err != nil {
		t.Fatalf("RegisterContributions: %v", err)
	}
	mux := http.NewServeMux()
	handler.Register(mux, a, st, plan)
	for _, route := range set.Routes {
		mux.Handle(route.Method+" "+route.Pattern, route.Handler)
	}

	login := httptest.NewRecorder()
	mux.ServeHTTP(login, httptest.NewRequest(http.MethodPost, "/api/auth/login",
		strings.NewReader(`{"username":"admin","password":"test-password"}`)))
	if login.Code != http.StatusOK {
		t.Fatalf("login = %d, want 200", login.Code)
	}
	var body struct {
		AccessToken string `json:"accessToken"`
	}
	if err := json.NewDecoder(login.Body).Decode(&body); err != nil || body.AccessToken == "" {
		t.Fatalf("login body missing accessToken: %v", err)
	}
	authReq := func(method, path, payload string) *httptest.ResponseRecorder {
		var req *http.Request
		if payload != "" {
			req = httptest.NewRequest(method, path, strings.NewReader(payload))
		} else {
			req = httptest.NewRequest(method, path, nil)
		}
		req.Header.Set("Authorization", "Bearer "+body.AccessToken)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		return rr
	}

	// create → 201
	created := authReq(http.MethodPost, "/api/users", `{"username":"c3crud","name":"C3 CRUD","password":"passw0rd-ok"}`)
	if created.Code != http.StatusCreated {
		t.Fatalf("create = %d, want 201", created.Code)
	}
	var row struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(created.Body).Decode(&row); err != nil || row.ID == "" {
		t.Fatalf("create body missing id: %v", err)
	}

	// list → 200 + contains the user
	list := authReq(http.MethodGet, "/api/users", "")
	if list.Code != http.StatusOK || !strings.Contains(list.Body.String(), "c3crud") {
		t.Fatalf("list = %d, want 200 containing c3crud", list.Code)
	}

	// detail → 200
	detail := authReq(http.MethodGet, "/api/users/"+row.ID, "")
	if detail.Code != http.StatusOK {
		t.Fatalf("detail = %d, want 200", detail.Code)
	}

	// patch → 200
	patch := authReq(http.MethodPatch, "/api/users/"+row.ID, `{"name":"C3 Renamed"}`)
	if patch.Code != http.StatusOK {
		t.Fatalf("patch = %d, want 200", patch.Code)
	}

	// delete → 204
	del := authReq(http.MethodDelete, "/api/users/"+row.ID, "")
	if del.Code != http.StatusNoContent {
		t.Fatalf("delete = %d, want 204", del.Code)
	}

	// operationlog rows appended for the successful writes.
	ops, err := st.ListOperations(20)
	if err != nil {
		t.Fatalf("ListOperations: %v", err)
	}
	var userOps []string
	for _, op := range ops {
		if strings.HasPrefix(op.Event, "users.") {
			userOps = append(userOps, op.Event)
		}
	}
	for _, want := range []string{"users.create", "users.update", "users.delete"} {
		if !slices.Contains(userOps, want) {
			t.Fatalf("operationlog missing %q (got %v)", want, userOps)
		}
	}
}
