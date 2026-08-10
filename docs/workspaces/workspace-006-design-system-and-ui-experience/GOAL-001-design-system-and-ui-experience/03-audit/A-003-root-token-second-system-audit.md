---
id: GOAL-001-design-system-and-ui-experience
doc: audit-entry
record_id: A-003
status: active
parent: null
created: 2026-08-09
updated: 2026-08-09
version: 0.1.0
---

## A-003 · 独立审计：Root 与「不另建第二套 Token 系统」决策（2026-08-09）

- **source**：independent
- **auditor**：Grok 4.5（xAI /audit）
- **类型** / **scope**：goal-definition + design-plan / Root 整体定义与 S1 design-plan；**重点**复审 D-002「在现有 shadcn + Tailwind v4 CSS 变量上扩展、禁止并行第二套 token 真相源」
- **verdict**：conditional

### 范围与区间

| 项 | 值 |
|----|-----|
| 工作区 | `workspace-006-design-system-and-ui-experience`（`workspace.md` 已校验：root_goal、canonical、`plan_refs`/`primary_plan`=VP-005、`shared_materials_catalog: none`） |
| 被审目标 | `GOAL-001-design-system-and-ui-experience` |
| audit_type | goal-definition + design-plan |
| 显式关注点 | 不另建第二套 Token 系统是否合理 |
| 不在本 scope | S1～S5 实施完成质量、视觉验收、关门、跨工作区内容比较 |

只读证据：Root 五件套 + ledger；`I-S1-001`；D-001/D-002；E-001/E-002；A-001/A-002；VP-005（交付形态与 exit）；`apps/web/src/index.css`、`components.json`；Renderer 抽样硬编码（confirm/modal/render）。

### 成果（有证据）

1. **工作区与 Root 绑定一致**：`workspace.md` / `goal-tree.md` / Root `00-meta` 同指 `GOAL-001-design-system-and-ui-experience`，`active`，`progress: 0/5`，`primary_plan: VP-005`，delivery 角色；共享资料 `none`，未把跨区材料当成本区证据。
2. **愿景对齐与范围诚实**：Root 目的、S1–S5、非目标与 VP-005 一致；分母钉在 `I-PROTO-FULL-001` include + VP-005 type 表；明确激活/建区 ≠ 视觉产品化已交付。
3. **P-001 路线图与事实边界**：五个等权检查点均未勾选；执行台账仅 scaffold（E-001）与基线盘点（E-002），无虚假完成主张。
4. **信息门禁（S1）**：I-001 closed（E-002 + `attachments/I-S1-001-ui-baseline-inventory.md`）；I-002 closed（D-002 用户采纳全文）。I-003/I-004 为 non-blocking 且未伪称已验证。
5. **代码基线与 D-002 前提可核对**（2026-08-09 只读）：
   - `apps/web/src/index.css`：`:root` / `.dark` 语义色 + `@theme inline` 将 `--color-*` / radius 映射到 `var(--*)`；**无** Typography / Shadow 语义 token；**无** destructive/chart。
   - `components.json`：`cssVariables: true`，`css: src/index.css`，new-york / neutral。
   - 消费侧已依赖语义 utility；仍有硬编码：`bg-black/40`、`shadow-xl`（confirm/modal）、`emerald-*` 成功反馈、`hsl(...)` 图表 stroke（`render.tsx`）。
6. **D-002 原则与未选方案有记录**：权威落点 `index.css`；拒 Style Dictionary/JSON 作为 S1 必须项；拒 Material 换名；拒把 JWT `account/tokens.ts` 或 branding API 并入 Design Token。

### 对照成功标准 · 重点：不另建第二套 Token 系统

#### 判定结论（原则层）

**合理，且与当前产品边界高度匹配。** 不构成 fail 或要求推翻 D-002 原则。

