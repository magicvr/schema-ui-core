---
title: I-S1-001 · apps/web Token / 主题 / shadcn 基线盘点
status: active
doc_type: info-baseline
created: 2026-08-09
updated: 2026-08-09
parent: GOAL-001-design-system-and-ui-experience
version: 0.1.0
related_info: I-001
---

# I-S1-001 · UI 基线盘点（只读事实）

> 盘点日期：2026-08-09。路径相对仓库根。**不是** S1 完成证明；不宣称 Charter #3 已满足。

## 1. 栈与配置

| 项 | 事实 |
|----|------|
| CSS 框架 | **Tailwind CSS v4**（`tailwindcss` + `@tailwindcss/vite` ^4.1.11） |
| 入口 CSS | `apps/web/src/index.css`：`@import "tailwindcss"` |
| shadcn 配置 | `apps/web/components.json`：`style: new-york`，`baseColor: neutral`，`cssVariables: true`，`css: src/index.css`，`iconLibrary: lucide`，`rsc: false` |
| tailwind.config | **无**独立 config 文件（Tailwind v4 + `components.json` 中 `tailwind.config: ""`） |
| 工具函数 | `apps/web/src/lib/utils.ts`：`cn` = `clsx` + `tailwind-merge` |
| 相关依赖 | `@radix-ui/react-slot`、`class-variance-authority`、`lucide-react` |

## 2. Design Token / CSS 变量（现行）

### 2.1 语义色（`:root` / `.dark`）

权威文件：`apps/web/src/index.css`。

| Token | 层 | 说明 |
|-------|-----|------|
| `--background` / `--foreground` | Color | 页面底/字 |
| `--card` / `--card-foreground` | Color | 卡片面 |
| `--primary` / `--primary-foreground` | Color | 主操作 |
| `--secondary` / `--secondary-foreground` | Color | 次级 |
| `--muted` / `--muted-foreground` | Color | 弱化 |
| `--accent` / `--accent-foreground` | Color | 悬停/强调底 |
| `--border` / `--input` / `--ring` | Color | 边框、输入、焦点环 |
| `--radius` | Radius | 基半径 `0.5rem` |

色值使用 **oklch**；浅色 `:root` 与深色 `.dark` 成对定义。  
**未定义**：`--destructive`、success/warning、chart 系列、显式 shadow/spacing 语义 token、typography 字族/字号 token。

### 2.2 Tailwind v4 `@theme inline` 映射

同一文件将语义变量映射为 `--color-*` 与 radius 派生：

- `--color-background` … `--color-ring` → `var(--*)`
- `--radius-sm` / `--radius-md` / `--radius-lg` ← 基于 `--radius` 计算

由此可用 utility：`bg-background`、`text-primary`、`border-border`、`ring-ring`、`rounded-md` 等。

### 2.3 全局规则

- `@custom-variant dark (&:is(.dark *));`
- `* { @apply border-border; }`
- `body { @apply bg-background text-foreground antialiased; }`

## 3. 主题切换机制

| 项 | 事实 |
|----|------|
| 类策略 | `document.documentElement` 上切换 **`dark` class** |
| 持久化 | `localStorage` key **`theme`** = `"dark"` \| `"light"` |
| 启动应用 | `apps/web/src/main.tsx` `applyStoredTheme()`：读 localStorage，缺省跟 `prefers-color-scheme` |
| 切换 UI | `apps/web/src/components/theme-toggle.tsx`（Sun/Moon + Button outline） |
| FOUC 风险 | `index.html` **无**同步内联主题脚本；主题在 JS 模块执行后才应用 → 可能首帧闪烁 |
| 壳层挂载 | `App.tsx` 使用 `ThemeToggle` |

## 4. shadcn / primitives 资产

| 路径 | 状态 |
|------|------|
| `apps/web/src/components/ui/button.tsx` | **唯一** ui 原语；CVA variants：`default` / `outline` / `secondary` / `ghost`；size default/sm/lg |
| 其他 `ui/*` | **不存在**（无 input、dialog、dropdown、toast、card、tabs 等） |
| `apps/web/src/components/data-table.tsx` | 产品表组件（非 shadcn 生成路径；另有 renderer `schema-table.tsx`） |
| `apps/web/src/components/theme-toggle.tsx` | 主题切换（依赖 Button） |

## 5. Shell / branding（与 Token 相邻）

| 路径 | 职责 |
|------|------|
| `apps/web/src/app/App.tsx` | Admin Shell：侧栏/导航、主题切换、用户区、挂载 `RenderPage` |
| `apps/web/src/app/branding.ts` | 站点标题 + logo URL（API `/api/branding`）；**非** Design Token，是产品品牌投影 |
| `apps/web/src/app/LoginPage.tsx` | 登录页（独立于主 Shell） |

## 6. Renderer 消费方式（抽样）

- 广泛使用 Tailwind 语义类：`bg-background`、`text-muted-foreground`、`border-input`、`bg-primary` 等。
- 部分硬编码：`bg-black/40`（confirm 遮罩）、`shadow-xl`（confirm 卡片）——**未**走 shadow token。
- 表单控件：`form-controls.tsx` 内联 class，未统一为 shadcn `Input`/`Select` 组件。

## 7. 易混淆路径（排除）

| 路径 | 说明 |
|------|------|
| `apps/web/src/account/tokens.ts` | **JWT / 会话 token**，不是 Design Token |
| `apps/web/dist/**` | 构建产物，非权威源 |

## 8. 与 VP-005 S1 的差距摘要（事实 → 后续决策输入）

| 维度 | 已有 | 缺口（供 I-002 / S1） |
|------|------|----------------------|
| Color | shadcn 中性语义对 + 深浅色 | 无 destructive/success/chart；无 fork 品牌色文档化覆盖面 |
| Typography | 仅 `text-sm`/`text-xs` 等 utility 散用 | 无语义字族/字阶 token |
| Radius | `--radius` + sm/md/lg | 可保留；需约定组件消费规则 |
| Shadow | 散用 `shadow-xl` 等 | 无语义 shadow token |
| Spacing | Tailwind 默认尺度 | 无 density / 语义 spacing 层 |
| Primitives | 仅 Button | 缺 dialog/input/toast/card 等 S1/S3 常用原语 |
| 主题 | class + localStorage | 缺 HTML 同步引导（FOUC）；无「无缝」验收清单 |

## 9. 证据路径清单

- `apps/web/src/index.css`
- `apps/web/components.json`
- `apps/web/package.json`
- `apps/web/src/main.tsx`
- `apps/web/index.html`
- `apps/web/src/components/theme-toggle.tsx`
- `apps/web/src/components/ui/button.tsx`
- `apps/web/src/lib/utils.ts`
- `apps/web/src/app/App.tsx`
- `apps/web/src/app/branding.ts`
