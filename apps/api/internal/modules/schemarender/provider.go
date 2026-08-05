// Package schemarender provides the core.schema-render module contribution.
package schemarender

import (
	"context"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	coreschema "github.com/magicvr/schema-ui-core/apps/api/internal/modules/schemarender/schema"
)

const ModuleID = coreschema.ModuleID

type Provider struct{}

func New() *Provider { return &Provider{} }

func (p *Provider) Descriptor() kernel.Module {
	return kernel.Module{
		ID:             ModuleID,
		Version:        "2.0.0",
		KernelAPIRange: ">=2.0 <3.0",
		DependsOn:      []string{"core.server-registration"},
		Provides:       []kernel.Capability{kernel.CapabilitySchema, kernel.CapabilityValidation},
		Contributions: kernel.ContributionKeys{
			Pages: coreschema.PageIDs(),
		},
	}
}

func (p *Provider) CompiledPersistence() ([]kernel.MigrationContribution, error) {
	return nil, nil
}

func (p *Provider) Register(_ context.Context, reg kernel.Registrar) error {
	documents := coreschema.SchemaDocuments()
	for _, pageID := range coreschema.PageIDs() {
		if err := reg.Schema(kernel.PageContribution{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: pageID},
			PageID:               pageID,
			Owner:                ModuleID,
			Document:             documents[pageID],
		}); err != nil {
			return err
		}
	}
	return nil
}
