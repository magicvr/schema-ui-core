---
id: r1-time-column-inventory-v0.1
doc_type: evidence-attachment
title: R1 时间列与单位 inventory v0.1
status: collecting
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-002-r1-contract-and-denominator-freeze
version: 0.1.0
---

# R1 时间列与单位 inventory v0.1

> 本附件是 R1 C1 的第一版只读盘点，不是最终分母或迁移许可。编译迁移目录的历史基线为 VP-013 开区时的 48 migrations；VP-040 目标物理合同已由用户选择为 PostgreSQL `timestamptz(6)` + SQLite fixed-6 UTC RFC3339 `TEXT`。旧物理形状与单位只作为转换来源。

## 1. 全局合同与扫描结论

| 项 | 当前观察 | R1 处理 |
|----|----------|---------|
| SQLite 现状 | 绝大多数时间列 DDL 为 `INTEGER` | 转为 canonical UTC RFC3339 fixed-6 `TEXT`；nullable/默认值逐列冻结 |
| PostgreSQL 现状 | VP-013 现有 DDL 为 `BIGINT` | 转为 `timestamptz(6)`；秒/毫秒按旧列单位解释后转换 |
| 单位 | 主要为 Unix seconds；`core.jobs`、`core.operationlog`、`core.persistence.mail_*` 为 Unix milliseconds | 不做单位猜测；每列转换函数与回读测试必须点名 |
| 现存 TEXT 时间列 | 本轮 compiled DDL 扫描未发现时间列为 TEXT | 仍需全 catalog 复扫，确认运行时/历史表无例外 |
| sentinel | `locked_until`、`last_login_failure_at`、部分单例配置 `updated_at DEFAULT 0` 使用 0 | 用户裁决：语义为“未发生/无期限”的 0 → `NULL`；逐列记录 |
| 非时间整数 | ID/ID 前缀、duration、TOTP step、version、计数器、金额、flag 等 | 排除；不得因列名或整数类型误纳入 |

## 2. 当前已确认的 compiled catalog 分组

