---
doc_type: goal-execution
id: E-004-r2-closeout
parent: GOAL-001-telegram-channel-runtime
date: 2026-09-03
status: recorded
---

# E-004 · R2 阶段关门（Webhook 路由、分发、身份映射与入站限流）

## 1. 阶段事实

- 子目标 [GOAL-003-r2-webhook-dispatch-identity](../GOAL-003-r2-webhook-dispatch-identity/00-meta.md) 顺利关门（`status: done`，3/3）。
- 关键交付：
  1. Webhook 管道：`POST /api/channel/telegram/webhook`，`X-Telegram-Bot-Api-Secret-Token` 常时校验 fail-closed，未配置 token 返回 503，畸形包 400，合法包 200。
  2. 入站三桶限流：IP（60/m）、Chat（30/m）、User（20/m）独立 limiter，超限 429 + `Retry-After`，IP 洪水在 Secret 失败时仍记账。
  3. 主体映射：`GetOrCreateSubject("telegram", user_id)` 幂等映射至 `upd.SubjectID`，存储失败返回 500 fail-closed 促使重试，不写入 `admin.users`，不依赖 `admin.wallet` HTTP 路由。
  4. 分发引擎：命令去 `/` 与 `@BotName` 精确分发，未注册命令发送 `DefaultTelegramUnknownCommandText` 常量回落；Callback 精确匹配。
  5. 模块装配：`channel.telegram` 编入 `BuiltinModules()`，在 `composition.go` 中按 Plan 动态装配；默认 Profile 不含此模块，保持隔离。
- 独立交叉审计闭环：本地 grok build（grok-4.6 · reasoning high）出具 A-002（指出 F-001 模块候选集与装配缺口），已全部修复并在 A-003 合法闭合。
- Root 纲领 R2 达成，Root progress 推进至 **2/4**。
