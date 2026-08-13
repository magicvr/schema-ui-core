---
id: D-001
goal_id: GOAL-005-w4-long-content-presentation
title: 决策 · S2 方案冻结：协议呈现语义处置 + 截断/换行实现设计
status: accepted
created: 2026-08-13
updated: 2026-08-13
parent: GOAL-001-design-implementation-conformance
version: 0.1.0
---

# D-001 · S2 方案冻结（2026-08-13）

## 1. 呈现语义处置（I-001）

- **结论**：协议未定义列表单元格截断/换行/列宽语义（E-001 §3）→ **呈现自由**。
- **处置**：`explicitly-out`（对协议：不新增 capability、不向上游提案；对实现：本波只在 Host 呈现层做窄修复，沿用 D-004 §4 dual-end 结构，不引入任何本地协议方言或私有 schema 字段）。
- **依据**：E-001 §3 核对证据（capability-registry 2.8、upstream fixtures、`RenderTableNode` 列规范）。

## 2. 实现设计（S3）

### 2.1 列表：长列截断 + 全文可发现（`components/data-table.tsx`）

- `DataTableColumn<T>` 新增可选布尔 `truncate`（默认 `false`，**行为零变化**）。作用于 `cellContent` 的字符串兜底输出：包一层 `span`，类 `block max-w-[16rem] truncate`，并加原生 `title={text}`（全文 affordance；移动卡片列表已有同类 affordance）。
  - 设 `truncate` 的列无论值长短都走该包裹（短值视觉无变化，`title` 与文本相同，无害）。
- 该列机制不感知值类型（数组/长串统一经 `String(value)`），避免为数组引入类型分叉。
- 现有 `render` 自定义列不受影响；既有表格（所有 schema 页）默认无截断，**不改变现状**。

### 2.2 Schema 层触发（`renderer/schema-table.tsx`）

- `DataTableColumn` 映射处：`truncate: column.truncate === true`（`SchemaTableColumnSpec` 同步增加可选 `truncate`）。
- 未知字段 fail-closed 纪律：`isColumnSpec` 只验证 `field`（既有），新增布尔只作类型守卫后透传（不合法即视为 `undefined`）。

### 2.3 页面启用（`apps/api/internal/modules/roles/schema/roles.json`）

- 仅 `roles-table` 的 `permissions`、`menuItems` 两列加 `"truncate": true`。
- API/Web 结构校验器若对未知字段 fail-closed：**先核实再改**；若会失败，按 ADR-0034 纪律在 D-001 增补「结构校验器是否拦截」处置项，且不改 API 语义、不破坏 fixture 分母。处置证据随 S3 落盘。

### 2.4 详情：recordView 自动换行（`renderer/render.tsx`）

- 详情 grid 修复（`render.tsx:1529`）：`sm:grid-cols-[8rem_1fr]` → `sm:grid-cols-[8rem_minmax(0,1fr)]`，使 `1fr` 可收缩；`<dd>` 保留 `break-words`，并把 `String(value)` 改为可换行渲染：数组按 `", "` 连接（提升可读性），其余原样字符串化。
- 移动端单列 grid 亦受益（可收缩 + break-words）；`<aside>` 保持 `overflow-y-auto`，横向溢出被消除。

## 3. 验收标准（S4 对照）

1. roles 列表：permissions/menuItems 单元格存在截断容器（`truncate` 类）且 `title` 包含全文；其他列不因长值被挤出（jsdom 断言类与属性）。
2. recordView：长数组值渲染为 `", "` 连接且容器可收缩（断言 `minmax(0,1fr)` 类与文本连接符）。
3. 全量既有 `apps/web` vitest 通过；`apps/web` build（tsc + vite）通过；既有 conformance fixtures 回归通过（不触碰 fixture 分母）。
4. 除 roles 两列与上述组件/渲染层外，无其他行为变化（diff 面与 §2 一致）。

## 4. 信息项状态（S2 出口）

- I-001 → `verified`（E-001 §3：协议未定义 → 呈现自由 → explicitly-out）。
- I-002 → `verified`（E-001 §4 受影响面清单）。
- I-003（non-blocking，截断交互形态）→ 采用「单行截断 + 原生 `title`」，范围/复审触发：若后续用户反馈需悬浮/展开形态，另立决策。

## 5. 方案冻结声明

本决策为 S2 方案冻结件（`status: accepted`）；随 A-001（self 方案审视）组成 S2 出口；S3 实施以此为准，不默认扩大范围。S3 中的 §2.3 结构校验器核实结果以 E 条目落盘。
