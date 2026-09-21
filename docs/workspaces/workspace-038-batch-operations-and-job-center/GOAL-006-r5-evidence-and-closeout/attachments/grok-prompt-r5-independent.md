# Grok Build · 独立交叉审计提示词（GOAL-006 R5 关门审计 · workspace-038）

> 执行方式：在仓库根 `schema-ui-core` 以 grok build（模型 grok-4.6 · 思考强度 high）调用 `/audit` 技能（`.grok/skills/audit/SKILL.md`），出具 `source: independent` 的 **VP-038 关门复审意见**。本文件由编排器预置，供独立会话只读核验。

## 任务

- **目标**：`GOAL-006-r5-evidence-and-closeout`（workspace-038 的最后一个纲领阶段）
- **工作区**：`workspace-038-batch-operations-and-job-center`（canonical `docs/workspaces/workspace-038-batch-operations-and-job-center/`；root_goal `GOAL-001-batch-operations-and-job-center`；primary_plan `VP-038-batch-operations-and-job-center`）
- **scope**：R5 关门复审——① `VP-038` §方向级退出判据 **1～7** 的**证据充分性**（是否用叙事代替证据、是否有「声称达成但无证据」）；② e2e 首次失败被判定为**既有挂具缺陷**是否成立（是否可能掩盖本 VP 引入的真实回退）；③ 是否存在**范围外改动 / pinned 工件触碰 / 未登记变更**；④ R4 三条残余由 `[workspace-010] GOAL-044` 交付后回填 `fixed` 是否**名实相符**（通用能力 vs jobs 特例、`reloadList()` 语义是否真的未被削弱）；⑤ 是否仍有未合法闭合的 required finding 或未登记残余。
- **只读约束**：不修改 `00-meta` status/progress、不改 goal-tree/workspace.md、不改方案正文、不改实现代码；只可追加审计意见（`03-audit/A-NNN` + 索引）或输出文本供代贴。

## 必读文件（按序）

1. `docs/vision/plans/VP-038-batch-operations-and-job-center.md`（**判据基准**：§边界表、§方向级退出判据 1～7、§P-005 信息表）
2. `docs/workspaces/workspace-038-batch-operations-and-job-center/workspace.md` 与 `goal-tree.md`
3. `GOAL-006-r5-evidence-and-closeout/00-meta.md`、`02-execution/E-001-r5-exit-matrix.md`（**核心**）、`03-audit/A-001-r5-closeout-self.md`
4. `GOAL-005-r4-result-center-experience/00-meta.md`、`01-decision/D-001-…`、`02-execution/E-001-…`、`03-audit/A-001/002/003`
5. `GOAL-004-r3-async-batch-operation/01-decision/D-001-…`、`03-audit/A-002-…`；`GOAL-003-r2-generic-job-read-surface/01-decision/D-001-…`、`03-audit/A-002-…`；`GOAL-002-r1-denominator-and-contract-freeze/01-decision/D-001-…`
6. 跨区交付：`docs/workspaces/workspace-010-design-implementation-conformance/GOAL-044-w32-r4-residual-seams/`（`00-meta`、`01-decision/D-001`、`02-execution/E-001`、`03-audit/A-001/A-002`）
7. 登记处：`docs/vision/roadmap.md` §「未决项统一登记」

## 重点核验清单（请自行读代码/跑命令复验，不要只信文档）

