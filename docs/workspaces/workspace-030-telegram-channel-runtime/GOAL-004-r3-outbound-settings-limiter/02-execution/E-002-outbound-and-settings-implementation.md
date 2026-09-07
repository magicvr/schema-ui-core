---
doc_type: goal-execution
id: E-002-outbound-and-settings-implementation
parent: GOAL-004-r3-outbound-settings-limiter
date: 2026-09-03
status: recorded
---

# E-002 · 出站生产适配器、动态设置与限流核账实施（C2）

## 1. 实施范围

依据 [D-001-r3-outbound-and-settings-architecture.md](../01-decision/D-001-r3-outbound-and-settings-architecture.md) 与 R1 合同 [D-002](../../GOAL-002-r1-contract-freeze/01-decision/D-002-telegram-channel-contract.md)，落地 Telegram 通道 R3 出站与管理面组件：
1. `apps/api/internal/channel/telegram/runtime.go`：`RuntimeManager` 并发安全动态配置管理器，支持 Bot Token 与 Webhook Secret 热切换，提供脱敏状态视图 `RuntimeStatus`。
2. `apps/api/internal/channel/telegram/http_sender.go`：`HTTPSender` 生产适配器，基于标准库 `net/http` 发起 POST `https://api.telegram.org/bot<token>/sendMessage`，具备 10s 超时控制、`msg.Validate()` fail-closed，未配置 Token 时自动安全降级至 `CaptureSender`。
3. `apps/api/internal/channel/telegram/settings_handler.go`：`SettingsHandler` 提供 `GET /api/channel/telegram/settings`（脱敏状态只读）与 `PATCH /api/channel/telegram/settings`（热更新配置）。
4. `apps/api/modules/channel/telegram/provider.go`：向模块贡献表注册设置路由 `GET/PATCH /api/channel/telegram/settings` 及 Webhook 路由。
5. `apps/api/kernel/profile.go`：同步候选集 `BuiltinModules()` 中的路由声明。
6. `apps/api/internal/composition/composition.go`：装配 `RuntimeManager`、`HTTPSender` 与 `SettingsHandler`。
7. 测试落地：
   - `runtime_test.go`：测试热切换、并发安全、掩码脱敏与设置端点 GET/PATCH。
   - `http_sender_test.go`：测试合法消息与 InlineKeyboard 序列化、未配置 Token 降级 Mock、校验失败拦截、API 错误处理及 10s 超时预算。
   - `composition_telegram_test.go`：测试全套三路由端到端贡献与注册。

## 2. 限流核账结果

- 入站三桶限流已在 R2（`webhook.go`）落地，并经 `webhook_test.go` 严格验证：
  - IP 桶：`tg:webhook:{ip}`（60/min），超限 429 + `Retry-After`，Secret 错误仍记账防洪水。
  - Chat 桶：`tg:chat:{chat_id}`（30/min），超限 429 + `Retry-After`。
  - User 桶：`tg:user:{user_id}`（20/min），超限 429 + `Retry-After`。
- 出站端口为外部 API 客户端，依赖自身超时与重试控制，无内部计数桶积压或泄露风险，核账 PASS。

## 3. 回归测试验证

- `go test ./internal/channel/telegram/... ./modules/channel/telegram/... ./internal/composition/... ./kernel/...` 全部 PASS。
- 全仓 `go test ./...` 全绿。
