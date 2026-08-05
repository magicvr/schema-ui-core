// Package settings provides the admin.settings module surface as a kernel.Provider
// (R4 C4.1). It reuses the central settings route adapter so the provider
// surface is behavior-identical; the manifest fragment is module-owned.
package settings

import (
	"context"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/handler"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	authsessiondata "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession/systemdata"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/settings/manifest"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/settings/migration"
	settingsrepository "github.com/magicvr/schema-ui-core/apps/api/internal/modules/settings/repository"
)

const ModuleID = migration.ModuleID

// Provider implements kernel.Provider for admin.settings.
type Provider struct {
	a          *auth.Authenticator
	repository *settingsrepository.Repository
	operations operationlog.Recorder
}

// New constructs the settings provider with framework-agnostic dependencies.
func New(a *auth.Authenticator, repository *settingsrepository.Repository, operations operationlog.Recorder) *Provider {
	return &Provider{a: a, repository: repository, operations: operations}
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
				"GET /api/branding", "GET /api/settings", "GET /api/settings/{id}", "PATCH /api/settings/{id}",
			},
			Pages:       []string{"settings"},
			Navigation:  []string{"menu_settings"},
			Permissions: []string{"settings.read", "settings.write"},
			Fragments:   []string{"settings"},
		},
	}
}

func (p *Provider) CompiledPersistence() ([]kernel.MigrationContribution, error) {
	return migration.Descriptors(), nil
}

func (p *Provider) Register(ctx context.Context, reg kernel.Registrar) error {
	for _, route := range handler.SettingsRoutes(p.a, p.repository, p.operations, ModuleID) {
		if err := reg.HTTP(route); err != nil {
			return err
		}
	}
	if err := reg.Schema(kernel.PageContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "settings"},
		PageID:               "settings",
		Resources:            []string{"settings"},
		Actions:              []string{"list", "detail", "update"},
		DataSource:           "/api/settings",
		Owner:                ModuleID,
	}); err != nil {
		return err
	}
	for _, permission := range []kernel.PermissionContribution{
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "settings.read"}, Permission: "settings.read", Resource: "settings", Action: "read", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "settings.write"}, Permission: "settings.write", Resource: "settings", Action: "write", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
	} {
		if err := reg.Authorization(permission); err != nil {
			return err
		}
	}
	if err := reg.Navigation(kernel.NavigationContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "menu_settings"},
		NodeID:               "menu_settings",
		PageID:               "settings",
		Order:                1,
		Label:                "Settings",
		Visibility:           authsessiondata.PolicyAdmin,
		Permission:           "settings.read",
		SystemDataVersion:    authsessiondata.SystemDataVersion,
	}); err != nil {
		return err
	}
	return reg.Manifest(kernel.FragmentContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "settings"},
		FragmentID:           "settings",
		ProtocolVersion:      "2.7",
		RequiredCapabilities: []string{"manifest", "navigation"},
		JSON:                 manifest.FragmentJSON,
	})
}
