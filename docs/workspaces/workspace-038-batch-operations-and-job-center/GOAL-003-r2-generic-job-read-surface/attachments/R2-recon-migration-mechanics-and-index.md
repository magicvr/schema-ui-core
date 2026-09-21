---
title: R2 侦察 · 迁移机制与索引决策（I-038-007 / O-3）
status: draft
created: 2026-09-19
updated: 2026-09-19
parent: GOAL-003-r2-generic-job-read-surface
version: 0.1.0
---

# R2 侦察 · 迁移机制与索引决策

> **性质**：只读侦察记录，为 `I-038-007`（管理列表索引决策 = R1 `D-001` §1.3 O-3）提供 file:line 证据。
> 编排器**独立核对**（非仅采信子代理报告）。

## 1. 迁移机制（已核实）

### 1.1 贡献结构与双方言

`kernel.MigrationContribution`（`apps/api/kernel/contribution.go:118-134`）：

| 字段 | 作用 |
|------|------|
| `ContributionIdentity{ModuleID, Key}` | 贡献身份；`Key` 在同一 ModuleID 内唯一 |
| `Version int` | **全局**版本号（跨模块单调） |
| `Name string` | 与 ledger 行比对的名字 |
| `Checksum string` | `kernel.MigrationChecksum(DDL, "<版本>:<名>:v1")` |
| `Apply func(Tx) error` | sqlite/规范体 |
| `ApplyPostgres func(Tx) error` | 可选；`nil` 表示 `Apply` 可移植 |
| `Tombstone bool` | 仅记录 ledger，无 `Apply` |

**双方言机制**（`contribution.go:124-128` 注释 + 代码先例）：`ApplyPostgres` 非 nil 时 postgres runner 改用它；**ledger checksum 在两种方言下都绑定 sqlite/规范历史**。存在自动化转换助手先例：`pgTimeDDL(stmts []string) []string`（`apps/api/modules/operationlog/migration/migration.go:252-258`），用正则 `pgTimeColRe` 把时间列改写为 `BIGINT NOT NULL`——**不是手工复制两份 DDL**。

### 1.2 版本与 ledger 校验（决定「能否加索引」）

`internal/store/migrate.go`：

- `validateApplied`（`:171-198`）要求 ledger 是编译 catalog 的**连续前缀**：
  - 首行必须 `version == 1`（`:179-181`）；
  - 相邻行必须**逐 1 递增**（`:183-185`）——「missing intermediate version」直接报错；
  - 每行 `version` 必须在 catalog 中已知（`:186-189`），且 `Name`（`:190-192`）与 `Checksum`（`:193-195`）必须与代码一致，否则 **checksum drift** 报错。
- `pendingMigrations`（`:200-216`）：凡 `version` 不在 ledger 中的 catalog 项即 pending。
- 因此**新增一个迁移版本必须占用下一个连续整数**（当前最大 = **71**，owner `core.operationlog` / `operation_log_digitaloffer_events`，`apps/api/modules/operationlog/migration/migration.go:483-486`）⇒ 新贡献应取 **72**。

### 1.3 「后续迁移给既有表加索引/改列」是**既有成熟模式**

同 ModuleID 多贡献的模块（`ContributionIdentity` 计数）：

| 模块 | 贡献数 | 说明 |
|------|--------|------|
| `core.operationlog` | 22 | 反复**重建** `operation_log` 表以扩展 event CHECK，每次带新版本 |
| `core.auth-session` | 15 | — |
| `admin.settings` | 6 | — |
| `core.persistence` | 5 | — |
| `admin.wallet` | 5 | 含 `wallet_ledger_deduct` 的 CHECK 重建 |
| `admin.notifications` | 3 | — |
| `admin.account` | 2 | `account_enable_state`(13) + `account_avatar_url`(35) —— **给既有表加列** |
| `admin.data-dictionary` | 2 | — |
| **`core.jobs`** | **1** | `async_jobs`(42) —— 迄今未再增贡献 |

**关键先例**：`admin.account` 的第二个贡献（`Version: 35`，`Key: account_avatar_url`，`apps/api/modules/account/migration/migration.go`）就是「同一 ModuleID 的**后续版本**给既有表加列」——证明**同一模块可以追加后续迁移**，不必新建模块。

**operationlog 的 22 个贡献**则证明「同一表反复重建 + 每次新版本」是仓内接受的模式（SQLite 不支持改 CHECK，故走重建）。

### 1.4 冻结断言的实际语义

`internal/store/migrate_test.go:692-693` 冻结的是 **`{ModuleID, Name, Checksum}` 三元组列表**：

```go
// VP-012 R4: migration-only core.jobs durable state machine.
{"core.jobs", "async_jobs", "55e1d3f88de080bd0b6015841e76f1ce32604444619d180a3b228123f99dec68"},
```

