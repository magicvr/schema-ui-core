---
id: r2-v73-v87-generated-statements-v0.1
doc_type: evidence-attachment
title: R2 v73–v87 conversion descriptor 语句清单 v0.1（落码产物）
status: draft
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-003-r2-codec-and-descriptor-m1-m2
version: 0.1.0
---

# R2 v73–v87 conversion descriptor 语句清单 v0.1

> 本文件由生成器从 **v72 已 apply 库的 live `sqlite_master`** 机械导出（GOAL-003
> checkpoint B/C 落码产物），记录每个 descriptor 的 canonical 语句顺序（m0 预检 →
> m1–m3 重建 → m4 校验）与显式 PG DDL。**m1–m3 的 `CREATE TABLE` 正文 = live DDL
> 逐字 + 冻结的逐列类型编辑**（`INTEGER` → `TEXT`；D0/voucher 列去 `NOT NULL` /
> `DEFAULT 0`），其余 token 未改。
>
> 缺省约定：`INSERT … SELECT` 的列清单按 live `PRAGMA table_info` 的 cid 序展开；
> 标识符用双引号；全部 `CREATE INDEX` 在 `DROP TABLE <t>_old` / `<t>_bak` 之后。

## 0. FK 子表盘点（F-5 先决条件，A-032 §G / 逐表附件 §5.8）

机械解析 v72 DDL 的全部 `REFERENCES <parent>(` 后，**存在 FK 子表的父表**为：

- `dict_types` ← `dict_entries`
- `menu_items` ← `role_menu_items`
- `operation_log` ← `operation_log_correlation`, `operation_log_session`
- `permissions` ← `role_permissions`
- `roles` ← `role_menu_items`, `role_permissions`, `user_roles`
- `scheduled_tasks` ← `task_runs`
- `users` ← `email_verification_challenges`, `login_failures`, `mfa_proofs`, `notifications`, `password_recovery_challenges`, `refresh_tokens`, `user_invites`, `user_mfa`, `user_password_history`, `user_roles`

**无任何存留 DDL 引用（= 裸四步安全）的表（46 张）**：`captcha_challenges`, `captcha_config`, `data_scope_policies`, `dict_entries`, `digital_entitlements`, `digital_offers`, `digital_purchases`, `email_verification_challenges`, `jobs`, `login_failures`, `mail_config`, `mail_outbox`, `mfa_proofs`, `notifications`, `operation_log_archive`, `operation_log_archive_correlation`, `operation_log_archive_session`, `operation_log_correlation`, `operation_log_session`, `password_policy`, `password_recovery_challenges`, `recycle_items`, `refresh_tokens`, `role_menu_items`, `role_permissions`, `schema_migrations`, `service_credentials`, `site_settings`, `subjects`, `system_data_grants`, `system_data_reconcile`, `task_runs`, `telegram_config`, `telegram_inbound_messages`, `telegram_outbound_messages`, `telegram_sessions`, `user_data_scopes`, `user_invites`, `user_mfa`, `user_password_history`, `user_roles`, `voucher_batches`, `vouchers`, `wallet_accounts`, `wallet_ledger_entries`, `wallet_reconciliation_runs`

## v73 · `core.persistence` · `vp040_temporal_core_persistence`

- `transform_id`: `0073:vp040-temporal-core-persistence:v1`
- **`MigrationChecksum`（真实值，D-017 单 checksum / SQLite 切片）**：`4dd07092330cb3b344143f49ec57edd1f89e92cf2786446c0635c4a320eb8a9f`
- 表范围（descriptor 顺序）：`schema_migrations`, `mail_outbox`, `mail_config`
- FK 子女盘点（机械解析 v72 DDL 的 `REFERENCES`）：父表 （无）；摘除并建回的子表 （无）
- m0 预检列：`mail_config.updated_at`；m0 retired-table guard：`records`
- PG 语句数 5；SQLite 语句数（m1–m3）13；m4 断言 3

canonical（m1–m3）：

```sql
ALTER TABLE "schema_migrations" RENAME TO "schema_migrations_old";
CREATE TABLE schema_migrations (
  version    INTEGER PRIMARY KEY,
  name       TEXT NOT NULL UNIQUE,
  checksum   TEXT NOT NULL CHECK (length(checksum) = 64),
  applied_at TEXT NOT NULL
);
INSERT INTO "schema_migrations" ("version", "name", "checksum", "applied_at")
SELECT "version", "name", "checksum", strftime('%Y-%m-%dT%H:%M:%S', applied_at, 'unixepoch') || '.000000Z'
FROM "schema_migrations_old";
DROP TABLE "schema_migrations_old";
ALTER TABLE "mail_outbox" RENAME TO "mail_outbox_old";
CREATE TABLE mail_outbox (
  id         TEXT PRIMARY KEY,
  to_addr    TEXT NOT NULL,
  subject    TEXT NOT NULL,
  body       TEXT NOT NULL,
  created_at TEXT NOT NULL
, channel TEXT NOT NULL DEFAULT 'mock', delivery_status TEXT NOT NULL DEFAULT 'delivered');
INSERT INTO "mail_outbox" ("id", "to_addr", "subject", "body", "created_at", "channel", "delivery_status")
SELECT "id", "to_addr", "subject", "body", strftime('%Y-%m-%dT%H:%M:%S', CASE WHEN created_at >= 0 THEN created_at/1000 ELSE (created_at-999)/1000 END, 'unixepoch') || '.' || printf('%03d', (created_at%1000 + 1000) % 1000) || '000Z', "channel", "delivery_status"
FROM "mail_outbox_old";
DROP TABLE "mail_outbox_old";
ALTER TABLE "mail_config" RENAME TO "mail_config_old";
CREATE TABLE mail_config (
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
);
INSERT INTO "mail_config" ("id", "channel", "mock_retention", "resend_from", "resend_api_key_enc", "smtp_host", "smtp_port", "smtp_username", "smtp_password_enc", "smtp_from", "updated_at")
SELECT "id", "channel", "mock_retention", "resend_from", "resend_api_key_enc", "smtp_host", "smtp_port", "smtp_username", "smtp_password_enc", "smtp_from", CASE WHEN updated_at = 0 THEN NULL ELSE strftime('%Y-%m-%dT%H:%M:%S', CASE WHEN updated_at >= 0 THEN updated_at/1000 ELSE (updated_at-999)/1000 END, 'unixepoch') || '.' || printf('%03d', (updated_at%1000 + 1000) % 1000) || '000Z' END
FROM "mail_config_old";
DROP TABLE "mail_config_old";
CREATE INDEX idx_mail_outbox_created_at ON mail_outbox(created_at);
```

m0（guard/预检）/ m4（校验）语句文本：

```sql
SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = 'records';
	{Table: "mail_config", Column: "updated_at", Voucher: false},;
	{Table: "schema_migrations", Columns: []string{"applied_at"}},;
	{Table: "mail_outbox", Columns: []string{"created_at"}},;
	{Table: "mail_config", Columns: []string{"updated_at"}},;
```

<details><summary>postgres m1–m2</summary>

```sql
ALTER TABLE "schema_migrations" ALTER COLUMN "applied_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("applied_at"::double precision)));
ALTER TABLE "mail_outbox" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (TIMESTAMPTZ 'epoch' + (("created_at" - CASE WHEN "created_at" >= 0 THEN 0 ELSE 999 END) / 1000) * INTERVAL '1 second' + ((("created_at" % 1000) + 1000) % 1000) * INTERVAL '1 millisecond');
ALTER TABLE "mail_config" ALTER COLUMN "updated_at" DROP DEFAULT;
ALTER TABLE "mail_config" ALTER COLUMN "updated_at" DROP NOT NULL;
ALTER TABLE "mail_config" ALTER COLUMN "updated_at" TYPE timestamptz(6) USING (CASE WHEN "updated_at" = 0 THEN NULL ELSE TIMESTAMPTZ 'epoch' + (("updated_at" - CASE WHEN "updated_at" >= 0 THEN 0 ELSE 999 END) / 1000) * INTERVAL '1 second' + ((("updated_at" % 1000) + 1000) % 1000) * INTERVAL '1 millisecond' END);
```

</details>

## v74 · `core.auth-session` · `vp040_temporal_authsession`

- `transform_id`: `0074:vp040-temporal-authsession:v1`
- **`MigrationChecksum`（真实值，D-017 单 checksum / SQLite 切片）**：`ae1aefe89925f1759e6e0154f64ddf73b9eb04472c1111e1400b473d7c477f55`
- 表范围（descriptor 顺序）：`system_data_reconcile`, `users`, `roles`, `permissions`, `menu_items`, `refresh_tokens`, `email_verification_challenges`, `password_recovery_challenges`, `login_failures`, `user_password_history`, `user_invites`, `service_credentials`
- FK 子女盘点（机械解析 v72 DDL 的 `REFERENCES`）：父表 `menu_items`, `permissions`, `roles`, `users`；摘除并建回的子表 `email_verification_challenges`, `login_failures`, `mfa_proofs`, `notifications`, `password_recovery_challenges`, `refresh_tokens`, `role_menu_items`, `role_permissions`, `user_invites`, `user_mfa`, `user_password_history`, `user_roles`
- m0 预检列：`users.locked_until`, `users.last_login_failure_at`, `login_failures.locked_until`；m0 retired-table guard：（无）
- PG 语句数 37；SQLite 语句数（m1–m3）94；m4 断言 12

canonical（m1–m3）：

