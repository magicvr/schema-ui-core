---
id: A-001-r2-read-surface-self
doc: audit-entry
parent: GOAL-003-r2-generic-job-read-surface
status: recorded
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# A-001 · R2 通用作业读面自审（GOAL-003 C1～C3）

## A-001 · R2 C1～C3 自审（2026-09-19）

- **source**：self
- **auditor**：`/govern` 编排器（DeepSeek Harness 会话）
- **类型** / **scope**：`stage` · GOAL-003 C1～C3（`admin.jobs` 模块接线 / 查询与索引 / 读面 API 与作用域；含 wallet 结果 URL 字节等价重构）
- **verdict**：**pass**（0 required；3 recommended）
- **完整意见**：本文件

### 范围与区间

- 被审目标：`docs/workspaces/workspace-038-batch-operations-and-job-center/GOAL-003-r2-generic-job-read-surface/`
- 工作区校验：`root_goal` = `GOAL-001-batch-operations-and-job-center`、`canonical_scope` = 本区路径、`vision_role: delivery`、`plan_refs`/`primary_plan` = `VP-038-…` —— 绑定一致，未跨区。
- 审计区间：2026-09-19（R2 立项 → 方案冻结 → C1～C3 实施 → 验证）。**不含** C4 本身与 R3/R4。
- 审计材料：`00-meta.md`、`01-decision.md` + `01-decision/D-001-…`、`02-execution/E-001/E-002`、`attachments/R2-recon-*`；代码 `apps/api/internal/jobs/list.go`、`internal/handler/jobs.go`、`modules/jobs/**`、`kernel/profile.go`、`internal/composition/composition.go`。

### 成果（有证据）

| # | 成果 | 证据 |
|---|------|------|
| 1 | `admin.jobs` 模块建立并进 admin 默认集，五面贡献齐全（HTTP/Schema/Authorization/Navigation/Manifest），Persistence 空（表归 `core.jobs`） | `modules/jobs/provider.go`；`kernel/profile.go`（默认集 + BuiltinModules） |
| 2 | 新权限 `jobs.read`（`PolicyAdmin`）声明、接线、入 catalog | `provider.go`；`testsupport/store.go`；handler 每条路由 `requirePermission(w, r, "jobs.read")` |
| 3 | 管理作用域读面：跨 actor 列表 + 详情 + 结果，`resourceList` 共享信封 | `internal/handler/jobs.go`；`handler/jobs_test.go`（`TestJobsListIsManagementScope`） |
| 4 | **既有 actor 隔离未被放宽**：`GetForActor` 及其冻结测试逐字不变 | `internal/jobs/repository.go` 未改；`modules/wallet/jobs_test.go:153-155`、`repository_test.go:97-99,218-224` 全绿 |
| 5 | 查询方法：kind/status/actor/时间过滤 + 排序白名单 + COUNT 同 WHERE + `pagination.Offset` | `internal/jobs/list.go`；`list_test.go`（含 hostile sort 回退用例） |
| 6 | 索引决策实施：迁移 v72，`async_jobs`(42) 行**字节不变** | `modules/jobs/migration/migration.go`；重算 checksum 确认 42 行仍为 `55e1d3f8…` |
| 7 | 冻结断言 6 处同步（冻结列表、身份指纹头、`lockedHeadExtraTables`、3 处 applied 尾部） | `identity.go`、`identity_test.go`、`migrate_test.go`、`operations_test.go`、`restart_test.go` |
| 8 | **R-1 修复**：`jobRuntime.enabled` 改为 `admin.jobs` 或 `admin.wallet` 任一存在即启用 | `composition.go`（原仅 `admin.wallet`） |
| 9 | wallet 结果 URL 泛化且**输出逐字不变** | `handler/wallet.go`（`jobs.ResultURL(WalletJobsBasePath, id)`）；`list_test.go` 断言等于历史字面量 |
| 10 | 错误码**复用**既有冻结码，未新造 | `jobs.go` 用 `INVALID_SORT_FIELD`/`INVALID_SORT_ORDER`/`INVALID_DATE_FILTER`/`INVALID_STATUS_FILTER`；`error_contract_test.go` 全绿 |
| 11 | 未触碰 pinned 协议工件 | `git show --stat` 各 checkpoint 无 `docs/schemas/**`、`apps/web/src/protocol/upstream/**` |
| 12 | O-2 口径遵守：新页面**未**声明 `actions.batch.request` | `modules/jobs/schema/jobs.json`（仅 `app.manifest`/`app.navigation`/`table.sort`）；capability guard 75/75 通过 |

### 对照检查点

