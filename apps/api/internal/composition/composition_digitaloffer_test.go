// GOAL-005 A-006 F-003/F-006 closure: composition-root acceptance for
// biz.digital-offer — the module must be resolvable from the compiled registry,
// enabled through an explicit custom plan, and reachable end-to-end from the
// real NewApp mux (Manifest pages, schema documents, admin/public routes) with
// the Telegram-disabled dispatcher semantics frozen by D-002 §6.
package composition

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/config"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	"go.uber.org/fx"
)

// digitalOfferCoreModules is the minimal compiled-module set every
// digital-offer plan needs (core chain + a settings provider for the admin
// capability graph), mirroring the channel.telegram test plans.
var digitalOfferCoreModules = []string{
	"core.server-registration", "core.auth-session", "core.schema-render",
	"core.manifest-route", "core.navigation-capability", "core.operationlog",
	"admin.settings", "admin.account",
}

// TestDigitalOfferPlanResolution proves A-006 F-003: the compiled registry now
// resolves biz.digital-offer, both with and without the channel in the plan.
func TestDigitalOfferPlanResolution(t *testing.T) {
	cfgEnabled := &config.Config{
		ProfileName:    string(kernel.ProfileCustom),
		ModulesEnabled: append(append([]string(nil), digitalOfferCoreModules...), "channel.telegram", "biz.digital-offer"),
	}
	plan, err := ResolvePlan(cfgEnabled)
	if err != nil {
		t.Fatalf("ResolvePlan enabled: %v", err)
	}
	if !plan.HasModule("biz.digital-offer") {
		t.Fatal("enabled plan must contain biz.digital-offer")
	}
	if !plan.HasModule("channel.telegram") {
		t.Fatal("enabled plan must contain channel.telegram")
	}

	cfgDisabled := &config.Config{
		ProfileName:    string(kernel.ProfileCustom),
		ModulesEnabled: append(append([]string(nil), digitalOfferCoreModules...), "biz.digital-offer"),
	}
	planDisabled, err := ResolvePlan(cfgDisabled)
	if err != nil {
		t.Fatalf("ResolvePlan disabled: %v", err)
	}
	if !planDisabled.HasModule("biz.digital-offer") || planDisabled.HasModule("channel.telegram") {
		t.Fatalf("disabled plan = hasOffer:%v hasTelegram:%v", planDisabled.HasModule("biz.digital-offer"), planDisabled.HasModule("channel.telegram"))
	}
}

