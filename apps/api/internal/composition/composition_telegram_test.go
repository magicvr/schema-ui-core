package composition

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"go.uber.org/fx"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/config"
	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/modules/authsession"
	"github.com/magicvr/schema-ui-core/apps/api/modules/operationlog"
	settingsrepository "github.com/magicvr/schema-ui-core/apps/api/modules/settings/repository"
)

func TestTelegramChannelComposition(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "telegram_test.db")
	st, err := testsupport.OpenStore(dbPath, "admin", "hash", false)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })

	// 1. Verify default profile (mvp) does not have channel.telegram
	planMVP, err := kernel.ResolveProfile(string(kernel.ProfileMVP), nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range planMVP.Modules {
		if m == "channel.telegram" {
			t.Fatalf("channel.telegram must not be in default mvp profile")
		}
	}

	// 2. Custom profile with channel.telegram enabled
	cfg := &config.Config{
		ProfileName:           string(kernel.ProfileCustom),
		ModulesEnabled:        []string{"core.server-registration", "core.auth-session", "core.schema-render", "core.manifest-route", "core.navigation-capability", "core.operationlog", "admin.settings", "channel.telegram"},
		TelegramBotToken:      "test-token",
		TelegramWebhookSecret: "test-secret",
	}

	plan, err := ResolvePlan(cfg)
	if err != nil {
		t.Fatalf("ResolvePlan error: %v", err)
	}
	if !plan.HasModule("channel.telegram") {
		t.Fatalf("expected plan to have channel.telegram")
	}

	// 3. ResolvePlan with BuiltinModules resolution
	reg, err := kernel.NewRegistry(kernel.BuiltinModules())
	if err != nil {
		t.Fatalf("kernel.NewRegistry failed: %v", err)
	}

	resolvedPlan, err := reg.Resolve(cfg.ModulesEnabled)
	if err != nil {
		t.Fatalf("registry resolve with channel.telegram failed: %v", err)
	}
	if !resolvedPlan.HasModule("channel.telegram") {
		t.Fatalf("expected resolved plan to contain channel.telegram")
	}

	// 4. Test provider registration through RegisterContributions
	desc := kernel.Module{
		ID:             "channel.telegram",
		Version:        "2.0.0",
		KernelAPIRange: ">=2.0 <3.0",
		DependsOn:      []string{"core.server-registration", "core.schema-render", "core.navigation-capability", "admin.settings"},
		Requires:       []kernel.Capability{kernel.CapabilityHTTP, kernel.CapabilitySchema, kernel.CapabilityNavigation},
		Contributions: kernel.ContributionKeys{
			Routes: []string{
				"GET /api/channel/telegram/settings",
				"PATCH /api/channel/telegram/settings",
				"POST /api/channel/telegram/lease/acquire",
				"POST /api/channel/telegram/lease/heartbeat",
				"POST /api/channel/telegram/lease/release",
				"POST /api/channel/telegram/webhook",
			},
			Pages:      []string{"telegram-settings"},
			Navigation: []string{"menu_telegram"},
			Fragments:  []string{"telegram-settings"},
		},
	}

	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	provider := &dummyTelegramProvider{desc: desc, handler: dummyHandler}

	set, err := kernel.RegisterContributions(context.Background(), resolvedPlan, []kernel.Provider{provider})
	if err != nil {
		t.Fatalf("RegisterContributions failed: %v", err)
	}

	if len(set.Routes) != 6 {
		t.Fatalf("unexpected route contributions: %+v", set.Routes)
	}

	// 5. Test webhook route execution
	req := httptest.NewRequest(http.MethodPost, "/api/channel/telegram/webhook", nil)
	w := httptest.NewRecorder()
	set.Routes[5].Handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}
}

type dummyTelegramProvider struct {
	desc    kernel.Module
	handler http.Handler
}

func (d *dummyTelegramProvider) Descriptor() kernel.Module {
	return d.desc
}

func (d *dummyTelegramProvider) CompiledPersistence() ([]kernel.MigrationContribution, error) {
	return nil, nil
}

