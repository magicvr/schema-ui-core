package schemarender

import (
	"context"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// W1 (GOAL-002 / workspace-010): core.schema-render is now a capability-only
// shell. The 8 example page documents moved to the optional dev.examples module
// (see dev/examples/provider_test.go). Assert it publishes no page surfaces.
func TestProviderPublishesNoPageDocuments(t *testing.T) {
	provider := New()
	descriptor := provider.Descriptor()
	plan := kernel.Plan{
		Modules:      []kernel.Module{descriptor},
		Capabilities: []kernel.Capability{kernel.CapabilitySchema, kernel.CapabilityValidation},
	}
	set, err := kernel.RegisterContributions(context.Background(), plan, []kernel.Provider{provider})
	if err != nil {
		t.Fatal(err)
	}
	if len(set.Pages) != 0 {
		t.Fatalf("core.schema-render must not own page documents after W1 split, got %d: %v", len(set.Pages), set.Pages)
	}
	hasSchema := false
	for _, capability := range descriptor.Provides {
		if capability == kernel.CapabilitySchema {
			hasSchema = true
		}
	}
	if !hasSchema {
		t.Fatal("core.schema-render must still provide the schema capability")
	}
}
