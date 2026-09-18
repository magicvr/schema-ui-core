---
id: E-001-feasibility-and-geometry-baseline
doc: execution-entry
status: recorded
goal_id: GOAL-009-list-visual-e2e-guard
created: 2026-09-18
updated: 2026-09-18
parent: GOAL-001-admin-workflow-continuity
version: 1.0.0
---

# E-001 · 可达性确认与真实几何基线

2026-09-18，完成 C1：

## 挂具可跑性

- `go version go1.26.0 windows/amd64`；Chromium（`ms-playwright/chromium-1234`）在缓存中；端口 25080 / 25173 空闲。
- 既有规格冒烟：`APP_PROFILE=admin npx playwright test e2e/w4-long-content-spotcheck.spec.ts` → **1 passed**（19.5s）。挂具（临时 SQLite、`cmd/server`、`npm run dev`）可正常起停。

## 目标页与 profile 可达性

- 全仓 schema 检索：只有 `account.json` 声明 table `props.filters`；同时具备 `mode: search` 与 `props.toolbar` 的页面中，`roles.json` 是唯一能满足全部断言目标者（search form 有 `q` + `system` 两个字段，窄档恰好隐藏 1 项；toolbar 有 Export / New role 两个触发器，可用于高度一致性）。
- 两 profile 实跑探测：`admin` 与 `mvp` 展开侧栏后均可达 `/users`、`/roles`（`mvp` 下分组需展开才渲染链接）。
- 侧栏为 `hidden lg:block`：窄屏下无法通过侧栏导航，故规格先在桌面宽度导航、再收窄视口。

## 真实几何基线（`/roles`）

桌面 1440（lg 档；2 个筛选项 + 操作单元）：

| 观测 | 值 |
|------|-----|
| `[data-list-filter-panel]` / `[data-list-page-actions]` / `[data-table-surface]` | 均存在，自上而下顺序为筛选 → actions → 列表 |
| 页面 actions 子项高度 | columns 32、Export 32、New role 32（**全部相等**） |
| 折叠开关 | **不存在**（2 项 + 操作单元恰好容纳于 lg 行，无隐藏项） |
| 搜索配对 | 同一 `data-filter-item`；`submit.left - input.right = -1`；垂直偏移 0；按钮含 `-ml-px` / `rounded-l-none` |
| 列表 footer | `data-table-footer` 位于 `data-table-surface` 内部；`data-pagination-footer` 存在 |
| 触发器图标 | columns 与 save-view 触发器各含 1 个 `svg` |
| 视图表单 | `data-saved-view-management` 在 `data-saved-views-surface` **内部**，不在 `data-list-page-actions`；与 save 按钮垂直间隔 4px |

窄屏 700（sm 档）：

| 观测 | 值 |
|------|-----|
| 折叠开关 | **存在**；隐藏 1 项、可见 1 项 |
| 开关背景 | `oklch(0.955 0 0)`；重置按钮背景 `oklch(1 0 0)`（**不同**） |
| 开关≠重置 | `resetIsToggle = false` |

上述基线即 C2 断言的依据。C2 见 `E-002`。
