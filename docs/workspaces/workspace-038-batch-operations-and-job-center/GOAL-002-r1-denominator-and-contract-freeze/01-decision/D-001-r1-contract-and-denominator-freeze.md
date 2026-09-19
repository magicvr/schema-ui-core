---
id: D-001-r1-contract-and-denominator-freeze
doc: decision-entry
parent: GOAL-002-r1-denominator-and-contract-freeze
status: accepted
created: 2026-09-19
updated: 2026-09-19
version: 1.0.0
---

# D-001 · R1 分母与契约冻结

> 本决策关闭 `I-038-001` / `I-038-002` / `I-038-003`（GOAL-002 C1～C3）。证据来源：三份只读侦察报告（位于 Root 目标 `../../GOAL-001-batch-operations-and-job-center/attachments/R1-recon-I-038-00{1,2,3}-*.md`）与 Root `../../GOAL-001-batch-operations-and-job-center/02-execution/E-002-r1-recon.md`；本目标自有的两份冻结矩阵在 `../attachments/r1-*-matrix.md`。
> **用户 P-004 裁决（2026-09-19）**：C2 = **方案 B**；C3 首波 = **仅「新建批量导出所选」**；C1 = **管理作用域 + 新 `jobs.read` 权限**。以下 §1～§3 为该裁决的落盘；§4 记录未选方案；§5 记录**派生项**与**未定项**（后者不冒充已裁决）。

---

## 1. `I-038-002` · 批量异步契约形态冻结（C2）

### 1.1 用户裁决

**方案 B —— 另立本地模块自有异步契约；ADR-0022 同步批量语义完全冻结。**

| 项 | 冻结内容 |
|----|---------|
| 契约归属 | 异步批量操作由新模块 **`admin.jobs`** 自有的**本地端点**承载（202 + `jobId` → 轮询 → 结果），沿用 `admin.wallet` reconcile 已交付的先例形状 |
| ADR-0022 | **不改动**。同步 `POST {path}/batch-delete` 的请求构造（`batchMapping` / `$selection.keys`）、执行路径（`runBatchRequest`）、成功语义（`200 {"deleted": n}` → `reloadList()` + 清选）**逐字保持** |
| 协议 pin | **零改动**：不新增/修改 `docs/schemas/**` 任何 pinned 工件，不新增/修改 `apps/web/src/protocol/upstream/*.cases.json`，不触发 `provenance-v2.9.json` 重签或 `stage3-fixtures.test.ts` 哈希更新 |
| 上游兼容声明 | 本地扩展**不得**写成上游兼容声明（`docs/architecture/module-contribution-playbook.md:84`） |
| 排除 | 不采用「前端识别 202 并转入进度视图」的 ADR-0022 扩展（即方案 A）；不新增异步语义到 `batchRequest` 构造器 |

### 1.2 派生约束（由裁决 + 代码事实推出，供 R2/R3 遵循）

| # | 约束 | 依据 |
|---|------|------|
| K-1 | 异步提交**不得**经由 ADR-0022 `batchMapping` 的成功路径 | `render.tsx:765` 只判 `response.ok`，202 属 ok；随后 `:1244` `reloadList()` 并清空选择——202 的响应体（`jobId`）会被**丢弃**，无法进入进度/结果视图 |
| K-2 | 异步端点返回 `202 Accepted` + job 投影（与 `POST /api/wallet/reconcile` 同形） | `internal/handler/wallet.go:384`；`walletJobToMap` `:989-1006` |
| K-3 | 结果读取沿用三段语义：`409 JOB_RESULT_NOT_READY`（未就绪）/ `410 JOB_RESULT_EXPIRED`（已过期）/ 终态返回结果载荷 | `wallet.go:437-447`；错误码已冻结于 `internal/handler/error_contract_test.go:75` |
| K-4 | 取消/重试沿用 `POST {base}/{id}/cancel`、`POST {base}/{id}/retry`，**写操作须有写权限门** | `wallet.go:400,413`；`wallet_test.go:680-708` |
| K-5 | 新增 Job kind 必须在 runner **启动前**注册，且 kind 不可与既有重复（重复注册会启动失败） | `internal/jobs/runner.go:101-106`；`modules/wallet/jobs.go:39` |
| K-6 | 进度必须由**新 handler 自己**实现细粒度上报 | 唯一先例 reconcile 只报硬编码 10 → 100（`modules/wallet/jobs.go:121-126`）；`reporter.Progress(n)` 语义为 `0..99`（`repository.go:115-122`），终态由 `CompleteWithCommit` 置 100 |
| K-7 | 无需新建迁移：`jobs` 表（迁移 42）已存在，含六态状态机 CHECK 与三个索引 | `modules/jobs/migration/migration.go:14-48`；`migration_test.go:15` |

### 1.3 未定项（**不**冒充已裁决）

