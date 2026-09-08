---
id: GOAL-003-sidebar-engine-navigation
doc: decision
status: active
parent: GOAL-001-nav-group-collapsible
created: 2026-09-08
updated: 2026-09-08
version: 0.1.0
---

# 决策记录 · GOAL-003-sidebar-engine-navigation

## 信息需求与阶段门禁

| ID | 级别 | 所需信息 / 假设 | 影响门禁 | 最晚阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 决策 |
|---|---|---|---|---|---|---|---|---|
| I-003-001 | required | Dashboard 的默认分组语义 | P1 | P1 | Dashboard Provider 显式声明 `workspace` Group；删除聚合层对 Dashboard 顶层单例的强制保留 | verified（用户决策） | — | D-001 |
| I-003-002 | required | 组/页面副字符的注册与传输 | P1 | P1 | kernel `Secondary` → Manifest 组装 → Web `secondary` projection；空值不输出/不渲染 | verified | — | D-001；E-002；E-003 |
| I-003-003 | non-blocking | 详情抽屉的通用性边界 | P2 | P2 | 仅实现布局和通用值格式化；不按 `username/email/status` 等业务字段推断 summary/操作区 | verified（用户要求） | — | D-001 |
| I-003-004 | required | 抽屉交互可访问性兼容 | P2/P3 | P3 | 保留 static/selection 分支；补遮罩、Esc、focus trap/restore、body lock 与响应式断言 | verified | — | D-001；E-002；E-003 |
| I-003-005 | non-blocking | 上游正式协议发行 | 本目标边界 | P1 | 保持 v2.7/v2.8/v2.9 pinned artifacts 不变；本地扩展不宣称上游协议字段 | deferred | 理由：本轮仅交付本仓 Admin Shell；owner：protocol maintainer；复核触发：首个对外消费者要求正式 secondary 字段时，另立兼容性目标并复审 | D-001 |

## 决策索引

| D-ID | 日期 | 标题 | 状态 | 文件 |
|---|---|---|---|---|
| D-001 | 2026-09-08 | Sidebar Engine、Workspace 默认组、副字符与通用抽屉范围冻结 | accepted | `01-decision/D-001-sidebar-engine-and-record-drawer-scope.md` |

## 当前方案边界

- Dashboard 通过 `NavigationContribution.Group` 显式注册到 `workspace`，而不是由 Shell 根据 `pageRef` 猜测；因此未来 `Group=nil` 导航仍保持既有顶层平铺向后兼容。
- `secondary` 是可选、非空、trim 后的文字表现元数据：组使用 `NavigationGroup.Secondary`，页面使用 `NavigationContribution.Secondary`。它不参与菜单身份、权限、路由、折叠状态或系统数据 checksum。
- Manifest 聚合只对已注册且已有的导航 raw item 叠加副字符；副字符通过协议已允许的 literal `label` fallback（`主文案 · secondary`）承载，Web projection 再拆分，不增加 pinned schema 字段。` · ` 作为本地 carrier 分隔符保留给该表现扩展，模块不应把它用于未注册的复合 label。不得凭 contribution 合成缺失的链接。相同注册 key 的组元数据（含 secondary）不一致时 fail closed；未注册 secondary 不输出。
- 当前产品为默认分组注册稳定 secondary：`WORKSPACE`、`IAM`、`CMS`、`OPS`、`COMMS`、`COMMERCE`；页面只示范注册稳定的 `01`（Dashboard）与 `SQL`（Data dictionary），其它页面不显示副字符。
- Sidebar Engine 采用现有 semantic `muted/accent/border/focus` token，不引入业务颜色与闪烁光点；desktop/mobile 共用 projection 与分组组件。
- recordView 仅改通用 shell：右侧 Drawer、移动底部 Sheet、blur backdrop、字段分隔卡片、可换行值、通用关闭 footer 与 modal 可访问性；不复制参考页的用户身份、安全指标或 Quick Ops。
- static `props.record` 继续无 backdrop、无 `aria-modal`、无 body lock/focus trap；selection-driven `selectedRow` 继续负责 Drawer 的打开/关闭。

## 未选方案

1. **在 Shell 中中央维护 NodeID → 副字符映射**：不采用，会破坏模块注册者拥有表现元数据、Manifest 单一投影与模块解耦。
2. **让模块作者直接把副字符手工拼进业务 label 或 page title**：不采用，会损坏 i18n 主文案与注册时的可选语义；由聚合层根据独立的 `Secondary` 注册值写入保留的 literal fallback carrier，Web projection 再拆分。
3. **复制参考页的用户 summary、安全指标、Quick Ops**：不采用，这些内容不是通用 recordView 契约，且会将业务假设写入通用控件。
4. **将 Dashboard 继续硬编码为顶层单例**：不采用，本轮用户明确要求 Dashboard 注册到必有的 Workspace 默认组。

