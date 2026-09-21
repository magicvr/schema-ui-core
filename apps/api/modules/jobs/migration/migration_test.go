package migration_test

import (
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/modules/jobs/migration"
)

func TestDescriptor(t *testing.T) {
	descriptors := migration.Descriptors()
	if len(descriptors) != 3 {
		t.Fatalf("descriptor count = %d, want 3", len(descriptors))
	}
	// The slice order is not contractual (the kernel sorts by Version), so look
	// the descriptors up by version.
	byVersion := make(map[int]kernel.MigrationContribution, len(descriptors))
	for _, d := range descriptors {
		byVersion[d.Version] = d
	}

	d, ok := byVersion[42]
	if !ok || d.ModuleID != migration.ModuleID || d.Name != "async_jobs" || d.Checksum == "" || d.Apply == nil {
		t.Fatalf("descriptor v42 = %+v", d)
	}
	// GOAL-003 R2 (D-001 §1): the management-list index is a NEW contribution.
	// The 0042 row must stay byte-identical (its checksum is in the ledger), so
	// adding an index can only be done as a later version.
	index, ok := byVersion[72]
	if !ok || index.ModuleID != migration.ModuleID ||
		index.Name != "jobs_management_indexes" || index.Key != "jobs_management_indexes" ||
		index.Checksum == "" || index.Apply == nil {
		t.Fatalf("index descriptor = %+v", index)
	}
	if index.Tombstone {
		t.Fatalf("index descriptor must not be a tombstone")
	}
	// Index-only DDL has no time-column type difference, so the canonical Apply
	// is portable to postgres and ApplyPostgres stays nil.
	if index.ApplyPostgres != nil {
		t.Fatalf("index descriptor must not carry ApplyPostgres (portable DDL)")
	}

	// workspace-040 R2 (GOAL-003 M2): the v76 timestamp conversion descriptor.
	conversion, ok := byVersion[76]
	if !ok || conversion.ModuleID != migration.ModuleID ||
		conversion.Name != "vp040_temporal_jobs" || conversion.Key != "vp040_temporal_jobs" ||
		conversion.Checksum == "" || conversion.Apply == nil {
		t.Fatalf("conversion descriptor = %+v", conversion)
	}
	if conversion.ApplyPostgres == nil {
		t.Fatalf("conversion descriptor must carry the explicit postgres apply")
	}
}
