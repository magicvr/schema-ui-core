---
doc_type: vision-review
id: VRev-011
status: active
source: independent
created: 2026-08-08
updated: 2026-08-09
version: 0.1.1
parent: null
---

# VRev-011 · VP-005 设计系统与 UI/UX 体验意图合理性复审（2026-08-08）

| 字段 | 值 |
|------|-----|
| source | independent |
| auditor | Grok 4.5 · `/vision-audit` |
| scope | `VP-005-design-system-and-ui-experience`（`planned` v0.1.0）；用户关注：意图是否合理、是否存在问题 |
| audit_type | vision-plan |
| verdict | conditional |
| 建议 class | editorial |

### 范围与结论

只读核对：`docs/architecture/principles.md` **P-006**、`docs/vision/alignment.md`、`charter.md` `@0.2.0`、`plans/VP-005-design-system-and-ui-experience.md`（v0.1.0）、closed `VP-003`/`VP-004`、`roadmap.md`、`workspaces.md`、`revisions.md`、既有 `reviews.md`（至 VRev-010，0 open required）、`protocol-inventory-v2.7.0.md`（抽样）、`I-PROTO-001` 覆盖表 `v0.1.3`（路径引用）、前端可观察现状（`apps/web`：Tailwind CSS 变量主题、`components.json` shadcn 配置、`theme-toggle`、Renderer 白名单 node types）。未读 Goal 正文替代愿景证据；未把 `planned` 读成已交付；**未改** Charter / VP / Goal status。

**总判：conditional（1 open required · 2 open recommended）。**

VP-005 作为「同愿景下新纲领波次」的结构选型**合法且合理**：Charter 成功边界第 3 条明确要求前端产品化（Tailwind + shadcn/ui 风格、浅/深色、Linear/Vercel 参考）；VP-003/004 已关闭并留下单主线与贡献契约；在业务域 VP 之前补齐设计系统与 Schema 驱动体验，符合可 fork 优先与组合编排顺序。机读对齐链成立；`planned` 零区合法。

但方向级退出判据 **#2 与建议阶段 S2** 在协议节点词汇与覆盖范围上**未与 I-PROTO-001 冻结白名单对齐**，且「全量」与「常用」互相矛盾——这足以阻断「方向已稳、可无修订激活」的宣称。其余抬升项（Cmd+K、WCAG AA）与交付形态/过程关门缺口记为 recommended。

### 对用户两项关注的独立回答

| 关注 | 独立结论 | 说明 |
|------|----------|------|
| **是否合理** | **方向合理** | 服务 Charter 边界 3；继承 VP-003 单主线 + VP-004 playbook；Non-goals 排除业务域、协议私有扩展、架构重开；与 roadmap 后续业务 VP「继承设计系统」关系正确。 |
| **是否存在问题** | **有，且含 1 条 required** | 退出范围节点命名/覆盖口径错误与自相矛盾（`F-V018`）；Cmd+K/WCAG 等相对 Charter 抬升与可验证性偏弱（`F-V019`）；过程关门与交付形态定名偏薄（`F-V020`）。 |

### 核对事实

| 核对项 | 结论 | 证据 |
|--------|------|------|
| 单愿景 | **pass** | 唯一 `status: active` Charter；`schema-ui-core-admin-foundation@0.2.0` |
| VP→Charter 机读 | **pass** | `vision_ref: schema-ui-core-admin-foundation@0.2.0` 精确匹配 |
| 语义对齐（抽样） | **pass（方向）** | 意图落在成功边界 3（前端产品化）；继承边界 4–5（单主线/Manifest）；Non-goals 对齐 Charter 非目标（业务产品、协议重定义） |
| VP 最小完备（P-006 §6.5） | **pass（骨架）** | 意图、方向级退出判据、`vision_ref`、工作区绑定表（空）、关门记录占位、规划短史均在 |
| planned 零区 | **pass** | alignment §5 允许；`lead_workspace: null`；roadmap 一致 |
| 结构选型 | **pass** | 新可关门主题波次 → 新 VP（P-006 §6.6）；不改 Charter 边界；不吸收进 closed VP-003/004 工作区 |
| 前置关闭 | **pass** | VP-001～004 均 `closed`；roadmap 行 5 前置写 VP-003/004 |
| 组合编排同步 | **pass** | `roadmap.md` 已登记 VP-005 `planned`；Charter 仍写「无 active 交付 VP」与 planned 状态一致（未误标 active） |
| 过早交付主张 | **无** | `planned`、0 工作区；未把 Token/Shell/Palette 写成已交付事实 |
| 既有 UI 基线（可行性素材） | **pass（素材）** | `apps/web` 已有 Tailwind v4、CSS 变量浅/深色、`ThemeToggle`、shadcn `components.json`、`ui/button`、Renderer 与 Shell E2E 基础——不是从零 |
| 退出 #2 vs 协议白名单 | **fail → 见 F-V018** | 正文写 Table/Form/**Detail**/**Filter** 与 S2「**全量**协议 Node」；实现冻结白名单为 `form`/`section`/`table`/`grid`/`tabs`/`text`/`recordView`/`actionButton`（无 Detail/Filter type） |
| 既有 Vision required | **pass** | VRev-001～010 open required = 0 |

