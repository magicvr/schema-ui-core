---
doc_type: goal-execution
id: E-001-goal-opened
parent: GOAL-004-r3-outbound-settings-limiter
date: 2026-09-03
status: recorded
---

# E-001 · GOAL-004 目标建立与 R3 启动

## 1. 目标建立

- 子目标 [GOAL-004-r3-outbound-settings-limiter](../00-meta.md) 建立，承接 Root 纲领 R3 阶段。
- 核心范围：
  1. 出站生产适配器（基于 stdlib `net/http`，10s 超时，POST `https://api.telegram.org/bot<token>/sendMessage`，文本与 InlineButton，无 token 降级 mock）。
  2. Admin 设置面（I-030-005 热切换，密钥 fail-closed，脱敏防泄漏）。
  3. 入站三桶限流核账（核对 IP 60/m, Chat 30/m, User 20/m 均在 R2 落地，出站无阻断）。
- 启动 C1：向用户提出 I-030-005 热切换及 Admin 设置架构选型决策。
