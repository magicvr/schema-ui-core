---
id: GOAL-005-r2-backup-port-and-closeout
title: R2 · Backup Port 与阶段关门（M4）
status: done
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 1.0.0
progress: 3/3
plan_refs:
  - VP-040-timestamptz-persistence-contract
primary_plan: VP-040-timestamptz-persistence-contract
serves_summary: 承接 Root R2 的 M4：落码 D-016 第 9 项（kernel.RecoveryPointPort 类型表面 + apps/api/internal/backup provider，含 <recovery-artifact> 校验、错误分类与 restore harness），收口 D-021 residual 与 A-002 待复审事项 7（部分升级/批级快照）的用户裁决，并完成 R2 的 self + independent 关门审计。
---

# GOAL-005 · R2 · Backup Port 与阶段关门（M4）

## 概述

本子目标承接 `GOAL-001-timestamptz-persistence-contract` 的 **R2 阶段 M4 检查点**（`D-016` §4）。M1/M2 已由 `GOAL-003` 关门（`done · 4/4`），M3 已由 `GOAL-004` 关门（`done · 3/3`）。

**子目标编号与 slug 经用户 2026-09-20 确认**（AGENTS §11：禁止静默默认 slug）：`GOAL-005-r2-backup-port-and-closeout`。

## 范围

| # | 交付 | 判据来源 |
|--:|------|----------|
| 1 | `kernel.RecoveryPointPort` 类型表面（最小方法集） | `D-007`；`D-010`；`D-016` §2 第 9 项 |
| 2 | `apps/api/internal/backup/` provider：`<recovery-artifact>` 校验、错误分类、restore harness | `D-006`；C3 边界 §4.1；`D-016` §2 第 9 项 |
| 3 | PG 侧 provider 的可执行路径（本机无 `pg_dump`/`pg_restore` 二进制 → 可用固定版本 Docker 临时容器提供客户端工具，Root `D-017` §3 约束③） | `D-017` §3；C3 边界 |
| 4 | `A-002` 待复审事项 7 的用户裁决：descriptor 各自事务导致的**部分升级**是否需在 M4 前补批级快照 / 整批回滚 | `GOAL-004/03-audit/A-002` 待复审 7 |
| 5 | R2 关门：self + independent 关门审计、required 合法闭合、Root 判据 2/3/4 的 R2 侧证据归集 | `D-016` §4 M4；Root 成功标准 |

## 非目标（本子目标**不**做）

- **不**做备份的调度、鉴权、远端存储、保留策略、KMS/TLS、UI（C3 边界已排除）。
- **不**做 VP-020 展示/输入时区回归矩阵与其用例 ID（**R3**）。
- **不**实施公共 wire formatter（**R3**，`D-016`）。
- **不**改 v1–v87 的 canonical SQL / checksum / identity / Apply 语义。
- **不**引入 ORM / 第三库 / Redis / MQ / 多实例。
- **不**把 Docker 隔离 PG / 版本矩阵当作生产就绪证据（`D-017` §3 约束②；`I-041-004` 仍属 R3 前复核）。

## 红线

- 破坏性 migration / restore 演练**只能**作用于一次性或专用测试 database（Root `D-017` §3 约束①）。
- 不把本地临时容器结果当作生产就绪证据。
- 不以「测试通过」代替独立审计结论；R2 关门必须经 independent 关门审计。

## 成功标准（本子目标检查点，用于 `progress` 派生）

| 检查点 | 判据 | 状态 |
|--------|------|------|
| **A** | `kernel.RecoveryPointPort` + `internal/backup` provider 落码（含 `<recovery-artifact>` 校验与错误分类），SQLite 侧 restore harness 有可执行测试 | **completed**（E-002） |
| **B** | PG 侧可执行路径（Docker 固定版本客户端工具）验证；部分升级/批级快照的用户裁决落盘 | **completed**（E-002 PG 实测通过 + E-003 C3 §4.2 全部接线；`I-041-006` 经用户 2026-09-20 裁决为 accepted） |
| **C** | R2 关门：self + independent 关门审计落盘、required 闭合、`D-021` residual 收口状态复核 | **completed**（E-004：`A-001` self →`A-002` independent（required=3）→`A-003` 修复响应 →`A-004` 复审判定 **pass / open required = 0**；`D-021` residual 已于 `GOAL-002/03-audit/A-048` 闭合） |

`progress: 3/3` 由 A～C 等权派生；**不**放行阶段、**不**关闭 finding、**不**推导 `done`。

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 | 影响门禁 | 状态 | 证据 |
|----|------|----------|----------|------|------|
| I-041-004 | non-blocking（继承） | PG 15/16/17 跨版本 `pg_restore` 兼容矩阵 | R3 前 | deferred（R3 侧 residual；Root `D-017` §3 约束②） | `I-040-003` 登记的 R3 residual |
| I-041-006 | required（本目标新增） | 部分升级语义：descriptor 各自事务下 v85 类失败会留下 v73–v84 已提交 | B | **verified（用户 2026-09-20 P-004：接受现有语义——可续跑 + A/C 回滚，不做批级原子化）** | `01-decision/D-001-partial-upgrade-and-c3-anchor-scope.md` |

## 父目标

- `GOAL-001-timestamptz-persistence-contract`（R2 检查点）

## 关门条件

**A～C 全部完成**、`03-audit` 的 self 与 independent 关门审计意见落盘、`required` finding 合法闭合、`I-041-006` 由用户裁决关闭后，才可关门；关门即代表 **R2 阶段完成**（Root `progress` → 2/3），但**不**放行 R3 的展示/输入时区矩阵与跨版本矩阵。
