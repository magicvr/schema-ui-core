---
doc_type: goal-decision
id: D-004-r3-c2-user-decisions
parent: GOAL-004-r3-session-operator-console
date: 2026-09-04
source: user
status: done
version: 0.1.0
---

# D-004 · R3 C2 用户方案裁决

## 用户已裁决

| 信息项 | 裁决 | 影响 |
|--------|------|------|
| C2 持久化对象面 | **双表最小面**：建立会话表与入站消息表；出站状态表留给 C3 | C2 只承载 `chat_id` 会话与入站文本/幂等记录，不提前扩展出站状态机 |
| C2 原文留存 | **规范化字段**：不保存 Telegram raw JSON | 持久化文本与 chat/user/message/update 等必要元数据，减少 PII 保留和原文 schema 耦合 |
| C2 持久化后分发失败 | **兼容现有语义**：唯一 inbox 记录先落盘，重复 update 跳过分发；持久化失败可重试；handler 错误保持现有告警且不自动重试 | 保持当前 Dispatcher 的同步、无后台重试边界；不新增 inbound dispatch 状态机 |

## 裁决边界

- 会话主键继续承接 D-002 的 `chat_id` 选择；单 bot 边界下使用运行时已确认的稳定 bot identity 作为 bot-scoped 幂等范围。具体列类型、索引名和迁移版本属于实施参数，必须满足 `(bot_id, update_id)` 唯一约束。
- 入站消息表只保存 VP-033 首波允许的规范化文本和元数据；不把 raw JSON 作为权威或备用留存，不写媒体/文件/贴纸/历史回灌。
- “唯一 inbox 记录先落盘”不等于增加出站 outbox；C2 不建立 outbound 表、不实现人工发送、不改变 C3 的 `pending`/`sent`/`failed` 状态机。
- 当前 kernel `TelegramUpdate` 与 `TelegramSender` 薄端口保持不变；repository 通过现有方言无关 `kernel.Store` 事务边界接入。

## 未选方案

- 不选择单一方向字段事件表；不选择在 C2 一并建立出站表。
- 不选择带 raw JSON 的规范化表；不选择 raw JSON 为主的存储。
- 不选择新增 inbound dispatch 状态机或以内存锁替代数据库幂等。

## 后续门禁

C2 schema/迁移、共同入站接线与幂等测试仍须完成 self + Grok independent 审视后实施；本记录只冻结用户方案，不把代码或迁移写成已完成事实。
