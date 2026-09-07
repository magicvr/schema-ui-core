---
id: GOAL-001-nav-group-collapsible
title: Admin 导航分组折叠体验交付
status: active
parent: null
created: 2026-09-07
updated: 2026-09-07
version: 0.1.0
progress: 0%
plan_refs:
  - VP-034-nav-group-collapsible
primary_plan: VP-034-nav-group-collapsible
serves_summary: 在现行 Admin Shell 中交付与模块解耦的导航分组、折叠/展开和直接 URL 激活态展开，并将当前已注册 sidebar 导航按合理语义迁移到分组。
---

# GOAL-001 · Admin 导航分组折叠体验交付

## 概述

在 VP-034 的范围内改造 Admin 左侧导航及其模块注册/聚合逻辑：支持多个模块将导航条目注册到同一具名组；组可折叠/展开；用户直接通过 URL 进入已分组页面时，对应组自动展开；当前已注册 sidebar 导航全部纳入合理分组迁移与验证。

本 Root 是本工作区唯一的总目标，`parent: null`。它是有界交付目标，不是 VP-010 长期符合性程序的子目标。

## 愿景对齐

- Charter：`schema-ui-core-admin-foundation@0.4.0`
- VP：`VP-034-nav-group-collapsible`（`active` · v0.3.0）
- 工作区：`workspace-034-nav-group-collapsible`（`delivery`）
- 依据：[VRev-083](../../../vision/reviews/VRev-083-vp034-nav-group-collapsible-activation.md) activation `pass`、[VRev-084](../../../vision/reviews/VRev-084-vp034-existing-navigation-scope-correction.md) scope correction `pass`

## 范围

### 纳入

- Shell sidebar 的 group 渲染、折叠/展开与键盘可访问性
- 模块 NavigationContribution / Manifest 的可选 `group` 注册与跨模块共组聚合
- 当前已注册 sidebar 导航的合理分组迁移：
  - `identity-access`：Users、Roles、Data permission
  - `content-data`：File library、Data dictionary
  - `operations`：Activity、System monitoring、Scheduled tasks、Recycle bin
  - `communications`：Mail console、Outbound email log、Telegram channel
  - `commerce`：Wallet、Prepaid vouchers、Digital products、Digital entitlements、Digital orders
  - Dashboard 保留为顶层主入口单例
- `dev.examples` 既有 Examples 组的兼容回归
- 直接 URL 进入时所属组自动展开
- top/user slot 的兼容回归与模块贡献 playbook 更新

### 不纳入

- 强制所有未来模块必须使用分组
- 将 Settings、Account、My wallet 或通知面等 top/user slot 全部搬到 sidebar
- 分组服务端持久化、分组独立权限、多级嵌套、拖拽排序
- 修改 Charter、VP-010 范围或 VP-008 `go` 语义

## 纲领路线图

以下 5 个检查点是 progress 的唯一来源，默认等权；未完成项不能写成已完成事实：

| 检查点 | 目的 | 状态 |
|---------|------|------|
| R1 | 现有导航清单、Profile/slot 矩阵、分组信息架构和 group key/顺序冻结 | completed |
| R2 | 模块注册与聚合契约：可选 group、跨模块共组、无 group 向后兼容 | completed |
| R3 | Shell 交互：折叠/展开、键盘可访问、直接 URL 自动展开、状态保持 | completed |
| R4 | 当前已注册 sidebar 全量迁移；默认、optional、custom/demo 组合与 top/user slot 回归 | completed |
| R5 | 证据矩阵、全量回归、Goal 审计、required finding 闭合与关门准备 | pending |

`progress: 80%` = 4/5 个检查点完成（R1、R2、R3、R4）。progress 只作展示，不放行阶段、不关闭 finding、不推导 `done`。

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|------------------|------|-------------|-------------|
| I-034-001 | required | 当前已注册导航的完整清单、slot、默认/optional Profile 覆盖 | R1 分母 / R4 回归 | R1 | 对照 `apps/api/kernel/profile.go`、模块 provider 与 manifest fragment，冻结 NodeID/PageID/slot/profile 矩阵 | verified（静态 R1 分母） | R4 运行时 Manifest/Profile harness 核验 | [R1 导航 / Profile / slot 盘点](attachments/r1-navigation-profile-slot-matrix.md)；E-002；当前 HEAD `f2044cf3` |
| I-034-002 | required | 五个分组的最终标题、组内顺序及 Dashboard 单例例外 | R1 方案冻结 | R1 | 用户确认分组基线；实现后的 Manifest/Shell 行为另行验证 | verified（决策） | R1 行为证据继续核对 | [D-002](01-decision/D-002-r1-baseline-and-group-contract.md)；VP-034 初始基线 |
| I-034-003 | non-blocking | 折叠状态采用会话内状态还是浏览器持久化 | R3 交互实现 | R3 | 用户确认 sessionStorage；实现阶段补交互测试；不引入服务端存储 | verified（用户决策 + 测试） | R3 以全量回归继续核对 | D-005；E-006；nav-groups.test.tsx |
| I-034-004 | required | 已分组 NodeID 的直接 URL、动态路径和激活态展开矩阵 | 判据 5 / R4 验收 | R3 | 建立 NodeID → route projection 矩阵，覆盖 sidebar 直接 route、已登记内页/动态路径与折叠后直接进入 | verified（R3/R4 矩阵） | R5 仅做关门前再核对 | D-005；E-008；`apps/web/src/app/nav-groups-r4.test.ts` |
| I-034-005 | required | optional compiled modules 在 custom/demo profile 中的跨模块分组聚合是否成立 | 判据 3/4/6 / R4 回归 | R4 | custom/demo manifest harness + 权限/无丢失回归 | verified（R4 runtime matrix） | R5 仅做关门前再核对 | E-008；`apps/api/internal/composition/nav_group_r4_test.go` |

## 父目标

- Root：`parent: null`（本目标是本工作区 Root）

## 台账布局

本目标使用平铺五件套与三个 ledger 目录：`01-decision/`、`02-execution/`、`03-audit/`。