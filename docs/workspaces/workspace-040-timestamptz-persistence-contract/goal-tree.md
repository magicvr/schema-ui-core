---
title: 目标树 · workspace-040-timestamptz-persistence-contract
status: done
created: 2026-09-20
updated: 2026-09-21
parent: null
version: 1.0.1
workspace_id: workspace-040-timestamptz-persistence-contract
---

# 目标树 · DB 时间列 timestamptz 持久化合同

> 工作区：`workspace-040-timestamptz-persistence-contract`
> canonical：`docs/workspaces/workspace-040-timestamptz-persistence-contract/`
> Root：`GOAL-001-timestamptz-persistence-contract`（**`done · 3/3`**，2026-09-21 用户确认关门）
> primary_plan：`VP-040-timestamptz-persistence-contract`（**`closed`** v0.3.0；VRev-105 `pass`）

## 树

```text
GOAL-001-timestamptz-persistence-contract [done] (3/3) · 纲领容器
├── GOAL-002-r1-contract-and-denominator-freeze [done] (4/4) · R1 合同与分母冻结
├── GOAL-003-r2-codec-and-descriptor-m1-m2 [done] (4/4) · R2 M1/M2 codec 与 descriptor
├── GOAL-004-r2-repository-and-predicate-rewrites [done] (3/3) · R2 M3 仓储与谓词改造
├── GOAL-005-r2-backup-port-and-closeout [done] (3/3) · R2 M4 Backup Port 与阶段关门
├── GOAL-006-r3-wire-formatter-and-unit-family-matrix [done] (3/3) · R3 公共 wire formatter 与单位族/时区矩阵
├── GOAL-007-r3-pg-cross-version-restore-matrix [done] (3/3) · R3 PG 15/16/17 pg_restore 跨版本矩阵与升级后恢复核对
└── GOAL-008-r3-exit-matrix-and-root-closeout [done] (3/3) · R3 退出判据矩阵、关门审计与 Root 确认关门
```

## 纲领路线图

```text
R1 合同与分母冻结 [completed · GOAL-002 · 4/4]
   → R2 双方言迁移 + Store 编解码 [completed · M1/M2 = GOAL-003 (4/4)、M3 = GOAL-004 (3/3)、M4 = GOAL-005 (3/3) · 独立关门审计 A-004 pass / open required = 0]
      → R3 读写/时区回归、备份核对、证据与关门 [completed · 边界由 Root D-018 冻结 · R3-A/B = GOAL-006 (done · 3/3) · R3-C = GOAL-007 (done · 3/3) · R3-D = GOAL-008 (done · 3/3，用户 2026-09-21 确认关门)]
```

> R1 = `GOAL-002`（`done · 4/4`）；R2 = `GOAL-003`（4/4）+ `GOAL-004`（3/3）+ `GOAL-005`（3/3），M1–M4 全绿并经两轮独立关门审计（`A-004` **pass / open required = 0**）。Root `progress: 3/3` 由 Root `00-meta.md` 的三个显式检查点派生；**R1/R2/R3 全部完成**，Root 已于 **2026-09-21 经用户书面确认关门**（Root `D-020`；`GOAL-008` 检查点 C）。R3-D 的 independent 关门审计为 `A-002`（conditional/开放 required = 0）与 `A-005`（**pass**/开放 required = 0）；用户报告的 `dev.cmd start` 缺陷 `F-I-101` 已 `fixed` 并独立复现。

## 状态表

| id | title | parent | status | progress | updated |
|----|-------|--------|--------|----------|---------|
| GOAL-001-timestamptz-persistence-contract | DB 时间列 timestamptz 持久化合同 | null | **done** | 3/3 | 2026-09-21 |
| GOAL-002-r1-contract-and-denominator-freeze | R1 · 时间合同与分母冻结 | GOAL-001-timestamptz-persistence-contract | **done** | 4/4 | 2026-09-20 |
| GOAL-003-r2-codec-and-descriptor-m1-m2 | R2 · 共享 codec 与 15 个 conversion descriptor（M1/M2） | GOAL-001-timestamptz-persistence-contract | **done** | 4/4 | 2026-09-20 |
| GOAL-004-r2-repository-and-predicate-rewrites | R2 · 仓储读写与谓词改造（M3） | GOAL-001-timestamptz-persistence-contract | **done** | 3/3 | 2026-09-20 |
| GOAL-005-r2-backup-port-and-closeout | R2 · Backup Port 与阶段关门（M4） | GOAL-001-timestamptz-persistence-contract | **done** | 3/3 | 2026-09-20 |
| GOAL-006-r3-wire-formatter-and-unit-family-matrix | R3 · 公共 wire fixed-6 formatter 与单位族/时区矩阵（R3-A/B） | GOAL-001-timestamptz-persistence-contract | **done** | 3/3 | 2026-09-21 |
| GOAL-007-r3-pg-cross-version-restore-matrix | R3 · PostgreSQL 15/16/17 pg_restore 跨版本矩阵与升级后恢复有界核对（R3-C） | GOAL-001-timestamptz-persistence-contract | **done** | 3/3 | 2026-09-21 |
| GOAL-008-r3-exit-matrix-and-root-closeout | R3 · 退出判据证据矩阵、关门审计与 Root 确认关门（R3-D） | GOAL-001-timestamptz-persistence-contract | **done** | 3/3 | 2026-09-21 |

## 说明

- Root **`done · 3/3`**——**R1、R2、R3 全部完成并关门**：R1（2026-09-20 用户书面确认，关门向 independent = A-046，开放 required = 0）；R2（独立复审 `A-004` pass / open required = 0）；R3-A/B（`GOAL-006`）与 R3-C（`GOAL-007`）于 2026-09-21 静默关门（各 `done · 3/3`）；**R3-D（`GOAL-008`）经退出判据矩阵 + self/independent 关门审计后由用户 2026-09-21 书面确认关门**（Root `D-020`；`GOAL-008` `done · 3/3`）。六条成功标准已勾选；跨目标开放 required = 0。
- **R1 关门边界**：只放行**设计面**。F-I-005 曾为用户书面 `accepted-residual`（范围 + 复审触发 + 失效条件见 child `D-021`），其复审触发（R2 首次记录任一 v73+ 哈希）已发生并经由 `GOAL-002/A-048` 按 **`fixed`** 闭合；R2 的生产 schema 变更已经 R2 自身验收。
- **M2/M3 次序**（**Root `D-017`**，用户 2026-09-20 裁决）：M2 落码后按 `D-018` 无过渡期使全仓测试红 → **M2 不单独提交**；M3 `GOAL-004` 转绿后与 M2 一并提交。`I-041-003` 已按常驻 PostgreSQL 15.4 实测证据关闭（约束：破坏性 migration 只可作用于一次性/专用测试 database）。M4 `GOAL-005` 已按用户确认 slug 立项并 `done · 3/3`。
- R2/R3 的门禁不得由激活状态或 progress 投影替代；本次 Root `done` 的依据是**判据 6 的用户书面确认**，不是 progress 值。
- **VP-040 的波次关闭 / Vision Review 属决策层（`/vision`）动作**，不在本次 Root 关门范围（Root `D-020` §3）。
- 状态、progress、parent 或新增子目标发生变化时，必须同步本文件树与状态表。任何进度值都不放行阶段、不关闭 finding、不推导 `done`。
