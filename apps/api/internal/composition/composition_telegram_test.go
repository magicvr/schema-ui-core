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
