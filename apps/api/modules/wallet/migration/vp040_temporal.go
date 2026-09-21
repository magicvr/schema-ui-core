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
// ModuleID admin.wallet · version 85 · D-014 allocation / D-017 transform id.
package migration

import (
	"github.com/magicvr/schema-ui-core/apps/api/internal/temporalmigrate"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

const (
	vp040V85Name        = "vp040_temporal_wallet"
	vp040V85TransformID = "0085:vp040-temporal-wallet:v1"
)

// vp040V85Guards are the m0 preconditions on tables retired by earlier history
// (descriptor ledger §1: v73 asserts the retired `records` table is absent).
var vp040V85Guards = []temporalmigrate.Guard{}

// vp040V85Preflight is the m0 sentinel census (Root D-012 / D-015 policy).
var vp040V85Preflight = []temporalmigrate.Preflight{
	{Table: "vouchers", Column: "expires_at", Voucher: true},
	{Table: "vouchers", Column: "redeemed_at", Voucher: true},
}

// vp040V85Verify is the m4 post-rebuild assertion set.
var vp040V85Verify = []temporalmigrate.Verify{
	{Table: "wallet_accounts", Columns: []string{"created_at", "updated_at"}},
	{Table: "wallet_ledger_entries", Columns: []string{"created_at"}},
	{Table: "wallet_reconciliation_runs", Columns: []string{"created_at"}},
	{Table: "subjects", Columns: []string{"created_at"}},
	{Table: "vouchers", Columns: []string{"expires_at", "redeemed_at", "created_at", "updated_at"}},
	{Table: "voucher_batches", Columns: []string{"created_at", "updated_at"}},
}

// vp040V85Rebuild is the ordered m1–m3 slice: constraint handling, table
// rebuild with the frozen conversion expressions, then every index.
var vp040V85Rebuild = []string{
	`ALTER TABLE "wallet_accounts" RENAME TO "wallet_accounts_old"`,
	`CREATE TABLE wallet_accounts (
  id                TEXT PRIMARY KEY,
  owner_type        TEXT NOT NULL CHECK (owner_type IN ('user','business','system','subject')),
  owner_id          TEXT NOT NULL,
  currency          TEXT NOT NULL DEFAULT 'CNY',
  balance_total     INTEGER NOT NULL DEFAULT 0 CHECK (balance_total >= 0),
  balance_available INTEGER NOT NULL DEFAULT 0 CHECK (balance_available >= 0),
  balance_frozen    INTEGER NOT NULL DEFAULT 0 CHECK (balance_frozen >= 0),
  status            TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active','disabled')),
  version           INTEGER NOT NULL DEFAULT 0,
  created_at        TEXT NOT NULL,
  updated_at        TEXT NOT NULL,
  UNIQUE (owner_type, owner_id, currency),
  CHECK (balance_total = balance_available + balance_frozen)
)`,
	`INSERT INTO "wallet_accounts" ("id", "owner_type", "owner_id", "currency", "balance_total", "balance_available", "balance_frozen", "status", "version", "created_at", "updated_at")
SELECT "id", "owner_type", "owner_id", "currency", "balance_total", "balance_available", "balance_frozen", "status", "version", strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z', strftime('%Y-%m-%dT%H:%M:%S', updated_at, 'unixepoch') || '.000000Z'
FROM "wallet_accounts_old"`,
	`DROP TABLE "wallet_accounts_old"`,
	`ALTER TABLE "wallet_ledger_entries" RENAME TO "wallet_ledger_entries_old"`,
	`CREATE TABLE wallet_ledger_entries (
  id                      TEXT PRIMARY KEY,
  account_id              TEXT NOT NULL,
  entry_type              TEXT NOT NULL CHECK (entry_type IN ('adjust','freeze','unfreeze','deduct_frozen')),
  amount_delta            INTEGER NOT NULL CHECK (amount_delta != 0),
  balance_after_total     INTEGER NOT NULL CHECK (balance_after_total >= 0),
  balance_after_available INTEGER NOT NULL CHECK (balance_after_available >= 0),
  balance_after_frozen    INTEGER NOT NULL CHECK (balance_after_frozen >= 0),
  ref_type                TEXT,
  ref_id                  TEXT,
  idempotency_key         TEXT,
  memo                    TEXT NOT NULL,
  actor_id                TEXT NOT NULL,
  actor_name              TEXT NOT NULL,
  created_at              TEXT NOT NULL,
  UNIQUE (account_id, idempotency_key),
  CHECK (balance_after_total = balance_after_available + balance_after_frozen)
)`,
	`INSERT INTO "wallet_ledger_entries" ("id", "account_id", "entry_type", "amount_delta", "balance_after_total", "balance_after_available", "balance_after_frozen", "ref_type", "ref_id", "idempotency_key", "memo", "actor_id", "actor_name", "created_at")
SELECT "id", "account_id", "entry_type", "amount_delta", "balance_after_total", "balance_after_available", "balance_after_frozen", "ref_type", "ref_id", "idempotency_key", "memo", "actor_id", "actor_name", strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z'
FROM "wallet_ledger_entries_old"`,
	`DROP TABLE "wallet_ledger_entries_old"`,
	`ALTER TABLE "wallet_reconciliation_runs" RENAME TO "wallet_reconciliation_runs_old"`,
	`CREATE TABLE wallet_reconciliation_runs (
  id             TEXT PRIMARY KEY,
  account_id     TEXT,
  result         TEXT NOT NULL CHECK (result IN ('consistent','inconsistent')),
  mismatch_count INTEGER NOT NULL DEFAULT 0,
  details        TEXT NOT NULL DEFAULT '{}',
  actor_id       TEXT NOT NULL,
  created_at     TEXT NOT NULL
)`,
	`INSERT INTO "wallet_reconciliation_runs" ("id", "account_id", "result", "mismatch_count", "details", "actor_id", "created_at")
SELECT "id", "account_id", "result", "mismatch_count", "details", "actor_id", strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z'
FROM "wallet_reconciliation_runs_old"`,
	`DROP TABLE "wallet_reconciliation_runs_old"`,
	`ALTER TABLE "subjects" RENAME TO "subjects_old"`,
	`CREATE TABLE subjects (
  id          TEXT PRIMARY KEY,
  issuer      TEXT NOT NULL,
  external_id TEXT NOT NULL,
  created_at  TEXT NOT NULL,
  UNIQUE (issuer, external_id)
)`,
	`INSERT INTO "subjects" ("id", "issuer", "external_id", "created_at")
SELECT "id", "issuer", "external_id", strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z'
FROM "subjects_old"`,
	`DROP TABLE "subjects_old"`,
	`ALTER TABLE "vouchers" RENAME TO "vouchers_old"`,
	`CREATE TABLE vouchers (
  id           TEXT PRIMARY KEY,
  batch_id     TEXT NOT NULL,
  code_hash    TEXT NOT NULL,
  code_prefix  TEXT NOT NULL,
  amount       INTEGER NOT NULL CHECK (amount > 0),
  currency     TEXT NOT NULL DEFAULT 'CNY',
  status       TEXT NOT NULL DEFAULT 'unused' CHECK (status IN ('unused','redeemed','void')),
  expires_at   TEXT,
  redeemed_by  TEXT,
  redeemed_at  TEXT,
  created_at   TEXT NOT NULL,
  updated_at   TEXT NOT NULL,
  UNIQUE (code_hash)
)`,
	`INSERT INTO "vouchers" ("id", "batch_id", "code_hash", "code_prefix", "amount", "currency", "status", "expires_at", "redeemed_by", "redeemed_at", "created_at", "updated_at")
SELECT "id", "batch_id", "code_hash", "code_prefix", "amount", "currency", "status", CASE WHEN expires_at IS NULL OR expires_at = 0 THEN NULL ELSE strftime('%Y-%m-%dT%H:%M:%S', expires_at, 'unixepoch') || '.000000Z' END, "redeemed_by", CASE WHEN redeemed_at IS NULL OR redeemed_at = 0 THEN NULL ELSE strftime('%Y-%m-%dT%H:%M:%S', redeemed_at, 'unixepoch') || '.000000Z' END, strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z', strftime('%Y-%m-%dT%H:%M:%S', updated_at, 'unixepoch') || '.000000Z'
FROM "vouchers_old"`,
	`DROP TABLE "vouchers_old"`,
	`ALTER TABLE "voucher_batches" RENAME TO "voucher_batches_old"`,
	`CREATE TABLE voucher_batches (
  batch_id   TEXT PRIMARY KEY,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
)`,
	`INSERT INTO "voucher_batches" ("batch_id", "created_at", "updated_at")
SELECT "batch_id", strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z', strftime('%Y-%m-%dT%H:%M:%S', updated_at, 'unixepoch') || '.000000Z'
FROM "voucher_batches_old"`,
	`DROP TABLE "voucher_batches_old"`,
	`CREATE INDEX idx_subjects_issuer_external ON subjects(issuer, external_id)`,
	`CREATE INDEX idx_vouchers_batch ON vouchers(batch_id, created_at DESC)`,
	`CREATE INDEX idx_vouchers_status ON vouchers(status)`,
	`CREATE INDEX idx_wallet_ledger_account ON wallet_ledger_entries(account_id, created_at DESC)`,
}

// vp040V85Statements returns the canonical checksum input (D-017: m0 → m1–m3 → m4).
func vp040V85Statements() []string {
	return temporalmigrate.OrderedWithGuards(vp040V85Guards, vp040V85Preflight, vp040V85Rebuild, vp040V85Verify)
}

// applyvp040V85 is the SQLite Apply body.
func applyvp040V85(tx kernel.Tx) error {
	if err := temporalmigrate.RunGuards(tx, vp040V85Guards, "vp040 v85 sqlite guards"); err != nil {
		return err
	}
	if err := temporalmigrate.RunPreflight(tx, vp040V85Preflight); err != nil {
		return err
	}
	if err := temporalmigrate.Exec(tx, vp040V85Rebuild, "vp040 v85 sqlite rebuild"); err != nil {
		return err
	}
	return temporalmigrate.RunVerify(tx, vp040V85Verify, "vp040 v85 sqlite verify")
}

// vp040V85Postgres is the explicit postgres m1–m2 slice (D-019 §6: no regex
// derivation; the millisecond family uses the integer-split expression).
var vp040V85Postgres = []string{
	`ALTER TABLE "wallet_accounts" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)))`,
	`ALTER TABLE "wallet_accounts" ALTER COLUMN "updated_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("updated_at"::double precision)))`,
	`ALTER TABLE "wallet_ledger_entries" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)))`,
	`ALTER TABLE "wallet_reconciliation_runs" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)))`,
	`ALTER TABLE "subjects" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)))`,
	`ALTER TABLE "vouchers" ALTER COLUMN "expires_at" TYPE timestamptz(6) USING (CASE WHEN "expires_at" IS NULL OR "expires_at" = 0 THEN NULL ELSE date_trunc('microseconds', to_timestamp("expires_at"::double precision)) END)`,
	`ALTER TABLE "vouchers" ALTER COLUMN "redeemed_at" TYPE timestamptz(6) USING (CASE WHEN "redeemed_at" IS NULL OR "redeemed_at" = 0 THEN NULL ELSE date_trunc('microseconds', to_timestamp("redeemed_at"::double precision)) END)`,
	`ALTER TABLE "vouchers" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)))`,
	`ALTER TABLE "vouchers" ALTER COLUMN "updated_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("updated_at"::double precision)))`,
	`ALTER TABLE "voucher_batches" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)))`,
	`ALTER TABLE "voucher_batches" ALTER COLUMN "updated_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("updated_at"::double precision)))`,
}

// vp040V85PostgresVerify is the postgres m4 type/precision assertion set.
var vp040V85PostgresVerify = []temporalmigrate.PgVerify{
	{Table: "wallet_accounts", Column: "created_at", NonNull: true},
	{Table: "wallet_accounts", Column: "updated_at", NonNull: true},
	{Table: "wallet_ledger_entries", Column: "created_at", NonNull: true},
	{Table: "wallet_reconciliation_runs", Column: "created_at", NonNull: true},
	{Table: "subjects", Column: "created_at", NonNull: true},
	{Table: "vouchers", Column: "expires_at", NonNull: false},
	{Table: "vouchers", Column: "redeemed_at", NonNull: false},
	{Table: "vouchers", Column: "created_at", NonNull: true},
	{Table: "vouchers", Column: "updated_at", NonNull: true},
	{Table: "voucher_batches", Column: "created_at", NonNull: true},
	{Table: "voucher_batches", Column: "updated_at", NonNull: true},
}

// applyvp040V85Postgres is the postgres Apply body.
func applyvp040V85Postgres(tx kernel.Tx) error {
	if err := temporalmigrate.RunPostgresGuards(tx, vp040V85Guards, "vp040 v85 postgres guards"); err != nil {
		return err
	}
	if err := temporalmigrate.RunPreflight(tx, vp040V85Preflight); err != nil {
		return err
	}
	if err := temporalmigrate.Exec(tx, vp040V85Postgres, "vp040 v85 postgres convert"); err != nil {
		return err
	}
	return temporalmigrate.RunPostgresVerify(tx, vp040V85PostgresVerify, "vp040 v85 postgres verify")
}

// VP040TemporalDescriptor returns the workspace-040 R2 conversion descriptor
// (D-014 allocation, D-017 checksum convention).
func VP040TemporalDescriptor() kernel.MigrationContribution {
	return kernel.MigrationContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: vp040V85Name},
		Version:              85,
		Name:                 vp040V85Name,
		Checksum:             kernel.MigrationChecksum(vp040V85Statements(), vp040V85TransformID),
		Apply:                applyvp040V85,
		ApplyPostgres:        applyvp040V85Postgres,
	}
}