func (d *dummyTelegramProvider) Register(ctx context.Context, reg kernel.Registrar) error {
	for _, r := range d.desc.Contributions.Routes {
		var method, pattern string
		if r == "GET /api/channel/telegram/settings" {
			method, pattern = "GET", "/api/channel/telegram/settings"
		} else if r == "PATCH /api/channel/telegram/settings" {
			method, pattern = "PATCH", "/api/channel/telegram/settings"
		} else if r == "POST /api/channel/telegram/lease/acquire" {
			method, pattern = "POST", "/api/channel/telegram/lease/acquire"
		} else if r == "POST /api/channel/telegram/lease/heartbeat" {
			method, pattern = "POST", "/api/channel/telegram/lease/heartbeat"
		} else if r == "POST /api/channel/telegram/lease/release" {
			method, pattern = "POST", "/api/channel/telegram/lease/release"
		} else {
			method, pattern = "POST", "/api/channel/telegram/webhook"
		}
		if err := reg.HTTP(kernel.RouteContribution{
			ContributionIdentity: kernel.ContributionIdentity{
				ModuleID: "channel.telegram",
				Key:      r,
			},
			Method:  method,
			Pattern: pattern,
			Handler: d.handler,
			Public:  true,
		}); err != nil {
			return err
		}
	}
	return nil
}

func TestResolveTelegramPorts_EnabledAndDisabled(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "telegram_ports_test.db")
	st, err := testsupport.OpenStore(dbPath, "admin", "hash", false)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })

	cfg := &config.Config{
		TelegramBotToken:      "live-bot-token",
		TelegramWebhookSecret: "live-secret",
	}

	// 1. When module is disabled in plan: returns disabled no-op dispatcher and fail-closed sender
	planDisabled := kernel.Plan{Modules: []kernel.Module{}}
	dispDisabled, senderDisabled, err := ResolveTelegramPorts(planDisabled, cfg, st)
	if err != nil {
		t.Fatalf("ResolveTelegramPorts(disabled): %v", err)
	}

	msg := kernel.TelegramMessage{ChatID: "123", Text: "Hello"}
	if err := senderDisabled.Send(context.Background(), msg); err != kernel.ErrTelegramDisabled {
		t.Fatalf("expected ErrTelegramDisabled when channel disabled, got %v", err)
	}
	if err := dispDisabled.RegisterCommand("start", func(ctx context.Context, upd kernel.TelegramUpdate) error { return nil }); err != nil {
		t.Fatalf("expected no-op nil from disabled dispatcher, got %v", err)
	}

	// 2. When module is enabled in plan: returns live dispatcher and sender
	planEnabled := kernel.Plan{Modules: []kernel.Module{{ID: "channel.telegram"}}}
	dispEnabled, senderEnabled, err := ResolveTelegramPorts(planEnabled, cfg, st)
	if err != nil {
		t.Fatalf("ResolveTelegramPorts(enabled): %v", err)
	}

	if senderEnabled == nil || dispEnabled == nil {
		t.Fatalf("expected non-nil live sender and dispatcher")
	}
	// Live dispatcher registers commands properly
	called := false
	err = dispEnabled.RegisterCommand("ping", func(ctx context.Context, upd kernel.TelegramUpdate) error {
		called = true
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error registering command on live dispatcher: %v", err)
	}
	_ = called
}

