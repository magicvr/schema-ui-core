---
id: r1-c2-per-table-rebuild-ddl-v1.0-fc
doc_type: design-attachment
title: R1 C2 逐表 exact rebuild DDL v1.0（冻结候选）
status: freeze-candidate
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-002-r1-contract-and-denominator-freeze
version: 1.0.0
---

# R1 C2 逐表 exact rebuild DDL v1.0（冻结候选）

> **状态：`freeze-candidate`。** 本文件是 A-030/A-032 **F-I-002.1** 点名的**逐表 exact SQLite rebuild DDL 正文**（`CREATE TABLE` / `INSERT … SELECT` / `CREATE INDEX`），按用户 `D-019` 裁决的 **F-5 子女先行**模式与 `D-018` 选项 C 的行拷贝机制书写，表达式取 `r1-c2-sqlite-rebuild-mechanism-v1.0-fc.md` §2（已由 A-032 复证）。**不是实施证据**；R2 才落码。
>
> 覆盖：`D-014` allocation 的 v73–v87 共 15 个 descriptor / 20 张时间列表，加 `D-019` §2/§3 的 **FK-preserve 子表**（3 张联接表 + 3 张跨 descriptor 子表）。
>
> **PG 侧**由配套附件承载（`D-019` §6：必须**显式书写**，禁止 `pgTimeColRe` 派生）。

## 0. 全局约定

| 项 | 约定 |
|----|------|
| 重建模式 | 有 FK 子表的父表走 **F-5 十步**（`D-019` §1）；无 FK 子表的表走**裸四步** |
| 索引时机 | **全部** `CREATE INDEX` / `CREATE UNIQUE INDEX` 在 `DROP TABLE …_old`（及子表回填）**之后**；SQLite 把索引名留给被重命名的表，提前建会撞名 |
| 秒族表达式 | `strftime('%Y-%m-%dT%H:%M:%S', <col>, 'unixepoch') \|\| '.000000Z'` |
| 毫秒族表达式 | `strftime('%Y-%m-%dT%H:%M:%S', CASE WHEN <col> >= 0 THEN <col>/1000 ELSE (<col>-999)/1000 END, 'unixepoch') \|\| '.' \|\| printf('%03d', (<col>%1000 + 1000) % 1000) \|\| '000Z'` |
| NULL 包裹 | 可空列：`CASE WHEN <col> IS NULL THEN NULL ELSE <表达式> END` |
| D0 包裹 | `DEFAULT 0` sentinel 列：`CASE WHEN <col> = 0 THEN NULL ELSE <表达式> END`；新 DDL **去 `NOT NULL`、去 `DEFAULT 0`** |
| voucher 包裹 | `#72/#73`：`CASE WHEN <col> IS NULL OR <col> = 0 THEN NULL ELSE <表达式> END` |
| 负值 `< 0` | **不进任何表达式**；只由 `m0` 预检 fail closed |
| 列序权威 | live `PRAGMA table_info`（v72 已 apply 库）；ALTER 追加列在末尾。**禁止**按模块文件分组拼列序 |
| 子表 live DDL 权威 | v72 已 apply 库的 `sqlite_master.sql` + 索引清单（**含** ALTER 追加列） |
| DDL 附加校验（进 `m4`） | `foreign_key_check` 无行；`integrity_check` = `ok`；`<t>_old` 不存在；目标列类型 = `TEXT`；**父表另加**：全部子表 `REFERENCES` 指向同名新父表（不含 `_old`） |

## 1. v74（`core.auth-session`）的 F-5 编排

v74 **同时重建 4 张 FK 父表**（`users`、`roles`、`permissions`、`menu_items`），其子表**互相交织**（联接表引用两张父表），因此必须**先把全部子表摘掉，再重建全部父表，最后统一建回子表**。

### 1.1 子表分组

| 组 | 表 | 含时间列 | 建回时的形状 |
|----|----|---------|-------------|
| **A · 本 descriptor 转换组** | `refresh_tokens`、`email_verification_challenges`、`password_recovery_challenges`、`login_failures`、`user_password_history`、`user_invites` | 是 | 转换后（`TEXT`） |
| **B · 本 descriptor FK-preserve 组** | `user_roles`、`role_permissions`、`role_menu_items` | **否** | 原样 |
| **C · 跨 descriptor 组** | `notifications`（v81）、`user_mfa`、`mfa_proofs`（v80） | 是 | **原样（时间列仍 `INTEGER`）** |

### 1.2 v74 语句顺序（单事务）

```text
-- 阶段 1：备份并摘除全部 12 张子表
CREATE TEMP TABLE refresh_tokens_bak AS SELECT * FROM refresh_tokens;
CREATE TEMP TABLE email_verification_challenges_bak AS SELECT * FROM email_verification_challenges;
CREATE TEMP TABLE password_recovery_challenges_bak AS SELECT * FROM password_recovery_challenges;
CREATE TEMP TABLE login_failures_bak AS SELECT * FROM login_failures;
CREATE TEMP TABLE user_password_history_bak AS SELECT * FROM user_password_history;
CREATE TEMP TABLE user_invites_bak AS SELECT * FROM user_invites;
CREATE TEMP TABLE user_roles_bak AS SELECT * FROM user_roles;
CREATE TEMP TABLE role_permissions_bak AS SELECT * FROM role_permissions;
CREATE TEMP TABLE role_menu_items_bak AS SELECT * FROM role_menu_items;
CREATE TEMP TABLE notifications_bak AS SELECT * FROM notifications;
CREATE TEMP TABLE user_mfa_bak AS SELECT * FROM user_mfa;
CREATE TEMP TABLE mfa_proofs_bak AS SELECT * FROM mfa_proofs;
DROP TABLE refresh_tokens; DROP TABLE email_verification_challenges;
DROP TABLE password_recovery_challenges; DROP TABLE login_failures;
DROP TABLE user_password_history; DROP TABLE user_invites;
DROP TABLE user_roles; DROP TABLE role_permissions; DROP TABLE role_menu_items;
DROP TABLE notifications; DROP TABLE user_mfa; DROP TABLE mfa_proofs;

-- 阶段 2：重建 4 张父表（此时无任何子表可被改写）——正文见 §3
-- 阶段 3：建回并回填 A 组（转换）——正文见 §3
-- 阶段 4：按 live DDL 建回并回填 B 组（FK-preserve，无转换）——正文见 §3.10
-- 阶段 5：按 live DDL 建回并回填 C 组（FK 修复，时间列仍 INTEGER）——正文见 §3.11
-- 阶段 6：全部索引最后建
-- 阶段 7：校验（§0 末行）
```

