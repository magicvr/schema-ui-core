---
doc_type: goal-audit
id: A-003-r2-closure-response
parent: GOAL-003-r2-webhook-dispatch-identity
date: 2026-09-03
source: self
scope: R2 Webhook 路由、Update 分发、主体映射与入站限流（A-002 审计响应与关门）
audit_type: stage-closeout
verdict: pass
open_required: 0
---

# A-003 · R2 独立审计意见响应与必改闭合（合并响应）

## 1. 响应背景

编排器响应独立交叉审计意见 [A-002-r2-independent-audit.md](A-002-r2-independent-audit.md)（grok-build grok-4.6 · reasoning high · `verdict: fail`，1 required F-001，4 recommended R-001～R-004）。

## 2. 意见响应与闭合台账

| 发现项 | 严重度 | 闭合路径 | 闭合依据与实施事实 | 状态 |
|--------|--------|----------|-------------------|------|
| **F-001** | high / required | **fixed** | 1. `apps/api/kernel/profile.go`：将 `channel.telegram` 编入 `BuiltinModules()` 候选集，声明 `DependsOn: ["core.server-registration"]`，`Requires: [CapabilityHTTP]`，未进入任何默认 Profile。<br>2. `apps/api/internal/composition/composition.go`：在 `newServer` 装配时对 `plan.HasModule("channel.telegram")` 进行全套运行时装配，基于底层 `st` 独立构造 `subject.NewStore(st)`（不依赖 `admin.wallet` HTTP），构造 WebhookHandler 并挂入 providers。<br>3. `apps/api/internal/composition/composition_telegram_test.go`：新增 `TestTelegramChannelComposition` 验证注册、Profile 隔离与路由执行。 | **closed** |
| **R-001** | med / recommended | **fixed** | `webhook.go` 中，主体映射如果因底层存储报错失败，立即返回 **500 Internal Server Error**，fail-closed 触发 Telegram 重试，禁止携带空 SubjectID 投递。 | **closed** |
| **R-002** | med / recommended | **fixed** | `webhook_test.go` 新增 `TestWebhook_RateLimiting_ChatBucket`（测试 30/min 上限与 429 Retry-After）及 `TestWebhook_SubjectMappingIdempotency`。 | **closed** |
| **R-003** | low / recommended | **fixed** | `dispatcher_test.go` 新增 `TestDispatcher_InvalidRegistrations`，覆盖 nil handler、空命令、非法命令、空 callback 及超长 callback 边界错误断言。 | **closed** |
| **R-004** | med / recommended | **fixed** | `modules/channel/telegram/provider.go` 中 Descriptor 剔除 Admin 六面套用，精简为横切通道依赖。 | **closed** |

## 3. 关门判定

- 开放 required findings：**0**（F-001 已合法闭合：`fixed`）。
- 退出判据 #1（Webhook 合同与 secret fail-closed）、#2（分发端口与未知回落）、#4（主体映射幂等与隔离）及三桶限流映射已完整落地并通过端到端装配测试。
- 全仓回归测试全部 PASS。
- GOAL-003 检查点 C1、C2、C3 全部完成，达到关门条件。
