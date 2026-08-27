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
	authsession "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
	settingsconfiguration "github.com/magicvr/schema-ui-core/apps/api/internal/modules/settings/configuration"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/settings/manifest"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/settings/migration"
	settingsrepository "github.com/magicvr/schema-ui-core/apps/api/internal/modules/settings/repository"
	settingsschema "github.com/magicvr/schema-ui-core/apps/api/internal/modules/settings/schema"
)

const ModuleID = migration.ModuleID

// Provider implements kernel.Provider for admin.settings.
type Provider struct {
	a          *auth.Authenticator
	repository *settingsrepository.Repository
	operations operationlog.Recorder
	assets     *handler.BrandingAssetStore
	// authRepo feeds the password-policy configuration surface (workspace-019
	// R3 · GOAL-004 D-001 §2): the policy row is owned by core.auth-session.
	authRepo *authsession.Repository
}

// New constructs the settings provider with framework-agnostic dependencies.
// assets is the W9 dedicated brand-asset store (may be nil in narrow tests
// that never exercise brand uploads).
func New(a *auth.Authenticator, repository *settingsrepository.Repository, operations operationlog.Recorder, assets *handler.BrandingAssetStore, authRepo *authsession.Repository) *Provider {
	return &Provider{a: a, repository: repository, operations: operations, assets: assets, authRepo: authRepo}
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
				"GET /api/branding", "GET /api/settings", "GET /api/settings/{id}", "PATCH /api/settings/{id}", "POST /api/settings/{id}/reset",
				"POST /api/branding/assets", "GET /api/branding/assets/{id}",
				"GET /api/settings/password-policy", "PATCH /api/settings/password-policy",
			},
			// W26 (GOAL-038 D-001 §2.2): the mail console + outbound log are
			// standalone pages under the SAME settings.read gate (no new
			// permission keys); the settings page no longer hosts them.
			Pages:            []string{"settings", "mail", "mail-outbox"},
			Navigation:       []string{"menu_settings", "menu_mail", "menu_mail_outbox"},
			Permissions:      []string{"settings.read", "settings.write"},
			ConfigNamespaces: []string{settingsconfiguration.Namespace},
			Fragments:        []string{"settings"},
		},
	}
}

func (p *Provider) CompiledPersistence() ([]kernel.MigrationContribution, error) {
	return migration.Descriptors(), nil
}

func (p *Provider) Register(ctx context.Context, reg kernel.Registrar) error {
	configuration := settingsconfiguration.Contribution()
	if err := reg.Configuration(configuration); err != nil {
		return err
	}
	for _, route := range handler.SettingsRoutes(p.a, p.repository, p.operations, ModuleID, configuration.Namespace, p.assets) {
		if err := reg.HTTP(route); err != nil {
			return err
		}
	}
	// workspace-019 R3 (GOAL-004 D-001 §2): password-policy configuration —
	// GET/PATCH on the admin.settings tab extension.
	for _, route := range handler.PasswordPolicyRoutes(p.a, p.authRepo, ModuleID) {
		if err := reg.HTTP(route); err != nil {
			return err
		}
	}
	// W9 (GOAL-010): dedicated brand-asset surface (auth upload + public GET).
	for _, route := range handler.BrandingAssetRoutes(p.a, p.assets, ModuleID) {
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
		Document:             settingsschema.SchemaDocuments()["settings"],
	}); err != nil {
		return err
	}
	// W26 (GOAL-038 D-001 §2.2): standalone mail console (channel config +
	// test send via the shared custom component) and all-channel outbound
	// log. Both ride the settings.read gate through the menu visibility; the
	// underlying APIs keep their existing per-route gates.
	if err := reg.Schema(kernel.PageContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "mail"},
		PageID:               "mail",
		Resources:            []string{"mail"},
		Actions:              []string{"detail", "update"},
		DataSource:           "/api/mail/config",
		Owner:                ModuleID,
		Document:             settingsschema.SchemaDocuments()["mail"],
	}); err != nil {
		return err
	}
	if err := reg.Schema(kernel.PageContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "mail-outbox"},
		PageID:               "mail-outbox",
		Resources:            []string{"mail-outbox"},
		Actions:              []string{"list", "detail"},
		DataSource:           "/api/mail/outbox",
		Owner:                ModuleID,
		Document:             settingsschema.SchemaDocuments()["mail-outbox"],
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
	// W26 (GOAL-038 D-001 §2.2): standalone mail pages reuse the settings.read
	// permission — the user red line is NO new permission keys.
	for _, node := range []kernel.NavigationContribution{
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "menu_mail"},
			NodeID:               "menu_mail",
			PageID:               "mail",
			Label:                "Mail console",
			Visibility:           authsessiondata.PolicyAdmin,
			Permission:           "settings.read",
			SystemDataVersion:    authsessiondata.SystemDataVersion,
		},
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "menu_mail_outbox"},
			NodeID:               "menu_mail_outbox",
			PageID:               "mail-outbox",
			Label:                "Outbound email log",
			Visibility:           authsessiondata.PolicyAdmin,
			Permission:           "settings.read",
			SystemDataVersion:    authsessiondata.SystemDataVersion,
		},
	} {
		if err := reg.Navigation(node); err != nil {
			return err
		}
	}
	if err := reg.Manifest(kernel.FragmentContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "settings"},
		FragmentID:           "settings",
		ProtocolVersion:      "2.7",
		RequiredCapabilities: []string{"manifest", "navigation"},
		JSON:                 manifest.FragmentJSON,
	}); err != nil {
		return err
	}
	return nil
}
