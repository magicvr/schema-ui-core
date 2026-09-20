---
id: r1-c2-predicate-exact-sql-v1.0-fc
doc_type: design-attachment
title: R1 C2 谓词 exact old/new SQL 单表 v1.0（冻结候选）
status: freeze-candidate
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-002-r1-contract-and-denominator-freeze
version: 1.0.0
---

# R1 C2 谓词 exact old/new SQL 单表 v1.0（冻结候选）

> **状态：`freeze-candidate`。** 本文件是 A-029（现行 independent head）**F-I-006** 关闭要求点名的「**一张** exact old/new SQL + 迁移顺序 + 负例/NULL/sentinel 测试 ID」单表；它收口 `r1-c2-predicate-index-matrix-v0.1.md`（dependency 清单）与 `r1-c2-readwrite-predicate-spec-v0.1.md`（family 级）两源，并把 **Root** `D-012`（voucher 异常值政策）与 **Root** `D-013`（单调 `updated_at`）写进同一张表。
>
> **D-012/D-013 编号一律 Root/child 限定**（A-029 **F-I-018** 关闭要求）：本文件中 `Root D-012` = `GOAL-001/01-decision/D-012-voucher-invalid-value-policy.md`；`Root D-013` = `GOAL-001/01-decision/D-013-monotonic-updated-at-policy.md`；child `D-011` = `GOAL-002/01-decision/D-011-voucher-monotonic-policies.md`；child `D-012` = `GOAL-002/01-decision/D-012-v73-allocation-negative-truncation.md`（v73/截断承接，**不是** voucher 政策）。
>
> **old 列引用的代码事实**在本次落盘时逐处核对（行号见每行 `callsite`）；`apps/` 代码**未修改**。

## 0. 表结构

- `disposition` 取值：`DB`（迁移改列，运行时 SQL 同步改类型参数）/ `RT`（仅运行时读写改，迁移不改列形状）/ `DB+RT`（两者都改）。
- 迁移顺序列给出同一 owner descriptor 内的语句次序（同 descriptor 内按 `m0..mN` 递增执行）。
- 测试 ID 沿用 `r1-c2-per-column-conversion-contract-v1.0-fc.md` §2 的 `T-<#>-<kind>`，并追加谓词专用 `P-<#>-<kind>`（`NEG` 负例 / `NULL` 空值 / `SENT` sentinel / `ORD` 序 / `ROLLBACK` 回滚）。

## 1. Sentinel 0→NULL 列（5 列，**Root** `D-012` 政策；child `D-011` 承接）

| # | 列 | owner | disposition | old SQL（现行，逐字） | new SQL（目标） | callsite | 测试 |
|--:|----|-------|-------------|------------------------|-----------------|----------|------|
| 5 | `users.locked_until` | v74 | DB+RT | 读：`u.locked_until > ?` / `u.locked_until <= ?`（整数秒参数）；写：`UPDATE users SET locked_until = ?`，插入时 `DEFAULT 0` | 读：`(u.locked_until IS NULL OR u.locked_until > ?)`（未锁定分支）/ `(u.locked_until IS NOT NULL AND u.locked_until <= ?)`（已锁定分支）；写：`time.Time` 绑定或 NULL | `apps/api/modules/authsession/users_repository.go:498,500`；DDL `authsession/migration/migration.go:126,435` | `P-5-SENT`,`P-5-NULL`,`P-5-ORD` |
| 6 | `users.last_login_failure_at` | v74 | DB+RT | 读：失败衰减窗口比较（整数秒）；DDL `last_login_failure_at INTEGER/BIGINT NOT NULL DEFAULT 0` | 读：`(u.last_login_failure_at IS NULL OR u.last_login_failure_at < ?)`（NULL = 无近期失败）；写：NULL 或 `Truncate(µs)` 时刻 | `authsession/migration/migration.go:208,220`；`authsession/users_repository.go`（衰减比较） | `P-6-SENT`,`P-6-NULL` |
| 20 | `login_failures.locked_until` | v74 | DB+RT | `UPDATE login_failures SET fail_count = CASE WHEN updated_at < ? THEN 1 ELSE fail_count + 1 END, updated_at = ? WHERE user_id = ? AND ip = ?`；插入 `INSERT INTO login_failures (user_id, ip, fail_count, locked_until, updated_at) VALUES (?, ?, 1, 0, ?)`；读 `SELECT locked_until … ` 后 Go 侧 `lockedUntil > now.Unix()`；开锁 `UPDATE login_failures SET locked_until = ?, fail_count = 0, updated_at = ? …` | 插入：`locked_until` 绑定 NULL；读：`SELECT locked_until` 扫入 `sql.NullTime` → `lockedUntil.Valid && lockedUntil.Time.After(nowUTC)`；开锁：绑定 `time.Time`；窗口比较参数改 `time.Time` | `authsession/accounts_lock_source.go:67-71,78-79,101-102,115-128`；DDL `authsession/migration/migration.go:204,216` | `P-20-SENT`,`P-20-NULL`,`P-20-ORD` |
| 34 | `mail_config.updated_at` | v73 | DB+RT | DDL `updated_at INTEGER/BIGINT NOT NULL DEFAULT 0`；读模型把 `0` 当未初始化 | DDL 去 `DEFAULT 0` + 可空；读：`updated_at IS NULL` = 未初始化；写：NULL 或 `Truncate(µs)` | `corepersistence/migration/migration.go:105,123`；`apps/api/internal/mail/runtime.go:169-176,308,382-390` | `P-34-SENT`,`P-34-NULL` |
| 78 | `telegram_config.updated_at` | v86 | DB+RT | DDL `updated_at INTEGER/BIGINT NOT NULL DEFAULT 0` | 同 #34 | `channel/telegram/migration/migration.go:15,24`；`apps/api/internal/channel/telegram/runtime.go:166-177,395-405` | `P-78-SENT`,`P-78-NULL` |

