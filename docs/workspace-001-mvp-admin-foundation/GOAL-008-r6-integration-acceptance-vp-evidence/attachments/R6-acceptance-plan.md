---
title: R6 集成验收与 VP 证据计划
status: draft
created: 2026-08-01
updated: 2026-08-01
parent: GOAL-008-r6-integration-acceptance-vp-evidence
version: 0.2.0
---

# R6 集成验收与 VP 证据计划

> 状态：**冻结候选，尚未正式冻结**。本文件已把 VP 三条退出判据映射为验收矩阵、决定最低环境矩阵、验证 evidence schema 候选，并登记账号权限 oracle 候选；正式冻结仍须阶段 1 计划审视通过（A-002）且 D-002 由用户裁决接受。权威边界为 [Root](../../GOAL-001-mvp-admin-foundation/00-meta.md)、[VP-001](../../../vision/plans/VP-001-mvp-admin-foundation.md) 与 [I-PROTO-001 v0.1.3](../../GOAL-001-mvp-admin-foundation/attachments/I-PROTO-001-coverage-draft.md)。

## 1. 目的与证据原则

R6 要把“可运行、可 fork、MVP 协议边界内可验证、账号权限集成成立”从分散事实收敛为可复核的工作区证据。证据至少满足：

1. **主张绑定**：每份结果点名服务哪条 VP 退出判据与哪个 R6 成功标准。
2. **身份绑定**：记录 repo revision、工作树状态、OS/arch、runtime/dependency 版本与执行时间。
3. **命令绑定**：记录规范化命令、cwd、退出码、起止时间与结果文件；不只保留转述。
4. **边界诚实**：include-partial、exclude、未执行、平台缺口与 residual 明确列出；不得用总体 pass 隐藏。
5. **完整性**：机器可读 index 列举全部结果文件及 SHA-256；Markdown 仅作人类索引，不取代原始记录。
6. **状态分离**：证据齐备只能支持提出 R6/VP 决定，不自动改变 Goal 或 VP status。

## 2. VP 退出判据映射（冻结候选）

| VP 判据 | R6 验收主张 | 最低证据候选 | 计划工作区产物 | 阻断信息项 |
|---------|-------------|--------------|------------------|------------|
| 1 · React + Go 可运行、可 fork、固定协议边界 | 声明环境中的干净安装/构建/双服务启动可复现；Web 通过 API 与固定 manifest 运行 | revision/environment manifest；install/build/test 结果；API health/manifest HTTP；浏览器关键路径；协议 provenance/SHA | `attachments/evidence/index.json`；`runtime/`；`build/`；人类索引 | I-008-002、I-008-004、I-008-005 |
| 2 · 受控 MVP 覆盖中每项有实现、范例/场景与验证路径 | R2 v0.1.3 的 11 个纳入域均能从 R5 登记追到当前实现和可执行验证，回归结果不扩大边界 | 覆盖基线 hash/version；R5 registry snapshot；stage3 per-suite executed/excluded；schema/fixture SHA；范例路由 smoke | `coverage/coverage-map.json`；`coverage/conformance-results.json`；Q2 链接表 | I-008-001、I-008-004、I-008-005 |
| 3 · 核心账号权限前后端集成，不依赖未声明业务模块 | 账号上下文由 API 到 Web/Renderer/动作链可观察；允许路径与拒绝路径均符合冻结 oracle | API session/account result；Web context；D-PERM fixture/host 对照；浏览器正向/拒绝场景；依赖边界清单 | `account-permission/scenarios.json`；HTTP/browser 结果；依赖清单 | I-008-003、I-008-004、I-008-005 |

## 2b. 验收矩阵（I-008-001 · 冻结候选）

> 每条 R6 主张点名服务哪条 VP 判据 / 哪个 R6 成功标准，并给出执行入口、预期结果与证据路径；`exclude` 明确列出，禁止“测试全绿即关 VP”的隐含规则。本矩阵为**冻结候选**，正式冻结在阶段 1 审视通过后由 D-002 接受确定。

