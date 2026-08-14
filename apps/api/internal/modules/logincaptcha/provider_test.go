// Package logincaptcha provider tests (S-11 · GOAL-011 D-002 §4): the module
// registers the captcha routes, the page schema, permission keys,
// menu_captcha navigation and the fragment; the login gate stays off by
// default.
package logincaptcha

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
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/logincaptcha/store"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
)

func newCaptchaTestEnv(t *testing.T) (*auth.Authenticator, *Service, *operationlog.Repository) {
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
	return a, NewService(store.NewRepository(st)), operationlog.NewRepository(st)
}

func planWithCaptcha(t *testing.T) kernel.Plan {
	t.Helper()
	registry, err := kernel.NewRegistry(kernel.BuiltinModules())
	if err != nil {
		t.Fatalf("registry: %v", err)
	}
	plan, err := registry.Resolve([]string{
		"core.server-registration", "core.auth-session", "core.manifest-route",
		"core.navigation-capability", "core.schema-render", "core.operationlog",
		"admin.login-captcha",
	})
	if err != nil {
		t.Fatalf("resolve plan: %v", err)
	}
	return plan
}

func TestCaptchaProviderRegistersSurfaces(t *testing.T) {
	a, service, operations := newCaptchaTestEnv(t)
	provider := New(a, service, operations)
	set, err := kernel.RegisterContributions(context.Background(), planWithCaptcha(t), []kernel.Provider{provider})
	if err != nil {
		t.Fatalf("RegisterContributions: %v", err)
	}

	wantRoutes := []string{
		"GET /api/auth/captcha",
		"GET /api/captcha/settings", "PATCH /api/captcha/settings",
	}
	for _, key := range wantRoutes {
		if !slices.ContainsFunc(set.Routes, func(r kernel.RouteContribution) bool { return r.Key == key }) {
			t.Fatalf("route %q missing from provider set", key)
		}
	}
	if len(set.Pages) != 1 || set.Pages[0].PageID != "captcha" {
		t.Fatalf("pages = %+v, want captcha page", set.Pages)
	}
	for _, perm := range []string{"captcha.read", "captcha.write"} {
		if !slices.ContainsFunc(set.Permissions, func(p kernel.PermissionContribution) bool { return p.Permission == perm }) {
			t.Fatalf("permission %q missing", perm)
		}
	}
	if len(set.Navigation) != 1 || set.Navigation[0].NodeID != "menu_captcha" {
		t.Fatalf("navigation = %+v, want menu_captcha", set.Navigation)
	}
	if len(set.Fragments) != 1 || set.Fragments[0].FragmentID != "captcha" {
		t.Fatalf("fragments = %+v, want captcha fragment", set.Fragments)
	}
}

func TestCaptchaProviderServesSurfaces(t *testing.T) {
	a, service, operations := newCaptchaTestEnv(t)
	plan := planWithCaptcha(t)
	provider := New(a, service, operations)
	set, err := kernel.RegisterContributions(context.Background(), plan, []kernel.Provider{provider})
	if err != nil {
		t.Fatalf("RegisterContributions: %v", err)
	}
	mux := http.NewServeMux()
	handler.RegisterWithReadiness(mux, a, nil, operations, plan, nil, service)
	for _, route := range set.Routes {
		mux.Handle(route.Method+" "+route.Pattern, route.Handler)
	}

	preflight := httptest.NewRecorder()
	mux.ServeHTTP(preflight, httptest.NewRequest(http.MethodGet, "/api/auth/captcha", nil))
	if preflight.Code != http.StatusOK {
		t.Fatalf("preflight = %d: %s", preflight.Code, preflight.Body.String())
	}
	var pre struct {
		Enabled bool `json:"enabled"`
	}
	if err := json.NewDecoder(preflight.Body).Decode(&pre); err != nil {
		t.Fatalf("preflight decode: %v", err)
	}
	if pre.Enabled {
		t.Fatal("preflight must report disabled by default (D-001 §5)")
	}

	anon := httptest.NewRecorder()
	mux.ServeHTTP(anon, httptest.NewRequest(http.MethodGet, "/api/captcha/settings", nil))
	if anon.Code != http.StatusUnauthorized {
		t.Fatalf("anonymous settings = %d, want 401", anon.Code)
	}

	login := httptest.NewRecorder()
	mux.ServeHTTP(login, httptest.NewRequest(http.MethodPost, "/api/auth/login",
		strings.NewReader(`{"username":"admin","password":"test-password"}`)))
	if login.Code != http.StatusOK {
		t.Fatalf("login = %d: %s", login.Code, login.Body.String())
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
	settings := authReq(http.MethodGet, "/api/captcha/settings", "")
	if settings.Code != http.StatusOK || !strings.Contains(settings.Body.String(), `"enabled":false`) {
		t.Fatalf("settings = %d: %s", settings.Code, settings.Body.String())
	}
	patched := authReq(http.MethodPatch, "/api/captcha/settings", `{"enabled":true}`)
	if patched.Code != http.StatusOK || !strings.Contains(patched.Body.String(), `"enabled":true`) {
		t.Fatalf("patch = %d: %s", patched.Code, patched.Body.String())
	}
	gated := httptest.NewRecorder()
	mux.ServeHTTP(gated, httptest.NewRequest(http.MethodPost, "/api/auth/login",
		strings.NewReader(`{"username":"admin","password":"test-password"}`)))
	if gated.Code != http.StatusBadRequest || !strings.Contains(gated.Body.String(), "INVALID_CAPTCHA") {
		t.Fatalf("gated login = %d: %s", gated.Code, gated.Body.String())
	}
}
