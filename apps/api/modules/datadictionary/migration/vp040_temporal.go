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
// ModuleID admin.data-dictionary · version 77 · D-014 allocation / D-017 transform id.
package migration

import (
	"github.com/magicvr/schema-ui-core/apps/api/internal/temporalmigrate"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

const (
	vp040V77Name        = "vp040_temporal_dictionary"
	vp040V77TransformID = "0077:vp040-temporal-dictionary:v1"
)

// vp040V77Guards are the m0 preconditions on tables retired by earlier history
// (descriptor ledger §1: v73 asserts the retired `records` table is absent).
var vp040V77Guards = []temporalmigrate.Guard{}

// vp040V77Preflight is the m0 sentinel census (Root D-012 / D-015 policy).
var vp040V77Preflight = []temporalmigrate.Preflight{}

// vp040V77Verify is the m4 post-rebuild assertion set.
var vp040V77Verify = []temporalmigrate.Verify{
	{Table: "dict_types", Columns: []string{"created_at", "updated_at"}, Children: []string{"dict_entries"}},
	{Table: "dict_entries", Columns: []string{"created_at", "updated_at"}},
}

// vp040V77Rebuild is the ordered m1–m3 slice: constraint handling, table
// rebuild with the frozen conversion expressions, then every index.
var vp040V77Rebuild = []string{
	`CREATE TEMP TABLE "dict_entries_bak" AS SELECT * FROM "dict_entries"`,
	`DROP TABLE "dict_entries"`,
	`ALTER TABLE "dict_types" RENAME TO "dict_types_old"`,
	`CREATE TABLE dict_types (
  id         TEXT PRIMARY KEY,
  key        TEXT NOT NULL UNIQUE,
  name       TEXT NOT NULL,
  enabled    INTEGER NOT NULL DEFAULT 1,
  description TEXT,
  sort       INTEGER NOT NULL DEFAULT 0,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
)`,
	`INSERT INTO "dict_types" ("id", "key", "name", "enabled", "description", "sort", "created_at", "updated_at")
SELECT "id", "key", "name", "enabled", "description", "sort", strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z', strftime('%Y-%m-%dT%H:%M:%S', updated_at, 'unixepoch') || '.000000Z'
FROM "dict_types_old"`,
	`DROP TABLE "dict_types_old"`,
	`CREATE TABLE dict_entries (
  id         TEXT PRIMARY KEY,
  dict_key   TEXT NOT NULL REFERENCES dict_types(key) ON DELETE CASCADE,
  entry_key  TEXT NOT NULL,
  label      TEXT NOT NULL,
  enabled    INTEGER NOT NULL DEFAULT 1,
  sort       INTEGER NOT NULL DEFAULT 0,
  remark     TEXT,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL, badge_style TEXT NOT NULL DEFAULT 'default',
  UNIQUE (dict_key, entry_key)
)`,
	`INSERT INTO "dict_entries" ("id", "dict_key", "entry_key", "label", "enabled", "sort", "remark", "created_at", "updated_at", "badge_style")
SELECT "id", "dict_key", "entry_key", "label", "enabled", "sort", "remark", strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z', strftime('%Y-%m-%dT%H:%M:%S', updated_at, 'unixepoch') || '.000000Z', "badge_style"
FROM "dict_entries_bak"`,
	`DROP TABLE "dict_entries_bak"`,
	`CREATE INDEX idx_dict_entries_dict_key ON dict_entries(dict_key, sort)`,
}

// vp040V77Statements returns the canonical checksum input (D-017: m0 → m1–m3 → m4).
func vp040V77Statements() []string {
	return temporalmigrate.OrderedWithGuards(vp040V77Guards, vp040V77Preflight, vp040V77Rebuild, vp040V77Verify)
}

// applyvp040V77 is the SQLite Apply body.
func applyvp040V77(tx kernel.Tx) error {
	if err := temporalmigrate.RunGuards(tx, vp040V77Guards, "vp040 v77 sqlite guards"); err != nil {
		return err
	}
	if err := temporalmigrate.RunPreflight(tx, vp040V77Preflight); err != nil {
		return err
	}
	if err := temporalmigrate.Exec(tx, vp040V77Rebuild, "vp040 v77 sqlite rebuild"); err != nil {
		return err
	}
	return temporalmigrate.RunVerify(tx, vp040V77Verify, "vp040 v77 sqlite verify")
}

// vp040V77Postgres is the explicit postgres m1–m2 slice (D-019 §6: no regex
// derivation; the millisecond family uses the integer-split expression).
var vp040V77Postgres = []string{
	`ALTER TABLE "dict_types" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)))`,
	`ALTER TABLE "dict_types" ALTER COLUMN "updated_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("updated_at"::double precision)))`,
	`ALTER TABLE "dict_entries" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)))`,
	`ALTER TABLE "dict_entries" ALTER COLUMN "updated_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("updated_at"::double precision)))`,
}

// vp040V77PostgresVerify is the postgres m4 type/precision assertion set.
var vp040V77PostgresVerify = []temporalmigrate.PgVerify{
	{Table: "dict_types", Column: "created_at", NonNull: true},
	{Table: "dict_types", Column: "updated_at", NonNull: true},
	{Table: "dict_entries", Column: "created_at", NonNull: true},
	{Table: "dict_entries", Column: "updated_at", NonNull: true},
}

// applyvp040V77Postgres is the postgres Apply body.
func applyvp040V77Postgres(tx kernel.Tx) error {
	if err := temporalmigrate.RunPostgresGuards(tx, vp040V77Guards, "vp040 v77 postgres guards"); err != nil {
		return err
	}
	if err := temporalmigrate.RunPreflight(tx, vp040V77Preflight); err != nil {
		return err
	}
	if err := temporalmigrate.Exec(tx, vp040V77Postgres, "vp040 v77 postgres convert"); err != nil {
		return err
	}
	return temporalmigrate.RunPostgresVerify(tx, vp040V77PostgresVerify, "vp040 v77 postgres verify")
}

// VP040TemporalDescriptor returns the workspace-040 R2 conversion descriptor
// (D-014 allocation, D-017 checksum convention).
func VP040TemporalDescriptor() kernel.MigrationContribution {
	return kernel.MigrationContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: vp040V77Name},
		Version:              77,
		Name:                 vp040V77Name,
		Checksum:             kernel.MigrationChecksum(vp040V77Statements(), vp040V77TransformID),
		Apply:                applyvp040V77,
		ApplyPostgres:        applyvp040V77Postgres,
	}
}
