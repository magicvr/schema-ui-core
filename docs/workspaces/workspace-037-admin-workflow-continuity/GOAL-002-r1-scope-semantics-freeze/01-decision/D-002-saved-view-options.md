---
id: D-002-saved-view-options
doc: decision
goal_id: GOAL-002-r1-scope-semantics-freeze
status: draft
created: 2026-09-17
updated: 2026-09-17
parent: GOAL-002-r1-scope-semantics-freeze
version: 0.2.0
---

# D-002 · Saved View 持久化与失效候选方案

## 当前事实

- 当前 `SchemaCrudProvider` 只在挂载页面内保存 `q/filters/sort/order/page/pageSize`，没有 Saved View 存储或列可见性状态。
- 当前 renderer 没有 Saved View 权限 target；既有权限边界是 manifest 可见性与 `SchemaCrudProvider` 的 declared target evaluation。
- 当前认证上下文可提供用户边界；本记录不把“用户级”扩展为跨设备或跨用户共享。

## 待用户裁决的方案

| 方案 | 做法 | 收益 | 代价 / 边界 |
|------|------|------|-------------|
| A（推荐首波） | 浏览器 `localStorage`，键按 `user.id + pageId + tableId` 隔离；保存 allowlist 内的 `q/filters/sort/order/pageSize` 与后续列配置 | 无 API/迁移；首波可逆、范围小；刷新后可恢复 | 仅同浏览器/同设备；清理站点数据会丢失；不支持跨设备共享 |
| B | 复用现有数据库，新增 Saved View API、表/迁移和权限校验 | 可跨浏览器/设备；服务器可审计 | 新增持久化面、迁移与 API 合同；需要独立高影响审计 |
| C | 只写 URL/history | 无新增存储 | 不满足 Saved View 个人持久化；会暴露查询值并污染历史；不作为实现方案 |

## 已由 D-003 承接的共同语义

- 首波仅保存当前表格查询与列配置 allowlist，不保存当前页、选择集、record view 或 modal 草稿。
- 无跨用户共享、最近/收藏或协作权限；无效/越权视图必须 fail closed。
- 用户已选择方案 A；完整决策见 [D-003-saved-view-localstorage-accepted.md](D-003-saved-view-localstorage-accepted.md)。本候选方案记录保留为 `draft`，I-037-002 的 R1 信息冻结已由 D-003 verified；R2 仍需提供序列化、Schema/权限失效和异常路径的实现证据。
