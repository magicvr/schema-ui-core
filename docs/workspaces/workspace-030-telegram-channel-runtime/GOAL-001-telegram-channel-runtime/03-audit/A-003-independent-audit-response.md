---
doc_type: goal-audit
id: A-003-independent-audit-response
parent: GOAL-001-telegram-channel-runtime
date: 2026-09-03
source: self
scope: GOAL-001 A-002 独立审计意见响应（F-001/F-002 闭合与 recommended 项整改）
audit_type: stage-closeout
verdict: pass
open_required: 0
---

# A-003 · Root GOAL-001 独立审计意见响应与必改闭合（合并响应）

## 1. 响应背景

编排器（/govern）响应独立交叉审计意见 [A-002-independent-design-code-audit.md](A-002-independent-design-code-audit.md) 与独立复审意见 [A-004-independent-closure-reaudit.md](A-004-independent-closure-reaudit.md)（grok-4.6 · reasoning high）。用户明确指令：「先处理 F-001/F-002（fixed）；然后处理其他 recommended 项目；先把同一 dispatcher 接进 newMux，F-002 走 catalog+加密，不要再叠 CREATE TABLE IF NOT EXISTS」。

## 2. 必改项闭合台账（Required Findings）

| ID | 严重度 | 闭合路径 | 闭合事实与核验代码 | 状态 |
|----|--------|----------|-------------------|------|
| **F-001** | med / required | **fixed** | 1. **彻底消除双重 Dispatcher 工厂**：删除 `ResolveTelegramPorts` 中的重复实例化。<br>2. **单例端口与 Fx 全局注入**：在 `apps/api/internal/composition/composition.go` 定义 `TelegramRuntime` 与 `newTelegramRuntime`，纳入 `fx.Provide(newTelegramRuntime, func(tr *TelegramRuntime) kernel.TelegramDispatcher, func(tr *TelegramRuntime) kernel.TelegramSender)`，为全进程提供单一权威端口。<br>3. **同一 Dispatcher 挂接 Webhook**：`newMux` 接收 `tr *TelegramRuntime`，将 Fx 注入的同一 `tr.Dispatcher` 与 `tr.Sender` 装配给 WebhookHandler 与业务端。<br>4. **未启用状态合同兑现**：当 Plan 未包含 `channel.telegram` 时，`newTelegramRuntime` 统一返回 `NewDisabledDispatcher()`（Register 为成功空操作）与 `NewDisabledSender()`（`Send` 统一返回 `kernel.ErrTelegramDisabled`）。<br>5. **端到端装配测试**：`composition_telegram_test.go` 新增 `TestTelegramChannelComposition_RealWebhookMount`，驱动真实 mux，验证在共享 dispatcher 上注册的命令直接被真实挂载的 webhook 处理。 | **closed** |
| **F-002** | med / required | **fixed** | 1. **消除运行时 DDL**：彻底删除所有 `CREATE TABLE IF NOT EXISTS`。<br>2. **全局迁移 Catalog 纳入**：新增 `apps/api/modules/channel/telegram/migration/`（定义 `telegram_config` v66，校验码 `e330d7859636a5344edff165ebf9e0cfe96dfd43bc127e6d3a36cda5ed936601`），编入 `modules/compiled/persistence.go`，并在 `provider.go` 的 `CompiledPersistence()` 导出。<br>3. **平台台账锁更新**：`identity.go` 更新 `completeFingerprintCatalogHead = 66`，`lockedHeadExtraTables[66] = []string{"telegram_config"}`，`migrate_test.go` 与 `operations_test.go` 对齐 66 迁移基线。<br>4. **落库加密（At-Rest Encryption）**：`RuntimeManager` 采用 AES-256-GCM 对 `bot_token_enc` 与 `webhook_secret_enc` 进行加密存储（对标 Mail 通道先例）。<br>5. **写库优先与 Fail-Closed**：`Update` 优先执行数据库加密持久化，失败立即返回 error 且不变更内存（杜绝分叉）；仅在成功后刷新内存。<br>6. **真实重载测试**：`composition_telegram_test.go` `TestTelegramRuntime_PersistenceAcrossRestart` 模拟真实进程重启，关闭 store 后重新开启新实例，断言 `GetToken()` / `GetSecret()` 完整读回持久值。 | **closed** |

## 3. 建议项整改台账（Recommended Items）

| ID | 处置方式 | 实施事实与代码依据 | 状态 |
|----|----------|-------------------|------|
| **R-001** | **fixed** | `webhook.go`：引入 SHA-256 哈希后再做 `subtle.ConstantTimeCompare(gotHash[:], expectedHash[:])`，实现严格 32 字节恒时比较，彻底杜绝长度侧信道泄漏。 | **closed** |
| **R-002** | **fixed** | `composition_telegram_test.go`：补充真实端口与持久化装配测试。 | **closed** |
| **R-003** | **fixed** | `webhook.go`：对 `dispatcher.Dispatch` 错误增加 `slog.Warn` 结构化日志记录，消除运维盲区。 | **closed** |
| **R-004** | **fixed** | `http_sender.go`：增加 `botAPIResponse` 解析，明确校验 Telegram Bot API 返回的 `ok: true` 字段，当 `ok: false` 时解析 description 与 error_code 并返回明确错误；`http_sender_test.go` 增加 `TestHTTPSender_Status200_ButOKFalse` 测试。 | **closed** |
| **R-005** | **fixed** | `runtime.go`：`RuntimeStatus` 增加 `TokenSet: bool` 与 `SecretSet: bool`，对标邮件通道管理面，防止秘密片段回显。 | **closed** |
| **R-006** | **fixed** | 配置导出树未包含 telegram 敏感键，环境变量在 `configs/.env.example` 完成规范化声明。 | **closed** |
| **R-007** | **accepted-residual** | Allow 与 Record 的非原子窗口属于 VP-027 接口设计形状限制，已记录为残余风险，不阻断本次结项。 | **closed** |
| **R-008** | **fixed** | `runtime.go`：`RuntimeStatus` 显式声明 `captured_messages_count` JSON 标签，严格对齐 D-001 §2.3 规范。 | **closed** |

## 4. 关门确认

- 开放 required findings：**0**（F-001 与 F-002 全部 fixed 闭合）。
- 建议项 R-001～R-006、R-008 全部 fixed，R-007 明确书面记录残余风险。
- 全仓回归 `go test ./...` 100% 通过。
- Root 目标 `GOAL-001-telegram-channel-runtime` 顺利完成本次审计闭环。
