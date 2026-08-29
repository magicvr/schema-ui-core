---
doc_type: vision-plan
id: VP-023-productionization-cli-package
title: 包消费产线化（发布运营 + CLI + 形态细化 + 覆盖与迁移）
status: active
vision_ref: schema-ui-core-admin-foundation@0.3.0
lead_workspace: workspace-023-productionization-cli-package
created: 2026-08-29
updated: 2026-08-29
version: 0.2.0
---

# VP-023 · 包消费产线化

## 意图

把 workspace-022 已验证的「构建期包消费」**可行性闭环**（VP-022 closed · 判据 #1–5 有界口径 + go/no-go GO）升级为**可运营的 cli+包 分发路径**：真实发布通道实证（Go tag/proxy `go get` + npm registry 上传安装）、CLI（create/add/upgrade）、六包细化与 d.ts 自动化、PG external 实测与运维文档、下游上手与 fork→包迁移路径。

**战略语境（用户确认）**：维持 Charter `@0.3.0` 现表述——**fork 与包消费两种交付形态并存**、单主线、不维护平行代码线；本 VP **不**推进 fork 退役、不改 Charter。

## 方向级退出判据

在同时满足下列方向时，本 VP **可以**有界或完整关门（证据必须在工作区目标内）：

1. **真实发布通道闭环**：Go 侧 origin tag + 真实 `go get @vX`（或等价私有 proxy）消费实证；npm 侧 registry 上传 + `pnpm add @schema-ui/*@<ver>` 安装实证；实验下游仓全程 registry 语义消费（无 `replace`/`file:` 路径依赖）。
2. **CLI 闭环**：`create-schema-ui`（对标 dotnet new）生成下游骨架（Go 组合根 + Web 骨架 + config + 主题覆盖）；`add/upgrade`（对标 dotnet add package）完成一次模块装配与一次 registry 升级，**零冲突**可复现；CLI 输出与手工文档路径等价（golden-field 双轨对照）。
3. **形态细化与类型面**：六包（protocol/lib/theme/ui/renderer/shell）独立可发布；d.ts 自动化管线修复（TS5056 同类冲突解决）；契约冻结面升格 v1.3.0（六包导出 + peer 矩阵定稿）。
4. **覆盖与运维**：PG external 消费实测（生产权威方言，补 F-005 核销）；包形态下游的运维路径成文（启动/升级/迁移/备份/优雅停机与主仓契约一致）；Golden consumer 团队化（CI 槽位或常驻 repo）。
5. **上手与迁移**：QUICKSTART 增设「cli+包起步」为主路径章节；fork→包迁移指南（既有 fork 下游如何过渡）；golden-field 从零到上线走查（时长记录）。
6. **产线化报告**：真实发布往返耗时矩阵、CLI 上手实测、registry 升级演练（含 breaking 场景）结果、go 后清单核销表、及「cli+包是否提升为默认主路径」的建议（不改 Charter 的表述层建议）。

## 边界（不进本波）

- 不改 Charter（fork 并存表述维持）；不推进 fork 退役；不迁移既有 fork 消费者（指南先行）。
- 不引入 G2 多模块细 tag（单模块粗粒度保持）；不做运行时镜像/模块下载（④ 形态）；不做业务域。
- 不重开 VP-022 的 closed 记录；与 VP-009/010（持续程序）正交分流。
- npm registry 生产发布若需正式 scope/凭据，属外部动作——本波以可复现脚本 + 凭据注入点交付，真实上传以用户授权为界（同 I-003 先例）。

## 前置与消费基线

- 继承：VP-003 模块契约 · VP-005 主题覆盖 · VP-006 协议面（pin `v2.9.0`）· VP-008 `go` 消费基线 · **VP-022**（dist-lib 链路 / pack 脚本 / 双 golden 仓 / freeze-face v1.2.0 / go 后清单）。
- 激活门禁：VRev intent 审视 + **架构类轻量 freshness 复核**（VP-022 先例；候选基线随激活时点核对）。
- 实验下游仓：`golden-field`（用户建仓；命名见 VP 工作区绑定注记）。

## 工作区绑定

| workspace_id | root_goal | role | joined | notes |
|--------------|-----------|------|--------|-------|
| — | — | lead | — | `planned` 0 区；激活时开 `workspace-023-productionization-cli-package` 并填表 |

> **实验下游仓命名（2026-08-29 用户指定用途）**：主推 **`golden-field`**（产线化验证的试验田；与 golden-consumer/golden-web 系列一致）。备选：`golden-field-admin` / `petstore-admin`。用户建仓后于激活事务中登记 URL/路径。

## 关门记录

（仅 `closed` / `abandoned` 时填写。）

| date | outcome | summary | evidence_links | residuals |
|------|---------|---------|----------------|-----------|
| — | — | — | — | — |

## 规划修订短史

| date | change |
|------|--------|
| 2026-08-29 | 初创 v0.1.0（用户指令：维持 Charter 并存战略 + 立项产线化 VP-023；承接 VP-022 go 后清单；组合层平台波 · 与三分支正交；对标 .NET dotnet new + NuGet 产线化） |`n| 2026-08-29 | v0.2.0 · **激活**（用户建仓授权）：[VRev-051](reviews/VRev-051-vp023-activation.md) self `pass`（0 required · 架构类轻量 freshness PASS `5c168070`→`041744b3` 不暂挂 `go`）；`planned → active`；lead `workspace-023-productionization-cli-package`；实验仓 `golden-field` 初始化随开区 |