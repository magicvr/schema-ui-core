# VP-030 / Workspace-030 退出判据全量证据矩阵（R4）

> 对应计划：[VP-030-telegram-channel-runtime](../../../../vision/plans/VP-030-telegram-channel-runtime.md) v0.2.0  
> 责任工作区：`workspace-030-telegram-channel-runtime` · Root [GOAL-001](../../GOAL-001-telegram-channel-runtime/00-meta.md) · 子目标 [GOAL-005](../00-meta.md)  
> 编制日期：2026-09-03  

---

## 1. 八项退出判据对照矩阵

| 判据编号 | 判据定义 | 达成目标 | 落地代码位置 | 验证测试与证据链 | 结论 |
|---------|---------|----------|-------------|-----------------|------|
| **#1** | **Webhook 合同**：secret 校验 fail-closed；无/错 secret 不可被当成合法 Update；有测试 | GOAL-003 | `apps/api/internal/channel/telegram/webhook.go` | `apps/api/internal/channel/telegram/webhook_test.go`：<br>- `TestWebhook_UnconfiguredToken_Returns503`<br>- `TestWebhook_SecretValidation_FailClosed`<br>- `TestWebhook_MalformedJSON_Returns400` | **PASS** |
| **#2** | **分发端口**：命令与 callback 的 Register/Unregister + 分发有测试；未知命令有确定回落（不 panic） | GOAL-002<br>GOAL-003 | `apps/api/kernel/telegram.go`<br>`apps/api/internal/channel/telegram/dispatcher.go` | `apps/api/internal/channel/telegram/dispatcher_test.go`：<br>- `TestDispatcher_RegisterAndDispatchCommand`<br>- `TestDispatcher_RegisterAndDispatchCallback`<br>- `TestDispatcher_InvalidRegistrations`<br>`apps/api/internal/channel/telegram/webhook_test.go`：<br>- `TestWebhook_UnknownCommand_SendsFallbackMessage`（发 `DefaultTelegramUnknownCommandText`） | **PASS** |
| **#3** | **出站端口**：`SendMessage` 文本可测（mock 供应商）；生产供应商不把 Bot 客户端类型漏进模块公共契约 | GOAL-002<br>GOAL-004 | `apps/api/kernel/telegram.go`<br>`apps/api/internal/channel/telegram/http_sender.go`<br>`apps/api/internal/channel/telegram/capture_sender.go` | `apps/api/internal/channel/telegram/http_sender_test.go`：<br>- `TestHTTPSender_ValidMessageDelivery`<br>- `TestHTTPSender_UnconfiguredToken_DowngradesToMock`<br>- `TestHTTPSender_ValidationFailClosed`<br>- `TestHTTPSender_TimeoutBudget`（10s 超时）<br>- 公共签名仅出现 `kernel.TelegramMessage`/`kernel.TelegramSender`，无 SDK 泄露 | **PASS** |
| **#4** | **身份映射**：同一 `telegram_user_id` 多次 get-or-create 得到同一 `subject_id`；不写 `admin.users`；不依赖 `admin.wallet` HTTP 路由 | GOAL-003 | `apps/api/internal/channel/telegram/webhook.go`<br>`apps/api/internal/composition/composition.go` | `apps/api/internal/channel/telegram/webhook_test.go`：<br>- `TestWebhook_CommandDispatch_AndSubjectMapping`<br>- `TestWebhook_SubjectMappingIdempotency`<br>- `composition.go` 独立构造 `subject.NewStore(st)`，即使在未启用 `admin.wallet` HTTP 的环境中仍可独立解析 | **PASS** |
| **#5** | **设置与密钥**：Admin 可配置 token/secret；密钥 fail-closed；不进配置包明文 | GOAL-004 | `apps/api/internal/channel/telegram/runtime.go`<br>`apps/api/internal/channel/telegram/settings_handler.go` | `apps/api/internal/channel/telegram/runtime_test.go`：<br>- `TestRuntimeManager_HotSwitch`（并发安全热切换）<br>- `TestSettingsHandler`（GET 输出掩码星号脱敏字段，PATCH 成功热切换） | **PASS** |
| **#6** | **限流评估落盘**：评估 VP-027 limiter 对 webhook/chat_id/user_id 足够；三桶请求计数映射落地 | GOAL-001<br>GOAL-002<br>GOAL-003 | `apps/api/internal/channel/telegram/webhook.go` | VRev-070 §6 书面评估通过。<br>`apps/api/internal/channel/telegram/webhook_test.go`：<br>- `TestWebhook_RateLimiting_IPBucket`（60/min 429+Retry-After）<br>- `TestWebhook_RateLimiting_ChatBucket`（30/min 429+Retry-After）<br>- `TestWebhook_RateLimiting_UserBucket`（20/min 429+Retry-After）<br>- `TestWebhook_RateLimiting_IPRecordsOnSecretFailure`（防洪水） | **PASS** |
| **#7** | **边界保持**：未改 Charter；未进默认集；未做 Mini App/Stars/FSM/付费命令；未引入 Redis/SDK；未重开历史 VP | GOAL-001～005 | `apps/api/kernel/profile.go`<br>`apps/api/internal/composition/composition_telegram_test.go` | - `profileDefaults["mvp"]` / `["admin"]` / `["demo"]` 均不含 `channel.telegram`<br>- `git status` / `git diff` 核实无 Charter 变更<br>- 无 `github.com/go-telegram-bot-api` 等依赖<br>- 无 Redis 依赖 | **PASS** |
| **#8** | **审计闭合**：开放 required finding = 0 | GOAL-002～005 | 工作区各目标 `03-audit/` | - GOAL-002: A-001 self pass (0 required)<br>- GOAL-003: A-001 self pass, A-002 independent fail (1 required F-001) -> A-003 fixed 闭合 (0 required)<br>- GOAL-004: A-001 self pass (0 required)<br>- GOAL-005: 关门双审 (0 required) | **PASS** |

---

## 2. 边界红线合规核查清单

1. **默认集隔离**：`channel.telegram` 作为 `BuiltinModules()` 候选模块，未被任何预置 Profile（`mvp`、`admin`、`demo`）默认引入。仅在显式声明的自定义 Plan 中通过 `plan.HasModule("channel.telegram")` 装配（`TestTelegramChannelComposition` 证实）。
2. **零第三方 SDK**：出站与入站均基于 Go 标准库 `net/http` 与 `encoding/json`，未引入任何 Telegram Bot 第三方封装库。
3. **零 Redis 依赖**：入站三桶限流完全基于 VP-027 进程内 `kernel.RateLimiterProvider`。
4. **零业务/Mini App 越界**：出站仅支持文本与 `callback_data` InlineKeyboard；分发仅支持命令与回调查询；无 Mini App、无 Stars 付费、无复杂对话状态机。
5. **用户隔离**：Telegram 用户主体直接写入 `subjects` 表（`issuer=telegram`），绝对不写入 `admin.users` 系统后台管理员表。
