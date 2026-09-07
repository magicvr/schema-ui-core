# R3 C1 · 方案候选分析（非决策）

日期：2026-09-04

状态：候选方案，等待用户书面裁决

范围：`I-033-009`、`I-033-010`、`I-033-019`、`I-033-020`、`I-033-021`、`I-033-022`

## 当前证据

- VP-033 要求人工台只在“未绑定且连接成功”时开放，管理员经现有 Telegram sender 以 bot 身份回复；无发言权时 composer 禁用。见 `docs/vision/plans/VP-033-telegram-operator-console.md` 的意图、首波冻结和方向级退出判据。
- 当前 webhook 与 polling 共用 `WebhookHandler.dispatchPayload`，可取得 `UpdateID`、`MessageID`、chat/user 和文本；但当前没有 Telegram session/message 表、人工台 API 或 `getChatMember` 接口。见 `apps/api/internal/channel/telegram/types.go`、`webhook.go`、`apps/api/kernel/telegram.go`。
- 当前 Admin Telegram 页面复用 `settings.read`/`settings.write`，已有约 10 秒 heartbeat 和 `useSchemaCrud().fetcher` 接缝；没有人工台专用权限或 transcript API。见 `apps/web/src/components/telegram-admin-tab.tsx`、`apps/api/modules/channel/telegram/provider.go`。

## A. 发言权探测与缓存失效（I-033-010）

### A1 · 混合策略（AI 推荐候选）

进入或刷新会话时做有限 TTL 的 `getChatMember` 预检；发送时仍以真实发送结果为最终权威。Telegram 403 立即使该会话权限缓存失效、禁用 composer，并要求下一次显式刷新/进入会话重新探测。

- 优点：大多数情况下可在发送前反馈，且不把预检结果误当成最终发送授权。
- 代价：需要定义 TTL、并发去重、预检失败状态、403 识别和 UI 状态机；Bot API 请求较多。

### A2 · 发送前预检

只在发送前调用 `getChatMember`，预检通过后发送；预检失败不发送并禁用 composer。

- 优点：失败不会产生首次出站请求；行为相对容易解释。
- 代价：增加每次发送的延迟和请求；预检与实际发送之间仍有时序窗口，必须处理发送阶段 403。

### A3 · 发送后失效

直接调用 sender；收到 Telegram 403 后把该会话 composer 灰掉并使缓存失效。

- 优点：以真实发送结果为唯一权威，请求少。
- 代价：首次发送可能失败；进入会话时无法提前证明发言权，UI 反馈较晚。

## B. 会话主键与分栏（I-033-019）

### B1 · `chat_id`（AI 推荐候选）

以 Telegram chat 作为会话边界，保存 chat type/title/username；sender user 作为消息元数据。它直接表达“按私聊与群分栏”，群内不同参与者不会被错误拆成多个会话。

### B2 · `subject_id`

以已映射的 Telegram subject 作为会话边界。私聊直观，但群聊参与者会被合并或需要额外群维度，偏离当前产品语义。

### B3 · `chat_id + subject_id`

以 chat 与参与者组合隔离。审计粒度细，但会将一个群拆成多个“会话”，列表/未读/权限语义更复杂。

## C. 人工台权限（I-033-021）

### C1 · 专用 operator 权限（AI 推荐候选）

新增 `telegram.operator.read` 与 `telegram.operator.write`，读取 transcript 使用 read，发送文本使用 write；settings 权限继续只控制连接配置。

- 优点：最小授权清晰，避免能改连接密钥的用户自动获得人工消息读取/发送能力。
- 代价：需要 Provider/RBAC contribution、角色种子/迁移、前后端权限矩阵和更多验证。

### C2 · 复用 `settings.read/write`

沿用现有 Telegram settings 的权限边界。

- 优点：接线少，沿用既有 Admin 页面授权。
- 代价：连接配置管理员同时获得会话内容读取和 bot 发言能力，最小权限较弱。

## D. 低歧义但仍需写入合同的边界

- **I-033-009**：AI 推荐 10 秒短轮询；单页面只允许一个 in-flight 请求，页面失焦/隐藏时暂停，恢复时立即刷新。不解除 SSE/WebSocket。
- **I-033-020**：AI 推荐以 bot 维度的 `update_id` 唯一约束作为 webhook 重试与 polling 重复投递的主幂等键；消息记录同时保存 `message_id`，在 update 缺失或历史补录场景下作为辅助约束。重复 update 不重复写消息，也不重复触发人工/业务分发；事务失败则返回可重试错误。
- **I-033-022**：AI 推荐先建立/锁定会话，再调用 sender；只有 sender 成功后写入 `sent` 消息，发送失败不伪装为已发送，并返回可重试错误。若用户需要失败成绩单或客户端幂等键，应在 C1 明确，否则不扩展为 outbox/重试队列。

以上推荐仅是工程候选，不是用户决策，也不解除 C1 门禁。用户裁决应写入新的 R3 decision ledger，并同步更新信息项状态、路线图和后续实施范围。
