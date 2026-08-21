---
id: E-004
goal_id: GOAL-005-w4-long-content-presentation
title: 执行 · S6 关门响应：A-003 F-3 浏览器版面点验 + 截断列宽修复
created: 2026-08-13
updated: 2026-08-13
parent: GOAL-001-design-implementation-conformance
version: 0.1.0
---

# E-004 · S6 关门响应：F-3 浏览器点验与修复（2026-08-13）

## 背景

A-003（independent · grok build · grok 4.6 · high）F-3（P2 recommended）：挤列消失仅有 jsdom 类/属性断言，无真实版面核验；且内层 `span` 的 `max-w` 在 `table-layout: auto` 下可能不压低列宽，存在残余风险。建议关门前在 `/roles` 桌面视口点验 + 打开详情确认换行；若内层 span 无效，把 `max-w-[16rem]` 或 `min-w-0` 放到 `td`。

## 点验过程与发现

1. 新增 e2e spec `apps/web/e2e/w4-long-content-spotcheck.spec.ts`（`APP_PROFILE=admin`，Playwright 自动拉起 API + Web，桌面视口 1440×900，登录 admin/admin → Roles 页）。
2. **首轮点验失败**：桌面表格容器 `scrollWidth(1134) > clientWidth(1118)`，且当时截断只在内层 span（`max-w-[16rem]`）——证实 grok 预判：`table-layout: auto` 下内层 span 的 max-width **不约束**列宽，长串仍按整串 min-content 撑列（Permissions/Menus 列宽实测远超 256px）。
3. **修复（按 F-3 建议）**：`data-table.tsx` 桌面 `<td>` 在 `column.truncate === true` 时追加 `max-w-[16rem]`（width cap 落在 td 层；span 保留 `truncate` 做省略号与 title）。
4. **复验通过**（e2e 全绿，实测数据）：
   - Permissions 列宽 = **256px**（≤257 断言）、Menus = **256px**；截断单元格 `title` === 全文文本；
   - 兄弟列 ID/Key/Name 各 ≈ **67px**（≥40 断言），不再被挤出；
   - 详情 Drawer：`dl.scrollWidth ≤ clientWidth`（无横向溢出）；permissions 值以 `", "` 连接呈现（`users.read, users.write, …`）并自动换行。
5. jsdom 层同步补断言：truncate 单元格的 `closest("td")` 带 `max-w-[16rem]`（`data-table.test.tsx`）。

## 门禁复验

- `apps/web` vitest 全量：**48 文件 879 通过**（含补强断言）。
- `apps/web` build（`tsc -b` + vite）：**通过**。
- e2e：`w4-long-content-spotcheck.spec.ts` **1 passed**（admin profile 真实浏览器）。

## 结论

F-3 → **fixed**（点验 + td 层修复 + 回归证据）；「挤列消失」「详情自动换行」现为真实浏览器版面断言，非仅 jsdom 类断言。e2e spec 作为本波验收证据**保留入库**（测试-only，无生产行为变化；diff 面在 D-001 §2 之外新增 1 个 e2e spec 与 data-table.test.tsx 补强断言，属 F-3 响应）。
