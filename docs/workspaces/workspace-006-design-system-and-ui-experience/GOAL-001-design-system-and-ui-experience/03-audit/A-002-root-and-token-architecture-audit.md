---
id: GOAL-001-design-system-and-ui-experience
doc: audit-entry
record_id: A-002
status: active
parent: null
created: 2026-08-09
updated: 2026-08-09
version: 0.1.0
---

## A-002 · 独立审计：Root 整体与 Token 架构取舍（2026-08-09）

- **source**：independent
- **auditor**：Codex GPT-5.1
- **类型** / **scope**：goal-definition + design-plan / Root 整体；重点复审 D-002「不另建第二套 Token 系统」及 S1 方案冻结门禁
- **verdict**：conditional

### 范围与区间

审计 [workspace-006 Root](../00-meta.md) 的愿景对齐、范围、S1～S5 路线图、信息门禁、现有决策/执行/审计证据，并重点复审 D-002：在现有 shadcn + Tailwind v4 CSS 变量体系上扩展，以 `apps/web/src/index.css` 作为唯一运行时 Token 权威，不为 S1 引入并行 JSON / Style Dictionary 真相源。

本意见只覆盖 Root 定义与当前 design-plan；S1～S5 尚无实施完成主张，因此不审计阶段完成、视觉质量或关门。

### 成果（有证据）

- `workspace.md`、`goal-tree.md` 与 Root `00-meta.md` 对同一 Root、canonical 范围、`active / 0/5`、`VP-005` 绑定保持一致；Root 为唯一 `parent: null` 目标。
- VP-005 `vision_ref` 精确匹配现行 Charter；Root 的目的、非目标、S1～S5 与 VP-005 退出方向一致，未扩张协议 disposition。
- Root 已按 P-001 给出五个串行纲领阶段，并将 progress 明确限制为派生展示。
- I-001 有可核对基线盘点，I-002 有用户采纳的 D-002；执行台账只记录 scaffold 与盘点，没有把 Token 方案写成已实施事实。
- 当前代码确实使用 `apps/web/src/index.css` 中的 `:root` / `.dark` 语义 CSS 变量与 `@theme inline` 映射；`components.json` 也指向该文件并启用 CSS variables。该结构与 Tailwind v4 的 CSS-first theme variables 及 shadcn Tailwind v4 映射方式相容。

### 对照成功标准

1. **不另建第二套 Token 系统的原则合理**：当前交付形态是单个 `apps/web` 内的可运行设计系统，不是独立 npm 包、跨平台设计工具同步或远程主题市场。此时再引入独立 JSON/Style Dictionary 作为必须项，会增加生成、同步和漂移面，却没有已登记的第二消费者需求。
2. **单一权威不等于把所有内容都写成同名变量**：`index.css` 可以同时保存原始语义值与 Tailwind utility 映射，但必须区分两者命名空间，避免映射自引用。
3. **扩展而非替换符合现状**：代码已广泛消费 `bg-primary`、`text-muted-foreground` 等语义 utility；保留 shadcn 命名能降低 S2～S4 迁移成本，并支持深/浅色与 fork 覆盖。
4. **未来升级边界尚未形成当前阻断**：若后续出现多前端包、非 Web 消费者、设计工具导出、跨仓同步或机器可验证发布契约，再评估生成式 token source；当前没有证据支持提前承担该复杂度。

### Findings（F-00N）

#### F-002 · required / medium · Shadow 命名与 Tailwind `@theme` 映射存在方案歧义

