---
id: r1-c2-per-table-pg-ddl-v1.0-fc
doc_type: design-attachment
title: R1 C2 逐表 PG 显式 DDL v1.0（冻结候选）
status: freeze-candidate
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-002-r1-contract-and-denominator-freeze
version: 1.0.0
---

# R1 C2 逐表 PG 显式 DDL v1.0（冻结候选）

> **状态：`freeze-candidate`。** 本文件是 `D-019` §6 要求的 **PostgreSQL 侧显式 DDL**：v73–v87 每个 descriptor 的 `ApplyPostgres` 语句序列。
>
> **为什么必须显式书写**：`operationlog/migration/migration.go:250` 的 `pgTimeColRe`（`(?m)^(\s*)(created_at|archived_at)\s+INTEGER NOT NULL`）与 `pgTimeDDL`（`:252-258`）从 **SQLite 字面**派生 PG 变体——一旦 SQLite 字面改成 `TEXT NOT NULL`，该正则**静默 no-op**，PG 侧会保留陈旧的 `BIGINT` 映射且**不报错**。A-032 §E 确认属实，并判定显式 PG DDL **属 R1 冻结交付**（不拖到 R2）。
>
> **与 `D-017` 的关系**：PG DDL **不进** `MigrationChecksum`（单 checksum 只哈希 SQLite `DDL` 切片 + `transform_id`）；PG 侧由测试断言覆盖。

## 0. 全局约定

### 0.1 两条 PG 转换表达式（与 SQLite 侧同族、同语义）

```sql
-- 秒族（legacy BIGINT 整数秒 → timestamptz(6)）
date_trunc('microseconds', to_timestamp(<col>::double precision))

-- 毫秒族（legacy BIGINT 整数毫秒 → timestamptz(6)）—— **整数拆分式**（见下方更正）
TIMESTAMPTZ 'epoch'
  + ((<col> - CASE WHEN <col> >= 0 THEN 0 ELSE 999 END) / 1000) * INTERVAL '1 second'
  + ((<col> % 1000 + 1000) % 1000) * INTERVAL '1 millisecond'
```

> **⚠️ 2026-09-20 实测更正（PG 临时容器，`postgres:16`）——原毫秒式有精度缺陷。**
>
> 原设计式为 `date_trunc('microseconds', TIMESTAMPTZ 'epoch' + <col> * INTERVAL '1 millisecond')`。实测在**大数值**上产生**非零误差**：
>
> | 输入（ms） | 原式输出 | 期望 |
> |-----------|---------|------|
> | `253402300799999`（公元 9999 年） | `9999-12-31T23:59:59.999008` | `…59.999000` |
>
> 误差 **+8 微秒**。根因：`BIGINT * INTERVAL` 路径内部经**浮点**转换，在 ~2.5×10¹⁴ 量级上丢失微秒精度。**`::numeric` 强制转换无效**（实测同为 `.999008`）。
>
> **正确的整数拆分式**（无浮点参与）实测结果（`postgres:16`，`SET TIME ZONE 'UTC'`）：
>
> | 输入（ms） | 输出 |
> |-----------|------|
> | `1758320000123` | `2025-09-19T22:13:20.123000` |
> | `1758320000999` | `2025-09-19T22:13:20.999000`（**无进位**） |
> | `-1` | `1969-12-31T23:59:59.999000` |
> | `-999` | `1969-12-31T23:59:59.001000` |
> | `-1000` | `1969-12-31T23:59:59.000000` |
> | `-1001` | `1969-12-31T23:59:58.999000` |
> | `-86400000` | `1969-12-31T00:00:00.000000` |
> | `-1758320000123` | `1914-04-14T01:46:39.877000` |
> | `253402300799999` | `9999-12-31T23:59:59.999000` ✅ |
>
> 该式与 SQLite 侧同构（同样 floor 秒 + 归一化余数）：PG 整数除法同样**向零截断**（实测 `-1/1000 = 0`、`-1%1000 = -1`），故负值需 `-999` 修正与 `(x%1000+1000)%1000` 归一化。
>
> **秒族无需更改**：`to_timestamp(value::double precision)` 在 `253402300799` 与 `-1` 上实测精确。
>
> `D-019` §6「PG DDL 必须显式书写」的要求因此**更有必要**——若沿用 `pgTimeColRe` 派生（它只把 `INTEGER` 换成 `BIGINT`），该精度缺陷不会被任何断言发现。

