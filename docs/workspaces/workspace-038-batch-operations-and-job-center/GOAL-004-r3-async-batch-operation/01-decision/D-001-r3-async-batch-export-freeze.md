---
id: D-001-r3-async-batch-export-freeze
doc: decision-entry
parent: GOAL-004-r3-async-batch-operation
status: accepted
created: 2026-09-19
updated: 2026-09-19
version: 1.0.0
---

# D-001 · R3 异步批量导出方案冻结

> 本决策关闭 `I-038-011` / `I-038-012`（GOAL-004 C1/C2 方案项）。
> **约束输入**：R1 `GOAL-002/01-decision/D-001-…`（C2 方案 B / C3 首波 1 条 / K-1～K-7）与 R2 `GOAL-003/01-decision/D-001-…`（O-1/O-2 口径、`jobs.write` 归属、`ResultURL` 派生）。
> **证据基线**：`attachments/R3-recon-frontend-batch-trigger.md`（前端触发侦察，编排器核对）+ 既有同步导出实现。

---

## 1. `I-038-011` · 导出数据面与权限口径（C1）

### 1.1 数据面分母（冻结）

| 项 | 冻结值 | 依据 |
|----|--------|------|
| 可导出资源 | **`users` / `roles`** —— 与同步导出**同一分母** | `internal/handler/export.go` 的 `export()` 仅接受这两个资源（其余 404 `RESOURCE_NOT_FOUND`） |
| 选择范围 | **仅选中行**（不是全量）；选择键归一化复用同步批量的 D3 不变式：仅标量、保序去重、空选择拒绝 | `resources.go:842-875`（同步 `batchDelete` 的同款归一化） |
| 单次上限 | **500** id（`BatchExportMaxIDs`） | 同步 batch-delete body 上限 4 KiB ≈ 数百 id；异步单元保持同数量级 |
| 列集 | **复用同步导出的列序**，不另立 | `exportHeaders`/`exportRow`（`export.go`，D-002 §3 冻结列序） |
| 转义 | 复用 `formulaSafe`（`= + - @ TAB CR` 前缀中和）+ UTF-8 BOM + RFC 4180 | 同上；避免两条导出面的安全口径分叉 |
| 结果形态 | Job result = JSON `{resource, rowCount, fileName, csv}` | Job `result` 列必须为合法 JSON（`CompleteWithCommit` 校验），故 CSV 作为字符串承载；`fileName` = `<resource>-selection.csv` |

### 1.2 权限口径（冻结，**双重门禁**）

批量导出提交路由要求**同时**持有：

| 键 | 作用 | 策略 |
|----|------|------|
| `jobs.write` | 授权**提交异步作业** | `PolicyAdmin` |
| `data.export` | 授权**把数据带出**（与同步导出同一键） | `PolicyAdminEditor`（既有） |

**理由（安全关键）**：异步路径导出的是同一份数据。若只要求 `jobs.write`，则新增 Job 运行时**扩大**了数据外带面——一个没有 `data.export` 的角色可经异步端点导出。双重门禁保证「加异步承接」不改变既有数据外带边界。

> 该配对在 provider 与 handler 两处都体现：provider 声明 `jobs.write`，handler 在路由内先 `jobs.write` 再 `data.export`。

### 1.3 未选方案

| 未选 | 理由 |
|------|------|
| 只要求 `jobs.write` | 会扩大数据外带面（见 §1.2） |
| 只要求 `data.export` | 异步作业是可被取消/重试的运行时对象，其提交属写操作，应有独立写键（与 R2 `jobs.read` 对称） |
| 新增导出资源（如 data-dictionary） | 会超出同步导出的分母，使两条导出面不一致；不在首波 |
| 导出全量而非选中行 | 首波冻结的是「批量导出**所选**」（R1 `D-001` §3.1）；全量导出已由同步端点提供 |

---

## 2. `I-038-012` · 前端触发机制实现口径（C2）

### 2.1 选择集获取路径（冻结，实测确认）

**自定义组件经 `useSchemaCrud()?.selection("<tableId>")` 读取当前选中行**（`renderer/render.tsx:271`、`:927-936`；经 `renderer/index.ts:10,15` 导出）。

