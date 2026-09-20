---
id: r1-c2-per-column-conversion-contract-v1.0-fc
doc_type: design-attachment
title: R1 C2 逐列转换合同 v1.0（冻结候选）
status: freeze-candidate
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-002-r1-contract-and-denominator-freeze
version: 1.0.0
---

# R1 C2 逐列转换合同 v1.0（冻结候选）

> **状态：`freeze-candidate`，不是已实施事实。** 本文件按 `r1-c2-owner-migration-spec-v0.1.md`「Descriptor acceptance checklist」第 1/2/3/4 项展开为 90 行逐列合同，用于闭合 A-029（现行 independent head，`conditional`，open required = 5）中 **F-I-002**（逐列 USING/rebuild/codec）、**F-I-003**（90 列 old→new→read/write）与 **F-I-005**（descriptor 名 / canonical SQL / checksum 记录结构）的剩余子项。它**不**声称任何 DDL、codec、迁移或测试已实现；R2 才落码。
>
> 生成时间点：A-029 为最新独立意见；本文件与同一轮落盘的 `r1-c2-predicate-exact-sql-v1.0-fc.md`、`r1-c2-descriptor-ledger-v1.0-fc.md` 均为**待 independent 复审**的冻结候选。

## 0. 与既有载体的关系（单一真相收口）

| 载体 | 角色 | 本文件的关系 |
|------|------|--------------|
| `r1-time-column-inventory-v0.3.md` | 分母与旧单位/空值事实（90 列，A-006 接受） | 本文件的 `#N` 与该文件逐列一一对应 |
| `r1-c2-column-contract-matrix-v0.2.md` | 6-key 映射族 | 本文件是其逐列展开；两者精度规则必须同一 |
| `r1-c2-owner-migration-spec-v0.1.md` | owner / version / descriptor 名 | 本文件是其验收清单第 1/2/3/4 项 |
| `r1-c2-predicate-exact-sql-v1.0-fc.md` | 谓词 exact old/new SQL | 本文件「read/write」列引用其 `P-*` 行号，不重复 SQL |
| `r1-c3-backup-recovery-boundary-v1.0-fc.md` | C3 备份/回滚边界 | **本文件不覆盖 C3。** 该 C3 附件**已落盘**（2026-09-20，`status: freeze-candidate`），为 C3 的**唯一权威**；原 `r1-c3-backup-restore-runbook-v0.1.md` 已标 `superseded`（响应 A-040 §G 第 1 项）。 |

## 1. 冻结表达式（唯一形态，全文件只允许这三种）

```sql
-- (E1) PG legacy seconds，非 sentinel，NOT NULL
ALTER TABLE "<table>" ALTER COLUMN "<col>" TYPE timestamptz(6)
  USING (date_trunc('microseconds', to_timestamp("<col>"::double precision)));

-- (E2) PG legacy milliseconds，非 sentinel，NOT NULL
ALTER TABLE "<table>" ALTER COLUMN "<col>" TYPE timestamptz(6)
  USING (date_trunc('microseconds', TIMESTAMPTZ 'epoch' + "<col>" * INTERVAL '1 millisecond'));

-- (E3) PG sentinel 0 → NULL 前置（D0 列专用；随后套用 E1 或 E2 的 ELSE 分支）
ALTER TABLE "<table>" ALTER COLUMN "<col>" DROP DEFAULT;
ALTER TABLE "<table>" ALTER COLUMN "<col>" DROP NOT NULL;
ALTER TABLE "<table>" ALTER COLUMN "<col>" TYPE timestamptz(6)
  USING (CASE WHEN "<col>" = 0 THEN NULL
              ELSE date_trunc('microseconds', to_timestamp("<col>"::double precision)) END);
-- 毫秒 D0 列（#34）的 ELSE 主体用 E2：TIMESTAMPTZ 'epoch' + "<col>" * INTERVAL '1 millisecond'
```

