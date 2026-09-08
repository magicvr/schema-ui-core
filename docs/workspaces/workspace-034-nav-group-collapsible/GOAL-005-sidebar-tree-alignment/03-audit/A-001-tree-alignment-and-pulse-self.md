---
id: A-001-tree-alignment-and-pulse-self
goal_id: GOAL-005-sidebar-tree-alignment
doc: audit-entry
source: self
auditor: current-session
date: 2026-09-08
scope: group content alignment and active secondary/pulse indicator semantics
verdict: pass
parent: GOAL-001-nav-group-collapsible
version: 0.1.0
---

# A-001 · 组内导航左翼与 active pulse self 审计

## 审计摘要

- 展开组内容已取消常态 `border-l`，并移除额外 `ml-2` / `pl-3`；叶子导航左翼与组标题基础列基本平齐。
- active 页面保留左侧静态高亮竖线；有 secondary 时显示副文本，不显示 pulse；无 secondary 时显示右侧 `animate-ping` 光点。
- inactive 页面无 secondary 时不显示 pulse；既有 active/deep-link、折叠、权限与 top/user slot 语义未改变。
- 定向导航测试 14/14、App 集成 11/11、Web 全量 1342/1342、TypeScript 与构建验证通过。

## Findings

| finding | level | 状态 | 证据 / 响应 |
|---|---|---|---|
| F-001 | required | fixed | `apps/web/src/app/nav-groups.test.tsx`：组内无 border/额外缩进；active marker 与左翼对齐断言通过 |
| F-002 | required | fixed | `nav-groups.test.tsx`：有 secondary 时无 `data-navigation-active-dot`；无 secondary 的 Roles 页面有 pulse 与 `animate-ping` |
| F-003 | recommended | fixed | `navigation.test.ts`、`nav-groups-r4.test.ts`、`App.integration.test.tsx`：active/deep-link/slot 回归通过 |

## 结论

本目标无开放 required finding。P1/P2 检查点均有实现与验证证据，可以将 `GOAL-005-sidebar-tree-alignment` 标记为 `done 2/2`。父目标与 GOAL-003/004 历史状态保持不变。

