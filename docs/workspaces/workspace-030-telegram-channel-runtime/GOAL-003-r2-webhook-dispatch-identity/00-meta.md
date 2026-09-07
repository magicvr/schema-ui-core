---
id: GOAL-003-r2-webhook-dispatch-identity
title: R2 Webhook 路由、Update 分发、主体映射与入站限流
status: done
parent: GOAL-001-telegram-channel-runtime
created: 2026-09-03
updated: 2026-09-03
version: 1.0.0
progress: 3/3
plan_refs:
  - VP-030-telegram-channel-runtime
primary_plan: VP-030-telegram-channel-runtime
serves_summary: 承载 VP-030 纲领 R2：落地 Telegram Webhook 路由（fail-closed 常时比较校验）、命令与 Callback Update 分发引擎、issuer=telegram 主体映射（不依赖 wallet HTTP 启用）及入站三桶限流（IP/chat_id/user_id）。A-002 F-001 必改项与建议整改项全部 closed 关门。
---

# GOAL-003 · R2 Webhook 路由、Update 分发、主体映射与入站限流

## 概述

承接 Root 纲领 **R2**（对应 VP-030 退出判据 #1/#2/#4 及 R1 合同 [GOAL-002 D-002](../GOAL-002-r1-contract-freeze/01-decision/D-002-telegram-channel-contract.md)）：
1. **Webhook HTTP 路由**（`POST /api/channel/telegram/webhook`）：Public 接口，无 token 时 503；Secret 头（`X-Telegram-Bot-Api-Secret-Token`）常时比较 fail-closed（401）；解析错误 400；合法 Update 返回 200 空包。
2. **入站三桶限流**（消费 VP-027 `RateLimiter`）：
   - IP 桶（`tg:webhook:{ip}`，60/min）
   - Chat 桶（`tg:chat:{chat_id}`，30/min）
   - User 桶（`tg:user:{user_id}`，20/min）
   - 顺序：无 token(503) → IP Allow/Record(429/401/400) → Chat Allow/Record(429) → User Allow/Record(429) → 身份映射 → 分发。
3. **主体映射**（I-030-007）：调用 `GetOrCreateSubject("telegram", user_id)`，获取并填入 `upd.SubjectID`；不依赖 `admin.wallet` HTTP 启用；不写入 `admin.users`。
4. **Update 分发**：
   - 命令分发：支持 `/cmd` 及 `/cmd@BotName` 规范化精准匹配，未注册命令返回 `DefaultTelegramUnknownCommandText`（通过 `TelegramSender` 回复）。
   - Callback 分发：`callback_data` 精确匹配；未注册 callback 静默 200。

## 纲领检查点（P-001）

| 检查点 | 内容 | 状态 |
|--------|------|------|
| C1 | **方案与信息裁决**：I-030-007 主体 Store 消费路径裁决；模块内部结构（Dispatcher/Webhook/Limiter）方案冻结（D-001） | **已关门**（2026-09-03 用户书面裁决：直接复用 subject.Store + internal/modules 分层，D-001） |
| C2 | **实现与回归**：Webhook 路由、Dispatcher 引擎、主体映射及限流装配落地与全量测试 | **已关门**（E-002 管道落地 + E-003 响应 A-002 修复 F-001 编入候选与 composition 装配） |
| C3 | **审计与关门**：自审与交叉审计（A-001 self + A-002 grok independent + A-003 闭合响应），无开放必改项关门 | **已关门**（A-001 pass + A-002 grok independent + A-003 fixed 闭合，0 required） |

`progress` = 已关门检查点数 / 3。当前 **3/3**。

## 信息就绪与未知项

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-030-007 | non-blocking | 主体 Store 消费路径：直接 import `modules/wallet/subject` vs 抽中性端口。无论哪条，不得要求 `admin.wallet` HTTP 已启。 | 方案冻结 | C1 | 用户裁决 | **verified** | — | 2026-09-03 用户裁决：直接复用 subject.Store，底层仅依赖 TxRunner，不依赖 admin.wallet HTTP（D-001） |

## 父目标

- `GOAL-001-telegram-channel-runtime`（Root · 纲领 R2）
