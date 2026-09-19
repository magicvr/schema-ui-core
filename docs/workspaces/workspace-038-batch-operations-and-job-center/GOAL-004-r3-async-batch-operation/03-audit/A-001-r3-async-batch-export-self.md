---
id: A-001-r3-async-batch-export-self
doc: audit-entry
parent: GOAL-004-r3-async-batch-operation
status: recorded
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# A-001 · R3 异步批量导出自审（GOAL-004 C1～C3）

## A-001 · R3 C1～C3 自审（2026-09-19）

- **source**：self
- **auditor**：`/govern` 编排器（DeepSeek Harness 会话）
- **类型** / **scope**：`stage` · GOAL-004 C1～C3（异步写面与进度 / 前端触发 / 同步路径回归）
- **verdict**：**pass**（0 required；3 recommended）
- **完整意见**：本文件

### 范围与区间

- 被审目标：`docs/workspaces/workspace-038-batch-operations-and-job-center/GOAL-004-r3-async-batch-operation/`
- 工作区校验：`root_goal` = `GOAL-001-batch-operations-and-job-center`、`canonical_scope` = 本区路径、`vision_role: delivery`、`plan_refs`/`primary_plan` = `VP-038-…` —— 绑定一致，未跨区。
- 审计区间：2026-09-19（R3 立项 → 方案冻结 → C1～C3 实施 → 验证）。**不含** C4 本身与 R4。
- 审计材料：`00-meta.md`、`01-decision.md` + `01-decision/D-001-…`、`02-execution/E-001/E-002`、`attachments/R3-recon-frontend-batch-trigger.md`；代码 `modules/jobs/export.go`、`internal/handler/jobs.go`、`internal/handler/export.go`、`internal/handler/jobs_export_source.go`、`apps/web/src/components/jobs-batch-export.tsx`、`modules/users/schema/users.json`、`internal/composition/composition.go`、`kernel/profile.go`。

### 成果（有证据）

| # | 成果 | 证据 |
|---|------|------|
| 1 | 首条真实批量操作以异步 Job 承接：kind `jobs.batch-export` + `POST /api/jobs/batch-export` → **202 + job 投影** | `modules/jobs/export.go`；`internal/handler/jobs.go`；`TestJobsBatchExportSubmitReturns202` |
| 2 | **进度真实细粒度**：按选中行线性上报（`done*85/total`）→ 90 渲染 → 99 收尾；100 由 `CompleteWithCommit` 置 | `export.go` 的 `runExport`；`TestBatchExportJobReportsRealProgress`（断言语义为非硬编码） |
| 3 | 结果经**既有** R2 端点可读（未另建结果面）；result = `{resource,rowCount,fileName,csv}` | `TestBatchExportJobProducesCSVResult`；复用 `GET /api/jobs/{id}/result` |
| 4 | **双重门禁**：`jobs.write` **且** `data.export`——异步路径不扩大数据外带面 | `jobs.go` 路由内两次 `requirePermission`；`TestJobsBatchExportRequiresBothGates`（且断言被拒请求**不得到达 submitter**） |
| 5 | 导出列集/转义与同步导出**同源**（抽出 `exportHeaders`/`SelectedExportRows`），两面不会漂移 | `internal/handler/export.go`；CSV 带 UTF-8 BOM 与冻结列序（测试断言） |
| 6 | 校验先于建行：被拒请求**不留 job 行** | `TestJobsBatchExportValidation`、`TestBatchExportSubmitValidation`（均断言 `ListJobs` total = 0） |
| 7 | K-1 落地：**未经** ADR-0022 `batchMapping`/`runBatchRequest` 提交 | `jobs-batch-export.tsx` 直接 `fetch` 本地端点；未引用 `batchMapping` |
| 8 | O-2 落地：`users.json` **未**声明 `table.selection` / `actions.batch.request` | `modules/users/schema/users.json`；`capability-declaration.guard.test.ts` 通过 |
| 9 | 选中集经 `useSchemaCrud().selection(tableId)` 读取（已实测的 seam），空选择禁用提交 | `jobs-batch-export.tsx`；`render.tsx:271,927-936` |
| 10 | **不调用 `reloadList()`**（否则会清空选择）且提交后保留选择 | `jobs-batch-export.tsx`；D-001 §2.4 |
| 11 | 同步 `batch-delete` 与其协议 fixture **未改动** | `git show --stat` 各 checkpoint 无 `resources.go`/`representative-pages.integration.test.tsx`/协议 fixture 变更；回归锚点 48/48 绿 |
| 12 | 未触碰 pinned 协议工件 | 各 checkpoint `--stat` 无 `docs/schemas/**`、`apps/web/src/protocol/upstream/**` |

### 对照检查点