与 `r1-c2-per-column-conversion-contract-v1.0-fc.md` §1 的 E1/E2 同一族；秒族保留 `to_timestamp(double)`（用户 2026-09-20 裁决 B，Root `D-015` 已按此修订）。**毫秒族的 E2 定义已按上述实测更正为整数拆分式。**

### 0.2 每个时间列的 PG 语句骨架

对每个 `NOT NULL` 时间列（**可空列**省去 DROP/SET NOT NULL 两行）：

```sql
ALTER TABLE "<table>" ALTER COLUMN "<col>" DROP NOT NULL;
ALTER TABLE "<table>" ALTER COLUMN "<col>" TYPE timestamptz(6)
  USING <对应表达式>;
ALTER TABLE "<table>" ALTER COLUMN "<col>" SET NOT NULL;
```

对 **D0 列**（`NOT NULL DEFAULT 0`，0 = 缺失）：

```sql
ALTER TABLE "<table>" ALTER COLUMN "<col>" DROP DEFAULT;
ALTER TABLE "<table>" ALTER COLUMN "<col>" DROP NOT NULL;
ALTER TABLE "<table>" ALTER COLUMN "<col>" TYPE timestamptz(6)
  USING (CASE WHEN "<col>" = 0 THEN NULL ELSE <对应表达式> END);
-- 不恢复 NOT NULL，也不恢复 DEFAULT
```

对 **voucher 列**（可空且 legacy `0` 亦为缺失）：

```sql
ALTER TABLE "<table>" ALTER COLUMN "<col>" TYPE timestamptz(6)
  USING (CASE WHEN "<col>" IS NULL OR "<col>" = 0 THEN NULL ELSE <对应表达式> END);
```

对**可空列**（NULL 保持 NULL）：

```sql
ALTER TABLE "<table>" ALTER COLUMN "<col>" TYPE timestamptz(6)
  USING (CASE WHEN "<col>" IS NULL THEN NULL ELSE <对应表达式> END);
```

### 0.3 通用约束

- **禁止** `::timestamptz(6)` 直接 typmod cast（默认 **round**，不等于向零截断）。
- 约束/索引/默认在类型转换后按原定义恢复；PG 无需重建表（与 SQLite 不同），因此 **FK 父表在 PG 侧没有 F-5 问题**。
- 重建后校验（进 `m4`）：`information_schema.columns.data_type = 'timestamp with time zone'` 且 `datetime_precision = 6`，对全部目标列成立；`pg_constraint` 中 CHECK/UNIQUE/PK 与迁移前一致。

## 1. v73 `core.persistence`（毫秒族）

```sql
-- schema_migrations.applied_at（NN，秒族；★ 见 §7 的写入路径联动）
ALTER TABLE "schema_migrations" ALTER COLUMN "applied_at" DROP NOT NULL;
ALTER TABLE "schema_migrations" ALTER COLUMN "applied_at" TYPE timestamptz(6)
  USING date_trunc('microseconds', to_timestamp("applied_at"::double precision));
ALTER TABLE "schema_migrations" ALTER COLUMN "applied_at" SET NOT NULL;

-- mail_outbox.created_at（NN，毫秒族）—— USING 用 §0.1 的 E2 整数拆分式（下方为展开）
ALTER TABLE "mail_outbox" ALTER COLUMN "created_at" DROP NOT NULL;
ALTER TABLE "mail_outbox" ALTER COLUMN "created_at" TYPE timestamptz(6)
  USING (TIMESTAMPTZ 'epoch'
         + (("created_at" - CASE WHEN "created_at" >= 0 THEN 0 ELSE 999 END) / 1000) * INTERVAL '1 second'
         + ((("created_at" % 1000) + 1000) % 1000) * INTERVAL '1 millisecond');
ALTER TABLE "mail_outbox" ALTER COLUMN "created_at" SET NOT NULL;

-- mail_config.updated_at（D0）—— ELSE 主体用 §0.1 的 E2 整数拆分式
ALTER TABLE "mail_config" ALTER COLUMN "updated_at" DROP DEFAULT;
ALTER TABLE "mail_config" ALTER COLUMN "updated_at" DROP NOT NULL;
ALTER TABLE "mail_config" ALTER COLUMN "updated_at" TYPE timestamptz(6)
  USING (CASE WHEN "updated_at" = 0 THEN NULL
              ELSE TIMESTAMPTZ 'epoch'
                   + (("updated_at" - CASE WHEN "updated_at" >= 0 THEN 0 ELSE 999 END) / 1000) * INTERVAL '1 second'
                   + ((("updated_at" % 1000) + 1000) % 1000) * INTERVAL '1 millisecond' END);
```

