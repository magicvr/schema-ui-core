package roles

import (
	"context"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/handler"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
)

func newTestEnv(t *testing.T) (*auth.Authenticator, *store.Store) {
	t.Helper()
	hash, err := auth.HashPassword("test-password", 4)
	if err != nil {
		t.Fatalf("hash seed password: %v", err)
	}
	st, err := store.Open(filepath.Join(t.TempDir(), "test.db"), "admin", hash, true)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	a := auth.New([]byte("test-secret"), 15*time.Minute, 30*24*time.Hour, st, false)
	return a, st
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
	a, st := newTestEnv(t)
	provider := New(a, st)
	set, err := kernel.RegisterContributions(context.Background(), planWithRoles(t), []kernel.Provider{provider})
	if err != nil {
		t.Fatalf("RegisterContributions: %v", err)
	}

	wantRoutes := []string{"GET /api/roles", "GET /api/roles/{id}", "POST /api/roles", "PATCH /api/roles/{id}", "DELETE /api/roles/{id}"}
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

// TestRolesProviderCompatWithCentral mirrors the users compat check.
func TestRolesProviderCompatWithCentral(t *testing.T) {
	a, st := newTestEnv(t)
	plan := planWithRoles(t)

	central := http.NewServeMux()
	handler.Register(central, a, st, plan)

	provider := New(a, st)
	set, err := kernel.RegisterContributions(context.Background(), plan, []kernel.Provider{provider})
	if err != nil {
		t.Fatalf("RegisterContributions: %v", err)
	}
	providerMux := http.NewServeMux()
	for _, route := range set.Routes {
		providerMux.Handle(route.Method+" "+route.Pattern, route.Handler)
	}

	for _, tc := range []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/roles"},
		{http.MethodGet, "/api/roles/unknown-id"},
		{http.MethodPost, "/api/roles"},
		{http.MethodPatch, "/api/roles/unknown-id"},
		{http.MethodDelete, "/api/roles/unknown-id"},
	} {
		rrCentral := httptest.NewRecorder()
		central.ServeHTTP(rrCentral, httptest.NewRequest(tc.method, tc.path, nil))
		rrProvider := httptest.NewRecorder()
		providerMux.ServeHTTP(rrProvider, httptest.NewRequest(tc.method, tc.path, nil))
		if rrCentral.Code != rrProvider.Code {
			t.Fatalf("%s %s: central=%d provider=%d, want identical", tc.method, tc.path, rrCentral.Code, rrProvider.Code)
		}
	}
}
