package users

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
	provider := New(a, st)
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

// TestUsersProviderCompatWithCentral verifies the provider-generated surface
// routes requests identically to the current central registration (freeze §7
// step 2 compat comparison, in-test, no production dual registration).
func TestUsersProviderCompatWithCentral(t *testing.T) {
	a, st := newTestEnv(t)
	plan := planWithUsers(t)

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
		{http.MethodGet, "/api/users"},
		{http.MethodGet, "/api/users/unknown-id"},
		{http.MethodPost, "/api/users"},
		{http.MethodPatch, "/api/users/unknown-id"},
		{http.MethodDelete, "/api/users/unknown-id"},
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
