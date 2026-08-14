package handler

import (
	"encoding/json"
	"errors"
	"fmt"
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
	notificationsschema "github.com/magicvr/schema-ui-core/apps/api/internal/modules/notifications/schema"
	authsession "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession"
	examplesschema "github.com/magicvr/schema-ui-core/apps/api/internal/modules/dev/examples/schema"
	filelibraryschema "github.com/magicvr/schema-ui-core/apps/api/internal/modules/filelibrary/schema"
	datadictionarystore "github.com/magicvr/schema-ui-core/apps/api/internal/modules/datadictionary/store"
	datadictionaryschema "github.com/magicvr/schema-ui-core/apps/api/internal/modules/datadictionary/schema"
	monitoringschema "github.com/magicvr/schema-ui-core/apps/api/internal/modules/systemmonitoring/schema"
	tasksschema "github.com/magicvr/schema-ui-core/apps/api/internal/modules/scheduledtasks/schema"
	tasksstore "github.com/magicvr/schema-ui-core/apps/api/internal/modules/scheduledtasks/store"
	recycleschema "github.com/magicvr/schema-ui-core/apps/api/internal/modules/recyclebin/schema"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
	recyclestore "github.com/magicvr/schema-ui-core/apps/api/internal/modules/recyclebin/store"
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
	captcha        *testCaptchaService
	recycle        *testRecycleService
}

const (
	testSeedUsername = "admin"
	testSeedPassword = "test-password"
	testBcryptCost   = 4 // cheap for tests
	testJWTSecret    = "test-secret"
)

// testUploadOpts lets tests override the upload policy (W7: config-driven
// limits) without touching package state; reset by the env cleanup.
var testUploadOpts []UploadOption

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
	// S-11 (GOAL-011): the fake is constructed BEFORE RegisterWithReadiness so
	// the login gate receives a live (non-nil) verifier — a typed-nil passed
	// through the variadic would satisfy the != nil check and panic on use.
	captchaService := newTestCaptchaService()
	recycleService := newTestRecycleService()
	plan := testAdminPlan(t)
	RegisterWithReadiness(mux, a, st, operations, plan, nil, captchaService)
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
	a.OnLockOpened = func(userID string) {
		NotifyAccountEvent(authRepository, userID, "account.locked", time.Now().UTC())
	}
	mountRoutes(AccountSelfRoutes(a, authRepository, operations, "admin.account", authRepository))
	mountRoutes(UserStateRoutes(a, authRepository, operations, "admin.account", authRepository))
	uploadDir := filepath.Join(t.TempDir(), "uploads")
	mountRoutes(NotificationRoutes(a, authRepository, "admin.notifications"))
	mountRoutes(ExportRoutes(a, authRepository, authRepository, operations, "admin.data-transfer"))
	mountRoutes(ImportRoutes(a, authRepository, operations, uploadDir, "admin.data-transfer"))
	mountRoutes(FileLibraryRoutes(a, uploadDir, operations, "admin.file-library"))
	mountRoutes(DictionaryRoutes(a, datadictionarystore.NewRepository(st), operations, "admin.data-dictionary"))
	mountRoutes(MonitoringRoutes(a, st, testAdminPlan(t), nil, filepath.Join(t.TempDir(), "monitor.db"), time.Now(), operations, "admin.system-monitoring"))
	scheduledTaskRunner := testTaskRunner{repository: tasksstore.NewRepository(st)}
	mountRoutes(ScheduledTaskRoutes(a, scheduledTaskRunner.repository, scheduledTaskRunner, operations, "admin.scheduled-tasks"))
	// S-11 (GOAL-011): the admin plan enables admin.login-captcha — the env
	// mounts its routes with the same fake service the login gate received
	// above (a test double: the module package imports handler, so handler
	// tests cannot import the real service; the module's own tests cover the
	// real store-backed service).
	mountRoutes(CaptchaRoutes(a, captchaService, operations, "admin.login-captcha"))
	// S-12 (GOAL-012): the admin plan enables admin.recycle-bin — the env
	// mounts its routes with a fake service; the real store-backed service is
	// covered by the module tests.
	recycleService = newTestRecycleService()
	mountRoutes(RecycleBinRoutes(a, recycleService, operations, "admin.recycle-bin"))
	mountRoutes(resourceRoutes(a, usersResourceWithNotifier(authRepository, operations, authRepository), "admin.users"))
	mountRoutes(resourceRoutes(a, rolesResource(authRepository, operations), "admin.roles"))
	RegisterSchemas(mux, testSchemaContributions())
	RegisterUpload(mux, a, uploadDir, testUploadOpts...)
	// testUploadOpts is reset after each test so per-test policy overrides
	// cannot leak into sibling tests.
	t.Cleanup(func() { testUploadOpts = nil })
	return &authTestEnv{
		mux: mux, a: a, st: st,
		authRepository: authRepository,
		operations:     operations,
		settings:       settings,
		uploadDir:      uploadDir,
		captcha:        captchaService,
		recycle:        recycleService,
	}
}