- **nullable 列**：`CASE WHEN "<col>" IS NULL THEN NULL ELSE <E1/E2 主体> END`。
- **秒族（E1）= 三份载体 A-020 已接受的 `to_timestamp(double)` + `date_trunc` 式，本次**保留不改**（用户 2026-09-20 裁决 B）；毫秒族（E2）= `TIMESTAMPTZ 'epoch' + n * INTERVAL '1 millisecond'`，全程不经二进制浮点（A-016 F-I-002.3 要求）。**
- **`/1000.0` 与 `::timestamptz(6)` typmod-only cast 全包禁止**（A-018 点名）。typmod 的 `p` 默认 **round**，不得被当成已实现 Root D-008 的向零截断；截断由新写入的 Go `UTC().Truncate(time.Microsecond)` 与 E1/E2 的 `date_trunc` 承担。
- **禁止**：`/1000.0`、`to_timestamp(n / 1000.0)`、raw fractional 值直接 `::timestamptz(6)` typmod cast、`+00:00` 或本地时区文本。
- **口径收口（对应 A-027/A-029 点名的 D-015 字面差）**：Root `D-015-negative-instant-truncation.md` 原写「legacy **sec/ms 均** `date_trunc` + 整数 interval」；实际秒族是 E1（`to_timestamp(double)`），只有毫秒族是整数 interval。本次按裁决 B 修订 **D-015 字面**（毫秒族整数 interval / 秒族保留 `to_timestamp(double)`），**不改三份载体的秒式**，因此不新增第二套秒列 SQL 家族。
- **SQLite**：不使用任何 SQLite 日期函数。逐列 `CREATE TABLE <t>_new (...)` + `INSERT INTO <t>_new SELECT`（经 `apps/api/internal/temporal` codec 读旧整数、绑定 canonical fixed-6 UTC TEXT）+ 重建索引/CHECK/FK + `DROP`/`RENAME`，全程在 store 事务内。

## 2. 逐列合同（90 行）

列含义：`#` = inventory v0.3 行号；`old` = 现行 SQLite `INTEGER` / PG `BIGINT` 单位与空值；`PG target` = 目标物理形状；`PG USING` = §1 表达式编号 + 该列特有点；`SQLite target` = 目标形状；`read/write` = codec 方向（`R:` 读、`W:` 写）；`tests` = 用例 ID（`T-<#>-<col>-<kind>`，`RT`=round-trip、`SORT`=词法/物理序、`NULL`=空值、`ZL`=sentinel 0、`NEG`=负值、`ORD`=索引/CHECK 重建）。