| 对照轴 | 证据 | 评估 |
|--------|------|------|
| VP-005 交付形态（F-V020） | 是：`apps/web` 内可运行 Token；**不是**独立 npm 包 / Figma 全量同步 / 远程主题市场 | 单一 CSS 运行时权威 = 交付形态正中；Style Dictionary 多源管线属于「不是」侧复杂度 |
| VP-005 exit 1 | 语义 Token + 深浅色 + fork 最小可复核路径；权威以**代码**为准 | 在现有 shadcn 变量上扩展即可满足方向，无需第二命名空间 |
| 已落地消费面 | `bg-primary` 等已广泛使用；Color 已是「`:root` 原始语义 → `@theme --color-*` 别名」 | 替换为并行 JSON/`ds.color.*` 会制造双写与迁移税，无第二消费者证明其收益 |
| 基线缺口形态 | 缺口在 Typography/Shadow/destructive/chart/FOUC/primitives **增量**，非「整栈无 token」 | 扩展优于重建 |
| 概念边界 | JWT tokens / branding API 已从 Design Token 排除 | 避免假「第二套」与真职责混淆 |
| 行业「规范 token 管线」反方 | DTCG JSON / Style Dictionary 在多包、设计工具、跨仓同步时价值高 | 当前 **无** 已登记第二消费者、跨平台或发布契约 required；提前引入 = 预支同步与漂移成本 |

#### 必须区分的两层含义（原则通过的前提）

「不另建第二套 **真相源**」**不等于**「文件里只能有一层名字」：

| 允许 / 正确 | 禁止 / 错误 |
|-------------|-------------|
| **同一权威文件**内：`:root`/`.dark` 存可主题覆盖的**原始语义值**；`@theme inline` 只做 **Tailwind 消费别名**（Color 现行：`--background` → `--color-background: var(--background)`） | 并行 `tokens.json` / Style Dictionary / `ds.color.primary.500` 与 `index.css` **双写且互不生成** |
| 可选 `docs/architecture/design-tokens.md` 作发现入口（D-002 已写明非第二真相源） | 文档或 JSON 反过来成为唯一权威、代码手抄漂移 |
| 未来若出现多包/设计工具等 required，再评估「**生成源 → 唯一产物 index.css**」（单向；仍非双真相） | 以「禁止第二套」为借口永久拒绝升级评估，或手抄第二份运行时表 |

**现行 Color 映射已证明「单文件、双层、单真相」可行。** 风险不在原则，而在 D-002 对 Shadow（及类似）是否复用该纪律——见既有 **F-002**。

#### 未选方案复核

| 未选 | 本审计是否同意未选 | 说明 |
|------|-------------------|------|
| Style Dictionary / JSON 作 S1 必须 | **同意未选** | 与 F-V020 最小交付一致；可后置 non-blocking |
| 全面 Material 等换名 | **同意未选** | 与 shadcn 消费面冲突，回归成本无收益 |
| JWT / branding 并入 Design Token | **同意未选** | 职责分离正确 |

### Findings（F-00N）

#### F-002 · required / medium · **仍开放**（本意见复审确认，非新开编号）

- **来源**：A-002 首报；本审计独立复核对代码与 D-002 后**维持 required**。
- **证据**：D-002 §5 新增 `--shadow-sm|md|lg`；Tailwind v4 中 `--shadow-*` 亦为 theme/utility 命名空间。Color 成功模式是 **不同名**（`--primary` vs `--color-primary`）。D-002 未写 Shadow 原始名与 `@theme` 别名的无自引用关系。
- **风险**：`@theme { --shadow-sm: var(--shadow-sm) }` 自引用；或仅 `:root` 同名导致 utility/主题意图不清。
- **影响门禁**：S1 **完成**（勾选检查点 / 宣称 exit 1）；**不**阻断有界实施与映射方案补强。
- **必改**：在 D-002 修订或后续决策中选定并验证一种无自引用结构（例：`--elevation-*` 原始语义 + `@theme --shadow-*: var(--elevation-*)`，或等价可核对方案）；实施证据证明 `shadow-sm|md|lg` 可生成且深/浅色取值符合预期。

#### F-005 · recommended / low · 将 Color 双层映射纪律写成 S1 实施模板

- **证据**：`index.css` 已对 Color 采用「语义原始值 + `--color-*` 别名」；D-002 原则正确但未把该模式提升为 Typography/Shadow/未来 token 的**显式模板**。
- **风险**：实施者误读「单一权威」=「全部挤进 `@theme` 同名变量」，放大 F-002 类问题；或误读为禁止任何 alias 层。
- **建议**：S1 实施说明或短文用一张表固定：权威值在 `:root`/`.dark`；`@theme` 仅 alias；禁止同名自引用；文档/JSON 不得成为第二运行时权威。

#### F-003 · recommended / medium · **仍有效**（复审确认）

