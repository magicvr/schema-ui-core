---
id: A-002-r5-closeout-independent
doc: audit-entry
parent: GOAL-006-r5-evidence-and-closeout
status: recorded
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# A-002 · R5 关门独立交叉审计（GOAL-006 / VP-038）

## A-002 · R5 close-out independent（2026-09-19）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high · `/audit`）
- **类型** / **scope**：`close-out` · workspace-038 R5 / VP-038 关门复审（判据 1～7 证据充分性；e2e 首次失败是否既有挂具缺陷；范围外改动 / pinned 工件；R4 三条残余回填 `fixed` 是否名实相符；开放 required / 未登记残余）
- **verdict**：**pass**（0 required；4 recommended）
- **完整意见**：本文件

### 范围与区间

- 被审目标：`docs/workspaces/workspace-038-batch-operations-and-job-center/GOAL-006-r5-evidence-and-closeout/`
- 工作区校验：`workspace.md` `id` = `workspace-038-batch-operations-and-job-center`，`root_goal` = `GOAL-001-batch-operations-and-job-center`，`canonical_scope` 与本区路径一致，`vision_role: delivery`，`plan_refs`/`primary_plan` = `VP-038-batch-operations-and-job-center`。`shared_materials_catalog: none`，本意见未把共享资料当事实或关闭证据。跨区交付仅按任务书以 **Q2 路径**核验 `[workspace-010-design-implementation-conformance] GOAL-044-w32-r4-residual-seams`，未读取该区其它目标状态。
- 审计区间：VP-038 激活 `e125d902` → HEAD `b7b08259`。判据 6 的「pinned 空 diff」按任务书以 R3 关门 `e1893a1a..HEAD` 复验；R2 授权的 v72 索引迁移单独标注。
- **不含**：不改 `00-meta` status/progress、goal-tree、workspace.md、方案正文或实现代码；不宣布 VP-038 `closed`（属 `/vision` + 用户书面确认）。
- 约束输入：任务书 `attachments/grok-prompt-r5-independent.md`；`VP-038` §方向级退出判据 1～7；self `A-001`。

### 范围与区间 · P-005

| ID | 登记处 | 级别 | 最晚阶段 | 本 scope 判定 |
|----|--------|------|----------|----------------|
| I-038-001～003 | Root `00-meta` **verified**；VP-038 信息表仍写 `open` | required | R1 | 实质已由 `GOAL-002` 关闭；VP 表滞后属 C4 投影，不构成本条 required |
| I-038-006 | VP-038 / Root | non-blocking · deferred | 关门后或容量触发 | 有界延期成立，不阻断 |
| I-038-017 | `GOAL-006` `00-meta` 仍 `open` | required | **C2 前** | **问题已被 `E-001` §3 回答**（e2e 不含 jobs 路径）；状态字段未翻 `verified`（见 F-004） |
| I-038-018 | `GOAL-006` `00-meta` `open` | non-blocking | C4 | 未到期 |
| I-038-019 | `GOAL-006` `00-meta` `open` | required | **C4** | 未到期；关门形态仍待用户书面确认 |

无「未收集且到期」的 required 信息。I-038-017 是状态滞后而非信息缺失。

---

### 成果（有证据）

