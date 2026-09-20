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
// ModuleID admin.mfa · version 80 · D-014 allocation / D-017 transform id.
package migration

import (
	"github.com/magicvr/schema-ui-core/apps/api/internal/temporalmigrate"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

const (
	vp040V80Name        = "vp040_temporal_mfa"
	vp040V80TransformID = "0080:vp040-temporal-mfa:v1"
)

// vp040V80Preflight is the m0 sentinel census (Root D-012 / D-015 policy).
var vp040V80Preflight = []temporalmigrate.Preflight{}

// vp040V80Verify is the m4 post-rebuild assertion set.
var vp040V80Verify = []temporalmigrate.Verify{
	{Table: "user_mfa", Columns: []string{"created_at", "updated_at"}},
	{Table: "mfa_proofs", Columns: []string{"expires_at", "created_at"}},
}

// vp040V80Rebuild is the ordered m1–m3 slice: constraint handling, table
// rebuild with the frozen conversion expressions, then every index.
var vp040V80Rebuild = []string{
	`ALTER TABLE "user_mfa" RENAME TO "user_mfa_old"`,
	`CREATE TABLE user_mfa (
  user_id                TEXT PRIMARY KEY REFERENCES users(id),
  status                 TEXT NOT NULL CHECK (status IN ('pending','active')),
  totp_secret_ciphertext TEXT NOT NULL,
  recovery_codes_hash    TEXT NOT NULL,
  last_used_step         INTEGER NOT NULL DEFAULT 0,
  created_at             TEXT NOT NULL,
  updated_at             TEXT NOT NULL
)`,
	`INSERT INTO "user_mfa" ("user_id", "status", "totp_secret_ciphertext", "recovery_codes_hash", "last_used_step", "created_at", "updated_at")
SELECT "user_id", "status", "totp_secret_ciphertext", "recovery_codes_hash", "last_used_step", strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z', strftime('%Y-%m-%dT%H:%M:%S', updated_at, 'unixepoch') || '.000000Z'
FROM "user_mfa_old"`,
	`DROP TABLE "user_mfa_old"`,
	`ALTER TABLE "mfa_proofs" RENAME TO "mfa_proofs_old"`,
	`CREATE TABLE mfa_proofs (
  id         TEXT PRIMARY KEY,
  user_id    TEXT NOT NULL REFERENCES users(id),
  fail_count INTEGER NOT NULL DEFAULT 0,
  expires_at TEXT NOT NULL,
  created_at TEXT NOT NULL
)`,
	`INSERT INTO "mfa_proofs" ("id", "user_id", "fail_count", "expires_at", "created_at")
SELECT "id", "user_id", "fail_count", strftime('%Y-%m-%dT%H:%M:%S', expires_at, 'unixepoch') || '.000000Z', strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z'
FROM "mfa_proofs_old"`,
	`DROP TABLE "mfa_proofs_old"`,
}

// vp040V80Statements returns the canonical checksum input (D-017: m0 → m1–m3 → m4).
func vp040V80Statements() []string {
	return temporalmigrate.Ordered(vp040V80Preflight, vp040V80Rebuild, vp040V80Verify)
}

// applyvp040V80 is the SQLite Apply body.
func applyvp040V80(tx kernel.Tx) error {
	if err := temporalmigrate.RunPreflight(tx, vp040V80Preflight); err != nil {
		return err
	}
	if err := temporalmigrate.Exec(tx, vp040V80Rebuild, "vp040 v80 sqlite rebuild"); err != nil {
		return err
	}
	return temporalmigrate.RunVerify(tx, vp040V80Verify, "vp040 v80 sqlite verify")
}

// vp040V80Postgres is the explicit postgres m1–m2 slice (D-019 §6: no regex
// derivation; the millisecond family uses the integer-split expression).
var vp040V80Postgres = []string{
	`ALTER TABLE "user_mfa" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)))`,
	`ALTER TABLE "user_mfa" ALTER COLUMN "updated_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("updated_at"::double precision)))`,
	`ALTER TABLE "mfa_proofs" ALTER COLUMN "expires_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("expires_at"::double precision)))`,
	`ALTER TABLE "mfa_proofs" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)))`,
}

// vp040V80PostgresVerify is the postgres m4 type/precision assertion set.
var vp040V80PostgresVerify = []temporalmigrate.PgVerify{
	{Table: "user_mfa", Column: "created_at", NonNull: true},
	{Table: "user_mfa", Column: "updated_at", NonNull: true},
	{Table: "mfa_proofs", Column: "expires_at", NonNull: true},
	{Table: "mfa_proofs", Column: "created_at", NonNull: true},
}

// applyvp040V80Postgres is the postgres Apply body.
func applyvp040V80Postgres(tx kernel.Tx) error {
	if err := temporalmigrate.RunPreflight(tx, vp040V80Preflight); err != nil {
		return err
	}
	if err := temporalmigrate.Exec(tx, vp040V80Postgres, "vp040 v80 postgres convert"); err != nil {
		return err
	}
	return temporalmigrate.RunPostgresVerify(tx, vp040V80PostgresVerify, "vp040 v80 postgres verify")
}

// VP040TemporalDescriptor returns the workspace-040 R2 conversion descriptor
// (D-014 allocation, D-017 checksum convention).
func VP040TemporalDescriptor() kernel.MigrationContribution {
	return kernel.MigrationContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: vp040V80Name},
		Version:              80,
		Name:                 vp040V80Name,
		Checksum:             kernel.MigrationChecksum(vp040V80Statements(), vp040V80TransformID),
		Apply:                applyvp040V80,
		ApplyPostgres:        applyvp040V80Postgres,
	}
}
