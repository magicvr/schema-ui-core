---
title: R6 集成验收与 VP 证据计划
status: draft
created: 2026-08-01
updated: 2026-08-01
parent: GOAL-008-r6-integration-acceptance-vp-evidence
version: 0.1.1
---

# R6 集成验收与 VP 证据计划

> 状态：**规划草案，尚未冻结**。本文件定义待验证的验收与证据结构，不证明任何 R6 主张已经满足。权威边界为 [Root](../../GOAL-001-mvp-admin-foundation/00-meta.md)、[VP-001](../../../vision/plans/VP-001-mvp-admin-foundation.md) 与 [I-PROTO-001 v0.1.3](../../GOAL-001-mvp-admin-foundation/attachments/I-PROTO-001-coverage-draft.md)。

## 1. 目的与证据原则

R6 要把“可运行、可 fork、MVP 协议边界内可验证、账号权限集成成立”从分散事实收敛为可复核的工作区证据。证据至少满足：

1. **主张绑定**：每份结果点名服务哪条 VP 退出判据与哪个 R6 成功标准。
2. **身份绑定**：记录 repo revision、工作树状态、OS/arch、runtime/dependency 版本与执行时间。
3. **命令绑定**：记录规范化命令、cwd、退出码、起止时间与结果文件；不只保留转述。
4. **边界诚实**：include-partial、exclude、未执行、平台缺口与 residual 明确列出；不得用总体 pass 隐藏。
5. **完整性**：机器可读 index 列举全部结果文件及 SHA-256；Markdown 仅作人类索引，不取代原始记录。
6. **状态分离**：证据齐备只能支持提出 R6/VP 决定，不自动改变 Goal 或 VP status。

## 2. VP 退出判据映射草案

| VP 判据 | R6 验收主张 | 最低证据候选 | 计划工作区产物 | 阻断信息项 |
|---------|-------------|--------------|------------------|------------|
| 1 · React + Go 可运行、可 fork、固定协议边界 | 声明环境中的干净安装/构建/双服务启动可复现；Web 通过 API 与固定 manifest 运行 | revision/environment manifest；install/build/test 结果；API health/manifest HTTP；浏览器关键路径；协议 provenance/SHA | `attachments/evidence/index.json`；`runtime/`；`build/`；人类索引 | I-008-002、I-008-004、I-008-005 |
| 2 · 受控 MVP 覆盖中每项有实现、范例/场景与验证路径 | R2 v0.1.3 的 11 个纳入域均能从 R5 登记追到当前实现和可执行验证，回归结果不扩大边界 | 覆盖基线 hash/version；R5 registry snapshot；stage3 per-suite executed/excluded；schema/fixture SHA；范例路由 smoke | `coverage/coverage-map.json`；`coverage/conformance-results.json`；Q2 链接表 | I-008-001、I-008-004、I-008-005 |
| 3 · 核心账号权限前后端集成，不依赖未声明业务模块 | 账号上下文由 API 到 Web/Renderer/动作链可观察；允许路径与拒绝路径均符合冻结 oracle | API session/account result；Web context；D-PERM fixture/host 对照；浏览器正向/拒绝场景；依赖边界清单 | `account-permission/scenarios.json`；HTTP/browser 结果；依赖清单 | I-008-003、I-008-004、I-008-005 |

## 3. 拟议证据包

以下为 `I-008-004` 的候选，不是已冻结 schema：

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
| 双服务 runtime | 待 `I-008-002` 冻结 | 启动/停止、端口/env、health、manifest、API proxy 与浏览器场景 |
| 协议回归 | 复用 Web stage3 / upstream / renderer tests | per-suite executed/excluded 结果与 v0.1.3 映射 |
| 账号权限 E2E | 待 `I-008-003` 冻结 | 正向、拒绝、缺上下文/能力与动作路径预期 |

## 4b. 规划期能力基线（2026-08-01）

以下是当前 revision 上为阶段 1 收集的本地能力事实，不是 R6 验收证据：

| 项 | 结果 | 边界 |
|----|------|------|
| revision / worktree | `7d20acc7702bcc0e514f787c455bf9c93d5b832f` / clean | 仅绑定本次规划复跑；尚未形成持久化 acceptance artifact |
| Web | `cd apps/web && npm test`：15 files / 395 tests passed；`npm run build` passed | 依赖树为当前工作副本；干净安装重跑规则仍待 I-008-002 |
| API | `cd apps/api && go test ./...` passed；`go build ./...` passed | 未覆盖 Linux/CI 等价性 |
| runtime entry | API `:8080/healthz`；Web `:5173`，`/api` proxy 到 API | 双服务启动/停止与浏览器关键路径尚未按 R6 contract 持久化 |
| environment | Windows/amd64；Node `v22.17.0`；npm `10.9.2`；Go `1.26.0`；Vitest `3.2.7` | 不代表最低支持矩阵已决定 |
| CI / browser / reporter | 未发现 `.github/workflows`、Playwright/Puppeteer/Cypress 等 runner 或 JSON/JUnit/evidence writer | I-008-005 仍需决定最低矩阵，不能把缺失当作已验证 |

`I-008-004` 的候选 schema 与 dry-run 分别见
[`evidence-index.schema.json`](evidence-index.schema.json) 和
[`evidence-index.dry-run.json`](evidence-index.dry-run.json)。二者都明确标为 draft/planning，dry-run 的结果没有持久化产物摘要，不能替代阶段 2 证据。

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

- 仓库未发现 CI workflow；Linux/CI 等价证据为空。
- 未发现现成浏览器 E2E 或统一 runtime smoke runner。
- Vitest/Go 当前命令未统一输出本计划拟议的 JSON evidence artifact。
- 账号权限已有 R4 单元/fixture/HTTP 证据，但 R6 的跨层正向/拒绝 oracle 尚未冻结。
- R5 recommended 项可作为整改候选；是否升级为 R6 required 必须有风险依据或用户/审计决定。
