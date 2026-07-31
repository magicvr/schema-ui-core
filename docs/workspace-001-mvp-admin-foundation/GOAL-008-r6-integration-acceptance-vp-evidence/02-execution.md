---
id: GOAL-008-r6-integration-acceptance-vp-evidence
doc: execution
status: active
parent: GOAL-001-mvp-admin-foundation
created: 2026-08-01
updated: 2026-08-01
version: 0.3.0
---

# 执行记录 · GOAL-008

## 时间线

### 2026-08-01 · 目标立项与 R6 规划

- 用户调用 `/govern 规划 R6 — 集成验收与 VP 证据`；创建本目标五件套与附件计划，`parent: GOAL-001-mvp-admin-foundation`，状态 `active`。
- [00-meta.md](00-meta.md) 写入四阶段路线图、六条规划成功标准与 `I-008-001`～`I-008-005` required 信息项。
- [01-decision.md](01-decision.md) 记录 D-001 立项与 D-002 规划草案；[R6-acceptance-plan.md](attachments/R6-acceptance-plan.md) v0.1.0 映射 VP 三条退出判据与拟议证据合同。
- 同步父目标 R6 为「规划中」并在工作区 [goal-tree.md](../goal-tree.md) 登记本目标；Root `progress` 保持 `5/6`。

### 2026-08-01 · 规划期能力盘点（非 R6 验收证据）

- 本轮只读/本地能力盘点复跑 package-defined Web tests 与 build：15 个测试文件、395 项测试通过，Vite build 通过；API `go test ./...` 通过。
- 盘点确认现有强输入包括：pinned upstream provenance/SHA、stage3 fixture coverage、Ajv schema validation、R3-R5 范例/集成测试与 Go handler/account tests。
- 盘点未发现仓库 CI workflow、现成浏览器 E2E、JSON/JUnit/coverage reporter 或统一 R6 evidence writer；本次命令输出未按拟议证据合同持久化。
- 因验收合同尚未冻结、revision/environment identity 与原始结果未按 R6 schema 落盘，上述结果**只作为规划输入**，不得计为阶段 2 完成、R6 关门或 VP 关门证据。

### 2026-08-01 · 阶段 1 本地能力基线盘点

- 在工作树 clean、revision `7d20acc7702bcc0e514f787c455bf9c93d5b832f` 上记录环境：Windows/amd64、Node `v22.17.0`、npm `10.9.2`、Go `1.26.0`。
- 从 `apps/web` 执行 `npm test`：15 个测试文件、395 项测试通过；执行 `npm run build`：TypeScript/Vite build 通过。
- 从 `apps/api` 执行 `go test ./...` 与 `go build ./...`：均退出码 0。
- 核对运行入口：`apps/api` 的 `go run ./cmd/server` / `make run` 默认服务 `:8080`，`GET /healthz` 返回健康结果；`apps/web` 的 `npm run dev` 默认 `:5173`，Vite 将 `/api` proxy 到 `http://127.0.0.1:8080`。
- 只读扫描未发现 `.github/workflows`、Playwright/Puppeteer/Cypress/Selenium/WebDriver runner、JSON/JUnit reporter 或统一 evidence writer；该事实只说明当前仓库能力边界，不决定最低验收矩阵。
- 新增候选记录形状：[evidence-index.schema.json](attachments/evidence-index.schema.json) 与 [evidence-index.dry-run.json](attachments/evidence-index.dry-run.json)。dry-run 使用 `mode: planning`、`overallOutcome: blocked`，结果 artifact 状态为 `not-captured`；它没有被计为 R6 验收 evidence。
- 以上动作缩小了 `I-008-002`、`I-008-004`、`I-008-005` 的未知范围；`I-008-001` 仍需冻结验收矩阵，`I-008-003` 仍需冻结 API→Web/Renderer/动作链 oracle。五项 required 尚未闭合。

### 2026-08-01 · 阶段 1 证据收集与最小 CI+浏览器矩阵（用户裁决 D-004）