```sql
CREATE TEMP TABLE "email_verification_challenges_bak" AS SELECT * FROM "email_verification_challenges";
CREATE TEMP TABLE "login_failures_bak" AS SELECT * FROM "login_failures";
CREATE TEMP TABLE "mfa_proofs_bak" AS SELECT * FROM "mfa_proofs";
CREATE TEMP TABLE "notifications_bak" AS SELECT * FROM "notifications";
CREATE TEMP TABLE "password_recovery_challenges_bak" AS SELECT * FROM "password_recovery_challenges";
CREATE TEMP TABLE "refresh_tokens_bak" AS SELECT * FROM "refresh_tokens";
CREATE TEMP TABLE "role_menu_items_bak" AS SELECT * FROM "role_menu_items";
CREATE TEMP TABLE "role_permissions_bak" AS SELECT * FROM "role_permissions";
CREATE TEMP TABLE "user_invites_bak" AS SELECT * FROM "user_invites";
CREATE TEMP TABLE "user_mfa_bak" AS SELECT * FROM "user_mfa";
CREATE TEMP TABLE "user_password_history_bak" AS SELECT * FROM "user_password_history";
CREATE TEMP TABLE "user_roles_bak" AS SELECT * FROM "user_roles";
DROP TABLE "email_verification_challenges";
DROP TABLE "login_failures";
DROP TABLE "mfa_proofs";
DROP TABLE "notifications";
DROP TABLE "password_recovery_challenges";
DROP TABLE "refresh_tokens";
DROP TABLE "role_menu_items";
DROP TABLE "role_permissions";
DROP TABLE "user_invites";
DROP TABLE "user_mfa";
DROP TABLE "user_password_history";
DROP TABLE "user_roles";
ALTER TABLE "users" RENAME TO "users_old";
CREATE TABLE users (
  id            TEXT PRIMARY KEY,
  username      TEXT NOT NULL UNIQUE,
  name          TEXT NOT NULL,
  roles         TEXT NOT NULL, -- JSON array; R3 normalizes
  password_hash TEXT NOT NULL,
  created_at    TEXT NOT NULL,
  updated_at    TEXT NOT NULL
, token_version INTEGER NOT NULL DEFAULT 0, failed_login_count INTEGER NOT NULL DEFAULT 0, locked_until TEXT, enabled INTEGER NOT NULL DEFAULT 1, notifications_enabled INTEGER NOT NULL DEFAULT 1, avatar_url TEXT NOT NULL DEFAULT '', must_change_password INTEGER NOT NULL DEFAULT 0, email TEXT, email_status TEXT CHECK (email_status IN ('pending','verified')), last_login_failure_at TEXT);
INSERT INTO "users" ("id", "username", "name", "roles", "password_hash", "created_at", "updated_at", "token_version", "failed_login_count", "locked_until", "enabled", "notifications_enabled", "avatar_url", "must_change_password", "email", "email_status", "last_login_failure_at")
SELECT "id", "username", "name", "roles", "password_hash", strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z', strftime('%Y-%m-%dT%H:%M:%S', updated_at, 'unixepoch') || '.000000Z', "token_version", "failed_login_count", CASE WHEN locked_until = 0 THEN NULL ELSE strftime('%Y-%m-%dT%H:%M:%S', locked_until, 'unixepoch') || '.000000Z' END, "enabled", "notifications_enabled", "avatar_url", "must_change_password", "email", "email_status", CASE WHEN last_login_failure_at = 0 THEN NULL ELSE strftime('%Y-%m-%dT%H:%M:%S', last_login_failure_at, 'unixepoch') || '.000000Z' END
FROM "users_old";
DROP TABLE "users_old";
ALTER TABLE "roles" RENAME TO "roles_old";
CREATE TABLE roles (
  id         TEXT PRIMARY KEY,
  key        TEXT NOT NULL UNIQUE CHECK (key <> ''),
  name       TEXT NOT NULL,
  system     INTEGER NOT NULL DEFAULT 0 CHECK (system IN (0, 1)),
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);
INSERT INTO "roles" ("id", "key", "name", "system", "created_at", "updated_at")
SELECT "id", "key", "name", "system", strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z', strftime('%Y-%m-%dT%H:%M:%S', updated_at, 'unixepoch') || '.000000Z'
FROM "roles_old";
DROP TABLE "roles_old";
ALTER TABLE "permissions" RENAME TO "permissions_old";
CREATE TABLE permissions (
  id          TEXT PRIMARY KEY,
  key         TEXT NOT NULL UNIQUE CHECK (key <> ''),
  description TEXT NOT NULL DEFAULT '',
  created_at  TEXT NOT NULL,
  updated_at  TEXT NOT NULL
);
INSERT INTO "permissions" ("id", "key", "description", "created_at", "updated_at")
SELECT "id", "key", "description", strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z', strftime('%Y-%m-%dT%H:%M:%S', updated_at, 'unixepoch') || '.000000Z'
FROM "permissions_old";
DROP TABLE "permissions_old";
ALTER TABLE "menu_items" RENAME TO "menu_items_old";
CREATE TABLE menu_items (
  id          TEXT PRIMARY KEY,
  page_ref    TEXT NOT NULL UNIQUE CHECK (page_ref <> ''),
  feature_key TEXT NOT NULL UNIQUE CHECK (feature_key <> ''),
  sort_order  INTEGER NOT NULL DEFAULT 0,
  enabled     INTEGER NOT NULL DEFAULT 1 CHECK (enabled IN (0, 1)),
  created_at  TEXT NOT NULL,
  updated_at  TEXT NOT NULL
);
INSERT INTO "menu_items" ("id", "page_ref", "feature_key", "sort_order", "enabled", "created_at", "updated_at")
SELECT "id", "page_ref", "feature_key", "sort_order", "enabled", strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z', strftime('%Y-%m-%dT%H:%M:%S', updated_at, 'unixepoch') || '.000000Z'
FROM "menu_items_old";
DROP TABLE "menu_items_old";
CREATE TABLE refresh_tokens (
  id         TEXT PRIMARY KEY,
  user_id    TEXT NOT NULL REFERENCES users(id),
  token_hash TEXT NOT NULL UNIQUE,
  expires_at TEXT NOT NULL,
  revoked_at TEXT,
  created_at TEXT NOT NULL
);
INSERT INTO "refresh_tokens" ("id", "user_id", "token_hash", "expires_at", "revoked_at", "created_at")
SELECT "id", "user_id", "token_hash", strftime('%Y-%m-%dT%H:%M:%S', expires_at, 'unixepoch') || '.000000Z', CASE WHEN revoked_at IS NULL THEN NULL ELSE strftime('%Y-%m-%dT%H:%M:%S', revoked_at, 'unixepoch') || '.000000Z' END, strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z'
FROM "refresh_tokens_bak";
DROP TABLE "refresh_tokens_bak";
CREATE TABLE email_verification_challenges (
  user_id       TEXT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  code_hash     TEXT NOT NULL,
  expires_at    TEXT NOT NULL,
  sent_at       TEXT NOT NULL,
  attempt_count INTEGER NOT NULL DEFAULT 0
);
INSERT INTO "email_verification_challenges" ("user_id", "code_hash", "expires_at", "sent_at", "attempt_count")
SELECT "user_id", "code_hash", strftime('%Y-%m-%dT%H:%M:%S', expires_at, 'unixepoch') || '.000000Z', strftime('%Y-%m-%dT%H:%M:%S', sent_at, 'unixepoch') || '.000000Z', "attempt_count"
FROM "email_verification_challenges_bak";
DROP TABLE "email_verification_challenges_bak";
CREATE TABLE password_recovery_challenges (
  user_id       TEXT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  code_hash     TEXT NOT NULL,
  expires_at    TEXT NOT NULL,
  sent_at       TEXT NOT NULL,
  attempt_count INTEGER NOT NULL DEFAULT 0
);
INSERT INTO "password_recovery_challenges" ("user_id", "code_hash", "expires_at", "sent_at", "attempt_count")
SELECT "user_id", "code_hash", strftime('%Y-%m-%dT%H:%M:%S', expires_at, 'unixepoch') || '.000000Z', strftime('%Y-%m-%dT%H:%M:%S', sent_at, 'unixepoch') || '.000000Z', "attempt_count"
FROM "password_recovery_challenges_bak";
DROP TABLE "password_recovery_challenges_bak";
CREATE TABLE login_failures (
  user_id      TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  ip          TEXT NOT NULL,
  fail_count  INTEGER NOT NULL DEFAULT 0,
  locked_until TEXT,
  updated_at  TEXT NOT NULL,
  PRIMARY KEY (user_id, ip)
);
INSERT INTO "login_failures" ("user_id", "ip", "fail_count", "locked_until", "updated_at")
SELECT "user_id", "ip", "fail_count", CASE WHEN locked_until = 0 THEN NULL ELSE strftime('%Y-%m-%dT%H:%M:%S', locked_until, 'unixepoch') || '.000000Z' END, strftime('%Y-%m-%dT%H:%M:%S', updated_at, 'unixepoch') || '.000000Z'
FROM "login_failures_bak";
DROP TABLE "login_failures_bak";
CREATE TABLE user_password_history (
  id            TEXT PRIMARY KEY,
  user_id       TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  password_hash TEXT NOT NULL,
  created_at    TEXT NOT NULL
);
INSERT INTO "user_password_history" ("id", "user_id", "password_hash", "created_at")
SELECT "id", "user_id", "password_hash", strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z'
FROM "user_password_history_bak";
DROP TABLE "user_password_history_bak";
CREATE TABLE user_invites (
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
);
INSERT INTO "user_invites" ("id", "token_hash", "roles", "invited_by", "email", "expires_at", "consumed_at", "revoked_at", "last_sent_at", "created_at")
SELECT "id", "token_hash", "roles", "invited_by", "email", strftime('%Y-%m-%dT%H:%M:%S', expires_at, 'unixepoch') || '.000000Z', CASE WHEN consumed_at IS NULL THEN NULL ELSE strftime('%Y-%m-%dT%H:%M:%S', consumed_at, 'unixepoch') || '.000000Z' END, CASE WHEN revoked_at IS NULL THEN NULL ELSE strftime('%Y-%m-%dT%H:%M:%S', revoked_at, 'unixepoch') || '.000000Z' END, strftime('%Y-%m-%dT%H:%M:%S', last_sent_at, 'unixepoch') || '.000000Z', strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z'
FROM "user_invites_bak";
DROP TABLE "user_invites_bak";
ALTER TABLE "system_data_reconcile" RENAME TO "system_data_reconcile_old";
CREATE TABLE system_data_reconcile (
  module_id        TEXT NOT NULL,
  kind             TEXT NOT NULL CHECK (kind IN ('base','authorization','navigation')),
  contribution_key TEXT NOT NULL,
  version          INTEGER NOT NULL CHECK (version > 0),
  checksum         TEXT NOT NULL CHECK (length(checksum) = 64),
  applied_at       TEXT NOT NULL,
  PRIMARY KEY (module_id, kind, contribution_key)
);
INSERT INTO "system_data_reconcile" ("module_id", "kind", "contribution_key", "version", "checksum", "applied_at")
SELECT "module_id", "kind", "contribution_key", "version", "checksum", strftime('%Y-%m-%dT%H:%M:%S', applied_at, 'unixepoch') || '.000000Z'
FROM "system_data_reconcile_old";
DROP TABLE "system_data_reconcile_old";
ALTER TABLE "service_credentials" RENAME TO "service_credentials_old";
CREATE TABLE service_credentials (
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
);
INSERT INTO "service_credentials" ("id", "name", "token_prefix", "token_hash", "scopes", "expires_at", "revoked_at", "last_used_at", "created_by", "created_at", "updated_at")
SELECT "id", "name", "token_prefix", "token_hash", "scopes", strftime('%Y-%m-%dT%H:%M:%S', expires_at, 'unixepoch') || '.000000Z', CASE WHEN revoked_at IS NULL THEN NULL ELSE strftime('%Y-%m-%dT%H:%M:%S', revoked_at, 'unixepoch') || '.000000Z' END, CASE WHEN last_used_at IS NULL THEN NULL ELSE strftime('%Y-%m-%dT%H:%M:%S', last_used_at, 'unixepoch') || '.000000Z' END, "created_by", strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z', strftime('%Y-%m-%dT%H:%M:%S', updated_at, 'unixepoch') || '.000000Z'
FROM "service_credentials_old";
DROP TABLE "service_credentials_old";
CREATE TABLE mfa_proofs (
  id         TEXT PRIMARY KEY,
  user_id    TEXT NOT NULL REFERENCES users(id),
  fail_count INTEGER NOT NULL DEFAULT 0,
  expires_at INTEGER NOT NULL,
  created_at INTEGER NOT NULL
);
INSERT INTO "mfa_proofs" ("id", "user_id", "fail_count", "expires_at", "created_at")
SELECT "id", "user_id", "fail_count", "expires_at", "created_at"
FROM "mfa_proofs_bak";
DROP TABLE "mfa_proofs_bak";
CREATE TABLE notifications (
  id         TEXT PRIMARY KEY,
  user_id    TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  event      TEXT NOT NULL CHECK (event IN ('account.locked','account.disabled','account.unlocked','account.password-changed')),
  title      TEXT NOT NULL,
  body       TEXT NOT NULL,
  read_at    INTEGER,
  created_at INTEGER NOT NULL
, title_key TEXT, body_key TEXT);
INSERT INTO "notifications" ("id", "user_id", "event", "title", "body", "read_at", "created_at", "title_key", "body_key")
SELECT "id", "user_id", "event", "title", "body", "read_at", "created_at", "title_key", "body_key"
FROM "notifications_bak";
DROP TABLE "notifications_bak";
CREATE TABLE role_menu_items (
  role_id      TEXT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
  menu_item_id TEXT NOT NULL REFERENCES menu_items(id) ON DELETE RESTRICT,
  PRIMARY KEY (role_id, menu_item_id)
);
INSERT INTO "role_menu_items" ("role_id", "menu_item_id")
SELECT "role_id", "menu_item_id"
FROM "role_menu_items_bak";
DROP TABLE "role_menu_items_bak";
CREATE TABLE role_permissions (
  role_id       TEXT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
  permission_id TEXT NOT NULL REFERENCES permissions(id) ON DELETE RESTRICT,
  PRIMARY KEY (role_id, permission_id)
);
INSERT INTO "role_permissions" ("role_id", "permission_id")
SELECT "role_id", "permission_id"
FROM "role_permissions_bak";
DROP TABLE "role_permissions_bak";
CREATE TABLE user_mfa (
  user_id                TEXT PRIMARY KEY REFERENCES users(id),
  status                 TEXT NOT NULL CHECK (status IN ('pending','active')),
  totp_secret_ciphertext TEXT NOT NULL,
  recovery_codes_hash    TEXT NOT NULL,
  last_used_step         INTEGER NOT NULL DEFAULT 0,
  created_at             INTEGER NOT NULL,
  updated_at             INTEGER NOT NULL
);
INSERT INTO "user_mfa" ("user_id", "status", "totp_secret_ciphertext", "recovery_codes_hash", "last_used_step", "created_at", "updated_at")
SELECT "user_id", "status", "totp_secret_ciphertext", "recovery_codes_hash", "last_used_step", "created_at", "updated_at"
FROM "user_mfa_bak";
DROP TABLE "user_mfa_bak";
CREATE TABLE user_roles (
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  role_id TEXT NOT NULL REFERENCES roles(id) ON DELETE RESTRICT,
  PRIMARY KEY (user_id, role_id)
);
INSERT INTO "user_roles" ("user_id", "role_id")
SELECT "user_id", "role_id"
FROM "user_roles_bak";
DROP TABLE "user_roles_bak";
CREATE INDEX idx_notifications_user_created ON notifications(user_id, created_at DESC);
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_role_menu_items_menu_item_id ON role_menu_items(menu_item_id);
CREATE INDEX idx_role_permissions_permission_id ON role_permissions(permission_id);
CREATE INDEX idx_service_credentials_created_at ON service_credentials(created_at DESC, id DESC);
CREATE INDEX idx_service_credentials_expires_at ON service_credentials(expires_at);
CREATE INDEX idx_user_invites_created ON user_invites(created_at DESC);
CREATE INDEX idx_user_password_history_user ON user_password_history(user_id, created_at DESC);
CREATE INDEX idx_user_roles_role_id ON user_roles(role_id);
CREATE UNIQUE INDEX idx_users_email_lower ON users(lower(email));
```

