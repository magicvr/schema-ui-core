---
id: GOAL-001-timestamptz-persistence-contract
title: DB 时间列 timestamptz 持久化合同
status: active
parent: null
created: 2026-09-20
updated: 2026-09-20
version: 0.6.0
progress: 2/3
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
| R1 | 合同与分母冻结：方言物理类型、列清单、零值/NULL、备份 residual | **completed**（GOAL-002 `done · 4/4`；关门向 independent = A-046，开放 required = 0；用户 2026-09-20 书面确认） |
| R2 | 双方言迁移 + Store 编解码 | **completed**（M1–M4 全绿：M1/M2 = `GOAL-003` `done · 4/4`、M3 = `GOAL-004` `done · 3/3`、M4 = `GOAL-005` `done · 3/3`；独立关门审计 `A-004` **pass / open required = 0**） |
| R3 | 读写/时区回归、备份有界核对、证据与关门 | **active**（边界与门禁已冻结：Root `D-018`；**R3-A/B = `GOAL-006` `active · 2/3`**：formatter+fixture+输入矩阵+单位族矩阵+VP-020 时区 round-trip 均已完成，检查点 C 待 self + grok independent 审计；**R3-C/D 待渐进立项**，编号/slug 已经用户预确认 = `GOAL-007-r3-pg-cross-version-restore-matrix` / `GOAL-008-r3-exit-matrix-and-root-closeout`，见 Root `D-019` §2） |

`progress: 2/3` = 2/3 个检查点完成（**R1、R2 completed**；R3 pending）。progress 只作展示，不放行阶段、不关闭 finding、不推导 `done`。

> **R1 关门边界**：本次关门只放行**设计面**。F-I-005 为 `accepted-residual`（**不得读作哈希已验证**，复审触发 = R2 首次记录任一 v73+ 哈希时）；R2 的生产 schema 变更仍须经 R2 自身验收。R1 关门**不等于**判据 2/3/4 已满足——它们分别属 R2/R3。
>
> **R2 边界**（Root `D-016`）：范围含 15 个 descriptor 落码、共享 codec、仓储读写改造、测试改写与金额列拆分、边界测试重定向、`rebuildOperationLog` 断言、Backup Port 类型表面与 provider；**非目标**含 VP-020 回归矩阵（R3）、备份调度/鉴权/远端存储/UI、ORM/第三库/Redis/MQ。完成判据 **M1–M4**；三项 required 信息项 `I-041-001`～`003` 在对应门禁前关闭。
>
> **M2/M3 次序**（Root `D-017`，用户 2026-09-20 P-004 裁决）：按 `D-018` 无过渡期，M2 落码后仓储层未改造 → 全仓测试红，故 **M2 不单独提交**；`GOAL-004-r2-repository-and-predicate-rewrites`（M3）转绿后与 M2 一并提交。`I-041-003` 已按常驻 PostgreSQL 15.4 + 实际 PG 集成执行**关闭**，并附「破坏性 migration 只可作用于一次性/专用测试 database」约束。
>
> **R3 信息门禁更新**（Root `D-019`，用户 2026-09-20 P-004 裁决）：`I-041-008`（公共 wire 输出的破坏性 / 兼容期）= **无破坏性、不需兼容期 → verified**；`I-040-004` 由 `GOAL-006` 检查点 B 的证据关闭。R3-C/R3-D 的编号与 slug 已经用户预确认，仍按 `D-018` §6 在 `GOAL-006` 关门后渐进立项。

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-040-001 | required | SQLite 用什么物理类型与 PG 合同平等？精度、UTC 表示、NULL/零值与编解码如何定义？ | R1 冻结、R2 迁移 | R1 | 对照 VP-013 合同与用户 D-002；R1 冻结前不得实施不可逆迁移 | collecting | R1 冻结前复核；责任人：R1 子目标 | 用户已选 PG `timestamptz(6)` + SQLite fixed-6 UTC RFC3339 TEXT；sentinel 0 → NULL；待逐列证据 |
| I-040-002 | required | 哪些绝对时刻列进入首波分母，哪些 INTEGER 明确排除？ | R1/R2 | R1 | 全仓扫描并冻结列清单 | collecting | R1 冻结前复核；责任人：R1 子目标 | 用户已选全部绝对时刻列；ID/duration/step/version/计数/金额/flag 排除 |
| I-040-003 | required | 存量库升级策略与备份 residual 是什么？ | R1/R3 | R1 | 对照 VP-013/016 dump/restore 路径；设计双方言原地转换失败策略 | collecting | R1 冻结前复核；责任人：R1 子目标 | 用户已选 SQLite/PG 各自原地转换；不提供 SQLite→PG 产品搬运器 |
| I-040-004 | required | 与 VP-020 展示合同的回归矩阵如何覆盖会话时区与 UTC 存储？ | R3 | R1 | 复用 VP-020 用例并补存储形状对照 | **verified** | — | `GOAL-006/02-execution/E-004`：Go wire/瞬时 round-trip（含 DST 边界、+05:45、负 epoch）+ Web 会话时区展示 round-trip（L1/L2/L3 层解析 + 秒粒度恒等） |
| I-040-005 | required | 激活前置：VP-039 波次、架构 freshness、激活 self Review、slug 是否满足？ | 激活 | 激活前 | `/vision` 核对 VP-039 `closed`、`6197e802` → `b0a6789b` freshness 与 VRev-104 | verified | — | VRev-104 `pass`；VP-039 `closed` v0.3.0；workspace/Root slug 已落盘 |

## 父目标

- Root：`parent: null`

## 台账布局

本目标使用平铺五件套与三个 ledger 目录：`01-decision/`、`02-execution/`、`03-audit/`。

## 备注

- 激活与开区只建立治理骨架，不代表 schema、迁移、编解码、备份或回归已完成。
- 新建阶段子目标前，先在本 Root 决策/路线图中冻结阶段边界；纲领阶段 R1 → R2 → R3 串行。
