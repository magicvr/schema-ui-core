package composition

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/config"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	authsession "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
	settingsrepository "github.com/magicvr/schema-ui-core/apps/api/internal/modules/settings/repository"
	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
)

func compositionCount(t *testing.T, st *store.Store, query string, args ...any) int {
	t.Helper()
	var count int
	if err := st.WithTx(context.Background(), func(tx *sql.Tx) error {
		return tx.QueryRow(query, args...).Scan(&count)
	}); err != nil {
		t.Fatalf("query %q: %v", query, err)
	}
	return count
}

func testMux(a *auth.Authenticator, st *store.Store, plan kernel.Plan, gate *readinessGate) (*http.ServeMux, error) {
	return newMux(
		&config.Config{DBPath: "test.db"},
		a,
		st,
		authsession.NewRepository(st),
		operationlog.NewRepository(st),
		settingsrepository.New(st),
		plan,
		gate,
	)
}

func TestResolvePlanUsesConfiguredProfileAndRejectsMissingDependencies(t *testing.T) {
	admin := config.Load()
	admin.ProfileName = "admin"
	admin.ModulesEnabled = nil
	admin.ProfileError = nil
	plan, err := ResolvePlan(admin)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Modules) < 8 {
		t.Fatalf("admin plan has %d modules, want the core and admin set", len(plan.Modules))
	}

	broken := &config.Config{
		ProfileName:    "custom",
		ModulesEnabled: []string{"admin.users"},
	}
	if _, err := ResolvePlan(broken); err == nil {
		t.Fatal("module graph with missing dependencies must fail closed")
	}
}

