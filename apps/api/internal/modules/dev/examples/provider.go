// Package examples provides the dev.examples demonstration module surface as a
// kernel.Provider. W1 (GOAL-002 / workspace-010): the protocol example pages and
// Examples navigation are a horizontal, dev-only demonstration surface — never
// enabled by mvp/admin defaults, explicitly enableable for development/dogfood
// via app.modules (config.yaml) or a dedicated profile (D-003 §3).
//
// dev.examples is exempt from the standard Admin six-face: it contributes only
// Schema documents + a Manifest fragment (no HTTP API / permissions / system-data
// navigation). Example pages may reuse real API endpoints (e.g. /api/users);
// there is no separate demo backend.
package examples

import (
	"context"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	examplesmanifest "github.com/magicvr/schema-ui-core/apps/api/internal/modules/dev/examples/manifest"
	examplesschema "github.com/magicvr/schema-ui-core/apps/api/internal/modules/dev/examples/schema"
)

const ModuleID = "dev.examples"

type Provider struct{}

func New() *Provider { return &Provider{} }

func (p *Provider) Descriptor() kernel.Module {
	return kernel.Module{
		ID:             ModuleID,
		Version:        "2.0.0",
		KernelAPIRange: ">=2.0 <3.0",
		DependsOn:      []string{"core.schema-render", "core.navigation-capability"},
		Contributions: kernel.ContributionKeys{
			Pages:     examplesschema.PageIDs(),
			Fragments: []string{"examples"},
		},
	}
}

func (p *Provider) CompiledPersistence() ([]kernel.MigrationContribution, error) {
	return nil, nil // examples own no persistence; schema migration is core.auth-session
}

func (p *Provider) Register(_ context.Context, reg kernel.Registrar) error {
	documents := examplesschema.SchemaDocuments()
	for _, pageID := range examplesschema.PageIDs() {
		if err := reg.Schema(kernel.PageContribution{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: pageID},
			PageID:               pageID,
			Owner:                ModuleID,
			Document:             documents[pageID],
		}); err != nil {
			return err
		}
	}
	return reg.Manifest(kernel.FragmentContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "examples"},
		FragmentID:           "examples",
		ProtocolVersion:      "2.7",
		RequiredCapabilities: []string{"manifest", "navigation"},
		JSON:                 examplesmanifest.FragmentJSON,
	})
}