- `CustomComponentProps.context` **不含**选择集（只有 host/nav 记录：`user` / `features` / `route.params` / `route.query`），故不能经 `context` 取。
- `context` 也**不含** tableId ⇒ 目标表须由节点 `props.targetTable` 显式声明（先例 `activity-export.tsx:26-29`）。
- 选择状态在 `SchemaCrudContext` provider 内（`render.tsx:818`、`:1501`，deps `:1537`），组件会随选择变化重渲染。

**硬前置**：目标表必须声明 `props.selection.mode === "multiple"`（`schema-table.tsx:932-935`），否则 `setSelection` 从不被调用、`selection()` 恒为 undefined。

### 2.2 触发入口位置（冻结）

**在 `users.json` body 的 section `users` 内、`users-table` 之旁新增一个 custom 节点**（`{"type":"custom","id":"...","component":"jobs-batch-export","props":{"targetTable":"users-table"}}`）。

| 备选位置 | 未选理由 |
|---------|---------|
| 表格 toolbar 项 | `schema-table.tsx:1326-1367` 的 toolbar **只渲染 `<button>`**，无 custom 分支；需改渲染器 |
| 表格节点 children | `schema-table.tsx` 从不读 `children`（零命中） |
| 页头 actions 宿主 | 单所有者 `claim()` 已被 `SchemaTable` 占用（`list-surface.tsx:23-30`） |

### 2.3 capability 声明口径（冻结，承接 R2 O-2）

- **不得**在 `users.json` 声明 `table.selection`：其 usage marker 是 `"requiresSelection"`（`capability-declaration.guard.test.ts:51`），而 users.json 不含该文本 ⇒ **守卫会变红**。
- **不得**声明 `actions.batch.request`（marker = `/batch-delete/`，本路径不经过该端点）。
- 即：加 `props.selection`（纯 props 变化）**不**新增 capability 声明。已在 `jobs.json` 验证过同款口径（R2，guard 75/75 通过）。

### 2.4 提交流程（冻结）

1. 组件读 `selection(targetTable)`；空选择时按钮禁用（与同步批量的 `requiresSelection` 行为对齐）。
2. 提交 `POST /api/jobs/batch-export`（body `{resource, ids}`），期望 **202**。
3. **不得**调用 `reloadList()`：任何 reload 成功都会清空选择（`render.tsx:950-957`），而异步提交后选择仍需保留给用户查看。
4. 取 202 体的 `id` 作为 jobId，轮询 `GET /api/jobs/{id}` 直到终态；轮询间隔沿用既有档位（5/10/30s，先例 `monitoring-auto-refresh.tsx:17-22`）。
5. 终态后展示结果入口（`resultUrl`）与行数。

**注意（K-1 的落地）**：本路径**不**经过 ADR-0022 的 `runBatchRequest`——后者的成功分支只判 `response.ok` 然后 `reloadList()`，会丢弃 202 体并抹掉选择。

---

## 3. 进度语义（冻结）

`reporter.Progress` 接受 `0..99`（`repository.go:115-122`），终态 100 由 `CompleteWithCommit` 置。本实现：

| 阶段 | 上报 |
|------|------|
| 每读一行 | `done * 85 / total`（跨选择线性） |
| 读完 | `90` |
| 渲染完成 | `99` |
| 提交终态 | 100（运行时置） |

**理由**：R1 `D-001` §1.2 K-6 指出既有先例（wallet reconcile）只报硬编码 10→100；VP-038 判据 3 要求「可观察进度」，故进度必须由本 handler 真实上报。

---

## 4. 对 C1～C4 的映射

| 检查点 | 覆盖 |
|--------|------|
| C1 异步写面与进度 | §1（数据面/权限/结果形态）+ §3（进度） |
| C2 前端触发 | §2（选择集/位置/capability/流程） |
| C3 同步路径回归 | §2.4 的 K-1 落地：不经 `batchMapping`；同步 `batch-delete` 与其协议 fixture 不动 |
| C4 R3 审计与投影 | 全篇作为被审对象；审计范围见 `03-audit.md` |
