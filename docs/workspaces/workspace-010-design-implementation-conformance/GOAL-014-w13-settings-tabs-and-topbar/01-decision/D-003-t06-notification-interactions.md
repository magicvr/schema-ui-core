---
id: D-003
doc: decision
status: accepted
parent: GOAL-014-w13-settings-tabs-and-topbar
created: 2026-08-16
updated: 2026-08-16
version: 0.1.0
---

# D-003 · T-06 通知中心交互修正（用户点名四项问题）

## 背景

2026-08-16 用户点名四项通知功能问题：① 顶栏铃铛下拉的通知条目不可点击——直觉是点击跳转通知列表页并展开对应详情、同时标记已读；② 通知列表页点击条目不打开详情也不标记已读；③ 列表行内 action【标为已读】多余（应由点击条目逻辑覆盖）；④ 通知页阅读后顶栏未读数要刷新页面才更新。

## 决策

### 交互模型（点击即读 + 展开详情）

- **铃铛下拉条目可点击**：点击条目 → best-effort POST /api/notifications/{id}/read（先标记）→ 关闭下拉 → 跳转 /notifications?open=<id>；通知页深链展开该条详情并标记已读。
- **列表页点击条目**：点击行 → 行内展开详情面板（标题/全文/事件/时间/已读状态）+ 未读条目自动 POST …/{id}/read；已读条目只展开不重复标记。
- **移除行内【标为已读】action**：由点击逻辑覆盖；保留工具栏【全部标为已读】与搜索/筛选表单。
- **未读数即时刷新**：POST /api/notifications/{id}/read 与 read-all 响应携带 X-Schema-UI-Config-Changed: notifications.read（复用既有 config-change 通道）；铃铛订阅该命名空间 → 徽标即时重查；下拉打开时同时重查条目。

### 实现形态

- 通知列表页由 schema 驱动的 table 节点改为 **GOAL-018 自定义组件 notification-center**（{type:custom, component:notification-center, props:{targetTable}}）：共享搜索表单的查询状态（targetTable 键），行点击/展开/标记/分页/全部已读全部内聚，避免为单个页面扩展协议 schema。
- 组件直接走 crud.fetcher（配置感知）标记已读：响应头自动触发铃铛刷新；不产生误导性 toast（通用 action 的 POST 成功文案是「已创建」，不适用于标已读）。
- **未选方案**：① 改 renderer 支持行点击动作（需改钉死协议 schema，拒绝）；② 用 recordView 抽屉 + 选中自动动作（renderer 无该钩子，且深链展开无法表达）；③ 后端新增「创建通知」管理端点仅为 e2e 造数（扩大攻击面，拒绝）。

## 影响

- **go 判定**：仅两条既有端点响应增加响应头（无新路由/无数据变更）→ Profile 默认集 / 模块矩阵 / Manifest 装配语义不变 → **VP-008 go 无影响、不暂挂**。
- **测试影响**：Go 头断言；web 铃铛可点击/徽标即时刷新单测 + notification-center 单测（行点击/已读不重复/深链）；e2e 通知页冒烟（空态渲染）；D-VAL/schema-keys 全绿。
