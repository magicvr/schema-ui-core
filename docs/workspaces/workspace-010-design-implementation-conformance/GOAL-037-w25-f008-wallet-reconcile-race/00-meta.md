---
id: GOAL-037-w25-f008-wallet-reconcile-race
title: W25 承接 · F-008 钱包对账竞态修复（池化+FK 时代偶发不一致）
status: active
created: 2026-08-23
updated: 2026-08-23
parent: GOAL-036-w25-page-performance-guardrails
version: 0.1.0
progress: 0/4
---

# GOAL-037 · W25 承接 · F-008 钱包对账竞态修复

## 概述

承接 GOAL-036（W25）A-001 响应连环发现（E-007 记录）的开放 finding **F-008**：`TestWalletLifecycleAndAdjustFlow` 在池化（池 4）+ WAL + FK=ON 配置下偶发 `reconcile result = inconsistent`（异步对账 job 回放链与账户余额不一致）。已知边界：

- 与 `_txlock=immediate` **无关**（A/B 实证：两种配置下均偶发）；曾被 SQLITE_BUSY 提前失败掩盖，BUSY 已由 `_txlock=immediate` 确定性修复。
- 基线（单连接 W24）该用例稳定，未在池化+FK 配置下做过归属判定——**暂按"本波配置暴露/引入"定位**，机制未定论。
- 与已修缺陷（F-001 FK 每连接 / F-002 回归 / F-007 时钟量化）正交。

**治理约定（用户书面指令 2026-08-23）**：本子目标承载 F-008 的治理上下文；**本子目标关门后**，再回归关门 GOAL-036（GOAL-036 的关门依赖本目标 done）。

## 成功标准

1. **C1 机制定性**：可复现的失败窗口 + 根因证据链（哪一层可见性/顺序假设被池化打破）；
2. **C2 方案冻结**：D-001 取舍（修复或测试语义修正，含未选方案）；
3. **C3 实施与回归**：修复落地 + `TestWalletLifecycleAndAdjustFlow` 高频连跑稳定 + `go test ./...` 全绿；
4. **C4 关门**：A-001（self）pass + 台账同步（goal-tree/workspace + GOAL-036 F-008 关闭引用）。

## 路线图（P-001 · 分母 = 4）

```text
S1 定性   → C1：复现窗口度量 + 根因（异步 job 对账事务与主链写入的可见性；回放读取隔离级别）
S2 方案   → C2：D-001 冻结
S3 实施回归 → C3：修复 + 高频连跑 + 全量
S4 关门   → C4：自审 + 台账（随后 GOAL-036 回归关门）
```

## 信息需求登记（P-005）

| 编号 | 问题 | 级别 | 影响门禁 | 状态 | 证据/结论 |
|------|------|------|----------|------|-----------|
| I-001 | reconcile 不一致的**完整机制**（复现条件、涉及哪层快照/顺序假设、为何池化+FK 暴露） | required | C1/C2（方案冻结前） | **collecting** | E-007 已有部分事实（与 txlock 无关、曾被 BUSY 掩盖）；需复现窗口 + 代码级因果链 |
| I-002 | 高频复现手段与窗口度量（`-count` 频率、是否与校验虚拟时钟/时序相关） | non-blocking | C3 验证 | **collecting** | 当前低频偶发（数十轮～个位数） |

## 边界与审计声明

- 范围仅限 wallet 异步对账（`ReconcileOnceTx` / jobs）与相关测试；不改协议/Profile/模块契约语义。
- 审计模式 `self`（数据一致性敏感面按 `independent` 备选，若根因涉及并发语义由用户裁定）。
- 本目标不关门 GOAL-036；GOAL-036 关门在本目标 done 之后（用户书面约定，见概述）。