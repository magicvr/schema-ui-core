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
// ModuleID admin.data-permission · version 78 · D-014 allocation / D-017 transform id.
package migration

import (
	"github.com/magicvr/schema-ui-core/apps/api/internal/temporalmigrate"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

const (
	vp040V78Name        = "vp040_temporal_data_permission"
	vp040V78TransformID = "0078:vp040-temporal-data-permission:v1"
)

// vp040V78Guards are the m0 preconditions on tables retired by earlier history
// (descriptor ledger §1: v73 asserts the retired `records` table is absent).
var vp040V78Guards = []temporalmigrate.Guard{}

// vp040V78Preflight is the m0 sentinel census (Root D-012 / D-015 policy).
var vp040V78Preflight = []temporalmigrate.Preflight{}

// vp040V78Verify is the m4 post-rebuild assertion set.
var vp040V78Verify = []temporalmigrate.Verify{
	{Table: "data_scope_policies", Columns: []string{"updated_at"}},
	{Table: "user_data_scopes", Columns: []string{"updated_at"}},
}

// vp040V78Rebuild is the ordered m1–m3 slice: constraint handling, table
// rebuild with the frozen conversion expressions, then every index.
var vp040V78Rebuild = []string{
	`ALTER TABLE "data_scope_policies" RENAME TO "data_scope_policies_old"`,
	`CREATE TABLE data_scope_policies (
  resource      TEXT PRIMARY KEY,
  owner_column  TEXT NOT NULL,
  default_scope TEXT NOT NULL CHECK (default_scope IN ('all','self')),
  enabled       INTEGER NOT NULL DEFAULT 1,
  updated_at    TEXT NOT NULL
)`,
	`INSERT INTO "data_scope_policies" ("resource", "owner_column", "default_scope", "enabled", "updated_at")
SELECT "resource", "owner_column", "default_scope", "enabled", strftime('%Y-%m-%dT%H:%M:%S', updated_at, 'unixepoch') || '.000000Z'
FROM "data_scope_policies_old"`,
	`DROP TABLE "data_scope_policies_old"`,
	`ALTER TABLE "user_data_scopes" RENAME TO "user_data_scopes_old"`,
	`CREATE TABLE user_data_scopes (
  user_id    TEXT NOT NULL,
  resource   TEXT NOT NULL,
  scope_type TEXT NOT NULL CHECK (scope_type IN ('all','self')),
  updated_at TEXT NOT NULL,
  PRIMARY KEY (user_id, resource)
)`,
	`INSERT INTO "user_data_scopes" ("user_id", "resource", "scope_type", "updated_at")
SELECT "user_id", "resource", "scope_type", strftime('%Y-%m-%dT%H:%M:%S', updated_at, 'unixepoch') || '.000000Z'
FROM "user_data_scopes_old"`,
	`DROP TABLE "user_data_scopes_old"`,
}

// vp040V78Statements returns the canonical checksum input (D-017: m0 → m1–m3 → m4).
func vp040V78Statements() []string {
	return temporalmigrate.OrderedWithGuards(vp040V78Guards, vp040V78Preflight, vp040V78Rebuild, vp040V78Verify)
}

// applyvp040V78 is the SQLite Apply body.
func applyvp040V78(tx kernel.Tx) error {
	if err := temporalmigrate.RunGuards(tx, vp040V78Guards, "vp040 v78 sqlite guards"); err != nil {
		return err
	}
	if err := temporalmigrate.RunPreflight(tx, vp040V78Preflight); err != nil {
		return err
	}
	if err := temporalmigrate.Exec(tx, vp040V78Rebuild, "vp040 v78 sqlite rebuild"); err != nil {
		return err
	}
	return temporalmigrate.RunVerify(tx, vp040V78Verify, "vp040 v78 sqlite verify")
}

// vp040V78Postgres is the explicit postgres m1–m2 slice (D-019 §6: no regex
// derivation; the millisecond family uses the integer-split expression).
var vp040V78Postgres = []string{
	`ALTER TABLE "data_scope_policies" ALTER COLUMN "updated_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("updated_at"::double precision)))`,
	`ALTER TABLE "user_data_scopes" ALTER COLUMN "updated_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("updated_at"::double precision)))`,
}

// vp040V78PostgresVerify is the postgres m4 type/precision assertion set.
var vp040V78PostgresVerify = []temporalmigrate.PgVerify{
	{Table: "data_scope_policies", Column: "updated_at", NonNull: true},
	{Table: "user_data_scopes", Column: "updated_at", NonNull: true},
}

// applyvp040V78Postgres is the postgres Apply body.
func applyvp040V78Postgres(tx kernel.Tx) error {
	if err := temporalmigrate.RunPostgresGuards(tx, vp040V78Guards, "vp040 v78 postgres guards"); err != nil {
		return err
	}
	if err := temporalmigrate.RunPreflight(tx, vp040V78Preflight); err != nil {
		return err
	}
	if err := temporalmigrate.Exec(tx, vp040V78Postgres, "vp040 v78 postgres convert"); err != nil {
		return err
	}
	return temporalmigrate.RunPostgresVerify(tx, vp040V78PostgresVerify, "vp040 v78 postgres verify")
}

// VP040TemporalDescriptor returns the workspace-040 R2 conversion descriptor
// (D-014 allocation, D-017 checksum convention).
func VP040TemporalDescriptor() kernel.MigrationContribution {
	return kernel.MigrationContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: vp040V78Name},
		Version:              78,
		Name:                 vp040V78Name,
		Checksum:             kernel.MigrationChecksum(vp040V78Statements(), vp040V78TransformID),
		Apply:                applyvp040V78,
		ApplyPostgres:        applyvp040V78Postgres,
	}
}
