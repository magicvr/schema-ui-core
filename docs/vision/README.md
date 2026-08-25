---
title: docs/vision · 规则面与本仓实例索引
status: active
created: 2026-07-29
updated: 2026-08-25
parent: null
version: 0.10.0
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
> **权威组合状态**以 [roadmap.md](roadmap.md) / [reviews.md](reviews.md) / 各 `plans/VP-*.md` 为准；本表为发现入口（VRev-013 · F-V024）。

| 文件 | 状态 / 说明 |
|------|-------------|
| [charter.md](charter.md) | **active** · `schema-ui-core-admin-foundation@0.2.0`；`primary_workspace` = workspace-001-mvp-admin-foundation |
| [plans/VP-001-mvp-admin-foundation.md](plans/VP-001-mvp-admin-foundation.md) | **closed** · lead: workspace-001-mvp-admin-foundation |
| [plans/VP-002-production-admin-foundation.md](plans/VP-002-production-admin-foundation.md) | **closed** · lead: workspace-002-production-admin-foundation |
| [plans/VP-003-modular-admin-architecture.md](plans/VP-003-modular-admin-architecture.md) | **closed** · lead: workspace-003-modular-admin-architecture |
| [plans/VP-004-module-contribution-readiness.md](plans/VP-004-module-contribution-readiness.md) | **closed** · lead: workspace-004-module-contribution-readiness |
| [plans/VP-006-full-protocol-contract-v2-7-0.md](plans/VP-006-full-protocol-contract-v2-7-0.md) | **closed**（2026-08-08 用户书面确认）· lead: workspace-005-full-protocol-contract-v2-7-0（整份 v2.7.0 契约；`I-PROTO-FULL-001` v1.0.1） |
| [plans/VP-005-design-system-and-ui-experience.md](plans/VP-005-design-system-and-ui-experience.md) | **closed**（2026-08-09 用户书面确认）· lead: workspace-006-design-system-and-ui-experience |
| [plans/VP-007-localization-and-system-settings.md](plans/VP-007-localization-and-system-settings.md) | **closed**（2026-08-09 用户书面确认）· lead: workspace-007-localization-and-system-settings |
| [plans/VP-008-admin-module-readiness-and-foundation-convergence.md](plans/VP-008-admin-module-readiness-and-foundation-convergence.md) | **closed**（2026-08-10 用户书面确认；`go` 签发）· lead: workspace-008-admin-module-readiness |
| [plans/VP-011-admin-functional-modules.md](plans/VP-011-admin-functional-modules.md) | **closed**（2026-08-18 有界）· lead: workspace-011-admin-functional-modules（四档能力地图上提 roadmap） |
| [plans/VP-012-shared-cross-module-contracts.md](plans/VP-012-shared-cross-module-contracts.md) | **closed**（2026-08-19 完整 · 首波）· lead: workspace-012-shared-cross-module-contracts |
| [plans/VP-013-store-dialects.md](plans/VP-013-store-dialects.md) | **closed**（2026-08-21 有界 · 架构 A1）· lead: workspace-013-store-dialects |
| [plans/VP-014-object-storage.md](plans/VP-014-object-storage.md) | **closed**（2026-08-21 有界 · 架构 A2）· lead: workspace-014-object-storage |
| [plans/VP-015-observability.md](plans/VP-015-observability.md) | **closed**（2026-08-22 有界 · 架构 A4）· lead: workspace-015-observability |
| [plans/VP-016-key-rotation-and-backup.md](plans/VP-016-key-rotation-and-backup.md) | **closed**（2026-08-22 有界 · 架构 A5）· lead: workspace-016-key-rotation-and-backup |
| [plans/VP-017-outbound-mail.md](plans/VP-017-outbound-mail.md) | **closed**（2026-08-24 现行渠道分母再关门 · v0.5.0）· lead: workspace-017-outbound-mail |
| [plans/VP-018-account-email-identity.md](plans/VP-018-account-email-identity.md) | **closed**（2026-08-24 解冻当日关门 · v1.0.0）· lead: workspace-018-account-email-identity |
| [plans/VP-019-iam-recovery.md](plans/VP-019-iam-recovery.md) | **active**（2026-08-25 激活 · Admin 功能 IAM：密码策略 / 邀请 / 自助恢复；VRev-043 independent `pass`） |
| [plans/VP-009-production-hardening.md](plans/VP-009-production-hardening.md) | **active** · 共享基架**持续安全与健壮性程序** · lead: workspace-009-production-hardening（Root 长期容器） |
| [plans/VP-010-design-implementation-conformance.md](plans/VP-010-design-implementation-conformance.md) | **active** · 设计意图—实现符合性**持续对齐程序** · lead: workspace-010-design-implementation-conformance（Root 长期容器；与 VP-009 正交） |
| [../architecture/module-architecture.md](../architecture/module-architecture.md) | VP-003 终态架构权威（Fx、Profile、Manifest、数据与生命周期边界） |
| [../architecture/module-contribution-playbook.md](../architecture/module-contribution-playbook.md) | VP-004 一方模块贡献 playbook |
| [dual-track-contract.md](dual-track-contract.md) | **done / historical** · Charter `@0.1.0` 双线意图记录；已由 VP-003 取代 |
| [roadmap.md](roadmap.md) | 组合编排索引 |
| [revisions.md](revisions.md) | Charter 修订台账（`VR-*`） |
| [reviews.md](reviews.md) | Vision Review 稳定索引（`VRev-001`～`VRev-043`；open required = 0） |
| [reviews/](reviews/) | 正式报告目录 |
| [workspaces.md](workspaces.md) | 工作区贡献图（1 primary + 17 delivery/lead） |
| [protocol-inventory-v2.7.0.md](protocol-inventory-v2.7.0.md) | 固定上游协议全量实施清单（`F-V001` 证据；整份契约收口见 VP-006） |
| [../workspace-001-mvp-admin-foundation/](../workspaces/workspace-001-mvp-admin-foundation/) | primary · Root `GOAL-001-mvp-VP-006 lead（closed 历史绑定） |
| [../workspace-006-design-system-and-ui-experience/](../workspaces/workspace-006-design-system-and-ui-experience/) | delivery · VP-005 lead（closed 历史绑定） |
| [../workspace-007-localization-and-system-settings/](../workspaces/workspace-007-localization-and-system-settings/) | delivery · VP-007 lead（closed 历史绑定） |
| [../workspace-008-admin-module-readiness/](../workspaces/workspace-008-admin-module-readiness/) | delivery · VP-008 lead（closed 历史绑定；`go` 已签发） |
| [../workspace-009-production-hardening/](../workspaces/workspace-009-production-hardening/) | lead · **VP-009**（active 长期程序） |
| [../workspace-010-design-implementation-conformance/](../workspaces/workspace-010-design-implementation-conformance/) | lead · **VP-010**（active 长期程序· VP-001 |
| [../workspace-002-production-admin-foundation/](../workspaces/workspace-002-production-admin-foundation/) | delivery · VP-002 lead（closed 历史绑定） |
| [../workspace-003-modular-admin-architecture/](../workspaces/workspace-003-modular-admin-architecture/) | delivery · VP-003 lead（closed 历史绑定） |
| [../workspace-004-module-contribution-readiness/](../workspaces/workspace-004-module-contribution-readiness/) | delivery · VP-004 lead（closed 历史绑定） |
| [../workspace-005-full-protocol-contract-v2-7-0/](../workspaces/workspace-005-full-protocol-contract-v2-7-0/) | delivery · **VP-006 lead**（closed 历史绑定） |
| [../workspace-011-admin-functional-modules/](../workspaces/workspace-011-admin-functional-modules/) | lead · VP-011（closed 历史绑定） |
| [../workspace-012-shared-cross-module-contracts/](../workspaces/workspace-012-shared-cross-module-contracts/) | lead · VP-012（closed 历史绑定） |
| [../workspace-013-store-dialects/](../workspaces/workspace-013-store-dialects/) | lead · VP-013（closed 历史绑定） |
| [../workspace-014-object-storage/](../workspaces/workspace-014-object-storage/) | delivery · VP-014（closed 历史绑定） |
| [../workspace-015-observability/](../workspaces/workspace-015-observability/) | delivery · VP-015（closed 历史绑定） |
| [../workspace-016-key-rotation-and-backup/](../workspaces/workspace-016-key-rotation-and-backup/) | delivery · VP-016（closed 历史绑定） |
| [../workspace-017-outbound-mail/](../workspaces/workspace-017-outbound-mail/) | delivery · VP-017（closed · 现行分母再关门历史绑定） |
| [../workspace-018-account-email-identity/](../workspaces/workspace-018-account-email-identity/) | delivery · VP-018（closed 历史绑定） |
| [../workspace-019-iam-recovery/](../workspaces/workspace-019-iam-recovery/) | lead · **VP-019**（active · 2026-08-25 开区） |

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
| `/vision-audit` | 愿景交叉审 | independent VRev 报告 + 索引；不改 status |
| `/govern` | 实现 | 工作区目标推进、审计响应、放行/关门 |
| `/audit` | 目标交叉审 | Goal `03-audit` independent |
