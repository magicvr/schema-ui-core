---
id: E-001
goal_id: GOAL-005-w4-long-content-presentation
title: 执行 · S1 基线与最小复现核实
created: 2026-08-13
updated: 2026-08-13
parent: GOAL-001-design-implementation-conformance
version: 0.1.0
---

# E-001 · S1 基线与最小复现核实（2026-08-13）

## 事实

1. **列表挤列机制（代码事实，S1 核实）**：
   - `apps/api/internal/modules/roles/schema/roles.json` 的 `roles-table` 列 `permissions`、`menuItems` 为数组字段；API 层原样返回数组（`apps/api/internal/handler/roles.go` `permissions`/`menuItems`）。
   - `DataTable.cellContent` 对无 `render` 的列走 `String(value)`（`apps/web/src/components/data-table.tsx:50`）：数组经 `Array.prototype.toString` 变成 `"a,b,c"`（逗号连接、无空格）。
   - 桌面表格单元格（`data-table.tsx:316-327`）无 `max-width`/截断类；表格 `w-full min-w-[32rem]` 自动布局下，长串列的 min-content = 全串宽度 → 把 ID/Key/Name 等列挤出可视范围，只能横向滚动。
   - 移动卡片列表已有 `truncate`（`data-table.tsx:126,133`），不受本问题影响。

2. **详情横向滚动机制（代码事实，S1 核实）**：
   - recordView 值渲染（`apps/web/src/renderer/render.tsx:1528-1534`）：`<dd className="break-words">` 已有 `overflow-wrap: break-word`，但所在 grid 为 `sm:grid-cols-[8rem_1fr]`：`1fr` 即 `minmax(auto, 1fr)`，min-content 不被 `break-word` 降低，长无空格串会把列撑宽；`<aside>` 固定 `max-w-md` 且仅 `overflow-y-auto` → 长值横向溢出，出现向右滚屏。
   - 移动端（`max-md` 全宽 Sheet）为单列 grid，同样存在长值横向溢出风险。

3. **协议呈现语义核对（I-001）**：
   - 本地 `docs/schemas/capability-registry.json`（`protocolVersion: "2.8"`）无任何截断/换行/列宽呈现 capability。
   - 已 vendor 的上游 fixtures（`table-sort`、`component-format`、`search-table` 等）与 `RenderTableNode` 列规范（`field`/`label`/`sortable`）均未定义单元格呈现语义。
   - 结论：协议**未定义**列呈现语义 → 属 Host/Renderer 呈现自由，不构成协议缺口；本波为呈现层整改，不新增/变更 capability。

4. **受影响面（I-002）**：
   - 共享路径：`DataTable.cellContent` 兜底（唯一生产消费方为 `SchemaTable`）；recordView 详情值渲染（所有页面的 `recordView` 抽屉）。
   - 长内容列盘点：`roles`（permissions/menuItems）、`users`（roles 数组）、`activity`（detail 长文本）；其余短标量列不受影响。

## 结论

I-001、I-002 证据齐备；最小复现机制已按代码事实确认（未启动浏览器，jsdom 测试将按同一机制断言）。
