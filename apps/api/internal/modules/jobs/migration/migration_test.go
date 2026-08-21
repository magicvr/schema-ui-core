package migration_test

import (
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/jobs/migration"
)

func TestDescriptor(t *testing.T) {
	descriptors := migration.Descriptors()
	if len(descriptors) != 1 {
		t.Fatalf("descriptor count = %d, want 1", len(descriptors))
	}
	d := descriptors[0]
	if d.ModuleID != migration.ModuleID || d.Version != 42 || d.Name != "async_jobs" || d.Checksum == "" || d.Apply == nil {
		t.Fatalf("descriptor = %+v", d)
	}
}