func TestNewMuxPublishesOnlySelectedProfileManifestPages(t *testing.T) {
	cfg := &config.Config{ProfileName: "mvp"}
	plan, err := ResolvePlan(cfg)
	if err != nil {
		t.Fatal(err)
	}
	st, err := testsupport.OpenStore(":memory:", "admin", "hash", false)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	a := auth.New([]byte("test-secret"), 0, 0, st, false)
	mux, err := testMux(a, st, plan, &readinessGate{})
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/.well-known/schema-ui/app-manifest.json", nil))
	if response.Code != http.StatusOK {
		t.Fatalf("manifest status = %d, body=%s", response.Code, response.Body.String())
	}
	var document struct {
		Pages []struct {
			PageID string `json:"pageId"`
		} `json:"pages"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &document); err != nil {
		t.Fatal(err)
	}
	pageIDs := map[string]bool{}
	for _, page := range document.Pages {
		pageIDs[page.PageID] = true
	}
	if !pageIDs["users"] || !pageIDs["roles"] {
		t.Fatalf("mvp admin pages missing: %v", pageIDs)
	}
	if pageIDs["settings"] || pageIDs["activity"] {
		t.Fatalf("disabled admin pages leaked: %v", pageIDs)
	}
	for pageID := range pageIDs {
		schemaResponse := httptest.NewRecorder()
		mux.ServeHTTP(schemaResponse, httptest.NewRequest(http.MethodGet, "/api/schema/"+pageID, nil))
		if schemaResponse.Code != http.StatusOK {
			t.Fatalf("manifest page %q schema status = %d, body=%s", pageID, schemaResponse.Code, schemaResponse.Body.String())
		}
	}
}

func TestNewMuxProjectsProfileRoutesAndSchemasFromOnePlan(t *testing.T) {
	cases := []struct {
		name             string
		profile          string
		settingsRoute    int
		operationsRoute  int
		brandingRoute    int
		settingsSchema   int
		activitySchema   int
		wantSettingsPage bool
		wantActivityPage bool
	}{
		{
			name:             "mvp disables optional settings and activity",
			profile:          "mvp",
			settingsRoute:    http.StatusNotFound,
			operationsRoute:  http.StatusNotFound,
			brandingRoute:    http.StatusOK, // public bootstrap remains available without admin.settings
			settingsSchema:   http.StatusNotFound,
			activitySchema:   http.StatusNotFound,
			wantSettingsPage: false,
			wantActivityPage: false,
		},
		{
			name:             "admin enables optional settings and activity",
			profile:          "admin",
			settingsRoute:    http.StatusUnauthorized,
			operationsRoute:  http.StatusUnauthorized,
			brandingRoute:    http.StatusOK,
			settingsSchema:   http.StatusOK,
			activitySchema:   http.StatusOK,
			wantSettingsPage: true,
			wantActivityPage: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &config.Config{ProfileName: tc.profile}
			plan, err := ResolvePlan(cfg)
			if err != nil {
				t.Fatal(err)
			}
			st, err := testsupport.OpenStore(":memory:", "admin", "hash", false)
			if err != nil {
				t.Fatal(err)
			}
			defer st.Close()
			a := auth.New([]byte("test-secret"), 0, 0, st, false)
			mux, err := testMux(a, st, plan, &readinessGate{})
			if err != nil {
				t.Fatal(err)
			}

			assertStatus := func(path string, want int) {
				t.Helper()
				recorder := httptest.NewRecorder()
				mux.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, path, nil))
				if recorder.Code != want {
					t.Fatalf("GET %s status = %d, want %d; body=%s", path, recorder.Code, want, recorder.Body.String())
				}
			}
			assertStatus("/api/settings", tc.settingsRoute)
			assertStatus("/api/operations", tc.operationsRoute)
			assertStatus("/api/branding", tc.brandingRoute)
			assertStatus("/api/schema/settings", tc.settingsSchema)
			assertStatus("/api/schema/activity", tc.activitySchema)

			manifestResponse := httptest.NewRecorder()
			mux.ServeHTTP(manifestResponse, httptest.NewRequest(http.MethodGet, "/.well-known/schema-ui/app-manifest.json", nil))
			if manifestResponse.Code != http.StatusOK {
				t.Fatalf("manifest status = %d, want 200", manifestResponse.Code)
			}
			var document struct {
				Pages []struct {
					PageID string `json:"pageId"`
				} `json:"pages"`
			}
			if err := json.Unmarshal(manifestResponse.Body.Bytes(), &document); err != nil {
				t.Fatal(err)
			}
			pageIDs := map[string]bool{}
			for _, page := range document.Pages {
				pageIDs[page.PageID] = true
			}
			if pageIDs["settings"] != tc.wantSettingsPage || pageIDs["activity"] != tc.wantActivityPage {
				t.Fatalf("optional manifest pages = settings:%v activity:%v, want settings:%v activity:%v", pageIDs["settings"], pageIDs["activity"], tc.wantSettingsPage, tc.wantActivityPage)
			}
		})
	}
}

// TestManifestHomePageRefDerivation pins the W1 home-page derivation table
// (D-003 §2) and the product-surface hygiene of default profiles (S4/S5):
// mvp/admin publish no dev.examples pages and home -> the first admin page;
// explicitly enabling dev.examples restores the overview home and examples.
func TestManifestHomePageRefDerivation(t *testing.T) {
	type manifestApp struct {
		HomePageRef string `json:"homePageRef"`
	}
	type manifestDoc struct {
		App   manifestApp `json:"app"`
		Pages []struct {
			PageID string `json:"pageId"`
		} `json:"pages"`
		Navigation struct {
			Sidebar []json.RawMessage `json:"sidebar"`
		} `json:"navigation"`
	}

	fetch := func(t *testing.T, plan kernel.Plan) (manifestDoc, *http.ServeMux) {
		t.Helper()
		st, err := testsupport.OpenStore(":memory:", "admin", "hash", false)
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = st.Close() })
		a := auth.New([]byte("test-secret"), 0, 0, st, false)
		mux, err := testMux(a, st, plan, &readinessGate{})
		if err != nil {
			t.Fatal(err)
		}
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/.well-known/schema-ui/app-manifest.json", nil))
		if rec.Code != http.StatusOK {
			t.Fatalf("manifest status = %d: %s", rec.Code, rec.Body.String())
		}
		var doc manifestDoc
		if err := json.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
			t.Fatal(err)
		}
		return doc, mux
	}

	// Default mvp: first admin functional page (users), no examples surface.
	mvpPlan, err := ResolvePlan(&config.Config{ProfileName: "mvp"})
	if err != nil {
		t.Fatal(err)
	}
	mvpDoc, mvpMux := fetch(t, mvpPlan)
	if mvpDoc.App.HomePageRef != "users" {
		t.Fatalf("mvp homePageRef = %q, want users", mvpDoc.App.HomePageRef)
	}
	for _, page := range mvpDoc.Pages {
		if page.PageID == "overview" || page.PageID == "data-table" || page.PageID == "form-controls" {
			t.Fatalf("mvp default manifest leaks example page %q (S5)", page.PageID)
		}
	}
	// Disabled examples: example schema endpoint 404s (S5).
	overviewSchema := httptest.NewRecorder()
	mvpMux.ServeHTTP(overviewSchema, httptest.NewRequest(http.MethodGet, "/api/schema/overview", nil))
	if overviewSchema.Code != http.StatusNotFound {
		t.Fatalf("mvp /api/schema/overview status = %d, want 404 (S5)", overviewSchema.Code)
	}
	// Disabled examples: no "Examples" sidebar navigation group (S5, F-003).
	for _, raw := range mvpDoc.Navigation.Sidebar {
		var item struct {
			Label string `json:"label"`
		}
		if err := json.Unmarshal(raw, &item); err != nil {
			t.Fatal(err)
		}
		if item.Label == "Examples" {
			t.Fatalf("mvp default manifest exposes Examples navigation group (S5)")
		}
	}

	// Default admin: still users (users precedes roles in declaration order).
	adminPlan, err := ResolvePlan(&config.Config{ProfileName: "admin"})
	if err != nil {
		t.Fatal(err)
	}
	adminDoc, _ := fetch(t, adminPlan)
	if adminDoc.App.HomePageRef != "users" {
		t.Fatalf("admin homePageRef = %q, want users", adminDoc.App.HomePageRef)
	}

	// mvp + dev.examples (dogfood): home -> overview, examples restored.
	registry, err := kernel.NewRegistry(kernel.BuiltinModules())
	if err != nil {
		t.Fatal(err)
	}
	examplesPlan, err := registry.Resolve(append(mvpPlan.IDs(), "dev.examples"))
	if err != nil {
		t.Fatal(err)
	}
	examplesDoc, examplesMux := fetch(t, examplesPlan)
	if examplesDoc.App.HomePageRef != "overview" {
		t.Fatalf("dev.examples enabled homePageRef = %q, want overview", examplesDoc.App.HomePageRef)
	}
	seen := map[string]bool{}
	for _, page := range examplesDoc.Pages {
		seen[page.PageID] = true
	}
	if !seen["overview"] || !seen["data-table"] {
		t.Fatalf("dev.examples enabled manifest missing example pages: %v", seen)
	}
	// Enabled examples: example schema endpoint serves again.
	overviewSchemaOn := httptest.NewRecorder()
	examplesMux.ServeHTTP(overviewSchemaOn, httptest.NewRequest(http.MethodGet, "/api/schema/overview", nil))
	if overviewSchemaOn.Code != http.StatusOK {
		t.Fatalf("dev.examples enabled /api/schema/overview status = %d, want 200", overviewSchemaOn.Code)
	}
}

// TestDeriveHomePageRefBranches pins the remaining decision-table branches of
// deriveHomePageRef (D-003 §2, F-002 of the independent wave audit): non-users
// admin order, no-admin-with-page fallback, and the no-page omit case.
func TestDeriveHomePageRefBranches(t *testing.T) {
	mkPlan := func(modules ...kernel.Module) kernel.Plan {
		return kernel.Plan{Modules: modules}
	}
	pages := func(id string, pageIDs ...string) kernel.Module {
		return kernel.Module{ID: id, Contributions: kernel.ContributionKeys{Pages: pageIDs}}
	}
	cases := []struct {
		name string
		plan kernel.Plan
		want string
	}{
		{name: "dev.examples wins", plan: mkPlan(pages("dev.examples", "overview", "data-table"), pages("admin.users", "users")), want: "overview"},
		{name: "users precedes roles", plan: mkPlan(pages("admin.roles", "roles"), pages("admin.users", "users")), want: "users"},
		{name: "roles only", plan: mkPlan(pages("admin.roles", "roles")), want: "roles"},
		{name: "activity only", plan: mkPlan(pages("admin.activity", "activity")), want: "activity"},
		{name: "no admin, page-bearing module", plan: mkPlan(pages("custom.foo", "foo")), want: "foo"},
		{name: "no pages omits home", plan: mkPlan(pages("custom.empty")), want: ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := deriveHomePageRef(tc.plan); got != tc.want {
				t.Fatalf("deriveHomePageRef = %q, want %q", got, tc.want)
			}
		})
	}
}

// TestDemoProfileManifest pins the W2 demo Profile (GOAL-003): demo = mvp +
// dev.examples, so the manifest exposes both the mvp pages and the examples,
// home -> overview, while mvp/admin production defaults stay examples-free.
func TestDemoProfileManifest(t *testing.T) {
	plan, err := ResolvePlan(&config.Config{ProfileName: "demo"})
	if err != nil {
		t.Fatal(err)
	}
	if !plan.HasModule("dev.examples") || !plan.HasModule("admin.users") {
		t.Fatalf("demo plan must include dev.examples and mvp admin modules: %v", plan.IDs())
	}
	st, err := testsupport.OpenStore(":memory:", "admin", "hash", false)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	a := auth.New([]byte("test-secret"), 0, 0, st, false)
	mux, err := testMux(a, st, plan, &readinessGate{})
	if err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/.well-known/schema-ui/app-manifest.json", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("manifest status = %d: %s", rec.Code, rec.Body.String())
	}
	var doc struct {
		App struct {
			HomePageRef string `json:"homePageRef"`
		} `json:"app"`
		Pages []struct {
			PageID string `json:"pageId"`
		} `json:"pages"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &doc); err != nil {
		t.Fatal(err)
	}
	if doc.App.HomePageRef != "overview" {
		t.Fatalf("demo homePageRef = %q, want overview", doc.App.HomePageRef)
	}
	seen := map[string]bool{}
	for _, page := range doc.Pages {
		seen[page.PageID] = true
	}
	for _, want := range []string{"users", "roles", "overview", "data-table", "form-controls"} {
		if !seen[want] {
			t.Fatalf("demo manifest missing page %q: %v", want, seen)
		}
	}
	if seen["settings"] || seen["activity"] {
		t.Fatalf("demo manifest leaks non-mvp pages: %v", seen)
	}
	// Example schema endpoint serves under demo.
	overviewSchema := httptest.NewRecorder()
	mux.ServeHTTP(overviewSchema, httptest.NewRequest(http.MethodGet, "/api/schema/overview", nil))
	if overviewSchema.Code != http.StatusOK {
		t.Fatalf("demo /api/schema/overview status = %d, want 200", overviewSchema.Code)
	}
	// mvp/admin still exclude examples (S3, W1 hygiene).
	for _, profile := range []string{"mvp", "admin"} {
		p, err := ResolvePlan(&config.Config{ProfileName: profile})
		if err != nil {
			t.Fatal(err)
		}
		for _, id := range p.IDs() {
			if id == "dev.examples" {
				t.Fatalf("%s default must not include dev.examples (S3)", profile)
			}
		}
	}
}

