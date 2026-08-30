// Package schemarender provides the core.schema-render module contribution:
// the schema / validation capability shell. W1 (GOAL-002 / workspace-010) moved
// the 8 example page documents out to the optional dev.examples module, so this
// module provides the capability but no page contributions of its own.
package schemarender

import (
	"context"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

const ModuleID = "core.schema-render"

type Provider struct{}

func New() *Provider { return &Provider{} }

func (p *Provider) Descriptor() kernel.Module {
	return kernel.Module{
		ID:             ModuleID,
		Version:        "2.0.0",
		KernelAPIRange: ">=2.0 <3.0",
		DependsOn:      []string{"core.server-registration"},
		Provides:       []kernel.Capability{kernel.CapabilitySchema, kernel.CapabilityValidation},
	}
}

func (p *Provider) CompiledPersistence() ([]kernel.MigrationContribution, error) {
	return nil, nil
}

// Register contributes no surfaces: the schema capability is provided by the
// module graph, and example/page documents are owned by their own modules.
func (p *Provider) Register(context.Context, kernel.Registrar) error {
	return nil
}
