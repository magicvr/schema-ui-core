// Package datadictionary provides the admin.data-dictionary module surface as
// a kernel.Provider (S-01 · GOAL-008 D-002): two-level dictionary types and
// entries with schema-driven pages, dictionary.read / dictionary.write
// permission keys, menu_dictionary navigation and dictionary.* audit events.
package datadictionary

import (
	"context"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/handler"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	authsessiondata "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession/systemdata"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/datadictionary/manifest"
	datadictionarystore "github.com/magicvr/schema-ui-core/apps/api/internal/modules/datadictionary/store"
	datadictionaryschema "github.com/magicvr/schema-ui-core/apps/api/internal/modules/datadictionary/schema"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
)

// ModuleID is the stable admin.data-dictionary module identifier.
const ModuleID = "admin.data-dictionary"

// Provider implements kernel.Provider for admin.data-dictionary.
type Provider struct {
	a          *auth.Authenticator
	repository *datadictionarystore.Repository
	operations operationlog.Recorder
	trash      handler.TrashRecorder
}

// New constructs the dictionary provider with framework-agnostic
// dependencies. trash (S-12 · GOAL-012 D-002 §2), when non-nil, opts the
// dictionary resources into recycle-bin snapshot recording on delete.
func New(a *auth.Authenticator, repository *datadictionarystore.Repository, operations operationlog.Recorder, trash ...handler.TrashRecorder) *Provider {
	var recorder handler.TrashRecorder
	if len(trash) > 0 {
		recorder = trash[0]
	}
	return &Provider{a: a, repository: repository, operations: operations, trash: recorder}
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
				"GET /api/data-dictionary/types", "GET /api/data-dictionary/types/{id}",
				"POST /api/data-dictionary/types", "PATCH /api/data-dictionary/types/{id}",
				"DELETE /api/data-dictionary/types/{id}", "POST /api/data-dictionary/types/batch-delete",
				"GET /api/data-dictionary/entries", "GET /api/data-dictionary/entries/{id}",
				"POST /api/data-dictionary/entries", "PATCH /api/data-dictionary/entries/{id}",
				"DELETE /api/data-dictionary/entries/{id}", "POST /api/data-dictionary/entries/batch-delete",
			},
			Pages:       []string{"data-dictionary", "dictionary-entries"},
			Navigation:  []string{"menu_dictionary"},
			Permissions: []string{"dictionary.read", "dictionary.write"},
			Fragments:   []string{"data-dictionary"},
		},
	}
}

func (p *Provider) CompiledPersistence() ([]kernel.MigrationContribution, error) {
	return nil, nil // tables are owned by the datadictionary/migration provider (0019)
}

func (p *Provider) Register(ctx context.Context, reg kernel.Registrar) error {
	for _, route := range handler.DictionaryRoutes(p.a, p.repository, p.operations, ModuleID, p.trash) {
		if err := reg.HTTP(route); err != nil {
			return err
		}
	}
	for _, pageID := range []string{"data-dictionary", "dictionary-entries"} {
		if err := reg.Schema(kernel.PageContribution{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: pageID},
			PageID:               pageID,
			Resources:            []string{"dict-types", "dict-entries"},
			Actions:              []string{"list", "create", "detail", "update", "delete"},
			DataSource:           "/api/data-dictionary/types",
			Owner:                ModuleID,
			Document:             datadictionaryschema.SchemaDocuments()[pageID],
		}); err != nil {
			return err
		}
	}
	for _, permission := range []kernel.PermissionContribution{
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "dictionary.read"}, Permission: "dictionary.read", Resource: "dictionary", Action: "read", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "dictionary.write"}, Permission: "dictionary.write", Resource: "dictionary", Action: "write", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
	} {
		if err := reg.Authorization(permission); err != nil {
			return err
		}
	}
	if err := reg.Navigation(kernel.NavigationContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "menu_dictionary"},
		NodeID:               "menu_dictionary",
		PageID:               "data-dictionary",
		Order:                4,
		Label:                "Data dictionary",
		Visibility:           authsessiondata.PolicyAdmin,
		Permission:           "dictionary.read",
		SystemDataVersion:    authsessiondata.SystemDataVersion,
	}); err != nil {
		return err
	}
	return reg.Manifest(kernel.FragmentContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "data-dictionary"},
		FragmentID:           "data-dictionary",
		ProtocolVersion:      "2.7",
		RequiredCapabilities: []string{"manifest", "navigation"},
		JSON:                 manifest.FragmentJSON,
	})
}
