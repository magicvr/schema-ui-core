---
doc_type: vision-plan
id: VP-005-design-system-and-ui-experience
title: 现代设计系统与 Schema 驱动 UI/UX 体验产品化
status: closed
vision_ref: schema-ui-core-admin-foundation@0.2.0
lead_workspace: workspace-006-design-system-and-ui-experience
created: 2026-08-08
updated: 2026-08-10
version: 0.5.1
parent: null
---

# VP-005 · 现代设计系统与 Schema 驱动 UI/UX 体验产品化

## 状态与门闩（2026-08-09 更新 · 已关门）

| 项 | 值 |
|----|-----|
| status | **`closed`**（2026-08-09 用户书面确认关门；VRev-015 `F-V027`/`F-V028` → `fixed`） |
| **lead_workspace** | **`workspace-006-design-system-and-ui-experience`**（同日 `/govern` scaffold；slug 用户确认；Root `GOAL-001-design-system-and-ui-experience` `done 5/5`） |
| **硬前置 VP-006** | **已满足**：[VP-006](VP-006-full-protocol-contract-v2-7-0.md) 于 2026-08-08 用户书面确认 **`closed`** |
| **Vision required** | **已满足**：VRev-011 `F-V018`/`F-V019`/`F-V020` → `fixed`（v0.3.0 editorial）；VRev-015 `F-V027`/`F-V028` → `fixed`（v0.5.0 editorial） |
| **关门门闩（现行）** | 已 `closed`（2026-08-09 用户书面确认）；保留 workspace-006 历史绑定，默认不接新区；reopen 须用户确认。**禁止**借本 VP 扩张 `I-PROTO-FULL-001` disposition |
| **关门 ≠ 视觉无残余** | residual 点名：F-VUI-007/010/011 `accepted-residual`（Root 台账 A-012）、I-004 open non-blocking（WCAG AA 路径 b）；不阻断关门 |

历史门闩（2026-08-08）：VP-006 未 closed 前禁止视觉实施——**已因 VP-006 closed 失效**。  
「须用户确认激活」门闩：2026-08-09 用户选择「现在激活」后**已解除**。

## 意图

在 [VP-006](VP-006-full-protocol-contract-v2-7-0.md)（整份 v2.7.0 契约，**已 closed**）与既有 [VP-003](VP-003-modular-admin-architecture.md) / [VP-004](VP-004-module-contribution-readiness.md) 底座之上，以 **Linear** 与 **Vercel Dashboard** 的克制、高密度、工作导向现代体验为标杆，建立规范化的 **Design Token** 体系与 **shadcn/ui** 风格基础组件资产，并升级 **Schema Renderer** 与全局 Shell 的视觉与交互层。

使 Schema 驱动的 Admin 页面与 Shell，在 **`I-PROTO-FULL-001` 已 include 的契约面** 上具备现代、精致、高信息密度且支持深/浅色模式的产品级体验，且不破坏协议与单主线架构。

### 交付形态定名（F-V020）

| 是 | 不是 |
|----|------|
| `apps/web` 内可运行的 Token / 主题机制 + shadcn 风格 primitives（权威落点以代码为准；可选 `docs/architecture/` 短文作发现入口） | 独立设计系统 npm 包、Figma 组件库全量同步、营销站动效 |
| Renderer / Shell 对 **已登记 type 与能力面** 的视觉与状态工效升级 | 协议覆盖扩张、新 registry type 冒充「视觉顺带」 |
| fork 品牌定制的**最小示例**（换 Token / 主题变量即可换品牌色的可复核路径） | 多租户主题市场、运行时远程主题商店 |
| 过程台账挂本 VP 的 delivery 工作区（激活后） | 重开 VP-003 架构迁移；修订 Goal Governance 元规则 |

## 继承边界

| 来源 | 本 VP 继承 |
|------|------------|
| Charter `@0.2.0` | 方向级成功边界第 3 条（前端产品化：Tailwind、shadcn/ui 风格、浅/深色、Linear/Vercel **参考**）；非目标保持排除。 |
| **VP-006 / `I-PROTO-FULL-001`** | **硬前置已 closed**。视觉升级范围 = 关闭时覆盖表 **include** 面（现行权威：workspace-005 Root `attachments/I-PROTO-FULL-001-coverage-v2-7-0.md` v1.0.1）。**禁止**回退到仅 `I-PROTO-001 v0.1.3` 子集叙事；**禁止**借本 VP 改写覆盖 disposition 或新增 exclude。 |
| VP-003 / `module-architecture.md` | 单主线模块化、薄内核、后端聚合 Manifest 与 Profile；不破坏单主线。 |
| VP-004 / playbook | 一方模块与 AI 操作契约；UI 扩展须符合 playbook。 |

