---
doc_type: goal-execution
id: E-002-webhook-dispatcher-implementation
parent: GOAL-003-r2-webhook-dispatch-identity
date: 2026-09-03
status: recorded
---

# E-002 · Telegram Webhook 路由、Update 分发、主体映射与入站限流实施（C2）

## 1. 实施范围

依据 [D-001-r2-architecture-and-subject-store.md](../01-decision/D-001-r2-architecture-and-subject-store.md) 与 R1 合同 [D-002](../../GOAL-002-r1-contract-freeze/01-decision/D-002-telegram-channel-contract.md)，落地 Telegram 通道 R2 运行时：
- `apps/api/internal/channel/telegram/types.go`：Update JSON 数据结构。
- `apps/api/internal/channel/telegram/dispatcher.go`：Dispatcher 调度引擎与未注册命令回落回复。
- `apps/api/internal/channel/telegram/capture_sender.go`：内存 Capture / Mock Sender。
- `apps/api/internal/channel/telegram/webhook.go`：Webhook HTTP 处理管道（503/429/401/400/200，常时 secret 比对，IP/Chat/User 三桶限流，`issuer=telegram` 主体映射）。
- `apps/api/internal/channel/telegram/dispatcher_test.go`：Dispatcher 单元测试。
- `apps/api/internal/channel/telegram/webhook_test.go`：Webhook 全链路测试。
- `apps/api/modules/channel/telegram/provider.go`：`kernel.Provider` 模块提供者（ID = `channel.telegram`）。
- `apps/api/modules/channel/telegram/provider_test.go`：模块提供者集成测试。

## 2. 实施与验证结果

1. **测试覆盖**：
   - `internal/channel/telegram` 11 项测试通过：
     - 未配置 Token 返回 503
     - 缺失/错误/未配置 Secret 均 fail-closed 返回 401
     - 畸形 JSON 返回 400
     - 正常命令分发 + 主体自动映射（`subject.Store.GetOrCreateSubject`）
     - 未注册命令发送 `DefaultTelegramUnknownCommandText` 默认回落
     - Callback 分发与未注册 Callback 静默 200
     - IP 限流（60/m）、Chat 限流（30/m）、User 限流（20/m）429 + Retry-After 验证
     - IP 限流在 Secret 失败时仍然正确记账（防洪水攻击）
   - `modules/channel/telegram` 模块注册测试通过。
2. **回归验证**：
   - `go test ./...` 全量 PASS。
