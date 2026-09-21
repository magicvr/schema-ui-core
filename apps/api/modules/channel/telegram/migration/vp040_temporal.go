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
// ModuleID channel.telegram · version 86 · D-014 allocation / D-017 transform id.
package migration

import (
	"github.com/magicvr/schema-ui-core/apps/api/internal/temporalmigrate"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

const (
	vp040V86Name        = "vp040_temporal_telegram"
	vp040V86TransformID = "0086:vp040-temporal-telegram:v1"
)

// vp040V86Guards are the m0 preconditions on tables retired by earlier history
// (descriptor ledger §1: v73 asserts the retired `records` table is absent).
var vp040V86Guards = []temporalmigrate.Guard{}

// vp040V86Preflight is the m0 sentinel census (Root D-012 / D-015 policy).
var vp040V86Preflight = []temporalmigrate.Preflight{
	{Table: "telegram_config", Column: "updated_at", Voucher: false},
}

// vp040V86Verify is the m4 post-rebuild assertion set.
var vp040V86Verify = []temporalmigrate.Verify{
	{Table: "telegram_config", Columns: []string{"updated_at"}},
	{Table: "telegram_sessions", Columns: []string{"last_message_at", "created_at", "updated_at"}},
	{Table: "telegram_inbound_messages", Columns: []string{"received_at"}},
	{Table: "telegram_outbound_messages", Columns: []string{"created_at", "updated_at"}},
}

// vp040V86Rebuild is the ordered m1–m3 slice: constraint handling, table
// rebuild with the frozen conversion expressions, then every index.
var vp040V86Rebuild = []string{
	`ALTER TABLE "telegram_config" RENAME TO "telegram_config_old"`,
	`CREATE TABLE telegram_config (
  id                 INTEGER PRIMARY KEY CHECK (id = 1),
  bot_token_enc      TEXT    NOT NULL DEFAULT '',
  webhook_secret_enc TEXT    NOT NULL DEFAULT '',
  updated_at         TEXT
, mode TEXT NOT NULL DEFAULT 'polling', webhook_public_base_url TEXT NOT NULL DEFAULT '')`,
	`INSERT INTO "telegram_config" ("id", "bot_token_enc", "webhook_secret_enc", "updated_at", "mode", "webhook_public_base_url")
SELECT "id", "bot_token_enc", "webhook_secret_enc", CASE WHEN updated_at = 0 THEN NULL ELSE strftime('%Y-%m-%dT%H:%M:%S', updated_at, 'unixepoch') || '.000000Z' END, "mode", "webhook_public_base_url"
FROM "telegram_config_old"`,
	`DROP TABLE "telegram_config_old"`,
	`ALTER TABLE "telegram_sessions" RENAME TO "telegram_sessions_old"`,
	`CREATE TABLE telegram_sessions (
  bot_id          INTEGER NOT NULL,
  chat_id         INTEGER NOT NULL,
  chat_type       TEXT    NOT NULL DEFAULT '',
  title           TEXT    NOT NULL DEFAULT '',
  username        TEXT    NOT NULL DEFAULT '',
  last_message_at TEXT NOT NULL,
  created_at      TEXT NOT NULL,
  updated_at      TEXT NOT NULL,
  PRIMARY KEY (bot_id, chat_id)
)`,
	`INSERT INTO "telegram_sessions" ("bot_id", "chat_id", "chat_type", "title", "username", "last_message_at", "created_at", "updated_at")
SELECT "bot_id", "chat_id", "chat_type", "title", "username", strftime('%Y-%m-%dT%H:%M:%S', last_message_at, 'unixepoch') || '.000000Z', strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z', strftime('%Y-%m-%dT%H:%M:%S', updated_at, 'unixepoch') || '.000000Z'
FROM "telegram_sessions_old"`,
	`DROP TABLE "telegram_sessions_old"`,
	`ALTER TABLE "telegram_inbound_messages" RENAME TO "telegram_inbound_messages_old"`,
	`CREATE TABLE telegram_inbound_messages (
  bot_id            INTEGER NOT NULL,
  update_id         INTEGER NOT NULL,
  chat_id           INTEGER NOT NULL,
  user_id           INTEGER,
  message_id        INTEGER,
  callback_query_id TEXT,
  direction         TEXT    NOT NULL DEFAULT 'inbound' CHECK (direction = 'inbound'),
  message_kind      TEXT    NOT NULL,
  text              TEXT,
  callback_data     TEXT,
  sender_username   TEXT,
  received_at       TEXT NOT NULL,
  PRIMARY KEY (bot_id, update_id)
)`,
	`INSERT INTO "telegram_inbound_messages" ("bot_id", "update_id", "chat_id", "user_id", "message_id", "callback_query_id", "direction", "message_kind", "text", "callback_data", "sender_username", "received_at")
SELECT "bot_id", "update_id", "chat_id", "user_id", "message_id", "callback_query_id", "direction", "message_kind", "text", "callback_data", "sender_username", strftime('%Y-%m-%dT%H:%M:%S', received_at, 'unixepoch') || '.000000Z'
FROM "telegram_inbound_messages_old"`,
	`DROP TABLE "telegram_inbound_messages_old"`,
	`ALTER TABLE "telegram_outbound_messages" RENAME TO "telegram_outbound_messages_old"`,
	`CREATE TABLE telegram_outbound_messages (
  bot_id        INTEGER NOT NULL,
  request_id    TEXT    NOT NULL,
  retry_root    TEXT    NOT NULL,
  retry_of      TEXT,
  chat_id       INTEGER NOT NULL,
  text          TEXT    NOT NULL,
  status        TEXT    NOT NULL CHECK (status IN ('pending', 'sent', 'failed')),
  error_message TEXT,
  created_at    TEXT NOT NULL,
  updated_at    TEXT NOT NULL,
  PRIMARY KEY (bot_id, request_id)
 )`,
	`INSERT INTO "telegram_outbound_messages" ("bot_id", "request_id", "retry_root", "retry_of", "chat_id", "text", "status", "error_message", "created_at", "updated_at")
SELECT "bot_id", "request_id", "retry_root", "retry_of", "chat_id", "text", "status", "error_message", strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z', strftime('%Y-%m-%dT%H:%M:%S', updated_at, 'unixepoch') || '.000000Z'
FROM "telegram_outbound_messages_old"`,
	`DROP TABLE "telegram_outbound_messages_old"`,
	`CREATE INDEX idx_telegram_inbound_messages_chat_received
  ON telegram_inbound_messages (bot_id, chat_id, received_at DESC, update_id DESC)`,
	`CREATE INDEX idx_telegram_outbound_messages_chat_created
  ON telegram_outbound_messages (bot_id, chat_id, created_at DESC, request_id DESC)`,
	`CREATE UNIQUE INDEX idx_telegram_outbound_messages_pending_root
  ON telegram_outbound_messages (bot_id, retry_root) WHERE status = 'pending'`,
	`CREATE INDEX idx_telegram_sessions_activity
  ON telegram_sessions (bot_id, last_message_at DESC, chat_id DESC)`,
}

// vp040V86Statements returns the canonical checksum input (D-017: m0 → m1–m3 → m4).
func vp040V86Statements() []string {
	return temporalmigrate.OrderedWithGuards(vp040V86Guards, vp040V86Preflight, vp040V86Rebuild, vp040V86Verify)
}

// applyvp040V86 is the SQLite Apply body.
func applyvp040V86(tx kernel.Tx) error {
	if err := temporalmigrate.RunGuards(tx, vp040V86Guards, "vp040 v86 sqlite guards"); err != nil {
		return err
	}
	if err := temporalmigrate.RunPreflight(tx, vp040V86Preflight); err != nil {
		return err
	}
	if err := temporalmigrate.Exec(tx, vp040V86Rebuild, "vp040 v86 sqlite rebuild"); err != nil {
		return err
	}
	return temporalmigrate.RunVerify(tx, vp040V86Verify, "vp040 v86 sqlite verify")
}

// vp040V86Postgres is the explicit postgres m1–m2 slice (D-019 §6: no regex
// derivation; the millisecond family uses the integer-split expression).
var vp040V86Postgres = []string{
	`ALTER TABLE "telegram_config" ALTER COLUMN "updated_at" DROP DEFAULT`,
	`ALTER TABLE "telegram_config" ALTER COLUMN "updated_at" DROP NOT NULL`,
	`ALTER TABLE "telegram_config" ALTER COLUMN "updated_at" TYPE timestamptz(6) USING (CASE WHEN "updated_at" = 0 THEN NULL ELSE date_trunc('microseconds', to_timestamp("updated_at"::double precision)) END)`,
	`ALTER TABLE "telegram_sessions" ALTER COLUMN "last_message_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("last_message_at"::double precision)))`,
	`ALTER TABLE "telegram_sessions" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)))`,
	`ALTER TABLE "telegram_sessions" ALTER COLUMN "updated_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("updated_at"::double precision)))`,
	`ALTER TABLE "telegram_inbound_messages" ALTER COLUMN "received_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("received_at"::double precision)))`,
	`ALTER TABLE "telegram_outbound_messages" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)))`,
	`ALTER TABLE "telegram_outbound_messages" ALTER COLUMN "updated_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("updated_at"::double precision)))`,
}

// vp040V86PostgresVerify is the postgres m4 type/precision assertion set.
var vp040V86PostgresVerify = []temporalmigrate.PgVerify{
	{Table: "telegram_config", Column: "updated_at", NonNull: false},
	{Table: "telegram_sessions", Column: "last_message_at", NonNull: true},
	{Table: "telegram_sessions", Column: "created_at", NonNull: true},
	{Table: "telegram_sessions", Column: "updated_at", NonNull: true},
	{Table: "telegram_inbound_messages", Column: "received_at", NonNull: true},
	{Table: "telegram_outbound_messages", Column: "created_at", NonNull: true},
	{Table: "telegram_outbound_messages", Column: "updated_at", NonNull: true},
}

// applyvp040V86Postgres is the postgres Apply body.
func applyvp040V86Postgres(tx kernel.Tx) error {
	if err := temporalmigrate.RunPostgresGuards(tx, vp040V86Guards, "vp040 v86 postgres guards"); err != nil {
		return err
	}
	if err := temporalmigrate.RunPreflight(tx, vp040V86Preflight); err != nil {
		return err
	}
	if err := temporalmigrate.Exec(tx, vp040V86Postgres, "vp040 v86 postgres convert"); err != nil {
		return err
	}
	return temporalmigrate.RunPostgresVerify(tx, vp040V86PostgresVerify, "vp040 v86 postgres verify")
}

// VP040TemporalDescriptor returns the workspace-040 R2 conversion descriptor
// (D-014 allocation, D-017 checksum convention).
func VP040TemporalDescriptor() kernel.MigrationContribution {
	return kernel.MigrationContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: vp040V86Name},
		Version:              86,
		Name:                 vp040V86Name,
		Checksum:             kernel.MigrationChecksum(vp040V86Statements(), vp040V86TransformID),
		Apply:                applyvp040V86,
		ApplyPostgres:        applyvp040V86Postgres,
	}
}