| # | 未定项 | 为什么未定 | 归属 |
|---|--------|-----------|------|
| O-1 | **前端触发机制**：批量工具栏如何提交到本地异步端点 | 用户裁决 B（契约归本地）+ C3（在列表页加批量入口）**不唯一确定**该机制。候选：(a) 自定义组件自持按钮与轮询（先例 `components/monitoring-auto-refresh.tsx`）；(b) `props` 本地扩展键（`props` 无 `additionalProperties:false`）+ 渲染器分支 | R2 方案（须在 R2 `01-decision` 冻结） |
| O-2 | **capability 声明口径** | 若新页面声明 `actions.batch.request`，`capability-declaration.guard.test.ts:110-114` 会因缺少 usage marker 直接失败——其现行 marker 是 `/\/batch-delete/`（`:57`）。⇒ 新页面**不应**声明该上游 capability，除非同时登记新 marker 并有据说明其语义 | R2 方案（与 O-1 同批冻结） |
| O-3 | **管理列表索引**：是否为跨 actor 的 `ORDER BY updated_at DESC` 新增索引 | 三个既有索引均为运行期状态机索引，该排序**不被覆盖**（`idx_jobs_actor` 以 `actor_id` 打头）。新增索引须以**新迁移版本号**进行，并同步更新两处冻结断言（`store/migrate_test.go:692-693` checksum、`jobs/migration/migration_test.go:15` Version 42） | R2 方案 |

---

## 2. `I-038-001` · 分母与作用域矩阵冻结（C1）

### 2.1 用户裁决

**管理作用域 + 新增 `jobs.read` 权限。**

| 项 | 冻结内容 |
|----|---------|
| 可见作用域 | `admin.jobs` 提供**跨 actor** 的作业列表与详情（管理作用域），由**新权限 `jobs.read`** 门控（策略 `PolicyAdmin`，对齐 `admin.activity` 的 `operations.read` 先例） |
| 写操作权限 | 取消 / 重试等写操作须有独立写权限（R2 方案确定键名，建议 `jobs.write`），不得仅靠 `jobs.read` |
| **既有 actor 隔离不动** | `GetForActor(id, kind, actorID)` 的三重限定语义**逐字保持**；`admin.wallet` 的 `/api/wallet/jobs/*` 路由与其 actor 隔离测试**不得放宽**。通用读面走**新方法 + 新路由** |
| 冻结的回归面 | `modules/wallet/jobs_test.go:153-155`（他人作业 → `ErrNotFound`）、`internal/jobs/repository_test.go:97-99`、`:218-224` 不得改动语义 |

### 2.2 冻结矩阵（机器可核对）

见 `../attachments/r1-job-kind-scope-matrix.md`。

### 2.3 冻结的读面缺口

`internal/jobs/repository.go` **无任何通用列表方法**（唯一 `ListRunnable` 是 worker 队列扫描，硬编码 queued/running 谓词、无 offset/无 total、永不返回终态）。管理列表所需而**当前缺失**的能力：按 kind / status / actor / 时间范围过滤、分页（page/pageSize + offset）、总数 COUNT、排序选项。

⇒ **冻结结论**：R2 **必须新增 repository 查询方法**；查询形状与索引决策见 §1.3 O-3。

---

## 3. `I-038-003` · 首波分母冻结（C3）

### 3.1 用户裁决

**首波异步操作 = 恰好一条：新建「批量导出所选」。**

| 项 | 冻结内容 |
|----|---------|
| 首波分母 | **1 条**：在列表页新增「批量导出所选」批量操作，选中行经 Job 异步导出 CSV，可在终态读取/下载结果 |
| 操作性质 | **新建**（不是改造既有同步操作）。理由：生产页面 schema 的批量 UI 分母为 **0**（`attachments/R1-recon-I-038-003-*.md` §3.1），既有后端 `batch-delete` 页面未使用；首波需要一条**真实**批量操作来兑现 VP-038 判据 3 |
| 同步 `batch-delete` | **保持同步**。见 §3.3 |
| 结果形态 | Job `result` 承载导出产物（可下载），与 `Result` + `ResultExpiresAt`（`model.go:50,58`）契合 |

### 3.2 冻结矩阵（机器可核对）

见 `../attachments/r1-first-wave-denominator-matrix.md`。

### 3.3 派生项（**由 VP-038 非目标推出，非本轮用户裁决**）

**`users` / `roles` 的同步 `batch-delete` 保持同步。** 依据：VP-038 已把它写成**显式非目标**（`docs/vision/plans/VP-038-batch-operations-and-job-center.md:51`：「把已交付的同步 `batch-delete` 改成 breaking 异步语义」），且侦察确认改为异步属 BREAKING（前端 `onSuccess.behavior="reload"` 依赖响应即终态；`users_batch_test.go:119-184` 断言 409 `LAST_ADMIN` 且**零删除**；协议 fixture 钉住批量构造）。

