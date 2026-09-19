# Grok Build · 独立交叉审计（GOAL-004 R3）

仓库根 `/audit`。grok-4.6 · reasoning high。`source: independent`。不改 status/progress/goal-tree。

## 目标

`[workspace-039-version-maintenance-diagnostics] GOAL-004-r3-version-and-diagnostics`

对照 GOAL-002 D-001 T-4 与本目标 D-001。

## 必读

1. GOAL-004 五件套、`D-001-r3-implementation-scope`、`E-002`、`A-001`
2. `apps/web/src/app/version-chip.tsx`、`version-chip.test.tsx`、`App.tsx` 接线
3. 确认未改 `kernel/profile.go` 默认集、未新建模块、未改 `protocol/upstream/**`

## 核验

- 无 `monitoring.read` 不 fetch
- 只展示 version 不展示 commit
- QUICKSTART URL 为冻结 GitHub blob
- 诊断入口仅在有 system-monitoring 页面时出现
- mvp 默认无 chip

## 输出

verdict + Findings。可写 `03-audit/A-002-*.md` 并更新索引。