func TestTelegramRuntime_PersistenceAcrossRestart(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "telegram_persist_test.db")
	st1, err := testsupport.OpenStore(dbPath, "admin", "hash", false)
	if err != nil {
		t.Fatal(err)
	}

	cfg1 := &config.Config{
		TelegramBotToken:      "seed-token",
		TelegramWebhookSecret: "seed-secret",
		TelegramMasterKey:     "test-master-key",
	}
	plan := kernel.Plan{Modules: []kernel.Module{{ID: "channel.telegram"}}}

	// Start first process instance and hot-switch via RuntimeManager.Update (F-002)
	tr1, err := newTelegramRuntime(plan, cfg1, st1, newRateLimiters())
	if err != nil {
		t.Fatalf("newTelegramRuntime: %v", err)
	}
	if tr1.Manager == nil {
		t.Fatalf("expected non-nil manager when channel.telegram enabled")
	}

	ctx := context.Background()
	if err := tr1.Manager.Update(ctx, "persisted-live-token", "persisted-live-secret"); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	_ = st1.Close()

	// Reopen database to simulate restart (F-002)
	st2, err := testsupport.OpenStore(dbPath, "admin", "hash", false)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st2.Close() })

	cfg2 := &config.Config{
		TelegramBotToken:      "seed-token", // seed is old, but DB has updated value!
		TelegramWebhookSecret: "seed-secret",
		TelegramMasterKey:     "test-master-key",
	}

	// New process runtime initialized with old seed should reload persisted token from DB
	tr2, err := newTelegramRuntime(plan, cfg2, st2, newRateLimiters())
	if err != nil {
		t.Fatalf("newTelegramRuntime(restart): %v", err)
	}
	if tr2.Manager == nil {
		t.Fatalf("expected non-nil manager after restart")
	}

	if tr2.Manager.GetToken() != "persisted-live-token" {
		t.Fatalf("expected persisted token to survive restart, got %q", tr2.Manager.GetToken())
	}
	if tr2.Manager.GetSecret() != "persisted-live-secret" {
		t.Fatalf("expected persisted secret to survive restart, got %q", tr2.Manager.GetSecret())
	}
}

// TestTelegramRuntime_ClearSurvivesRestart (A-008 informational): an admin
// clearing the token/secret must survive restart instead of reverting to the
// env seed — the persisted row is authoritative even when its values are empty.
func TestTelegramRuntime_ClearSurvivesRestart(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "telegram_clear_test.db")
	st1, err := testsupport.OpenStore(dbPath, "admin", "hash", false)
	if err != nil {
		t.Fatal(err)
	}

	cfg1 := &config.Config{
		TelegramBotToken:      "seed-token",
		TelegramWebhookSecret: "seed-secret",
		TelegramMasterKey:     "test-master-key",
	}
	plan := kernel.Plan{Modules: []kernel.Module{{ID: "channel.telegram"}}}

	tr1, err := newTelegramRuntime(plan, cfg1, st1, newRateLimiters())
	if err != nil {
		t.Fatalf("newTelegramRuntime: %v", err)
	}
	if tr1.Manager == nil {
		t.Fatalf("expected non-nil manager when channel.telegram enabled")
	}

	// Admin clears both secrets.
	ctx := context.Background()
	if err := tr1.Manager.Update(ctx, "", ""); err != nil {
		t.Fatalf("Update(clear) failed: %v", err)
	}
	if tr1.Manager.GetToken() != "" || tr1.Manager.GetSecret() != "" {
		t.Fatalf("expected cleared in-memory state, got token=%q secret=%q", tr1.Manager.GetToken(), tr1.Manager.GetSecret())
	}

	_ = st1.Close()

	// Reopen with a NON-EMPTY env seed: the persisted (cleared) row must win.
	st2, err := testsupport.OpenStore(dbPath, "admin", "hash", false)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st2.Close() })

	cfg2 := &config.Config{
		TelegramBotToken:      "seed-token", // must NOT resurrect after clear
		TelegramWebhookSecret: "seed-secret",
		TelegramMasterKey:     "test-master-key",
	}

	tr2, err := newTelegramRuntime(plan, cfg2, st2, newRateLimiters())
	if err != nil {
		t.Fatalf("newTelegramRuntime(restart): %v", err)
	}
	if tr2.Manager.GetToken() != "" {
		t.Fatalf("cleared token resurrected to env seed %q", tr2.Manager.GetToken())
	}
	if tr2.Manager.GetSecret() != "" {
		t.Fatalf("cleared secret resurrected to env seed %q", tr2.Manager.GetSecret())
	}
}

