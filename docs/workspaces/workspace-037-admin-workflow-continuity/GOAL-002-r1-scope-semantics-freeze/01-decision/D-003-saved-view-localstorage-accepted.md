---
id: D-003-saved-view-localstorage-accepted
doc: decision
goal_id: GOAL-002-r1-scope-semantics-freeze
status: accepted
created: 2026-09-17
updated: 2026-09-17
parent: GOAL-002-r1-scope-semantics-freeze
version: 0.1.0
---

# D-003 · 用户确认 Saved View 使用 localStorage

## 决策

用户于 2026-09-17 明确选择 D-002 的方案 A：首波 Saved View 使用浏览器 `localStorage`，按 `user.id + pageId + tableId` 建立用户与页面表格边界。

## 冻结边界

- 保存范围是当前表格查询与列配置 allowlist：`q`、`filters`、`sort`、`order`、`pageSize` 与列可见性；同一存储记录可额外保存当前选中的视图 ID，作为 UI 恢复指针，不是业务查询字段。
- 不保存当前页、选择集、record view 或 modal 草稿；恢复从第一页开始。
- 只保证同一浏览器、同一设备上的刷新后恢复；不提供跨设备同步、跨用户共享、最近/收藏或协作权限。
- 存储键的各段按稳定编码处理（包括显式编码段内 `.`）；读取到 malformed、未知字段、超出 allowlist 或当前权限/Schema 边界的视图时 fail closed，不污染当前页面状态。
- 浏览器存储不可用、读写异常或容量不足时，必须通过统一反馈呈现错误；不得伪装成已持久化。
- 本决策不新增 API、数据库表、迁移或服务端持久化；序列化、校验与失效实现证据仍归 R2 留存。R1 的信息冻结依据是本记录、矩阵和用户取舍；R2 未完成前不得宣称 Saved View 退出判据已满足。

## 未选方案

- B（数据库 + API）：跨设备能力更强，但会扩大持久化、迁移和权限合同，不属于本首波。
- C（只写 URL/history）：不满足个人持久化，且会把查询值带入历史；不采用。

## 关联

- 选项登记：`D-002-saved-view-options.md`
- 信息项：I-037-002
- 后续实现：R2 Saved Views
- UI 恢复指针：同一存储记录中的 `activeViewId` 只用于刷新后的视图选择恢复，不进入资源请求 allowlist。