> **`users` 的 FK 子表完整清单（10 张）**：`refresh_tokens`、`email_verification_challenges`、`password_recovery_challenges`、`login_failures`、`user_password_history`、`user_invites`、`user_roles`、`notifications`、`user_mfa`、`mfa_proofs`——阶段 1 必须**全部**摘除，漏一张即悬挂 FK 或丢行。

### 1.3 C 组的时间列转换留给 v80/v81

- `notifications` 由 **v81** 转换；`user_mfa`/`mfa_proofs` 由 **v80**；届时**无人引用它们**，可用裸四步。
- 代价（已知并接受）：这三张表被重建两次。

## 2. `core.persistence` / `core.jobs` / `admin.*` 逐表 DDL

### 2.1 `schema_migrations`（v73；无 FK 子表）

```sql
-- legacy（internal/store/identity.go:58；authsession/migration.go:18 为 v1，禁止改）
CREATE TABLE IF NOT EXISTS schema_migrations (
  version    INTEGER PRIMARY KEY,
  name       TEXT NOT NULL UNIQUE,
  checksum   TEXT NOT NULL CHECK (length(checksum) = 64),
  applied_at INTEGER NOT NULL
)
-- new
CREATE TABLE schema_migrations (
  version    INTEGER PRIMARY KEY,
  name       TEXT NOT NULL UNIQUE,
  checksum   TEXT NOT NULL CHECK (length(checksum) = 64),
  applied_at TEXT NOT NULL
)
```
```sql
ALTER TABLE schema_migrations RENAME TO schema_migrations_old;
CREATE TABLE schema_migrations (
  version    INTEGER PRIMARY KEY,
  name       TEXT NOT NULL UNIQUE,
  checksum   TEXT NOT NULL CHECK (length(checksum) = 64),
  applied_at TEXT NOT NULL
);
INSERT INTO schema_migrations (version, name, checksum, applied_at)
SELECT version, name, checksum, <秒表达式(applied_at)> FROM schema_migrations_old;
DROP TABLE schema_migrations_old;
-- 无索引
```

### 2.2 `mail_outbox`（v73；毫秒族；`channel`/`delivery_status` 为 v60 ALTER 追加）

```sql
-- legacy（corepersistence/migration/migration.go:53-59 + v60 ALTERs）
CREATE TABLE mail_outbox (
  id         TEXT PRIMARY KEY,
  to_addr    TEXT NOT NULL,
  subject    TEXT NOT NULL,
  body       TEXT NOT NULL,
  created_at INTEGER NOT NULL
)
ALTER TABLE mail_outbox ADD COLUMN channel TEXT NOT NULL DEFAULT 'mock'
ALTER TABLE mail_outbox ADD COLUMN delivery_status TEXT NOT NULL DEFAULT 'delivered'
-- new
CREATE TABLE mail_outbox (
  id              TEXT PRIMARY KEY,
  to_addr         TEXT NOT NULL,
  subject         TEXT NOT NULL,
  body            TEXT NOT NULL,
  created_at      TEXT NOT NULL,
  channel         TEXT NOT NULL DEFAULT 'mock',
  delivery_status TEXT NOT NULL DEFAULT 'delivered'
)
```
```sql
ALTER TABLE mail_outbox RENAME TO mail_outbox_old;
CREATE TABLE mail_outbox ( …new，见上… );
INSERT INTO mail_outbox (id, to_addr, subject, body, created_at, channel, delivery_status)
SELECT id, to_addr, subject, body, <毫秒表达式(created_at)>, channel, delivery_status
FROM mail_outbox_old;
DROP TABLE mail_outbox_old;
CREATE INDEX idx_mail_outbox_created_at ON mail_outbox(created_at);
```

### 2.3 `mail_config`（v73；毫秒族；`#34` D0）

```sql
-- legacy（:94-106）
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
  updated_at         INTEGER NOT NULL DEFAULT 0
)
-- new：仅 updated_at 改 TEXT（可空、无默认）
```
```sql
ALTER TABLE mail_config RENAME TO mail_config_old;
CREATE TABLE mail_config ( …同上，updated_at TEXT… );
INSERT INTO mail_config (id, channel, mock_retention, resend_from, resend_api_key_enc, smtp_host, smtp_port,
                         smtp_username, smtp_password_enc, smtp_from, updated_at)
SELECT id, channel, mock_retention, resend_from, resend_api_key_enc, smtp_host, smtp_port,
       smtp_username, smtp_password_enc, smtp_from,
       CASE WHEN updated_at = 0 THEN NULL ELSE <毫秒表达式(updated_at)> END
FROM mail_config_old;
DROP TABLE mail_config_old;
-- 无索引
```

### 2.4 `jobs`（v76；毫秒族；无 FK 子表）

- legacy 见 `jobs/migration/migration.go:15-44`（含六态 CHECK）；new 仅把 `lease_expires_at`/`created_at`/`updated_at`/`finished_at`/`expires_at` 改 `TEXT`（前两者按可空性），**六态 CHECK 与其余列逐字保留**。