| # | v | owner | table.column | old | PG target | PG USING | SQLite target | read/write | tests |
|--:|---|-------|--------------|-----|-----------|----------|---------------|-----------|-------|
| 1 | 73 | `core.persistence` | `schema_migrations.applied_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-1-RT`,`T-1-SORT` |
| 2 | 74 | `core.auth-session` | `system_data_reconcile.applied_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-2-RT`,`T-2-SORT` |
| 3 | 74 | `core.auth-session` | `users.created_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-3-RT`,`T-3-SORT` |
| 4 | 74 | `core.auth-session` | `users.updated_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:**Root** D-013 单调 `max(truncNow, old+1µs)` | `T-4-RT`,`T-4-MONO` |
| 5 | 74 | `core.auth-session` | `users.locked_until` | sec D0 | `timestamptz(6) NULL`，去 `DEFAULT 0` | E3 | `TEXT NULL` | R:NULL→未锁定；W:NULL 写入 | `T-5-ZL`,`T-5-NULL` |
| 6 | 74 | `core.auth-session` | `users.last_login_failure_at` | sec D0 | `timestamptz(6) NULL`，去 `DEFAULT 0` | E3 | `TEXT NULL` | R:NULL→无近期失败；W:NULL 写入 | `T-6-ZL`,`T-6-NULL` |
| 7 | 74 | `core.auth-session` | `refresh_tokens.expires_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-7-RT`,`T-7-SORT` |
| 8 | 74 | `core.auth-session` | `refresh_tokens.revoked_at` | sec N | `timestamptz(6) NULL` | E1（NULL 分支） | `TEXT NULL` | R:NULL→未撤销；W:NULL 保持 | `T-8-NULL`,`T-8-RT` |
| 9 | 74 | `core.auth-session` | `refresh_tokens.created_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-9-RT`,`T-9-SORT` |
| 10 | 74 | `core.auth-session` | `roles.created_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-10-RT`,`T-10-SORT` |
| 11 | 74 | `core.auth-session` | `roles.updated_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:**Root** D-013 单调 | `T-11-RT`,`T-11-MONO` |
| 12 | 74 | `core.auth-session` | `permissions.created_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-12-RT`,`T-12-SORT` |
| 13 | 74 | `core.auth-session` | `permissions.updated_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-13-RT`,`T-13-SORT` |
| 14 | 74 | `core.auth-session` | `menu_items.created_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-14-RT`,`T-14-SORT` |
| 15 | 74 | `core.auth-session` | `menu_items.updated_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-15-RT`,`T-15-SORT` |
| 16 | 74 | `core.auth-session` | `email_verification_challenges.expires_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-16-RT`,`T-16-SORT` |
| 17 | 74 | `core.auth-session` | `email_verification_challenges.sent_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-17-RT`,`T-17-SORT` |
| 18 | 74 | `core.auth-session` | `password_recovery_challenges.expires_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-18-RT`,`T-18-SORT` |
| 19 | 74 | `core.auth-session` | `password_recovery_challenges.sent_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-19-RT`,`T-19-SORT` |
| 20 | 74 | `core.auth-session` | `login_failures.locked_until` | sec D0 | `timestamptz(6) NULL`，去 `DEFAULT 0` | E3 | `TEXT NULL` | R:NULL→无源锁；W:NULL 写入 | `T-20-ZL`,`T-20-NULL` |
| 21 | 74 | `core.auth-session` | `login_failures.updated_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-21-RT`,`T-21-SORT` |
| 22 | 74 | `core.auth-session` | `user_password_history.created_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-22-RT`,`T-22-SORT` |
| 23 | 74 | `core.auth-session` | `user_invites.expires_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-23-RT`,`T-23-SORT` |
| 24 | 74 | `core.auth-session` | `user_invites.consumed_at` | sec N | `timestamptz(6) NULL` | E1（NULL 分支） | `TEXT NULL` | R:NULL→未消费；W:NULL 保持 | `T-24-NULL`,`T-24-RT` |
| 25 | 74 | `core.auth-session` | `user_invites.revoked_at` | sec N | `timestamptz(6) NULL` | E1（NULL 分支） | `TEXT NULL` | R:NULL→未撤销；W:NULL 保持 | `T-25-NULL`,`T-25-RT` |
| 26 | 74 | `core.auth-session` | `user_invites.last_sent_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-26-RT`,`T-26-SORT` |
| 27 | 74 | `core.auth-session` | `user_invites.created_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-27-RT`,`T-27-SORT` |
| 28 | 74 | `core.auth-session` | `service_credentials.expires_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-28-RT`,`T-28-SORT` |
| 29 | 74 | `core.auth-session` | `service_credentials.revoked_at` | sec N | `timestamptz(6) NULL` | E1（NULL 分支） | `TEXT NULL` | R:NULL→未撤销；W:NULL 保持 | `T-29-NULL`,`T-29-RT` |
| 30 | 74 | `core.auth-session` | `service_credentials.last_used_at` | sec N | `timestamptz(6) NULL` | E1（NULL 分支） | `TEXT NULL` | R:NULL→从未使用；W:NULL 保持 | `T-30-NULL`,`T-30-RT` |
| 31 | 74 | `core.auth-session` | `service_credentials.created_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-31-RT`,`T-31-SORT` |
| 32 | 74 | `core.auth-session` | `service_credentials.updated_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-32-RT`,`T-32-SORT` |
| 33 | 73 | `core.persistence` | `mail_outbox.created_at` | ms NN | `timestamptz(6) NOT NULL` | E2 | `TEXT NOT NULL` | R:`FromUnixMilli`；W:`Truncate(µs)` | `T-33-RT`,`T-33-SORT` |
| 34 | 73 | `core.persistence` | `mail_config.updated_at` | ms D0 | `timestamptz(6) NULL`，去 `DEFAULT 0` | E3（毫秒 ELSE 主体 = **E2**：`date_trunc('microseconds', TIMESTAMPTZ 'epoch' + "<col>" * INTERVAL '1 millisecond')`） | `TEXT NULL` | R:NULL→未初始化；W:NULL 写入 | `T-34-ZL`,`T-34-NULL` |
| 35 | 75 | `core.operationlog` | `operation_log.created_at` | ms NN | `timestamptz(6) NOT NULL` | E2 | `TEXT NOT NULL` | R:`FromUnixMilli`；W:`Truncate(µs)` | `T-35-RT`,`T-35-SORT` |
| 36 | 75 | `core.operationlog` | `operation_log_archive.created_at` | ms NN | `timestamptz(6) NOT NULL` | E2 | `TEXT NOT NULL` | R:`FromUnixMilli`；W:`Truncate(µs)` | `T-36-RT`,`T-36-SORT` |
| 37 | 75 | `core.operationlog` | `operation_log_archive.archived_at` | ms NN | `timestamptz(6) NOT NULL` | E2 | `TEXT NOT NULL` | R:`FromUnixMilli`；W:`Truncate(µs)` | `T-37-RT`,`T-37-SORT` |
| 38 | 81 | `admin.notifications` | `notifications.read_at` | sec N | `timestamptz(6) NULL` | E1（NULL 分支） | `TEXT NULL` | R:NULL→未读；W:NULL 保持 | `T-38-NULL`,`T-38-RT` |
| 39 | 81 | `admin.notifications` | `notifications.created_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-39-RT`,`T-39-SORT` |
| 40 | 84 | `admin.settings` | `site_settings.updated_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-40-RT`,`T-40-SORT` |
| 41 | 76 | `core.jobs` | `jobs.lease_expires_at` | ms N | `timestamptz(6) NULL` | E2（NULL 分支） | `TEXT NULL` | R:NULL→无租约；W:NULL 保持 | `T-41-NULL`,`T-41-RT` |
| 42 | 76 | `core.jobs` | `jobs.created_at` | ms NN | `timestamptz(6) NOT NULL` | E2 | `TEXT NOT NULL` | R:`FromUnixMilli`；W:`Truncate(µs)` | `T-42-RT`,`T-42-SORT` |
| 43 | 76 | `core.jobs` | `jobs.updated_at` | ms NN | `timestamptz(6) NOT NULL` | E2 | `TEXT NOT NULL` | R:`FromUnixMilli`；W:`Truncate(µs)` | `T-43-RT`,`T-43-SORT` |
| 44 | 76 | `core.jobs` | `jobs.finished_at` | ms N | `timestamptz(6) NULL` | E2（NULL 分支） | `TEXT NULL` | R:NULL→未完成；W:NULL 保持 | `T-44-NULL`,`T-44-RT` |
| 45 | 76 | `core.jobs` | `jobs.expires_at` | ms N | `timestamptz(6) NULL` | E2（NULL 分支） | `TEXT NULL` | R:NULL→无过期；W:NULL 保持 | `T-45-NULL`,`T-45-RT` |
| 46 | 77 | `admin.data-dictionary` | `dict_types.created_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-46-RT`,`T-46-SORT` |
| 47 | 77 | `admin.data-dictionary` | `dict_types.updated_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-47-RT`,`T-47-SORT` |
| 48 | 77 | `admin.data-dictionary` | `dict_entries.created_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-48-RT`,`T-48-SORT` |
| 49 | 77 | `admin.data-dictionary` | `dict_entries.updated_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-49-RT`,`T-49-SORT` |
| 50 | 78 | `admin.data-permission` | `data_scope_policies.updated_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-50-RT`,`T-50-SORT` |
| 51 | 78 | `admin.data-permission` | `user_data_scopes.updated_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-51-RT`,`T-51-SORT` |
| 52 | 79 | `admin.login-captcha` | `captcha_challenges.expires_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-52-RT`,`T-52-SORT` |
| 53 | 79 | `admin.login-captcha` | `captcha_challenges.created_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-53-RT`,`T-53-SORT` |
| 54 | 79 | `admin.login-captcha` | `captcha_config.created_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-54-RT`,`T-54-SORT` |
| 55 | 79 | `admin.login-captcha` | `captcha_config.updated_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-55-RT`,`T-55-SORT` |
| 56 | 82 | `admin.recycle-bin` | `recycle_items.deleted_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-56-RT`,`T-56-SORT` |
| 57 | 82 | `admin.recycle-bin` | `recycle_items.restored_at` | sec N | `timestamptz(6) NULL` | E1（NULL 分支） | `TEXT NULL` | R:NULL→未恢复；W:NULL 保持 | `T-57-NULL`,`T-57-ORD` |
| 58 | 83 | `admin.scheduled-tasks` | `scheduled_tasks.created_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-58-RT`,`T-58-SORT` |
| 59 | 83 | `admin.scheduled-tasks` | `scheduled_tasks.updated_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-59-RT`,`T-59-SORT` |
| 60 | 83 | `admin.scheduled-tasks` | `task_runs.started_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-60-RT`,`T-60-SORT` |
| 61 | 83 | `admin.scheduled-tasks` | `task_runs.finished_at` | sec N | `timestamptz(6) NULL` | E1（NULL 分支） | `TEXT NULL` | R:NULL→未完成；W:去掉 `COALESCE(…,0)` 与写 0 | `T-61-NULL`,`T-61-RT` |
| 62 | 83 | `admin.scheduled-tasks` | `task_runs.created_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-62-RT`,`T-62-SORT` |
| 63 | 80 | `admin.mfa` | `user_mfa.created_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-63-RT`,`T-63-SORT` |
| 64 | 80 | `admin.mfa` | `user_mfa.updated_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-64-RT`,`T-64-SORT` |
| 65 | 80 | `admin.mfa` | `mfa_proofs.expires_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-65-RT`,`T-65-SORT` |
| 66 | 80 | `admin.mfa` | `mfa_proofs.created_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-66-RT`,`T-66-SORT` |
| 67 | 85 | `admin.wallet` | `wallet_accounts.created_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-67-RT`,`T-67-SORT` |
| 68 | 85 | `admin.wallet` | `wallet_accounts.updated_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-68-RT`,`T-68-SORT` |
| 69 | 85 | `admin.wallet` | `wallet_ledger_entries.created_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)`；序 `(created_at,id)` | `T-69-RT`,`T-69-SORT` |
| 70 | 85 | `admin.wallet` | `wallet_reconciliation_runs.created_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-70-RT`,`T-70-SORT` |
| 71 | 85 | `admin.wallet` | `subjects.created_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-71-RT`,`T-71-SORT` |
| 72 | 85 | `admin.wallet` | `vouchers.expires_at` | sec N（legacy ≤0 视为缺失） | `timestamptz(6) NULL` | E3（`= 0` 单分支；**`< 0` 不进 USING**，只走 `m0` 预检 fail closed） | `TEXT NULL` | R:NULL→无过期；W:负值 fail closed | `T-72-ZL`,`T-72-NEG`,`T-72-NULL` |
| 73 | 85 | `admin.wallet` | `vouchers.redeemed_at` | sec N（legacy ≤0 视为缺失） | `timestamptz(6) NULL` | E3（`= 0` 单分支；**`< 0` 不进 USING**，只走 `m0` 预检 fail closed） | `TEXT NULL` | R:NULL→未兑换；W:负值 fail closed | `T-73-ZL`,`T-73-NEG`,`T-73-NULL` |
| 74 | 85 | `admin.wallet` | `vouchers.created_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-74-RT`,`T-74-SORT` |
| 75 | 85 | `admin.wallet` | `vouchers.updated_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-75-RT`,`T-75-SORT` |
| 76 | 85 | `admin.wallet` | `voucher_batches.created_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-76-RT`,`T-76-SORT` |
| 77 | 85 | `admin.wallet` | `voucher_batches.updated_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-77-RT`,`T-77-SORT` |
| 78 | 86 | `admin.channel.telegram` | `telegram_config.updated_at` | sec D0 | `timestamptz(6) NULL`，去 `DEFAULT 0` | E3 | `TEXT NULL` | R:NULL→未配置；W:NULL 写入 | `T-78-ZL`,`T-78-NULL` |
| 79 | 86 | `admin.channel.telegram` | `telegram_sessions.last_message_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-79-RT`,`T-79-SORT` |
| 80 | 86 | `admin.channel.telegram` | `telegram_sessions.created_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-80-RT`,`T-80-SORT` |
| 81 | 86 | `admin.channel.telegram` | `telegram_sessions.updated_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-81-RT`,`T-81-SORT` |
| 82 | 86 | `admin.channel.telegram` | `telegram_inbound_messages.received_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-82-RT`,`T-82-SORT` |
| 83 | 86 | `admin.channel.telegram` | `telegram_outbound_messages.created_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-83-RT`,`T-83-SORT` |
| 84 | 86 | `admin.channel.telegram` | `telegram_outbound_messages.updated_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-84-RT`,`T-84-SORT` |
| 85 | 87 | `biz.digital-offer` | `digital_offers.created_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-85-RT`,`T-85-SORT` |
| 86 | 87 | `biz.digital-offer` | `digital_offers.updated_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-86-RT`,`T-86-SORT` |
| 87 | 87 | `biz.digital-offer` | `digital_purchases.created_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-87-RT`,`T-87-SORT` |
| 88 | 87 | `biz.digital-offer` | `digital_entitlements.expires_at` | sec N | `timestamptz(6) NULL` | E1（NULL 分支） | `TEXT NULL` | R:NULL→长期有效；W:NULL 保持 | `T-88-NULL`,`T-88-ORD` |
| 89 | 87 | `biz.digital-offer` | `digital_entitlements.created_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-89-RT`,`T-89-SORT` |
| 90 | 87 | `biz.digital-offer` | `digital_entitlements.updated_at` | sec NN | `timestamptz(6) NOT NULL` | E1 | `TEXT NOT NULL` | R:`FromUnix`；W:`Truncate(µs)` | `T-90-RT`,`T-90-SORT` |

