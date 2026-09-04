---
doc_type: goal-execution
id: E-008-r2-c3-finding-remediation
parent: GOAL-003-r2-connection-settings
date: 2026-09-04
status: done
version: 0.1.0
---

# E-008 · R2 C3 A-010 必改项修复事实

## 已发生事实

- A-010 Grok independent 原文保留，记录的 F-001～F-003 `required` 不被静默降级。
- 在 `4cc96b06` 中修复 F-001：Telegram `getUpdates` payload 仍为 `timeout=30s`，本地 polling request context 改为 `35s`，独立 HTTP client 保持 `40s`，并增加 timeout ordering/deadline grace 测试。
- 在 `4cc96b06` 中修复 F-002：Bot API create/execute request errors 不再包装含 token 的认证 URL；增加 transport error 脱敏测试，并核对 manager `LastError` 路径。
- 在 `4cc96b06` 中修复 F-003：watcher 在拿到 operation lock 后复核 lifecycle context 与 `started`，`startPolling` 同时拒绝 stopped、缺 lifecycle 或已取消的上下文；增加 stopped watcher 回归测试。
- 已新增 A-011 `source: self` 修复复核；required closure 仍等待 Grok independent re-audit，不提前关闭 C3。

## 验证

- `go test ./internal/channel/telegram ./internal/composition -count=1 -timeout=120s`：通过。
- `go test -race ./internal/channel/telegram -count=1 -timeout=120s`：通过。
- `git diff --check`：通过；代码修复提交 `4cc96b06`。
