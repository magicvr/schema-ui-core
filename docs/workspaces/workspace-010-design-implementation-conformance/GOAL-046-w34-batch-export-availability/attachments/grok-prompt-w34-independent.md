# Grok Build · 独立交叉审计提示词（GOAL-046 W34 · workspace-010）

> 执行方式：在仓库根 `schema-ui-core` 以 grok build（模型 grok-4.6 · 思考强度 high）调用 `/audit` 技能（`.grok/skills/audit/SKILL.md`），出具 `source: independent` 的 **W34 交叉审计意见**。本文件由编排器预置，供独立会话只读核验。

## 任务

- **目标**：`GOAL-046-w34-batch-export-availability`（workspace-010 的一个波次子目标）
- **工作区**：`workspace-010-design-implementation-conformance`（canonical `docs/workspaces/workspace-010-design-implementation-conformance/`；root_goal `GOAL-001-design-implementation-conformance`；primary_plan `VP-010-design-implementation-conformance`）
- **背景（用户报告）**：用户在 VP-038 关门后报告「无论是用户列表页还是角色列表页，选中列表项并点击『导出所选』的时候，都会报错宣示『未找到』，而并不会正常导出。」
- **scope**：① 根因诊断是否成立、是否有**更简单的解释被忽略**（例如数据缺失、权限、代理）；② 修复面是否最小且未越界（是否动了服务端门禁语义 / Job 六态 / 导出契约 / pinned 工件 / VP-038 状态）；③ 「不可用时可见但禁用」的取舍是否被正确论证（与 W33 `D-001` §3 fail-open 冻结取向、以及 page-actions 高度契约的关系）；④ 两层守卫（Go config 守卫 + 组件可用性用例 + e2e 双 profile 契约）是否**真的能捕获**该缺陷，变异验证是否名实相符；⑤ 是否有未合法闭合的 required finding 或未登记残余。
- **只读约束**：不修改 `00-meta` status/progress、不改 goal-tree/workspace.md、不改方案正文、不改实现代码；只可追加审计意见（`03-audit/A-NNN` + 索引）或输出文本供代贴。

## 必读文件（按序）

1. `GOAL-046-w34-batch-export-availability/00-meta.md`、`01-decision/D-001-w34-availability-freeze.md`（**判据基准**：两道根因、fail-open 取舍、守卫设计、授权边界）
2. `GOAL-046-w34-batch-export-availability/02-execution/E-001-w34-implementation.md`（**核心证据**：复现表、产物、变异、全量回归、实施中发现的回归）
3. `GOAL-046-w34-batch-export-availability/03-audit/A-001-w34-self.md`（self 意见，含 3 条 recommended）
4. 变更实现：`apps/api/configs/config.yaml`、`apps/api/internal/config/operator_config_coverage_test.go`、`apps/web/src/components/jobs-batch-export.tsx`、`apps/web/src/components/jobs-batch-export.test.tsx`、`apps/web/e2e/batch-export-availability.spec.ts`、`apps/web/src/i18n/messages/{en-US,zh-CN}.json`
5. 上游缺陷面：`apps/api/kernel/profile.go`（`ProfileMVP` / `ProfileAdmin` 集合）、`apps/api/modules/jobs/provider.go`（路由与权限贡献）、`apps/api/internal/handler/jobs.go`（两道门禁与 404 码）、`apps/api/internal/handler/route_envelope.go`（未挂载路由的 JSON 回落）
6. 冻结取向来源：`docs/workspaces/workspace-010-design-implementation-conformance/GOAL-045-w33-list-actions-slot-and-roles-trigger/01-decision/D-001-w33-slot-and-roles-trigger-freeze.md` §3
7. 登记处：`docs/vision/roadmap.md`（「未决项统一登记」/「最近更新」）

## 重点核验清单（请自行读代码/跑命令复验，不要只信文档）

