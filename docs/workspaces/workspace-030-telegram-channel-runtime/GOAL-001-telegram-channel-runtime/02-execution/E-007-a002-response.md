---
doc_type: goal-execution
id: E-007-a002-response
parent: GOAL-001-telegram-channel-runtime
date: 2026-09-03
status: recorded
---

# E-007 · Root 独立审计 A-002 响应、F-001/F-002 修复与全项整改

## 1. 响应与实施事实

依据用户指令 `/govern 响应 GOAL-001 A-002：先处理 F-001/F-002（fixed）；然后处理其他recommended项目`，全面落实整改代码与测试：

1. **F-001（进程级端口装配与 disabled stub · fixed）**：
   - 新增 `apps/api/internal/channel/telegram/disabled.go`，实现 `DisabledSender`（`Send` 统一返回 `kernel.ErrTelegramDisabled`）与 `DisabledDispatcher`（`RegisterCommand`/`Callback` 返回 nil 成功空操作）。
   - 在 `apps/api/internal/composition/composition.go` 新增 `ResolveTelegramPorts(plan, cfg, st)`，在模块未启用时提供进程级 disabled stub，启用时注入共享的 live Dispatcher 与 HTTPSender。
   - `composition_telegram_test.go` 新增 `TestResolveTelegramPorts_EnabledAndDisabled` 覆盖全量断言。
2. **F-002（Admin 设置数据库持久化 · fixed）**：
   - `apps/api/internal/channel/telegram/runtime.go`：`RuntimeManager` 引入 `TxRunner`，在数据库中建立 `telegram_config` 表持久化存储 `bot_token` 与 `webhook_secret`，启动时自动从数据库重载，`Update` 时自动更新数据库。
   - `apps/api/internal/composition/composition.go`：装配时将 `st` 传入 `NewRuntimeManager`。
   - `composition_telegram_test.go` 新增 `TestTelegramRuntime_PersistenceAcrossRestart` 验证配置跨进程重启持久留存。
3. **Recommended 建议项整改**：
   - **R-001**：`webhook.go` 引入 SHA-256 哈希后再做 `subtle.ConstantTimeCompare`，实现严格 32 字节恒时比对。
   - **R-002**：`composition_telegram_test.go` 增加真实端口与真实数据库持久化装配测试。
   - **R-003**：`webhook.go` 对 `dispatcher.Dispatch` 错误增加 `slog.Warn` 结构化日志。
   - **R-004**：`http_sender.go` 增加 `botAPIResponse` 校验 Telegram Bot API 返回的 `ok: true` 字段，若 `ok: false` 返回明确错误描述与错误码，`http_sender_test.go` 增加测试覆盖。
   - **R-005**：`runtime.go` `RuntimeStatus` 增加 `TokenSet: bool` 与 `SecretSet: bool`，对标邮件通道管理面。
   - **R-006**：敏感环境变量在 `configs/.env.example` 完成声明。
   - **R-007**：记录 Allow/Record 非原子性为 VP-027 端口形状残余风险。
   - **R-008**：`RuntimeStatus` 字段名修正为 `captured_messages_count`。

## 2. 验证结果

- `go test ./internal/channel/telegram/... ./modules/channel/telegram/... ./internal/composition/... ./kernel/...` 全部 PASS。
- 全仓 `go test ./...` 100% 通过。
- 开放必改项归零，A-003 闭合响应完成。
