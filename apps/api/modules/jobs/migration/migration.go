// Package migration owns the compiled-global async Job schema. core.jobs is a
// migration-only owner and is deliberately absent from runtime profiles.
package migration

import (
	"context"
	"fmt"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

const ModuleID = "core.jobs"

var jobsDDL = []string{
	`CREATE TABLE jobs (
  id               TEXT PRIMARY KEY,
  kind             TEXT NOT NULL CHECK (length(trim(kind)) > 0),
  status           TEXT NOT NULL CHECK (status IN ('queued','running','succeeded','failed','cancelled','expired')),
  payload          TEXT NOT NULL DEFAULT '{}',
  progress         INTEGER NOT NULL DEFAULT 0 CHECK (progress BETWEEN 0 AND 100),
  cancel_requested INTEGER NOT NULL DEFAULT 0 CHECK (cancel_requested IN (0,1)),
  attempt          INTEGER NOT NULL DEFAULT 0 CHECK (attempt >= 0),
  max_attempts     INTEGER NOT NULL DEFAULT 3 CHECK (max_attempts > 0 AND attempt <= max_attempts),
  lease_owner      TEXT,
  lease_version    INTEGER NOT NULL DEFAULT 0 CHECK (lease_version >= 0),
  lease_expires_at INTEGER,
  result           TEXT,
  error_code       TEXT,
  error_message    TEXT,
  actor_id         TEXT NOT NULL CHECK (length(trim(actor_id)) > 0),
  correlation_id   TEXT NOT NULL CHECK (length(trim(correlation_id)) > 0),
  created_at       INTEGER NOT NULL,
  updated_at       INTEGER NOT NULL,
  finished_at      INTEGER,
  expires_at       INTEGER,
  CHECK (
    (status = 'queued' AND lease_owner IS NULL AND lease_expires_at IS NULL AND result IS NULL AND error_code IS NULL AND finished_at IS NULL AND expires_at IS NULL)
    OR (status = 'running' AND lease_owner IS NOT NULL AND lease_expires_at IS NOT NULL AND result IS NULL AND error_code IS NULL AND finished_at IS NULL AND expires_at IS NULL)
    OR (status = 'succeeded' AND lease_owner IS NULL AND lease_expires_at IS NULL AND result IS NOT NULL AND error_code IS NULL AND progress = 100 AND finished_at IS NOT NULL AND expires_at IS NOT NULL)
    OR (status = 'failed' AND lease_owner IS NULL AND lease_expires_at IS NULL AND result IS NULL AND error_code IS NOT NULL AND finished_at IS NOT NULL AND expires_at IS NULL)
    OR (status = 'cancelled' AND lease_owner IS NULL AND lease_expires_at IS NULL AND result IS NULL AND error_code IS NULL AND finished_at IS NOT NULL AND expires_at IS NULL)
    OR (status = 'expired' AND lease_owner IS NULL AND lease_expires_at IS NULL AND result IS NULL AND error_code IS NULL AND progress = 100 AND finished_at IS NOT NULL AND expires_at IS NOT NULL)
  )
)`,
	`CREATE INDEX idx_jobs_runnable ON jobs(status, cancel_requested, lease_expires_at, created_at)`,
	`CREATE INDEX idx_jobs_actor ON jobs(actor_id, kind, updated_at DESC)`,
	`CREATE INDEX idx_jobs_expiry ON jobs(status, expires_at)`,
}

// jobsPGDDL is the postgres variant of jobsDDL: Unix time columns
// (lease_expires_at / created_at / updated_at / finished_at / expires_at) are
// BIGINT (R1 v1.4 §3).
var jobsPGDDL = []string{
	`CREATE TABLE jobs (
  id               TEXT PRIMARY KEY,
  kind             TEXT NOT NULL CHECK (length(trim(kind)) > 0),
  status           TEXT NOT NULL CHECK (status IN ('queued','running','succeeded','failed','cancelled','expired')),
  payload          TEXT NOT NULL DEFAULT '{}',
  progress         INTEGER NOT NULL DEFAULT 0 CHECK (progress BETWEEN 0 AND 100),
  cancel_requested INTEGER NOT NULL DEFAULT 0 CHECK (cancel_requested IN (0,1)),
  attempt          INTEGER NOT NULL DEFAULT 0 CHECK (attempt >= 0),
  max_attempts     INTEGER NOT NULL DEFAULT 3 CHECK (max_attempts > 0 AND attempt <= max_attempts),
  lease_owner      TEXT,
  lease_version    INTEGER NOT NULL DEFAULT 0 CHECK (lease_version >= 0),
  lease_expires_at BIGINT,
  result           TEXT,
  error_code       TEXT,
  error_message    TEXT,
  actor_id         TEXT NOT NULL CHECK (length(trim(actor_id)) > 0),
  correlation_id   TEXT NOT NULL CHECK (length(trim(correlation_id)) > 0),
  created_at       BIGINT NOT NULL,
  updated_at       BIGINT NOT NULL,
  finished_at      BIGINT,
  expires_at       BIGINT,
  CHECK (
    (status = 'queued' AND lease_owner IS NULL AND lease_expires_at IS NULL AND result IS NULL AND error_code IS NULL AND finished_at IS NULL AND expires_at IS NULL)
    OR (status = 'running' AND lease_owner IS NOT NULL AND lease_expires_at IS NOT NULL AND result IS NULL AND error_code IS NULL AND finished_at IS NULL AND expires_at IS NULL)
    OR (status = 'succeeded' AND lease_owner IS NULL AND lease_expires_at IS NULL AND result IS NOT NULL AND error_code IS NULL AND progress = 100 AND finished_at IS NOT NULL AND expires_at IS NOT NULL)
    OR (status = 'failed' AND lease_owner IS NULL AND lease_expires_at IS NULL AND result IS NULL AND error_code IS NOT NULL AND finished_at IS NOT NULL AND expires_at IS NULL)
    OR (status = 'cancelled' AND lease_owner IS NULL AND lease_expires_at IS NULL AND result IS NULL AND error_code IS NULL AND finished_at IS NOT NULL AND expires_at IS NULL)
    OR (status = 'expired' AND lease_owner IS NULL AND lease_expires_at IS NULL AND result IS NULL AND error_code IS NULL AND progress = 100 AND finished_at IS NOT NULL AND expires_at IS NOT NULL)
  )
)`,
	`CREATE INDEX idx_jobs_runnable ON jobs(status, cancel_requested, lease_expires_at, created_at)`,
	`CREATE INDEX idx_jobs_actor ON jobs(actor_id, kind, updated_at DESC)`,
	`CREATE INDEX idx_jobs_expiry ON jobs(status, expires_at)`,
}

func Descriptors() []kernel.MigrationContribution {
	return []kernel.MigrationContribution{
		// workspace-040 R2 (GOAL-003 M2): v73–v87 timestamp conversions.
		// The kernel orders the compiled catalog by Version; source order
		// is not significant.
		VP040TemporalDescriptor(),
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "async_jobs"},
			Version:              42,
			Name:                 "async_jobs",
			Checksum:             kernel.MigrationChecksum(jobsDDL, "0042:async-jobs:v1"),
			Apply:                migrateJobs,
			ApplyPostgres:        migrateJobsPG,
		},
		{
			// GOAL-003 R2 (D-001 §1): management-scope job list index. The
			// runtime state-machine indexes cannot serve a cross-actor
			// `ORDER BY created_at DESC`, so the admin list gets its own
			// index — the same shape admin.activity uses for
			// operation_log(created_at DESC). Index-only DDL is portable
			// (no time-column type difference), so ApplyPostgres is omitted.
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "jobs_management_indexes"},
			Version:              72,
			Name:                 "jobs_management_indexes",
			Checksum:             kernel.MigrationChecksum(jobsManagementIndexDDL, "0072:jobs-management-indexes:v1"),
			Apply:                migrateJobsManagementIndexes,
		},
	}
}

