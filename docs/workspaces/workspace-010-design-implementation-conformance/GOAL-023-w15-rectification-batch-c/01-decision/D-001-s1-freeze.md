---
id: GOAL-023-w15-rectification-batch-c
doc: decision
status: active
parent: GOAL-020-w15-user-perspective-findings
created: 2026-08-17
updated: 2026-08-17
version: 0.1.0
---

# D-001 · 批 C 冻结

- F06：changePassword onSuccess messageKey `schema.account.passwordChangedReauth` + reload（token_version 已立即作废）。
- F08：字段 reason 用 `required` / `string` 规则码。
- F09：error toast 不自动消失。
- F13：`X-Refresh-Token` 对应当前会话 `current: true`，并带本次请求 UA/IP。
- F14：MFA 取消；空表区分 listEmpty/noItemsMatch；钱包 `decimals`；删除 `nginx.conf;C`。菜单方向键留到同文件 App/notification-bell（本批先落地 MFA/空态/遗留物）。