**行数与分配自检**：90 行 = `sec 80 + ms 10` = `NN 71 + N 14 + D0 5`；owner 分配 = v73 3 / v74 31 / v75 3 / v76 5 / v77 4 / v78 2 / v79 4 / v80 4 / v81 2 / v82 2 / v83 5 / v84 1 / v85 11 / v86 7 / v87 6 = **90**。v74 = inventory `#2`–`#32`（31 列），**不含** `#1 schema_migrations.applied_at`（归 v73）——即 **v74 表范围 = `system_data_reconcile` + auth 表，`system_data_reconcile` 归 v74，`schema_migrations` 归 v73**。这与 `r1-v73-owner-allocation-draft-v0.1.md` L18/L19 的「#1 归 73 / #2–#32 excluding #1 归 74」一致，消掉 A-027 F-I-005 点名的 v74「ledger/reconcile」 vs 「schema/system ledger」歧义。

## 3. C2 冻结前必须收口的 3 项（本候选仍未满足）

1. **D-015 字面已按裁决 B 收口并落盘**：秒族保留三份载体已接受的 E1（`to_timestamp(double)` + `date_trunc`），毫秒族为整数 interval。Root `D-015-negative-instant-truncation.md` 与 child `D-012-v73-allocation-negative-truncation.md` 已改为分族表述（E-026）。**A-030 已确认该子项 `fixed`**；本文件仍为冻结候选（整条 F-I-002 未闭合，见 §3.4 与 A-030）。
2. **不可逆点清单**：已在 `r1-c2-predicate-exact-sql-v1.0-fc.md` §6 单列（0→NULL、精度截断、非规范 TEXT fail closed）；本文件 §2 的 `ZL`/`NEG`/`NULL` 用例族与之对应。
3. **per-owner `MigrationChecksum` 尚未记录**：见 `r1-c2-descriptor-ledger-v1.0-fc.md`（结构、名称、transform ID 已定，哈希值需 R2 落码后填充）。
4. **逐表 exact SQLite rebuild DDL 仍未写出**（A-030 F-I-002.1）：本文件 §1 给的是共享模板与逐列目标形状；每张受影响表的完整 `<t>_new` DDL + `INSERT SELECT` 正文仍是 C2 冻结前必交项。
5. **非法值单路径**：`< 0` 只走 `m0` 预检 fail closed，不写 USING 分支——该唯一机制见 `r1-c2-predicate-exact-sql-v1.0-fc.md` §1 说明块与 §6 `m0`。
6. **双方言 checksum 约定二选一**：未选定（P-004 待用户/编排器裁决；见 descriptor 台账 §5.3）。

## 4. 声明

- 本文件为**冻结候选**，`status: freeze-candidate`；不改变任何 `status` / `progress` / goal-tree。
- 本文件**不**声称 `apps/` 下任何 DDL、codec、迁移或测试已实现。A-027 F-I-002/F-I-003/F-I-005 的闭合需 independent 复审本文件 + 其配套 `P-*` 谓词表与 descriptor 台账。
- 引用 D-012 时一律限定 Root/child（Root `D-012` = voucher 异常值政策；child `D-012` = v73 allocation / 负瞬间承接），执行 F-I-018 的收口要求。