| # | 主张 | 独立复核 |
|---|------|----------|
| 1 | 判据 1：种类×作用域与首波分母是可核对矩阵，排除项点名 | `GOAL-002/attachments/r1-job-kind-scope-matrix.md`、`r1-first-wave-denominator-matrix.md`。生产 kind 抽查：`wallet.reconcile`（`modules/wallet/jobs.go:19` `RegisterWithTerminalHook`）+ `jobs.batch-export`（`modules/jobs/export.go:18`）。排除清单 X-1～X-12 / S-1～S-5 与代码位置对应 |
| 2 | 判据 2：列表/详情/结果 `jobs.read` 门控；`jobToMap` 不嵌入 payload/result；越权 fail-closed | `internal/handler/jobs.go` 三条 GET 均 `requirePermission(..., "jobs.read")`；`jobToMap` 字段集无 payload/result（succeeded 仅派生 `resultUrl`）。本会话 `go test ./internal/handler/ -run 'TestJobs' -count=1` **PASS**（含 `TestJobsListIsManagementScope`、`TestJobsRoutesGates` 匿名 401 / editor 403、`TestJobsDetailAndResult`）；`go test ./internal/jobs/ -count=1` **PASS**（含 `TestGetJobIsManagementScopeButGetForActorIsNot`：`GetForActor` 仍 `WHERE id=? AND kind=? AND actor_id=?`） |
| 3 | 判据 3：至少一条真实异步批量 + 同步 `batch-delete` 不退化 | R3 区间 `c434e34a..e1893a1a` 对 `apps/api/internal/handler/resources.go`、`apps/web/src/protocol/upstream/**`、`apps/web/src/protocol/conformance/request-construction.ts` **空 diff**。`POST /api/jobs/batch-export` 双重门禁 `jobs.write`+`data.export`，202 + jobId。本会话 `npx`/本地 vitest：`capability-declaration.guard.test.ts` 36 例 + `jobs-result-center.test.tsx` 13 例 **PASS** |
| 4 | 判据 4：六态可读、逐值本地化、行操作来自服务端派生字段 | `modules/jobs/schema/jobs.json` status 列 `valueLabels` 覆盖 queued/running/succeeded/failed/cancelled/expired；`jobToMap` 派生 `cancellable`/`retryable`/`downloadable`。本会话 vitest：`jobs-result-center` 13 + `jobs-batch-export` 4 + `jobs-auto-refresh` 4 + `job-result-download` 4 + `table-refresh-seam` 3 = **64/64 PASS**。六态例断言 en-US `Running` / zh-CN `执行中`，未映射值 `quiesced` 原样保留 |
| 5 | 判据 5：`admin.jobs` 仅进 admin 默认集；写路由 `jobs.write`；真反例成立 | `kernel/profile.go`：`admin.jobs` 只出现在 `ProfileAdmin`（mvp/demo 默认集均无）。`jobs.go` cancel/retry 均 `jobs.write`。本会话 `go test ./internal/handler/ -run 'TestJobWriteGate\|TestJobActionsRequireJobsWrite\|TestJobActionRoutesAbsent' -count=1` **PASS**（`TestJobWriteGateRequiresJobsWriteNotJobsRead`：持 `jobs.read` 无 `jobs.write` → 读 200 / 写 403） |
| 6 | 判据 6：pinned 路径与 Job 合同未动；v42 checksum 未变 | `git diff --stat e1893a1a..HEAD -- docs/schemas apps/web/src/protocol/upstream apps/api/modules/jobs/migration apps/api/internal/jobs/repository.go apps/api/internal/jobs/runner.go` **为空**。`async_jobs` v42 checksum 仍为 `55e1d3f88de080bd0b6015841e76f1ce32604444619d180a3b228123f99dec68`（`migrate_test.go:693`）；本会话 `go test ./internal/store/ -run TestCompiledMigrationCatalogOwnership` 与 `go test ./modules/jobs/migration/` **PASS**。从 VP 激活起 `jobs/migration` 仅有 R2 授权的 v72 索引增量，未改 v42 DDL |
| 7 | 判据 7 前置：R1～R4 台账 cross 链开放 required = 0 | `GOAL-002` A-001 pass → A-002 conditional（3 required）→ A-003 全 `fixed`；`GOAL-003`/`GOAL-004`/`GOAL-005` 均为 self pass + independent pass + 响应 开放 required = 0。R4 第 6～8 条回填见下节 |
| 8 | e2e 首次失败 = 既有挂具顺序缺陷，**不是**本 VP 回退 | 见「e2e 根因」 |
| 9 | R4 三条残余回填 `fixed` **名实相符** | 见「回填名实相符」 |

本会话独立复跑（均 exit 0，工作树无残留）：

| 命令 | 结果 |
|------|------|
| `cd apps/api && go test ./internal/handler/ -run 'TestJobs' -count=1` | PASS |
| `cd apps/api && go test ./internal/jobs/ -count=1` | PASS |
| `cd apps/api && go test ./internal/handler/ -run 'TestJobWriteGate\|TestJobActionsRequireJobsWrite\|TestJobActionRoutesAbsent' -count=1` | PASS |
| `cd apps/api && go test ./internal/store/ -run TestCompiledMigrationCatalogOwnership -count=1` | PASS |
| `cd apps/api && go test ./modules/jobs/migration/ -count=1` | PASS |
| `cd apps/web && npm exec -- vitest run`（任务书 6 个文件） | **6 files / 64 tests PASS** |
| `cd apps/web && npx playwright test 00-force-password-change.spec.ts` | **1 passed (10.8s)** |
| `cd apps/web && npm run test:e2e` | **16 passed / 4 skipped / 0 failed（3.9m，exit 0）** |

