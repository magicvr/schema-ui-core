---
id: E-003-verification
goal_id: GOAL-004-sidebar-active-indicator
doc: execution-entry
status: recorded
date: 2026-09-08
parent: GOAL-001-nav-group-collapsible
version: 0.1.0
---

# E-003 · 导航回归与构建验证

## 已发生事实

- `apps/web/src/app/nav-groups.test.tsx`、`navigation.test.ts`、`nav-groups-r4.test.ts` 定向回归：3 个文件、13 个测试通过。
- Web 全量 Vitest：99 个测试文件、1342 个测试通过。
- TypeScript：`tsc -b` 通过。
- Web 构建：`npm run build` 通过，Vite 产物构建完成；构建生成的 conformance claim 工作树副作用未纳入本目标变更。
- 既有 active/deep-link、折叠状态、权限过滤与 protocol fixture 回归保持通过。

## 结果边界

本轮未修改 API、Manifest、recordView 或分组成员/顺序；只调整共用 Web NavigationLink 的视觉指示与间距。

