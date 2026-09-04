---
doc_type: goal-audit
id: A-029-r3-c4-ui-foundation-self
parent: GOAL-004-r3-session-operator-console
date: 2026-09-05
source: self
auditor: Codex
audit_type: implementation-slice
scope: R3 C4 UI 基础切片与 I-033-009 非阻断刷新行为；不含 getChatMember capability API 方案或 C4 关门
verdict: conditional
open_required: 0
version: 0.1.0
---

# A-029 · R3 C4 UI 基础切片 self 审视（2026-09-05）

## 结论

本条只审视已实现的 C4 基础切片，不作为 C4 或 R3 的关门意见。现有 Telegram
Admin 页已接入 C3 operator 的会话列表与文本成绩单，并实现 10 秒单飞、页面隐藏
暂停、恢复可见时立即刷新；在 capability 尚未确认前，composer 与 retry 均保持
fail-closed 禁用。

## 已核对事实

- `apps/web/src/components/telegram-admin-tab.tsx` 复用现有认证 `fetcher`，按
  C3 已冻结的 `{items,total,page,pageSize}` 合同读取会话和成绩单，保持 chat ID
  为字符串，按选中 chat 加载文本时间线。
- operator refresh 以 promise flight 去重；页面隐藏时清理 10 秒 timer，恢复
  可见时立即启动新的 refresh；同一 chat 的成绩单请求也不会重复并发。
- 未获得 capability 结果时，发送 composer 和失败消息 retry 控件均禁用；本条
  没有假设发言权，也没有新增 Bot API 路由或持久化字段。
- `telegram-admin-tab.test.tsx` 的定向测试 8/8 通过；本次 Web 全量测试为
  92 个测试文件、1203 个测试通过。
- `npm run build` 仍被既有 `src/renderer/form-controls.tsx:946-947` 的
  `number | undefined` 到 `string | undefined` 类型错误阻断；该错误不在本条
  写集。prebuild 改写的三个 conformance 生成物已恢复，未纳入本条。

## 未决边界

`I-033-023` 记录的 capability API 形状与缓存所有权仍待用户裁决。A-029 不选择
独立 capability 路由、成绩单附带 capability 或会话列表附带 capability；因此
本条为切片级 `conditional`，不放行 C4 关门，也不把 composer 视为已交付。

## 后续门禁

用户裁决后，C4 仍需实现 `getChatMember`/`can_send`、60 秒 bot/chat 缓存、403
失效和显式重探，接通发送/失败/retry 状态，并进行 self + Grok independent
验证后才能关闭 C4。