m0（guard/预检）/ m4（校验）语句文本：

```sql
	{Table: "users", Column: "locked_until", Voucher: false},;
	{Table: "users", Column: "last_login_failure_at", Voucher: false},;
	{Table: "login_failures", Column: "locked_until", Voucher: false},;
	{Table: "system_data_reconcile", Columns: []string{"applied_at"}},;
	{Table: "users", Columns: []string{"created_at", "updated_at", "locked_until", "last_login_failure_at"}, Children: []string{"email_verification_challenges", "login_failures", "mfa_proofs", "notifications", "password_recovery_challenges", "refresh_tokens", "user_invites", "user_mfa", "user_password_history", "user_roles"}},;
	{Table: "roles", Columns: []string{"created_at", "updated_at"}, Children: []string{"role_menu_items", "role_permissions", "user_roles"}},;
	{Table: "permissions", Columns: []string{"created_at", "updated_at"}, Children: []string{"role_permissions"}},;
	{Table: "menu_items", Columns: []string{"created_at", "updated_at"}, Children: []string{"role_menu_items"}},;
	{Table: "refresh_tokens", Columns: []string{"expires_at", "revoked_at", "created_at"}},;
	{Table: "email_verification_challenges", Columns: []string{"expires_at", "sent_at"}},;
	{Table: "password_recovery_challenges", Columns: []string{"expires_at", "sent_at"}},;
	{Table: "login_failures", Columns: []string{"locked_until", "updated_at"}},;
	{Table: "user_password_history", Columns: []string{"created_at"}},;
	{Table: "user_invites", Columns: []string{"expires_at", "consumed_at", "revoked_at", "last_sent_at", "created_at"}},;
	{Table: "service_credentials", Columns: []string{"expires_at", "revoked_at", "last_used_at", "created_at", "updated_at"}},;
```

<details><summary>postgres m1–m2</summary>

```sql
ALTER TABLE "system_data_reconcile" ALTER COLUMN "applied_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("applied_at"::double precision)));
ALTER TABLE "users" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)));
ALTER TABLE "users" ALTER COLUMN "updated_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("updated_at"::double precision)));
ALTER TABLE "users" ALTER COLUMN "locked_until" DROP DEFAULT;
ALTER TABLE "users" ALTER COLUMN "locked_until" DROP NOT NULL;
ALTER TABLE "users" ALTER COLUMN "locked_until" TYPE timestamptz(6) USING (CASE WHEN "locked_until" = 0 THEN NULL ELSE date_trunc('microseconds', to_timestamp("locked_until"::double precision)) END);
ALTER TABLE "users" ALTER COLUMN "last_login_failure_at" DROP DEFAULT;
ALTER TABLE "users" ALTER COLUMN "last_login_failure_at" DROP NOT NULL;
ALTER TABLE "users" ALTER COLUMN "last_login_failure_at" TYPE timestamptz(6) USING (CASE WHEN "last_login_failure_at" = 0 THEN NULL ELSE date_trunc('microseconds', to_timestamp("last_login_failure_at"::double precision)) END);
ALTER TABLE "roles" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)));
ALTER TABLE "roles" ALTER COLUMN "updated_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("updated_at"::double precision)));
ALTER TABLE "permissions" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)));
ALTER TABLE "permissions" ALTER COLUMN "updated_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("updated_at"::double precision)));
ALTER TABLE "menu_items" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)));
ALTER TABLE "menu_items" ALTER COLUMN "updated_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("updated_at"::double precision)));
ALTER TABLE "refresh_tokens" ALTER COLUMN "expires_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("expires_at"::double precision)));
ALTER TABLE "refresh_tokens" ALTER COLUMN "revoked_at" TYPE timestamptz(6) USING (CASE WHEN "revoked_at" IS NULL THEN NULL ELSE date_trunc('microseconds', to_timestamp("revoked_at"::double precision)) END);
ALTER TABLE "refresh_tokens" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)));
ALTER TABLE "email_verification_challenges" ALTER COLUMN "expires_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("expires_at"::double precision)));
ALTER TABLE "email_verification_challenges" ALTER COLUMN "sent_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("sent_at"::double precision)));
ALTER TABLE "password_recovery_challenges" ALTER COLUMN "expires_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("expires_at"::double precision)));
ALTER TABLE "password_recovery_challenges" ALTER COLUMN "sent_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("sent_at"::double precision)));
ALTER TABLE "login_failures" ALTER COLUMN "locked_until" DROP DEFAULT;
ALTER TABLE "login_failures" ALTER COLUMN "locked_until" DROP NOT NULL;
ALTER TABLE "login_failures" ALTER COLUMN "locked_until" TYPE timestamptz(6) USING (CASE WHEN "locked_until" = 0 THEN NULL ELSE date_trunc('microseconds', to_timestamp("locked_until"::double precision)) END);
ALTER TABLE "login_failures" ALTER COLUMN "updated_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("updated_at"::double precision)));
ALTER TABLE "user_password_history" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)));
ALTER TABLE "user_invites" ALTER COLUMN "expires_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("expires_at"::double precision)));
ALTER TABLE "user_invites" ALTER COLUMN "consumed_at" TYPE timestamptz(6) USING (CASE WHEN "consumed_at" IS NULL THEN NULL ELSE date_trunc('microseconds', to_timestamp("consumed_at"::double precision)) END);
ALTER TABLE "user_invites" ALTER COLUMN "revoked_at" TYPE timestamptz(6) USING (CASE WHEN "revoked_at" IS NULL THEN NULL ELSE date_trunc('microseconds', to_timestamp("revoked_at"::double precision)) END);
ALTER TABLE "user_invites" ALTER COLUMN "last_sent_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("last_sent_at"::double precision)));
ALTER TABLE "user_invites" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)));
ALTER TABLE "service_credentials" ALTER COLUMN "expires_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("expires_at"::double precision)));
ALTER TABLE "service_credentials" ALTER COLUMN "revoked_at" TYPE timestamptz(6) USING (CASE WHEN "revoked_at" IS NULL THEN NULL ELSE date_trunc('microseconds', to_timestamp("revoked_at"::double precision)) END);
ALTER TABLE "service_credentials" ALTER COLUMN "last_used_at" TYPE timestamptz(6) USING (CASE WHEN "last_used_at" IS NULL THEN NULL ELSE date_trunc('microseconds', to_timestamp("last_used_at"::double precision)) END);
ALTER TABLE "service_credentials" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)));
ALTER TABLE "service_credentials" ALTER COLUMN "updated_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("updated_at"::double precision)));
```

</details>

## v75 · `core.operationlog` · `vp040_temporal_operationlog`

- `transform_id`: `0075:vp040-temporal-operationlog:v1`
- **`MigrationChecksum`（真实值，D-017 单 checksum / SQLite 切片）**：`5b038c6cf66f0e01f2446721a246917fc4c0b459dc877242f374aea971322b3a`
- 表范围（descriptor 顺序）：`operation_log`, `operation_log_archive`
- FK 子女盘点（机械解析 v72 DDL 的 `REFERENCES`）：父表 `operation_log`；摘除并建回的子表 `operation_log_correlation`, `operation_log_session`
- m0 预检列：（无）；m0 retired-table guard：（无）
- PG 语句数 3；SQLite 语句数（m1–m3）22；m4 断言 2

canonical（m1–m3）：