```sql
ALTER TABLE jobs RENAME TO jobs_old;
CREATE TABLE jobs ( …new… );
INSERT INTO jobs (id, kind, status, payload, progress, cancel_requested, attempt, max_attempts, lease_owner, lease_version,
                  lease_expires_at, result, error_code, error_message, actor_id, correlation_id, created_at, updated_at, finished_at, expires_at)
SELECT id, kind, status, payload, progress, cancel_requested, attempt, max_attempts, lease_owner, lease_version,
  CASE WHEN lease_expires_at IS NULL THEN NULL ELSE <毫秒表达式(lease_expires_at)> END,
  result, error_code, error_message, actor_id, correlation_id,
  <毫秒表达式(created_at)>, <毫秒表达式(updated_at)>,
  CASE WHEN finished_at IS NULL THEN NULL ELSE <毫秒表达式(finished_at)> END,
  CASE WHEN expires_at IS NULL THEN NULL ELSE <毫秒表达式(expires_at)> END
FROM jobs_old;
DROP TABLE jobs_old;
CREATE INDEX idx_jobs_runnable ON jobs(status, cancel_requested, lease_expires_at, created_at);
CREATE INDEX idx_jobs_actor ON jobs(actor_id, kind, updated_at DESC);
CREATE INDEX idx_jobs_expiry ON jobs(status, expires_at);
CREATE INDEX IF NOT EXISTS idx_jobs_created_at ON jobs(created_at DESC, id DESC);
```

### 2.5 `dict_types` + `dict_entries`（v77；父 + 子，F-5）

```sql
-- legacy：dict_types :20-29；dict_entries :30-41 + v39 ALTER :76（badge_style 在末尾）
CREATE TABLE dict_types (
  id         TEXT PRIMARY KEY,
  key        TEXT NOT NULL UNIQUE,
  name       TEXT NOT NULL,
  enabled    INTEGER NOT NULL DEFAULT 1,
  description TEXT,
  sort       INTEGER NOT NULL DEFAULT 0,
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL
)
CREATE TABLE dict_entries (
  id         TEXT PRIMARY KEY,
  dict_key   TEXT NOT NULL REFERENCES dict_types(key) ON DELETE CASCADE,
  entry_key  TEXT NOT NULL,
  label      TEXT NOT NULL,
  enabled    INTEGER NOT NULL DEFAULT 1,
  sort       INTEGER NOT NULL DEFAULT 0,
  remark     TEXT,
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL,
  UNIQUE (dict_key, entry_key)
)
ALTER TABLE dict_entries ADD COLUMN badge_style TEXT NOT NULL DEFAULT 'default'
```
```sql
-- F-5（已实测：2/2 行保留、FK 文本恢复、foreign_key_check 空）
CREATE TEMP TABLE dict_entries_bak AS SELECT * FROM dict_entries;
DROP TABLE dict_entries;
ALTER TABLE dict_types RENAME TO dict_types_old;
CREATE TABLE dict_types ( …new：created_at/updated_at → TEXT NOT NULL… );
INSERT INTO dict_types (id, key, name, enabled, description, sort, created_at, updated_at)
SELECT id, key, name, enabled, description, sort, <秒表达式(created_at)>, <秒表达式(updated_at)> FROM dict_types_old;
DROP TABLE dict_types_old;
CREATE TABLE dict_entries ( …new：created_at/updated_at → TEXT NOT NULL；badge_style 保留… );
INSERT INTO dict_entries (id, dict_key, entry_key, label, enabled, sort, remark, created_at, updated_at, badge_style)
SELECT id, dict_key, entry_key, label, enabled, sort, remark,
       <秒表达式(created_at)>, <秒表达式(updated_at)>, badge_style
FROM dict_entries_bak;
DROP TABLE dict_entries_bak;
CREATE INDEX idx_dict_entries_dict_key ON dict_entries(dict_key, sort);
```

### 2.6 `data_scope_policies` / `user_data_scopes`（v78；**两表均无 FK 子表** → 裸四步）

> 响应 A-034 **F-I-023**（v78 整段此前缺席）。以下 legacy/新 DDL 均为 **live `sqlite_master` 实测逐字**（非推断）。

```sql
-- legacy（datadictionary 之外的 datapermission 模块；live sqlite_master 实测）
CREATE TABLE data_scope_policies (
  resource      TEXT PRIMARY KEY,
  owner_column  TEXT NOT NULL,
  default_scope TEXT NOT NULL CHECK (default_scope IN ('all','self')),
  enabled       INTEGER NOT NULL DEFAULT 1,
  updated_at    INTEGER NOT NULL
)
CREATE TABLE user_data_scopes (
  user_id    TEXT NOT NULL,
  resource   TEXT NOT NULL,
  scope_type TEXT NOT NULL CHECK (scope_type IN ('all','self')),
  updated_at INTEGER NOT NULL,
  PRIMARY KEY (user_id, resource)
)
-- new（仅 updated_at → TEXT NOT NULL；CHECK / PK / DEFAULT 逐字保留）
CREATE TABLE data_scope_policies (
  resource      TEXT PRIMARY KEY,
  owner_column  TEXT NOT NULL,
  default_scope TEXT NOT NULL CHECK (default_scope IN ('all','self')),
  enabled       INTEGER NOT NULL DEFAULT 1,
  updated_at    TEXT NOT NULL
)
CREATE TABLE user_data_scopes (
  user_id    TEXT NOT NULL,
  resource   TEXT NOT NULL,
  scope_type TEXT NOT NULL CHECK (scope_type IN ('all','self')),
  updated_at TEXT NOT NULL,
  PRIMARY KEY (user_id, resource)
)
```

```sql
-- data_scope_policies（裸四步；无显式索引，仅 PK 自动索引）
ALTER TABLE data_scope_policies RENAME TO data_scope_policies_old;
CREATE TABLE data_scope_policies (
  resource      TEXT PRIMARY KEY,
  owner_column  TEXT NOT NULL,
  default_scope TEXT NOT NULL CHECK (default_scope IN ('all','self')),
  enabled       INTEGER NOT NULL DEFAULT 1,
  updated_at    TEXT NOT NULL
);
INSERT INTO data_scope_policies (resource, owner_column, default_scope, enabled, updated_at)
SELECT resource, owner_column, default_scope, enabled, <秒表达式(updated_at)>
FROM data_scope_policies_old;
DROP TABLE data_scope_policies_old;

-- user_data_scopes（裸四步；无显式索引，仅 PK 自动索引）
ALTER TABLE user_data_scopes RENAME TO user_data_scopes_old;
CREATE TABLE user_data_scopes (
  user_id    TEXT NOT NULL,
  resource   TEXT NOT NULL,
  scope_type TEXT NOT NULL CHECK (scope_type IN ('all','self')),
  updated_at TEXT NOT NULL,
  PRIMARY KEY (user_id, resource)
);
INSERT INTO user_data_scopes (user_id, resource, scope_type, updated_at)
SELECT user_id, resource, scope_type, <秒表达式(updated_at)>
FROM user_data_scopes_old;
DROP TABLE user_data_scopes_old;
-- 两表均无显式 CREATE INDEX 语句（PK 自动索引由 PRIMARY KEY 子句自动重建）
```

