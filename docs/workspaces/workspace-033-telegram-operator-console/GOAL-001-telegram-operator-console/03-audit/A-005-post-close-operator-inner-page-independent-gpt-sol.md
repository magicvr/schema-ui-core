---
doc_type: goal-audit
id: A-005-post-close-operator-inner-page-independent-gpt-sol
parent: GOAL-001-telegram-operator-console
date: 2026-09-05
source: independent
auditor: subagent (gpt-5.6-sol · reasoning medium)
audit_type: post-close-refinement-independent
scope: 关门后 Telegram 设置页/人工会话内页分离（修复前工作树）
verdict: conditional
open_required: 1
version: 0.1.0
---

# A-005 · 关门后人工会话内页分离 independent audit（2026-09-05）

## 独立结论

本条由一次性 `subagent (gpt-5.6-sol · reasoning medium)` 对修复前工作树独立执行；
不直接沿用主线程的实现判断，也未修改文件、未调用 Grok。结论为
`verdict: conditional`、`open_required: 1`。

## Findings

### F-001 · required · 设置页泄漏运行态计数

- **证据（审计时）**：`apps/web/src/components/telegram-admin-tab.tsx:797-801`。
- `captured_messages_count` 当时位于 settings surface 的配置渲染块之外，因此 Telegram
  设置页在加载状态后仍会显示运行态统计。
- 这违反“通道页只提供配置与人工会话入口，实际 chat 只在内页”的边界，属于 required
  finding；合法处理路径为移除或限制到 `surface: "operator"`，不能保留为 residual。

### F-002 · recommended · 契约与负向 surface 断言不足

- **证据（审计时）**：`apps/api/modules/channel/telegram/provider_test.go:186-192`；
  `apps/web/src/components/telegram-admin-tab.test.tsx:114-138`。
- 审计建议 provider 测试明确验证 settings/operator 的 Resources、Actions、ModuleID、
  Key、Owner，并补充设置页不渲染运行态/人工内容的负向断言。
- 本项为 recommended，不单独阻断，但与 F-001 同属本次边界回归的验证缺口。

## 其它范围

manifest/sidebar、operator route/schema、provider 两页 registration、App breadcrumb 与
已有 operator chat 行为未发现其它 required finding。A-005 不修改 Root status/progress，
等待编排器响应与修复后 independent re-audit。
