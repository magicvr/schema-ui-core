package composition

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/config"
	"github.com/magicvr/schema-ui-core/apps/api/internal/ratelimit"
	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/modules/authsession"
	"github.com/magicvr/schema-ui-core/apps/api/modules/operationlog"
	settingsrepository "github.com/magicvr/schema-ui-core/apps/api/modules/settings/repository"
)

// S5 (GOAL-041 W29): acceptance-matrix HTTP manifest snapshots.
//
// D-002 §3.2 defines the S5 verification matrix: default profiles (mvp 6 /
// admin 22 / demo 14) x explicit-custom module combos (Digital Offer, Telegram).
// This test builds the REAL composition mux per combo (plan -> store -> newMux,
// the same assembly the server uses), serves GET /.well-known/schema-ui/
// app-manifest.json through httptest, asserts the frozen page counts, and —
// when S5_SNAPSHOT_DIR is set — writes the served bytes as canonical evidence
// snapshots. This proves the served manifest is a real profile projection, not
// the 35-page static universe, and that explicit-custom combos aggregate
// correctly.

type manifestSnapshotCase struct {
	name        string
	profile     string
	extra       []string // explicit-custom module additions (ProfileCustom)
	withTelegram bool
	wantPages   int
}

func s5SnapshotCases(t *testing.T) []manifestSnapshotCase {
	t.Helper()
	adminRes, err := kernel.ResolveProfile("admin", nil)
	if err != nil {
		t.Fatal(err)
	}
	adminModules := append([]string(nil), adminRes.Modules...)
	// Fresh backing arrays per case: append() on a shared slice aliases the
	// backing array, so the second append would overwrite the first's tail.
	withExtra := func(items ...string) []string {
		out := make([]string, 0, len(adminModules)+len(items))
		out = append(out, adminModules...)
		return append(out, items...)
	}
	return []manifestSnapshotCase{
		{name: "mvp", profile: "mvp", wantPages: 6},
		{name: "admin", profile: "admin", wantPages: 22},
		{name: "demo", profile: "demo", wantPages: 14},
		{name: "admin+digitaloffer", profile: "custom", extra: withExtra("biz.digital-offer"), wantPages: 25},
		{name: "admin+telegram", profile: "custom", extra: withExtra("channel.telegram"), withTelegram: true, wantPages: 24},
	}
}

// s5Mux mirrors testMux but additionally wires the real TelegramRuntime for
// telegram-enabled combos (testMux hardcodes tr == nil).
func s5Mux(t *testing.T, plan kernel.Plan, withTelegram bool) *http.ServeMux {
	t.Helper()
	st, err := testsupport.OpenStore(":memory:", "admin", "hash", false)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	a := auth.New([]byte("test-secret"), 0, 0, st, true)
	jobRuntime, err := newJobRuntime(st)
	if err != nil {
		t.Fatal(err)
	}
	cachePort, err := newCache(&config.Config{DBPath: "test.db"})
	if err != nil {
		t.Fatal(err)
	}
	eventBusPort := newEventBus(&config.Config{DBPath: "test.db"}, slog.New(slog.NewTextHandler(io.Discard, nil)))
	var tr *TelegramRuntime
	if withTelegram {
		tr, err = newTelegramRuntime(plan, &config.Config{TelegramBotToken: "s5-test-token", TelegramWebhookSecret: "s5-test-secret"}, st, ratelimit.NewProvider())
		if err != nil {
			t.Fatal(err)
		}
	}
	mux, err := newMux(
		&config.Config{DBPath: "test.db"},
		a,
		st,
		authsession.NewRepository(st),
		operationlog.NewRepository(st),
		settingsrepository.New(st),
		plan,
		&readinessGate{},
		jwtSecret("test-secret"),
		jobRuntime,
		nil,
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		cachePort,
		eventBusPort,
		ratelimit.NewProvider(),
		tr,
	)
	if err != nil {
		t.Fatal(err)
	}
	return mux
}

func TestS5ManifestSnapshots(t *testing.T) {
	snapshotDir := os.Getenv("S5_SNAPSHOT_DIR")
	for _, tc := range s5SnapshotCases(t) {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &config.Config{ProfileName: tc.profile, ModulesEnabled: tc.extra}
			plan, err := ResolvePlan(cfg)
			if err != nil {
				t.Fatalf("ResolvePlan(%s): %v", tc.name, err)
			}
			t.Logf("%s plan modules: %v", tc.name, plan.IDs())
			for _, want := range tc.extra {
				if !plan.HasModule(want) {
					t.Fatalf("%s: plan missing explicit module %q", tc.name, want)
				}
			}
			mux := s5Mux(t, plan, tc.withTelegram)

			response := httptest.NewRecorder()
			mux.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/.well-known/schema-ui/app-manifest.json", nil))
			if response.Code != http.StatusOK {
				t.Fatalf("GET manifest = %d, want 200: %s", response.Code, response.Body.String())
			}

			var document struct {
				ProtocolVersion     string `json:"protocolVersion"`
				RequiredCapabilities []string `json:"requiredCapabilities"`
				Pages               []struct {
					PageID string `json:"pageId"`
				} `json:"pages"`
			}
			if err := json.Unmarshal(response.Body.Bytes(), &document); err != nil {
				t.Fatal(err)
			}
			if len(document.Pages) != tc.wantPages {
				t.Fatalf("%s pages = %d, want %d (%v)", tc.name, len(document.Pages), tc.wantPages, pageIDs(document.Pages))
			}
			if document.ProtocolVersion != "2.7" {
				t.Fatalf("%s protocolVersion = %q, want 2.7", tc.name, document.ProtocolVersion)
			}
			hasManifestCap := false
			for _, cap := range document.RequiredCapabilities {
				if cap == "app.manifest" {
					hasManifestCap = true
				}
			}
			if !hasManifestCap {
				t.Fatalf("%s manifest missing requiredCapabilities app.manifest", tc.name)
			}

			if snapshotDir != "" {
				dir := filepath.Join(snapshotDir, tc.name)
				if err := os.MkdirAll(dir, 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(dir, "app-manifest.json"), response.Body.Bytes(), 0o644); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

func pageIDs(pages []struct{ PageID string `json:"pageId"` }) []string {
	out := make([]string, 0, len(pages))
	for _, p := range pages {
		out = append(out, p.PageID)
	}
	return out
}
