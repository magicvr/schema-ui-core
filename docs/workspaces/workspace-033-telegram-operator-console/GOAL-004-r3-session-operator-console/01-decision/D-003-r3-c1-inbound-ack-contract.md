---
doc_type: goal-decision
id: D-003-r3-c1-inbound-ack-contract
parent: GOAL-004-r3-session-operator-console
date: 2026-09-04
source: Codex govern
status: done
version: 0.1.0
---

# D-003 · R3 C1 入站持久化与确认顺序补充合同

## 触发与边界

A-003（Grok independent）指出，D-002 已选定 bot 维度 `update_id` 幂等，但尚未把持久化成功与 webhook 确认 / polling offset 推进的先后关系写成可实施合同。该记录是对用户已选 I-033-020 的治理补全，不改变 D-002 的七项选择，也不扩大 VP-033 首波边界。

## 已冻结合同

1. 仅符合 VP-033 首波范围的 Telegram 实际投递文本进入会话/消息持久化；非命令文本及 bot 可见的群消息仍是本 VP 的入站对象，历史回灌、媒体和领域事件总线不在范围内。
2. 入站持久化必须先成功：会话/消息及 `message_id` 等元数据的持久化事务成功后，webhook 才能返回成功确认，polling 才能推进对应 `offset`。已存在的同一 `(bot, update_id)` 重复投递可作为幂等成功处理，但不得再次落盘或分发。
3. 持久化失败必须返回可重试错误；webhook 不返回成功确认，polling 不推进该 update 的 offset，并在重试时继续使用同一 bot-scoped `update_id` 幂等边界。不得把“已向 Telegram 确认”置于持久化之前。
4. `(bot_id, update_id)` 是入站唯一幂等键。webhook 重试与 polling 重复投递共享该键；重复 update 不重复持久化、不重复分发。C2 必须将 offset 推进放在 handler/持久化成功之后，不能沿用当前 polling 代码的先推进时序。
5. C2 可在内部 `UpdatePayload` 上实现该合同，不为幂等扩张已交付 VP-030 kernel `TelegramUpdate` 端口；当前代码尚未因本记录而被宣称已实现。

## 门禁结论

本记录闭合 A-003 F-001 的合同缺口，闭合路径为 `fixed`；A-004 self 响应已记录。C1 仍须等待本记录后的 Grok independent re-audit，确认后才可放行 C2；该记录不把尚未实现的代码写成成功事实。