| 检查点 | 状态 | 证据 |
|--------|------|------|
| C1 异步写面与进度 | **达成** | 成果 1～6 |
| C2 前端触发 | **达成** | 成果 7～10 |
| C3 同步路径回归 | **达成** | 成果 11；Go 全绿 + web 1440 全绿 + 回归锚点 48/48 |
| C4 审计与投影 | 进行中 | 本条即 self 腿；independent 腿待跑 |

### 独立复核抽查（编排器对自身实现的复验）

| 主张 | 复核结果 |
|------|---------|
| 双重门禁真的生效（不是只声明） | **成立**：路由内两次 `requirePermission`；editor 403 且 submitter 未被调用 |
| 被拒请求不留 job 行 | **成立**：两组测试均断言 `ListJobs` total = 0（校验在建行之前） |
| 进度非硬编码 | **成立**：慢行源下观察到中间采样；实现按行线性计算 |
| 列集与同步导出同源 | **成立**：`SelectedExportRows` 调用 `exportHeaders`/`exportRow`（同一函数），非复制 |
| 计划描述符与 provider 一致 | **成立**：初次不一致被 `MODULE_API_MISMATCH` 挡下，已同步；`TestNewMuxProjectsProfileRoutesAndSchemasFromOnePlan` 绿 |
| `users.json` 未新增 capability 声明 | **成立**：diff 仅 `props.selection` 与 custom 节点 |
| 全量回归 | **成立**：Go `go test ./...` 全绿；web `1440/1440`；`tsc -b` exit 0 |

未发现与证据矛盾的陈述。

### Findings

#### F-001 · 双重门禁缺「只有 jobs.write、没有 data.export」的正向反例测试

- 严重度：med
- 建议：**recommended**
- 描述：现有测试覆盖匿名 401 与 editor（两键皆无）403，但**没有**一个「持有 `jobs.write` 而不持有 `data.export`」的角色样例。若将来有人误删第二条 `requirePermission`，测试不会变红——而这正是本阶段最关键的安全口径（异步路径不得成为数据外带旁路）。
- 证据：`TestJobsBatchExportRequiresBothGates`（仅匿名与 editor）；`testsupport/store.go` 的 `jobs.write` 为 `PolicyAdmin`、`data.export` 为 `PolicyAdminEditor`，当前角色集下二者对 admin 同时成立、对 editor 同时不成立，故现有断言无法区分两条门。
- 状态：open（recommended）

#### F-002 · 前端组件的行为缺组件级测试

- 严重度：med
- 建议：**recommended**
- 描述：`jobs-batch-export.tsx` 的关键契约（空选择禁用、只接受 202、轮询至终态、**不调用 `reloadList()`**、终态后下载）目前**无组件测试**；`representative-pages.integration.test.tsx` 只渲染页面、不驱动该组件的交互。回归锚点保护的是**同步**批量路径，不保护本组件。
- 证据：`apps/web/src/components/` 下同族组件（如 `activity-export.test.tsx`、`monitoring-auto-refresh` 相关）均有测试文件，本组件没有；`crud.selection()` 亦无 custom-component 消费者测试（前端侦察 U-8 已提示）。
- 状态：open（recommended）

#### F-003 · `users` 页导出能力与既有全量导出的关系未在 UI 上说明

- 严重度：low
- 建议：**recommended**
- 描述：users 页现同时存在「Export（全量，`export.users` 自定义处理器）」与「导出所选（异步）」两个入口，语义不同但 UI 文案未说明差异；`I-038-013`（文件名/格式一致性）仍为 open。
- 证据：`users.json` toolbar 的 `export`（`export.users` → `GET /api/export/users`）与新 custom 节点；`I-038-013` 状态 open。
- 状态：open（recommended，R4 可一并收敛）

### 必改项汇总（required）

**无。** 未发现 high 级未关闭 required；`I-038-011`/`012` 已 `verified`，`I-038-013` 为 non-blocking（最晚 R4）。

### 结论 + 建议下一步

C1～C3 交付物**如实、可核对、边界干净**：首条真实批量操作确实由 Job 运行时承接（202 + jobId + 真实进度 + 可读结果），既有同步批量路径逐字未动，数据外带面经双重门禁保持不扩大。三条 recommended 均为测试加固与 UI 文案，不阻断。

**verdict = pass**，C4 的 self 腿通过。**建议下一步**：按项目级决策 [independent-audit-execution.md](../../../../architecture/independent-audit-execution.md) 调用本地 grok build（grok 4.6 · high · `/audit`）执行 independent 腿——**须专门核验双重门禁的真实性、进度非硬编码、以及同步 `batch-delete` 与既有 Job 六态合同逐字未退化**。

**本条不修改** `status`、检查点或派生 `progress`。
