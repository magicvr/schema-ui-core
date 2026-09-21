# Grok Build · 独立交叉审计提示词（GOAL-004 R3 C1～C3）

> 执行方式：在仓库根 `schema-ui-core` 以 grok build（模型 grok-4.6 · 思考强度 high）调用 `/audit` 技能（`.grok/skills/audit/SKILL.md`），出具 `source: independent` 的 **R3 实施复审意见**。本文件由编排器预置，供独立会话只读核验。

## 任务

- **目标**：`GOAL-004-r3-async-batch-operation`
- **工作区**：`workspace-038-batch-operations-and-job-center`（canonical `docs/workspaces/workspace-038-batch-operations-and-job-center/`；root_goal `GOAL-001-batch-operations-and-job-center`）
- **scope**：R3 C1～C3 实施复审——① `jobs.write` + `data.export` **双重门禁**是否真的生效（异步路径不得成为既有导出权限的旁路）；② 进度是否**真实细粒度**而非硬编码；③ **同步 `batch-delete` 与既有 Job 六态合同是否逐字未退化**；④ 导出**数据面**是否越界（列集、资源分母、敏感字段）且与同步导出一致；⑤ 提交校验是否**先于建行**（被拒不留 job 行）；⑥ 新 Job kind 注册时机与重复冲突；⑦ 前端是否**未**声明 `actions.batch.request`、是否**未**调用 `reloadList()`、是否正确读取选择集；⑧ 计划描述符与 provider 声明是否一致；⑨ 是否越界改 pinned 工件。
- **只读约束**：不修改 `00-meta` status/progress、不改 goal-tree、不改方案正文；只可追加审计意见（`03-audit/A-NNN` + 索引）或输出文本供代贴。

## 必读文件（按序）

1. `docs/workspaces/workspace-038-batch-operations-and-job-center/workspace.md` 与 `goal-tree.md`
2. `GOAL-004-r3-async-batch-operation/00-meta.md`（4 检查点 / 审计模式 cross / 信息就绪）
3. `GOAL-004-…/01-decision.md` 与 `01-decision/D-001-r3-async-batch-export-freeze.md`（**核心**：§1 数据面与双重门禁 + 未选方案、§2 前端触发与 capability 口径、§3 进度语义）
4. `GOAL-004-…/02-execution/E-001-r3-establishment.md`、`E-002-r3-implementation.md`
5. `GOAL-004-…/attachments/R3-recon-frontend-batch-trigger.md`
6. `GOAL-004-…/03-audit.md` 与 `03-audit/A-001-r3-async-batch-export-self.md`
7. 约束输入：`GOAL-002-…/01-decision/D-001-…`（R1 冻结：C2 方案 B / C3 首波 1 条 / K-1～K-7）与 `GOAL-003-…/01-decision/D-001-…`（R2：O-1/O-2、`jobs.write` 归属、`ResultURL`）

## 重点核验清单（请自行读代码复验，不要只信文档）

