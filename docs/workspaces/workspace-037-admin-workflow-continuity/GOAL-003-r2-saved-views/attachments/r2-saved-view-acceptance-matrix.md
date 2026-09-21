---
id: r2-saved-view-acceptance-matrix
doc: attachment
status: active
created: 2026-09-17
updated: 2026-09-17
parent: GOAL-003-r2-saved-views
version: 0.4.0
---

# R2 Saved View 验收矩阵

| 检查项 | 目标 | 当前状态 | 证据 |
|--------|------|----------|------|
| 分母 | 24 个 `type: table` Schema 表面；custom 列表型表面排除 | verified | R1 matrix v0.2.0 + D-001 + E-002 |
| 纠偏 | custom 隐藏 `telegram-operator`、data-permission filter、排除理由 | verified | E-002；`notifications`/`mail-admin-tab`/`telegram-operator` 明确不进入首波 |
| 键隔离 | 编码后的 `user.id + pageId + tableId` | verified | `saved-views.test.ts` |
| 状态 allowlist | q/filters/sort/order/pageSize/columns；无 page/selection/modal | verified | `saved-views.test.ts` |
| 保存/选择/恢复 | 新建、active pointer、刷新/重挂载恢复第一页 | verified | `saved-views.ui.test.tsx`（5 UI tests） |
| 更新/删除 | 保留 id/createdAt；删除确认与 active pointer 清理 | verified | `saved-views.test.ts` + `saved-views.ui.test.tsx` |
| 失效/异常 | malformed、duplicate、Schema/权限失效、storage read/write failure fail closed | verified | `saved-views.test.ts` + UI feedback |
| 安全反馈 | success/error 可见且不回显原始 payload | verified（R2 scope） | `saved-views.ui.test.tsx`；跨页面统一反馈仍由 R4 复核 |