- **判据 1（可机器核对）**：`GOAL-002/attachments/r1-job-kind-scope-matrix.md` 与 `r1-first-wave-denominator-matrix.md` 是否真的把「Job 种类 × 作用域」「首波分母」写成可核对矩阵；排除项是否点名。抽查 `apps/api/internal/jobs/*` 与各模块 `Register*JobKind` 消费点，看矩阵是否与代码一致。
- **判据 2（可见性 fail-closed）**：读 `apps/api/internal/handler/jobs.go` 与 `apps/api/internal/jobs/list.go`，确认列表/详情/结果三条路由的门控与作用域；跑 `cd apps/api && go test ./internal/handler/ -run 'TestJobs' -count=1 -v` 与 `go test ./internal/jobs/ -count=1`，**独立判断**是否存在越作用域读取他人作业的路径（例如 `jobToMap` 是否泄露 payload/result）。
- **判据 3（异步承接 + 同步不退化）**：`git log --oneline` 找 R3 checkpoint，用 `git show --stat` 核对是否**未**改 `internal/handler/resources.go`、`apps/web/src/protocol/upstream/**`、`apps/web/src/protocol/conformance/request-construction.ts`；跑 `cd apps/web && npx vitest run src/protocol/capability-declaration.guard.test.ts src/renderer/jobs-result-center.test.tsx`。
- **判据 4（结果中心）**：读 `apps/api/modules/jobs/schema/jobs.json` 与 `apps/web/src/renderer/schema-table.tsx`，确认六态**逐值本地化**（`valueLabels`）与行操作可用性来自服务端派生字段；跑 `cd apps/web && npx vitest run src/renderer/jobs-result-center.test.tsx src/components/jobs-batch-export.test.tsx src/components/jobs-auto-refresh.test.tsx src/lib/job-result-download.test.ts src/renderer/table-refresh-seam.test.tsx`。**独立判断**：是否有「文档声称六态可读但 UI 实际不可读」的情形。
- **判据 5（权限/Profile）**：核对 `admin.jobs` 是否仅进 admin 默认集（`apps/api/kernel/profile.go`）；两条写路由是否 **`jobs.write`** 门控（`internal/handler/jobs.go`）；只读组合下是否不挂载（`TestJobActionRoutesAbsentWithoutActions`）；**真反例**是否成立（`TestJobWriteGateRequiresJobsWriteNotJobsRead`：自定义只读角色持 `jobs.read` 无 `jobs.write` → 读 200 / 写 403）。跑 `cd apps/api && go test ./internal/handler/ -run 'TestJobWriteGate|TestJobActionsRequireJobsWrite|TestJobActionRoutesAbsent' -count=1 -v`。
- **判据 6（范围保持，重点）**：`git diff --stat e1893a1a..HEAD -- docs/schemas apps/web/src/protocol/upstream apps/api/modules/jobs/migration apps/api/internal/jobs/repository.go apps/api/internal/jobs/runner.go` 是否**为空**；`async_jobs` v42 checksum 是否仍为 `55e1d3f88de080bd0b6015841e76f1ce32604444619d180a3b228123f99dec68`（`apps/api/modules/jobs/migration/migration.go`）；同步 `batch-delete` 是否零改动。并**独立判断**是否存在未在上述清单内的越界改动（例如 renderer 里被删掉的既有行为）。
- **判据 7（证据与审计）**：本目标 `03-audit` 台账是否覆盖 R1～R5 全部阶段；每个 A 条目的 findings 是否按 P-003 三路径合法闭合（`fixed`/`accepted-residual`/`user-overruled`），尤其 R4 `A-003` 第 6～8 条从 `accepted-residual` **回填 `fixed`** 是否有 `GOAL-044` 的可核对证据（实现 + 变异 + 测试）；是否仍有开放 required。
- **e2e 根因判定（重点，请自行复跑）**：`apps/web/e2e/00-force-password-change.spec.ts` 的注释声称「共用库 + 文件顺序假设被 VP-036 的 `command-palette.spec.ts` 破坏」。请**独立复跑**验证：① 单跑该 spec 是否通过；② 全量 `cd apps/web && npm run test:e2e` 是否全绿（约 4 分钟，允许按需只跑关键两项）。若你判断该失败**其实掩盖了本 VP 引入的回退**，请给出反证（例如 `git stash` 到 R1 前基线后仍失败 / 或指出具体代码路径）。
- **回填名实相符（重点）**：`GOAL-044` 的三项能力是否**通用**（非 jobs 特例）——`valueLabels` 是否对任意列生效、`refreshTable` 是否对任意表格生效且**不清空选择**、`activeStatuses` 是否可配置；并确认 `reloadList()` 的 ADR-0022 D2 语义**未被削弱**（`renderer/table-refresh-seam.test.tsx` 第 2 例）。
- **残余登记**：`docs/vision/roadmap.md`「未决项统一登记」是否覆盖 R4 三条残余（`fixed`）与「本地扩展登记」；R5 自身的 F-001/F-002（e2e 顺序契约 / jobs e2e 覆盖）是否已登记或需登记。

## 输出要求

按 `.grok/skills/audit/SKILL.md` 与 `skills/prompts/05-independent-audit.md` 结构：verdict（pass | conditional | fail，附尺度）、Findings（F-00N；required | recommended；严重度；evidence 路径）、必改项汇总、与 self `A-001` 的异同、**给用户/编排器的关门建议**（VP-038 是否具备关门条件、缺什么）。若可写入：追加 `03-audit/A-002-*.md`（source: independent）并更新 `03-audit.md` 索引（**不改** status/progress/goal-tree/workspace.md/方案正文/实现代码；若为复跑 e2e 而临时改动，必须还原）。
