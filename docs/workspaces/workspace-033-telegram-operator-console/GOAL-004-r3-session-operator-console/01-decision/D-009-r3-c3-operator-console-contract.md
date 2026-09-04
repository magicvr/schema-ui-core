---
doc_type: goal-decision
id: D-009-r3-c3-operator-console-contract
parent: GOAL-004-r3-session-operator-console
date: 2026-09-04
source: self
status: done
version: 0.2.0
---

# D-009 · R3 C3 人工台 API、权限与运行时实施合同

## 合同范围

本合同把 D-002 的专用权限、未绑定运行时边界和同步 sender 选择，与 D-008 的
显式重试身份细化为可测试的 C3 API 合同。C3 只交付会话列表、统一成绩单、人工
文本发送、失败重试、权限贡献和运行时接线；页面、菜单、`getChatMember` TTL
缓存、composer 反馈和端到端 UI 留在 C4。合同不引入历史回灌、媒体、FSM、群发、
频道、多 bot、多实例 polling、独立进程或 SSE/WebSocket。

## 运行时与权限门禁

- 所有 operator 路由先经过认证和专用权限检查，再检查 Telegram 运行时；这些路由
  在 composition 中必须像 settings/lease 一样由 `a.Middleware` 预包装后才注入
  Provider；`Public: false` 只是路由声明，不是认证实现。匿名请求
  返回 `401 UNAUTHENTICATED`，缺权限返回 `403 FORBIDDEN`，不得通过错误差异泄漏
  运行时或会话信息。
- 读取权限为 `telegram.operator.read`，默认策略为
  `system.admin-editor-viewer`；发送和重试权限为 `telegram.operator.write`，默认
  策略为 `system.admin`。两者均由 `channel.telegram` Provider 贡献，独立于
  `settings.read/write`，并进入现有 system-data reconcile。
- 只有运行时状态为 `running`、`bot_id > 0`、receiver 为 `webhook` 或 `polling`，
  且 `Dispatcher.HasBusinessHandlers()` 为 false 时，operator surface 才可用。
  状态由服务端从 `RuntimeManager.ConnectionStatus()` 和 dispatcher 派生；客户端
  不得提交或覆盖 `bot_id`。不满足时返回 `409 TELEGRAM_OPERATOR_UNAVAILABLE`，
  不调用 sender。
- 未绑定 polling 的 lease acquire/heartbeat/release 保持既有三条路径，但其鉴权
  接受 `settings.read` 或 `telegram.operator.read`；这让专用 operator reader 能
  将 `idle/none` 的 receiver 按既有心跳需求推到 `running`。settings API 仍只受
  `settings.read/write` 控制，operator API 不自行创建 lease。
- API 只从已确认的当前 bot identity 查询数据；不会把 token、secret、Bot API URL
  或完整下游错误写入响应或持久化记录。

## HTTP surface

Provider 新增以下四条非 Public operator 路由；composition 必须把 operator handler
用 `a.Middleware` 包装后再传入 Provider；现有 Telegram settings 的 page、
navigation 与 fragment 贡献保持不变，C3 不新增 operator page、navigation 或
fragment：

| 方法与路径 | 权限 | 用途 |
|------------|------|------|
| `GET /api/channel/telegram/operator/sessions` | `telegram.operator.read` | 列出有入站活动的会话 |
| `GET /api/channel/telegram/operator/sessions/{chat_id}/messages` | `telegram.operator.read` | 读取该会话的统一文本成绩单 |
| `POST /api/channel/telegram/operator/sessions/{chat_id}/messages` | `telegram.operator.write` | 建立一次新的人工发送尝试 |
| `POST /api/channel/telegram/operator/sessions/{chat_id}/messages/{request_id}/retry` | `telegram.operator.write` | 对失败尝试发起显式重试 |

`chat_id` 和 Telegram 的大整数标识在 JSON 中均以十进制字符串返回，避免浏览器
数字精度损失。列表和成绩单统一使用 `{items,total,page,pageSize}` 响应，默认
`page=1,pageSize=20`，`pageSize` 最大 100；非法或非正参数返回 400，SQL offset
必须使用已有的溢出安全分页辅助。会话按 `last_message_at DESC, chat_id DESC`
稳定排序，成绩单按 `occurred_at DESC` 再按来源标识稳定排序，均不接受客户端
覆盖 bot 范围或排序列。

所有时间字段（`lastMessageAt`、`occurredAt` 及 outbound 的状态时间，如有返回）
统一为 UTC RFC3339 字符串；所有 Telegram 数字标识（`chatId`、`userId`、
`messageId`、`updateId`）均以十进制字符串返回。

会话条目返回 `chatId`、`chatType`、`title`、`username`、`lastMessageAt`；仅有
入站文本/命令活动的 v68 session 可见。成绩单是入站文本/命令与出站人工文本的
统一时间线：入站条目带 `direction=inbound`、`status=received`、`updateId`、
`messageId`、`userId`/`senderUsername`（有值时）和 `text`；出站条目带
`direction=outbound`、`status=pending|sent|failed`、`requestId`、可空的
`retryOf` 和 `text`。callback、空文本和未建模媒体不进入成绩单；失败原因只能
是短的、经过脱敏和截断的可诊断文本。

## 发送与重试状态机