- **证据**：D-002 同时规定唯一权威采用 `:root` / `.dark` + `@theme inline`，并要求新增 `--shadow-sm`、`--shadow-md`、`--shadow-lg` CSS 变量；Tailwind 的 `--shadow-*` 本身就是生成 `shadow-*` utility 的 theme namespace。
- **风险**：若实施为 `@theme inline { --shadow-sm: var(--shadow-sm) }`，会形成自引用/无效映射；若只在 `:root` / `.dark` 定义同名变量，则 utility 生成与主题覆盖意图也不够明确。当前 D-002 尚未唯一决定原始值名与 Tailwind alias 的关系，却已将 I-002 标为 closed 并允许 S1 方案冻结。
- **影响门禁**：S1 方案冻结与 S1 完成；在该映射被明确并用构建/渲染证据验证前，不应勾选 S1。
- **必改**：由 `/govern` 在 D-002 响应或后续决策中明确一种无自引用方案，例如以 `--elevation-sm|md|lg` 保存可被 `.dark` 覆盖的原始语义值，再在 `@theme inline` 映射 `--shadow-sm|md|lg: var(--elevation-*)`；或选择另一种同样可核对的单一权威结构。实施证据须证明 `shadow-sm|md|lg` utility 可生成并在深/浅色下取到预期值。

#### F-003 · recommended / medium · 状态与图表颜色的消费闭环需纳入 S1 证据

- **证据**：I-001 已识别 success/warning/chart 缺口；D-002 只明确 destructive 与 chart 1～5。当前 Renderer 仍有 `emerald-500` 成功反馈与按 HSL 公式生成的图表 stroke；confirm/modal 仍有 `bg-black/40` 与 `shadow-xl`。
- **风险**：若只声明 token 而不迁移这些消费点，S2/S4 仍会保留与品牌/深浅色脱节的视觉值，“统一 Design Token”会成为定义层而非消费闭环。
- **建议**：S1 实施计划明确 success/warning/info 是否进入公共语义面；至少把当前已存在的 success、chart、overlay、shadow 硬编码列入消费矩阵与阶段验收。若某类颜色有意保留为非 Token，记录理由和边界。

#### F-004 · recommended / low · 为单一 CSS 权威声明升级触发条件

- **证据**：D-002 将独立 JSON / Style Dictionary 延为 non-blocking，但未写重新评估条件。
- **风险**：未来出现跨包/跨平台消费者时，团队可能因“禁止第二套”而继续把 Web CSS 当作不适合的交换格式，或另行手抄形成真正的第二真相源。
- **建议**：在 S1 发现入口或后续架构短文中声明：多包/非 Web 消费、设计工具导出、跨仓同步或机器发布契约任一成为 required 时，重新评估“生成源 → index.css”单向产物架构；这不改变当前不引入第二套真相源的结论。

### 必改项汇总

- **F-002**：明确 Shadow 原始语义变量与 Tailwind `--shadow-*` theme namespace 的无自引用映射，并以 S1 实施证据验证。该项在合法闭合前阻断 S1 完成，不阻断继续做有界实施/验证。

### 与既有意见的异同

- 与 A-001 一致：不为当前单应用 S1 引入第二套 JSON / Style Dictionary Token 真相源是合理选择；保留 shadcn + Tailwind v4 消费面优于推翻重建。
- 比 A-001 更严格：A-001 将 `index.css` 可读性列为唯一 recommended finding；本次结合 Tailwind `--shadow-*` namespace 与当前硬编码消费点，识别出 F-002 required 实施歧义和 F-003 消费闭环建议。因此本次 verdict 为 `conditional`，并不推翻“不另建第二套”的原则。
- 两条意见不构成 P-004 结论冲突：A-001 的 `pass` 针对原则合理性；本意见同样认可原则，仅对可实施的方案完整性增加门禁条件。开放 finding 取并集。

### 结论 + 建议给编排器/用户的下一步

Root 的愿景对齐、范围、路线图和事实诚实性总体合格；“不另建第二套 Token 系统”的方向是合理且适配当前单应用交付形态的。问题不在是否需要 JSON/Style Dictionary，而在单一 CSS 权威内部仍须把命名映射与实际消费闭环设计完整。

建议用 `/govern` 响应 A-002：优先修正 F-002 的 Shadow 映射决策，再把 F-003 纳入 S1 实施/验收清单。未闭合 F-002 前不得勾选 S1；本意见不要求改变 Root `status: active` 或 `progress: 0/5`。

### 声明

本意见不修改 status/progress、方案正文或 goal-tree；响应由 /govern 处理。
