---
doc_type: goal-execution
id: E-009-r3-c2-implementation
parent: GOAL-004-r3-session-operator-console
date: 2026-09-04
source: Codex govern
status: done
version: 0.1.0
---

# E-009 · R3 C2 入站持久化实现

## 已发生事实

- 在 `apps/api/modules/channel/telegram/migration/migration.go` 增加 v68 `telegram_ingress`，为 SQLite 与 PostgreSQL 创建 `telegram_sessions`、`telegram_inbound_messages` 及活动索引；同步 catalog checksum、restore fingerprint、fresh/upgrade/restart 断言。
- 在 `apps/api/modules/channel/telegram/store/repository.go` 增加模块内 repository：同一 `Store.Run` 中使用 `ON CONFLICT DO NOTHING` + `RowsAffected` 记录 `(bot_id, update_id)`，新收据才 upsert 会话；重复路径不改会话、不分发。
- Webhook 与 polling 共用规范化入站路径；支持非空 text/command/callback，跳过空文本/未建模媒体；运行时 bot identity 来自已确认的 `getMe` 状态；主体映射位于唯一收据之前。
- polling 已将成功接受更新的 offset 推进移到 handler 成功之后；限流仍按既有兼容语义推进，其他入站持久化失败进入 error 状态并退出当前 polling 循环。
- 增加 SQLite repository/共同 polling/webhook/error-path 测试，以及 gated PostgreSQL 首次、重复和 8 路并发唯一竞争测试；私聊缺少 chat title 时用发送者姓名补齐会话 title。

## 验证与产物

- 定向 `go test ./internal/channel/telegram ./modules/channel/telegram/... ./internal/composition ./internal/store -count=1` 通过。
- `go test -race ./internal/channel/telegram ./modules/channel/telegram/store` 通过。
- gated `go test ./internal/store -run '^TestPostgresTelegramIngressRepositoryIdempotency$' -count=1 -v` 通过；未配置 PG 时该测试遵守既有 skip 门禁。
- `go vet ./...` 与 `go build ./...` 通过。
- 实现 checkpoint：`72486d59`（`feat(telegram): persist inbound sessions and receipts`）。
- 全 API 回归首次并行运行中仅出现一次 `TestShutdownDrainHarnessPostgres` in-flight `EOF`；随后同一测试隔离重跑通过。该一次性现象不被记录为 C2 代码通过证据，最终关门前需再跑完整回归。

## 状态边界

C2 代码与定向验证已完成，但尚未完成 A-012 之后的 Grok independent 实现审计；本条不修改 C2 检查点进度或 R3 状态。