### 合理性总评（独立立场）

| 维度 | 立场 |
|------|------|
| 为何现在做 | **同意**：架构与贡献契约已闭；Charter #3 产品化体验尚未作为独立可关门波次收口；业务模块前统一体验可降低后续分叉成本。 |
| 波次粒度 | **大体合适，边界需收紧**：S1–S5 建议阶段合理；但退出 #2/#3 若按字面「全量 + Cmd+K + WCAG AA」会变成超大 UI 波次，应在激活前把退出分母钉死到冻结能力集 + 可核验证据。 |
| 是否应用 Charter strategic | **否**：不改目的/非目标；属边界 3 的意图展开，**editorial** 收紧 VP 即可。 |
| 是否可保持 planned | **是**；可继续规划。**不建议**在 `F-V018` 未 editorial 闭合前宣称方向已稳或直接激活开区。 |

### Findings

#### F-V018 · 退出判据 #2 / S2 节点范围与 I-PROTO-001 冻结白名单不对齐，且「全量」与「常用」自相矛盾

- level: `required`
- status: `fixed`（2026-08-09 `/vision`；范围权威升为 VP-006 关闭时 `I-PROTO-FULL-001`，非再钉 v0.1.3）
- severity: high
- impact: 激活后 Root/实施会在「只美化已实现白名单」与「补全协议未实现 UI 面」之间漂移；关门时无法客观判定 exit 2；也可能被读成授权扩张 `I-PROTO-001 v0.1.3` 的 D-COMP/D-TABLE/D-FORM partial 边界。
- finding: |
  1. 退出判据 2 写「Table, Form, **Detail**, **Filter** 等协议常用 Node」；建议阶段 S2 写「Table、Form、Detail、Filter 等**全量**协议 Node」。
  2. 仓库实现与冻结证据中的 Renderer type 白名单为：`form`, `section`, `table`, `grid`, `tabs`, `text`, `recordView`, `actionButton`（见 `apps/web/src/renderer/render.test.ts` 与 I-PROTO-001 §5 引用路径）。**不存在**名为 `Detail` 或 `Filter` 的白名单 node type；详情更接近 `recordView`，筛选/搜索更接近 table 能力面（`search-table` 等），不是独立 Node 标签。
  3. 「全量协议 Node」与「常用」及继承声明「不破坏 `I-PROTO-001 v0.1.3` / 不重新定义协议」**不能同时按字面成立**：inventory 能力全集 ≠ 冻结覆盖子集；partial 域（D-COMP/D-TABLE/D-FORM）仍有明确排除子面。
  4. VP 继承表虽提及 I-PROTO-001，**未**像 VP-003 那样把升级范围钉在冻结白名单 / disposition 表上，导致 exit 2 可解释空间过大。
- evidence:
  - `docs/vision/plans/VP-005-design-system-and-ui-experience.md` §方向级退出判据 2、§建议实现阶段 S2、§继承边界
  - `apps/web/src/renderer/render.test.ts`（`isWhitelistedNodeType`）
  - `docs/workspaces/workspace-001-mvp-admin-foundation/GOAL-001-mvp-admin-foundation/attachments/I-PROTO-001-coverage-draft.md` v0.1.3（D-COMP/D-TABLE/D-FORM include-partial）
  - Charter 非目标：不在本项目内重新定义上游协议语义