## 2. v74 `core.auth-session`（秒族；31 列）

> PG 侧**不需要** F-5：`ALTER COLUMN TYPE` 不重命名表，故子表 FK 不受影响。

```sql
-- system_data_reconcile.applied_at（NN）
ALTER TABLE "system_data_reconcile" ALTER COLUMN "applied_at" DROP NOT NULL;
ALTER TABLE "system_data_reconcile" ALTER COLUMN "applied_at" TYPE timestamptz(6)
  USING date_trunc('microseconds', to_timestamp("applied_at"::double precision));
ALTER TABLE "system_data_reconcile" ALTER COLUMN "applied_at" SET NOT NULL;

-- ── users ─────────────────────────────────────────────────────────────
-- created_at / updated_at（NN）
ALTER TABLE "users" ALTER COLUMN "created_at" DROP NOT NULL;
ALTER TABLE "users" ALTER COLUMN "created_at" TYPE timestamptz(6)
  USING date_trunc('microseconds', to_timestamp("created_at"::double precision));
ALTER TABLE "users" ALTER COLUMN "created_at" SET NOT NULL;
-- updated_at 同形
ALTER TABLE "users" ALTER COLUMN "updated_at" DROP NOT NULL;
ALTER TABLE "users" ALTER COLUMN "updated_at" TYPE timestamptz(6)
  USING date_trunc('microseconds', to_timestamp("updated_at"::double precision));
ALTER TABLE "users" ALTER COLUMN "updated_at" SET NOT NULL;
-- locked_until / last_login_failure_at（D0）
ALTER TABLE "users" ALTER COLUMN "locked_until" DROP DEFAULT;
ALTER TABLE "users" ALTER COLUMN "locked_until" DROP NOT NULL;
ALTER TABLE "users" ALTER COLUMN "locked_until" TYPE timestamptz(6)
  USING (CASE WHEN "locked_until" = 0 THEN NULL
              ELSE date_trunc('microseconds', to_timestamp("locked_until"::double precision)) END);
ALTER TABLE "users" ALTER COLUMN "last_login_failure_at" DROP DEFAULT;
ALTER TABLE "users" ALTER COLUMN "last_login_failure_at" DROP NOT NULL;
ALTER TABLE "users" ALTER COLUMN "last_login_failure_at" TYPE timestamptz(6)
  USING (CASE WHEN "last_login_failure_at" = 0 THEN NULL
              ELSE date_trunc('microseconds', to_timestamp("last_login_failure_at"::double precision)) END);

-- ── 其余 NN 秒列（同形：DROP NOT NULL → TYPE … USING date_trunc(to_timestamp(…)) → SET NOT NULL）──
-- refresh_tokens: expires_at, created_at
-- roles:          created_at, updated_at
-- permissions:    created_at, updated_at
-- menu_items:     created_at, updated_at
-- email_verification_challenges:          expires_at, sent_at
-- password_recovery_challenges:           expires_at, sent_at
-- login_failures: updated_at
-- user_password_history:                  created_at
-- user_invites:   expires_at, last_sent_at, created_at
-- service_credentials:                    expires_at, created_at, updated_at

-- ── 可空秒列（NULL 包裹）──
-- refresh_tokens.revoked_at
ALTER TABLE "refresh_tokens" ALTER COLUMN "revoked_at" TYPE timestamptz(6)
  USING (CASE WHEN "revoked_at" IS NULL THEN NULL
              ELSE date_trunc('microseconds', to_timestamp("revoked_at"::double precision)) END);
-- user_invites.consumed_at / revoked_at；service_credentials.revoked_at / last_used_at 同形

-- ── login_failures.locked_until（D0）──
ALTER TABLE "login_failures" ALTER COLUMN "locked_until" DROP DEFAULT;
ALTER TABLE "login_failures" ALTER COLUMN "locked_until" DROP NOT NULL;
ALTER TABLE "login_failures" ALTER COLUMN "locked_until" TYPE timestamptz(6)
  USING (CASE WHEN "locked_until" = 0 THEN NULL
              ELSE date_trunc('microseconds', to_timestamp("locked_until"::double precision)) END);
```

