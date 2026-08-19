---
id: D-004-r6-use-audit-fail-closed
goal: GOAL-007-r6-api-token-service-credential
doc: decision-entry
record_id: D-004
status: accepted
parent: null
created: 2026-08-19
updated: 2026-08-19
version: 0.1.0
---

# D-004 · 使用审计失败时 fail closed

## 决定

采纳 A-010 F-010 的 `fixed` 路径：service credential 认证成功后，使用审计事件与 `last_used_at` 更新都属于 R6 成功标准 3 的必需持久化。任一写入失败时，认证中间件返回 HTTP 503 `STORAGE_UNAVAILABLE`，不调用下游 handler。

生产组合根通过调用方事务同时提交使用审计与 `last_used_at` 更新；任一写入失败都会回滚另一项。兼容旧 Repository seam 时仍先执行 recorder 再更新元数据，并在任一错误时 fail closed，避免在存储故障期间继续放行不可追责的机器凭据调用。

## 取舍

- 未采用 D-003 的 best-effort 放行语义；该语义只保留为历史决策，不再覆盖本次 F-010 响应范围。
- 未选择 `accepted-residual` 或 `user-overruled`；用户已书面确认 `fixed`。
- 本次不改变 Profile、Manifest、协议 pin、credential scope 或管理 mutation 的事务语义。

## 验证要求

必须有回归测试证明 recorder 失败与 `last_used_at` 更新失败均返回 503 且不执行下游 handler；正常使用仍写入 audit 与 `last_used_at`。
