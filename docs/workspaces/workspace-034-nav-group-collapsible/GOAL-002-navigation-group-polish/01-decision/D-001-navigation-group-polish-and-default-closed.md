---
id: D-001-navigation-group-polish-and-default-closed
doc: decision-entry
parent_goal: GOAL-002-navigation-group-polish
source: user-confirmed + /govern
status: accepted
date: 2026-09-07
created: 2026-09-07
updated: 2026-09-07
version: 0.1.0
---

# D-001 · 导航分组语义样式与默认折叠行为

## 用户要求

用户明确提出两项增量：

1. 优化导航菜单分组设计的样式，保持现有视觉语义工程化，并为导航分组增加语义样式；
2. 默认分组关闭，除非当前访问的是分组内的直接页面或该分组下的内页/动态路径。

## 方案

- 复用现有 Tailwind/shadcn 风格的 `muted`、`accent`、`border`、focus ring、间距与圆角语义，不引入新的主题变量或硬编码业务色。
- 分组 header 使用层级容器和原生 button：提供明确的边界、hover/focus、active 状态、折叠指示和子项左侧层级导轨。
- 非 active 分组的默认 open 状态为 `false`；已有 sessionStorage 偏好优先于默认值。
- 当前 child 或 D-005 已登记的父级内页/动态路径 active 时，自动 open 优先于关闭偏好；离开 active 路由后不强制覆盖用户手动关闭偏好。
- 继续使用 `schema-ui:nav-groups:v1`，不改 Manifest 协议、不增加 NavGroup key/id、不改变 top/user/Dashboard 语义。

## 未选方案

- 不使用新的主题 token 或模块专属颜色，避免把单一导航增量变成设计系统分叉。
- 不使用 localStorage 或服务端状态，避免引入跨用户/Profile 的陈旧 UI 状态。
- 不把所有分组强制展开；默认关闭是本次明确的产品行为。

## 依据与边界

- 父目标 R3 已完成的 active/deep-link/sessionStorage 机制与 D-005。
- 当前请求是同一 VP-034 目标边界内的后续 UI/默认行为增量；父目标 `GOAL-001-nav-group-collapsible` 保持历史 `done 5/5`，不重开。
