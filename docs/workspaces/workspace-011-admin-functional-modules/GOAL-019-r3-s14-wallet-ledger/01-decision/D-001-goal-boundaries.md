---
id: D-001
goal: GOAL-019-r3-s14-wallet-ledger
title: 立项边界：模块身份、Profile 归属与审计策略
date: 2026-08-16
status: accepted
parent: GOAL-019-r3-s14-wallet-ledger
created: 2026-08-16
updated: 2026-08-16
version: 1.0.0
---

# D-001 · 立项边界（S-14 钱包/账务）

## 决策

1. **模块身份**：候选 **admin.wallet**（标准 Admin 模块，R3 第四批次；S1 方案冻结确认最终名与 Descriptor 依赖，预期 core.schema-render / core.navigation-capability / core.operationlog）。账本/账务语义以「余额 + 流水 + 对账」为最小领域面（I-011-001 §4 S-14）。
2. **Profile 归属（I-004）**：进入 **admin 默认集**候选（内容扩展先例 S-01/S-02），S1 确认；mvp/demo 默认不启用。
3. **审计策略**：钱包/账务属 **data 门禁**（余额变动审计 + 迁移基建）→ S1 方案冻结与 S5 关门必须 **grok build independent**（用户书面偏好：grok-4.6 · reasoning high）；本会话无 provider，由用户安排在 grok build /audit 会话执行；S2/S3 以 self 审计。立项 A-002 independent 同样待安排。
4. **无越界**：不改变 Profile 默认集语义 / 模块矩阵 / Manifest 装配语义 / 协议 pin；共享基架问题回流 VP-009/VP-010；go 失效触发时门闩自动关闭；不引入支付通道 / 外部资金结算 / 多租户（Charter 非目标）。

## 未选方案

- 新建独立「账务引擎」共享模块而非 admin 模块 → 范围过大；按分档 S-14 以标准模块落地，S1 再判共享面。
- 本波次引入外部支付/结算集成 → 超出 Admin 功能模块边界与 Charter 非目标。
