---
id: workspace-024-distribution-formalization
title: 分发形态正式化工作区（cli+包 对外服务化收口 · 已结项）
status: done
root_goal: GOAL-001-distribution-formalization
canonical_scope: docs/workspaces/workspace-024-distribution-formalization/
shared_materials_catalog: none
vision_role: delivery
plan_refs:
  - VP-024-distribution-formalization
primary_plan: VP-024-distribution-formalization
created: 2026-08-29
updated: 2026-08-30
version: 0.2.0
parent: null
---

# 工作区上下文 · 分发形态正式化（cli+包 对外服务化收口）

本工作区是 [VP-024-distribution-formalization](../../vision/plans/VP-024-distribution-formalization.md)（**`closed`** v0.3.0 · 2026-08-29 用户书面确认）的唯一 lead delivery workspace。**组合层平台波**：收口 VP-022/023 go 后合并残余 7 项，把「cli+包 分发路径」升级为对外正式化（serve 壳 / npmjs 公开发布 / compose CI 实跑 / fork 对照计时 / 六包形态细化 / 迁移工具化 / 方法 B 置顶）；不改 Charter、不弃 fork。**工作区已结项**（2026-08-29）：Root `done · 7/7`，关门双审闭合，残余四项书面登记（GOAL-008 D-001）。

- **Root** `GOAL-001-distribution-formalization`：**`done`** · 7/7（R1 serve 壳 → R2 公开发布 → R3 compose/CI 实跑 → R4 fork 对照计时 → R5 六包形态细化 → R6 迁移工具化 → R7 置顶与收口报告，全部关门；2026-08-29 用户书面确认）。
- 激活门禁已满足（2026-08-29）：[VRev-052](../../vision/reviews/VRev-052-vp024-activation.md) self `pass`（0 required；V-F087/V-F088 recommended → 开区事务内 fixed）；**平台/架构类轻量 freshness PASS**（`041744b3` → `c9122478`：协议 pin / 依赖锁 / 迁移台账 / Profile 默认集 / provenance 无变更；breaking v0.3.0 = 冻结面流程内变更）不暂挂 `go`。
- 不改变 Charter `primary_workspace`（仍为 workspace-001）。
- 消费基线（VP-022/023 已交付）：CLI `schema-ui create/add/upgrade` · GH Packages 六包 · freeze-face v1.3.0 · ops-playbook/compose/workflow · QUICKSTART 方法 B · fork→包迁移指南 · `apps/api/v0.1.0/v0.2.0/v0.3.0` tag。
- 实验下游仓：`github.com/magicvr/golden-field`（registry 语义消费 · 已钉 `apps/api v0.3.0`）。
- 边界（与 VP-024 冻结一致）：不改 Charter（fork 并存维持）；不重开 workspace-022/023；不做运行时模块下载/热插拔；不引入 G2 细版本；npmjs 正式 scope/凭据 = 外部动作，R2 前置门禁（用户授权为界）。

## 绑定

| 字段 | 当前值 | 说明 |
|------|--------|------|
| 工作区 ID | `workspace-024-distribution-formalization` | 与本区目标及资料引用的 `workspace_id` 一致 |
| Root Goal | `GOAL-001-distribution-formalization` | `parent: null`；**done** · 7/7（2026-08-29 结项） |
| canonical 范围 | `docs/workspaces/workspace-024-distribution-formalization/` | 本区唯一目标状态范围 |
| 共享资料目录 | `none` | 暂无固定共享资料 |
| 愿景角色 | `delivery` | VP-024 lead（active）；不改变 Charter primary workspace |
| 规划对齐 | `primary_plan` = `VP-024-distribution-formalization`（`closed` v0.3.0） | 2026-08-29 激活/开区（VRev-052 self `pass`）；同日用户书面确认关门（VRev-053 independent `pass`） |

## 愿景对齐

Charter：`schema-ui-core-admin-foundation@0.3.0`（2026-08-29 strategic：构建期包消费 + pin 2.9.0）。
VP-024：分发形态正式化 · cli+包 对外服务化收口（vision_ref @0.3.0）——八条方向级退出判据 = VP-022/023 go 后合并残余 7 项 + 方法 B 置顶与收口报告；外部动作边界（npmjs 授权）以用户授权为界。

## 纲领阶段（Root 路线图指针）

| 阶段 | 内容 | 状态 |
|------|------|------|
| R1 | **serve 壳闭环**（判据 #1）：`schema-ui serve` 子命令（HTTP 壳 + config 装载 + assembly 服务器面）；create 骨架直接 `serve` 启动；RT-D02 合同接线 | **已关门**（GOAL-002 done 5/5） |
| R2 | **公开发布通道**（判据 #2）：npmjs.com 六包 + CLI 公开发布；golden-field 免凭据消费实证；发布流程成文（**R2 前置门禁 I-024-001：用户授权正式 scope/凭据**） | **已关门**（GOAL-003 done 4/4 · scope @magicvr 定稿） |
| R3 | **compose/CI 实跑**（判据 #3）：compose/Dockerfile + consumer-regression workflow 真实 CI（或用户环境等价）实跑 PASS；workflow 补 pnpm setup 与 token 注记 | **已关门**（GOAL-004 done 4/4 · hosted 实触发 = 登记） |
| R4 | **fork 对照计时**（判据 #4）：同一演进集 fork 同步 vs 包路径 bump 实测对比（耗时/冲突/迁移成本）定量结论 | **已关门**（GOAL-005 done 4/4） |
| R5 | **六包形态细化**（判据 #5/#6）：renderer 依赖图 external 化（ui 包可消费）；纯原子拆分；冻结面升格 v1.4.0 | **已关门**（GOAL-006 done 5/5） |
| R6 | **迁移工具化**（判据 #7）：fork→包迁移指南配套工具（`schema-ui migrate-fork` 或等价） | **已关门**（GOAL-007 done 4/4） |
| R7 | **置顶与收口报告**（判据 #8）：QUICKSTART 方法 B 置顶（fork 第二路径·Charter 不动）+ 收口报告（核销表/公开消费往返实证/fork 对照结论/残余清零） | **已关门**（GOAL-008 done 4/4 · Root 关门审计 A-002 grok `pass` · 用户书面确认 · VP-024 `closed` v0.3.0（VRev-053）） |

## 固定共享资料引用

> `shared-materials/index.json` 只能提供候选路径与摘要。缺完整引用字段的行无效。

| reference_id | workspace_id | material_id | source | version | sha256 | purpose | local_record | status |
|--------------|--------------|-------------|--------|---------|--------|---------|--------------|--------|
| — | — | — | — | — | — | — | none | — |