---

### 对照成功标准（VP-038 判据 1～7）

| # | 标准 | 状态 | 证据 |
|---|------|------|------|
| 1 | 分母与契约冻结，排除项点名 | **达成** | 两份冻结矩阵 + kind 注册点抽查 |
| 2 | 通用作业可见性 fail-closed，不泄露他人作业 | **达成** | 管理作用域读面按设计跨 actor（`jobs.read`=`PolicyAdmin`）；无权限 401/403；不存在 404；actor 路径 `GetForActor` 未放宽；投影不带 payload/result |
| 3 | 至少一条异步批量 + 同步 batch-delete 不退化 | **达成** | `jobs.batch-export`；R3 对 resources.go / 协议 fixture **空 diff** |
| 4 | 结果中心六类呈现可用 | **达成** | 真实 `jobs.json` + 生产渲染链 13 例；六态本地化双语言断言。**浏览器 e2e 未覆盖本面**（F-002，已登记为 recommended） |
| 5 | 权限与 Profile 安全 | **达成** | admin-only 模块；写门禁真反例；只读组合不挂写路由 |
| 6 | 基础设施与范围保持 | **达成** | pinned / runner / repository / v42 checksum 按任务书基线为空或冻结；R4/R5 `apps/` 变更均为 jobs 写面、结果中心、GOAL-044 渲染器 seam、e2e 重命名。未见 Redis/broker/多实例/新业务域 |
| 7 | 退出矩阵 + 回归 + 独立意见 + 开放 required=0 + **用户书面确认** | **进行中** | 矩阵与回归已落盘；本条即独立意见；开放 required = 0（本条）。**用户书面确认与组合投影属 C4，尚未发生** → VP-038 **尚不具备关门条件** |

> 注：`E-001` §1 判据 6 写「R1～R5 区间 migration 空 diff」。从 VP 激活 `e125d902` 起 `apps/api/modules/jobs/migration/**` **有** v72 索引增量（R2 授权）。任务书指定的 `e1893a1a..HEAD` 确为空。这是口径不精确，不是越界改动。

---

### e2e 根因（独立复跑）

**结论：首次失败是既有挂具缺陷，不掩盖本 VP 引入的回退。**

1. **机制成立**：全量套件共用一块 scratch 库（`playwright.config.ts` `e2eDbPath`），`workers: 1`。`sign-in.ts` `signInAsAdmin` 在 fresh seed 上会**自己走完强制改密**（`E2E_INITIAL_PASSWORD` → `E2E_PASSWORD`）。`command-palette.spec.ts`（VP-036）`import { signInAsAdmin }`，按文件名排在原 `force-password-change.spec.ts` 之前，会消费「fresh seed」前提。其后该 spec 用 `admin`/`admin` 登录，失败形态只能是「invalid username or password」——与 `E-001` 记录一致。
2. **两向复跑**：隔离 `00-force-password-change.spec.ts` → **1 passed (10.8s)**；全量 `npm run test:e2e` → **16 passed / 4 skipped / 0 failed（3.9m）**，且 `00-force-password-change` 为第 1 条（2.4s）。
3. **为何不是本 VP 回退**：`e1893a1a..HEAD` 的 `apps/` 变更不触及 auth / 强制改密路径；唯一 e2e 改动是该 spec 的 `git mv` + 注释。登录失败发生在任何业务 API 之前，与 jobs 模块无关。
4. **反证未做**：无需 `git stash` 回 R1 前——机制与两向结果已足够排除「本 VP 引入回退」；若强制改密实现被破坏，隔离跑也会失败。

修复（`00-` 前缀）只让「必须先跑」的注释变成真的，未改断言。同类缺陷仍可能复发（F-001）。

