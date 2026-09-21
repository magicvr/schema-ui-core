---
id: D-017-m3-staging-and-pg-verification-environment
doc: decision-entry
status: accepted
parent: null
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# D-017 · M2/M3 交付次序与 PG 验证环境裁决

## 决定的来源

- **用户 2026-09-20 P-004 裁决**（`/govern` 询问，三问同轮）：
  1. M2 落码后仓储层仍扫 `int64` → `go test ./...` 红（23 个包）；用户选**方案 1**：**立即开设 M3 子目标**，M2 与 M3 的落码**一起跑到绿再提交**，M2 不单独提交红态。
  2. M3 子目标编号/slug：**`GOAL-004-r2-repository-and-predicate-rewrites`**（AGENTS §11：slug 须用户确认）。
  3. `I-041-003` 关闭：**接受**现有常驻 PostgreSQL 15.4 + 已实际执行的 PG 集成测试作为「M3 存在真实可执行 PG 环境」的证据；并附三条约束（见 §3）。
- 触发原因（事实）：按 `D-018` **无过渡期**，v73–v87 把列形状改成 `TEXT`/`timestamptz(6)` 后，仓储层尚未改造 → 全仓测试红。AGENTS §6b 规定「验证失败 fail closed」，故**不得**单独提交该中间红态。

## 1. 交付次序（本决策冻结）

| 阶段 | 交付 | 提交条件 |
|------|------|----------|
| M2 | 15 个 conversion descriptor 的 SQLite `Apply` + PG `ApplyPostgres`、真实 `MigrationChecksum`、ledger 写入路径改造、`rebuildOperationLog` fail-closed 断言、catalog/指纹断言更新 | **不单独提交**（红态 fail closed）；代码留在工作树 |
| M3 | 仓储读写/谓词改造、双方言回归、边界测试重定向、金额列断言拆分、leftover 21 列补入 | 全仓测试**转绿**后与 M2 一并以显式路径提交 |

- 追溯性由 `GOAL-003` 的 `E-00N` 事实记录 + `GOAL-004` 的 `E-00N` + `attachments/r2-v73-v87-generated-statements-v0.1.md`（逐语句清单）承担；**M2/M3 仍是两个独立检查点**，不因合并提交而合并。
- 本决策**不**改变 `D-016` §4 的检查点判据，只决定**提交次序**。

## 2. 未选方案

- **方案 2（先单独提交 M2 红态）**：违反 AGENTS §6b「验证失败 fail closed」，且会让仓库在 M3 完成前处于测试红态。**未采用**（用户亦未授权破例）。
- **方案 3（在 GOAL-003 内追加最小读取改造先转绿）**：与 M3 完整范围重叠、造成返工，且会让 GOAL-003 越过 `D-016` 冻结的 M2 边界。**未采用**。

## 3. `I-041-003` 关闭（用户书面裁决）

**结论**：`I-041-003` → **verified / closed**。

| 项 | 内容 |
|----|------|
| 证据 | `apps/api/configs/.env`（gitignored）配置 `PG_TEST_*` 指向**常驻 PostgreSQL 15.4**（Debian 15.4-2.pgdg120+1, x86_64；host 192.168.31.213:5432，user `sa`，db `postgres`）；PG 侧集成测试**实际连上并执行**（本轮实测：catalog 应用到 v73 时由服务器返回真实 `SQLSTATE 42804`，随后按预期修复） |
| 用户约束 ① | M3 的**破坏性 migration / round-trip 必须运行于一次性或专用测试 database**，**不得**作用于共享业务 schema |
| 用户约束 ② | Docker 隔离 PG / 版本矩阵**不作为** `I-041-003` 的关闭条件，另作为 **CI / release reproducibility** 验证 |
| 用户约束 ③ | M4 的 `pg_dump` / `pg_restore` **可使用固定版本 Docker 临时容器**提供客户端工具（本机无 `psql`/`pg_dump`/`pg_restore` 二进制，`docker` 可用） |

- 该裁决**修订** R1 期间记录的「本机无常驻 PG」前提：**可达的常驻 PG 15.4 现已存在**并被 PG 集成测试实际使用。
- 约束 ① 对 M3 可执行测试的载体（`D-020` 边界测试重定向）**即刻生效**：任何会重建/破坏 schema 的用例必须指向一次性/专用 database。

## 影响与边界

- 本决策只决定**提交次序**与**PG 验证环境**的关闭；**不**放行任何 schema 变更，**不**关闭任何 finding，**不**修改 `D-016` 的检查点判据。
- `I-041-004`（PG 15/16/17 跨版本 `pg_restore` 矩阵）仍为 R3 前复核的 `non-blocking` 项；与本决策的用户约束 ② 一致。