| 模块 / 表 | 绝对时刻列 | SQLite 旧形状 | PG 旧形状 | 当前单位 / 证据 | R1 目标 |
|-----------|------------|---------------|------------|----------------|---------|
| `core.auth-session.users` | `created_at`, `updated_at`, `locked_until`, `last_login_failure_at` | INTEGER | BIGINT | seconds；`apps/api/modules/authsession/accounts.go`, `accounts_lock_source.go`；0 sentinel 于 lock/failure | `timestamptz(6)` / fixed-6 TEXT；0 → NULL（逐列确认） |
| `core.auth-session.refresh_tokens` | `expires_at`, `revoked_at`, `created_at` | INTEGER | BIGINT | seconds；`account_operations.go`, `accounts.go` | 同上 |
| `core.auth-session.roles/permissions/menu_items` | `created_at`, `updated_at` | INTEGER | BIGINT | seconds；`roles_repository.go`, `systemdata` | 同上 |
| `core.auth-session.schema_migrations` | `applied_at` | INTEGER | BIGINT | seconds；`authsession/migration/migration.go` / `store/migrate.go` | 同上 |
| `core.auth-session.system_data_reconcile` | `applied_at` | INTEGER | BIGINT | seconds；migration DDL | 同上 |
| `core.auth-session.service_credentials` | `expires_at`, `revoked_at`, `last_used_at`, `created_at`, `updated_at` | INTEGER | BIGINT | seconds；`service_credentials.go` | nullable `revoked_at/last_used_at` → NULL |
| `core.auth-session.email_verification_challenges` | `expires_at`, `sent_at` | INTEGER | BIGINT | seconds；migration comments + `email_identity.go` | same |
| `core.auth-session.password_recovery_challenges` | `expires_at`, `sent_at` | INTEGER | BIGINT | seconds；migration comments + `recovery.go` | same |
| `core.auth-session.user_password_history` | `created_at` | INTEGER | BIGINT | seconds；migration DDL | same |
| `core.auth-session.user_invites` | `expires_at`, `consumed_at`, `revoked_at`, `last_sent_at`, `created_at` | INTEGER | BIGINT | seconds；`invites.go` | nullable fields → NULL |
| `core.auth-session.login_failures` | `locked_until`, `updated_at` | INTEGER | BIGINT | seconds；`accounts_lock_source.go`; `locked_until=0` sentinel | 0 → NULL; verify domain handling |
| `core.persistence.mail_outbox` | `created_at` | INTEGER | BIGINT | **milliseconds**; `internal/mail/outbox.go:99-102, 233-237` | fixed-6 TEXT / `timestamptz(6)` |
| `core.persistence.mail_config` | `updated_at` | INTEGER DEFAULT 0 | BIGINT DEFAULT 0 | **milliseconds** at runtime; `internal/mail/runtime.go:382-390`, read `UnixMilli` | 0 sentinel → NULL; singleton schema/default needs explicit R1 decision |
| `core.persistence.records` | `updated_at` | INTEGER | BIGINT | historical v3/v6; table is retired by later catalog; no current repository | mark historical/retired; verify current DB migration path before including conversion |
| `core.jobs.jobs` | `lease_expires_at`, `created_at`, `updated_at`, `finished_at`, `expires_at` | INTEGER | BIGINT | **milliseconds**; `internal/jobs/model.go:144-146`, repository/actions/list | nullable lease/finished/expires → NULL |
| `core.operationlog.operation_log` | `created_at` | INTEGER | BIGINT | **milliseconds**; `operationlog/repository.go:173-177, 312-337` | same |
| `core.operationlog.operation_log_archive` | `created_at`, `archived_at` | INTEGER | BIGINT | **milliseconds**; `operationlog/retention.go:21-32` | same |
| `admin.data-dictionary.dict_types/dict_entries` | `created_at`, `updated_at` | INTEGER | BIGINT | seconds; `datadictionary/store/repository.go` | same |
| `admin.data-permission.data_scope_policies/user_data_scopes` | `updated_at` | INTEGER | BIGINT | seconds; `datapermission/store/repository.go` | same |
| `admin.login-captcha.captcha_challenges` | `expires_at`, `created_at` | INTEGER | BIGINT | seconds; `logincaptcha/store/repository.go` | same |
| `admin.login-captcha.captcha_config` | `created_at`, `updated_at` | INTEGER | BIGINT | seconds; `logincaptcha/store/repository.go` | same |
| `admin.mfa.user_mfa` | `created_at`, `updated_at` | INTEGER | BIGINT | seconds; `mfa/store/repository.go`; `last_used_step` excluded | same |
| `admin.mfa.mfa_proofs` | `expires_at`, `created_at` | INTEGER | BIGINT | seconds; `mfa/store/repository.go` | same |
| `admin.notifications.notifications` | `read_at`, `created_at` | INTEGER nullable / NOT NULL | BIGINT nullable / NOT NULL | seconds; `notifications/migration/migration.go`, `notifications_repository.go` | `read_at NULL` preserved |
| `admin.recycle-bin.recycle_items` | `deleted_at`, `restored_at` | INTEGER / nullable | BIGINT / nullable | seconds; `recyclebin/store/repository.go` | `restored_at NULL` preserved |
| `admin.scheduled-tasks.scheduled_tasks` | `created_at`, `updated_at` | INTEGER | BIGINT | seconds; `scheduledtasks/store/repository.go` | same |
| `admin.scheduled-tasks.task_runs` | `started_at`, `finished_at`, `created_at` | INTEGER / nullable | BIGINT / nullable | seconds; `scheduledtasks/store/repository.go` | `finished_at NULL` preserved; verify nil scan/write |
| `admin.settings.site_settings` | `updated_at` | INTEGER | BIGINT | seconds; `settings/migration/migration.go:52-73` | same |
| `admin.wallet.wallet_accounts` | `created_at`, `updated_at` | INTEGER | BIGINT | seconds; `wallet/store/repository.go` | same |
| `admin.wallet.wallet_ledger_entries` | `created_at` | INTEGER | BIGINT | seconds; `wallet/store/repository.go` | same |
| `admin.wallet.wallet_reconciliation_runs` | `created_at` | INTEGER | BIGINT | seconds; `wallet/store/repository.go` | same |
| `admin.wallet.subjects` | `created_at` | INTEGER | BIGINT | seconds; `wallet/subject/subject.go` | same |
| `admin.wallet.vouchers` | `expires_at`, `redeemed_at`, `created_at`, `updated_at` | INTEGER / nullable | BIGINT / nullable | seconds; `wallet/voucher/service.go` | nullable expiration/redemption → NULL |
| `admin.wallet.voucher_batches` | `created_at`, `updated_at` | INTEGER | BIGINT | seconds; migration batch registry | same |
| `admin.channel.telegram.telegram_config` | `updated_at` | INTEGER DEFAULT 0 | BIGINT DEFAULT 0 | seconds; `telegram/migration/migration.go:10-25`; 0 sentinel | 0 → NULL; verify singleton read/write |
| `admin.channel.telegram.telegram_sessions` | `last_message_at`, `created_at`, `updated_at` | INTEGER | BIGINT | seconds; `telegram/store/repository.go` | same |
| `admin.channel.telegram.telegram_inbound_messages` | `received_at` | INTEGER | BIGINT | seconds; `telegram/store/repository.go` | same |
| `admin.channel.telegram.telegram_outbound_messages` | `created_at`, `updated_at` | INTEGER | BIGINT | seconds; `telegram/store/repository.go` | same |
| `biz.digital-offer.digital_offers` | `created_at`, `updated_at` | INTEGER | BIGINT | seconds; `digitaloffer/store/store.go` | same |
| `biz.digital-offer.digital_purchases` | `created_at` | INTEGER | BIGINT | seconds; `digitaloffer/store/store.go` | same |
| `biz.digital-offer.digital_entitlements` | `expires_at`, `created_at`, `updated_at` | INTEGER / nullable | BIGINT / nullable | seconds; `digitaloffer/store/store.go` | nullable duration expiry → NULL |

