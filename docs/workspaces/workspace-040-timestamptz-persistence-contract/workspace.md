---
id: workspace-040-timestamptz-persistence-contract
title: DB 时间列 timestamptz 持久化合同工作区
status: done
root_goal: GOAL-001-timestamptz-persistence-contract
canonical_scope: docs/workspaces/workspace-040-timestamptz-persistence-contract/
shared_materials_catalog: none
vision_role: delivery
plan_refs:
  - VP-040-timestamptz-persistence-contract
primary_plan: VP-040-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-21
version: 0.2.0
parent: null
---

# 工作区上下文 · DB 时间列 timestamptz 持久化合同

本工作区是 [VP-040-timestamptz-persistence-contract](../../vision/plans/VP-040-timestamptz-persistence-contract.md) 的唯一 `delivery` workspace，承接架构分支 C1：Store 时间列的 UTC 绝对时刻语义、PostgreSQL 生产权威物理类型、SQLite 合同平等物理类型、双方言迁移台账与读写编解码。

- VP-040 于 2026-09-20 经 `/vision` 激活为 **`active` v0.2.0**；激活就绪 self Review = [VRev-104](../../vision/reviews/VRev-104-vp040-timestamptz-persistence-contract-activation.md) `pass`，open required = 0。
- **Root `[workspace-040-timestamptz-persistence-contract] GOAL-001-timestamptz-persistence-contract` 已于 2026-09-21 经用户书面确认关门（`done · 3/3`）**：R1（`GOAL-002` `done · 4/4`）→ R2（`GOAL-003` 4/4 + `GOAL-004` 3/3 + `GOAL-005` 3/3）→ R3（`GOAL-006` 3/3 + `GOAL-007` 3/3 + `GOAL-008` 3/3）全部完成；六条成功标准已勾选；跨目标开放 required = 0。关门依据：Root `D-020`、`GOAL-008/attachments/r3d-root-exit-criteria-matrix-v0.1.md`、independent 关门审计 `GOAL-008/A-002`（conditional/0）与 `A-005`（pass/0）。
- 用户已将 R1 合同定为 PostgreSQL `timestamptz(6)` + SQLite fixed-6 UTC RFC3339 `TEXT`；所有绝对时刻列纳入分母（90 列 / 44 表）；双方言各自原地转换；sentinel `0` 按语义转为 `NULL`；公共 wire 输出统一固定 6 位微秒 UTC `Z`。
- 红线：不引入 ORM/第三库；不把 SQLite 假装成原生 `timestamptz`；不把驱动时间类型泄漏到 handler/模块公共契约；不消耗 Redis/MQ/多实例/A3 trigger；不重开 VP-013/VP-020。
- **本工作区的实现层工作已收口**；VP-040 自身的波次关闭 / Vision Review 属决策层（`/vision`）动作，不在本工作区 Root 关门范围内。

## 绑定

| 字段 | 当前值 | 说明 |
|------|--------|------|
| 工作区 ID | `workspace-040-timestamptz-persistence-contract` | 本区唯一稳定标识 |
| Root Goal | `GOAL-001-timestamptz-persistence-contract` | `parent: null`；**`done · 3/3`**（2026-09-21 用户确认关门） |
| canonical 范围 | `docs/workspaces/workspace-040-timestamptz-persistence-contract/` | 本区唯一目标状态范围 |
| 共享资料目录 | `none` | 当前无固定共享资料；不得声明共享资料引用 |
| 愿景角色 | `delivery` | VP-040 唯一 delivery workspace；不改变 Charter primary workspace |
| 规划对齐 | `primary_plan` = `VP-040-timestamptz-persistence-contract` | `plan_refs` 必填且精确绑定 |

## 愿景对齐

- Charter：`schema-ui-core-admin-foundation@0.4.0`
- VP：[VP-040-timestamptz-persistence-contract](../../vision/plans/VP-040-timestamptz-persistence-contract.md)（`active` · v0.2.3；本区交付已完成，VP 波次关闭待 `/vision`）
- 计划阶段审视：[VRev-101](../../vision/reviews/VRev-101-vp039-vp040-planned.md) self `pass`
- 激活审视：[VRev-104](../../vision/reviews/VRev-104-vp040-timestamptz-persistence-contract-activation.md) self `pass`
- Vision open required：0；R1/R2/R3 信息门禁均已在实现层关闭（各目标 `00-meta` 与审计台账为证）

## 纲领阶段

| 阶段 | 目的 | 状态 |
|------|------|------|
| R1 | 合同与分母冻结：方言物理类型、列清单、零值/NULL、备份 residual | **completed**（`GOAL-002` `done · 4/4`；关门向 independent = A-046，开放 required = 0；`F-I-005` residual 经 `A-048` 按 `fixed` 闭合） |
| R2 | 双方言迁移 + Store 编解码，沿用不可变 checksum 台账 | **completed**（M1/M2 = `GOAL-003` 4/4、M3 = `GOAL-004` 3/3、M4 = `GOAL-005` 3/3；15/15 v73–v87 checksum 冻结于 `internal/store/migrate_test.go`） |
| R3 | 读写/时区回归 + 备份有界核对 + 证据与关门 | **completed**（R3-A/B = `GOAL-006` 3/3、R3-C = `GOAL-007` 3/3、R3-D = `GOAL-008` 3/3；用户 2026-09-21 确认关门） |

纲领阶段按 R1 → R2 → R3 串行，三段均已完成；新建阶段子目标前必须先在 Root 决策/路线图中冻结边界与信息门禁（本区已无待建阶段）。

## 固定共享资料引用

> `shared_materials_catalog: none`；当前无可作为事实依据的共享资料引用。

| reference_id | workspace_id | material_id | source | version | sha256 | purpose | local_record | status |
|--------------|--------------|-------------|--------|---------|--------|---------|--------------|--------|
| — | — | — | — | — | — | — | — | none |

## 备注

本工作区只保存 VP-040 的实现层目标、决策、执行事实与 Goal 审计；不得把 Vision Review 或 VP 状态复制成第二套目标状态源。愿景组合编排仍以 `docs/vision/` 为准；本区 `goal-tree.md` 与目标五件套是实现层状态真相源。
