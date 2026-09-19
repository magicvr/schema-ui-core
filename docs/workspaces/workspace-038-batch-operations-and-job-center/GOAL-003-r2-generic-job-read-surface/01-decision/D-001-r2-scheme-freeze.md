---
id: D-001-r2-scheme-freeze
doc: decision-entry
parent: GOAL-003-r2-generic-job-read-surface
status: accepted
created: 2026-09-19
updated: 2026-09-19
version: 1.0.0
---

# D-001 · R2 方案冻结

> 本决策关闭 `I-038-007` / `I-038-008` / `I-038-009`（GOAL-003 C2/C3/C4 方案项）。
> **用户 P-004 裁决（2026-09-19）**：索引 = **新增单列 `created_at DESC`**；权限策略 = **`jobs.read` / `jobs.write` 均 `PolicyAdmin`**；结果 URL = **采纳泛化 + 共享 helper + 各模块自申 base path**，并**登记为字节等价重构交审计复核**。
> 证据基线：`attachments/R2-recon-migration-mechanics-and-index.md`（编排器独立核对）+ R1 `GOAL-002/01-decision/D-001-…`。

---

## 1. `I-038-007` · 管理列表索引决策（O-3）

### 1.1 用户裁决

**新增单列索引 `idx_jobs_created_at ON jobs(created_at DESC)`**，对齐 `admin.activity` 的 `idx_operation_log_created_at` 先例。不建 `kind`/`status` 复合索引（首波默认列表为全量 + 时间倒序）。

### 1.2 实施规格（冻结）

| 项 | 值 |
|----|-----|
| 新贡献 Version | **72**（当前最大 = 71，`core.operationlog` / `operation_log_digitaloffer_events`；`validateApplied` 要求连续前缀，见侦察 §1.2） |
| ModuleID | **`core.jobs`**（同模块追加后续贡献，先例 `admin.account` 13→35；见侦察 §1.3） |
| Key / Name | `jobs_management_indexes` |
| DDL（sqlite 与 postgres 文本相同） | `CREATE INDEX idx_jobs_created_at ON jobs(created_at DESC)` |
| `ApplyPostgres` | **省略（nil）** —— 纯索引 DDL 无时间列类型差异，`Apply` 可移植（`contribution.go:124-128`） |
| `async_jobs`(42) 行 | **不动**（其 checksum 冻结，改则 checksum drift 拒绝启动） |

**必须同步的测试期望**：

| 测试 | 变更 |
|------|------|
| `internal/store/migrate_test.go:692-693` 冻结列表 | **追加**新三元组 `{"core.jobs", "jobs_management_indexes", <新 checksum>}`（既有行不变） |
| `modules/jobs/migration/migration_test.go:15` | `len(descriptors)` 期望 **1 → 2**；并断言新描述符的 Version/Name |

### 1.3 查询形状（冻结）

- 默认排序 = **`created_at DESC, id DESC`**（对齐 `operationsSortSQL` 的 tie-breaker 约定）。
- 排序白名单 = `{createdAt, updatedAt}`（须经白名单函数映射为列名，禁止拼接用户输入）。
- 过滤 = kind / status / actorId / from / to（均等值或范围，参数化）。
- 分页 = `COUNT(*)` 与列表**同 WHERE**，`LIMIT ? OFFSET ?` 用 `pagination.Offset(page, pageSize, total)`。
- 响应信封 = `resourceList{items,total,page,pageSize}`（`resources.go:253-258`）。

---

## 2. `I-038-008` · 结果 URL 泛化口径

### 2.1 用户裁决

**采纳泛化**：字段 `resultUrl` **保留**（VP-012 `D-002` §6 冻结其存在），但值的**推导机制**改为共享 helper + 各模块自申 base path；**wallet 的输出字符串逐字不变**。

### 2.2 契约边界（已核实，不可越界）

| 约束 | 来源 | 本决策的处置 |
|------|------|-------------|
| `resultUrl` 字段存在（GET Job 返回，不内嵌 payload/result） | VP-012 `D-002` §6（`D-002-r4-precise-contract.md:85`） | **保持** |
| `GET /api/wallet/jobs/{id}/result` 为四个冻结 route key 之一，权限 `wallet.read` + actor predicate | VP-012 `D-002` §1（`:21-25`）、§6（`:82`） | **保持**；wallet 路由与权限零改动 |
| wallet 的 `walletJobToMap` 硬编码 `/api/wallet/jobs/{id}/result` | `handler/wallet.go:1003` | **输出字符串不变**；仅去掉硬编码 |
| 无任何测试钉住该字面量 | `grep resultUrl apps/api` 仅命中 `wallet.go:1003` | 重构不破坏测试 |

### 2.3 机制（冻结）

- `internal/jobs` 提供 `ResultURL(basePath, id string) string`（或等价命名），由 base path 与 job id 拼出结果地址。
- **各模块自申 base path**：`admin.wallet` 申 `/api/wallet/jobs`（输出与现状**逐字一致**）；`admin.jobs` 申 `/api/jobs`。
- **不采用** kind → base path 集中登记表（会引入 kind 与路由的隐式耦合：新 kind 忘登记则 URL 静默错误）。

