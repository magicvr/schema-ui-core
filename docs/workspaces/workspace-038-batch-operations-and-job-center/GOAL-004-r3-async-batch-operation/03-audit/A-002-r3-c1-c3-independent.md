---
id: A-002-r3-c1-c3-independent
doc: audit-entry
parent: GOAL-004-r3-async-batch-operation
status: recorded
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# A-002 · R3 C1～C3 独立交叉审计（GOAL-004）

## A-002 · R3 C1～C3 independent（2026-09-19）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high · `/audit`）
- **类型** / **scope**：`stage` · GOAL-004 C1～C3 实施复审——① `jobs.write` + `data.export` **双重门禁**是否真的生效（是否有路径让只持 `jobs.write` 者导出数据）；② 进度是否**真实细粒度**而非硬编码；③ **同步 `batch-delete` 与既有 Job 六态合同是否逐字未退化**；以及数据面/校验先于建行/kind 注册/前端口径/描述符/pinned 边界
- **verdict**：**pass**（0 required；5 recommended）
- **完整意见**：本文件

### 范围与区间

- 被审目标：`docs/workspaces/workspace-038-batch-operations-and-job-center/GOAL-004-r3-async-batch-operation/`
- 工作区校验：`workspace.md` `id` = `workspace-038-batch-operations-and-job-center`，`root_goal` = `GOAL-001-batch-operations-and-job-center`，`canonical_scope` 与本区路径一致，`vision_role: delivery`，`plan_refs`/`primary_plan` = `VP-038-batch-operations-and-job-center`。`shared_materials_catalog: none`，本意见未把共享资料当事实或关闭证据。未读取其他工作区目标状态。
- 审计区间：R3 五个 checkpoint `c434e34a` → `e306da93` → `25546b17` → `6d801015` → `8e2d0a28` 及其后工作区现状（`8e2d0a28..HEAD` 仅工作区文档：`621164ec` 实施事实、`a5af9f06` self 审计；无 `apps/api` / `apps/web` / `docs/schemas` / `apps/web/src/protocol/upstream` 后续提交）。
- **不含**：C4 本身、Root R3 投影、R4 结果中心、愿景层 VRev。
- 约束输入：R1 `GOAL-002/01-decision/D-001-r1-contract-and-denominator-freeze.md`（方案 B / 首波 1 条 / K-1～K-7）；R2 `GOAL-003/01-decision/D-001-r2-scheme-freeze.md`（O-1/O-2、`jobs.write` 归属、`ResultURL`）。
- 执行索引声称 `02-execution/E-001-r3-establishment.md` 存在，磁盘上**无此文件**；本条以 `E-002`、D-001、代码与本会话复跑为准。

### 范围与区间 · P-005

| ID | 00-meta | 01-decision | 最晚阶段 | 本 scope 判定 |
|----|---------|-------------|----------|----------------|
| I-038-011 | 仍写 `open` | **verified**（D-001 §1，users/roles 同分母 + 双重门禁） | C1 前 | 决策与实施已落地；00-meta 状态未同步（见 F-004） |
| I-038-012 | 仍写 `open` | **verified**（D-001 §2，`selection(tableId)` seam） | C2 前 | 同上 |
| I-038-013 | `open` · non-blocking | open | R4 前 | 未到期，不阻断 C1～C3 |

无到期且影响本 scope 的 required 信息项处于「未裁决/未实施」状态。`accepted-residual` 不适用。

### 成果（有证据）

