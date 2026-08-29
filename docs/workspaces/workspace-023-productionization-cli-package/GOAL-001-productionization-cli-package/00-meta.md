---
id: GOAL-001-productionization-cli-package
title: 包消费产线化（cli+包 分发路径可运营化）
status: done
parent: null
created: 2026-08-29
updated: 2026-08-29
version: 0.7.0
progress: 5/5
plan_refs:
  - VP-023-productionization-cli-package
primary_plan: VP-023-productionization-cli-package
serves_summary: 把 VP-022 验证的构建期包消费可行性闭环升级为可运营的 cli+包 分发路径（真实发布/CLI/六包细化/覆盖运维/报告）
---

# GOAL-001 · 包消费产线化

## 概述

承接 [VP-023-productionization-cli-package](../../vision/plans/VP-023-productionization-cli-package.md)（**closed** v0.3.0 · 2026-08-29）：六条方向级退出判据全部落地（真实发布通道 / CLI / 六包细化+d.ts / PG+运维 / 上手迁移 / 产线化报告）。**不改 Charter**（fork 与包消费并存维持，Charter @0.3.0 表述未动）。实验下游仓 = `github.com/magicvr/golden-field`（已初始化 · Go 无 replace · 六包 GH Packages registry 语义消费 · consumer-regression workflow 槽位）。breaking 实演以 v0.3.0 真实执行（用户 P-004 裁决）。

## 成功标准（对应 VP-023 六条判据）

- [x] 判据 #1：真实发布通道闭环（Go `go get @vX` + npm registry 安装；golden-field 无 replace/file:）
- [x] 判据 #2：CLI 闭环（create/add/upgrade · 一次 registry 升级零冲突 · 双轨对照）
- [x] 判据 #3：六包独立发布 + d.ts 自动化（TS5056 修复）→ 冻结面 v1.3.0
- [x] 判据 #4：PG external 实测 + 运维路径文档 + golden 仓团队化
- [x] 判据 #5：QUICKSTART cli+包 主路径章节 + fork→包迁移指南 + golden-field 从零上线走查
- [x] 判据 #6：产线化报告（往返耗时/CLI 实测/breaking 演练/核销表/默认主路径建议）

> 六条全部达成（2026-08-29 关门）；breaking 演练 = 真实实演（v0.3.0 `kernel.JoinKeys → JoinIdentifiers` · 用户 P-004 裁决）；grok 独立双审 F-001～F-008 全闭合。残余 = go 后清单 7 项（serve 壳 / 六包 external 化 / 纯原子拆分 / fork 对照计时 / 迁移工具化 / 包公开可见性 / compose CI 实跑）→ 已立项收口（VP-024 · planned）。

## 纲领路线图（P-001）

| 阶段 | 内容 | 检查点/状态 |
|------|------|-------------|
| R1 | 真实发布通道：Go tag+`go get`（或私有 proxy）实效；npm registry 上传+安装实效；golden-field 初始化并移除 replace/file: 依赖 | **已关门**（2026-08-29 · GOAL-002 done 4/4 · A-001 self `pass` · 判据 #1 满足；golden-field 全程 registry 语义；升级演练绑定 R2 发布） |
| R2 | CLI 闭环：create-schema-ui / add / upgrade（对标 dotnet new + NuGet）——golden-field 双轨对照 | 依赖 R1 | **已关门**（2026-08-29 · GOAL-003 done 4/4 · A-001 self `pass` · 判据 #2 满足；CLI 实现 + create 双端全绿 + registry 升级演练零冲突） |
| R3 | 六包细化（protocol/lib/theme/ui/renderer/shell）+ d.ts 自动化（TS5056）→ 冻结面 v1.3.0 | 依赖 R1 | **已关门**（2026-08-29 · GOAL-004 done 4/4 · A-001 self `pass` · 判据 #3 满足；六包 registry 发布 + TS5056 根治 · F-006 核销） |
| R4 | PG external 实测（F-005 核销）+ 运维路径文档（启动/升级/迁移/备份/drain）+ golden 仓 CI 槽位 | 依赖 R1–R3 | **已关门**（2026-08-29 · GOAL-005 done 4/4 · A-001 self `pass` · 判据 #4 满足；PG 双方言实证 · ops-playbook+compose+workflow；F-001 推荐 compose 实跑留 CI） |
| R5 | 产线化报告（判据 #6）+ 核销表 + 建议 → independent 审计（grok）→ 关门 | 依赖 R1–R4 | **已关门**（2026-08-29 · GOAL-006 done 4/4 · 判据 #5/#6 满足 · breaking 实演 v0.3.0 · grok 独立双审 F-001～F-008 全闭合 · Root done 5/5） |

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-023-001 | required | Go 发布通道形态（origin tag + 公共 proxy）与凭据可用性 | R1 发布实证 | R1 | `go get` 实测 | **verified**（`apps/api/v0.1.0/v0.2.0` tag + 公共 proxy 下载实证） | — | GOAL-002 E-001 |
| I-023-002 | required | npm registry 目标（GitHub Packages）与发布凭据 | R1 发布实证 | R1 | 上传样例 | **verified**（GH Packages 六包发布 + golden-field registry 安装实证） | — | GOAL-002 E-002 / GOAL-004 E-001 |
| I-023-003 | required | PG 可用实例（docker postgres:16 容器 gf-pg） | R4 实测 | R4 | 容器就绪检查 | **verified**（15432→5432 · 63 迁移 apply · 幂等重入） | — | GOAL-005 E-001 |

## 父目标

- `null`（Root）

## 台账布局

`01-decision/`、`02-execution/`、`03-audit/` 平铺 ledger；条目见各索引。

## 备注

- 审计模式（D-001）：阶段关门 default self；R5 与 Root 关门 = independent（grok build 先例）。
- freshness 三字段（VP-022 先例）：见 `01-decision/D-001`.
- golden-field 初始化 = R1 前置动作（用户授权「用的时候你自己初始化」）。