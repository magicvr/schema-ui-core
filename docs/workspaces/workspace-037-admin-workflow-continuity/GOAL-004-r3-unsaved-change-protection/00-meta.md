---
id: GOAL-004-r3-unsaved-change-protection
title: R3 未保存变更保护与离开确认
status: done
parent: GOAL-001-admin-workflow-continuity
created: 2026-09-17
updated: 2026-09-17
version: 0.3.0
progress: 4/4
plan_refs:
  - VP-037-admin-workflow-continuity
primary_plan: VP-037-admin-workflow-continuity
vision_ref: schema-ui-core-admin-foundation@0.4.0
---

# GOAL-004 · R3 未保存变更保护与离开确认

## 概述

在 R1 D-004 冻结语义和 R2 已建立的 dirty-state 基础切片之上，完成默认表单的 dirty 注册、页面内导航/浏览器离开保护、提交成功/失败与 reset/cancel 的状态闭环。目标是证明 dirty guard 的行为可验证，不把 R4 的统一反馈或 R5 的组合关门提前计入。

## 范围与边界

- 默认模式表单以挂载后的初始化值作为 baseline；结构性变化标记 dirty，回到 baseline 清 dirty。
- search form 的 q、筛选、排序、分页属于查询/视图状态，不阻断离开。
- App 内部导航（菜单、面包屑、Schema navigate）在 dirty 时复用统一确认；取消保持当前页面，确认后丢弃并继续导航。
- `popstate` 取消时恢复最近已提交的 URL；确认时提交浏览器目标；`beforeunload` 只遵循浏览器原生提示合同。
- 提交成功清 dirty；客户端/服务端/传输失败保留草稿和 dirty；不自动重复提交。modal close/cancel 与 reset 遵循同一保护语义。

## 成功检查点

- [x] C1：dirty registry、默认表单 baseline 和 search 非 dirty 边界有实现与单元/Renderer 证据（R3 UI 4 项 + registry 3 项）。
- [x] C2：内部导航和 `popstate` 的确认、取消、URL 恢复与确认后切换有 App 集成证据（App integration）。
- [x] C3：`beforeunload` dirty/clean 两条路径，以及 modal close/cancel 的确认行为有自动化证据（R3 UI + App integration）。
- [x] C4：提交成功/失败、reset/cancel 结果、自审 + independent audit、required finding 响应与 Git checkpoint 完成；R3 已关闭并投影 Root（A-003/A-004，checkpoint `d2b39189`）。

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|----------------|----------|--------------|-----------------|------|-------------|-------------|
| R3-I-001 | required | 默认表单的 baseline、结构比较与成功提交后的清理是否覆盖 inline/modal 生命周期？ | C1/C4 | C1 | 读取 FormInner 生命周期，补 Renderer 回归 | verified | 2026-09-17；E-002 | `r3-dirty-state.ui.test.tsx` + `render.tsx` |
| R3-I-002 | required | App 内部导航、popstate 取消/确认与 committed URL 恢复是否满足 D-004？ | C2 | C2 | App jsdom 集成测试 + 路径断言 | verified | 2026-09-17；E-002 | `App.integration.test.tsx` |
| R3-I-003 | required | beforeunload、modal close/cancel 和 reset 路径是否可观察且不误放行？ | C3/C4 | C3 | 事件测试、modal/表单回归与浏览器合同核对 | verified | 2026-09-17；E-002 | `r3-dirty-state.ui.test.tsx` + `App.integration.test.tsx` |
| R3-I-004 | non-blocking | 浏览器对 beforeunload 文案的具体呈现是否需要定制？ | R3 UX 细节 | R5 或真实触发 | 沿用 D-004，使用浏览器原生文案；真实需求走 `/vision` | deferred | 浏览器不保证自定义文案；owner=`/vision` | D-004 |

## 父目标

- `GOAL-001-admin-workflow-continuity`（Root 当前 `active · 3/5`；R1/R2/R3 已完成，本目标为 R3 已关闭阶段）。

## 台账布局

本目标从第一条记录起使用平铺 ledger：`01-decision/`、`02-execution/`、`03-audit/`，并保留 `attachments/`。C4 已由 A-003 independent recheck、A-004 self close-out 与 Git checkpoint `d2b39189` 共同关闭；R4/R5 不在本目标范围。
