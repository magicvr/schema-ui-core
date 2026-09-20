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
// ModuleID core.operationlog · version 75 · D-014 allocation / D-017 transform id.
package migration

import (
	"github.com/magicvr/schema-ui-core/apps/api/internal/temporalmigrate"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

const (
	vp040V75Name        = "vp040_temporal_operationlog"
	vp040V75TransformID = "0075:vp040-temporal-operationlog:v1"
)

// vp040V75Guards are the m0 preconditions on tables retired by earlier history
// (descriptor ledger §1: v73 asserts the retired `records` table is absent).
var vp040V75Guards = []temporalmigrate.Guard{}

// vp040V75Preflight is the m0 sentinel census (Root D-012 / D-015 policy).
var vp040V75Preflight = []temporalmigrate.Preflight{}

// vp040V75Verify is the m4 post-rebuild assertion set.
var vp040V75Verify = []temporalmigrate.Verify{
	{Table: "operation_log", Columns: []string{"created_at"}, Children: []string{"operation_log_correlation", "operation_log_session"}},
	{Table: "operation_log_archive", Columns: []string{"created_at", "archived_at"}},
}

// vp040V75Rebuild is the ordered m1–m3 slice: constraint handling, table
// rebuild with the frozen conversion expressions, then every index.
var vp040V75Rebuild = []string{
	`CREATE TEMP TABLE "operation_log_correlation_bak" AS SELECT * FROM "operation_log_correlation"`,
	`CREATE TEMP TABLE "operation_log_session_bak" AS SELECT * FROM "operation_log_session"`,
	`DROP TABLE "operation_log_correlation"`,
	`DROP TABLE "operation_log_session"`,
	`ALTER TABLE "operation_log" RENAME TO "operation_log_old"`,
	`CREATE TABLE operation_log (
  id         TEXT PRIMARY KEY,
  event      TEXT NOT NULL CHECK (event IN ('records.create','records.update','records.delete','auth.login','auth.logout','auth.refresh','users.create','users.update','users.delete','roles.create','roles.update','roles.delete','settings.update','users.enable','users.disable','users.unlock','account.password-change','account.session-revoke','data.export','data.import','files.upload','files.download','files.delete','dictionary.create','dictionary.update','dictionary.delete','scheduled-tasks.create','scheduled-tasks.update','scheduled-tasks.delete','captcha.settings-update','recycle.restore','recycle.purge','data-permission.policy-update','data-permission.scope-update','mfa.enroll','mfa.confirm','mfa.disable','mfa.recovery-rotate','mfa.admin-reset','mfa.login','wallet.account-create','wallet.account-update','wallet.adjust','wallet.freeze','wallet.unfreeze','wallet.reconcile','wallet.deduct-frozen','account.avatar-change','wallet.reconcile.queued','wallet.reconcile.failed','wallet.reconcile.cancelled','service-credentials.create','service-credentials.use','service-credentials.revoke','mail.channel-update','mail.test-send','bizoffer.offer.create','bizoffer.offer.update','bizoffer.offer.status','bizoffer.entitlement.void')),
  actor_id   TEXT NOT NULL,
  actor_name TEXT NOT NULL,
  record_id  TEXT,
  detail     TEXT,
  created_at TEXT NOT NULL
)`,
	`INSERT INTO "operation_log" ("id", "event", "actor_id", "actor_name", "record_id", "detail", "created_at")
SELECT "id", "event", "actor_id", "actor_name", "record_id", "detail", strftime('%Y-%m-%dT%H:%M:%S', CASE WHEN created_at >= 0 THEN created_at/1000 ELSE (created_at-999)/1000 END, 'unixepoch') || '.' || printf('%03d', (created_at%1000 + 1000) % 1000) || '000Z'
FROM "operation_log_old"`,
	`DROP TABLE "operation_log_old"`,
	`ALTER TABLE "operation_log_archive" RENAME TO "operation_log_archive_old"`,
	`CREATE TABLE operation_log_archive (
  id          TEXT PRIMARY KEY,
  event       TEXT NOT NULL,
  actor_id    TEXT NOT NULL,
  actor_name  TEXT NOT NULL,
  record_id   TEXT,
  detail      TEXT,
  created_at  TEXT NOT NULL,
  archived_at TEXT NOT NULL
)`,
	`INSERT INTO "operation_log_archive" ("id", "event", "actor_id", "actor_name", "record_id", "detail", "created_at", "archived_at")
SELECT "id", "event", "actor_id", "actor_name", "record_id", "detail", strftime('%Y-%m-%dT%H:%M:%S', CASE WHEN created_at >= 0 THEN created_at/1000 ELSE (created_at-999)/1000 END, 'unixepoch') || '.' || printf('%03d', (created_at%1000 + 1000) % 1000) || '000Z', strftime('%Y-%m-%dT%H:%M:%S', CASE WHEN archived_at >= 0 THEN archived_at/1000 ELSE (archived_at-999)/1000 END, 'unixepoch') || '.' || printf('%03d', (archived_at%1000 + 1000) % 1000) || '000Z'
FROM "operation_log_archive_old"`,
	`DROP TABLE "operation_log_archive_old"`,
	`CREATE TABLE operation_log_correlation (
  operation_id   TEXT PRIMARY KEY REFERENCES operation_log(id) ON DELETE CASCADE,
  correlation_id TEXT NOT NULL
)`,
	`INSERT INTO "operation_log_correlation" ("operation_id", "correlation_id")
SELECT "operation_id", "correlation_id"
FROM "operation_log_correlation_bak"`,
	`DROP TABLE "operation_log_correlation_bak"`,
	`CREATE TABLE operation_log_session (
  operation_id TEXT PRIMARY KEY REFERENCES operation_log(id) ON DELETE CASCADE,
  session_id   TEXT NOT NULL
)`,
	`INSERT INTO "operation_log_session" ("operation_id", "session_id")
SELECT "operation_id", "session_id"
FROM "operation_log_session_bak"`,
	`DROP TABLE "operation_log_session_bak"`,
	`CREATE INDEX idx_operation_log_created_at ON operation_log(created_at DESC)`,
	`CREATE INDEX idx_operation_log_archive_created_at ON operation_log_archive(created_at DESC)`,
	`CREATE INDEX idx_operation_log_correlation_id ON operation_log_correlation(correlation_id)`,
	`CREATE INDEX idx_operation_log_session_id ON operation_log_session(session_id)`,
}

// vp040V75Statements returns the canonical checksum input (D-017: m0 → m1–m3 → m4).
func vp040V75Statements() []string {
	return temporalmigrate.OrderedWithGuards(vp040V75Guards, vp040V75Preflight, vp040V75Rebuild, vp040V75Verify)
}

// applyvp040V75 is the SQLite Apply body.
func applyvp040V75(tx kernel.Tx) error {
	if err := temporalmigrate.RunGuards(tx, vp040V75Guards, "vp040 v75 sqlite guards"); err != nil {
		return err
	}
	if err := temporalmigrate.RunPreflight(tx, vp040V75Preflight); err != nil {
		return err
	}
	if err := temporalmigrate.Exec(tx, vp040V75Rebuild, "vp040 v75 sqlite rebuild"); err != nil {
		return err
	}
	return temporalmigrate.RunVerify(tx, vp040V75Verify, "vp040 v75 sqlite verify")
}

// vp040V75Postgres is the explicit postgres m1–m2 slice (D-019 §6: no regex
// derivation; the millisecond family uses the integer-split expression).
var vp040V75Postgres = []string{
	`ALTER TABLE "operation_log" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (TIMESTAMPTZ 'epoch' + (("created_at" - CASE WHEN "created_at" >= 0 THEN 0 ELSE 999 END) / 1000) * INTERVAL '1 second' + ((("created_at" % 1000) + 1000) % 1000) * INTERVAL '1 millisecond')`,
	`ALTER TABLE "operation_log_archive" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (TIMESTAMPTZ 'epoch' + (("created_at" - CASE WHEN "created_at" >= 0 THEN 0 ELSE 999 END) / 1000) * INTERVAL '1 second' + ((("created_at" % 1000) + 1000) % 1000) * INTERVAL '1 millisecond')`,
	`ALTER TABLE "operation_log_archive" ALTER COLUMN "archived_at" TYPE timestamptz(6) USING (TIMESTAMPTZ 'epoch' + (("archived_at" - CASE WHEN "archived_at" >= 0 THEN 0 ELSE 999 END) / 1000) * INTERVAL '1 second' + ((("archived_at" % 1000) + 1000) % 1000) * INTERVAL '1 millisecond')`,
}

// vp040V75PostgresVerify is the postgres m4 type/precision assertion set.
var vp040V75PostgresVerify = []temporalmigrate.PgVerify{
	{Table: "operation_log", Column: "created_at", NonNull: true},
	{Table: "operation_log_archive", Column: "created_at", NonNull: true},
	{Table: "operation_log_archive", Column: "archived_at", NonNull: true},
}

// applyvp040V75Postgres is the postgres Apply body.
func applyvp040V75Postgres(tx kernel.Tx) error {
	if err := temporalmigrate.RunPostgresGuards(tx, vp040V75Guards, "vp040 v75 postgres guards"); err != nil {
		return err
	}
	if err := temporalmigrate.RunPreflight(tx, vp040V75Preflight); err != nil {
		return err
	}
	if err := temporalmigrate.Exec(tx, vp040V75Postgres, "vp040 v75 postgres convert"); err != nil {
		return err
	}
	return temporalmigrate.RunPostgresVerify(tx, vp040V75PostgresVerify, "vp040 v75 postgres verify")
}

// VP040TemporalDescriptor returns the workspace-040 R2 conversion descriptor
// (D-014 allocation, D-017 checksum convention).
func VP040TemporalDescriptor() kernel.MigrationContribution {
	return kernel.MigrationContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: vp040V75Name},
		Version:              75,
		Name:                 vp040V75Name,
		Checksum:             kernel.MigrationChecksum(vp040V75Statements(), vp040V75TransformID),
		Apply:                applyvp040V75,
		ApplyPostgres:        applyvp040V75Postgres,
	}
}
