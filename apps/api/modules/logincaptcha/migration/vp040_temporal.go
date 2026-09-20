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
// ModuleID admin.login-captcha · version 79 · D-014 allocation / D-017 transform id.
package migration

import (
	"github.com/magicvr/schema-ui-core/apps/api/internal/temporalmigrate"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

const (
	vp040V79Name        = "vp040_temporal_captcha"
	vp040V79TransformID = "0079:vp040-temporal-captcha:v1"
)

// vp040V79Guards are the m0 preconditions on tables retired by earlier history
// (descriptor ledger §1: v73 asserts the retired `records` table is absent).
var vp040V79Guards = []temporalmigrate.Guard{}

// vp040V79Preflight is the m0 sentinel census (Root D-012 / D-015 policy).
var vp040V79Preflight = []temporalmigrate.Preflight{}

// vp040V79Verify is the m4 post-rebuild assertion set.
var vp040V79Verify = []temporalmigrate.Verify{
	{Table: "captcha_challenges", Columns: []string{"expires_at", "created_at"}},
	{Table: "captcha_config", Columns: []string{"created_at", "updated_at"}},
}

// vp040V79Rebuild is the ordered m1–m3 slice: constraint handling, table
// rebuild with the frozen conversion expressions, then every index.
var vp040V79Rebuild = []string{
	`ALTER TABLE "captcha_challenges" RENAME TO "captcha_challenges_old"`,
	`CREATE TABLE captcha_challenges (
  id          TEXT PRIMARY KEY,
  answer_hash TEXT NOT NULL,
  expires_at  TEXT NOT NULL,
  created_at  TEXT NOT NULL
)`,
	`INSERT INTO "captcha_challenges" ("id", "answer_hash", "expires_at", "created_at")
SELECT "id", "answer_hash", strftime('%Y-%m-%dT%H:%M:%S', expires_at, 'unixepoch') || '.000000Z', strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z'
FROM "captcha_challenges_old"`,
	`DROP TABLE "captcha_challenges_old"`,
	`ALTER TABLE "captcha_config" RENAME TO "captcha_config_old"`,
	`CREATE TABLE captcha_config (
  id         INTEGER PRIMARY KEY CHECK (id = 1),
  enabled    INTEGER NOT NULL DEFAULT 0,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
)`,
	`INSERT INTO "captcha_config" ("id", "enabled", "created_at", "updated_at")
SELECT "id", "enabled", strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z', strftime('%Y-%m-%dT%H:%M:%S', updated_at, 'unixepoch') || '.000000Z'
FROM "captcha_config_old"`,
	`DROP TABLE "captcha_config_old"`,
}

// vp040V79Statements returns the canonical checksum input (D-017: m0 → m1–m3 → m4).
func vp040V79Statements() []string {
	return temporalmigrate.OrderedWithGuards(vp040V79Guards, vp040V79Preflight, vp040V79Rebuild, vp040V79Verify)
}

// applyvp040V79 is the SQLite Apply body.
func applyvp040V79(tx kernel.Tx) error {
	if err := temporalmigrate.RunGuards(tx, vp040V79Guards, "vp040 v79 sqlite guards"); err != nil {
		return err
	}
	if err := temporalmigrate.RunPreflight(tx, vp040V79Preflight); err != nil {
		return err
	}
	if err := temporalmigrate.Exec(tx, vp040V79Rebuild, "vp040 v79 sqlite rebuild"); err != nil {
		return err
	}
	return temporalmigrate.RunVerify(tx, vp040V79Verify, "vp040 v79 sqlite verify")
}

// vp040V79Postgres is the explicit postgres m1–m2 slice (D-019 §6: no regex
// derivation; the millisecond family uses the integer-split expression).
var vp040V79Postgres = []string{
	`ALTER TABLE "captcha_challenges" ALTER COLUMN "expires_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("expires_at"::double precision)))`,
	`ALTER TABLE "captcha_challenges" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)))`,
	`ALTER TABLE "captcha_config" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)))`,
	`ALTER TABLE "captcha_config" ALTER COLUMN "updated_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("updated_at"::double precision)))`,
}

// vp040V79PostgresVerify is the postgres m4 type/precision assertion set.
var vp040V79PostgresVerify = []temporalmigrate.PgVerify{
	{Table: "captcha_challenges", Column: "expires_at", NonNull: true},
	{Table: "captcha_challenges", Column: "created_at", NonNull: true},
	{Table: "captcha_config", Column: "created_at", NonNull: true},
	{Table: "captcha_config", Column: "updated_at", NonNull: true},
}

// applyvp040V79Postgres is the postgres Apply body.
func applyvp040V79Postgres(tx kernel.Tx) error {
	if err := temporalmigrate.RunPostgresGuards(tx, vp040V79Guards, "vp040 v79 postgres guards"); err != nil {
		return err
	}
	if err := temporalmigrate.RunPreflight(tx, vp040V79Preflight); err != nil {
		return err
	}
	if err := temporalmigrate.Exec(tx, vp040V79Postgres, "vp040 v79 postgres convert"); err != nil {
		return err
	}
	return temporalmigrate.RunPostgresVerify(tx, vp040V79PostgresVerify, "vp040 v79 postgres verify")
}

// VP040TemporalDescriptor returns the workspace-040 R2 conversion descriptor
// (D-014 allocation, D-017 checksum convention).
func VP040TemporalDescriptor() kernel.MigrationContribution {
	return kernel.MigrationContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: vp040V79Name},
		Version:              79,
		Name:                 vp040V79Name,
		Checksum:             kernel.MigrationChecksum(vp040V79Statements(), vp040V79TransformID),
		Apply:                applyvp040V79,
		ApplyPostgres:        applyvp040V79Postgres,
	}
}