// jobsManagementIndexDDL adds the management-list index. The 0042 table and
// its indexes are untouched: that contribution's checksum is frozen in the
// migration ledger, so a new contribution is the only legal way to add DDL.
//
// The index carries the `id` tiebreak column because the admin list pages with
// `ORDER BY created_at DESC, id DESC` — `created_at` is millisecond precision,
// so without a deterministic tiebreak rows sharing a timestamp could repeat or
// vanish across pages (the same reason operationlog and wallet add `, id DESC`).
// EXPLAIN QUERY PLAN confirms the tiebreak column is what lets the index order
// the result directly instead of falling back to a temp B-tree.
// `IF NOT EXISTS` matches the pure-index precedent (admin.settings 0063).
var jobsManagementIndexDDL = []string{
	`CREATE INDEX IF NOT EXISTS idx_jobs_created_at ON jobs(created_at DESC, id DESC)`,
}

func migrateJobsManagementIndexes(tx kernel.Tx) error {
	for _, stmt := range jobsManagementIndexDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("create jobs management index: %w", err)
		}
	}
	return nil
}

func migrateJobs(tx kernel.Tx) error {
	for _, stmt := range jobsDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("create async jobs: %w", err)
		}
	}
	return nil
}

func migrateJobsPG(tx kernel.Tx) error {
	for _, stmt := range jobsPGDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("create async jobs (postgres): %w", err)
		}
	}
	return nil
}
