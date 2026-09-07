package activity

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/handler"
	"github.com/magicvr/schema-ui-core/apps/api/internal/ratelimit"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/modules/operationlog"
	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
)

func newTestEnv(t *testing.T) (*auth.Authenticator, *store.Store, *operationlog.Repository) {
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
	return a, st, operationlog.NewRepository(st)
}

func planWithActivity(t *testing.T) kernel.Plan {
	t.Helper()
	registry, err := kernel.NewRegistry(kernel.BuiltinModules())
	if err != nil {
		t.Fatalf("registry: %v", err)
	}
	plan, err := registry.Resolve([]string{
		"core.server-registration", "core.auth-session", "core.manifest-route",
		"core.navigation-capability", "core.schema-render", "core.operationlog",
		"admin.activity",
	})
	if err != nil {
		t.Fatalf("resolve plan: %v", err)
	}
	return plan
}

func TestActivityProviderRegistersSurfaces(t *testing.T) {
	a, _, operations := newTestEnv(t)
	set, err := kernel.RegisterContributions(context.Background(), planWithActivity(t), []kernel.Provider{New(a, operations)})
	if err != nil {
		t.Fatalf("RegisterContributions: %v", err)
	}
	// Read-only: list/detail + F-03 export route.
	if len(set.Routes) != 3 {
		t.Fatalf("routes = %d, want 3 (read-only list/detail/export)", len(set.Routes))
	}
	if len(set.Pages) != 1 || set.Pages[0].PageID != "activity" {
		t.Fatalf("pages = %+v", set.Pages)
	}
	if len(set.Permissions) != 1 || set.Permissions[0].Permission != "operations.read" {
		t.Fatalf("permissions = %+v", set.Permissions)
	}
	if len(set.Fragments) != 1 || set.Fragments[0].FragmentID != "activity" {
		t.Fatalf("fragments = %+v", set.Fragments)
	}
}

// TestActivityProviderServesReadOnly verifies the operations list works for an
// admin and the activity surface is read-only (no write routes) on the provider
// finalize path (C4.2).
func TestActivityProviderServesReadOnly(t *testing.T) {
	a, st, operations := newTestEnv(t)
	plan := planWithActivity(t)
	set, err := kernel.RegisterContributions(context.Background(), plan, []kernel.Provider{New(a, operations)})
	if err != nil {
		t.Fatalf("RegisterContributions: %v", err)
	}
	mux := http.NewServeMux()
	handler.Register(mux, a, st, operations, plan, ratelimit.NewProvider())
	for _, route := range set.Routes {
		mux.Handle(route.Method+" "+route.Pattern, route.Handler)
	}

	anon := httptest.NewRecorder()
	mux.ServeHTTP(anon, httptest.NewRequest(http.MethodGet, "/api/operations", nil))
	if anon.Code != http.StatusUnauthorized {
		t.Fatalf("anonymous operations = %d, want 401", anon.Code)
	}

	login := httptest.NewRecorder()
	mux.ServeHTTP(login, httptest.NewRequest(http.MethodPost, "/api/auth/login",
		strings.NewReader(`{"username":"admin","password":"test-password"}`)))
	var body struct {
		AccessToken string `json:"accessToken"`
	}
	if err := json.NewDecoder(login.Body).Decode(&body); err != nil || body.AccessToken == "" {
		t.Fatalf("login body missing accessToken: %v", err)
	}

	list := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/operations", nil)
	req.Header.Set("Authorization", "Bearer "+body.AccessToken)
	mux.ServeHTTP(list, req)
	if list.Code != http.StatusOK {
		t.Fatalf("authenticated operations = %d, want 200", list.Code)
	}

	// No write routes registered: POST on a GET-only path → 405 (method not
	// allowed), proving read-only.
	post := httptest.NewRecorder()
	mux.ServeHTTP(post, httptest.NewRequest(http.MethodPost, "/api/operations", nil))
	if post.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST /api/operations = %d, want 405 (read-only)", post.Code)
	}
}
