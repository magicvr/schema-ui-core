---
id: E-003
goal: GOAL-013-nav-order-config
date: 2026-08-14
status: recorded
parent: GOAL-013-nav-order-config
created: 2026-08-14
updated: 2026-08-14
version: 1.0.0
---

# E-003 · S2 实现完成（导航排序）

## 事实

- 2026-08-14：S2 实现完成（提交 b839595）。

## 实施中发现的关键事实（P-005 回流）

- **UI 侧边栏顺序的真实载体是 manifest 聚合**（`manifest.ForModulesWithFragments`，按 ModuleID 字母序合并各模块 fragment 的 navigation 槽），**不是** kernel 的 `sortNavigation`（后者只驱动 system-data reconcile 的 menu_items 表顺序）。
- 因此排序需要在**两层**应用：kernel（系统数据顺序）+ manifest（发布文档顺序，UI 实际渲染）。D-001 原描述只覆盖 kernel 层；S2 补充 manifest 层。

## 实现内容

- kernel：`DefaultNavigationOrder` 常量（12 项，用户确认清单）；`Plan.NavigationOrder` 字段；`sortNavigation(nodes, order)` 三层排序（清单优先 → 未列 NodeID 追加末尾 → Parent 分组仍优先）；`resolveNavigationOrder` 未知 NodeID 整体回退默认 + 告警。
- manifest：`SortNavigation(data, order)` 对 top/sidebar/user 槽按 NodeID 重排（从 visibleWhen `features.<nodeID>` 提取，id/pageRef 兜底；未列项稳定追加末尾）；`ForModulesWithFragments` 变参 order。
- config：`Config.NavigationOrder` + YAML `navigation.order` 段（yaml.Node 宽松解析：非列表/非字符串 → 回退默认 + 告警）+ `NAVIGATION_ORDER` env 逗号分隔覆盖。
- composition：`ResolvePlan` 把 cfg.NavigationOrder 挂到 plan；manifest 构建传 navOrder（plan 覆盖或 kernel 默认）。
- 测试：kernel 快照 + 排序语义（默认/自定义/非法回退/部分覆盖）；config navigation 解析 5 项；manifest SortNavigation 排序。

## 实测（admin profile 真实 manifest）

- 排序前 sidebar：Activity | Dashboard | Data dictionary | File library | Recycle bin | Roles | Scheduled tasks | System monitoring | Users（字母序）。
- 排序后 sidebar：Dashboard | Users | Roles | Activity | File library | Data dictionary | System monitoring | Scheduled tasks | Recycle bin；user 槽：Settings | Account —— 与默认清单一致。

## 遗留

- S3 验证（含覆盖路径实测、全量回归、web/e2e）与 S4/S5。