- closure: |
  `/vision` **editorial**（建议激活前完成）：
  1. 将 exit 2 / S2 的节点范围**改写为** I-PROTO-001 冻结 Renderer 白名单 type（或显式「当前已实现白名单 ∪ 本 VP 书面新增且同步覆盖表版本的 type」）。
  2. 删除或降级「全量协议 Node」措辞；统一为「冻结白名单内已支持节点的视觉/交互升级」，禁止静默扩张 registry。
  3. 若「Filter/Detail」指 table 搜索面 / `recordView`，用协议真实 type 与能力面命名，避免杜撰 Node 名。
  4. 在继承边界或 exit 5 明示：本 VP **不**扩大 `include-partial` 子集、不恢复 `D-UPLOAD` 等 exclude 域。
  不改 `vision_ref`、不要求 Charter strategic。
- 建议 class: `editorial`

#### F-V019 · 部分退出项相对 Charter #3 抬升，且缺少方向级可验证口径

- level: `recommended`
- status: `fixed`（2026-08-09 `/vision` · 路径 b）
- severity: medium
- impact: 把「参考 Linear/Vercel 的克制体验」抬成硬门禁后，Cmd+K / WCAG AA /「无闪烁」等可能在无测量定义时阻塞关门，或反过来被主观宣称完成。
- finding: |
  Charter 成功边界 3 要求：Tailwind + shadcn/ui 风格、浅/深色、Linear/Vercel **参考**。VP-005 将其产品化为 Token/Renderer/Shell/状态，**方向正确**。
  但 exit 1 将 **WCAG AA** 与「无缝无闪烁」写入方向级退出；exit 3 将 **Command Palette（Cmd+K）** 或「统一快捷操作」写入退出分母。二者均非 Charter 明文成功条件。
  本 finding **不**否定它们作为优质产品目标，而是要求：要么（a）保留为方向退出但补「最小可核验证据」（例如关键表面对比度抽检清单、主题切换无 FOUC 的页面级检查、Palette 的最小命令集）；要么（b）降为 S3/S4 建议阶段或可选加分，**默认不进退出分母**（与 VP-004 可选脚手架先例同构）。
- evidence:
  - `docs/vision/charter.md` 方向级成功边界 3
  - `docs/vision/plans/VP-005-design-system-and-ui-experience.md` exit 1、exit 3
- closure: |
  `/vision` editorial 或激活后 Root 方案冻结时二选一写清路径 (a)/(b)；闭合不要求改 Charter。
- 建议 class: `editorial`

#### F-V020 · 缺少过程可关门 exit 与交付形态定名；Non-goals 笔误

- level: `recommended`
- status: `fixed`（2026-08-09 `/vision`）
- severity: low
- impact: 协作者/AI 可能误判交付仅为「改样式」而无文档/fork 示例；或关门时缺过程门闩表述；笔误降低契约清晰度。
- finding: |
  1. VP-004 exit 5 显式要求 lead Root 完成、开放 required=0、Vision Review 无阻断、用户确认。VP-005 五条退出均偏产品能力，**无**对等的过程可关门条款（虽 alignment §7 仍全局适用，但 VP 正文自洽性弱于前序）。
  2. 未像 VRev-010 后的 VP-004 那样做**交付形态定名**：Token/主题机制的权威落点（仅 `apps/web` 代码 vs 另有 architecture/design 说明）、fork 品牌定制的最小示例形态、是否产出可发现的设计系统短文——均未方向级一锤定音。
  3. Non-goals：「不在前端引入违背…的**私有私有** Schema 扩展」疑似笔误（重复「私有」）。
- evidence:
  - `docs/vision/plans/VP-005-design-system-and-ui-experience.md` 退出 1–5、Non-goals、工作区绑定
  - `docs/vision/plans/VP-004-module-contribution-readiness.md` exit 5 + 交付形态定名（对照）
- closure: |
  `/vision` editorial：补一条过程可关门方向（或声明「过程门闩仅适用 alignment §7 / 默认门禁」）；意图节 1 段交付形态定名；修正「私有私有」。
