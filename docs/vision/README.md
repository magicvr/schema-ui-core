---
title: docs/vision · 规则面与本仓实例索引
status: active
created: 2026-07-29
updated: 2026-08-04
parent: null
version: 0.5.1
---

# docs/vision · 愿景层

本目录同时承载两类内容，**请区分**：

| 类别 | 角色 | 是否随 core 预装 |
|------|------|------------------|
| **规则权威** | 对齐门禁与目录约定 | 是（Skills core 镜像） |
| **本仓实例** | Charter / VP / Review 过程 | 否（冷启动后由 `/vision` 写入） |

愿景体系**不是** goal-tree、progress% 或 Goal 审计台账。

## 规则权威（必读）

| 文件 | 角色 |
|------|------|
| [alignment.md](alignment.md) | 愿景对齐契约与门禁（P-006 操作细则权威） |
| [consumer-checklist.md](consumer-checklist.md) | 与 alignment MUST 表同表的操作勾选 |
| [../standalone-bootstrap.md](../standalone-bootstrap.md) | MUST 表第三同步点（standalone / core-only 核对） |
| （本 README） | 规则面说明 + **本仓实例索引** |

原则全文：`docs/architecture/principles.md`（P-001～**P-006**）。

## 本仓实例索引（schema-ui-core）

> 下列为本仓库**已落盘**的愿景实例，不是 core 预装模板。

| 文件 | 状态 / 说明 |
|------|-------------|
| [charter.md](charter.md) | **active** · `schema-ui-core-admin-foundation@0.1.0`；`primary_workspace` = workspace-001-mvp-admin-foundation |
| [plans/VP-001-mvp-admin-foundation.md](plans/VP-001-mvp-admin-foundation.md) | **closed** · lead: workspace-001-mvp-admin-foundation |
| [plans/VP-002-production-admin-foundation.md](plans/VP-002-production-admin-foundation.md) | **closed**（2026-08-04）· lead: workspace-002-production-admin-foundation |
| [dual-track-contract.md](dual-track-contract.md) | 双线分支维护契约（`F-V003` → `fixed`；方向 3 VP 前置约束） |
| [roadmap.md](roadmap.md) | 组合编排索引 |
| [revisions.md](revisions.md) | Charter 修订台账（`VR-*`） |
| [reviews.md](reviews.md) | Vision Review 台账（`VRev-*`） |
| [workspaces.md](workspaces.md) | 工作区贡献图（1 primary + 1 delivery） |
| [protocol-inventory-v2.7.0.md](protocol-inventory-v2.7.0.md) | 固定上游协议实施清单（`F-V001` 证据） |
| [../workspace-001-mvp-admin-foundation/](../workspace-001-mvp-admin-foundation/) | 实现层 primary 工作区 · Root `GOAL-001-mvp-admin-foundation` |
| [../workspace-002-production-admin-foundation/](../workspace-002-production-admin-foundation/) | 实现层 delivery 工作区 · Root `GOAL-001-production-admin-foundation` · VP-002 lead |

模板（冷启动复制源）：`docs/templates/vision/charter.md`、`vision-plan.md`。

## 完整安装

缺 active Charter 或本目录 `alignment.md` → **不完整安装**（仅允许引导补齐）。  
分发 Skills 时还须存在 canonical [`docs/contracts/`](../contracts/)（见 alignment MUST）。  
当前安装状态以 [consumer-checklist.md](consumer-checklist.md) 为准，**不要**仅凭本目录有文件就宣称完整安装。

## 入口分工

| 入口 | 层 | 用途 |
|------|----|------|
| `/vision` | 决策 | Charter / VP / Review 响应 / re-align / 结构选型 |
| `/vision-audit` | 交叉 | 独立 Vision Review（只写 `reviews.md`） |
| `/govern` | 实现 | 开区、Root、子目标、Goal finding 响应 |
| `/audit` | 交叉 | Goal `03-audit` independent |
