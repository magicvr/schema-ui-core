# Grok Build · 独立交叉审计提示词（GOAL-003 R2 C1～C3）

> 执行方式：在仓库根 `schema-ui-core` 以 grok build（模型 grok-4.6 · 思考强度 high）调用 `/audit` 技能（`.grok/skills/audit/SKILL.md`），出具 `source: independent` 的 **R2 实施复审意见**。本文件由编排器预置，供独立会话只读核验。

## 任务

- **目标**：`GOAL-003-r2-generic-job-read-surface`
- **工作区**：`workspace-038-batch-operations-and-job-center`（canonical `docs/workspaces/workspace-038-batch-operations-and-job-center/`；root_goal `GOAL-001-batch-operations-and-job-center`）
- **scope**：R2 C1～C3 实施复审——① `admin.jobs` 模块接线是否完整、`Descriptor.Contributions` 与实际 `reg.*` 调用是否逐键一致；② `jobs.read` 是否真正 fail-closed，是否存在越权读取他人作业的路径；③ **既有 `GetForActor` actor 隔离语义与冻结测试是否未被放宽**；④ **wallet 结果 URL 的「字节等价重构」是否真的等价**（`D-001` §2.4 已登记该跨 VP 触碰交本条复核）；⑤ 查询方法的注入/越界/性能缺陷（分页、排序白名单、COUNT 与列表同 WHERE）；⑥ 索引决策（迁移 v72）与冻结断言同步是否正确；⑦ 是否越界改 pinned 工件或重开既有 VP。
- **只读约束**：不修改 `00-meta` status/progress、不改 goal-tree、不改方案正文；只可追加审计意见（`03-audit/A-NNN` + 索引）或输出文本供代贴。

## 必读文件（按序）

1. `docs/workspaces/workspace-038-batch-operations-and-job-center/workspace.md` 与 `goal-tree.md`
2. `GOAL-003-r2-generic-job-read-surface/00-meta.md`（4 检查点 / 审计模式 cross / 信息就绪）
3. `GOAL-003-…/01-decision.md` 与 `01-decision/D-001-r2-scheme-freeze.md`（核心：§1 索引与 §1.2a 细化、§2 结果 URL 与 §2.4 跨 VP 登记、§3 前端口径、§4 其余冻结、§5 未选方案）
4. `GOAL-003-…/02-execution/E-001-r2-establishment.md`、`E-002-r2-c1-c3-implementation.md`
5. `GOAL-003-…/attachments/R2-recon-index-and-query-shape.md`
6. `GOAL-003-…/03-audit.md` 与 `03-audit/A-001-r2-read-surface-self.md`
7. 约束输入：`GOAL-002-r1-denominator-and-contract-freeze/01-decision/D-001-r1-contract-and-denominator-freeze.md`（R1 冻结：C1 管理作用域 + jobs.read；C2 方案 B；K-1～K-7；R-1）

## 重点核验清单（请自行读代码复验，不要只信文档）

- **模块接线**：`apps/api/modules/jobs/provider.go` 的 `Descriptor().Contributions` 是否与实际 `reg.HTTP/Schema/Authorization/Navigation/Manifest` 调用**逐键一致**（`kernel/provider.go` 的 `stringSetEqual` 会 fail closed——请确认它确实覆盖；`kernel.BuiltinModules()` 与 `profileDefaults[ProfileAdmin]` 的描述符是否与 provider 一致）。
- **权限 fail-closed**：`apps/api/internal/handler/jobs.go` 三条路由是否都先 `requirePermission(w, r, "jobs.read")`；`jobs.read` 在 `internal/testsupport/store.go` 的策略是否 `PolicyAdmin`（editor 应当 403）。**是否存在任何未门控路径？**
- **actor 隔离未被放宽（重点）**：`apps/api/internal/jobs/repository.go` 的 `GetForActor`/`RequestCancel`/`Retry` 是否**未被修改**；`modules/wallet/jobs.go` 是否未改；`GetJob`（新方法）是否只被 `admin.jobs` 使用而**没有**被接到 wallet 路由上。请用 `git log -p` 检查这几个文件在 R2 期间的改动。
- **wallet 结果 URL 字节等价（重点）**：`apps/api/internal/handler/wallet.go` 的 `walletJobToMap` 现在经 `jobs.ResultURL(WalletJobsBasePath, job.ID)`。请**独立复算**：对任意 job id，新表达式与历史字面量 `"/api/wallet/jobs/" + job.ID + "/result"` 是否**逐字节相同**（含 `WalletJobsBasePath` 常量值、`ResultURL` 的 `TrimSuffix` 行为）。这是 VP-012 `D-002` §6 契约不被破坏的判据。
- **查询正确性与安全**：`internal/jobs/list.go` 的 `jobSortSQL` 是否只返回函数内字面量（不可注入）；`jobsWhere` 是否全参数化；`ListJobs` 的 COUNT 与列表是否同 WHERE 同 args；分页是否用 `pagination.Offset`；页码越界是否空页而非报错。
- **索引与冻结断言**：`modules/jobs/migration/migration.go` 的 v72 是否**未改动** `jobsDDL`/`jobsPGDDL` 与 `async_jobs`(42) 描述符（其 checksum 必须仍为 `55e1d3f88de080bd0b6015841e76f1ce32604444619d180a3b228123f99dec68`）；`internal/store` 的 6 处冻结断言是否与代码一致；`lockedHeadExtraTables[72]` 是否为空（纯索引无新对象）。可运行 `cd apps/api && go test ./internal/store/ -count=1` 复验。
- **R-1 修复**：`internal/composition/composition.go` 中 `jobRuntime.enabled` 现在是否在 `admin.jobs` **或** `admin.wallet` 任一存在时置 true；旧行为是否确实只在 `admin.wallet` 下置 true（对照 git 历史）。注意：self A-001 F-002 指出该修复**缺显式回归测试**。
- **O-2 口径**：`modules/jobs/schema/jobs.json` 是否**未**声明 `actions.batch.request`（声明了会触发 `capability-declaration.guard.test.ts` 的 marker 缺失失败）。
- **边界**：`git log`/`git show --stat` 复核 R2 的 5 个 checkpoint（`d8532c68`、`0624b808`、`c456cbe2`、`aa21c411`、`bd471b3b`）是否只含 `apps/api/**`、`apps/web/**` 与 `docs/workspaces/workspace-038-…/**`，**无** `docs/schemas/**`、`apps/web/src/protocol/upstream/**` 改动。
- **回归证据可信度**：`E-002` §6 声称 Go 全绿 / web 1440 全绿 / tsc exit 0，并留痕 `TestShutdownDrainHarnessPostgres` 一次 flake。请判断该留痕是否诚实（既有 flake vs 本次引入）。

## 输出要求

按 `.grok/skills/audit/SKILL.md` 与 `skills/prompts/05-independent-audit.md` 结构：verdict（pass | conditional | fail，附尺度）、Findings（F-00N；required | recommended；严重度；evidence 路径）、必改项汇总、与 self A-001 的异同、结论与给编排器/用户的下一步。若可写入：追加 `03-audit/A-002-*.md`（source: independent）并更新 `03-audit.md` 索引（不改 status/progress/goal-tree）。
