---
doc_type: goal-audit
id: A-001-r2-self-audit
parent: GOAL-003-r2-webhook-dispatch-identity
date: 2026-09-03
source: self
scope: R2 Webhook 路由、Update 分发、主体映射与入站限流（C1～C3 全量范围）
audit_type: stage-closeout
verdict: pass
open_required: 0
---

# A-001 · R2 Webhook 路由、Update 分发、主体映射与入站限流自审（self）

## 1. 审计基本信息

- **目标**：[GOAL-003-r2-webhook-dispatch-identity](../00-meta.md)
- **审视范围**：
  - C1 用户裁决（D-001：I-030-007 直接复用 `subject.Store` + `internal/modules` 代码布局）。
  - C2 代码实施（`internal/channel/telegram/`、`modules/channel/telegram/`）与单元/全量测试（E-002）。
  - 退出判据对照：判据 #1（Webhook 合同 secret fail-closed）、判据 #2（分发端口 Register/未知回落）、判据 #4（身份映射 `issuer=telegram` 幂等且不写 `admin.users`）。
  - 限流映射（判据 #6 落地）：IP 60/m, Chat 30/m, User 20/m 三桶请求计数与 Record 永不 Clear。
- **审计模式**：`cross`（自审 A-001 + 本地 grok build 独立审计 A-002）。
- **结论**：**PASS**（开放必改 0，建议 0）。

## 2. 判据与合同对照

| 检查项 | 契约与标准 | 实际落地 | 判定 |
|--------|------------|----------|------|
| Webhook 路由与 Secret | `POST /api/channel/telegram/webhook`，`X-Telegram-Bot-Api-Secret-Token` 常时比较 fail-closed（401），无 Token 503 | `webhook.go` 严格按顺序实现；`subtle.ConstantTimeCompare`；测试全覆盖 | PASS |
| Update 分发 | 命令剥离 `/` 与 `@BotName` 精确匹配；未注册命令回落 `DefaultTelegramUnknownCommandText`；Callback 精确匹配 | `dispatcher.go` 实现线程安全分发与回落；测试全覆盖 | PASS |
| 主体映射 | `GetOrCreateSubject("telegram", userID)` 幂等创建；不写 `admin.users`；不依赖 `admin.wallet` HTTP | 复用 `modules/wallet/subject.Store`，纯 TxRunner 依赖；测试验证幂等性 | PASS |
| 入站限流 | IP（60/m）、Chat（30/m）、User（20/m）；超限 429 + Retry-After；IP 在 Secret 失败时仍 Record | 消费 `kernel.RateLimiterProvider`；三桶独立实例；测试覆盖超限与记账 | PASS |
| 模块边界 | ModuleID = `channel.telegram`，豁免业务导航，不进默认 Profile | `modules/channel/telegram/provider.go` 干净提供 Route | PASS |

## 3. Findings

- **Required findings**: 0
- **Recommended findings**: 0

## 4. 下一步

自审通过。进入独立交叉审计（调用本地 grok build 执行 `/audit`）。
