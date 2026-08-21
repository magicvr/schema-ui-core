---
id: E-002
goal_id: GOAL-005-w4-long-content-presentation
title: 执行 · S3 实现整改事实
created: 2026-08-13
updated: 2026-08-13
parent: GOAL-001-design-implementation-conformance
version: 0.1.0
---

# E-002 · S3 实现整改（2026-08-13）

## 改动清单（与 D-001 §2 一致）

| 文件 | 改动 |
|------|------|
| `apps/web/src/components/data-table.tsx` | `DataTableColumn<T>` 新增可选 `truncate`；`cellContent` 字符串兜底在 `truncate === true` 时包 `span.block.max-w-[16rem].truncate` + 原生 `title` 全文 + `data-table-cell="truncated"` 标记；`render` 自定义列与 `—` 空值路径不受影响；默认 `false` 时行为与改动前逐字节一致 |
| `apps/web/src/renderer/schema-table.tsx` | `SchemaTableColumnSpec` 新增可选 `truncate`；`DataTableColumn` 映射处透传 `column.truncate === true` |
| `apps/web/src/renderer/render.tsx` | recordView 详情 grid：`sm:grid-cols-[8rem_1fr]` → `sm:grid-cols-[8rem_minmax(0,1fr)]`（值列可收缩）；`<dd>` 保留 `break-words`；数组值按 `", "` 连接渲染，其余 `String(value)` |
| `apps/api/internal/modules/roles/schema/roles.json` | `roles-table` 的 `permissions`、`menuItems` 两列加 `"truncate": true`（D-001 §2.3 页面启用） |

## D-001 §2.3 前置核实（结构校验器是否拦截未知列字段）

- 本仓 schema 校验事实（E-001 基线 + S3 复核）：
  1. 浏览器加载路径 `loadPageDocument` → `validatePageDocument` → `runtime-schema-validate.ts` 只校验 `page`/`node`/`action`/`reaction` 四类 schema（`page.schema.json` / `node.schema.json` / `action.schema.json` / `reaction.schema.json`）。`node.schema.json` 的 `props` 为宽松对象（`additionalProperties: false` 仅约束 node 顶级键），`columns` 属 `props` 内部，**不在结构校验范围内**；对 `truncate` 未知字段无拦截。
  2. `component-registry.json` 的 `table.props.columns.items` 是 `additionalProperties: false` 的严格列定义（仅 field/label/format/tagMap/labelKey/sortable/sortField/visibleWhen/reactions/permissions）。该文件未被 `loadPageDocument`/`validatePageDocument` 校验路径引用（全仓代码无 `component-registry` 的校验消费点）；`roles.json` 是本地业务页面文档，页面文档校验按 (1) 执行。
  3. 上游 `schema-ui-docs` 的页面文档校验是否有 L2/L3 列语义层：本仓 vendor 的 fixture 分母与校验器不含列严格校验层（E-001 §3）。
- **结论**：`"truncate": true` 属「本地页面文档的呈现扩展字段」，不触发本仓结构校验失败、不修改任何校验器、不触碰上游 fixture 分母（上游无此字段，本仓不向上游提案）。不构成 D-001 §1 的协议方言——协议字段集未变，本地 schema 页文档多一个被 Renderer 白名单解析器（`parseRenderNode` table 分支透传 `props`）携带的呈现开关。
- 兜底确认：`parseRenderNode` 的 table 分支 `...(isRecord(value.props) ? { props: value.props } : {})` 原样透传 `truncate` 到 `SchemaTable`，无需改动解析器。

## 设计要点

- 列表截断放在 `DataTable`（唯一生产消费方 `SchemaTable` 走同一路径）；`roles` 页启用，其他页面默认行为不变（D-001 §2.1/§2.2）。
- 详情换行只动 grid 轨道与数组字符串化；`<aside>` 结构、`overflow-y-auto`、Drawer/Sheet 模式不变。
