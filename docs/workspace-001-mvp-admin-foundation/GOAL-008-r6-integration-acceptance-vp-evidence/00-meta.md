---
id: GOAL-008-r6-integration-acceptance-vp-evidence
title: R6 · 集成验收与 VP 证据
status: active
parent: GOAL-001-mvp-admin-foundation
created: 2026-08-01
updated: 2026-08-01
version: 0.2.3
---

# GOAL-008 · R6 · 集成验收与 VP 证据

## 概述

承接 Root 已完成的 R1-R5，围绕 [VP-001 的三条方向级退出判据](../../vision/plans/VP-001-mvp-admin-foundation.md#方向级退出判据)建立可复核的集成验收合同、执行证据与工作区证据索引。目标是回答“当前工作区的 React + Go Admin MVP 是否有足够证据支持提出 VP 关门”，而不是仅重述既有测试通过或自动修改 VP 状态。

本目标处于**规划期**。当前只固化范围、路线、信息门禁与证据形状；R6 验收执行、Root `6/6`、Root `done` 与 VP `closed` 均未发生。

## 范围边界

### 纳入

- 将 VP-001 三条退出判据映射为可执行的验收主张、命令/运行态检查、证据文件与工作区内 Q2 链接。
- 验证 React Web 与 Go API 在声明环境中的可运行、可 fork 基线，包括干净依赖安装/构建、双服务启动、HTTP 健康与浏览器关键路径。
- 对 R2 冻结的 MVP 覆盖边界与 R3-R5 证据做集成级回归：协议 pin、逐域实现/范例/验证入口、结构/行为 conformance 与明确排除项保持一致。
- 验证核心账号与权限的前后端集成正向/负向路径，且不依赖未声明业务领域模块。
- 形成带 revision、环境、命令、退出码、时间、结果、排除/残余与 SHA-256 的机器可读证据索引，并在关门审计后向 `/vision` 提供 VP 关门提案输入。

### 排除

- 不扩大 [I-PROTO-001 v0.1.3](../GOAL-001-mvp-admin-foundation/attachments/I-PROTO-001-coverage-draft.md) 的 include / include-partial 边界，不主张完整协议支持。
- 不在本目标中定义新的业务领域模块、完整 component registry、上传域或多选批量语义。
- 不把本地一次测试输出、截图或 Root `progress` 单独当作 VP 关门证据。
- 不自动修改 `docs/vision/plans/VP-001-mvp-admin-foundation.md` 的 status；VP 关门属于 `/vision` + 用户确认。
- 不在规划回合实施 CI、发布、部署或产品化增强；它们是否进入 R6 验收最低集由信息项与后续决策确定。

## 高层路线图（P-001）

| 阶段 | 名称 | 状态 | 退出条件 |
|------|------|------|----------|
| 1 | 验收合同与证据计划冻结 | **已冻结**（2026-08-01） | `I-008-001`～`I-008-005` 均有证据结论或合规 residual；验收矩阵、环境矩阵、账号权限 oracle 与证据格式经计划阶段审视（A-002），无开放 required finding |
| 2 | 集成验收执行 | **已完成**（2026-08-01） | 在已声明 revision/环境运行 Web、API、协议回归与账号权限集成检查；原始/机器可读结果落盘；失败与排除不被隐藏 |
| 3 | VP 证据汇编与缺口整改 | 未开始 | VP 三条退出判据逐条指向工作区 Q2 证据；所有 required 缺口 fixed 或经用户书面 residual/overruled；边界主张一致 |
| 4 | R6 关门审计与 VP 提案输入 | 未开始 | R6 close-out 审计结论可核对、开放 required=0；用户另行授权 Root R6/`progress`/status 变化，并由 `/vision` 决定是否提出 VP 关门 |

> 阶段 2 → 3 门禁已过（A-003 pass，2026-08-01）：C-001～C-008 全执行、evidence-index（mode: acceptance）经 schema 校验、排除显式。阶段 3 可开始。

阶段通常串行；若阶段 2 内的 Web/API、协议回归、账号权限与证据打包已由阶段 1 冻结为独立执行块，可并行收集，但不得越过同一 required 信息门禁。

## 成功标准（规划基线）

- [ ] 一份受控验收矩阵把 VP-001 三条退出判据映射为明确主张、执行入口、预期结果、证据路径、排除与 residual；不存在“测试全绿即可关 VP”的隐含规则。
- [ ] React + Go 基架在声明的干净环境完成可复现启动与浏览器/API 关键路径验证；证据绑定 repo revision、依赖/runtime 版本与工作树状态。
- [ ] R2 v0.1.3 的每个纳入域均可从 R5 登记表追到实现、范例/场景与可执行验证，且 R6 回归没有越过 include-partial / exclude 边界。
- [ ] 核心账号与权限链路至少包含可核对的前后端正向与拒绝路径，证明集成不依赖未声明业务模块。
- [ ] 机器可读证据索引完整记录命令、退出码、时间、环境、结果、排除/残余与文件摘要；证据可由工作区目标记录稳定寻址。
- [ ] R6 关门审计无开放 required finding；Root R6 完成、Root `done` 与 VP 关门仍分别等待用户和 `/vision` 的后续受控决定。

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-008-001 | required | VP-001 三条退出判据分别需要哪些最小验收主张、证据与允许的有界排除？ | 阶段 1 计划冻结 / 阶段 4 VP 提案输入 | 阶段 1 结束前 | 对照 VP-001、Root 成功标准、R2/R5 登记与历史 residual，完成并审视验收矩阵 | **verified**（A-002） | — | [R6-acceptance-plan.md](attachments/R6-acceptance-plan.md) v0.2.0 §2b 验收矩阵 C-001～C-008 |
| I-008-002 | required | “可运行、可 fork”的干净环境、依赖安装、双服务启动、端口/env 与最小浏览器/API 路径是什么？ | 阶段 1 计划冻结 / 阶段 2 执行 | 阶段 1 结束前 | 在隔离/干净工作副本演练启动路径，记录 OS、runtime、命令、env 与失败模式 | **verified**（A-002） | — | 本地双服务/health/proxy/账号上下文/records 实测落盘 `evidence/planning/results/runtime-probes.log`；干净安装（`npm ci`）+ Linux/CI 等价由 **GitHub Actions 首跑 green**（run `30666932343`，api/web/browser-e2e 全 pass）闭合 |
| I-008-003 | required | 账号与权限前后端集成的正向/拒绝 oracle、身份载体与期望错误/可见性结果是什么？ | 阶段 2 账号权限验收 / 关门 | 阶段 1 结束前 | 对照 GOAL-006 设计/fixtures 与当前 host/API，冻结端到端场景和预期 | **verified**（A-002） | — | [account-permission-oracle.md](attachments/account-permission-oracle.md) v0.1.0（P-1～P-4 正向、D-1～D-6 拒绝） |
| I-008-004 | required | R6 机器可读证据包的 schema、revision identity、目录、哈希与重跑规则是什么？ | 阶段 2 证据采集 / 阶段 3 汇编 | 阶段 1 结束前 | 定义 evidence index + result 记录格式；用一次 dry-run 验证可解析性和文件摘要 | **verified**（A-002） | — | `evidence-index.schema.json` + `validate-evidence-dry-run.mjs` 校验通过：**可解析、5 artifact SHA-256 可重算**；正式 acceptance index 属阶段 2 |
| I-008-005 | required | R6 的验收环境矩阵如何处理 Windows 本地、Linux/CI 等价、浏览器 E2E 与缺失平台证据？ | 阶段 1 计划冻结 / 阶段 4 关门主张 | 阶段 1 结束前 | 盘点现有 CI/浏览器能力；决定最低矩阵，或由用户书面接受有范围与复审触发的 residual | **verified**（A-002） | — | 用户裁决“搭建最小 CI+浏览器矩阵”（D-004）；workflow + Playwright 已建、本机 E2E 通过；**GitHub Actions 首跑 green**（run `30666932343`，browser-e2e job 53s pass） |

## 阶段 1 当前收集结论（2026-08-01）

- **阶段 1 已冻结**：A-002 计划审视 pass，D-002 冻结为 accepted，F-008-001 关闭（`fixed`），五项 `I-008-001`～`I-008-005` 均 `verified`。阶段 2 可开始。
- 本轮按用户裁决（D-004）搭建最小 CI + 浏览器矩阵：`.github/workflows/r6-basic-matrix.yml`（web/api/browser-e2e 三 job）、Playwright `apps/web/e2e/shell.spec.ts`、`playwright.config.ts`；本地 Playwright E2E 通过（shell 渲染 + `/api/accounts/me` 经 proxy 返回 dev session + `/api/records`）。
- `I-008-004` 的 draft schema 已用真实 artifact dry-run 验证：**可解析、5 个 artifact SHA-256 可重算**；dry-run 持久化为 `evidence/planning/evidence-index.dry-run.json`，校验脚本 `validate-evidence-dry-run.mjs` 可重跑。
- **GitHub Actions 首跑已 green**（run `30666932343`）：api/web/browser-e2e 三 job 全 pass，`I-008-002`（`npm ci` 干净安装 + Linux 等价）与 `I-008-005`（浏览器矩阵 CI 证据）已由实际运行闭合；非阻断注解为 Node 20 弃用与 `go.sum` 缺失（API 无外部依赖）。
- 阶段 2 执行尚未开始；Root `progress: 5/6`、VP-001 `active` 不变。

## 父目标

- [GOAL-001-mvp-admin-foundation](../GOAL-001-mvp-admin-foundation/00-meta.md)

## 备注

- 工作区：`workspace-001-mvp-admin-foundation`；canonical 范围：`docs/workspace-001-mvp-admin-foundation/`。
- 本目标继承 Root 的 `primary_plan: VP-001-mvp-admin-foundation` 语境，不扩写第二套愿景边界。
- R4/R5 的 recommended 跟踪项可作为 R6 缺口输入，但不会在没有 required 升级或用户裁决时自动阻断。