| # | 主张 | 独立核验 |
|---|------|----------|
| 1 | **双重门禁在提交路径真实生效**（重点） | `JobsRoutes` 在 `submitter != nil` 时挂 `POST /api/jobs/batch-export`；handler **先** `requirePermission(..., "jobs.write")`，`!ok` 则 `return`；**再** `requirePermission(..., "data.export")`，`!ok` 则 `return`，其后才 `submitBatchExport`。`requirePermission`：无身份 401 / 权限集不含该键 403。全仓 `SubmitBatchExport` 的生产调用点仅此 handler。本会话 `TestJobsBatchExportRequiresBothGates` 绿（匿名 401、editor 403，且 submitter 未被调用） |
| 2 | **不存在「只持 jobs.write 者导出数据」的现行路径**（重点） | 内置策略：`jobs.write` = `PolicyAdmin` → 仅 `admin`；`data.export` = `PolicyAdminEditor` → `admin`+`editor`。故 **admin 同时持有两键**，**不存在内置「有 jobs.write、无 data.export」主体**。editor 有 `data.export` 无 `jobs.write`，在第一道门被拒，测不到第二道。自定义角色若只被授予 `jobs.write`，仍会被第二道 `return` 挡住。`GET /api/jobs/{id}/result` 只门 `jobs.read`（R2 管理读面）；`jobs.read` 同为 `PolicyAdmin`，内置 admin 仍有 `data.export`。结论：**现行内置矩阵下无 jobs.write-only 旁路**；第二道门的回归锁见 F-001 |
| 3 | **进度真实细粒度，非硬编码 10→100**（重点） | `runExport`：`stepToPercent(done) = done * 85 / total`（按 `len(payload.IDs)`），读完 `Progress(90)`，渲染完 `Progress(99)`。`reporter.Progress` → `UpdateProgress`，合法区间 **0..99**（`repository.go:115-122`）；终态 100 由 `CompleteWithCommit` 的 `progress=100` 写入。handler **从不**上报 ≥100。空选择在除法前被拒，无除零。本会话 `TestBatchExportJobReportsRealProgress` 绿 |
| 4 | **同步 `batch-delete` 与 Job 六态合同未退化**（重点） | `git diff c434e34a^..8e2d0a28` 对 `internal/handler/resources.go`、`apps/web/src/renderer/render.tsx`、`apps/web/src/protocol/upstream/request-construction.cases.json`、`apps/web/src/app/representative-pages.integration.test.tsx`、`internal/jobs/model.go`、`internal/jobs/repository.go`、`modules/jobs/migration/migration.go` **空 diff**。六态仍为 queued/running/succeeded/failed/cancelled/expired。v42 `async_jobs` 冻结 checksum 仍为 `55e1d3f88de080bd0b6015841e76f1ce32604444619d180a3b228123f99dec68`（`migrate_test.go:693`）。本会话 `go test ./modules/jobs/migration/ -count=1` 绿 |
| 5 | 导出数据面与同步导出同源；无密码哈希 | `exportHeaders`：users 8 列 `id,username,name,roles,enabled,locked,createdAt,updatedAt`；roles 11 列 `id,key,name,system,permissions,menuItems,assignedUsers,editable,deletable,createdAt,updatedAt`。`SelectedExportRows` 调 `exportRow`（含 `formulaSafe` 的 `= + - @ TAB CR`）。`userToMap` 含 `mustChangePassword`/`mfaEnabled`/`email`，**未**进入 users 导出列。分母 `ExportableResources` = `users`/`roles`。同步 `export()` **仍内联**同一组表头字面量（见 F-005）；当前逐项相同 |
| 6 | 校验先于建行 | `BatchExportService.Submit`：先 `Supported` + `NormalizeExportIDs`（空/全空白/>500），其后才 `jobs.NewID` / `runner.Submit`。`TestJobsBatchExportValidation` 与 `TestBatchExportSubmitValidation` 均断言 `ListJobs` total = 0。本会话两组测试绿 |
| 7 | K-5：kind 在 runner.Start 前注册；重复 fail-closed | `composition.go` 在装配 providers 时 `NewBatchExportService(jobRuntime.runner, ...)`（`Register`）；`jobs.Start()` 在 `registerLifecycle` 的 `OnStart`、`runtime.Ready` **之后**。`runner.go:104-106`：已启动或 kind 已存在则返回 error。本会话 composition 描述符测试绿 |
| 8 | 前端：只接受 202+id；不调用 `reloadList`；选择集 seam 正确；终态停轮询 | `jobs-batch-export.tsx`：`response.status !== 202 \|\| body?.id === undefined` 则报错；全文无 `reloadList`/`refreshList`/`runBatchRequest`/`batchMapping`。`crud?.selection(targetTable)` 对齐 `render.tsx:271,927-936`。`TERMINAL` = succeeded/failed/cancelled/expired；effect 依赖 `job?.id`/`job?.status`，终态提前 return。`users.json` **未**声明 `table.selection` / `actions.batch.request`；仅 `props.selection.mode=multiple` + custom 节点。轮询间隔与冻结档位不一致（见 F-003）；组件无测试且 representative-pages **未注册**该组件（见 F-002） |
| 9 | 计划描述符与 provider 逐键一致 | `kernel/profile.go` `admin.jobs` Contributions：Routes `GET /api/jobs`、`GET /api/jobs/{id}`、`GET /api/jobs/{id}/result`、`POST /api/jobs/batch-export`；Permissions `jobs.read`、`jobs.write`；Pages/Navigation/Fragments `jobs`/`menu_jobs`/`jobs`。`provider.go` 在 `submitter != nil` 时追加同一 POST 与 `jobs.write`。本会话 `TestNewMuxProjectsProfileRoutesAndSchemasFromOnePlan` 绿 |
| 10 | 未触碰 pinned 协议工件 | 五 checkpoint `--stat` 仅 `apps/api/**`、`apps/web/**`、`docs/workspaces/workspace-038-…/**`。无 `docs/schemas/**`、无 `apps/web/src/protocol/upstream/**` |
| 11 | 回归证据本会话独立复跑 | 见下表。E-002 声称的 Go 全绿与 web 抽验 **本条已复证** |

