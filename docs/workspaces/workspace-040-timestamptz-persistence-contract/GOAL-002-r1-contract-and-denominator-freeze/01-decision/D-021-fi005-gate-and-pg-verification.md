---
id: D-021-fi005-gate-and-pg-verification
doc: decision-entry
status: accepted
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# D-021 · F-I-005 关门口径与 PG 侧经验验证（用户 P-004 裁决）

## 决定的来源

- **F-I-005** 的关闭要求含「已记录的 canonical SQL / `MigrationChecksum`」；而 `MigrationChecksum` 只能对**真实存在的语句切片**求值 → 结构性依赖 R2 代码。这与 F-I-002 曾出现的死锁同类（R1 关门前置需要 R2 产物）。
- **PG 侧从未经验验证**（实测）：本机 `psql` / `pg_dump` / `pg_restore` **均不存在**；5432 **无监听**；`compose.yaml` 仅 `api` + `web` 两个服务，**无 PG 服务**。而 `r1-c2-per-table-pg-ddl-v1.0-fc.md` 的表达式（`date_trunc` + `to_timestamp` / 整数 interval）**全部只是设计**。
- 用户 2026-09-20 经 P-004 裁决两项：**F-I-005 = 选项 A（拆分）**；**PG 验证 = 选项 A（允许临时起容器）**。

## 决定 1 · F-I-005 拆分为「R1 关约定、R2 关哈希」

### R1 关门所需（均**已有**或本轮可完成）

| 项 | 载体 | 状态 |
|----|------|------|
| checksum 计算**约定**（单 checksum / SQLite DDL 切片 / PG 不进哈希 / `transform_id` 无方言后缀） | child `D-017-v73-checksum-convention.md` | 已落盘 |
| 15 个 descriptor 的 `Name` 与 `transform_id` | `r1-c2-descriptor-ledger-v1.0-fc.md` §1 | 已落盘 |
| `MigrationChecksum` **算法与输入结构**（有序语句清单 + `m0–m5` 序位） | 同台账 §2 | 已落盘 |
| 唯一表范围（v73/v74/v85 不相交；v74 = `system_data_reconcile` + auth） | 同台账 §1.1 | 已落盘 |
| append-only 边界（v1–v72 不可变；`rebuildOperationLog` 断言不改 stmts） | `D-019` §5 | 已落盘 |

### 移交 R2 的**显式验收项**（不得静默当已完成）

1. 15 个 descriptor 的 **canonical SQL 切片落码**，并按上述约定计算、**记录真实 `MigrationChecksum` 哈希**到 ledger；
2. `migrate_test.go` / `postgres_test.go` 的 v73+ 追加行与**金额列断言拆分**（`wallet_accounts.balance_total` / `wallet_ledger_entries.amount_delta` 保持 `bigint`）；
3. leftover 21 名补入 PG 断言集合。

- **范围**：本拆分只适用于「哈希值」这一子项；F-I-005 的其余子项（约定、名、算法、唯一范围、append-only 边界）**在 R1 内闭合**。
- **复审触发条件**：R2 落码后首次记录哈希时，须由 independent 复审该哈希与约定一致；若约定被修订，本拆分的 R1 部分随之回退。
- **不得**把本拆分读作「F-I-005 已完成」——它是**有范围的移交**，不是闭合声明。

## 决定 2 · 允许临时起 PostgreSQL 容器做经验验证

- 使用**本地已有镜像**（`postgres:16` / `postgres:15-alpine` / `postgres:17-alpine` 已在本地；Docker 29.7.2 守护进程可用）。
- **临时**容器，**非 5432** 端口映射（避免影响其他服务）；验证完成后**销毁**容器。
- 验证目标：`r1-c2-per-table-pg-ddl-v1.0-fc.md` §0 的两条 PG 表达式——秒族 `date_trunc('microseconds', to_timestamp(col::double precision))` 与毫秒族 `date_trunc('microseconds', TIMESTAMPTZ 'epoch' + col * INTERVAL '1 millisecond')`——在边界值（历元 0、999 ms 无进位、负毫秒、公元 9999 年）上的输出，以及 `timestamptz(6)` 的精度断言。
- **边界**：本决定只授权**本地临时验证容器**；**不**授权改动 `compose.yaml`、**不**授权常驻服务、**不**授权把验证结果当作生产就绪证据。

## 影响与边界

- 本决定**不闭合任何 required**；F-I-005 的 R1 侧仍需 independent 复审确认「约定/名/算法」确实完整。
- PG 验证结果若与设计不符，必须回到 `r1-c2-per-table-pg-ddl-v1.0-fc.md` 修正并留痕，**不得**静默调整表达式。
- 本决定不改变 `D-017` / `D-018` / `D-019` / `D-020`。
