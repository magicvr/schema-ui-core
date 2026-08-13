package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	accountschema "github.com/magicvr/schema-ui-core/apps/api/internal/modules/account/schema"
	activityschema "github.com/magicvr/schema-ui-core/apps/api/internal/modules/activity/schema"
	authsession "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession"
	examplesschema "github.com/magicvr/schema-ui-core/apps/api/internal/modules/dev/examples/schema"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
	rolesschema "github.com/magicvr/schema-ui-core/apps/api/internal/modules/roles/schema"
	settingsconfiguration "github.com/magicvr/schema-ui-core/apps/api/internal/modules/settings/configuration"
	settingsrepository "github.com/magicvr/schema-ui-core/apps/api/internal/modules/settings/repository"
	settingsschema "github.com/magicvr/schema-ui-core/apps/api/internal/modules/settings/schema"
	usersschema "github.com/magicvr/schema-ui-core/apps/api/internal/modules/users/schema"
	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
)

// authTestEnv wires a full handler mux backed by a temp SQLite store and an
// Authenticator. Tests that only exercise public routes can use env.mux
// directly; write-route tests log in first to obtain a Bearer access token.
type authTestEnv struct {
	mux            *http.ServeMux
	a              *auth.Authenticator
	st             *store.Store
	authRepository *authsession.Repository
	operations     *operationlog.Repository
	settings       *settingsrepository.Repository
	uploadDir      string
}

const (
	testSeedUsername = "admin"
	testSeedPassword = "test-password"
	testBcryptCost   = 4 // cheap for tests
	testJWTSecret    = "test-secret"
)

// newAuthTestEnv seeds the admin user and mounts the complete Admin plan (no
// dev-session fallback).
func newAuthTestEnv(t *testing.T) *authTestEnv {
	t.Helper()
	return newAuthTestEnvWith(t, false)
}

// newDevSessionTestEnv mounts the complete Admin plan with the explicit
// dev-session fallback enabled (acceptance M9: production must not enable this).
func newDevSessionTestEnv(t *testing.T) *authTestEnv {
	t.Helper()
	return newAuthTestEnvWith(t, true)
}

func newAuthTestEnvWith(t *testing.T, devSession bool) *authTestEnv {
	t.Helper()
	hash, err := auth.HashPassword(testSeedPassword, testBcryptCost)
	if err != nil {
		t.Fatalf("hash seed password: %v", err)
	}
	st, err := testsupport.OpenStore(filepath.Join(t.TempDir(), "test.db"), testSeedUsername, hash, true)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	authRepository := authsession.NewRepository(st)
	operations := operationlog.NewRepository(st)
	settings := settingsrepository.New(st)
	a := auth.NewWithRepository([]byte(testJWTSecret), 15*time.Minute, 30*24*time.Hour, authRepository, devSession)
	mux := http.NewServeMux()
	plan := testAdminPlan(t)
	Register(mux, a, st, operations, plan)
	// R6 C6.1: test env mounts the same resource-factory routes the module
	// providers register (behavior-identical to the production finalize path);
	// dead handler adapters MountProviderRoutes/RegisterSettings/RegisterActivity
	// are removed. Full RegisterContributions contract is covered by kernel and
	// composition tests.
	mountRoutes := func(routes []kernel.RouteContribution) {
		for _, r := range routes {
			mux.Handle(r.Method+" "+r.Pattern, r.Handler)
		}
	}
	mountRoutes(SettingsRoutes(a, settings, operations, "admin.settings", settingsconfiguration.Namespace))
	mountRoutes(ResourceRoutes(a, operationsResource(operations), "admin.activity"))
	mountRoutes(AccountSelfRoutes(a, authRepository, operations, "admin.account"))
	mountRoutes(UserStateRoutes(a, authRepository, operations, "admin.account"))
	uploadDir := filepath.Join(t.TempDir(), "uploads")
	mountRoutes(ExportRoutes(a, authRepository, authRepository, operations, "admin.data-transfer"))
	mountRoutes(ImportRoutes(a, authRepository, operations, uploadDir, "admin.data-transfer"))
	mountRoutes(resourceRoutes(a, usersResource(authRepository, operations), "admin.users"))
	mountRoutes(resourceRoutes(a, rolesResource(authRepository, operations), "admin.roles"))
	RegisterSchemas(mux, testSchemaContributions())
	RegisterUpload(mux, a, uploadDir)
	return &authTestEnv{
		mux: mux, a: a, st: st,
		authRepository: authRepository,
		operations:     operations,
		settings:       settings,
		uploadDir:      uploadDir,
	}
}

