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
// ModuleID admin.recycle-bin · version 82 · D-014 allocation / D-017 transform id.
package migration

import (
	"github.com/magicvr/schema-ui-core/apps/api/internal/temporalmigrate"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

const (
	vp040V82Name        = "vp040_temporal_recycle"
	vp040V82TransformID = "0082:vp040-temporal-recycle:v1"
)

// vp040V82Preflight is the m0 sentinel census (Root D-012 / D-015 policy).
var vp040V82Preflight = []temporalmigrate.Preflight{}

// vp040V82Verify is the m4 post-rebuild assertion set.
var vp040V82Verify = []temporalmigrate.Verify{
	{Table: "recycle_items", Columns: []string{"deleted_at", "restored_at"}},
}

// vp040V82Rebuild is the ordered m1–m3 slice: constraint handling, table
// rebuild with the frozen conversion expressions, then every index.
var vp040V82Rebuild = []string{
	`ALTER TABLE "recycle_items" RENAME TO "recycle_items_old"`,
	`CREATE TABLE recycle_items (
  id          TEXT PRIMARY KEY,
  resource    TEXT NOT NULL,
  resource_id TEXT NOT NULL,
  payload     TEXT NOT NULL,
  actor_id    TEXT NOT NULL,
  actor_name  TEXT NOT NULL,
  deleted_at  TEXT NOT NULL,
  restored_at TEXT
)`,
	`INSERT INTO "recycle_items" ("id", "resource", "resource_id", "payload", "actor_id", "actor_name", "deleted_at", "restored_at")
SELECT "id", "resource", "resource_id", "payload", "actor_id", "actor_name", strftime('%Y-%m-%dT%H:%M:%S', deleted_at, 'unixepoch') || '.000000Z', CASE WHEN restored_at IS NULL THEN NULL ELSE strftime('%Y-%m-%dT%H:%M:%S', restored_at, 'unixepoch') || '.000000Z' END
FROM "recycle_items_old"`,
	`DROP TABLE "recycle_items_old"`,
	`CREATE INDEX idx_recycle_items_deleted_at ON recycle_items(deleted_at DESC)`,
	`CREATE UNIQUE INDEX idx_recycle_items_active ON recycle_items(resource, resource_id) WHERE restored_at IS NULL`,
}

// vp040V82Statements returns the canonical checksum input (D-017: m0 → m1–m3 → m4).
func vp040V82Statements() []string {
	return temporalmigrate.Ordered(vp040V82Preflight, vp040V82Rebuild, vp040V82Verify)
}

// applyvp040V82 is the SQLite Apply body.
func applyvp040V82(tx kernel.Tx) error {
	if err := temporalmigrate.RunPreflight(tx, vp040V82Preflight); err != nil {
		return err
	}
	if err := temporalmigrate.Exec(tx, vp040V82Rebuild, "vp040 v82 sqlite rebuild"); err != nil {
		return err
	}
	return temporalmigrate.RunVerify(tx, vp040V82Verify, "vp040 v82 sqlite verify")
}

// vp040V82Postgres is the explicit postgres m1–m2 slice (D-019 §6: no regex
// derivation; the millisecond family uses the integer-split expression).
var vp040V82Postgres = []string{
	`ALTER TABLE "recycle_items" ALTER COLUMN "deleted_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("deleted_at"::double precision)))`,
	`ALTER TABLE "recycle_items" ALTER COLUMN "restored_at" TYPE timestamptz(6) USING (CASE WHEN "restored_at" IS NULL THEN NULL ELSE date_trunc('microseconds', to_timestamp("restored_at"::double precision)) END)`,
}

// vp040V82PostgresVerify is the postgres m4 type/precision assertion set.
var vp040V82PostgresVerify = []temporalmigrate.PgVerify{
	{Table: "recycle_items", Column: "deleted_at", NonNull: true},
	{Table: "recycle_items", Column: "restored_at", NonNull: false},
}

// applyvp040V82Postgres is the postgres Apply body.
func applyvp040V82Postgres(tx kernel.Tx) error {
	if err := temporalmigrate.RunPreflight(tx, vp040V82Preflight); err != nil {
		return err
	}
	if err := temporalmigrate.Exec(tx, vp040V82Postgres, "vp040 v82 postgres convert"); err != nil {
		return err
	}
	return temporalmigrate.RunPostgresVerify(tx, vp040V82PostgresVerify, "vp040 v82 postgres verify")
}

// VP040TemporalDescriptor returns the workspace-040 R2 conversion descriptor
// (D-014 allocation, D-017 checksum convention).
func VP040TemporalDescriptor() kernel.MigrationContribution {
	return kernel.MigrationContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: vp040V82Name},
		Version:              82,
		Name:                 vp040V82Name,
		Checksum:             kernel.MigrationChecksum(vp040V82Statements(), vp040V82TransformID),
		Apply:                applyvp040V82,
		ApplyPostgres:        applyvp040V82Postgres,
	}
}
