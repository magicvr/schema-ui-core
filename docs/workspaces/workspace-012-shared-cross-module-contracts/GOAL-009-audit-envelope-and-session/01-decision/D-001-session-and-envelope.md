---
id: D-001-session-and-envelope
doc: decision-entry
status: accepted
parent: GOAL-009-audit-envelope-and-session
created: 2026-08-19
updated: 2026-08-19
version: 0.1.0
---

# D-001 · session 语义与 envelope 范围

用户要求两项一起做：D-003 外写路径改结构化 envelope，以及 session 关联。

## 决定

1. **session**：用户登录会话 = refresh token id，写入 access JWT `sid`，中间件灌到 `account.User.SessionID`，operation_log 侧表 `operation_log_session`（同 correlation 模式）。机器凭据 session = credential id。旧 token 无 sid 时 session 可空。
2. **effective actor** = 当前 `actor_id`。无 impersonation，不另开列。
3. **envelope**：所有生产 mutation 写路径用 `NewDetail`。无 body 的事件仍可无 detail。
4. 不做归档查询 UI、不做 impersonation、不改 Profile/协议 pin。
