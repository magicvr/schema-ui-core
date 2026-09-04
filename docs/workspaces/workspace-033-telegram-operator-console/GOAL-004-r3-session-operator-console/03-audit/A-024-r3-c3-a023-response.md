---
doc_type: goal-audit
id: A-024-r3-c3-a023-response
parent: GOAL-004-r3-session-operator-console
date: 2026-09-05
source: self
auditor: Codex
audit_type: finding-response
scope: 响应 A-023 F-001/F-002 recommended；补测试钉与 HTTPSender/operator fail-closed 接缝
verdict: pass
open_required: 0
version: 0.1.0
---

# A-024 · R3 C3 A-023 推荐项响应（2026-09-05）

## 响应结论

A-023 为 Grok `grok-4.6 · reasoning high` independent `pass`，`open_required: 0`。
本条响应其两项 low/recommended finding；原始 A-023 与 A-018/A-020 原文均保留，
不把推荐项升级为 required，也不接受 residual 或 overrule。修复提交为
`fa0caa70 fix(telegram): close C3 recommended audit gaps`。

## F-001 · 验证分母测试钉

状态：**fixed**。

- `composition_telegram_test.go` 通过真实 composition mux 验证 operator 匿名请求
  为 `401 UNAUTHENTICATED`，锁定 `Public:false` 不替代 Middleware。
- `telegram_operator_test.go` 通过真实 `auth.Authenticator.Middleware` 验证缺少
  `telegram.operator.*` 的 service credential 为 `403 FORBIDDEN`；同文件补齐
  `INVALID_PAGE`、`INVALID_PAGE_SIZE`、webhook/unknown receiver/invalid bot id、
  空 token、未知 retry request 的失败即红断言。
- `postgres_telegram_test.go` 增加同 request 与同 retry root 的 8 路并发竞争，
  各自只允许一个 pending 创建者；现有 SQLite 并发覆盖继续保留。

## F-002 · 空 token 与 Send nil 接缝

状态：**fixed**。

- `http_sender.go` 增加 `ErrTelegramTokenMissing`：无 runtime 或无显式 fallback
  时不再以 `nil` 表示成功；既有显式 `CaptureSender` 降级行为保留，避免扩大
  非 operator 调用方的行为变更。
- `telegram_operator.go` 在 sender 返回 `nil` 后再次确认 runtime/token。若状态
  在发送窗口内消失，不调用 `MarkSent`，保持 durable row 为 `pending` 并返回
  `TELEGRAM_OPERATOR_UNAVAILABLE`；pending 仍阻止后续重复外发。
- `http_sender_test.go` 锁定无 runtime fail-closed；
  `telegram_operator_test.go` 锁定 token 在 sender 窗口消失时 pending 与无重复
  sender。空 token 进入发送前的 handler 失败路径也有 `409` 与零 sender 调用
  断言。

## 验证事实

- 通过：HTTPSender、Telegram operator、composition 的新增/既有专项测试。
- 通过：上述三组隔离 `-race` 测试。
- 通过：相关包回归 `go test ./internal/store ./internal/handler ./internal/composition ./internal/channel/telegram ./modules/channel/telegram/... ./kernel -count=1`。
- 通过：`TestPostgresTelegramIngressRepositoryIdempotency`、
  `TestPostgresTelegramOutboundConflictAndRetryState`、
  `TestPostgresTelegramOutboundConcurrentRequestAndRootIdempotency`；本机 gated
  PostgreSQL 路径实际执行并通过，不以 skip 作为证据。
- `git diff --check` 与本批 owned 文件 trailing-whitespace 检查通过；A-023
  re-audit 尚待执行。

## 后续门禁

本条只响应 recommended findings，不修改 C3/R3 status 或 progress。由于本条涉及
实现修复，下一步再次调用 Grok independent 对 `fa0caa70` 当前 HEAD 做复审；复审
通过并由 `/govern` 响应后，才可关闭 C3 并放行 C4。
