---
id: D-001-w33-slot-and-roles-trigger-freeze
doc: decision-entry
parent: GOAL-045-w33-list-actions-slot-and-roles-trigger
status: active
created: 2026-09-19
updated: 2026-09-19
version: 1.0.0
---

# D-001 · W33 列表页 actions 左侧插槽与 roles 触发面方案冻结

用户 2026-09-19 选型：位置方案 **A**（表格 page-actions 行左侧插槽）、治理载体 **workspace-010 W33**。本决策冻结实现口径。

## 1. 插槽的声明与语义

**声明（页面 schema，custom 节点 props）**：

```json
{ "type": "custom", "id": "roles-batch-export", "component": "jobs-batch-export",
  "props": { "targetTable": "roles-table", "resource": "roles",
             "slot": "list-page-actions" } }
```

- `slot: "list-page-actions"`：把该 custom 节点投放进 `targetTable` 表格的**列表控件行左侧**，而不是在文档流中原地渲染。
- 未声明 `slot` 的 custom 节点**行为逐字不变**（原地渲染）——这是本扩展的兼容边界。
- 节点**仍留在 section/body 子节点位置**（不改变 pinned 校验下的合法形态；R3 侦察确认 custom 节点只能是容器子节点）。
- 插槽值当前只定义一个：`list-page-actions`；其它值按「未声明」处理（原地渲染），避免拼写错误导致控件消失。

## 2. 宿主与可见性（`I-045-001`）

- `SchemaTable` 渲染该行时提供 `data-list-page-actions-left` 宿主 `div`；宿主元素通过页面级 CRUD seam 发布：`publishListActionsHost(tableId, el)` / `listActionsHost(tableId)`（**状态支撑**，宿主出现即触发消费者重渲染；参照 `App.tsx` 既有 `setListActionsHost` + 保存视图 portal 先例）。
- 表格是否需要渲染左侧宿主：由 `hasListActionsSlot(tableId)` 决定——声明了插槽的节点在挂载时 `registerListActionsSlot(targetTable)`（effect，返回注销；参照既有 `registerDirtySource` seam），因此**没有插槽消费者的表格不会多出一段空行**。
- 投放：节点渲染为 `createPortal(content, host)`；宿主尚未出现的那一次渲染按 §3 回落。

## 3. 回落策略（`I-045-002`）：fail-open 原地渲染

- `targetTable` 缺失、非字符串、或对应表格不存在 → 节点**原地渲染**（等同未声明 slot）。
- 理由：本扩展是**布局**能力，不是权限或数据门禁。若采用「隐藏」，一处 schema 笔误就会让「导出所选」这个操作入口**静默消失**——这是比「位置不对」严重得多的失败模式；原地渲染则退化成 R3 的既有观感，用户仍能完成操作。
- 与既有约定一致：保存视图在无页头宿主时也走 `inlineSavedViewHeader` 原地回落，而非隐藏。

## 4. 行布局改动

- `data-list-page-actions` 行：`justify-end` → `justify-between`，内部拆为两段：
  - **左段** `data-list-page-actions-left`（插槽宿主，`flex flex-wrap items-center gap-2`）；
  - **右段** 既有内容（列配置 + `props.toolbar` 按钮，保持 `justify-end`）。
- **不改** `props.toolbar` 的语义、权限门控（`effectivePermission`）、批量/选择禁用规则与按钮样式；右段视觉不变。
- 行渲染条件：`列配置存在 || toolbar 非空 || hasListActionsSlot(tableId)`（三者皆无时不渲染该行，避免空行）。

## 5. 页面改动清单

| 页面 | 改动 |
|------|------|
| `users` | `users-batch-export` 节点补 `"slot": "list-page-actions"`（`targetTable` 已有）；按钮即从「筛选栏上方」落到列表控件行左侧。其余不动 |
| `roles` | `roles-table` 增 `"selection": { "mode": "multiple" }`；body 增 `roles-batch-export` 节点（`component: jobs-batch-export`、`targetTable: roles-table`、`resource: roles`、同插槽声明）。**后端零改动** |

- 两页共用同一个组件实例类型（`jobs-batch-export`），`resource` 由节点 props 传入 → 提交体分别为 `{resource:"users"|"roles", ids:[…]}`；服务端分母本就含两者。
- i18n：复用既有键（`schema.jobs.batchExport.*`），预期零新增键。

## 6. 授权与可写范围

用户 2026-09-19 选型授权本波次承载，可写：

- `apps/web/src/renderer/**`（共享渲染器：插槽 seam 与 custom 分支）；
- `apps/api/modules/users/schema/users.json`、`apps/api/modules/roles/schema/roles.json`（VP-038 已交付面的**页面覆盖补齐**）；
- 相关测试与 `docs/vision/roadmap.md` 登记节、本目标台账。

**仍不可写**：任何 pinned 工件（`docs/schemas/**`、`apps/web/src/protocol/upstream/**`）；Job 六态合同与后端导出分母；VP-038 的 `status`（保持 `closed`）与 workspace-038 台账正文（只在 roadmap 交叉登记）。

## 7. 非目标

- 不实现「toolbar 组件条目」（方案 B）与「页头插槽」（方案 C）。
- 不改右侧按钮对齐方式、不改列配置交互、不重排列表页其它区域。
- 不为插槽新增 pinned 能力 id；不实现通用「任意节点任意插槽」框架（只定义 `list-page-actions` 一个值）。