**覆盖诚实边界（加强 self F-002）**：本会话全量 e2e **未设 `APP_PROFILE`**，`playwright.config.ts` 默认 **`mvp`**。`admin.jobs` 不在 mvp/demo 默认集，故即使用例去点 jobs 页也不会挂载。`E-001` §2 写「在 admin profile 上跑通」，与所记录命令 `npm run test:e2e` **不一致**。4 个 skip 与 mvp 条件跳过相符（`localization` admin-only 段 + telegram 三例）。

---

### 回填名实相符（`[workspace-010] GOAL-044` → 038 `A-003` 第 6～8 条）

三项均以**通用能力**落地，不是 jobs 特例；`reloadList()` 的 ADR-0022 D2 **未被削弱**。

| 残余 | 声称 | 独立核对 |
|------|------|----------|
| self F-002 列值本地化 | 列级 `valueLabels`，任意列 | `schema-table.tsx` `labeledCellValue`：badge 列与普通列均走映射；键缺失 `translate(mapped, undefined, raw)` fail-open。jobs 页只是消费者（`jobs.json` 六态映射）。测试含未映射值 + **缺目录键**（`schema.jobs.status.gone` → 仍显示 `running`） |
| self F-003 定向刷新 | `refreshTable(tableId)` 不清空选择 | `render.tsx` `refreshTable` 只 bump token，**不** `setSelections`。`table-refresh-seam.test.tsx` 第 1 例：refresh 后选择仍为 2。组件已改 `crud.refreshTable(targetTable)`，源码无 `reloadList()` 调用 |
| self F-004 空闲不轮询 | `activeStatuses` 可配置 | `jobs-auto-refresh.tsx` 从 props 读数组；`jobs.json` 声明 `["queued","running"]`。缺省 `[]` = 不启用空闲判定。行不可得时保守刷新（组件级测试钉住） |
| D2 未削弱 | `reloadList()` 仍清空全部选择 | `render.tsx:989-993` 仍 `setSelections({})`。`table-refresh-seam.test.tsx` **第 2 例**标题即为 `reloadList still CLEARS every selection (ADR-0022 D2 unchanged)`，本会话该例 PASS |

GOAL-044 自身 3 条 recommended 已由该区 `A-002` 闭合（缺键断言 / 保守分支断言 / 本地扩展登记进 roadmap）。roadmap「未决项统一登记」已把 038 R4 三条标 `fixed`，并另登「本地扩展登记（未做）」——与「实现已交付、清单文档仍缺」一致，不构成名实不符。

---

### Findings

#### F-001 · e2e 顺序假设仍依赖文件名排序

- 严重度：med
- 建议：**recommended**
- 描述：`00-` 前缀让 fresh-seed 用例本次先跑，但契约仍是 lexicographic 隐式约定。后续波次只要再新增更靠前的文件（或改 Playwright 排序），失败形态仍是「登录 401」，与真实缺陷难以区分。
- 证据：`apps/web/e2e/00-force-password-change.spec.ts` 注释；`e2e/sign-in.ts` 自行完成强制改密；`command-palette.spec.ts` 调用 `signInAsAdmin`。
- 状态：**open（recommended）** —— 与 self `A-001` F-001 同向。处置：登记到 roadmap「未决项统一登记」（e2e 挂具顺序契约），由后续符合性波次改为 project 依赖或独立 fresh 库。

#### F-002 · 浏览器回归未驱动 jobs 结果中心；默认 e2e profile 甚至不挂载 `admin.jobs`

- 严重度：med
- 建议：**recommended**
- 描述：11 个 spec **没有**「提交批量导出 → 观察进度 → 下载」路径。判据 4/5 的浏览器侧证据来自 jsdom 交互级测试 + HTTP 契约，不是真实浏览器。独立加强：`npm run test:e2e` 默认 `APP_PROFILE=mvp`，mvp **不含** `admin.jobs`，因此本 VP 新增面在默认 e2e 中**模块级缺席**。`E-001` 「admin profile」与所记命令不符。
- 证据：`apps/web/e2e/**` 清单；`playwright.config.ts:18` 默认 mvp；`kernel/profile.go` mvp/demo 无 `admin.jobs`；本会话全量 e2e 16/4/0 且 localization 走了 mvp 分支。
- 状态：**open（recommended）** —— 与 self `A-001` F-002 同向并加强。不静默扩本波次范围。处置：登记为后续波次（admin profile 下补一条 jobs e2e），或由用户书面接受为有界残余。

