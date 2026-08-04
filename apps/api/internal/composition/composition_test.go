package composition

import (
	"context"
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
	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
)

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
	st, err := store.Open(":memory:", "admin", "hash", false)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	a := auth.New([]byte("test-secret"), 0, 0, st, false)
	mux, err := newMux(a, st, plan)
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
			brandingRoute:    http.StatusNotFound,
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
			st, err := store.Open(":memory:", "admin", "hash", false)
			if err != nil {
				t.Fatal(err)
			}
			defer st.Close()
			a := auth.New([]byte("test-secret"), 0, 0, st, false)
			mux, err := newMux(a, st, plan)
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

func TestMVPRecoveryRestoresOptionalModuleDataAndCoreReadiness(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "r3-recovery.db")
	snapshotPath := filepath.Join(dir, "pre-r3-snapshot.db")
	restoredPath := filepath.Join(dir, "r3-restored.db")
	st, err := store.Open(path, "admin", "hash", true)
	if err != nil {
		t.Fatal(err)
	}
	updatedAt := time.Date(2026, 8, 5, 0, 0, 0, 0, time.UTC)
	if _, err := st.UpdateSiteSettings("Recovered Admin", "/assets/logo.svg", updatedAt); err != nil {
		t.Fatalf("update settings: %v", err)
	}
	recordID := "default"
	detail := `{"siteTitle":"Recovered Admin"}`
	if err := st.RecordOperation(store.Operation{
		ID:        "r3-settings-update",
		Event:     store.EventSettingsUpdate,
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
	st, err = store.Open(path, "admin", "different-seed", true)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := st.UpdateSiteSettings("Corrupted Admin", "/assets/bad.svg", updatedAt.Add(time.Minute)); err != nil {
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
	st, err = store.Open(restoredPath, "admin", "different-seed", true)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	settings, err := st.GetSiteSettings()
	if err != nil {
		t.Fatalf("read retained settings: %v", err)
	}
	if settings.SiteTitle != "Recovered Admin" || settings.LogoURL != "/assets/logo.svg" {
		t.Fatalf("retained settings = %+v", settings)
	}
	operations, total, err := st.ListOperationsFiltered(store.OperationFilter{
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
	if err := st.RecordOperation(store.Operation{
		ID:        "r3-recovery-check",
		Event:     store.EventSettingsUpdate,
		ActorID:   "user-admin",
		ActorName: "Admin",
		RecordID:  &recoveryRecordID,
		Detail:    &recoveryDetail,
		CreatedAt: updatedAt.Add(time.Minute),
	}); err != nil {
		t.Fatalf("write operation log after recovery: %v", err)
	}
	operations, total, err = st.ListOperationsFiltered(store.OperationFilter{
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
	mux, err := newMux(a, st, plan)
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

func TestAppStartFailsClosedWithStableLifecycleErrorWhenPortIsUnavailable(t *testing.T) {
	blocker, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer blocker.Close()
	cfg := &config.Config{
		AppName:      "test",
		AppEnv:       "development",
		HTTPAddr:     blocker.Addr().String(),
		DBPath:       filepath.Join(t.TempDir(), "composition.db"),
		ProfileName:  "mvp",
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
		IdleTimeout:  time.Second,
	}
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
}
