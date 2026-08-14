// Package recyclebin provides the admin.recycle-bin module surface as a
// kernel.Provider (S-12 · GOAL-012 D-002): deleted-row snapshots with
// browse/restore/purge, recycle.read / recycle.write permission keys and the
// recycle.restore / recycle.purge audit events. Delete hooks are wired through
// the resource factory (Resource.Trash → handler.TrashRecorder).
package recyclebin

import (
	"context"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/handler"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	authsessiondata "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession/systemdata"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/recyclebin/manifest"
	recycleschema "github.com/magicvr/schema-ui-core/apps/api/internal/modules/recyclebin/schema"
	recyclestore "github.com/magicvr/schema-ui-core/apps/api/internal/modules/recyclebin/store"
)

// ModuleID is the stable admin.recycle-bin module identifier.
const ModuleID = "admin.recycle-bin"

// Provider implements kernel.Provider for admin.recycle-bin.
type Provider struct {
	a          *auth.Authenticator
	service    *Service
	operations operationlog.Recorder
}

// New constructs the recycle provider.
func New(a *auth.Authenticator, service *Service, operations operationlog.Recorder) *Provider {
	return &Provider{a: a, service: service, operations: operations}
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
				"GET /api/recycle-bin", "GET /api/recycle-bin/{id}",
				"POST /api/recycle-bin/{id}/restore", "DELETE /api/recycle-bin/{id}",
			},
			Pages:       []string{"recycle-bin"},
			Navigation:  []string{"menu_recycle_bin"},
			Permissions: []string{"recycle.read", "recycle.write"},
			Fragments:   []string{"recycle-bin"},
		},
	}
}

func (p *Provider) CompiledPersistence() ([]kernel.MigrationContribution, error) {
	return nil, nil // tables are owned by the recyclebin/migration provider (0025)
}

func (p *Provider) Register(ctx context.Context, reg kernel.Registrar) error {
	for _, route := range handler.RecycleBinRoutes(p.a, p.service, p.operations, ModuleID) {
		if err := reg.HTTP(route); err != nil {
			return err
		}
	}
	for _, pageID := range []string{"recycle-bin"} {
		if err := reg.Schema(kernel.PageContribution{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: pageID},
			PageID:               pageID,
			Resources:            []string{"recycle-bin"},
			Actions:              []string{"list", "detail", "restore", "purge"},
			DataSource:           "/api/recycle-bin",
			Owner:                ModuleID,
			Document:             recycleschema.SchemaDocuments()[pageID],
		}); err != nil {
			return err
		}
	}
	for _, permission := range []kernel.PermissionContribution{
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "recycle.read"}, Permission: "recycle.read", Resource: "recycle-bin", Action: "read", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "recycle.write"}, Permission: "recycle.write", Resource: "recycle-bin", Action: "write", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
	} {
		if err := reg.Authorization(permission); err != nil {
			return err
		}
	}
	if err := reg.Navigation(kernel.NavigationContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "menu_recycle_bin"},
		NodeID:               "menu_recycle_bin",
		PageID:               "recycle-bin",
		Order:                8,
		Label:                "Recycle bin",
		Visibility:           authsessiondata.PolicyAdmin,
		Permission:           "recycle.read",
		SystemDataVersion:    authsessiondata.SystemDataVersion,
	}); err != nil {
		return err
	}
	return reg.Manifest(kernel.FragmentContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "recycle-bin"},
		FragmentID:           "recycle-bin",
		ProtocolVersion:      "2.7",
		RequiredCapabilities: []string{"manifest", "navigation"},
		JSON:                 manifest.FragmentJSON,
	})
}

var _ = recyclestore.Item{}