func TestTelegramRuntime_ConnectionSettingsPersistenceAndAuthority(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "telegram_connection_settings_test.db")
	st1, err := testsupport.OpenStore(dbPath, "admin", "hash", false)
	if err != nil {
		t.Fatal(err)
	}

	plan := kernel.Plan{Modules: []kernel.Module{{ID: "channel.telegram"}}}
	cfg1 := &config.Config{
		TelegramBotToken:             "seed-token",
		TelegramWebhookSecret:        "seed-secret",
		TelegramMode:                 "webhook",
		TelegramWebhookPublicBaseURL: "https://seed.example",
		TelegramMasterKey:            "test-master-key",
	}
	tr1, err := newTelegramRuntime(plan, cfg1, st1, newRateLimiters())
	if err != nil {
		t.Fatalf("newTelegramRuntime: %v", err)
	}
	if got := tr1.Manager.Status(); got.Mode != "webhook" || got.WebhookPublicBaseURL != "https://seed.example" {
		t.Fatalf("seed connection settings = %+v", got)
	}
	if err := tr1.Manager.UpdateSettings(context.Background(), "live-token", "live-secret", "webhook", "https://live.example"); err != nil {
		t.Fatalf("UpdateSettings: %v", err)
	}
	if err := st1.Close(); err != nil {
		t.Fatal(err)
	}

	st2, err := testsupport.OpenStore(dbPath, "admin", "hash", false)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st2.Close() })
	cfg2 := &config.Config{
		TelegramBotToken:             "stale-seed-token",
		TelegramWebhookSecret:        "stale-seed-secret",
		TelegramMode:                 "polling",
		TelegramWebhookPublicBaseURL: "https://stale.example",
		TelegramMasterKey:            "test-master-key",
	}
	tr2, err := newTelegramRuntime(plan, cfg2, st2, newRateLimiters())
	if err != nil {
		t.Fatalf("newTelegramRuntime(restart): %v", err)
	}
	got := tr2.Manager.Status()
	if got.Mode != "webhook" || got.WebhookPublicBaseURL != "https://live.example" ||
		tr2.Manager.GetToken() != "live-token" || tr2.Manager.GetSecret() != "live-secret" {
		t.Fatalf("persisted Telegram settings were not authoritative: status=%+v token=%q secret=%q", got, tr2.Manager.GetToken(), tr2.Manager.GetSecret())
	}
}

func TestTelegramRuntime_EmptyConnectionSettingsRemainAuthoritative(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "telegram_empty_connection_settings_test.db")
	st1, err := testsupport.OpenStore(dbPath, "admin", "hash", false)
	if err != nil {
		t.Fatal(err)
	}

	plan := kernel.Plan{Modules: []kernel.Module{{ID: "channel.telegram"}}}
	cfg1 := &config.Config{
		TelegramBotToken:             "seed-token",
		TelegramWebhookSecret:        "seed-secret",
		TelegramMode:                 "webhook",
		TelegramWebhookPublicBaseURL: "https://seed.example",
		TelegramMasterKey:            "test-master-key",
	}
	_, err = newTelegramRuntime(plan, cfg1, st1, newRateLimiters())
	if err != nil {
		t.Fatalf("newTelegramRuntime: %v", err)
	}
	if err := st1.Run(context.Background(), func(tx kernel.Tx) error {
		_, err := tx.Exec(context.Background(), `UPDATE telegram_config SET mode = '', webhook_public_base_url = '' WHERE id = 1`)
		return err
	}); err != nil {
		t.Fatalf("clear persisted connection settings: %v", err)
	}
	if err := st1.Close(); err != nil {
		t.Fatal(err)
	}

	st2, err := testsupport.OpenStore(dbPath, "admin", "hash", false)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st2.Close() })
	cfg2 := &config.Config{
		TelegramBotToken:             "stale-seed-token",
		TelegramWebhookSecret:        "stale-seed-secret",
		TelegramMode:                 "webhook",
		TelegramWebhookPublicBaseURL: "https://stale.example",
		TelegramMasterKey:            "test-master-key",
	}
	tr2, err := newTelegramRuntime(plan, cfg2, st2, newRateLimiters())
	if err != nil {
		t.Fatalf("newTelegramRuntime(restart): %v", err)
	}
	got := tr2.Manager.Status()
	if got.Mode != "polling" || got.WebhookPublicBaseURL != "" {
		t.Fatalf("empty persisted connection settings were not authoritative: %+v", got)
	}
}

