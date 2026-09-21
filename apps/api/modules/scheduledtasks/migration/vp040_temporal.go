// Code generated for workspace-040 R2 (GOAL-003 checkpoints B/C) from the
// v72-applied schema. DO NOT EDIT BY HAND.
//
// Every CREATE TABLE body below is the live sqlite_master text of the frozen
// v1–v72 history (live PRAGMA table_info column order) with only the frozen
// per-column edits of r1-c2-per-column-conversion-contract-v1.0-fc.md applied:
// the converted column becomes TEXT and, for the D0 / voucher columns, loses
// NOT NULL and DEFAULT 0. No other token changed.
//
// Checksum input (D-017): m0 preflight → m1–m3 rebuild → m4 verification. The
// postgres variant is explicit DDL and is not hashed (D-019 §6).
//
// ModuleID admin.scheduled-tasks · version 83 · D-014 allocation / D-017 transform id.
package migration

import (
	"github.com/magicvr/schema-ui-core/apps/api/internal/temporalmigrate"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

const (
	vp040V83Name        = "vp040_temporal_scheduled_tasks"
	vp040V83TransformID = "0083:vp040-temporal-scheduled-tasks:v1"
)

// vp040V83Guards are the m0 preconditions on tables retired by earlier history
// (descriptor ledger §1: v73 asserts the retired `records` table is absent).
var vp040V83Guards = []temporalmigrate.Guard{}

// vp040V83Preflight is the m0 sentinel census (Root D-012 / D-015 policy).
var vp040V83Preflight = []temporalmigrate.Preflight{}

// vp040V83Verify is the m4 post-rebuild assertion set.
var vp040V83Verify = []temporalmigrate.Verify{
	{Table: "scheduled_tasks", Columns: []string{"created_at", "updated_at"}, Children: []string{"task_runs"}},
	{Table: "task_runs", Columns: []string{"started_at", "finished_at", "created_at"}},
}

// vp040V83Rebuild is the ordered m1–m3 slice: constraint handling, table
// rebuild with the frozen conversion expressions, then every index.
var vp040V83Rebuild = []string{
	`CREATE TEMP TABLE "task_runs_bak" AS SELECT * FROM "task_runs"`,
	`DROP TABLE "task_runs"`,
	`ALTER TABLE "scheduled_tasks" RENAME TO "scheduled_tasks_old"`,
	`CREATE TABLE scheduled_tasks (
  id          TEXT PRIMARY KEY,
  key         TEXT NOT NULL UNIQUE,
  cron        TEXT NOT NULL,
  name        TEXT NOT NULL,
  enabled     INTEGER NOT NULL DEFAULT 1,
  description TEXT,
  handler     TEXT NOT NULL DEFAULT 'system.noop',
  created_at  TEXT NOT NULL,
  updated_at  TEXT NOT NULL
)`,
	`INSERT INTO "scheduled_tasks" ("id", "key", "cron", "name", "enabled", "description", "handler", "created_at", "updated_at")
SELECT "id", "key", "cron", "name", "enabled", "description", "handler", strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z', strftime('%Y-%m-%dT%H:%M:%S', updated_at, 'unixepoch') || '.000000Z'
FROM "scheduled_tasks_old"`,
	`DROP TABLE "scheduled_tasks_old"`,
	`CREATE TABLE task_runs (
  id          TEXT PRIMARY KEY,
  task_id     TEXT NOT NULL REFERENCES scheduled_tasks(id) ON DELETE CASCADE,
  status      TEXT NOT NULL CHECK (status IN ('ran','failed')),
  started_at  TEXT NOT NULL,
  finished_at TEXT,
  detail      TEXT,
  created_at  TEXT NOT NULL
)`,
	`INSERT INTO "task_runs" ("id", "task_id", "status", "started_at", "finished_at", "detail", "created_at")
SELECT "id", "task_id", "status", strftime('%Y-%m-%dT%H:%M:%S', started_at, 'unixepoch') || '.000000Z', CASE WHEN finished_at IS NULL THEN NULL ELSE strftime('%Y-%m-%dT%H:%M:%S', finished_at, 'unixepoch') || '.000000Z' END, "detail", strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z'
FROM "task_runs_bak"`,
	`DROP TABLE "task_runs_bak"`,
	`CREATE INDEX idx_task_runs_task_started ON task_runs(task_id, started_at DESC)`,
}

// vp040V83Statements returns the canonical checksum input (D-017: m0 → m1–m3 → m4).
func vp040V83Statements() []string {
	return temporalmigrate.OrderedWithGuards(vp040V83Guards, vp040V83Preflight, vp040V83Rebuild, vp040V83Verify)
}

// applyvp040V83 is the SQLite Apply body.
func applyvp040V83(tx kernel.Tx) error {
	if err := temporalmigrate.RunGuards(tx, vp040V83Guards, "vp040 v83 sqlite guards"); err != nil {
		return err
	}
	if err := temporalmigrate.RunPreflight(tx, vp040V83Preflight); err != nil {
		return err
	}
	if err := temporalmigrate.Exec(tx, vp040V83Rebuild, "vp040 v83 sqlite rebuild"); err != nil {
		return err
	}
	return temporalmigrate.RunVerify(tx, vp040V83Verify, "vp040 v83 sqlite verify")
}

// vp040V83Postgres is the explicit postgres m1–m2 slice (D-019 §6: no regex
// derivation; the millisecond family uses the integer-split expression).
var vp040V83Postgres = []string{
	`ALTER TABLE "scheduled_tasks" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)))`,
	`ALTER TABLE "scheduled_tasks" ALTER COLUMN "updated_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("updated_at"::double precision)))`,
	`ALTER TABLE "task_runs" ALTER COLUMN "started_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("started_at"::double precision)))`,
	`ALTER TABLE "task_runs" ALTER COLUMN "finished_at" TYPE timestamptz(6) USING (CASE WHEN "finished_at" IS NULL THEN NULL ELSE date_trunc('microseconds', to_timestamp("finished_at"::double precision)) END)`,
	`ALTER TABLE "task_runs" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)))`,
}

// vp040V83PostgresVerify is the postgres m4 type/precision assertion set.
var vp040V83PostgresVerify = []temporalmigrate.PgVerify{
	{Table: "scheduled_tasks", Column: "created_at", NonNull: true},
	{Table: "scheduled_tasks", Column: "updated_at", NonNull: true},
	{Table: "task_runs", Column: "started_at", NonNull: true},
	{Table: "task_runs", Column: "finished_at", NonNull: false},
	{Table: "task_runs", Column: "created_at", NonNull: true},
}

// applyvp040V83Postgres is the postgres Apply body.
func applyvp040V83Postgres(tx kernel.Tx) error {
	if err := temporalmigrate.RunPostgresGuards(tx, vp040V83Guards, "vp040 v83 postgres guards"); err != nil {
		return err
	}
	if err := temporalmigrate.RunPreflight(tx, vp040V83Preflight); err != nil {
		return err
	}
	if err := temporalmigrate.Exec(tx, vp040V83Postgres, "vp040 v83 postgres convert"); err != nil {
		return err
	}
	return temporalmigrate.RunPostgresVerify(tx, vp040V83PostgresVerify, "vp040 v83 postgres verify")
}

// VP040TemporalDescriptor returns the workspace-040 R2 conversion descriptor
// (D-014 allocation, D-017 checksum convention).
func VP040TemporalDescriptor() kernel.MigrationContribution {
	return kernel.MigrationContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: vp040V83Name},
		Version:              83,
		Name:                 vp040V83Name,
		Checksum:             kernel.MigrationChecksum(vp040V83Statements(), vp040V83TransformID),
		Apply:                applyvp040V83,
		ApplyPostgres:        applyvp040V83Postgres,
	}
}