- **双服务运行时实测（I-008-002）**：在 revision `f3e04f6bd5c1f4ba6b7b72444fd9a0a0ab52d4d5` 上启动 `go run ./cmd/server`（:8080）与 `npm run dev`（:5173），实测 `GET /healthz`（200）、`GET /`（200）、`GET /api/accounts/me`（200，dev-001 session，roles=[admin,editor]，features.beta=true）、`GET /api/records`（200）、`GET /api/records/rec-1`（200）、`GET /api/records?sort=INVALID`（400 拒绝路径）、manifest `.well-known`（200）。结果落盘 `attachments/evidence/planning/results/runtime-probes.log`。
- **Evidence schema dry-run 验证（I-008-004）**：将 web test/build 与 api test/build 输出持久化到 `attachments/evidence/planning/results/{web-test,web-build,api-test,api-build}.log`，新增校验脚本 [validate-evidence-dry-run.mjs](attachments/validate-evidence-dry-run.mjs)（ajv 2020，含 `date-time` format），验证 `evidence-index.schema.json` 可解析且 5 个 artifact SHA-256 可重算；dry-run 持久化为 [evidence-index.dry-run.json](attachments/evidence/planning/evidence-index.dry-run.json)。结论：**SCHEMA_VALIDATION_OK**，但正式 acceptance index 尚未持久化，阶段 1 未冻结。
- **最小 CI + 浏览器矩阵（I-008-005，D-004）**：新增 `.github/workflows/r6-basic-matrix.yml`（web / api / browser-e2e 三 job，Linux + Node 22 + Go 1.26 + Playwright Chromium）；`apps/web` 新增 `@playwright/test` 与 `test:e2e` script、`playwright.config.ts`、`e2e/shell.spec.ts`；`vite.config.ts` 固定 `host:127.0.0.1`。本地 `npm run test:e2e` 通过（shell 渲染、manifest 导航、`/api/accounts/me` proxy 返回 dev session、`/api/records` 非空；截图 `test-results/r6-overview.png`）。`git status` 与 workflow YAML 语法校验通过；**GitHub Actions 实际首跑未发生**（待推送）。
- **验收矩阵与 oracle 升级为冻结候选（I-008-001/003）**：[R6-acceptance-plan.md](attachments/R6-acceptance-plan.md) v0.2.0 新增 §2b 验收矩阵 C-001～C-008 与 §4c 环境矩阵；新增 [account-permission-oracle.md](attachments/account-permission-oracle.md) v0.1.0（正向 P-1～P-4、拒绝 D-1～D-6）。D-002 仍为 proposed，阶段 1 未冻结。
- 本轮后五项 required 均进入「有证据/已决定」的冻结候选状态；阶段 1 正式冻结、阶段 2 开放仍须计划审视（A-002）与用户确认。Root `progress: 5/6`、VP-001 `active` 不变。

### 2026-08-01 · CI 首跑（GitHub Actions · r6-basic-matrix）

- 推送 `f3e04f6..43369fb` 到 `origin/dev` 触发 `.github/workflows/r6-basic-matrix.yml`；run `30666932343` **success**。
- 三 job 全绿：**api**（Linux Go 1.26）22s；**web**（Linux Node 22）27s；**browser E2E**（Linux Node 22，Playwright Chromium）53s。浏览器 E2E 在 CI 上通过 shell 渲染 + `/api` proxy 账号上下文场景，与本地一致。
- 非阻断注解（不影响成功）：`actions/checkout`/`setup-node`/`setup-go` 触发 Node 20 弃用强制跑在 Node 24（GitHub 侧行为）；`setup-go` 缓存因 `apps/api/go.sum` 不存在而 skip（API 无外部依赖，`go.mod` 仅 module + go 版本行，无 requires；`go test ./...` 与 `go build ./...` 仍通过）。
- 意义：`I-008-002` 的干净安装（`npm ci`）与 Linux/CI 等价证据、`I-008-005` 的浏览器矩阵 CI 证据已由**实际首跑**闭合；不再停留在「已配置未跑通」。

### 2026-08-01 · 阶段 1 冻结（A-002 计划审视 + D-002 冻结 + F-008-001 关闭）

