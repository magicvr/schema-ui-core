// Package filelibrary provider tests (S-02 · GOAL-007 D-002 `7): the module
// registers the full HTTP surface, page schema, permission keys, navigation
// node and manifest fragment; the registered surface serves the library
// endpoints with the frozen auth behavior.
package filelibrary

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

func newFileLibraryTestEnv(t *testing.T) (*auth.Authenticator, *store.Store, *authsession.Repository, *operationlog.Repository, string) {
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
	return a, st, repository, operationlog.NewRepository(st), filepath.Join(t.TempDir(), "uploads")
}

func planWithFileLibrary(t *testing.T) kernel.Plan {
	t.Helper()
	registry, err := kernel.NewRegistry(kernel.BuiltinModules())
	if err != nil {
		t.Fatalf("registry: %v", err)
	}
	plan, err := registry.Resolve([]string{
		"core.server-registration", "core.auth-session", "core.manifest-route",
		"core.navigation-capability", "core.schema-render", "core.operationlog",
		"admin.file-library",
	})
	if err != nil {
		t.Fatalf("resolve plan: %v", err)
	}
	return plan
}

func TestFileLibraryProviderRegistersSurfaces(t *testing.T) {
	a, _, _, operations, uploadDir := newFileLibraryTestEnv(t)
	provider := New(a, operations, uploadDir)
	set, err := kernel.RegisterContributions(context.Background(), planWithFileLibrary(t), []kernel.Provider{provider})
	if err != nil {
		t.Fatalf("RegisterContributions: %v", err)
	}

	wantRoutes := []string{
		"GET /api/library/files", "GET /api/library/files/{id}",
		"GET /api/library/files/{id}/download",
		"DELETE /api/library/files/{id}",
		"POST /api/library/files/upload",
	}
	if len(set.Routes) != len(wantRoutes) {
		t.Fatalf("routes = %d, want %d", len(set.Routes), len(wantRoutes))
	}
	for _, key := range wantRoutes {
		if !slices.ContainsFunc(set.Routes, func(r kernel.RouteContribution) bool { return r.Key == key }) {
			t.Fatalf("route %q missing from provider set", key)
		}
	}
	if len(set.Pages) != 1 || set.Pages[0].PageID != "file-library" || set.Pages[0].Owner != ModuleID {
		t.Fatalf("pages = %+v, want single file-library page owned by %s", set.Pages, ModuleID)
	}
	wantPerms := []string{"files.read", "files.delete"}
	if len(set.Permissions) != len(wantPerms) {
		t.Fatalf("permissions = %d, want %d", len(set.Permissions), len(wantPerms))
	}
	for _, perm := range wantPerms {
		if !slices.ContainsFunc(set.Permissions, func(p kernel.PermissionContribution) bool { return p.Permission == perm }) {
			t.Fatalf("permission %q missing", perm)
		}
	}
	if len(set.Navigation) != 1 || set.Navigation[0].NodeID != "menu_files" {
		t.Fatalf("navigation = %+v, want menu_files", set.Navigation)
	}
	if len(set.Fragments) != 1 || set.Fragments[0].FragmentID != "file-library" {
		t.Fatalf("fragments = %+v, want file-library fragment", set.Fragments)
	}
	// The page schema document parses and carries the custom download handler.
	var doc struct {
		Actions map[string]struct {
			Type    string `json:"type"`
			Handler string `json:"handler"`
		} `json:"actions"`
	}
	if err := json.Unmarshal(set.Pages[0].Document, &doc); err != nil {
		t.Fatalf("page schema parse: %v", err)
	}
	if doc.Actions["downloadFile"].Type != "custom" || doc.Actions["downloadFile"].Handler != "library.download" {
		t.Fatalf("downloadFile action = %+v, want custom library.download", doc.Actions["downloadFile"])
	}
}

// The provider surface serves the library endpoints end-to-end (anonymous 401,
// admin list 200, unknown detail 404) through the same finalize path the
// composition root uses.
func TestFileLibraryProviderServesSurfaces(t *testing.T) {
	a, st, _, operations, uploadDir := newFileLibraryTestEnv(t)
	plan := planWithFileLibrary(t)
	provider := New(a, operations, uploadDir)
	set, err := kernel.RegisterContributions(context.Background(), plan, []kernel.Provider{provider})
	if err != nil {
		t.Fatalf("RegisterContributions: %v", err)
	}
	mux := http.NewServeMux()
	handler.Register(mux, a, st, operations, plan)
	handler.RegisterUpload(mux, a, uploadDir)
	for _, route := range set.Routes {
		mux.Handle(route.Method+" "+route.Pattern, route.Handler)
	}

	// Anonymous list fails closed.
	anon := httptest.NewRecorder()
	mux.ServeHTTP(anon, httptest.NewRequest(http.MethodGet, "/api/library/files", nil))
	if anon.Code != http.StatusUnauthorized {
		t.Fatalf("anonymous GET /api/library/files = %d, want 401", anon.Code)
	}

	// Admin login → empty library 200.
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
	req := httptest.NewRequest(http.MethodGet, "/api/library/files", nil)
	req.Header.Set("Authorization", "Bearer "+body.AccessToken)
	mux.ServeHTTP(list, req)
	if list.Code != http.StatusOK || !strings.Contains(list.Body.String(), `"total":0`) {
		t.Fatalf("admin list = %d, want 200 empty: %s", list.Code, list.Body.String())
	}

	// Unknown detail → 404 FILE_NOT_FOUND.
	detail := httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/library/files/00000000000000000000000000000000", nil)
	req.Header.Set("Authorization", "Bearer "+body.AccessToken)
	mux.ServeHTTP(detail, req)
	if detail.Code != http.StatusNotFound {
		t.Fatalf("unknown detail = %d, want 404", detail.Code)
	}
}
