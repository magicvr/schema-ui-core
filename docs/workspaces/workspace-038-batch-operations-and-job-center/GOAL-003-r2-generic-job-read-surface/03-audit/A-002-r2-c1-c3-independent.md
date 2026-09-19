---
id: A-002-r2-c1-c3-independent
doc: audit-entry
parent: GOAL-003-r2-generic-job-read-surface
status: recorded
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# A-002 · R2 C1～C3 独立交叉审计（GOAL-003）

## A-002 · R2 C1～C3 independent（2026-09-19）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high · `/audit`）
- **类型** / **scope**：`stage` · GOAL-003 C1～C3 实施复审（`admin.jobs` 模块接线 / `jobs.read` fail-closed / **既有 `GetForActor` actor 隔离未被放宽** / **wallet 结果 URL 字节等价** / 查询安全 / 索引与冻结断言 / 边界与 pinned 工件）
- **verdict**：**pass**（0 required；3 recommended）
- **完整意见**：本文件

### 范围与区间

- 被审目标：`docs/workspaces/workspace-038-batch-operations-and-job-center/GOAL-003-r2-generic-job-read-surface/`
- 工作区校验：`workspace.md` `id` = `workspace-038-batch-operations-and-job-center`，`root_goal` = `GOAL-001-batch-operations-and-job-center`，`canonical_scope` 与本区路径一致，`vision_role: delivery`，`plan_refs`/`primary_plan` = `VP-038-batch-operations-and-job-center`。`shared_materials_catalog: none`，本意见未把共享资料当事实或关闭证据。未读取其他工作区目标状态；flake 先例仅以 Q2 路径引用既有审计文件。
- 审计区间：R2 五个 checkpoint `d8532c68` → `0624b808` → `c456cbe2` → `aa21c411` → `bd471b3b` 及其后工作区现状（`bd471b3b..HEAD` 无 `apps/api` / `apps/web` / `docs/schemas` / `apps/web/src/protocol/upstream` 后续提交）。
- **不含**：C4 本身、Root R2 投影、R3/R4 实现、愿景层 VRev。
- 约束输入：R1 `GOAL-002/01-decision/D-001-r1-contract-and-denominator-freeze.md`（C1 管理作用域 + `jobs.read`；既有 actor 隔离不动；K-1～K-7；R-1）。

### 范围与区间 · P-005

| ID | 00-meta | 01-decision | 最晚阶段 | 本 scope 判定 |
|----|---------|-------------|----------|----------------|
| I-038-007 | 仍写 `open` | **verified**（D-001 §1，v72 复合索引） | C2 前 | 决策与实施已落地；00-meta 状态未同步（见 F-002） |
| I-038-008 | 仍写 `open` | **verified**（D-001 §2，字节等价重构） | C3 前 | 同上；等价性本条独立复算成立 |
| I-038-009 | 仍写 `open` | **verified**（D-001 §3，R2 只冻方案） | R2 方案冻结前 | 方案已冻；实现属 R3/R4，本 scope 只核 O-2 未声明 `actions.batch.request` |
| I-038-010 | `open` · non-blocking | open | R4 前 | 未到期，不阻断 C1～C3 |

无到期且影响本 scope 的 required 信息项处于「未裁决/未实施」状态。`accepted-residual` 不适用。

### 成果（有证据）

