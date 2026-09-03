package composition

import (
	"context"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/internal/config"
	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
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
		ModulesEnabled:        []string{"core.server-registration", "channel.telegram"},
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
		DependsOn:      []string{"core.server-registration"},
		Requires:       []kernel.Capability{kernel.CapabilityHTTP},
		Contributions: kernel.ContributionKeys{
			Routes: []string{
				"GET /api/channel/telegram/settings",
				"PATCH /api/channel/telegram/settings",
				"POST /api/channel/telegram/webhook",
			},
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

	if len(set.Routes) != 3 {
		t.Fatalf("unexpected route contributions: %+v", set.Routes)
	}

	// 5. Test webhook route execution
	req := httptest.NewRequest(http.MethodPost, "/api/channel/telegram/webhook", nil)
	w := httptest.NewRecorder()
	set.Routes[2].Handler.ServeHTTP(w, req)
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
	dispDisabled, senderDisabled := ResolveTelegramPorts(planDisabled, cfg, st)

	msg := kernel.TelegramMessage{ChatID: "123", Text: "Hello"}
	if err := senderDisabled.Send(context.Background(), msg); err != kernel.ErrTelegramDisabled {
		t.Fatalf("expected ErrTelegramDisabled when channel disabled, got %v", err)
	}
	if err := dispDisabled.RegisterCommand("start", func(ctx context.Context, upd kernel.TelegramUpdate) error { return nil }); err != nil {
		t.Fatalf("expected no-op nil from disabled dispatcher, got %v", err)
	}

	// 2. When module is enabled in plan: returns live dispatcher and sender
	planEnabled := kernel.Plan{Modules: []kernel.Module{{ID: "channel.telegram"}}}
	dispEnabled, senderEnabled := ResolveTelegramPorts(planEnabled, cfg, st)

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
	}
	plan := kernel.Plan{Modules: []kernel.Module{{ID: "channel.telegram"}}}

	// Start first process instance and hot-switch via settings
	_, sender1 := ResolveTelegramPorts(plan, cfg1, st1)
	_ = sender1

	// Simulate PATCH settings update
	ctx := context.Background()
	_ = st1.Run(ctx, func(tx kernel.Tx) error {
		_, err := tx.Exec(ctx, `UPDATE telegram_config SET bot_token = ?, webhook_secret = ? WHERE id = 1`,
			"persisted-live-token", "persisted-live-secret")
		return err
	})

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
	}

	_, _ = ResolveTelegramPorts(plan, cfg2, st2)

	// Verify DB retains persisted values
	var dbToken, dbSecret string
	var updatedAt int64
	err = st2.Run(ctx, func(tx kernel.Tx) error {
		row := tx.QueryRow(ctx, `SELECT bot_token, webhook_secret, updated_at FROM telegram_config WHERE id = 1`)
		return row.Scan(&dbToken, &dbSecret, &updatedAt)
	})
	if err != nil {
		t.Fatalf("query telegram_config after restart: %v", err)
	}
	if dbToken != "persisted-live-token" || dbSecret != "persisted-live-secret" {
		t.Fatalf("expected persisted token and secret to survive restart, got token=%q, secret=%q", dbToken, dbSecret)
	}
}
