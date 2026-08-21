---
id: GOAL-022-w15-rectification-batch-b
doc: decision
status: active
parent: GOAL-020-w15-user-perspective-findings
created: 2026-08-17
updated: 2026-08-17
version: 0.1.0
---

# D-001 · 批 B 冻结

- **F03**：`formatRFC3339Milli`；scheduled-tasks / dictionary 时间字段改为 RFC3339 串；filelibrary `created` **不改名**，值改为 RFC3339。
- **F11**：GET by-owner / me / me/entries 只读，缺失 404；创建走 POST。POST adjust 仍可 GetOrCreate。go-impact：改变 workspace-011 自动开户读路径。
- **F10**：登录 429 加 `Retry-After`；配额改 413。
- **F12**：`DefaultPageSize = 20`。
