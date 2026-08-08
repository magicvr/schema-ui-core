---
id: GOAL-001-design-system-and-ui-experience
title: 现代设计系统与 Schema 驱动 UI/UX 体验产品化
status: active
parent: null
created: 2026-08-09
updated: 2026-08-09
version: 0.3.0
progress: 2/5
plan_refs:
  - VP-005-design-system-and-ui-experience
primary_plan: VP-005-design-system-and-ui-experience
serves_summary: 在 I-PROTO-FULL-001 已 include 的契约面上交付 Design Token、shadcn/ui 风格 primitives、Schema Renderer 与 Shell 的产品级视觉与状态工效；不扩张协议覆盖、不做业务域模块。
---

# GOAL-001 · 现代设计系统与 Schema 驱动 UI/UX 体验产品化

## 概述

本 Root 承接 [VP-005 · 设计系统与 UI/UX](../../vision/plans/VP-005-design-system-and-ui-experience.md)（`active` v0.4.0+）。在已关闭的 VP-003 单主线架构、VP-004 贡献契约与 **VP-006 / `I-PROTO-FULL-001` 整份契约面** 之上，交付：

1. 统一 **Design Token** 与深/浅色主题切换（fork 可经 Token 做最小品牌定制）；
2. **Schema Renderer** 对 VP-005 钉死 type / 能力面的视觉与交互升级；
3. **Admin Shell** 与框架级体验（侧栏、面包屑、用户区、Dialog/Toast 一致语言）；
4. 状态全生命周期（Skeleton / Empty / 错误 / 表单异步反馈）；
5. 回归不回退 + 过程可关门。

**范围分母（F-V018 fixed）**：仅 `I-PROTO-FULL-001` include 面与 VP-005 真实 registry type 表；详情=`recordView`，筛选=`table`+search；**禁止** `Detail`/`Filter` 杜撰 Node 名；**禁止**扩张覆盖 disposition。

**默认不进退出分母（F-V019 路径 b）**：WCAG AA 全站、Cmd+K（可选质量 / S3 增强）。

**视觉方向（D-004 accepted · 2026-08-09）**：Stitch 定稿为过程输入（本地 `raw/stitch-vp005-visual-refs/exports/stitch_schema_ui_core_admin_console/`）；仓库摘要 [attachments/visual-direction-stitch-summary.md](./attachments/visual-direction-stitch-summary.md)。**不**因定稿勾选 S1–S5；**不**把 `code.html` 当生产源。

**重开（D-006 · 2026-08-09）**：A-006 认定 S2/S3 视觉 fidelity 未达标、过早 `done`；用户书面要求回退。Root / 工作区回 `active`；开放 required = **F-VUI-001、F-VUI-002**。

权威覆盖表（只读）：[I-PROTO-FULL-001 v1.0.0](../../workspace-005-full-protocol-contract-v2-7-0/GOAL-001-full-protocol-contract-v2-7-0/attachments/I-PROTO-FULL-001-coverage-v2-7-0.md)。

## 愿景对齐

| 字段 | 值 |
|------|----|
| `plan_refs` | `VP-005-design-system-and-ui-experience` |
| `primary_plan` | `VP-005-design-system-and-ui-experience` |
| Charter | `schema-ui-core-admin-foundation@0.2.0` |
| 工作区角色 | `delivery`（lead of VP-005） |
| 工作区 | `workspace-006-design-system-and-ui-experience` |

## 成功边界

### 阶段层（可验收 · 等权检查点 · 对应 VP-005 S1–S5）

- [x] **S1**：Token / 主题 / shadcn primitives — 语义化 Token（D-002/D-003）；深/浅色可切换且关键壳层无持续 FOUC；primitives 可发现；F-002 fixed（A-005）。**说明（A-006）**：主路径尚未广泛消费 Card/Input 等（F-VUI-004 recommended）；S1 勾选仅代表基建，不代表产品观感已对齐 Stitch。
- [ ] **S2**：Renderer 视觉重构 — 须对照 **D-004**：桌面密表 / 移动卡片列表；`recordView` 右栏或 Drawer（移动 Sheet）；Modal 仅短编辑与 Confirm；钉死 type 面可观察视觉升级。**A-006 F-VUI-001 open**：先前「仅 chart Token」不得勾选。
- [ ] **S3**：Shell 与工作流交互 — 须对照 **D-004** 壳气质 + 登录等；移动汉堡抽屉仅为子项，**不足**单独完成 S3。**A-006 F-VUI-002 open**。
- [x] **S4**：状态与反馈 — Skeleton / Empty / 错误与 async 判定统一（GOAL-004）。
- [ ] **S5**：视觉回归 + fork Token 示例 + 过程关门 — fork 示例与回归绿保留在 GOAL-005；**过程关门因 A-006 失效，本检查点取消勾选**，直至 S2/S3 诚实完成且开放 required 闭合后再议。

