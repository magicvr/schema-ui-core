---
doc_type: goal-decision
id: D-001-r3-outbound-and-settings-architecture
parent: GOAL-004-r3-outbound-settings-limiter
date: 2026-09-03
status: accepted
---

# D-001 · R3 出站生产适配器与动态设置架构决策（2026-09-03）

## 1. 用户裁决事实

用户于 2026-09-03 书面裁决：
1. **I-030-005 配置生效方式**：**热切换（沿用 Mail 先例）**。
   - Telegram 通道设置（Bot Token 与 Webhook Secret）支持运行时并发安全动态更新与读取。
   - 新的 Token / Secret 更新后对后续请求即时生效，无需重启 API 进程。
2. **出站生产与 Mock 切换**：**基于 Token 自动降级**。
   - 配置有效 Bot Token 时：出站请求通过标准库 `net/http` 发送至 Telegram Bot API（`https://api.telegram.org/bot<token>/sendMessage`），严格应用 10s 超时控制与 `msg.Validate()`。
   - 未配置 Bot Token 时：出站请求自动安全降级至内存 `CaptureSender` 记录，保证开发、快测与无外部网络环境下单测零阻断。

## 2. 详细技术方案（冻结）

### 2.1 运行时通道管理器（`RuntimeManager`）

- 结构定义在 `apps/api/internal/channel/telegram/runtime.go`：
  - 管理动态配置：`BotToken` 与 `WebhookSecret`。
  - 使用 `sync.RWMutex` 保证并发安全读取与热切换。
  - 提供 `GetConfig() (token, secret string)` 与 `UpdateConfig(token, secret string)`。
  - `WebhookHandler` 的 `tokenGetter` 与 `secretGetter` 直接绑定到 `RuntimeManager`。
  - 支持从环境变量或初始配置（`config.Config`）完成第一代种子注入（Seed）。

### 2.2 出站生产适配器（`HTTPSender` / `AdaptiveSender`）

- 结构定义在 `apps/api/internal/channel/telegram/http_sender.go`：
  - 实现 `kernel.TelegramSender` 接口。
  - 内部组合 `RuntimeManager`、`http.Client`（带 10s 超时）与 `CaptureSender`。
  - `Send(ctx context.Context, msg kernel.TelegramMessage) error` 执行逻辑：
    1. 调用 `msg.Validate()`，违规直接 fail-closed 返回错误。
    2. 读取当前 Token：若 Token 为空，调用 `CaptureSender.Send` 记录，返回 nil。
    3. 若 Token 非空，构建 Telegram Bot API 请求：
       - URL: `https://api.telegram.org/bot<token>/sendMessage`
       - Payload JSON:
         ```json
         {
           "chat_id": "<chatID>",
           "text": "<text>",
           "reply_markup": {
             "inline_keyboard": [[{"text": "...", "callback_data": "..."}]]
           }
         }
         ```
       - `reply_markup` 仅在 `len(msg.Buttons) > 0` 时序列化输出。
       - 使用 `ctx` 并附加 10s 截止时间，发起 POST 请求。
       - 若 HTTP 响应状态码 >= 400，解析或提取错误并返回格式化 error。

### 2.3 Admin 设置与状态接口（判据 #5）

- 在 `apps/api/internal/channel/telegram/settings_handler.go` 中提供管理端点（或挂入模块路由）：
  - `GET /api/channel/telegram/settings`：返回当前通道状态（`configured: bool`、`token_masked: string`、`secret_masked: string`、`captured_messages_count: int`）。严禁输出完整密钥明文。
  - `PATCH /api/channel/telegram/settings`：接收新的 `bot_token` 与 `webhook_secret`，热更新至 `RuntimeManager`。
- 权限保护：需要 Admin 权限或在 Public 外部网络下默认不暴露。
- 保证密钥 fail-closed，且不随配置导出包导出明文。

### 2.4 限流核账

- 入站三桶限流已在 R2 严格落地于 `webhook.go`：
  - IP 桶（`tg:webhook:{ip}`，60/min）
  - Chat 桶（`tg:chat:{chat_id}`，30/min）
  - User 桶（`tg:user:{user_id}`，20/min）
- 出站端口为纯出站主动调用，由调用方驱动，核验确认无内部限流残留要求。
