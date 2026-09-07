// Package migration exposes the compiled biz.digital-offer migration owner
// (0070) to the global ledger registry (compiled/persistence.go).
package migration

import (
	"context"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// Provider exposes the compiled digital-offer migration owner (0070).
type Provider struct{}

func (Provider) Descriptor() kernel.Module {
	return kernel.Module{ID: ModuleID, Version: "1.0.0", KernelAPIRange: ">=2.0 <3.0"}
}

func (Provider) CompiledPersistence() ([]kernel.MigrationContribution, error) {
	return Descriptors(), nil
}

func (Provider) Register(ctx context.Context, reg kernel.Registrar) error { return nil }