- 建议 class: `editorial`

### 不构成 fail / 不新开额外 required 的诚实边界

1. 本 `conditional` **不是**对「做设计系统」方向的否决；结构选型与 Charter 对齐成立。
2. 既有 Token/主题脚手架不构成 exit 1 已满足；本意见未做视觉验收。
3. Charter「与工作区/VP 的关系」未逐条枚举 planned VP-005 可接受（权威组合索引在 `roadmap.md`）；不升格 finding。
4. `workspaces.md` 未出现 VP-005 绑定正确（planned 零区）。
5. 独立 Vision Review **不**激活 VP、**不**开区、**不**闭合 finding。

### 声明

本意见不修改 Charter / VP / Goal status 或 progress；required finding 的响应由 `/vision` 协调；实现层执行仍交 `/govern`。独立 Vision Review **不**自行闭合 finding。

### 门禁含义

- Vision Review **open required = 1**（`F-V018`）。
- recommended open：`F-V019`、`F-V020`。
- **允许**：保持 `planned`；只读规划与 editorial 修订讨论。
- **禁止（在 F-V018 合法闭合前）**：以本审查宣称「方向已稳」；将本 VP 作为已就绪 `primary_plan` 直接开区/激活交付；将 exit 2 读成授权扩张协议覆盖。
- 闭合 `F-V018` 后：仍须 `/vision` + 用户确认方可 `active` 与 `/govern` 开区；本 `conditional` 不自动变为实施放行。

### 响应（对独立意见 · 组合裁决追加 · 不闭合 finding）

| date | actor | summary |
|------|-------|---------|
| 2026-08-08 | `/vision` · 用户书面 | **不回溯改写**本报告原 verdict/findings。用户裁决：协议目标为 `schema-ui-docs@v2.7.0` **整份契约**；新建 [VP-006](../plans/VP-006-full-protocol-contract-v2-7-0.md) 为当前组合主意图；**VP-006 closed 前禁止 VP-005 激活/视觉实施**。VP-005 → `0.2.0` 写入硬门闩并改继承 VP-006。`F-V018`/`F-V019`/`F-V020` **仍 open**（VP-005 解冻前须按全量覆盖表重写 exit 再闭合）；本追加 **不是** fixed / residual / overruled。 |

### 响应（对独立意见 · VRev-011 findings 闭合 · 2026-08-09）

| date | actor | summary |
|------|-------|---------|
| 2026-08-09 | `/vision` · 用户指令「处理 F-V018（及 F-V019/F-V020），再决定是否解冻 VP-005」 | **不回溯改写**原 verdict `conditional` 与 finding 正文。**F-V018 → `fixed`**：VP-005 → **v0.3.0** editorial——exit 2 / S2 钉死 `I-PROTO-FULL-001` 真实 registry type 表（布局/数据/表单控件）；明确 **详情=`recordView`、筛选=`table`+search 能力面**，禁止 `Detail`/`Filter` 杜撰 Node 名；禁止借视觉波次扩张覆盖 disposition；删除模糊「全量协议 Node」作为分母。范围权威 = VP-006 closed 时覆盖表（非 v0.1.3 子集）。**F-V019 → `fixed`（路径 b）**：WCAG AA 与 Cmd+K **默认不进**方向级退出分母（可选质量/S3 增强）；Charter #3 参考体验保留。**F-V020 → `fixed`**：交付形态定名表；exit 6 过程可关门；Non-goals 无「私有私有」笔误。VP-006 硬前置标为**已满足**；本 scope **0 open required、0 open recommended**。 |
| 2026-08-09 | `/vision` · 用户选择「现在激活」 | 在 F-V018/019/020 已 fixed 前提下，VP-005 **`planned` → `active`（v0.4.0）**。同步 roadmap / Charter 关系节（VR-009 editorial）/ workspaces 注 / overview。`lead_workspace` 仍 `null`；物理 scaffold 交 **`/govern`**（建议 `workspace-006-design-system-and-ui-experience`，slug 须用户确认）。**禁止**在 closed workspace-003/004/005 吸收本意图。激活 **不**宣称视觉产品化已交付。 |
