---
id: D-001-list-page-visual-contract
doc: decision-entry
status: accepted
goal_id: GOAL-007-list-page-visual-alignment
created: 2026-09-18
updated: 2026-09-18
parent: GOAL-001-admin-workflow-continuity
version: 1.0.0
---

# D-001 · R6 列表页视觉与筛选实施合同

## 决定

采用“共享筛选展示层 + 保留现有查询归属 + 主内容标题局部插槽”的实现方向：

- `render.tsx` 的 search form 继续持有草稿并调用现有查询/重置 handlers；`schema-table.tsx` 继续持有表级筛选、排序、分页和 Saved View 状态。
- 新增纯展示职责的筛选面板，默认折叠、按响应式网格只展示第一行；不截断字段值、不发送请求、不改变校验/重置/查询归属。隐藏但已提交的有效条件通过可访问提示表达。
- Saved View 视图选择与保存/更新/删除控件通过 page-scoped 标题插槽放在主内容右上角；无插槽时保留表级回退位置。页面标题插槽是局部主内容能力，不触碰 `data-shell-region="topbar"`、`sidenav` 或导航抽屉。
- 默认视图文案使用集中对象名称解析与参数化“全部{对象}”文案；不得从 URL、接口路径或标题字符串做不可靠猜测。清空 Saved View 选择仍保持现有“不清除筛选”的行为。
- “列配置”作为页面级工具栏第一个 DOM 子节点；保留原有权限、至少一列约束和业务按钮行为。
- 只有有效列表响应时显示分页区域；总页数统一至少为 1，单页/零结果时显示第 1 页和禁用的前后/跳页操作；初始加载/失败且无有效列表时不伪造分页。
- 视觉调整只作用于列表区域，复用现有语义 token；未实装的多选批量功能不新增。

## 理由

该方向能直接覆盖用户列出的六项视觉/布局要求，同时保持 R1～R4 已冻结的查询、Saved View 和反馈边界。范例页是静态 HTML，无法直接复用组件；抽出展示层可以让所有 schema table 页面共享折叠和样式，避免逐页复制。标题插槽比绝对定位可靠，能在窄屏、多表页和路由切换时保持结构清晰。

## 关键证据

- 范例布局：`raw/new-table/schema_ui_core_2/code.html`、`raw/new-table/monochrome_technical/DESIGN.md`。
- 通用调用链：`apps/web/src/app/App.tsx` → `apps/web/src/renderer/schema-table.tsx` → `apps/web/src/components/data-table.tsx`；search form 桥接在 `apps/web/src/renderer/render.tsx`。
- shell 边界：`apps/web/src/app/App.tsx` 的 `data-shell-region="topbar"`、`data-shell-region="sidenav"`。
- 单页分页旧条件与测试：`apps/web/src/renderer/schema-table.tsx`、`apps/web/src/renderer/schema-table.test.tsx`。

## 未选方案

- 不逐页定制：会复制样式与对象名称逻辑，后续 schema table 容易遗漏。
- 不统一重写查询控制器：本需求是展示与布局变化，重写会扩大范围并触碰已验证的查询/重置语义。
- 不用绝对定位把视图操作“抬到”页面顶部：容易与标题、多表页和移动布局重叠。
- 不新增范例中的多选批量动作：用户已明确未实装时直接忽略。

## 复查触发

若发现同一表必须把两类筛选合并为一张卡片、多个可见列表争抢唯一页级视图区、自定义表单组件不支持折叠参数，或对象名无法可靠本地化，则暂停受影响实现并重新评估本决策。
