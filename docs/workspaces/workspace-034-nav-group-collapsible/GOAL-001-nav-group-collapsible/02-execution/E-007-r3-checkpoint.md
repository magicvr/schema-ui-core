---
id: E-007-r3-checkpoint
doc: execution-entry
parent_goal: GOAL-001-nav-group-collapsible
status: recorded
date: 2026-09-07
created: 2026-09-07
updated: 2026-09-07
version: 0.1.0
---

# E-007 · R3 Git checkpoint

## 已发生事实

1. 在 A-008 R3 self `pass`、Web 全量回归与构建通过后，创建 Git checkpoint：`6e581ca9`（`feat(workspace-034): add collapsible navigation groups`）。
2. checkpoint scope：Web navigation active/deep-link projection、统一父级映射、Shell collapsible group button、aria/Enter/Space、sessionStorage 状态与容错、R3 专项测试、D-005/E-006/A-008 及 workspace-034 状态投影。
3. checkpoint 前验证：Web Vitest `98/98` files / `1337/1337` tests、`tsc -b`、`vite build`、R3 专项 3/3 与 navigation 7/7、`git diff --check` 均通过。
4. R3 实现不修改协议 NavGroup key/id，不移动 top/user slot；I-034-004 的全量 Profile/route matrix 仍由 R4 负责。

## 追溯边界

该 commit 是 R3 恢复点，不替代 R4 全量迁移回归或 Root 关门审计。