```sql
CREATE TEMP TABLE "operation_log_correlation_bak" AS SELECT * FROM "operation_log_correlation";
CREATE TEMP TABLE "operation_log_session_bak" AS SELECT * FROM "operation_log_session";
DROP TABLE "operation_log_correlation";
DROP TABLE "operation_log_session";
ALTER TABLE "operation_log" RENAME TO "operation_log_old";
CREATE TABLE operation_log (
  id         TEXT PRIMARY KEY,
  event      TEXT NOT NULL CHECK (event IN ('records.create','records.update','records.delete','auth.login','auth.logout','auth.refresh','users.create','users.update','users.delete','roles.create','roles.update','roles.delete','settings.update','users.enable','users.disable','users.unlock','account.password-change','account.session-revoke','data.export','data.import','files.upload','files.download','files.delete','dictionary.create','dictionary.update','dictionary.delete','scheduled-tasks.create','scheduled-tasks.update','scheduled-tasks.delete','captcha.settings-update','recycle.restore','recycle.purge','data-permission.policy-update','data-permission.scope-update','mfa.enroll','mfa.confirm','mfa.disable','mfa.recovery-rotate','mfa.admin-reset','mfa.login','wallet.account-create','wallet.account-update','wallet.adjust','wallet.freeze','wallet.unfreeze','wallet.reconcile','wallet.deduct-frozen','account.avatar-change','wallet.reconcile.queued','wallet.reconcile.failed','wallet.reconcile.cancelled','service-credentials.create','service-credentials.use','service-credentials.revoke','mail.channel-update','mail.test-send','bizoffer.offer.create','bizoffer.offer.update','bizoffer.offer.status','bizoffer.entitlement.void')),
  actor_id   TEXT NOT NULL,
  actor_name TEXT NOT NULL,
  record_id  TEXT,
  detail     TEXT,
  created_at TEXT NOT NULL
);
INSERT INTO "operation_log" ("id", "event", "actor_id", "actor_name", "record_id", "detail", "created_at")
SELECT "id", "event", "actor_id", "actor_name", "record_id", "detail", strftime('%Y-%m-%dT%H:%M:%S', CASE WHEN created_at >= 0 THEN created_at/1000 ELSE (created_at-999)/1000 END, 'unixepoch') || '.' || printf('%03d', (created_at%1000 + 1000) % 1000) || '000Z'
FROM "operation_log_old";
DROP TABLE "operation_log_old";
ALTER TABLE "operation_log_archive" RENAME TO "operation_log_archive_old";
CREATE TABLE operation_log_archive (
  id          TEXT PRIMARY KEY,
  event       TEXT NOT NULL,
  actor_id    TEXT NOT NULL,
  actor_name  TEXT NOT NULL,
  record_id   TEXT,
  detail      TEXT,
  created_at  TEXT NOT NULL,
  archived_at TEXT NOT NULL
);
INSERT INTO "operation_log_archive" ("id", "event", "actor_id", "actor_name", "record_id", "detail", "created_at", "archived_at")
SELECT "id", "event", "actor_id", "actor_name", "record_id", "detail", strftime('%Y-%m-%dT%H:%M:%S', CASE WHEN created_at >= 0 THEN created_at/1000 ELSE (created_at-999)/1000 END, 'unixepoch') || '.' || printf('%03d', (created_at%1000 + 1000) % 1000) || '000Z', strftime('%Y-%m-%dT%H:%M:%S', CASE WHEN archived_at >= 0 THEN archived_at/1000 ELSE (archived_at-999)/1000 END, 'unixepoch') || '.' || printf('%03d', (archived_at%1000 + 1000) % 1000) || '000Z'
FROM "operation_log_archive_old";
DROP TABLE "operation_log_archive_old";
CREATE TABLE operation_log_correlation (
  operation_id   TEXT PRIMARY KEY REFERENCES operation_log(id) ON DELETE CASCADE,
  correlation_id TEXT NOT NULL
);
INSERT INTO "operation_log_correlation" ("operation_id", "correlation_id")
SELECT "operation_id", "correlation_id"
FROM "operation_log_correlation_bak";
DROP TABLE "operation_log_correlation_bak";
CREATE TABLE operation_log_session (
  operation_id TEXT PRIMARY KEY REFERENCES operation_log(id) ON DELETE CASCADE,
  session_id   TEXT NOT NULL
);
INSERT INTO "operation_log_session" ("operation_id", "session_id")
SELECT "operation_id", "session_id"
FROM "operation_log_session_bak";
DROP TABLE "operation_log_session_bak";
CREATE INDEX idx_operation_log_created_at ON operation_log(created_at DESC);
CREATE INDEX idx_operation_log_archive_created_at ON operation_log_archive(created_at DESC);
CREATE INDEX idx_operation_log_correlation_id ON operation_log_correlation(correlation_id);
CREATE INDEX idx_operation_log_session_id ON operation_log_session(session_id);
```

m0（guard/预检）/ m4（校验）语句文本：

```sql
	{Table: "operation_log", Columns: []string{"created_at"}, Children: []string{"operation_log_correlation", "operation_log_session"}},;
	{Table: "operation_log_archive", Columns: []string{"created_at", "archived_at"}},;
```

<details><summary>postgres m1–m2</summary>

```sql
ALTER TABLE "operation_log" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (TIMESTAMPTZ 'epoch' + (("created_at" - CASE WHEN "created_at" >= 0 THEN 0 ELSE 999 END) / 1000) * INTERVAL '1 second' + ((("created_at" % 1000) + 1000) % 1000) * INTERVAL '1 millisecond');
ALTER TABLE "operation_log_archive" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (TIMESTAMPTZ 'epoch' + (("created_at" - CASE WHEN "created_at" >= 0 THEN 0 ELSE 999 END) / 1000) * INTERVAL '1 second' + ((("created_at" % 1000) + 1000) % 1000) * INTERVAL '1 millisecond');
ALTER TABLE "operation_log_archive" ALTER COLUMN "archived_at" TYPE timestamptz(6) USING (TIMESTAMPTZ 'epoch' + (("archived_at" - CASE WHEN "archived_at" >= 0 THEN 0 ELSE 999 END) / 1000) * INTERVAL '1 second' + ((("archived_at" % 1000) + 1000) % 1000) * INTERVAL '1 millisecond');
```

</details>

## v76 · `core.jobs` · `vp040_temporal_jobs`

- `transform_id`: `0076:vp040-temporal-jobs:v1`
- **`MigrationChecksum`（真实值，D-017 单 checksum / SQLite 切片）**：`6b3649579cc6aedc713fab8c7f9f6dafbec548317f7395082d9e5ddb730aca56`
- 表范围（descriptor 顺序）：`jobs`
- FK 子女盘点（机械解析 v72 DDL 的 `REFERENCES`）：父表 （无）；摘除并建回的子表 （无）
- m0 预检列：（无）；m0 retired-table guard：（无）
- PG 语句数 5；SQLite 语句数（m1–m3）8；m4 断言 1

canonical（m1–m3）：

```sql
ALTER TABLE "jobs" RENAME TO "jobs_old";
CREATE TABLE jobs (
  id               TEXT PRIMARY KEY,
  kind             TEXT NOT NULL CHECK (length(trim(kind)) > 0),
  status           TEXT NOT NULL CHECK (status IN ('queued','running','succeeded','failed','cancelled','expired')),
  payload          TEXT NOT NULL DEFAULT '{}',
  progress         INTEGER NOT NULL DEFAULT 0 CHECK (progress BETWEEN 0 AND 100),
  cancel_requested INTEGER NOT NULL DEFAULT 0 CHECK (cancel_requested IN (0,1)),
  attempt          INTEGER NOT NULL DEFAULT 0 CHECK (attempt >= 0),
  max_attempts     INTEGER NOT NULL DEFAULT 3 CHECK (max_attempts > 0 AND attempt <= max_attempts),
  lease_owner      TEXT,
  lease_version    INTEGER NOT NULL DEFAULT 0 CHECK (lease_version >= 0),
  lease_expires_at TEXT,
  result           TEXT,
  error_code       TEXT,
  error_message    TEXT,
  actor_id         TEXT NOT NULL CHECK (length(trim(actor_id)) > 0),
  correlation_id   TEXT NOT NULL CHECK (length(trim(correlation_id)) > 0),
  created_at       TEXT NOT NULL,
  updated_at       TEXT NOT NULL,
  finished_at      TEXT,
  expires_at       TEXT,
  CHECK (
    (status = 'queued' AND lease_owner IS NULL AND lease_expires_at IS NULL AND result IS NULL AND error_code IS NULL AND finished_at IS NULL AND expires_at IS NULL)
    OR (status = 'running' AND lease_owner IS NOT NULL AND lease_expires_at IS NOT NULL AND result IS NULL AND error_code IS NULL AND finished_at IS NULL AND expires_at IS NULL)
    OR (status = 'succeeded' AND lease_owner IS NULL AND lease_expires_at IS NULL AND result IS NOT NULL AND error_code IS NULL AND progress = 100 AND finished_at IS NOT NULL AND expires_at IS NOT NULL)
    OR (status = 'failed' AND lease_owner IS NULL AND lease_expires_at IS NULL AND result IS NULL AND error_code IS NOT NULL AND finished_at IS NOT NULL AND expires_at IS NULL)
    OR (status = 'cancelled' AND lease_owner IS NULL AND lease_expires_at IS NULL AND result IS NULL AND error_code IS NULL AND finished_at IS NOT NULL AND expires_at IS NULL)
    OR (status = 'expired' AND lease_owner IS NULL AND lease_expires_at IS NULL AND result IS NULL AND error_code IS NULL AND progress = 100 AND finished_at IS NOT NULL AND expires_at IS NOT NULL)
  )
);
INSERT INTO "jobs" ("id", "kind", "status", "payload", "progress", "cancel_requested", "attempt", "max_attempts", "lease_owner", "lease_version", "lease_expires_at", "result", "error_code", "error_message", "actor_id", "correlation_id", "created_at", "updated_at", "finished_at", "expires_at")
SELECT "id", "kind", "status", "payload", "progress", "cancel_requested", "attempt", "max_attempts", "lease_owner", "lease_version", CASE WHEN lease_expires_at IS NULL THEN NULL ELSE strftime('%Y-%m-%dT%H:%M:%S', CASE WHEN lease_expires_at >= 0 THEN lease_expires_at/1000 ELSE (lease_expires_at-999)/1000 END, 'unixepoch') || '.' || printf('%03d', (lease_expires_at%1000 + 1000) % 1000) || '000Z' END, "result", "error_code", "error_message", "actor_id", "correlation_id", strftime('%Y-%m-%dT%H:%M:%S', CASE WHEN created_at >= 0 THEN created_at/1000 ELSE (created_at-999)/1000 END, 'unixepoch') || '.' || printf('%03d', (created_at%1000 + 1000) % 1000) || '000Z', strftime('%Y-%m-%dT%H:%M:%S', CASE WHEN updated_at >= 0 THEN updated_at/1000 ELSE (updated_at-999)/1000 END, 'unixepoch') || '.' || printf('%03d', (updated_at%1000 + 1000) % 1000) || '000Z', CASE WHEN finished_at IS NULL THEN NULL ELSE strftime('%Y-%m-%dT%H:%M:%S', CASE WHEN finished_at >= 0 THEN finished_at/1000 ELSE (finished_at-999)/1000 END, 'unixepoch') || '.' || printf('%03d', (finished_at%1000 + 1000) % 1000) || '000Z' END, CASE WHEN expires_at IS NULL THEN NULL ELSE strftime('%Y-%m-%dT%H:%M:%S', CASE WHEN expires_at >= 0 THEN expires_at/1000 ELSE (expires_at-999)/1000 END, 'unixepoch') || '.' || printf('%03d', (expires_at%1000 + 1000) % 1000) || '000Z' END
FROM "jobs_old";
DROP TABLE "jobs_old";
CREATE INDEX idx_jobs_actor ON jobs(actor_id, kind, updated_at DESC);
CREATE INDEX idx_jobs_created_at ON jobs(created_at DESC, id DESC);
CREATE INDEX idx_jobs_expiry ON jobs(status, expires_at);
CREATE INDEX idx_jobs_runnable ON jobs(status, cancel_requested, lease_expires_at, created_at);
```