- 推送 `43369fb..11d0166`（CI 结果记录）到 `origin/dev`。
- **A-002（self · pass）**：阶段 1 冻结候选计划审视——五项 `I-008` 均 verified，验收矩阵 C-001～C-008、环境矩阵（D-004）、账号权限 oracle、evidence schema dry-run、CI 首跑 green 全部有证据；无开放 required。见 [03-audit.md](03-audit.md)。
- **D-002 冻结为 accepted**：四阶段 R6 计划正式成为执行计划；阶段 2 可开始。
- **F-008-001 关闭（fixed）**：A-001 的必改项满足；响应节已落盘。
- 阶段 1 冻结完成；阶段 2 集成验收执行未开始。Root `progress: 5/6`、VP-001 `active` 不变。

### 2026-08-01 · 阶段 2 集成验收执行

- 在声明 revision `a941bedb1fc2cd4859a408df50653e867da35ff2`（worktree clean）执行验收矩阵 C-001～C-008：
  - **C-001**：web test/build（15 files / 395 tests pass；Vite build pass）
  - **C-002**：api test/build（`go test ./...`、`go build ./...` pass）
  - **C-003**：双服务启动 + health/proxy/账号上下文/records 探测（`GET /healthz` 200、`GET /` 200、`GET /api/accounts/me` 200 dev-001 session、`GET /api/records` 200、`GET /api/records/rec-1` 200、`GET /api/records?sort=INVALID` 400、manifest 200）
  - **C-004**：浏览器 E2E（Playwright Chromium）1 passed / 0 unexpected；截图 `r6-overview.png`
  - **C-005**：账号权限正向（`/api/accounts/me` dev session roles=[admin,editor]）+ D-PERM fixtures（`permissions-inheritance.test.ts` 在 web-test 内）
  - **C-006**：拒绝路径 oracle 依赖（renderer/组件层，见 oracle D-1～D-6；浏览器层不伪造，已列入 exclusion）
  - **C-007**：stage3 conformance（`stage3-fixtures.test.ts` 在 web-test 内，含 request-construction non-batch）
  - **C-008**：机器可读 evidence index 生成并经 schema 校验
- 结果按 `evidence-index.schema.json` 持久化为 **`attachments/evidence/acceptance/evidence-index.json`**（`mode: acceptance`，**7 artifact 全部 SHA-256 verified，overallOutcome=pass**）；原始输出在 `attachments/evidence/acceptance/results/*`。生成脚本 [build-acceptance-index.mjs](attachments/build-acceptance-index.mjs) 可重跑。
- **失败 / 未执行 / 排除显式记录**（不隐藏）：reactions multi-round 16/16（D-008）、request-construction batch 11（D-010 Q1=否）、D-UPLOAD 整域（v0.1.3）、本地非干净安装（`npm ci` 干净安装由 CI run `30667596846` 覆盖）、浏览器级拒绝未断言（真实 manifest 无权限门控项，拒绝以 renderer/组件层断言）——均列入 evidence-index exclusions。
- 阶段 2 全部 required execution item 已运行并落盘；evidence index 可解析、文件摘要可重算。阶段 2 → 阶段 3 门禁审视待做。

### 2026-08-01 · 阶段 2 → 3 门禁通过（A-003）

- **A-003（self · pass）**：阶段 2 退出条件核对——C-001～C-008 全执行并落盘、evidence-index（mode: acceptance）经 ajv 校验（7 SHA-256 verified、overallOutcome=pass）、失败/排除显式记录、无新关键未知。见 [03-audit.md](03-audit.md)。
- 阶段 2 → 3 门禁通过；阶段 3「VP 证据汇编与缺口整改」可开始。Root `progress: 5/6`、VP-001 `active` 不变。

## 待办（计划 · 非完成事实）

1. 阶段 3 VP 证据汇编：三条退出判据各指向 Q2 工作区证据；required 缺口按 P-003 合法闭合。
2. 阶段 3 → 4 门禁审视；完成 R6 close-out 审计后，再由用户决定 Root R6 / `progress` / status。
3. VP 关门另走 `/vision`（读取 R6 工作区证据、形成关门提案并获得用户确认）。

## 进度评估

阶段 1 冻结、阶段 2 执行完成、阶段 2→3 门禁通过（A-003 pass）：C-001～C-008 全执行、正式 `evidence-index.json`（mode: acceptance）经 schema 校验（7 artifact SHA-256 verified、overallOutcome=pass）、排除显式。阶段 3 未开始。Root `progress` 仍 `5/6`，没有 R6 关门完成事实。