- **证据**：`render.tsx` 仍用 `emerald-*` 与 `hsl(...)` 图表色；confirm/modal 仍用 `bg-black/40`、`shadow-xl`；D-002 已定 destructive/chart，未闭合 success/warning/overlay 是否进公共语义面。
- **建议**：S1 消费矩阵至少列入既有硬编码点；不进 Token 的须写边界理由，避免「定义有、消费无」。

#### F-004 · recommended / low · **仍有效**（复审确认）

- **证据**：D-002 将独立 JSON 管线延后，未写重新评估触发条件。
- **建议**：多包/非 Web 消费、设计工具导出、跨仓同步、机器发布契约任一成为 required 时，评估「生成源 → 唯一 `index.css` 产物」；**不**改变当前「不并行双真相」结论。

#### F-006 · recommended / low · Typography 字阶路径的 OR 表述宜在实施前收口

- **证据**：D-002 §3「`@theme` 映射 `--text-xs`… 或依赖 Tailwind 默认 scale」。
- **风险**：低于 F-002（两路径皆可工作），但 S1 验收「禁止 Renderer 硬编码 px」时路径不清会导致证据形态不一。
- **建议**：S1 实施计划二选一写死，并附 1～2 个消费示例；非 S1 完成阻断项。

### 必改项汇总

| ID | 级别 | 摘要 | 阻断 |
|----|------|------|------|
| **F-002** | required | Shadow 原始语义 vs Tailwind `--shadow-*` 无自引用映射 + 实施验证 | **是**（S1 完成） |
| F-005 | recommended | Color 双层纪律作实施模板 | 否 |
| F-003 | recommended | 硬编码消费闭环 / success·overlay 边界 | 否 |
| F-004 | recommended | 单一权威升级触发条件 | 否 |
| F-006 | recommended | Typography 字阶路径收口 | 否 |

**开放 required 计数（本目标意见并集）**：仍为 **1**（F-002）。本意见**不**新增 required。

### 与既有意见的异同

| 意见 | 对「不另建第二套」 | 差异 |
|------|-------------------|------|
| A-001 | pass；原则高度合理 | 仅 recommended 可读性；未深入 Shadow namespace |
| A-002 | 原则合理；verdict conditional（F-002） | 识别映射歧义与消费闭环 |
| **A-003（本）** | **原则独立复审仍通过**；verdict **conditional** 仅因 F-002 仍开放 | 用 VP-005 F-V020 与**现行 Color 双层代码**加固「真相源 ≠ 禁止 alias 层」；新增 F-005/F-006（recommended）；**不**与 A-001/A-002 构成 P-004 冲突 |

开放 findings **取并集**；原则层三审一致：**不另建第二套 Token 真相源 = 合理决策**。

### 信息门禁核对（P-005）

| I-ID | 级别 | 状态 | 与本 scope |
|------|------|------|------------|
| I-001 | required | closed | 决策与审计基线充分 |
| I-002 | required | closed | 命名约定已采纳；**≠** F-002 实施映射已闭合 |
| I-003 | non-blocking | open | 不影响 S1 design-plan |
| I-004 | non-blocking | open | 默认不进 exit；无伪关闭 |

到期且影响本 scope 的 required 信息项：**0**。  
影响 S1 **完成** 的开放 required finding：**F-002**。

### 结论 + 建议给编排器/用户的下一步

1. **重点结论**：「不另建第二套 Token 系统」（不以 Style Dictionary/独立 JSON/`ds.*` 作为并行真相源；在 shadcn + Tailwind v4 CSS 变量上扩展）**合理**，与 VP-005 交付形态、现网消费面和基线缺口形态一致；**无需推翻**。
2. **条件**在于方案可实施性：必须把「单一真相源」落实为与 Color 相同的**双层映射纪律**（F-002），否则 S1 有实施歧义。
3. Root 定义、对齐、诚实性与 S1 信息门禁总体合格；`status: active` / `0/5` 与事实一致，本意见不要求改状态。

**建议 `/govern` 响应句**：

> 响应 A-003（并并集 A-001/A-002）：确认「不另建第二套 Token 真相源」原则维持；优先闭合 **F-002**（Shadow 无自引用映射决策 + 验证计划）；将 F-005/F-003/F-006 纳入 S1 实施清单；F-002 合法闭合前不得勾选 S1。

### 声明

本意见 `source: independent`，不修改 `status` / `progress` / 方案正文 / goal-tree；响应、finding 闭合与推进由 **`/govern`** 处理。
