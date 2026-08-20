// Package migration owns the admin.scheduled-tasks schema (S-04 · GOAL-010
// D-002 `1): scheduled_tasks definitions plus task_runs execution history.
package migration

import (
	"context"
	"fmt"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
)

// ModuleID is the S-04 tasks module owner.
const ModuleID = "admin.scheduled-tasks"

// tasksDDL (0021): task definitions and run history. Deleting a task cascades
// to its run rows.
var tasksDDL = []string{
	`CREATE TABLE scheduled_tasks (
  id          TEXT PRIMARY KEY,
  key         TEXT NOT NULL UNIQUE,
  cron        TEXT NOT NULL,
  name        TEXT NOT NULL,
  enabled     INTEGER NOT NULL DEFAULT 1,
  description TEXT,
  handler     TEXT NOT NULL DEFAULT 'system.noop',
  created_at  INTEGER NOT NULL,
  updated_at  INTEGER NOT NULL
)`,
	`CREATE TABLE task_runs (
  id          TEXT PRIMARY KEY,
  task_id     TEXT NOT NULL REFERENCES scheduled_tasks(id) ON DELETE CASCADE,
  status      TEXT NOT NULL CHECK (status IN ('ran','failed')),
  started_at  INTEGER NOT NULL,
  finished_at INTEGER,
  detail      TEXT,
  created_at  INTEGER NOT NULL
)`,
	`CREATE INDEX idx_task_runs_task_started ON task_runs(task_id, started_at DESC)`,
}

// tasksPGDDL is the postgres variant of tasksDDL: Unix time columns
// (created_at / updated_at / started_at / finished_at) are BIGINT (R1 v1.4 §3).
var tasksPGDDL = []string{
	`CREATE TABLE scheduled_tasks (
  id          TEXT PRIMARY KEY,
  key         TEXT NOT NULL UNIQUE,
  cron        TEXT NOT NULL,
  name        TEXT NOT NULL,
  enabled     INTEGER NOT NULL DEFAULT 1,
  description TEXT,
  handler     TEXT NOT NULL DEFAULT 'system.noop',
  created_at  BIGINT NOT NULL,
  updated_at  BIGINT NOT NULL
)`,
	`CREATE TABLE task_runs (
  id          TEXT PRIMARY KEY,
  task_id     TEXT NOT NULL REFERENCES scheduled_tasks(id) ON DELETE CASCADE,
  status      TEXT NOT NULL CHECK (status IN ('ran','failed')),
  started_at  BIGINT NOT NULL,
  finished_at BIGINT,
  detail      TEXT,
  created_at  BIGINT NOT NULL
)`,
	`CREATE INDEX idx_task_runs_task_started ON task_runs(task_id, started_at DESC)`,
}

// Descriptors returns the immutable 0021 tasks history.
func Descriptors() []kernel.MigrationContribution {
	return []kernel.MigrationContribution{
		{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "scheduled_tasks"},
			Version:              21,
			Name:                 "scheduled_tasks",
			Checksum:             kernel.MigrationChecksum(tasksDDL, "0021:scheduled-tasks:v1"),
			Apply:                migrateTasks,
			ApplyPostgres:        migrateTasksPG,
		},
	}
}

func migrateTasks(tx kernel.Tx) error {
	for _, stmt := range tasksDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("create scheduled task tables: %w", err)
		}
	}
	return nil
}

func migrateTasksPG(tx kernel.Tx) error {
	for _, stmt := range tasksPGDDL {
		if _, err := tx.Exec(context.Background(), stmt); err != nil {
			return fmt.Errorf("create scheduled task tables (postgres): %w", err)
		}
	}
	return nil
}