### 对照检查点

| 检查点 | 状态 | 证据 |
|--------|------|------|
| C1 异步写面与进度 | **达成** | 成果 1～7、9；`TestBatchExport*` / `TestJobsBatchExport*` 本会话绿 |
| C2 前端触发 | **达成** | 成果 8；capability guard 36/36 绿；生产 `main.tsx` 有 side-effect import。测试缺口见 F-002/F-003 |
| C3 同步路径回归 | **达成** | 成果 4/10/11；representative-pages 12/12、stage3-fixtures 278/278、Go `./...` 全绿 |
| C4 审计与投影 | 进行中 | 本条为 independent 腿；不改 status/progress |

### 回归证据可信度

本会话**独立复跑**（均 `-count=1`，exit 0）：

| 命令 | 结果 |
|------|------|
| `cd apps/api && go test ./modules/jobs/ -run TestBatchExport` | ok 0.428s |
| `cd apps/api && go test ./internal/handler/ -run TestJobsBatchExport\|TestJobsRoutesIncludeReadAndSubmit` | ok 0.589s |
| `cd apps/api && go test ./internal/composition/ -run TestNewMuxProjectsProfileRoutesAndSchemasFromOnePlan` | ok 0.500s |
| `cd apps/api && go test ./modules/jobs/migration/` | ok 0.475s |
| `cd apps/api && go test ./...` | **全绿**（含 `internal/store` 64.9s、`internal/handler` 50.6s、`internal/composition` 30.7s、`internal/jobs` 2.0s） |
| `cd apps/web && npx vitest run src/app/representative-pages.integration.test.tsx src/protocol/capability-declaration.guard.test.ts src/protocol/conformance/stage3-fixtures.test.ts` | **3 files / 326 tests 全绿**（guard 36、stage3 278、representative-pages 12）。stderr 见 F-002 |

未在本会话复跑 web 全量 1440 与 `tsc -b`。E-002 的 1440/`tsc` 数字标为「执行记录自称」；本条独立复证的是 Go 全量 + 任务书点名的 web 抽验。

