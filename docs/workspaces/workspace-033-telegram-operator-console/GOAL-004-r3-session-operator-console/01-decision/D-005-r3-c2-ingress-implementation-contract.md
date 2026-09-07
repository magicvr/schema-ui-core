---
doc_type: goal-decision
id: D-005-r3-c2-ingress-implementation-contract
parent: GOAL-004-r3-session-operator-console
date: 2026-09-04
source: Codex govern
status: done
version: 0.2.0
---

# D-005 · R3 C2 入站落盘实施合同

## 触发与边界

本合同把用户在 D-004 中已裁决的三项方向细化为可实现、可测试的 C2 参数：双表最小面、规范化字段、不建立 inbound dispatch 状态机，并保持现有 handler 错误语义。它不新增产品方案，不进入 C3 出站状态机、人工发送、列表权限或 UI。

## 持久化对象

### 1. `telegram_sessions`

- 一个运行时确认的 `bot_id` 与一个 Telegram `chat_id` 构成会话边界；复合主键为 `(bot_id, chat_id)`。
- 保存规范化会话资料：`chat_type`、`title`、`username`、最近入站活动时间以及创建/更新时间。Telegram ID 与 Unix 时间使用方言对应的 SQLite `INTEGER` / PostgreSQL `BIGINT`。
- 会话 upsert 与入站收据写入处于同一个 `kernel.Store` 事务；C3 按 `last_message_at` 与 `chat_id` 的稳定降序索引读取，不在 C2 预置未决的未读语义。

### 2. `telegram_inbound_messages`

- 复合主键为 `(bot_id, update_id)`，这是 webhook 重试与 polling 重复投递共享的唯一幂等边界。
- 保存规范化入站字段：固定 `direction=inbound`、`message_kind`（`text` / `command` / `callback`）、`chat_id`、可空 `user_id`、可空 `message_id`、可空 `callback_query_id`、可空 `text`、可空 `callback_data`、发送者 username 与接收时间；不保存 Telegram raw JSON、媒体/文件/贴纸或历史回灌内容。
- `text` 消息按既有 Dispatcher 规则区分 `command` 与普通 `text`。普通 `text` 是 C2 成绩单对象；`command` 与 `callback` 只作为已规范化的入站收据和一次性分发依据，不在 C2 成绩单中展示。
- 通过 `(bot_id, chat_id, received_at, update_id)` 索引支持后续成绩单读取；不增加 outbound 表或发送状态字段。

## 共同入站路径与顺序

1. webhook 与 polling 都继续进入同一个内部 `UpdatePayload` 路径；不扩张 kernel `TelegramUpdate`。先完成支持范围分类：文本消息（含命令）和 callback 记录；媒体、空 payload 或缺少会话 chat 的不支持更新不进入 C2 成绩单。
2. 先沿用当前 IP/chat/user 限流边界。被限流的 webhook 保持 `429 + Retry-After`；被限流的 polling update 保持当前“拒绝并跳过、不自动重试”的兼容语义，由 connection manager 在明确识别该拒绝后推进一次 offset，不能把它冒充为已持久化。
3. 对支持的更新，使用运行时经 `getMe` 确认的稳定 `bot_id`；bot identity 不可用时视为持久化错误，不能用 token、username 或零值替代幂等范围。
4. 在唯一收据事务之前，先沿用既有 `subjectStore.GetOrCreateSubject`（有 user identity 时调用）。该 repository 自己使用一个独立的 `Store.Run`，不得嵌套进 inbox 事务；主体映射错误属于可重试持久化错误，不能先铸造唯一收据。重复 update 也必须先完成这项既有可重试预分发工作，不能用重复短路径吞掉后续 Dispatcher 语义。
5. 在同一个 `Store` 事务中以方言无关的 `INSERT ... ON CONFLICT DO NOTHING`（占位符保持 `?`）尝试插入 `(bot_id, update_id)` 入站收据；必须读取 `RowsAffected()`：`0` 表示已成功接受的重复投递，同一事务不得再执行会话 upsert，调用方不得分发；`1` 才 upsert 对应会话。禁止在 PostgreSQL 唯一冲突后于同一 Tx 继续查询或写入。新收据路径任一写入失败都回滚整个事务并返回错误。
6. 新收据事务提交成功后才调用现有 Dispatcher。Dispatcher handler 错误沿用当前告警/日志且不自动重试；因为收据已经成功，webhook 仍可成功确认，polling 才可推进 offset。重复收据直接作为幂等成功返回，不再次调用 Dispatcher。
7. webhook 只有共同路径返回 nil（新收据成功、重复收据成功或明确不支持更新）才返回 2xx；落盘或主体映射错误返回 5xx。polling 只有共同路径成功后才推进成功接受更新的 offset；落盘或主体映射错误保持当前 offset，供同一 `update_id` 重试。限流拒绝是第 2 步声明的唯一非持久化跳过路径。

## 实现与验证约束

- Telegram channel module 自有 v68 migration、repository 与测试；repository 只依赖已有方言无关 `kernel.Store` / `TxRunner`，不改 kernel Telegram port。
- 必须覆盖 SQLite migration 与 restart、PostgreSQL DDL 形状和 gated runtime duplicate path、webhook/polling 共同路径、首次写入、`ON CONFLICT DO NOTHING` 重复投递、并发唯一竞争、主体映射失败后仍可重试且不先铸造 inbox、事务失败不确认/不推进 offset、handler 错误不自动重试，以及命令/回调收据不重复分发。
- 本合同只冻结 C2 实施语义；完成代码、迁移和测试后必须另行 self + Grok independent 审计，不能以本合同或自审代替运行时证据。

## 未选与禁止扩展

- 不采用单一方向事件表，不在 C2 建立 outbound 表、`pending/sent/failed` 状态机、后台重试、人工发送或 kernel 新端口。
- 不以进程内 map/mutex 代替数据库唯一约束；不把 raw JSON 作为主存储或兜底审计资料。
