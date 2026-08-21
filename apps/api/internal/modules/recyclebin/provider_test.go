// Package recyclebin provider tests (S-12 · GOAL-012 D-002 §3): the module
// registers routes/permissions/nav/fragment, and the real store-backed service
// serves the conflict path end-to-end over HTTP (grok A-003 F-002).
package recyclebin

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

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/handler"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	authsession "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession"
	datadictionarystore "github.com/magicvr/schema-ui-core/apps/api/internal/modules/datadictionary/store"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
	recyclestore "github.com/magicvr/schema-ui-core/apps/api/internal/modules/recyclebin/store"
	tasksstore "github.com/magicvr/schema-ui-core/apps/api/internal/modules/scheduledtasks/store"
	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
)

func newRecycleTestEnv(t *testing.T) (*auth.Authenticator, *Service, *operationlog.Repository, *datadictionarystore.Repository) {
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
	authRepository := authsession.NewRepository(st)
	a := auth.NewWithRepository([]byte("test-secret"), 15*time.Minute, 30*24*time.Hour, authRepository, false)
	dictionary := datadictionarystore.NewRepository(st)
	return a, NewService(recyclestore.NewRepository(st), dictionary, tasksstore.NewRepository(st)), operationlog.NewRepository(st), dictionary
}

func planWithRecycle(t *testing.T) kernel.Plan {
	t.Helper()
	registry, err := kernel.NewRegistry(kernel.BuiltinModules())
	if err != nil {
		t.Fatalf("registry: %v", err)
	}
	plan, err := registry.Resolve([]string{
		"core.server-registration", "core.auth-session", "core.manifest-route",
		"core.navigation-capability", "core.schema-render", "core.operationlog",
		"admin.recycle-bin",
	})
	if err != nil {
		t.Fatalf("resolve plan: %v", err)
	}
	return plan
}

func TestRecycleProviderRegistersSurfaces(t *testing.T) {
	a, service, operations, _ := newRecycleTestEnv(t)
	provider := New(a, service, operations)
	set, err := kernel.RegisterContributions(context.Background(), planWithRecycle(t), []kernel.Provider{provider})
	if err != nil {
		t.Fatalf("RegisterContributions: %v", err)
	}
	for _, key := range []string{
		"GET /api/recycle-bin", "GET /api/recycle-bin/{id}",
		"POST /api/recycle-bin/{id}/restore", "DELETE /api/recycle-bin/{id}",
		"POST /api/recycle-bin/purge-all",
	} {
		if !slices.ContainsFunc(set.Routes, func(r kernel.RouteContribution) bool { return r.Key == key }) {
			t.Fatalf("route %q missing", key)
		}
	}
	if len(set.Pages) != 1 || set.Pages[0].PageID != "recycle-bin" {
		t.Fatalf("pages = %+v", set.Pages)
	}
	for _, perm := range []string{"recycle.read", "recycle.write"} {
		if !slices.ContainsFunc(set.Permissions, func(p kernel.PermissionContribution) bool { return p.Permission == perm }) {
			t.Fatalf("permission %q missing", perm)
		}
	}
	if len(set.Navigation) != 1 || set.Navigation[0].NodeID != "menu_recycle_bin" {
		t.Fatalf("navigation = %+v", set.Navigation)
	}
	if len(set.Fragments) != 1 || set.Fragments[0].FragmentID != "recycle-bin" {
		t.Fatalf("fragments = %+v", set.Fragments)
	}
}

