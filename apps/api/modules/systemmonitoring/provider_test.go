// Package systemmonitoring provider tests (S-03 · GOAL-009 D-002 `5): the
// module registers the status + errors routes, the page schema, the
// monitoring.read permission, menu_monitoring navigation and the fragment.
package systemmonitoring

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
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	authsession "github.com/magicvr/schema-ui-core/apps/api/modules/authsession"
	"github.com/magicvr/schema-ui-core/apps/api/modules/operationlog"
	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
)

func newMonitoringTestEnv(t *testing.T) (*auth.Authenticator, *store.Store, kernel.Plan, *operationlog.Repository, string) {
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
	registry, err := kernel.NewRegistry(kernel.BuiltinModules())
	if err != nil {
		t.Fatalf("registry: %v", err)
	}
	plan, err := registry.Resolve([]string{
		"core.server-registration", "core.auth-session", "core.manifest-route",
		"core.navigation-capability", "core.schema-render", "core.operationlog",
		"admin.system-monitoring",
	})
	if err != nil {
		t.Fatalf("resolve plan: %v", err)
	}
	return a, st, plan, operationlog.NewRepository(st), filepath.Join(t.TempDir(), "monitor.db")
}

func TestMonitoringProviderRegistersSurfaces(t *testing.T) {
	a, st, plan, operations, dbPath := newMonitoringTestEnv(t)
	provider := New(a, st, plan, func() bool { return true }, dbPath, time.Now(), operations)
	set, err := kernel.RegisterContributions(context.Background(), plan, []kernel.Provider{provider})
	if err != nil {
		t.Fatalf("RegisterContributions: %v", err)
	}

	wantRoutes := []string{
		"GET /api/system-monitoring/status",
		"GET /api/system-monitoring/errors", "GET /api/system-monitoring/errors/{id}",
	}
	for _, key := range wantRoutes {
		if !slices.ContainsFunc(set.Routes, func(r kernel.RouteContribution) bool { return r.Key == key }) {
			t.Fatalf("route %q missing from provider set", key)
		}
	}
	if len(set.Pages) != 1 || set.Pages[0].PageID != "system-monitoring" || set.Pages[0].Owner != ModuleID {
		t.Fatalf("pages = %+v, want single system-monitoring page owned by %s", set.Pages, ModuleID)
	}
	if len(set.Permissions) != 1 || set.Permissions[0].Permission != "monitoring.read" {
		t.Fatalf("permissions = %+v, want monitoring.read", set.Permissions)
	}
	if len(set.Navigation) != 1 || set.Navigation[0].NodeID != "menu_monitoring" {
		t.Fatalf("navigation = %+v, want menu_monitoring", set.Navigation)
	}
	if len(set.Fragments) != 1 || set.Fragments[0].FragmentID != "system-monitoring" {
		t.Fatalf("fragments = %+v, want system-monitoring fragment", set.Fragments)
	}
}

// The provider surface serves the status endpoint end-to-end (anonymous 401,
// admin 200 with the status summary).
func TestMonitoringProviderServesStatus(t *testing.T) {
	a, st, plan, operations, dbPath := newMonitoringTestEnv(t)
	provider := New(a, st, plan, func() bool { return true }, dbPath, time.Now(), operations)
	set, err := kernel.RegisterContributions(context.Background(), plan, []kernel.Provider{provider})
	if err != nil {
		t.Fatalf("RegisterContributions: %v", err)
	}
	mux := http.NewServeMux()
	handler.Register(mux, a, nil, operations, plan)
	for _, route := range set.Routes {
		mux.Handle(route.Method+" "+route.Pattern, route.Handler)
	}

	anon := httptest.NewRecorder()
	mux.ServeHTTP(anon, httptest.NewRequest(http.MethodGet, "/api/system-monitoring/status", nil))
	if anon.Code != http.StatusUnauthorized {
		t.Fatalf("anonymous status = %d, want 401", anon.Code)
	}

	login := httptest.NewRecorder()
	mux.ServeHTTP(login, httptest.NewRequest(http.MethodPost, "/api/auth/login",
		strings.NewReader(`{"username":"admin","password":"test-password"}`)))
	if login.Code != http.StatusOK {
		t.Fatalf("login = %d", login.Code)
	}
	var body struct {
		AccessToken string `json:"accessToken"`
	}
	if err := json.NewDecoder(login.Body).Decode(&body); err != nil || body.AccessToken == "" {
		t.Fatalf("login body missing accessToken")
	}
	req := httptest.NewRequest(http.MethodGet, "/api/system-monitoring/status", nil)
	req.Header.Set("Authorization", "Bearer "+body.AccessToken)
	status := httptest.NewRecorder()
	mux.ServeHTTP(status, req)
	if status.Code != http.StatusOK || !strings.Contains(status.Body.String(), `"status":"ok"`) {
		t.Fatalf("admin status = %d: %s", status.Code, status.Body.String())
	}
}