// TestDigitalOfferCompositionRoot proves the full chain from config to the
// served surface: the plan resolves, NewApp mounts the module, the manifest
// carries both pages, the schema documents answer under authentication, the
// admin surface is permission-gated and the public catalog is open. The plan
// deliberately excludes channel.telegram so the assembled dispatcher is the
// DisabledDispatcher no-op (D-002 §6) — the app starts with zero Bot API
// dependency.
func TestDigitalOfferCompositionRoot(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "digitaloffer_composition_test.db")
	seedHash, err := auth.HashPassword("admin-password", 4)
	if err != nil {
		t.Fatal(err)
	}
	cfg := &config.Config{
		ProfileName:    string(kernel.ProfileCustom),
		ModulesEnabled: append(append([]string(nil), digitalOfferCoreModules...), "biz.digital-offer"),
		AuthAccessTTL:  15 * time.Minute,
		AuthRefreshTTL: 30 * 24 * time.Hour,
		HTTPAddr:       "127.0.0.1:0",
		DBPath:         dbPath,
	}

	var mux *http.ServeMux
	app, err := newAppWithOptions(
		cfg,
		"test-secret",
		seedHash,
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		fx.Populate(&mux),
	)
	if err != nil {
		t.Fatalf("newAppWithOptions: %v", err)
	}
	startCtx, startCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer startCancel()
	if err := app.Start(startCtx); err != nil {
		t.Fatalf("app.Start: %v", err)
	}
	stopCtx, stopCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer stopCancel()
	defer func() { _ = app.Stop(stopCtx) }()

	// 1. Manifest carries both digital-offer pages with their route/schemaUrl
	// and sidebar navigation refs (structured assertion — A-014 F-001 aligns
	// this disabled-branch test with the Telegram-enabled one).
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/.well-known/schema-ui/app-manifest.json", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("manifest = %d", rr.Code)
	}
	var manifestDoc struct {
		Pages []struct {
			PageID    string `json:"pageId"`
			Route     string `json:"route"`
			SchemaURL string `json:"schemaUrl"`
		} `json:"pages"`
		Navigation struct {
			Sidebar []struct {
				PageRef string `json:"pageRef"`
			} `json:"sidebar"`
		} `json:"navigation"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &manifestDoc); err != nil {
		t.Fatalf("manifest is not valid JSON: %v", err)
	}
	manifestPages := map[string]struct{ route, schemaURL string }{}
	for _, p := range manifestDoc.Pages {
		manifestPages[p.PageID] = struct{ route, schemaURL string }{p.Route, p.SchemaURL}
	}
	for _, want := range []struct {
		id, route, schemaURL string
	}{
		{"digitaloffer-offers", "/digitaloffer-offers", "/api/schema/digitaloffer-offers"},
		{"digitaloffer-entitlements", "/digitaloffer-entitlements", "/api/schema/digitaloffer-entitlements"},
	} {
		got, ok := manifestPages[want.id]
		if !ok || got.route != want.route || got.schemaURL != want.schemaURL {
			t.Fatalf("manifest page %s = %+v, want route %s schemaUrl %s", want.id, got, want.route, want.schemaURL)
		}
	}
	manifestRefs := map[string]bool{}
	for _, n := range manifestDoc.Navigation.Sidebar {
		manifestRefs[n.PageRef] = true
	}
	for _, ref := range []string{"digitaloffer-offers", "digitaloffer-entitlements"} {
		if !manifestRefs[ref] {
			t.Fatalf("manifest sidebar missing pageRef %s (refs %v)", ref, manifestRefs)
		}
	}

	// 2. Schema documents: unauthenticated 401, admin 200 with the right
	// pageId; an unknown pageId stays fail-closed (404).
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/schema/digitaloffer-offers", nil))
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("unauthenticated schema = %d, want 401", rr.Code)
	}
	token := loginCompositionAdmin(t, mux)
	for _, pageID := range []string{"digitaloffer-offers", "digitaloffer-entitlements"} {
		req := httptest.NewRequest(http.MethodGet, "/api/schema/"+pageID, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rr = httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("schema %s = %d %s, want 200", pageID, rr.Code, rr.Body.String())
		}
		if !strings.Contains(rr.Body.String(), `"pageId"`) || !strings.Contains(rr.Body.String(), pageID) {
			t.Fatalf("schema %s body = %s", pageID, rr.Body.String())
		}
	}
	// Unknown pageId probes fail closed: 401 for anonymous callers (the
	// middleware guards before any page enumeration) and 404 for
	// authenticated ones (no seeded document).
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/schema/digitaloffer-unknown", nil))
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("unknown schema anonymous = %d, want 401", rr.Code)
	}
	unknownReq := httptest.NewRequest(http.MethodGet, "/api/schema/digitaloffer-unknown", nil)
	unknownReq.Header.Set("Authorization", "Bearer "+token)
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, unknownReq)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("unknown schema authenticated = %d, want 404", rr.Code)
	}

	// 3. Admin surface: permission-gated (401 unauthenticated, 200 admin).
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/digitaloffer/offers", nil))
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("unauthenticated admin offers = %d, want 401", rr.Code)
	}
	req := httptest.NewRequest(http.MethodGet, "/api/digitaloffer/offers", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("admin offers = %d %s, want 200", rr.Code, rr.Body.String())
	}

	// 4. Public catalog: open read of the (empty) on-sale list.
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/biz/offers", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("public catalog = %d %s, want 200", rr.Code, rr.Body.String())
	}
}

// loginCompositionAdmin logs in with the seeded initial password, rotates it
// (the seed carries must_change_password=1, which gates every other surface),
// and returns a fully privileged access token.
func loginCompositionAdmin(t *testing.T, mux *http.ServeMux) string {
	t.Helper()
	token := postLogin(t, mux, "admin-password")
	req := httptest.NewRequest(http.MethodPost, "/api/account/password",
		strings.NewReader(`{"currentPassword":"admin-password","newPassword":"admin-rotated-pw1"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("password rotation = %d %s", rr.Code, rr.Body.String())
	}
	return postLogin(t, mux, "admin-rotated-pw1")
}

func postLogin(t *testing.T, mux *http.ServeMux, password string) string {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{"username":"admin","password":"`+password+`"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("login = %d %s", rr.Code, rr.Body.String())
	}
	var out struct {
		AccessToken string `json:"accessToken"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &out); err != nil || out.AccessToken == "" {
		t.Fatalf("login body = %s (code %d)", rr.Body.String(), rr.Code)
	}
	return out.AccessToken
}
