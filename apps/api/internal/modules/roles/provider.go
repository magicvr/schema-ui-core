// Package roles provides the admin.roles module surface as a kernel.Provider
// (R4 C3.2). It reuses the generic schema-driven resource factory so the
// provider-generated HTTP surface is byte-compatible with the current central
// registration (freeze package §7 step 2). Framework-agnostic: no go.uber.org/fx
// import; the composition root constructs the provider with plain dependencies.
package roles

import (
	"context"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/handler"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	authsessiondata "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession/systemdata"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/roles/manifest"
	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
)

const ModuleID = "admin.roles"

// Provider implements kernel.Provider for admin.roles.
type Provider struct {
	a  *auth.Authenticator
	st *store.Store
}

// New constructs the roles provider with framework-agnostic dependencies.
func New(a *auth.Authenticator, st *store.Store) *Provider {
	return &Provider{a: a, st: st}
}

func (p *Provider) Descriptor() kernel.Module {
	return kernel.Module{
		ID:             ModuleID,
		Version:        "2.0.0",
		KernelAPIRange: ">=2.0 <3.0",
		DependsOn:      []string{"core.auth-session", "core.navigation-capability", "core.schema-render", "core.operationlog"},
		Requires:       kernel.StandardAdminCapabilities(),
		Contributions: kernel.ContributionKeys{
			Routes: []string{
				"GET /api/roles", "GET /api/roles/{id}", "POST /api/roles",
				"PATCH /api/roles/{id}", "DELETE /api/roles/{id}",
			},
			Pages:       []string{"roles"},
			Navigation:  []string{"menu_roles"},
			Permissions: []string{"roles.read", "roles.write", "roles.assign"},
			Fragments:   []string{"roles"},
		},
	}
}

func (p *Provider) CompiledPersistence() ([]kernel.MigrationContribution, error) {
	return nil, nil // no module-owned migrations in C3; existing ledger is store-owned
}

func (p *Provider) Register(ctx context.Context, reg kernel.Registrar) error {
	for _, route := range handler.ResourceRoutes(p.a, handler.RolesResource(p.st), ModuleID) {
		if err := reg.HTTP(route); err != nil {
			return err
		}
	}
	if err := reg.Schema(kernel.PageContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "roles"},
		PageID:               "roles",
		Resources:            []string{"roles"},
		Actions:              []string{"list", "create", "detail", "update", "delete"},
		DataSource:           "/api/roles",
		Owner:                ModuleID,
	}); err != nil {
		return err
	}
	for _, permission := range []kernel.PermissionContribution{
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "roles.read"}, Permission: "roles.read", Resource: "roles", Action: "read", PolicyID: authsessiondata.PolicyAdminEditorViewer, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "roles.write"}, Permission: "roles.write", Resource: "roles", Action: "write", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "roles.assign"}, Permission: "roles.assign", Resource: "roles", Action: "assign", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
	} {
		if err := reg.Authorization(permission); err != nil {
			return err
		}
	}
	if err := reg.Navigation(kernel.NavigationContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "menu_roles"},
		NodeID:               "menu_roles",
		PageID:               "roles",
		Order:                2,
		Label:                "Roles",
		Visibility:           authsessiondata.PolicyAdmin,
		Permission:           "roles.read",
		SystemDataVersion:    authsessiondata.SystemDataVersion,
	}); err != nil {
		return err
	}
	return reg.Manifest(kernel.FragmentContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "roles"},
		FragmentID:           "roles",
		ProtocolVersion:      "2.7",
		RequiredCapabilities: []string{"manifest", "navigation"},
		JSON:                 manifest.FragmentJSON,
	})
}