> 子表盘点（实测）：**无任何表**的存留 DDL 含 `REFERENCES data_scope_policies` 或 `REFERENCES user_data_scopes`（count = 0）→ 裸四步安全。

### 2.7 `captcha_challenges` / `captcha_config`（v79；均无 FK）

```sql
-- legacy（logincaptcha/migration.go:19-24 / :25-30）
CREATE TABLE captcha_challenges (
  id          TEXT PRIMARY KEY,
  answer_hash TEXT NOT NULL,
  expires_at  INTEGER NOT NULL,
  created_at  INTEGER NOT NULL
)
CREATE TABLE captcha_config (
  id         INTEGER PRIMARY KEY CHECK (id = 1),
  enabled    INTEGER NOT NULL DEFAULT 0,
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL
)
```
```sql
-- 两表均裸四步；时间列 → TEXT；captcha_config.id 保持 INTEGER PRIMARY KEY CHECK (id = 1)；均无索引
```

### 2.8 `user_mfa` / `mfa_proofs`（v80；裸四步）

- legacy：`mfa/migration.go:21-29` / `:30-36`（**FK 子句逐字保留**；`last_used_step` 非时间列）。
- new：`created_at`/`updated_at` / `expires_at`/`created_at` → `TEXT NOT NULL`；均无索引。

### 2.9 `notifications`（v81；裸四步）

```sql
-- legacy 有效形状（notifications/migration.go:16-24 + v37 ALTERs :53-54）
CREATE TABLE notifications (
  id         TEXT PRIMARY KEY,
  user_id    TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  event      TEXT NOT NULL CHECK (event IN ('account.locked','account.disabled','account.unlocked','account.password-changed')),
  title      TEXT NOT NULL,
  body       TEXT NOT NULL,
  read_at    INTEGER,
  created_at INTEGER NOT NULL,
  title_key  TEXT,
  body_key   TEXT
)
```
```sql
ALTER TABLE notifications RENAME TO notifications_old;
CREATE TABLE notifications ( …同上，read_at → TEXT 可空，created_at → TEXT NOT NULL… );
INSERT INTO notifications (id, user_id, event, title, body, read_at, created_at, title_key, body_key)
SELECT id, user_id, event, title, body,
       CASE WHEN read_at IS NULL THEN NULL ELSE <秒表达式(read_at)> END,
       <秒表达式(created_at)>, title_key, body_key
FROM notifications_old;
DROP TABLE notifications_old;
CREATE INDEX idx_notifications_user_created ON notifications(user_id, created_at DESC);
```

### 2.10 `recycle_items`（v82；无 FK）

```sql
-- legacy（recyclebin/migration.go:20-29）；索引 :30-31
CREATE TABLE recycle_items (
  id          TEXT PRIMARY KEY,
  resource    TEXT NOT NULL,
  resource_id TEXT NOT NULL,
  payload     TEXT NOT NULL,
  actor_id    TEXT NOT NULL,
  actor_name  TEXT NOT NULL,
  deleted_at  INTEGER NOT NULL,
  restored_at INTEGER
)
CREATE UNIQUE INDEX idx_recycle_items_active ON recycle_items(resource, resource_id) WHERE restored_at IS NULL
CREATE INDEX idx_recycle_items_deleted_at ON recycle_items(deleted_at DESC)
```
```sql
-- 裸四步；deleted_at → TEXT NOT NULL；restored_at → TEXT 可空（NULL 包裹）；索引最后逐字重建
```

### 2.11 `scheduled_tasks` + `task_runs`（v83；父 + 子，F-5）

```sql
-- legacy（scheduledtasks/migration.go:18-28 / :29-37；索引 :38）
CREATE TABLE scheduled_tasks (
  id          TEXT PRIMARY KEY,
  key         TEXT NOT NULL UNIQUE,
  cron        TEXT NOT NULL,
  name        TEXT NOT NULL,
  enabled     INTEGER NOT NULL DEFAULT 1,
  description TEXT,
  handler     TEXT NOT NULL DEFAULT 'system.noop',
  created_at  INTEGER NOT NULL,
  updated_at  INTEGER NOT NULL
)
CREATE TABLE task_runs (
  id          TEXT PRIMARY KEY,
  task_id     TEXT NOT NULL REFERENCES scheduled_tasks(id) ON DELETE CASCADE,
  status      TEXT NOT NULL CHECK (status IN ('ran','failed')),
  started_at  INTEGER NOT NULL,
  finished_at INTEGER,
  detail      TEXT,
  created_at  INTEGER NOT NULL
)
CREATE INDEX idx_task_runs_task_started ON task_runs(task_id, started_at DESC)
```
```sql
CREATE TEMP TABLE task_runs_bak AS SELECT * FROM task_runs;
DROP TABLE task_runs;
-- rebuild scheduled_tasks（裸四步，时间列 → TEXT）
-- create task_runs（时间列 → TEXT；finished_at 可空）
INSERT INTO task_runs (id, task_id, status, started_at, finished_at, detail, created_at)
SELECT id, task_id, status, <秒表达式(started_at)>,
       CASE WHEN finished_at IS NULL THEN NULL ELSE <秒表达式(finished_at)> END,
       detail, <秒表达式(created_at)>
FROM task_runs_bak;
DROP TABLE task_runs_bak;
CREATE INDEX idx_task_runs_task_started ON task_runs(task_id, started_at DESC);
```

### 2.12 `site_settings`（v84；无 FK；15 列）

