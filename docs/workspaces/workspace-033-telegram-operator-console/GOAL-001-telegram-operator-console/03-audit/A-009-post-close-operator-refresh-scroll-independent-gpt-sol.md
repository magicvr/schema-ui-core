---
doc_type: goal-audit
id: A-009-post-close-operator-refresh-scroll-independent-gpt-sol
parent: GOAL-001-telegram-operator-console
date: 2026-09-05
source: independent
auditor: subagent (gpt-5.6-sol · reasoning medium)
audit_type: post-remediation-independent
scope: 代码 checkpoint dc7ac5e5 的轮询刷新稳定性与 Telegram 人工会话滚动布局
verdict: pass
open_required: 0
open_recommended: 1
version: 0.1.0
---

# A-009 · 轮询刷新与滚动布局 independent audit（2026-09-05）

## 独立结论

一次性 `subagent (gpt-5.6-sol · reasoning medium)` 对当前 checkpoint 只读核对，未修改文件、
未调用 Grok。结论为 `verdict: pass`、`open_required: 0`，没有发现新的 required finding。

## 核验范围与证据

- `apps/web/src/components/telegram-admin-tab.tsx:404-425` 的 operator refresh 仍通过
  `operatorRefreshRef` 合并同一时刻的刷新；`:289-324` 的 `loadTimeline` 保留已有 timeline，
  只在当前 chat 仍匹配时替换数据。
- `apps/web/src/components/telegram-admin-tab.tsx:871-881` 在已有消息刷新时保留 message list，
  只显示 `timelineRefreshing` 状态；初次无消息时仍显示 loading，未牺牲首次加载反馈。
- `apps/web/src/components/telegram-admin-tab.tsx:845-860` 与 `:303-314,349-354` 在切换真实
  chat 时清空旧 timeline，并以当前 chat 守卫阻止旧响应写入新会话。
- `apps/web/src/components/telegram-admin-tab.tsx:807,840-842,871,881,930-958`：operator
  面板具有 viewport max-height、`min-h-0`、flex/grid 剩余空间分配；sessions nav 与 message
  list 具有 `overflow-y-auto`/`overflow-x-hidden`，transcript 与 message list 可收缩，composer
  为 `shrink-0`，长文本/时间戳/textarea 有横向与高度约束。
- `apps/web/src/components/telegram-admin-tab.test.tsx:173-194,517-594` 覆盖布局契约与后台
  refresh pending 时保留已有消息、刷新完成后更新消息；既有 lease、visibility、请求合并、
  capability race、send/retry 测试继续通过。
- 定向测试通过：`telegram-admin-tab.test.tsx` 16 tests、相关 catalog 一并执行共 28 tests；
  本审计范围内未见失败或 whitespace error。

## 推荐项与边界

### R-001 · recommended · 缺少浏览器计算布局证据

当前测试主要断言 class 字符串，没有在真实浏览器测量 `clientHeight`、`scrollHeight`、
`scrollWidth` 或验证特定 viewport 下 composer 的像素级可见性。该缺口不构成 required finding，
因为本次目标的代码契约、构建产物和行为测试已经覆盖核心实现，且未要求浏览器 E2E 作为门禁。

## 结论

当前变更满足“后台刷新不造成消息区闪烁”和“会话列表/消息区受控滚动”的代码级要求；建议后续
若建立浏览器 E2E 体系，再补 R-001 的实际计算布局验证。独立意见不修改 Root status/progress。