> **Root** `D-012-voucher-invalid-value-policy.md` 的三桶预检要求（0 / 负值 / 正值分别计数）对 **#5/#6/#20/#34/#72/#73/#78** 强制；预检脚本计入 owner descriptor 的 `preflight` 步骤（§6 `m0`），计数不达标即 fail closed。**注意**：child `D-012-v73-allocation-negative-truncation.md` 是 v73 allocation / 负瞬间承接，与本节 voucher 政策无关（F-I-018）。

## 2. voucher / entitlement 异常值（**Root** `D-012` + child `D-011`，负值 fail closed）

| # | 列 | owner | disposition | old SQL（现行，逐字） | new SQL（目标） | callsite | 测试 |
|--:|----|-------|-------------|------------------------|-----------------|----------|------|
| 72 | `vouchers.expires_at` | v85 | DB+RT | `SELECT id, batch_id, code_prefix, amount, currency, status, expires_at, redeemed_by, redeemed_at, created_at, updated_at …` 后 Go 侧 `if exp.Valid && exp.Int64 > 0 {`（`<=0` 一律当缺失） | 迁移：`=0 → NULL`；`<0` → 迁移期报错 fail closed（不静默转 NULL）。运行时：`if exp.Valid {`（不再用 `>0` 判存在），过期判定用 `exp.Time.Before(nowUTC)` | `wallet/voucher/service.go:330,340` | `P-72-SENT`,`P-72-NEG` |
| 73 | `vouchers.redeemed_at` | v85 | DB+RT | 同上，`if redAt.Valid && redAt.Int64 > 0 {` | 迁移同 #72；运行时 `if redAt.Valid {` | `wallet/voucher/service.go:347` | `P-73-SENT`,`P-73-NEG` |
| 88 | `digital_entitlements.expires_at` | v87 | DB+RT | form CHECK：duration 要求 `expires_at IS NOT NULL`、count 要求 `expires_at IS NULL` | CHECK 谓词逐字保留，仅列类型改 `TEXT`/`timestamptz(6)`；重建后重放两类 form 断言 | `digitaloffer/migration/migration.go:57-73`；`digitaloffer/store/store.go:440-601` | `P-88-NULL`,`P-88-ORD` |

## 3. 单调 `updated_at`（**Root** `D-013` + child `D-011`）

| # | 列 | owner | disposition | old SQL（现行，逐字） | new SQL（目标） | callsite | 测试 |
|--:|----|-------|-------------|------------------------|-----------------|----------|------|
| 4 | `users.updated_at` | v74 | DB+RT | `nextUpdatedAt := now.Unix()`；`if nextUpdatedAt <= current.UpdatedAt.Unix() { nextUpdatedAt = current.UpdatedAt.Unix() + 1 }`；随后 `UPDATE users SET … updated_at = ?` | `nextUpdatedAt := now.UTC().Truncate(time.Microsecond)`；`if !nextUpdatedAt.After(current.UpdatedAt) { nextUpdatedAt = current.UpdatedAt.Add(time.Microsecond) }`；绑定 `time.Time` | `authsession/users_repository.go:232-234,253` | `P-4-MONO`（墙钟回拨 + 快速连续写） |
| 11 | `roles.updated_at` | v74 | DB+RT | 同上整数秒 `+1` 语义 | `max(Truncate(µs) now, old + 1µs)` | `authsession/roles_repository.go:132-137` | `P-11-MONO` |

