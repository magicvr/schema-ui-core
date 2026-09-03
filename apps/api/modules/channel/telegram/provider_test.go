package telegram_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	moduletg "github.com/magicvr/schema-ui-core/apps/api/modules/channel/telegram"
)

type mockRegistrar struct {
	routes []kernel.RouteContribution
}

func (m *mockRegistrar) HTTP(r kernel.RouteContribution) error {
	m.routes = append(m.routes, r)
	return nil
}
func (m *mockRegistrar) Schema(p kernel.PageContribution) error                 { return nil }
func (m *mockRegistrar) Authorization(a kernel.PermissionContribution) error    { return nil }
func (m *mockRegistrar) Navigation(n kernel.NavigationContribution) error       { return nil }
func (m *mockRegistrar) Manifest(f kernel.FragmentContribution) error           { return nil }
func (m *mockRegistrar) Configuration(c kernel.ConfigurationContribution) error { return nil }

func TestTelegramModuleProvider(t *testing.T) {
	dummyWebhook := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	dummySettings := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	p := moduletg.New(dummyWebhook, dummySettings)

	// Check Descriptor
	desc := p.Descriptor()
	if desc.ID != moduletg.ModuleID {
		t.Fatalf("expected module ID %q, got %q", moduletg.ModuleID, desc.ID)
	}
	if len(desc.Contributions.Routes) != 3 {
		t.Fatalf("expected 3 route contributions, got %+v", desc.Contributions.Routes)
	}

	// Check Persistence
	contribs, err := p.CompiledPersistence()
	if err != nil || len(contribs) != 0 {
		t.Fatalf("expected nil persistence, got %v, err=%v", contribs, err)
	}

	// Check Register
	reg := &mockRegistrar{}
	if err := p.Register(context.Background(), reg); err != nil {
		t.Fatalf("Register failed: %v", err)
	}
	if len(reg.routes) != 3 {
		t.Fatalf("expected 3 routes registered, got %d", len(reg.routes))
	}

	// Verify handler invocation through registered route
	req := httptest.NewRequest(http.MethodPost, "/api/channel/telegram/webhook", nil)
	w := httptest.NewRecorder()
	dummyWebhook.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK from registered route handler, got %d", w.Code)
	}
}

func TestTelegramModule_RegisterContributionsIntegration(t *testing.T) {
	dummyWebhook := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	dummySettings := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	p := moduletg.New(dummyWebhook, dummySettings)

	// Verify BuiltinModules includes channel.telegram
	builtin := kernel.BuiltinModules()
	found := false
	for _, m := range builtin {
		if m.ID == moduletg.ModuleID {
			found = true
			if m.Version != p.Descriptor().Version || m.KernelAPIRange != p.Descriptor().KernelAPIRange {
				t.Fatalf("mismatched version/kernel API range between builtin and provider: builtin=%+v, provider=%+v", m, p.Descriptor())
			}
		}
	}
	if !found {
		t.Fatalf("expected channel.telegram to be in kernel.BuiltinModules()")
	}

	// Verify default profiles do NOT contain channel.telegram
	for _, prof := range []string{"mvp", "admin", "demo"} {
		res, err := kernel.ResolveProfile(prof, nil)
		if err != nil {
			t.Fatalf("ResolveProfile(%q) error: %v", prof, err)
		}
		for _, modID := range res.Modules {
			if modID == moduletg.ModuleID {
				t.Fatalf("profile %q must NOT include %q by default", prof, moduletg.ModuleID)
			}
		}
	}

	// Build a plan with channel.telegram enabled
	plan := kernel.Plan{
		Capabilities: []kernel.Capability{kernel.CapabilityHTTP},
		Modules: []kernel.Module{
			p.Descriptor(),
		},
	}

	set, err := kernel.RegisterContributions(context.Background(), plan, []kernel.Provider{p})
	if err != nil {
		t.Fatalf("RegisterContributions failed: %v", err)
	}

	if len(set.Routes) != 3 {
		t.Fatalf("expected 3 routes in ContributionSet, got %d", len(set.Routes))
	}
}