> 响应 A-034 **F-I-024**：本节列序改为 **live `PRAGMA table_info` 实测**（cid 0–14），并**明确声明实测结论与按文件行号推断一致**——`operation_log_retention_days`(12) / `operation_log_expiration_action`(13) 在 `default_currency`(14) **之前**。给出完整 15 列 `CREATE TABLE` 正文（此前只有差异说明）。实测 `sqlite_master` 的存储文本显示 ALTER 追加段内 `default_currency` 位于末尾，故 live cid 序 = 基座 4 列 + 追加 11 列的**追加顺序**。

```sql
-- legacy（live 实测，15 列，cid 0-14）
CREATE TABLE site_settings (
  id                              TEXT PRIMARY KEY CHECK (id = 'default'),
  site_title                      TEXT NOT NULL,
  logo_url                        TEXT NOT NULL DEFAULT '',
  updated_at                      INTEGER NOT NULL,
  logo_url_light                  TEXT NOT NULL DEFAULT '',
  logo_url_dark                   TEXT NOT NULL DEFAULT '',
  favicon_url                     TEXT NOT NULL DEFAULT '',
  default_locale                  TEXT NOT NULL DEFAULT '',
  site_timezone                   TEXT NOT NULL DEFAULT '',
  default_theme                   TEXT NOT NULL DEFAULT '',
  copyright_text                  TEXT NOT NULL DEFAULT '',
  icp_number                      TEXT NOT NULL DEFAULT '',
  operation_log_retention_days    INTEGER NOT NULL DEFAULT 90,
  operation_log_expiration_action TEXT NOT NULL DEFAULT 'archive',
  default_currency                TEXT NOT NULL DEFAULT ''
)
-- new（仅 updated_at → TEXT NOT NULL；其余列、DEFAULT、CHECK 逐字保留）
CREATE TABLE site_settings (
  id                              TEXT PRIMARY KEY CHECK (id = 'default'),
  site_title                      TEXT NOT NULL,
  logo_url                        TEXT NOT NULL DEFAULT '',
  updated_at                      TEXT NOT NULL,
  logo_url_light                  TEXT NOT NULL DEFAULT '',
  logo_url_dark                   TEXT NOT NULL DEFAULT '',
  favicon_url                     TEXT NOT NULL DEFAULT '',
  default_locale                  TEXT NOT NULL DEFAULT '',
  site_timezone                   TEXT NOT NULL DEFAULT '',
  default_theme                   TEXT NOT NULL DEFAULT '',
  copyright_text                  TEXT NOT NULL DEFAULT '',
  icp_number                      TEXT NOT NULL DEFAULT '',
  operation_log_retention_days    INTEGER NOT NULL DEFAULT 90,
  operation_log_expiration_action TEXT NOT NULL DEFAULT 'archive',
  default_currency                TEXT NOT NULL DEFAULT ''
)
```

```sql
ALTER TABLE site_settings RENAME TO site_settings_old;
CREATE TABLE site_settings ( …new，见上，15 列… );
INSERT INTO site_settings (id, site_title, logo_url, updated_at, logo_url_light, logo_url_dark, favicon_url,
                           default_locale, site_timezone, default_theme, copyright_text, icp_number,
                           operation_log_retention_days, operation_log_expiration_action, default_currency)
SELECT id, site_title, logo_url, <秒表达式(updated_at)>, logo_url_light, logo_url_dark, favicon_url,
       default_locale, site_timezone, default_theme, copyright_text, icp_number,
       operation_log_retention_days, operation_log_expiration_action, default_currency
FROM site_settings_old;
DROP TABLE site_settings_old;
CREATE INDEX IF NOT EXISTS idx_site_settings_updated_at ON site_settings (updated_at);
```

> **Go seeder 联动**：`settings/migration/migration.go:69-72`（SQLite）与 `:52-55`（PG）以 Unix 秒插入 `updated_at`，**必须与 v84 同批**改为目标格式（`D-018` 无过渡期）。

## 3. `authsession`（v74）逐表 DDL 正文

> §1.2 的 F-5 编排已**整链实测通过**，见 §8。

### 3.1 `system_data_reconcile`（NN；裸四步）

- legacy `authsession/migration.go:92-100`；new 仅 `applied_at` → `TEXT NOT NULL`；无索引。

### 3.2 `users`（FK 父表 → 阶段 2；**17 列**，列序 = live `PRAGMA table_info`）

```sql
CREATE TABLE users (
  id                    TEXT PRIMARY KEY,
  username              TEXT NOT NULL UNIQUE,
  name                  TEXT NOT NULL,
  roles                 TEXT NOT NULL,
  password_hash         TEXT NOT NULL,
  created_at            TEXT NOT NULL,
  updated_at            TEXT NOT NULL,
  token_version         INTEGER NOT NULL DEFAULT 0,
  failed_login_count    INTEGER NOT NULL DEFAULT 0,
  locked_until          TEXT,
  enabled               INTEGER NOT NULL DEFAULT 1,
  notifications_enabled INTEGER NOT NULL DEFAULT 1,
  avatar_url            TEXT NOT NULL DEFAULT '',
  must_change_password  INTEGER NOT NULL DEFAULT 0,
  email                 TEXT,
  email_status          TEXT CHECK (email_status IN ('pending','verified')),
  last_login_failure_at TEXT
)
```
```sql
ALTER TABLE users RENAME TO users_old;
CREATE TABLE users ( …new，见上… );
INSERT INTO users (id, username, name, roles, password_hash, created_at, updated_at, token_version,
                   failed_login_count, locked_until, enabled, notifications_enabled, avatar_url,
                   must_change_password, email, email_status, last_login_failure_at)
SELECT id, username, name, roles, password_hash,
       <秒表达式(created_at)>, <秒表达式(updated_at)>,
       token_version, failed_login_count,
       CASE WHEN locked_until = 0 THEN NULL ELSE <秒表达式(locked_until)> END,
       enabled, notifications_enabled, avatar_url, must_change_password, email, email_status,
       CASE WHEN last_login_failure_at = 0 THEN NULL ELSE <秒表达式(last_login_failure_at)> END
FROM users_old;
DROP TABLE users_old;
CREATE UNIQUE INDEX idx_users_email_lower ON users(lower(email));
```

