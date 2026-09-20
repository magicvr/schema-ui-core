---
id: r1-time-column-inventory-v0.2
doc_type: evidence-attachment
title: R1 时间列与单位 inventory v0.2 · 完整只读盘点
status: recorded
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-002-r1-contract-and-denominator-freeze
version: 0.2.0
---

# R1 时间列与单位 inventory v0.2

## 结论摘要

- 当前 live schema 共 **90 个绝对时刻列**；compiled historical DDL 另有 `core.persistence.records.updated_at`（v3 建、v6 retire/drop），当前无 repository/handler，不计 live 分母，但保留为历史来源。
- 全仓 compiled migration DDL 未发现现存时间列 TEXT；当前形状全部为 SQLite `INTEGER` / PostgreSQL `BIGINT`。
- 运行时单位确实混用：
  - **milliseconds**：`core.jobs` 全部时间列；`core.operationlog` `created_at`/archive；`core.persistence.mail_outbox.created_at` 与 `mail_config.updated_at`。
  - **seconds**：authsession、dictionary、permission、captcha、MFA、notifications、recycle-bin、scheduled-tasks、wallet、Telegram、digital-offer 等主要面。
- 当前唯一明确的 0 sentinel：auth lock/failure 列、mail/telegram singleton config 的 `updated_at DEFAULT 0`；`task_runs.finished_at` 另有 runtime 写 0、读时 `>0` 才恢复 nil 的历史 sentinel。
- 排除：ID 前缀、`duration_seconds`、TOTP `last_used_step`、version、计数器、金额与 flag。

## 完整分组清单（当前旧形状 → VP-040 目标）

> 所有下列旧形状均为 SQLite `INTEGER` / PostgreSQL `BIGINT`，除非标注历史/特殊；目标统一为 PostgreSQL `timestamptz(6)` + SQLite fixed-6 UTC RFC3339 `TEXT`。旧单位必须按列转换，不能从列名猜测。

