---
id: D-001-sidebar-engine-and-record-drawer-scope
goal_id: GOAL-003-sidebar-engine-navigation
doc: decision-entry
source: user-request
status: accepted
date: 2026-09-08
parent: GOAL-001-nav-group-collapsible
version: 0.1.0
---

# D-001 · Sidebar Engine、Workspace 默认组、副字符与通用抽屉范围冻结

## 决策

1. Dashboard 由 `admin.dashboard` Provider 显式注册到 `workspace` 分组；`workspace` 使用稳定的中英文主/副显示语义并排在其它产品分组之前。聚合层不再把 Dashboard 强制保留为顶层单例。
2. 组与页面副字符统一作为可选的 `secondary` 表现元数据：组由 `NavigationGroup.Secondary` 注册，页面由 `NavigationContribution.Secondary` 注册；未注册时完全不输出/不渲染。
3. 为保持 `schema-ui-docs@v2.9.0` 的 strict manifest fieldset 不变，聚合层把已注册 secondary 叠加到协议允许的 literal `label` fallback（`主文案 · secondary`），Web projection 再拆分为独立显示值；本目标不宣称完成正式上游协议发行。Manifest 只对已存在的注册导航 raw item 叠加元数据。
4. Sidebar Engine 采用现有 semantic token，分组标题使用紧凑层级与可选英文 secondary，页面使用稳定 active/hover/focus 高亮、树形边界与可选右侧 secondary；不实现参考页中的闪烁选中光点。
5. recordView 只复用参考页的通用 layout：固定右侧 Drawer、遮罩模糊、字段 section/card/divide、响应式底部 Sheet、长值可换行与通用关闭操作；不添加用户、邮箱、ROOT、MFA、安全指标、Quick Ops 或 SESSION-REF 等业务硬编码。

## 理由

- Dashboard 的分组归属必须由注册者声明，才能与模块解耦契约、权限/slot 边界和未来未分组导航兼容。
- `secondary` 与主 label 分离，避免把英文缩写/编号混入 i18n 主文案；注册缺省可保持当前简洁导航。
- 本轮是已完成 VP-034 的产品视觉增量；保留上游协议 pin，避免把局部视觉扩展伪装成上游兼容版本升级。
- 详情控件是通用 renderer surface，参考页的用户信息仅为样例数据，不应成为 renderer 的字段名推断或固定操作。

## 验收指向

- API：secondary 的 trim/冲突/注册投影与 `menu_items` checksum 不变；Workspace 组、Dashboard 入组、未注册缺省与 profile 回归。
- Web：manifest/parser/projector 传递 secondary；desktop/mobile 两处渲染同一语义；无 secondary 时无空占位；active 状态无闪烁点。
- Renderer：generic record fixture 只显示声明字段/通用对象值；selection Drawer 与 static panel 的 modal、遮罩、焦点、滚动锁和响应式 classes 可核对。