m0（guard/预检）/ m4（校验）语句文本：

```sql
	{Table: "jobs", Columns: []string{"lease_expires_at", "created_at", "updated_at", "finished_at", "expires_at"}},;
```

<details><summary>postgres m1–m2</summary>

```sql
ALTER TABLE "jobs" ALTER COLUMN "lease_expires_at" TYPE timestamptz(6) USING (CASE WHEN "lease_expires_at" IS NULL THEN NULL ELSE TIMESTAMPTZ 'epoch' + (("lease_expires_at" - CASE WHEN "lease_expires_at" >= 0 THEN 0 ELSE 999 END) / 1000) * INTERVAL '1 second' + ((("lease_expires_at" % 1000) + 1000) % 1000) * INTERVAL '1 millisecond' END);
ALTER TABLE "jobs" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (TIMESTAMPTZ 'epoch' + (("created_at" - CASE WHEN "created_at" >= 0 THEN 0 ELSE 999 END) / 1000) * INTERVAL '1 second' + ((("created_at" % 1000) + 1000) % 1000) * INTERVAL '1 millisecond');
ALTER TABLE "jobs" ALTER COLUMN "updated_at" TYPE timestamptz(6) USING (TIMESTAMPTZ 'epoch' + (("updated_at" - CASE WHEN "updated_at" >= 0 THEN 0 ELSE 999 END) / 1000) * INTERVAL '1 second' + ((("updated_at" % 1000) + 1000) % 1000) * INTERVAL '1 millisecond');
ALTER TABLE "jobs" ALTER COLUMN "finished_at" TYPE timestamptz(6) USING (CASE WHEN "finished_at" IS NULL THEN NULL ELSE TIMESTAMPTZ 'epoch' + (("finished_at" - CASE WHEN "finished_at" >= 0 THEN 0 ELSE 999 END) / 1000) * INTERVAL '1 second' + ((("finished_at" % 1000) + 1000) % 1000) * INTERVAL '1 millisecond' END);
ALTER TABLE "jobs" ALTER COLUMN "expires_at" TYPE timestamptz(6) USING (CASE WHEN "expires_at" IS NULL THEN NULL ELSE TIMESTAMPTZ 'epoch' + (("expires_at" - CASE WHEN "expires_at" >= 0 THEN 0 ELSE 999 END) / 1000) * INTERVAL '1 second' + ((("expires_at" % 1000) + 1000) % 1000) * INTERVAL '1 millisecond' END);
```

</details>

## v77 · `admin.data-dictionary` · `vp040_temporal_dictionary`

- `transform_id`: `0077:vp040-temporal-dictionary:v1`
- **`MigrationChecksum`（真实值，D-017 单 checksum / SQLite 切片）**：`0266f2937603f3cbed07b33ee460964ad1fe083e2c0f426c4946b0f922104edb`
- 表范围（descriptor 顺序）：`dict_types`, `dict_entries`
- FK 子女盘点（机械解析 v72 DDL 的 `REFERENCES`）：父表 `dict_types`；摘除并建回的子表 `dict_entries`
- m0 预检列：（无）；m0 retired-table guard：（无）
- PG 语句数 4；SQLite 语句数（m1–m3）10；m4 断言 2

canonical（m1–m3）：

```sql
CREATE TEMP TABLE "dict_entries_bak" AS SELECT * FROM "dict_entries";
DROP TABLE "dict_entries";
ALTER TABLE "dict_types" RENAME TO "dict_types_old";
CREATE TABLE dict_types (
  id         TEXT PRIMARY KEY,
  key        TEXT NOT NULL UNIQUE,
  name       TEXT NOT NULL,
  enabled    INTEGER NOT NULL DEFAULT 1,
  description TEXT,
  sort       INTEGER NOT NULL DEFAULT 0,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);
INSERT INTO "dict_types" ("id", "key", "name", "enabled", "description", "sort", "created_at", "updated_at")
SELECT "id", "key", "name", "enabled", "description", "sort", strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z', strftime('%Y-%m-%dT%H:%M:%S', updated_at, 'unixepoch') || '.000000Z'
FROM "dict_types_old";
DROP TABLE "dict_types_old";
CREATE TABLE dict_entries (
  id         TEXT PRIMARY KEY,
  dict_key   TEXT NOT NULL REFERENCES dict_types(key) ON DELETE CASCADE,
  entry_key  TEXT NOT NULL,
  label      TEXT NOT NULL,
  enabled    INTEGER NOT NULL DEFAULT 1,
  sort       INTEGER NOT NULL DEFAULT 0,
  remark     TEXT,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL, badge_style TEXT NOT NULL DEFAULT 'default',
  UNIQUE (dict_key, entry_key)
);
INSERT INTO "dict_entries" ("id", "dict_key", "entry_key", "label", "enabled", "sort", "remark", "created_at", "updated_at", "badge_style")
SELECT "id", "dict_key", "entry_key", "label", "enabled", "sort", "remark", strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z', strftime('%Y-%m-%dT%H:%M:%S', updated_at, 'unixepoch') || '.000000Z', "badge_style"
FROM "dict_entries_bak";
DROP TABLE "dict_entries_bak";
CREATE INDEX idx_dict_entries_dict_key ON dict_entries(dict_key, sort);
```

m0（guard/预检）/ m4（校验）语句文本：

```sql
	{Table: "dict_types", Columns: []string{"created_at", "updated_at"}, Children: []string{"dict_entries"}},;
	{Table: "dict_entries", Columns: []string{"created_at", "updated_at"}},;
```

<details><summary>postgres m1–m2</summary>

```sql
ALTER TABLE "dict_types" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)));
ALTER TABLE "dict_types" ALTER COLUMN "updated_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("updated_at"::double precision)));
ALTER TABLE "dict_entries" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)));
ALTER TABLE "dict_entries" ALTER COLUMN "updated_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("updated_at"::double precision)));
```

</details>

## v78 · `admin.data-permission` · `vp040_temporal_data_permission`

- `transform_id`: `0078:vp040-temporal-data-permission:v1`
- **`MigrationChecksum`（真实值，D-017 单 checksum / SQLite 切片）**：`b1fa8aa94597a8f48a061efe43d7285300016070d391e38b8df95f49dbcdf7ed`
- 表范围（descriptor 顺序）：`data_scope_policies`, `user_data_scopes`
- FK 子女盘点（机械解析 v72 DDL 的 `REFERENCES`）：父表 （无）；摘除并建回的子表 （无）
- m0 预检列：（无）；m0 retired-table guard：（无）
- PG 语句数 2；SQLite 语句数（m1–m3）8；m4 断言 2

canonical（m1–m3）：

```sql
ALTER TABLE "data_scope_policies" RENAME TO "data_scope_policies_old";
CREATE TABLE data_scope_policies (
  resource      TEXT PRIMARY KEY,
  owner_column  TEXT NOT NULL,
  default_scope TEXT NOT NULL CHECK (default_scope IN ('all','self')),
  enabled       INTEGER NOT NULL DEFAULT 1,
  updated_at    TEXT NOT NULL
);
INSERT INTO "data_scope_policies" ("resource", "owner_column", "default_scope", "enabled", "updated_at")
SELECT "resource", "owner_column", "default_scope", "enabled", strftime('%Y-%m-%dT%H:%M:%S', updated_at, 'unixepoch') || '.000000Z'
FROM "data_scope_policies_old";
DROP TABLE "data_scope_policies_old";
ALTER TABLE "user_data_scopes" RENAME TO "user_data_scopes_old";
CREATE TABLE user_data_scopes (
  user_id    TEXT NOT NULL,
  resource   TEXT NOT NULL,
  scope_type TEXT NOT NULL CHECK (scope_type IN ('all','self')),
  updated_at TEXT NOT NULL,
  PRIMARY KEY (user_id, resource)
);
INSERT INTO "user_data_scopes" ("user_id", "resource", "scope_type", "updated_at")
SELECT "user_id", "resource", "scope_type", strftime('%Y-%m-%dT%H:%M:%S', updated_at, 'unixepoch') || '.000000Z'
FROM "user_data_scopes_old";
DROP TABLE "user_data_scopes_old";
```

m0（guard/预检）/ m4（校验）语句文本：

```sql
	{Table: "data_scope_policies", Columns: []string{"updated_at"}},;
	{Table: "user_data_scopes", Columns: []string{"updated_at"}},;
```

<details><summary>postgres m1–m2</summary>

```sql
ALTER TABLE "data_scope_policies" ALTER COLUMN "updated_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("updated_at"::double precision)));
ALTER TABLE "user_data_scopes" ALTER COLUMN "updated_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("updated_at"::double precision)));
```

</details>

## v79 · `admin.login-captcha` · `vp040_temporal_captcha`

- `transform_id`: `0079:vp040-temporal-captcha:v1`
- **`MigrationChecksum`（真实值，D-017 单 checksum / SQLite 切片）**：`c6a661ebfb90158f3a712ad149084f3e84f996ca6773f04cb8743b0e2e41142f`
- 表范围（descriptor 顺序）：`captcha_challenges`, `captcha_config`
- FK 子女盘点（机械解析 v72 DDL 的 `REFERENCES`）：父表 （无）；摘除并建回的子表 （无）
- m0 预检列：（无）；m0 retired-table guard：（无）
- PG 语句数 4；SQLite 语句数（m1–m3）8；m4 断言 2

canonical（m1–m3）：

```sql
ALTER TABLE "captcha_challenges" RENAME TO "captcha_challenges_old";
CREATE TABLE captcha_challenges (
  id          TEXT PRIMARY KEY,
  answer_hash TEXT NOT NULL,
  expires_at  TEXT NOT NULL,
  created_at  TEXT NOT NULL
);
INSERT INTO "captcha_challenges" ("id", "answer_hash", "expires_at", "created_at")
SELECT "id", "answer_hash", strftime('%Y-%m-%dT%H:%M:%S', expires_at, 'unixepoch') || '.000000Z', strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z'
FROM "captcha_challenges_old";
DROP TABLE "captcha_challenges_old";
ALTER TABLE "captcha_config" RENAME TO "captcha_config_old";
CREATE TABLE captcha_config (
  id         INTEGER PRIMARY KEY CHECK (id = 1),
  enabled    INTEGER NOT NULL DEFAULT 0,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);
INSERT INTO "captcha_config" ("id", "enabled", "created_at", "updated_at")
SELECT "id", "enabled", strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z', strftime('%Y-%m-%dT%H:%M:%S', updated_at, 'unixepoch') || '.000000Z'
FROM "captcha_config_old";
DROP TABLE "captcha_config_old";
```

m0（guard/预检）/ m4（校验）语句文本：

```sql
	{Table: "captcha_challenges", Columns: []string{"expires_at", "created_at"}},;
	{Table: "captcha_config", Columns: []string{"created_at", "updated_at"}},;
```

