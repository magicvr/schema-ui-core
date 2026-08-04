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
