---
doc_type: goal-audit
id: A-005-a004-closure-response
parent: GOAL-001-telegram-channel-runtime
date: 2026-09-03
source: self
scope: GOAL-001 A-004 独立复审意见响应（F-001 单例挂载与 F-002 全局 Catalog 加密落地）
audit_type: stage-closeout
verdict: pass
open_required: 0
---

# A-005 · Root GOAL-001 A-004 独立复审意见响应与必改闭合（合并响应）

## 1. 响应背景

编排器（/govern）响应独立复审意见 [A-004-independent-closure-reaudit.md](A-004-independent-closure-reaudit.md)（grok-4.6 · reasoning high · `verdict: fail`，2 required F-001/F-002 维持 open）。依据用户明确指令：「先把同一 dispatcher 接进 newMux，F-002 走 catalog+加密，不要再叠 CREATE TABLE IF NOT EXISTS」，完成彻底整改。

## 2. 必改项彻底闭合台账（Required Findings）

| ID | 严重度 | 闭合路径 | 闭合事实与核验代码 | 状态 |
|----|--------|----------|-------------------|------|
| **F-001** | med / required | **fixed** | 1. **单例端口与 Fx 进程级注入**：在 `apps/api/internal/composition/composition.go` 中定义 `TelegramRuntime` 与 `newTelegramRuntime`，纳入 `fx.Provide(newTelegramRuntime, func(tr *TelegramRuntime) kernel.TelegramDispatcher, func(tr *TelegramRuntime) kernel.TelegramSender)`。<br>2. **消除双重 Dispatcher**：`newMux` 接收 `tr *TelegramRuntime`，直接将 Fx 注入的同一 `tr.Dispatcher` 挂载到 WebhookHandler，彻底杜绝第二套 Dispatcher 实例。<br>3. **未启用状态规范化**：未启用时统一返回 `NewDisabledDispatcher()` 与 `NewDisabledSender()`（`Send` 统一返回 `kernel.ErrTelegramDisabled`）。<br>4. **端到端装配测试**：`composition_telegram_test.go` 新增 `TestTelegramChannelComposition_RealWebhookMount`，驱动真实 mux，验证在共享 dispatcher 上注册的命令直接被真实挂载的 webhook 处理。 | **closed** |
| **F-002** | med / required | **fixed** | 1. **消除运行时 DDL**：彻底删除所有 `CREATE TABLE IF NOT EXISTS`。<br>2. **全局迁移 Catalog 纳入**：新增 `apps/api/modules/channel/telegram/migration/`（定义 `telegram_config` v66，校验码 `e330d7859636a5344edff165ebf9e0cfe96dfd43bc127e6d3a36cda5ed936601`），编入 `modules/compiled/persistence.go`，并在 `provider.go` 的 `CompiledPersistence()` 导出。<br>3. **平台台账锁更新**：`identity.go` 更新 `completeFingerprintCatalogHead = 66`，`lockedHeadExtraTables[66] = []string{"telegram_config"}`，`migrate_test.go` 与 `operations_test.go` 对齐 66 迁移基线。<br>4. **落库加密（At-Rest Encryption）**：`RuntimeManager` 采用 AES-256-GCM 对 `bot_token_enc` 与 `webhook_secret_enc` 进行加密存储（对标 Mail 通道先例）。<br>5. **写库优先与 Fail-Closed**：`Update` 优先执行数据库加密持久化，失败立即返回 error 且不变更内存（杜绝分叉）；仅在成功后刷新内存。<br>6. **真实重载测试**：`composition_telegram_test.go` `TestTelegramRuntime_PersistenceAcrossRestart` 模拟真实进程重启，关闭 store 后重新开启新实例，断言 `GetToken()` / `GetSecret()` 完整读回持久值。 | **closed** |

## 3. 建议项状态核销

- **R-001（恒时比较）**：SHA-256 + 32-byte constant-time compare，已确认闭合。
- **R-002（真实 Webhook 装配测试）**：`TestTelegramChannelComposition_RealWebhookMount` 真实断言 401、200 与命令分发，已闭合。
- **R-003（Dispatch 错误日志）**：`slog.Warn` 结构化日志，已确认闭合。
- **R-004（Bot API ok 校验）**：`http_sender.go` 严格解析 `botAPIResponse.OK`，已闭合。
- **R-005（秘密片段不回显）**：`RuntimeStatus` 仅输出 `TokenSet: bool` 与 `SecretSet: bool`，已闭合。
- **R-006（敏感键导出排除）**：`configpkg.go` 纳入 `sensitiveFields` 登记表并自动剔除，已闭合。
- **R-007（Allow/Record 窗口）**：记录为 VP-027 端口形状残余风险，用户已书面知悉并接受。
- **R-008（字段名）**：对齐 `captured_messages_count`，已闭合。

## 4. 结论与关门事实

- A-004 中指出的 F-001 与 F-002 必改项现已通过真实代码修改、全局 Catalog 迁移纳入及集成测试全面闭合。
- 开放 required findings：**0**。
- 全仓回归测试 `go test ./...` 100% 通过。
- Root 目标 `GOAL-001-telegram-channel-runtime` 顺利完成治理闭环。
