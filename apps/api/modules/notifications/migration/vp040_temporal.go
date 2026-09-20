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
// ModuleID admin.notifications · version 81 · D-014 allocation / D-017 transform id.
package migration

import (
	"github.com/magicvr/schema-ui-core/apps/api/internal/temporalmigrate"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

const (
	vp040V81Name        = "vp040_temporal_notifications"
	vp040V81TransformID = "0081:vp040-temporal-notifications:v1"
)

// vp040V81Preflight is the m0 sentinel census (Root D-012 / D-015 policy).
var vp040V81Preflight = []temporalmigrate.Preflight{}

// vp040V81Verify is the m4 post-rebuild assertion set.
var vp040V81Verify = []temporalmigrate.Verify{
	{Table: "notifications", Columns: []string{"read_at", "created_at"}},
}

// vp040V81Rebuild is the ordered m1–m3 slice: constraint handling, table
// rebuild with the frozen conversion expressions, then every index.
var vp040V81Rebuild = []string{
	`ALTER TABLE "notifications" RENAME TO "notifications_old"`,
	`CREATE TABLE notifications (
  id         TEXT PRIMARY KEY,
  user_id    TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  event      TEXT NOT NULL CHECK (event IN ('account.locked','account.disabled','account.unlocked','account.password-changed')),
  title      TEXT NOT NULL,
  body       TEXT NOT NULL,
  read_at    TEXT,
  created_at TEXT NOT NULL
, title_key TEXT, body_key TEXT)`,
	`INSERT INTO "notifications" ("id", "user_id", "event", "title", "body", "read_at", "created_at", "title_key", "body_key")
SELECT "id", "user_id", "event", "title", "body", CASE WHEN read_at IS NULL THEN NULL ELSE strftime('%Y-%m-%dT%H:%M:%S', read_at, 'unixepoch') || '.000000Z' END, strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z', "title_key", "body_key"
FROM "notifications_old"`,
	`DROP TABLE "notifications_old"`,
	`CREATE INDEX idx_notifications_user_created ON notifications(user_id, created_at DESC)`,
}

// vp040V81Statements returns the canonical checksum input (D-017: m0 → m1–m3 → m4).
func vp040V81Statements() []string {
	return temporalmigrate.Ordered(vp040V81Preflight, vp040V81Rebuild, vp040V81Verify)
}

// applyvp040V81 is the SQLite Apply body.
func applyvp040V81(tx kernel.Tx) error {
	if err := temporalmigrate.RunPreflight(tx, vp040V81Preflight); err != nil {
		return err
	}
	if err := temporalmigrate.Exec(tx, vp040V81Rebuild, "vp040 v81 sqlite rebuild"); err != nil {
		return err
	}
	return temporalmigrate.RunVerify(tx, vp040V81Verify, "vp040 v81 sqlite verify")
}

// vp040V81Postgres is the explicit postgres m1–m2 slice (D-019 §6: no regex
// derivation; the millisecond family uses the integer-split expression).
var vp040V81Postgres = []string{
	`ALTER TABLE "notifications" ALTER COLUMN "read_at" TYPE timestamptz(6) USING (CASE WHEN "read_at" IS NULL THEN NULL ELSE date_trunc('microseconds', to_timestamp("read_at"::double precision)) END)`,
	`ALTER TABLE "notifications" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)))`,
}

// vp040V81PostgresVerify is the postgres m4 type/precision assertion set.
var vp040V81PostgresVerify = []temporalmigrate.PgVerify{
	{Table: "notifications", Column: "read_at", NonNull: false},
	{Table: "notifications", Column: "created_at", NonNull: true},
}

// applyvp040V81Postgres is the postgres Apply body.
func applyvp040V81Postgres(tx kernel.Tx) error {
	if err := temporalmigrate.RunPreflight(tx, vp040V81Preflight); err != nil {
		return err
	}
	if err := temporalmigrate.Exec(tx, vp040V81Postgres, "vp040 v81 postgres convert"); err != nil {
		return err
	}
	return temporalmigrate.RunPostgresVerify(tx, vp040V81PostgresVerify, "vp040 v81 postgres verify")
}

// VP040TemporalDescriptor returns the workspace-040 R2 conversion descriptor
// (D-014 allocation, D-017 checksum convention).
func VP040TemporalDescriptor() kernel.MigrationContribution {
	return kernel.MigrationContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: vp040V81Name},
		Version:              81,
		Name:                 vp040V81Name,
		Checksum:             kernel.MigrationChecksum(vp040V81Statements(), vp040V81TransformID),
		Apply:                applyvp040V81,
		ApplyPostgres:        applyvp040V81Postgres,
	}
}
