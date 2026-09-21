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
// ModuleID core.persistence · version 73 · D-014 allocation / D-017 transform id.
package migration

import (
	"github.com/magicvr/schema-ui-core/apps/api/internal/temporalmigrate"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

const (
	vp040V73Name        = "vp040_temporal_core_persistence"
	vp040V73TransformID = "0073:vp040-temporal-core-persistence:v1"
)

// vp040V73Guards are the m0 preconditions on tables retired by earlier history
// (descriptor ledger §1: v73 asserts the retired `records` table is absent).
var vp040V73Guards = []temporalmigrate.Guard{
	{Table: "records"},
}

// vp040V73Preflight is the m0 sentinel census (Root D-012 / D-015 policy).
var vp040V73Preflight = []temporalmigrate.Preflight{
	{Table: "mail_config", Column: "updated_at", Voucher: false},
}

// vp040V73Verify is the m4 post-rebuild assertion set.
var vp040V73Verify = []temporalmigrate.Verify{
	{Table: "schema_migrations", Columns: []string{"applied_at"}},
	{Table: "mail_outbox", Columns: []string{"created_at"}},
	{Table: "mail_config", Columns: []string{"updated_at"}},
}

// vp040V73Rebuild is the ordered m1–m3 slice: constraint handling, table
// rebuild with the frozen conversion expressions, then every index.
var vp040V73Rebuild = []string{
	`ALTER TABLE "schema_migrations" RENAME TO "schema_migrations_old"`,
	`CREATE TABLE schema_migrations (
  version    INTEGER PRIMARY KEY,
  name       TEXT NOT NULL UNIQUE,
  checksum   TEXT NOT NULL CHECK (length(checksum) = 64),
  applied_at TEXT NOT NULL
)`,
	`INSERT INTO "schema_migrations" ("version", "name", "checksum", "applied_at")
SELECT "version", "name", "checksum", strftime('%Y-%m-%dT%H:%M:%S', applied_at, 'unixepoch') || '.000000Z'
FROM "schema_migrations_old"`,
	`DROP TABLE "schema_migrations_old"`,
	`ALTER TABLE "mail_outbox" RENAME TO "mail_outbox_old"`,
	`CREATE TABLE mail_outbox (
  id         TEXT PRIMARY KEY,
  to_addr    TEXT NOT NULL,
  subject    TEXT NOT NULL,
  body       TEXT NOT NULL,
  created_at TEXT NOT NULL
, channel TEXT NOT NULL DEFAULT 'mock', delivery_status TEXT NOT NULL DEFAULT 'delivered')`,
	`INSERT INTO "mail_outbox" ("id", "to_addr", "subject", "body", "created_at", "channel", "delivery_status")
SELECT "id", "to_addr", "subject", "body", strftime('%Y-%m-%dT%H:%M:%S', CASE WHEN created_at >= 0 THEN created_at/1000 ELSE (created_at-999)/1000 END, 'unixepoch') || '.' || printf('%03d', (created_at%1000 + 1000) % 1000) || '000Z', "channel", "delivery_status"
FROM "mail_outbox_old"`,
	`DROP TABLE "mail_outbox_old"`,
	`ALTER TABLE "mail_config" RENAME TO "mail_config_old"`,
	`CREATE TABLE mail_config (
  id                 INTEGER PRIMARY KEY CHECK (id = 1),
  channel            TEXT    NOT NULL DEFAULT 'mock',
  mock_retention     INTEGER NOT NULL DEFAULT 500,
  resend_from        TEXT    NOT NULL DEFAULT '',
  resend_api_key_enc TEXT    NOT NULL DEFAULT '',
  smtp_host          TEXT    NOT NULL DEFAULT '',
  smtp_port          INTEGER NOT NULL DEFAULT 0,
  smtp_username      TEXT    NOT NULL DEFAULT '',
  smtp_password_enc  TEXT    NOT NULL DEFAULT '',
  smtp_from          TEXT    NOT NULL DEFAULT '',
  updated_at         TEXT
)`,
	`INSERT INTO "mail_config" ("id", "channel", "mock_retention", "resend_from", "resend_api_key_enc", "smtp_host", "smtp_port", "smtp_username", "smtp_password_enc", "smtp_from", "updated_at")
SELECT "id", "channel", "mock_retention", "resend_from", "resend_api_key_enc", "smtp_host", "smtp_port", "smtp_username", "smtp_password_enc", "smtp_from", CASE WHEN updated_at = 0 THEN NULL ELSE strftime('%Y-%m-%dT%H:%M:%S', CASE WHEN updated_at >= 0 THEN updated_at/1000 ELSE (updated_at-999)/1000 END, 'unixepoch') || '.' || printf('%03d', (updated_at%1000 + 1000) % 1000) || '000Z' END
FROM "mail_config_old"`,
	`DROP TABLE "mail_config_old"`,
	`CREATE INDEX idx_mail_outbox_created_at ON mail_outbox(created_at)`,
}

// vp040V73Statements returns the canonical checksum input (D-017: m0 → m1–m3 → m4).
func vp040V73Statements() []string {
	return temporalmigrate.OrderedWithGuards(vp040V73Guards, vp040V73Preflight, vp040V73Rebuild, vp040V73Verify)
}

// applyvp040V73 is the SQLite Apply body.
func applyvp040V73(tx kernel.Tx) error {
	if err := temporalmigrate.RunGuards(tx, vp040V73Guards, "vp040 v73 sqlite guards"); err != nil {
		return err
	}
	if err := temporalmigrate.RunPreflight(tx, vp040V73Preflight); err != nil {
		return err
	}
	if err := temporalmigrate.Exec(tx, vp040V73Rebuild, "vp040 v73 sqlite rebuild"); err != nil {
		return err
	}
	return temporalmigrate.RunVerify(tx, vp040V73Verify, "vp040 v73 sqlite verify")
}

// vp040V73Postgres is the explicit postgres m1–m2 slice (D-019 §6: no regex
// derivation; the millisecond family uses the integer-split expression).
var vp040V73Postgres = []string{
	`ALTER TABLE "schema_migrations" ALTER COLUMN "applied_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("applied_at"::double precision)))`,
	`ALTER TABLE "mail_outbox" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (TIMESTAMPTZ 'epoch' + (("created_at" - CASE WHEN "created_at" >= 0 THEN 0 ELSE 999 END) / 1000) * INTERVAL '1 second' + ((("created_at" % 1000) + 1000) % 1000) * INTERVAL '1 millisecond')`,
	`ALTER TABLE "mail_config" ALTER COLUMN "updated_at" DROP DEFAULT`,
	`ALTER TABLE "mail_config" ALTER COLUMN "updated_at" DROP NOT NULL`,
	`ALTER TABLE "mail_config" ALTER COLUMN "updated_at" TYPE timestamptz(6) USING (CASE WHEN "updated_at" = 0 THEN NULL ELSE TIMESTAMPTZ 'epoch' + (("updated_at" - CASE WHEN "updated_at" >= 0 THEN 0 ELSE 999 END) / 1000) * INTERVAL '1 second' + ((("updated_at" % 1000) + 1000) % 1000) * INTERVAL '1 millisecond' END)`,
}

// vp040V73PostgresVerify is the postgres m4 type/precision assertion set.
var vp040V73PostgresVerify = []temporalmigrate.PgVerify{
	{Table: "schema_migrations", Column: "applied_at", NonNull: true},
	{Table: "mail_outbox", Column: "created_at", NonNull: true},
	{Table: "mail_config", Column: "updated_at", NonNull: false},
}

// applyvp040V73Postgres is the postgres Apply body.
func applyvp040V73Postgres(tx kernel.Tx) error {
	if err := temporalmigrate.RunPostgresGuards(tx, vp040V73Guards, "vp040 v73 postgres guards"); err != nil {
		return err
	}
	if err := temporalmigrate.RunPreflight(tx, vp040V73Preflight); err != nil {
		return err
	}
	if err := temporalmigrate.Exec(tx, vp040V73Postgres, "vp040 v73 postgres convert"); err != nil {
		return err
	}
	return temporalmigrate.RunPostgresVerify(tx, vp040V73PostgresVerify, "vp040 v73 postgres verify")
}

// VP040TemporalDescriptor returns the workspace-040 R2 conversion descriptor
// (D-014 allocation, D-017 checksum convention).
func VP040TemporalDescriptor() kernel.MigrationContribution {
	return kernel.MigrationContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: vp040V73Name},
		Version:              73,
		Name:                 vp040V73Name,
		Checksum:             kernel.MigrationChecksum(vp040V73Statements(), vp040V73TransformID),
		Apply:                applyvp040V73,
		ApplyPostgres:        applyvp040V73Postgres,
	}
}
