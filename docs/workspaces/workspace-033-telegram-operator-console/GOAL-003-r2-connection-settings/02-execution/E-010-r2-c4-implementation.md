---
doc_type: goal-execution
id: E-010-r2-c4-implementation
parent: GOAL-003-r2-connection-settings
date: 2026-09-04
status: done
version: 0.1.0
---

# E-010 · R2 C4 Admin 设置页与 polling lease 实现事实

## 已发生事实

- 提交 `d95f7544` 接入了 Telegram console lease：`POST /api/channel/telegram/lease/acquire`、`heartbeat`、`release`；三个路由均保持 module-gated、`Public: false`，并由认证身份的 `SessionID` 作为唯一租约键，不接受客户端提供的 session id。
- lease handler 使用既有 `settings.read` 权限；acquire/heartbeat/release 分别调用同一个 `ConnectionManager` 的对应方法，响应只返回租约状态、TTL、当前 lease 数和非 secret connection 状态。
- composition 将 lease handler 与 settings handler 一起包在认证 middleware 中，并传入同一个 Fx 构造的 `TelegramRuntime.Connection`；未启用 `channel.telegram` 的 profile 不注册 settings、lease 或 webhook routes。
- Telegram Admin custom component 已支持 mode、显式 `webhook_public_base_url`、write-only token/secret、connection state/receiver/bot username/error 展示；polling 模式页面打开时 acquire，每 10 秒 heartbeat，组件清理时串行 release。
- 前端 lease action 通过本地 promise queue 串行化，避免卸载时 release 与尚未完成的 heartbeat 逆序导致租约被重新创建。
- 增补了 lease 认证/会话隔离测试、provider route declaration/registration 测试、真实 composition mux lease 接线测试、disabled profile 全 Telegram surface 404 测试，以及 UI polling 状态、mode/URL PATCH 测试。

## 验证

- `go test ./internal/channel/telegram ./modules/channel/telegram ./internal/composition -count=1 -timeout=180s` 通过。
- `go test -race ./internal/channel/telegram -count=1 -timeout=180s` 通过。
- `npm test -- --run src/components/telegram-admin-tab.test.tsx` 通过（5 tests）。
- 两份 i18n JSON 可解析，`git diff --check` 通过；构建曾触发既有 `apps/web/src/renderer/form-controls.tsx:946-947` 的 `number | undefined` → `string | undefined` 基线类型错误，未由本次改动引入。

## 决策边界

- 本轮已向用户询问 lease HTTP route/permission 方案，但交互未返回选择；依据当前运行时的继续执行规则，实际采用 AI 推荐的“三端点 + `settings.read` + 服务端派生 session”作为实现默认。此事实不冒充用户已接受的产品裁决；如用户后续选择其他契约，应回到 `/govern` 记录变更并同步实现。
- C4 实现事实已完成，独立审计尚未执行；GOAL-003 仍保持 `active · 3/5`，C5 尚未关闭。
