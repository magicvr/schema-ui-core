---
title: 执行记录 · R1 · 默认 Renderer 主路径与示例降级
status: active
created: 2026-08-01
updated: 2026-08-01
parent: null
version: 0.3.0
---

# 执行记录 · GOAL-003

## 2026-08-01 · 立项

- 用户确认按 Root D-004 创建 R1 子目标；本目标对应「默认 Renderer 主路径 + 示例降级」。
- 建立五件套；`parent` = `GOAL-001-production-admin-foundation`；`status` = `active`；`progress` = `0/4`。
- 硬依赖 GOAL-002；与 GOAL-004 协作验收可演示页。

> 本节仅记录立项；尚未修改 `App.tsx` 或示例注册表。

## 2026-08-01 · 记录决策 D-003（示例兼容策略）

- 经 `/govern` 回应用户：`I-003-001` 定为**迁移为 Schema**（用户否决双分支与仅测试保留）。
- 写入 D-003：5 个 `EXAMPLE_PAGES` 语义迁移为 Schema 文档；`EXAMPLE_PAGES` 移出渲染路径（源文件暂留参考，GOAL-004 落地后清理）；可测试显式入口 = schemaUrl 链（`page.route` → `page.schemaUrl` → `loadPageDocument` → `RenderPage`）；5 份文档由 GOAL-004 作者；成功标准 2/4 同步修订。
- `I-003-001` → closed（证据 D-003）。
- **尚未改动 `App.tsx` / `registry.tsx`**；默认分支切换实施与测试留待下一轮（届时按 P-004.3.1 询问是否补 self 审计）。

## 2026-08-01 · 默认分支切换实施（P-004.3.1 门禁闭合后）

- **P-004.3.1 补审**：用户选择补 self 审计；A-002（source=self，覆盖 D-003 修订后定义）已落盘 `03-audit`，verdict = pass，与 A-001 汇总同向无冲突。门禁闭合后开始实施。
- **`apps/web/src/app/App.tsx`**：
  - 移除 `EXAMPLE_PAGES[pageId]` 默认查找与 `registry` 导入；移除「renderer remains a later protocol boundary」占位及其统计卡 / “Manifest-driven shell is ready” 区块。
  - 新增 `SchemaPageSurface`：按 `page.route → page.schemaUrl → loadPageDocument → RenderPage` 渲染，含 loading / 统一错误 / ready 三态；`App` 新增可注入 `schemaFetcher` prop（测试注入 fixture，`loadPageDocument.fetcher` 可注入，无需 API 先行）。
  - 新增 `PageSchemaErrorSurface`：`PageSchemaError` 统一 fail-closed 错误面（code / message / url / issues）。
- **`apps/web/src/renderer/render.tsx`**：table 节点在无 `tableRenderer` 时的 fallback 文案由「wired by the example page」改为「table data injection is pending the R1 list-data contract」（示例不再参与渲染路径；数据注入契约属 I-003-002）。
- **测试**：
  - `apps/web/src/app/app-examples.test.tsx` 重写为 schema-driven 默认路径断言：(a) 示例路由渲染注入 Schema 内容；(b) 手写示例标记（Context snapshot / evaluateExpression）不再出现；(c) 404 → `PAGE_NOT_FOUND`；(d) 非法文档 → `PAGE_SCHEMA_INVALID`。
  - `apps/web/src/app/App.integration.test.tsx` 注入 schema fixtures（home/catalog/catalog-detail），断言默认路径渲染 Schema 内容、参数化路由（/catalog/42）经 schemaUrl 展开加载；新增「页面文档缺失 → 统一 `PAGE_NOT_FOUND`」fail-closed 用例。
- **验证**：`npm test` **407/407 全绿**（含 app-examples 4、App.integration 6、load-page 10）；`npm run build`（`tsc -b` + vite）通过。
- **成功标准对照**：1 / 3 / 4 已勾选（meta `progress` 0/4 → 3/4）；标准 2 的 5 份迁移 Schema 文档由 GOAL-004 作为页面资产落地后闭合。
- **`I-003-002`**（required · 列表页数据注入）保持 **open**；列表页验收前与 GOAL-004 / `I-004-002` 复核。
- **手写示例源文件处置**：`registry.tsx` + 5 个示例组件**暂留**为参考（D-003），GOAL-004 文档落地后清理；本目标已不再引用渲染路径。
