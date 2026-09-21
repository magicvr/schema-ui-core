---
id: r1-c2-descriptor-ledger-v1.0-fc
doc_type: design-attachment
title: R1 C2 v73–v87 conversion descriptor 台账 v1.0（冻结候选）
status: freeze-candidate
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-002-r1-contract-and-denominator-freeze
version: 1.0.0
---

# R1 C2 v73–v87 conversion descriptor 台账 v1.0（冻结候选）

> **状态：`freeze-candidate`。** 本文件闭合 A-029（现行 independent head，`conditional`，open required = 5）**F-I-005** 关闭要求中可**在 R1 内确定**的子项：**唯一表范围**（消掉 v74 ledger 歧义）、**descriptor 名**、**transform ID 约定**、**`MigrationChecksum` 计算输入（有序语句清单 + 规范化算法）**。**不**闭合「已迁移的 checksum 值」与「可执行测试改写」——`MigrationChecksum` 的哈希值只能在 R2 落码后计算记录（见 §5）。

## 1. 台账核心（v73–v87，共 15 个 descriptor）

`version` 由 D-014 接受的未发布 baseline 固定；`name` / `transform_id` 本次按仓库既有约定（`kernel.MigrationChecksum(stmts, transformID)`；既有形如 `"0051:mail-outbox:v1"`）定为 v73–v87 的确定性值，**待 C2 复审接受**。

| v | ModuleID | descriptor `Name` | `transform_id` | 表范围（唯一） | 列（inventory `#N`） |
|--:|----------|-------------------|----------------|----------------|----------------------|
| 73 | `core.persistence` | `vp040_temporal_core_persistence` | `0073:vp040-temporal-core-persistence:v1` | `schema_migrations`、`mail_outbox`、`mail_config`；断言 retired `records` 不存在 | #1, #33, #34 |
| 74 | `core.auth-session` | `vp040_temporal_authsession` | `0074:vp040-temporal-authsession:v1` | `system_data_reconcile`、`users`、`refresh_tokens`、`roles`、`permissions`、`menu_items`、`email_verification_challenges`、`password_recovery_challenges`、`login_failures`、`user_password_history`、`user_invites`、`service_credentials`（**不含 `schema_migrations`，归 v73**） | #2–#32（31 列） |
| 75 | `core.operationlog` | `vp040_temporal_operationlog` | `0075:vp040-temporal-operationlog:v1` | `operation_log`、`operation_log_archive` | #35–#37 |
| 76 | `core.jobs` | `vp040_temporal_jobs` | `0076:vp040-temporal-jobs:v1` | `jobs` | #41–#45 |
| 77 | `admin.data-dictionary` | `vp040_temporal_dictionary` | `0077:vp040-temporal-dictionary:v1` | `dict_types`、`dict_entries` | #46–#49 |
| 78 | `admin.data-permission` | `vp040_temporal_data_permission` | `0078:vp040-temporal-data-permission:v1` | `data_scope_policies`、`user_data_scopes` | #50–#51 |
| 79 | `admin.login-captcha` | `vp040_temporal_captcha` | `0079:vp040-temporal-captcha:v1` | `captcha_challenges`、`captcha_config` | #52–#55 |
| 80 | `admin.mfa` | `vp040_temporal_mfa` | `0080:vp040-temporal-mfa:v1` | `user_mfa`、`mfa_proofs`（`last_used_step` 排除） | #63–#66 |
| 81 | `admin.notifications` | `vp040_temporal_notifications` | `0081:vp040-temporal-notifications:v1` | `notifications` | #38–#39 |
| 82 | `admin.recycle-bin` | `vp040_temporal_recycle` | `0082:vp040-temporal-recycle:v1` | `recycle_items` | #56–#57 |
| 83 | `admin.scheduled-tasks` | `vp040_temporal_scheduled_tasks` | `0083:vp040-temporal-scheduled-tasks:v1` | `scheduled_tasks`、`task_runs` | #58–#62 |
| 84 | `admin.settings` | `vp040_temporal_settings` | `0084:vp040-temporal-settings:v1` | `site_settings` | #40 |
| 85 | `admin.wallet` | `vp040_temporal_wallet` | `0085:vp040-temporal-wallet:v1` | `wallet_accounts`、`wallet_ledger_entries`、`wallet_reconciliation_runs`、`subjects`、`vouchers`、`voucher_batches` | #67–#77 |
| 86 | `admin.channel.telegram` | `vp040_temporal_telegram` | `0086:vp040-temporal-telegram:v1` | `telegram_config`、`telegram_sessions`、`telegram_inbound_messages`、`telegram_outbound_messages` | #78–#84 |
| 87 | `biz.digital-offer` | `vp040_temporal_digital_offer` | `0087:vp040-temporal-digital-offer:v1` | `digital_offers`、`digital_purchases`、`digital_entitlements` | #85–#90 |

