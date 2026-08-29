// Package account provides the admin.account module surface as a
// kernel.Provider (F-03 · GOAL-005 D-002): self-service profile/password/
// sessions plus admin enable/disable/unlock. Framework-agnostic: no
// go.uber.org/fx import; the composition root constructs the provider with
// plain dependencies.
package account

import (
	"context"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/handler"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/modules/account/manifest"
	accountschema "github.com/magicvr/schema-ui-core/apps/api/modules/account/schema"
	authsession "github.com/magicvr/schema-ui-core/apps/api/modules/authsession"
	authsessiondata "github.com/magicvr/schema-ui-core/apps/api/modules/authsession/systemdata"
	"github.com/magicvr/schema-ui-core/apps/api/modules/operationlog"
)

const ModuleID = "admin.account"

// Provider implements kernel.Provider for admin.account.
type Provider struct {
	a          *auth.Authenticator
	repository *authsession.Repository
	operations operationlog.Recorder
	// avatarAssets is the account avatar store (W13 T-05).
	avatarAssets *handler.RasterAssetStore
	// mailSender is THE composed kernel.MailSender (workspace-018 R3): the
	// email identity surface sends verification codes through it.
	mailSender kernel.MailSender
}

// New constructs the account provider with framework-agnostic dependencies.
func New(a *auth.Authenticator, repository *authsession.Repository, operations operationlog.Recorder, avatarAssets *handler.RasterAssetStore, mailSender kernel.MailSender) *Provider {
	return &Provider{a: a, repository: repository, operations: operations, avatarAssets: avatarAssets, mailSender: mailSender}
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
				"GET /api/account/profile", "PATCH /api/account/profile",
				"POST /api/account/avatar", "GET /api/account/avatars/{id}",
				"POST /api/account/password", "GET /api/account/sessions",
				"POST /api/account/sessions/{id}/revoke", "POST /api/account/sessions/revoke-others",
				"POST /api/account/email/bind", "POST /api/account/email/verify", "POST /api/account/email/resend",
				"POST /api/users/{id}/enable", "POST /api/users/{id}/disable",
				"POST /api/users/{id}/unlock",
			},
			Pages:       []string{"account"},
			Navigation:  []string{"menu_account"},
			Permissions: []string{"users.enable", "users.disable"},
			Fragments:   []string{"account"},
		},
	}
}

func (p *Provider) CompiledPersistence() ([]kernel.MigrationContribution, error) {
	return nil, nil // migration 0013 is contributed via the migration package provider
}

func (p *Provider) Register(ctx context.Context, reg kernel.Registrar) error {
	for _, route := range handler.AccountSelfRoutes(p.a, p.repository, p.operations, p.avatarAssets, ModuleID, p.repository) {
		if err := reg.HTTP(route); err != nil {
			return err
		}
	}
	// workspace-018 R3: self-service email bind/verify/resend through the ONE
	// composed MailSender (GOAL-004 D-001 §3).
	for _, route := range handler.EmailIdentityRoutes(p.a, p.repository, p.mailSender, ModuleID) {
		if err := reg.HTTP(route); err != nil {
			return err
		}
	}
	if p.avatarAssets != nil {
		for _, route := range handler.AccountAvatarRoutes(p.a, p.avatarAssets, p.repository, p.operations, ModuleID) {
			if err := reg.HTTP(route); err != nil {
				return err
			}
		}
	}
	for _, route := range handler.UserStateRoutes(p.a, p.repository, p.operations, ModuleID, p.repository) {
		if err := reg.HTTP(route); err != nil {
			return err
		}
	}
	if err := reg.Schema(kernel.PageContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "account"},
		PageID:               "account",
		Resources:            []string{"account"},
		Actions:              []string{"detail", "update", "list"},
		DataSource:           "/api/account/sessions",
		Owner:                ModuleID,
		Document:             accountschema.SchemaDocuments()["account"],
	}); err != nil {
		return err
	}
	for _, permission := range []kernel.PermissionContribution{
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "users.enable"}, Permission: "users.enable", Resource: "users", Action: "enable", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "users.disable"}, Permission: "users.disable", Resource: "users", Action: "disable", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
	} {
		if err := reg.Authorization(permission); err != nil {
			return err
		}
	}
	if err := reg.Navigation(kernel.NavigationContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "menu_account"},
		NodeID:               "menu_account",
		PageID:               "account",
		Order:                1,
		Label:                "Account",
		Visibility:           authsessiondata.PolicyAdminEditorViewer,
		Permission:           "",
		SystemDataVersion:    authsessiondata.SystemDataVersion,
	}); err != nil {
		return err
	}
	return reg.Manifest(kernel.FragmentContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "account"},
		FragmentID:           "account",
		ProtocolVersion:      "2.7",
		RequiredCapabilities: []string{"manifest", "navigation"},
		JSON:                 manifest.FragmentJSON,
	})
}
