---
doc_type: goal-decision
id: D-011-r3-c4-independent-capability-route
parent: GOAL-004-r3-session-operator-console
date: 2026-09-05
source: user
status: done
version: 0.1.0
---

# D-011 · R3 C4 独立 capability 路由裁决

## 用户已裁决

用户选择 **独立 capability 路由（推荐）**，作为 C4 `getChatMember` 发言权反馈
的唯一 HTTP 承载方式。实现按本合同细化，不混入成绩单附带或会话列表附带的
capability 载荷。

## API 合同

- 新增 `GET /api/channel/telegram/operator/sessions/{chat_id}/capability`，由
  `telegram.operator.read` 保护，并沿用 D-009 的认证、运行时、未绑定和
  `bot_id` 门禁。
- 成功响应为 `{ "chatId": "<十进制字符串>", "canSend": true|false }`；不返回
  token、secret、Bot API URL、原始 Telegram JSON、完整下游错误或缓存内部细节。
- `refresh=1` 是显式强制重探参数。缺省请求可命中缓存；其他非空 `refresh` 值
  返回 `400 INVALID_BODY`。成绩单的 10 秒轮询不得隐式触发 capability 重探。
- capability 只查询当前 bot identity 在目标 `chat_id` 的成员能力，不能由客户端
  提交或覆盖 `bot_id`/`user_id`。探测调用 Telegram `getChatMember`，`user_id` 固定
  为当前 bot id。

## 缓存、并发与状态映射

- 缓存 owner 是 channel.telegram 内的 server-side capability service；它消费
  composition 注入的 `kernel.Cache`，使用独立 namespace
  `telegram-operator-capability`，键按 `(bot_id, chat_id)` 固定，采用 60 秒
  absolute TTL。允许与拒绝结果均缓存。
- 同一 bot/chat 的并发探测使用 single-flight；强制重探绕过并替换已有缓存，但
  仍与同键在途请求合并。无后台刷新、自动重试或跨 bot/chat 复用。
- Telegram member 状态映射固定为：`creator`/`member` 默认允许；
  `administrator` 默认允许；若响应明确给出 `can_send_messages=false` 或
  `can_post_messages=false`，则拒绝；`restricted` 仅在明确
  `can_send_messages=true` 时允许；`left`、`kicked`、未知状态拒绝。缺字段不能
  推导为允许。
- `getChatMember` 的 Telegram 403 视为拒绝并可缓存为 `canSend=false`；其他
  非成功探测错误返回 cataloged `502 TELEGRAM_CAPABILITY_UNAVAILABLE`，不写入
  允许缓存，客户端保持 fail-closed。

## 真实发送的最终权威

- `HTTPSender` 的 Telegram HTTP 403 或 API `error_code=403` 必须结构化识别；真实
  `sendMessage` 403 经过 capability-invalidating sender 立即删除同一
  `(bot_id, chat_id)` 缓存键。非 403 发送错误不删除 capability 缓存。
- 发送/重试的 durable `pending` → `sent|failed`、`request_id` 幂等、显式新
  request 与 `retry_of` 仍完全遵守 D-008/D-009；403 仍是发送失败，不伪装为成功，
  不自动重试。
- Web UI 按当前选中 chat 保存 capability 状态；`unknown`、`denied`、`error`
  和请求中的状态都禁用 composer/retry。重新进入会话或手动刷新才带
  `refresh=1` 重探；普通 10 秒刷新只更新会话/成绩单。发送或重试返回失败时 UI
  立即 fail-closed 并刷新成绩单，不自动重探。

## 取舍与边界

- **未选：成绩单附带**。它会把权限探测与时间线读取耦合，使只读刷新频率影响
  Bot API 请求和缓存语义。
- **未选：会话列表附带**。它不能自然表达选中 chat 的显式重探，也会让列表刷新
  载荷承担逐会话 capability 扇出。
- 本裁决只关闭 I-033-023 的 API 形状与缓存 owner；不扩展历史回灌、FSM、群发、
  频道、多 bot、多实例 polling、独立进程或 SSE/WebSocket。

## 依据

D-002 已冻结混合 `getChatMember` 预检、真实发送结果最终权威、60 秒 bot/chat
缓存、403 失效和显式重探；D-008/D-009 已冻结 outbound 状态机和重试身份。本条
记录用户对 I-033-023 三种互斥 API 形状的书面选择。
