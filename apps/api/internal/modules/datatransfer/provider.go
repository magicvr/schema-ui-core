// Package datatransfer provides the admin.data-transfer module surface as a
// kernel.Provider (F-02 · GOAL-004 D-002): schema-driven CSV export and CSV
// import with per-row error reporting. No manifest fragment (shared capability
// consumed via toolbar actions on the users/roles pages).
package datatransfer

import (
	"context"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/handler"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	authsession "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession"
	authsessiondata "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession/systemdata"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
)

const ModuleID = "admin.data-transfer"

// Provider implements kernel.Provider for admin.data-transfer.
type Provider struct {
	a          *auth.Authenticator
	users      *authsession.Repository
	roles      *authsession.Repository
	operations operationlog.Recorder
	uploadDir  string
}

// New constructs the data-transfer provider with framework-agnostic
// dependencies. uploadDir is the shared uploads storage root (same directory
// as the central upload endpoint).
func New(a *auth.Authenticator, repository *authsession.Repository, operations operationlog.Recorder, uploadDir string) *Provider {
	return &Provider{a: a, users: repository, roles: repository, operations: operations, uploadDir: uploadDir}
}

func (p *Provider) Descriptor() kernel.Module {
	return kernel.Module{
		ID:             ModuleID,
		Version:        "2.0.0",
		KernelAPIRange: ">=2.0 <3.0",
		DependsOn:      []string{"core.auth-session", "core.schema-render", "core.operationlog"},
		Requires:       kernel.StandardAdminCapabilities(),
		Contributions: kernel.ContributionKeys{
			Routes:      []string{"GET /api/export/{resource}", "POST /api/import/{resource}"},
			Permissions: []string{"data.export", "data.import"},
		},
	}
}

func (p *Provider) CompiledPersistence() ([]kernel.MigrationContribution, error) {
	return nil, nil // no migrations owned by this module
}

func (p *Provider) Register(ctx context.Context, reg kernel.Registrar) error {
	for _, route := range handler.ExportRoutes(p.a, p.users, p.roles, p.operations, ModuleID) {
		if err := reg.HTTP(route); err != nil {
			return err
		}
	}
	for _, route := range handler.ImportRoutes(p.a, p.users, p.operations, p.uploadDir, ModuleID) {
		if err := reg.HTTP(route); err != nil {
			return err
		}
	}
	for _, permission := range []kernel.PermissionContribution{
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "data.export"}, Permission: "data.export", Resource: "data", Action: "export", PolicyID: authsessiondata.PolicyAdminEditor, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "data.import"}, Permission: "data.import", Resource: "data", Action: "import", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
	} {
		if err := reg.Authorization(permission); err != nil {
			return err
		}
	}
	return nil
}