### Renderer 视觉范围钉死（F-V018 · exit 2 分母）

权威 type 命名必须与 **component registry / 现行 Renderer 白名单** 一致；**禁止**使用非协议 type 名（如 `Detail`、`Filter`）作为退出或阶段范围标签。

| 子面 | 协议 type（视觉升级分母） | 说明 |
|------|---------------------------|------|
| 布局 | `grid`, `section`, `tabs` | 布局容器 |
| 数据与操作 | `text`, `table`, `recordView`, `actionButton`, `statCard`, `chart` | **详情** = `recordView`（不是 `Detail`）；**筛选/搜索** = `table` 能力面 + `form.mode=search` / search-table 行为（不是 `Filter` Node） |
| 表单壳 | `form` | 字段白名单见下行 |
| 表单控件 | `input`, `select`, `inputNumber`, `datePicker`, `dateRangePicker`, `textarea`, `switch`, `checkbox`, `radio`, `cascader`, `checkboxGroup`, `richText`, `password`, `upload` | 与 `form-controls` 白名单 + capability 门禁一致 |
| 能力面（非独立 type） | table selection/toolbar/batch、upload action、reactions 显隐态等 | 随已 include 的 D-TABLE / D-ACT / D-UPLOAD / D-EXPR 做**视觉与反馈**升级，不改协议语义 |

**明确不做（本 VP）**：

- 扩张 `I-PROTO-FULL-001` disposition（不新增 include / 不改 exclude）。
- 引入 registry 外 type 或私有 Schema 扩展冒充上游。
- 用「全量协议 Node」等模糊措辞替代上表分母。

## 方向级退出判据

在 **VP-006 已 closed**、本 VP **已 active** 且 lead 工作区有 Q2 证据时，同时满足下列方向可提议 `closed`：

1. **统一 Design Token 与主题切换机制**  
   语义化 Token（Color, Typography, Radius, Shadow, Spacing）；深/浅色可切换且关键壳层无持续 FOUC；fork 可通过 Token 做品牌定制（最小示例可复核）。  
   **默认不进退出分母**：**WCAG AA** 全站合规（见 F-V019 路径 b）；可作为 S1/S5 质量建议或 Root 可选加分，若用户日后书面升格为退出项须补「关键表面对比度抽检清单」。

2. **Schema Renderer 在 `I-PROTO-FULL-001` include 面上的视觉与交互升级**  
   上表 **已钉死 type / 能力面** 均接入设计系统（可观察的视觉与交互一致性）；**禁止**杜撰 type 名；**禁止**「仅 MVP 旧白名单」偷换范围。具体验收清单激活后由 Root 方案枚举，但分母不得窄于上表、也不得宽于覆盖表 include。

3. **Admin Shell 与框架级体验**  
   侧边栏、面包屑、用户区；Dialog / Toast 等一致语言。  
   **默认不进退出分母**：**Command Palette（Cmd+K）**（F-V019 路径 b）；可作为 S3 可选增强，不阻塞关门。

4. **状态全生命周期标准化**  
   Skeleton / Empty / 错误页与表单异步反馈的一致性（覆盖主范例路径）。

5. **不破坏协议与架构**  
   与 `schema-ui-docs@v2.7.0` + `I-PROTO-FULL-001` 及单主线契约兼容；既有 E2E/Smoke 与关键回归不回退。

6. **过程可关门**  
   lead Root 完成约定范围、开放 required findings = 0、Vision Review 无阻断本 VP 的 open required、用户确认关门。

## 建议实现阶段（非退出判据；**仅本 VP `active` 后**供 `/govern` 参考）

| 阶段 | 阶段目的 |
|------|----------|
| S1 | Token / 主题 / shadcn primitives；可选对比度抽检 |
| S2 | Renderer：**上表钉死 type + 能力面** 的视觉重构（非「全量模糊 Node」） |
| S3 | Shell 与工作流交互；**可选** Cmd+K |
| S4 | 状态与反馈 |
| S5 | 视觉回归、fork Token 定制示例、过程关门 |

## Non-goals（非目标）

