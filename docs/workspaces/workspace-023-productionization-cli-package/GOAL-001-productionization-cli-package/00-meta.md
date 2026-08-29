---
id: GOAL-001-productionization-cli-package
title: 包消费产线化（cli+包 分发路径可运营化）
status: active
parent: null
created: 2026-08-29
updated: 2026-08-29
version: 0.1.0
progress: 0/5
plan_refs:
  - VP-023-productionization-cli-package
primary_plan: VP-023-productionization-cli-package
serves_summary: 把 VP-022 验证的构建期包消费可行性闭环升级为可运营的 cli+包 分发路径（真实发布/CLI/六包细化/覆盖运维/报告）
---

# GOAL-001 · 包消费产线化

## 概述

承接 [VP-023-productionization-cli-package](../../vision/plans/VP-023-productionization-cli-package.md)（active v0.2.0）：六条方向级退出判据落地（真实发布通道 / CLI / 六包细化+d.ts / PG+运维 / 上手迁移 / 产线化报告）。**不改 Charter**（fork 与包消费并存）。实验下游仓 = `github.com/magicvr/golden-field`（本仓平级 · 空仓待初始化）。

## 成功标准（对应 VP-023 六条判据）

- [ ] 判据 #1：真实发布通道闭环（Go `go get @vX` + npm registry 安装；golden-field 无 replace/file:）
- [ ] 判据 #2：CLI 闭环（create/add/upgrade · 一次 registry 升级零冲突 · 双轨对照）
- [ ] 判据 #3：六包独立发布 + d.ts 自动化（TS5056 修复）→ 冻结面 v1.3.0
- [ ] 判据 #4：PG external 实测 + 运维路径文档 + golden 仓团队化
- [ ] 判据 #5：QUICKSTART cli+包 主路径章节 + fork→包迁移指南 + golden-field 从零上线走查
- [ ] 判据 #6：产线化报告（往返耗时/CLI 实测/breaking 演练/核销表/默认主路径建议）

## 纲领路线图（P-001）

| 阶段 | 内容 | 检查点/状态 |
|------|------|-------------|
| R1 | 真实发布通道：Go tag+`go get`（或私有 proxy）实效；npm registry 上传+安装实效；golden-field 初始化并移除 replace/file: 依赖 | 未开 |
| R2 | CLI 闭环：create-schema-ui / add / upgrade（对标 dotnet new + NuGet）——golden-field 双轨对照 | 依赖 R1 |
| R3 | 六包细化（protocol/lib/theme/ui/renderer/shell）+ d.ts 自动化（TS5056）→ 冻结面 v1.3.0 | 依赖 R1 |
| R4 | PG external 实测（F-005 核销）+ 运维路径文档（启动/升级/迁移/备份/drain）+ golden 仓 CI 槽位 | 依赖 R1–R3 |
| R5 | 产线化报告（判据 #6）+ 核销表 + 建议 → independent 审计（grok）→ 关门 | 依赖 R1–R4 |

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-023-001 | required | Go 发布通道形态（origin tag + 公共 proxy vs 私有 proxy/file proxy）与凭据可用性 | R1 发布实证 | R1 | 通道探测 + `go get` 实测 | open | — | 待确认 |
| I-023-002 | required | npm registry 目标（npmjs scope vs GitHub Packages）与发布凭据 | R1 发布实证 | R1 | 通道探测 + 上传样例 | open | — | 待确认 |
| I-023-003 | required | PG 可用实例（本机 CI 或容器）供 external 消费实测 | R4 实测 | R4 | 环境探测（docker/pg 服务） | open | — | 待确认 |

## 父目标

- `null`（Root）

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger；条目见各索引。

## 备注

- 审计模式（D-001）：阶段关门 default self；R5 与 Root 关门 = independent（grok build 先例）。
- freshness 三字段（VP-022 先例）：见 `01-decision/D-001`.
- golden-field 初始化 = R1 前置动作（用户授权「用的时候你自己初始化」）。