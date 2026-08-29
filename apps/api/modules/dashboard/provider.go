// Package dashboard provides the admin.dashboard module surface as a
// kernel.Provider (F-01 · GOAL-003 D-002): the production home dashboard page
// (registry display nodes only; no routes/permissions — card data sources are
// the existing users/roles list endpoints readable by every standard role).
package dashboard

import (
	"context"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/modules/dashboard/manifest"
	dashboardschema "github.com/magicvr/schema-ui-core/apps/api/modules/dashboard/schema"
	authsessiondata "github.com/magicvr/schema-ui-core/apps/api/modules/authsession/systemdata"
)

const ModuleID = "admin.dashboard"

// Provider implements kernel.Provider for admin.dashboard.
type Provider struct{}

// New constructs the dashboard provider (no framework dependencies).
func New() *Provider {
	return &Provider{}
}

func (p *Provider) Descriptor() kernel.Module {
	return kernel.Module{
		ID:             ModuleID,
		Version:        "2.0.0",
		KernelAPIRange: ">=2.0 <3.0",
		DependsOn:      []string{"core.auth-session", "core.navigation-capability", "core.schema-render", "core.manifest-route"},
		Requires:       kernel.StandardAdminCapabilities(),
		Contributions: kernel.ContributionKeys{
			Pages:      []string{"dashboard"},
			Navigation: []string{"menu_dashboard"},
			Fragments:  []string{"dashboard"},
		},
	}
}

func (p *Provider) CompiledPersistence() ([]kernel.MigrationContribution, error) {
	return nil, nil
}

func (p *Provider) Register(ctx context.Context, reg kernel.Registrar) error {
	if err := reg.Schema(kernel.PageContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "dashboard"},
		PageID:               "dashboard",
		Resources:            []string{"users", "roles"},
		Actions:              []string{"list"},
		DataSource:           "/api/users",
		Owner:                ModuleID,
		Document:             dashboardschema.SchemaDocuments()["dashboard"],
	}); err != nil {
		return err
	}
	if err := reg.Navigation(kernel.NavigationContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "menu_dashboard"},
		NodeID:               "menu_dashboard",
		PageID:               "dashboard",
		Order:                0,
		Label:                "Dashboard",
		Visibility:           authsessiondata.PolicyAdminEditorViewer,
		Permission:           "",
		SystemDataVersion:    authsessiondata.SystemDataVersion,
	}); err != nil {
		return err
	}
	return reg.Manifest(kernel.FragmentContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "dashboard"},
		FragmentID:           "dashboard",
		ProtocolVersion:      "2.7",
		RequiredCapabilities: []string{"manifest", "navigation"},
		JSON:                 manifest.FragmentJSON,
	})
}