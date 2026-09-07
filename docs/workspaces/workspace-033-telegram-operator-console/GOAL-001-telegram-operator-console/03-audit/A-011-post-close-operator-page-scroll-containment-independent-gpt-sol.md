---
doc_type: goal-audit
id: A-011-post-close-operator-page-scroll-containment-independent-gpt-sol
parent: GOAL-001-telegram-operator-console
date: 2026-09-05
source: independent
auditor: subagent (gpt-5.6-sol · reasoning medium)
audit_type: post-remediation-independent
scope: 代码 checkpoint 9e9102cb 的 Telegram operator 页面级滚动隔离、刷新稳定性与前次浏览器测量建议复核
verdict: pass
open_required: 0
open_recommended: 0
version: 0.1.0
---

# A-011 · 页面级滚动隔离 independent audit（2026-09-05）

## 独立结论

一次性 `subagent (gpt-5.6-sol · reasoning medium)` 对当前代码只读核对，未修改文件、未改变目标
状态、未调用 Grok。结论为 `verdict: pass`、`open_required: 0`、`open_recommended: 0`，没有发现
新的 required 或 recommended finding。

## 核验范围与证据

- `apps/web/src/app/App.tsx:698-729,899-904,1033-1064` 对 `telegram-operator` 建立独立
  route-specific contained shell：固定 `100dvh` 根 shell、`min-h-0` flex body、`main` 的
  `overflow-y-hidden`，以及 operator page region 的 `h-full` 高度链；普通页面仍使用
  `data-shell-scroll-mode="page"` 与 `overflow-y-auto`。
- `apps/web/src/components/telegram-admin-tab.tsx:810-938` 的 operator surface、grid、transcript
  和 message list 具有可收缩的 `min-h-0`/flex/overflow 链；sessions 与 message list 各自
  `overflow-y-auto`，composer 为 `shrink-0`。
- 同文件 `:296-311,366-395,404-445,875-885` 保留后台刷新期间已有 timeline/session 内容，
  仅在真实会话切换时清空旧 timeline，并通过 in-flight 合并、可见性暂停和轻量 refreshing 状态
  避免周期性闪烁与竞态覆盖。
- `apps/api/modules/channel/telegram/schema/telegram-operator.json:11-17` 直接使用
  `type: "custom"`、`component: "telegram-admin-tab"`、`props.surface: "operator"`；
  `schema_test.go:45-59` 对该形状作精确断言。
- 定向 Vitest 共 27 个测试通过，Telegram schema Go test 通过；真实 Chromium 命令
  `$env:APP_PROFILE='custom'; npm run test:e2e -- telegram-operator-layout.spec.ts` 返回 2 个测试
  通过。E2E 使用 80 个 sessions、120 条长消息验证 document/body/root/main 不溢出、sessions/message
  list 内部滚动、composer 仍在 operator 内，并验证内层滚动不改变 window/document scroll；另以
  100 个 Users 验证普通页面仍可 page-level 滚动。

## 结论

本次变更满足“页面自身不因 operator 数据量增长而出现滚动条”的新增要求，同时保留“轮询刷新不
闪烁”“会话列表/消息区独立滚动”和“普通页面仍可滚动”的既有行为。A-009 的 R-001 浏览器
计算布局证据边界已被真实 Chromium E2E 覆盖；本审计不保留开放 recommended finding。