**表范围唯一性自检**：15 个 descriptor 的表集合两两不交；`schema_migrations` 仅出现在 v73、`system_data_reconcile` 仅出现在 v74。列分配合计 **3 + 31 + 3 + 5 + 4 + 2 + 4 + 4 + 2 + 2 + 5 + 1 + 11 + 7 + 6 = 90**，与 inventory v0.3 逐列一一对应。

### 1.1 v74 范围消歧（本次收口）

| 载体 | 原文 | 判定 |
|------|------|------|
| `r1-v73-owner-allocation-draft-v0.1.md` L19 | v74 表范围写「… `system_data_reconcile` …」（「schema/system ledger」措辞的来源） | **以列分配为准**：#2 `system_data_reconcile.applied_at` 归 v74，#1 `schema_migrations.applied_at` 归 v73 |
| `r1-c2-owner-migration-spec-v0.1.md` L28 | v74 assigned scope 写「ledger/reconcile」 | **同上**；「ledger」指 `system_data_reconcile`，不是 `schema_migrations` |

结论：**v74 表范围 = `system_data_reconcile` + auth 表；`schema_migrations` 归 v73**。两份既有载体的措辞按本表统一；本表生效后「v74 占用 v73 或 v85 ledger」的两种读法均被排除（v85 wallet ledger 为 `wallet_ledger_entries` / `wallet_reconciliation_runs`，与 v74 无交集）。

## 2. `MigrationChecksum` 计算输入（逐字算法）

现行实现（`apps/api/kernel/persistence.go:14-17`）：

```go
func MigrationChecksum(stmts []string, transformID string) string {
    input := normalizeSQL(strings.Join(stmts, "\n")) + "\n" + transformID
    sum := sha256.Sum256([]byte(input))
    return hex.EncodeToString(sum[:])
}
// normalizeSQL: 逐行 TrimSpace，丢弃空行，以 "\n" 重新连接
```

因此每个 descriptor 的 checksum 由 **(a) 有序语句切片**、**(b) `transform_id`** 唯一决定。**(a) 的顺序即 `migration.go` 中 `DDL` 切片的字面顺序**，按 `r1-c2-predicate-exact-sql-v1.0-fc.md` §6 的 `m0 → m5` 展开：

| 序位 | 语句来源 |
|------|----------|
| `m0` | preflight 计数语句（每个 sentinel 列一条 SELECT 计数；无 sentinel 的 descriptor 省略） |
| `m1` | 约束放开（`ALTER TABLE … DROP DEFAULT` / `DROP NOT NULL` / `DROP CONSTRAINT` / `DROP INDEX`） |
| `m2` | 逐列转换：PG `ALTER TABLE … ALTER COLUMN … TYPE timestamptz(6) USING …`（E1/E2/E3，逐列一条，按 `#N` 升序）；SQLite `CREATE TABLE <t>_new …` + `INSERT INTO <t>_new SELECT …` + `DROP TABLE <t>` + `ALTER TABLE <t>_new RENAME TO <t>` |
| `m3` | 约束/索引重建（`CREATE INDEX` / `CREATE UNIQUE INDEX` / 表级 CHECK；列序与谓词逐字见 predicate 表） |
| `m4` | 校验语句（`information_schema` 类型断言 / 完整性与 FK 检查） |
| `m5` | descriptor 元数据落盘（version/name/checksum 追加，append-only） |

