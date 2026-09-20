---
id: GOAL-001-timestamptz-persistence-contract
title: DB 时间列 timestamptz 持久化合同
status: active
parent: null
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
progress: 0/3
plan_refs:
  - VP-040-timestamptz-persistence-contract
primary_plan: VP-040-timestamptz-persistence-contract
serves_summary: 在 Charter 0.4.0 与 VP-013 双方言 Store 合同之上，冻结并交付 DB 时间列的 UTC 绝对时刻、PG/SQLite 合同平等物理类型、双方言迁移与编解码；不引入 ORM/第三库，不把驱动类型泄漏进公共面。
---

# GOAL-001 · DB 时间列 timestamptz 持久化合同

## 概述

承接 [VP-040-timestamptz-persistence-contract](../../../vision/plans/VP-040-timestamptz-persistence-contract.md)（`active` v0.2.0；2026-09-20 激活，激活 self Review = [VRev-104](../../../vision/reviews/VRev-104-vp040-timestamptz-persistence-contract-activation.md) `pass`）。本 Root 是本工作区唯一总目标，`parent: null`。

本目标处理 C1 / `RT-T03` 的持久化层合同：生产权威 PostgreSQL 的时间物理类型、SQLite 合同平等物理类型、UTC 绝对时刻语义、纳入分母的 `*_at` 列、读写编解码、不可变迁移台账与备份/恢复边界。它不重做 VP-020 的展示/输入时区能力，不重开 VP-013 的 Store 方言决策，也不实现 Redis/MQ/A3/多实例。

## 愿景对齐

- Charter：`schema-ui-core-admin-foundation@0.4.0`
- VP：`VP-040-timestamptz-persistence-contract`（`active` · v0.2.0）
- 工作区：`workspace-040-timestamptz-persistence-contract`（`delivery`）
- 激活依据：[VRev-104](../../../vision/reviews/VRev-104-vp040-timestamptz-persistence-contract-activation.md) self `pass`；Vision open required = 0

## 红线（激活即生效）

- 不引入 ORM、第三数据库或第二套迁移真相。
- 不把 SQLite 假装成原生 `timestamptz`；物理类型与逻辑语义必须形成可核对的合同平等关系。
- 不让 `pgtype`、驱动时间类型或 `*sql.Tx` 泄漏进 handler/模块公共契约。
- 不把所有 `INTEGER` 静默当作时间列；金额、flag、version、bot_id 等必须显式排除或分类。
- 不消耗 Redis/MQ/多实例/A3 trigger；不重开 VP-013、VP-020、VP-009 或 VP-010。

## 成功标准（对应 VP-040 六条方向级退出判据）

- [ ] 判据 1：PG 物理类型、SQLite 合同平等物理类型、UTC 语义、NULL/零值、编解码与公共面禁止泄漏均已书面冻结（`I-040-001`）。
- [ ] 判据 2：R1 纳入分母的时间列双方言迁移完成并通过 checksum；非时间 `INTEGER` 显式排除。
- [ ] 判据 3：写入绝对时刻与读回一致；VP-020 展示/输入时区不漂移；双方言回归含至少一条 PG 路径。
- [ ] 判据 4：R1 范围内 SQLite 快照 / PG dump 路径可升级后恢复，或用户书面 residual 明确点名。
- [ ] 判据 5：未引入 ORM/第三库/Redis/MQ/多实例；未把 Admin 维护提示或业务域混入。
- [ ] 判据 6：退出矩阵与必要独立意见落盘，开放 required = 0，并经用户确认关门。

## 纲领路线图

以下 3 个检查点是 progress 的唯一来源，默认等权：

| 检查点 | 目的 | 状态 |
|---------|------|------|
| R1 | 合同与分母冻结：方言物理类型、列清单、零值/NULL、备份 residual | pending |
| R2 | 双方言迁移 + Store 编解码 | pending |
| R3 | 读写/时区回归、备份有界核对、证据与关门 | pending |

`progress: 0/3` = 0/3 个检查点完成。progress 只作展示，不放行阶段、不关闭 finding、不推导 `done`。

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-040-001 | required | SQLite 用什么物理类型与 PG 合同平等？精度、UTC 表示、NULL/零值与编解码如何定义？ | R1 冻结、R2 迁移 | R1 | 对照 VP-013 合同平等原则与现有 `INTEGER`/`BIGINT` 时间列；R1 冻结前不得实施不可逆迁移 | collecting | R1 冻结前复核；责任人：本 Root 编排 | 激活默认候选：SQLite `INTEGER` / PostgreSQL `BIGINT`，暂按 UTC Unix seconds；非最终冻结；现有毫秒字段须单独分类 |
| I-040-002 | required | 哪些 `*_at` 列进入首波分母，哪些 `INTEGER` 明确排除？ | R1/R2 | R1 | 全仓扫描时间列与非时间整数列，形成可核对列清单 | open | — | 待 R1 |
| I-040-003 | required | 存量库升级策略与备份 residual 是什么？ | R1/R3 | R1 | 对照 VP-013/016 dump/restore 路径；必要时用户书面接受有界 residual | open | — | 待 R1 |
| I-040-004 | required | 与 VP-020 展示合同的回归矩阵如何覆盖会话时区与 UTC 存储？ | R3 | R1 | 复用 VP-020 用例并补存储形状变更对照 | open | — | 待 R1/R3 |
| I-040-005 | required | 激活前置：VP-039 波次、架构 freshness、激活 self Review、slug 是否满足？ | 激活 | 激活前 | `/vision` 核对 VP-039 `closed`、`6197e802` → `b0a6789b` freshness 与 VRev-104 | verified | — | VRev-104 `pass`；VP-039 `closed` v0.3.0；workspace/Root slug 已落盘 |

## 父目标

- Root：`parent: null`

## 台账布局

本目标使用平铺五件套与三个 ledger 目录：`01-decision/`、`02-execution/`、`03-audit/`。

## 备注

- 激活与开区只建立治理骨架，不代表 schema、迁移、编解码、备份或回归已完成。
- 新建阶段子目标前，先在本 Root 决策/路线图中冻结阶段边界；纲领阶段 R1 → R2 → R3 串行。
