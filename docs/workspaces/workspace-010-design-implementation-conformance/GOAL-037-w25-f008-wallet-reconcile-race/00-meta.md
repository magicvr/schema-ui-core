---
id: GOAL-037-w25-f008-wallet-reconcile-race
title: W25 承接 · F-008 钱包对账竞态修复（池化+FK 时代偶发不一致）
status: done
created: 2026-08-23
updated: 2026-08-23
parent: GOAL-036-w25-page-performance-guardrails
version: 0.3.0
progress: 4/4
---

# GOAL-037 · W25 承接 · F-008 钱包对账竞态修复

## 概述

承接 GOAL-036（W25）A-001 响应连环发现（E-007 记录）的开放 finding **F-008**：`TestWalletLifecycleAndAdjustFlow` 在池化（池 4）+ WAL + FK=ON 配置下偶发 `reconcile result = inconsistent`。

**已闭环结论**：根因 = **流水 id 同毫秒乱序**（`newID` 随机后缀在"同一毫秒内"决定回放字典序 ≠ 写入序，回放先遇 `freeze` 未入账 → `insufficient balance`）；与池化无直接因果（池化 + `synchronous=NORMAL` 提高同毫秒连发概率而暴露）。修复：产品与测试替身 id 均加入**同毫秒单调计数**；测试对异步 terminal 审计事件改**轮询断言**。**根治（用户指令 2026-08-23「不留残余」）**：0050 `wallet_ledger_order_repair` 数据修复迁移（既库乱序重排、fail-closed）+ reconcile 成功审计原子化（job 事务内 `RecordOperationTx`）+ 失败路径可观测（E-003）。治理约定（用户书面）：**本目标关门后回归关门 GOAL-036**——已完成。

## 成功标准（全部达成）

1. **C1 机制定性**：失败帧 details + 代码因果链（E-001）；
2. **C2 方案冻结**：D-001（id 内嵌计数；rowid/seq 列方案拒绝；残余书面）；
3. **C3 实施与回归**：修复落地，`-count=100` 100/100，全量 go 两轮全绿，vitest 1097/tsc 0（E-002）；
4. **C4 关门**：A-001（self）pass，台账同步，GOAL-036 回归关门（2026-08-23）。

## 路线图（分母 = 4）——全部完成

```text
S1 定性   → 机制闭环（E-001）✓
S2 方案   → D-001 冻结 ✓
S3 实施回归 → E-002（100/100 + 全量）✓
S4 关门   → A-001 pass + GOAL-036 回归关门 ✓
```

## 信息需求登记（P-005）

| 编号 | 问题 | 级别 | 影响门禁 | 状态 | 证据/结论 |
|------|------|------|----------|------|-----------|
| I-001 | reconcile 不一致的完整机制 | required | C1/C2 | **closed** | 同毫秒流水 id 随机后缀乱序；失败帧 `replay apply failed: insufficient balance`（E-001/E-002） |
| I-002 | 高频复现手段与窗口度量 | non-blocking | C3 | **closed** | `-count=100` 修复前混合失败、修复后 100/100（E-002） |

## 边界与审计声明

- 范围仅限 wallet 异步对账与相关测试；协议/Profile 契约零改动。
- 审计模式 `self`（E-001 证据链完整，无需升级）；残余与复审触发见 D-001（既有库旧流水重对账、terminal 事件 best-effort 完整性课题）。
- 本目标 `done`；GOAL-036 已按用户书面约定回归关门。