func TestTelegramChannelComposition_RealWebhookMount(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "telegram_real_mount_test.db")
	st, err := testsupport.OpenStore(dbPath, "admin", "hash", false)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })
	telegramAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/botlive-bot-token/getMe":
			_, _ = w.Write([]byte(`{"ok":true,"result":{"id":101,"is_bot":true,"username":"mount_test_bot"}}`))
		case "/botlive-bot-token/deleteWebhook":
			_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(telegramAPI.Close)

	cfg := &config.Config{
		ProfileName:           string(kernel.ProfileCustom),
		ModulesEnabled:        []string{"core.server-registration", "core.auth-session", "core.schema-render", "core.manifest-route", "core.navigation-capability", "core.operationlog", "admin.settings", "channel.telegram"},
		TelegramBotToken:      "live-bot-token",
		TelegramWebhookSecret: "correct-secret",
		TelegramMasterKey:     "test-master-key",
	}

	// Resolve through the registry so channel.telegram's DependsOn pulls in
	// admin.settings (whose settings.read permission menu_telegram references)
	// exactly as production does (R-001 / A-002).
	reg, err := kernel.NewRegistry(kernel.BuiltinModules())
	if err != nil {
		t.Fatal(err)
	}
	plan, err := reg.Resolve([]string{
		"core.server-registration", "core.auth-session", "core.schema-render",
		"core.manifest-route", "core.navigation-capability", "core.operationlog", "admin.settings", "channel.telegram",
	})
	if err != nil {
		t.Fatalf("registry.Resolve: %v", err)
	}

	// Construct single shared TelegramRuntime (F-001)
	rateLimiters := newRateLimiters()
	tr, err := buildTelegramRuntime(plan, cfg, st, rateLimiters, telegramRuntimeOptions{APIBaseURL: telegramAPI.URL})
	if err != nil {
		t.Fatalf("newTelegramRuntime: %v", err)
	}
	if err := tr.Connection.Start(context.Background()); err != nil {
		t.Fatalf("telegram connection start: %v", err)
	}
	t.Cleanup(func() { _ = tr.Connection.Stop(context.Background()) })

	// Register command on the SAME dispatcher that webhook dispatches to
	commandCalled := false
	err = tr.Dispatcher.RegisterCommand("status", func(ctx context.Context, upd kernel.TelegramUpdate) error {
		commandCalled = true
		return nil
	})
	if err != nil {
		t.Fatalf("RegisterCommand failed: %v", err)
	}

	// Mount real mux with real tr (R-002)
	a := auth.New([]byte("test-secret-at-least-32-chars-long!!"), 0, 0, st, false)
	authRepo := authsession.NewRepository(st)
	opRepo := operationlog.NewRepository(st)
	setRepo := settingsrepository.New(st)
	jobs, err := newJobRuntime(st)
	if err != nil {
		t.Fatal(err)
	}

	mux, err := newMuxWithExtraProviders(
		cfg,
		a,
		st,
		authRepo,
		opRepo,
		setRepo,
		plan,
		&readinessGate{},
		jwtSecret("test-secret-at-least-32-chars-long!!"),
		jobs,
		nil,
		nil,
		slog.Default(),
		nil,
		nil,
		rateLimiters,
		tr,
	)
	if err != nil {
		t.Fatalf("newMuxWithExtraProviders failed: %v", err)
	}

	// 1. Webhook with wrong secret -> 401
	reqWrongSec := httptest.NewRequest(http.MethodPost, "/api/channel/telegram/webhook", strings.NewReader(`{}`))
	reqWrongSec.Header.Set("X-Telegram-Bot-Api-Secret-Token", "wrong-secret")
	wWrongSec := httptest.NewRecorder()
	mux.ServeHTTP(wWrongSec, reqWrongSec)
	if wWrongSec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 on wrong secret, got %d", wWrongSec.Code)
	}

	// 2. Webhook with correct secret calling /status command -> 200 and dispatches to registered handler
	reqValid := httptest.NewRequest(http.MethodPost, "/api/channel/telegram/webhook", strings.NewReader(`{
		"update_id": 1,
		"message": {
			"chat": {"id": 12345},
			"from": {"id": 67890},
			"text": "/status"
		}
	}`))
	reqValid.Header.Set("X-Telegram-Bot-Api-Secret-Token", "correct-secret")
	wValid := httptest.NewRecorder()
	mux.ServeHTTP(wValid, reqValid)
	if wValid.Code != http.StatusOK {
		t.Fatalf("expected 200 on valid webhook, got %d", wValid.Code)
	}
	if !commandCalled {
		t.Fatalf("expected command handler on shared dispatcher to be called via mounted webhook")
	}
	var persistedInbound int
	if err := st.Run(context.Background(), func(tx kernel.Tx) error {
		return tx.QueryRow(context.Background(), `SELECT COUNT(*) FROM telegram_inbound_messages WHERE bot_id = ? AND update_id = ?`, 101, 1).Scan(&persistedInbound)
	}); err != nil {
		t.Fatalf("query mounted webhook receipt: %v", err)
	}
	if persistedInbound != 1 {
		t.Fatalf("mounted webhook persisted receipts = %d, want 1", persistedInbound)
	}

	// 3. Settings endpoint without auth -> 401 Unauthorized
	reqSettings := httptest.NewRequest(http.MethodGet, "/api/channel/telegram/settings", nil)
	wSettings := httptest.NewRecorder()
	mux.ServeHTTP(wSettings, reqSettings)
	if wSettings.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 on unauthenticated settings request, got %d", wSettings.Code)
	}

	// The real provider also mounts the authenticated console lease. Use the
	// explicit dev-session test mode to exercise the complete mux middleware
	// and prove the route reaches the same process-level manager as the webhook.
	aDev := auth.New([]byte("test-secret-at-least-32-chars-long!!"), 0, 0, st, true)
	muxLease, err := newMuxWithExtraProviders(
		cfg,
		aDev,
		st,
		authRepo,
		opRepo,
		setRepo,
		plan,
		&readinessGate{},
		jwtSecret("test-secret-at-least-32-chars-long!!"),
		jobs,
		nil,
		nil,
		slog.Default(),
		nil,
		nil,
		rateLimiters,
		tr,
	)
	if err != nil {
		t.Fatalf("newMuxWithExtraProviders(lease): %v", err)
	}
	reqLease := httptest.NewRequest(http.MethodPost, "/api/channel/telegram/lease/acquire", nil)
	wLease := httptest.NewRecorder()
	muxLease.ServeHTTP(wLease, reqLease)
	if wLease.Code != http.StatusOK || tr.Connection.ActiveLeaseCount() != 1 {
		t.Fatalf("lease acquire through composed mux status=%d active=%d body=%s", wLease.Code, tr.Connection.ActiveLeaseCount(), wLease.Body.String())
	}
	reqRelease := httptest.NewRequest(http.MethodPost, "/api/channel/telegram/lease/release", nil)
	wRelease := httptest.NewRecorder()
	muxLease.ServeHTTP(wRelease, reqRelease)
	if wRelease.Code != http.StatusOK || tr.Connection.ActiveLeaseCount() != 0 {
		t.Fatalf("lease release through composed mux status=%d active=%d body=%s", wRelease.Code, tr.Connection.ActiveLeaseCount(), wRelease.Body.String())
	}
}

