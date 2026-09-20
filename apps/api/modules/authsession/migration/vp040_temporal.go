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
// ModuleID core.auth-session · version 74 · D-014 allocation / D-017 transform id.
package migration

import (
	"github.com/magicvr/schema-ui-core/apps/api/internal/temporalmigrate"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

const (
	vp040V74Name        = "vp040_temporal_authsession"
	vp040V74TransformID = "0074:vp040-temporal-authsession:v1"
)

// vp040V74Preflight is the m0 sentinel census (Root D-012 / D-015 policy).
var vp040V74Preflight = []temporalmigrate.Preflight{
	{Table: "users", Column: "locked_until", Voucher: false},
	{Table: "users", Column: "last_login_failure_at", Voucher: false},
	{Table: "login_failures", Column: "locked_until", Voucher: false},
}

// vp040V74Verify is the m4 post-rebuild assertion set.
var vp040V74Verify = []temporalmigrate.Verify{
	{Table: "system_data_reconcile", Columns: []string{"applied_at"}},
	{Table: "users", Columns: []string{"created_at", "updated_at", "locked_until", "last_login_failure_at"}, Children: []string{"email_verification_challenges", "login_failures", "mfa_proofs", "notifications", "password_recovery_challenges", "refresh_tokens", "user_invites", "user_mfa", "user_password_history", "user_roles"}},
	{Table: "roles", Columns: []string{"created_at", "updated_at"}, Children: []string{"role_menu_items", "role_permissions", "user_roles"}},
	{Table: "permissions", Columns: []string{"created_at", "updated_at"}, Children: []string{"role_permissions"}},
	{Table: "menu_items", Columns: []string{"created_at", "updated_at"}, Children: []string{"role_menu_items"}},
	{Table: "refresh_tokens", Columns: []string{"expires_at", "revoked_at", "created_at"}},
	{Table: "email_verification_challenges", Columns: []string{"expires_at", "sent_at"}},
	{Table: "password_recovery_challenges", Columns: []string{"expires_at", "sent_at"}},
	{Table: "login_failures", Columns: []string{"locked_until", "updated_at"}},
	{Table: "user_password_history", Columns: []string{"created_at"}},
	{Table: "user_invites", Columns: []string{"expires_at", "consumed_at", "revoked_at", "last_sent_at", "created_at"}},
	{Table: "service_credentials", Columns: []string{"expires_at", "revoked_at", "last_used_at", "created_at", "updated_at"}},
}

// vp040V74Rebuild is the ordered m1–m3 slice: constraint handling, table
// rebuild with the frozen conversion expressions, then every index.
var vp040V74Rebuild = []string{
	`CREATE TEMP TABLE "email_verification_challenges_bak" AS SELECT * FROM "email_verification_challenges"`,
	`CREATE TEMP TABLE "login_failures_bak" AS SELECT * FROM "login_failures"`,
	`CREATE TEMP TABLE "mfa_proofs_bak" AS SELECT * FROM "mfa_proofs"`,
	`CREATE TEMP TABLE "notifications_bak" AS SELECT * FROM "notifications"`,
	`CREATE TEMP TABLE "password_recovery_challenges_bak" AS SELECT * FROM "password_recovery_challenges"`,
	`CREATE TEMP TABLE "refresh_tokens_bak" AS SELECT * FROM "refresh_tokens"`,
	`CREATE TEMP TABLE "role_menu_items_bak" AS SELECT * FROM "role_menu_items"`,
	`CREATE TEMP TABLE "role_permissions_bak" AS SELECT * FROM "role_permissions"`,
	`CREATE TEMP TABLE "user_invites_bak" AS SELECT * FROM "user_invites"`,
	`CREATE TEMP TABLE "user_mfa_bak" AS SELECT * FROM "user_mfa"`,
	`CREATE TEMP TABLE "user_password_history_bak" AS SELECT * FROM "user_password_history"`,
	`CREATE TEMP TABLE "user_roles_bak" AS SELECT * FROM "user_roles"`,
	`DROP TABLE "email_verification_challenges"`,
	`DROP TABLE "login_failures"`,
	`DROP TABLE "mfa_proofs"`,
	`DROP TABLE "notifications"`,
	`DROP TABLE "password_recovery_challenges"`,
	`DROP TABLE "refresh_tokens"`,
	`DROP TABLE "role_menu_items"`,
	`DROP TABLE "role_permissions"`,
	`DROP TABLE "user_invites"`,
	`DROP TABLE "user_mfa"`,
	`DROP TABLE "user_password_history"`,
	`DROP TABLE "user_roles"`,
	`ALTER TABLE "users" RENAME TO "users_old"`,
	`CREATE TABLE users (
  id            TEXT PRIMARY KEY,
  username      TEXT NOT NULL UNIQUE,
  name          TEXT NOT NULL,
  roles         TEXT NOT NULL, -- JSON array; R3 normalizes
  password_hash TEXT NOT NULL,
  created_at    TEXT NOT NULL,
  updated_at    TEXT NOT NULL
, token_version INTEGER NOT NULL DEFAULT 0, failed_login_count INTEGER NOT NULL DEFAULT 0, locked_until TEXT, enabled INTEGER NOT NULL DEFAULT 1, notifications_enabled INTEGER NOT NULL DEFAULT 1, avatar_url TEXT NOT NULL DEFAULT '', must_change_password INTEGER NOT NULL DEFAULT 0, email TEXT, email_status TEXT CHECK (email_status IN ('pending','verified')), last_login_failure_at TEXT)`,
	`INSERT INTO "users" ("id", "username", "name", "roles", "password_hash", "created_at", "updated_at", "token_version", "failed_login_count", "locked_until", "enabled", "notifications_enabled", "avatar_url", "must_change_password", "email", "email_status", "last_login_failure_at")
SELECT "id", "username", "name", "roles", "password_hash", strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z', strftime('%Y-%m-%dT%H:%M:%S', updated_at, 'unixepoch') || '.000000Z', "token_version", "failed_login_count", CASE WHEN locked_until = 0 THEN NULL ELSE strftime('%Y-%m-%dT%H:%M:%S', locked_until, 'unixepoch') || '.000000Z' END, "enabled", "notifications_enabled", "avatar_url", "must_change_password", "email", "email_status", CASE WHEN last_login_failure_at = 0 THEN NULL ELSE strftime('%Y-%m-%dT%H:%M:%S', last_login_failure_at, 'unixepoch') || '.000000Z' END
FROM "users_old"`,
	`DROP TABLE "users_old"`,
	`ALTER TABLE "roles" RENAME TO "roles_old"`,
	`CREATE TABLE roles (
  id         TEXT PRIMARY KEY,
  key        TEXT NOT NULL UNIQUE CHECK (key <> ''),
  name       TEXT NOT NULL,
  system     INTEGER NOT NULL DEFAULT 0 CHECK (system IN (0, 1)),
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
)`,
	`INSERT INTO "roles" ("id", "key", "name", "system", "created_at", "updated_at")
SELECT "id", "key", "name", "system", strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z', strftime('%Y-%m-%dT%H:%M:%S', updated_at, 'unixepoch') || '.000000Z'
FROM "roles_old"`,
	`DROP TABLE "roles_old"`,
	`ALTER TABLE "permissions" RENAME TO "permissions_old"`,
	`CREATE TABLE permissions (
  id          TEXT PRIMARY KEY,
  key         TEXT NOT NULL UNIQUE CHECK (key <> ''),
  description TEXT NOT NULL DEFAULT '',
  created_at  TEXT NOT NULL,
  updated_at  TEXT NOT NULL
)`,
	`INSERT INTO "permissions" ("id", "key", "description", "created_at", "updated_at")
SELECT "id", "key", "description", strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z', strftime('%Y-%m-%dT%H:%M:%S', updated_at, 'unixepoch') || '.000000Z'
FROM "permissions_old"`,
	`DROP TABLE "permissions_old"`,
	`ALTER TABLE "menu_items" RENAME TO "menu_items_old"`,
	`CREATE TABLE menu_items (
  id          TEXT PRIMARY KEY,
  page_ref    TEXT NOT NULL UNIQUE CHECK (page_ref <> ''),
  feature_key TEXT NOT NULL UNIQUE CHECK (feature_key <> ''),
  sort_order  INTEGER NOT NULL DEFAULT 0,
  enabled     INTEGER NOT NULL DEFAULT 1 CHECK (enabled IN (0, 1)),
  created_at  TEXT NOT NULL,
  updated_at  TEXT NOT NULL
)`,
	`INSERT INTO "menu_items" ("id", "page_ref", "feature_key", "sort_order", "enabled", "created_at", "updated_at")
SELECT "id", "page_ref", "feature_key", "sort_order", "enabled", strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z', strftime('%Y-%m-%dT%H:%M:%S', updated_at, 'unixepoch') || '.000000Z'
FROM "menu_items_old"`,
	`DROP TABLE "menu_items_old"`,
	`CREATE TABLE refresh_tokens (
  id         TEXT PRIMARY KEY,
  user_id    TEXT NOT NULL REFERENCES users(id),
  token_hash TEXT NOT NULL UNIQUE,
  expires_at TEXT NOT NULL,
  revoked_at TEXT,
  created_at TEXT NOT NULL
)`,
	`INSERT INTO "refresh_tokens" ("id", "user_id", "token_hash", "expires_at", "revoked_at", "created_at")
SELECT "id", "user_id", "token_hash", strftime('%Y-%m-%dT%H:%M:%S', expires_at, 'unixepoch') || '.000000Z', CASE WHEN revoked_at IS NULL THEN NULL ELSE strftime('%Y-%m-%dT%H:%M:%S', revoked_at, 'unixepoch') || '.000000Z' END, strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z'
FROM "refresh_tokens_bak"`,
	`DROP TABLE "refresh_tokens_bak"`,
	`CREATE TABLE email_verification_challenges (
  user_id       TEXT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  code_hash     TEXT NOT NULL,
  expires_at    TEXT NOT NULL,
  sent_at       TEXT NOT NULL,
  attempt_count INTEGER NOT NULL DEFAULT 0
)`,
	`INSERT INTO "email_verification_challenges" ("user_id", "code_hash", "expires_at", "sent_at", "attempt_count")
SELECT "user_id", "code_hash", strftime('%Y-%m-%dT%H:%M:%S', expires_at, 'unixepoch') || '.000000Z', strftime('%Y-%m-%dT%H:%M:%S', sent_at, 'unixepoch') || '.000000Z', "attempt_count"
FROM "email_verification_challenges_bak"`,
	`DROP TABLE "email_verification_challenges_bak"`,
	`CREATE TABLE password_recovery_challenges (
  user_id       TEXT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  code_hash     TEXT NOT NULL,
  expires_at    TEXT NOT NULL,
  sent_at       TEXT NOT NULL,
  attempt_count INTEGER NOT NULL DEFAULT 0
)`,
	`INSERT INTO "password_recovery_challenges" ("user_id", "code_hash", "expires_at", "sent_at", "attempt_count")
SELECT "user_id", "code_hash", strftime('%Y-%m-%dT%H:%M:%S', expires_at, 'unixepoch') || '.000000Z', strftime('%Y-%m-%dT%H:%M:%S', sent_at, 'unixepoch') || '.000000Z', "attempt_count"
FROM "password_recovery_challenges_bak"`,
	`DROP TABLE "password_recovery_challenges_bak"`,
	`CREATE TABLE login_failures (
  user_id      TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  ip          TEXT NOT NULL,
  fail_count  INTEGER NOT NULL DEFAULT 0,
  locked_until TEXT,
  updated_at  TEXT NOT NULL,
  PRIMARY KEY (user_id, ip)
)`,
	`INSERT INTO "login_failures" ("user_id", "ip", "fail_count", "locked_until", "updated_at")
SELECT "user_id", "ip", "fail_count", CASE WHEN locked_until = 0 THEN NULL ELSE strftime('%Y-%m-%dT%H:%M:%S', locked_until, 'unixepoch') || '.000000Z' END, strftime('%Y-%m-%dT%H:%M:%S', updated_at, 'unixepoch') || '.000000Z'
FROM "login_failures_bak"`,
	`DROP TABLE "login_failures_bak"`,
	`CREATE TABLE user_password_history (
  id            TEXT PRIMARY KEY,
  user_id       TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  password_hash TEXT NOT NULL,
  created_at    TEXT NOT NULL
)`,
	`INSERT INTO "user_password_history" ("id", "user_id", "password_hash", "created_at")
SELECT "id", "user_id", "password_hash", strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z'
FROM "user_password_history_bak"`,
	`DROP TABLE "user_password_history_bak"`,
	`CREATE TABLE user_invites (
  id           TEXT PRIMARY KEY,
  token_hash   TEXT NOT NULL UNIQUE,
  roles        TEXT NOT NULL,
  invited_by   TEXT NOT NULL REFERENCES users(id),
  email        TEXT,
  expires_at   TEXT NOT NULL,
  consumed_at  TEXT,
  revoked_at   TEXT,
  last_sent_at TEXT NOT NULL,
  created_at   TEXT NOT NULL
)`,
	`INSERT INTO "user_invites" ("id", "token_hash", "roles", "invited_by", "email", "expires_at", "consumed_at", "revoked_at", "last_sent_at", "created_at")
SELECT "id", "token_hash", "roles", "invited_by", "email", strftime('%Y-%m-%dT%H:%M:%S', expires_at, 'unixepoch') || '.000000Z', CASE WHEN consumed_at IS NULL THEN NULL ELSE strftime('%Y-%m-%dT%H:%M:%S', consumed_at, 'unixepoch') || '.000000Z' END, CASE WHEN revoked_at IS NULL THEN NULL ELSE strftime('%Y-%m-%dT%H:%M:%S', revoked_at, 'unixepoch') || '.000000Z' END, strftime('%Y-%m-%dT%H:%M:%S', last_sent_at, 'unixepoch') || '.000000Z', strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z'
FROM "user_invites_bak"`,
	`DROP TABLE "user_invites_bak"`,
	`ALTER TABLE "system_data_reconcile" RENAME TO "system_data_reconcile_old"`,
	`CREATE TABLE system_data_reconcile (
  module_id        TEXT NOT NULL,
  kind             TEXT NOT NULL CHECK (kind IN ('base','authorization','navigation')),
  contribution_key TEXT NOT NULL,
  version          INTEGER NOT NULL CHECK (version > 0),
  checksum         TEXT NOT NULL CHECK (length(checksum) = 64),
  applied_at       TEXT NOT NULL,
  PRIMARY KEY (module_id, kind, contribution_key)
)`,
	`INSERT INTO "system_data_reconcile" ("module_id", "kind", "contribution_key", "version", "checksum", "applied_at")
SELECT "module_id", "kind", "contribution_key", "version", "checksum", strftime('%Y-%m-%dT%H:%M:%S', applied_at, 'unixepoch') || '.000000Z'
FROM "system_data_reconcile_old"`,
	`DROP TABLE "system_data_reconcile_old"`,
	`ALTER TABLE "service_credentials" RENAME TO "service_credentials_old"`,
	`CREATE TABLE service_credentials (
  id           TEXT PRIMARY KEY CHECK (length(id) = 32),
  name         TEXT NOT NULL COLLATE NOCASE UNIQUE CHECK (length(trim(name)) BETWEEN 1 AND 100),
  token_prefix TEXT NOT NULL CHECK (length(token_prefix) = 15),
  token_hash   TEXT NOT NULL UNIQUE CHECK (length(token_hash) = 64),
  scopes       TEXT NOT NULL,
  expires_at   TEXT NOT NULL,
  revoked_at   TEXT,
  last_used_at TEXT,
  created_by   TEXT NOT NULL,
  created_at   TEXT NOT NULL,
  updated_at   TEXT NOT NULL
)`,
	`INSERT INTO "service_credentials" ("id", "name", "token_prefix", "token_hash", "scopes", "expires_at", "revoked_at", "last_used_at", "created_by", "created_at", "updated_at")
SELECT "id", "name", "token_prefix", "token_hash", "scopes", strftime('%Y-%m-%dT%H:%M:%S', expires_at, 'unixepoch') || '.000000Z', CASE WHEN revoked_at IS NULL THEN NULL ELSE strftime('%Y-%m-%dT%H:%M:%S', revoked_at, 'unixepoch') || '.000000Z' END, CASE WHEN last_used_at IS NULL THEN NULL ELSE strftime('%Y-%m-%dT%H:%M:%S', last_used_at, 'unixepoch') || '.000000Z' END, "created_by", strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z', strftime('%Y-%m-%dT%H:%M:%S', updated_at, 'unixepoch') || '.000000Z'
FROM "service_credentials_old"`,
	`DROP TABLE "service_credentials_old"`,
	`CREATE TABLE mfa_proofs (
  id         TEXT PRIMARY KEY,
  user_id    TEXT NOT NULL REFERENCES users(id),
  fail_count INTEGER NOT NULL DEFAULT 0,
  expires_at INTEGER NOT NULL,
  created_at INTEGER NOT NULL
)`,
	`INSERT INTO "mfa_proofs" ("id", "user_id", "fail_count", "expires_at", "created_at")
SELECT "id", "user_id", "fail_count", "expires_at", "created_at"
FROM "mfa_proofs_bak"`,
	`DROP TABLE "mfa_proofs_bak"`,
	`CREATE TABLE notifications (
  id         TEXT PRIMARY KEY,
  user_id    TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  event      TEXT NOT NULL CHECK (event IN ('account.locked','account.disabled','account.unlocked','account.password-changed')),
  title      TEXT NOT NULL,
  body       TEXT NOT NULL,
  read_at    INTEGER,
  created_at INTEGER NOT NULL
, title_key TEXT, body_key TEXT)`,
	`INSERT INTO "notifications" ("id", "user_id", "event", "title", "body", "read_at", "created_at", "title_key", "body_key")
SELECT "id", "user_id", "event", "title", "body", "read_at", "created_at", "title_key", "body_key"
FROM "notifications_bak"`,
	`DROP TABLE "notifications_bak"`,
	`CREATE TABLE role_menu_items (
  role_id      TEXT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
  menu_item_id TEXT NOT NULL REFERENCES menu_items(id) ON DELETE RESTRICT,
  PRIMARY KEY (role_id, menu_item_id)
)`,
	`INSERT INTO "role_menu_items" ("role_id", "menu_item_id")
SELECT "role_id", "menu_item_id"
FROM "role_menu_items_bak"`,
	`DROP TABLE "role_menu_items_bak"`,
	`CREATE TABLE role_permissions (
  role_id       TEXT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
  permission_id TEXT NOT NULL REFERENCES permissions(id) ON DELETE RESTRICT,
  PRIMARY KEY (role_id, permission_id)
)`,
	`INSERT INTO "role_permissions" ("role_id", "permission_id")
SELECT "role_id", "permission_id"
FROM "role_permissions_bak"`,
	`DROP TABLE "role_permissions_bak"`,
	`CREATE TABLE user_mfa (
  user_id                TEXT PRIMARY KEY REFERENCES users(id),
  status                 TEXT NOT NULL CHECK (status IN ('pending','active')),
  totp_secret_ciphertext TEXT NOT NULL,
  recovery_codes_hash    TEXT NOT NULL,
  last_used_step         INTEGER NOT NULL DEFAULT 0,
  created_at             INTEGER NOT NULL,
  updated_at             INTEGER NOT NULL
)`,
	`INSERT INTO "user_mfa" ("user_id", "status", "totp_secret_ciphertext", "recovery_codes_hash", "last_used_step", "created_at", "updated_at")
SELECT "user_id", "status", "totp_secret_ciphertext", "recovery_codes_hash", "last_used_step", "created_at", "updated_at"
FROM "user_mfa_bak"`,
	`DROP TABLE "user_mfa_bak"`,
	`CREATE TABLE user_roles (
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  role_id TEXT NOT NULL REFERENCES roles(id) ON DELETE RESTRICT,
  PRIMARY KEY (user_id, role_id)
)`,
	`INSERT INTO "user_roles" ("user_id", "role_id")
SELECT "user_id", "role_id"
FROM "user_roles_bak"`,
	`DROP TABLE "user_roles_bak"`,
	`CREATE INDEX idx_notifications_user_created ON notifications(user_id, created_at DESC)`,
	`CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id)`,
	`CREATE INDEX idx_role_menu_items_menu_item_id ON role_menu_items(menu_item_id)`,
	`CREATE INDEX idx_role_permissions_permission_id ON role_permissions(permission_id)`,
	`CREATE INDEX idx_service_credentials_created_at ON service_credentials(created_at DESC, id DESC)`,
	`CREATE INDEX idx_service_credentials_expires_at ON service_credentials(expires_at)`,
	`CREATE INDEX idx_user_invites_created ON user_invites(created_at DESC)`,
	`CREATE INDEX idx_user_password_history_user ON user_password_history(user_id, created_at DESC)`,
	`CREATE INDEX idx_user_roles_role_id ON user_roles(role_id)`,
	`CREATE UNIQUE INDEX idx_users_email_lower ON users(lower(email))`,
}

// vp040V74Statements returns the canonical checksum input (D-017: m0 → m1–m3 → m4).
func vp040V74Statements() []string {
	return temporalmigrate.Ordered(vp040V74Preflight, vp040V74Rebuild, vp040V74Verify)
}

// applyvp040V74 is the SQLite Apply body.
func applyvp040V74(tx kernel.Tx) error {
	if err := temporalmigrate.RunPreflight(tx, vp040V74Preflight); err != nil {
		return err
	}
	if err := temporalmigrate.Exec(tx, vp040V74Rebuild, "vp040 v74 sqlite rebuild"); err != nil {
		return err
	}
	return temporalmigrate.RunVerify(tx, vp040V74Verify, "vp040 v74 sqlite verify")
}

// vp040V74Postgres is the explicit postgres m1–m2 slice (D-019 §6: no regex
// derivation; the millisecond family uses the integer-split expression).
var vp040V74Postgres = []string{
	`ALTER TABLE "system_data_reconcile" ALTER COLUMN "applied_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("applied_at"::double precision)))`,
	`ALTER TABLE "users" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)))`,
	`ALTER TABLE "users" ALTER COLUMN "updated_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("updated_at"::double precision)))`,
	`ALTER TABLE "users" ALTER COLUMN "locked_until" DROP DEFAULT`,
	`ALTER TABLE "users" ALTER COLUMN "locked_until" DROP NOT NULL`,
	`ALTER TABLE "users" ALTER COLUMN "locked_until" TYPE timestamptz(6) USING (CASE WHEN "locked_until" = 0 THEN NULL ELSE date_trunc('microseconds', to_timestamp("locked_until"::double precision)) END)`,
	`ALTER TABLE "users" ALTER COLUMN "last_login_failure_at" DROP DEFAULT`,
	`ALTER TABLE "users" ALTER COLUMN "last_login_failure_at" DROP NOT NULL`,
	`ALTER TABLE "users" ALTER COLUMN "last_login_failure_at" TYPE timestamptz(6) USING (CASE WHEN "last_login_failure_at" = 0 THEN NULL ELSE date_trunc('microseconds', to_timestamp("last_login_failure_at"::double precision)) END)`,
	`ALTER TABLE "roles" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)))`,
	`ALTER TABLE "roles" ALTER COLUMN "updated_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("updated_at"::double precision)))`,
	`ALTER TABLE "permissions" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)))`,
	`ALTER TABLE "permissions" ALTER COLUMN "updated_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("updated_at"::double precision)))`,
	`ALTER TABLE "menu_items" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)))`,
	`ALTER TABLE "menu_items" ALTER COLUMN "updated_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("updated_at"::double precision)))`,
	`ALTER TABLE "refresh_tokens" ALTER COLUMN "expires_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("expires_at"::double precision)))`,
	`ALTER TABLE "refresh_tokens" ALTER COLUMN "revoked_at" TYPE timestamptz(6) USING (CASE WHEN "revoked_at" IS NULL THEN NULL ELSE date_trunc('microseconds', to_timestamp("revoked_at"::double precision)) END)`,
	`ALTER TABLE "refresh_tokens" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)))`,
	`ALTER TABLE "email_verification_challenges" ALTER COLUMN "expires_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("expires_at"::double precision)))`,
	`ALTER TABLE "email_verification_challenges" ALTER COLUMN "sent_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("sent_at"::double precision)))`,
	`ALTER TABLE "password_recovery_challenges" ALTER COLUMN "expires_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("expires_at"::double precision)))`,
	`ALTER TABLE "password_recovery_challenges" ALTER COLUMN "sent_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("sent_at"::double precision)))`,
	`ALTER TABLE "login_failures" ALTER COLUMN "locked_until" DROP DEFAULT`,
	`ALTER TABLE "login_failures" ALTER COLUMN "locked_until" DROP NOT NULL`,
	`ALTER TABLE "login_failures" ALTER COLUMN "locked_until" TYPE timestamptz(6) USING (CASE WHEN "locked_until" = 0 THEN NULL ELSE date_trunc('microseconds', to_timestamp("locked_until"::double precision)) END)`,
	`ALTER TABLE "login_failures" ALTER COLUMN "updated_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("updated_at"::double precision)))`,
	`ALTER TABLE "user_password_history" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)))`,
	`ALTER TABLE "user_invites" ALTER COLUMN "expires_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("expires_at"::double precision)))`,
	`ALTER TABLE "user_invites" ALTER COLUMN "consumed_at" TYPE timestamptz(6) USING (CASE WHEN "consumed_at" IS NULL THEN NULL ELSE date_trunc('microseconds', to_timestamp("consumed_at"::double precision)) END)`,
	`ALTER TABLE "user_invites" ALTER COLUMN "revoked_at" TYPE timestamptz(6) USING (CASE WHEN "revoked_at" IS NULL THEN NULL ELSE date_trunc('microseconds', to_timestamp("revoked_at"::double precision)) END)`,
	`ALTER TABLE "user_invites" ALTER COLUMN "last_sent_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("last_sent_at"::double precision)))`,
	`ALTER TABLE "user_invites" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)))`,
	`ALTER TABLE "service_credentials" ALTER COLUMN "expires_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("expires_at"::double precision)))`,
	`ALTER TABLE "service_credentials" ALTER COLUMN "revoked_at" TYPE timestamptz(6) USING (CASE WHEN "revoked_at" IS NULL THEN NULL ELSE date_trunc('microseconds', to_timestamp("revoked_at"::double precision)) END)`,
	`ALTER TABLE "service_credentials" ALTER COLUMN "last_used_at" TYPE timestamptz(6) USING (CASE WHEN "last_used_at" IS NULL THEN NULL ELSE date_trunc('microseconds', to_timestamp("last_used_at"::double precision)) END)`,
	`ALTER TABLE "service_credentials" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)))`,
	`ALTER TABLE "service_credentials" ALTER COLUMN "updated_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("updated_at"::double precision)))`,
}

// vp040V74PostgresVerify is the postgres m4 type/precision assertion set.
var vp040V74PostgresVerify = []temporalmigrate.PgVerify{
	{Table: "system_data_reconcile", Column: "applied_at", NonNull: true},
	{Table: "users", Column: "created_at", NonNull: true},
	{Table: "users", Column: "updated_at", NonNull: true},
	{Table: "users", Column: "locked_until", NonNull: false},
	{Table: "users", Column: "last_login_failure_at", NonNull: false},
	{Table: "roles", Column: "created_at", NonNull: true},
	{Table: "roles", Column: "updated_at", NonNull: true},
	{Table: "permissions", Column: "created_at", NonNull: true},
	{Table: "permissions", Column: "updated_at", NonNull: true},
	{Table: "menu_items", Column: "created_at", NonNull: true},
	{Table: "menu_items", Column: "updated_at", NonNull: true},
	{Table: "refresh_tokens", Column: "expires_at", NonNull: true},
	{Table: "refresh_tokens", Column: "revoked_at", NonNull: false},
	{Table: "refresh_tokens", Column: "created_at", NonNull: true},
	{Table: "email_verification_challenges", Column: "expires_at", NonNull: true},
	{Table: "email_verification_challenges", Column: "sent_at", NonNull: true},
	{Table: "password_recovery_challenges", Column: "expires_at", NonNull: true},
	{Table: "password_recovery_challenges", Column: "sent_at", NonNull: true},
	{Table: "login_failures", Column: "locked_until", NonNull: false},
	{Table: "login_failures", Column: "updated_at", NonNull: true},
	{Table: "user_password_history", Column: "created_at", NonNull: true},
	{Table: "user_invites", Column: "expires_at", NonNull: true},
	{Table: "user_invites", Column: "consumed_at", NonNull: false},
	{Table: "user_invites", Column: "revoked_at", NonNull: false},
	{Table: "user_invites", Column: "last_sent_at", NonNull: true},
	{Table: "user_invites", Column: "created_at", NonNull: true},
	{Table: "service_credentials", Column: "expires_at", NonNull: true},
	{Table: "service_credentials", Column: "revoked_at", NonNull: false},
	{Table: "service_credentials", Column: "last_used_at", NonNull: false},
	{Table: "service_credentials", Column: "created_at", NonNull: true},
	{Table: "service_credentials", Column: "updated_at", NonNull: true},
}

// applyvp040V74Postgres is the postgres Apply body.
func applyvp040V74Postgres(tx kernel.Tx) error {
	if err := temporalmigrate.RunPreflight(tx, vp040V74Preflight); err != nil {
		return err
	}
	if err := temporalmigrate.Exec(tx, vp040V74Postgres, "vp040 v74 postgres convert"); err != nil {
		return err
	}
	return temporalmigrate.RunPostgresVerify(tx, vp040V74PostgresVerify, "vp040 v74 postgres verify")
}

// VP040TemporalDescriptor returns the workspace-040 R2 conversion descriptor
// (D-014 allocation, D-017 checksum convention).
func VP040TemporalDescriptor() kernel.MigrationContribution {
	return kernel.MigrationContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: vp040V74Name},
		Version:              74,
		Name:                 vp040V74Name,
		Checksum:             kernel.MigrationChecksum(vp040V74Statements(), vp040V74TransformID),
		Apply:                applyvp040V74,
		ApplyPostgres:        applyvp040V74Postgres,
	}
}
