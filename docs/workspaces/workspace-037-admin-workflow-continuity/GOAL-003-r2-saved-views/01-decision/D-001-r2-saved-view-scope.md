---
id: D-001-r2-saved-view-scope
doc: decision
goal_id: GOAL-003-r2-saved-views
status: accepted
created: 2026-09-17
updated: 2026-09-17
parent: GOAL-003-r2-saved-views
version: 0.1.0
---

# D-001 · R2 Saved Views 分母与运行时 allowlist

## 决定

1. 首波 Saved Views 只覆盖 R1 `r1-denominator-matrix.json` 中的 24 个 `type: table` Schema 表面；`notifications`/`mail-admin-tab`/`telegram-operator` 等自定义列表表现明确排除，不因“看起来像列表”扩张分母。
2. R1 F-002～F-004 的精度修订在 R2 入口完成：custom profile 隐藏路由集合包含 `telegram-operator`；`data-permission/policies` 的 `tableFilterFields` 以活 Schema 为准；附件记录首波排除理由。
3. 运行时 allowlist 从当前表格的 Schema columns、sortable fields 与 search/filter form 绑定派生。历史 Saved View 读入时逐条校验，未知/失效字段的视图丢弃并反馈，当前查询保持不变。
4. 持久化合同沿用 R1 D-003：用户、页面、表格 ID 均做分段编码（含显式 `.` 编码），同浏览器同设备，`activeViewId` 仅作 UI 恢复指针；不保存 page/selection/record/modal draft，不新增服务端持久化。

## 取舍

不把自定义组件的内部列表状态强行接入 Schema table contract；这样首波边界可复核，未来若要覆盖自定义表面需另开决策并补独立证据。

## 关联

- R1 independent：`GOAL-002.../03-audit/A-002-r1-independent-semantics-freeze.md` F-002～F-004
- R1 用户决策：`GOAL-002.../01-decision/D-003-saved-view-localstorage-accepted.md`
- R2 验收附件：[r2-saved-view-acceptance-matrix.md](../attachments/r2-saved-view-acceptance-matrix.md)