| # | 主张 | 独立核验 |
|---|------|----------|
| 1 | `admin.jobs` 五面贡献与 `Descriptor.Contributions` 逐键一致，并进 admin 默认集 | Provider `Contributions` = Routes `GET /api/jobs`×3、Pages `jobs`、Navigation `menu_jobs`、Permissions `jobs.read`、Fragments `jobs`。`Register()` 对应 `reg.HTTP`（`kernel.RouteKey` → `GET method + " " + pattern`）、`reg.Schema` Key `jobs`、`reg.Authorization` Key `jobs.read`、`reg.Navigation` Key `menu_jobs`、`reg.Manifest` Key `jobs`。`kernel.BuiltinModules()` 与 `profileDefaults[ProfileAdmin]` 描述符同值；mvp/demo 不含 `admin.jobs`。`composition_test.go` admin `wantPermissions: 35` / `wantNavigation: 19`。本会话 `go test ./kernel/ ./internal/composition/` 绿 |
| 2 | kernel 对贡献键 fail-closed 的真实覆盖面 | `stringSetEqual` / `contributionsEqual` 比较的是 **Descriptor ↔ Plan**（`descriptorsMatch`），不是 Descriptor ↔ 实际 `reg.*` 调用。实际 `reg.*` 的 **undeclared key** 由 `validatingRegistrar.declare` → `contributionDeclared` fail-closed（`provider_test.go` `TestRegisterContributionsUndeclaredKey`）。**declared-but-unregistered** 在 `RegisterContributions` 末尾**没有**自动完整性门。对本模块手工 1:1 核对通过，故不升格为对本实施的缺陷 |
| 3 | `jobs.read` 三条路由均先门控，editor 403 | `handler/jobs.go` 三处 `requirePermission(..., "jobs.read")` 且均包在 `a.Middleware`。`requirePermission`：无身份 401 / 无权限 403。`testsupport/store.go` 与 provider 均为 `PolicyID: PolicyAdmin`；`rolesForPolicy(PolicyAdmin)` 只含 `admin`。`TestJobsRoutesGates` 覆盖匿名 401 与 editor 403。`ListJobs`/`GetJob` 的 HTTP 入口仅 `JobsRoutes` |
| 4 | **既有 actor 隔离未被放宽**（重点） | R2 五 commit 对 `internal/jobs/repository.go`、`modules/wallet/jobs.go` 的 `git log -p d8532c68^..bd471b3b` **空 diff**。现行 `GetForActor` 仍为 `WHERE id=? AND kind=? AND actor_id=?`；`RequestCancel`/`Retry` 仍带 `actor_id`。wallet `JobService.Get/Cancel/Retry` 仍走 `GetForActor`；wallet 路由仍 `jobService.Get(..., user.ID)` + `wallet.read`/`wallet.write`。`GetJob` 只被 `admin.jobs` handler/测试使用。冻结测试仍在 `repository_test.go:97-99,218-224` 与 `modules/wallet/jobs_test.go:153-155`。本会话 `go test ./internal/jobs/ ./modules/wallet/` 绿 |
| 5 | **wallet 结果 URL 字节等价**（重点） | 见下方独立复算。`d8532c68` 对 `wallet.go` 的实质 diff 仅：新增 `const WalletJobsBasePath = "/api/wallet/jobs"`，以及 `walletJobToMap` 把 `"/api/wallet/jobs/" + job.ID + "/result"` 换成 `jobs.ResultURL(WalletJobsBasePath, job.ID)`。其余投影字段、wallet 路由/权限/actor predicate **未改** |
| 6 | 查询：白名单字面量、参数化 WHERE、COUNT 同 WHERE、越界空页 | `jobSortSQL` 只返回函数内 `created_at`/`updated_at` + `ASC`/`DESC`。`jobsWhere` 全 `?`。`ListJobs` 一次 `jobsWhere` 供 COUNT 与列表，再 `append` `pageSize` 与 `pagination.Offset`。handler 对未知 sort/非法 order/日期/status/page 先 400。`list_test.go` 含 hostile sort 回退与 page 99 空页。本会话 `go test ./internal/jobs/ ./internal/handler/` 绿 |
| 7 | 索引 v72 与冻结断言；`async_jobs`(42) checksum 未变 | 现场 `Descriptors()`：v42 checksum = `55e1d3f88de080bd0b6015841e76f1ce32604444619d180a3b228123f99dec68`（与 `migrate_test.go:693` 冻结行逐字一致）；独立 `MigrationChecksum(v72 DDL, "0072:jobs-management-indexes:v1")` = `d946e1dbba9311586db8ea29a63d337f38ca7e5927974eae2b64f0f53787cb67`（与冻结行一致）。`completeFingerprintCatalogHead = 72`；`lockedHeadExtraTables[72] = {}`。R2 diff 未改 `jobsDDL`/`jobsPGDDL` 的 CREATE TABLE / 既有三索引。`ApplyPostgres == nil`。本会话 `go test ./internal/store/ ./modules/jobs/migration/` 绿 |
| 8 | R-1：`jobRuntime.enabled` 在 `admin.jobs` **或** `admin.wallet` 下置 true | `composition.go`：wallet 分支与 jobs 分支各自 `enabled.Store(true)`，语义为或。旧行为仅 wallet 分支置 true。**缺**「仅 jobs、不含 wallet」回归测试（F-001） |
| 9 | O-2：新页面未声明 `actions.batch.request` | `modules/jobs/schema/jobs.json` `requiredCapabilities` 仅 `app.manifest` / `app.navigation` / `table.sort`；`modules/jobs/**` 零命中 `actions.batch` |
| 10 | 未触碰 pinned 协议工件 | 五 checkpoint `--stat` 仅 `apps/api/**`、`apps/web/**`（i18n/分母测试）、`docs/workspaces/workspace-038-…/**`。无 `docs/schemas/**`、无 `apps/web/src/protocol/upstream/**` |
| 11 | `jobs.write` 本目标未声明 | 全仓 `jobs.write` 仅 `handler/jobs.go` 注释；符合 D-001「写操作归 R4」 |

