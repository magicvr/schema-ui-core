---
id: A-008-r3-shell-self
doc: audit-entry
parent_goal: GOAL-001-nav-group-collapsible
source: self
auditor: /govern（schema-ui-core 编排器）
type: execution-facts
scope: R3 Shell 分组折叠、键盘交互、sessionStorage 与深链自动展开
date: 2026-09-07
verdict: pass
created: 2026-09-07
updated: 2026-09-07
version: 0.1.0
---

# A-008 · R3 Shell 分组交互自审

## 范围与证据

本自审覆盖 D-005 与 E-006 的 R3 实施：

- `apps/web/src/app/navigation.ts` 的稳定 group key、group active 与页面父级深链激活；
- `apps/web/src/app/App.tsx` 的折叠 button、`aria-expanded`/`aria-controls`、Enter/Space、sessionStorage 容错；
- desktop sidebar/mobile drawer 共用组件；top/user slot 不纳入分组组件；
- `navigation.test.ts`、`nav-groups.test.tsx`、App 集成、全量 Web 回归与构建。

## 事实核对

1. 原生 button 提供可访问焦点与语义；显式 Enter/Space handler 预防默认行为重复并切换状态。
2. `schema-ui:nav-groups:v1` 只保存布尔 open/closed 状态；损坏/不可用存储回退默认展开，不阻断导航。
3. group active 来自可见 child active；当前页面为已登记内页/动态路径时，父 NavLink 与父组 active，直接 URL 会自动展开。
4. Dashboard/top/user 语义未被 R3 组件改变；UserMenu 继续只消费 link，mobile drawer 继续只合并 top/sidebar。
5. 现行协议 `NavGroup` 未增加 key/id；UI key 仅为 labelKey/literal 派生的本地稳定标识。

## 验证结果

- `navigation.test.ts`：7/7 pass。
- `nav-groups.test.tsx`：3/3 pass（Enter/Space、sessionStorage、损坏存储、inner deep link）。
- Web 全量：98 files / 1337 tests pass。
- `tsc -b` pass；`vite build` pass（仅既有 chunk size warning）。

## Findings

本 self scope 无新的 required/recommended finding。I-034-004 的完整 default/optional/custom/demo route matrix 仍按 R4 门禁 open，不阻断 R3 检查点的 Shell 代码完成；R4 必须继续验证所有分组 NodeID、动态/内页、权限过滤、slot 与 Profile 组合。

## 结论

**verdict: `pass`**。R3 Shell 行为代码与当前专项验证足以标记 R3 检查点完成；不把 R4 全量迁移/组合回归或 R5 关门写成已完成事实。

## 声明

本意见为 `source: self`，R3 为可逆 UI/本地状态变更，不触发 independent provider 门禁；不修改目标 status/progress，不替代后续 R4/R5 审计。
