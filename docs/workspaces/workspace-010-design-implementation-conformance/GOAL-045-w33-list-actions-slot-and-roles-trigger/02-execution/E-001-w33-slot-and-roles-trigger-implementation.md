---
id: E-001-w33-slot-and-roles-trigger-implementation
doc: execution-entry
parent: GOAL-045-w33-list-actions-slot-and-roles-trigger
status: recorded
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# E-001 · W33 C1～C3 实施（列表 actions 左侧插槽 + roles 触发面）

## 事实（2026-09-19）

### 1. C1 · 方案冻结

`01-decision/D-001-w33-slot-and-roles-trigger-freeze.md` v1.0.0：插槽声明（`props.slot = "list-page-actions"` + `props.targetTable`，未知值按未声明处理）、宿主与可见性（表格提供 `data-list-page-actions-left`，经 CRUD seam 发布；`registerListActionsSlot` 决定该行是否渲染）、**fail-open 原地渲染**回落、行布局 `justify-between` 两段、两页改动清单、授权与可写范围。`I-045-001`/`002` 由该决策关闭为 `verified`。

### 2. C2 · 实施产物

| 产物 | 位置 |
|------|------|
| 插槽 seam：`registerListActionsSlot` / `hasListActionsSlot` / `publishListActionsHost` / `listActionsHost`（状态支撑，宿主出现即触发消费者重渲染） | `apps/web/src/renderer/render.tsx` |
| slot 分支：`declaresListPageActionsSlot`（**纯谓词**，未声明者走原路径）+ `ListActionsSlotNode`（注册 → portal，宿主缺失时原地渲染） | `apps/web/src/renderer/render.tsx` |
| 行拆两段：左段为插槽宿主（`data-list-page-actions-left`），右段为既有列配置 + `props.toolbar`（样式/权限/禁用语义不变）；行渲染条件含 `hasListActionsSlot` | `apps/web/src/renderer/schema-table.tsx` |
| `users` 页：`users-batch-export` 节点补 `"slot": "list-page-actions"` | `apps/api/modules/users/schema/users.json` |
| `roles` 页：`roles-table` 增 `props.selection.mode = multiple`；新增 `roles-batch-export` 节点（`component: jobs-batch-export` / `resource: roles` / 同插槽声明） | `apps/api/modules/roles/schema/roles.json` |

- **后端零改动**：导出分母本就含 `roles`（R3 `D-001` §1），提交体为 `{resource:"roles", ids:[…]}`。
- **i18n 零新增键**：复用 `schema.jobs.batchExport.*`。

### 3. 实施中修正的真实问题

| # | 问题 | 处置 |
|---|------|------|
| 1 | 内联 ref 回调（`ref={(el) => crud?.publishListActionsHost(...)}`）导致 React 每次 commit 先 detach(null) 再 attach(el)，`null↔element` 震荡把发布函数打进**无限更新循环**（测试报 `Maximum update depth exceeded`） | 改为**身份稳定**的 `useCallback` ref + 经 ref 读取 CRUD 值（回调本身不再变化） |
| 2 | `SlotAwareCustomNode` 的 effect 依赖整个 `crud` 对象 → 任何 provider 状态变化都重注册，再次形成循环（测试挂起 10 分钟超时） | 改为只依赖**稳定的** `registerListActionsSlot` 函数（仓库既有写法：`crud?.notifyFeedback` 等同族 seam） |
| 3 | 新 hook 落在组件**早退分支之后**（`columns.length === 0` / `dataSource === null` / `keyCheck` 失败）→ `Rendered fewer hooks than expected`（3 例失败） | 把 hook 上移到组件顶部 hook 区，并注释说明 React 的稳定 hook 计数要求 |
| 4 | 初版让**所有** custom 节点都经过新的包装组件（订阅 CRUD context） | 改为纯谓词分流：仅声明插槽的节点走新分支，其余走**逐字原路径**（不新增订阅、不改变渲染特征） |
| 5 | **浏览器 e2e 回归**（`schema-crud.spec.ts`）：users → roles → users 往返后，users 表**只剩两 schema 列的并集**（`ID/NAME/UPDATED`），用户名列消失 → 创建的用户行匹配不到 | 根因：roles 页新增 custom 节点后，**两页的 table 节点落在同一子节点索引**，React 复用了同一个 `SchemaTable` 实例，把上一张表的 `visibleColumns` 带进另一张表（`columns` 收窄 effect 取交集后非空即保留）。修复：**渲染器按 table 节点 id 作 key**（`<Fragment key={node.id}>`）——不同表即不同实例，跨页不再串状态。**归因证据**：干净树基线跑该 spec **通过**（1 passed），带 W33 改动**失败**；临时 instrumented spec 打印出 `headers#1` 8 列 → roles 8 列 → `headers#2` 3 列；修复后两页各自 8 列 |
| 6 | 上述修复同时消除了一个**既有**隐患：任何两页若把 table 排在相同索引，都会互相串列状态（本次只是第一次被触发） | 已补渲染器回归用例（同树两次切换 users/roles 文档，断言列集合各自完整、不是交集） |

