// Package logincaptcha provider tests (S-11 · GOAL-011 D-002 §4): the module
// registers the captcha routes, the page schema, permission keys,
// menu_captcha navigation and the fragment; the login gate stays off by
// default.
package logincaptcha

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/handler"
	"github.com/magicvr/schema-ui-core/apps/api/internal/ratelimit"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	authsession "github.com/magicvr/schema-ui-core/apps/api/modules/authsession"
	"github.com/magicvr/schema-ui-core/apps/api/modules/logincaptcha/store"
	"github.com/magicvr/schema-ui-core/apps/api/modules/operationlog"
	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
)

func newCaptchaTestEnv(t *testing.T) (*auth.Authenticator, *Service, *operationlog.Repository, *authsession.Repository) {
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
	return a, NewService(store.NewRepository(st)), operationlog.NewRepository(st), repository
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
	a, service, operations, _ := newCaptchaTestEnv(t)
	provider := New(a, service, operations, ratelimit.NewProvider())
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
	if len(set.Pages) != 0 {
		t.Fatalf("pages = %+v, want none (D-003: switch lives in settings)", set.Pages)
	}
	for _, perm := range []string{"captcha.read", "captcha.write"} {
		if !slices.ContainsFunc(set.Permissions, func(p kernel.PermissionContribution) bool { return p.Permission == perm }) {
			t.Fatalf("permission %q missing", perm)
		}
	}
	if len(set.Navigation) != 0 {
		t.Fatalf("navigation = %+v, want none (D-003: no menu item)", set.Navigation)
	}
	if len(set.Fragments) != 0 {
		t.Fatalf("fragments = %+v, want none (D-003)", set.Fragments)
	}
}

func TestCaptchaProviderServesSurfaces(t *testing.T) {
	a, service, operations, _ := newCaptchaTestEnv(t)
	plan := planWithCaptcha(t)
	provider := New(a, service, operations, ratelimit.NewProvider())
	set, err := kernel.RegisterContributions(context.Background(), plan, []kernel.Provider{provider})
	if err != nil {
		t.Fatalf("RegisterContributions: %v", err)
	}
	mux := http.NewServeMux()
	handler.RegisterWithReadiness(mux, a, nil, operations, plan, nil, ratelimit.NewProvider(), service)
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
	if settings.Code != http.StatusOK || !strings.Contains(settings.Body.String(), `"enabled":"false"`) {
		t.Fatalf("settings = %d: %s", settings.Code, settings.Body.String())
	}
	patched := authReq(http.MethodPatch, "/api/captcha/settings", `{"enabled":"true"}`)
	if patched.Code != http.StatusOK || !strings.Contains(patched.Body.String(), `"enabled":"true"`) {
		t.Fatalf("patch = %d: %s", patched.Code, patched.Body.String())
	}
	gated := httptest.NewRecorder()
	mux.ServeHTTP(gated, httptest.NewRequest(http.MethodPost, "/api/auth/login",
		strings.NewReader(`{"username":"admin","password":"test-password"}`)))
	if gated.Code != http.StatusBadRequest || !strings.Contains(gated.Body.String(), "INVALID_CAPTCHA") {
		t.Fatalf("gated login = %d: %s", gated.Code, gated.Body.String())
	}
}

// F-005 (grok A-003): real store-backed service end-to-end — challenge
// generate → answer → login 200; settings 403 for a non-admin; captcha
// failures never count toward the lockout budget.
func TestCaptchaRealServiceChallengeLogin(t *testing.T) {
	a, service, operations, _ := newCaptchaTestEnv(t)
	plan := planWithCaptcha(t)
	provider := New(a, service, operations, ratelimit.NewProvider())
	set, err := kernel.RegisterContributions(context.Background(), plan, []kernel.Provider{provider})
	if err != nil {
		t.Fatalf("RegisterContributions: %v", err)
	}
	mux := http.NewServeMux()
	handler.RegisterWithReadiness(mux, a, nil, operations, plan, nil, ratelimit.NewProvider(), service)
	for _, route := range set.Routes {
		mux.Handle(route.Method+" "+route.Pattern, route.Handler)
	}
	if err := service.SetEnabled(true, time.Now().UTC()); err != nil {
		t.Fatalf("enable: %v", err)
	}
	// Preflight issues a real persisted challenge.
	pre := httptest.NewRecorder()
	mux.ServeHTTP(pre, httptest.NewRequest(http.MethodGet, "/api/auth/captcha", nil))
	if pre.Code != http.StatusOK {
		t.Fatalf("preflight = %d: %s", pre.Code, pre.Body.String())
	}
	var body struct {
		Enabled   bool `json:"enabled"`
		Challenge *struct {
			ID       string `json:"id"`
			Question string `json:"question"`
		} `json:"challenge"`
	}
	if err := json.NewDecoder(pre.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if !body.Enabled || body.Challenge == nil {
		t.Fatalf("preflight = %+v, want enabled challenge", body)
	}
	// Solve the arithmetic question.
	var aNum, bNum int
	var op string
	if _, err := fmt.Sscanf(body.Challenge.Question, "%d %s %d = ?", &aNum, &op, &bNum); err != nil {
		t.Fatalf("parse question %q: %v", body.Challenge.Question, err)
	}
	answer := ""
	if op == "+" {
		answer = strconv.Itoa(aNum + bNum)
	} else {
		answer = strconv.Itoa(aNum - bNum)
	}
	login := httptest.NewRecorder()
	payload := fmt.Sprintf(`{"username":"admin","password":"test-password","captchaId":%q,"captchaAnswer":%q}`, body.Challenge.ID, answer)
	mux.ServeHTTP(login, httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(payload)))
	if login.Code != http.StatusOK {
		t.Fatalf("login with solved challenge = %d: %s", login.Code, login.Body.String())
	}
}

