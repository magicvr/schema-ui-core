---
id: GOAL-007-list-page-visual-alignment
title: R6 列表页视觉与筛选体验收敛
status: done
parent: GOAL-001-admin-workflow-continuity
created: 2026-09-18
updated: 2026-09-18
version: 0.3.0
progress: 4/4
plan_refs:
  - VP-037-admin-workflow-continuity
primary_plan: VP-037-admin-workflow-continuity
vision_ref: schema-ui-core-admin-foundation@0.4.0
---

# GOAL-007 · R6 列表页视觉与筛选体验收敛

## 概述

在 Root 尚未关门、R5 用户确认门禁保持开放的前提下，新增 R6 子目标，参考 `raw/new-table` 范例页收敛现有通用列表页的视觉和布局。R6 沿用 VP-037 与当前工作区，不新建 VP 或 workspace，不改变顶部功能栏、左侧导航、后端 schema、Saved View 存储格式或现有查询/重置语义。

## 范围与边界

- 通用列表、筛选卡片/控件、页面级按钮与表格/分页区域采用范例页的布局层级和视觉层次。
- 视图切换与保存视图移到主内容标题行右上角；默认视图文案使用“全部{对象}”，对象名由集中解析逻辑提供。
- “列配置”位于页面级按钮组最左端。
- 筛选栏默认折叠，折叠后只展示第一行字段；字段值、查询、重置和已有即时筛选/查询提交逻辑保持不变。
- 单页或零结果的有效列表响应也显示分页区域；不可用的前后页/跳页操作保持禁用。
- 当前尚未实现的多选批量动作不新增、不扩张本目标范围。

明确非目标：顶部功能栏与左侧导航结构/样式、全局视觉 token 值、后端 schema、查询控制器重写、Saved View 持久化格式、逐业务页面定制和新的批量操作。

## 高层路线图

1. **C1 · 范例与现状基线**：记录范例页布局、通用列表调用链、token 映射、查询/重置/分页现状与 shell 边界。已完成，证据见 `D-001` 与 `E-001`。
2. **C2 · 列表视觉合同**：冻结共享筛选展示层、对象语义文案、标题行插槽、按钮顺序、折叠可访问性和单页分页合同。已完成，证据见 `D-001` 与实现核对。
3. **C3 · 通用实现与回归**：在不改变查询归属的前提下完成实现，补齐单元/集成回归并核对顶部栏与左侧导航未变。已完成，证据见 `E-002`。
4. **C4 · 阶段审计与交付**：完成 self 审计，响应 required finding；不借 R6 审计关闭 R5 的 Root/VP 用户确认门禁。已完成，证据见 `A-001`。

## 成功检查点

- [x] C1：范例页、受影响实现、token、查询/重置、分页与 shell 边界已形成可核对基线。
- [x] C2：列表视觉合同与信息项达到实施就绪，未把未知项伪装为已确认事实。
- [x] C3：通用列表页完成视觉/交互调整，查询/重置语义保持，单页分页显示，相关测试通过。
- [x] C4：R6 self 审计与响应完成；R5-I-004 仍按原台账等待用户对 Root/VP 的书面关门确认。

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-007-001 | required | 范例页的列表、筛选、视图、页面级按钮和分页布局是什么？ | C2 | C2 | 读取 `raw/new-table` HTML/DESIGN.md 并登记映射 | verified | 2026-09-18 已完成；若范例与实现冲突回到 C2 | `raw/new-table/schema_ui_core_2/code.html`、`raw/new-table/monochrome_technical/DESIGN.md`、D-001 |
| I-007-002 | required | 当前所有通用列表的调用链及顶部栏/左侧导航不可变边界是什么？ | C2/C3 | C2 | 盘点 `App.tsx`、`schema-table.tsx`、`data-table.tsx`、`render.tsx` 与相关测试 | verified | 2026-09-18 已完成；实现不得扩展到 shell | SCOUT 基线、D-001 |
| I-007-003 | required | 查询、重置、表级筛选与单页分页的现有行为和测试约束是什么？ | C2/C3 | C2 | 对照 handlers、query bridge、分页条件与测试 | verified | 2026-09-18 已完成；实现只改展示条件 | SCOUT 基线、D-001 |
| I-007-004 | required | 每个列表页可稳定使用的语义对象名称来源是什么？ | C3 | C3 前 | 集中解析 table/page title 与现有翻译；补充缺口测试 | verified | 2026-09-18 已由集中解析器与 3 项对象标签测试验证；未知对象回退为表/页标题或记录 | D-001、`list-surface.test.ts`、E-002 |
| I-007-005 | non-blocking | 当前分母是否已有可复用的多选批量动作？ | C3/C4 | C3 | 复用现有能力；未实装则记录忽略 | deferred | 责任人：R6；触发：用户明确要求批量操作 | 用户已明确未实装时直接忽略 |

## 父目标

- `[workspace-037-admin-workflow-continuity]` `GOAL-001-admin-workflow-continuity`。

## 台账布局

本目标从第一条记录起使用平铺 ledger：`01-decision/`、`02-execution/`、`03-audit/`，并保留 `attachments/`。

## 备注

- R5 `GOAL-006-r5-composition-acceptance` 仍为 `active · 3/4`，其 R5-I-004 用户书面确认未被本目标替代、关闭或推断。
- `progress: 4/4` 只由上述四个显式检查点派生；它不放行方案、不关闭 finding，也不改变 Root/VP 状态。