- 新发送 body 为 `{"requestId":"...","text":"..."}`。服务端校验 chat 已在
  当前 bot 的 session 中，`requestId` 必须是 1～128 字节的安全可打印标识，文本
  非空且不超过 4096 个 UTF-8 字节；不接受客户端的 status、bot_id、retry_of 或
  任意 Telegram markup。
- 新尝试先在事务中写入 `pending` outbound 记录并提交，成功后才调用现有
  `kernel.TelegramSender`；发送成功更新为 `sent`，sender 返回错误更新为
  `failed`，并返回可识别的失败记录。pending 写入失败时绝不调用 sender。状态更新
  失败时不得伪装成 `sent`，同一 request_id 也不得因此再次外发。
- outbound 的 request-id 冲突和 retry-root pending 冲突都必须在事务内以
  `INSERT ... ON CONFLICT DO NOTHING`（占位符 `?`）处理；随后以
  `RowsAffected()==0` 判断冲突，再在同一事务中读取并比较 request payload，或读取
  root 的 pending 状态映射为 `TELEGRAM_REQUEST_IN_PROGRESS`。禁止先触发 PostgreSQL
  唯一异常、再在同一已中止事务中 SELECT；不得把 `kernel.IsUniqueViolation` 当作
  幂等接受结果。
- `(bot_id, request_id)` 是幂等键。相同 request_id 且 payload（chat、text、
  retry_of）一致时返回已有记录：terminal 状态不重复外发，pending 状态返回
  `409 TELEGRAM_REQUEST_IN_PROGRESS`；payload 不一致返回 `409
  TELEGRAM_REQUEST_CONFLICT`。客户端在网络重试时必须复用同一 request_id。
- 每个 outbound 记录保存 `retry_root`；首发的 root 是自身且 `retry_of` 为空。重试
  路由只接受 failed 尝试，服务端为其创建新的 request_id 和新的 `pending` 记录，
  `retry_of` 固定指向该发送链的原始 root；不原地改写失败行。若该 root 已有
  pending 尝试，返回 `409 TELEGRAM_REQUEST_IN_PROGRESS`；若任一尝试已 sent，
  返回 `409 TELEGRAM_RETRY_NOT_ALLOWED`。不同 request_id 代表不同的显式尝试，
  但同一 root 同时最多一个 pending，数据库约束与事务共同保证该边界。
- 重试 body 只接受新的 `requestId`；正文和 chat 从失败记录服务端派生。新 request_id
  必须遵守同一幂等规则；无后台 worker、定时器或自动重试。所有 pending/failed
  记录都可在成绩单中核对，不能用“未落盘”隐藏失败。

## 持久化与验证分母

C3 新增 v69 `telegram_outbound_messages`，SQLite 与 PostgreSQL 具有相同语义：
`bot_id`、`request_id`、`retry_root`、可空 `retry_of`、`chat_id`、纯文本、
`status`、脱敏错误、`created_at`/`updated_at`；唯一键为 `(bot_id,request_id)`，
并对同一 `(bot_id,retry_root)` 的 `pending` 建立数据库级部分唯一索引（SQLite 与
PostgreSQL 均使用 `WHERE status = 'pending'`）。不得保存 Bot token、secret、raw
Telegram JSON 或完整 HTTP 请求。SQLite 使用 `INTEGER`、PostgreSQL 使用
`BIGINT`，两套 DDL 的状态约束和索引语义必须一致。

Provider `Descriptor()`、`kernel/profile.go` 与实际 `reg.HTTP`/`reg.Authorization`
贡献必须同步声明四条 operator 路由和两个权限键；`channel.telegram` 仍不进入
mvp/admin/demo 默认 Profile。错误码 `TELEGRAM_OPERATOR_UNAVAILABLE`、
`TELEGRAM_REQUEST_IN_PROGRESS`、`TELEGRAM_REQUEST_CONFLICT`、
`TELEGRAM_RETRY_NOT_ALLOWED` 及未知 chat/request 的稳定选择必须登记到冻结
error catalog，并有中英文条目；`request_id` 收紧为 mux-safe 的
`[A-Za-z0-9._-]{1,128}`。4096 字节正文上限是有意的实施参数，不宣称 Telegram
字符数上限。

如果 sender 成功但 `sent` 状态更新失败，API 返回 5xx、保留 `pending`、不允许
同 request 或 pending root 再外发；该 fail-closed 卡住状态是有意合同，后续只可
通过受治理的运维 reconcile 处理。调用 sender 前再次确认当前 bot identity/token
仍就绪；`HTTPSender` 在 token 缺失且没有明确的测试 CaptureSender 时返回的 nil
不得被当作外部发送成功。

C3 self 与 independent 验证至少覆盖：SQLite 与 gated PostgreSQL 的迁移/读写；
read/write 权限及默认角色 grant；运行时未配置、已绑定、非 running 的 fail-closed；
分页及 chat/bot scope；callback/空文本排除；pending 先于 sender；成功、失败、
持久化失败；同 request 并发和 payload 冲突；retry_of/root/并发 pending/无自动重试；
lease 的 `settings.read OR telegram.operator.read` 授权；匿名/无权限/服务凭据缺
scope 的 401/403 顺序；Descriptor/profile/catalog 同步；以及 `go test`、相关
`-race`、`git diff --check` 和显式未跟踪文件空白检查。PG 未配置时只能标注 gated
skip，不得把 skip 当作通过。