func testSchemaContributions() []kernel.PageContribution {
	contributors := []struct {
		moduleID  string
		documents map[string][]byte
	}{
		{accountschema.ModuleID, accountschema.SchemaDocuments()},
		{examplesschema.ModuleID, examplesschema.SchemaDocuments()},
		{usersschema.ModuleID, usersschema.SchemaDocuments()},
		{rolesschema.ModuleID, rolesschema.SchemaDocuments()},
		{settingsschema.ModuleID, settingsschema.SchemaDocuments()},
		{activityschema.ModuleID, activityschema.SchemaDocuments()},
	}
	var pages []kernel.PageContribution
	for _, contributor := range contributors {
		for pageID, document := range contributor.documents {
			pages = append(pages, kernel.PageContribution{
				ContributionIdentity: kernel.ContributionIdentity{ModuleID: contributor.moduleID, Key: pageID},
				PageID:               pageID, Owner: contributor.moduleID, Document: document,
			})
		}
	}
	return pages
}

func testAdminPlan(t *testing.T) kernel.Plan {
	t.Helper()
	resolution, err := kernel.ResolveProfile("admin", nil)
	if err != nil {
		t.Fatalf("resolve admin profile: %v", err)
	}
	registry, err := kernel.NewRegistry(kernel.BuiltinModules())
	if err != nil {
		t.Fatalf("build module registry: %v", err)
	}
	plan, err := registry.Resolve(resolution.Modules)
	if err != nil {
		t.Fatalf("resolve admin module plan: %v", err)
	}
	return plan
}

// login performs POST /api/auth/login and returns the access token.
func (e *authTestEnv) login(t *testing.T, username, password string) string {
	t.Helper()
	body := `{"username":` + quote(username) + `,"password":` + quote(password) + `}`
	code, out := sendJSON(t, e.mux, http.MethodPost, "/api/auth/login", body)
	if code != http.StatusOK {
		t.Fatalf("login status = %d, want 200: %v", code, out)
	}
	tok, _ := out["accessToken"].(string)
	if tok == "" {
		t.Fatalf("accessToken missing in %v", out)
	}
	return tok
}

// addUser inserts a user directly into the store for permission tests.
func (e *authTestEnv) addUser(t *testing.T, username, password string, roles []string) {
	t.Helper()
	hash, err := auth.HashPassword(password, testBcryptCost)
	if err != nil {
		t.Fatalf("hash %s password: %v", username, err)
	}
	now := time.Now().UTC()
	if err := e.authRepository.CreateUser(authsession.User{
		ID:           "user-" + username,
		Username:     username,
		Name:         username,
		Roles:        roles,
		PasswordHash: hash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}); err != nil {
		t.Fatalf("create user %s: %v", username, err)
	}
}

// bearer returns a request with the given access token attached.
func bearer(t *testing.T, token, method, path string, body string) *http.Request {
	t.Helper()
	var reader *strings.Reader
	if body == "" {
		reader = strings.NewReader("")
	} else {
		reader = strings.NewReader(body)
	}
	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Authorization", "Bearer "+token)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	return req
}

func quote(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}

func getJSON(t *testing.T, mux *http.ServeMux, path string) (int, map[string]any) {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	var body map[string]any
	if rr.Body.Len() > 0 {
		if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
			t.Fatalf("decode %q: %v", rr.Body.String(), err)
		}
	}
	return rr.Code, body
}

// adminToken logs in as the seeded admin and returns the Bearer access token.
func adminToken(t *testing.T, env *authTestEnv) string {
	t.Helper()
	return env.login(t, testSeedUsername, testSeedPassword)
}

// getResource fetches a resource path as the seeded admin (GOAL-006 S4: reads
// are authenticated and permission-gated, so a Bearer token is required).
func getResource(t *testing.T, env *authTestEnv, path string) (int, map[string]any) {
	t.Helper()
	req := bearer(t, adminToken(t, env), http.MethodGet, path, "")
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	var body map[string]any
	if rr.Body.Len() > 0 {
		if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
			t.Fatalf("decode %q: %v", rr.Body.String(), err)
		}
	}
	return rr.Code, body
}

func sendJSON(t *testing.T, mux *http.ServeMux, method, path, body string) (int, map[string]any) {
	t.Helper()
	var reader *strings.Reader
	if body == "" {
		reader = strings.NewReader("")
	} else {
		reader = strings.NewReader(body)
	}
	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	var out map[string]any
	if rr.Body.Len() > 0 && rr.Header().Get("Content-Type") != "" {
		_ = json.NewDecoder(rr.Body).Decode(&out)
	}
	return rr.Code, out
}

func containsString(haystack, needle string) bool {
	return strings.Contains(haystack, needle)
}