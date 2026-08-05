// Package activity provides the admin.activity module surface as a kernel.Provider
// (R4 C4.2). It reuses the read-only operations resource factory; the manifest
// fragment is module-owned. Operation-log writes remain owned by core.operationlog;
// Activity is only the query/UI surface and disabling it must not stop the writer.
package activity

import (
	"context"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/handler"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/activity/manifest"
	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
)

const ModuleID = "admin.activity"

// Provider implements kernel.Provider for admin.activity.
type Provider struct {
	a  *auth.Authenticator
	st *store.Store
}

// New constructs the activity provider with framework-agnostic dependencies.
func New(a *auth.Authenticator, st *store.Store) *Provider {
	return &Provider{a: a, st: st}
}

func (p *Provider) Descriptor() kernel.Module {
	return kernel.Module{
		ID:             ModuleID,
		Version:        "2.0.0",
		KernelAPIRange: ">=2.0 <3.0",
		DependsOn:      []string{"core.auth-session", "core.navigation-capability", "core.manifest-route", "core.operationlog"},
		Requires:       kernel.StandardAdminCapabilities(),
		Contributions: kernel.ContributionKeys{
			Routes:      []string{"GET /api/operations", "GET /api/operations/{id}"},
			Pages:       []string{"activity"},
			Navigation:  []string{"menu_activity"},
			Permissions: []string{"operations.read"},
			Fragments:   []string{"activity"},
		},
	}
}

func (p *Provider) CompiledPersistence() ([]kernel.MigrationContribution, error) {
	return nil, nil
}

func (p *Provider) Register(ctx context.Context, reg kernel.Registrar) error {
	for _, route := range handler.ResourceRoutes(p.a, handler.OperationsResource(p.st), ModuleID) {
		if err := reg.HTTP(route); err != nil {
			return err
		}
	}
	if err := reg.Schema(kernel.PageContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "activity"},
		PageID:               "activity",
		Resources:            []string{"operations"},
		Actions:              []string{"list", "detail"},
		DataSource:           "/api/operations",
		Owner:                ModuleID,
	}); err != nil {
		return err
	}
	if err := reg.Authorization(kernel.PermissionContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "operations.read"},
		Permission:           "operations.read", Resource: "operations", Action: "read",
	}); err != nil {
		return err
	}
	if err := reg.Navigation(kernel.NavigationContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "menu_activity"},
		NodeID:               "menu_activity",
		Order:                2,
		Label:                "Activity",
		Visibility:           "operations.read",
		Permission:           "operations.read",
	}); err != nil {
		return err
	}
	return reg.Manifest(kernel.FragmentContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "activity"},
		FragmentID:           "activity",
		ProtocolVersion:      "2.7",
		RequiredCapabilities: []string{"manifest", "navigation"},
		JSON:                 manifest.FragmentJSON,
	})
}
