---
id: D-001-w32-seams-freeze
doc: decision-entry
parent: GOAL-044-w32-r4-residual-seams
status: active
created: 2026-09-19
updated: 2026-09-19
version: 1.0.0
---

# D-001 · W32 三项通用 seam 方案冻结

承接 `[workspace-038-batch-operations-and-job-center] GOAL-005` 的 `A-001` F-002/F-003/F-004（用户 2026-09-19 指令移交）。本决策冻结三项的落点、语义与边界；用户授权与跨区可写范围见 §5。

## 1. ① 通用列值本地化：列级 `valueLabels`（本地扩展）

**形态**：表格列新增本地扩展属性

```json
{ "field": "status", "labelKey": "schema.jobs.column.status",
  "valueLabels": { "queued": "schema.jobs.status.queued", "…": "…" } }
```

- 语义：单元格**展示值** = `t(valueLabels[raw], undefined, raw)`——即「有映射则用目录文案，无映射/键缺失则**回落到原始值**」。
- **fail-open 的理由**：这是**展示**层映射，不是权限/门禁；回落原始值等于今天的既定行为（`wallet.status`、`scheduledTasks.enabled` 等页面同样显示原始值），因此映射缺失不会让页面变得不可用，只会保留旧观感。反之若 fail-closed 显示空单元格，会把「翻译缺失」升级成「数据缺失」，是更差的失败模式。
- **不新增 pinned 能力 id**、不改 `docs/schemas/**`：`valueLabels` 与既有本地扩展（`badgeStyleField`/`truncate`/`width`/`minWidth`）同族——`node.schema.json` 的 `props` 是不受约束对象（仅禁止 CSS 属性名），仓库既有页面已在用未见于 pinned `component-registry.json` 的列属性。
- **与 pinned `format:"tag"`+`tagMap` 的关系**（`I-044-001`）：pinned 的 `tagMap` 是「值 → {text, tone}」的**字面量**映射，不支持 i18n 键，且本仓库渲染器**未实现** `tag`/`tagMap`（全仓 0 命中）。本扩展只解决「值 → i18n 键」；**不**实现 `tagMap`、不与其竞争。若上游未来实现 `tag`，页面可迁移，两者语义可共存。

## 2. ② 表格定向刷新 seam：`refreshTable(tableId)`

**新增 seam**（`SchemaCrudValue`）：

| 方法 | 语义 |
|------|------|
| `refreshTable(tableId)` | 用**当前查询**重新取该表格的数据；**不触碰选择集**（与 `reloadList()` 的关键差别） |
| `tableRefreshToken(tableId)` | 该表格的刷新令牌（表格的取数 effect 依赖它） |
| `publishTableRows(tableId, rows)` | 表格把自己的当前行**发布**到页面级注册表（ref 支撑，不触发重渲染） |
| `tableRows(tableId)` | 读取该注册表（供 ③ 判定） |

**与既有刷新的分工（`I-044-002`）**：

| seam | 作用面 | 选择集 | 备注 |
|------|--------|--------|------|
| `reloadList()` | 全页列表波次 | **清空全部**（ADR-0022 D2） | 变更后语义，保持不变 |
| `refreshList(dataSource)` | 仅 display 数据源（statCard/chart） | 不涉及 | W25 既有 |
| `refreshTable(tableId)` | **单个 table 节点** | **保留** | 本决策新增 |

**与 in-flight 合并的关系（明确取舍）**：`reloadList()`/`refreshList()` 会从其 in-flight 表里删除目标键，理由是「变更后重拉不得与变更前请求合流」。`refreshTable()` **不删除** in-flight 键——它是**只读轮询**语义，没有变更前/后之分，合流只会省一次请求且返回的仍是服务端最新数据。**该差异是有意的**，避免为轮询再引入一套缓存失效规则；若将来出现「变更后定向刷新」需求，必须改为删除 in-flight 键，届时另立决策。

## 3. ③ 空闲不轮询：行可见性 + 显式活跃态

- 数据来源（`I-044-003`）：② 的 `publishTableRows`/`tableRows`。
- 组件参数（schema 节点 `props`）：
  - `statusField`（默认 `"status"`）：判定字段；
  - `activeStatuses`（字符串数组）：**声明哪些状态需要继续轮询**；缺省 = 不启用空闲判定（保持既有行为，向后兼容）。
- 判定：tick 时若 `tableRows(targetTable)` 中**没有任何行**的 `statusField` 属于 `activeStatuses` → **跳过本次刷新**（不打请求）。行数据不可得（未发布/非表格页）→ 按「有活跃行」处理（保守：宁多刷不漏刷）。

## 4. `jobs` 页接入

- `status` 列加 `valueLabels` 六态映射（复用**既有** `schema.jobs.status.*` 键，两目录已有，零新增键）。
- `jobs-auto-refresh` 节点加 `props: { "targetTable": "jobs-table", "statusField": "status", "activeStatuses": ["queued","running"] }`，组件改用 `refreshTable("jobs-table")`（不再 `reloadList()`）。

## 5. 用户授权与跨工作区可写范围

用户 2026-09-19 指令（原文见 `00-meta.md` §概述）授权：在本区开承载子目标、由 `[workspace-038…]` 有界接受后**转由本子目标执行修正**。据此本目标可写：

- `apps/web/src/renderer/**`（`render.tsx`、`schema-table.tsx`）——**多工作区共享的渲染器**，本区为符合性程序，属其正常工作面；
- `apps/web/src/components/jobs-auto-refresh.tsx`（workspace-038 R4 交付的组件）——按指令修改其内部实现以接入 ②③；
- `apps/api/modules/jobs/schema/jobs.json`（列 `valueLabels` + 组件 props）与 `apps/web/src/i18n/**`（如需）；
- 本目标台账 + `[workspace-038…] GOAL-005` `A-003` 的**闭合回填**（按 P-003 只加闭合注记，不改其 status/progress）。

**仍不可写**：任何 pinned 工件（`docs/schemas/**`、`apps/web/src/protocol/upstream/**`）；Job 六态合同；同步 `batch-delete`；其他工作区台账正文。

## 6. 非目标

- 不实现 pinned `format:"tag"`/`tagMap`/`datetime`；不为 `valueLabels` 申请 pinned 能力 id。
- 不改变 `reloadList()` 的 ADR-0022 D2 语义；不重写 `refreshList`。
- 不做逐页全量审视；不新增页面、不改导航、不改权限键。