> `users.updated_at` 运行时写入按 **Root** `D-013` 单调规则 `max(Truncate(µs) now, old + 1µs)`。

### 3.3 `roles` / `permissions` / `menu_items`（FK 父表 → 阶段 2）

- legacy：`roles` `:47-54`、`permissions` `:61-67`、`menu_items` `:74-82`；三表均无索引。
- new：仅 `created_at`/`updated_at` → `TEXT NOT NULL`，其余列与 CHECK 逐字保留。

### 3.4 `refresh_tokens`（A 组）

- legacy `:35-42`；new：`expires_at`/`created_at` → `TEXT NOT NULL`，`revoked_at` → `TEXT`（NULL 包裹）；索引 `idx_refresh_tokens_user_id`。

### 3.5 `email_verification_challenges` / `password_recovery_challenges`（A 组；两表同构）

- legacy `:171-177` / `:231-237`；new：`expires_at`/`sent_at` → `TEXT NOT NULL`；均无索引。

### 3.6 `login_failures`（A 组；`locked_until` = **D0**）

- legacy `:200-207`；new：`locked_until` → `TEXT`（去 `NOT NULL`/`DEFAULT 0`，copy 用 `= 0 → NULL`），`updated_at` → `TEXT NOT NULL`；无索引。

### 3.7 `user_password_history`（A 组）

- legacy `:268-273`；new：`created_at` → `TEXT NOT NULL`；索引 `idx_user_password_history_user`（`:274`）。

### 3.8 `user_invites`（A 组）

- legacy `:294-305`；new：`expires_at`/`last_sent_at`/`created_at` → `TEXT NOT NULL`，`consumed_at`/`revoked_at` → `TEXT`（NULL 包裹）；索引 `idx_user_invites_created`（`:306`）。

### 3.9 `service_credentials`（NN ×3 + N ×2；无 FK 子表 → 裸四步）

- legacy `:462-474`；索引 `:475-476`。new：`expires_at`/`created_at`/`updated_at` → `TEXT NOT NULL`，`revoked_at`/`last_used_at` → `TEXT`。

### 3.10 B 组三张联接表的建回正文（阶段 4；无时间列）

```sql
CREATE TABLE user_roles (
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  role_id TEXT NOT NULL REFERENCES roles(id) ON DELETE RESTRICT,
  PRIMARY KEY (user_id, role_id)
);
INSERT INTO user_roles (user_id, role_id) SELECT user_id, role_id FROM user_roles_bak;
DROP TABLE user_roles_bak;

CREATE TABLE role_permissions (
  role_id       TEXT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
  permission_id TEXT NOT NULL REFERENCES permissions(id) ON DELETE RESTRICT,
  PRIMARY KEY (role_id, permission_id)
);
INSERT INTO role_permissions (role_id, permission_id) SELECT role_id, permission_id FROM role_permissions_bak;
DROP TABLE role_permissions_bak;

CREATE TABLE role_menu_items (
  role_id      TEXT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
  menu_item_id TEXT NOT NULL REFERENCES menu_items(id) ON DELETE RESTRICT,
  PRIMARY KEY (role_id, menu_item_id)
);
INSERT INTO role_menu_items (role_id, menu_item_id) SELECT role_id, menu_item_id FROM role_menu_items_bak;
DROP TABLE role_menu_items_bak;

CREATE INDEX idx_user_roles_role_id ON user_roles(role_id);
CREATE INDEX idx_role_permissions_permission_id ON role_permissions(permission_id);
CREATE INDEX idx_role_menu_items_menu_item_id ON role_menu_items(menu_item_id);
```

### 3.11 C 组三张跨 descriptor 子表的建回正文（阶段 5；**时间列仍 INTEGER**）

```sql
-- notifications（live 形状，含 v37 追加的 title_key / body_key）
CREATE TABLE notifications (
  id         TEXT PRIMARY KEY,
  user_id    TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  event      TEXT NOT NULL CHECK (event IN ('account.locked','account.disabled','account.unlocked','account.password-changed')),
  title      TEXT NOT NULL,
  body       TEXT NOT NULL,
  read_at    INTEGER,
  created_at INTEGER NOT NULL,
  title_key  TEXT,
  body_key   TEXT
);
INSERT INTO notifications (id, user_id, event, title, body, read_at, created_at, title_key, body_key)
SELECT id, user_id, event, title, body, read_at, created_at, title_key, body_key FROM notifications_bak;
DROP TABLE notifications_bak;
CREATE INDEX idx_notifications_user_created ON notifications(user_id, created_at DESC);

-- user_mfa（时间列仍 INTEGER）
CREATE TABLE user_mfa (
  user_id                TEXT PRIMARY KEY REFERENCES users(id),
  status                 TEXT NOT NULL CHECK (status IN ('pending','active')),
  totp_secret_ciphertext TEXT NOT NULL,
  recovery_codes_hash    TEXT NOT NULL,
  last_used_step         INTEGER NOT NULL DEFAULT 0,
  created_at             INTEGER NOT NULL,
  updated_at             INTEGER NOT NULL
);
INSERT INTO user_mfa (user_id, status, totp_secret_ciphertext, recovery_codes_hash, last_used_step, created_at, updated_at)
SELECT user_id, status, totp_secret_ciphertext, recovery_codes_hash, last_used_step, created_at, updated_at FROM user_mfa_bak;
DROP TABLE user_mfa_bak;

-- mfa_proofs（时间列仍 INTEGER）
CREATE TABLE mfa_proofs (
  id         TEXT PRIMARY KEY,
  user_id    TEXT NOT NULL REFERENCES users(id),
  fail_count INTEGER NOT NULL DEFAULT 0,
  expires_at INTEGER NOT NULL,
  created_at INTEGER NOT NULL
);
INSERT INTO mfa_proofs (id, user_id, fail_count, expires_at, created_at)
SELECT id, user_id, fail_count, expires_at, created_at FROM mfa_proofs_bak;
DROP TABLE mfa_proofs_bak;
```

> v80/v81 之后再按 §2.7/§2.8 对这三张表做时间列类型转换（裸四步）。

## 4. `operationlog`（v75）逐表 DDL 正文（毫秒族）