> 语义等价性要求：单调不变量、cache/ETag 行为、`ORDER BY updated_at` 次序在三组用例下必须与迁移前一致（`P-4-MONO`/`P-11-MONO` 各含 wall-clock rollback、同微秒并发写、跨秒边界三例）。

## 4. 运行时空值/状态谓词（RT 为主）

| # | 列 | owner | disposition | old SQL（现行，逐字） | new SQL（目标） | callsite | 测试 |
|--:|----|-------|-------------|------------------------|-----------------|----------|------|
| 8 | `refresh_tokens.revoked_at` | v74 | RT | `UPDATE refresh_tokens SET revoked_at = ? WHERE user_id = ? AND revoked_at IS NULL`；`… WHERE id = ? AND revoked_at IS NULL` | 谓词逐字保留；参数改 `time.Time` | `authsession/accounts.go:257,271,356`；`account_operations.go:64` | `P-8-NULL` |
| 24/25 | `user_invites.consumed_at` / `revoked_at` | v74 | DB+RT | `consumed_at IS NULL AND revoked_at IS NULL AND expires_at > ?`（过期分支 `<= ?`）；`UPDATE user_invites SET revoked_at = ? WHERE id = ? AND consumed_at IS NULL AND revoked_at IS NULL` | 谓词逐字保留；`expires_at` 参数改 `time.Time` | `authsession/invites.go:200,206,297` | `P-24-NULL`,`P-25-NULL` |
| 29/30 | `service_credentials.revoked_at` / `last_used_at` | v74 | DB+RT | `UPDATE service_credentials SET revoked_at = ?, updated_at = ? WHERE id = ? AND revoked_at IS NULL`；`SET last_used_at = ?, updated_at = ? …` | 谓词逐字保留；参数改 `time.Time` | `authsession/service_credentials.go:165,185,198` | `P-29-NULL`,`P-30-NULL` |
| 38 | `notifications.read_at` | v81 | DB+RT | `where += " AND read_at IS NULL"`；`SELECT id FROM notifications WHERE user_id = ? AND read_at IS NOT NULL ORDER BY created_at ASC LIMIT ?`；`UPDATE notifications SET read_at = ? WHERE user_id = ? AND read_at IS NULL`；`SELECT COUNT(*) … WHERE user_id = ? AND read_at IS NULL` | 谓词逐字保留；`read_at` 参数改 `time.Time`；`ORDER BY created_at ASC` 在 fixed-6 TEXT 上等价 | `authsession/notifications_repository.go:106,150,156,164,220,240` | `P-38-NULL`,`P-38-ORD` |
| 57 | `recycle_items.restored_at` | v82 | DB+RT | 部分唯一索引 `CREATE UNIQUE INDEX idx_recycle_items_active ON recycle_items(resource, resource_id) WHERE restored_at IS NULL`；`where := "WHERE restored_at IS NULL"`；`UPDATE recycle_items SET restored_at = ? WHERE id = ? AND restored_at IS NULL`；`DELETE FROM recycle_items WHERE restored_at IS NULL` | 部分索引与三条谓词**逐字保留**，重建表后按同名重建索引；`restored_at` 参数改 `time.Time` | `recyclebin/migration/migration.go:30,48`；`recyclebin/store/repository.go:100,215,230` | `P-57-ORD`,`P-57-NULL` |
| 61 | `task_runs.finished_at` | v83 | DB+RT | `SELECT id, task_id, status, started_at, COALESCE(finished_at, 0), COALESCE(detail, ''), created_at FROM task_runs … ORDER BY started_at DESC`；运行时写 `0` | 去掉 `COALESCE(finished_at, 0)`，直接选 `finished_at` 扫入 `sql.NullTime`；运行时写 NULL；`ORDER BY started_at DESC` 保留 | `scheduledtasks/store/repository.go:289-290,343-344`；`259-309,343-362` | `P-61-NULL` |
| 41/44/45 | `jobs.lease_expires_at` / `finished_at` / `expires_at` | v76 | DB+RT | 六态 CHECK（逐字见下）；`OR (status='running' AND cancel_requested=0 AND lease_expires_at <= ? AND attempt < max_attempts)`；`WHERE id=? AND status='succeeded' AND expires_at <= ?`；`… AND expires_at <= ?`（清扫）；`ORDER BY created_at, id` | **CHECK 谓词逐字保留**（NULL 分支在类型转换后语义不变）；时间参数由 `toMillis(now)` 改 `time.Time`；`ORDER BY created_at, id` 保留 | `jobs/migration/migration.go:36-43,75-82`；`apps/api/internal/jobs/repository.go:89,244,256,266,271,282,288,300-301` | `P-41-ORD`,`P-44-NULL`,`P-45-NULL`,`P-41-ROLLBACK` |
| 65 | `mfa_proofs.expires_at` | v80 | DB+RT | `DELETE FROM mfa_proofs WHERE user_id = ? AND expires_at <= ?` | 参数改 `time.Time` | `mfa/store/repository.go:273` | `P-65-ORD` |
| 52 | `captcha_challenges.expires_at` | v79 | DB+RT | `DELETE FROM captcha_challenges WHERE expires_at <= ?`；`DELETE FROM captcha_challenges WHERE id = ? AND expires_at > ? AND answer_hash = ?` | 参数改 `time.Time` | `logincaptcha/store/repository.go:40,69` | `P-52-ORD` |
| 18 | `password_recovery_challenges.expires_at` | v74 | DB+RT | `DELETE FROM password_recovery_challenges WHERE user_id = ? AND expires_at <= ?` | 参数改 `time.Time` | `authsession/recovery.go:287` | `P-18-ORD` |
| 16 | `email_verification_challenges.expires_at` | v74 | DB+RT | `DELETE FROM email_verification_challenges WHERE user_id = ? AND expires_at <= ?` | 参数改 `time.Time` | `authsession/email_identity.go:259` | `P-16-ORD` |

