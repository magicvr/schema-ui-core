---
doc_type: goal-execution
id: E-008-a004-final-closure
parent: GOAL-001-telegram-channel-runtime
date: 2026-09-03
status: recorded
---

# E-008 · 响应 A-004 复审意见：单例 Dispatcher 接进 newMux 与 Catalog 迁移加密持久化

## 1. 实施事实

落实用户指令 `/govern 响应 GOAL-001 A-004：A-003 对 F-001/F-002 的 closed 不成立；先把同一 dispatcher 接进 newMux，F-002 走 catalog+加密，不要再叠 CREATE TABLE IF NOT EXISTS`：

1. **单例 Dispatcher 与 Fx/newMux 统一注入（F-001 彻底闭合）**：
   - 在 `apps/api/internal/composition/composition.go` 定义 `TelegramRuntime` 与 `newTelegramRuntime`，纳入 Fx 依赖注入列表。
   - `newMux` 直接接收 `*TelegramRuntime`，把 Fx 注入的同一 `tr.Dispatcher` 与 `tr.Sender` 装配给 WebhookHandler。
   - 删除第二套 dispatcher 工厂，消除双实例分叉。
   - `composition_telegram_test.go` 新增 `TestTelegramChannelComposition_RealWebhookMount`，验证注册在共享 dispatcher 上的命令直接被 real mux 上的 webhook 处理。
2. **Catalog 迁移 66 与 AES-GCM 加密落库（F-002 彻底闭合）**：
   - 彻底删除 `apps/api/internal/channel/telegram/runtime.go` 中所有 `CREATE TABLE IF NOT EXISTS` 运行时 DDL。
   - 新增 `apps/api/modules/channel/telegram/migration/`，定义版本 66 的 `telegram_config` 迁移（sqlite 与 postgres），校验码 `e330d7859636a5344edff165ebf9e0cfe96dfd43bc127e6d3a36cda5ed936601`，编入 `modules/compiled/persistence.go`。
   - 平台指纹与迁移基线更新：`identity.go` 更新 `completeFingerprintCatalogHead = 66`，`lockedHeadExtraTables[66] = []string{"telegram_config"}`，`migrate_test.go` 与 `operations_test.go` 对齐 66 迁移。
   - `RuntimeManager` 引入 AES-256-GCM 对 `bot_token_enc` 与 `webhook_secret_enc` 进行加密存储与解密读取。
   - `Update` 先执行数据库加密持久化，失败立即 fail-closed 报错，不改内存；仅在落库成功后刷新内存。
   - `composition_telegram_test.go` `TestTelegramRuntime_PersistenceAcrossRestart` 验证进程重启后重新开启 store 能读回加密持久化的配置。

## 2. 验证结果

- `go test ./...` 100% 通过（包含 `store`、`composition`、`kernel`、`internal/channel/telegram`、`modules/channel/telegram` 等全部包）。
- 开放必改项归零，A-005 合并响应落盘。