## 3. Explicit exclusions requiring R1 evidence

| Category | Examples | Why excluded |
|----------|----------|--------------|
| Identifier timestamps | `jobs`/wallet/digital-offer ID hex prefixes from `UnixMilli()` | ID generation, not persisted absolute-time columns |
| Durations | `duration_seconds`, TTL/config durations | Relative durations, not instants |
| Counters/steps | `last_used_step`, `attempt`, `version`, `progress`, `fail_count`, `sort` | Numeric state, not time |
| Money/quantities | wallet `balance_*`, `amount_delta`, offer `price_amount`, voucher `amount` | Financial/quantity semantics; not time |
| Flags | `enabled`, `cancel_requested`, `must_change_password` | Boolean state |
| Retired records | `core.persistence.records.updated_at` | Historical DDL is later retired; must verify no current table before any conversion |

## 4. R1 unresolved evidence

1. Complete compiled-catalog table/column enumeration must be mechanically checked against the 48-migration catalog and current schema introspection; this document is v0.1, not a closure claim.
2. Every runtime scan/write path must gain a target codec mapping: old seconds → `timestamptz`/fixed-6 TEXT, old milliseconds → same, nullable and sentinel branches.
3. SQLite in-place conversion needs table-rebuild strategy for NOT NULL/default/sentinel changes; PostgreSQL needs `USING` conversion and nullable/default ordering. Failure/rollback must be tested before R2.
4. `I-040-004` (VP-020 display/input regression matrix) remains open for R3.

## 5. Evidence index

- Historical contract: `docs/workspaces/workspace-013-store-dialects/GOAL-002-r1-tx-port-and-config/attachments/r1-tx-port-and-config-freeze.md` §3/§4.
- Historical full dual-dialect catalog: `docs/workspaces/workspace-013-store-dialects/GOAL-004-r3-dual-dialect-ledger/00-meta.md` §R3 / I-001/I-002.
- Runtime millisecond evidence: `apps/api/internal/jobs/model.go:144-146`; `apps/api/internal/mail/outbox.go:99-102,233-237`; `apps/api/internal/mail/runtime.go:308,382-390`; `apps/api/modules/operationlog/repository.go:173-177,312-337`; `apps/api/modules/operationlog/retention.go:21-32`.
- Runtime second evidence: representative authsession, wallet, dictionary, MFA, notification, recycle-bin, scheduled-task, Telegram and digital-offer repositories listed above.

## Status

`I-040-001`～`I-040-003` remain `collecting` pending C1/C2/C3 evidence. R1 is not yet ready for self or independent close-out.
