package migration

import (
	"context"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// ModuleID is the owning module for the telegram_config migration.
const ModuleID = "channel.telegram"

// Provider owns the compiled-global migration history for channel.telegram.
type Provider struct{}

func (Provider) Descriptor() kernel.Module {
	return kernel.Module{ID: ModuleID, Version: "2.0.0", KernelAPIRange: ">=2.0 <3.0"}
}

func (Provider) CompiledPersistence() ([]kernel.MigrationContribution, error) {
	return Descriptors()
}

func (Provider) Register(ctx context.Context, reg kernel.Registrar) error {
	return nil
}
