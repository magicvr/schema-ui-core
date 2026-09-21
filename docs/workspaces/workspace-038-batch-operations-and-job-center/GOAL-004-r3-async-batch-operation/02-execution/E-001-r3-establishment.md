---
id: E-001-r3-establishment
doc: execution-entry
parent: GOAL-004-r3-async-batch-operation
status: recorded
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# E-001 · R3 立项与前端触发侦察

## 事实（2026-09-19）

### 1. 立项

按 P-001，Root `GOAL-001` 的纲领路线图已就位且 R2 已关门，故在 R3 阶段创建子目标 `GOAL-004-r3-async-batch-operation`（五件套 + 三个 ledger 目录齐全），承接 R1 `D-001` §5 的 **T-5/T-6** 与 `jobs.write` 声明。

`GOAL-004` 的 4 个显式检查点：

| 检查点 | 内容 | 承接 |
|--------|------|------|
| C1 | 异步写面与进度 | T-5、`jobs.write`（`I-038-011`） |
| C2 | 前端触发 | T-6 / O-1 / O-2 的实现（`I-038-012`） |
| C3 | 同步路径回归 | 同步 `batch-delete` 不退化 |
| C4 | R3 审计与投影 | — |

审计模式在实施前按 P-002 判定为 **`cross`**（新增写面 + 数据外带面 + 既有同步路径不回退保证 = 权限/数据高影响门禁）。

### 2. 前端触发机制侦察（只读）

以只读侦察确认「自定义组件能否取到当前表格选择集」，产出报告 `attachments/R3-recon-frontend-batch-trigger.md`，关键结论：

| # | 结论 | 证据 |
|---|------|------|
| 1 | 选择集**可**从自定义组件读取，但**不经** `context`：seam 是 `useSchemaCrud()?.selection("<tableId>")` → `{keys,count} \| undefined` | `render.tsx:271`、`:927-936`；`renderer/index.ts:10,15` |
| 2 | `CustomComponentProps.context` 只含 host/nav 记录（`user`/`features`/`route.*`），**不含**选择集也不含 tableId | `AuthGate.tsx:83-86`；`App.tsx:786-793`；`render.tsx:3240,3351` |
| 3 | 硬前置：目标表必须声明 `props.selection.mode === "multiple"`，否则 `setSelection` 从不被调用 | `schema-table.tsx:932-935` |
| 4 | 任何 `reloadList()` 成功都会清空全部选择 ⇒ 异步提交后**不得** reload | `render.tsx:950-957` |
| 5 | 触发入口不能放 toolbar（只渲染 `<button>`）、不能放表格 children（从不读）、不能放页头宿主（已被 `claim()` 占用）⇒ 放在 section 内表格之旁 | `schema-table.tsx:1326-1367`；`list-surface.tsx:23-30` |
| 6 | 声明 `table.selection` 会使 guard 变红（marker = `"requiresSelection"`，users.json 无该文本）⇒ 只加 `props.selection`，不加 capability | `capability-declaration.guard.test.ts:51` |
| 7 | 既有同步批量路径的成功分支只判 `response.ok` 后 `reloadList()`，会丢弃 202 体 ⇒ 复用不可行（K-1 落地） | `render.tsx:765-778`、`:1240-1248` |
| 8 | `users.json` 当前无 `props.selection`、无 `actions.batch.request`；table id = `users-table` | `modules/users/schema/users.json` |

### 3. 本轮**未**做的事（边界）

- **未**改动 `apps/**`：侦察为纯只读。
- **未**冻结方案：`I-038-011`/`012` 的方案口径在 D-001 落盘（见 `E-002`）。
- **未**执行审计（C4 待 C1～C3 完成后进行）。

### 4. Git checkpoint

见提交记录（R3 立项 + 方案冻结 + 侦察报告）。
