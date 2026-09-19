# Grok Build · 独立交叉审计提示词（GOAL-005 R4 C1～C3）

> 执行方式：在仓库根 `schema-ui-core` 以 grok build（模型 grok-4.6 · 思考强度 high）调用 `/audit` 技能（`.grok/skills/audit/SKILL.md`），出具 `source: independent` 的 **R4 实施复审意见**。本文件由编排器预置，供独立会话只读核验。

## 任务

- **目标**：`GOAL-005-r4-result-center-experience`
- **工作区**：`workspace-038-batch-operations-and-job-center`（canonical `docs/workspaces/workspace-038-batch-operations-and-job-center/`；root_goal `GOAL-001-batch-operations-and-job-center`）
- **scope**：R4 C1～C3 实施复审——① 新增的管理作用域写面（`POST /api/jobs/{id}/cancel`、`POST /api/jobs/{id}/retry`）是否真的 `jobs.write` 门控、fail-closed、且不越权；② **既有 actor 作用域写路径（`RequestCancel`/`Retry`）是否逐字未放宽**（隔离语义不得被"新增通用路径"稀释）；③ 可取消/可重试状态集合与错误码是否与既有合同和 Job 六态合同一致（是否偷偷改了合同）；④ 服务端派生字段（`cancellable`/`retryable`/`downloadable`）是否与写合同**同源**而非各自实现；⑤ 前端是否真的打到冻结路由、下载是否单一实现、`jobs-auto-refresh` 是否引入未声明副作用；⑥ 前端交互级测试（含 R3 遗留 4 例）是否**真的能失败**（可判别性），而非恒真断言；⑦ 是否越界改 pinned 工件 / Job 六态 / 同步 `batch-delete` / 既有 VP；⑧ 文档登记（`D-001`/`E-001`/`A-001`/`00-meta`/`goal-tree`）是否与代码事实一致。
- **只读约束**：不修改 `00-meta` status/progress、不改 goal-tree、不改方案正文；只可追加审计意见（`03-audit/A-NNN` + 索引）或输出文本供代贴。

## 必读文件（按序）

1. `docs/workspaces/workspace-038-batch-operations-and-job-center/workspace.md` 与 `goal-tree.md`
2. `GOAL-005-r4-result-center-experience/00-meta.md`（4 检查点 / 审计模式 cross / 信息就绪）
3. `GOAL-005-…/01-decision.md` 与 `01-decision/D-001-r4-structure-and-write-contract-freeze.md`（**核心**：§0 两项用户 P-004 裁决原文 + §1 写面派生合同 + §2 投影派生字段 + §3 前端结构与 download 门控口径 + §5 自动刷新 + §7 残余 R-1）
4. `GOAL-005-…/02-execution.md` 与 `02-execution/E-001-r4-implementation.md`（产物清单、4 项实施中修正的问题、验证矩阵、checkpoint hash `ea6e4006`/`f1351534`/`bafb766f`/`215ebc`）
5. `GOAL-005-…/03-audit.md` 与 `03-audit/A-001-r4-result-center-self.md`（self 腿：0 required + 4 recommended）
6. 约束输入：`GOAL-002-…/01-decision/D-001-…`（R1 冻结：方案 B / 首波 1 条 / K-1～K-7）、`GOAL-003-…/01-decision/D-001-…`（R2：`jobs.read`/`jobs.write` 归属、索引形状）、`GOAL-004-…/01-decision/D-001-…`（R3：双重门禁、进度、前端口径）

## 重点核验清单（请自行读代码复验，不要只信文档）

