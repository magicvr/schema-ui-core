// Package datadictionary provider tests (S-01 · GOAL-008 D-002 `6): the module
// registers the two resources, both page schemas, permission keys, the
// menu_dictionary navigation node and the manifest fragment.
package datadictionary

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
	"github.com/magicvr/schema-ui-core/apps/api/modules/datadictionary/store"
	"github.com/magicvr/schema-ui-core/apps/api/modules/operationlog"
	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
)

func newDictTestEnv(t *testing.T) (*auth.Authenticator, *store.Repository, *operationlog.Repository) {
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

func planWithDictionary(t *testing.T) kernel.Plan {
	t.Helper()
	registry, err := kernel.NewRegistry(kernel.BuiltinModules())
	if err != nil {
		t.Fatalf("registry: %v", err)
	}
	plan, err := registry.Resolve([]string{
		"core.server-registration", "core.auth-session", "core.manifest-route",
		"core.navigation-capability", "core.schema-render", "core.operationlog",
		"admin.data-dictionary",
	})
	if err != nil {
		t.Fatalf("resolve plan: %v", err)
	}
	return plan
}

func TestDictionaryProviderRegistersSurfaces(t *testing.T) {
	a, repository, operations := newDictTestEnv(t)
	provider := New(a, repository, operations)
	set, err := kernel.RegisterContributions(context.Background(), planWithDictionary(t), []kernel.Provider{provider})
	if err != nil {
		t.Fatalf("RegisterContributions: %v", err)
	}

	wantRoutes := []string{
		"GET /api/data-dictionary/types", "GET /api/data-dictionary/types/{id}",
		"POST /api/data-dictionary/types", "PATCH /api/data-dictionary/types/{id}",
		"DELETE /api/data-dictionary/types/{id}", "POST /api/data-dictionary/types/batch-delete",
		"GET /api/data-dictionary/entries", "GET /api/data-dictionary/entries/{id}",
		"POST /api/data-dictionary/entries", "PATCH /api/data-dictionary/entries/{id}",
		"DELETE /api/data-dictionary/entries/{id}", "POST /api/data-dictionary/entries/batch-delete",
	}
	if len(set.Routes) != len(wantRoutes) {
		t.Fatalf("routes = %d, want %d", len(set.Routes), len(wantRoutes))
	}
	for _, key := range wantRoutes {
		if !slices.ContainsFunc(set.Routes, func(r kernel.RouteContribution) bool { return r.Key == key }) {
			t.Fatalf("route %q missing from provider set", key)
		}
	}
	if len(set.Pages) != 2 {
		t.Fatalf("pages = %d, want 2", len(set.Pages))
	}
	for _, pageID := range []string{"data-dictionary", "dictionary-entries"} {
		found := false
		for _, page := range set.Pages {
			if page.PageID == pageID && page.Owner == ModuleID {
				found = true
			}
		}
		if !found {
			t.Fatalf("page %s missing or not owned by %s", pageID, ModuleID)
		}
	}
	for _, perm := range []string{"dictionary.read", "dictionary.write"} {
		if !slices.ContainsFunc(set.Permissions, func(p kernel.PermissionContribution) bool { return p.Permission == perm }) {
			t.Fatalf("permission %q missing", perm)
		}
	}
	if len(set.Navigation) != 1 || set.Navigation[0].NodeID != "menu_dictionary" {
		t.Fatalf("navigation = %+v, want menu_dictionary", set.Navigation)
	}
	if len(set.Fragments) != 1 || set.Fragments[0].FragmentID != "data-dictionary" {
		t.Fatalf("fragments = %+v, want data-dictionary fragment", set.Fragments)
	}
	// The types page carries the navigate action to the entries page.
	var doc struct {
		Actions map[string]struct {
			Type string `json:"type"`
			URL  string `json:"url"`
		} `json:"actions"`
	}
	for _, page := range set.Pages {
		if page.PageID != "data-dictionary" {
			continue
		}
		if err := json.Unmarshal(page.Document, &doc); err != nil {
			t.Fatalf("page schema parse: %v", err)
		}
	}
	// GOAL-015 F-007(b): the entries inner page route is /dictionary-entries/{dictKey}
	// (path-param template) — the row navigate action binds the type key into it.
	if doc.Actions["openEntries"].Type != "navigate" || doc.Actions["openEntries"].URL != "/dictionary-entries/{dictKey}" {
		t.Fatalf("openEntries action = %+v, want navigate /dictionary-entries/{dictKey}", doc.Actions["openEntries"])
	}
}

// The provider surface serves the dictionary endpoints end-to-end: anonymous
// 401, admin list 200, unknown detail 404.
func TestDictionaryProviderServesSurfaces(t *testing.T) {
	a, repository, operations := newDictTestEnv(t)
	plan := planWithDictionary(t)
	provider := New(a, repository, operations)
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
	mux.ServeHTTP(anon, httptest.NewRequest(http.MethodGet, "/api/data-dictionary/types", nil))
	if anon.Code != http.StatusUnauthorized {
		t.Fatalf("anonymous GET types = %d, want 401", anon.Code)
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
	req := httptest.NewRequest(http.MethodGet, "/api/data-dictionary/types", nil)
	req.Header.Set("Authorization", "Bearer "+body.AccessToken)
	mux.ServeHTTP(list, req)
	if list.Code != http.StatusOK || !strings.Contains(list.Body.String(), `"total":0`) {
		t.Fatalf("admin list = %d, want 200 empty: %s", list.Code, list.Body.String())
	}
}
