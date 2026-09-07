---
doc_type: goal-audit
id: A-013-post-close-operator-im-chat-independent-gpt-sol
parent: GOAL-001-telegram-operator-console
date: 2026-09-05
source: independent
auditor: subagent (gpt-5.6-sol · reasoning medium)
audit_type: post-remediation-independent
scope: 代码 checkpoint 6ccef765 的 IM 消息排序、滚动、发送者标签与 composer 快捷键
verdict: conditional
open_required: 1
open_recommended: 2
version: 0.1.0
---

# A-013 · IM 聊天交互 independent audit（2026-09-05）

## 独立结论

一次性 `subagent (gpt-5.6-sol · reasoning medium)` 对 checkpoint `6ccef765` 只读核对，未修改
文件、未改变目标状态、未调用 Grok。结论为 `verdict: conditional`，`open_required: 1`、
`open_recommended: 2`。

## 核验范围与证据

- `apps/web/src/components/telegram-admin-tab.tsx:91-100,315-330` 将正常时间戳的 newest-first
  timeline 转为 oldest-to-newest；`useLayoutEffect` 与 `timelineStickToBottomRef` 在初始加载、
  切换会话及用户停留底部时跟随，用户上翻超过阈值后保留滚动位置：`:139-145,469-489,917-924`。
- 同文件 `:947-956` 以 `justify-end/bg-primary` 与 `justify-start/bg-background` 区分出站/入站
  气泡；`:938-943` 的 sender label 优先级为选中会话 `title` → `item.senderUsername` →
  选中会话 `username` → fallback，因此私聊可显示昵称，出站显示 `Bot`。
- 同文件 `:996-1008` 实现默认 Enter 发送、Ctrl+Enter 换行及 checkbox 反转；发送流程和
  capability/disabled 门禁沿用既有 `sendMessage()`：`:504-541,990-993,1040-1047`。
- `apps/web/src/components/telegram-admin-tab.test.tsx:275-390` 覆盖正序、左右气泡、sender 标签、
  自动跟随、上翻保持位置、默认/反转快捷键；`apps/web/e2e/telegram-operator-layout.spec.ts:91-169`
  覆盖真实 Chromium 的排序、标签、初始底部和页面/内层滚动。

## Findings

### F-001 · required · open · 群组/频道入站 sender label 可能显示 chat 标题

- `telegram-admin-tab.tsx:938-943` 将 `session.title` 放在 `item.senderUsername` 之前。
- 后端将 `Session.Title` 与 `TimelineEntry.SenderUsername` 分开建模：
  `apps/api/modules/channel/telegram/store/repository.go:61-89`；因此群组/频道会话的 title 可能是
  chat 标题，而不是当前入站消息发言人的 nickname/username。
- 影响：广义 operator 会话中，入站消息可能显示群名而非当前发言人，不满足用户要求的 sender
  语义。需要修正实现，或由用户书面限定为 private chat；本意见不替用户作范围裁决。

### F-002 · recommended · open · 缺少真实 Chromium composer 快捷键证据

`telegram-operator-layout.spec.ts:79-170` 未实际触发 Enter、Ctrl+Enter 或 checkbox 反转；目前
仅有组件级测试覆盖。

### F-003 · recommended · open · 缺少真实 Chromium 上翻后刷新保持位置证据

组件测试验证 `scrollTop` 保持，但 E2E 只验证内层手动滚动不改变页面级滚动，未验证上翻后刷新
不会抢回底部。

## 结论

实现的正序展示、底部跟随、气泡样式和快捷键代码路径基本成立，但在修正 F-001 并补足浏览器
行为证据前，不应将本轮交互改造视为无条件收束。意见不修改 Root/workspace 状态；响应由
`/govern` 处理。
