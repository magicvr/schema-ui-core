---
id: GOAL-001-design-system-and-ui-experience
title: 现代设计系统与 Schema 驱动 UI/UX 体验产品化
status: active
parent: null
created: 2026-08-09
updated: 2026-08-09
version: 0.1.7
progress: 5/5
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

- [x] **S1**：Token / 主题 / shadcn primitives — 语义化 Token（D-002/D-003）；深/浅色可切换且关键壳层无持续 FOUC；primitives 可发现；暗色主文案对比度不低于定稿 Overview dark 可读性；可选对比度抽检（不默认进退出分母）。**F-002 已 fixed（A-005 · 2026-08-09）。**
- [x] **S2**：Renderer 视觉重构 — VP-005 钉死 type + 能力面接入设计系统；**呈现约束（D-004）**：桌面密表 / 移动卡片列表；`recordView` 为右栏或 Drawer（移动 Sheet）；Modal 仅短编辑与 Confirm；分母不得窄于 type 表、不得宽于 I-PROTO-FULL-001 include。
- [x] **S3**：Shell 与工作流交互 — 桌面侧栏+顶栏；移动汉堡+导航抽屉；用户区；Dialog/Toast 一致语言；与 D-004 壳气质对齐；**可选** Cmd+K（默认不进退出分母）。
- [x] **S4**：状态与反馈 — Skeleton / Empty / 错误页与表单异步反馈在主范例路径一致。statCard/chart/list-table 的 loading 态统一改用 `Skeleton` primitive（`role="status"`），共享纯函数 `resolveAsyncDisplayState` 直接单测；GOAL-004 五件套（C1–C3 全绿，A-001 self 审计 pass）。
- [x] **S5**：视觉回归 + fork Token 示例 + 过程关门 — fork 品牌定制最小示例（`brand.example.css` + 结构测试 + README）；vitest 616/616 全绿、build exit 0、Playwright e2e 2/2 真实通过（`schema-crud.spec.ts`/`shell.spec.ts`，非诚实退化证据）；独立交叉审计 A-002（grok build/grok-4.5/高思考强度）已落盘响应（A-003，2 条 required finding 均 fixed）；GOAL-005 五件套（C1–C2 全绿）。open required（本 Root scope）= 0；**关门本身仍待用户书面确认**，见下方"关门证据摘要"。

### 阶段 ↔ VP 退出判据映射

| 阶段 | 主要服务的 VP 退出判据 | 证据（建区后随进展回填） |
|------|------------------------|--------------------------|
| S1 | exit 1 Token / 主题 | 待 S1 产物；方向输入 = D-004/E-004 |
| S2 | exit 2 Renderer 纳入面视觉 | 待 S2 产物；呈现 = D-004 |
| S3 | exit 3 Shell | 待 S3 产物；壳 = D-004 |
| S4 | exit 4 状态生命周期 | GOAL-004 E-001（Skeleton 统一 + async-state 纯函数 + 单测） |
| S5 | exit 5–6 回归诚实 + 过程可关门 | GOAL-005 E-001/E-002；A-002（independent）+ A-003（响应，fixed）；vitest 616 + build + e2e 2/2 |

## 纲领路线图

| 阶段 | 名称 | 状态 | 说明 |
|------|------|------|------|
| S1 | Token / 主题 / primitives | **完成** | D-002/D-003 accepted；D-004 方向冻；F-002 fixed（A-005）；C1–C6 全绿 |
| S2 | Renderer 钉死 type 视觉重构 | **完成** | chart pie → `--chart-*` Token；overlay/shadow S1 已完成 |
| S3 | Shell 与工作流 | **完成** | 移动汉堡+抽屉（bg-overlay/shadow-lg）；navigate 关闭；shell.test.ts 通过 |
| S4 | 状态与反馈 | **完成** | Skeleton 统一 + 纯判定函数（GOAL-004；A-001 self pass） |
| S5 | 回归 / fork 示例 / 关门 | **完成** | fork brand.example.css + 结构测试；vitest 616 + build + e2e 2/2；A-002 independent（conditional）→ A-003 响应（fixed）；开放 required = 0 |

阶段通常串行；同一纲领阶段内允许并行子目标。建区与视觉定稿 **均不**勾选任何检查点。

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 现行 `apps/web` Token/主题/shadcn 基线清单（可观察路径） | S1 方案冻结 | S1 前 | 盘点 CSS 变量、`components.json`、ThemeToggle、ui/* | **closed** | — | E-002 + `attachments/I-S1-001-ui-baseline-inventory.md`（2026-08-09） |
| I-002 | required | S1 Token 语义分层与命名约定（Color/Typography/Radius/Shadow/Spacing） | S1 方案冻结；S2–S4 消费 | S1 决策时 | Root 决策落盘 | **closed** | — | D-002 **accepted**（2026-08-09）；D-003 修订 §3/§5 |
| I-003 | non-blocking | 主范例页清单（S4 状态验收路径） | S4 验收 | S4 前 | 对照 schemarender 8 页 + Shell 路由 | **closed** | — | S4 覆盖 statCard/chart/list-table（`render.tsx`/`data-table.tsx`）+ 默认 8 范例 + 登录/壳层（GOAL-004 E-001） |
| I-004 | non-blocking | 是否将对比度抽检升格为退出分母 | 仅当用户书面升格 | 任意 | P-004 用户裁决 | **open** | 默认 **否**（F-V019 路径 b） | 默认不进 exit 1；S1 仍以 dark 可读为质量下限（D-004） |
| I-005 | required | 目标态视觉方向是否冻结（可复核参考） | S1 实施对照；S2/S3 呈现 | S1 实施前宜齐 | Stitch 定稿 + Root 决策 | **closed** | — | D-004 **accepted** + E-004 + `attachments/visual-direction-stitch-summary.md`（2026-08-09） |

## 非目标（本 Root）

- **不**扩张 `I-PROTO-FULL-001` disposition；不新增 registry type 冒充视觉交付。
- **不**交付订单、钱包、类目、通知等业务领域模块。
- **不**重开 VP-003 架构迁移、不恢复长期双线、不引入热插拔插件市场。
- **不**在本项目内重新定义或替代上游协议语义。
- **不**修订 Goal Governance 核心方法论。
- **不**把建区、Token 脚手架或 **Stitch 定稿** 写成「Charter #3 已全部满足」或 S1–S5 已完成。
- **不**将 Stitch 导出 HTML 接入生产主线。
- 不为 VP 在 `docs/vision/` 建立 Goal 五件套或 progress% 权威。

## 派生进度展示

`progress: 5/5` 由上方 S1～S5 五个等权检查点派生（S1–S5 均已完成）。progress 仅为展示；不放行阶段、不关闭 finding、不覆盖信息门禁，也不自动推导 `status: done`。

**关门状态（P-004 · 待用户裁决）**：S1–S5 检查点齐备，GOAL-002~GOAL-005 四个子目标 `status: done`，`03-audit.md` 台账（Root 层）与 GOAL-004/GOAL-005 台账开放 required findings 均为 0（GOAL-005 的 A-002 independent conditional 已由 A-003 self 响应 fixed）。**Root `status` 本次不单方面置为 `done`**——按 AGENTS §6b P-004，关门须用户书面确认；本条仅记录"证据齐备，可提请关门"的事实，不代为裁决。

## 台账布局

本目标使用 ledger 目录：`01-decision/`、`02-execution/`、`03-audit/`。索引文件保留 frontmatter 与条目表；独立记录使用 `D-NNN-*`、`E-NNN-*`、`A-NNN-*`。
