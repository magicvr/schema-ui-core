---
id: E-005-r2-checkpoint
doc: execution-entry
parent_goal: GOAL-001-nav-group-collapsible
status: recorded
date: 2026-09-07
created: 2026-09-07
updated: 2026-09-07
version: 0.1.0
---

# E-005 · R2 Git checkpoint

## 已发生事实

1. 在 R2 self A-005 `pass` 与本地 grok build A-006 `pass` 合并后，创建 Git checkpoint：`41e89f47`（`feat(workspace-034): implement navigation group contract`）。
2. checkpoint scope：kernel NavigationGroup 契约与 fail-closed 冲突、17 个模块 group 声明、Manifest sidebar 归一化、composition/serve 双 assembly、R2 回归测试、en-US/zh-CN group labelKey、workspace-034 治理决策/执行/审计台账。
3. checkpoint 前验证：`go test ./... -count=1`、Web Vitest `97/97` files / `1332/1332` tests、`tsc -b`、`vite build`、`git diff --check` 均通过。pnpm wrapper 的 esbuild ignored-build-script 保护为环境入口问题，未作为代码失败证据。
4. R1/R2 派生进度同步为 `2/5 = 40%`；R3/R4/R5 仍 pending。R3 尚未实施 Shell 折叠/展开、键盘交互、直接 URL 自动展开或状态保持。

## 追溯边界

该 commit 是恢复点，不替代 Goal Audit 或验收结论；后续 R3 方案与实现应在新的 owned paths / checkpoint 中继续记录。
