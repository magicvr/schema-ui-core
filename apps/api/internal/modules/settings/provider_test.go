package settings

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
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
	settingsrepository "github.com/magicvr/schema-ui-core/apps/api/internal/modules/settings/repository"
	authsession "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession"
	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
)

func newTestEnv(t *testing.T) (*auth.Authenticator, *store.Store, *settingsrepository.Repository, *operationlog.Repository, *authsession.Repository) {
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
	return a, st, settingsrepository.New(st), operationlog.NewRepository(st), authsession.NewRepository(st)
}

func planWithSettings(t *testing.T) kernel.Plan {
	t.Helper()
	registry, err := kernel.NewRegistry(kernel.BuiltinModules())
	if err != nil {
		t.Fatalf("registry: %v", err)
	}
	plan, err := registry.Resolve([]string{
		"core.server-registration", "core.auth-session", "core.manifest-route",
		"core.navigation-capability", "core.schema-render", "core.operationlog",
		"admin.settings",
	})
	if err != nil {
		t.Fatalf("resolve plan: %v", err)
	}
	return plan
}

func TestSettingsProviderRegistersSurfaces(t *testing.T) {
	a, _, settings, operations, authRepo := newTestEnv(t)
	set, err := kernel.RegisterContributions(context.Background(), planWithSettings(t), []kernel.Provider{New(a, settings, operations, nil, authRepo)})
	if err != nil {
		t.Fatalf("RegisterContributions: %v", err)
	}
	wantRoutes := []string{"GET /api/branding", "GET /api/settings", "GET /api/settings/{id}", "PATCH /api/settings/{id}", "POST /api/settings/{id}/reset", "POST /api/branding/assets", "GET /api/branding/assets/{id}", "GET /api/settings/password-policy", "PATCH /api/settings/password-policy"}
	if len(set.Routes) != len(wantRoutes) {
		t.Fatalf("routes = %d, want %d", len(set.Routes), len(wantRoutes))
	}
	// W26 (GOAL-038 D-001 §2.2): the settings module now owns three pages —
	// settings + the standalone mail console / outbound log (sorted by PageID).
	if len(set.Pages) != 3 {
		t.Fatalf("pages = %d, want 3", len(set.Pages))
	}
	for i, wantID := range []string{"mail", "mail-outbox", "settings"} {
		if set.Pages[i].PageID != wantID || set.Pages[i].Owner != ModuleID {
			t.Fatalf("pages[%d] = %+v, want PageID %q owned by %s", i, set.Pages[i], wantID, ModuleID)
		}
	}
	if len(set.Navigation) != 3 {
		t.Fatalf("navigation = %d, want 3 (menu_settings + menu_mail + menu_mail_outbox)", len(set.Navigation))
	}
	if len(set.Fragments) != 1 || set.Fragments[0].FragmentID != "settings" {
		t.Fatalf("fragments = %+v", set.Fragments)
	}
	if len(set.Configurations) != 1 || set.Configurations[0].Namespace != "settings.branding" {
		t.Fatalf("configurations = %+v", set.Configurations)
	}
	if err := set.Configurations[0].Validate(json.RawMessage(set.Configurations[0].Defaults)); err != nil {
		t.Fatalf("configuration defaults: %v", err)
	}
}

// TestSettingsProviderServesBrandingAndAuth verifies public branding works and
// settings list requires auth on the provider finalize path (C4.1).
func TestSettingsProviderServesBrandingAndAuth(t *testing.T) {
	a, st, settings, operations, authRepo := newTestEnv(t)
	plan := planWithSettings(t)
	set, err := kernel.RegisterContributions(context.Background(), plan, []kernel.Provider{New(a, settings, operations, nil, authRepo)})
	if err != nil {
		t.Fatalf("RegisterContributions: %v", err)
	}
	mux := http.NewServeMux()
	handler.Register(mux, a, st, operations, plan)
	for _, route := range set.Routes {
		mux.Handle(route.Method+" "+route.Pattern, route.Handler)
	}

	// Public branding (no auth) → 200.
	branding := httptest.NewRecorder()
	mux.ServeHTTP(branding, httptest.NewRequest(http.MethodGet, "/api/branding", nil))
	if branding.Code != http.StatusOK {
		t.Fatalf("branding = %d, want 200", branding.Code)
	}

	// Settings list anonymous → 401.
	anon := httptest.NewRecorder()
	mux.ServeHTTP(anon, httptest.NewRequest(http.MethodGet, "/api/settings", nil))
	if anon.Code != http.StatusUnauthorized {
		t.Fatalf("anonymous settings = %d, want 401", anon.Code)
	}

	// Login + list → 200.
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
	req := httptest.NewRequest(http.MethodGet, "/api/settings", nil)
	req.Header.Set("Authorization", "Bearer "+body.AccessToken)
	mux.ServeHTTP(list, req)
	if list.Code != http.StatusOK {
		t.Fatalf("authenticated settings = %d, want 200", list.Code)
	}
}