// F-002 (grok A-003): the real store-backed service must serve the restore
// conflict over HTTP — a snapshot whose key is taken again returns 409
// RECYCLE_RESTORE_CONFLICT and stays unrestored for a later retry.
func TestRecycleRealServiceRestoreConflictHTTP(t *testing.T) {
	a, service, operations, dictionary := newRecycleTestEnv(t)
	plan := planWithRecycle(t)
	provider := New(a, service, operations)
	set, err := kernel.RegisterContributions(context.Background(), plan, []kernel.Provider{provider})
	if err != nil {
		t.Fatalf("RegisterContributions: %v", err)
	}
	mux := http.NewServeMux()
	handler.Register(mux, a, nil, operations, plan)
	for _, route := range set.Routes {
		mux.Handle(route.Method+" "+route.Pattern, route.Handler)
	}
	// Login.
	login := httptest.NewRecorder()
	mux.ServeHTTP(login, httptest.NewRequest(http.MethodPost, "/api/auth/login",
		strings.NewReader(`{"username":"admin","password":"test-password"}`)))
	if login.Code != http.StatusOK {
		t.Fatalf("login = %d: %s", login.Code, login.Body.String())
	}
	var loginBody struct {
		AccessToken string `json:"accessToken"`
	}
	if err := json.NewDecoder(login.Body).Decode(&loginBody); err != nil || loginBody.AccessToken == "" {
		t.Fatalf("login body missing accessToken")
	}
	authReq := func(method, path, payload string) *httptest.ResponseRecorder {
		var req *http.Request
		if payload != "" {
			req = httptest.NewRequest(method, path, strings.NewReader(payload))
		} else {
			req = httptest.NewRequest(method, path, nil)
		}
		req.Header.Set("Authorization", "Bearer "+loginBody.AccessToken)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		return rr
	}

	now := time.Now().UTC()
	// Create a dict type, record a snapshot as the factory would (post-delete).
	if err := dictionary.CreateType(datadictionarystore.DictType{ID: "t1", Key: "status", Name: "Status", Enabled: true, CreatedAt: now, UpdatedAt: now}); err != nil {
		t.Fatalf("create type: %v", err)
	}
	if err := service.Record(t.Context(), "dict-types", "t1", map[string]any{
		"id": "t1", "key": "status", "name": "Status", "enabled": true, "description": "", "sort": 0,
		"createdAt": float64(now.Unix()), "updatedAt": float64(now.Unix()),
	}, account.User{ID: "user-admin", Name: "Admin"}, now); err != nil {
		t.Fatalf("record: %v", err)
	}
	if _, err := dictionary.DeleteType("t1"); err != nil {
		t.Fatalf("delete type: %v", err)
	}
	// Re-create the same key with a NEW id: the snapshot's key is taken again.
	if err := dictionary.CreateType(datadictionarystore.DictType{ID: "t2", Key: "status", Name: "Status v2", Enabled: true, CreatedAt: now, UpdatedAt: now}); err != nil {
		t.Fatalf("recreate type: %v", err)
	}
	items, _, err := service.List(recyclestore.ListFilter{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("items = %d, want 1", len(items))
	}
	// Restore → 409 RECYCLE_RESTORE_CONFLICT (localized envelope).
	conflict := authReq(http.MethodPost, "/api/recycle-bin/"+items[0].ID+"/restore", "")
	if conflict.Code != http.StatusConflict {
		t.Fatalf("conflict restore = %d: %s", conflict.Code, conflict.Body.String())
	}
	if !strings.Contains(conflict.Body.String(), "RECYCLE_RESTORE_CONFLICT") {
		t.Fatalf("conflict body = %s", conflict.Body.String())
	}
	if !strings.Contains(conflict.Body.String(), "messageKey") {
		t.Fatalf("conflict body must carry the catalog messageKey (F-003): %s", conflict.Body.String())
	}
	// The snapshot is retained and unrestored (retryable).
	detail := authReq(http.MethodGet, "/api/recycle-bin/"+items[0].ID, "")
	var detailBody map[string]any
	if err := json.Unmarshal(detail.Body.Bytes(), &detailBody); err != nil {
		t.Fatalf("decode detail: %v", err)
	}
	if detailBody["restored"] != false {
		t.Fatalf("snapshot must stay unrestored after conflict: %v", detailBody)
	}
	// Resolve the conflict (delete the occupying row) → restore succeeds.
	if _, err := dictionary.DeleteType("t2"); err != nil {
		t.Fatalf("delete occupying: %v", err)
	}
	ok := authReq(http.MethodPost, "/api/recycle-bin/"+items[0].ID+"/restore", "")
	if ok.Code != http.StatusOK {
		t.Fatalf("restore after conflict resolution = %d: %s", ok.Code, ok.Body.String())
	}
}