| # | 主张（claim） | VP 判据 | 执行入口（命令 / 运行态） | 预期结果 | 证据路径 | 排除 / residual |
|----|--------------|---------|---------------------------|----------|----------|-----------------|
| C-001 | Web 测试与构建在声明环境可复现 | 1 | `cd apps/web && npm test`；`cd apps/web && npm run build` | 15 files / 395 tests pass；Vite build pass | `attachments/evidence/planning/results/{web-test,web-build}.log`（draft dry-run 已哈希） | 干净安装重跑见 CI |
| C-002 | API 测试与构建在声明环境可复现 | 1 | `cd apps/api && go test ./...`；`cd apps/api && go build ./...` | go test/build pass | `attachments/evidence/planning/results/{api-test,api-build}.log` | — |
| C-003 | 双服务启动、health 与 Web→API proxy 成立 | 1、3 | 启动 `go run ./cmd/server`（:8080）与 `npm run dev`（:5173）；`GET /healthz`、`GET /`、`GET /api/accounts/me`、`GET /api/records` | 全部 HTTP 200；`/api` proxy 命中 Go API | `attachments/evidence/planning/results/runtime-probes.log` | 端口为本地声明值 |
| C-004 | 浏览器关键路径（shell 渲染 + 账号上下文经 proxy） | 1、3 | `cd apps/web && npm run test:e2e`（Playwright Chromium，webServer 双服务） | shell 渲染、manifest 导航、`/api/accounts/me` 返回 dev-001 session、`/api/records` 非空 | `apps/web/e2e/shell.spec.ts`；`test-results/r6-overview.png` | Windows 本地已验证；Linux/CI 由 workflow 承接 |
| C-005 | 账号权限允许路径（dev session 有 admin+editor）符合 oracle | 3 | `GET /api/accounts/me`；`apps/web/src/protocol/conformance/permissions-inheritance`（17 例 fixture） | session 含 `roles:[admin,editor]`；permission-inheritance 17 例 pass | runtime-probes；GOAL-006 `dperm/cases.json` | 拒绝路径见 oracle（C-006） |
| C-006 | 账号权限拒绝路径 / fail-closed 可见 | 3 | oracle 冻结场景：无能力时权限表达式求值 false、项目隐藏 | 拒绝路径按冻结 oracle 显式可见（不靠总体 pass 掩盖） | `attachments/account-permission-oracle.md`（I-008-003 候选） | 见 oracle 排除项 |
| C-007 | R2 v0.1.3 纳入域可追溯且回归不扩大边界 | 2 | `cd apps/web && npm test`（stage3 222 项，含 request-construction non-batch 64） | 每个纳入域从 R5 登记追到实现与可执行验证；include-partial/exclude 不变 | `I-007-001-registry.md`；`I-PROTO-001 v0.1.3` | batch request-construction 与 reactions multi-round 为既有排除（D-008/D-010） |
| C-008 | 证据包可解析、文件摘要可重算 | 1、2、3 | `node validate-evidence-dry-run.mjs`（ajv 2020 校验） | schema 校验通过；5 个 artifact SHA-256 可重算 | `evidence-index.schema.json`；`evidence-index.dry-run.json` | 正式 acceptance index 尚未持久化（阶段 2） |

## 3. 拟议证据包

以下为 `I-008-004` 的候选，schema 已通过 dry-run 校验（可解析、哈希可重算），**尚未**冻结为正式 acceptance contract：

```text
attachments/evidence/
├── index.json
├── environment.json
├── results/
│   ├── web-test.json
│   ├── web-build.json
│   ├── api-test.json
│   └── api-build.json
├── runtime/
│   ├── api-health.json
│   ├── app-manifest.json
│   └── browser-scenarios.json
├── coverage/
│   ├── coverage-map.json
│   └── conformance-results.json
└── account-permission/
    └── scenarios.json
```

`index.json` 候选最小字段：

- `schemaVersion`
- `goalId` / `workspaceId` / `vpId`
- `repositoryRevision` / `worktreeState`
- `startedAt` / `completedAt`
- `environmentRef`
- `results[]`：`id`、`claimRefs[]`、`command`、`cwd`、`exitCode`、`outcome`、`artifact`、`sha256`
- `exclusions[]` / `residuals[]`
- `overallOutcome`（不得覆盖单项失败或排除）

## 4. 命令与运行态入口候选

> 下列仅是当前仓库可发现入口；阶段 1 需在 `I-008-002` / `I-008-005` 中决定干净安装方式、平台矩阵、browser runner 与结构化 reporter。

| 轨道 | 当前入口 | R6 需补的身份/结果 |
|------|----------|--------------------|
| Web test | `cd apps/web && npm test` | 结构化结果、revision/environment、失败与排除明细 |
| Web build | `cd apps/web && npm run build` | 依赖/runtime 版本、输出摘要、可 fork 干净安装步骤 |
| API test | `cd apps/api && go test ./...` | 结构化或可解析结果、Go/OS 版本、race/平台范围决定 |
| API build | `cd apps/api && go build ./...` | binary/package identity 与构建环境 |
| 双服务 runtime | 已冻结候选（C-003）：`go run ./cmd/server` + `npm run dev`，health/proxy/账号上下文/records | 按 evidence index 持久化启动/停止、端口/env、结果 |
| 协议回归 | 复用 Web stage3 / upstream / renderer tests | per-suite executed/excluded 结果与 v0.1.3 映射 |
| 账号权限 E2E | 已冻结候选（C-004/C-005/C-006）：Playwright shell.spec + oracle | 正向、拒绝、缺上下文/能力与动作路径预期 |

## 4b. 规划期能力基线（2026-08-01）

以下是阶段 1 收集的本地能力事实；`I-008-005` 已由用户裁决采用「本轮搭建最小 CI + 浏览器矩阵」而非接受平台 residual：

