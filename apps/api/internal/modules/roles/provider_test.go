package roles

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
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
)

func newTestEnv(t *testing.T) (*auth.Authenticator, *store.Store, *authsession.Repository, *operationlog.Repository) {
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
	repository := authsession.NewRepository(st)
	a := auth.NewWithRepository([]byte("test-secret"), 15*time.Minute, 30*24*time.Hour, repository, false)
	return a, st, repository, operationlog.NewRepository(st)
}

func planWithRoles(t *testing.T) kernel.Plan {
	t.Helper()
	registry, err := kernel.NewRegistry(kernel.BuiltinModules())
	if err != nil {
		t.Fatalf("registry: %v", err)
	}
	plan, err := registry.Resolve([]string{
		"core.server-registration", "core.auth-session", "core.manifest-route",
		"core.navigation-capability", "core.schema-render", "core.operationlog",
		"admin.roles",
	})
	if err != nil {
		t.Fatalf("resolve plan: %v", err)
	}
	return plan
}

func TestRolesProviderRegistersSurfaces(t *testing.T) {
	a, _, repository, operations := newTestEnv(t)
	provider := New(a, repository, operations)
	set, err := kernel.RegisterContributions(context.Background(), planWithRoles(t), []kernel.Provider{provider})
	if err != nil {
		t.Fatalf("RegisterContributions: %v", err)
	}

	wantRoutes := []string{"GET /api/roles", "GET /api/roles/{id}", "POST /api/roles", "PATCH /api/roles/{id}", "DELETE /api/roles/{id}", "POST /api/roles/batch-delete"}
	if len(set.Routes) != len(wantRoutes) {
		t.Fatalf("routes = %d, want %d", len(set.Routes), len(wantRoutes))
	}
	for _, key := range wantRoutes {
		if !slices.ContainsFunc(set.Routes, func(r kernel.RouteContribution) bool { return r.Key == key }) {
			t.Fatalf("route %q missing from provider set", key)
		}
	}
	if len(set.Pages) != 1 || set.Pages[0].PageID != "roles" || set.Pages[0].Owner != ModuleID {
		t.Fatalf("pages = %+v, want single roles page owned by %s", set.Pages, ModuleID)
	}
	wantPerms := []string{"roles.read", "roles.write", "roles.assign"}
	if len(set.Permissions) != len(wantPerms) {
		t.Fatalf("permissions = %d, want %d", len(set.Permissions), len(wantPerms))
	}
	for _, perm := range wantPerms {
		if !slices.ContainsFunc(set.Permissions, func(p kernel.PermissionContribution) bool { return p.Permission == perm }) {
			t.Fatalf("permission %q missing", perm)
		}
	}
	if len(set.Navigation) != 1 || set.Navigation[0].NodeID != "menu_roles" {
		t.Fatalf("navigation = %+v, want menu_roles", set.Navigation)
	}
	if len(set.Fragments) != 1 || set.Fragments[0].FragmentID != "roles" {
		t.Fatalf("fragments = %+v, want roles fragment", set.Fragments)
	}
}

// TestRolesProviderServesAuthenticatedCRUD mirrors the users cutover check: the
// provider surface serves roles with frozen auth/permission behavior.
func TestRolesProviderServesAuthenticatedCRUD(t *testing.T) {
	a, st, repository, operations := newTestEnv(t)
	plan := planWithRoles(t)
	provider := New(a, repository, operations)
	set, err := kernel.RegisterContributions(context.Background(), plan, []kernel.Provider{provider})
	if err != nil {
		t.Fatalf("RegisterContributions: %v", err)
	}
	mux := http.NewServeMux()
	handler.Register(mux, a, st, operations, plan)
	for _, route := range set.Routes {
		mux.Handle(route.Method+" "+route.Pattern, route.Handler)
	}

	anon := httptest.NewRecorder()
	mux.ServeHTTP(anon, httptest.NewRequest(http.MethodGet, "/api/roles", nil))
	if anon.Code != http.StatusUnauthorized {
		t.Fatalf("anonymous GET /api/roles = %d, want 401", anon.Code)
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

	list := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/roles", nil)
	req.Header.Set("Authorization", "Bearer "+body.AccessToken)
	mux.ServeHTTP(list, req)
	if list.Code != http.StatusOK {
		t.Fatalf("authenticated GET /api/roles = %d, want 200", list.Code)
	}

	detail := httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/roles/unknown-id", nil)
	req.Header.Set("Authorization", "Bearer "+body.AccessToken)
	mux.ServeHTTP(detail, req)
	if detail.Code != http.StatusNotFound {
		t.Fatalf("authenticated GET /api/roles/unknown-id = %d, want 404", detail.Code)
	}
}