### 阶段 ↔ VP 退出判据映射

| 阶段 | 主要服务的 VP 退出判据 | 证据（现行） |
|------|------------------------|--------------|
| S1 | exit 1 Token / 主题 | GOAL-002 E-001；A-005 F-002 fixed |
| S2 | exit 2 Renderer 纳入面视觉 | **未达标**（A-006 F-VUI-001）；GOAL-003 已重开 |
| S3 | exit 3 Shell | **未达标**（A-006 F-VUI-002）；GOAL-003 已重开 |
| S4 | exit 4 状态生命周期 | GOAL-004 E-001 |
| S5 | exit 5–6 回归诚实 + 过程可关门 | GOAL-005 局部（fork/回归）有效；**Root 关门无效**（D-005 superseded / D-006） |

## 纲领路线图

| 阶段 | 名称 | 状态 | 说明 |
|------|------|------|------|
| S1 | Token / 主题 / primitives | **完成（基建）** | 勾选保留；F-VUI-004 recommended 待 S2 消费 |
| S2 | Renderer 钉死 type 视觉重构 | **未完成 · 重开** | F-VUI-001 open |
| S3 | Shell 与工作流 | **未完成 · 重开** | F-VUI-002 open；移动抽屉可保留为已交付子项 |
| S4 | 状态与反馈 | **完成** | GOAL-004 done |
| S5 | 回归 / fork / 关门 | **过程未完成** | 禁止再次 done 直至 F-VUI-001/002 闭合 + 用户确认 |

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 现行 `apps/web` Token/主题/shadcn 基线清单 | S1 方案冻结 | S1 前 | 盘点 | **closed** | — | E-002 + `attachments/I-S1-001-ui-baseline-inventory.md` |
| I-002 | required | S1 Token 语义分层与命名 | S1 方案冻结 | S1 决策时 | Root 决策 | **closed** | — | D-002 **accepted**；D-003 |
| I-003 | non-blocking | 主范例页清单（S4） | S4 验收 | S4 前 | 对照范例 | **closed** | — | GOAL-004 |
| I-004 | non-blocking | 对比度是否升格退出分母 | 可选 | 任意 | P-004 | **open** | 默认否 | F-V019 路径 b |
| I-005 | required | 目标态视觉方向是否冻结 | S1 对照；S2/S3 呈现 | S1 实施前宜齐 | Stitch + Root 决策 | **closed** | — | D-004 + E-004 |

## 非目标（本 Root）

- **不**扩张 `I-PROTO-FULL-001` disposition；不新增 registry type 冒充视觉交付。
- **不**交付订单、钱包、类目、通知等业务领域模块。
- **不**重开 VP-003 架构迁移、不恢复长期双线、不引入热插拔插件市场。
- **不**在本项目内重新定义或替代上游协议语义。
- **不**修订 Goal Governance 核心方法论。
- **不**把 Token 脚手架、Stitch 定稿或过窄自审写成「Charter #3 已全部满足」或 S2–S3 已完成。
- **不**将 Stitch 导出 HTML 接入生产主线。

## 派生进度展示

`progress: 2/5` 由上方 S1～S5 五个等权检查点派生（仅 S1、S4 勾选）。progress 仅为展示；不放行阶段、不关闭 finding、不覆盖信息门禁，也不自动推导 `status: done`。

**关门状态（2026-08-09 · D-006）**：D-005 已 **superseded**。Root `status: **active**`。开放 required：**F-VUI-001、F-VUI-002**（A-006）。在其合法闭合（fixed / accepted-residual / user-overruled）并再次经用户书面确认前，**禁止** `status: done`。

## 台账布局

本目标使用 ledger 目录：`01-decision/`、`02-execution/`、`03-audit/`。索引文件保留 frontmatter 与条目表；独立记录使用 `D-NNN-*`、`E-NNN-*`、`A-NNN-*`。