### wallet 结果 URL · 独立复算（D-001 §2.4）

历史字面量：`"/api/wallet/jobs/" + job.ID + "/result"`。

现行：

```text
WalletJobsBasePath = "/api/wallet/jobs"          // 无尾斜杠
ResultURL(base, id) = TrimSuffix(base, "/") + "/" + id + "/result"
```

对当前常量，`TrimSuffix("/api/wallet/jobs", "/")` 是恒等，故

```text
ResultURL(WalletJobsBasePath, id) == "/api/wallet/jobs/" + id + "/result"
```

对**任意** `id`（含空串、内嵌 `/`、空白、UUID）两边逐字节相同。本会话用 Go `jobs.ResultURL` 与 PowerShell 拼接两路复算，抽样 id 全部 `equal=true`。

`encoding/json` 对 `map[string]any` **按键名排序**，且 `walletJobToMap` 除 `resultUrl` 推导外字段未改，故同一 job 的 JSON 输出在 `resultUrl` 字符串不变时与历史逐字节相同（D-001 §2.4 判据）。VP-012 `D-002` 的字段存在性、wallet 路由键、`wallet.read` + actor predicate **未被这次重构改写**。

测试钉的是 helper 字面量而非 `WalletJobsBasePath` 常量，且 `wallet_test.go` 无 `resultUrl` 断言（F-003）；**等价性本身由代码与复算成立**，不因此降为 fail。

### 对照检查点

| 检查点 | 状态 | 证据 |
|--------|------|------|
| C1 模块与权限接线 | **达成** | 成果 1/2/8/9/10；`go test ./kernel/ ./internal/composition/` 绿 |
| C2 查询与索引 | **达成** | 成果 6/7；`go test ./internal/jobs/ ./modules/jobs/migration/ ./internal/store/` 绿 |
| C3 读面 API 与作用域 | **达成** | 成果 3/4/5/11；`go test ./internal/handler/ ./modules/wallet/` 绿 |
| C4 审计与投影 | 进行中 | 本条为 independent 腿；不改 status/progress |

### 回归证据可信度

本会话**独立复跑**（均 `-count=1`，exit 0）：

| 包 | 结果 |
|----|------|
| `./internal/store/` | ok 59.9s |
| `./internal/jobs/` | ok 1.7s |
| `./modules/jobs/migration/` | ok 0.5s |
| `./internal/handler/` | ok 46.4s |
| `./modules/wallet/` | ok 0.4s |
| `./internal/composition/` | ok 22.6s（含 drain harness，本跑未 flake） |
| `./kernel/` | ok 0.6s |

**未**在本会话复跑 `apps/api` 全量 `go test ./...`、web vitest 1440、`tsc -b`。E-002 §6 的全量数字因此标为「执行记录自称 + 相关包本会话绿」，不是本条独立全量复证。

`TestShutdownDrainHarnessPostgres` 一次 flake 的留痕：**诚实**。该 harness 在 Q2 `docs/workspaces/workspace-009-production-hardening/GOAL-016-w15-api-web-audit-remediation/03-audit/A-003-w15-s6-independent.md` N-003 与 workspace-022 GOAL-003 A-001 F-003 已多次记为全量并发 EOF flake；隔离复跑绿。R2 五 commit 未改该测试。本会话 composition 包一次通过，与「既有 flake、非本次引入」一致。

### Findings

#### F-001 · R-1 修复缺「仅 admin.jobs、不含 admin.wallet」显式回归测试