// TestTelegramFxInjection_SameRuntime proves F-001 / A-006 closure THROUGH the
// real Fx graph: the *TelegramRuntime that NewApp's fx graph injects into
// newMux (and exposes as kernel.TelegramDispatcher) is the SAME instance whose
// dispatcher the mounted webhook uses. It does NOT hand-pass tr to newMux —
// the previous test's flaw the audit rejected.
func TestTelegramFxInjection_SameRuntime(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "telegram_fx_test.db")
	apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/botlive-bot-token/getMe":
			_, _ = w.Write([]byte(`{"ok":true,"result":{"id":101,"is_bot":true,"username":"fx_test_bot"}}`))
		case "/botlive-bot-token/deleteWebhook":
			_, _ = w.Write([]byte(`{"ok":true,"result":true}`))
		default:
			t.Errorf("unexpected Telegram API request: %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer apiServer.Close()
	cfg := &config.Config{
		ProfileName:           string(kernel.ProfileCustom),
		ModulesEnabled:        []string{"core.server-registration", "core.auth-session", "core.schema-render", "core.manifest-route", "core.navigation-capability", "core.operationlog", "admin.settings", "channel.telegram"},
		TelegramBotToken:      "live-bot-token",
		TelegramWebhookSecret: "correct-secret",
		TelegramMasterKey:     "test-master-key",
		HTTPAddr:              "127.0.0.1:0",
		DBPath:                dbPath,
	}

	// Populate the injected graph values: the single *TelegramRuntime and the
	// *http.ServeMux that NewApp builds. No tr is passed by hand anywhere.
	var injected *TelegramRuntime
	var mux *http.ServeMux
	app, err := newAppWithOptions(
		cfg,
		"test-secret",
		"hash",
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		fx.Populate(&injected),
		fx.Populate(&mux),
		fx.Supply(&telegramRuntimeOptions{APIBaseURL: apiServer.URL}),
	)
	if err != nil {
		t.Fatalf("newAppWithOptions: %v", err)
	}

	startCtx, startCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer startCancel()
	if err := app.Start(startCtx); err != nil {
		t.Fatalf("app.Start: %v", err)
	}
	stopCtx, stopCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer stopCancel()
	defer func() { _ = app.Stop(stopCtx) }()

	if injected == nil || injected.Webhook == nil {
		t.Fatalf("expected live *TelegramRuntime injected by Fx graph")
	}
	if mux == nil {
		t.Fatalf("expected *http.ServeMux injected by Fx graph")
	}

	// Register a command on the INJECTED dispatcher, then hit the webhook route
	// served by the SAME mux — the command must be dispatched. If Fx had built a
	// second runtime for the webhook, this registration would be invisible.
	commandCalled := false
	if err := injected.Dispatcher.RegisterCommand("status", func(ctx context.Context, upd kernel.TelegramUpdate) error {
		commandCalled = true
		return nil
	}); err != nil {
		t.Fatalf("RegisterCommand on injected dispatcher: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/channel/telegram/webhook", strings.NewReader(`{
		"update_id": 1,
		"message": {
			"chat": {"id": 12345},
			"from": {"id": 67890},
			"text": "/status"
		}
	}`))
	req.Header.Set("X-Telegram-Bot-Api-Secret-Token", "correct-secret")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on valid webhook through Fx-built mux, got %d", w.Code)
	}
	if !commandCalled {
		t.Fatalf("Fx-injected dispatcher did NOT serve the webhook command — runtime duplication still present")
	}
}

