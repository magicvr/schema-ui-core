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
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	authsession "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession"
	authsessiondata "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession/systemdata"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/users/manifest"
)

const ModuleID = "admin.users"

// Provider implements kernel.Provider for admin.users.
type Provider struct {
	a          *auth.Authenticator
	repository *authsession.Repository
	operations handler.OperationRecorder
}

// New constructs the users provider with framework-agnostic dependencies.
func New(a *auth.Authenticator, repository *authsession.Repository, operations handler.OperationRecorder) *Provider {
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
				"GET /api/users", "GET /api/users/{id}", "POST /api/users",
				"PATCH /api/users/{id}", "DELETE /api/users/{id}",
			},
			Pages:       []string{"users"},
			Navigation:  []string{"menu_users"},
			Permissions: []string{"users.read", "users.write"},
			Fragments:   []string{"users"},
		},
	}
}

func (p *Provider) CompiledPersistence() ([]kernel.MigrationContribution, error) {
	return nil, nil // account/RBAC migrations are owned by core.auth-session
}

func (p *Provider) Register(ctx context.Context, reg kernel.Registrar) error {
	for _, route := range handler.ResourceRoutes(p.a, handler.UsersResource(p.repository, p.operations), ModuleID) {
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
	}); err != nil {
		return err
	}
	for _, permission := range []kernel.PermissionContribution{
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "users.read"}, Permission: "users.read", Resource: "users", Action: "read", PolicyID: authsessiondata.PolicyAdminEditorViewer, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "users.write"}, Permission: "users.write", Resource: "users", Action: "write", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
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
