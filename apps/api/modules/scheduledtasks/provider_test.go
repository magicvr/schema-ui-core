// Package scheduledtasks provider tests (S-04 · GOAL-010 D-002 `6): the module
// registers the task resources, both page schemas, permission keys,
// menu_scheduled_tasks navigation and the fragment.
package scheduledtasks

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
	"github.com/magicvr/schema-ui-core/apps/api/internal/ratelimit"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	authsession "github.com/magicvr/schema-ui-core/apps/api/modules/authsession"
	"github.com/magicvr/schema-ui-core/apps/api/modules/operationlog"
	"github.com/magicvr/schema-ui-core/apps/api/modules/scheduledtasks/store"
	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
)

func newTasksTestEnv(t *testing.T) (*auth.Authenticator, *store.Repository, *operationlog.Repository) {
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
	return a, store.NewRepository(st), operationlog.NewRepository(st)
}

func planWithTasks(t *testing.T) kernel.Plan {
	t.Helper()
	registry, err := kernel.NewRegistry(kernel.BuiltinModules())
	if err != nil {
		t.Fatalf("registry: %v", err)
	}
	plan, err := registry.Resolve([]string{
		"core.server-registration", "core.auth-session", "core.manifest-route",
		"core.navigation-capability", "core.schema-render", "core.operationlog",
		"admin.scheduled-tasks",
	})
	if err != nil {
		t.Fatalf("resolve plan: %v", err)
	}
	return plan
}

func TestTasksProviderRegistersSurfaces(t *testing.T) {
	a, repository, operations := newTasksTestEnv(t)
	provider := New(a, repository, operations)
	set, err := kernel.RegisterContributions(context.Background(), planWithTasks(t), []kernel.Provider{provider})
	if err != nil {
		t.Fatalf("RegisterContributions: %v", err)
	}

	wantRoutes := []string{
		"GET /api/scheduled-tasks", "GET /api/scheduled-tasks/{id}",
		"POST /api/scheduled-tasks", "PATCH /api/scheduled-tasks/{id}",
		"DELETE /api/scheduled-tasks/{id}", "POST /api/scheduled-tasks/batch-delete",
		"POST /api/scheduled-tasks/{id}/run", "GET /api/scheduled-tasks/{id}/runs",
		"GET /api/task-runs", "GET /api/task-runs/{id}",
	}
	for _, key := range wantRoutes {
		if !slices.ContainsFunc(set.Routes, func(r kernel.RouteContribution) bool { return r.Key == key }) {
			t.Fatalf("route %q missing from provider set", key)
		}
	}
	if len(set.Pages) != 2 {
		t.Fatalf("pages = %d, want 2", len(set.Pages))
	}
	for _, perm := range []string{"tasks.read", "tasks.write"} {
		if !slices.ContainsFunc(set.Permissions, func(p kernel.PermissionContribution) bool { return p.Permission == perm }) {
			t.Fatalf("permission %q missing", perm)
		}
	}
	if len(set.Navigation) != 1 || set.Navigation[0].NodeID != "menu_scheduled_tasks" {
		t.Fatalf("navigation = %+v, want menu_scheduled_tasks", set.Navigation)
	}
	if len(set.Fragments) != 1 || set.Fragments[0].FragmentID != "scheduled-tasks" {
		t.Fatalf("fragments = %+v, want scheduled-tasks fragment", set.Fragments)
	}
}

// The provider surface serves the tasks endpoints end-to-end (anonymous 401,
// admin list 200, cron-valid create 201, invalid cron 400).
func TestTasksProviderServesSurfaces(t *testing.T) {
	a, repository, operations := newTasksTestEnv(t)
	plan := planWithTasks(t)
	provider := New(a, repository, operations)
	set, err := kernel.RegisterContributions(context.Background(), plan, []kernel.Provider{provider})
	if err != nil {
		t.Fatalf("RegisterContributions: %v", err)
	}
	mux := http.NewServeMux()
	handler.Register(mux, a, nil, operations, plan, ratelimit.NewProvider())
	for _, route := range set.Routes {
		mux.Handle(route.Method+" "+route.Pattern, route.Handler)
	}

	anon := httptest.NewRecorder()
	mux.ServeHTTP(anon, httptest.NewRequest(http.MethodGet, "/api/scheduled-tasks", nil))
	if anon.Code != http.StatusUnauthorized {
		t.Fatalf("anonymous list = %d, want 401", anon.Code)
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
	created := authReq(http.MethodPost, "/api/scheduled-tasks",
		`{"key":"hourly","cron":"0 * * * *","name":"Hourly"}`)
	if created.Code != http.StatusCreated {
		t.Fatalf("create = %d: %s", created.Code, created.Body.String())
	}
	bad := authReq(http.MethodPost, "/api/scheduled-tasks",
		`{"key":"bad","cron":"not a cron","name":"Bad"}`)
	if bad.Code != http.StatusBadRequest || !strings.Contains(bad.Body.String(), "INVALID_CRON") {
		t.Fatalf("invalid cron = %d: %s", bad.Code, bad.Body.String())
	}
}