// TestTelegramSettingsSchema_MountAndDisable (R-003 / A-002): through the real
// composition-root mux, the telegram-settings schema document answers 200 when
// channel.telegram is enabled and 404 when it is not (the module is not in the
// default profiles). This probes the RegisterSchemas(set.Pages) wiring path the
// same way s2_access_drill does for probe modules.
func TestTelegramSettingsSchema_MountAndDisable(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "telegram_schema_mount_test.db")

	// Enabled: schema document served through the composition root.
	st, err := testsupport.OpenStore(dbPath, "admin", "hash", false)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })

	cfgEnabled := &config.Config{
		ProfileName:           string(kernel.ProfileCustom),
		ModulesEnabled:        []string{"core.server-registration", "core.auth-session", "core.schema-render", "core.manifest-route", "core.navigation-capability", "core.operationlog", "admin.settings", "channel.telegram"},
		TelegramBotToken:      "live-bot-token",
		TelegramWebhookSecret: "correct-secret",
		TelegramMasterKey:     "test-master-key",
	}
	planEnabled, err := ResolvePlan(cfgEnabled)
	if err != nil {
		t.Fatalf("ResolvePlan(enabled): %v", err)
	}
	// W13 F-010: schemas are authenticated — dev-session authenticator so the
	// probe reaches the handler (s2_access_drill precedent).
	a := auth.New([]byte("test-secret-at-least-32-chars-long!!"), 0, 0, st, true)
	authRepo := authsession.NewRepository(st)
	opRepo := operationlog.NewRepository(st)
	setRepo := settingsrepository.New(st)
	jobs, err := newJobRuntime(st)
	if err != nil {
		t.Fatal(err)
	}
	cachePort, err := newCache(&config.Config{ProfileName: string(kernel.ProfileCustom)})
	if err != nil {
		t.Fatal(err)
	}
	eventBusPort := newEventBus(&config.Config{ProfileName: string(kernel.ProfileCustom)}, slog.New(slog.NewTextHandler(io.Discard, nil)))
	rateLimiters := newRateLimiters()
	tr, err := newTelegramRuntime(planEnabled, cfgEnabled, st, rateLimiters)
	if err != nil {
		t.Fatalf("newTelegramRuntime: %v", err)
	}
	muxEnabled, err := newMux(
		cfgEnabled, a, st, authRepo, opRepo, setRepo, planEnabled, &readinessGate{},
		jwtSecret("test-secret-at-least-32-chars-long!!"), jobs, nil, slog.New(slog.NewTextHandler(io.Discard, nil)),
		cachePort, eventBusPort, rateLimiters, tr,
	)
	if err != nil {
		t.Fatalf("newMux(enabled): %v", err)
	}
	reqEnabled := httptest.NewRequest(http.MethodGet, "/api/schema/telegram-settings", nil)
	wEnabled := httptest.NewRecorder()
	muxEnabled.ServeHTTP(wEnabled, reqEnabled)
	if wEnabled.Code != http.StatusOK {
		t.Fatalf("enabled module: /api/schema/telegram-settings status = %d, want 200", wEnabled.Code)
	}

	// Disabled: default mvp profile has no channel.telegram → schema 404.
	planDisabled, err := ResolvePlan(&config.Config{ProfileName: string(kernel.ProfileMVP)})
	if err != nil {
		t.Fatalf("ResolvePlan(disabled): %v", err)
	}
	a2 := auth.New([]byte("test-secret-at-least-32-chars-long!!"), 0, 0, st, true)
	jobs2, err := newJobRuntime(st)
	if err != nil {
		t.Fatal(err)
	}
	muxDisabled, err := newMux(
		&config.Config{ProfileName: string(kernel.ProfileMVP)}, a2, st, authRepo, opRepo, setRepo,
		planDisabled, &readinessGate{}, jwtSecret("test-secret-at-least-32-chars-long!!"), jobs2, nil,
		slog.New(slog.NewTextHandler(io.Discard, nil)), cachePort, eventBusPort, newRateLimiters(), nil,
	)
	if err != nil {
		t.Fatalf("newMux(disabled): %v", err)
	}
	reqDisabled := httptest.NewRequest(http.MethodGet, "/api/schema/telegram-settings", nil)
	wDisabled := httptest.NewRecorder()
	muxDisabled.ServeHTTP(wDisabled, reqDisabled)
	if wDisabled.Code != http.StatusNotFound {
		t.Fatalf("disabled module: /api/schema/telegram-settings status = %d, want 404", wDisabled.Code)
	}
	for _, probe := range []struct {
		method string
		path   string
	}{
		{method: http.MethodGet, path: "/api/channel/telegram/settings"},
		{method: http.MethodPost, path: "/api/channel/telegram/lease/acquire"},
		{method: http.MethodPost, path: "/api/channel/telegram/lease/heartbeat"},
		{method: http.MethodPost, path: "/api/channel/telegram/lease/release"},
		{method: http.MethodPost, path: "/api/channel/telegram/webhook"},
	} {
		req := httptest.NewRequest(probe.method, probe.path, nil)
		w := httptest.NewRecorder()
		muxDisabled.ServeHTTP(w, req)
		if w.Code != http.StatusNotFound {
			t.Fatalf("disabled module: %s %s status = %d, want 404", probe.method, probe.path, w.Code)
		}
	}
}
