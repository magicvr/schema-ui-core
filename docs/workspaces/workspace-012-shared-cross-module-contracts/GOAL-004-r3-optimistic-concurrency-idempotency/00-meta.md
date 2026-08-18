---
id: GOAL-004-r3-optimistic-concurrency-idempotency
title: R3 · 乐观并发与幂等契约
status: active
parent: GOAL-001-shared-cross-module-contracts
created: 2026-08-18
updated: 2026-08-18
version: 0.1.0
plan_refs:
  - VP-012-shared-cross-module-contracts
primary_plan: VP-012-shared-cross-module-contracts
serves_summary: 将 wallet 已有版本锁与账本幂等能力提升为可复用、可观测的 expectedVersion/ETag/409 与 operation replay 契约。
---

# GOAL-004 · R3 · 乐观并发与幂等契约

## 概述

R3 以 `admin.wallet` 为首个真实消费模块，固化跨模块可复用的版本前置条件与幂等结果语义：调用方可通过 `expectedVersion` 或 `If-Match` 防止覆盖，并通过稳定 `operationId` 识别同 key replay。

## 范围

- 提供独立的版本 ETag / `If-Match` / expected version 解析契约。
- wallet account 读取和写入返回版本 ETag；status 更新接受 `expectedVersion`，兼容既有 `version`。
- stale version 统一返回 409；缺失/非法/互相矛盾的前置条件 fail closed。
- wallet ledger mutation 返回稳定 operation 结果：`operationId`、`state`、`idempotencyKey`、`replayed`、`resourceVersion`。
- 同账户同 key 同 payload 重放返回同一 operation；同 key 异载荷返回 409；无 key 的旧调用保持可用。

## 非目标

- 不在 R3 批量改造所有 repository，也不为 settings/scheduledtasks/dictionary 新增迁移。
- 不实现 R4 的异步 Job 状态机、队列、进度或取消。
- 不改变 Profile、Manifest 装配、模块矩阵或 Tier D 边界。

## 纲领路线图（P-001）

| 阶段 | 内容 | 状态 |
|------|------|------|
| S0 | 扫描现有 version/ETag/idempotency/transaction 语义，冻结边界 | ✅ 已完成（D-001/E-001；待 A-001） |
| S1 | 实现共享版本前置条件与 wallet ETag/expectedVersion | 未开始 |
| S2 | 实现 operation replay 结果、冲突语义与回归测试 | 未开始 |
| S3 | 全量验证、自审、independent 审计与关门 | 未开始 |

## 成功标准

1. wallet 单资源响应具有稳定 ETag，`If-Match` 与 `expectedVersion` 可互换且不一致时拒绝。
2. stale 写返回 409；缺失或非法前置条件不被当作 version 0 静默接受。
3. 同 key 同 payload replay 返回相同 `operationId` 且 `replayed=true`，不重复写账本；异载荷返回 409。
4. shared contract、wallet repository/service/HTTP 与兼容路径均有测试；API 全量验证通过。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 现有 version/ETag/idempotency 实现与真实消费切片是否足以冻结契约 | S1 方案冻结 | S0 结束前 | 扫描 API repository/handler/migration/tests | verified | 2026-08-18 | E-001：wallet 已有 version CAS、409、账户作用域唯一 key 与事务；Manifest 仅有 If-None-Match |
| I-002 | required | 审计模式与 provider | S3 关门 | S1 实施前 | 按 data/compatibility 风险分级 | verified | 2026-08-18 | 模式 independent；provider 为项目级 grok-build (grok-4.6 reasoning high) |

## 父目标

- `GOAL-001-shared-cross-module-contracts`（Root；依赖已关闭 R1）

## 台账布局

五件套 + `01-decision/`、`02-execution/`、`03-audit/`、`attachments/`。
