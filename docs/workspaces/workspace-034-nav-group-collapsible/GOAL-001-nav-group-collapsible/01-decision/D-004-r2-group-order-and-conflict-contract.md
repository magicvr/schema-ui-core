---
id: D-004-r2-group-order-and-conflict-contract
doc: decision-entry
parent_goal: GOAL-001-nav-group-collapsible
source: user-confirmed + /govern response to A-002
status: accepted
date: 2026-09-07
created: 2026-09-07
updated: 2026-09-07
version: 0.1.0
---

# D-004 · R2 组序、混排与冲突契约

## 决策

用户确认以下规则，用于闭合 A-002 的 F-002/F-003/F-004，并作为 R2 实施契约：

### 1. Sidebar 输出顺序

1. 先按现有 NodeID 排序规则确定每个普通导航叶子的组内相对顺序。
2. `menu_dashboard` 始终作为 sidebar 第一项的顶层单例；它不进入任何结构化组。
3. 五个结构化 Admin 组按显式 `GroupOrder` 升序输出：
   - `identity-access` = 10
   - `content-data` = 20
   - `operations` = 30
   - `communications` = 40
   - `commerce` = 50
4. 未声明 group 的普通 sidebar 链接保持自身相对顺序，并排在结构化组之后；不得因缺 group 而丢失。
5. 已有协议组（如 demo 的 `Examples`）保持原组对象，不与五个结构化 Admin 组合并；输出在结构化组之后。
6. top/user slot 完全不参与结构化分组归一化。普通 sidebar 节点带有 group 元数据但无法在 sidebar 中匹配时，聚合 fail closed。

### 2. 同 key 元数据冲突

同一 `group key` 的跨模块声明必须对以下完整元组精确全等：

`(key, order, label, labelKey, icon)`

- literal `label` 与 `labelKey` 不互相替代；一侧缺失而另一侧存在视为冲突。
- 空 icon 只与空 icon 相等；不得由 first-writer 静默覆盖。
- 冲突在 kernel `ContributionSet.finalize` 阶段 fail closed，使用稳定错误码 `CodeModuleNavigationGroupConflict`；不得等到 Manifest 聚合后才发现。

### 3. 分层边界与协议兼容

- Group 是公开 Manifest 的展示聚合语义，不是 `Parent` 层级；不得用 Parent 表示产品组。
- Group 不进入 `menu_items` 身份或 system-data checksum，不触发 DB migration；权限/角色授权仍以 NodeID 为准。
- 现行协议 `NavGroup` 不新增 `key`/`id` 字段，继续兼容 2.7/2.8/2.9 strict schema；R3 使用 `labelKey` 与 child active 作为 UI 稳定依据。
- 空组不输出；只由已启用且实际有成员的贡献形成组。

### 4. R2 必须具备的契约测试

- 同 key 一致元数据通过；任一字段不一致 fail closed 并产生 `CodeModuleNavigationGroupConflict`。
- 不同模块向同组注册后只产生一个协议组，成员完整且组内顺序稳定。
- 未声明 group 的普通 sidebar 链接仍出现并保持相对顺序。
- Dashboard 顶层单例第一项；top/user slot 不被移动或归组。
- `dev.examples` 的 Examples 组保留，不与 Admin 五组合并。
- optional `channel.telegram` / `biz.digital-offer` 启用时分别加入 communications/commerce；未启用时不产生空组。
- 组序不随 `DefaultNavigationOrder` 首成员位次改变。

## 依据

- 用户本轮对 `r2-mixed-order` 与 `r2-metadata-conflict` 的明确选择。
- A-002 independent（grok-4.6 · high）F-002/F-003/F-004 与 recommended F-006/F-008/F-009/F-010。
- D-002 用户确认的五组基线与契约优先方案。