| 项 | 结果 | 边界 |
|----|------|------|
| revision / worktree | `f3e04f6bd5c1f4ba6b7b72444fd9a0a0ab52d4d5` / clean（阶段 1 本轮运行时） | 证据 dry-run 绑定该 revision；阶段 2 以冻结 revision 重跑 |
| Web | `cd apps/web && npm test`：15 files / 395 tests passed；`npm run build` passed | 本地依赖树；干净安装重跑由 CI workflow 承接（`npm ci`） |
| API | `cd apps/api && go test ./...` passed；`go build ./...` passed | Linux/CI 等价由 workflow 承接 |
| runtime entry | API `:8080/healthz`；Web `:5173`（`host:127.0.0.1`），`/api` proxy 到 API | 双服务启动/health/proxy/账号上下文/records 已实测（C-003/C-004） |
| environment | Windows/amd64；Node `v22.17.0`；npm `10.9.2`；Go `1.26.0`；Vitest `3.2.7`；Playwright `1.62.1` | 最低矩阵已决定（见 §4c），Linux/CI 证据待 workflow 首跑 |
| CI / browser / reporter | 新增 `.github/workflows/r6-basic-matrix.yml`（web/api/browser-e2e 三 job）；Playwright Chromium 已本地跑通 | CI 实际首跑尚未发生（需推送到远端），不能把「已配置」写成「已跑绿」 |

`I-008-004` 的候选 schema 与 dry-run 分别见
[`evidence-index.schema.json`](evidence-index.schema.json) 和
[`evidence-index.dry-run.json`](evidence-index.dry-run.json)。schema 已通过 [validate-evidence-dry-run.mjs](validate-evidence-dry-run.mjs)（ajv 2020）校验：**可解析、5 个 artifact SHA-256 可重算**；但仍为 draft，正式 acceptance index 属阶段 2。

## 4c. 最低环境矩阵决定（I-008-005 · 冻结候选）

用户于 2026-08-01 裁决：**本轮搭建最小 CI + 浏览器矩阵**，不接受“Windows-only + 平台 residual”。

| 平台 / 轨道 | 机制 | 执行点 | 状态 |
|-------------|------|--------|------|
| Windows/amd64 本地 | 开发机命令（npm test/build、go test/build、双服务、Playwright） | 本机 | **已验证**（阶段 1 dry-run + E2E 通过） |
| Linux/amd64 CI 等价 | `.github/workflows/r6-basic-matrix.yml`：web job（Node 22 + npm ci + test/build）、api job（Go 1.26 + test/build）、browser-e2e job（Playwright Chromium） | GitHub Actions | **首跑 green**（run `30666932343`：api 22s / web 27s / browser-e2e 53s） |
| 浏览器 E2E | Playwright Chromium，webServer 启动双服务；`apps/web/e2e/shell.spec.ts` | 本机 + CI | **本机 + CI 均 pass** |

边界：若后续 CI 失败或超时，属阶段 2 执行事实，须在验收记录中显式列失败/排除，不得用本地 pass 掩盖；不静默降级回 residual。非阻断注解：`setup-go` 缓存因 `apps/api/go.sum` 缺失而 skip（API 无外部依赖，`go test`/`go build` 仍通过）。

## 5. 阶段门禁

### 阶段 1 → 阶段 2

- `I-008-001`～`I-008-005` 均 `verified`，或对明确范围已有用户书面 `accepted-residual`；
- D-002 已由后续 accepted 决策冻结，而非继续停在 proposed；
- 计划审视覆盖验收矩阵、命令、环境、oracle 与 evidence schema；开放 required finding = 0；
- 若审计意见冲突或 independent-only，按 P-004 先裁决。

### 阶段 2 → 阶段 3

- 冻结矩阵中全部 required execution item 已运行并落盘；
- 失败、未执行、排除与平台缺口均有显式结果，不用总体 pass 掩盖；
- evidence index 可解析，文件摘要可重算；
- 新发现关键未知已回流信息表，未静默扩域。

### 阶段 3 → 阶段 4

- VP 三条判据各有 Q2 工作区证据链接；
- required 缺口均按 P-003 合法闭合；
- residual 点名具体 workspace/goal、范围、期限/触发与缓解；
- 证据边界与 Root/VP/覆盖基线无明显冲突。

### R6 关门后

- 只有用户可授权 Root R6 检查点、`progress` 或 Root status 的变化；
- VP status 变化必须由 `/vision` 读取 R6 工作区证据、形成关门提案并获得用户确认；
- R6 `pass` 不等于完整协议支持、发布就绪或外部鉴证。

## 6. 当前已知缺口

- CI 首跑已 green（run `30666932343`）；Linux/CI 等价与浏览器 E2E 的 CI 证据已闭合。遗留非阻断注解：Node 20 弃用（GitHub 强制 Node 24）与 `apps/api/go.sum` 缺失（API 无外部依赖，setup-go 缓存 skip，`go test`/`go build` 仍通过）。
- Vitest/Go 当前命令未统一输出 JSON evidence artifact；阶段 2 需要证据 writer 或结构化结果落盘（`I-008-004` 已通过 draft schema dry-run 证明形状可解析）。
- 账号权限跨层正向已实测（C-005）；**拒绝路径 oracle 已登记候选**（[account-permission-oracle.md](account-permission-oracle.md)，I-008-003），阶段 2 执行前需审视通过。
- R5 recommended 项可作为整改候选；是否升级为 R6 required 必须有风险依据或用户/审计决定。
