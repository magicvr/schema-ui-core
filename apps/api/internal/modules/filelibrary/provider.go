// Package filelibrary provides the admin.file-library module surface as a
// kernel.Provider (S-02 · GOAL-007 D-002): a unified file/attachment library
// over the shared upload store (C-09) — schema-driven list/detail, download,
// hard delete and the upload-confirmation endpoint, with files.read /
// files.write / files.delete permission keys and operation-log events.
package filelibrary

import (
	"context"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/handler"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	authsessiondata "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession/systemdata"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/filelibrary/manifest"
	filelibraryschema "github.com/magicvr/schema-ui-core/apps/api/internal/modules/filelibrary/schema"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
)

// ModuleID is the stable admin.file-library module identifier.
const ModuleID = "admin.file-library"

// Provider implements kernel.Provider for admin.file-library.
type Provider struct {
	a          *auth.Authenticator
	operations operationlog.Recorder
	uploadDir  string
}

// New constructs the file-library provider with framework-agnostic
// dependencies. uploadDir is the shared uploads storage root (the same
// directory as the central upload endpoint).
func New(a *auth.Authenticator, operations operationlog.Recorder, uploadDir string) *Provider {
	return &Provider{a: a, operations: operations, uploadDir: uploadDir}
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
				"GET /api/library/files", "GET /api/library/files/{id}",
				"GET /api/library/files/{id}/download",
				"DELETE /api/library/files/{id}",
				"POST /api/library/files/upload",
			},
			Pages:       []string{"file-library"},
			Navigation:  []string{"menu_files"},
			Permissions: []string{"files.read", "files.delete"},
			Fragments:   []string{"file-library"},
		},
	}
}

func (p *Provider) CompiledPersistence() ([]kernel.MigrationContribution, error) {
	return nil, nil // no tables owned by this module; operation-log events ride 0018
}

func (p *Provider) Register(ctx context.Context, reg kernel.Registrar) error {
	for _, route := range handler.FileLibraryRoutes(p.a, p.uploadDir, p.operations, ModuleID) {
		if err := reg.HTTP(route); err != nil {
			return err
		}
	}
	if err := reg.Schema(kernel.PageContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "file-library"},
		PageID:               "file-library",
		Resources:            []string{"files"},
		Actions:              []string{"list", "detail"},
		DataSource:           "/api/library/files",
		Owner:                ModuleID,
		Document:             filelibraryschema.SchemaDocuments()["file-library"],
	}); err != nil {
		return err
	}
	for _, permission := range []kernel.PermissionContribution{
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "files.read"}, Permission: "files.read", Resource: "files", Action: "read", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "files.delete"}, Permission: "files.delete", Resource: "files", Action: "delete", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
	} {
		if err := reg.Authorization(permission); err != nil {
			return err
		}
	}
	if err := reg.Navigation(kernel.NavigationContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "menu_files"},
		NodeID:               "menu_files",
		PageID:               "file-library",
		Order:                3,
		Label:                "File library",
		Visibility:           authsessiondata.PolicyAdmin,
		Permission:           "files.read",
		SystemDataVersion:    authsessiondata.SystemDataVersion,
	}); err != nil {
		return err
	}
	return reg.Manifest(kernel.FragmentContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "file-library"},
		FragmentID:           "file-library",
		ProtocolVersion:      "2.7",
		RequiredCapabilities: []string{"manifest", "navigation"},
		JSON:                 manifest.FragmentJSON,
	})
}
