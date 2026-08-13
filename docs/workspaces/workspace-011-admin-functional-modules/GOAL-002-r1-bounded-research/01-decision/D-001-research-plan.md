---
id: D-001
doc: decision-entry
goal: GOAL-002-r1-bounded-research
status: accepted
created: 2026-08-14
updated: 2026-08-14
version: 0.1.0
---

# D-001 · 调研计划：来源 / 三档判据 / 对照基线 / 产出形态

## 决策（S1）

1. **候选池来源**：业界 admin 框架与 SaaS 后台的共性功能清单（如 Ant Design Pro / vue-element-admin / react-admin / Strapi / Directus / Supabase Studio / 企业后台模板）作为通用能力样本；常用业务领域（订单 / 钱包 / 类目 / 通知）来自用户点名 + 业界电商/平台后台惯例。来源逐项登记（S2 证据）。
2. **三档判据**（与 VP-011 `三档方法论一致）：
   - 一等公民：业界普遍存在、几乎所有 Admin 都需要、且当前基架尚未覆盖；
   - 常用：高频使用但非普遍必备；
   - 增补：低频、按需、可由 fork 项目按需启用。
3. **对照基线**：已交付基架（users/roles/settings/activity/operationlog 模块、既有页面集、协议面 I-PROTO-FULL-001 v2.8.0）；「已覆盖」= 既有模块/页面/协议已提供等价能力。
4. **产出形态**：分档清单（附件 `I-011-001-tiered-inventory.md`）+ Root 纲领路线图 R2/R3/R4 细化回写。
5. **边界**：本目标只产出清单与判定，不实现；不改变 Profile/模块矩阵/Manifest/协议 pin/共同门禁语义。

## 审计模式

**self**（信息收集；来源可复核；R2 立项时可补独立复审）。

## 未选方案

- 以用户经验直接判档不收集证据：分档将影响后续多波次立项，需可复核来源。
- 调研扩大为协议能力扩展提案：超出本目标；协议缺口单独回流（上游提案或 /vision）。
