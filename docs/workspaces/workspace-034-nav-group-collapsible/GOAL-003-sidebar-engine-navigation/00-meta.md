---
id: GOAL-003-sidebar-engine-navigation
title: Sidebar Engine 导航与通用详情抽屉视觉优化
status: done
parent: GOAL-001-nav-group-collapsible
created: 2026-09-08
updated: 2026-09-08
version: 0.1.0
progress: 100%
plan_refs:
  - VP-034-nav-group-collapsible
primary_plan: VP-034-nav-group-collapsible
serves_summary: 在已完成的 Admin 导航分组能力上，补齐 Sidebar Engine 风格的导航层级与注册副字符语义，并将通用 recordView 详情抽屉调整为参考页的高密度侧滑显示，而不引入用户域硬编码。
---

# GOAL-003 · Sidebar Engine 导航与通用详情抽屉视觉优化

## 概述

这是 workspace-034 在 `GOAL-001-nav-group-collapsible` 与 `GOAL-002-navigation-group-polish` 之后新增的 sibling 增量目标。父目标和既有增量目标的完成事实保持不变；本目标只承接新的视觉语义与通用详情显示优化。

本目标对照：

- `raw/stitch_schema_ui_core_admin_console` 的 Sidebar Engine 范例控件：分组标题、英文副字符、叶子项层级、紧凑间距、边界/悬停/激活语义；不复制范例中的闪烁选中光点。
- `raw/schema_ui_core_2` 的详情侧栏：右侧固定抽屉、遮罩模糊、分区标题、分隔字段卡片、响应式底部 Sheet；不复制其中的用户、ROOT、MFA、Quick Ops、SESSION-REF 等业务内容。

## 成功标准

- [x] `workspace` 作为默认产品分组由 Dashboard 注册，Dashboard 位于该组；既有五个产品分组与 Examples 作者组继续保持可用。
- [x] 组与页面均支持由注册者提供可选的 `secondary` 副字符；未注册时不输出、不渲染，且不影响权限、slot、路由和 active 状态。
- [x] Desktop sidebar 与 mobile navigation drawer 采用统一的 Sidebar Engine 层级视觉；页面激活使用稳定高亮，不增加闪烁光点。
- [x] 通用 `recordView` 详情抽屉采用参考页的侧滑/遮罩/字段卡片视觉，数据、字段声明、翻译回退和静态/selection 模式边界保持通用，不硬编码用户信息。
- [x] 相关 API/Web 单测、构建与 scope-specific 浏览器回归通过，并形成可核对的执行与 self 审计证据。

## 纲领路线图

以下 3 个检查点是本目标 `progress` 的唯一来源，默认等权：

| 检查点 | 目的 | 状态 |
|---|---|---|
| P1 | 注册副字符、Workspace 默认组与 Sidebar Engine desktop/mobile 视觉实现 | completed |
| P2 | 通用 recordView 抽屉视觉、响应式与可访问性实现 | completed |
| P3 | API/Web 回归、构建、事实审视与 self 审计完成 | completed |

`progress: 100%` = 3/3 个检查点完成（P1、P2、P3）。progress 仅作展示，不放行阶段、不关闭 finding、不自动推导 `done`。

## 信息需求与门禁（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|---|---|---|---|---|---|---|---|---|
| I-003-001 | required | Dashboard 是否由注册者显式加入 `workspace` 默认组，还是由 Shell 猜测未分组页面？ | P1 方案冻结 | P1 | 对照用户指令、父目标的 Group=nil 向后兼容边界与 Dashboard Provider；冻结显式注册方案 | verified（本轮决策） | — | D-001 |
| I-003-002 | required | 组/页面副字符如何从注册层传到可见导航，且未注册时保持缺省？ | P1 实施 | 检查 kernel contribution、Manifest 聚合、Web parser/projector 的全链路回归 | verified | — | D-001；E-002；E-003；API/Web 回归 |
| I-003-003 | non-blocking | 详情抽屉是否需要复制范例页的用户 summary、安全指标与 Quick Ops？ | P2 方案冻结 | P2 | 保持 recordView 通用；只复用布局/色彩/密度，不按字段名推断业务块 | verified（用户要求 + 本轮决策） | — | D-001 |
| I-003-004 | required | 选择驱动详情抽屉的遮罩、Esc、焦点循环/恢复与 body scroll lock 是否仍保持现有通用语义？ | P2 验收 | P2/P3 | 扩展现有 Drawer/Sheet 测试；静态 `props.record` 保持非 modal | verified | — | D-001；E-002；E-003；visual-fidelity/render tests |
| I-003-005 | non-blocking | `secondary` 是否需要升级上游 `schema-ui-docs` 的正式协议字段？ | 本目标交付边界 | P1 | 本目标采用 schema-ui-core 的本地、可选表现扩展；不修改已 pin 的上游 schema/provenance，后续正式上游发布另立目标 | deferred | 理由：本轮仅交付本仓 Admin Shell；owner：protocol maintainer；复核触发：首个对外消费者要求正式 secondary 字段时，另建兼容性目标并复审 | D-001 |

## 愿景对齐

- Charter：`schema-ui-core-admin-foundation@0.4.0`
- VP：`VP-034-nav-group-collapsible`（`active` · v0.3.0）
- 工作区：`workspace-034-nav-group-collapsible`（`delivery`）
- 父目标：`GOAL-001-nav-group-collapsible`（`done 5/5`）
- 既有增量：`GOAL-002-navigation-group-polish`（`done 2/2`，不重开、不改写）

## 范围边界

### 纳入

- Dashboard Provider 显式注册 `workspace` 分组；默认产品分组的 secondary 文本与页面 secondary 文本由注册者提供。
- API kernel/Manifest 聚合与 Web manifest 消费链的可选 `secondary` 表现元数据（由 literal `label` fallback 承载，保持 pinned manifest fieldset）；呈现元数据不进入 `menu_items`、权限或 navigation checksum。
- Sidebar Engine 风格的 sidebar 与 mobile drawer：分组标题层级、英文副字符、叶子项右侧副字符、紧凑树形导轨、激活/hover/focus 状态。
- 通用 recordView 抽屉的 header/backdrop/body/字段卡片/footer、响应式 Sheet、对象/数组/长值显示与选择模式可访问性。
- 默认/optional/custom/demo profile 与既有 recordView/active/deep-link 回归。

### 不纳入

- 不重开或改写父目标、GOAL-002、VP-034 的历史 status/progress/既有分组成员事实。
- 不复制参考页中的用户姓名、邮箱、ROOT、MFA、Quick Ops、SESSION-REF 或任何用户域字段推断。
- 不增加闪烁选中光点；不改变 top/user slot 语义、权限、路由、sessionStorage 折叠契约。
- 不修改 `schema-ui-docs@v2.9.0` pinned artifacts、provenance 或对外正式协议版本；本目标的 `secondary` 仅作为 schema-ui-core 本地可选表现扩展。

## 父级关系

- `parent: GOAL-001-nav-group-collapsible`；父目标完整 id 与当前 workspace 相同。
- 目标文件夹在 workspace 根平铺；层级只由 `parent` 表达。

## 台账布局

本目标使用五件套与 `01-decision/`、`02-execution/`、`03-audit/` 三个平铺 ledger 目录。

## 备注

本目标的审计模式按当前风险暂定为 `self`：范围边界清楚、实现可回滚；若实施过程中将 `secondary` 升级为正式跨项目协议或出现兼容性门禁，再按 P-004 重新判定审计模式。