> **注意**：authsession 的 `PGDDL` 变体里 `users` 的时间列在若干历史版本被判为 `BIGINT`（如 `migration.go:220,435`）。转换前须以 **live `information_schema` / `pg_attribute`** 核对实际类型，不得假定。

## 3. v75 `core.operationlog`（毫秒族）——**本 descriptor 禁用 `pgTimeDDL`**

```sql
-- operation_log.created_at（NN，毫秒族）—— USING 用 §0.1 的 E2 整数拆分式
ALTER TABLE "operation_log" ALTER COLUMN "created_at" DROP NOT NULL;
ALTER TABLE "operation_log" ALTER COLUMN "created_at" TYPE timestamptz(6)
  USING (TIMESTAMPTZ 'epoch'
         + (("created_at" - CASE WHEN "created_at" >= 0 THEN 0 ELSE 999 END) / 1000) * INTERVAL '1 second'
         + ((("created_at" % 1000) + 1000) % 1000) * INTERVAL '1 millisecond');
ALTER TABLE "operation_log" ALTER COLUMN "created_at" SET NOT NULL;

-- operation_log_archive.created_at / archived_at（NN，毫秒族，同形）
```

- **不重建表**：PG 侧无需处理 `operation_log_correlation` / `operation_log_session` 的 FK（这是 SQLite 专属问题）。
- **禁止** `pgRebuild`（`operationlog/migration/migration.go:296-301`）——它把 `pgTimeDDL` 的派生结果送进 rename 重建，在 SQLite 字面改 `TEXT` 后派生静默失效。
- 索引不变：`idx_operation_log_created_at`、`idx_operation_log_archive_created_at` 在 PG 侧随列类型自动适配，**无需重建**。

## 4. v76 `core.jobs` / v77 `core.data-dictionary` / **v78 `admin.data-permission`** / v79 captcha / v82 recycle / v84 settings

### 4.1 v78 `admin.data-permission`（**秒族 NN ×2**，响应 A-036 **F-I-026**）

```sql
-- data_scope_policies.updated_at（NN，秒族）
ALTER TABLE "data_scope_policies" ALTER COLUMN "updated_at" DROP NOT NULL;
ALTER TABLE "data_scope_policies" ALTER COLUMN "updated_at" TYPE timestamptz(6)
  USING date_trunc('microseconds', to_timestamp("updated_at"::double precision));
ALTER TABLE "data_scope_policies" ALTER COLUMN "updated_at" SET NOT NULL;

-- user_data_scopes.updated_at（NN，秒族）
ALTER TABLE "user_data_scopes" ALTER COLUMN "updated_at" DROP NOT NULL;
ALTER TABLE "user_data_scopes" ALTER COLUMN "updated_at" TYPE timestamptz(6)
  USING date_trunc('microseconds', to_timestamp("updated_at"::double precision));
ALTER TABLE "user_data_scopes" ALTER COLUMN "updated_at" SET NOT NULL;
```

- 对应 ledger v78（`#50` = `data_scope_policies.updated_at`、`#51` = `user_data_scopes.updated_at`），均 `S-NN`。
- **不**属于 F-5：PG 侧本无 F-5 需求（§0.3）；这两表在 SQLite 侧亦无 FK 子表（SQLite 附件 §2.6 实测 REFERENCES count = 0）。
- PG 侧该 descriptor **无需**重建表、无需处理 FK；两表在 PG 下无显式索引（仅 PK）。

### 4.2 其余秒族 NN 列（骨架同 §0.2）

```sql
-- jobs（毫秒族！lease_expires_at / created_at / updated_at / finished_at / expires_at）
-- 前三个按 NN 或可空性：created_at/updated_at NN；lease_expires_at/finished_at/expires_at 可空（NULL 包裹）
-- 六态 CHECK 依赖列类型，PG 在 ALTER TYPE 后自动重校验；无需重建

-- dict_types.created_at / updated_at（NN，秒族）
-- dict_entries.created_at / updated_at（NN，秒族）
-- captcha_challenges.expires_at / created_at（NN，秒族）
-- captcha_config.created_at / updated_at（NN，秒族）
-- recycle_items.deleted_at（NN）/ restored_at（可空，秒族）
-- site_settings.updated_at（NN，秒族）
-- 全部按 §0.2 骨架逐列执行
```

> v77/v78/v82/v83/v86/v87 的**部分唯一索引与部分 CHECK 在 PG 侧不因列类型变化而改变定义**，无需重建；SQLite 侧因表重建才需逐字重建。