<details><summary>postgres m1–m2</summary>

```sql
ALTER TABLE "captcha_challenges" ALTER COLUMN "expires_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("expires_at"::double precision)));
ALTER TABLE "captcha_challenges" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)));
ALTER TABLE "captcha_config" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)));
ALTER TABLE "captcha_config" ALTER COLUMN "updated_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("updated_at"::double precision)));
```

</details>

## v80 · `admin.mfa` · `vp040_temporal_mfa`

- `transform_id`: `0080:vp040-temporal-mfa:v1`
- **`MigrationChecksum`（真实值，D-017 单 checksum / SQLite 切片）**：`2088626f0bdf5c2b9aba9a5faf297762cba553c0cf750c1255e1c8633ed6c2fc`
- 表范围（descriptor 顺序）：`user_mfa`, `mfa_proofs`
- FK 子女盘点（机械解析 v72 DDL 的 `REFERENCES`）：父表 （无）；摘除并建回的子表 （无）
- m0 预检列：（无）；m0 retired-table guard：（无）
- PG 语句数 4；SQLite 语句数（m1–m3）8；m4 断言 2

canonical（m1–m3）：

```sql
ALTER TABLE "user_mfa" RENAME TO "user_mfa_old";
CREATE TABLE user_mfa (
  user_id                TEXT PRIMARY KEY REFERENCES users(id),
  status                 TEXT NOT NULL CHECK (status IN ('pending','active')),
  totp_secret_ciphertext TEXT NOT NULL,
  recovery_codes_hash    TEXT NOT NULL,
  last_used_step         INTEGER NOT NULL DEFAULT 0,
  created_at             TEXT NOT NULL,
  updated_at             TEXT NOT NULL
);
INSERT INTO "user_mfa" ("user_id", "status", "totp_secret_ciphertext", "recovery_codes_hash", "last_used_step", "created_at", "updated_at")
SELECT "user_id", "status", "totp_secret_ciphertext", "recovery_codes_hash", "last_used_step", strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z', strftime('%Y-%m-%dT%H:%M:%S', updated_at, 'unixepoch') || '.000000Z'
FROM "user_mfa_old";
DROP TABLE "user_mfa_old";
ALTER TABLE "mfa_proofs" RENAME TO "mfa_proofs_old";
CREATE TABLE mfa_proofs (
  id         TEXT PRIMARY KEY,
  user_id    TEXT NOT NULL REFERENCES users(id),
  fail_count INTEGER NOT NULL DEFAULT 0,
  expires_at TEXT NOT NULL,
  created_at TEXT NOT NULL
);
INSERT INTO "mfa_proofs" ("id", "user_id", "fail_count", "expires_at", "created_at")
SELECT "id", "user_id", "fail_count", strftime('%Y-%m-%dT%H:%M:%S', expires_at, 'unixepoch') || '.000000Z', strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z'
FROM "mfa_proofs_old";
DROP TABLE "mfa_proofs_old";
```

m0（guard/预检）/ m4（校验）语句文本：

```sql
	{Table: "user_mfa", Columns: []string{"created_at", "updated_at"}},;
	{Table: "mfa_proofs", Columns: []string{"expires_at", "created_at"}},;
```

<details><summary>postgres m1–m2</summary>

```sql
ALTER TABLE "user_mfa" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)));
ALTER TABLE "user_mfa" ALTER COLUMN "updated_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("updated_at"::double precision)));
ALTER TABLE "mfa_proofs" ALTER COLUMN "expires_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("expires_at"::double precision)));
ALTER TABLE "mfa_proofs" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)));
```

</details>

## v81 · `admin.notifications` · `vp040_temporal_notifications`

- `transform_id`: `0081:vp040-temporal-notifications:v1`
- **`MigrationChecksum`（真实值，D-017 单 checksum / SQLite 切片）**：`a7565dc641f3c3291ff25cbefa52199ca06a8f94f7a978efac5b907615442b37`
- 表范围（descriptor 顺序）：`notifications`
- FK 子女盘点（机械解析 v72 DDL 的 `REFERENCES`）：父表 （无）；摘除并建回的子表 （无）
- m0 预检列：（无）；m0 retired-table guard：（无）
- PG 语句数 2；SQLite 语句数（m1–m3）5；m4 断言 1

canonical（m1–m3）：

```sql
ALTER TABLE "notifications" RENAME TO "notifications_old";
CREATE TABLE notifications (
  id         TEXT PRIMARY KEY,
  user_id    TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  event      TEXT NOT NULL CHECK (event IN ('account.locked','account.disabled','account.unlocked','account.password-changed')),
  title      TEXT NOT NULL,
  body       TEXT NOT NULL,
  read_at    TEXT,
  created_at TEXT NOT NULL
, title_key TEXT, body_key TEXT);
INSERT INTO "notifications" ("id", "user_id", "event", "title", "body", "read_at", "created_at", "title_key", "body_key")
SELECT "id", "user_id", "event", "title", "body", CASE WHEN read_at IS NULL THEN NULL ELSE strftime('%Y-%m-%dT%H:%M:%S', read_at, 'unixepoch') || '.000000Z' END, strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z', "title_key", "body_key"
FROM "notifications_old";
DROP TABLE "notifications_old";
CREATE INDEX idx_notifications_user_created ON notifications(user_id, created_at DESC);
```

m0（guard/预检）/ m4（校验）语句文本：

```sql
	{Table: "notifications", Columns: []string{"read_at", "created_at"}},;
```

<details><summary>postgres m1–m2</summary>

```sql
ALTER TABLE "notifications" ALTER COLUMN "read_at" TYPE timestamptz(6) USING (CASE WHEN "read_at" IS NULL THEN NULL ELSE date_trunc('microseconds', to_timestamp("read_at"::double precision)) END);
ALTER TABLE "notifications" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)));
```

</details>

## v82 · `admin.recycle-bin` · `vp040_temporal_recycle`

- `transform_id`: `0082:vp040-temporal-recycle:v1`
- **`MigrationChecksum`（真实值，D-017 单 checksum / SQLite 切片）**：`0132f6a873dd427b42c3a668bc88badc9b50a6c6729601e7cd1c5d5e4a1568fe`
- 表范围（descriptor 顺序）：`recycle_items`
- FK 子女盘点（机械解析 v72 DDL 的 `REFERENCES`）：父表 （无）；摘除并建回的子表 （无）
- m0 预检列：（无）；m0 retired-table guard：（无）
- PG 语句数 2；SQLite 语句数（m1–m3）6；m4 断言 1

canonical（m1–m3）：

```sql
ALTER TABLE "recycle_items" RENAME TO "recycle_items_old";
CREATE TABLE recycle_items (
  id          TEXT PRIMARY KEY,
  resource    TEXT NOT NULL,
  resource_id TEXT NOT NULL,
  payload     TEXT NOT NULL,
  actor_id    TEXT NOT NULL,
  actor_name  TEXT NOT NULL,
  deleted_at  TEXT NOT NULL,
  restored_at TEXT
);
INSERT INTO "recycle_items" ("id", "resource", "resource_id", "payload", "actor_id", "actor_name", "deleted_at", "restored_at")
SELECT "id", "resource", "resource_id", "payload", "actor_id", "actor_name", strftime('%Y-%m-%dT%H:%M:%S', deleted_at, 'unixepoch') || '.000000Z', CASE WHEN restored_at IS NULL THEN NULL ELSE strftime('%Y-%m-%dT%H:%M:%S', restored_at, 'unixepoch') || '.000000Z' END
FROM "recycle_items_old";
DROP TABLE "recycle_items_old";
CREATE INDEX idx_recycle_items_deleted_at ON recycle_items(deleted_at DESC);
CREATE UNIQUE INDEX idx_recycle_items_active ON recycle_items(resource, resource_id) WHERE restored_at IS NULL;
```

m0（guard/预检）/ m4（校验）语句文本：

```sql
	{Table: "recycle_items", Columns: []string{"deleted_at", "restored_at"}},;
```

<details><summary>postgres m1–m2</summary>

```sql
ALTER TABLE "recycle_items" ALTER COLUMN "deleted_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("deleted_at"::double precision)));
ALTER TABLE "recycle_items" ALTER COLUMN "restored_at" TYPE timestamptz(6) USING (CASE WHEN "restored_at" IS NULL THEN NULL ELSE date_trunc('microseconds', to_timestamp("restored_at"::double precision)) END);
```

</details>

## v83 · `admin.scheduled-tasks` · `vp040_temporal_scheduled_tasks`

- `transform_id`: `0083:vp040-temporal-scheduled-tasks:v1`
- **`MigrationChecksum`（真实值，D-017 单 checksum / SQLite 切片）**：`b6c4f115e54d163a2ec9a1dfab0c74ca5d10c09c9cf1c41c14b02ad5d1186e91`
- 表范围（descriptor 顺序）：`scheduled_tasks`, `task_runs`
- FK 子女盘点（机械解析 v72 DDL 的 `REFERENCES`）：父表 `scheduled_tasks`；摘除并建回的子表 `task_runs`
- m0 预检列：（无）；m0 retired-table guard：（无）
- PG 语句数 5；SQLite 语句数（m1–m3）10；m4 断言 2

canonical（m1–m3）：

```sql
CREATE TEMP TABLE "task_runs_bak" AS SELECT * FROM "task_runs";
DROP TABLE "task_runs";
ALTER TABLE "scheduled_tasks" RENAME TO "scheduled_tasks_old";
CREATE TABLE scheduled_tasks (
  id          TEXT PRIMARY KEY,
  key         TEXT NOT NULL UNIQUE,
  cron        TEXT NOT NULL,
  name        TEXT NOT NULL,
  enabled     INTEGER NOT NULL DEFAULT 1,
  description TEXT,
  handler     TEXT NOT NULL DEFAULT 'system.noop',
  created_at  TEXT NOT NULL,
  updated_at  TEXT NOT NULL
);
INSERT INTO "scheduled_tasks" ("id", "key", "cron", "name", "enabled", "description", "handler", "created_at", "updated_at")
SELECT "id", "key", "cron", "name", "enabled", "description", "handler", strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z', strftime('%Y-%m-%dT%H:%M:%S', updated_at, 'unixepoch') || '.000000Z'
FROM "scheduled_tasks_old";
DROP TABLE "scheduled_tasks_old";
CREATE TABLE task_runs (
  id          TEXT PRIMARY KEY,
  task_id     TEXT NOT NULL REFERENCES scheduled_tasks(id) ON DELETE CASCADE,
  status      TEXT NOT NULL CHECK (status IN ('ran','failed')),
  started_at  TEXT NOT NULL,
  finished_at TEXT,
  detail      TEXT,
  created_at  TEXT NOT NULL
);
INSERT INTO "task_runs" ("id", "task_id", "status", "started_at", "finished_at", "detail", "created_at")
SELECT "id", "task_id", "status", strftime('%Y-%m-%dT%H:%M:%S', started_at, 'unixepoch') || '.000000Z', CASE WHEN finished_at IS NULL THEN NULL ELSE strftime('%Y-%m-%dT%H:%M:%S', finished_at, 'unixepoch') || '.000000Z' END, "detail", strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z'
FROM "task_runs_bak";
DROP TABLE "task_runs_bak";
CREATE INDEX idx_task_runs_task_started ON task_runs(task_id, started_at DESC);
```