> **诚实标注**：本轮 P-004 提问中该条为确认项，用户**未**勾选。本决策按 VP-038 显式非目标记为**派生结论**，不冒充用户本轮裁决；若用户意图改变，须回到 VP-038 层修订非目标后再改 R1。

### 3.4 不进首波（本轮用户裁决明确排除）

| 操作 | 位置 | 处置 | 备注 |
|------|------|------|------|
| 回收站 `purge-all` | `recyclebin.go:131` | **不进首波** | 唯一无界 `DELETE`（`repository.go:230`），不可逆；保持同步 |
| CSV 导入 | `import.go:45` | **不进首波** | 保持同步（逐行 no-rollback 语义） |
| **既有 data-transfer 导出** `GET /api/export/{resource}` | `export.go:44` | **不进首波** | 侦察列为长操作（非流式、10000 行上限）；首波只做 W-1「所选批量导出」，既有全量导出保持同步 |
| 操作日志导出 | `operations_export.go:19-24` | **不进首波** | 与 data-transfer 导出同形态同上限；保持同步 |
| `settings` 重置 | `settings.go:41` | **不进首波** | 保持同步 |
| `scheduled-tasks` 手动触发 | `scheduledtasks.go:446` | **不进首波** | 改异步即 BREAKING（现 204，测试与前端依赖） |
| 其余（notifications read-all、代金券批量生成、单目标操作、后台周期任务） | — | **不进首波** | 行数有界或非用户触发 |

以上操作保持**现状**；它们仍是**后续波次的候选**（不构成本 VP 的承诺）。

### 3.5 承接 `V-F126`

Vision Review `V-F126`（`open · recommended`，由 `I-038-003` 承接）的意图是「首波批量操作分母必须明确」。本决策以 §3.1～§3.4 的逐项矩阵给出明确分母（首波 1 条 + 明确排除清单），即已承接。`V-F126` 的闭合登记属愿景层（`/vision`），不在本目标台账内自行改判。

---

## 4. 未选方案（用户裁决的备选）

| 决策点 | 未选 | 未选理由（证据） |
|--------|------|-----------------|
| C2 契约形态 | **方案 A**（扩展 ADR-0022 增加异步变体） | 上游 `action.schema.json` / `node.schema.json` 为 `additionalProperties:false` 严格校验；新增字段 = 改 pinned 工件 = 上游协议变更，须重 pin 并更新 `stage3-fixtures.test.ts` 哈希；且直接改写唯一已交付批量路径（5 个非只读资源共用），回归风险中高 |
| C2 契约形态 | **方案 C**（本地扩展位挂在批量工具栏） | 与 B 的差别仅在前端触发形态；用户选择以 B 的**本地契约**为准，触发形态作为 R2 方案项（§1.3 O-1）另行冻结 |
| C3 首波 | `purge-all` / 导入 / 操作日志导出 改异步 | 三者均为「改造既有同步端点」或「全量而非所选批量」；首波聚焦「真实**所选批量**操作 + 可下载结果」这一条最贴合 VP-038 判据 3 的路径 |
| C1 作用域 | **复用既有权限**（如 `operations.read`） | 语义不精确：作业中心与操作日志是不同产品面，权限键混用会使授权矩阵无法表达「可读作业但不可读操作日志」 |
| C1 作用域 | **仅 actor 作用域** | 会收窄 VP-038 判据 2（「已注册作业可按权限列出」）与判据 5（管理 vs actor 作用域由 R1 冻结）的意图，使「通用作业读面」退化为 wallet 现状 |

---

## 5. 对 R2/R3 的移交清单

| # | 移交项 | 来源 |
|---|--------|------|
| T-1 | `admin.jobs` 模块建立（`kernel.Profile` 内容扩展，进 admin 默认集）+ provider/schema/manifest/navigation/权限接线 | `I-038-004` 用户裁决；Root `D-001` |
| T-2 | 新权限 `jobs.read`（+ 写权限）在 `admin.jobs` 声明并接线；`composition_test.go:529` 的 `wantPermissions/wantNavigation` 计数须同步 | C1 |
| T-3 | 新增 repository 查询方法（kind/status/actor/时间过滤 + 分页 + total + 排序） | §2.3 |
| T-4 | 管理列表索引决策（O-3）与两处冻结断言的同步更新 | §1.3 |
| T-5 | 异步批量导出端点（202 + jobId）+ Job kind + handler（含细粒度 `reporter.Progress`）+ 结果读取 | C2 §1.2、C3 |
| T-6 | 前端触发机制与 capability 声明口径（O-1 / O-2） | §1.3 |
| T-7 | 结果中心页面与体验收敛（R4） | Root 路线图 |
| T-8 | 冻结回归面：同步 `batch-delete` 全部测试、wallet job 路由/权限/错误码测试、provider 路由清单测试 | §1.1、§2.1 |