**双方言约定（用户 2026-09-20 P-004 裁决 = 选项 A，见 child `01-decision/D-017-v73-checksum-convention.md`）**：沿用 v1–v72 现行约定——每个 descriptor **只有一个** checksum，`stmts` = **SQLite canonical DDL 切片**（`m0 → m5` 顺序），**PG `ApplyPostgres` 变体不进哈希**，`transform_id` **不加** `:sqlite`/`:pg` 后缀。PG 侧差异由 `information_schema` 类型/精度测试断言承担，不由 checksum 承担。此项**已选定**，不再是开放二选一（原 §5.3 关闭）。

## 3. 冻结后仍不得改动（append-only 边界）

- v1–v72 的 version / name / checksum / canonical SQL **不可变**；`migrate_test.go:124`（`len(applied) != 72`）、`:765`（`len(catalog) != len(want)`）与逐条冻结哈希表只允许**追加** v73–v87 行。
- `apps/api/internal/store/postgres_test.go:291-307` 的 PG `bigint` 硬断言必须**拆开**：时间列改判 `timestamp with time zone`（精度 6），**金额列 `wallet_accounts.balance_total` / `wallet_ledger_entries.amount_delta` 保持 `bigint`**（A-027 F-I-005 点名项）。
- leftover 列名表（21 名，`r1-v73-owner-allocation-draft-v0.1.md` L37–L38）逐名进入 PG 断言集合；`postgres_test.go:312-316` 缺失的七列必须补齐。

## 4. 与 R2 的接口

R2 落码时的**唯一**允许输入：本台账的 `name` / `transform_id` / 表范围 / 列分配 + `r1-c2-per-column-conversion-contract-v1.0-fc.md` 的 E1/E2/E3 + `r1-c2-predicate-exact-sql-v1.0-fc.md` 的谓词与 `m0–m5` 顺序。任何偏离须回到 C2 决策，禁止在 R2 内静默改写。

## 5. 本候选**未**闭合的项（诚实边界）

1. **已迁移的 checksum 值**：`MigrationChecksum` 只能对**已存在的语句切片**求值；v73–v87 的语句切片属 R2 落码产物。本文件记录的是**计算输入的结构与算法**，不是哈希值本身。A-029 F-I-005 的「已记录的 canonical SQL / `MigrationChecksum`」子项在 R2 首次落码并写入 ledger 后才可闭合。
2. **可执行测试改写**：`migrate_test.go` / `postgres_test.go` / `operations_test.go` / `restart_test.go` 的 v73+ 断言与金额/时间拆分为 R2 产物；本文件只冻结「改成什么」与「哪些不可改」。
3. ~~**双方言 checksum 约定二选一**（§2 末）需 independent 明确选择。~~ **已由用户 2026-09-20 P-004 裁决为选项 A**（见 §2 与 child `D-017`）；本项不再是开放项。
4. ~~**唯一表范围** … **接受与否由 independent 复审判定**。~~ **已获独立接受**：唯一表范围（v74 消歧）经 **A-034** 接受（并独立复算列号不相交），**descriptor 名 / transform ID** 经 **A-036 §C** 接受（15/15 载体）。本项不再是开放项。
5. ~~**`r1-c2-owner-migration-spec-v0.1.md` L28 的 v74「ledger/reconcile」措辞同轮已改写** … **是否接受由 independent 复审判定**。~~ **已获独立接受**：**A-036 §B/§C** 确认 v74 与已接受 allocation 同文。本项不再是开放项。

> **本文件剩余开放项仅 1、2 两项**（已迁移哈希值 / 可执行测试改写），二者均由用户 2026-09-20 P-004 书面裁定为 **`accepted-residual` 移交 R2**（child `D-021`：范围穷举三项 + 复审触发 + 失效条件）。**不得读作哈希已验证。**

## 6. 声明

- 本文件 `status: freeze-candidate`；不改变 `status` / `progress` / goal-tree；**不**声称任何 descriptor 已创建、任何 checksum 已记录、任何测试已改写。
- 引用 D-012 一律限定 Root/child（Root `D-012` = voucher 异常值政策；child `D-012` = v73 allocation / 负瞬间承接）。
