---
id: GOAL-001-design-system-and-ui-experience
doc: decision-entry
record_id: D-002
status: accepted
parent: null
created: 2026-08-09
updated: 2026-08-09
version: 0.1.2
---

## D-002 · S1 Token 语义分层与命名约定

### 触发

- I-001 基线已盘点（E-002 / `attachments/I-S1-001-ui-baseline-inventory.md`）。
- I-002 required：S1 方案冻结前须约定 Color / Typography / Radius / Shadow / Spacing 分层与命名。
- 用户 `/govern` 要求提议；2026-08-09 用户选择 **「采纳 D-002 全文」** → 本决策 **`accepted`**。

### 已采纳决定

**原则：在现有 shadcn + Tailwind v4 CSS 变量体系上扩展，禁止并行第二套 token 真相源。**

| # | 决定项 | 建议内容 |
|---|--------|----------|
| 1 | **权威落点** | Design Token 唯一运行时权威 = `apps/web/src/index.css`（`:root` / `.dark` + `@theme inline`）。可选短文发现入口（S1 末或 S5）：`docs/architecture/design-tokens.md`（非第二真相源）。 |
| 2 | **Color** | **保留并文档化**现有 shadcn 语义名：`--background`、`--foreground`、`--card`、`--primary`、`--secondary`、`--muted`、`--accent`、`--border`、`--input`、`--ring` 及其 `*-foreground`。S1 实施时**增量**增加：`--destructive` / `--destructive-foreground`；展示用 `--chart-1`…`--chart-5`（供 `statCard`/`chart`）。色值继续 **oklch**。 |
| 3 | **Typography** | 新增语义层：`--font-sans`、`--font-mono`。字阶 **依赖 Tailwind 默认 `text-*` scale**（S1 不重映射 `--text-xs`…）；**禁止** Renderer 硬编码 px。标题层级：Shell `text-sm`/`font-medium`，页面标题 `text-base`/`text-lg`。**收口依据**：D-003 / F-006。 |
| 4 | **Radius** | **保留** `--radius` + 现有 `--radius-sm|md|lg`。组件默认：控件 `rounded-md`，卡片/对话框 `rounded-lg`。 |
| 5 | **Shadow** | 原始语义：`--elevation-sm|md|lg`（`:root`/`.dark`，深色降不透明度）；`@theme` 映射 `--shadow-*: var(--elevation-*)` 生成 `shadow-*` utility。**禁止**同名自引用。替换散落 `shadow-xl` 为语义消费。**修订依据**：D-003 / F-002（A-002/A-003）。 |
| 6 | **Spacing** | **不**另建完整 spacing 矩阵；消费 **Tailwind 默认 spacing scale**（4px 基）。可选单一密度 token：`--density` 不强制 S1；若引入，仅文档化，默认 `comfortable`。 |
| 7 | **命名纪律** | 公共 API = **CSS 自定义属性**（`--*`）+ Tailwind 语义 utility（`bg-primary`）。禁止在业务/Renderer 中引入 `ds.color.primary.500` 式第二命名空间，除非同时生成到 `index.css`。 |
| 8 | **深/浅色** | 继续 **class 策略**（`html.dark`）+ `localStorage.theme`。S1 须补 `index.html` **同步**内联引导脚本（读 theme / prefers-color-scheme），目标：关键壳层无持续 FOUC。 |
| 9 | **fork 品牌定制最小面** | fork 只覆盖 Color（+ 可选 `--radius`）在 `:root`/`.dark`；Typography/Shadow 默认跟随。示例形态：一份「覆盖 3～5 个主色变量即可换品牌」的注释块或 `index.css` 旁 `brand.example.css` 片段（S5 完整示例）。 |
| 10 | **与 branding API 边界** | `GET /api/branding`（siteTitle/logoUrl）**不属于** Design Token；S1 不把 logo 色解析进 token 表。 |
| 11 | **不进退出分母** | WCAG AA 全站合规仍按 F-V019 路径 b；S1 可选 3～5 表面对比度抽检清单，**默认不**阻塞 S1 完成。 |

### 为什么（相对基线）

1. 基线已是 shadcn 语义变量 + Tailwind v4 `@theme`——重建 token 系统成本高且破坏现有 `bg-primary` 消费面。
2. 缺口集中在 **Typography / Shadow / 少数语义色 / FOUC / primitives 数量**，扩展比替换更符合 VP-005 交付形态。
3. fork 面收敛到 Color（+ radius）可复核、可文档化。

### 未选方案

| 方案 | 未选原因 |
|------|----------|
| 引入 Style Dictionary / 独立 JSON token 管线作为 S1 必须项 | 超出「可运行代码权威」最小交付；可后续 non-blocking |
| 全面换成 Material / 其他命名 | 与已落地 shadcn new-york 冲突；回归成本高 |
| 把 JWT `account/tokens.ts` 或 branding API 并入 Design Token | 概念混淆；职责分离 |

### 闭合 I-002

| 日期 | 动作 |
|------|------|
| 2026-08-09 | 用户书面采纳全文 → D-002 **accepted**；I-002 **closed**（证据：本文件 + 会话确认「采纳 D-002 全文」） |

**允许**：进入 S1 方案冻结与实施。  
**仍禁止**：在无实施证据前勾选 S1 检查点或宣称 exit 1 已满足。

### 后续修订

| 日期 | 动作 |
|------|------|
| 2026-08-09 | D-003 合并响应 A-001/A-002/A-003：修订 §3 Typography 字阶路径、§5 Shadow elevation 双层映射；F-002 决策半程锁定，完整闭合仍待实施证据。 |

### 依赖证据

- I-001 / E-002 / `attachments/I-S1-001-ui-baseline-inventory.md`
- VP-005 exit 1、交付形态定名、F-V019 路径 b
- D-003；A-001 / A-002 / A-003 / A-004
