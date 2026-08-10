---
id: GOAL-001-design-system-and-ui-experience
doc: execution-entry
record_id: E-002
status: recorded
parent: null
created: 2026-08-09
updated: 2026-08-09
version: 0.1.0
---

# E-002 · I-001 UI 基线盘点（Token / 主题 / shadcn）

## 2026-08-09 · 只读盘点事实

### 已发生事实

1. 用户 `/govern` 指令：焦点 `[workspace-006] GOAL-001`，盘点 `apps/web` Token/主题/shadcn 基线（I-001），并提议 S1 Token 命名决策（I-002）。
2. 对 `apps/web` 源码与配置做了只读扫描（未改业务代码）。
3. 落盘基线附件：`attachments/I-S1-001-ui-baseline-inventory.md`（v0.1.0）。
4. **结论摘要**（详见附件）：
   - Tailwind **v4** + shadcn **new-york / neutral / cssVariables**；权威 CSS = `src/index.css`。
   - 已有 **Color + Radius** 语义 CSS 变量（oklch）与 `.dark` 配对；`@theme inline` 映射到 Tailwind color/radius。
   - 主题：`html.dark` + `localStorage.theme` + `ThemeToggle`；`main.tsx` 启动应用；**无** index.html 同步脚本（FOUC 风险）。
   - shadcn **仅** `ui/button.tsx`；无独立 tailwind.config。
   - **缺** Typography / Shadow 语义 token、destructive/chart 色、统一 spacing 语义层。
   - `account/tokens.ts` 为 JWT，与 Design Token 无关。
5. **未**实施 Token 重构、**未**勾选 S1、**未**改 `progress`（仍 `0/5`）。

### 证据

| 主张 | 路径 |
|------|------|
| 基线全文 | `attachments/I-S1-001-ui-baseline-inventory.md` |
| CSS 变量 | `apps/web/src/index.css` |
| shadcn 配置 | `apps/web/components.json` |
| 主题启动/切换 | `apps/web/src/main.tsx`、`components/theme-toggle.tsx` |