| 模块 / 表 | 绝对时刻列 | 旧单位 | nullable / sentinel | 主要证据 |
|-----------|------------|---------|---------------------|----------|
| `core.auth-session` schema ledger | `schema_migrations.applied_at`, `system_data_reconcile.applied_at` | sec | NN | `authsession/migration/migration.go:18-23,91-100,332-337,410-420`; `internal/store/migrate.go:121-124` |
| `users` | `created_at`, `updated_at`, `locked_until`, `last_login_failure_at` | sec | lock/failure D0=0 | `authsession/migration/migration.go:26-34,121-127,199-220`; `accounts.go:191-298`; `accounts_lock_source.go:65-128` |
| auth sessions | `refresh_tokens.expires_at`, `revoked_at`, `created_at` | sec | revoked N | `authsession/migration/migration.go:35-42,353-360`; `account_operations.go:27-65` |
| RBAC | `roles/permissions/menu_items.created_at`, `updated_at` | sec | NN | `authsession/migration/migration.go:47-82,366-401`; `roles_repository.go:93-137,404-417` |
| auth challenges | `email_verification_challenges.expires_at/sent_at`; `password_recovery_challenges.expires_at/sent_at` | sec | NN | `authsession/migration/migration.go:170-247`; `email_identity.go:174-204`; `recovery.go:156-240` |
| auth history/invites | `user_password_history.created_at`; `user_invites.expires_at/consumed_at/revoked_at/last_sent_at/created_at` | sec | consumed/revoked N | `authsession/migration/migration.go:268-321`; `password_policy.go:169-185`; `invites.go:73-100,297-342` |
| service credentials | `expires_at/revoked_at/last_used_at/created_at/updated_at` | sec | revoked/last_used N | `authsession/migration/migration.go:438-476`; `service_credentials.go:145-249` |
| `core.persistence.mail_outbox` | `created_at` | **ms** | NN | `corepersistence/migration/migration.go:48-74`; `internal/mail/outbox.go:95-102,222-240` |
| `core.persistence.mail_config` | `updated_at` | **ms** | D0=0 | `corepersistence/migration/migration.go:87-124`; `internal/mail/runtime.go:169-176,308,382-390` |
| historical `records` | `records.updated_at` | unknown/legacy | retired v6 | `corepersistence/migration/migration.go:12-46`; no current repository; exclude live after current-schema verification |
| `core.jobs.jobs` | `lease_expires_at/created_at/updated_at/finished_at/expires_at` | **ms** | lease/finished/expires N; status CHECK | `jobs/migration/migration.go:14-87`; `internal/jobs/model.go:136-146`; repository/actions/list |
| `core.operationlog.operation_log` | `created_at` | **ms** | NN | `operationlog/migration/migration.go:13-24,245-280`; `operationlog/repository.go:172-177,286-337` |
| operation archive | `operation_log_archive.created_at/archived_at` | **ms** | NN | `operationlog/migration/migration.go:525-537`; `operationlog/retention.go:21-42,66-73` |
| `admin.data-dictionary` | `dict_types/dict_entries.created_at/updated_at` | sec | NN | `datadictionary/migration/migration.go:16-71`; `datadictionary/store/repository.go:109-187,292-393` |
| `admin.data-permission` | `data_scope_policies.updated_at`, `user_data_scopes.updated_at` | sec | NN | `datapermission/migration/migration.go:16-52`; repository `68-211` |
| `admin.login-captcha` | `captcha_challenges.expires_at/created_at`; `captcha_config.created_at/updated_at` | sec | NN | `logincaptcha/migration/migration.go:15-47`; repository `36-122` |
| `admin.mfa` | `user_mfa.created_at/updated_at`; `mfa_proofs.expires_at/created_at` | sec | NN | `mfa/migration/migration.go:16-57`; `mfa/store/repository.go:59-307` |
| `admin.notifications` | `notifications.read_at/created_at` | sec | `read_at` N = unread | `notifications/migration/migration.go:13-40`; `notifications_repository.go:85-95,186-275` |
| `admin.recycle-bin` | `deleted_at/restored_at` | sec | `restored_at` N = restorable | `recyclebin/migration/migration.go:16-49`; repository `81-216` |
| `admin.scheduled-tasks` | `scheduled_tasks.created_at/updated_at`; `task_runs.started_at/finished_at/created_at` | sec | `finished_at` N but historical runtime 0 sentinel | `scheduledtasks/migration/migration.go:15-64`; repository `114-193,259-367` |
| `admin.settings` | `site_settings.updated_at` | sec | NN | `settings/migration/migration.go:22-42`; repository `204-244,331-363` |
| `admin.wallet` accounts/ledger | `wallet_accounts.created_at/updated_at`; `wallet_ledger_entries.created_at`; `wallet_reconciliation_runs.created_at` | sec | NN | `wallet/migration/migration.go:16-114`; `wallet/store/repository.go:250-718,886-1003` |
| wallet subjects/vouchers | `subjects.created_at`; `vouchers.expires_at/redeemed_at/created_at/updated_at`; `voucher_batches.created_at/updated_at` | sec | expires/redeemed N | `wallet/migration/migration.go:275-347,388-407`; `wallet/subject/subject.go:86-191`; `wallet/voucher/service.go:79-419` |
| `admin.channel.telegram` config/sessions | `telegram_config.updated_at`; `telegram_sessions.last_message_at/created_at/updated_at` | sec | config D0=0; sessions NN | `channel/telegram/migration/migration.go:10-46,66-79`; repository `161-284` |
| Telegram messages | `telegram_inbound_messages.received_at`; `telegram_outbound_messages.created_at/updated_at` | sec | NN | `channel/telegram/migration/migration.go:47-137`; repository `127-179,550-724` |
| `biz.digital-offer` | `digital_offers.created_at/updated_at`; `digital_purchases.created_at`; `digital_entitlements.expires_at/created_at/updated_at` | sec | entitlement `expires_at` N | `digitaloffer/migration/migration.go:16-138`; `digitaloffer/store/store.go:154-676` |

## Non-time exclusions

- ID prefixes generated with `UnixMilli`/`UnixNano` are identifiers, not persisted absolute-time columns.
- `duration_seconds`, `last_used_step`, `version`, `attempt`, `progress`, `fail_count`, `sort`, `count_per_purchase`, `remaining_count` are not instants.
- `balance_*`, `amount_delta`, `price_amount`, voucher `amount`, `enabled`/flags are not instants.
- `records.updated_at` is a historical retired table and requires current-schema proof before any conversion scope.

## C1 closure status

The read-only scan now covers the 90 live absolute-time columns plus one historical retired source. C1 inventory is ready for audit review. C2/C3 remain open: target codec/DDL, table rebuild/`USING` conversions, NULL/default backfill, failure rollback and backup verification are not yet implemented or frozen.
