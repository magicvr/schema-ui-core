---
id: workspace-023-productionization-cli-package
title: 包消费产线化工作区
status: active
root_goal: GOAL-001-productionization-cli-package
canonical_scope: docs/workspaces/workspace-023-productionization-cli-package/
shared_materials_catalog: none
vision_role: delivery
plan_refs:
  - VP-023-productionization-cli-package
primary_plan: VP-023-productionization-cli-package
created: 2026-08-29
updated: 2026-08-29
version: 0.1.0
parent: null
---

# 工作区上下文 · 包消费产线化

本工作区是 [VP-023-productionization-cli-package](../../vision/plans/VP-023-productionization-cli-package.md)（**`active`** v0.2.0 · 2026-08-29 激活）的唯一 lead delivery workspace。**组合层平台波**：把 VP-022 的可行性闭环升级为可运营的 cli+包 分发路径（真实发布通道 / CLI / 六包细化+d.ts / PG+运维 / 上手迁移）。**不改 Charter**（fork 与包消费并存表述维持）。

- **Root** `GOAL-001-productionization-cli-package`：**`active`** · 0/5（R1 发布通道 → R2 CLI → R3 六包细化+d.ts → R4 覆盖运维 → R5 报告与关门）。
- 激活门禁（2026-08-29）：VRev-051 self `pass`（0 required）；**架构类轻量 freshness PASS**（`5c168070` → `041744b3`，不暂挂 `go`）；VP-009/010 无阻断。
- **实验下游仓**：`github.com/magicvr/golden-field`（克隆于 `../golden-field`，空仓待初始化）。
- 消费基线：VP-022 交付（dist-lib 链路 / pack 脚本 / golden-consumer/golden-web / 冻结面 v1.2.0 / go 后清单）。
- 不改变 Charter `primary_workspace`。

## 纲领阶段（Root 路线图指针）

| 阶段 | 内容 | 状态 |
|------|------|------|
| R1 | **真实发布通道**：Go origin tag + `go get @vX`（或私有 proxy）实效；npm registry 上传 + `pnpm add @ver` 实效；golden-field 移除 replace/file: 依赖 | 未开 |
| R2 | **CLI 闭环**：create/add/upgrade（对标 dotnet new + NuGet）双轨对照 | 依赖 R1 |
| R3 | **六包细化 + d.ts 自动化**（TS5056 修复）→ 冻结面 v1.3.0 | 依赖 R1 |
| R4 | **覆盖运维**：PG external 实测 + 运维路径文档 + golden 仓团队化（CI 槽位） | 依赖 R1–R3 |
| R5 | **报告与关门**：产线化报告 + 核销表 + 建议 → independent 审计（grok）→ 关门 | 依赖 R1–R4 |