⇒ **已应用的 `async_jobs`(42) 行的 checksum 不可变**（改它就是 checksum drift，会让既有库拒绝启动）。但**新增**一条 `core.jobs` 贡献（新 Version/Key/Name/Checksum）是**追加**，不违反该冻结——只需把新三元组加入该测试的期望列表。

`modules/jobs/migration/migration_test.go:15` 断言 `len(descriptors) != 1` 与 `Version != 42` / `Name != "async_jobs"` —— **该测试必须随新贡献更新**（`len` 变 2）。

## 2. 索引决策（O-3）分析

### 2.1 既有三索引的前导列

| 索引 | 列序 | 可服务 |
|------|------|--------|
| `idx_jobs_runnable` | `(status, cancel_requested, lease_expires_at, created_at)` | `status` 等值 + 后续列等值/范围；**不含** `updated_at` |
| `idx_jobs_actor` | `(actor_id, kind, updated_at DESC)` | `actor_id`（±`kind`）等值 + `updated_at` 排序 |
| `idx_jobs_expiry` | `(status, expires_at)` | 过期扫描 |

### 2.2 候选管理列表查询 × 索引覆盖

| 查询形状 | 既有索引覆盖 | 结论 |
|---------|-------------|------|
| (a) 无过滤 `ORDER BY created_at DESC` | ❌ 无 `created_at` 前导索引 | **需新索引** |
| (b) `WHERE status=? ORDER BY created_at DESC` | ⚠️ `idx_jobs_runnable` 前导列是 `status`，但排序列 `created_at` 是**第 4 列**且中间有 `cancel_requested`/`lease_expires_at` | 部分；排序仍需 sort |
| (c) `WHERE kind=? ORDER BY created_at DESC` | ❌ `kind` 仅在 `idx_jobs_actor` 第 2 列（前导是 `actor_id`） | **需新索引** |
| (d) `WHERE actor_id=? ORDER BY created_at DESC` | ⚠️ `idx_jobs_actor` 覆盖 `actor_id` 等值，但排序列是 `updated_at DESC` 而非 `created_at` | 等值可走索引，排序不匹配 |
| (e) `created_at` 范围 + `ORDER BY created_at DESC` | ❌ | **需新索引** |

`COUNT(*)` 同 WHERE：SQLite 可用覆盖索引做 index-only count（若索引包含 WHERE 列）。

### 2.3 建议（待用户裁决）

若默认列表为「全量 + 时间倒序」（形状 a，最贴近 `admin.activity` 的 `idx_operation_log_created_at ON operation_log(created_at DESC)` 先例）：

```sql
-- sqlite / 规范体
CREATE INDEX idx_jobs_created_at ON jobs(created_at DESC)
-- postgres（时间列为 BIGINT，DDL 文本相同）
CREATE INDEX idx_jobs_created_at ON jobs(created_at DESC)
```

若默认列表需要 `kind` / `status` 过滤，建议复合：

```sql
CREATE INDEX idx_jobs_created_at ON jobs(created_at DESC)
CREATE INDEX idx_jobs_kind_created ON jobs(kind, created_at DESC)
CREATE INDEX idx_jobs_status_created ON jobs(status, created_at DESC)
```

**实施代价**：新贡献 `Version: 72`、`Key: "jobs_management_indexes"`；`modules/jobs/migration/migration_test.go:15` 的 `len(descriptors)` 期望 1 → 2；`internal/store/migrate_test.go` 冻结列表追加新三元组。**`async_jobs`(42) 行本身不动。**

## 3. EXPLAIN 可用性

仓内**未发现** `EXPLAIN QUERY PLAN` 的使用（grep `EXPLAIN` 在 `apps/api` 无命中）。⇒ 索引覆盖结论为**静态列序推理**，不可机器验证；R2 若需实测，须新增工具或测试。

## 4. 待确认 / 未知

| # | 项 | 影响 |
|---|----|------|
| U-01 | 列表默认排序是否应为 `created_at DESC` 还是 `updated_at DESC`（后者与 `idx_jobs_actor` 的排序一致，但仅限 actor 作用域） | 决定索引形状 |
| U-02 | 是否需要 kind / status 复合索引，取决于 UI 是否默认过滤 | 决定索引数量 |
| U-03 | 管理列表是否需要「结果过期」的惰性语义（R1 矩阵 R-4） | 影响查询是否要 join/更新 |
| U-04 | Postgres 侧是否真的用 `ApplyPostgres`（`pgTimeDDL` 仅改时间列类型；纯索引 DDL 两方言文本相同，可能无需 `ApplyPostgres`） | 决定是否省略 `ApplyPostgres` |
