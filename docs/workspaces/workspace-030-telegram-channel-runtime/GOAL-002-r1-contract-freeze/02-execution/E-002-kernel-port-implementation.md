---
doc_type: goal-execution
id: E-002-kernel-port-implementation
parent: GOAL-002-r1-contract-freeze
date: 2026-09-03
status: recorded
---

# E-002 · Telegram 通道内核端口与测试落地（C2）

## 1. 目标与范围

依据 [D-002-telegram-channel-contract.md](../01-decision/D-002-telegram-channel-contract.md) 冻结规范，落地 Go 内核级 Telegram 通道端口定义与合同级快测：
- `apps/api/kernel/telegram.go`
- `apps/api/kernel/telegram_test.go`

## 2. 实施事实

1. **类型与接口**：
   - `TelegramInlineButton`：定义 `Text` 与 `CallbackData`（≤64 字节约束）。
   - `TelegramMessage`：定义 `ChatID`、`Text`、`Buttons [][]TelegramInlineButton`。实现 `Validate() error` 表驱动校验。
   - `TelegramSender`：`Send(ctx context.Context, msg TelegramMessage) error` 同步接口。
   - `TelegramUpdate`：`ChatID`、`UserID`、`SubjectID`、`Command`、`Text`、`CallbackData` 薄入站视图。
   - `TelegramHandler`：`func(ctx context.Context, upd TelegramUpdate) error`。
   - `TelegramDispatcher`：`RegisterCommand` / `UnregisterCommand` / `RegisterCallback` / `UnregisterCallback` 静态注册接口。
2. **校验与规范化辅助工具**：
   - `NormalizeTelegramCommand(raw string) (string, error)`：剥离前导 `/` 与 `@BotName` 后缀，精确匹配校验。
   - `ValidateTelegramCallback(data string) error`：校验 non-empty 与 ≤64 字节上限。
   - 错误常量及默认回落文案：`DefaultTelegramUnknownCommandText`（"Sorry, unrecognized command."）与哨兵错误（`ErrTelegramDisabled`、`ErrTelegramHandlerNil` 等）。
3. **合同级快测**：
   - `TestTelegramMessage_Validate`：覆盖缺少 ChatID、非数字 ChatID、合法负 ChatID、空文本、空白文本、按钮空字段、超长 CallbackData（64 字节边界）等 12 组表驱动测试用例。
   - `TestNormalizeTelegramCommand`：覆盖前导斜杠、@Bot 标头、空串、仅斜杠、内部空格/斜杠等 9 组用例。
   - `TestValidateTelegramCallback`：覆盖空串、64 字节、65 字节边界断言。
   - `TestStubDispatcher`：验证 stub 实现的注册、冲突拦截、执行与注销。
4. **测试结果**：
   - `go test -v ./kernel/...` 全部 PASS。

## 3. 产物清单

- `apps/api/kernel/telegram.go`
- `apps/api/kernel/telegram_test.go`