- **既有隔离零退化（重点）**：`git diff --stat e1893a1a..HEAD -- apps/api/internal/jobs/repository.go apps/api/internal/jobs/runner.go` 是否**为空**；`RequestCancel`/`Retry` 的 `actor_id` 谓词是否原样保留。并**独立判断**：新增 `RequestCancelAny`/`RetryAny` 是否在任何路径上被 actor 作用域入口间接调用（若被间接复用，等于放宽了既有隔离）。
- **写门禁真实性（重点）**：`apps/api/internal/handler/jobs.go` 的 cancel/retry 路由是否 `requirePermission(w, r, "jobs.write")` 且返回 false 时**真的 return**。self A-001 F-001 已承认：现有"editor 403"用例**无法区分** `jobs.write` 与 `jobs.read`（两键同为 `PolicyAdmin`），并以嵌套守卫 + web 可判别对照加固。请**独立评估**该处置是否充分，以及是否还有第三条可构造的判别路径（例如直接检查路由贡献或权限贡献的声明来源）。
- **状态集合与错误码（重点）**：`internal/jobs/actions.go` 的可取消集合是否为 `{queued, running}`、可重试是否为 `failed && attempt < max_attempts`；错误码是否复用冻结的 `JOB_NOT_FOUND`/`JOB_NOT_CANCELLABLE`/`JOB_NOT_RETRYABLE`（**未新增**错误码）；`attempt` 是否**未**被重置；running 取消是否只置 `cancel_requested=1` 而由既有 `FinalizeCancel` 收尾。跑 `cd apps/api && go test ./internal/jobs/ -run 'Any' -count=1 -v` 复验。
- **派生字段同源**：`internal/handler/jobs.go` 的 `jobToMap` 中 `cancellable`/`retryable`/`downloadable` 是否与 `actions.go` 的状态判定**同一条规则**（若两处各写一遍即为漂移风险）；`downloadable = succeeded` 是否与结果路由的三段语义（409/410/200）一致。
- **投影为追加不破坏 R2 冻结面**：新增字段是否**只增不改**（`resultUrl`/`progress`/`error` 等原字段名与含义是否逐字不变）；R2 审计曾判 `pass` 的读面语义是否仍成立。
- **前端真实打点（重点）**：`apps/web/src/renderer/jobs-result-center.test.tsx` 断言的 URL 是否就是冻结路由（`POST /api/jobs/{id}/cancel`、`/retry`、`GET /api/jobs/{id}/result`）；该测试是否**渲染出厂 `jobs.json` 本身**（而非手写替身）；`download` 是否真的走 `lib/job-result-download.ts` 单一实现。**请自行判断这些断言是否可能恒真**（例如：断言"没有 POST"是否因为根本没渲染出行操作）。
- **R3 遗留测试可判别性**：`apps/web/src/components/jobs-batch-export.test.tsx` 的「不调用 `reloadList()`」断言方式是否真的能发现一次 reload（它断的是 `/api/users` 请求次数与选择计数——**若列表请求本就被缓存/合并，断言会不会漏判？**）；「只接受 202」用例是否真的把 200 排除在轮询之外。
- **变异可判别性（重点）**：self 已做三次变异（cancel 状态集合扩宽 → 终态用例变红；移除 `permissionCascade` → 只读主体用例变红；download 门控改 `jobs.write` → 只读主体用例变红）。请**独立复跑至少一次**你自己选的变异（例如移除 handler 的 `jobs.write` 检查、或把 `retryable` 判定改为忽略 `max_attempts`），确认相应测试确实变红，并核对变异后**还原**。
- **自动刷新的副作用**：`components/jobs-auto-refresh.tsx` 是否只调用已声明的 `reloadList`（页面级 seam），是否引入新的全局监听/定时器泄漏（effect 清理是否完备）；`jobs` 表是否真的未声明 `props.selection`（否则自动刷新会清空用户选择）。
- **schema 合法性**：`apps/api/modules/jobs/schema/jobs.json` 是否只用 pinned 允许的属性 + 仓库既有本地扩展（`badgeStyleField`）；`permissionCascade` 是否落在允许的节点类型上且同名 `permissions` 来源存在；`CustomAction` 是否**未**带 `onSuccess`（pinned schema `additionalProperties:false`）。跑 `cd apps/web && npx vitest run src/protocol/all-module-schemas-dval.test.ts src/protocol/capability-declaration.guard.test.ts src/i18n/schema-keys.structural.test.ts` 复验。
- **边界**：R4 的 checkpoint（`ea6e4006`、`f1351534`、`bafb766f`、`ab215ebc`）的 `--stat` 是否只含 `apps/api/**`、`apps/web/**` 与 `docs/workspaces/workspace-038-…/**`，**无** `docs/schemas/**`、`apps/web/src/protocol/upstream/**`、`apps/api/modules/jobs/migration/**`。
- **回归证据可信度**：`E-001` §3 声称 Go 全绿 / web 117 files / 1457 tests 全绿 / `tsc -b` 与 `vite build` exit 0。请独立复跑 `cd apps/api && go test ./... -count=1` 与 `cd apps/web && npx vitest run src/renderer/jobs-result-center.test.tsx src/components/jobs-batch-export.test.tsx src/lib/job-result-download.test.ts`，并抽验至少一项 web 断言。
- **文档一致性**：`00-meta` 的 `progress: 3/4` 是否只由 4 个显式检查点派生且 C4 仍未勾选；`goal-tree.md`/`workspace.md`/Root `00-meta` 的子目标表是否与事实一致（**注意**：Root 仍 `active · 3/5`，R4 检查点**不应**在 C4 通过前被勾选）。

## 输出要求

按 `.grok/skills/audit/SKILL.md` 与 `skills/prompts/05-independent-audit.md` 结构：verdict（pass | conditional | fail，附尺度）、Findings（F-00N；required | recommended；严重度；evidence 路径）、必改项汇总、与 self A-001 的异同、结论与给编排器/用户的下一步。若可写入：追加 `03-audit/A-002-*.md`（source: independent）并更新 `03-audit.md` 索引（不改 status/progress/goal-tree）。
