---
doc_type: goal-execution
id: E-004-r2-c2-implementation
parent: GOAL-003-r2-connection-settings
date: 2026-09-04
status: done
version: 0.1.0
---

# E-004 · R2 C2 配置与持久化实现事实

## 已发生事实

- 在 `245e763d` 中完成 C2 生产接缝：新增 Telegram v67 additive migration，为既有 `telegram_config` 增加 `mode` 与 `webhook_public_base_url`；保留 v66 migration，并更新 migration catalog、head/重启断言与 provider persistence 贡献测试。
- `internal/config` 已接入 Telegram mode、webhook public base URL 的 YAML/env/default/export 配置与 fail-closed 校验；配置导出只包含非敏感连接字段，bot token、webhook secret 与 master key 不进入导出内容。
- Telegram runtime 通过 `NewRuntimeManagerWithSettings` 读取或初始化单一持久化行；无行时使用 YAML/env seed，已有行（包括新字段为空）为权威来源。token/secret 仍加密持久化；`UpdateSettings` 在持久化成功后才发布内存快照，settings PATCH 支持完整连接设置的部分更新。
- 新增 runtime/settings、config、composition authority 测试；在 `4cec07f` 中补充并固定“既有行的 mode/URL 为空时，重启不得被 stale seed 复活”的回归测试。
- C2 未实现 Bot API client、polling/webhook connection manager、会话租约/heartbeat、Fx 连接生命周期或管理 UI；这些仍属于后续 C3/C4 事实边界。

## 验证

- `apps/api`：`go test ./internal/channel/telegram ./internal/config ./internal/composition ./modules/channel/telegram ./internal/store ./cmd/schema-ui ./server`：通过（C2 实现 checkpoint）。
- `apps/api`：`go test ./internal/composition -run 'TestTelegramRuntime_(ConnectionSettingsPersistenceAndAuthority|EmptyConnectionSettingsRemainAuthoritative)$' -count=1`：通过（C2 authority 回归补充）。
- `git diff --check`：通过；提交前工作树仅包含本次 C2 测试变更，未使用全量暂存。
