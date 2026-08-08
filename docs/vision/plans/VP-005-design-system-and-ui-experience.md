---
doc_type: vision-plan
id: VP-005-design-system-and-ui-experience
title: 现代设计系统与 Schema 驱动 UI/UX 体验产品化
status: planned
vision_ref: schema-ui-core-admin-foundation@0.2.0
lead_workspace: null
created: 2026-08-08
updated: 2026-08-08
version: 0.2.0
parent: null
---

# VP-005 · 现代设计系统与 Schema 驱动 UI/UX 体验产品化

## 状态与硬门闩（2026-08-08 用户裁决）

| 项 | 值 |
|----|-----|
| status | `planned`（**未激活**） |
| **实施门闩** | **[VP-006](VP-006-full-protocol-contract-v2-7-0.md) 未 `closed` 前，禁止**将本 VP 标为 `active`、禁止 `/govern` 为本 VP 开区、禁止启动任何视觉优化 / Design Token / Shell 产品化 / Renderer 换皮实施 |
| 允许 | 仅文档层 editorial（对齐门闩、继承边界、退出判据措辞）；不得当作当前交付焦点 |

本门闩由用户 2026-08-08 书面确认：须先完成 `schema-ui-docs@v2.7.0` **整份契约**可验证兼容（VP-006），再进行本 VP 的视觉产品化，避免在不完整协议面上强化设计债与审计跑偏。

## 意图

在 [VP-006](VP-006-full-protocol-contract-v2-7-0.md)（整份 v2.7.0 契约）与既有 [VP-003](VP-003-modular-admin-architecture.md) / [VP-004](VP-004-module-contribution-readiness.md) 底座之上，以 **Linear** 与 **Vercel Dashboard** 的克制、高密度、工作导向现代体验为标杆，建立规范化的 **Design Token** 体系与 **shadcn/ui** 风格基础组件资产，并全面升级 **Schema Renderer** 与全局 Shell 的视觉与交互层。

使得由 Schema 驱动生成的 Admin 页面以及 Shell 框架，在 **已支持的全量契约面** 上具备现代、精致、高信息密度且支持深/浅色模式的产品级体验，且不破坏协议与单主线架构。

## 继承边界

| 来源 | 本 VP 继承 |
|------|------------|
| Charter `@0.2.0` | 方向级成功边界第 3 条（前端产品化：Tailwind、shadcn/ui 风格、浅/深色、Linear/Vercel 参考）；非目标保持排除。 |
| **VP-006** | **硬前置**：整份 `schema-ui-docs@v2.7.0` 契约覆盖与实现；本 VP 在其 `closed` 后方可激活。视觉升级范围对齐 VP-006 关闭时的覆盖表，**不得**回退到仅 `I-PROTO-001 v0.1.3` 子集叙事。 |
| VP-003 / `module-architecture.md` | 单主线模块化、薄内核、后端聚合 Manifest 与 Profile；不破坏单主线。 |
| VP-004 / playbook | 一方模块与 AI 操作契约；UI 扩展须符合 playbook。 |

## 方向级退出判据

在 **VP-006 已 closed**、且同时满足下列方向、均有工作区 Q2 证据时，本 VP **可以**提议 `closed`：

1. **统一 Design Token 与主题切换机制**  
   语义化 Token（Color, Typography, Radius, Shadow, Spacing）；深/浅色可切换；fork 可通过 Token 做品牌定制。可验证口径在激活后 Root 方案中冻结（含是否将 WCAG AA 计入退出分母）。

2. **Schema Renderer 在全量契约纳入面上的视觉与交互升级**  
   覆盖表 `include` / `include-partial` 所纳入的 node / 控件 / table·form 能力面接入设计系统；禁止用杜撰 type 名或「仅 MVP 白名单」偷换范围。具体 type 列表以 **VP-006 关闭时覆盖表** 为准。

3. **Admin Shell 与框架级体验**  
   侧边栏、面包屑、用户/租户区；Dialog/Toast 等一致语言。Command Palette 等增强项是否计入退出分母在激活后书面冻结。

4. **状态全生命周期标准化**  
   Skeleton / Empty / 错误页与表单异步反馈的一致性。

5. **不破坏协议与架构**  
   保持与 `schema-ui-docs@v2.7.0`（VP-006 关闭时覆盖表）及单主线契约兼容；E2E/Smoke 与关键回归通过。

6. **过程可关门**  
   lead Root 完成约定范围、开放 required = 0、Vision Review 无阻断、用户确认。

## 建议实现阶段（非退出判据；**仅 VP-006 closed 且本 VP active 后**供 `/govern` 参考）

| 阶段 | 阶段目的 |
|------|----------|
| S1 | Token / 主题 / shadcn primitives |
| S2 | Renderer 全量纳入面视觉重构 |
| S3 | Shell 与工作流交互 |
| S4 | 状态与反馈 |
| S5 | 视觉回归与 fork 定制示例 |

## Non-goals（非目标）

- **在 VP-006 完成前实施本 VP**（硬禁止）。
- 不制作偏离工作导向的花哨营销动效。
- 不在本项目内开发特定业务领域功能（订单、钱包、通知等）。
- 不引入违背 `schema-ui-docs@v2.7.0` 的私有 Schema 扩展。
- 不破坏 Go 后端聚合 Manifest 与 Profile；不重开已关闭架构迁移。
- 不把本 VP 当作协议覆盖扩张的载体（协议扩张 = VP-006）。

## 与前后 VP 的关系

| VP | 关系 |
|----|------|
| VP-001～VP-004 | 历史基线；已关闭。 |
| **VP-006** | **硬前置**；未 closed 前本 VP 不得激活/实施。 |
| 后续业务 VP | 可继承本 VP 设计系统 + VP-006 协议面。 |

## 工作区绑定

| workspace_id | root_goal | role | joined | notes |
|--------------|-----------|------|--------|-------|
| — | — | lead | — | **冻结**：不得开区直至 VP-006 `closed` 且用户确认激活本 VP |

## 关门记录

（仅 `closed` 或 `abandoned` 时填写。）

| date | outcome | summary | evidence_links | residuals |
|------|---------|---------|----------------|-----------|
| — | — | — | — | — |

## 规划修订短史

| 日期 | 版本 | 变更 |
|------|------|------|
| 2026-08-08 | `0.1.0` | 初创：设计系统与 UI/UX 体验产品化；`planned`。 |
| 2026-08-08 | `0.2.0` | 用户裁决：协议目标纠正为整份 v2.7.0 契约（VP-006）；**硬阻塞**本 VP 激活/视觉实施直至 VP-006 closed；继承与退出改为依赖 VP-006 覆盖表；修正「仅 I-PROTO-001 v0.1.3」作为视觉范围的叙事。 |