#### F-003 · R5 自身两条 recommended 尚未进入 roadmap「未决项统一登记」

- 严重度：low
- 建议：**recommended**
- 描述：roadmap 已覆盖 R4 三条 `fixed` 与「本地扩展登记」。self F-001（e2e 顺序契约）与 F-002（jobs e2e 覆盖）**尚未**出现在该节。C3 检查点要求「残余在 roadmap 登记节可查」。
- 证据：`docs/vision/roadmap.md` §未决项统一登记（2026-09-19 最近更新只列 R4 回填 + 本地扩展登记）；self `A-001` F-001/F-002 状态 open。
- 状态：**open（recommended）** —— 由 `/govern` 在 C3 响应中登记，不必在本波次实现。

#### F-004 · 关门台账/信息项索引滞后（不否定实质证据）

- 严重度：low
- 建议：**recommended**
- 描述：若干索引与状态字段落后于已落盘证据，C4 投影前应一次纠偏，避免「用过期索引放行」：① `GOAL-006/03-audit.md` 在 A-001 已存在时仍写「暂无」（本条写入时一并补索引，不改 status）；② `I-038-017` 仍 `open`，而 `E-001` §3 已回答该 required 问题；③ `GOAL-005/03-audit.md` 说明段仍写三条残余「等用户书面接受」，与 `A-003` 回填 `fixed` 及表格「8/8 fixed」不一致；④ VP-038 信息表 `I-038-001`～`003` 仍 `open`，Root 已 `verified`。
- 证据：上述文件；本条核对日 HEAD `b7b08259`。
- 状态：**open（recommended）** —— 属 C3 响应 / C4 投影，不是新的实现缺口。

### 必改项汇总（required）

**无。**

---

### 与既有意见的异同（self `A-001`）

| 项 | self A-001 | 本条 |
|----|------------|------|
| verdict | pass（0 required + 2 recommended） | **pass**（0 required + 4 recommended） |
| 判据 1～6 证据充分 | 同意 | **独立读代码 + 复跑测试后同意**，非转述 |
| e2e 首次失败 = 挂具缺陷 | 同意 | **同意**；隔离 + 全量复跑 + sign-in 机制核对；未发现本 VP 回退 |
| F-001 顺序契约 / F-002 jobs e2e | recommended | **同向采纳**；F-002 增加「默认 mvp 不挂载 admin.jobs」与「E-001 admin profile 用词无命令支撑」 |
| R4 回填名实 | 引用 GOAL-044 | **独立核对实现 + 第 2 例 D2 对照测试**，名实相符 |
| 新增 | — | F-003 未登记；F-004 索引/信息项滞后 |
| 冲突 | — | **无**。未把 self 任一项升级为 required，未否定其 pass |

未触发 P-004 冲突裁决。

---

### 结论 + 建议给编排器/用户的下一步

**verdict = pass。** 判据 1～6 的证据可独立核对，不是叙事代替证据；e2e 挂具判定成立；R4 三条残余回填名实相符；R1～R4 无开放 required；本条 0 required。

#### VP-038 是否具备关门条件？

**否。** 判据 7 要求「独立意见已落盘 **且** 组合投影同步 **且经用户书面确认**」。本条只满足独立意见落盘与开放 required = 0 的审计侧条件。C4（goal-tree / workspace / Root / VP 投影 + 用户书面确认）尚未发生；`I-038-019` 仍为 C4 的 required 信息项。**不得**把本条 pass 读成 VP-038 `closed`。

建议 `/govern`：

1. **C3 响应本条**：登记 F-001/F-002 到 roadmap「未决项统一登记」；将 `I-038-017` 标 `verified`（证据 = `E-001` §3 + 本条复跑）；纠偏 F-004 所列索引（不改本目标 status/progress，直至 C4）。
2. **C4**：同步投影；按 VP-037 先例准备关门提请文本；**等用户书面确认**后再改 VP-038 status。
3. F-001/F-002 保持 recommended：不在本波次补 Playwright project 依赖、不补 jobs e2e，除非用户当场要求扩范围。

**本条不修改** `status`、检查点或派生 `progress`。
