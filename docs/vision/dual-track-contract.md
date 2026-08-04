---
doc_type: vision-contract
title: 双线分支维护契约（MVP 基架线 / 完整 Admin 能力线）
status: active
created: 2026-08-04
updated: 2026-08-04
parent: null
version: 0.1.0
vision_ref: schema-ui-core-admin-foundation@0.1.0
---

# 双线分支维护契约

本契约响应 Vision Review `F-V003`（recommended，VRev-001～005）：Charter 成功边界第 4 条确认维护两条可 fork 演进线，但命名、协议兼容、回合并方向与变更发布方式此前未落盘。本文件在方向 3（业务能力）VP 建立前固化该策略；它是**组合编排级文档**，不自动派生 VP / 工作区实现细节，也不替代 [roadmap.md](roadmap.md) 的编排索引。

## 1. 两条线（命名）

| 线 | 名称（建议默认，可修订） | 定位 | 历史载体 | 当前状态 |
|----|--------------------------|------|----------|----------|
| A 线 | **MVP 基架线**（发布前缀 `schema-ui-core-mvp`） | 最小可扩展、低成本 fork 起点；只承载冻结协议子集与骨架 | VP-001 / workspace-001（closed） | 已交付；作为可裁剪基线维护 |
| B 线 | **完整 Admin 能力线**（发布前缀 `schema-ui-core-admin`） | 可直接接业务的 Schema 驱动 Admin 基架 | VP-002 / workspace-002（closed，2026-08-04） | 已交付；作为活跃主线维护 |

命名与发布前缀为**建议默认值**，用户可修订；修订只需更新本表并留痕，不构成 Charter strategic。

## 2. 协议兼容策略

- 两线共享**同一协议固定点**：`schema-ui-docs` v2.7.0，pinned commit `ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b`；冻结覆盖基线 `I-PROTO-001 v0.1.3`（workspace-001 Root `D-009` 冻结）。
- 任何**升级上游协议版本、扩大冻结子集、改变 `D-UPLOAD` 排除**，都必须：追加新决策 → 递增覆盖表版本 → 在受影响 `/govern` 信息门禁前完成验证（沿用 VP-002「继承的协议基线」规则；本契约不重开该规则）。
- B 线是 A 线的**边界内超集**：B 线不得破坏 A 线已冻结的协议语义；A 线新能力（若有）须评估对 B 线兼容影响后回灌。

## 3. 回合并方向

- **活跃主线 = B 线（完整 Admin 能力线）**：缺陷修复与加固默认落在 B 线。
- **A 线（MVP 基架线）为可裁剪基线**：接收 B 线的协议兼容性修复（cherry-pick / 等价回灌），不反向承载 B 线新能力。
- 方向 3（订单、钱包、类目、通知等业务能力）属于 B 线之上的**扩展实现线**（Charter 非目标 1 明确不属本愿景成功条件），不改变两线边界。

## 4. 变更发布方式

- **版本语义**：A 线按协议兼容修复发布 patch/minor；B 线按能力波次发布 minor。发布必须附迁移说明——既有 SQLite 迁移链 `0001`～`0008` 已建立惯例：升级走迁移、降级走 `pre-vNNNN` 快照。
- **fork 指引**：以根 [QUICKSTART.md](../../QUICKSTART.md) 为唯一入口（GOAL-008 交付），随 B 线维护；A 线 fork 从历史 closed 载体（VP-001 基线）裁剪。
- **发布一致性**：本仓库为 Skills 消费仓；契约矩阵 runtime 证据以生成仓发布溯源为准（沿用 VRev-003 `F-V007` accepted-residual 边界），不在本仓复跑。

## 5. 生效与维护

- 方向 3 VP 建立**前**必须复核本契约（F-V003 closure 路径）。
- 后续任何修订追加「修订史」节并更新 `updated`。
- 本契约不改变 Charter 目的/边界/非目标（属 editorial 级决策记录）；组合编排索引见 [roadmap.md](roadmap.md)。

## 修订史

| 日期 | 版本 | 变更 |
|------|------|------|
| 2026-08-04 | `0.1.0` | 响应 `F-V003`（recommended）落盘：两线命名、协议兼容、回合并方向与发布方式四项策略；命名与回合并方向为建议默认值（用户可修订）。 |
