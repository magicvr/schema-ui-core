package migration_test

import (
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/modules/jobs/migration"
)

func TestDescriptor(t *testing.T) {
	descriptors := migration.Descriptors()
	if len(descriptors) != 2 {
		t.Fatalf("descriptor count = %d, want 2", len(descriptors))
	}
	d := descriptors[0]
	if d.ModuleID != migration.ModuleID || d.Version != 42 || d.Name != "async_jobs" || d.Checksum == "" || d.Apply == nil {
		t.Fatalf("descriptor = %+v", d)
	}
	// GOAL-003 R2 (D-001 §1): the management-list index is a NEW contribution.
	// The 0042 row must stay byte-identical (its checksum is in the ledger), so
	// adding an index can only be done as a later version.
	index := descriptors[1]
	if index.ModuleID != migration.ModuleID || index.Version != 72 ||
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
}