| 检查点 | 状态 | 证据 |
|--------|------|------|
| C1 模块与权限接线 | **达成** | 成果 1/2/11/12；`go test ./internal/composition/ ./kernel/` 全绿 |
| C2 查询与索引 | **达成** | 成果 5/6/7；`internal/jobs`、`modules/jobs/migration`、`internal/store` 全绿 |
| C3 读面 API 与作用域 | **达成** | 成果 3/4/9/10；`internal/handler` 全绿 |
| C4 审计与投影 | 进行中 | 本条即 self 腿；independent 腿待跑 |

### 独立复核抽查（编排器对自身实现的复验）

| 主张 | 复核结果 |
|------|---------|
| 索引细化后 `async_jobs`(42) checksum 未变 | **成立**：迁移描述符重算输出 `55e1d3f8…` 与冻结值逐字一致 |
| `(created_at DESC, id DESC)` 消除临时 B 树 | **成立**：`sqlite3` EXPLAIN 实测（单列 → `USE TEMP B-TREE`；复合 → 无） |
| 分页需要 tiebreak | **成立**：`created_at` 为 `UnixMilli`；仓内 `operationlog`/`wallet` 先例均带 `, id DESC` |
| wallet 结果地址逐字不变 | **成立**：`ResultURL("/api/wallet/jobs", id)` 与历史 `"/api/wallet/jobs/" + id + "/result"` 同构；专项测试断言 |
| 详情不内嵌 payload/result | **成立**：投影无这两个字段；测试断言其不存在 |
| `jobs.read` 真正 fail-closed | **成立**：匿名 401、editor 403（`PolicyAdmin` 不含 editor）；三条路由全覆盖 |
| 排序白名单防注入 | **成立**：`jobSortSQL` 只返回函数内字面量；未知键与 SQL 形状值均 400 |
| 全量回归 | **成立**：Go `go test ./...` 全绿；web `1440/1440`；`tsc -b` exit 0 |

未发现与证据矛盾的陈述。

### Findings

#### F-001 · R2 实施过程中曾出现「加了迁移贡献但未同步冻结断言」的真实回归

- 严重度：med
- 建议：**recommended**
- 描述：新增 v72 后 `internal/store` 变红（`TestCompleteFingerprintTracksCatalogHead`、`TestMigrateFreshDB`、`TestMigrateExistingV3ToV4`、`TestRestartPersistence`）。根因是更新链不完整——**迁移贡献与 6 处冻结断言必须一次同步**。已修复且全绿；`E-002` §3 已把「必须一次同步 6 处」登记为 R3/R4 可复用的教训。
- 证据：`E-002` §3；`identity.go:93`、`identity_test.go`、`migrate_test.go`、`operations_test.go`、`restart_test.go`。
- 状态：open（recommended；事实性缺口已闭合，保留为流程提醒）

#### F-002 · `jobRuntime.enabled` 的 R-1 修复缺一条显式回归测试

- 严重度：med
- 建议：**recommended**
- 描述：R-1 的修复（`admin.jobs` 存在即启用 runner）是行为变更，但现有测试只覆盖 admin profile（同时含 wallet），**无法区分**「修复生效」与「旧逻辑碰巧也对」。若将来有人回退该行，测试不会变红。
- 证据：`composition.go` 的 `plan.HasModule("admin.jobs")` 分支；`s5_manifest_snapshot_test.go` 等均使用含 wallet 的 admin 集合。
- 状态：open（recommended）

#### F-003 · `V-F126` 交接项与本目标无显式链接

- 严重度：low
- 建议：**recommended**
- 描述：`V-F126` 的交接项登记在 R1 首波矩阵 §7，R2 未引用。R5 关门时需要跨两个目标回看。
- 证据：`GOAL-002/attachments/r1-first-wave-denominator-matrix.md` §7。
- 状态：open（recommended）

### 必改项汇总（required）

**无。** 未发现 high 级未关闭 required，也未发现到期且影响本 scope 的 required 信息项（`I-038-007`～`009` 均已 `verified`；`I-038-010` 为 non-blocking，最晚 R4）。

### 结论 + 建议下一步

C1～C3 交付物**如实、可核对、边界干净**：管理读面按用户裁决落地（`jobs.read` = `PolicyAdmin`），既有 actor 隔离逐字保持，索引决策经真实 EXPLAIN 二次细化并留痕，wallet 结果 URL 泛化保持字节等价。三项 recommended 均为加固/流程项，不阻断。

**verdict = pass**，C4 的 self 腿通过。**建议下一步**：按项目级决策 [independent-audit-execution.md](../../../../architecture/independent-audit-execution.md) 调用本地 grok build（grok 4.6 · high · `/audit`）执行 independent 腿——**须专门核验 wallet 结果 URL 的字节等价性与 actor 隔离未被放宽**（`D-001` §2.4 已登记该触碰交审计复核）。

**本条不修改** `status`、检查点或派生 `progress`。
