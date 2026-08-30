// Package users provides the admin.users module surface as a kernel.Provider
// (R4 C3.2). It reuses the generic schema-driven resource factory so the
// provider-generated HTTP surface is byte-compatible with the current central
// registration (freeze package §7 step 2). Framework-agnostic: no go.uber.org/fx
// import; the composition root constructs the provider with plain dependencies.
package users

import (
	"context"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/handler"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	authsession "github.com/magicvr/schema-ui-core/apps/api/modules/authsession"
	authsessiondata "github.com/magicvr/schema-ui-core/apps/api/modules/authsession/systemdata"
	"github.com/magicvr/schema-ui-core/apps/api/modules/operationlog"
	"github.com/magicvr/schema-ui-core/apps/api/modules/users/manifest"
	usersschema "github.com/magicvr/schema-ui-core/apps/api/modules/users/schema"
)

const ModuleID = "admin.users"

// Provider implements kernel.Provider for admin.users.
type Provider struct {
	a          *auth.Authenticator
	repository *authsession.Repository
	operations operationlog.Recorder
	// mailSender is THE composed kernel.MailSender (workspace-019 R3):
	// invitation letters ride it when an invite carries a target email.
	mailSender kernel.MailSender
	// publicBaseURL is the optional canonical external origin for emailed
	// invite links (W13 F-006 · GOAL-013 A-001); empty keeps the
	// request-derived fallback.
	publicBaseURL string
}

// New constructs the users provider with framework-agnostic dependencies.
func New(a *auth.Authenticator, repository *authsession.Repository, operations operationlog.Recorder, mailSender kernel.MailSender, publicBaseURL string) *Provider {
	return &Provider{a: a, repository: repository, operations: operations, mailSender: mailSender, publicBaseURL: publicBaseURL}
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
				"GET /api/users", "GET /api/users/{id}", "POST /api/users",
				"PATCH /api/users/{id}", "DELETE /api/users/{id}",
				"POST /api/users/batch-delete",
				"GET /api/users/invites", "POST /api/users/invites",
				"DELETE /api/users/invites/{id}", "POST /api/users/invites/{id}/resend",
			},
			Pages:       []string{"users", "users-invites"},
			Navigation:  []string{"menu_users"},
			Permissions: []string{"users.read", "users.write", "users.invite"},
			Fragments:   []string{"users"},
		},
	}
}

func (p *Provider) CompiledPersistence() ([]kernel.MigrationContribution, error) {
	return nil, nil // account/RBAC migrations are owned by core.auth-session
}

func (p *Provider) Register(ctx context.Context, reg kernel.Registrar) error {
	for _, route := range handler.ResourceRoutes(p.a, handler.UsersResourceWithNotifier(p.repository, p.operations, p.repository), ModuleID) {
		if err := reg.HTTP(route); err != nil {
			return err
		}
	}
	// workspace-019 R3 (GOAL-004 D-001 §3): invitation management quartet,
	// gated by the users.invite permission inside the handler.
	for _, route := range handler.InviteAdminRoutes(p.a, p.repository, p.mailSender, p.operations, ModuleID, p.publicBaseURL) {
		if err := reg.HTTP(route); err != nil {
			return err
		}
	}
	if err := reg.Schema(kernel.PageContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "users"},
		PageID:               "users",
		Resources:            []string{"users"},
		Actions:              []string{"list", "create", "detail", "update", "delete"},
		DataSource:           "/api/users",
		Owner:                ModuleID,
		Document:             usersschema.SchemaDocuments()["users"],
	}); err != nil {
		return err
	}
	// workspace-019 UX polish: invitation management lives on its own child
	// page (data-dictionary → dictionary-entries precedent); the users page
	// toolbar holds the users.invite-gated entry navigation.
	if err := reg.Schema(kernel.PageContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "users-invites"},
		PageID:               "users-invites",
		Resources:            []string{"users", "invites"},
		Actions:              []string{"list", "create"},
		DataSource:           "/api/users/invites",
		Owner:                ModuleID,
		Document:             usersschema.SchemaDocuments()["users-invites"],
	}); err != nil {
		return err
	}
	for _, permission := range []kernel.PermissionContribution{
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "users.read"}, Permission: "users.read", Resource: "users", Action: "read", PolicyID: authsessiondata.PolicyAdminEditorViewer, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "users.write"}, Permission: "users.write", Resource: "users", Action: "write", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "users.invite"}, Permission: "users.invite", Resource: "users", Action: "invite", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
	} {
		if err := reg.Authorization(permission); err != nil {
			return err
		}
	}
	if err := reg.Navigation(kernel.NavigationContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "menu_users"},
		NodeID:               "menu_users",
		PageID:               "users",
		Order:                1,
		Label:                "Users",
		Visibility:           authsessiondata.PolicyAdmin,
		Permission:           "users.read",
		SystemDataVersion:    authsessiondata.SystemDataVersion,
	}); err != nil {
		return err
	}
	return reg.Manifest(kernel.FragmentContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "users"},
		FragmentID:           "users",
		ProtocolVersion:      "2.7",
		RequiredCapabilities: []string{"manifest", "navigation"},
		JSON:                 manifest.FragmentJSON,
	})
}