func TestCaptchaRealServiceSettingsForbiddenForNonAdmin(t *testing.T) {
	a, service, operations, authRepository := newCaptchaTestEnv(t)
	plan := planWithCaptcha(t)
	provider := New(a, service, operations, ratelimit.NewProvider())
	set, err := kernel.RegisterContributions(context.Background(), plan, []kernel.Provider{provider})
	if err != nil {
		t.Fatalf("RegisterContributions: %v", err)
	}
	mux := http.NewServeMux()
	handler.RegisterWithReadiness(mux, a, nil, operations, plan, nil, ratelimit.NewProvider(), service)
	for _, route := range set.Routes {
		mux.Handle(route.Method+" "+route.Pattern, route.Handler)
	}
	// Create an editor user (no captcha.read/write in the mvp system data set).
	hash, err := auth.HashPassword("editor-pass", 4)
	if err != nil {
		t.Fatalf("hash: %v", err)
	}
	now := time.Now().UTC()
	if err := authRepository.CreateUser(authsession.User{
		ID: "user-editor", Username: "editor", Name: "Editor",
		Roles: []string{"editor"}, PasswordHash: hash, CreatedAt: now, UpdatedAt: now,
	}); err != nil {
		t.Fatalf("create editor: %v", err)
	}
	login := httptest.NewRecorder()
	mux.ServeHTTP(login, httptest.NewRequest(http.MethodPost, "/api/auth/login",
		strings.NewReader(`{"username":"editor","password":"editor-pass"}`)))
	if login.Code != http.StatusOK {
		t.Fatalf("editor login = %d: %s", login.Code, login.Body.String())
	}
	var body struct {
		AccessToken string `json:"accessToken"`
	}
	if err := json.NewDecoder(login.Body).Decode(&body); err != nil || body.AccessToken == "" {
		t.Fatalf("editor login body missing accessToken")
	}
	req := httptest.NewRequest(http.MethodGet, "/api/captcha/settings", nil)
	req.Header.Set("Authorization", "Bearer "+body.AccessToken)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("editor settings = %d, want 403 (F-005)", rec.Code)
	}
}

func TestCaptchaRealServiceFailuresDoNotLock(t *testing.T) {
	a, service, operations, _ := newCaptchaTestEnv(t)
	plan := planWithCaptcha(t)
	provider := New(a, service, operations, ratelimit.NewProvider())
	set, err := kernel.RegisterContributions(context.Background(), plan, []kernel.Provider{provider})
	if err != nil {
		t.Fatalf("RegisterContributions: %v", err)
	}
	mux := http.NewServeMux()
	handler.RegisterWithReadiness(mux, a, nil, operations, plan, nil, ratelimit.NewProvider(), service)
	for _, route := range set.Routes {
		mux.Handle(route.Method+" "+route.Pattern, route.Handler)
	}
	if err := service.SetEnabled(true, time.Now().UTC()); err != nil {
		t.Fatalf("enable: %v", err)
	}
	// Many captcha-rejected attempts (above the 20-attempt rate limit and the
	// lockout budget) must never count: they never reach credential validation.
	for i := 0; i < 25; i++ {
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/auth/login",
			strings.NewReader(`{"username":"admin","password":"test-password","captchaId":"cap-missing","captchaAnswer":"1"}`)))
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("attempt %d = %d: %s", i, rec.Code, rec.Body.String())
		}
	}
	// Disable the gate: the account must not be locked or rate-limited.
	if err := service.SetEnabled(false, time.Now().UTC()); err != nil {
		t.Fatalf("disable: %v", err)
	}
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/auth/login",
		strings.NewReader(`{"username":"admin","password":"test-password"}`)))
	if rec.Code != http.StatusOK {
		t.Fatalf("login after captcha-only failures = %d: %s (F-005)", rec.Code, rec.Body.String())
	}
}
