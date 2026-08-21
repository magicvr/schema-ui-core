// Package migration owns the compiled-global async Job schema. core.jobs is a
// migration-only owner and is deliberately absent from runtime profiles.
package migration

import (
	"context"
	"fmt"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
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
	return []kernel.MigrationContribution{{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "async_jobs"},
		Version:              42,
		Name:                 "async_jobs",
		Checksum:             kernel.MigrationChecksum(jobsDDL, "0042:async-jobs:v1"),
		Apply:                migrateJobs,
		ApplyPostgres:        migrateJobsPG,
	}}
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