## 5. v80 `mfa` / v81 `notifications` / v83 `scheduled-tasks` / v85 `wallet` / v86 `telegram` / v87 `digital-offer`（秒族）

```sql
-- mfa: user_mfa.created_at / updated_at（NN）；mfa_proofs.expires_at / created_at（NN）
-- notifications: created_at（NN）；read_at（可空）
-- scheduled_tasks: created_at / updated_at（NN）；task_runs.started_at（NN）/ created_at（NN）/ finished_at（可空）
-- wallet: wallet_accounts.created_at / updated_at（NN）；
--         wallet_ledger_entries.created_at（NN）；wallet_reconciliation_runs.created_at（NN）；
--         subjects.created_at（NN）；voucher_batches.created_at / updated_at（NN）；
--         vouchers.created_at / updated_at（NN）+ expires_at / redeemed_at（★ voucher 规则：IS NULL OR = 0）
-- telegram: telegram_config.updated_at（★ D0）；telegram_sessions.last_message_at / created_at / updated_at（NN）；
--           telegram_inbound_messages.received_at（NN）；telegram_outbound_messages.created_at / updated_at（NN）
-- digital-offer: digital_offers.created_at / updated_at（NN）；digital_purchases.created_at（NN）；
--                digital_entitlements.created_at / updated_at（NN）+ expires_at（可空，普通 NULL 包裹）
```

```sql
-- vouchers 两个 voucher 列（示例）
ALTER TABLE "vouchers" ALTER COLUMN "expires_at" TYPE timestamptz(6)
  USING (CASE WHEN "expires_at" IS NULL OR "expires_at" = 0 THEN NULL
              ELSE date_trunc('microseconds', to_timestamp("expires_at"::double precision)) END);
-- telegram_config.updated_at（D0，示例）
ALTER TABLE "telegram_config" ALTER COLUMN "updated_at" DROP DEFAULT;
ALTER TABLE "telegram_config" ALTER COLUMN "updated_at" DROP NOT NULL;
ALTER TABLE "telegram_config" ALTER COLUMN "updated_at" TYPE timestamptz(6)
  USING (CASE WHEN "updated_at" = 0 THEN NULL
              ELSE date_trunc('microseconds', to_timestamp("updated_at"::double precision)) END);
```

## 6. PG 侧仍需联动的**非 DDL** 项

| 项 | 位置 | 目标 |
|----|------|------|
| `applied_at` 写入 | `internal/store/postgres.go:165-167` | 绑定 `time.Time`（`Truncate(µs)`），而非 `Unix()` |
| `postgresLedgerDDL` | `internal/store/identity.go:65` | `applied_at timestamptz(6) NOT NULL` |
| `settings` seeder | `modules/settings/migration/migration.go:52-55` | 以目标格式写入 `updated_at` |
| 金额列断言 | `internal/store/postgres_test.go:291-307` | 时间列改判 `timestamp with time zone`（精度 6）；**金额列 `wallet_accounts.balance_total` / `wallet_ledger_entries.amount_delta` 保持 `bigint`** |
| leftover 列名 | 同文件 `:312-316` | 补齐 21 名 leftover 清单中缺失的七列 |

## 7. v73 的 `applied_at` 双向边界（与 SQLite 附件 §6 一致）

- `identity.go:58` / `:65` 两个 restore 字面：改。
- `identity.go:317,319-321`（`stampCatalog`）、`migrate.go:121-123`（`applyMigration`）、`postgres.go:165-167`（`applyMigrationPG`）：改。
- `authsession/migration/migration.go:18-23` 的 `schemaMigrationsDDL`：**禁止改**（属 v1；其 PG 变体 `:332` 同理）。
- **不影响 v1 checksum**：`identity.go` 的两个 restore 字面不在 `0001:r2-baseline` 的 `r2BaselineDDL` 哈希输入内。

## 8. 声明

- 本文件 `status: freeze-candidate`；**不**声称任何 PG DDL/迁移/测试已实施；`apps/` 本轮未修改。
- §2/§4/§5 对**同形重复语句**采用「骨架 + 列清单」写法（§0.2/§0.3 已给出唯一形态），逐列展开见 `r1-c2-per-column-conversion-contract-v1.0-fc.md` §2 的 90 行合同；若 independent 复审要求逐列逐字展开，下一轮补齐。
- 是否构成 **F-I-002**（PG 侧）的合法闭合，由 independent 复审判定。
- 引用 `D-0NN` 一律限定 Root/child。
