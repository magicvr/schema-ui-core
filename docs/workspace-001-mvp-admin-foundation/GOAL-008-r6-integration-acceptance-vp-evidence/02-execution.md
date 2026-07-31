---
id: GOAL-008-r6-integration-acceptance-vp-evidence
doc: execution
status: active
parent: GOAL-001-mvp-admin-foundation
created: 2026-08-01
updated: 2026-08-01
version: 0.2.0
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

## 待办（计划 · 非完成事实）

1. 阶段 1 计划审视（A-002）：核对验收矩阵、环境矩阵、账号权限 oracle 与 evidence schema 候选；D-002 由 proposed 冻结为 accepted；开放 required finding = 0。
2. 用户确认阶段 1 冻结后进入阶段 2：按冻结合同运行集成验收、按 evidence index 持久化机器可读证据、显式记录失败/排除/平台缺口。
3. 推送远端触发 GitHub Actions 首跑，复核 Linux/CI 等价与浏览器 E2E 在 CI 上的结果；若失败，按执行事实记录而非用本地 pass 掩盖。
4. 完成 R6 close-out 审计后，再由用户决定 Root R6 / `progress` / status；VP 关门另走 `/vision`。

## 进度评估

阶段 1 的验收矩阵（C-001～C-008）、环境矩阵（D-004）、账号权限 oracle 与 evidence schema 已升级为**冻结候选**：本地双服务实测、schema dry-run（5 artifact SHA-256 可重算）与浏览器 E2E 已通过。阶段 1 **尚未正式冻结**（计划审视 A-002 未做、D-002 仍 proposed），阶段 2 未开放；GitHub Actions 首跑未发生。没有 R6 验收完成事实，Root `progress` 仍 `5/6`。