- 严重度：med
- 建议：**recommended**
- 描述：`jobRuntime.enabled` 现由两个独立 `HasModule` 分支置 true，语义正确。但测试套件没有一条「plan 含 `admin.jobs`、不含 `admin.wallet` → enabled」的断言（`*_test.go` 对 `HasModule("admin.jobs")` / `jobRuntime.enabled` 零命中）。admin 默认集同时含 wallet，现有 composition 测试无法区分「修复生效」与「旧逻辑碰巧也对」。回退该行不会变红。
- 证据：`internal/composition/composition.go:586-614`；本会话对 `*_test.go` 的检索。
- 状态：open
- 与 self：对应 A-001 F-002，独立复核后**同意**。

#### F-002 · `00-meta` 信息项状态与 `01-decision` 不一致

- 严重度：low
- 建议：**recommended**
- 描述：`I-038-007`/`008`/`009` 在 `01-decision.md` 为 `verified`（D-001 关闭），`00-meta.md` 信息表仍写 `open`。P-005 允许两项之一维护台账，但双源冲突会误导后续编排。不否定 C2/C3 决策与实施已完成。
- 证据：`00-meta.md` 信息表；`01-decision.md` 信息表；`01-decision/D-001-r2-scheme-freeze.md`。
- 状态：open

#### F-003 · wallet `resultUrl` 缺把常量钉在历史字面量上的集成断言

- 严重度：low
- 建议：**recommended**
- 描述：`TestResultURLReproducesTheHistoricalWalletAddress` 使用字面量 `"/api/wallet/jobs"`，不引用 `handler.WalletJobsBasePath`；`wallet_test.go` 对 `resultUrl` 零命中。若有人只改常量、不改 helper，现有测试不会红。当前常量值与 helper 对任意 id 的输出已独立复算为历史字面量，故不构成 D-001 §2.4 失败。
- 证据：`internal/jobs/list_test.go:183-195`；`internal/handler/wallet.go:105,1015`；`grep resultUrl apps/api/internal/handler/wallet_test.go` 无命中。
- 状态：open

### 必改项汇总（required）

**无。** 无 high 级未关闭 required；无到期且影响本 scope、仍未裁决/未实施的 required 信息项。

### 与既有意见的异同（self A-001）

| 项 | self A-001 | 本条 independent |
|----|------------|------------------|
| verdict | pass（0 required，3 recommended） | **pass**（0 required，3 recommended） |
| actor 隔离 | 称 repository 未改、冻结测试绿 | **独立用 git 空 diff + 现行 SQL + wallet 调用链复核，同意** |
| wallet 结果 URL | 称同构 + 专项测试 | **独立代数复算 + 任意 id 抽样 + wallet.go 单行 diff，同意等价成立** |
| 模块接线 / `jobs.read` / 查询 / 索引 checksum | 达成 | 同意；并澄清 `stringSetEqual` 覆盖的是 Descriptor↔Plan，本模块 `reg.*` 另由手工 1:1 + `contributionDeclared` 覆盖 |
| F-001 冻结断言一次同步教训 | recommended | 同意为流程提醒；R2 已修且本会话 store 绿，不另开 finding |
| F-002 R-1 缺测试 | recommended | **同意并落为本条 F-001** |
| F-003 `V-F126` 交接链接 | recommended | 不在 C1～C3 实施正确性门禁内；不纳入本条 findings |
| 00-meta 信息项状态 | 未报 | **新增 F-002** |
| resultUrl 测试钉常量 | 未报 | **新增 F-003** |
| 全量 Go/web/tsc | 称全绿 | 相关包本会话绿；全量与 web/tsc **未独立复跑** |

无与 self 在 required / 结论上的冲突。无需 P-004 裁决。

### 结论 + 建议给编排器/用户的下一步

C1～C3 实施与冻结方案一致：管理读面由 `jobs.read`（`PolicyAdmin`）门控；**`GetForActor` / wallet 路由的 actor 隔离未被放宽**；**wallet `resultUrl` 对任意 job id 与历史字面量逐字节相同**，不构成重开 VP-012。三项 recommended 均为测试/台账加固，不阻断。

**verdict = pass。** 建议用 `/govern` 响应本条与 A-001：可选择补 F-001 回归测试、同步 `00-meta` 信息项状态、把 `WalletJobsBasePath` 钉进测试；然后在开放 required = 0 的前提下投影 Root R2 检查点。

### 声明

本意见不修改 status/progress/goal-tree/方案正文；响应由 `/govern` 处理。
