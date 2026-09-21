---
id: E-002-r3-implementation
doc: execution-entry
parent: GOAL-004-r3-async-batch-operation
status: recorded
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# E-002 · R3 C1～C3 实施（异步批量导出）

## 事实（2026-09-19）

### 1. 落地产物

**C1 · 异步写面与进度**

| 产物 | 位置 |
|------|------|
| `BatchExportJobKind = "jobs.batch-export"`；`Submit`（先校验后建行）+ `runExport`（**真实细粒度进度**）+ `RenderExportCSV` | `modules/jobs/export.go` |
| `POST /api/jobs/batch-export` → **202 + job 投影**；**双重门禁** `jobs.write` **且** `data.export`；选择键校验复用同步批量的 D3 不变式 | `internal/handler/jobs.go` |
| `ExportHeaders`/`ExportableResources`/`SelectedExportRows` 抽出，使异步与同步导出**共用列序与转义** | `internal/handler/export.go` |
| `BatchExportRowSource`（users/roles 实体适配）+ `BatchExportSubmitter`（模块服务适配） | `internal/handler/jobs_export_source.go` |
| `jobs.write`（`PolicyAdmin`）声明与接线；submitter 为 nil 时不挂路由也不声明该键 | `modules/jobs/provider.go`；`kernel/profile.go` |
| 组合根：runner 启动前注册 kind；按 plan 只接线已启用资源的实体 | `internal/composition/composition.go` |

**C2 · 前端触发**

| 产物 | 位置 |
|------|------|
| `jobs-batch-export` 自定义组件：经 `useSchemaCrud().selection(tableId)` 读选择集；空选择禁用；提交只接受 202 + id；轮询至终态；终态后下载 CSV；**不调用 `reloadList()`** | `apps/web/src/components/jobs-batch-export.tsx` |
| `users-table` 增 `props.selection.mode=multiple`；body section 增 custom 节点 | `modules/users/schema/users.json` |
| i18n 6 键 + `main.tsx` 与 3 个测试夹具的 side-effect import | `i18n/messages/*.json`；`main.tsx`；三个测试文件 |

**C3 · 同步路径回归**

既有同步 `batch-delete` 与其协议 fixture **未改动**；回归锚点全绿（见 §3）。

### 2. 实施中修正的真实契约问题

| # | 问题 | 处置 |
|---|------|------|
| 1 | 新增 R3 路由/权限后，**计划描述符**（`kernel.BuiltinModules`）与 **provider 声明** 不一致 → `MODULE_API_MISMATCH`（`descriptorsMatch` 报 `contribution declaration keys` 不一致） | 同步 `kernel/profile.go` 的 `admin.jobs` 描述符（补 `POST /api/jobs/batch-export` 与 `jobs.write`）。**教训**：provider 的 `Descriptor().Contributions` 与计划描述符必须同步改动，否则装配 fail closed |
| 2 | `internal/handler` 与 `modules/jobs` 之间存在 import cycle，无法在同一测试文件里既用 handler 夹具又用模块服务 | 测试拆两处：HTTP 契约在 handler（用忠实替身 submitter，仍走真 repository）；作业行为在 `modules/jobs`（真 runner + 真 repository + 假行源） |
| 3 | `users.json` 加 `props.selection` 后，多个渲染测试夹具报 `unknown custom component` | 按既有先例在 `main.tsx` 与 3 个测试夹具补 side-effect import |

### 3. 验证证据

| 验证 | 结果 |
|------|------|
| `go build ./...`（apps/api） | exit 0 |
| `go test ./...`（apps/api） | **全绿** |
| `vitest`（apps/web） | **114 files / 1440 tests 全绿** |
| `npm run typecheck`（tsc -b + e2e tsconfig） | exit 0 |
| `capability-declaration.guard` + `representative-pages.integration`（同步批量回归锚点） | **48/48 通过** |

**新增测试覆盖**：

| 测试 | 断言 |
|------|------|
| `TestJobsBatchExportSubmitReturns202` | 202 + 投影（queued 不带 `resultUrl`）；选择去重保序 |
| `TestJobsBatchExportRequiresBothGates` | 匿名 401 / editor 403；被拒请求**不得到达 submitter** |
| `TestJobsBatchExportValidation` | 5 类非法 body → 400 冻结码；不支持资源 → 404；**被拒不留 job 行** |
| `TestJobsRoutesIncludeReadAndSubmit` | R2 读路由不受影响；GET 提交路径不建作业 |
| `TestBatchExportJobProducesCSVResult` | succeeded；result = `{resource,rowCount,fileName,csv}`；CSV 带 BOM 与冻结列序 |
| `TestBatchExportJobReportsRealProgress` | 慢行源下可观察中间进度（至少一个非 0/非 100 采样）——证明不是硬编码 10→100 |
| `TestBatchExportSubmitValidation` | 空/全空白/不支持/超上限（>500）均拒绝且不留行 |

### 4. 边界

- **未**触碰任何 pinned 协议工件（`docs/schemas/**`、`apps/web/src/protocol/upstream/**`）。
- **未**改同步 `batch-delete` 的路由、权限、原子语义或其协议 fixture。
- **未**经 ADR-0022 `batchMapping`/`runBatchRequest` 提交（K-1 落地）。
- **未**在 `users.json` 声明 `table.selection` 或 `actions.batch.request`（O-2 落地）。
- **未**声明或实现 R3 之外的资源导出（仅 users/roles，与同步导出同分母）。
- 结果中心的完整体验（列表/详情/过期/重试/取消）仍归 R4。

### 5. Git checkpoints

| hash | 内容 |
|------|------|
| `e306da93` | R3 后端（模块/端点/权限/组合根） |
| `6d801015` | R3 前端（自定义组件 + users.json + i18n + 夹具） |
| `8e2d0a28` | R3 测试（HTTP 契约 + 作业行为） |
| `25546b17` | R3 方案冻结与前端侦察 |
| `c434e34a` | R3 子目标立项 |

### 6. 未做（移交 C4）

- **未**执行审计：C4 待跑（self + grok build independent）。
- **未**投影 Root R3 检查点。
