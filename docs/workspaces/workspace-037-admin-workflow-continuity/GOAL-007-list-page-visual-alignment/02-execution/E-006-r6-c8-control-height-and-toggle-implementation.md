---
id: E-006-r6-c8-control-height-and-toggle-implementation
doc: execution-entry
status: recorded
goal_id: GOAL-007-list-page-visual-alignment
created: 2026-09-18
updated: 2026-09-18
parent: GOAL-001-admin-workflow-continuity
version: 1.0.0
---

# E-006 · 完成 C8 控制高度、折叠开关 token 与空展开抑制

2026-09-18，按 D-004 完成 R6 C8 的三项 UI 修正：

1. **统一页面 actions 高度**（`apps/web/src/renderer/schema-table.tsx`）：schema toolbar 触发器由 `h-9 rounded-md bg-primary px-3 text-sm` 改为 `h-8 ... text-xs`；“列配置”触发器（`data-saved-view-columns-trigger`）由 `px-2.5 py-2` 改为显式 `h-8`，其下拉面板偏移由 `top-10` 同步为 `top-9`。两者实测高度均为 32px。
2. **折叠开关语义 token**（`apps/web/src/index.css`、`apps/web/src/components/list-filter-panel.tsx`）：新增 `--control` / `--control-foreground` 语义值（`:root` 浅色 `oklch(0.955 0 0)` / `oklch(0.205 0 0)`；`.dark` 深色 `oklch(0.17 0 0)` / `oklch(0.87 0 0)`，深色取比 `--card` 更暗的“下沉”值以贴合范例页 `#121215` 相对容器 `#141414` 的关系），并在 `@theme inline` 增加 `--color-control` / `--color-control-foreground` 别名。展开/收起按键改用 `bg-control` / `text-control-foreground`，与“重置”的透明底描边样式明确区分。
3. **空展开抑制**（`apps/web/src/components/list-filter-panel.tsx`）：引入响应式槽位表 `COLLAPSED_SLOTS`（base/sm/md/lg × withActions/withoutActions）作为折叠容量与可见性类的**同一真相源**；仅当 `items.length > collapsedCapacity` 时渲染展开/收起按键与操作单元。新增 `useTier()` 通过 `matchMedia` 监听断点并在无 `matchMedia` 时回退 `innerWidth`。

回归更新：`schema-table.test.tsx` 新增“单行筛选不渲染展开按键与操作单元”用例、原操作单元用例改为 5 个筛选项（超过 lg 容量 4）；`search-form-filters.test.tsx` 固定为窄屏 tier 以继续覆盖折叠交互；`theme.test.ts` 新增 `--control`/`--control-foreground` 结构守卫（双层声明 + `@theme` 别名 + 非自引用）。

验证事实：

- `npm exec -- tsc --noEmit --pretty false`：通过。
- 定向 Vitest：5 个文件、80 个测试通过。
- 全量 Web Vitest：111 个文件、1420 个测试通过。
- `git diff --check`：通过。
- 真实 Chromium 核对（临时预览页，核对后已删除）：
  - 高度：列配置触发器与 toolbar 触发器在 480/700/900/1280px 均为 32px（相等）。
  - token：展开按键底色 `oklch(0.955 0 0)` 对重置按键 `oklch(1 0 0)`；深色下重置为 `oklch(0.145 0 0)`，token 生效。
  - 空展开抑制：2 个筛选项时 900px 与 1280px 无展开按键、480px 与 700px 有（确有隐藏项）；5 个筛选项时 1280px 有展开按键。
  - 无 console 错误。
- Git checkpoint：本轮 C8 以 `dc4b5c5b` 提交；只暂存本轮 owned paths，未使用 `git add -A`。

C8 实现与回归已完成；C6 修订审计与 Root/VP 的 R6 完成投影待后续审计记录。R5-I-004、Root 和 VP 仍保持开放。