- **根因是否成立（重点）**：`configs/config.yaml` 的 `app.modules.list` 是否真的缺 `admin.jobs`（修复前）；`apps/api/kernel/profile.go` 的 admin preset 是否含 `admin.jobs` 而 mvp/demo 不含；`apps/api/internal/handler/jobs.go` 的 `POST /api/jobs/batch-export` 是否只在 `submitter != nil` 时挂载（`modules/jobs/provider.go` 的 `Descriptor`）；未挂载时 `route_envelope.go` 是否产出 `NOT_FOUND`/`error.notFound`（zh = 「未找到」）。**独立判断**是否存在别的更可能解释（例如 `data.export` 缺失、代理、`targetTable` 选择集为空）。
- **修复是否越界**：`git show --stat e1013c8b` 是否只覆盖 `D-001` §6 列出的路径；是否**未**改 `apps/api/internal/handler/jobs.go`、`apps/api/modules/jobs/export.go`、`apps/web/src/protocol/upstream/**`、`docs/schemas/**`、VP-038 文件 `status`、workspace-038 台账正文。跑 `git diff --stat 784e9152..e1013c8b`。
- **可用性判据是否等价于「路由存在且被授权」**：核对 `jobs.write` 是否**只**由 `admin.jobs` 贡献、`data.export` 是否**只**由 `admin.data-transfer` 贡献（`grep -rn '"jobs.write"\|"data.export"' apps/api --include=*.go`，排除测试）；`D-001` §3 的等价性论证是否成立。**独立判断**：是否会出现「两权限齐备但路由仍未挂载」或「路由挂载但两权限不齐备」的组合，从而使客户端判据误判。
- **fail-open 取舍（重点，可能与 self 有分歧）**：读 W33 `D-001` §3 原文，判断「可见 + disabled + 说明」是否确实优于「隐藏」；并核对若改为隐藏是否会破坏 page-actions 高度契约（`apps/web/e2e/list-visual-surface.spec.ts:82` 的断言）。跑 `cd apps/web && npx playwright test e2e/list-visual-surface.spec.ts --reporter=line`（约 1 分钟）。
- **守卫是否真的有效**：读 `operator_config_coverage_test.go` 的扫描逻辑，判断「超集」判据是否会在**其它**漂移形态下失效（例如改用 `preset:`、把 list 写成内联数组 `[a, b]`、或注释掉模块行）。跑 `cd apps/api && go test ./internal/config/ -run TestOperatorConfigCoversAdminPreset -count=1 -v`；如条件允许，自行制造一次漂移确认其失败（**必须还原**）。
- **组件用例**：跑 `cd apps/web && npx vitest run src/components/jobs-batch-export.test.tsx src/components/jobs-batch-export-roles.test.tsx src/renderer/list-actions-slot.test.tsx`；判断 4 例可用性用例是否覆盖了「门禁缺一」「两门禁齐备」「上下文未知 fail-open」「404 文案」四种形态，以及是否存在**应由服务端门禁保证却被客户端断言越权替代**的风险（即客户端判据被误当成授权边界）。
- **e2e 契约**：跑 `cd apps/web && npx playwright test e2e/batch-export-availability.spec.ts --reporter=line`（mvp 分支：可见但禁用）与 `APP_PROFILE=admin npx playwright test e2e/batch-export-availability.spec.ts --reporter=line`（admin 分支：可用）。**独立判断**该 spec 的 `adminCapable` 推断（`admin`/`custom`）是否有误判面（例如 harness 的 custom 列表变化）。
- **回归数字**：核对 `E-001` §4 的数字（vitest 121 files / 1481 tests、e2e mvp 17/5/0、admin 18/4/0、`go test ./...` 全绿）。至少独立复跑 `cd apps/web && npx vitest run` 与 `cd apps/api && go test ./internal/config/ ./modules/jobs/ -count=1`。
- **实施中发现的回归是否如实登记**：`E-001` §5 声称第一版「隐藏」实现被 `list-visual-surface` e2e 捕获（`got 0, 32, 32, 32`）。请判断该叙述是否可信、是否应升级为 finding（而非仅「成果」），以及修正后的形态是否留下空插槽宿主。
- **残余登记**：`docs/vision/roadmap.md` 是否登记了本缺陷（`[workspace-038]` 相关行与「最近更新」）；self `A-001` 的 `F-001`（手写 YAML 扫描）/`F-002`（客户端镜像漂移面）是否有界且关闭要求明确。

## 输出要求

按 `.grok/skills/audit/SKILL.md` 与 `skills/prompts/05-independent-audit.md` 结构：verdict（pass | conditional | fail，附尺度）、Findings（F-00N；required | recommended；严重度；evidence 路径）、必改项汇总、与 self `A-001` 的异同、**给用户/编排器的关门建议**（W34 是否具备关门条件、缺什么）。若可写入：追加 `03-audit/A-002-*.md`（source: independent）并更新 `03-audit.md` 索引（**不改** status/progress/goal-tree/workspace.md/方案正文/实现代码；若为复跑而临时改动，必须还原）。
