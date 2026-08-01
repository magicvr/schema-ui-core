---
title: 决策 · R1 · Schema 加载、校验与统一错误面
status: active
created: 2026-08-01
updated: 2026-08-01
parent: null
version: 0.1.0
---

# 决策 · GOAL-002

## D-001 · 本目标只交付加载管线，不切换默认路由分支

- **日期**：2026-08-01
- **状态**：accepted
- **决定**：本目标交付 `schemaUrl` → fetch/解析 → 结构校验 → 错误面；**不**修改 `App.tsx` 的 `EXAMPLE_PAGES` 优先逻辑（属 GOAL-003）。
- **理由**：差量矩阵将「加载/校验」与「默认主路径」拆开，便于独立验收与并行准备页面资产。
- **未选**：与默认路径合并为一个大目标——范围过大、门禁混杂。

## D-002 · 校验复用已 pin 的 schemas，不新写语义

- **日期**：2026-08-01
- **状态**：accepted
- **决定**：结构校验优先复用 `docs/schemas/` 与 `protocol/conformance/schema-validate.ts`（或等价抽取）；**不**重新定义上游 node/page 语义，**不**扩大 `v0.1.3`。
- **理由**：I-PROTO-004 vendor 与 provenance 已 pin；R1 产品化是串联，不是重做协议。