### 4.1 `operation_log`（FK 父表 → **必须走 WithSessions 舞步**）

- 当前有效 DDL = `operationLogDigitalOfferDDL[0]`（`operationlog/migration.go:513-521`），索引 `idx_operation_log_created_at`（`:522`）。
- `event` 列 CHECK 为**超长枚举**，**逐字保留**（完整清单以源码为准，本文件不重抄以降低转录风险）。
- 子表：`operation_log_correlation`（`:782-785` + 索引 `:786`）、`operation_log_session`（`:554-557` + 索引 `:558`）。`operation_log_archive_session`（`:559-562`）**无 FK**，不参与舞步。

```text
 1. CREATE TEMP TABLE operation_log_correlation_backup AS SELECT operation_id, correlation_id FROM operation_log_correlation;
 2. CREATE TEMP TABLE operation_log_session_backup       AS SELECT operation_id, session_id     FROM operation_log_session;
 3. DROP TABLE operation_log_session;
 4. DROP TABLE operation_log_correlation;
 5. ALTER TABLE operation_log RENAME TO operation_log_old;
 6. CREATE TABLE operation_log ( ...new：created_at → TEXT NOT NULL；event CHECK 逐字... );
 7. INSERT INTO operation_log (id, event, actor_id, actor_name, record_id, detail, created_at)
    SELECT id, event, actor_id, actor_name, record_id, detail, <毫秒表达式(created_at)> FROM operation_log_old;
 8. DROP TABLE operation_log_old;
 9. CREATE INDEX idx_operation_log_created_at ON operation_log(created_at DESC);
10. CREATE TABLE operation_log_correlation (operation_id TEXT PRIMARY KEY REFERENCES operation_log(id) ON DELETE CASCADE, correlation_id TEXT NOT NULL);
    CREATE INDEX idx_operation_log_correlation_id ON operation_log_correlation(correlation_id);
11. CREATE TABLE operation_log_session (operation_id TEXT PRIMARY KEY REFERENCES operation_log(id) ON DELETE CASCADE, session_id TEXT NOT NULL);
    CREATE INDEX idx_operation_log_session_id ON operation_log_session(session_id);
12. INSERT INTO operation_log_correlation (operation_id, correlation_id) SELECT operation_id, correlation_id FROM operation_log_correlation_backup;
13. INSERT INTO operation_log_session (operation_id, session_id) SELECT operation_id, session_id FROM operation_log_session_backup;
14. DROP TABLE operation_log_correlation_backup; DROP TABLE operation_log_session_backup;
```

> 与既有 helper 的关系：步骤 1–14 与 `rebuildOperationLogWithSessions`（`:709-754`）同序；差别仅在步骤 6/7 的类型与表达式。既有 `rebuildOperationLog` 的 copy 是**硬编码直通**（`:763-765`），**不可复用**——须参数化（见 §7）。

### 4.2 `operation_log_archive`（无 FK 子表）

- legacy `:527-536`；new：`created_at`/`archived_at` → `TEXT NOT NULL`；索引 `idx_operation_log_archive_created_at`（`:537`）。
- 旁表 `operation_log_archive_correlation`（`:538-541`）与 `operation_log_archive_session`（`:559-562`）**无 FK** 指向它，不参与重建。

## 5. `wallet` / `telegram` / `digital-offer` 逐表 DDL 正文

### 5.1 `wallet_accounts`（NN；金额列保持 INTEGER）

- legacy 有效形状 `wallet/migration.go:304-318`（v64 rebuild 后）；new 仅 `created_at`/`updated_at` → `TEXT NOT NULL`；表级 `UNIQUE` 与 `CHECK (balance_total = balance_available + balance_frozen)` 逐字保留；无索引。

### 5.2 `wallet_ledger_entries`（NN；14 列）

- legacy 有效形状 `:145-162`（v33 rebuild 后）；new 仅 `created_at` → `TEXT NOT NULL`；索引 `idx_wallet_ledger_account`（`:163`）。
- 金额列 `amount_delta` 及三个 `balance_after_*` **保持 `INTEGER`**。

### 5.3 `wallet_reconciliation_runs` / `subjects` / `voucher_batches`

- legacy `:56-64` / `:279-285` / `:391-395`；均裸四步，时间列 → `TEXT NOT NULL`。
- 索引：`subjects` 用 `idx_subjects_issuer_external`（`:286`）；另两表无索引。

### 5.4 `vouchers`（`expires_at`/`redeemed_at` = N + legacy 0 → NULL）

- legacy `:287-301`；索引 `idx_vouchers_batch`、`idx_vouchers_status`（`:302-303`）。
- new：`expires_at`/`redeemed_at` → `TEXT`（**voucher 包裹** `IS NULL OR = 0`），`created_at`/`updated_at` → `TEXT NOT NULL`。

### 5.5 `telegram_config`（v86；**D0**）

```sql
-- legacy 有效形状（telegram/migration.go:11-16 + v67 ALTERs :29-30）
CREATE TABLE telegram_config (
  id                      INTEGER PRIMARY KEY CHECK (id = 1),
  bot_token_enc           TEXT    NOT NULL DEFAULT '',
  webhook_secret_enc      TEXT    NOT NULL DEFAULT '',
  updated_at              INTEGER NOT NULL DEFAULT 0,
  mode                    TEXT    NOT NULL DEFAULT 'polling',
  webhook_public_base_url TEXT    NOT NULL DEFAULT ''
)
```
```sql
-- new：updated_at → TEXT（可空、无默认）；copy 用 CASE WHEN updated_at = 0 THEN NULL ELSE <秒表达式> END；无索引
```

### 5.6 `telegram_sessions` / `telegram_inbound_messages` / `telegram_outbound_messages`（v86）

- legacy 有效形状 `:34-44` / `:47-61` / `:100-112`；均裸四步，时间列 → `TEXT NOT NULL`。
- 索引（最后建）：`idx_telegram_sessions_activity`、`idx_telegram_inbound_messages_chat_received`、`idx_telegram_outbound_messages_chat_created`、**部分唯一** `idx_telegram_outbound_messages_pending_root … WHERE status = 'pending'`。

