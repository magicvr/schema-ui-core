---
id: r1-time-column-inventory-v0.3
doc_type: evidence-attachment
title: R1 时间列与单位 inventory v0.3 · 逐列可加总清单
status: recorded
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-002-r1-contract-and-denominator-freeze
version: 0.3.0
---

# R1 时间列与单位 inventory v0.3

## 计数与边界

- **live denominator = 90 columns**，以下逐列列出，可机械计数。
- `core.persistence.records.updated_at` 是 v3 历史表、v6 `records_retire` drop 的 retired source，**不计 live denominator**；转换前仍需 current-schema 无表证明。
- 当前 compiled DDL 的时间物理形状：SQLite `INTEGER` / PostgreSQL `BIGINT`；无现存时间列 TEXT。
- 目标：PostgreSQL `timestamptz(6)`；SQLite fixed-6 UTC RFC3339 `TEXT`。
- `sec` / `ms` 是现有 runtime 绑定单位；`I/B` = current SQLite INTEGER / PG BIGINT；`NN` = non-null、无 default；`N` = nullable；`D0` = non-null default 0。

## 逐列清单（90）

| # | 模块 / 表 / 列 | 旧单位 | 旧 null/default | runtime / DDL evidence |
|---:|---|---|---|---|
| 1 | `core.auth-session.schema_migrations.applied_at` | sec | NN | `authsession/migration/migration.go:18-23,332-337`; `store/migrate.go:121-124` |
| 2 | `core.auth-session.system_data_reconcile.applied_at` | sec | NN | `authsession/migration/migration.go:91-100,410-420`; `systemdata/reconcile.go:46-64` |
| 3 | `users.created_at` | sec | NN | `authsession/migration/migration.go:26-34`; `accounts.go:24-29,280-298` |
| 4 | `users.updated_at` | sec | NN | same as #3; `users_repository.go:101-148` |
| 5 | `users.locked_until` | sec | D0=0 | `authsession/migration/migration.go:121-127,431-436`; `accounts.go:219-237`; `auth.go:211-216` |
| 6 | `users.last_login_failure_at` | sec | D0=0 | `authsession/migration/migration.go:190-220`; `accounts.go:194-201,234-237` |
| 7 | `refresh_tokens.expires_at` | sec | NN | `authsession/migration/migration.go:35-42,353-360`; `account_operations.go:27-47` |
| 8 | `refresh_tokens.revoked_at` | sec | N | same; `account_operations.go:42-47` |
| 9 | `refresh_tokens.created_at` | sec | NN | same; `account_operations.go:47` |
| 10 | `roles.created_at` | sec | NN | `authsession/migration/migration.go:47-54`; `roles_repository.go:93-97,404-417` |
| 11 | `roles.updated_at` | sec | NN | same; `roles_repository.go:132-137` |
| 12 | `permissions.created_at` | sec | NN | `authsession/migration/migration.go:61-67`; system-data reconcile |
| 13 | `permissions.updated_at` | sec | NN | same |
| 14 | `menu_items.created_at` | sec | NN | `authsession/migration/migration.go:74-82`; system-data reconcile |
| 15 | `menu_items.updated_at` | sec | NN | same |
| 16 | `email_verification_challenges.expires_at` | sec | NN | `authsession/migration/migration.go:170-187`; `email_identity.go:174-204` |
| 17 | `email_verification_challenges.sent_at` | sec | NN | same |
| 18 | `password_recovery_challenges.expires_at` | sec | NN | `authsession/migration/migration.go:231-247`; `recovery.go:156-173` |
| 19 | `password_recovery_challenges.sent_at` | sec | NN | same |
| 20 | `login_failures.locked_until` | sec | D0=0 | `authsession/migration/migration.go:199-220`; `accounts_lock_source.go:67-128` |
| 21 | `login_failures.updated_at` | sec | NN | same; `accounts_lock_source.go:67-81` |
| 22 | `user_password_history.created_at` | sec | NN | `authsession/migration/migration.go:268-284`; `password_policy.go:169-185` |
| 23 | `user_invites.expires_at` | sec | NN | `authsession/migration/migration.go:287-322`; `invites.go:73-100,136-142` |
| 24 | `user_invites.consumed_at` | sec | N | same; `invites.go:90-97` |
| 25 | `user_invites.revoked_at` | sec | N | same |
| 26 | `user_invites.last_sent_at` | sec | NN | same; `invites.go:336-342` |
| 27 | `user_invites.created_at` | sec | NN | same |
| 28 | `service_credentials.expires_at` | sec | NN | `authsession/migration/migration.go:438-476`; `service_credentials.go:52-58,226-249` |
| 29 | `service_credentials.revoked_at` | sec | N | same; `service_credentials.go:165-199,243-249` |
| 30 | `service_credentials.last_used_at` | sec | N | same |
| 31 | `service_credentials.created_at` | sec | NN | same |
| 32 | `service_credentials.updated_at` | sec | NN | same |
| 33 | `core.persistence.mail_outbox.created_at` | ms | NN | `corepersistence/migration/migration.go:48-74`; `internal/mail/outbox.go:95-102,222-240` |
| 34 | `core.persistence.mail_config.updated_at` | ms | D0=0 | `corepersistence/migration/migration.go:87-124`; `internal/mail/runtime.go:169-176,308,382-390` |
| 35 | `core.operationlog.operation_log.created_at` | ms | NN | `operationlog/migration/migration.go:13-24`; `operationlog/repository.go:172-177,312-337` |
| 36 | `core.operationlog.operation_log_archive.created_at` | ms | NN | `operationlog/migration/migration.go:525-537`; `retention.go:21-42` |
| 37 | `core.operationlog.operation_log_archive.archived_at` | ms | NN | same; `retention.go:31-42` |
| 38 | `admin.notifications.notifications.read_at` | sec | N | `notifications/migration/migration.go:13-40`; `notifications_repository.go:85-95,186-275` |
| 39 | `admin.notifications.notifications.created_at` | sec | NN | same |
| 40 | `admin.settings.site_settings.updated_at` | sec | NN | `settings/migration/migration.go:22-42`; `settings/repository.go:204-244,331-363` |
| 41 | `core.jobs.jobs.lease_expires_at` | ms | N | `jobs/migration/migration.go:14-87`; `internal/jobs/model.go:136-146`; repository `83-112` |
| 42 | `core.jobs.jobs.created_at` | ms | NN | same; repository `40-53` |
| 43 | `core.jobs.jobs.updated_at` | ms | NN | same; repository `110-121` |
| 44 | `core.jobs.jobs.finished_at` | ms | N | same; `nullableTime` |
| 45 | `core.jobs.jobs.expires_at` | ms | N | same; actions/list `38-75,243-301` |
| 46 | `admin.data-dictionary.dict_types.created_at` | sec | NN | `datadictionary/migration/migration.go:16-29`; repository `109-187` |
| 47 | `admin.data-dictionary.dict_types.updated_at` | sec | NN | same |
| 48 | `admin.data-dictionary.dict_entries.created_at` | sec | NN | `datadictionary/migration/migration.go:30-41`; repository `292-393` |
| 49 | `admin.data-dictionary.dict_entries.updated_at` | sec | NN | same |
| 50 | `admin.data-permission.data_scope_policies.updated_at` | sec | NN | `datapermission/migration/migration.go:16-26`; repository `68-136` |
| 51 | `admin.data-permission.user_data_scopes.updated_at` | sec | NN | `datapermission/migration/migration.go:27-33`; repository `148-211` |
| 52 | `admin.login-captcha.captcha_challenges.expires_at` | sec | NN | `logincaptcha/migration/migration.go:15-24`; repository `36-49` |
| 53 | `admin.login-captcha.captcha_challenges.created_at` | sec | NN | same |
| 54 | `admin.login-captcha.captcha_config.created_at` | sec | NN | `logincaptcha/migration/migration.go:25-30`; repository `115-122` |
| 55 | `admin.login-captcha.captcha_config.updated_at` | sec | NN | same |
| 56 | `admin.recycle-bin.recycle_items.deleted_at` | sec | NN | `recyclebin/migration/migration.go:16-31`; repository `81-89,141-196` |
| 57 | `admin.recycle-bin.recycle_items.restored_at` | sec | N | same; partial index `WHERE restored_at IS NULL` |
| 58 | `admin.scheduled-tasks.scheduled_tasks.created_at` | sec | NN | `scheduledtasks/migration/migration.go:15-28`; repository `114-193` |
| 59 | `admin.scheduled-tasks.scheduled_tasks.updated_at` | sec | NN | same |
| 60 | `admin.scheduled-tasks.task_runs.started_at` | sec | NN | `scheduledtasks/migration/migration.go:29-38`; repository `259-309` |
| 61 | `admin.scheduled-tasks.task_runs.finished_at` | sec | N + runtime 0 sentinel | same; repository `259-309,343-362` |
| 62 | `admin.scheduled-tasks.task_runs.created_at` | sec | NN | same |
| 63 | `admin.mfa.user_mfa.created_at` | sec | NN | `mfa/migration/migration.go:16-29`; repository `59-237` |
| 64 | `admin.mfa.user_mfa.updated_at` | sec | NN | same |
| 65 | `admin.mfa.mfa_proofs.expires_at` | sec | NN | `mfa/migration/migration.go:30-57`; repository `261-307` |
| 66 | `admin.mfa.mfa_proofs.created_at` | sec | NN | same |
| 67 | `admin.wallet.wallet_accounts.created_at` | sec | NN | `wallet/migration/migration.go:16-36`; repository `250-520` |
| 68 | `admin.wallet.wallet_accounts.updated_at` | sec | NN | same |
| 69 | `admin.wallet.wallet_ledger_entries.created_at` | sec | NN | `wallet/migration/migration.go:37-64`; repository `637-799` |
| 70 | `admin.wallet.wallet_reconciliation_runs.created_at` | sec | NN | same; repository `886-1003` |
| 71 | `admin.wallet.subjects.created_at` | sec | NN | `wallet/migration/migration.go:275-285`; `wallet/subject/subject.go:86-191` |
| 72 | `admin.wallet.vouchers.expires_at` | sec | N; legacy <=0 treated absent | `wallet/migration/migration.go:287-347`; `wallet/voucher/service.go:79-108,325-352` |
| 73 | `admin.wallet.vouchers.redeemed_at` | sec | N; legacy <=0 treated absent | same; `service.go:219-223,399-419` |
| 74 | `admin.wallet.vouchers.created_at` | sec | NN | same |
| 75 | `admin.wallet.vouchers.updated_at` | sec | NN | same |
| 76 | `admin.wallet.voucher_batches.created_at` | sec | NN | `wallet/migration/migration.go:388-407`; voucher service batch writes |
| 77 | `admin.wallet.voucher_batches.updated_at` | sec | NN | same |
| 78 | `admin.channel.telegram.telegram_config.updated_at` | sec | D0=0 | `telegram/migration/migration.go:10-25`; runtime `internal/channel/telegram/runtime.go:166-177,395-405` |
| 79 | `admin.channel.telegram.telegram_sessions.last_message_at` | sec | NN | `telegram/migration/migration.go:33-46`; repository `161-284` |
| 80 | `admin.channel.telegram.telegram_sessions.created_at` | sec | NN | same |
| 81 | `admin.channel.telegram.telegram_sessions.updated_at` | sec | NN | same |
| 82 | `admin.channel.telegram.telegram_inbound_messages.received_at` | sec | NN | `telegram/migration/migration.go:47-64`; repository `125-179,375-384` |
| 83 | `admin.channel.telegram.telegram_outbound_messages.created_at` | sec | NN | `telegram/migration/migration.go:99-137`; repository `608-724` |
| 84 | `admin.channel.telegram.telegram_outbound_messages.updated_at` | sec | NN | same |
| 85 | `biz.digital-offer.digital_offers.created_at` | sec | NN | `digitaloffer/migration/migration.go:16-40`; `digitaloffer/store/store.go:154-314` |
| 86 | `biz.digital-offer.digital_offers.updated_at` | sec | NN | same |
| 87 | `biz.digital-offer.digital_purchases.created_at` | sec | NN | `digitaloffer/migration/migration.go:41-56`; store `351-432` |
| 88 | `biz.digital-offer.digital_entitlements.expires_at` | sec | N | `digitaloffer/migration/migration.go:57-73`; store `440-601` |
| 89 | `biz.digital-offer.digital_entitlements.created_at` | sec | NN | same |
| 90 | `biz.digital-offer.digital_entitlements.updated_at` | sec | NN | same |

## Historical retired source / excluded categories

- `core.persistence.records.updated_at` (historical v3 create, v6 `records_retire` drop; no current repository/handler): excluded from live 90, but C3 must assert current schema has no table before conversion.
- Excluded non-time integers: ID prefixes, `duration_seconds`, TOTP `last_used_step`, version, counters, amounts, balances, flags.
- TEXT/JSON payloads that may contain audit timestamps (`recycle_items.payload`, `operation_log.detail`) are not named time columns and remain outside this VP unless a later scope decision says otherwise.

## Current status

This artifact fixes the mechanical count and restores `login_failures.locked_until` and `login_failures.updated_at`. C1 inventory is now independently addable by column; C2/C3 codec, constraints, conversion, backup and predicate design remain open.