- **双重门禁（重点）**：`apps/api/internal/handler/jobs.go` 的 `POST /api/jobs/batch-export` 是否**确实**先 `requirePermission(w, r, "jobs.write")` 再 `requirePermission(w, r, "data.export")`——注意第二个 `requirePermission` 已在第一个通过后执行，若它返回 false 是否真的 `return`（未被忽略）。**独立判断**：是否存在任何路径让只持 `jobs.write` 者导出数据？self A-001 F-001 指出**缺**「只有 jobs.write、没有 data.export」的反例测试，请评估该缺口的真实风险（提示：`testsupport/store.go` 的 `jobs.write` = `PolicyAdmin`、`data.export` = `PolicyAdminEditor`）。
- **进度非硬编码（重点）**：`apps/api/modules/jobs/export.go` 的 `runExport` 是否按 `len(ids)` 真实计算（`done*85/total` 等）而非常量；`reporter.Progress` 的合法区间是 `0..99`（`internal/jobs/repository.go:115-122`）——确认实现**从不**上报 >=100，且终态 100 由运行时置。跑 `cd apps/api && go test ./modules/jobs/ -run TestBatchExport -count=1` 复验。
- **同步路径未退化（重点）**：R3 的 checkpoints 是否**未**改动 `internal/handler/resources.go`、`apps/web/src/renderer/render.tsx` 的批量路径、`apps/web/src/protocol/upstream/request-construction.cases.json`、`apps/web/src/app/representative-pages.integration.test.tsx`。用 `git log -p` 核对区间。并跑 `cd apps/web && npx vitest run src/app/representative-pages.integration.test.tsx src/protocol/capability-declaration.guard.test.ts src/protocol/conformance/stage3-fixtures.test.ts` 复验。
- **Job 六态合同**：`internal/jobs/model.go` 的六态与 `repository.go` 的迁移/summary 是否未改；`modules/jobs/migration/migration.go` 的 v42 `async_jobs` 描述符与 R2 v72 是否未动（v42 checksum 必须仍为 `55e1d3f88de080bd0b6015841e76f1ce32604444619d180a3b228123f99dec68`）。
- **数据面**：`internal/handler/export.go` 的 `exportHeaders` 列集是否与**同步导出**的原内联列集**逐项相同**（users 8 列 / roles 11 列）——即抽出重构没有改变任何一列或顺序；`SelectedExportRows` 是否复用 `exportRow`（含 `formulaSafe` 公式注入中和）。是否有敏感字段（如密码哈希）出现在导出里？
- **校验先于建行**：`BatchExportService.Submit` 是否在任何 `jobs.NewID`/`runner.Submit` 之前完成资源白名单与 `NormalizeExportIDs`（空选择、全空白、>500 上限）；两处测试是否真的断言了 `ListJobs` total = 0。
- **kind 注册时机**：`internal/composition/composition.go` 是否在 `runner.Start()` **之前**构造 `NewBatchExportService`（K-5）；重复注册是否会 fail closed（`runner.go:104-106`）。
- **前端（重点）**：`apps/web/src/components/jobs-batch-export.tsx` 是否只接受 `response.status === 202` 且带 `id`；是否**从不**调用 `crud.refreshList`/`reloadList`；`crud.selection(targetTable)` 是否是选择集的真实 seam（对照 `render.tsx:271,927-936`）；轮询是否在终态停止（`TERMINAL` 集合与 effect 依赖）。
- **capability 口径**：`apps/api/modules/users/schema/users.json` 是否**未**新增 `table.selection` 或 `actions.batch.request`（新增声明会因 marker 不匹配使 guard 变红）；custom 节点与 `props.selection` 是否是本次唯一 schema 变化。
- **描述符一致性**：`apps/api/kernel/profile.go` 的 `admin.jobs` 计划描述符是否与 `modules/jobs/provider.go` 的 `Descriptor().Contributions` **逐键一致**（含 R3 的 `POST /api/jobs/batch-export` 与 `jobs.write`）；`go test ./internal/composition/ -run TestNewMuxProjectsProfileRoutesAndSchemasFromOnePlan` 是否绿。
- **边界**：R3 的 5 个 checkpoint（`c434e34a`、`25546b17`、`e306da93`、`6d801015`、`8e2d0a28`）的 `--stat` 是否只含 `apps/api/**`、`apps/web/**` 与 `docs/workspaces/workspace-038-…/**`，**无** `docs/schemas/**`、`apps/web/src/protocol/upstream/**`。
- **回归证据可信度**：`E-002` §3 声称 Go 全绿 / web 1440 全绿 / tsc exit 0 / 回归锚点 48/48。请独立复跑 Go 全量（`cd apps/api && go test ./... -count=1`）并抽验至少一项 web 断言。

## 输出要求

按 `.grok/skills/audit/SKILL.md` 与 `skills/prompts/05-independent-audit.md` 结构：verdict（pass | conditional | fail，附尺度）、Findings（F-00N；required | recommended；严重度；evidence 路径）、必改项汇总、与 self A-001 的异同、结论与给编排器/用户的下一步。若可写入：追加 `03-audit/A-002-*.md`（source: independent）并更新 `03-audit.md` 索引（不改 status/progress/goal-tree）。
