---
id: A-001-active-marker-self
goal_id: GOAL-004-sidebar-active-indicator
doc: audit-entry
source: self
auditor: current-session
date: 2026-09-08
scope: Sidebar group spacing and active page marker/secondary behavior
verdict: pass
parent: GOAL-001-nav-group-collapsible
version: 0.1.0
---

# A-001 · active 指示与组间距 self 审计

## 审计摘要

- `NavigationItems` 的非 horizontal 组间距已从 `space-y-6` 收紧为 `space-y-2`，desktop sidebar 与 mobile drawer 共用同一渲染组件。
- active 页面保留最左侧静态 `bg-primary` 高亮竖线；不使用动画。
- 页面 secondary 仍由注册 projection 提供；有 secondary 时显示在右侧，无 secondary 时不渲染右侧元素，不显示闪烁光点或空占位。
- active/deep-link、折叠状态和 permission projection 未改变；top/user slot 的 projection 未被改写。

## Findings

| finding | level | 状态 | 证据 / 响应 |
|---|---|---|---|
| F-001 | required | fixed | `apps/web/src/app/nav-groups.test.tsx` 验证 `space-y-2`、active marker、secondary 与 `.animate-ping` 缺省；3 个导航测试文件 13/13 通过 |
| F-002 | recommended | fixed | `tsc -b` 与 `npm run build` 通过；Web 全量 99 files / 1342 tests 通过 |

## 结论

本目标无开放 required finding。P1/P2 检查点均有实现与验证证据，可以将 `GOAL-004-sidebar-active-indicator` 标记为 `done 2/2`。父目标与 GOAL-003 历史状态保持不变。

