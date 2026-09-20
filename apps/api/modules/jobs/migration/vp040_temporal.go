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
// ModuleID core.jobs · version 76 · D-014 allocation / D-017 transform id.
package migration

import (
	"github.com/magicvr/schema-ui-core/apps/api/internal/temporalmigrate"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

const (
	vp040V76Name        = "vp040_temporal_jobs"
	vp040V76TransformID = "0076:vp040-temporal-jobs:v1"
)

// vp040V76Guards are the m0 preconditions on tables retired by earlier history
// (descriptor ledger §1: v73 asserts the retired `records` table is absent).
var vp040V76Guards = []temporalmigrate.Guard{}

// vp040V76Preflight is the m0 sentinel census (Root D-012 / D-015 policy).
var vp040V76Preflight = []temporalmigrate.Preflight{}

// vp040V76Verify is the m4 post-rebuild assertion set.
var vp040V76Verify = []temporalmigrate.Verify{
	{Table: "jobs", Columns: []string{"lease_expires_at", "created_at", "updated_at", "finished_at", "expires_at"}},
}

// vp040V76Rebuild is the ordered m1–m3 slice: constraint handling, table
// rebuild with the frozen conversion expressions, then every index.
var vp040V76Rebuild = []string{
	`ALTER TABLE "jobs" RENAME TO "jobs_old"`,
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
  lease_expires_at TEXT,
  result           TEXT,
  error_code       TEXT,
  error_message    TEXT,
  actor_id         TEXT NOT NULL CHECK (length(trim(actor_id)) > 0),
  correlation_id   TEXT NOT NULL CHECK (length(trim(correlation_id)) > 0),
  created_at       TEXT NOT NULL,
  updated_at       TEXT NOT NULL,
  finished_at      TEXT,
  expires_at       TEXT,
  CHECK (
    (status = 'queued' AND lease_owner IS NULL AND lease_expires_at IS NULL AND result IS NULL AND error_code IS NULL AND finished_at IS NULL AND expires_at IS NULL)
    OR (status = 'running' AND lease_owner IS NOT NULL AND lease_expires_at IS NOT NULL AND result IS NULL AND error_code IS NULL AND finished_at IS NULL AND expires_at IS NULL)
    OR (status = 'succeeded' AND lease_owner IS NULL AND lease_expires_at IS NULL AND result IS NOT NULL AND error_code IS NULL AND progress = 100 AND finished_at IS NOT NULL AND expires_at IS NOT NULL)
    OR (status = 'failed' AND lease_owner IS NULL AND lease_expires_at IS NULL AND result IS NULL AND error_code IS NOT NULL AND finished_at IS NOT NULL AND expires_at IS NULL)
    OR (status = 'cancelled' AND lease_owner IS NULL AND lease_expires_at IS NULL AND result IS NULL AND error_code IS NULL AND finished_at IS NOT NULL AND expires_at IS NULL)
    OR (status = 'expired' AND lease_owner IS NULL AND lease_expires_at IS NULL AND result IS NULL AND error_code IS NULL AND progress = 100 AND finished_at IS NOT NULL AND expires_at IS NOT NULL)
  )
)`,
	`INSERT INTO "jobs" ("id", "kind", "status", "payload", "progress", "cancel_requested", "attempt", "max_attempts", "lease_owner", "lease_version", "lease_expires_at", "result", "error_code", "error_message", "actor_id", "correlation_id", "created_at", "updated_at", "finished_at", "expires_at")
SELECT "id", "kind", "status", "payload", "progress", "cancel_requested", "attempt", "max_attempts", "lease_owner", "lease_version", CASE WHEN lease_expires_at IS NULL THEN NULL ELSE strftime('%Y-%m-%dT%H:%M:%S', CASE WHEN lease_expires_at >= 0 THEN lease_expires_at/1000 ELSE (lease_expires_at-999)/1000 END, 'unixepoch') || '.' || printf('%03d', (lease_expires_at%1000 + 1000) % 1000) || '000Z' END, "result", "error_code", "error_message", "actor_id", "correlation_id", strftime('%Y-%m-%dT%H:%M:%S', CASE WHEN created_at >= 0 THEN created_at/1000 ELSE (created_at-999)/1000 END, 'unixepoch') || '.' || printf('%03d', (created_at%1000 + 1000) % 1000) || '000Z', strftime('%Y-%m-%dT%H:%M:%S', CASE WHEN updated_at >= 0 THEN updated_at/1000 ELSE (updated_at-999)/1000 END, 'unixepoch') || '.' || printf('%03d', (updated_at%1000 + 1000) % 1000) || '000Z', CASE WHEN finished_at IS NULL THEN NULL ELSE strftime('%Y-%m-%dT%H:%M:%S', CASE WHEN finished_at >= 0 THEN finished_at/1000 ELSE (finished_at-999)/1000 END, 'unixepoch') || '.' || printf('%03d', (finished_at%1000 + 1000) % 1000) || '000Z' END, CASE WHEN expires_at IS NULL THEN NULL ELSE strftime('%Y-%m-%dT%H:%M:%S', CASE WHEN expires_at >= 0 THEN expires_at/1000 ELSE (expires_at-999)/1000 END, 'unixepoch') || '.' || printf('%03d', (expires_at%1000 + 1000) % 1000) || '000Z' END
FROM "jobs_old"`,
	`DROP TABLE "jobs_old"`,
	`CREATE INDEX idx_jobs_actor ON jobs(actor_id, kind, updated_at DESC)`,
	`CREATE INDEX idx_jobs_created_at ON jobs(created_at DESC, id DESC)`,
	`CREATE INDEX idx_jobs_expiry ON jobs(status, expires_at)`,
	`CREATE INDEX idx_jobs_runnable ON jobs(status, cancel_requested, lease_expires_at, created_at)`,
}

// vp040V76Statements returns the canonical checksum input (D-017: m0 → m1–m3 → m4).
func vp040V76Statements() []string {
	return temporalmigrate.OrderedWithGuards(vp040V76Guards, vp040V76Preflight, vp040V76Rebuild, vp040V76Verify)
}

// applyvp040V76 is the SQLite Apply body.
func applyvp040V76(tx kernel.Tx) error {
	if err := temporalmigrate.RunGuards(tx, vp040V76Guards, "vp040 v76 sqlite guards"); err != nil {
		return err
	}
	if err := temporalmigrate.RunPreflight(tx, vp040V76Preflight); err != nil {
		return err
	}
	if err := temporalmigrate.Exec(tx, vp040V76Rebuild, "vp040 v76 sqlite rebuild"); err != nil {
		return err
	}
	return temporalmigrate.RunVerify(tx, vp040V76Verify, "vp040 v76 sqlite verify")
}

// vp040V76Postgres is the explicit postgres m1–m2 slice (D-019 §6: no regex
// derivation; the millisecond family uses the integer-split expression).
var vp040V76Postgres = []string{
	`ALTER TABLE "jobs" ALTER COLUMN "lease_expires_at" TYPE timestamptz(6) USING (CASE WHEN "lease_expires_at" IS NULL THEN NULL ELSE TIMESTAMPTZ 'epoch' + (("lease_expires_at" - CASE WHEN "lease_expires_at" >= 0 THEN 0 ELSE 999 END) / 1000) * INTERVAL '1 second' + ((("lease_expires_at" % 1000) + 1000) % 1000) * INTERVAL '1 millisecond' END)`,
	`ALTER TABLE "jobs" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (TIMESTAMPTZ 'epoch' + (("created_at" - CASE WHEN "created_at" >= 0 THEN 0 ELSE 999 END) / 1000) * INTERVAL '1 second' + ((("created_at" % 1000) + 1000) % 1000) * INTERVAL '1 millisecond')`,
	`ALTER TABLE "jobs" ALTER COLUMN "updated_at" TYPE timestamptz(6) USING (TIMESTAMPTZ 'epoch' + (("updated_at" - CASE WHEN "updated_at" >= 0 THEN 0 ELSE 999 END) / 1000) * INTERVAL '1 second' + ((("updated_at" % 1000) + 1000) % 1000) * INTERVAL '1 millisecond')`,
	`ALTER TABLE "jobs" ALTER COLUMN "finished_at" TYPE timestamptz(6) USING (CASE WHEN "finished_at" IS NULL THEN NULL ELSE TIMESTAMPTZ 'epoch' + (("finished_at" - CASE WHEN "finished_at" >= 0 THEN 0 ELSE 999 END) / 1000) * INTERVAL '1 second' + ((("finished_at" % 1000) + 1000) % 1000) * INTERVAL '1 millisecond' END)`,
	`ALTER TABLE "jobs" ALTER COLUMN "expires_at" TYPE timestamptz(6) USING (CASE WHEN "expires_at" IS NULL THEN NULL ELSE TIMESTAMPTZ 'epoch' + (("expires_at" - CASE WHEN "expires_at" >= 0 THEN 0 ELSE 999 END) / 1000) * INTERVAL '1 second' + ((("expires_at" % 1000) + 1000) % 1000) * INTERVAL '1 millisecond' END)`,
}

// vp040V76PostgresVerify is the postgres m4 type/precision assertion set.
var vp040V76PostgresVerify = []temporalmigrate.PgVerify{
	{Table: "jobs", Column: "lease_expires_at", NonNull: false},
	{Table: "jobs", Column: "created_at", NonNull: true},
	{Table: "jobs", Column: "updated_at", NonNull: true},
	{Table: "jobs", Column: "finished_at", NonNull: false},
	{Table: "jobs", Column: "expires_at", NonNull: false},
}

// applyvp040V76Postgres is the postgres Apply body.
func applyvp040V76Postgres(tx kernel.Tx) error {
	if err := temporalmigrate.RunPostgresGuards(tx, vp040V76Guards, "vp040 v76 postgres guards"); err != nil {
		return err
	}
	if err := temporalmigrate.RunPreflight(tx, vp040V76Preflight); err != nil {
		return err
	}
	if err := temporalmigrate.Exec(tx, vp040V76Postgres, "vp040 v76 postgres convert"); err != nil {
		return err
	}
	return temporalmigrate.RunPostgresVerify(tx, vp040V76PostgresVerify, "vp040 v76 postgres verify")
}

// VP040TemporalDescriptor returns the workspace-040 R2 conversion descriptor
// (D-014 allocation, D-017 checksum convention).
func VP040TemporalDescriptor() kernel.MigrationContribution {
	return kernel.MigrationContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: vp040V76Name},
		Version:              76,
		Name:                 vp040V76Name,
		Checksum:             kernel.MigrationChecksum(vp040V76Statements(), vp040V76TransformID),
		Apply:                applyvp040V76,
		ApplyPostgres:        applyvp040V76Postgres,
	}
}
