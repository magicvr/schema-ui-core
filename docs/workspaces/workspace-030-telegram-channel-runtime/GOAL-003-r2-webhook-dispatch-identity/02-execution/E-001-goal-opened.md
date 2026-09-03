---
doc_type: goal-execution
id: E-001-goal-opened
parent: GOAL-003-r2-webhook-dispatch-identity
date: 2026-09-03
status: recorded
---

# E-001 · GOAL-003 目标建立与 R2 启动

## 1. 目标建立

- 子目标 [GOAL-003-r2-webhook-dispatch-identity](../00-meta.md) 建立，承接 Root 纲领 R2 阶段。
- 核心范围：
  1. `POST /api/channel/telegram/webhook` 路由与 Secret 常时比较（fail-closed 401/503/400/200）。
  2. 入站三桶限流（IP 60/m, Chat 30/m, User 20/m）与 Record 永不 Clear 请求计数映射。
  3. `issuer=telegram` 外部主体映射（调用 `subject.Store.GetOrCreateSubject`），不依赖 `admin.wallet` HTTP 启用。
  4. 命令规范化分发（支持 `/cmd` 及 `/cmd@BotName`、冲突拦截、未注册命令发送 `DefaultTelegramUnknownCommandText` 回复）与 Callback 分发（`callback_data` 精确匹配）。
- 启动 C1：向用户提出 I-030-007 主体 Store 消费路径及内部实现方案裁决。