六态 CHECK（`jobs`，迁移后逐字保留，列类型改 `TEXT`/`timestamptz(6)`）：

```sql
CHECK (
  (status = 'queued'    AND lease_owner IS NULL     AND lease_expires_at IS NULL     AND result IS NULL     AND error_code IS NULL     AND finished_at IS NULL     AND expires_at IS NULL)
  OR (status = 'running'   AND lease_owner IS NOT NULL AND lease_expires_at IS NOT NULL AND result IS NULL     AND error_code IS NULL     AND finished_at IS NULL     AND expires_at IS NULL)
  OR (status = 'succeeded' AND lease_owner IS NULL     AND lease_expires_at IS NULL     AND result IS NOT NULL AND error_code IS NULL     AND progress = 100 AND finished_at IS NOT NULL AND expires_at IS NOT NULL)
  OR (status = 'failed'    AND lease_owner IS NULL     AND lease_expires_at IS NULL     AND result IS NULL     AND error_code IS NOT NULL AND finished_at IS NOT NULL AND expires_at IS NULL)
  OR (status = 'cancelled' AND lease_owner IS NULL     AND lease_expires_at IS NULL     AND result IS NULL     AND error_code IS NULL     AND finished_at IS NOT NULL AND expires_at IS NULL)
  OR (status = 'expired'   AND lease_owner IS NULL     AND lease_expires_at IS NULL     AND result IS NULL     AND error_code IS NULL     AND progress = 100 AND finished_at IS NOT NULL AND expires_at IS NOT NULL)
)
```

## 5. 全 90 列谓词/索引/CHECK 变更清单（无遗漏口径）

下列列**当前没有任何时间谓词/索引/CHECK**（纯读写列），迁移后同样无谓词；本表以「none」显式登记，避免「未列出 = 未检查」：

| 处置 | 列（`#N`） | 计数 |
|------|------------|-----:|
| Sentinel 0→NULL（§1） | #5, #6, #20, #34, #78 | 5 |
| voucher/entitlement 异常值（§2） | #72, #73, #88 | 3 |
| 单调 updated_at（§3） | #4, #11 | 2 |
| 运行时空值/状态谓词（§4） | #8, #16, #18, #24, #25, #29, #30, #38, #41, #44, #45, #52, #57, #61, #65 | 15 |
| 仅 `ORDER BY`（保留，无类型相关改写） | #1, #2, #35, #36, #37, #39, #56, #62, #69, #70, #76, #82, #83, #87, #89 | 15 |
| **none**（无谓词/索引/CHECK） | #3, #7, #9, #10, #12, #13, #14, #15, #17, #19, #21, #22, #23, #26, #27, #28, #31, #32, #33, #40, #42, #43, #46, #47, #48, #49, #50, #51, #53, #54, #55, #58, #59, #60, #63, #64, #66, #67, #68, #71, #74, #75, #77, #79, #80, #81, #84, #85, #86, #90 | 50 |

