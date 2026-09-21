---
id: E-002-r2-boundary-corrections
doc: execution
goal_id: GOAL-003-r2-saved-views
status: recorded
created: 2026-09-17
updated: 2026-09-17
parent: GOAL-003-r2-saved-views
version: 0.1.0
---

# E-002 · R2 分母与边界纠偏

2026-09-17，响应 R1 independent A-002 的 F-002～F-004：

- 修订 R1 分母矩阵至 v0.2.0：custom `hiddenRoutePageIds` 补入 `telegram-operator`，derived hidden count 更新为 5；admin 保持 4。
- 将 `data-permission/policies` 的 `tableFilterFields` 更正为空数组，因为 `resource` 是 `rowKey` 而不是 Schema filter。
- 在 R2 D-001 与验收矩阵中明确 Saved View 首波只覆盖 24 个 `type: table` Schema 表面；`notifications`、`mail-admin-tab`、`telegram-operator` 等自定义列表表面排除。

矩阵经 JSON 解析核对，分母仍为 24 个列表表面、58 个表单节点；R2-I-001 可标记为 `verified`。