> 说明：问题 4 一度被怀疑是 2 例 `s5-denominator-render` 超时的原因；经**干净树基线实验**排除（见 §4）。

### 4. C3 · 验证证据

| 验证 | 结果 |
|------|------|
| `npm test`（vitest，默认超时） | **121 files / 1476 tests 全绿**（W33 前 119/1468 → +2 files / +8 tests） |
| `npm run typecheck`（`tsc -b` + e2e tsconfig） | exit 0 |
| `npm run build`（vite） | exit 0 |
| `go build ./...` / `go test ./internal/docscheck/ ./modules/users/ ./modules/roles/` | exit 0 / 全绿 |
| e2e 双 profile（`npm run test:e2e`） | 见 §4 末行（修复 table 身份后） |
| **干净树基线实验**（关键归因，两次） | ① 2 例 `s5-denominator-render` 超时：干净树（无任何 W33 改动）**同样复现**（118 files / 1466 tests，2 failed）→ 与 W33 无关，属既有负载敏感（已按 §4 末段补显式超时）；② `schema-crud` e2e 失败：干净树**通过**、带 W33 **失败** → 确认由 W33 引入，已定位并修复（§3 问题 5） |

**新增测试**

| 测试 | 断言 |
|------|------|
| `renderer/list-actions-slot.test.tsx`（6） | 声明插槽的节点**确实**落在目标表格左段；未声明者仍在文档流（且不产生空左段）；目标表不存在 → **原地渲染**（fail-open，控件不消失）；未知 slot 值等同未声明；**两张表跨页切换时列状态互不串扰**（回归锁）；两页 shipped 声明（id/table/resource/slot）被固定 |
| `components/jobs-batch-export-roles.test.tsx`（2） | 真实 `roles.json`：触发面在左段（不在筛选栏）；空选择禁用；勾选 2 行 → 提交体 `{resource:"roles", ids:["role-admin","role-editor"]}`；**不 reloadList**（`/api/roles` 仅 1 次、计数仍在） |

**变异验证（均实测变红后还原）**

| 变异 | 结果 |
|------|------|
| `roles.json` 的 `slot` 值改成 `list-page-action`（拼写错） | `list-actions-slot` 的 shipped 声明例 **红**；roles 页「触发面在左段」例 **红**（2 failed / 5 passed） |

**既有 flake 的处置（本次一并修复）**：`i18n/s5-denominator-render.test.tsx` 的 3 个真实 App 渲染用例补显式超时（20s/30s/20s）并注释证据（干净树复现 + 提升超时转绿），使全套件在默认超时下稳定绿——该文件属 workspace-010 W29 交付面，改动为**测试可靠性**，不含产品行为变化。

### 5. 边界

- **未**改后端（导出分母、job 运行时、权限键零改动）；**未**改 `props.toolbar` 的语义/权限/禁用规则；**未**触碰 pinned 工件。
- **未**改 VP-038 `status`（保持 `closed`）与 workspace-038 台账正文；仅在 `roadmap.md` 登记。
- 右侧按钮组视觉与列配置交互未变；仅新增左侧段与 `justify-between`。

### 6. Git checkpoint

| hash | 内容 |
|------|------|
| 见 `03-audit/A-002` | W33 实施 + 测试 + 台账（C4 条目登记） |

### 7. 未做（移交 C4）

- 未把 roles 页触发面加入 e2e（`jobs-result-center.spec.ts` 仍只走 users 路径）；未新增 independent 审计（模式 `self`，理由见 `03-audit/A-001`）。
