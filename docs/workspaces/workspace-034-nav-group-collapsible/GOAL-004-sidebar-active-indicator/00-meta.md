---
id: GOAL-004-sidebar-active-indicator
title: Sidebar 分组间距与页面激活指示微调
status: done
parent: GOAL-001-nav-group-collapsible
created: 2026-09-08
updated: 2026-09-08
version: 0.1.0
progress: 100%
plan_refs:
  - VP-034-nav-group-collapsible
primary_plan: VP-034-nav-group-collapsible
serves_summary: 在已完成的 Sidebar Engine 视觉基础上收紧分组间距，并将页面激活态统一为左侧高亮竖线与可选 secondary 副文本，移除右侧闪烁光点。
---

# GOAL-004 · Sidebar 分组间距与页面激活指示微调

## 概述

这是 workspace-034 中继承 `GOAL-001-nav-group-collapsible` 的新 sibling 增量目标，不重开或改写 GOAL-002/GOAL-003 的历史事实。

本轮只处理两项窄范围 UI 调整：

1. 参照 `raw/stitch_schema_ui_core_admin_console` Sidebar Engine 范例的紧凑节奏，缩小导航组与组之间当前过大的垂直间距。
2. 页面 active 状态使用最左侧高亮竖线；右侧不再使用闪烁光点，若注册了 secondary 则在原右侧位置显示副文本，未注册时不显示任何占位。

## 成功标准

- [x] desktop sidebar 与 mobile navigation drawer 的组间距采用参考页相近的紧凑间距。
- [x] active 页面显示左侧高亮竖线；active 页面有 secondary 时右侧显示文字，没有 secondary 时无光点、无空占位。
- [x] inactive 页面、top/user slot 与分组折叠、深链自动展开语义不受影响。
- [x] 相关单测、TypeScript 检查与构建通过，并完成 self 审计记录。

## 纲领路线图

以下 2 个检查点是本目标 progress 的唯一来源，默认等权：

| 检查点 | 目的 | 状态 |
|---|---|---|
| P1 | Sidebar 组间距与页面 active 左竖线/secondary 替换光点实现 | completed |
| P2 | 导航状态回归、构建与 self 审计 | completed |

`progress: 100%` = 2/2 个检查点完成（P1、P2）。progress 仅作展示，不放行阶段、不关闭 finding、不自动推导 `done`。

## 信息需求与门禁（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|---|---|---|---|---|---|---|---|---|
| I-004-001 | required | 组间距应收紧到何种参考节奏？ | P1 | P1 | 对照 Sidebar Engine 范例 `nav space-y-sm`，采用 Tailwind `space-y-2`（8px） | verified（用户指令 + 参考页） | — | D-001 |
| I-004-002 | required | active 页面右侧闪烁光点与 secondary 的替换关系是什么？ | P1 | P1 | 有 secondary 显示文字；无 secondary 不显示光点或占位；左侧 active 竖线始终保留 | verified（用户指令） | — | D-001 |
| I-004-003 | non-blocking | active 竖线是否需在 desktop/mobile 两个 sidebar surface 对齐？ | P2 | P2 | 两个 surface 共用 NavigationLink，并补 DOM/class 回归 | verified | — | D-001；E-002；E-003 |

## 愿景对齐

- Charter：`schema-ui-core-admin-foundation@0.4.0`
- VP：`VP-034-nav-group-collapsible`（`active` · v0.3.0）
- 工作区：`workspace-034-nav-group-collapsible`（`delivery`）
- 父目标：`GOAL-001-nav-group-collapsible`（`done 5/5`）
- 前置增量：`GOAL-003-sidebar-engine-navigation`（`done 3/3`）

## 范围边界

### 纳入

- `NavigationItems` 非 horizontal surface 的组间垂直间距。
- `NavigationLink` 页面 active 左侧竖线、secondary 右侧显示和无 secondary 缺省行为。
- desktop sidebar、mobile drawer、导航单测与既有 active/deep-link/折叠回归。

### 不纳入

- 不改导航分组成员、顺序、注册协议、权限、路由或 sessionStorage。
- 不改 recordView、top/user slot、Dashboard Workspace 注册或其他已完成目标。
- 不恢复闪烁光点；不为无 secondary 页面保留视觉占位。

## 父级关系

- `parent: GOAL-001-nav-group-collapsible`；父目标完整 id 与当前 workspace 相同。
- 目标文件夹在 workspace 根平铺；层级只由 `parent` 表达。

## 台账布局

本目标使用五件套与 `01-decision/`、`02-execution/`、`03-audit/` 三个平铺 ledger 目录。