合计 **5 + 3 + 2 + 15 + 15 + 50 = 90**。`ORDER BY` 行清单（逐字现行）：`wallet_accounts … ORDER BY created_at DESC, id DESC`（`wallet/store/repository.go:200`）、`wallet_ledger_entries … ORDER BY created_at DESC, id DESC`（`:782`）与 `ORDER BY created_at ASC, id ASC`（`:922`）、`wallet_reconciliation_runs … ORDER BY created_at DESC, id DESC`（`:987`）、`vouchers … ORDER BY created_at DESC`（`voucher/service.go:391`）、`mail_outbox … ORDER BY created_at DESC, id DESC`（`internal/mail/outbox.go:107`）、`operation_log` retention/archive 序（`operationlog/repository.go:172-177,312-337`；`retention.go:21-42`）、`telegram_sessions … last_message_at DESC` 与 inbound/outbound 序（`channel/telegram/store/repository.go:161-284,125-179,375-384,608-724`）、`digital_offers/purchases/entitlements` 序（`digitaloffer/store/store.go:233,292,419,480,573,651`）、`task_runs … ORDER BY started_at DESC`（`scheduledtasks/store/repository.go:289-290,343-344`）、`user_password_history … ORDER BY created_at DESC, id DESC`（`authsession/password_policy.go:148,184`）、`refresh_tokens` 序（`authsession/account_operations.go:29`）、`service_credentials … ORDER BY created_at DESC, id DESC`（`authsession/service_credentials.go:91`）、`notifications … ORDER BY created_at DESC, id DESC`（`authsession/notifications_repository.go:164`）。

**序语义前提**：SQLite 目标为定宽 fixed-6 UTC TEXT，其**词法序 = 时刻序**（`T-*-SORT` 用例族逐列验证）；PG 为原生 `timestamptz`。所有 `ORDER BY` 的 id tie-break 保持不变；`idx_jobs_runnable` / `idx_jobs_actor` / `idx_jobs_expiry` / `idx_jobs_created_at`、`idx_recycle_items_active` 重建后列序与部分谓词逐字不变（`P-41-ORD`/`P-57-ORD`）。

## 6. 迁移顺序（每个 owner descriptor 内）

```text
m0  preflight：逐列 0 / 负值 / 正值分桶计数（Root D-012）；非法值或计数不符 → 事务回滚，fail closed
m1  放开受影响约束：DROP DEFAULT / DROP NOT NULL / DROP CHECK / DROP 部分索引（仅被本 descriptor 触及者）
m2  逐列转换：PG 用 §1 E1/E2/E3 的 ALTER … TYPE … USING；SQLite 用 <t>_new 重建 + codec 拷贝
m3  重建约束：CHECK / NOT NULL / 部分唯一索引 / 普通索引（列序与谓词逐字不变）
m4  校验：information_schema 类型与精度（PG）/ fixed-6 TEXT 形状与词法序（SQLite）；FK/完整性检查
m5  写 ledger：descriptor version/name/checksum 追加（append-only；v1–v72 不变）
```

- **不可逆点**（`m2` 内，回滚只能靠事务或 m0 的转换前 snapshot）：0→NULL、微秒截断、非规范 TEXT fail closed。
- **owner 边界**：`core.persistence` 拥有 `schema_migrations.applied_at` 的转换，但 **Store runner 仍是该列的唯一写入者**（Root `D-011`）；v73 与 v74 的边界即 `schema_migrations`（v73）vs `system_data_reconcile` + auth 表（v74）。

## 7. 声明

- 本文件为冻结候选；**不**声称 `apps/` 下任何代码/迁移/测试已修改。A-029 点名「现行谓词仍为整数」在本文件中是 **old 列**的事实陈述，不是已修复声明。
- 本文件闭合 F-I-006 的关闭要求中「一张 exact old/new SQL 表 + 迁移顺序 + 负例/NULL/sentinel 测试 ID + voucher/monotonic 两族」「Root D-012/D-013 进入该表」四项；并按其要求以 Root/child 限定 D-012/D-013（F-I-018 的冻结载体部分）。是否接受由 independent 复审判定。
