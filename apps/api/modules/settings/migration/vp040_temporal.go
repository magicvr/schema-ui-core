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
// ModuleID admin.settings · version 84 · D-014 allocation / D-017 transform id.
package migration

import (
	"github.com/magicvr/schema-ui-core/apps/api/internal/temporalmigrate"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

const (
	vp040V84Name        = "vp040_temporal_settings"
	vp040V84TransformID = "0084:vp040-temporal-settings:v1"
)

// vp040V84Preflight is the m0 sentinel census (Root D-012 / D-015 policy).
var vp040V84Preflight = []temporalmigrate.Preflight{}

// vp040V84Verify is the m4 post-rebuild assertion set.
var vp040V84Verify = []temporalmigrate.Verify{
	{Table: "site_settings", Columns: []string{"updated_at"}},
}

// vp040V84Rebuild is the ordered m1–m3 slice: constraint handling, table
// rebuild with the frozen conversion expressions, then every index.
var vp040V84Rebuild = []string{
	`ALTER TABLE "site_settings" RENAME TO "site_settings_old"`,
	`CREATE TABLE site_settings (
  id         TEXT PRIMARY KEY CHECK (id = 'default'),
  site_title TEXT NOT NULL,
  logo_url   TEXT NOT NULL DEFAULT '',
  updated_at TEXT NOT NULL
, logo_url_light TEXT NOT NULL DEFAULT '', logo_url_dark TEXT NOT NULL DEFAULT '', favicon_url TEXT NOT NULL DEFAULT '', default_locale TEXT NOT NULL DEFAULT '', site_timezone TEXT NOT NULL DEFAULT '', default_theme TEXT NOT NULL DEFAULT '', copyright_text TEXT NOT NULL DEFAULT '', icp_number TEXT NOT NULL DEFAULT '', operation_log_retention_days INTEGER NOT NULL DEFAULT 90, operation_log_expiration_action TEXT NOT NULL DEFAULT 'archive', default_currency TEXT NOT NULL DEFAULT '')`,
	`INSERT INTO "site_settings" ("id", "site_title", "logo_url", "updated_at", "logo_url_light", "logo_url_dark", "favicon_url", "default_locale", "site_timezone", "default_theme", "copyright_text", "icp_number", "operation_log_retention_days", "operation_log_expiration_action", "default_currency")
SELECT "id", "site_title", "logo_url", strftime('%Y-%m-%dT%H:%M:%S', updated_at, 'unixepoch') || '.000000Z', "logo_url_light", "logo_url_dark", "favicon_url", "default_locale", "site_timezone", "default_theme", "copyright_text", "icp_number", "operation_log_retention_days", "operation_log_expiration_action", "default_currency"
FROM "site_settings_old"`,
	`DROP TABLE "site_settings_old"`,
	`CREATE INDEX idx_site_settings_updated_at ON site_settings (updated_at)`,
}

// vp040V84Statements returns the canonical checksum input (D-017: m0 → m1–m3 → m4).
func vp040V84Statements() []string {
	return temporalmigrate.Ordered(vp040V84Preflight, vp040V84Rebuild, vp040V84Verify)
}

// applyvp040V84 is the SQLite Apply body.
func applyvp040V84(tx kernel.Tx) error {
	if err := temporalmigrate.RunPreflight(tx, vp040V84Preflight); err != nil {
		return err
	}
	if err := temporalmigrate.Exec(tx, vp040V84Rebuild, "vp040 v84 sqlite rebuild"); err != nil {
		return err
	}
	return temporalmigrate.RunVerify(tx, vp040V84Verify, "vp040 v84 sqlite verify")
}

// vp040V84Postgres is the explicit postgres m1–m2 slice (D-019 §6: no regex
// derivation; the millisecond family uses the integer-split expression).
var vp040V84Postgres = []string{
	`ALTER TABLE "site_settings" ALTER COLUMN "updated_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("updated_at"::double precision)))`,
}

// vp040V84PostgresVerify is the postgres m4 type/precision assertion set.
var vp040V84PostgresVerify = []temporalmigrate.PgVerify{
	{Table: "site_settings", Column: "updated_at", NonNull: true},
}

// applyvp040V84Postgres is the postgres Apply body.
func applyvp040V84Postgres(tx kernel.Tx) error {
	if err := temporalmigrate.RunPreflight(tx, vp040V84Preflight); err != nil {
		return err
	}
	if err := temporalmigrate.Exec(tx, vp040V84Postgres, "vp040 v84 postgres convert"); err != nil {
		return err
	}
	return temporalmigrate.RunPostgresVerify(tx, vp040V84PostgresVerify, "vp040 v84 postgres verify")
}

// VP040TemporalDescriptor returns the workspace-040 R2 conversion descriptor
// (D-014 allocation, D-017 checksum convention).
func VP040TemporalDescriptor() kernel.MigrationContribution {
	return kernel.MigrationContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: vp040V84Name},
		Version:              84,
		Name:                 vp040V84Name,
		Checksum:             kernel.MigrationChecksum(vp040V84Statements(), vp040V84TransformID),
		Apply:                applyvp040V84,
		ApplyPostgres:        applyvp040V84Postgres,
	}
}