m0（guard/预检）/ m4（校验）语句文本：

```sql
	{Table: "scheduled_tasks", Columns: []string{"created_at", "updated_at"}, Children: []string{"task_runs"}},;
	{Table: "task_runs", Columns: []string{"started_at", "finished_at", "created_at"}},;
```

<details><summary>postgres m1–m2</summary>

```sql
ALTER TABLE "scheduled_tasks" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)));
ALTER TABLE "scheduled_tasks" ALTER COLUMN "updated_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("updated_at"::double precision)));
ALTER TABLE "task_runs" ALTER COLUMN "started_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("started_at"::double precision)));
ALTER TABLE "task_runs" ALTER COLUMN "finished_at" TYPE timestamptz(6) USING (CASE WHEN "finished_at" IS NULL THEN NULL ELSE date_trunc('microseconds', to_timestamp("finished_at"::double precision)) END);
ALTER TABLE "task_runs" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)));
```

</details>

## v84 · `admin.settings` · `vp040_temporal_settings`

- `transform_id`: `0084:vp040-temporal-settings:v1`
- **`MigrationChecksum`（真实值，D-017 单 checksum / SQLite 切片）**：`bb3041a3d3fbeb5b3d706209f53cc578dc0e5d15016502919aac040b6bec2112`
- 表范围（descriptor 顺序）：`site_settings`
- FK 子女盘点（机械解析 v72 DDL 的 `REFERENCES`）：父表 （无）；摘除并建回的子表 （无）
- m0 预检列：（无）；m0 retired-table guard：（无）
- PG 语句数 1；SQLite 语句数（m1–m3）5；m4 断言 1

canonical（m1–m3）：

```sql
ALTER TABLE "site_settings" RENAME TO "site_settings_old";
CREATE TABLE site_settings (
  id         TEXT PRIMARY KEY CHECK (id = 'default'),
  site_title TEXT NOT NULL,
  logo_url   TEXT NOT NULL DEFAULT '',
  updated_at TEXT NOT NULL
, logo_url_light TEXT NOT NULL DEFAULT '', logo_url_dark TEXT NOT NULL DEFAULT '', favicon_url TEXT NOT NULL DEFAULT '', default_locale TEXT NOT NULL DEFAULT '', site_timezone TEXT NOT NULL DEFAULT '', default_theme TEXT NOT NULL DEFAULT '', copyright_text TEXT NOT NULL DEFAULT '', icp_number TEXT NOT NULL DEFAULT '', operation_log_retention_days INTEGER NOT NULL DEFAULT 90, operation_log_expiration_action TEXT NOT NULL DEFAULT 'archive', default_currency TEXT NOT NULL DEFAULT '');
INSERT INTO "site_settings" ("id", "site_title", "logo_url", "updated_at", "logo_url_light", "logo_url_dark", "favicon_url", "default_locale", "site_timezone", "default_theme", "copyright_text", "icp_number", "operation_log_retention_days", "operation_log_expiration_action", "default_currency")
SELECT "id", "site_title", "logo_url", strftime('%Y-%m-%dT%H:%M:%S', updated_at, 'unixepoch') || '.000000Z', "logo_url_light", "logo_url_dark", "favicon_url", "default_locale", "site_timezone", "default_theme", "copyright_text", "icp_number", "operation_log_retention_days", "operation_log_expiration_action", "default_currency"
FROM "site_settings_old";
DROP TABLE "site_settings_old";
CREATE INDEX idx_site_settings_updated_at ON site_settings (updated_at);
```

m0（guard/预检）/ m4（校验）语句文本：

```sql
	{Table: "site_settings", Columns: []string{"updated_at"}},;
```

<details><summary>postgres m1–m2</summary>

```sql
ALTER TABLE "site_settings" ALTER COLUMN "updated_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("updated_at"::double precision)));
```

</details>

## v85 · `admin.wallet` · `vp040_temporal_wallet`

- `transform_id`: `0085:vp040-temporal-wallet:v1`
- **`MigrationChecksum`（真实值，D-017 单 checksum / SQLite 切片）**：`5c004e0c643f46de5e72035022642f9674ca92306784782014f4e285c07efeb8`
- 表范围（descriptor 顺序）：`wallet_accounts`, `wallet_ledger_entries`, `wallet_reconciliation_runs`, `subjects`, `vouchers`, `voucher_batches`
- FK 子女盘点（机械解析 v72 DDL 的 `REFERENCES`）：父表 （无）；摘除并建回的子表 （无）
- m0 预检列：`vouchers.expires_at`, `vouchers.redeemed_at`；m0 retired-table guard：（无）
- PG 语句数 11；SQLite 语句数（m1–m3）28；m4 断言 6

canonical（m1–m3）：

```sql
ALTER TABLE "wallet_accounts" RENAME TO "wallet_accounts_old";
CREATE TABLE wallet_accounts (
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
);
INSERT INTO "wallet_accounts" ("id", "owner_type", "owner_id", "currency", "balance_total", "balance_available", "balance_frozen", "status", "version", "created_at", "updated_at")
SELECT "id", "owner_type", "owner_id", "currency", "balance_total", "balance_available", "balance_frozen", "status", "version", strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z', strftime('%Y-%m-%dT%H:%M:%S', updated_at, 'unixepoch') || '.000000Z'
FROM "wallet_accounts_old";
DROP TABLE "wallet_accounts_old";
ALTER TABLE "wallet_ledger_entries" RENAME TO "wallet_ledger_entries_old";
CREATE TABLE wallet_ledger_entries (
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
);
INSERT INTO "wallet_ledger_entries" ("id", "account_id", "entry_type", "amount_delta", "balance_after_total", "balance_after_available", "balance_after_frozen", "ref_type", "ref_id", "idempotency_key", "memo", "actor_id", "actor_name", "created_at")
SELECT "id", "account_id", "entry_type", "amount_delta", "balance_after_total", "balance_after_available", "balance_after_frozen", "ref_type", "ref_id", "idempotency_key", "memo", "actor_id", "actor_name", strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z'
FROM "wallet_ledger_entries_old";
DROP TABLE "wallet_ledger_entries_old";
ALTER TABLE "wallet_reconciliation_runs" RENAME TO "wallet_reconciliation_runs_old";
CREATE TABLE wallet_reconciliation_runs (
  id             TEXT PRIMARY KEY,
  account_id     TEXT,
  result         TEXT NOT NULL CHECK (result IN ('consistent','inconsistent')),
  mismatch_count INTEGER NOT NULL DEFAULT 0,
  details        TEXT NOT NULL DEFAULT '{}',
  actor_id       TEXT NOT NULL,
  created_at     TEXT NOT NULL
);
INSERT INTO "wallet_reconciliation_runs" ("id", "account_id", "result", "mismatch_count", "details", "actor_id", "created_at")
SELECT "id", "account_id", "result", "mismatch_count", "details", "actor_id", strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z'
FROM "wallet_reconciliation_runs_old";
DROP TABLE "wallet_reconciliation_runs_old";
ALTER TABLE "subjects" RENAME TO "subjects_old";
CREATE TABLE subjects (
  id          TEXT PRIMARY KEY,
  issuer      TEXT NOT NULL,
  external_id TEXT NOT NULL,
  created_at  TEXT NOT NULL,
  UNIQUE (issuer, external_id)
);
INSERT INTO "subjects" ("id", "issuer", "external_id", "created_at")
SELECT "id", "issuer", "external_id", strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z'
FROM "subjects_old";
DROP TABLE "subjects_old";
ALTER TABLE "vouchers" RENAME TO "vouchers_old";
CREATE TABLE vouchers (
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
);
INSERT INTO "vouchers" ("id", "batch_id", "code_hash", "code_prefix", "amount", "currency", "status", "expires_at", "redeemed_by", "redeemed_at", "created_at", "updated_at")
SELECT "id", "batch_id", "code_hash", "code_prefix", "amount", "currency", "status", CASE WHEN expires_at IS NULL OR expires_at = 0 THEN NULL ELSE strftime('%Y-%m-%dT%H:%M:%S', expires_at, 'unixepoch') || '.000000Z' END, "redeemed_by", CASE WHEN redeemed_at IS NULL OR redeemed_at = 0 THEN NULL ELSE strftime('%Y-%m-%dT%H:%M:%S', redeemed_at, 'unixepoch') || '.000000Z' END, strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z', strftime('%Y-%m-%dT%H:%M:%S', updated_at, 'unixepoch') || '.000000Z'
FROM "vouchers_old";
DROP TABLE "vouchers_old";
ALTER TABLE "voucher_batches" RENAME TO "voucher_batches_old";
CREATE TABLE voucher_batches (
  batch_id   TEXT PRIMARY KEY,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);
INSERT INTO "voucher_batches" ("batch_id", "created_at", "updated_at")
SELECT "batch_id", strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z', strftime('%Y-%m-%dT%H:%M:%S', updated_at, 'unixepoch') || '.000000Z'
FROM "voucher_batches_old";
DROP TABLE "voucher_batches_old";
CREATE INDEX idx_subjects_issuer_external ON subjects(issuer, external_id);
CREATE INDEX idx_vouchers_batch ON vouchers(batch_id, created_at DESC);
CREATE INDEX idx_vouchers_status ON vouchers(status);
CREATE INDEX idx_wallet_ledger_account ON wallet_ledger_entries(account_id, created_at DESC);
```

m0（guard/预检）/ m4（校验）语句文本：

```sql
	{Table: "vouchers", Column: "expires_at", Voucher: true},;
	{Table: "vouchers", Column: "redeemed_at", Voucher: true},;
	{Table: "wallet_accounts", Columns: []string{"created_at", "updated_at"}},;
	{Table: "wallet_ledger_entries", Columns: []string{"created_at"}},;
	{Table: "wallet_reconciliation_runs", Columns: []string{"created_at"}},;
	{Table: "subjects", Columns: []string{"created_at"}},;
	{Table: "vouchers", Columns: []string{"expires_at", "redeemed_at", "created_at", "updated_at"}},;
	{Table: "voucher_batches", Columns: []string{"created_at", "updated_at"}},;
```

<details><summary>postgres m1–m2</summary>

```sql
ALTER TABLE "wallet_accounts" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)));
ALTER TABLE "wallet_accounts" ALTER COLUMN "updated_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("updated_at"::double precision)));
ALTER TABLE "wallet_ledger_entries" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)));
ALTER TABLE "wallet_reconciliation_runs" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)));
ALTER TABLE "subjects" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)));
ALTER TABLE "vouchers" ALTER COLUMN "expires_at" TYPE timestamptz(6) USING (CASE WHEN "expires_at" IS NULL OR "expires_at" = 0 THEN NULL ELSE date_trunc('microseconds', to_timestamp("expires_at"::double precision)) END);
ALTER TABLE "vouchers" ALTER COLUMN "redeemed_at" TYPE timestamptz(6) USING (CASE WHEN "redeemed_at" IS NULL OR "redeemed_at" = 0 THEN NULL ELSE date_trunc('microseconds', to_timestamp("redeemed_at"::double precision)) END);
ALTER TABLE "vouchers" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)));
ALTER TABLE "vouchers" ALTER COLUMN "updated_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("updated_at"::double precision)));
ALTER TABLE "voucher_batches" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)));
ALTER TABLE "voucher_batches" ALTER COLUMN "updated_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("updated_at"::double precision)));
```