- **在用户确认激活前**将本 VP 当作交付焦点实施（硬禁止静默开工）。
- 不制作偏离工作导向的花哨营销动效。
- 不在本项目内开发特定业务领域功能（订单、钱包、通知等）。
- 不引入违背 `schema-ui-docs@v2.7.0` 的私有 Schema 扩展。
- 不破坏 Go 后端聚合 Manifest 与 Profile；不重开已关闭架构迁移。
- 不把本 VP 当作协议覆盖扩张的载体（协议扩张已由 VP-006 收口；变更覆盖表须新决策）。
- 不为 VP 在 `docs/vision/` 建立 Goal 五件套或 progress% 权威。

## 与前后 VP 的关系

| VP | 关系 |
|----|------|
| VP-001～VP-004 | 历史基线；已关闭。 |
| **VP-006** | **硬前置已 closed**；本 VP 在其整份契约面之上做产品化视觉。 |
| 后续业务 VP | 可继承本 VP 设计系统 + VP-006 协议面。 |

## 工作区绑定

| workspace_id | root_goal | role | joined | notes |
|--------------|-----------|------|--------|-------|
| workspace-006-design-system-and-ui-experience | GOAL-001-design-system-and-ui-experience | lead | 2026-08-09 | 用户确认 slug；`/govern` scaffold Root + S1–S5 纲领；`vision_role: delivery` |

## 关门记录

| date | outcome | summary | evidence_links | residuals |
|------|---------|---------|----------------|-----------|
| 2026-08-09 | **closed**（用户书面确认） | 用户书面「确认关门」（VRev-015 F-V027 fixed）；S1–S5 全部交付；vitest 616/616 + build exit 0 + Playwright e2e 2/2 回归全绿；独立审 A-001/A-004/A-011/A-002(S4+S5) 与编排响应闭环 | `docs/workspace-006-design-system-and-ui-experience/goal-tree.md`（Root done 5/5；GOAL-002 6/6、003 2/2、004 3/3、005 2/2）；Root D-008 / A-012 / E-010；GOAL-005 E-001/E-002（616 tests + e2e 2/2）；`reviews/VRev-015-vp005-closeout-readiness.md` | F-VUI-007/010/011 `accepted-residual`（Root 台账 A-012）；I-004 `open non-blocking`（F-V019 路径 b，WCAG AA 不进退出分母） |

## 规划修订短史

| 日期 | 版本 | 变更 |
|------|------|------|
| 2026-08-08 | `0.1.0` | 初创：设计系统与 UI/UX 体验产品化；`planned`。 |
| 2026-08-08 | `0.2.0` | 用户裁决：协议目标纠正为整份 v2.7.0 契约（VP-006）；**硬阻塞**本 VP 激活/视觉实施直至 VP-006 closed；继承与退出改为依赖 VP-006 覆盖表。 |
| 2026-08-09 | `0.3.0` | `/vision` 响应 VRev-011：**F-V018** exit 2/S2 钉死 `I-PROTO-FULL-001` 真实 type（禁 Detail/Filter 杜撰名）；**F-V019** 选路径 b（WCAG AA / Cmd+K 默认不进退出分母）；**F-V020** 补交付形态定名 + 过程 exit 6 自洽 + Non-goals 笔误已清。VP-006 硬前置标为已满足；实施门闩改为「须用户书面确认激活」。仍 `planned`。 |
| 2026-08-09 | `0.4.0` | 用户书面确认「现在激活」：`planned` → **`active`**；解除激活门闩。`lead_workspace` 仍 `null`；物理 scaffold 交 `/govern`（slug 须用户确认）。未宣称视觉产品化已交付。 |
| 2026-08-09 | `0.4.1` | `/govern` 开区：用户确认 slug `workspace-006-design-system-and-ui-experience`；Root `GOAL-001-design-system-and-ui-experience`；`lead_workspace` 绑定；S1–S5 纲领落盘（`0/5`）。激活/开区 **不**宣称视觉已交付。 |
| 2026-08-09 | `0.5.0` | `/vision` 响应 VRev-015：用户书面「确认关门」；`active` → **`closed`**；关门记录落盘（exit 1–6 ↔ 证据映射 + residual 点名）；roadmap / workspaces / Charter 关系节原子同步（VR-011 editorial）。`F-V027`/`F-V028` → `fixed`。 |
| 2026-08-10 | `0.5.1` | editorial 同步 VP-006 现行覆盖权威指针至 `I-PROTO-FULL-001` v1.0.1；保留本 VP 已关门状态与原始协议范围。 |