### Findings

#### F-001 · 双重门禁缺「只有 jobs.write、没有 data.export」的区分测试

- 严重度：med
- 建议：**recommended**
- 描述：提交路径的两道 `requirePermission` **均真实 `return`**，现行内置角色矩阵下**没有** jobs.write-only 旁路（见成果 1/2）。但 `TestJobsBatchExportRequiresBothGates` 只覆盖匿名 401 与 editor 403——editor 在**第一道** `jobs.write` 失败，永远走不到 `data.export`。若将来误删第二条 `requirePermission`，现有测试仍绿。这是对本阶段最关键安全口径的**回归锁缺口**，不是现行漏洞。
- 证据：`internal/handler/jobs.go:73-80`；`jobs_export_test.go:113-135`；`testsupport/store.go` 中 `jobs.write`=`PolicyAdmin`、`data.export`=`PolicyAdminEditor`；`modules/authsession/systemdata/policy.go` `rolesForPolicy`
- 状态：open
- 与 self：对应 A-001 F-001，独立复核后**同意**；并补上「现行无旁路、缺口是回归检测」的判定

#### F-002 · 前端组件无测试，且 representative-pages 未注册该组件

- 严重度：med
- 建议：**recommended**
- 描述：`jobs-batch-export.tsx` 的关键契约（空选择禁用、只接受 202、不调用 `reloadList`、终态停轮询、下载）无组件测试。更进一步：`representative-pages.integration.test.tsx` **零命中** `jobs-batch-export` 导入；本会话该文件渲染 users 页时 stderr 为 `[schema-ui] unknown custom component "jobs-batch-export" (node id: users-batch-export)`，12 条测试仍绿。E-002 所称「回归锚点 48/48」保护的是同步批量与 capability 守卫，**不保护**本组件是否被注册、更不保护其交互。生产 `main.tsx` 已 import，故不构成 C2 失败。
- 证据：本会话 vitest stderr；对 `representative-pages.integration.test.tsx` 的 `jobs-batch-export` 检索零命中；`main.tsx:19` 有 import；同族 `activity-export.test.tsx` / `wallet-ensure.test.tsx` 有组件测试
- 状态：open
- 与 self：对应 A-001 F-002，独立复核后**同意并加强**（不仅缺交互测试，集成夹具未注册组件）

#### F-003 · 轮询间隔未遵循 D-001 冻结的 5/10/30s 档位

- 严重度：low
- 建议：**recommended**
- 描述：D-001 §2.4 冻结「轮询间隔沿用既有档位（5/10/30s，先例 `monitoring-auto-refresh.tsx:17-22`）」。实现为常量 `POLL_INTERVAL_MS = 2000`（2s 固定），无 5/10/30 档。C2 成功标准是「轮询至终态」，该行为成立，故不升格 required；但与已冻结实现口径不一致。
- 证据：`apps/web/src/components/jobs-batch-export.tsx:25,90`；`01-decision/D-001-r3-async-batch-export-freeze.md` §2.4 第 4 步；`monitoring-auto-refresh.tsx:17-22`
- 状态：open

#### F-004 · `00-meta` 信息项状态与 `01-decision` 不一致

- 严重度：low
- 建议：**recommended**
- 描述：`I-038-011`/`012` 在 `01-decision.md` 为 `verified`（D-001 关闭），`00-meta.md` 信息表仍写 `open`。P-005 允许两项之一维护台账，但双源冲突会误导后续编排。不否定 C1/C2 决策与实施已完成。
- 证据：`00-meta.md` 信息表；`01-decision.md` 信息表；`01-decision/D-001-r3-async-batch-export-freeze.md`
- 状态：open

#### F-005 · 同步导出路径仍内联表头，未调用抽出的 `exportHeaders`