### 5.7 `digital_offers` / `digital_purchases` / `digital_entitlements`

- legacy `digitaloffer/migration.go:21-39` / `:41-54` / `:57-73`；均裸四步，时间列 → `TEXT`。
- 两组 form CHECK **逐字保留**；`digital_entitlements.expires_at` → `TEXT` 可空（普通 NULL 包裹，非 voucher 规则）。
- 索引：`idx_digital_offers_status`、`idx_digital_purchases_subject`、`idx_digital_purchases_offer`、`idx_digital_entitlements_subject`、`idx_digital_entitlements_purchase`。

### 5.8 `wallet`（v85）子表盘点前置（A-032 §G）

- v85 落码前**必须**先跑子表盘点：确认无任何存留 DDL 含 `REFERENCES wallet_accounts` / `wallet_ledger_entries` / `wallet_reconciliation_runs` / `subjects` / `vouchers` / `voucher_batches`。
- 若发现子表 → 按 §1 的 F-5 处理并把子表写入本冻结包；若为空 → 裸四步即可。

## 6. v73 的 `applied_at` 写入路径（F-I-022 关闭要求）

v73 设计**必须**同时改「DDL 字面」与「四处写入」，否则 restore 库与 migrate 库形状分叉，且 v73 行本身把整数写进 TEXT 列（**静默破坏合同**）：

| # | 位置 | 现状 | 目标 |
|--:|------|------|------|
| 1 | `internal/store/identity.go:58` `sqliteLedgerDDL` | `applied_at INTEGER NOT NULL` | `TEXT NOT NULL` |
| 2 | `internal/store/identity.go:65` `postgresLedgerDDL` | `applied_at BIGINT NOT NULL` | `timestamptz(6) NOT NULL` |
| 3 | `internal/store/identity.go:317,319-321` `stampCatalog` | `now := time.Now().UTC().Unix()` | `now := time.Now().UTC().Truncate(time.Microsecond)`；绑定 `time.Time`（**保留**「整批共用一个 now」语义） |
| 4 | `internal/store/migrate.go:121-123` `applyMigration` | `time.Now().UTC().Unix()` | 同上 |
| 5 | `internal/store/postgres.go:165-167` `applyMigrationPG` | `time.Now().UTC().Unix()` | 同上 |

- **禁止改** `authsession/migration.go:18-23` 的 `schemaMigrationsDDL`（属 **v1**）。
- **不影响 v1 checksum**：`identity.go` 的两个 restore 字面不在 `0001:r2-baseline` 的 `r2BaselineDDL` 哈希输入内。

## 7. 既有 `rebuildOperationLog` 的最小安全改动（设计在 R1、代码在 R2）

- **改动 1**：`rebuildOperationLog`（`operationlog/migration.go:756-776`）在 `ALTER TABLE operation_log RENAME TO operation_log_old` **之前**加 **fail-closed 断言**：`operation_log_correlation` 与 `operation_log_session` **均不存在**（SQLite 查 `sqlite_master`，PG 查 `pg_class`）；存在则返回错误而不是继续。
- **改动 2**：**不**把子表舞步内嵌进 `rebuildOperationLog`（会与 `WithSessions` 双重 drop）；`WithCorrelation`/`WithSessions` 先 drop 再调用，断言自然通过。
- **改动 3**：v75 **必须**调用 `rebuildOperationLogWithSessions`（或与 F-5 同一通用 helper）；**禁止** `pgRebuild`。
- **append-only 边界**：改 Go 控制流/断言**不进入** `MigrationChecksum`（哈希的是 `stmts` 切片，不是 `Apply` 函数体）；**禁止**改 `0004–0071` 的 DDL 切片字面与 `transform_id` → **0001–0072 canonical SQL/checksum 不变性不受影响**。

## 8. 整链实测验证记录（本机 SQLite 3.51.2）

对 §1.2 的 v74 F-5 编排做了端到端探针（`PRAGMA foreign_keys=ON`），覆盖 4 张父表（`users`/`roles`/`permissions`/`menu_items`）、A 组 1 张（`refresh_tokens`）、B 组 1 张（`user_roles`，双父表 CASCADE+RESTRICT）、C 组 2 张（`notifications`、`user_mfa`，时间列保持 INTEGER）。

| 阶段 | 检查项 | 结果 |
|------|--------|------|
| v74 后 | `PRAGMA foreign_key_check` | **0 行** |
| v74 后 | 行数（users / roles / user_roles / notifications / user_mfa） | **2 / 1 / 1 / 1 / 1**（无丢行） |
| v74 后 | `users.locked_until`（D0） | `0` → **NULL**；非 0 → `2025-09-19T22:13:22.000000Z` |
| v74 后 | C 组 `notifications.created_at` 列类型 | 仍为 **INTEGER**（按 `D-019` §3 设计） |
| v74 后 | B/C 组子表 `REFERENCES` 文本 | 指向 `users(id)`，**不含** `_old` → **OK** |
| v81 后 | `notifications` 二次重建（裸四步） | `fk_check = 0`；`created_at` = 27 字符 TEXT；FK 文本仍 **OK**；`integrity_check = ok` |

**结论**：`D-019` 的 F-5 编排（含「v74 先修 FK、v80/v81 再转类型」两次重建）在本机可执行，且**不丢行、不悬挂 FK**。

## 9. 声明

- 本文件 `status: freeze-candidate`；**不**声称任何 DDL/迁移/测试已实施；`apps/` 本轮未修改。
- §2/§3 给出可直接落码的完整正文；§4/§5 对**超长既有 DDL**（`operation_log.event` 枚举、已由 v33/v64 重建过的 wallet 表）以**源码行号 + 差异说明**代替重抄，以降低转录风险；R2 落码时以源码为准。
- `<秒表达式(col)>` / `<毫秒表达式(col)>` 为 §0 两条表达式的占位写法，展开以 §0 为准。
- 是否构成 **F-I-002.1 / F-I-021 / F-I-022** 的合法闭合，由 independent 复审判定；本编排器不自证。
- 引用 `D-0NN` 一律限定 Root/child。
