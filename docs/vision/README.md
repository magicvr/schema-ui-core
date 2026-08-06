---
title: docs/vision · 规则面与本仓实例索引
status: active
created: 2026-07-29
updated: 2026-08-07
parent: null
version: 0.7.0
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

### Vision Review 台账（v0.13+）

| 路径 | 角色 |
|------|------|
| [reviews.md](reviews.md) | **稳定索引** + 当前 open required 投影 + 条目表 |
| [reviews/](reviews/) | 平铺正式报告 `VRev-NNN-<slug>.md`（self / independent 共用编号序列） |
| `docs/templates/vision/reviews-index.md` | 新索引复制源 |
| `docs/templates/vision/review.md` | 单条 VRev 报告复制源 |

**2026-08-07**：`VRev-001`～`VRev-010` 已自 legacy inline **无重编号**迁移到 `reviews/`；新记录只写报告并更新索引，禁止再 inline 追加正文。

## 本仓实例索引（schema-ui-core）

> 下列为本仓库**已落盘**的愿景实例，不是 core 预装模板。

| 文件 | 状态 / 说明 |
|------|-------------|
| [charter.md](charter.md) | **active** · `schema-ui-core-admin-foundation@0.2.0`；`primary_workspace` = workspace-001-mvp-admin-foundation |
| [plans/VP-001-mvp-admin-foundation.md](plans/VP-001-mvp-admin-foundation.md) | **closed** · lead: workspace-001-mvp-admin-foundation |
| [plans/VP-002-production-admin-foundation.md](plans/VP-002-production-admin-foundation.md) | **closed**（2026-08-04）· lead: workspace-002-production-admin-foundation |
| [plans/VP-003-modular-admin-architecture.md](plans/VP-003-modular-admin-architecture.md) | **active** · lead: workspace-003-modular-admin-architecture；完整单主线模块化终态，建区不等于实现完成 |
| [../architecture/module-architecture.md](../architecture/module-architecture.md) | VP-003 终态架构权威（Fx、Profile、Manifest、数据与生命周期边界） |
| [dual-track-contract.md](dual-track-contract.md) | **done / historical** · Charter `@0.1.0` 双线意图记录；已由 VP-003 取代 |
| [roadmap.md](roadmap.md) | 组合编排索引 |
| [revisions.md](revisions.md) | Charter 修订台账（`VR-*`） |
| [reviews.md](reviews.md) | Vision Review 稳定索引（`VRev-001`～`VRev-010`；0 open required） |
| [reviews/](reviews/) | 正式报告目录（已迁入 VRev-001～010；新条目只写此处） |
| [workspaces.md](workspaces.md) | 工作区贡献图（1 primary + 2 delivery；VP-003 已绑定 lead） |
| [protocol-inventory-v2.7.0.md](protocol-inventory-v2.7.0.md) | 固定上游协议实施清单（`F-V001` 证据） |
| [../workspace-001-mvp-admin-foundation/](../workspace-001-mvp-admin-foundation/) | 实现层 primary 工作区 · Root `GOAL-001-mvp-admin-foundation` |
| [../workspace-002-production-admin-foundation/](../workspace-002-production-admin-foundation/) | 实现层 delivery 工作区 · Root `GOAL-001-production-admin-foundation` · VP-002 lead |
| [../workspace-003-modular-admin-architecture/](../workspace-003-modular-admin-architecture/) | 实现层 delivery 工作区 · Root `GOAL-001-modular-admin-architecture` · VP-003 lead |
| [../workspace-004-module-contribution-readiness/](../workspace-004-module-contribution-readiness/) | 实现层 delivery 工作区 · Root `GOAL-001-module-contribution-readiness` |

模板（冷启动 / Review 复制源）：`docs/templates/vision/charter.md`、`vision-plan.md`、`reviews-index.md`、`review.md`。

## 完整安装

缺 active Charter 或本目录 `alignment.md` → **不完整安装**（仅允许引导补齐）。  
分发 Skills 时还须存在 canonical [`docs/contracts/`](../contracts/)（见 alignment MUST）。  
当前安装状态以 [consumer-checklist.md](consumer-checklist.md) 为准，**不要**仅凭本目录有文件就宣称完整安装。

Skills 消费包当前 pin：**goal-governance `v0.13.0`**（见 `skills/.goal-governance-install.json`）。

## 入口分工

| 入口 | 层 | 用途 |
|------|----|------|
| `/vision` | 决策 | Charter / VP / Review 响应 / re-align / 结构选型 |
| `/vision-audit` | 交叉 | 独立 Vision Review（写 `reviews/VRev-*.md` + 更新 `reviews.md` 索引） |
| `/govern` | 实现 | 开区、Root、子目标、Goal finding 响应 |
| `/audit` | 交叉 | Goal `03-audit` independent |