- 严重度：low
- 建议：**recommended**
- 描述：`exportHeaders` 被声明为「同步与异步共用的唯一列序来源」，但同步 `export()` 仍在 `export.go:213` 与 `:234` 内联同一组字面量，**未**调用 `exportHeaders()`。当前两组逐项相同（users 8 / roles 11），异步 `SelectedExportRows` 已走 `exportHeaders`+`exportRow`。若只改函数、不改内联，两面会静默分叉。不构成「抽出改变了列」——当前列序未变。
- 证据：`internal/handler/export.go:71-79` vs `:213` / `:234`；`e306da93` 对 `export.go` 的 diff 只**新增**函数，未改同步 `export()` 的表头赋值
- 状态：open

### 必改项汇总（required）

**无。** 无 high 级未关闭 required；无到期且影响本 scope、仍未裁决/未实施的 required 信息项。

### 与既有意见的异同（self A-001）

| 项 | self A-001 | 本条 independent |
|----|------------|------------------|
| verdict | pass（0 required，3 recommended） | **pass**（0 required，5 recommended） |
| 双重门禁是否生效 | 称两道 `requirePermission` + editor 403 | **同意生效**；独立判定**无** jobs.write-only 现行旁路；测试缺口同意为 recommended（本条 F-001） |
| 进度非硬编码 | 称 `done*85/total` + 慢行源采样 | **同意**；本会话复跑 `TestBatchExportJobReportsRealProgress`；并核 `UpdateProgress` 拒 ≥100、终态由 `CompleteWithCommit` 置 100 |
| 同步 batch-delete / 六态 | 称 checkpoint 未改、48/48 | **独立用 git 空 diff + v42 checksum 冻结行 + Go 全量 + stage3 278** 复核，同意未退化 |
| 数据面 / 校验先于建行 / K-5 / 描述符 | 达成 | 同意；本会话相关测试绿 |
| F-001 缺 jobs.write-only 反例 | recommended | **同意** → 本条 F-001 |
| F-002 缺组件测试 | recommended | **同意并加强**（representative-pages 未注册组件，本会话 stderr 可见）→ 本条 F-002 |
| F-003 两套导出 UI 文案 | recommended，可交 R4 | **同意**不阻断 C1～C3；不另开 finding（`I-038-013` 最晚 R4） |
| 轮询 2s vs 冻结 5/10/30s | 未报 | **新增 F-003** |
| 00-meta 信息项状态 | 未报 | **新增 F-004** |
| 同步路径未改用 `exportHeaders` | 未报（称抽出后不会漂移） | **新增 F-005**（当前列相同，漂移锁不完整） |
| 全量 Go / web 抽验 | 称 Go 全绿 / 1440 / tsc / 48/48 | **本会话 Go `./...` 全绿**；web 抽验 326/326；1440 与 tsc **未**独立复跑 |

无与 self 在 required / 结论上的冲突。无需 P-004 裁决。

### 结论 + 建议给编排器/用户的下一步

C1～C3 实施与冻结方案一致：**提交路径的双重门禁真实 fail-closed**，现行内置矩阵下只持 `jobs.write` 者**不能**经异步端点导出；进度按选中行线性上报且从不经 `reporter.Progress` 写 ≥100；同步 `batch-delete`、Job 六态与 v42 checksum **逐字未动**。五项 recommended 均为测试/台账/冻结口径加固，不阻断。

**verdict = pass。** 建议用 `/govern` 响应本条与 A-001：可选择补 jobs.write-without-data.export 反例测试、给 `jobs-batch-export` 补组件测试并在 representative-pages 夹具注册该组件、对齐或改写轮询档位冻结、同步 `00-meta` 信息项、让同步 `export()` 调用 `exportHeaders()`；然后在开放 required = 0 的前提下投影 Root R3 检查点。

### 声明

本意见不修改 status/progress/goal-tree/方案正文；响应由 `/govern` 处理。
