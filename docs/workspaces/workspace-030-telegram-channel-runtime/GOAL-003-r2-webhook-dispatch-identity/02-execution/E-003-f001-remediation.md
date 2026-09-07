---
doc_type: goal-execution
id: E-003-f001-remediation
parent: GOAL-003-r2-webhook-dispatch-identity
date: 2026-09-03
status: recorded
---

# E-003 · 响应独立审计 A-002（F-001 必改闭合与建议整改）

## 1. 响应与整改事实

针对独立交叉审计 [A-002-r2-independent-audit.md](../03-audit/A-002-r2-independent-audit.md) 指出的问题完成全面整改：

1. **F-001（high / required · 闭合：fixed）**：
   - `apps/api/kernel/profile.go`：将 `channel.telegram` 编入 `BuiltinModules()` 候选集，声明 `DependsOn: ["core.server-registration"]`，`Requires: [CapabilityHTTP]`，`Contributions.Routes: ["POST /api/channel/telegram/webhook"]`。严格未加入 `mvp`/`admin`/`demo` 默认集。
   - `apps/api/internal/composition/composition.go`：在 `newServer` 装配管道中增加 `plan.HasModule("channel.telegram")` 分支，构造独立三桶限流器、直接基于平台事务构造 `subject.NewStore(st)`（不依赖 `admin.wallet` HTTP 是否启用），构造 WebhookHandler 并向 providers 登记 `telegrammodule.New(tgWebhook)`。
   - `apps/api/internal/composition/composition_telegram_test.go`：增加端到端装配测试 `TestTelegramChannelComposition`，验证内置模块注册、默认 Profile 隔离、自定义 Profile 激活、Plan 展开与 Webhook 路由注册及调用。
2. **R-001（recommended · 实施：fixed）**：
   - `apps/api/internal/channel/telegram/webhook.go`：当 `subjectStore.GetOrCreateSubject` 返回错误时，不再静默忽略，而是立即返回 **500 Internal Server Error**，保证在数据库或存储故障时 fail-closed，促使 Telegram Bot API 重试投递，避免以空 SubjectID 执行业务。
3. **R-002（recommended · 实施：fixed）**：
   - `apps/api/internal/channel/telegram/webhook_test.go`：新增 `TestWebhook_RateLimiting_ChatBucket`，验证 Chat 桶 30/min 上限与第 31 次请求触发 429 Too Many Requests + `Retry-After` 头。
   - 同时增加 `TestWebhook_SubjectMappingIdempotency`，验证同一用户多次调用得到相同 SubjectID。
4. **R-003（recommended · 实施：fixed）**：
   - `apps/api/internal/channel/telegram/dispatcher_test.go`：新增 `TestDispatcher_InvalidRegistrations`，覆盖 nil handler、空命令、非法命令（含空格/斜杠）、空 callback、超长 callback（>64字节）等边界报错断言。
5. **R-004（recommended · 实施：fixed）**：
   - `apps/api/modules/channel/telegram/provider.go`：将 Descriptor 从标准 Admin 六面模式裁剪为横切通道模式（依赖 `core.server-registration`，要求 `CapabilityHTTP`，无 navigation/schema 依赖）。
6. **环境配置同步**：
   - `apps/api/internal/config/config.go`：增加 `TelegramBotToken` 与 `TelegramWebhookSecret`（支持 YAML 与 `TELEGRAM_BOT_TOKEN`/`TELEGRAM_WEBHOOK_SECRET` 环境变量）。
   - `apps/api/configs/.env.example`：补充环境变量说明文档，通过 `TestCanonicalEnvExample` 测试。

## 2. 验证结果

- `go test ./internal/channel/telegram/... ./modules/channel/telegram/... ./internal/composition/... ./kernel/... ./internal/config/...` 全部 PASS。
- 全仓 `go test ./...` 全绿。
