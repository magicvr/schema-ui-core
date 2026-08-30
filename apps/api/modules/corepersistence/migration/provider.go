package migration

import (
	"context"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// Provider exposes retired migration history to the compiled-global catalog.
type Provider struct{}

func (Provider) Descriptor() kernel.Module {
	return kernel.Module{ID: ModuleID, Version: "2.0.0", KernelAPIRange: ">=2.0 <3.0"}
}

func (Provider) CompiledPersistence() ([]kernel.MigrationContribution, error) {
	return Descriptors(), nil
}

func (Provider) Register(context.Context, kernel.Registrar) error { return nil }