func TestSystemDataReconcileUsesFinalizedProfileContributions(t *testing.T) {
	tests := []struct {
		profile         string
		wantPermissions int
		wantNavigation  int
	}{
		{profile: "mvp", wantPermissions: 5, wantNavigation: 2},
		{profile: "admin", wantPermissions: 8, wantNavigation: 4},
	}
	for _, tt := range tests {
		t.Run(tt.profile, func(t *testing.T) {
			plan, err := ResolvePlan(&config.Config{ProfileName: tt.profile})
			if err != nil {
				t.Fatal(err)
			}
			st, err := testsupport.OpenStore(":memory:", "admin", "hash", false)
			if err != nil {
				t.Fatal(err)
			}
			defer st.Close()
			a := auth.New([]byte("test-secret"), 0, 0, st, false)
			if _, err := testMux(a, st, plan, &readinessGate{}); err != nil {
				t.Fatal(err)
			}
			if got := compositionCount(t, st, `SELECT COUNT(*) FROM permissions`); got != tt.wantPermissions {
				t.Fatalf("permissions = %d, want %d", got, tt.wantPermissions)
			}
			if got := compositionCount(t, st, `SELECT COUNT(*) FROM menu_items`); got != tt.wantNavigation {
				t.Fatalf("navigation = %d, want %d", got, tt.wantNavigation)
			}
			if got := compositionCount(t, st, `SELECT COUNT(*) FROM system_data_reconcile`); got != 1+tt.wantPermissions+tt.wantNavigation {
				t.Fatalf("system-data ledger = %d, want %d", got, 1+tt.wantPermissions+tt.wantNavigation)
			}
			if err := st.SystemDataReady(); err != nil {
				t.Fatalf("system-data readiness: %v", err)
			}
		})
	}
}

