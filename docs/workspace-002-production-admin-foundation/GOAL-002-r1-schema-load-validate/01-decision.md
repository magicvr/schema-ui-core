---
title: 决策 · R1 · Schema 加载、校验与统一错误面
status: active
created: 2026-08-01
updated: 2026-08-01
parent: null
version: 0.2.0
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

## D-003 · 页面文档由 Go 端点 `/api/schema/{pageId}` 提供（I-002-001 裁决）

- **日期**：2026-08-01
- **状态**：accepted
- **决定**：
  1. 页面文档以 JSON fixtures 由 Go 提供 `GET /api/schema/{pageId}`；manifest 的 `schemaUrl=/api/schema/<pageId>` 声明**保持不变**，不做静态化改写。
  2. 前端加载器 `loadPageDocument` 运行时按已解析 `schemaUrl`（含路由参数模板展开，复用 `resolveSchemaUrl`）`fetch`，加载路径强制**浏览器安全 Ajv 结构校验**（D-VAL，page/node），非法文档 fail-closed；网络/解析/校验/ID 不符失败统一暴露 `PageSchemaError`。
  3. 浏览器安全校验复用 pinned `docs/schemas/*`（经 `@schemas` 别名构建期导入），**不**重写语义、**不**扩 `v0.1.3`。
  4. 本轮 Go 仅提供加载器所需的最小页面文档子集（`overview`/`catalog`）；代表性页面全集属 GOAL-004。仍**不**切换 App 默认渲染分支（GOAL-003）。
- **理由**：用户对 `I-002-001` 裁决采用 Go 端点方案。Vite dev 已 proxy `/api/*`→`127.0.0.1:8080`；`/api/schema/*` 本就是 manifest 已声明契约，保持契约比改写 schemaUrl 更少漂移；与 R4「Schema 驱动 CRUD 走真实 API」形态一致。
- **未选方案**：
  - **A · Vite `public/` 静态 JSON + 改写 schemaUrl**：偏离已声明 `/api/schema/*`，R4 换真实 API 时需再改 schemaUrl；静态资产与后端契约双轨漂移。
  - **C · 构建内嵌 import JSON**：无运行时网络/解析失败路径，与成功标准「网络/解析失败暴露统一错误」冲突。
- **关联信息项**：`I-002-001` → `verified`（证据：`schema.go`、`load-page.ts`、`schema_test.go`、`load-page.test.ts`）。
- **边界**：`shared_materials_catalog: none`，无共享资料固定引用；页面文档为本地产品数据，不是 `I-PROTO-001` 冻结契约的一部分。
