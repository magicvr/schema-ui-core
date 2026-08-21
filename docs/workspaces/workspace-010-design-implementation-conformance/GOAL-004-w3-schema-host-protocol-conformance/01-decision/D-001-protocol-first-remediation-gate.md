---
id: D-001
goal_id: GOAL-004-w3-schema-host-protocol-conformance
title: 协议优先与实现停止线
status: accepted
created: 2026-08-12
updated: 2026-08-12
parent: GOAL-004-w3-schema-host-protocol-conformance
version: 0.1.0
---

# D-001 · 协议优先与实现停止线

## 决定

1. 先完成 Host/App 协议候选盘点与上游协议增补，再修正本仓实现。
2. 在上游增补已合并并形成可核对版本/commit、本仓拿到并固定新协议工件之前，禁止修改 `apps/api` 与 `apps/web` 来“解决”本波发现。
3. 停止线内允许：只读证据核验、协议提案、compatibility matrix、schemas/fixtures 设计与测试计划；不允许以临时本地字段、未登记 extension 或手写分支替代协议。
4. 新协议必须为附件中的每个业务候选给出 `adopt-now`、`reserve-extension` 或 `explicitly-out` 之一。候选目录不是要求所有能力都成为同一版本的 mandatory surface，但禁止无说明遗漏。
5. 实现阶段只修两类项目：协议已覆盖但实现偏离的项目，以及新协议明确采纳且本产品需要的 Host/App 契约。Host adapter 本身不因手写而判错。

## 原因

- 先按旧协议或本地猜测修复，容易把协议缺口固化为新的私有方言。
- auth/bootstrap/branding/shell/error 属于 Host/App 跨层契约，不能只从单个 React 组件或 Go handler 反推标准。
- 明确停止线可以让后续 validator、fixtures、API/Web 修改共享同一权威基线。

## 未选方案

| 方案 | 未选原因 |
|------|----------|
| 直接按 2.7.0 修复全部发现 | 无法处理 Host/App 缺口，且会把未定义语义误判为实现错误 |
| 先在本仓增加私有 `x-*` 字段 | 若无上游 extension 治理，会产生第二套不可移植协议 |
| 协议设计与产品修复同时推进 | 会破坏 S3 实施门禁，导致实现反向约束协议 |
| 把所有 Host UI 改写成 page schema | Host 壳层有独立生命周期和安全边界；正确目标是规范 Host 契约，而非消灭 Host |

## 影响

- I-003 是 S4 的 required 硬门禁；未满足时目标仍可推进 S1～S3，但不得进入正式实现。
- 协议变更按 `cross` 审计；provider 在首次方案审视前由用户指定。
- 若实施中发现协议仍不足，受影响范围回流 S2/I-002，不得临时旁路。