</details>

## v86 · `channel.telegram` · `vp040_temporal_telegram`

- `transform_id`: `0086:vp040-temporal-telegram:v1`
- **`MigrationChecksum`（真实值，D-017 单 checksum / SQLite 切片）**：`07a9a0ba61110b94acd4f4f7e87144fd297e8a99e44a273059080e07f1d172f1`
- 表范围（descriptor 顺序）：`telegram_config`, `telegram_sessions`, `telegram_inbound_messages`, `telegram_outbound_messages`
- FK 子女盘点（机械解析 v72 DDL 的 `REFERENCES`）：父表 （无）；摘除并建回的子表 （无）
- m0 预检列：`telegram_config.updated_at`；m0 retired-table guard：（无）
- PG 语句数 9；SQLite 语句数（m1–m3）20；m4 断言 4

canonical（m1–m3）：

```sql
ALTER TABLE "telegram_config" RENAME TO "telegram_config_old";
CREATE TABLE telegram_config (
  id                 INTEGER PRIMARY KEY CHECK (id = 1),
  bot_token_enc      TEXT    NOT NULL DEFAULT '',
  webhook_secret_enc TEXT    NOT NULL DEFAULT '',
  updated_at         TEXT
, mode TEXT NOT NULL DEFAULT 'polling', webhook_public_base_url TEXT NOT NULL DEFAULT '');
INSERT INTO "telegram_config" ("id", "bot_token_enc", "webhook_secret_enc", "updated_at", "mode", "webhook_public_base_url")
SELECT "id", "bot_token_enc", "webhook_secret_enc", CASE WHEN updated_at = 0 THEN NULL ELSE strftime('%Y-%m-%dT%H:%M:%S', updated_at, 'unixepoch') || '.000000Z' END, "mode", "webhook_public_base_url"
FROM "telegram_config_old";
DROP TABLE "telegram_config_old";
ALTER TABLE "telegram_sessions" RENAME TO "telegram_sessions_old";
CREATE TABLE telegram_sessions (
  bot_id          INTEGER NOT NULL,
  chat_id         INTEGER NOT NULL,
  chat_type       TEXT    NOT NULL DEFAULT '',
  title           TEXT    NOT NULL DEFAULT '',
  username        TEXT    NOT NULL DEFAULT '',
  last_message_at TEXT NOT NULL,
  created_at      TEXT NOT NULL,
  updated_at      TEXT NOT NULL,
  PRIMARY KEY (bot_id, chat_id)
);
INSERT INTO "telegram_sessions" ("bot_id", "chat_id", "chat_type", "title", "username", "last_message_at", "created_at", "updated_at")
SELECT "bot_id", "chat_id", "chat_type", "title", "username", strftime('%Y-%m-%dT%H:%M:%S', last_message_at, 'unixepoch') || '.000000Z', strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z', strftime('%Y-%m-%dT%H:%M:%S', updated_at, 'unixepoch') || '.000000Z'
FROM "telegram_sessions_old";
DROP TABLE "telegram_sessions_old";
ALTER TABLE "telegram_inbound_messages" RENAME TO "telegram_inbound_messages_old";
CREATE TABLE telegram_inbound_messages (
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
);
INSERT INTO "telegram_inbound_messages" ("bot_id", "update_id", "chat_id", "user_id", "message_id", "callback_query_id", "direction", "message_kind", "text", "callback_data", "sender_username", "received_at")
SELECT "bot_id", "update_id", "chat_id", "user_id", "message_id", "callback_query_id", "direction", "message_kind", "text", "callback_data", "sender_username", strftime('%Y-%m-%dT%H:%M:%S', received_at, 'unixepoch') || '.000000Z'
FROM "telegram_inbound_messages_old";
DROP TABLE "telegram_inbound_messages_old";
ALTER TABLE "telegram_outbound_messages" RENAME TO "telegram_outbound_messages_old";
CREATE TABLE telegram_outbound_messages (
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
 );
INSERT INTO "telegram_outbound_messages" ("bot_id", "request_id", "retry_root", "retry_of", "chat_id", "text", "status", "error_message", "created_at", "updated_at")
SELECT "bot_id", "request_id", "retry_root", "retry_of", "chat_id", "text", "status", "error_message", strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z', strftime('%Y-%m-%dT%H:%M:%S', updated_at, 'unixepoch') || '.000000Z'
FROM "telegram_outbound_messages_old";
DROP TABLE "telegram_outbound_messages_old";
CREATE INDEX idx_telegram_inbound_messages_chat_received
  ON telegram_inbound_messages (bot_id, chat_id, received_at DESC, update_id DESC);
CREATE INDEX idx_telegram_outbound_messages_chat_created
  ON telegram_outbound_messages (bot_id, chat_id, created_at DESC, request_id DESC);
CREATE UNIQUE INDEX idx_telegram_outbound_messages_pending_root
  ON telegram_outbound_messages (bot_id, retry_root) WHERE status = 'pending';
CREATE INDEX idx_telegram_sessions_activity
  ON telegram_sessions (bot_id, last_message_at DESC, chat_id DESC);
```

m0（guard/预检）/ m4（校验）语句文本：

```sql
	{Table: "telegram_config", Column: "updated_at", Voucher: false},;
	{Table: "telegram_config", Columns: []string{"updated_at"}},;
	{Table: "telegram_sessions", Columns: []string{"last_message_at", "created_at", "updated_at"}},;
	{Table: "telegram_inbound_messages", Columns: []string{"received_at"}},;
	{Table: "telegram_outbound_messages", Columns: []string{"created_at", "updated_at"}},;
```

<details><summary>postgres m1–m2</summary>

```sql
ALTER TABLE "telegram_config" ALTER COLUMN "updated_at" DROP DEFAULT;
ALTER TABLE "telegram_config" ALTER COLUMN "updated_at" DROP NOT NULL;
ALTER TABLE "telegram_config" ALTER COLUMN "updated_at" TYPE timestamptz(6) USING (CASE WHEN "updated_at" = 0 THEN NULL ELSE date_trunc('microseconds', to_timestamp("updated_at"::double precision)) END);
ALTER TABLE "telegram_sessions" ALTER COLUMN "last_message_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("last_message_at"::double precision)));
ALTER TABLE "telegram_sessions" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)));
ALTER TABLE "telegram_sessions" ALTER COLUMN "updated_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("updated_at"::double precision)));
ALTER TABLE "telegram_inbound_messages" ALTER COLUMN "received_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("received_at"::double precision)));
ALTER TABLE "telegram_outbound_messages" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)));
ALTER TABLE "telegram_outbound_messages" ALTER COLUMN "updated_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("updated_at"::double precision)));
```

</details>

## v87 · `biz.digital-offer` · `vp040_temporal_digital_offer`

- `transform_id`: `0087:vp040-temporal-digital-offer:v1`
- **`MigrationChecksum`（真实值，D-017 单 checksum / SQLite 切片）**：`753b22027027066bd54b8909974b2f867ac5340bd513ce861bf1e2553f1fc7e4`
- 表范围（descriptor 顺序）：`digital_offers`, `digital_purchases`, `digital_entitlements`
- FK 子女盘点（机械解析 v72 DDL 的 `REFERENCES`）：父表 （无）；摘除并建回的子表 （无）
- m0 预检列：（无）；m0 retired-table guard：（无）
- PG 语句数 6；SQLite 语句数（m1–m3）17；m4 断言 3

canonical（m1–m3）：

```sql
ALTER TABLE "digital_offers" RENAME TO "digital_offers_old";
CREATE TABLE digital_offers (
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
);
INSERT INTO "digital_offers" ("id", "name", "description", "price_amount", "currency", "entitlement_form", "duration_seconds", "count_per_purchase", "status", "version", "created_at", "updated_at")
SELECT "id", "name", "description", "price_amount", "currency", "entitlement_form", "duration_seconds", "count_per_purchase", "status", "version", strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z', strftime('%Y-%m-%dT%H:%M:%S', updated_at, 'unixepoch') || '.000000Z'
FROM "digital_offers_old";
DROP TABLE "digital_offers_old";
ALTER TABLE "digital_purchases" RENAME TO "digital_purchases_old";
CREATE TABLE digital_purchases (
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
);
INSERT INTO "digital_purchases" ("id", "subject_id", "offer_id", "offer_name", "amount", "currency", "freeze_entry_id", "deduct_entry_id", "request_id", "status", "created_at")
SELECT "id", "subject_id", "offer_id", "offer_name", "amount", "currency", "freeze_entry_id", "deduct_entry_id", "request_id", "status", strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z'
FROM "digital_purchases_old";
DROP TABLE "digital_purchases_old";
ALTER TABLE "digital_entitlements" RENAME TO "digital_entitlements_old";
CREATE TABLE digital_entitlements (
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
);
INSERT INTO "digital_entitlements" ("id", "subject_id", "offer_id", "purchase_id", "form", "expires_at", "remaining_count", "status", "created_at", "updated_at")
SELECT "id", "subject_id", "offer_id", "purchase_id", "form", CASE WHEN expires_at IS NULL THEN NULL ELSE strftime('%Y-%m-%dT%H:%M:%S', expires_at, 'unixepoch') || '.000000Z' END, "remaining_count", "status", strftime('%Y-%m-%dT%H:%M:%S', created_at, 'unixepoch') || '.000000Z', strftime('%Y-%m-%dT%H:%M:%S', updated_at, 'unixepoch') || '.000000Z'
FROM "digital_entitlements_old";
DROP TABLE "digital_entitlements_old";
CREATE INDEX idx_digital_entitlements_purchase ON digital_entitlements(purchase_id);
CREATE INDEX idx_digital_entitlements_subject ON digital_entitlements(subject_id, status);
CREATE INDEX idx_digital_offers_status ON digital_offers(status, created_at DESC);
CREATE INDEX idx_digital_purchases_offer ON digital_purchases(offer_id);
CREATE INDEX idx_digital_purchases_subject ON digital_purchases(subject_id, created_at DESC);
```

m0（guard/预检）/ m4（校验）语句文本：

```sql
	{Table: "digital_offers", Columns: []string{"created_at", "updated_at"}},;
	{Table: "digital_purchases", Columns: []string{"created_at"}},;
	{Table: "digital_entitlements", Columns: []string{"expires_at", "created_at", "updated_at"}},;
```

<details><summary>postgres m1–m2</summary>

```sql
ALTER TABLE "digital_offers" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)));
ALTER TABLE "digital_offers" ALTER COLUMN "updated_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("updated_at"::double precision)));
ALTER TABLE "digital_purchases" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)));
ALTER TABLE "digital_entitlements" ALTER COLUMN "expires_at" TYPE timestamptz(6) USING (CASE WHEN "expires_at" IS NULL THEN NULL ELSE date_trunc('microseconds', to_timestamp("expires_at"::double precision)) END);
ALTER TABLE "digital_entitlements" ALTER COLUMN "created_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("created_at"::double precision)));
ALTER TABLE "digital_entitlements" ALTER COLUMN "updated_at" TYPE timestamptz(6) USING (date_trunc('microseconds', to_timestamp("updated_at"::double precision)));
```

</details>
