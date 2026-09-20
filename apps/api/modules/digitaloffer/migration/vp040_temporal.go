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
// ModuleID biz.digital-offer · version 87 · D-014 allocation / D-017 transform id.
package migration

import (
	"github.com/magicvr/schema-ui-core/apps/api/internal/temporalmigrate"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

const (
	vp040V87Name        = "vp040_temporal_digital_offer"
	vp040V87TransformID = "0087:vp040-temporal-digital-offer:v1"
)

// vp040V87Guards are the m0 preconditions on tables retired by earlier history
// (descriptor ledger §1: v73 asserts the retired `records` table is absent).
var vp040V87Guards = []temporalmigrate.Guard{}

// vp040V87Preflight is the m0 sentinel census (Root D-012 / D-015 policy).
var vp040V87Preflight = []temporalmigrate.Preflight{}

// vp040V87Verify is the m4 post-rebuild assertion set.
var vp040V87Verify = []temporalmigrate.Verify{
	{Table: "digital_offers", Columns: []string{"created_at", "updated_at"}},
	{Table: "digital_purchases", Columns: []string{"created_at"}},
	{Table: "digital_entitlements", Columns: []string{"expires_at", "created_at", "updated_at"}},
}

// vp040V87Rebuild is the ordered m1–m3 slice: constraint handling, table
// rebuild with the frozen conversion expressions, then every index.
var vp040V87Rebuild = []string{
	`ALTER TABLE "digital_offers" RENAME TO "digital_offers_old"`,
	`CREATE TABLE digital_offers (
  id                 TEXT PRIMARY KEY,
  name               TEXT NOT NULL,
  description        TEXT NOT NULL DEFAULT '',
  price_amount       INTEGER NOT NULL CHECK (price_amount > 0),
  currency           TEXT NOT NULL,
  entitlement_form   TEXT NOT NULL CHECK (entitlement_form IN ('duration','count')),
  duration_seconds   INTEGER,
  count_per_purchase INTEGER,
  status             TEXT NOT NULL DEFAULT 'draft' CHECK (status IN ('draft','on_sale','off_sale')),
  version            INTEGER NOT NULL DEFAULT 0,
  created_at         TEXT NOT NULL,
  updated_at         TEXT NOT NULL,
  CHECK (
    (entitlement_form = 'duration' AND duration_seconds IS NOT NULL AND duration_seconds >= 1 AND count_per_purchase IS NULL)
    OR
    (entitlement_form = 'count' AND count_per_purchase IS NOT NULL AND count_per_purchase >= 1 AND duration_seconds IS NULL)
  )
)`,
	`INSERT INTO "digital_offers" ("id", "name", "description", "price_amount", "currency", "entitlement_form", "duration_seconds", "count_per_purchase", "status", "version", "created_at", "updated_at")
SELECT "id", "name", "description", "price_amount", "currency", "entitlement_form", "duration_seconds", "count_per_purchase", "status", "version", strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z', strftime('%Y-%m-%dT%H:%M:%S', updated_at, 'unixepoch') || '.000000Z'
FROM "digital_offers_old"`,
	`DROP TABLE "digital_offers_old"`,
	`ALTER TABLE "digital_purchases" RENAME TO "digital_purchases_old"`,
	`CREATE TABLE digital_purchases (
  id              TEXT PRIMARY KEY,
  subject_id      TEXT NOT NULL,
  offer_id        TEXT NOT NULL,
  offer_name      TEXT NOT NULL,
  amount          INTEGER NOT NULL CHECK (amount > 0),
  currency        TEXT NOT NULL,
  freeze_entry_id TEXT NOT NULL,
  deduct_entry_id TEXT NOT NULL,
  request_id      TEXT NOT NULL,
  status          TEXT NOT NULL DEFAULT 'fulfilled' CHECK (status = 'fulfilled'),
  created_at      TEXT NOT NULL,
  UNIQUE (subject_id, request_id)
)`,
	`INSERT INTO "digital_purchases" ("id", "subject_id", "offer_id", "offer_name", "amount", "currency", "freeze_entry_id", "deduct_entry_id", "request_id", "status", "created_at")
SELECT "id", "subject_id", "offer_id", "offer_name", "amount", "currency", "freeze_entry_id", "deduct_entry_id", "request_id", "status", strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z'
FROM "digital_purchases_old"`,
	`DROP TABLE "digital_purchases_old"`,
	`ALTER TABLE "digital_entitlements" RENAME TO "digital_entitlements_old"`,
	`CREATE TABLE digital_entitlements (
  id              TEXT PRIMARY KEY,
  subject_id      TEXT NOT NULL,
  offer_id        TEXT NOT NULL,
  purchase_id     TEXT NOT NULL,
  form            TEXT NOT NULL CHECK (form IN ('duration','count')),
  expires_at      TEXT,
  remaining_count INTEGER,
  status          TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active','voided')),
  created_at      TEXT NOT NULL,
  updated_at      TEXT NOT NULL,
  CHECK (
    (form = 'duration' AND expires_at IS NOT NULL AND remaining_count IS NULL)
    OR
    (form = 'count' AND remaining_count IS NOT NULL AND remaining_count >= 0 AND expires_at IS NULL)
  )
)`,
	`INSERT INTO "digital_entitlements" ("id", "subject_id", "offer_id", "purchase_id", "form", "expires_at", "remaining_count", "status", "created_at", "updated_at")
SELECT "id", "subject_id", "offer_id", "purchase_id", "form", CASE WHEN expires_at IS NULL THEN NULL ELSE strftime('%Y-%m-%dT%H:%M:%S', expires_at, 'unixepoch') || '.000000Z' END, "remaining_count", "status", strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z', strftime('%Y-%m-%dT%H:%M:%S', updated_at, 'unixepoch') || '.000000Z'
FROM "digital_entitlements_old"`,
	`DROP TABLE "digital_entitlements_old"`,
	`CREATE INDEX idx_digital_entitlements_purchase ON digital_entitlements(purchase_id)`,
	`CREATE INDEX idx_digital_entitlements_subject ON digital_entitlements(subject_id, status)`,
	`CREATE INDEX idx_digital_offers_status ON digital_offers(status, created_at DESC)`,
	`CREATE INDEX idx_digital_purchases_offer ON digital_purchases(offer_id)`,
	`CREATE INDEX idx_digital_purchases_subject ON digital_purchases(subject_id, created_at DESC)`,
}

// vp040V87Statements returns the canonical checksum input (D-017: m0 → m1–m3 → m4).
func vp040V87Statements() []string {
	return temporalmigrate.OrderedWithGuards(vp040V87Guards, vp040V87Preflight, vp040V87Rebuild, vp040V87Verify)
}

// applyvp040V87 is the SQLite Apply body.
func applyvp040V87(tx kernel.Tx) error {
	if err := temporalmigrate.RunGuards(tx, vp040V87Guards, "vp040 v87 sqlite guards"); err != nil {
		return err
	}
	if err := temporalmigrate.RunPreflight(tx, vp040V87Preflight); err != nil {
		return err
	}
	if err := temporalmigrate.Exec(tx, vp040V87Rebuild, "vp040 v87 sqlite rebuild"); err != nil {
		return err
	}
	return temporalmigrate.RunVerify(tx, vp040V87Verify, "vp040 v87 sqlite verify")
}

// vp040V87Postgres is the explicit postgres m1–m2 slice (D-019 §6: no regex
// derivation; the millisecond family uses the integer-split expression).
var vp040V87Postgres = []string{
	`ALTER TABLE "digital_offers" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)))`,
	`ALTER TABLE "digital_offers" ALTER COLUMN "updated_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("updated_at"::double precision)))`,
	`ALTER TABLE "digital_purchases" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)))`,
	`ALTER TABLE "digital_entitlements" ALTER COLUMN "expires_at" TYPE timestamptz(6) USING (CASE WHEN "expires_at" IS NULL THEN NULL ELSE date_trunc('microseconds', to_timestamp("expires_at"::double precision)) END)`,
	`ALTER TABLE "digital_entitlements" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)))`,
	`ALTER TABLE "digital_entitlements" ALTER COLUMN "updated_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("updated_at"::double precision)))`,
}

// vp040V87PostgresVerify is the postgres m4 type/precision assertion set.
var vp040V87PostgresVerify = []temporalmigrate.PgVerify{
	{Table: "digital_offers", Column: "created_at", NonNull: true},
	{Table: "digital_offers", Column: "updated_at", NonNull: true},
	{Table: "digital_purchases", Column: "created_at", NonNull: true},
	{Table: "digital_entitlements", Column: "expires_at", NonNull: false},
	{Table: "digital_entitlements", Column: "created_at", NonNull: true},
	{Table: "digital_entitlements", Column: "updated_at", NonNull: true},
}

// applyvp040V87Postgres is the postgres Apply body.
func applyvp040V87Postgres(tx kernel.Tx) error {
	if err := temporalmigrate.RunPostgresGuards(tx, vp040V87Guards, "vp040 v87 postgres guards"); err != nil {
		return err
	}
	if err := temporalmigrate.RunPreflight(tx, vp040V87Preflight); err != nil {
		return err
	}
	if err := temporalmigrate.Exec(tx, vp040V87Postgres, "vp040 v87 postgres convert"); err != nil {
		return err
	}
	return temporalmigrate.RunPostgresVerify(tx, vp040V87PostgresVerify, "vp040 v87 postgres verify")
}

// VP040TemporalDescriptor returns the workspace-040 R2 conversion descriptor
// (D-014 allocation, D-017 checksum convention).
func VP040TemporalDescriptor() kernel.MigrationContribution {
	return kernel.MigrationContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: vp040V87Name},
		Version:              87,
		Name:                 vp040V87Name,
		Checksum:             kernel.MigrationChecksum(vp040V87Statements(), vp040V87TransformID),
		Apply:                applyvp040V87,
		ApplyPostgres:        applyvp040V87Postgres,
	}
}
