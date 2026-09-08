---
id: GOAL-005-sidebar-tree-alignment
title: Sidebar 组内导航左翼与 active 光点语义调整
status: done
parent: GOAL-001-nav-group-collapsible
created: 2026-09-08
updated: 2026-09-08
version: 0.1.0
progress: 100%
plan_refs:
  - VP-034-nav-group-collapsible
primary_plan: VP-034-nav-group-collapsible
serves_summary: 细调 Sidebar Engine 组内导航树的左侧对齐与 active 状态：取消展开组的常态竖向导轨，左移组内菜单，并让无 secondary 的 active 页面显示右侧闪烁光点。
---

# GOAL-005 · Sidebar 组内导航左翼与 active 光点语义调整

## 概述

这是 workspace-034 中继承 `GOAL-001-nav-group-collapsible` 的新 sibling 增量目标，不重开或改写 GOAL-002、GOAL-003、GOAL-004 的历史事实。

本轮响应新的视觉反馈：

- 取消分组展开后长期显示在组内导航左侧的竖向边界线。
- 将组内导航整体左移，使 active 页面最左侧视觉效果与组名称左侧基本平齐。
- active 页面有已注册 secondary 时继续显示 secondary；没有 secondary 时显示右侧闪烁光点，二者不并列。

## 成功标准

- [x] 展开组不再显示常态 `border-l` 竖线，组内导航左翼与组名称左侧对齐。
- [x] active 页面有 secondary 时显示副文本且不显示闪烁光点；无 secondary 时显示闪烁光点。
- [x] inactive 页面、组折叠/展开、active deep-link、desktop/mobile 共用导航语义不受影响。
- [x] 相关测试、TypeScript 检查与构建通过，并完成 self 审计。

## 纲领路线图

以下 2 个检查点是本目标 progress 的唯一来源，默认等权：

| 检查点 | 目的 | 状态 |
|---|---|---|
| P1 | 组内导航左翼对齐与 active secondary/pulse 互斥实现 | completed |
| P2 | 导航回归、TypeScript/构建验证与 self 审计 | completed |

`progress: 100%` = 2/2 个检查点完成（P1、P2）。progress 仅作展示，不放行阶段、不关闭 finding、不自动推导 `done`。

## 信息需求与门禁（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|---|---|---|---|---|---|---|---|---|
| I-005-001 | required | 组内常态竖线是否应取消，菜单左翼如何对齐？ | P1 | P1 | 对照用户指令；移除 group content border/indent，并保留组标题 padding 作为对齐基准 | verified（用户指令） | — | D-001 |
| I-005-002 | required | active secondary 与闪烁光点的互斥规则是什么？ | P1 | P1 | secondary 存在时显示文字；缺失时显示 pulse；两者互斥 | verified（用户指令） | — | D-001 |
| I-005-003 | non-blocking | pulse 是否只作用于 sidebar/mobile active 页面？ | P2 | P2 | 复用非 horizontal `NavigationLink`；horizontal top navigation 保持紧凑无 pulse | verified（实现边界） | — | D-001 |

## 愿景对齐

- Charter：`schema-ui-core-admin-foundation@0.4.0`
- VP：`VP-034-nav-group-collapsible`（`active` · v0.3.0）
- 工作区：`workspace-034-nav-group-collapsible`（`delivery`）
- 父目标：`GOAL-001-nav-group-collapsible`（`done 5/5`）
- 前置增量：`GOAL-004-sidebar-active-indicator`（`done 2/2`）

## 范围边界

### 纳入

- group content 的常态竖线与左侧缩进。
- sidebar/mobile active 页面左竖线、secondary/pulse 互斥显示。
- 导航 DOM/class 回归、deep-link、折叠状态和构建验证。

### 不纳入

- 不改变 group key、组序、组成员、注册协议、权限、路由或 sessionStorage。
- 不修改 recordView、API Manifest 或页面 secondary 注册值。
- 不把 pulse 添加到 desktop top navigation 或 user slot horizontal surface。

## 父级关系

- `parent: GOAL-001-nav-group-collapsible`；父目标完整 id 与当前 workspace 相同。
- 目标文件夹在 workspace 根平铺；层级只由 `parent` 表达。

## 台账布局

本目标使用五件套与 `01-decision/`、`02-execution/`、`03-audit/` 三个平铺 ledger 目录。

