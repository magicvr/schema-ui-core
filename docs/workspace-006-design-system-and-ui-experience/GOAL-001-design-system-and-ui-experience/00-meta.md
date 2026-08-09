---
id: GOAL-001-design-system-and-ui-experience
title: 现代设计系统与 Schema 驱动 UI/UX 体验产品化
status: active
parent: null
created: 2026-08-09
updated: 2026-08-09
version: 0.4.1
progress: 4/5
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

**视觉方向（D-004 accepted · 2026-08-09）**：Stitch 定稿为过程输入；仓库摘要 [attachments/visual-direction-stitch-summary.md](./attachments/visual-direction-stitch-summary.md)。**不**把 `code.html` 当生产源。

**重开与 closeout-ready**：A-006/D-006 曾废止过早 D-005。S2/S3 按 D-004 重做（GOAL-003 E-002；commits `f16dc9f` / `5716df9`）；独立审 A-008 **pass**；编排 A-009 闭合 F-VUI-001/002。**D-007 已 superseded**（不得以目标意图冒充关门签字）。Root 现为 **closeout-ready**：`active`，开放 required = 0，待用户**显式**书面确认后立 D-008 再 `done`。

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

- [x] **S1**：Token / 主题 / shadcn primitives — GOAL-002；F-002 fixed（A-005）；主路径消费见 S2（F-VUI-004 fixed）。
- [x] **S2**：Renderer 视觉重构 — 桌面密表 / 移动卡片；`recordView` Drawer/Sheet；表单/展示 primitives（GOAL-003 C1；A-008 F-VUI-001 fixed）。
- [x] **S3**：Shell 与工作流 — topbar + ~256 sidenav + 登录设计系统面（GOAL-003 C2；A-008 F-VUI-002 fixed）；移动汉堡为子能力。
- [x] **S4**：状态与反馈 — GOAL-004。
- [ ] **S5**：视觉回归 + fork Token 示例 + 过程关门 — GOAL-005 局部（fork/回归）有效；**过程关门待用户显式确认（D-008）**，故本检查点暂不勾选。

### 阶段 ↔ VP 退出判据映射

| 阶段 | 主要服务的 VP 退出判据 | 证据（现行） |
|------|------------------------|--------------|
| S1 | exit 1 Token / 主题 | GOAL-002；A-005 |
| S2 | exit 2 Renderer 纳入面视觉 | GOAL-003 E-002；A-008 |
| S3 | exit 3 Shell | GOAL-003 E-002；A-008 |
| S4 | exit 4 状态生命周期 | GOAL-004 |
| S5 | exit 5–6 回归诚实 + 过程可关门 | GOAL-005 + 本轮回归绿；过程关门待 D-008 |

## 纲领路线图

| 阶段 | 名称 | 状态 | 说明 |
|------|------|------|------|
| S1 | Token / 主题 / primitives | **完成** | 基建 + 主路径消费 |
| S2 | Renderer 钉死 type 视觉重构 | **完成** | F-VUI-001 fixed |
| S3 | Shell 与工作流 | **完成** | F-VUI-002 fixed |
| S4 | 状态与反馈 | **完成** | GOAL-004 |
| S5 | 回归 / fork / 关门 | **过程待确认** | closeout-ready；待 D-008 |

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 现行 `apps/web` Token/主题/shadcn 基线清单 | S1 方案冻结 | S1 前 | 盘点 | **closed** | — | E-002 + inventory |
| I-002 | required | S1 Token 语义分层与命名 | S1 方案冻结 | S1 决策时 | Root 决策 | **closed** | — | D-002/D-003 |
| I-003 | non-blocking | 主范例页清单（S4） | S4 验收 | S4 前 | 对照范例 | **closed** | — | GOAL-004 |
| I-004 | non-blocking | 对比度是否升格退出分母 | 可选 | 任意 | P-004 | **open** | 默认否 | F-V019 路径 b |
| I-005 | required | 目标态视觉方向是否冻结 | S1 对照；S2/S3 呈现 | S1 实施前宜齐 | Stitch + Root 决策 | **closed** | — | D-004 + E-004 |

## 非目标（本 Root）

- **不**扩张 `I-PROTO-FULL-001` disposition；不新增 registry type 冒充视觉交付。
- **不**交付订单、钱包、类目、通知等业务领域模块。
- **不**重开 VP-003 架构迁移、不恢复长期双线、不引入热插拔插件市场。
- **不**在本项目内重新定义或替代上游协议语义。
- **不**修订 Goal Governance 核心方法论。
- **不**将 Stitch 导出 HTML 接入生产主线。

## 派生进度展示

`progress: 4/5` 由上方 S1～S5 五个等权检查点派生（S1–S4 勾选）。progress 仅为展示；不放行阶段、不关闭 finding、不覆盖信息门禁，也不自动推导 `status: done`。

**closeout-ready（2026-08-09 · E-007）**：开放 required = **0**（F-VUI-001/002 fixed via A-008/A-009）。Root `status: **active**`。D-007 **superseded**。再次 `done` 仅当用户**显式**书面确认（将落盘 D-008）。F-VUI-007 = accepted-residual（非阻断）。

## 台账布局

本目标使用 ledger 目录：`01-decision/`、`02-execution/`、`03-audit/`。
