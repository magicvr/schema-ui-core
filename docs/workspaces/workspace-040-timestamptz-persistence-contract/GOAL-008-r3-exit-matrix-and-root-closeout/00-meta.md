---
id: GOAL-008-r3-exit-matrix-and-root-closeout
title: R3 · 退出判据证据矩阵、关门审计与 Root 确认关门（R3-D）
status: active
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-21
updated: 2026-09-21
version: 0.1.0
progress: 0/3
plan_refs:
  - VP-040-timestamptz-persistence-contract
primary_plan: VP-040-timestamptz-persistence-contract
serves_summary: 承接 Root R3 的 R3-D：把 Root 六条成功标准逐条落成可核对证据矩阵（引用 GOAL-002～GOAL-007 的已关门结论与审计台账），完成 self + grok independent 关门审计并合法闭合 required，最后按判据 6 取得用户确认后把 Root 标为 done。
---

# GOAL-008 · 退出判据证据矩阵、关门审计与 Root 确认关门

## 概述

承接 `GOAL-001-timestamptz-persistence-contract` 的 **R3 阶段 R3-D 检查点**（Root `D-018` §2 第 5 项）。R3-A/B（`GOAL-006`）与 R3-C（`GOAL-007`）均已 `done · 3/3`；本目标只做**退出与关门**，不新增实现范围。

**子目标编号与 slug 经用户 2026-09-20 预确认**（AGENTS §11；Root `D-019` §2）：`GOAL-008-r3-exit-matrix-and-root-closeout`。

## 范围

| # | 交付 | 依据 |
|--:|------|------|
| 1 | **退出矩阵**：Root 六条成功标准逐条落成「判据 → 证据 → 结论」，每条指向**已关门目标**的可核对产物（台账/附件/测试） | Root `00-meta.md` 成功标准；`D-018` §2 第 5 项 |
| 2 | **判据 5 的反向核验**：无 ORM/第三数据库/Redis/MQ/多实例；无 Admin 维护提示或业务域混入；驱动类型未泄漏进公共面 | 同上；Root 红线 |
| 3 | **残留与例外清账**：`F-I-005`（R1 `D-021` residual，已 `fixed`）、`I-041-004`/`I-041-008`/`I-040-004` 等已关闭项，以及 R3-C 矩阵的**证据定位边界**（容器 = CI/reproducibility）在退出矩阵中显式点名 | 各目标 `03-audit` |
| 4 | self + grok independent **关门审计**；required 合法闭合 | `D-018` §4 R3-D |
| 5 | **用户确认关门**（判据 6）→ Root `status: done`、`progress: 3/3` | Root `00-meta.md`；`D-018` §2 第 5 项 |

## 非目标（本子目标**不**做）

- **不**重开 R1/R2/R3-A/B/R3-C 的冻结决策、canonical SQL/checksum（v1–v87 不可变）、descriptor、Store codec、C3 Port、wire 合同或跨版本矩阵结论。
- **不**新增功能、迁移、依赖或测试范围；本目标只做**证据汇总与关门**。
- **不**把容器矩阵结果改写成生产就绪证据（`D-017` §3 约束②）。
- **不**在没有用户确认的情况下把 Root 标为 `done`（判据 6 是硬门禁）。

## 红线

- 退出矩阵的每一行必须有**可核对产物**；禁止用「测试通过」代替审计结论，禁止把 `progress` 当作判据满足。
- 存在未合法闭合的 required / 必改 findings 时**不得**关门（三路径闭合：`fixed` / `accepted-residual` / `user-overruled`）。
- **不得**静默代替用户作出 Root 关门裁决。

## 成功标准（本子目标检查点，用于 `progress` 派生）

| 检查点 | 判据 | 状态 |
|--------|------|------|
| **A** | 退出矩阵落盘（六条判据逐条证据 + 判据 5 反向核验 + 残留/例外清账）；每条指向可核对产物 | pending |
| **B** | self + grok independent 关门审计落盘；开放 required = 0（required 按三路径合法闭合） | pending |
| **C** | **用户确认关门**（判据 6）；Root `status: done` / `progress: 3/3`；goal-tree 与 `docs/vision` 投影同步 | pending |

`progress: 0/3` 由 A～C 等权派生；**不**放行阶段、**不**关闭 finding、**不**推导 `done`。

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 | 影响门禁 | 状态 | 证据 |
|----|------|----------|----------|------|------|
| I-041-010 | required（本目标新增） | Root 关门的**用户确认**（判据 6）：是否接受退出矩阵与残留清账后关门 | C | **open**（P-004 必须由用户裁决） | Root `00-meta.md` 判据 6；`D-018` §2 第 5 项 |
| I-041-011 | non-blocking（本目标新增） | 产品级「PG 客户端/服务端支持组合」是否写入发布说明或退出矩阵（R3-C `A-002` 移交的建议 P-004 点） | A | **open** | `GOAL-007/03-audit/A-002` §「明确结论」；`GOAL-007/D-002` §4 |

## 父目标

- `GOAL-001-timestamptz-persistence-contract`（R3 检查点）

## 关门条件

**A～C 全部完成**、`03-audit` 的 self 与 independent 意见落盘、required 合法闭合、且**用户已确认 Root 关门**；本节门条件即 Root 关门条件本身。
