---
id: D-001
goal: GOAL-016-r3-s09-data-permission
title: 立项边界：模块身份、Profile 归属与审计策略
date: 2026-08-15
status: accepted
parent: GOAL-016-r3-s09-data-permission
created: 2026-08-15
updated: 2026-08-15
version: 1.0.0
---

# D-001 · 立项边界（S-09 数据权限）

## 决策

1. **模块身份**：候选 admin.data-permission（标准 Admin 模块，R3 第三批次；S1 方案冻结确认最终名与 Descriptor 依赖，预期 core.auth-session / core.schema-render / core.navigation-capability / core.operationlog）。
2. **Profile 归属（I-003）**：进入 **admin 默认集**候选（内容扩展先例 S-01/S-02），S1 确认；mvp/demo 默认不启用。
3. **审计策略**：数据权限属 **data 门禁** → S1 方案冻结与 S5 关门必须 **grok build independent**（用户书面偏好：grok-4.6 · reasoning high）；S2/S3 以 self 审计。
4. **无越界**：不改变 Profile 默认集语义 / 模块矩阵 / Manifest 装配语义 / 协议 pin；共享基架问题回流 VP-009/VP-010；go 失效触发时门闩自动关闭。

## 未选方案

- 新建独立「数据权限服务」共享模块而非 admin 模块 → 范围过大，按分档 S-09 以标准模块落地，S1 再判共享面。
