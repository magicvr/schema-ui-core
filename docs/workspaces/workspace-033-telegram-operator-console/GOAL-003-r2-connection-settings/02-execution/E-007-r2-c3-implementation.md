---
doc_type: goal-execution
id: E-007-r2-c3-implementation
parent: GOAL-003-r2-connection-settings
date: 2026-09-04
status: done
version: 0.1.0
---

# E-007 · R2 C3 Bot API、connection manager 与 Fx 生命周期实现事实

## 已发生事实

- 在 `9bc825ba` 中新增内部 `BotAPIClient`：实现 `getMe`、`setWebhook`、`deleteWebhook`、`getUpdates`；管理调用采用现有短超时，polling 使用独立 `30s` request timeout 与 `40s` HTTP client timeout；响应体有上限，HTTP/protocol/`ok=false` 均 fail closed，错误不包含 bot token。
- 在 `9bc825ba` 中新增 `ConnectionManager`：单 owner 串行化 Start/Reconcile/Stop/lease 操作；webhook 按 `getMe → setWebhook` 建立，polling 按 `getMe → deleteWebhook` 后按业务 handler/lease demand 启动；模式切换先 drain；缺 secret/URL、Bot API 错误和 polling handler 非限流错误均进入诊断 error，不伪造 `running`。
- `Dispatcher.HasBusinessHandlers`、polling update 共用 webhook 的限流/subject/dispatch 路径、20 秒 lease 与 1 秒 sweep 已接入；polling 空结果继续等待，正常 context cancel 完成 drain；异步 polling 错误退出会清理 receiver handles。
- `RuntimeManager` 新增非敏感连接状态、settings patch 的 updateMu 内合并及 settings-changed callback；Admin settings 持久化成功后触发 manager reconcile，避免并发 PATCH 在锁外读合并。
- Fx composition 将同一个 Telegram runtime 的 manager 接入 `channel.telegram` Start/Ready/Stop，并保留统一 runtime shutdown drain；组合测试通过注入本地 Fake Bot API endpoint 验证启动接缝。
- 新增 Bot API、连接建立、空结果、dispatch/drain、异步错误、lease 过期、失败切换无双活、热切换与 Fx 相关测试；C3 不包含 Admin settings UI、HTTP heartbeat lease surface 或 R3 会话落盘。

## 验证

- `go test ./internal/channel/telegram ./internal/composition -count=1 -timeout=120s`：通过。
- `go test -race ./internal/channel/telegram -count=1 -timeout=120s`：通过。
- `gofmt` 与 `git diff --check`：通过；实现提交 `9bc825ba`。

## 状态边界

C3 实施事实已完成，但因该 scope 属于 connection/production lifecycle 高影响门禁，独立审计尚未完成；本条不把 C3 checkpoint 或 GOAL-003 标为完成，当前仍为 `active · 2/5`。
