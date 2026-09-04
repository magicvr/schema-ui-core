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
	routes      []kernel.RouteContribution
	permissions []kernel.PermissionContribution
	pages       []kernel.PageContribution
	navigation  []kernel.NavigationContribution
}

func (m *mockRegistrar) HTTP(r kernel.RouteContribution) error {
	m.routes = append(m.routes, r)
	return nil
}
func (m *mockRegistrar) Schema(p kernel.PageContribution) error {
	m.pages = append(m.pages, p)
	return nil
}
func (m *mockRegistrar) Authorization(a kernel.PermissionContribution) error {
	m.permissions = append(m.permissions, a)
	return nil
}
func (m *mockRegistrar) Navigation(n kernel.NavigationContribution) error {
	m.navigation = append(m.navigation, n)
	return nil
}
func (m *mockRegistrar) Manifest(f kernel.FragmentContribution) error           { return nil }
func (m *mockRegistrar) Configuration(c kernel.ConfigurationContribution) error { return nil }

func TestTelegramModuleProvider(t *testing.T) {
	dummyWebhook := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	dummySettings := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	dummyLease := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	dummyOperator := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	p := moduletg.New(dummyWebhook, dummySettings, dummyLease, dummyOperator)

	// Check Descriptor
	desc := p.Descriptor()
	if desc.ID != moduletg.ModuleID {
		t.Fatalf("expected module ID %q, got %q", moduletg.ModuleID, desc.ID)
	}
	if len(desc.Contributions.Routes) != 10 {
		t.Fatalf("expected 10 route contributions, got %+v", desc.Contributions.Routes)
	}
	if len(desc.Contributions.Permissions) != 2 || desc.Contributions.Permissions[0] != "telegram.operator.read" || desc.Contributions.Permissions[1] != "telegram.operator.write" {
		t.Fatalf("unexpected operator permissions: %+v", desc.Contributions.Permissions)
	}
	// GOAL-006 R5: telegram-settings page + menu_telegram navigation declared.
	if len(desc.Contributions.Pages) != 1 || desc.Contributions.Pages[0] != "telegram-settings" {
		t.Fatalf("expected telegram-settings page contribution, got %+v", desc.Contributions.Pages)
	}
	if len(desc.Contributions.Navigation) != 1 || desc.Contributions.Navigation[0] != "menu_telegram" {
		t.Fatalf("expected menu_telegram navigation contribution, got %+v", desc.Contributions.Navigation)
	}

	// Check Persistence
	contribs, err := p.CompiledPersistence()
	if err != nil || len(contribs) != 4 || contribs[0].Version != 66 || contribs[0].Name != "telegram_config" || contribs[1].Version != 67 || contribs[1].Name != "telegram_config_connection" || contribs[2].Version != 68 || contribs[2].Name != "telegram_ingress" || contribs[3].Version != 69 || contribs[3].Name != "telegram_outbound" {
		t.Fatalf("unexpected persistence: %+v, err=%v", contribs, err)
	}

	// Check Register
	reg := &mockRegistrar{}
	if err := p.Register(context.Background(), reg); err != nil {
		t.Fatalf("Register failed: %v", err)
	}
	if len(reg.routes) != 10 {
		t.Fatalf("expected 10 routes registered, got %d", len(reg.routes))
	}
	if len(reg.permissions) != 2 || reg.permissions[0].Permission != "telegram.operator.read" || reg.permissions[1].Permission != "telegram.operator.write" {
		t.Fatalf("unexpected registered permissions: %+v", reg.permissions)
	}
	if len(reg.pages) != 1 || reg.pages[0].PageID != "telegram-settings" || reg.pages[0].DataSource != "/api/channel/telegram/settings" {
		t.Fatalf("unexpected page contributions: %+v", reg.pages)
	}
	if len(reg.navigation) != 1 || reg.navigation[0].NodeID != "menu_telegram" || reg.navigation[0].PageID != "telegram-settings" {
		t.Fatalf("unexpected navigation contributions: %+v", reg.navigation)
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
	dummyLease := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	dummyOperator := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	p := moduletg.New(dummyWebhook, dummySettings, dummyLease, dummyOperator)

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

	// Resolve a plan with channel.telegram enabled through the registry: its
	// DependsOn pulls in admin.settings (whose settings.read permission the
	// menu_telegram nav references — R-001 / A-002).
	registry, err := kernel.NewRegistry(kernel.BuiltinModules())
	if err != nil {
		t.Fatalf("NewRegistry: %v", err)
	}
	plan, err := registry.Resolve([]string{
		"core.server-registration", "core.auth-session", "core.schema-render",
		"core.manifest-route", "core.navigation-capability", "core.operationlog", "admin.settings", moduletg.ModuleID,
	})
	if err != nil {
		t.Fatalf("registry.Resolve: %v", err)
	}

	set, err := kernel.RegisterContributions(context.Background(), plan, []kernel.Provider{p, settingsPermissionStub{}})
	if err != nil {
		t.Fatalf("RegisterContributions failed: %v", err)
	}

	if len(set.Routes) != 10 {
		t.Fatalf("expected 10 routes in ContributionSet, got %d", len(set.Routes))
	}
	operatorPermissions := map[string]bool{}
	for _, permission := range set.Permissions {
		if permission.ModuleID == moduletg.ModuleID {
			operatorPermissions[permission.Permission] = true
		}
	}
	if len(operatorPermissions) != 2 || !operatorPermissions["telegram.operator.read"] || !operatorPermissions["telegram.operator.write"] {
		t.Fatalf("expected operator permissions in ContributionSet, got %+v", set.Permissions)
	}
	if len(set.Pages) != 1 || set.Pages[0].PageID != "telegram-settings" {
		t.Fatalf("expected telegram-settings page in ContributionSet, got %+v", set.Pages)
	}
	if len(set.Navigation) != 1 || set.Navigation[0].NodeID != "menu_telegram" {
		t.Fatalf("expected menu_telegram navigation in ContributionSet, got %+v", set.Navigation)
	}
	if set.Navigation[0].Permission != "settings.read" {
		t.Fatalf("expected menu_telegram to ride settings.read, got %q", set.Navigation[0].Permission)
	}
	if len(set.Fragments) != 1 || set.Fragments[0].FragmentID != "telegram-settings" {
		t.Fatalf("expected telegram-settings fragment in ContributionSet, got %+v", set.Fragments)
	}
}

// settingsPermissionStub contributes the settings.read permission on behalf of
// admin.settings so the module unit test can assert menu_telegram rides it
// without wiring the whole settings module (production wires the real provider).
type settingsPermissionStub struct{}

func (settingsPermissionStub) Descriptor() kernel.Module {
	// Match the canonical admin.settings descriptor from BuiltinModules so
	// descriptorsMatch passes when both are in the same ContributionSet.
	for _, m := range kernel.BuiltinModules() {
		if m.ID == "admin.settings" {
			return m
		}
	}
	panic("admin.settings missing from BuiltinModules")
}

func (settingsPermissionStub) CompiledPersistence() ([]kernel.MigrationContribution, error) {
	return nil, nil
}

func (settingsPermissionStub) Register(_ context.Context, reg kernel.Registrar) error {
	return reg.Authorization(kernel.PermissionContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: "admin.settings", Key: "settings.read"},
		Permission:           "settings.read",
		Resource:             "settings",
		Action:               "read",
		PolicyID:             "system.admin",
		SystemDataVersion:    1,
	})
}