func testSchemaContributions() []kernel.PageContribution {
	contributors := []struct {
		moduleID  string
		documents map[string][]byte
	}{
		{accountschema.ModuleID, accountschema.SchemaDocuments()},
		{notificationsschema.ModuleID, notificationsschema.SchemaDocuments()},
		{examplesschema.ModuleID, examplesschema.SchemaDocuments()},
		{usersschema.ModuleID, usersschema.SchemaDocuments()},
		{rolesschema.ModuleID, rolesschema.SchemaDocuments()},
		{settingsschema.ModuleID, settingsschema.SchemaDocuments()},
		{activityschema.ModuleID, activityschema.SchemaDocuments()},
		{filelibraryschema.ModuleID, filelibraryschema.SchemaDocuments()},
		{datadictionaryschema.ModuleID, datadictionaryschema.SchemaDocuments()},
		{monitoringschema.ModuleID, monitoringschema.SchemaDocuments()},
		{tasksschema.ModuleID, tasksschema.SchemaDocuments()},
		{recycleschema.ModuleID, recycleschema.SchemaDocuments()},
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
// testTaskRunner records a manual run row for the scheduled-tasks surface
// (the real scheduler tick loop is exercised in the module package tests).
type testTaskRunner struct {
	repository *tasksstore.Repository
}

func (r testTaskRunner) Execute(task tasksstore.Task, now time.Time) error {
	finished := now
	return r.repository.RecordRun(tasksstore.TaskRun{
		ID: "run-test-" + fmt.Sprint(now.UnixNano()), TaskID: task.ID, Status: "ran",
		StartedAt: now, FinishedAt: &finished, Detail: "manual", CreatedAt: now,
	})
}

func (r testTaskRunner) HandlerKeys() []string { return []string{"system.noop"} }

// testCaptchaService is a controllable captcha double implementing both
// handler.CaptchaVerifier (login gate) and handler.CaptchaService (routes).
type testCaptchaService struct {
	required   bool
	generateID string
	verifyErr  error
}

func newTestCaptchaService() *testCaptchaService {
	return &testCaptchaService{generateID: "cap-test-1"}
}

func (s *testCaptchaService) Required() bool { return s.required }

func (s *testCaptchaService) Generate() (string, string, int64, error) {
	return s.generateID, "1 + 1 = ?", 300, nil
}

func (s *testCaptchaService) SetEnabled(enabled bool, now time.Time) error {
	s.required = enabled
	return nil
}

func (s *testCaptchaService) Verify(captchaID, answer string, now time.Time) error {
	if !s.required {
		return nil
	}
	if captchaID == s.generateID && answer == "2" {
		return nil
	}
	if s.verifyErr != nil {
		return s.verifyErr
	}
	return errors.New("captcha verification failed")
}

// testRecycleService is a controllable recycle double implementing
// handler.RecycleBinService. The real store-backed service is covered by the
// module tests; the handler tests use this fake to exercise the routes.
type testRecycleService struct {
	items   []RecycleItem
	restore func(id string) (map[string]any, error)
}

func newTestRecycleService() *testRecycleService {
	return &testRecycleService{}
}

func (s *testRecycleService) ListItems(resource, q string, page, pageSize int) ([]RecycleItem, int, error) {
	var out []RecycleItem
	for _, item := range s.items {
		if resource != "" && item.Resource != resource {
			continue
		}
		if q != "" && !strings.Contains(item.ResourceID, q) {
			continue
		}
		out = append(out, item)
	}
	return out, len(out), nil
}

func (s *testRecycleService) GetItem(id string) (*RecycleItem, error) {
	for i := range s.items {
		if s.items[i].ID == id {
			item := s.items[i]
			return &item, nil
		}
	}
	return nil, recyclestore.ErrItemNotFound
}

func (s *testRecycleService) Restore(id string, now time.Time) (map[string]any, error) {
	if s.restore != nil {
		return s.restore(id)
	}
	for i := range s.items {
		if s.items[i].ID == id {
			s.items[i].RestoredAt = now
			return s.items[i].Payload, nil
		}
	}
	return nil, recyclestore.ErrItemNotFound
}

func (s *testRecycleService) Purge(id string) error {
	for i := range s.items {
		if s.items[i].ID == id {
			s.items = append(s.items[:i], s.items[i+1:]...)
			return nil
		}
	}
	return recyclestore.ErrItemNotFound
}

func (s *testRecycleService) PurgeAll() (int, error) {
	before := len(s.items)
	s.items = nil
	return before, nil
}

func (s *testRecycleService) add(item RecycleItem) {
	s.items = append(s.items, item)
}