### 2.4 跨 VP 触碰登记（用户裁决）

对 `apps/api/internal/handler/wallet.go` 的改动**登记为字节等价重构**（byte-identical refactor）：

- **不重开 VP-012**；**不改** `D-002` 契约（字段、路由、权限、actor predicate 全不变）。
- 判定标准：`walletJobToMap` 对同一 job 的输出 **JSON 逐字节相同**（`resultUrl` 字符串不变）。
- **交 cross 审计复核**（C4 的 independent 腿须专门核验该等价性）。

---

## 3. `I-038-009` · 前端触发机制与 capability 声明口径（O-1 / O-2）

### 3.1 冻结（R2 只冻结方案，不实现）

| 项 | 冻结值 | 依据 |
|----|--------|------|
| **O-1 前端触发机制** | **自定义组件**（`registerCustomComponent`）自持提交与轮询按钮，形态照 `components/monitoring-auto-refresh.tsx`（`setInterval` + 定向刷新） | R1 `D-001` §1.3 O-1；`custom-components.ts:21-31`；`monitoring-auto-refresh.tsx:36-38` |
| **O-2 capability 声明** | 新页面**不声明** `actions.batch.request`；改用既有 `table.selection`（marker = `"requiresSelection"`）与 `actions.page.trigger`（marker = `"actionRef"`） | `capability-declaration.guard.test.ts:38-58,110-114`：未登记 marker 的 capability 会直接失败 |
| 轮询形态 | 前端轮询 `GET /api/jobs/{id}` 直至终态；间隔采用既有先例档位（5/10/30s） | R1 `D-001` §1.2 K-2；`monitoring-auto-refresh.tsx:17-22` |
| 202 处理 | 异步提交**不得**走 ADR-0022 `batchMapping` 成功路径（`render.tsx:765` 只判 `response.ok`，会丢弃 jobId） | R1 `D-001` §1.2 K-1 |

### 3.2 实现归属

O-1 / O-2 的**实现**属 R3（批量异步承接）与 R4（结果中心）；本目标只冻结口径。

---

## 4. 其余方案项冻结

| 项 | 冻结值 | 依据 |
|----|--------|------|
| **写权限键名** | `jobs.write` | 用户裁决（`PolicyAdmin`） |
| **策略** | `jobs.read` / `jobs.write` 均 `PolicyID: authsessiondata.PolicyAdmin` | 用户裁决；对齐 `admin.scheduled-tasks` 的 `tasks.read`/`tasks.write` |
| **`jobRuntime.enabled`（R-1）** | 改为 `plan.HasModule("admin.jobs") \|\| plan.HasModule("admin.wallet")` 任一存在即启用 | R1 矩阵 §6 R-1：现行仅在 `admin.wallet` 下置 true（`composition.go:588`），含 `admin.jobs` 不含 `admin.wallet` 的 Profile 下 runner 不启动 |
| **模块路由前缀** | `/api/jobs`（列表 / 详情 / 结果；写操作 `/api/jobs/{id}/cancel`、`/api/jobs/{id}/retry`） | §2.3 base path |
| **页面** | 首波只提供 `jobs` 一个页面（结果中心体验属 R4） | Root 路线图 R2/R4 边界 |
| **导航分组** | `operations` 组（与 `admin.activity` / `admin.scheduled-tasks` 同组） | `scheduledtasks/provider.go:102-114` 先例 |

---

## 5. 未选方案

| 决策点 | 未选 | 理由 |
|--------|------|------|
| 索引 | `created_at` + `kind`/`status` 复合索引 | 首波默认列表无强制过滤；复合索引增加写入开销。若后续出现真实过滤需求，按同模式追加新贡献（版本 73+） |
| 索引 | 不新增索引、改用 `updated_at DESC` 排序 | 跨 actor 时 `idx_jobs_actor` 前导列是 `actor_id`，仍不能覆盖排序；且与 `admin.activity` 的 `created_at` 先例不一致 |
| 权限 | `jobs.read` 用 `PolicyAdminEditor` | 作业含其他 actor 的执行详情（错误消息、correlation），属管理面；与 `tasks.read` 一致更保守 |
| 结果 URL | 集中 kind → base path 登记表 | 隐式耦合：新 kind 忘登记则 URL 静默错误 |
| 结果 URL | 不动 wallet（仅 admin.jobs 自建） | 用户选择泛化；字节等价重构风险已由「输出不变 + 审计复核」约束 |

---

## 6. 对 C1～C4 的映射

| 检查点 | 覆盖 |
|--------|------|
| C1 模块与权限接线 | §4（模块/权限/策略/R-1/路由/页面/导航） |
| C2 查询与索引 | §1（索引规格 + 查询形状 + 测试期望同步） |
| C3 读面 API 与作用域 | §2（结果 URL）+ §1.3（信封/分页/排序白名单） |
| C4 R2 审计与投影 | §3（O-1/O-2 冻结）+ §2.4（跨 VP 触碰登记，交审计复核） |