func TestSystemDataReconcilePreservesDisabledProfileData(t *testing.T) {
	st, err := testsupport.OpenStore(filepath.Join(t.TempDir(), "profile-downgrade.db"), "admin", "hash", false)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	a := auth.New([]byte("test-secret"), 0, 0, st, false)

	adminPlan, err := ResolvePlan(&config.Config{ProfileName: "admin"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := testMux(a, st, adminPlan, &readinessGate{}); err != nil {
		t.Fatal(err)
	}
	mvpPlan, err := ResolvePlan(&config.Config{ProfileName: "mvp"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := testMux(a, st, mvpPlan, &readinessGate{}); err != nil {
		t.Fatal(err)
	}

	if got := compositionCount(t, st, `SELECT COUNT(*) FROM permissions WHERE key IN ('settings.read', 'settings.write', 'operations.read')`); got != 3 {
		t.Fatalf("disabled-profile permissions retained = %d, want 3", got)
	}
	if got := compositionCount(t, st, `SELECT COUNT(*) FROM menu_items WHERE feature_key IN ('menu_settings', 'menu_activity')`); got != 2 {
		t.Fatalf("disabled-profile navigation retained = %d, want 2", got)
	}
	if got := compositionCount(t, st, `SELECT COUNT(*) FROM system_data_reconcile WHERE module_id IN ('admin.settings', 'admin.activity')`); got != 5 {
		t.Fatalf("disabled-profile ledger retained = %d, want 5", got)
	}
	if got := compositionCount(t, st, `SELECT COUNT(*) FROM system_data_grants WHERE module_id IN ('admin.settings', 'admin.activity')`); got != 7 {
		t.Fatalf("disabled-profile managed grants retained = %d, want 7", got)
	}
}

func TestMVPRecoveryRestoresOptionalModuleDataAndCoreReadiness(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "r3-recovery.db")
	snapshotPath := filepath.Join(dir, "pre-r3-snapshot.db")
	restoredPath := filepath.Join(dir, "r3-restored.db")
	st, err := testsupport.OpenStore(path, "admin", "hash", true)
	if err != nil {
		t.Fatal(err)
	}
	updatedAt := time.Date(2026, 8, 5, 0, 0, 0, 0, time.UTC)
	settingsRepository := settingsrepository.New(st)
	operationsRepository := operationlog.NewRepository(st)
	if _, err := settingsRepository.UpdateSiteSettings("Recovered Admin", "/assets/logo.svg", updatedAt); err != nil {
		t.Fatalf("update settings: %v", err)
	}
	recordID := "default"
	detail := `{"siteTitle":"Recovered Admin"}`
	if err := operationsRepository.RecordOperation(operationlog.Operation{
		ID:        "r3-settings-update",
		Event:     operationlog.EventSettingsUpdate,
		ActorID:   "user-admin",
		ActorName: "Admin",
		RecordID:  &recordID,
		Detail:    &detail,
		CreatedAt: updatedAt,
	}); err != nil {
		t.Fatalf("record operation: %v", err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}

	// Capture the last known-good database after the Admin write. The recovery
	// drill below mutates the live copy, then restores this exact snapshot.
	snapshot, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read recovery snapshot: %v", err)
	}
	if err := os.WriteFile(snapshotPath, snapshot, 0o644); err != nil {
		t.Fatalf("write recovery snapshot: %v", err)
	}

	// Simulate a failed optional-module rollout by changing the live copy after
	// the snapshot. The restored path must retain the known-good fields instead.
	st, err = testsupport.OpenStore(path, "admin", "different-seed", true)
	if err != nil {
		t.Fatal(err)
	}
	settingsRepository = settingsrepository.New(st)
	if _, err := settingsRepository.UpdateSiteSettings("Corrupted Admin", "/assets/bad.svg", updatedAt.Add(time.Minute)); err != nil {
		t.Fatalf("mutate failed rollout: %v", err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(restoredPath, snapshot, 0o644); err != nil {
		t.Fatalf("restore snapshot: %v", err)
	}

	// Start the recovered database with the MVP plan. Settings and Activity
	// surfaces are disabled, but their persisted data and operation log remain.
	st, err = testsupport.OpenStore(restoredPath, "admin", "different-seed", true)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	settingsRepository = settingsrepository.New(st)
	operationsRepository = operationlog.NewRepository(st)
	settings, err := settingsRepository.GetSiteSettings()
	if err != nil {
		t.Fatalf("read retained settings: %v", err)
	}
	if settings.SiteTitle != "Recovered Admin" || settings.LogoURL != "/assets/logo.svg" {
		t.Fatalf("retained settings = %+v", settings)
	}
	operations, total, err := operationsRepository.ListOperationsFiltered(operationlog.OperationFilter{
		Q: "settings.update", Sort: "createdAt", Order: "desc", Page: 1, PageSize: 10,
	})
	if err != nil {
		t.Fatalf("read retained operation log: %v", err)
	}
	if total != 1 || len(operations) != 1 || operations[0].ID != "r3-settings-update" {
		t.Fatalf("retained operations = %d / %+v, want one named row", total, operations)
	}
	recoveryRecordID := "recovery-check"
	recoveryDetail := `{"profile":"mvp","restored":true}`
	if err := operationsRepository.RecordOperation(operationlog.Operation{
		ID:        "r3-recovery-check",
		Event:     operationlog.EventSettingsUpdate,
		ActorID:   "user-admin",
		ActorName: "Admin",
		RecordID:  &recoveryRecordID,
		Detail:    &recoveryDetail,
		CreatedAt: updatedAt.Add(time.Minute),
	}); err != nil {
		t.Fatalf("write operation log after recovery: %v", err)
	}
	operations, total, err = operationsRepository.ListOperationsFiltered(operationlog.OperationFilter{
		Q: "settings.update", Sort: "createdAt", Order: "desc", Page: 1, PageSize: 10,
	})
	if err != nil {
		t.Fatalf("read operation log after recovery: %v", err)
	}
	if total != 2 || len(operations) != 2 || operations[0].ID != "r3-recovery-check" {
		t.Fatalf("recovered operations = %d / %+v, want two rows with recovery write first", total, operations)
	}

	resolution, err := kernel.ResolveProfile("mvp", nil)
	if err != nil {
		t.Fatal(err)
	}
	registry, err := kernel.NewRegistry(kernel.BuiltinModules())
	if err != nil {
		t.Fatal(err)
	}
	plan, err := registry.Resolve(resolution.Modules)
	if err != nil {
		t.Fatal(err)
	}
	a := auth.New([]byte("test-secret"), 0, 0, st, false)
	gate := &readinessGate{}
	gate.setReady()
	mux, err := testMux(a, st, plan, gate)
	if err != nil {
		t.Fatal(err)
	}
	ready := httptest.NewRecorder()
	mux.ServeHTTP(ready, httptest.NewRequest(http.MethodGet, "/readyz", nil))
	if ready.Code != http.StatusOK {
		t.Fatalf("readyz after MVP reopen = %d, want 200: %s", ready.Code, ready.Body.String())
	}
	manifest := httptest.NewRecorder()
	mux.ServeHTTP(manifest, httptest.NewRequest(http.MethodGet, "/.well-known/schema-ui/app-manifest.json", nil))
	if manifest.Code != http.StatusOK {
		t.Fatalf("manifest after recovery = %d: %s", manifest.Code, manifest.Body.String())
	}
	var document struct {
		Pages []struct {
			PageID string `json:"pageId"`
		} `json:"pages"`
	}
	if err := json.Unmarshal(manifest.Body.Bytes(), &document); err != nil {
		t.Fatal(err)
	}
	for _, page := range document.Pages {
		if page.PageID == "settings" || page.PageID == "activity" {
			t.Fatalf("disabled page %q leaked after recovery", page.PageID)
		}
	}
	for _, path := range []string{"/api/settings", "/api/operations", "/api/schema/settings", "/api/schema/activity"} {
		response := httptest.NewRecorder()
		mux.ServeHTTP(response, httptest.NewRequest(http.MethodGet, path, nil))
		if response.Code != http.StatusNotFound {
			t.Fatalf("disabled route %s status = %d, want 404", path, response.Code)
		}
	}
}

func lifecycleAppConfig(t *testing.T, profile, addr string) *config.Config {
	t.Helper()
	return &config.Config{
		AppName:      "test",
		AppEnv:       "development",
		HTTPAddr:     addr,
		DBPath:       filepath.Join(t.TempDir(), "composition.db"),
		ProfileName:  profile,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
		IdleTimeout:  time.Second,
	}
}

func TestAppStartsAndStopsDualProfiles(t *testing.T) {
	for _, profile := range []string{"mvp", "admin"} {
		t.Run(profile, func(t *testing.T) {
			cfg := lifecycleAppConfig(t, profile, "127.0.0.1:0")
			app, err := NewApp(cfg, "test-secret", "hash", slog.New(slog.NewTextHandler(io.Discard, nil)))
			if err != nil {
				t.Fatal(err)
			}
			startCtx, startCancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer startCancel()
			if err := app.Start(startCtx); err != nil {
				t.Fatalf("start %s profile: %v", profile, err)
			}
			stopCtx, stopCancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer stopCancel()
			if err := app.Stop(stopCtx); err != nil {
				t.Fatalf("stop %s profile: %v", profile, err)
			}
		})
	}
}

func TestAppStartsCustomProfileWithExplicitModules(t *testing.T) {
	admin, err := kernel.ResolveProfile("admin", nil)
	if err != nil {
		t.Fatal(err)
	}
	cfg := lifecycleAppConfig(t, "custom", "127.0.0.1:0")
	cfg.ModulesEnabled = append([]string(nil), admin.Modules...)
	app, err := NewApp(cfg, "test-secret", "hash", slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatal(err)
	}
	startCtx, startCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer startCancel()
	if err := app.Start(startCtx); err != nil {
		t.Fatalf("start custom profile: %v", err)
	}
	stopCtx, stopCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer stopCancel()
	if err := app.Stop(stopCtx); err != nil {
		t.Fatalf("stop custom profile: %v", err)
	}
}

func TestAppStartFailsClosedWithStableLifecycleErrorWhenPortIsUnavailable(t *testing.T) {
	for _, profile := range []string{"mvp", "admin"} {
		t.Run(profile, func(t *testing.T) {
			blocker, err := net.Listen("tcp", "127.0.0.1:0")
			if err != nil {
				t.Fatal(err)
			}
			defer blocker.Close()
			cfg := lifecycleAppConfig(t, profile, blocker.Addr().String())
			app, err := NewApp(cfg, "test-secret", "hash", slog.New(slog.NewTextHandler(io.Discard, nil)))
			if err != nil {
				t.Fatal(err)
			}
			startCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if err := app.Start(startCtx); err == nil {
				t.Fatal("occupied HTTP port must fail startup")
			} else {
				var kernelErr *kernel.Error
				if !errors.As(err, &kernelErr) || kernelErr.Code != kernel.CodeLifecycleStartFailed {
					t.Fatalf("startup error = %v, want stable lifecycle error", err)
				}
			}
			stopCtx, stopCancel := context.WithTimeout(context.Background(), time.Second)
			defer stopCancel()
			_ = app.Stop(stopCtx)
		})
	}
}

// stubProvider is a minimal kernel.Provider used to prove register-gate
// fail-closed under each profile (freeze §3 dual-profile matrix).
type stubProvider struct {
	desc kernel.Module
}

func (p *stubProvider) Descriptor() kernel.Module { return p.desc }
func (p *stubProvider) CompiledPersistence() ([]kernel.MigrationContribution, error) {
	return nil, nil
}
func (p *stubProvider) Register(context.Context, kernel.Registrar) error { return nil }

// TestDualProfileRegisterValidationFailClosed runs the register/conflict gate
// under both compiled profiles: a provider whose descriptor does not exactly
// match the plan module must fail closed (no partial surface), and neither
// profile's success is used as evidence for the other (freeze §3).
func TestDualProfileRegisterValidationFailClosed(t *testing.T) {
	for _, profile := range []string{"mvp", "admin"} {
		cfg := &config.Config{ProfileName: profile}
		plan, err := ResolvePlan(cfg)
		if err != nil {
			t.Fatalf("%s plan: %v", profile, err)
		}
		// Mismatched descriptor version → MODULE_API_MISMATCH (register gate).
		bad := &stubProvider{desc: kernel.Module{ID: "admin.users", Version: "9.9.9", KernelAPIRange: ">=2.0 <3.0"}}
		if _, err := kernel.RegisterContributions(context.Background(), plan, []kernel.Provider{bad}); err == nil {
			t.Fatalf("%s: mismatched provider descriptor must fail closed", profile)
		}
	}
}

// TestReadyzGatedOnModuleReadiness verifies /readyz returns unavailable until the
// module graph Start+Ready succeeds (R5 real readiness), not just store ping.
func TestReadyzGatedOnModuleReadiness(t *testing.T) {
	cfg := &config.Config{ProfileName: "admin"}
	plan, err := ResolvePlan(cfg)
	if err != nil {
		t.Fatal(err)
	}
	st, err := testsupport.OpenStore(":memory:", "admin", "hash", false)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	a := auth.New([]byte("test-secret"), 0, 0, st, false)

	// Gate unset → readyz unavailable even though the store pings.
	notReady := httptest.NewRecorder()
	mux, err := testMux(a, st, plan, &readinessGate{})
	if err != nil {
		t.Fatal(err)
	}
	mux.ServeHTTP(notReady, httptest.NewRequest(http.MethodGet, "/readyz", nil))
	if notReady.Code != http.StatusServiceUnavailable {
		t.Fatalf("readyz with unset gate = %d, want 503", notReady.Code)
	}

	// Gate set → readyz ok.
	readyGate := &readinessGate{}
	readyGate.setReady()
	mux, err = testMux(a, st, plan, readyGate)
	if err != nil {
		t.Fatal(err)
	}
	ready := httptest.NewRecorder()
	mux.ServeHTTP(ready, httptest.NewRequest(http.MethodGet, "/readyz", nil))
	if ready.Code != http.StatusOK {
		t.Fatalf("readyz with set gate = %d, want 200", ready.Code)
	}
}
